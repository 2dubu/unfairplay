package recovery

import (
	"encoding/binary"
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
		lifted.x[0] = uint64(byte(random.Uint32()))
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
		{"b", bootstrapAccumulatorPassB, runBootstrapAccumulatorStepBCircuit},
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

func TestBootstrapAccumulatorStepBSemanticMatchesCircuit(t *testing.T) {
	random := rand.New(rand.NewSource(int64(bootstrapAccumulatorPassB)))
	for trial := 0; trial < 512; trial++ {
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
		binary.LittleEndian.PutUint64(pages[0x200000][:8], 0x100000000+uint64(trial))
		registers := [31]uint64{}
		for index := range registers {
			registers[index] = random.Uint64()
		}
		registers[2], registers[30] = 0x224ea33b, 0x224ea33b
		registers[7], registers[21] = 0xc0c02ce911b20ca0, 0x818059d223641940
		registers[3], registers[23] = 0x33fcaf9b, 0xc0c02ce9
		registers[24], registers[25] = 0x19c0683c464360c5, 0xe63f97c3b9bc9f3b
		registers[26], registers[9] = 0x4042601c6f58ed3f, 0xbfbd9fe390a712c1
		registers[16], registers[0] = 0x598208200b552ec1, 0xa67df7dff4aad13f
		registers[15] = uint64(random.Uint32())
		registers[8], registers[19], registers[4] = 0, 0, 0
		registers[28], registers[14] = 0x300000, 0x400000
		registers[11], registers[12] = 0, 0
		registers[5], registers[6], registers[13] = 0, 0xffffffff, 0xffffffff

		circuit := newState(&model{pages: clonePages(pages)})
		semantic := newState(&model{pages: clonePages(pages)})
		circuit.x, semantic.x = registers, registers
		circuit.sp, semantic.sp = 0x100100, 0x100100
		if err := circuit.storeValue(circuit.sp+0xa8, 0x200000, 8); err != nil {
			t.Fatal(err)
		}
		semantic.writable = clonePages(circuit.writable)
		semantic.dirty = map[uint64]bool{0x100000: true}
		circuitNext, circuitErr := runBootstrapAccumulatorStepBCircuit(circuit)
		semanticNext, semanticErr := runBootstrapAccumulatorStepB(semantic)
		if circuitErr != nil || semanticErr != nil {
			t.Fatalf("trial %d errors: circuit=%v semantic=%v", trial, circuitErr, semanticErr)
		}
		if circuitNext != semanticNext || circuit.x != semantic.x {
			t.Fatalf("trial %d scalar state differs: next circuit=%#x semantic=%#x\ncircuit=%#v\nsemantic=%#v",
				trial, circuitNext, semanticNext, circuit.x, semantic.x)
		}
		if circuit.n != semantic.n || circuit.z != semantic.z || circuit.c != semantic.c || circuit.v != semantic.v {
			t.Fatalf("trial %d flags differ", trial)
		}
		if !reflect.DeepEqual(circuit.writable, semantic.writable) {
			t.Fatalf("trial %d memory writes differ", trial)
		}
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

func TestBootstrapForwardShuffleIterationMatchesLiftedBlocks(t *testing.T) {
	random := rand.New(rand.NewSource(int64(bootstrapForwardShuffleSeed)))
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
		lifted.x[12] = uint64(1 + random.Intn(271))
		lifted.x[13] = 0x100100
		direct.x = lifted.x
		direct.n, direct.z, direct.c, direct.v = lifted.n, lifted.z, lifted.c, lifted.v

		seed := compact26A428BootstrapBlocks[bootstrapForwardShuffleSeed]
		next, stopped, err := runCompactBlock(lifted, nil, seed)
		if err != nil || stopped || next != bootstrapForwardShuffleStep {
			t.Fatalf("trial %d lifted seed: next=%#x stopped=%v err=%v", trial, next, stopped, err)
		}
		step := compact26A428BootstrapBlocks[bootstrapForwardShuffleStep]
		nextLifted, stopped, err := runCompactBlock(lifted, nil, step)
		if err != nil || stopped {
			t.Fatalf("trial %d lifted step: next=%#x stopped=%v err=%v", trial, nextLifted, stopped, err)
		}
		nextDirect, err := runBootstrapForwardShuffleIteration(direct)
		if err != nil {
			t.Fatalf("trial %d direct iteration: %v", trial, err)
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

func TestBootstrapReverseByteStepMatchesLiftedBlock(t *testing.T) {
	random := rand.New(rand.NewSource(int64(bootstrapReverseBytePass)))
	for trial := 0; trial < 512; trial++ {
		pages := map[uint64][]byte{
			0x100000: make([]byte, pageSize),
			0x200000: make([]byte, pageSize),
			0x300000: make([]byte, pageSize),
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
		lifted.x[9] = uint64(uint32(random.Intn(33) - 16))
		lifted.x[10] = uint64(1 + random.Intn(255))
		lifted.x[11] = uint64(uint32(random.Intn(33) - 16))
		lifted.x[13] = 0x100100
		lifted.x[14] = uint64(uint32(random.Intn(33) - 16))
		lifted.x[21] = 0x200800
		lifted.x[24] = 0x300100
		direct.x = lifted.x
		direct.n, direct.z, direct.c, direct.v = lifted.n, lifted.z, lifted.c, lifted.v

		block := compact26A428BootstrapBlocks[bootstrapReverseBytePass]
		nextLifted, stopped, err := runCompactBlock(lifted, nil, block)
		if err != nil || stopped {
			t.Fatalf("trial %d lifted block: next=%#x stopped=%v err=%v", trial, nextLifted, stopped, err)
		}
		nextDirect, err := runBootstrapReverseByteStep(direct)
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

func TestBootstrapLookupCopyIterationMatchesLiftedBlocks(t *testing.T) {
	random := rand.New(rand.NewSource(int64(bootstrapLookupCopySeed)))
	for trial := 0; trial < 256; trial++ {
		pages := map[uint64][]byte{
			0x100000: make([]byte, pageSize),
			0x200000: make([]byte, pageSize),
			0x300000: make([]byte, pageSize),
			0x400000: make([]byte, pageSize),
			0x500000: make([]byte, pageSize),
			0x600000: make([]byte, pageSize),
			0x601000: make([]byte, pageSize),
			0x700000: make([]byte, pageSize),
			0x800000: make([]byte, pageSize),
			0x900000: make([]byte, pageSize),
		}
		for _, page := range pages {
			if _, err := random.Read(page); err != nil {
				t.Fatal(err)
			}
		}
		sourceOffset := random.Uint64()
		for index := 0x1b9; index <= 0x2b8; index++ {
			address := 0x600000 + uint64(index)*8
			binary.LittleEndian.PutUint64(pages[address&^(pageSize-1)][address&(pageSize-1):], 0x800100-sourceOffset)
		}
		lifted := newState(&model{pages: clonePages(pages)})
		direct := newState(&model{pages: clonePages(pages)})
		for register := range lifted.x {
			lifted.x[register] = random.Uint64()
		}
		lifted.sp = 0x100000
		lifted.x[0] = sourceOffset
		lifted.x[1] = 0
		lifted.x[3] = 0x300100
		lifted.x[4] = 0x400100
		lifted.x[5] = 0x600000
		lifted.x[9] = uint64(uint32(random.Intn(33) - 16))
		lifted.x[11] = 0x900100
		lifted.x[12] = uint64(random.Intn(16))
		lifted.x[13] = 1
		lifted.x[16] = uint64(byte(random.Uint32()))
		if err := lifted.storeValue(lifted.sp+0xde8, 0x200100, 8); err != nil {
			t.Fatal(err)
		}
		if err := lifted.storeValue(lifted.sp+0xde0, 0x200108, 8); err != nil {
			t.Fatal(err)
		}
		if err := lifted.storeValue(lifted.sp+0xdf0, 0x700800, 8); err != nil {
			t.Fatal(err)
		}
		direct.sp = lifted.sp
		direct.x = lifted.x
		direct.n, direct.z, direct.c, direct.v = lifted.n, lifted.z, lifted.c, lifted.v
		direct.writable = clonePages(lifted.writable)
		direct.dirty = map[uint64]bool{}
		for page := range direct.writable {
			direct.dirty[page] = true
		}

		seed := compact26A428BootstrapBlocks[bootstrapLookupCopySeed]
		next, stopped, err := runCompactBlock(lifted, nil, seed)
		if err != nil || stopped || next != bootstrapLookupCopyStep {
			t.Fatalf("trial %d lifted seed: next=%#x stopped=%v err=%v", trial, next, stopped, err)
		}
		step := compact26A428BootstrapBlocks[bootstrapLookupCopyStep]
		nextLifted, stopped, err := runCompactBlock(lifted, nil, step)
		if err != nil || stopped {
			t.Fatalf("trial %d lifted step: next=%#x stopped=%v err=%v", trial, nextLifted, stopped, err)
		}
		nextDirect, err := runBootstrapLookupCopyIteration(direct)
		if err != nil {
			t.Fatalf("trial %d direct iteration: %v", trial, err)
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

func TestBootstrapPointerWalkIterationMatchesLiftedBlocks(t *testing.T) {
	random := rand.New(rand.NewSource(int64(bootstrapPointerWalkSeed)))
	for trial := 0; trial < 512; trial++ {
		counter := random.Uint32()
		removed := uint64(0x6f164ec4) & (uint64(counter) << 1)
		lookupOffset := (uint64(counter) - removed + uint64(0xf22d1348378b2762)) ^ uint64(0xf22d1348378b2762)
		lookupBase := uint64(0x100000000000)
		lookupAddress := lookupBase + lookupOffset
		pages := map[uint64][]byte{
			lookupAddress &^ (pageSize - 1): make([]byte, pageSize),
			0x200000:                        make([]byte, pageSize),
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
		lifted.x[10] = uint64(counter)
		lifted.x[16] = uint64(uint32(random.Intn(33) - 16))
		lifted.x[21] = 0x200800
		lifted.x[24] = lookupBase
		direct.x = lifted.x
		direct.n, direct.z, direct.c, direct.v = lifted.n, lifted.z, lifted.c, lifted.v

		seed := compact26A428BootstrapBlocks[bootstrapPointerWalkSeed]
		next, stopped, err := runCompactBlock(lifted, nil, seed)
		if err != nil || stopped || next != bootstrapPointerWalkDispatch {
			t.Fatalf("trial %d lifted seed: next=%#x stopped=%v err=%v", trial, next, stopped, err)
		}
		dispatch := compact26A428BootstrapBlocks[bootstrapPointerWalkDispatch]
		nextLifted, stopped, err := runCompactBlock(lifted, nil, dispatch)
		if err != nil || stopped {
			t.Fatalf("trial %d lifted dispatch: next=%#x stopped=%v err=%v", trial, nextLifted, stopped, err)
		}
		nextDirect, err := runBootstrapPointerWalkIteration(direct)
		if err != nil {
			t.Fatalf("trial %d direct iteration: %v", trial, err)
		}
		if nextLifted != nextDirect || lifted.x != direct.x {
			t.Fatalf("trial %d state differs: next lifted=%#x direct=%#x\nlifted=%#v\ndirect=%#v", trial, nextLifted, nextDirect, lifted.x, direct.x)
		}
		if lifted.n != direct.n || lifted.z != direct.z || lifted.c != direct.c || lifted.v != direct.v {
			t.Fatalf("trial %d flags differ", trial)
		}
	}
}

func TestBootstrapPairCopyStepMatchesLiftedBlock(t *testing.T) {
	random := rand.New(rand.NewSource(int64(bootstrapPairCopyPass)))
	for trial := 0; trial < 256; trial++ {
		counter := uint32(2 + random.Intn(1024))
		selector := uint32(random.Intn(33) - 16)
		firstIndex := (selector ^ 0xeb62) + counter + 0xffff8181
		secondIndex := counter - 2
		const sourceBase = uint64(0x100000000000)
		const destinationBase = uint64(0x200000000000)
		pages := map[uint64][]byte{
			0x300000: make([]byte, pageSize),
			0x400000: make([]byte, pageSize),
		}
		for _, item := range []struct {
			base  uint64
			index uint32
		}{{sourceBase, firstIndex}, {sourceBase, secondIndex}, {destinationBase, firstIndex}, {destinationBase, secondIndex}} {
			address := addOffset(item.base, int64(int32(item.index))*4)
			if pages[address&^(pageSize-1)] == nil {
				pages[address&^(pageSize-1)] = make([]byte, pageSize)
			}
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
		lifted.x[8] = 0x300100
		lifted.x[9] = 0x300200
		lifted.x[10] = 0x400800
		lifted.x[11] = uint64(selector)
		lifted.x[13] = uint64(counter)
		if err := lifted.storeValue(lifted.x[8]+8, sourceBase, 8); err != nil {
			t.Fatal(err)
		}
		if err := lifted.storeValue(lifted.x[9]+8, destinationBase, 8); err != nil {
			t.Fatal(err)
		}
		direct.x = lifted.x
		direct.n, direct.z, direct.c, direct.v = lifted.n, lifted.z, lifted.c, lifted.v
		direct.writable = clonePages(lifted.writable)
		direct.dirty = map[uint64]bool{}
		for page := range direct.writable {
			direct.dirty[page] = true
		}

		block := compact26A428BootstrapBlocks[bootstrapPairCopyPass]
		nextLifted, stopped, err := runCompactBlock(lifted, nil, block)
		if err != nil || stopped {
			t.Fatalf("trial %d lifted block: next=%#x stopped=%v err=%v", trial, nextLifted, stopped, err)
		}
		nextDirect, err := runBootstrapPairCopyStep(direct)
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

func TestBootstrapPairCopyPassMatchesRepeatedSteps(t *testing.T) {
	random := rand.New(rand.NewSource(0x100176524))
	for trial := 0; trial < 64; trial++ {
		const (
			sourceObject      = uint64(0x300100)
			destinationObject = uint64(0x300200)
			dispatchBase      = uint64(0x400000)
			sourceBase        = uint64(0x100000000000)
			destinationBase   = uint64(0x200000000000)
			selector          = uint32(0x951c)
			exitTarget        = uint64(0x123456789)
		)
		counter := uint32(2 * (1 + random.Intn(64)))
		pages := map[uint64][]byte{
			0x300000:        make([]byte, pageSize),
			sourceBase:      make([]byte, pageSize),
			destinationBase: make([]byte, pageSize),
		}
		for _, index := range []uint32{selector, selector ^ 0x24} {
			address := dispatchBase + uint64(index)*8
			if pages[address&^(pageSize-1)] == nil {
				pages[address&^(pageSize-1)] = make([]byte, pageSize)
			}
		}
		for _, page := range pages {
			if _, err := random.Read(page); err != nil {
				t.Fatal(err)
			}
		}
		base := newState(&model{pages: clonePages(pages)})
		for register := range base.x {
			base.x[register] = random.Uint64()
		}
		base.x[8] = sourceObject
		base.x[9] = destinationObject
		base.x[10] = dispatchBase
		base.x[11] = uint64(selector)
		base.x[13] = uint64(counter)
		if err := base.storeValue(sourceObject+8, sourceBase, 8); err != nil {
			t.Fatal(err)
		}
		if err := base.storeValue(destinationObject+8, destinationBase, 8); err != nil {
			t.Fatal(err)
		}
		if err := base.storeValue(dispatchBase+uint64(selector)*8, exitTarget, 8); err != nil {
			t.Fatal(err)
		}
		if err := base.storeValue(dispatchBase+uint64(selector^0x24)*8, bootstrapPairCopyPass, 8); err != nil {
			t.Fatal(err)
		}

		repeated := newState(base.model)
		bulk := newState(base.model)
		for _, item := range []*state{repeated, bulk} {
			item.x = base.x
			item.n, item.z, item.c, item.v = base.n, base.z, base.c, base.v
			item.writable = clonePages(base.writable)
			item.dirty = make(map[uint64]bool, len(item.writable))
			for page := range item.writable {
				item.dirty[page] = true
			}
		}

		var repeatedNext uint64
		for pair := uint32(0); pair < counter/2; pair++ {
			var err error
			repeatedNext, err = runBootstrapPairCopyStep(repeated)
			if err != nil {
				t.Fatalf("trial %d pair %d: %v", trial, pair, err)
			}
		}
		bulkNext, err := runBootstrapPairCopyPass(bulk)
		if err != nil {
			t.Fatalf("trial %d bulk pass: %v", trial, err)
		}
		if repeatedNext != bulkNext || repeated.x != bulk.x ||
			repeated.n != bulk.n || repeated.z != bulk.z || repeated.c != bulk.c || repeated.v != bulk.v {
			t.Fatalf("trial %d scalar state differs: repeated=%#x bulk=%#x", trial, repeatedNext, bulkNext)
		}
		if !reflect.DeepEqual(repeated.writable, bulk.writable) {
			t.Fatalf("trial %d memory differs", trial)
		}
	}
}

func TestBootstrapConstantFillIterationMatchesLiftedBlocks(t *testing.T) {
	random := rand.New(rand.NewSource(int64(bootstrapConstantFillSeed)))
	for trial := 0; trial < 256; trial++ {
		counter := uint32(random.Intn(1024))
		selector := uint32(random.Intn(33) - 16)
		first := (selector + 0xffffe51b) | 0x4561
		index := counter + 0xde3d1312 + first
		const stack = uint64(0x100000000000)
		output := addOffset(stack+0x2d0, int64(int32(index))*4)
		pages := map[uint64][]byte{output &^ (pageSize - 1): make([]byte, pageSize)}
		for page := uint64(0x4fe000); page <= 0x502000; page += pageSize {
			pages[page] = make([]byte, pageSize)
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
		lifted.sp = stack
		lifted.x[9] = uint64(selector)
		lifted.x[15] = uint64(counter)
		lifted.x[25] = 0x500800
		direct.sp = lifted.sp
		direct.x = lifted.x
		direct.n, direct.z, direct.c, direct.v = lifted.n, lifted.z, lifted.c, lifted.v

		seed := compact26A428BootstrapBlocks[bootstrapConstantFillSeed]
		next, stopped, err := runCompactBlock(lifted, nil, seed)
		if err != nil || stopped || next != bootstrapConstantFillDispatch {
			t.Fatalf("trial %d lifted seed: next=%#x stopped=%v err=%v", trial, next, stopped, err)
		}
		dispatch := compact26A428BootstrapBlocks[bootstrapConstantFillDispatch]
		nextLifted, stopped, err := runCompactBlock(lifted, nil, dispatch)
		if err != nil || stopped {
			t.Fatalf("trial %d lifted dispatch: next=%#x stopped=%v err=%v", trial, nextLifted, stopped, err)
		}
		nextDirect, err := runBootstrapConstantFillIteration(direct)
		if err != nil {
			t.Fatalf("trial %d direct iteration: %v", trial, err)
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
