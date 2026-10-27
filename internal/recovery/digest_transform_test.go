package recovery

import (
	"encoding/binary"
	"testing"
)

func compactKernelThroughSPStores(kernel compactKernel, count int) compactKernel {
	for index, instruction := range kernel.instructions {
		if instruction.op != compactStoreSP32 {
			continue
		}
		count--
		if count == 0 {
			kernel.instructions = kernel.instructions[:index+1]
			return kernel
		}
	}
	panic("compact kernel has too few stack stores")
}

func TestDigestLanes26A428MatchLiftedKernel(t *testing.T) {
	const (
		stack  = uint64(0x7000000e0000)
		tableA = uint64(0x710000000000)
		tableB = uint64(0x710000001000)
	)
	pages := map[uint64][]byte{
		stack:  make([]byte, pageSize),
		tableA: make([]byte, pageSize),
		tableB: make([]byte, pageSize),
	}
	for index := range 256 {
		pages[tableA][index] = byte(index*73 + 19)
		pages[tableB][index] = byte(index*29 + 7)
	}
	binary.LittleEndian.PutUint32(pages[stack][0xd48:], 0x13579bdf)
	binary.LittleEndian.PutUint32(pages[stack][0xd50:], 0x2468ace0)
	binary.LittleEndian.PutUint32(pages[stack][0xd58:], 0x0badf00d)

	s := newState(&model{pages: clonePages(pages)})
	s.sp = stack
	s.x[5] = 0x10203040
	s.x[8] = 0x50607080
	s.x[9] = 0x4ec4ec4f
	s.x[11] = 6
	s.x[12] = 13
	s.x[13] = tableA
	s.x[14] = 2
	s.x[16] = 0x14
	s.x[17] = 7
	s.x[30] = tableB

	a := uint32(pages[tableA][7])
	bIndex := (binary.LittleEndian.Uint32(pages[stack][0xd58:]) + uint32(s.x[8])) % uint32(s.x[12])
	b := uint32(pages[tableA][bIndex])
	cIndex := (binary.LittleEndian.Uint32(pages[stack][0xd50:]) + uint32(s.x[8])) % uint32(s.x[12])
	c := uint32(pages[tableA][cIndex])
	want := digestLanes26A428(a, b, c, uint32(s.x[16]))

	if err := runCompactKernel(s, compactKernelThroughSPStores(compact26A428Digest, 4)); err != nil {
		t.Fatal(err)
	}
	for index, offset := range []int{0xd78, 0xd80, 0xd90, 0xd98} {
		got, err := s.load(stack+uint64(offset), 4)
		if err != nil || uint32(got) != want[index] {
			t.Fatalf("lane %d: lifted=%#08x semantic=%#08x err=%v", index, uint32(got), want[index], err)
		}
	}
}
