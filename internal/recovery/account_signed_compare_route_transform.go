package recovery

import "fmt"

const (
	accountSignedCompareRoute     = uint64(0x100a19fc0)
	accountSignedCompareRouteTail = uint64(0x100a1a000)
)

func runAccountSignedCompareRoute(s *state) (uint64, error) {
	s.traceBlock(accountSignedCompareRoute)
	s.x[11] = 0x3d756422
	s.x[11] = uint64(uint32(s.x[8]) + uint32(s.x[11]))
	s.x[12] = 0xc28acfbe
	s.x[13] = uint64(uint32(s.x[11]) & uint32(s.x[12]))
	s.x[11] = 0xffffcbc5
	s.x[11] = uint64(uint32(s.x[8]) + uint32(s.x[11]))
	s.x[12] = 0x815
	s.x[12] = uint64(uint32(s.x[11]) | uint32(s.x[12]))
	value, err := s.load(s.x[10], 1)
	if err != nil {
		return 0, fmt.Errorf("account signed-compare first input: %w", err)
	}
	s.x[11] = uint64(uint32(int32(int8(value))))
	value, err = s.load(s.x[9], 1)
	if err != nil {
		return 0, fmt.Errorf("account signed-compare second input: %w", err)
	}
	s.x[14] = uint64(uint32(int32(int8(value))))
	s.x[15] = uint64(uint32(s.x[14]) ^ uint32(s.x[11]))
	s.x[15] &= 0xff
	s.x[16] = uint64(uint32(s.x[12]) - 0x22)
	s.compare32(uint32(s.x[15]), uint32(byte(s.x[16])))

	s.traceBlock(accountSignedCompareRouteTail)
	s.x[15] = uint64(uint32(boolValue(s.condition("eq"))))
	s.x[16] = 0xffff72a6
	s.x[16] = uint64(uint32(s.x[13]) + uint32(s.x[16]))
	s.x[15] = uint64(uint32(s.x[15]) * uint32(s.x[16]))
	s.x[16] = 0xffffd0a8
	s.x[8] = uint64(uint32(s.x[8]) + uint32(s.x[16]))
	s.x[8] = uint64(uint32(s.x[15]) ^ uint32(s.x[8]))
	next, err := loadBootstrapPointer(s, s.x[30], uint32(s.x[8]))
	if err != nil {
		return 0, fmt.Errorf("account signed-compare target: %w", err)
	}
	s.x[8] = next
	s.x[17] = 0x3b8d
	return next, nil
}
