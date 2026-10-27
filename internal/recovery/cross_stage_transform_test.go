package recovery

import (
	"math/rand"
	"reflect"
	"testing"
)

func runCrossStageBitCycleCircuit(s *state) (uint64, error) {
	pc := crossStageBitCycleStart
	for count := 0; count < 64; count++ {
		var next uint64
		var err error
		switch pc {
		case crossStageBitCycleStart:
			next, err = runCrossStageBitCycleStart(s)
		case crossStageBitCyclePlace:
			next, err = runCrossStageBitCyclePlace(s)
		case crossStageBitCycleTest:
			next, err = runCrossStageBitCycleTest(s)
		case crossStageBitCycleWrite:
			next, err = runCrossStageBitCycleWrite(s)
		default:
			return pc, nil
		}
		if err != nil {
			return pc, err
		}
		pc = next
	}
	return 0, nil
}

func TestCrossStageBitCycleMatchesCircuit(t *testing.T) {
	random := rand.New(rand.NewSource(int64(crossStageBitCycleStart)))
	const descriptor, source, output, dispatch = uint64(0x100000), uint64(0x200000), uint64(0x300000), uint64(0x400000)
	const exitTarget = uint64(0x123456789)
	selectors := map[uint32]uint64{
		0x0100: crossStageBitCyclePlace,
		0x0200: crossStageBitCycleTest,
		0x0201: crossStageBitCycleWrite,
		0x0a1a: crossStageBitCycleStart,
		0x0b1a: exitTarget,
	}
	for trial := 0; trial < 1024; trial++ {
		pages := map[uint64][]byte{
			descriptor: make([]byte, pageSize),
			source:     make([]byte, pageSize),
			output:     make([]byte, pageSize),
		}
		for selector := range selectors {
			address := addOffset(dispatch, int64(int32(selector))*8)
			pages[address&^(pageSize-1)] = make([]byte, pageSize)
		}
		for _, page := range pages {
			if _, err := random.Read(page); err != nil {
				t.Fatal(err)
			}
		}
		model := &model{pages: clonePages(pages)}
		base := newState(model)
		for register := range base.x {
			base.x[register] = random.Uint64()
		}
		base.x[0] = descriptor
		base.x[1] = uint64(random.Uint32())
		shiftStep := uint32(8)
		base.x[2] = uint64(shiftStep - 0x0a1a)
		base.x[3] = 1
		base.x[4] = 3
		base.x[7] = 1
		base.x[8] = dispatch
		base.x[9] = descriptor + 0x100
		base.x[12] = 0
		base.x[14] = 0x100
		base.x[15] = 0x300
		base.x[16] = 0
		base.x[19] = 0
		base.x[20] = 0
		base.x[21] = 0
		base.x[22] = 0
		base.x[23] = 0x100
		base.x[24] = uint64(random.Uint32())
		base.x[26] = 0x100
		base.x[27] = 0
		base.x[30] = 0x100
		if err := base.storeValue(descriptor+0x18, source, 8); err != nil {
			t.Fatal(err)
		}
		if err := base.storeValue(descriptor+0x108, output-4, 8); err != nil {
			t.Fatal(err)
		}
		for selector, target := range selectors {
			if err := base.storeValue(addOffset(dispatch, int64(int32(selector))*8), target, 8); err != nil {
				t.Fatal(err)
			}
		}

		circuit := newState(model)
		semantic := newState(model)
		for _, item := range []*state{circuit, semantic} {
			item.x = base.x
			item.writable = clonePages(base.writable)
			item.dirty = map[uint64]bool{}
			for page := range item.writable {
				item.dirty[page] = true
			}
		}
		circuitNext, err := runCrossStageBitCycleCircuit(circuit)
		if err != nil {
			t.Fatal(err)
		}
		semanticNext, err := runCrossStageBitCycle(semantic)
		if err != nil {
			t.Fatal(err)
		}
		if circuitNext != semanticNext || circuit.x != semantic.x || circuit.n != semantic.n || circuit.z != semantic.z || circuit.c != semantic.c || circuit.v != semantic.v || !reflect.DeepEqual(circuit.writable, semantic.writable) {
			t.Fatalf("trial %d state differs", trial)
		}
	}
}
