package recovery

import (
	"math/rand"
	"reflect"
	"testing"
)

func runBootstrapByteWriterCircuit(s *state) (uint64, error) {
	pc := bootstrapByteWriterControl
	for count := 0; count < 32; count++ {
		var next uint64
		var err error
		switch pc {
		case bootstrapByteWriterControl:
			next, err = runBootstrapByteWriterControl(s)
		case bootstrapByteWriterMix:
			next, err = runBootstrapByteWriterMix(s)
		case bootstrapByteWriterLoad:
			next, err = runBootstrapByteWriterLoad(s)
		case bootstrapByteWriterStore:
			next, err = runBootstrapByteWriterStore(s)
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

func TestBootstrapByteWriterCycleMatchesCircuit(t *testing.T) {
	random := rand.New(rand.NewSource(int64(bootstrapByteWriterControl)))
	const descriptor, table, outputObject, output, dispatch = uint64(0x100000), uint64(0x200000), uint64(0x300000), uint64(0x400000), uint64(0x500000)
	const exitTarget = uint64(0x123456789)
	for trial := 0; trial < 1024; trial++ {
		withLoad := trial&1 == 0
		pages := map[uint64][]byte{
			descriptor:   make([]byte, pageSize),
			table:        make([]byte, pageSize),
			outputObject: make([]byte, pageSize),
			output:       make([]byte, pageSize),
			dispatch:     make([]byte, pageSize),
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
		base.x[0] = outputObject
		base.x[1] = 0
		base.x[2] = 1
		base.x[3] = 0
		base.x[4] = 0
		base.x[5] = 2
		base.x[6] = 0
		base.x[7] = 4
		base.x[8] = descriptor
		base.x[9] = 1
		base.x[10] = 0
		base.x[11] = dispatch
		base.x[12] = 0
		base.x[13] = uint64(random.Uint32())
		base.x[14] = 4
		base.x[15] = 1
		base.x[16] = ^uint64(0)
		base.x[19] = ^uint64(0)
		base.x[20] = 0
		base.x[21] = 0
		base.x[22] = 0
		base.x[24] = 0
		base.x[19] = ^uint64(0)
		base.x[20] = 0
		base.x[21] = 0
		base.x[22] = 0
		base.x[23] = 0
		base.x[24] = 0
		base.x[25] = ^uint64(0)
		base.x[27] = uint64(random.Uint32())
		if withLoad {
			base.x[14] = 0
			base.x[15] = 0
			base.x[16] = 0
		}
		if err := base.storeValue(descriptor+8, table, 8); err != nil {
			t.Fatal(err)
		}
		if err := base.storeValue(outputObject+0x18, output, 8); err != nil {
			t.Fatal(err)
		}
		targets := map[uint32]uint64{0: bootstrapByteWriterStore, 4: exitTarget}
		if withLoad {
			targets[1] = bootstrapByteWriterMix
			targets[2] = bootstrapByteWriterLoad
		}
		for selector, target := range targets {
			if err := base.storeValue(dispatch+uint64(selector)*8, target, 8); err != nil {
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
		circuitNext, err := runBootstrapByteWriterCircuit(circuit)
		if err != nil {
			t.Fatal(err)
		}
		semanticNext, err := runBootstrapByteWriterCycle(semantic)
		if err != nil {
			t.Fatal(err)
		}
		if circuitNext != semanticNext || circuit.x != semantic.x || circuit.n != semantic.n || circuit.z != semantic.z || circuit.c != semantic.c || circuit.v != semantic.v || !reflect.DeepEqual(circuit.writable, semantic.writable) {
			t.Fatalf("trial %d state differs", trial)
		}
	}
}
