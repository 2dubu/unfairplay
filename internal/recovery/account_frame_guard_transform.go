package recovery

import "fmt"

const (
	accountFrameGuardRoute   = uint64(0x1000365d4)
	accountReverseWriteRoute = uint64(0x1000b0690)
)

func runAccountFrameGuardRoute(s *state) (uint64, error) {
	s.traceBlock(accountFrameGuardRoute)
	s.x[9] = 0x1447f9dd
	s.x[23] = uint64(uint32(s.x[8]) + uint32(s.x[9]))
	s.x[9] = uint64(uint32(s.x[8]) + 0x71f4)
	s.x[10] = 0x20c
	s.x[11] = uint64(uint32(s.x[9]) | uint32(s.x[10]))
	value, err := s.load(s.sp+0xa0, 4)
	if err != nil {
		return 0, fmt.Errorf("account frame-guard stack word: %w", err)
	}
	s.x[10] = uint64(uint32(value))
	s.x[9] = 0x1cd95c3a
	s.compare32(uint32(s.x[10]), uint32(s.x[9]))
	s.x[9] = 0x038685c1
	s.x[9] = uint64(uint32(s.x[11]) ^ uint32(s.x[9]))
	if s.condition("eq") {
		s.compare32(uint32(s.x[28]), uint32(s.x[9]))
	} else {
		s.setNZCV(0)
	}
	s.x[9] = uint64(uint32(boolValue(s.condition("ne"))))
	s.x[12] = 0x980e
	s.x[9] = uint64(uint32(s.x[9]) * uint32(s.x[12]))
	s.x[8] = uint64(uint32(s.x[9]) ^ uint32(s.x[8]))
	next, err := loadBootstrapPointer(s, s.x[20], uint32(s.x[8]))
	if err != nil {
		return 0, fmt.Errorf("account frame-guard target: %w", err)
	}
	s.x[8] = next
	s.x[12] = 0xce757bc5
	s.x[17] = 0x85f4
	return next, nil
}

func runAccountReverseWriteRoute(s *state) (uint64, error) {
	s.traceBlock(accountReverseWriteRoute)
	s.x[13] = s.subtractWithFlags(s.x[13], 1, 64)
	value, err := s.load(s.x[11]+s.x[13], 1)
	if err != nil {
		return 0, fmt.Errorf("account reverse-write input: %w", err)
	}
	s.x[16] = uint64(byte(value))
	s.x[17] = uint64(uint32(s.x[14]) & uint32(s.x[16]<<1))
	s.x[16] = uint64(uint32(s.x[16]) - uint32(s.x[17]))
	s.x[16] = uint64(uint32(s.x[16]) + 0x1b)
	if err := s.storeValue(s.x[10]+s.x[13], uint64(byte(s.x[16])), 1); err != nil {
		return 0, fmt.Errorf("account reverse-write output: %w", err)
	}
	s.x[16] = uint64(uint32(boolValue(s.condition("eq"))))
	s.x[16] = uint64(uint32(s.x[16]) * uint32(s.x[15]))
	s.x[16] = uint64(uint32(s.x[16]) ^ uint32(s.x[12]))
	next, err := loadBootstrapPointer(s, s.x[9], uint32(s.x[16]))
	if err != nil {
		return 0, fmt.Errorf("account reverse-write target: %w", err)
	}
	s.x[0] = next
	s.x[17] = 0x2a6a
	return next, nil
}
