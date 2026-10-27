package recovery

import "fmt"

const accountByteTableRoute = uint64(0x1008688f8)

func runAccountByteTableRoute(s *state) (uint64, error) {
	s.traceBlock(accountByteTableRoute)
	s.x[12] = uint64(uint32(s.x[11]) ^ uint32(s.x[21]))
	s.x[12] = uint64(uint32(s.x[12]) * uint32(s.x[24]))
	value, err := s.load(s.x[9]+uint64(uint32(s.x[10])), 1)
	if err != nil {
		return 0, fmt.Errorf("account byte-table input: %w", err)
	}
	s.x[14] = uint64(byte(value))
	s.x[15] = uint64(uint32(s.x[11]) - 0x19)
	s.x[15] = uint64(uint32(s.x[15]) | uint32(s.x[26]))
	s.x[15] = uint64(uint32(s.x[15]) ^ uint32(s.x[27]))
	s.x[15] = uint64(uint32(s.x[15]) + uint32(s.x[14]))
	if err := s.storeValue(s.x[19]+s.x[8], uint64(byte(s.x[15])), 1); err != nil {
		return 0, fmt.Errorf("account byte-table output: %w", err)
	}
	s.compare32(uint32(s.x[14]), 0x43)
	s.x[14] = uint64(uint32(boolValue(s.condition("eq"))))
	s.x[15] = uint64(uint32(s.x[12]) + uint32(s.x[28]))
	s.x[14] = uint64(uint32(s.x[14]) * uint32(s.x[15]))
	s.x[11] = uint64(uint32(s.x[14]) ^ uint32(s.x[11]))
	value, err = s.load(s.sp+0xdf0, 8)
	if err != nil {
		return 0, fmt.Errorf("account byte-table dispatch: %w", err)
	}
	s.x[14] = value
	next, err := loadBootstrapPointer(s, s.x[14], uint32(s.x[11]))
	if err != nil {
		return 0, fmt.Errorf("account byte-table target: %w", err)
	}
	s.x[11] = next
	s.x[17] = 0xd9db
	return next, nil
}
