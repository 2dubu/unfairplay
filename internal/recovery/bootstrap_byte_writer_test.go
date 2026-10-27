package recovery

import (
	"math/rand"
	"reflect"
	"testing"
)

func runBootstrapByteWriterCircuit(s *state) (uint64, error) {
	pc := bootstrapByteWriterControl
	for count := 0; count < 4096; count++ {
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
	for trial := 0; trial < 512; trial++ {
		wordCount := uint32(1 + random.Intn(64))
		targets := map[uint32]uint64{
			0x82a7: bootstrapByteWriterMix,
			0x4dcf: bootstrapByteWriterStore,
			0x4de1: bootstrapByteWriterLoad,
			0x208a: bootstrapByteWriterControl,
			0x2072: exitTarget,
		}
		pages := map[uint64][]byte{
			descriptor:   make([]byte, pageSize),
			table:        make([]byte, pageSize),
			outputObject: make([]byte, pageSize),
			output:       make([]byte, pageSize),
			dispatch:     make([]byte, pageSize),
		}
		for selector := range targets {
			address := dispatch + uint64(selector)*8
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
		base.x[0] = outputObject
		base.x[1] = 0xffffdfae
		base.x[2] = 0xcf68
		base.x[3] = 0xa7aea628
		base.x[4] = 0x58512c7a
		base.x[5] = 0x205c
		base.x[6] = 0x4dcf
		base.x[7] = 0x5320ebab
		base.x[8] = descriptor
		base.x[9] = uint64(wordCount)
		base.x[10] = uint64(wordCount * 4)
		base.x[11] = dispatch
		base.x[12] = 0
		base.x[13] = 0x224ea33b
		base.x[14] = 0xac1367ab
		base.x[15] = 0x20
		base.x[16] = 0x53ec6af7
		base.x[19] = 0xacdee6f7
		base.x[20] = 0x2073
		base.x[21] = 0x9f945f5d
		base.x[22] = 0x606bcdff
		base.x[23] = 0xca
		base.x[24] = 0x87483ef6
		base.x[25] = 0
		base.x[27] = uint64(random.Uint32())
		if err := base.storeValue(descriptor+8, table, 8); err != nil {
			t.Fatal(err)
		}
		storeOffset := uint32(0x87483ef6 + 0x4dcf)
		firstWriteOffset := wordCount*4 + storeOffset
		outputPointer := output + uint64(wordCount*4-1) - uint64(firstWriteOffset)
		if err := base.storeValue(outputObject+0x18, outputPointer, 8); err != nil {
			t.Fatal(err)
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

func TestBootstrapByteWriterDecodeIdentity(t *testing.T) {
	for value := uint32(0); value < 256; value++ {
		old := byte((value - (0xca & (value << 1)) - 0x1b) ^ 0xfe)
		if got, want := byte(value)^0x1b, old; got != want {
			t.Fatalf("value %#x decoded to %#x, want %#x", value, got, want)
		}
	}
}
