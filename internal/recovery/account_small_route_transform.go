package recovery

import "fmt"

const (
	accountRecordGuardRoute = uint64(0x100952a40)
	accountStackGateSetup   = uint64(0x1001034e0)
	accountWideCompareRoute = uint64(0x1001a147c)
	accountMaskTargetRoute  = uint64(0x100246ee0)
)

func runAccountRecordGuardRoute(s *state) (uint64, error) {
	s.traceBlock(accountRecordGuardRoute)
	s.x[9] = 0x5af8e3dc
	s.x[9] = uint64(uint32(s.x[10]) + uint32(s.x[9]))
	s.x[11] = 0xa506dfdd
	s.x[9] = uint64(uint32(s.x[9]) & uint32(s.x[11]))
	s.x[11] = uint64(uint32(s.x[25]) + 0xfb5830fc)
	s.x[11] = uint64(int64(int32(s.x[11])) * int64(int32(s.x[0])))
	s.x[12] = s.x[19] + s.x[11]
	value, err := s.load(s.x[12]+0x24, 4)
	if err != nil {
		return 0, fmt.Errorf("account record-guard input: %w", err)
	}
	s.x[11] = uint64(uint32(value))
	s.compare32(uint32(s.x[11]), uint32(s.x[25]))
	s.x[11] = uint64(uint32(boolValue(s.condition("eq"))))
	s.x[12] = 0xffffa642
	s.x[12] = uint64(uint32(s.x[9]) + uint32(s.x[12]))
	s.x[11] = uint64(uint32(s.x[12]) * uint32(s.x[11]))
	s.x[10] = uint64(uint32(s.x[11]) ^ uint32(s.x[10]))
	next, err := loadBootstrapPointer(s, s.x[30], uint32(s.x[10]))
	if err != nil {
		return 0, fmt.Errorf("account record-guard target: %w", err)
	}
	s.x[11] = next
	s.x[10] = 0x04a7cf03
	s.x[17] = 0x5313
	return next, nil
}

func runAccountStackGateSetup(s *state) (uint64, error) {
	s.traceBlock(accountStackGateSetup)
	value, err := s.load(s.sp+0x38, 4)
	if err != nil {
		return 0, fmt.Errorf("account stack-gate word: %w", err)
	}
	s.x[9] = uint64(uint32(value))
	s.x[12] = 0x3d4d8efa
	s.compare32(uint32(s.x[9]), uint32(s.x[12]))
	s.x[8] = uint64(uint32(boolValue(s.condition("lo"))))
	s.x[21] = 0xb904
	s.x[13] = 0x3d4d8f3a
	s.x[13] = uint64(uint32(s.x[11]) + uint32(s.x[13]))
	s.compare32(uint32(s.x[13]), uint32(s.x[12]))
	s.x[15] = uint64(uint32(boolValue(s.condition("lo"))))
	s.x[8] = uint64(uint32(s.x[8]) ^ uint32(s.x[15]))
	s.compare32(uint32(s.x[13]), uint32(s.x[9]))
	s.x[13] = uint64(uint32(boolValue(s.condition("hi"))))
	s.compare32(uint32(s.x[8]), 0)
	if s.condition("ne") {
		s.x[8] = uint64(uint32(s.x[15]))
	} else {
		s.x[8] = uint64(uint32(s.x[13]))
	}
	s.x[8] = uint64(uint32(s.x[8]) ^ 1)
	s.x[9] = 0x7d
	s.x[8] = uint64(uint32(s.x[8]) * uint32(s.x[9]))
	s.x[8] = uint64(uint32(s.x[8]) ^ 0x7922)
	value, err = s.load(s.x[20]+uint64(uint32(s.x[8]))*8, 8)
	if err != nil {
		return 0, fmt.Errorf("account stack-gate target: %w", err)
	}
	s.x[15] = value
	s.x[8] = 0x7ca22ead
	s.x[13] = 0x0edffc52
	s.x[9] = 0x795f
	value, err = s.load(s.sp+0x20, 8)
	if err != nil {
		return 0, fmt.Errorf("account stack-gate saved pointer: %w", err)
	}
	s.x[12] = value
	s.x[0] = 0x5691bcf2
	s.x[1] = 0xfffffffff12003ae
	s.x[2] = 0x0edffc92
	s.x[5] = 0xdb51
	s.x[17] = 0xd90f
	return s.x[15], nil
}

func runAccountWideCompareRoute(s *state) (uint64, error) {
	s.traceBlock(accountWideCompareRoute)
	s.x[9] = 0x1b52ea91
	s.x[9] = uint64(uint32(s.x[15]) ^ uint32(s.x[9]))
	s.x[11] = 0x7d6f77d4
	s.x[9] ^= s.x[11]
	s.x[9] += s.x[14]
	s.compare64(s.x[9], s.x[10])
	s.x[11] = uint64(uint32(boolValue(s.condition("lo"))))
	next, err := loadBootstrapPointer(s, s.x[27], uint32(s.x[15]))
	if err != nil {
		return 0, fmt.Errorf("account wide-compare target: %w", err)
	}
	s.x[9] = next
	s.x[0] = s.x[28]
	s.x[2] = 0xcfdb
	s.x[17] = 0x34ca
	return next, nil
}

func runAccountMaskTargetRoute(s *state) (uint64, error) {
	s.traceBlock(accountMaskTargetRoute)
	s.x[13] = 0xfca75ab2
	s.x[13] = uint64(uint32(s.x[12]) + uint32(s.x[13]))
	s.x[14] = 0x03589ebf
	s.x[13] = uint64(uint32(s.x[13]) & uint32(s.x[14]))
	next, err := loadBootstrapPointer(s, s.x[27], uint32(s.x[12]))
	if err != nil {
		return 0, fmt.Errorf("account mask-target target: %w", err)
	}
	s.x[12] = next
	s.x[17] = 0x32b3
	return next, nil
}
