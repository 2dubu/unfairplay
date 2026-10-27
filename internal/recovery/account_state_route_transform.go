package recovery

import "fmt"

const (
	accountStateCompareRoute = uint64(0x100943d90)
	accountStateGuardRoute   = uint64(0x1009972b8)
)

func runAccountStateCompareRoute(s *state) (uint64, error) {
	s.traceBlock(accountStateCompareRoute)
	s.x[10] = uint64(uint32(s.x[10]) - 0x7d5)
	s.x[8] = 0x2ed5
	s.x[11] = uint64(uint32(s.x[10]) ^ uint32(s.x[8]))
	s.x[8] = 0xffff9ccb
	s.x[8] = uint64(uint32(s.x[10]) + uint32(s.x[8]))
	s.x[12] = 0x51e0
	s.x[8] = uint64(uint32(s.x[8]) | uint32(s.x[12]))
	s.x[12] = uint64(uint32(s.x[23]) ^ 0xa7089364)
	s.x[13] = uint64(uint32(s.x[23]) ^ 0x58f76c9c)
	s.x[12] = uint64(uint32(s.x[12]) + uint32(s.x[13]))
	s.x[13] = uint64(uint32(s.x[23]) ^ 0x575a0eff)
	s.x[14] = 0xa8a5f101
	s.x[13] = uint64(uint32(s.x[13]) + uint32(s.x[14]))
	s.x[12] = uint64(uint32(s.x[12]) ^ uint32(s.x[13]))
	s.x[13] = uint64(uint32(s.x[11]) + 0x3e27867e)
	s.x[13] = uint64(uint32(s.x[13]) ^ uint32(s.x[23]))
	s.x[14] = 0xc1d819a2
	s.x[13] = uint64(uint32(s.x[13]) + uint32(s.x[14]))
	s.x[12] = uint64(uint32(s.x[12]) ^ uint32(s.x[13]))
	s.x[13] = 0xffffa01c
	s.x[13] = uint64(uint32(s.x[8]) ^ uint32(s.x[13]))
	s.x[12] = uint64(uint32(s.x[12]) & uint32(s.x[13]))
	s.x[13] = 0xce757bc8
	s.x[12] = uint64(uint32(s.x[12]) + uint32(s.x[13]))
	value, err := s.load(s.sp+0x1d8, 4)
	if err != nil {
		return 0, fmt.Errorf("account state-compare input: %w", err)
	}
	s.x[13] = uint64(uint32(value))
	s.x[12] = uint64(uint32(s.x[12]) ^ uint32(s.x[13]))
	s.compare32(uint32(s.x[12]), uint32(s.x[3]))
	s.x[12] = uint64(uint32(boolValue(s.condition("ne"))))
	s.x[13] = 0x74
	s.x[12] = uint64(uint32(s.x[12]) * uint32(s.x[13]))
	s.x[10] = uint64(uint32(s.x[12]) ^ uint32(s.x[10]))
	next, err := loadBootstrapPointer(s, s.x[30], uint32(s.x[10]))
	if err != nil {
		return 0, fmt.Errorf("account state-compare target: %w", err)
	}
	s.x[10] = next
	s.x[17] = 0x2ee6
	return next, nil
}

func runAccountStateGuardRoute(s *state) (uint64, error) {
	s.traceBlock(accountStateGuardRoute)
	s.x[10] = 0xd59bc80c
	s.x[22] = uint64(uint32(s.x[8]) + uint32(s.x[10]))
	s.x[10] = uint64(uint32(s.x[22]) ^ 0xd13be8ef)
	value, err := s.load(s.sp+0x3f4, 4)
	if err != nil {
		return 0, fmt.Errorf("account state-guard input: %w", err)
	}
	s.x[11] = uint64(uint32(value))
	s.x[12] = uint64(uint32(s.x[10]) + 0x7e5b9b83)
	s.x[12] = uint64(uint32(s.x[12]) ^ uint32(s.x[23]))
	s.x[13] = uint64(uint32(s.x[23]) ^ 0x696550fa)
	s.x[14] = 0x969aaf06
	s.x[13] = uint64(uint32(s.x[13]) + uint32(s.x[14]))
	s.x[14] = uint64(uint32(s.x[23]) ^ 0x241341b9)
	s.x[15] = 0xdbecbe47
	s.x[14] = uint64(uint32(s.x[14]) + uint32(s.x[15]))
	s.x[13] = uint64(uint32(s.x[14]) ^ uint32(s.x[13]))
	s.x[14] = 0x7cfc957a
	s.x[12] = uint64(uint32(s.x[12]) + uint32(s.x[14]))
	s.x[12] = uint64(uint32(s.x[13]) ^ uint32(s.x[12]))
	s.x[11] = uint64(uint32(s.x[11]) - uint32(s.x[12]))
	s.x[12] = 0x785de4ec
	s.compare32(uint32(s.x[9]), uint32(s.x[12]))
	s.x[12] = 0xf529f7ad
	if s.condition("ne") {
		s.compare32(uint32(s.x[11]), uint32(s.x[12]))
	} else {
		s.setNZCV(0)
	}
	s.x[11] = uint64(uint32(boolValue(s.condition("ne"))))
	s.x[12] = 0x3fef
	s.x[11] = uint64(uint32(s.x[11]) * uint32(s.x[12]))
	s.x[8] = uint64(uint32(s.x[11]) ^ uint32(s.x[8]))
	next, err := loadBootstrapPointer(s, s.x[30], uint32(s.x[8]))
	if err != nil {
		return 0, fmt.Errorf("account state-guard target: %w", err)
	}
	s.x[8] = next
	s.x[5] = 0x318a201b
	s.x[17] = 0xfd9d
	return next, nil
}
