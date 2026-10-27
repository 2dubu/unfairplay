package recovery

import "fmt"

const (
	bootstrapDescriptorCall = uint64(0x10081db08)
	bootstrapNestedCallA    = uint64(0x10084fa0c)
	bootstrapNestedCallB    = uint64(0x10084fb14)
)

func runBootstrapDescriptorCall(s *state) (uint64, error) {
	s.traceBlock(bootstrapDescriptorCall)
	s.x[8] = 0xffff99d7
	s.x[9] = s.x[29] - 0x78
	s.x[9] = uint64(uint32(s.x[9]) ^ 0xf9d49974)
	s.x[10] = 0xffffc0c7
	s.x[9] = uint64(uint32(s.x[9]) * 0xda51a827)
	s.x[10] = uint64(uint32(s.x[23]) + uint32(s.x[10]))
	s.x[9] = uint64(uint32(s.x[10]) ^ uint32(s.x[9]))
	if err := s.storeValue(addOffset(s.x[29], -0x68), uint64(uint32(s.x[9])), 4); err != nil {
		return 0, fmt.Errorf("bootstrap descriptor-call state: %w", err)
	}
	value, err := s.load(s.sp+0x18, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap descriptor-call stack value: %w", err)
	}
	s.x[10] = value
	s.x[9] = s.sp + 0xb0
	if err := storeBootstrapPair(s, addOffset(s.x[29], -0x78), s.x[9], s.x[10]); err != nil {
		return 0, fmt.Errorf("bootstrap descriptor-call pair: %w", err)
	}
	first, err := s.load(s.x[24], 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap descriptor-call first word: %w", err)
	}
	second, err := s.load(s.x[26], 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap descriptor-call second word: %w", err)
	}
	s.x[9] = uint64(uint32(first) ^ uint32(second) ^ 0xa6)
	s.x[9] = uint64(uint32(s.x[9]) * 0x0a5cd4cf)
	s.x[10] = s.x[27] + uint64(byte(s.x[9]))
	value, err = s.load(s.x[10], 1)
	if err != nil {
		return 0, fmt.Errorf("bootstrap descriptor-call first byte: %w", err)
	}
	s.x[10] = uint64(byte(value) ^ 0x5d)
	value, err = s.load(s.x[28]+uint64(uint32(s.x[10])), 1)
	if err != nil {
		return 0, fmt.Errorf("bootstrap descriptor-call second byte: %w", err)
	}
	s.x[10] = uint64(byte(value))
	s.x[19] = uint64(uint32(s.x[23]) + uint32(s.x[8]))
	s.x[8] = uint64(uint32(s.x[9]) ^ uint32(s.x[10]))
	s.x[9] = 0xcfa4
	s.x[8] = uint64(uint32(s.x[9]) + uint32(byte(s.x[8])))
	next, err := loadBootstrapPointer(s, s.x[25], uint32(s.x[8]))
	if err != nil {
		return 0, fmt.Errorf("bootstrap descriptor-call target: %w", err)
	}
	s.x[8] = next
	s.x[0] = s.x[29] - 0x78
	s.x[17] = 0x120f
	s.x[30] = bootstrapCallContinuation
	return next, nil
}

func runBootstrapNestedCallA(s *state) (uint64, error) {
	s.traceBlock(bootstrapNestedCallA)
	s.x[8] = uint64(uint32(s.x[19]) + 0xffff727c)
	s.x[21] = uint64(uint32(s.x[8]) | 0x434a)
	s.x[8] = uint64(uint32(s.x[19]) ^ 0xb299)
	s.x[9] = uint64(uint32(s.x[9]) + 0x2578d536)
	if err := s.storeValue(s.sp+0xb4, uint64(uint32(s.x[9])), 4); err != nil {
		return 0, fmt.Errorf("bootstrap nested-call A selector: %w", err)
	}
	s.x[9] = uint64(uint32(s.x[19]) - 0xf26)
	s.x[22] = s.x[29] - 0x78
	s.x[10] = 0x944c37a9 &^ uint64(uint32(s.x[22]))
	s.x[11] = uint64(uint32(s.x[22]) & 0x6bb3c850)
	s.x[10] = uint64(uint32(s.x[10]) | uint32(s.x[11]))
	s.x[10] = uint64(uint32(s.x[10]) ^ 0x6d98aedd)
	s.x[10] = uint64(uint32(s.x[10]) * 0xda51a827)
	s.x[11] = s.sp + 0xa0
	value, err := s.load(s.sp+0x98, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap nested-call A saved value: %w", err)
	}
	s.x[27] = value
	if err := storeBootstrapPair(s, addOffset(s.x[29], -0x78), s.x[27], s.x[11]); err != nil {
		return 0, fmt.Errorf("bootstrap nested-call A pair: %w", err)
	}
	s.x[9] = uint64(uint32(s.x[9]) ^ uint32(s.x[10]))
	if err := s.storeValue(addOffset(s.x[29], -0x68), uint64(uint32(s.x[9])), 4); err != nil {
		return 0, fmt.Errorf("bootstrap nested-call A state: %w", err)
	}
	value, err = s.load(s.sp+0x20, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap nested-call A register 26: %w", err)
	}
	s.x[26] = value
	value, err = s.load(s.sp+0x28, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap nested-call A register 24: %w", err)
	}
	s.x[24] = value
	first, err := s.load(s.x[24], 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap nested-call A first word: %w", err)
	}
	second, err := s.load(s.x[26], 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap nested-call A second word: %w", err)
	}
	s.x[9] = uint64(uint32(first) ^ uint32(second) ^ 0xa6)
	s.x[28] = 0x0a5cd4cf
	s.x[9] = uint64(uint32(s.x[9]) * uint32(s.x[28]))
	const tableRoot = uint64(0x10105d1e8)
	s.x[20] = tableRoot
	value, err = s.load(tableRoot+0x5c58, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap nested-call A first table: %w", err)
	}
	s.x[10] = value + 0xffffffffed6f2dad
	if err := s.storeValue(s.sp+0x80, s.x[10], 8); err != nil {
		return 0, fmt.Errorf("bootstrap nested-call A first table state: %w", err)
	}
	s.x[10] += uint64(byte(s.x[9]))
	value, err = s.load(s.x[10], 1)
	if err != nil {
		return 0, fmt.Errorf("bootstrap nested-call A first byte: %w", err)
	}
	s.x[10] = uint64(byte(value) ^ 0x29)
	value, err = s.load(tableRoot+0x8580, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap nested-call A second table: %w", err)
	}
	s.x[11] = value + 0xffffffff997cf512
	if err := s.storeValue(s.sp+0x78, s.x[11], 8); err != nil {
		return 0, fmt.Errorf("bootstrap nested-call A second table state: %w", err)
	}
	value, err = s.load(s.x[11]+uint64(uint32(s.x[10])), 1)
	if err != nil {
		return 0, fmt.Errorf("bootstrap nested-call A second byte: %w", err)
	}
	s.x[10] = uint64(byte(value))
	s.x[23] = uint64(uint32(s.x[8]) + uint32(s.x[8]<<1))
	s.x[8] = uint64(uint32(s.x[10]) ^ uint32(s.x[9]))
	s.x[9] = 0xcef0
	s.x[8] = uint64(uint32(s.x[9]) + uint32(byte(s.x[8])))
	next, err := loadBootstrapPointer(s, s.x[25], uint32(s.x[8]))
	if err != nil {
		return 0, fmt.Errorf("bootstrap nested-call A target: %w", err)
	}
	s.x[8] = next
	s.x[0] = s.x[29] - 0x78
	s.x[17] = 0x120f
	s.x[30] = bootstrapNestedCallB
	return next, nil
}

