package recovery

import "fmt"

const (
	accountBooleanTargetRoute = uint64(0x1001dfb5c)
	accountZeroFillRoute      = uint64(0x10021e7fc)
	accountZeroFillRouteTail  = uint64(0x10021e800)
	accountDirectReturnSecond = uint64(0x100222afc)
)

func runAccountBooleanTargetRoute(s *state) (uint64, error) {
	s.traceBlock(accountBooleanTargetRoute)
	s.x[11] = 0x5859
	s.x[8] = 0x22f5
	s.x[12] = 0x8d58a2ab
	s.x[14] = uint64(uint32(s.x[16]) ^ uint32(s.x[12]))
	s.x[13] = 0x72a75d55
	s.x[17] = uint64(uint32(s.x[14]) + uint32(s.x[13]))
	s.x[14] = 0x8e6b57ae
	s.x[2] = uint64(uint32(s.x[16]) ^ uint32(s.x[14]))
	s.x[15] = 0x7194a852
	s.x[2] = uint64(uint32(s.x[2]) + uint32(s.x[15]))
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[2]))
	s.x[2] = uint64(uint32(s.x[16]) ^ 0xcd468ec0)
	s.x[16] = 0x32b97140
	s.x[2] = uint64(uint32(s.x[2]) + uint32(s.x[16]))
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[2]))
	s.x[2] = uint64(uint32(s.x[17]) + 0xce757bc5)
	s.x[3] = 0x318a843b
	s.x[17] = uint64(uint32(s.x[3]) - uint32(s.x[17]))
	s.x[17] = uint64(uint32(s.x[2]) | uint32(s.x[17]))
	s.x[17] = uint64(uint32(s.x[17]) >> 31)
	s.x[2] = 0x1e6
	s.x[17] = uint64(uint32(s.x[17]) * uint32(s.x[2]))
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[8]))
	value, err := s.load(s.x[10]+uint64(uint32(s.x[17]))*8, 8)
	if err != nil {
		return 0, fmt.Errorf("account boolean-target target: %w", err)
	}
	s.x[4] = value
	s.x[3] = 0x10af29d9
	s.x[17] = 0x45bc
	return value, nil
}

func runAccountZeroFillRoute(s *state) (uint64, error) {
	s.traceBlock(accountZeroFillRoute)
	s.x[8] = 0
	s.traceBlock(accountZeroFillRouteTail)
	s.x[13] = s.x[0] + 4
	if err := s.storeValue(s.x[13]+uint64(uint32(s.x[11])), uint64(byte(s.x[12])), 1); err != nil {
		return 0, fmt.Errorf("account zero-fill first output: %w", err)
	}
	s.x[11] = uint64(uint32(s.x[14]) + (uint32(0x7ed) << 12))
	s.x[11] = uint64(uint32(s.x[11]) + 0x8f2)
	s.x[12] = 0x8080602a
	s.x[11] = uint64(uint32(s.x[11]) | uint32(s.x[12]))
	s.x[12] = 0x2f
	if err := s.storeValue(s.x[13]+s.x[9], uint64(byte(s.x[12])), 1); err != nil {
		return 0, fmt.Errorf("account zero-fill second output: %w", err)
	}
	s.x[12] = 0x7f009993
	s.x[9] = uint64(uint32(s.x[15]) + uint32(s.x[9]))
	s.x[12] = uint64(uint32(s.x[9]) + uint32(s.x[12]))
	s.x[15] = 0x3d
	if err := s.storeValue(s.x[13]+uint64(uint32(s.x[9])), uint64(byte(s.x[15])), 1); err != nil {
		return 0, fmt.Errorf("account zero-fill third output: %w", err)
	}
	s.x[9] = uint64(uint32(s.x[11]) + uint32(s.x[12]))
	if err := s.storeValue(s.x[13]+uint64(uint32(s.x[9])), 0, 1); err != nil {
		return 0, fmt.Errorf("account zero-fill terminator: %w", err)
	}
	if err := s.storeValue(s.x[0], uint64(uint32(s.x[9])), 4); err != nil {
		return 0, fmt.Errorf("account zero-fill count: %w", err)
	}
	value, err := loadBootstrapPointer(s, s.x[10], uint32(s.x[14]))
	if err != nil {
		return 0, fmt.Errorf("account zero-fill target: %w", err)
	}
	s.x[9] = value
	s.x[17] = 0x9478
	return value, nil
}
