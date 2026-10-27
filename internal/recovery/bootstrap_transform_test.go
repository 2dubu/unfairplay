package recovery

import (
	"math/rand"
	"reflect"
	"testing"
)

func TestBootstrapRecurrenceMatchesLiftedKernel(t *testing.T) {
	const (
		leftPage        = uint64(0x100000)
		rightPage       = uint64(0x200000)
		coefficientPage = uint64(0x300000)
		dispatchPage    = uint64(0x400000)
		baseOffset      = uint64(0x800)
	)
	random := rand.New(rand.NewSource(0x26a428))
	for trial := 0; trial < 512; trial++ {
		pages := map[uint64][]byte{
			leftPage:        make([]byte, pageSize),
			rightPage:       make([]byte, pageSize),
			coefficientPage: make([]byte, pageSize),
			dispatchPage:    make([]byte, pageSize),
		}
		for _, page := range pages {
			if _, err := random.Read(page); err != nil {
				t.Fatal(err)
			}
		}

		lifted := newState(&model{pages: clonePages(pages)})
		direct := newState(&model{pages: clonePages(pages)})
		for register := range lifted.x {
			lifted.x[register] = random.Uint64()
		}
		index := uint32(int32(random.Intn(65) - 32))
		counter := uint32(random.Intn(65) - 32)
		lifted.x[4] = leftPage + baseOffset
		lifted.x[5] = rightPage + baseOffset
		lifted.x[22] = coefficientPage + baseOffset
		lifted.x[25] = dispatchPage + baseOffset
		lifted.x[20] = uint64(counter)
		lifted.x[13] = uint64(index - counter)
		lifted.x[21] = uint64(uint32(int32(random.Intn(65)-32)) - counter)
		// Keep the computed dispatch index inside the synthetic page while still
		// exercising both sides of the conditional compare.
		lifted.x[0] = 0
		lifted.x[9] = uint64(uint32(random.Intn(65) - 32))
		lifted.n, lifted.z = random.Intn(2) != 0, random.Intn(2) != 0
		lifted.c, lifted.v = random.Intn(2) != 0, random.Intn(2) != 0
		direct.x = lifted.x
		direct.n, direct.z, direct.c, direct.v = lifted.n, lifted.z, lifted.c, lifted.v

		if err := runCompactKernel(lifted, compact26A428BootstrapKernel); err != nil {
			t.Fatalf("trial %d lifted kernel: %v", trial, err)
		}
		if err := runBootstrapRecurrenceStep(direct, fault{}); err != nil {
			t.Fatalf("trial %d direct recurrence: %v", trial, err)
		}
		if lifted.x != direct.x {
			t.Fatalf("trial %d register state differs\nlifted=%#v\ndirect=%#v", trial, lifted.x, direct.x)
		}
		if lifted.n != direct.n || lifted.z != direct.z || lifted.c != direct.c || lifted.v != direct.v {
			t.Fatalf("trial %d flags differ: lifted=%v%v%v%v direct=%v%v%v%v", trial,
				lifted.n, lifted.z, lifted.c, lifted.v, direct.n, direct.z, direct.c, direct.v)
		}
		if !reflect.DeepEqual(lifted.writable, direct.writable) {
			t.Fatalf("trial %d memory writes differ", trial)
		}
	}
}

func TestBootstrapRecurrenceRejectsMissingMemory(t *testing.T) {
	s := newState(&model{pages: make(map[uint64][]byte)})
	if err := runBootstrapRecurrenceStep(s, fault{}); err == nil {
		t.Fatal("bootstrap recurrence accepted missing tables")
	}
}

