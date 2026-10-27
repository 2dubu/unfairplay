package recovery

import "fmt"

const (
	accountByteTargetRoute     = uint64(0x1002723e8)
	accountByteTargetRouteTail = uint64(0x100272400)
	accountByteFourRoute       = uint64(0x100276c4c)
	accountComputedTargetRoute = uint64(0x10027a6a4)
)

func runAccountByteTargetRoute(s *state) (uint64, error) {
	s.traceBlock(accountByteTargetRoute)
	s.x[8] = 0x7f011a25
	s.x[12] = uint64(uint32(s.x[14]) + uint32(s.x[8]))
	value, err := s.load(s.x[13]+2, 1)
	if err != nil {
		return 0, fmt.Errorf("account byte-target input: %w", err)
	}
	s.x[8] = uint64(byte(value))
	next, err := loadBootstrapPointer(s, s.x[10], uint32(s.x[12]))
	if err != nil {
		return 0, fmt.Errorf("account byte-target target: %w", err)
	}
	s.x[12] = next
	s.x[17] = 0xd68d
	s.traceBlock(accountByteTargetRouteTail)
	return next, nil
}

func runAccountByteFourRoute(s *state) (uint64, error) {
	s.traceBlock(accountByteFourRoute)
	s.x[8] = 0xbcffc92c
	s.x[14] = 0x81ff6fef
	s.x[8] = uint64(uint32(s.x[12]) + uint32(s.x[8]))
	s.x[14] = uint64(uint32(s.x[8]) & uint32(s.x[14]))
	value, err := s.load(s.x[13]+1, 1)
	if err != nil {
		return 0, fmt.Errorf("account byte-four input: %w", err)
	}
	s.x[8] = uint64(byte(value))
	s.compare32(uint32(s.x[8]), 4)
	s.x[8] = uint64(uint32(boolValue(s.condition("eq"))))
	s.x[15] = uint64(uint32(s.x[12]) + 0x7f006c74)
	s.x[15] = uint64(uint32(s.x[15]) + uint32(s.x[14]))
	s.x[8] = uint64(uint32(s.x[8]) * uint32(s.x[15]))
	s.x[8] = uint64(uint32(s.x[8]) ^ uint32(s.x[12]))
	next, err := loadBootstrapPointer(s, s.x[10], uint32(s.x[8]))
	if err != nil {
		return 0, fmt.Errorf("account byte-four target: %w", err)
	}
	s.x[12] = next
	s.x[8] = 0xffff5bde
	s.x[17] = 0x9478
	return next, nil
}

func runAccountComputedTargetRoute(s *state) (uint64, error) {
	s.traceBlock(accountComputedTargetRoute)
	s.x[13] = uint64(uint32(s.x[12]) + (uint32(0xa3d) << 12))
	s.x[13] = uint64(uint32(s.x[13]) + 0x2f6)
	s.x[14] = 0x805b6626
	s.x[13] = uint64(uint32(s.x[13]) | uint32(s.x[14]))
	s.x[14] = 0x4d761558
	s.x[13] = uint64(uint32(s.x[13]) + uint32(s.x[14]))
	s.compare32(uint32(s.x[8]), uint32(s.x[13]))
	s.x[8] = uint64(uint32(boolValue(s.condition("eq"))))
	s.x[13] = 0xa083
	s.x[8] = uint64(uint32(s.x[8]) * uint32(s.x[13]))
	s.x[8] = uint64(uint32(s.x[8]) ^ uint32(s.x[12]))
	value, err := s.load(s.x[10]+uint64(uint32(s.x[8]))*8, 8)
	if err != nil {
		return 0, fmt.Errorf("account computed-target target: %w", err)
	}
	s.x[13] = value
	s.x[8] = 0xffff5bde
	s.x[17] = 0x9478
	return value, nil
}
