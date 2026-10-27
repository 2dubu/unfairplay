package recovery

import "fmt"

const (
	accountLowBitTarget       = uint64(0x100240a88)
	accountFrameCompareTarget = uint64(0x100263e98)
)

func runAccountLowBitTarget(s *state) (uint64, error) {
	s.traceBlock(accountLowBitTarget)
	s.x[8] = 0xe4adc9e7
	s.x[9] = 0x1af7
	s.x[20] = uint64(uint32(s.x[28]) + uint32(s.x[8]))
	s.x[19] = uint64(uint32(s.x[20]) ^ uint32(s.x[9]))
	s.x[8] = 0xe4ae305f
	s.x[9] = uint64(uint32(s.x[11]) ^ uint32(s.x[19]))
	s.x[9] &= 1
	s.x[10] = 0xe1
	s.x[9] = uint64(uint32(s.x[9]) * uint32(s.x[10]))
	s.x[8] = uint64(uint32(s.x[28]) + uint32(s.x[8]))
	s.x[8] = uint64(uint32(s.x[9]) ^ uint32(s.x[8]))
	next, err := loadBootstrapPointer(s, s.x[27], uint32(s.x[8]))
	if err != nil {
		return 0, fmt.Errorf("account low-bit target: %w", err)
	}
	s.x[8] = next
	value, err := s.load(s.sp+0x30, 4)
	if err != nil {
		return 0, fmt.Errorf("account low-bit stack word: %w", err)
	}
	s.x[25] = uint64(uint32(value))
	s.x[15] = s.x[25]
	s.x[17] = 0xaf7d
	return next, nil
}

func runAccountFrameCompareTarget(s *state) (uint64, error) {
	s.traceBlock(accountFrameCompareTarget)
	s.x[8] = 0xffffe428
	s.x[8] = uint64(uint32(s.x[12]) + uint32(s.x[8]))
	s.x[8] = uint64(uint32(s.x[8]) | 0x18)
	s.x[9] = 0xffffddec
	s.x[10] = 0x2018
	s.x[9] = uint64(uint32(s.x[12]) + uint32(s.x[9]))
	s.x[9] = uint64(uint32(s.x[9]) | uint32(s.x[10]))
	value, err := s.load(s.x[30]+0x48, 4)
	if err != nil {
		return 0, fmt.Errorf("account frame-compare input: %w", err)
	}
	s.x[10] = uint64(uint32(value))
	s.x[11] = uint64(uint32(s.x[8]) ^ 0xce756a96)
	s.compare32(uint32(s.x[10]), uint32(s.x[11]))
	s.x[11] = uint64(uint32(boolValue(s.condition("eq"))))
	s.x[13] = uint64(uint32(s.x[9]) ^ 0xa372)
	s.x[11] = uint64(uint32(s.x[11]) * uint32(s.x[13]))
	s.x[11] = uint64(uint32(s.x[11]) ^ uint32(s.x[12]))
	next, err := loadBootstrapPointer(s, s.x[27], uint32(s.x[11]))
	if err != nil {
		return 0, fmt.Errorf("account frame-compare target: %w", err)
	}
	s.x[11] = next
	s.x[17] = 0xa968
	return next, nil
}