func TestBootstrapWordStepsMatchLiftedBlocks(t *testing.T) {
	tests := []struct {
		name    string
		address uint64
		direct  func(*state) (uint64, error)
	}{
		{"short", bootstrapShortWordPass, runBootstrapShortWordStep},
		{"long", bootstrapLongWordPass, runBootstrapLongWordStep},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			random := rand.New(rand.NewSource(int64(test.address)))
			for trial := 0; trial < 256; trial++ {
				pages := map[uint64][]byte{
					0x100000: make([]byte, pageSize),
					0x200000: make([]byte, pageSize),
					0x300000: make([]byte, pageSize),
					0x400000: make([]byte, pageSize),
				}
				for _, page := range pages {
					if _, err := random.Read(page); err != nil {
						t.Fatal(err)
					}
				}
				lifted := newState(&model{pages: clonePages(pages)})
				direct := newState(&model{pages: clonePages(pages)})
				for register := range lifted.x {
					lifted.x[register] = random.Uint64()
				}
				lifted.x[0] = 0x100400
				lifted.x[1] = 0
				lifted.x[8] = uint64(uint32(random.Intn(32)) * 4)
				lifted.x[9] = uint64(uint32(random.Intn(32)) * 4)
				lifted.x[10] = uint64(uint32(random.Intn(32) - 16))
				lifted.x[11] = 0x300100
				lifted.x[19] = 0x200800
				lifted.x[21] = 0x400100
				direct.x = lifted.x
				direct.n, direct.z, direct.c, direct.v = lifted.n, lifted.z, lifted.c, lifted.v

				block := compact26A428BootstrapBlocks[test.address]
				nextLifted, stopped, err := runCompactBlock(lifted, nil, block)
				if err != nil || stopped {
					t.Fatalf("trial %d lifted block: next=%#x stopped=%v err=%v", trial, nextLifted, stopped, err)
				}
				nextDirect, err := test.direct(direct)
				if err != nil {
					t.Fatalf("trial %d direct step: %v", trial, err)
				}
				if nextLifted != nextDirect || lifted.x != direct.x {
					t.Fatalf("trial %d state differs: next lifted=%#x direct=%#x\nlifted=%#v\ndirect=%#v", trial, nextLifted, nextDirect, lifted.x, direct.x)
				}
				if lifted.n != direct.n || lifted.z != direct.z || lifted.c != direct.c || lifted.v != direct.v {
					t.Fatalf("trial %d flags differ", trial)
				}
				if !reflect.DeepEqual(lifted.writable, direct.writable) {
					t.Fatalf("trial %d writes differ", trial)
				}
			}
		})
	}
}