func runBootstrapNestedCallB(s *state) (uint64, error) {
	s.traceBlock(bootstrapNestedCallB)
	s.x[8] = uint64(uint32(s.x[22]) & 0x60e71988)
	s.x[8] = uint64(uint32(s.x[8] << 1))
	s.x[8] = uint64(uint32(s.x[8]) - uint32(s.x[22]))
	s.x[8] = uint64(uint32(s.x[8]) + 0x9f18e673)
	s.x[8] = uint64(uint32(s.x[8]) ^ 0x782a3411)
	s.x[10] = 0xba209ef7
	if err := s.storeValue(addOffset(s.x[29], -0x68), s.x[27], 8); err != nil {
		return 0, fmt.Errorf("bootstrap nested-call B saved value: %w", err)
	}
	s.x[9] = s.sp + 0xb0
	if err := s.storeValue(addOffset(s.x[29], -0x78), s.x[9], 8); err != nil {
		return 0, fmt.Errorf("bootstrap nested-call B stack pointer: %w", err)
	}
	s.x[8] = uint64(uint32(s.x[8])*uint32(s.x[10]) + uint32(s.x[19]))
	s.x[9] = 0xffffb18a
	s.x[8] = uint64(uint32(s.x[8]) + uint32(s.x[9]))
	if err := s.storeValue(addOffset(s.x[29], -0x6c), uint64(uint32(s.x[8])), 4); err != nil {
		return 0, fmt.Errorf("bootstrap nested-call B state: %w", err)
	}
	first, err := s.load(s.x[24], 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap nested-call B first word: %w", err)
	}
	second, err := s.load(s.x[26], 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap nested-call B second word: %w", err)
	}
	s.x[8] = uint64(uint32(first) ^ uint32(second) ^ 0xa6)
	s.x[8] = uint64(uint32(s.x[8]) * uint32(s.x[28]))
	value, err := s.load(s.x[20], 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap nested-call B first table: %w", err)
	}
	s.x[9] = value + uint64(byte(s.x[8])) + 0xffffffffcd114862
	value, err = s.load(s.x[9], 1)
	if err != nil {
		return 0, fmt.Errorf("bootstrap nested-call B first byte: %w", err)
	}
	s.x[9] = uint64(byte(value) ^ 0xcd)
	value, err = s.load(s.x[20]+0xf60, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap nested-call B second table: %w", err)
	}
	s.x[10] = value + 0xffffffffe5b0614e
	s.x[9] += s.x[10]
	value, err = s.load(s.x[9], 1)
	if err != nil {
		return 0, fmt.Errorf("bootstrap nested-call B second byte: %w", err)
	}
	s.x[9] = uint64(byte(value))
	s.x[8] = uint64(uint32(s.x[9]) ^ uint32(s.x[8]))
	s.x[9] = 0xcf67
	s.x[8] = uint64(uint32(s.x[9]) + uint32(byte(s.x[8])))
	next, err := loadBootstrapPointer(s, s.x[25], uint32(s.x[8]))
	if err != nil {
		return 0, fmt.Errorf("bootstrap nested-call B target: %w", err)
	}
	s.x[8] = next
	s.x[0] = s.x[29] - 0x78
	s.x[17] = 0x425b
	s.x[30] = bootstrapNestedReturnStart
	return next, nil
}

func storeBootstrapPair(s *state, address, first, second uint64) error {
	if err := s.storeValue(address, first, 8); err != nil {
		return err
	}
	return s.storeValue(address+8, second, 8)
}
