package recovery

import "fmt"

const (
	accountLessThanRoute  = uint64(0x100900b5c)
	accountOffsetRoute    = uint64(0x1009851c8)
	accountSevenFoldRoute = uint64(0x100a7b90c)
	accountWideLimitRoute = uint64(0x10067049c)
)

func runAccountLessThanRoute(s *state) (uint64, error) {
	s.traceBlock(accountLessThanRoute)
	s.x[11] = 0x61c8
	s.x[14] = uint64(uint32(s.x[12]) + uint32(s.x[11]))
	s.x[11] = uint64(uint32(s.x[14]) ^ 0x296d)
	s.x[13] = 0x04e043ad
	s.compare32(uint32(s.x[9]), uint32(s.x[13]))
	s.x[13] = uint64(uint32(boolValue(s.condition("lt"))))
	s.x[13] = uint64(uint32(s.x[11]) ^ uint32(s.x[13]))
	s.x[13] = uint64(uint32(s.x[0]) &^ uint32(s.x[13]))
	s.x[15] = 0x57
	s.x[13] = uint64(uint32(s.x[13]) * uint32(s.x[15]))
	s.x[12] = uint64(uint32(s.x[12]) + 0x2751)
	s.x[12] = uint64(uint32(s.x[13]) ^ uint32(s.x[12]))
	next, err := loadBootstrapPointer(s, s.x[30], uint32(s.x[12]))
	if err != nil {
		return 0, fmt.Errorf("account less-than target: %w", err)
	}
	s.x[12] = next
	s.x[17] = 0xdb0
	return next, nil
}

func runAccountOffsetRoute(s *state) (uint64, error) {
	s.traceBlock(accountOffsetRoute)
	s.x[11] = uint64(uint32(s.x[12]) + 0x15aa)
	next, err := loadBootstrapPointer(s, s.x[30], uint32(s.x[11]))
	if err != nil {
		return 0, fmt.Errorf("account offset target: %w", err)
	}
	s.x[11] = next
	s.x[17] = 0x21a3
	return next, nil
}

func runAccountSevenFoldRoute(s *state) (uint64, error) {
	s.traceBlock(accountSevenFoldRoute)
	s.x[11] = 0x3fa
	s.x[14] = uint64(uint32(s.x[13]) ^ uint32(s.x[11]))
	s.x[11] = uint64(uint32(s.x[14]) + 0x692d)
	s.x[12] = uint64(uint32(s.x[14]) << 3)
	s.x[12] = uint64(uint32(s.x[12]) - uint32(s.x[14]))
	s.x[15] = uint64(uint32(s.x[12]) + 0x04dfb657)
	s.compare32(uint32(s.x[9]), uint32(s.x[15]))
	s.x[15] = uint64(uint32(boolValue(s.condition("eq"))))
	s.x[16] = 0xffffebe0
	s.x[14] = uint64(uint32(s.x[14]) + uint32(s.x[16]))
	s.x[14] = uint64(uint32(s.x[15]) * uint32(s.x[14]))
	s.x[13] = uint64(uint32(s.x[14]) ^ uint32(s.x[13]))
	next, err := loadBootstrapPointer(s, s.x[30], uint32(s.x[13]))
	if err != nil {
		return 0, fmt.Errorf("account seven-fold target: %w", err)
	}
	s.x[13] = next
	s.x[17] = 0x7b0c
	return next, nil
}

func runAccountWideLimitRoute(s *state) (uint64, error) {
	s.traceBlock(accountWideLimitRoute)
	s.x[8] = uint64(uint32(s.x[10]) + 0xc84)
	s.x[10] = uint64(uint32(s.x[8]) + uint32(s.x[13]))
	s.x[2] = uint64(uint32(s.x[10]) | uint32(s.x[15]))
	s.x[10] = uint64(uint32(s.x[8]) ^ uint32(s.x[16]))
	s.x[9] += s.x[0]
	s.x[9] += s.x[2]
	s.compare64(s.x[9], 0x3ff)
	s.x[17] = uint64(uint32(boolValue(s.condition("eq"))))
	s.x[3] = uint64(uint32(s.x[10]) + uint32(s.x[1]))
	s.x[17] = uint64(uint32(s.x[17]) * uint32(s.x[3]))
	s.x[8] = uint64(uint32(s.x[17]) ^ uint32(s.x[8]))
	value, err := s.load(s.sp+0xdf0, 8)
	if err != nil {
		return 0, fmt.Errorf("account wide-limit table: %w", err)
	}
	s.x[17] = value
	next, err := loadBootstrapPointer(s, s.x[17], uint32(s.x[8]))
	if err != nil {
		return 0, fmt.Errorf("account wide-limit target: %w", err)
	}
	s.x[8] = next
	s.x[17] = 0xe9f8
	return next, nil
}
