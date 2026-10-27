package recovery

import (
	"encoding/binary"
	"testing"
)

func compactKernelThroughFirstStore(kernel compactKernel) compactKernel {
	for index, instruction := range kernel.instructions {
		if instruction.op == compactStoreUXTW8 {
			kernel.instructions = kernel.instructions[:index+1]
			return kernel
		}
	}
	panic("compact kernel has no byte store")
}

func TestPairTransform26A428MatchesLiftedHookDomain(t *testing.T) {
	const (
		stack      = uint64(0x7000000e0000)
		tableBase  = uint64(0x710000000000)
		outputBase = uint64(0x710000001000)
		evenSeed   = uint32(0xb7fcbf9e)
		oddSeed    = uint32(0xd439142d)
	)
	table := make([]uint32, pairTransformTableWords)
	for index := range table {
		table[index] = uint32(index)*0x9e3779b9 ^ 0xa5a55a5a
	}
	want, ok := pairTransform26A428(table, evenSeed, oddSeed, 165, 120)
	if !ok {
		t.Fatal("semantic pair transform rejected its proven domain")
	}

	pages := map[uint64][]byte{
		stack + 0x4000: make([]byte, pageSize),
		tableBase:      make([]byte, pageSize),
		outputBase:     make([]byte, pageSize),
	}
	binary.LittleEndian.PutUint64(pages[stack+0x4000][0x3d8:], tableBase)
	for index, value := range table {
		binary.LittleEndian.PutUint32(pages[tableBase][index*4:], value)
	}

	evenKernel := compactKernelThroughFirstStore(compact26A428Even)
	oddKernel := compactKernelThroughFirstStore(compact26A428Odd)
	evenIndex, oddIndex := uint32(165), uint32(120)
	for index := range 50 {
		s := newState(&model{pages: clonePages(pages)})
		s.sp = stack
		s.x[1] = 23
		s.x[10] = pairTransformTableWords
		s.x[12] = uint64(index * 2)
		s.x[14] = uint64(evenSeed)
		s.x[15] = outputBase
		s.x[22] = 0x0263747b
		s.x[23] = uint64(evenIndex - 0xf587d089)
		s.x[26] = 0xf587d089
		s.x[30] = 0x621b97c3
		if err := runCompactKernel(s, evenKernel); err != nil {
			t.Fatalf("even iteration %d: %v", index, err)
		}
		got, err := s.load(outputBase+uint64(index*2), 1)
		if err != nil || byte(got) != want[index*2] {
			t.Fatalf("even iteration %d: lifted=%#02x semantic=%#02x err=%v", index, byte(got), want[index*2], err)
		}
		evenIndex = (evenIndex + 41) % pairTransformTableWords

		s = newState(&model{pages: clonePages(pages)})
		s.sp = stack
		s.x[0] = outputBase
		s.x[3] = uint64(oddIndex - 0xf587d089)
		s.x[4] = uint64(oddSeed)
		s.x[5] = uint64(uint32(0x60319d84 + index*2))
		s.x[9] = 0xf587d089
		s.x[10] = 0x2c94
		s.x[12] = 0x621b97c3
		s.x[13] = pairTransformTableWords
		s.x[16] = 0x0263747b
		if err := runCompactKernel(s, oddKernel); err != nil {
			t.Fatalf("odd iteration %d: %v", index, err)
		}
		got, err = s.load(outputBase+uint64(index*2+1), 1)
		if err != nil || byte(got) != want[index*2+1] {
			t.Fatalf("odd iteration %d: lifted=%#02x semantic=%#02x err=%v", index, byte(got), want[index*2+1], err)
		}
		oddIndex = (2*oddIndex + 41) % pairTransformTableWords
	}
}

func TestPairTransform26A428RejectsUnknownDomain(t *testing.T) {
	if _, ok := pairTransform26A428(make([]uint32, pairTransformTableWords-1), 0, 0, 0, 0); ok {
		t.Fatal("accepted a short table")
	}
	if _, ok := pairTransform26A428(make([]uint32, pairTransformTableWords), 0, 0, pairTransformTableWords, 0); ok {
		t.Fatal("accepted an out-of-range table index")
	}
}

func clonePages(input map[uint64][]byte) map[uint64][]byte {
	output := make(map[uint64][]byte, len(input))
	for address, page := range input {
		output[address] = append([]byte(nil), page...)
	}
	return output
}
