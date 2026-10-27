package recovery

import "fmt"

const (
	accountFrameRouteSetup = uint64(0x10005c474)
	accountFrameRouteCheck = uint64(0x100067abc)
)

func runAccountFrameRouteSetup(s *state) (uint64, error) {
	s.traceBlock(accountFrameRouteSetup)
	if err := s.storeValue(s.sp+0x18, uint64(uint32(s.x[0])), 4); err != nil {
		return 0, fmt.Errorf("account frame-route first stack word: %w", err)
	}
	if err := s.storeValue(s.sp+0x2c, uint64(uint32(s.x[25])), 4); err != nil {
		return 0, fmt.Errorf("account frame-route second stack word: %w", err)
	}
	if err := s.storeValue(s.sp+0x20, s.x[24], 8); err != nil {
		return 0, fmt.Errorf("account frame-route stack pointer: %w", err)
	}
	s.x[9] = 0xd27f8721
	s.x[11] = uint64(uint32(s.x[13]) ^ uint32(s.x[9]))
	value, err := s.load(addOffset(s.x[29], -0x7c), 1)
	if err != nil {
		return 0, fmt.Errorf("account frame-route byte: %w", err)
	}
	s.x[9] = uint64(byte(value))
	s.x[12] = uint64(uint32(s.x[13]) ^ uint32(s.x[3]))
	s.x[12] = uint64(uint32(shiftLeft(s.x[9], int(s.x[12]), 32)))
	s.x[12] = uint64(uint32(s.x[12]) & uint32(s.x[7]))
	s.x[9] = uint64(uint32(s.x[9]) - uint32(s.x[12]))
	s.x[9] = uint64(uint32(s.x[9]) - 0x16)
	s.x[16] = uint64(uint32(int32(int8(s.x[9]))))
	value, err = s.load(addOffset(s.x[29], -0x94), 4)
	if err != nil {
		return 0, fmt.Errorf("account frame-route first pair word: %w", err)
	}
	s.x[26] = uint64(uint32(value))
	value, err = s.load(addOffset(s.x[29], -0x90), 4)
	if err != nil {
		return 0, fmt.Errorf("account frame-route second pair word: %w", err)
	}
	s.x[0] = uint64(uint32(value))
	value, err = s.load(addOffset(s.x[29], -0x8c), 4)
	if err != nil {
		return 0, fmt.Errorf("account frame-route third pair word: %w", err)
	}
	s.x[22] = uint64(uint32(value))
	value, err = s.load(addOffset(s.x[29], -0x88), 4)
	if err != nil {
		return 0, fmt.Errorf("account frame-route fourth pair word: %w", err)
	}
	s.x[3] = uint64(uint32(value))
	value, err = s.load(addOffset(s.x[29], -0x84), 4)
	if err != nil {
		return 0, fmt.Errorf("account frame-route fifth word: %w", err)
	}
	s.x[9] = uint64(uint32(value))
	if err := s.storeValue(s.sp+0x7c, uint64(uint32(s.x[9])), 4); err != nil {
		return 0, fmt.Errorf("account frame-route fifth output: %w", err)
	}
	next, err := loadBootstrapPointer(s, s.x[21], uint32(s.x[13]))
	if err != nil {
		return 0, fmt.Errorf("account frame-route target: %w", err)
	}
	s.x[9] = next
	if err := s.storeValue(s.sp+0x38, uint64(uint32(s.x[16])), 4); err != nil {
		return 0, fmt.Errorf("account frame-route result word: %w", err)
	}
	if err := s.storeValue(s.sp+0x3c, uint64(uint32(s.x[4])), 4); err != nil {
		return 0, fmt.Errorf("account frame-route companion word: %w", err)
	}
	s.x[17] = 0xb54d
	return next, nil
}

func runAccountFrameRouteCheck(s *state) (uint64, error) {
	s.traceBlock(accountFrameRouteCheck)
	s.x[11] = 0xffff9690
	s.x[11] = uint64(uint32(s.x[9]) + uint32(s.x[11]))
	s.x[12] = 0x8201
	s.x[11] = uint64(uint32(s.x[11]) | uint32(s.x[12]))
	value, err := s.load(s.sp+0x38, 4)
	if err != nil {
		return 0, fmt.Errorf("account frame-route check first word: %w", err)
	}
	s.x[16] = uint64(uint32(value))
	value, err = s.load(s.sp+0x3c, 4)
	if err != nil {
		return 0, fmt.Errorf("account frame-route check second word: %w", err)
	}
	s.x[4] = uint64(uint32(value))
	s.x[12] = 0x1baa7ad1
	s.compare32(uint32(s.x[4]), uint32(s.x[12]))
	s.x[13] = uint64(uint32(boolValue(s.condition("lt"))))
	s.x[12] = 0xffff7970
	s.x[12] = uint64(uint32(s.x[11]) + uint32(s.x[12]))
	s.x[15] = uint64(uint32(s.x[12]) ^ 0x1bc6)
	s.x[13] = uint64(uint32(s.x[15]) * uint32(s.x[13]))
	s.x[9] = uint64(uint32(s.x[13]) ^ uint32(s.x[9]))
	next, err := loadBootstrapPointer(s, s.x[21], uint32(s.x[9]))
	if err != nil {
		return 0, fmt.Errorf("account frame-route check target: %w", err)
	}
	s.x[9] = next
	s.x[27] = 0x3b74
	s.x[17] = 0x64ef
	return next, nil
}
