package recovery

import "fmt"

const (
	accountMaskSideTarget   = uint64(0x100265918)
	accountSavedWordTarget  = uint64(0x100265dbc)
	accountArithmeticTarget = uint64(0x1002b0c84)
	accountSavedPairTarget  = uint64(0x1002eee34)
)

func runAccountMaskSideTarget(s *state) (uint64, error) {
	s.traceBlock(accountMaskSideTarget)
	s.x[11] = 0
	s.x[8] = 0x83ab8b0a
	s.x[9] = 0x7c54dfd5
	s.x[8] = uint64(uint32(s.x[10]) + uint32(s.x[8]))
	s.x[9] = uint64(uint32(s.x[8]) & uint32(s.x[9]))
	value, err := s.load(s.x[27]+uint64(uint32(s.x[10]))*8, 8)
	if err != nil {
		return 0, fmt.Errorf("account mask-side target: %w", err)
	}
	s.x[8] = value
	s.x[17] = 0xdd9e
	return value, nil
}

func runAccountSavedWordTarget(s *state) (uint64, error) {
	s.traceBlock(accountSavedWordTarget)
	s.x[9] = 0xe4addce4
	value, err := s.load(s.sp+0x68, 4)
	if err != nil {
		return 0, fmt.Errorf("account saved-word input: %w", err)
	}
	s.x[8] = uint64(uint32(value))
	s.x[10] = uint64(uint32(s.x[8]) & uint32(s.x[11]))
	value, err = s.load(s.x[30]+0x48, 4)
	if err != nil {
		return 0, fmt.Errorf("account saved-word frame input: %w", err)
	}
	s.x[8] = uint64(uint32(value))
	s.x[10] = uint64(uint32(s.x[10]) ^ 1)
	s.x[11] = uint64(uint32(s.x[28]) + 0xe4ad8da3)
	s.x[10] = uint64(uint32(s.x[10]) * uint32(s.x[11]))
	s.x[9] = uint64(uint32(s.x[28]) + uint32(s.x[9]))
	s.x[9] = uint64(uint32(s.x[10]) ^ uint32(s.x[9]))
	next, err := loadBootstrapPointer(s, s.x[27], uint32(s.x[9]))
	if err != nil {
		return 0, fmt.Errorf("account saved-word target: %w", err)
	}
	s.x[9] = next
	value, err = s.load(s.sp+0x5c, 4)
	if err != nil {
		return 0, fmt.Errorf("account saved-word companion: %w", err)
	}
	s.x[19] = uint64(uint32(value))
	s.x[13] = s.x[19]
	s.x[17] = 0xb6f9
	return next, nil
}

func runAccountArithmeticTarget(s *state) (uint64, error) {
	s.traceBlock(accountArithmeticTarget)
	s.x[8] = 0x602f
	s.x[9] = uint64(uint32(s.x[10]) ^ uint32(s.x[8]))
	s.x[12] = 0x2d15
	s.x[8] = 0x5d1
	s.x[8] = uint64(uint32(s.x[9]) * uint32(s.x[8]))
	s.x[13] = 0xce753d43
	s.x[12] = uint64(uint32(s.x[9]) + uint32(s.x[12]))
	s.x[9] = uint64(uint32(s.x[12]) + uint32(s.x[13]))
	s.x[9] = uint64(uint32(s.x[9]) + uint32(s.x[8]))
	s.compare32(uint32(s.x[11]), uint32(s.x[9]))
	s.x[9] = uint64(uint32(boolValue(s.condition("eq"))))
	s.x[11] = 0xd276
	s.x[9] = uint64(uint32(s.x[9]) * uint32(s.x[11]))
	s.x[9] = uint64(uint32(s.x[9]) ^ uint32(s.x[10]))
	next, err := loadBootstrapPointer(s, s.x[27], uint32(s.x[9]))
	if err != nil {
		return 0, fmt.Errorf("account arithmetic target: %w", err)
	}
	s.x[9] = next
	s.x[17] = 0x2de9
	return next, nil
}

func runAccountSavedPairTarget(s *state) (uint64, error) {
	s.traceBlock(accountSavedPairTarget)
	if err := s.storeValue(s.sp+0x5c, uint64(uint32(s.x[15])), 4); err != nil {
		return 0, fmt.Errorf("account saved-pair first word: %w", err)
	}
	if err := s.storeValue(s.sp+0x30, uint64(uint32(s.x[25])), 4); err != nil {
		return 0, fmt.Errorf("account saved-pair second word: %w", err)
	}
	s.x[9] = 0x1b02428d
	s.x[10] = 0x00502001
	s.x[9] = uint64(uint32(s.x[13]) + uint32(s.x[9]))
	s.x[12] = uint64(uint32(s.x[9]) | uint32(s.x[10]))
	s.x[9] = 0x9b5beb5e
	s.x[10] = 0x7ff6772f
	s.x[9] = uint64(uint32(s.x[13]) + uint32(s.x[9]))
	s.x[28] = uint64(uint32(s.x[9]) & uint32(s.x[10]))
	s.x[9] = 0xce757bcf
	s.compare32(uint32(s.x[8]), uint32(s.x[9]))
	s.x[8] = uint64(uint32(boolValue(s.condition("eq"))))
	s.x[9] = uint64(uint32(s.x[12]) ^ 0xfffc2833)
	s.x[9] = uint64(uint32(s.x[9]) + uint32(s.x[28]))
	s.x[8] = uint64(uint32(s.x[9]) * uint32(s.x[8]))
	s.x[8] = uint64(uint32(s.x[8]) ^ uint32(s.x[13]))
	next, err := loadBootstrapPointer(s, s.x[27], uint32(s.x[8]))
	if err != nil {
		return 0, fmt.Errorf("account saved-pair target: %w", err)
	}
	s.x[8] = next
	s.x[11] = 1
	s.x[17] = 0x6f0a
	return next, nil
}