func TestBootstrapAccumulatorStepsMatchLiftedBlocks(t *testing.T) {
	tests := []struct {
		name    string
		address uint64
		direct  func(*state) (uint64, error)
	}{
		{"a", bootstrapAccumulatorPassA, runBootstrapAccumulatorStepA},
		{"b", bootstrapAccumulatorPassB, runBootstrapAccumulatorStepB},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			random := rand.New(rand.NewSource(int64(test.address)))
			for trial := 0; trial < 256; trial++ {
				pages := map[uint64][]byte{
					0x100000: make([]byte, pageSize),
					0x200000: make([]byte, pageSize),
					0x300000: make([]byte, pageSize),
					0x400000: make([]byte, pageSize),
					0x500000: make([]byte, pageSize),
					0x600000: make([]byte, pageSize),
					0x700000: make([]byte, pageSize),
				}
				for _, page := range pages {
					if _, err := random.Read(page); err != nil {
						t.Fatal(err)
					}
				}
				lifted := newState(&model{pages: clonePages(pages)})
				direct := newState(&model{pages: clonePages(pages)})
				for register := range lifted.x {
					lifted.x[register] = random.Uint64()
				}
				lifted.sp = 0x100100
				if err := lifted.storeValue(lifted.sp+0xa8, 0x200800, 8); err != nil {
					t.Fatal(err)
				}
				lifted.x[0] = 0
				lifted.x[2] = uint64(uint32(random.Intn(33) - 16))
				lifted.x[8] = uint64(uint32(random.Intn(32)))
				lifted.x[9] = 0x300800
				lifted.x[11] = uint64(uint32(random.Intn(32)))
				lifted.x[12] = uint64(uint32(random.Intn(17)-8) - uint32(lifted.x[11]))
				lifted.x[14] = uint64(uint32(random.Intn(17)-8) - uint32(lifted.x[11]))
				lifted.x[19] = uint64(uint32(random.Intn(17)-8) - uint32(lifted.x[8]))
				lifted.x[21] = 0
				lifted.x[26] = 0x400800
				lifted.x[27] = 0x500800
				lifted.x[28] = 0x600800
				lifted.x[4] = uint64(uint32(random.Intn(17)-8) - uint32(lifted.x[8]))
				if test.address == bootstrapAccumulatorPassB {
					lifted.x[11] = uint64(uint32(random.Intn(33) - 16))
					lifted.x[12] = 0
					lifted.x[14] = 0x700800
				}
				direct.sp = lifted.sp
				direct.x = lifted.x
				direct.n, direct.z, direct.c, direct.v = lifted.n, lifted.z, lifted.c, lifted.v
				direct.writable = clonePages(lifted.writable)
				direct.dirty = map[uint64]bool{}
				for page := range direct.writable {
					direct.dirty[page] = true
				}

				block := compact26A428BootstrapBlocks[test.address]
				nextLifted, stopped, err := runCompactBlock(lifted, nil, block)
				if err != nil || stopped {
					t.Fatalf("trial %d lifted block: next=%#x stopped=%v err=%v", trial, nextLifted, stopped, err)
				}
				nextDirect, err := test.direct(direct)
				if err != nil {
					t.Fatalf("trial %d direct step: %v", trial, err)
				}
				if nextLifted != nextDirect || lifted.x != direct.x {
					t.Fatalf("trial %d state differs: next lifted=%#x direct=%#x\nlifted=%#v\ndirect=%#v", trial, nextLifted, nextDirect, lifted.x, direct.x)
				}
				if lifted.n != direct.n || lifted.z != direct.z || lifted.c != direct.c || lifted.v != direct.v {
					t.Fatalf("trial %d flags differ", trial)
				}
				if !reflect.DeepEqual(lifted.writable, direct.writable) {
					t.Fatalf("trial %d writes differ", trial)
				}
			}
		})
	}
}

func TestBootstrapShuffleStepMatchesLiftedBlock(t *testing.T) {
	random := rand.New(rand.NewSource(int64(bootstrapShufflePass)))
	for trial := 0; trial < 512; trial++ {
		pages := map[uint64][]byte{0x100000: make([]byte, pageSize)}
		if _, err := random.Read(pages[0x100000]); err != nil {
			t.Fatal(err)
		}
		lifted := newState(&model{pages: clonePages(pages)})
		direct := newState(&model{pages: clonePages(pages)})
		for register := range lifted.x {
			lifted.x[register] = random.Uint64()
		}
		lifted.x[11] = uint64(2 + random.Intn(255))
		lifted.x[16] = 0x100100
		direct.x = lifted.x
		direct.n, direct.z, direct.c, direct.v = lifted.n, lifted.z, lifted.c, lifted.v

		block := compact26A428BootstrapBlocks[bootstrapShufflePass]
		nextLifted, stopped, err := runCompactBlock(lifted, nil, block)
		if err != nil || stopped {
			t.Fatalf("trial %d lifted block: next=%#x stopped=%v err=%v", trial, nextLifted, stopped, err)
		}
		nextDirect, err := runBootstrapShuffleStep(direct)
		if err != nil {
			t.Fatalf("trial %d direct step: %v", trial, err)
		}
		if nextLifted != nextDirect || lifted.x != direct.x {
			t.Fatalf("trial %d state differs: next lifted=%#x direct=%#x\nlifted=%#v\ndirect=%#v", trial, nextLifted, nextDirect, lifted.x, direct.x)
		}
		if lifted.n != direct.n || lifted.z != direct.z || lifted.c != direct.c || lifted.v != direct.v {
			t.Fatalf("trial %d flags differ", trial)
		}
		if !reflect.DeepEqual(lifted.writable, direct.writable) {
			t.Fatalf("trial %d writes differ", trial)
		}
	}
}
