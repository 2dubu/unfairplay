package recovery

import "fmt"

const (
	accountMultiplyTarget  = uint64(0x1001545c0)
	accountFrameByteTarget = uint64(0x100159f64)
)

func runAccountMultiplyTarget(s *state) (uint64, error) {
	s.traceBlock(accountMultiplyTarget)
	s.x[8] = 0x3b
	s.x[8] = uint64(uint32(s.x[11]) * uint32(s.x[8]))
	s.x[9] = 0x9bd3
	s.x[8] = uint64(uint32(s.x[8]) ^ uint32(s.x[9]))
	value, err := s.load(s.x[27]+uint64(uint32(s.x[8]))*8, 8)
	if err != nil {
		return 0, fmt.Errorf("account multiply-target target: %w", err)
	}
	s.x[8] = value
	s.x[10] = 0xce757bc5
	s.x[28] = 0x1b527405
	value, err = s.load(s.sp+0x2c, 4)
	if err != nil {
		return 0, fmt.Errorf("account multiply-target stack word: %w", err)
	}
	s.x[9] = uint64(uint32(value))
	s.x[11] = 0x26ef
	s.x[17] = 0xcdc6
	return s.x[8], nil
}

func runAccountFrameByteTarget(s *state) (uint64, error) {
	s.traceBlock(accountFrameByteTarget)
	if err := s.storeValue(s.sp+0x5c, uint64(uint32(s.x[19])), 4); err != nil {
		return 0, fmt.Errorf("account frame-byte saved word: %w", err)
	}
	s.x[8] = 0xffffc05c
	s.x[10] = uint64(uint32(s.x[9]) + uint32(s.x[8]))
	s.x[8] = 0x8c27722e
	s.x[12] = 0x73d7ff75
	s.x[8] = uint64(uint32(s.x[9]) + uint32(s.x[8]))
	s.x[8] = uint64(uint32(s.x[8]) & uint32(s.x[12]))
	value, err := s.load(s.x[30]+0x60, 8)
	if err != nil {
		return 0, fmt.Errorf("account frame-byte output pointer: %w", err)
	}
	s.x[12] = value
	if err := s.storeValue(s.x[12]+0x620, uint64(byte(s.x[11])), 1); err != nil {
		return 0, fmt.Errorf("account frame-byte output: %w", err)
	}
	value, err = s.load(s.x[30]+0x48, 4)
	if err != nil {
		return 0, fmt.Errorf("account frame-byte input: %w", err)
	}
	s.x[11] = uint64(uint32(value))
	s.x[12] = uint64(uint32(s.x[8]) ^ 0xce756abb)
	s.compare32(uint32(s.x[11]), uint32(s.x[12]))
	s.x[12] = uint64(uint32(boolValue(s.condition("eq"))))
	s.x[13] = 0x2b38
	s.x[12] = uint64(uint32(s.x[12]) * uint32(s.x[13]))
	s.x[9] = uint64(uint32(s.x[12]) ^ uint32(s.x[9]))
	value, err = s.load(s.x[27]+uint64(uint32(s.x[9]))*8, 8)
	if err != nil {
		return 0, fmt.Errorf("account frame-byte target: %w", err)
	}
	s.x[9] = value
	s.x[17] = 0x56fd
	return value, nil
}
