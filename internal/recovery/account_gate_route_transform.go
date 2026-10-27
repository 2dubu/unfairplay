package recovery

import "fmt"

const (
	accountGateSetupRoute = uint64(0x10001b158)
	accountMaskGateRoute  = uint64(0x10006b958)
	accountDualMaskRoute  = uint64(0x1000206ec)
	accountStackLoadRoute = uint64(0x1000b46bc)
)

func runAccountGateSetupRoute(s *state) (uint64, error) {
	s.traceBlock(accountGateSetupRoute)
	value, err := s.load(s.sp+0x18, 4)
	if err != nil {
		return 0, fmt.Errorf("account gate-setup first word: %w", err)
	}
	s.x[0] = uint64(uint32(value))
	s.x[16] = 0xab68f2e5
	s.compare32(uint32(s.x[0]), uint32(s.x[16]))
	s.x[11] = uint64(uint32(boolValue(s.condition("lo"))))
	s.x[13] = 0x147fd7d7
	s.x[15] = uint64(uint32(s.x[12]) ^ uint32(s.x[13]))
	s.x[13] = 0xab6866a4
	s.x[13] = uint64(uint32(s.x[9]) + uint32(s.x[13]))
	s.x[13] = uint64(uint32(s.x[15]) + uint32(s.x[13]))
	s.compare32(uint32(s.x[13]), uint32(s.x[16]))
	s.x[16] = uint64(uint32(boolValue(s.condition("lo"))))
	s.x[11] = uint64(uint32(s.x[11]) ^ uint32(s.x[16]))
	s.compare32(uint32(s.x[13]), uint32(s.x[0]))
	s.x[13] = uint64(uint32(boolValue(s.condition("hi"))))
	s.compare32(uint32(s.x[11]), 0)
	if s.condition("ne") {
		s.x[11] = uint64(uint32(s.x[16]))
	} else {
		s.x[11] = uint64(uint32(s.x[13]))
	}
	s.x[13] = 0x32
	s.x[11] = uint64(uint32(s.x[11]) * uint32(s.x[13]))
	s.x[13] = uint64(uint32(s.x[12]) + 0xeb81312b)
	s.x[11] = uint64(uint32(s.x[11]) ^ uint32(s.x[13]))
	next, err := loadBootstrapPointer(s, s.x[21], uint32(s.x[11]))
	if err != nil {
		return 0, fmt.Errorf("account gate-setup target: %w", err)
	}
	s.x[13] = next
	s.x[11] = 0x446b6361
	s.x[16] = 0x7d2c7c67
	value, err = s.load(s.sp, 8)
	if err != nil {
		return 0, fmt.Errorf("account gate-setup base: %w", err)
	}
	s.x[2] = value
	value, err = s.load(s.sp+0x20, 8)
	if err != nil {
		return 0, fmt.Errorf("account gate-setup pointer: %w", err)
	}
	s.x[24] = value
	value, err = s.load(s.sp+0x2c, 4)
	if err != nil {
		return 0, fmt.Errorf("account gate-setup stack word: %w", err)
	}
	s.x[25] = uint64(uint32(value))
	s.x[3] = 0xa9
	s.x[7] = 0xc6
	s.x[19] = 0xffff73a7
	s.x[20] = 0x5580
	s.x[22] = 0xffffffff82d32df1
	s.x[23] = 0x36
	s.x[26] = 0x7d2c7ca7
	s.x[17] = 0x71c7
	return next, nil
}

func runAccountMaskGateRoute(s *state) (uint64, error) {
	s.traceBlock(accountMaskGateRoute)
	s.x[9] = 0x69a3622f
	s.x[9] = uint64(uint32(s.x[12]) + uint32(s.x[9]))
	s.x[11] = 0x965cd7f9
	s.x[9] = uint64(uint32(s.x[9]) & uint32(s.x[11]))
	s.x[13] = uint64(uint32(s.x[12]) + 0x3a29)
	s.x[11] = uint64(uint32(s.x[13]) ^ 0x12d6a942)
	s.x[11] = uint64(uint32(s.x[11]) ^ uint32(s.x[9]))
	s.compare32(uint32(s.x[4]), uint32(s.x[11]))
	s.x[11] = uint64(uint32(boolValue(s.condition("ge"))))
	s.x[15] = 0xf7
	s.x[11] = uint64(uint32(s.x[11]) * uint32(s.x[15]))
	s.x[11] = uint64(uint32(s.x[11]) ^ uint32(s.x[12]))
	next, err := loadBootstrapPointer(s, s.x[21], uint32(s.x[11]))
	if err != nil {
		return 0, fmt.Errorf("account mask-gate target: %w", err)
	}
	s.x[11] = next
	s.x[17] = 0x2d4d
	return next, nil
}

func runAccountDualMaskRoute(s *state) (uint64, error) {
	s.traceBlock(accountDualMaskRoute)
	s.x[9] = 0x002b056a
	s.x[9] = uint64(uint32(s.x[13]) + uint32(s.x[9]))
	s.x[11] = 0x14540044
	s.x[12] = uint64(uint32(s.x[9]) | uint32(s.x[11]))
	s.x[9] = 0xd67f8a41
	s.x[9] = uint64(uint32(s.x[13]) + uint32(s.x[9]))
	s.x[11] = 0xfbfff29f
	s.x[11] = uint64(uint32(s.x[9]) & uint32(s.x[11]))
	s.x[9] = 0x4056d6b9
	s.x[9] = uint64(uint32(s.x[11]) + uint32(s.x[9]))
	s.compare32(uint32(s.x[4]), uint32(s.x[9]))
	s.x[9] = uint64(uint32(boolValue(s.condition("eq"))))
	s.x[15] = uint64(uint32(s.x[12]) + 0xeb816c6a)
	s.x[9] = uint64(uint32(s.x[9]) * uint32(s.x[15]))
	s.x[9] = uint64(uint32(s.x[9]) ^ uint32(s.x[13]))
	next, err := loadBootstrapPointer(s, s.x[21], uint32(s.x[9]))
	if err != nil {
		return 0, fmt.Errorf("account dual-mask target: %w", err)
	}
	s.x[9] = next
	s.x[17] = 0xb54d
	return next, nil
}

func runAccountStackLoadRoute(s *state) (uint64, error) {
	s.traceBlock(accountStackLoadRoute)
	s.x[9] = uint64(uint32(s.x[12]) + 0xeb80f63d)
	next, err := loadBootstrapPointer(s, s.x[21], uint32(s.x[9]))
	if err != nil {
		return 0, fmt.Errorf("account stack-load target: %w", err)
	}
	s.x[11] = next
	value, err := s.load(s.sp+0x2c, 4)
	if err != nil {
		return 0, fmt.Errorf("account stack-load word: %w", err)
	}
	s.x[9] = uint64(uint32(value))
	s.x[17] = 0x7894
	return next, nil
}
