package recovery

import (
	"encoding/binary"
	"testing"
)

func TestWriterTransform26A428MatchesLiftedHookDomain(t *testing.T) {
	const (
		stack       = uint64(0x7000000e0000)
		tableBase   = uint64(0x710000000000)
		outputBase  = uint64(0x710000001000)
		contextBase = uint64(0x710000002000)
		seed        = uint32(0x600767c4)
	)
	table := make([]uint32, pairTransformTableWords)
	for index := range table {
		table[index] = uint32(index)*0x85ebca6b ^ 0xc2b2ae35
	}
	controls := make([]uint32, 140)
	for index := range controls {
		controls[index] = uint32(index)*0x27d4eb2d ^ 0x165667b1
	}
	want, ok := writerTransform26A428(table, controls, seed, 129)
	if !ok {
		t.Fatal("semantic writer transform rejected its proven domain")
	}

	pages := map[uint64][]byte{
		stack + 0x4000:       make([]byte, pageSize),
		tableBase:            make([]byte, pageSize),
		outputBase:           make([]byte, pageSize),
		contextBase:          make([]byte, pageSize),
		contextBase + 0x1000: make([]byte, pageSize),
	}
	binary.LittleEndian.PutUint64(pages[stack+0x4000][0x3d8:], tableBase)
	binary.LittleEndian.PutUint32(pages[contextBase+0x1000][0x1b4:], 0x0263747b)
	for index, value := range table {
		binary.LittleEndian.PutUint32(pages[tableBase][index*4:], value)
	}

	kernel := compactKernelThroughFirstStore(compact26A428Writer)
	tableIndex := uint32(129)
	for index, control := range controls {
		s := newState(&model{pages: clonePages(pages)})
		s.sp = stack
		s.x[0] = uint64(seed)
		s.x[1] = uint64(uint32(0x60319dbf + index))
		s.x[2] = uint64(tableIndex - 0xf587d089)
		s.x[9] = 0x30b1
		s.x[11] = outputBase
		s.x[12] = pairTransformTableWords
		s.x[15] = 0x621b97c3
		s.x[16] = uint64(control)
		s.x[27] = contextBase
		if err := runCompactKernel(s, kernel); err != nil {
			t.Fatalf("writer iteration %d: %v", index, err)
		}
		got, err := s.load(outputBase+60+uint64(index), 1)
		if err != nil || byte(got) != want[index] {
			t.Fatalf("writer iteration %d: lifted=%#02x semantic=%#02x err=%v", index, byte(got), want[index], err)
		}
		tableIndex = (tableIndex + 36) % pairTransformTableWords
	}
}

func TestWriterTransform26A428RejectsUnknownDomain(t *testing.T) {
	if _, ok := writerTransform26A428(make([]uint32, pairTransformTableWords), make([]uint32, 139), 0, 0); ok {
		t.Fatal("accepted a short control vector")
	}
	if _, ok := writerTransform26A428(make([]uint32, pairTransformTableWords), make([]uint32, 140), 0, pairTransformTableWords); ok {
		t.Fatal("accepted an out-of-range table index")
	}
}
