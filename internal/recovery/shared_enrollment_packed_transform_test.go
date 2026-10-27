package recovery

import (
	"math/rand"
	"testing"
)

func TestSharedEnrollmentPackedPathMatchesCircuit(t *testing.T) {
	random := rand.New(rand.NewSource(int64(sharedEnrollmentPackedStart)))
	const (
		stack    = uint64(0x100100)
		work     = uint64(0x200000)
		dispatch = uint64(0x300000)
		target   = uint64(0x123456789)
	)
	dispatchAddress := dispatch + 0x19d6*8
	for trial := 0; trial < 4096; trial++ {
		pages := map[uint64][]byte{
			stack &^ (pageSize - 1):           make([]byte, pageSize),
			work:                              make([]byte, pageSize),
			dispatchAddress &^ (pageSize - 1): make([]byte, pageSize),
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
		iteration := random.Uint32()
		current := byte(random.Uint32())&0xfc | 2
		base.sp = stack
		base.x[0] = uint64(byte(iteration))
		base.x[1] = sharedEnrollmentPackedStart
		base.x[2] = 0x19d6
		base.x[3] = 0x060be5d0
		base.x[5] = 0xba1c1cae
		base.x[6] = 0xae
		base.x[7] = 0x52
		base.x[8] = 0x1b
		base.x[11] = uint64(0x45e3e363 + iteration)
		base.x[12] = uint64(0x22f1f257 + iteration)
		base.x[14] = uint64(0x514647f8 - 5*iteration)
		base.x[15] = uint64(0xd35cb3b4 + 5*iteration)
		base.x[16] = uint64(0x22f1f1a9 + iteration)
		base.x[17] = 0x60e8
		base.x[19] = 0x19d6
		base.x[20] = uint64(current)
		base.x[21] = dispatch
		base.x[22] = uint64(((uint32(current) * 2) & 0xffffffc8) - uint32(current))
		base.x[23] = 0xfffffff6
		base.x[24] = 0xaca8
		base.x[25] = 0x739b
		if err := base.storeValue(stack+0x58, work, 8); err != nil {
			t.Fatal(err)
		}
		if err := base.storeValue(dispatchAddress, target, 8); err != nil {
			t.Fatal(err)
		}

		circuit := newState(base.model)
		semantic := newState(base.model)
		for _, item := range []*state{circuit, semantic} {
			item.sp, item.x = base.sp, base.x
			item.writable = clonePages(base.writable)
			item.dirty = map[uint64]bool{}
			for page := range item.writable {
				item.dirty[page] = true
			}
		}
		next, err := runSelectionDirectBlock100264a8c(circuit)
		if err != nil || next != 0x100264c00 {
			t.Fatalf("trial %d A: next=%#x err=%v", trial, next, err)
		}
		circuitNext, err := runSelectionDirectBlock100264c00(circuit)
		if err != nil {
			t.Fatalf("trial %d B: %v", trial, err)
		}
		semanticNext, err := runSharedEnrollmentPackedPath(semantic)
		if err != nil {
			t.Fatalf("trial %d semantic: %v", trial, err)
		}
		if circuitNext != semanticNext {
			t.Fatalf("trial %d next: circuit=%#x semantic=%#x", trial, circuitNext, semanticNext)
		}
		for _, register := range []int{1, 2, 3, 17, 19, 20, 21, 22, 27, 28} {
			if circuit.x[register] != semantic.x[register] {
				t.Fatalf("trial %d x%d: circuit=%#x semantic=%#x", trial, register, circuit.x[register], semantic.x[register])
			}
		}
	}
}

func TestSharedEnrollmentPackedWeightsAreUnique(t *testing.T) {
	seen := map[int32]bool{}
	for value := 0; value < 256; value++ {
		output := sharedEnrollmentPackedWeight(byte(value))
		if seen[output] {
			t.Fatalf("weight repeats %d", output)
		}
		seen[output] = true
	}
}
