package recovery

import "fmt"

const (
	accountSavedCompareDispatch = uint64(0x100300904)
	accountPairCompareDispatch  = uint64(0x1003157f8)
	accountPairCompareTail      = uint64(0x100315800)
)

func runAccountSavedCompareDispatch(s *state) (uint64, error) {
	s.traceBlock(accountSavedCompareDispatch)
	if err := s.storeValue(s.sp+0x34, uint64(uint32(s.x[13])), 4); err != nil {
		return 0, fmt.Errorf("account saved-compare word: %w", err)
	}
	s.x[9] = 0x682442f2
	s.x[10] = 0x7c897df8
	s.x[9] = uint64(uint32(s.x[28]) + uint32(s.x[9]))
	s.x[10] = uint64(uint32(s.x[9]) & uint32(s.x[10]))
	s.x[9] = 0xe4ada466
	s.x[11] = 0xce7546de
	s.x[11] = uint64(uint32(s.x[10]) + uint32(s.x[11]))
	s.compare32(uint32(s.x[8]), uint32(s.x[11]))
	s.x[8] = uint64(uint32(boolValue(s.condition("eq"))))
	s.x[11] = 0x158
	s.x[8] = uint64(uint32(s.x[8]) * uint32(s.x[11]))
	s.x[9] = uint64(uint32(s.x[28]) + uint32(s.x[9]))
	s.x[8] = uint64(uint32(s.x[8]) ^ uint32(s.x[9]))
	next, err := loadBootstrapPointer(s, s.x[27], uint32(s.x[8]))
	if err != nil {
		return 0, fmt.Errorf("account saved-compare target: %w", err)
	}
	s.x[8] = next
	s.x[17] = 0xb1c4
	return next, nil
}

func runAccountPairCompareDispatch(s *state) (uint64, error) {
	s.traceBlock(accountPairCompareDispatch)
	if err := s.storeValue(s.sp+0x2c, uint64(uint32(s.x[9])), 4); err != nil {
		return 0, fmt.Errorf("account pair-compare first word: %w", err)
	}
	if err := s.storeValue(s.sp+0x30, uint64(uint32(s.x[10])), 4); err != nil {
		return 0, fmt.Errorf("account pair-compare second word: %w", err)
	}
	value, err := s.load(s.x[30]+0x48, 4)
	if err != nil {
		return 0, fmt.Errorf("account pair-compare input: %w", err)
	}
	s.x[8] = uint64(uint32(value))
	s.traceBlock(accountPairCompareTail)
	s.x[9] = 0xce757bc6
	s.compare32(uint32(s.x[8]), uint32(s.x[9]))
	s.x[8] = uint64(uint32(boolValue(s.condition("eq"))))
	s.x[9] = 0x5dc7
	s.x[9] = uint64(uint32(s.x[8]) * uint32(s.x[9]))
	s.x[8] = uint64(uint32(s.x[28]) + 0xe4ad9d7c)
	s.x[9] = uint64(uint32(s.x[9]) ^ uint32(s.x[8]))
	next, err := loadBootstrapPointer(s, s.x[27], uint32(s.x[9]))
	if err != nil {
		return 0, fmt.Errorf("account pair-compare target: %w", err)
	}
	s.x[9] = next
	s.x[11] = 1
	s.x[17] = 0x6f0a
	return next, nil
}
