package recovery

import "fmt"

const (
	accountSmallStoreReturn = uint64(0x10032feb4)
	accountEqualityTarget   = uint64(0x100049f84)
	accountBitWordTarget    = uint64(0x1000550a4)
)

func runAccountSmallStoreReturn(s *state) (uint64, error) {
	s.traceBlock(accountSmallStoreReturn)
	s.x[8] = 0x79
	if err := s.storeValue(s.x[0]+0x58, uint64(byte(s.x[8])), 1); err != nil {
		return 0, fmt.Errorf("account small-store byte: %w", err)
	}
	value, err := s.load(s.x[1]+4, 4)
	if err != nil {
		return 0, fmt.Errorf("account small-store word: %w", err)
	}
	s.x[8] = uint64(uint32(value))
	if err := s.storeValue(s.x[0]+0x48, uint64(uint32(s.x[8])), 4); err != nil {
		return 0, fmt.Errorf("account small-store word output: %w", err)
	}
	value, err = s.load(s.x[1]+8, 8)
	if err != nil {
		return 0, fmt.Errorf("account small-store pointer: %w", err)
	}
	s.x[8] = value
	if err := s.storeValue(s.x[0]+0x50, s.x[8], 8); err != nil {
		return 0, fmt.Errorf("account small-store pointer output: %w", err)
	}
	s.x[0] = 0
	return s.x[30], nil
}

func runAccountEqualityTarget(s *state) (uint64, error) {
	s.traceBlock(accountEqualityTarget)
	s.compare32(uint32(s.x[10]), uint32(s.x[1]))
	s.x[12] = uint64(uint32(boolValue(s.condition("eq"))))
	s.x[8] = 0xebb80623
	s.x[8] = uint64(uint32(s.x[9]) + uint32(s.x[8]))
	s.x[13] = uint64(uint32(s.x[8]) ^ 0x9bb6)
	s.x[12] = uint64(uint32(s.x[12]) * uint32(s.x[13]))
	s.x[12] = uint64(uint32(s.x[12]) ^ uint32(s.x[8]))
	next, err := loadBootstrapPointer(s, s.x[20], uint32(s.x[12]))
	if err != nil {
		return 0, fmt.Errorf("account equality target: %w", err)
	}
	s.x[13] = next
	s.x[28] = s.x[11]
	s.x[23] = s.x[9]
	s.x[12] = s.x[10]
	s.x[17] = 0x85f4
	return next, nil
}

func runAccountBitWordTarget(s *state) (uint64, error) {
	s.traceBlock(accountBitWordTarget)
	s.x[25] = 0x80000004
	s.x[21] = 0x7c79d912
	s.x[22] = 0xfc79d912
	s.x[9] = s.x[10] + (s.x[10] << 1)
	s.x[15] = s.x[26] + (s.x[9] << 4)
	value, err := s.load(s.x[15]+0x28, 4)
	if err != nil {
		return 0, fmt.Errorf("account bit-word input: %w", err)
	}
	s.x[10] = uint64(uint32(value))
	s.test32(uint32(s.x[10]) & 1)
	s.x[11] = uint64(uint32(boolValue(s.condition("eq"))))
	s.x[9] = 0xffff641e
	s.x[9] = uint64(uint32(s.x[12]) + uint32(s.x[9]))
	s.x[13] = uint64(uint32(s.x[9]) ^ 0x1907)
	s.x[11] = uint64(uint32(s.x[11]) * uint32(s.x[13]))
	s.x[11] = uint64(uint32(s.x[11]) ^ uint32(s.x[12]))
	next, err := loadBootstrapPointer(s, s.x[20], uint32(s.x[11]))
	if err != nil {
		return 0, fmt.Errorf("account bit-word target: %w", err)
	}
	s.x[11] = next
	if err := s.storeValue(s.sp+0x10, s.x[15], 8); err != nil {
		return 0, fmt.Errorf("account bit-word saved pointer: %w", err)
	}
	s.x[17] = 0xcd17
	return next, nil
}
