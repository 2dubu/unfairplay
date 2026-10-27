package recovery

import (
	"encoding/binary"
	"testing"
)

// These are the exact live instructions in the protected lane-5 value slice.
// Address-only dispatcher work is intentionally excluded.
var outerTableLane5ReferenceIndexes = [...]int{
	29558,
	29860, 29865,
	29888, 29889, 29890, 29891, 29892, 29893, 29894, 29895, 29896,
	29897, 29898, 29899, 29900, 29901, 29902, 29903, 29904, 29905, 29906,
	29907, 29908, 29909, 29910, 29911, 29912, 29913, 29914, 29915, 29916,
	29917, 29918, 29919, 29920, 29921,
}

func TestOuterTableLane5MaskMatchesProtectedSlice(t *testing.T) {
	const (
		stack     = uint64(0x7000000e0000)
		tableBase = uint64(0x710000000000)
	)
	instructions := make([]compactInstruction, 0, len(outerTableLane5ReferenceIndexes))
	for _, index := range outerTableLane5ReferenceIndexes {
		instructions = append(instructions, compact26A428Outer.instructions[index])
	}
	kernel := compactKernel{start: 0x100479704, instructions: instructions}

	value := uint32(0x243f6a88)
	for iteration := range 512 {
		value = value*1664525 + 1013904223
		pendingS1 := value
		value = value*1664525 + 1013904223
		s3 := value
		value = value*1664525 + 1013904223
		s0 := value
		value = value*1664525 + 1013904223
		old := value

		pages := map[uint64][]byte{
			stack + 0x1000: make([]byte, pageSize),
			stack + 0x4000: make([]byte, pageSize),
			stack + 0x8000: make([]byte, pageSize),
			tableBase:      make([]byte, pageSize),
		}
		binary.LittleEndian.PutUint32(pages[stack+0x1000][0x134:], s0)
		binary.LittleEndian.PutUint64(pages[stack+0x4000][0x3d8:], tableBase)
		binary.LittleEndian.PutUint32(pages[tableBase], old)
		s := newState(&model{pages: pages})
		s.sp = stack
		s.x[10] = uint64(pendingS1)
		s.x[14] = uint64(s3)
		if err := runCompactKernel(s, kernel); err != nil {
			t.Fatalf("iteration %d: %v", iteration, err)
		}
		got, err := s.load(tableBase, 4)
		if err != nil {
			t.Fatal(err)
		}
		want := old ^ outerTableLane5Mask(pendingS1, s3, s0)
		if uint32(got) != want {
			t.Fatalf("iteration %d: protected=%#08x semantic=%#08x", iteration, uint32(got), want)
		}
	}
}

func TestOuterTableRecoveredMaskVectors(t *testing.T) {
	const (
		a = uint32(0x01234567)
		b = uint32(0x89abcdef)
		c = uint32(0x0f1e2d3c)
		d = uint32(0x4b5a6978)
	)
	tests := []struct {
		name string
		got  uint32
		want uint32
	}{
		{"lane0", outerTableLane0Mask(a, b, c, d, 0xdeadbeef, 0x13579bdf, 0x2468ace0, 0xfedcba98), 0x0723b0a0},
		{"lane1", outerTableLane1Mask(a, b, c, d), 0xac3635c9},
		{"lane2", outerTableLane2Mask(a, b, c, d, 0xdeadbeef), 0x4770d3a2},
		{"lane3", outerTableLane3Mask(a, b), 0xd583636e},
		{"lane4", outerTableLane4Mask(a, b, c, d), 0x33333333},
		{"lane5", outerTableLane5Mask(a, b, c), 0x85a68f14},
		{"lane6", outerTableLane6Mask(a, b, c), 0x561d3c4e},
		{"lane7", outerTableLane7Mask(a, b, c, d, 0xdeadbeef), 0x80bcca61},
		{"lane8", outerTableLane8Mask(a, b, c, d), 0x33333333},
		{"lane9", outerTableLane9Mask(a, b, c), 0x9e05d60b},
		{"lane10", outerTableLane10Mask(a, b, c), 0xbf0a18a3},
		{"lane11", outerTableLane11Mask(a, b, c, d), 0xd19963b6},
		{"lane12", outerTableLane12Mask(a, b, c, d), 0xf5a79e13},
		{"lane13", outerTableLane13Mask(a, b, c), 0xc03a5e91},
		{"lane14", outerTableLane14Mask(a, b, c, d), 0xb6e33e77},
		{"lane15", outerTableLane15Mask(a, b, c, d), 0x33333333},
		{"lane16", outerTableLane16Mask(a, b, c), 0x6c8ae9c5},
		{"lane17", outerTableLane17Mask(a, b, c), 0xda80f414},
		{"lane18", outerTableLane18Mask(a, b, c, d, 0xdeadbeef), 0xfe677373},
		{"lane19", outerTableLane19Mask(a, b, c, d, 0xdeadbeef, 0x13579bdf), 0x169731cb},
		{"lane20", outerTableLane20Mask(a, b, c, d), 0x0a92601c},
	}
	for _, test := range tests {
		if test.got != test.want {
			t.Errorf("%s mask = %#08x, want %#08x", test.name, test.got, test.want)
		}
	}
}
