package recovery

import "fmt"

const (
	accountStackStateRoute      = uint64(0x10093d000)
	accountOuterCallSetupA      = uint64(0x10098af98)
	accountSignedThresholdRoute = uint64(0x100a9bde8)
	accountOuterCallSetupB      = uint64(0x100aa1d94)
)

func runAccountStackStateRoute(s *state) (uint64, error) {
	s.traceBlock(accountStackStateRoute)
	s.x[8] = 0xffffee3c
	s.x[13] = uint64(uint32(s.x[3]) + uint32(s.x[8]))
	s.x[8] = 0x3e17
	s.x[9] = uint64(uint32(s.x[3]) + uint32(s.x[8]))
	value, err := s.load(s.sp+0x4dc, 4)
	if err != nil {
		return 0, fmt.Errorf("account stack-state word load: %w", err)
	}
	s.x[8] = uint64(uint32(value))
	value, err = s.load(s.sp+0x570, 4)
	if err != nil {
		return 0, fmt.Errorf("account stack-state saved-register load: %w", err)
	}
	s.x[19] = uint64(uint32(value))
	value, err = s.load(s.x[22]+shiftLeft(signExtend32(uint32(s.x[9])), 3, 64), 8)
	if err != nil {
		return 0, fmt.Errorf("account stack-state target load: %w", err)
	}
	s.x[9] = value
	s.x[17] = 0xe224
	return s.x[9], nil
}

func runAccountOuterCallSetupA(s *state) (uint64, error) {
	s.traceBlock(accountOuterCallSetupA)
	s.x[8] = 0xffffbfe1
	s.x[23] = uint64(uint32(s.x[3]) + uint32(s.x[8]))
	if err := s.storeValue(s.sp+0x4dc, uint64(uint32(s.x[20])), 4); err != nil {
		return 0, fmt.Errorf("account outer-call A state store: %w", err)
	}
	value, err := s.load(s.sp+0x4e0, 8)
	if err != nil {
		return 0, fmt.Errorf("account outer-call A pointer load: %w", err)
	}
	s.x[8] = value
	value, err = s.load(s.x[8], 4)
	if err != nil {
		return 0, fmt.Errorf("account outer-call A word load: %w", err)
	}
	s.x[8] = uint64(uint32(value))
	value, err = s.load(s.sp+0x574, 4)
	if err != nil {
		return 0, fmt.Errorf("account outer-call A selector load: %w", err)
	}
	s.x[9] = uint64(uint32(value))
	s.x[10] = 0x0a895c19
	s.x[9] = uint64(uint32(s.x[9]) + uint32(s.x[10]))
	s.x[9] = uint64(uint32(s.x[8]) + uint32(s.x[9]))
	s.x[8] = 0xffffcf3b
	s.x[8] = uint64(uint32(s.x[3]) + uint32(s.x[8]))
	value, err = s.load(s.x[22]+shiftLeft(signExtend32(uint32(s.x[8])), 3, 64), 8)
	if err != nil {
		return 0, fmt.Errorf("account outer-call A target load: %w", err)
	}
	s.x[10] = value
	s.x[24] = s.x[3]
	s.x[8] = s.x[1]
	s.x[17] = 0x8178
	return s.x[10], nil
}

func runAccountSignedThresholdRoute(s *state) (uint64, error) {
	s.traceBlock(accountSignedThresholdRoute)
	s.x[9] = uint64(uint32(s.x[13]) ^ uint32(0x37ec))
	s.x[14] = 0x1171
	s.x[9] = uint64(uint32(s.x[9]) * uint32(s.x[14]))
	s.x[14] = 0xba071aab
	s.compare32(uint32(s.x[8]), uint32(s.x[14]))
	s.x[14] = uint64(uint32(boolValue(s.condition("lt"))))
	s.x[15] = 0x2f57
	s.x[15] = uint64(uint32(s.x[13]) + uint32(s.x[15]))
	s.x[15] = uint64(uint32(s.x[9]) ^ uint32(s.x[15]))
	s.x[14] = uint64(uint32(s.x[15]) * uint32(s.x[14]))
	s.x[14] = uint64(uint32(s.x[14]) ^ uint32(s.x[13]))
	value, err := s.load(s.x[22]+shiftLeft(signExtend32(uint32(s.x[14])), 3, 64), 8)
	if err != nil {
		return 0, fmt.Errorf("account signed-threshold target load: %w", err)
	}
	s.x[14] = value
	s.x[17] = 0xb2e8
	return s.x[14], nil
}

func runAccountOuterCallSetupB(s *state) (uint64, error) {
	s.traceBlock(accountOuterCallSetupB)
	s.x[8] = uint64(uint32(s.x[23]) - 0x111)
	s.x[9] = 0x412a
	s.x[24] = uint64(uint32(s.x[8]) | uint32(s.x[9]))
	value, err := s.load(s.sp+0x230, 8)
	if err != nil {
		return 0, fmt.Errorf("account outer-call B pointer load: %w", err)
	}
	s.x[8] = value
	if err := s.storeValue(s.sp+0x4e0, s.x[8], 8); err != nil {
		return 0, fmt.Errorf("account outer-call B pointer store: %w", err)
	}
	value, err = s.load(s.sp+0x860, 4)
	if err != nil {
		return 0, fmt.Errorf("account outer-call B state load: %w", err)
	}
	s.x[8] = uint64(uint32(value))
	s.x[9] = 0x0a895c15
	s.x[9] = uint64(uint32(s.x[10]) + uint32(s.x[9]))
	s.x[9] = uint64(uint32(s.x[9]) + uint32(s.x[8]))
	s.x[8] = 0xffffc5a6
	s.x[8] = uint64(uint32(s.x[23]) + uint32(s.x[8]))
	s.x[8] = uint64(uint32(s.x[8]) + uint32(s.x[24]))
	value, err = s.load(s.x[22]+shiftLeft(signExtend32(uint32(s.x[8])), 3, 64), 8)
	if err != nil {
		return 0, fmt.Errorf("account outer-call B target load: %w", err)
	}
	s.x[10] = value
	s.x[8] = s.x[1]
	s.x[17] = 0x8178
	return s.x[10], nil
}
