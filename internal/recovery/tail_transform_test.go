package recovery

import (
	"encoding/binary"
	"testing"
)

func TestTailTransform26A428MatchesLiftedHookDomain(t *testing.T) {
	const (
		stack       = uint64(0x7000000e0000)
		tableBase   = uint64(0x710000000000)
		outputBase  = uint64(0x710000001000)
		contextBase = uint64(0x710000002000)
		leftSeed    = uint32(0x449341de)
		rightSeed   = uint32(0x224)
	)
	table := make([]uint32, pairTransformTableWords)
	for index := range table {
		table[index] = uint32(index)*0x7f4a7c15 ^ 0x3c6ef372
	}
	want, ok := tailTransform26A428(table, leftSeed^rightSeed, 87)
	if !ok {
		t.Fatal("semantic tail transform rejected its proven domain")
	}

	pages := map[uint64][]byte{
		stack + 0x4000: make([]byte, pageSize),
		tableBase:      make([]byte, pageSize),
		outputBase:     make([]byte, pageSize),
		contextBase:    make([]byte, pageSize),
	}
	binary.LittleEndian.PutUint64(pages[stack+0x4000][0x3d8:], tableBase)
	binary.LittleEndian.PutUint32(pages[contextBase][0x874:], 0x0263747b)
	for index, value := range table {
		binary.LittleEndian.PutUint32(pages[tableBase][index*4:], value)
	}

	kernel := compactKernelThroughFirstStore(compact26A428Tail)
	tableIndex := uint32(87)
	for index := range 140 {
		s := newState(&model{pages: clonePages(pages)})
		s.sp = stack
		s.x[0] = uint64(leftSeed)
		s.x[1] = uint64(tableIndex - 0xf587d089)
		s.x[2] = uint64(rightSeed)
		s.x[4] = uint64(uint32(0x60319e23 + index))
		s.x[9] = 0x2d53
		s.x[10] = 0xf587d089
		s.x[12] = 0x621b97c3
		s.x[13] = pairTransformTableWords
		s.x[16] = outputBase
		s.x[27] = contextBase
		if err := runCompactKernel(s, kernel); err != nil {
			t.Fatalf("tail iteration %d: %v", index, err)
		}
		got, err := s.load(outputBase+160+uint64(index), 1)
		if err != nil || byte(got) != want[index] {
			t.Fatalf("tail iteration %d: lifted=%#02x semantic=%#02x err=%v", index, byte(got), want[index], err)
		}
		tableIndex = (tableIndex + pairTransformTableWords - 1) % pairTransformTableWords
	}
}

func TestTailTransform26A428RejectsUnknownDomain(t *testing.T) {
	if _, ok := tailTransform26A428(make([]uint32, pairTransformTableWords-1), 0, 0); ok {
		t.Fatal("accepted a short table")
	}
	if _, ok := tailTransform26A428(make([]uint32, pairTransformTableWords), 0, pairTransformTableWords); ok {
		t.Fatal("accepted an out-of-range table index")
	}
}
