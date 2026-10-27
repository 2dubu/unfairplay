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
