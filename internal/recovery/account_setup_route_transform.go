package recovery

import "fmt"

const (
	accountSetupRoute  = uint64(0x1005806f8)
	accountLengthRoute = uint64(0x1008f6390)
)

func runAccountSetupRoute(s *state) (uint64, error) {
	s.traceBlock(accountSetupRoute)
	s.x[9] = uint64(uint32(s.x[8]) ^ 0xe42)
	s.x[9] = uint64(uint32(s.x[9]) + uint32(s.x[9]<<2))
	value, err := s.load(s.sp+0x3ec, 4)
	if err != nil {
		return 0, fmt.Errorf("account setup-route word: %w", err)
	}
	s.x[23] = uint64(uint32(value))
	value, err = s.load(s.sp+0x2c0, 8)
	if err != nil {
		return 0, fmt.Errorf("account setup-route pointer: %w", err)
	}
	s.x[6] = value
	s.x[10] = 0x26b47be8
	if err := s.storeValue(s.sp+0x3f4, uint64(uint32(s.x[10])), 4); err != nil {
		return 0, fmt.Errorf("account setup-route output: %w", err)
	}
	next, err := loadBootstrapPointer(s, s.x[30], uint32(s.x[8]))
	if err != nil {
		return 0, fmt.Errorf("account setup-route target: %w", err)
	}
	s.x[11] = next
	s.x[8] = 0xce757bc5
	s.x[10] = s.x[25]
	s.x[17] = 0xa800
	return next, nil
}

func runAccountLengthRoute(s *state) (uint64, error) {
	s.traceBlock(accountLengthRoute)
	s.x[10] = 0xd59c47e3
	s.x[22] = uint64(uint32(s.x[9]) ^ uint32(s.x[10]))
	value, err := s.load(s.x[27], 4)
	if err != nil {
		return 0, fmt.Errorf("account length-route input: %w", err)
	}
	s.x[10] = uint64(uint32(value))
	s.x[11] = 0x87a21b14
	s.x[11] = uint64(uint32(s.x[10]) + uint32(s.x[11]))
	s.compare32(uint32(s.x[11]), 0x10)
	s.x[12] = uint64(uint32(boolValue(s.condition("lo"))))
	s.x[13] = uint64(uint32(s.x[9]) ^ 0xbd15)
	s.x[12] = uint64(uint32(s.x[12]) * uint32(s.x[13]))
	s.x[12] = uint64(uint32(s.x[12]) ^ uint32(s.x[9]))
	next, err := loadBootstrapPointer(s, s.x[30], uint32(s.x[12]))
	if err != nil {
		return 0, fmt.Errorf("account length-route target: %w", err)
	}
	s.x[12] = next
	s.x[5] = 0x318a201b
	s.x[1] = 0x04a7cf03
	s.x[0] = 0x30
	s.x[7] = 0x0a5cd4cf
	s.x[17] = 0xfd9d
	return next, nil
}
