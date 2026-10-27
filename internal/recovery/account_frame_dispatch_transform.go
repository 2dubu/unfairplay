package recovery

import "fmt"

const (
	accountCompareDispatchA = uint64(0x10015a8b4)
	accountFrameDispatch    = uint64(0x10015b824)
	accountFrameOffset      = uint64(0x10015c8c8)
	accountCompareDispatchB = uint64(0x100177e78)
)

func runAccountCompareDispatchA(s *state) (uint64, error) {
	s.traceBlock(accountCompareDispatchA)
	s.x[8] = uint64(uint32(s.x[9]) ^ 0x3a3d)
	s.x[11] = 0x4cef
	s.x[12] = 0xce757bcf
	s.compare32(uint32(s.x[10]), uint32(s.x[12]))
	s.x[12] = uint64(uint32(boolValue(s.condition("ne"))))
	s.x[13] = 0x5e61
	s.x[10] = uint64(uint32(s.x[8]) + uint32(s.x[11]))
	s.x[11] = uint64(uint32(s.x[10]) ^ uint32(s.x[13]))
	s.x[11] = uint64(uint32(s.x[12]) * uint32(s.x[11]))
	s.x[9] = uint64(uint32(s.x[11]) ^ uint32(s.x[9]))
	next, err := loadBootstrapPointer(s, s.x[27], uint32(s.x[9]))
	if err != nil {
		return 0, fmt.Errorf("account compare-dispatch A target: %w", err)
	}
	s.x[9] = next
	s.x[17] = 0x269d
	return next, nil
}

func runAccountFrameDispatch(s *state) (uint64, error) {
	s.traceBlock(accountFrameDispatch)
	if err := s.storeValue(s.sp+0x68, uint64(uint32(s.x[11])), 4); err != nil {
		return 0, fmt.Errorf("account frame-dispatch saved word: %w", err)
	}
	s.x[8] = 0xe4ad9d69
	value, err := s.load(s.x[30]+0x48, 4)
	if err != nil {
		return 0, fmt.Errorf("account frame-dispatch input: %w", err)
	}
	s.x[9] = uint64(uint32(value))
	s.x[10] = 0xce757bcf
	s.compare32(uint32(s.x[9]), uint32(s.x[10]))
	s.x[9] = uint64(uint32(boolValue(s.condition("eq"))))
	s.x[10] = 0xac1a
	s.x[9] = uint64(uint32(s.x[9]) * uint32(s.x[10]))
	s.x[8] = uint64(uint32(s.x[28]) + uint32(s.x[8]))
	s.x[9] = uint64(uint32(s.x[9]) ^ uint32(s.x[8]))
	next, err := loadBootstrapPointer(s, s.x[27], uint32(s.x[9]))
	if err != nil {
		return 0, fmt.Errorf("account frame-dispatch target: %w", err)
	}
	s.x[9] = next
	s.x[11] = 1
	s.x[17] = 0x6f0a
	return next, nil
}

func runAccountFrameOffset(s *state) (uint64, error) {
	s.traceBlock(accountFrameOffset)
	s.x[8] = 0xffffea5c
	s.x[13] = uint64(uint32(s.x[19]) + uint32(s.x[8]))
	value, err := s.load(s.x[30]+0x48, 4)
	if err != nil {
		return 0, fmt.Errorf("account frame-offset input: %w", err)
	}
	s.x[8] = uint64(uint32(value))
	next, err := loadBootstrapPointer(s, s.x[27], uint32(s.x[19]))
	if err != nil {
		return 0, fmt.Errorf("account frame-offset target: %w", err)
	}
	s.x[9] = next
	s.x[17] = 0xfed4
	return next, nil
}

func runAccountCompareDispatchB(s *state) (uint64, error) {
	s.traceBlock(accountCompareDispatchB)
	value, err := s.load(s.x[30]+0x48, 4)
	if err != nil {
		return 0, fmt.Errorf("account compare-dispatch B input: %w", err)
	}
	s.x[8] = uint64(uint32(value))
	s.x[9] = uint64(uint32(s.x[10]) ^ 0xce7525df)
	s.compare32(uint32(s.x[8]), uint32(s.x[9]))
	s.x[8] = uint64(uint32(boolValue(s.condition("eq"))))
	s.x[9] = uint64(uint32(s.x[10]) + 0x3adf)
	s.x[8] = uint64(uint32(s.x[8]) * uint32(s.x[9]))
	s.x[8] = uint64(uint32(s.x[8]) ^ uint32(s.x[10]))
	next, err := loadBootstrapPointer(s, s.x[27], uint32(s.x[8]))
	if err != nil {
		return 0, fmt.Errorf("account compare-dispatch B target: %w", err)
	}
	s.x[8] = next
	s.x[17] = 0x2ebf
	return next, nil
}
