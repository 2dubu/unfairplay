package recovery

import "fmt"

const (
	accountTripleCompareRoute = uint64(0x1002a7edc)
	accountDirectReturn       = uint64(0x100111414)
	accountZeroByteRoute      = uint64(0x10011a26c)
	accountLimitTargetRoute   = uint64(0x10019e4f0)
)

func runAccountTripleCompareRoute(s *state) (uint64, error) {
	s.traceBlock(accountTripleCompareRoute)
	s.x[15] = uint64(uint32(s.x[16]) ^ 0x3996)
	s.x[17] = 0x164f
	s.x[15] = uint64(uint32(s.x[15]) * uint32(s.x[17]))
	s.x[17] = 0xac3e9be8
	s.x[15] = uint64(uint32(s.x[15]) ^ uint32(s.x[17]))
	s.x[17] = uint64(uint32(s.x[1]) + uint32(s.x[13]))
	s.x[15] = uint64(uint32(s.x[17]) + uint32(s.x[15]))
	s.x[17] = uint64(uint32(s.x[9]) + 0x766ea525)
	s.x[2] = 0xd2f3335c
	s.compare32(uint32(s.x[15]), uint32(s.x[2]))
	s.x[2] = uint64(uint32(boolValue(s.condition("lo"))))
	s.x[3] = uint64(uint32(s.x[16]) + 0xd2f2f9c8)
	s.compare32(uint32(s.x[17]), uint32(s.x[3]))
	s.x[3] = uint64(uint32(boolValue(s.condition("lo"))))
	s.x[3] = uint64(uint32(s.x[3]) ^ uint32(s.x[2]))
	s.compare32(uint32(s.x[15]), uint32(s.x[17]))
	s.x[15] = uint64(uint32(boolValue(s.condition("hi"))))
	s.compare32(uint32(s.x[3]), 0)
	if s.condition("ne") {
		s.x[15] = uint64(uint32(s.x[2]))
	}
	s.x[15] = uint64(uint32(s.x[15]) ^ 1)
	s.x[17] = 0x4d
	s.x[15] = uint64(uint32(s.x[15]) * uint32(s.x[17]))
	s.x[15] = uint64(uint32(s.x[15]) ^ uint32(s.x[16]))
	next, err := loadBootstrapPointer(s, s.x[10], uint32(s.x[15]))
	if err != nil {
		return 0, fmt.Errorf("account triple-compare target: %w", err)
	}
	s.x[15] = next
	s.x[2] = 0x318a201b
	s.x[17] = 0x5f28
	return next, nil
}

func runAccountDirectReturn(s *state) (uint64, error) {
	return runAccountDirectReturnAt(s, accountDirectReturn)
}

func runAccountDirectReturnAt(s *state, pc uint64) (uint64, error) {
	s.traceBlock(pc)
	s.x[0] = s.x[8]
	s.sp += 0x10
	return s.x[30], nil
}

func runAccountZeroByteRoute(s *state) (uint64, error) {
	s.traceBlock(accountZeroByteRoute)
	s.x[8] = 0
	s.x[12] = s.x[0] + 4
	if err := s.storeValue(s.x[12]+uint64(uint32(s.x[9])), 0, 1); err != nil {
		return 0, fmt.Errorf("account zero-byte output: %w", err)
	}
	if err := s.storeValue(s.x[0], uint64(uint32(s.x[9])), 4); err != nil {
		return 0, fmt.Errorf("account zero-byte count: %w", err)
	}
	value, err := s.load(s.x[10]+uint64(uint32(s.x[11]))*8, 8)
	if err != nil {
		return 0, fmt.Errorf("account zero-byte target: %w", err)
	}
	s.x[9] = value
	s.x[17] = 0x30d6
	return value, nil
}

func runAccountLimitTargetRoute(s *state) (uint64, error) {
	s.traceBlock(accountLimitTargetRoute)
	s.x[12] = 0x2d52
	s.x[13] = 0xce757bcf
	s.compare32(uint32(s.x[8]), uint32(s.x[13]))
	s.x[13] = uint64(uint32(boolValue(s.condition("lt"))))
	s.x[14] = 0x58
	s.x[13] = uint64(uint32(s.x[13]) * uint32(s.x[14]))
	s.x[14] = 0x968d
	s.x[13] = uint64(uint32(s.x[13]) ^ uint32(s.x[14]))
	value, err := s.load(s.x[10]+uint64(uint32(s.x[13]))*8, 8)
	if err != nil {
		return 0, fmt.Errorf("account limit-target target: %w", err)
	}
	s.x[13] = value
	s.x[17] = 0xcf72
	return value, nil
}
