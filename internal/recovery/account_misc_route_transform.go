package recovery

import "fmt"

const (
	accountByteIndexRoute  = uint64(0x10024da1c)
	accountWideOffsetRoute = uint64(0x1002734dc)
	accountRouteConstant   = uint64(0x1002b4f80)
	accountCompareRoute    = uint64(0x10075a804)
	accountExpressionRoute = uint64(0x100c1ee1c)
	accountBitTestRoute    = uint64(0x100c94234)
)

func runAccountByteIndexRoute(s *state) (uint64, error) {
	s.traceBlock(accountByteIndexRoute)
	s.x[14] = uint64(uint32(s.x[13]) ^ uint32(s.x[2]))
	value, err := loadBootstrapPointer(s, s.x[27], uint32(s.x[8]))
	if err != nil {
		return 0, fmt.Errorf("account byte-index first target: %w", err)
	}
	s.x[8] = value
	value, err = s.load(s.x[11], 8)
	if err != nil {
		return 0, fmt.Errorf("account byte-index first input pointer: %w", err)
	}
	s.x[12] = value
	value, err = s.load(s.x[12], 1)
	if err != nil {
		return 0, fmt.Errorf("account byte-index first input: %w", err)
	}
	s.x[11] = uint64(byte(value))
	value, err = s.load(s.x[9], 1)
	if err != nil {
		return 0, fmt.Errorf("account byte-index second input: %w", err)
	}
	s.x[15] = uint64(byte(value))
	s.x[15] = uint64(uint32(s.x[15]) ^ uint32(s.x[24]))
	s.x[11] = uint64(uint32(s.x[11]) - uint32(s.x[15]))
	s.x[15] = uint64(uint32(s.x[11]) ^ uint32(s.x[25]))
	s.x[16] = uint64(uint32(s.x[14]) + 0x1b5222b6)
	s.x[17] = 0x1b527404
	s.x[16] = uint64(uint32(s.x[16]) ^ uint32(s.x[17]))
	s.x[11] = uint64(uint32(shiftLeft(s.x[11], int(s.x[16]), 32)))
	s.x[11] = uint64(uint32(s.x[11]) & uint32(s.x[23]))
	s.x[11] = uint64(uint32(s.x[11]) + uint32(s.x[15]))
	s.compare32(uint32(s.x[11]), uint32(s.x[25]))
	s.x[11] = uint64(uint32(boolValue(s.condition("eq"))))
	s.x[15] = 0xd
	s.x[11] = uint64(uint32(s.x[11]) * uint32(s.x[15]))
	s.x[11] = uint64(uint32(s.x[11]) ^ uint32(s.x[13]))
	next, err := loadBootstrapPointer(s, s.x[27], uint32(s.x[11]))
	if err != nil {
		return 0, fmt.Errorf("account byte-index target: %w", err)
	}
	s.x[13] = next
	s.x[11] = 1
	s.x[17] = 0x34ca
	return next, nil
}

func runAccountWideOffsetRoute(s *state) (uint64, error) {
	s.traceBlock(accountWideOffsetRoute)
	s.x[28] = s.x[0]
	s.x[16] = 0
	s.x[11] = 0x702e
	s.x[13] = uint64(uint32(s.x[14]) ^ uint32(s.x[11]))
	s.x[11] = 0xffffffff99c2dace
	s.x[11] = s.x[10] + s.x[11]
	s.x[11] += s.x[13]
	next, err := loadBootstrapPointer(s, s.x[27], uint32(s.x[14]))
	if err != nil {
		return 0, fmt.Errorf("account wide-offset target: %w", err)
	}
	s.x[14] = next
	s.x[17] = 0x30c
	return next, nil
}

func runAccountRouteConstant(s *state) (uint64, error) {
	s.traceBlock(accountRouteConstant)
	s.x[28] = 0x1b527405
	s.x[17] = 0x6f0a
	return s.x[8], nil
}

func runAccountCompareRoute(s *state) (uint64, error) {
	s.traceBlock(accountCompareRoute)
	s.x[14] = uint64(uint32(s.x[12]) + uint32(s.x[30]))
	s.x[11] = uint64(uint32(s.x[14]) ^ uint32(s.x[3]))
	s.x[10] = uint64(uint32(s.x[10]) + 1)
	s.x[15] = s.x[11] ^ s.x[4]
	s.x[8] = s.x[15] + s.x[8]
	s.compare32(uint32(s.x[13]), uint32(s.x[10]))
	s.x[15] = uint64(uint32(boolValue(s.condition("eq"))))
	s.x[15] = uint64(uint32(s.x[15]) * uint32(s.x[5]))
	s.x[12] = uint64(uint32(s.x[12]) ^ uint32(s.x[15]))
	value, err := s.load(s.sp+0xdf0, 8)
	if err != nil {
		return 0, fmt.Errorf("account compare-route table: %w", err)
	}
	s.x[15] = value
	next, err := loadBootstrapPointer(s, s.x[15], uint32(s.x[12]))
	if err != nil {
		return 0, fmt.Errorf("account compare-route target: %w", err)
	}
	s.x[12] = next
	s.x[17] = 0x529e
	return next, nil
}

func runAccountExpressionRoute(s *state) (uint64, error) {
	s.traceBlock(accountExpressionRoute)
	s.x[9] = uint64(uint32(s.x[3]) ^ uint32(s.x[30]))
	s.x[25] = uint64(uint32(s.x[9]) * uint32(s.x[8]))
	s.x[9] = uint64(uint32(s.x[3]) + uint32(s.x[11]))
	s.x[9] = uint64(uint32(s.x[9]) | uint32(s.x[15]))
	s.x[17] = uint64(uint32(s.x[2]) + uint32(s.x[16]))
	s.x[2] = uint64(uint32(s.x[17]) + uint32(s.x[9]))
	s.x[17] = uint64(uint32(s.x[25]) + uint32(s.x[0]))
	s.compare32(uint32(s.x[2]), uint32(s.x[17]))
	s.x[17] = uint64(uint32(boolValue(s.condition("eq"))))
	s.x[17] = uint64(uint32(s.x[17]) * uint32(s.x[1]))
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[3]))
	next, err := loadBootstrapPointer(s, s.x[28], uint32(s.x[17]))
	if err != nil {
		return 0, fmt.Errorf("account expression-route target: %w", err)
	}
	s.x[3] = next
	s.x[26] = s.x[4]
	s.x[17] = 0x119a
	return next, nil
}

func runAccountBitTestRoute(s *state) (uint64, error) {
	s.traceBlock(accountBitTestRoute)
	s.x[17] = uint64(uint32(s.x[9]) ^ uint32(s.x[12]))
	s.x[17] = uint64(uint32(s.x[17]) + uint32(s.x[17]<<3))
	s.x[3] = uint64(uint32(s.x[17]) << 1)
	s.x[25] = uint64(uint32(s.x[26]) >> 1)
	s.x[4] = uint64(uint32(s.x[25]) ^ uint32(s.x[13]))
	s.test32(uint32(s.x[26]) & 2)
	s.x[17] = uint64(uint32(boolValue(s.condition("eq"))))
	s.x[26] = uint64(uint32(s.x[3]) ^ uint32(s.x[14]))
	s.x[17] = uint64(uint32(s.x[17]) * uint32(s.x[26]))
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[9]))
	next, err := loadBootstrapPointer(s, s.x[28], uint32(s.x[17]))
	if err != nil {
		return 0, fmt.Errorf("account bit-test target: %w", err)
	}
	s.x[26] = next
	s.x[17] = 0xc750
	return next, nil
}
