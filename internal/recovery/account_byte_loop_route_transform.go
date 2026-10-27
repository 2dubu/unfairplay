package recovery

import "fmt"

const (
	accountTableOffsetRoute = uint64(0x100822168)
	accountByteLoopRoute    = uint64(0x10090fcdc)
	accountBitByteLoopRoute = uint64(0x100620688)
)

func runAccountTableOffsetRoute(s *state) (uint64, error) {
	s.traceBlock(accountTableOffsetRoute)
	s.x[10] = 0xffffd2f0
	s.x[10] = uint64(uint32(s.x[24]) + uint32(s.x[10]))
	next, err := loadBootstrapPointer(s, s.x[22], uint32(s.x[10]))
	if err != nil {
		return 0, fmt.Errorf("account table-offset target: %w", err)
	}
	s.x[10] = next
	s.x[17] = 0x111b
	return next, nil
}

func runAccountByteLoopRoute(s *state) (uint64, error) {
	s.traceBlock(accountByteLoopRoute)
	s.x[6] = s.x[10] + s.x[11]
	value, err := s.load(s.x[12]+s.x[6], 1)
	if err != nil {
		return 0, fmt.Errorf("account byte-loop input: %w", err)
	}
	s.x[7] = uint64(byte(value))
	s.x[17] = uint64(uint32(s.x[16]) + 0x75)
	s.compare32(uint32(s.x[7]), uint32(byte(s.x[17])))
	s.x[17] = uint64(uint32(boolValue(s.condition("ne"))))
	s.x[17] = uint64(uint32(s.x[17]) * uint32(s.x[0]))
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[15]))
	value, err = s.load(s.sp+0xdf0, 8)
	if err != nil {
		return 0, fmt.Errorf("account byte-loop table: %w", err)
	}
	s.x[2] = value
	next, err := loadBootstrapPointer(s, s.x[2], uint32(s.x[17]))
	if err != nil {
		return 0, fmt.Errorf("account byte-loop target: %w", err)
	}
	s.x[19] = next
	s.x[2] = s.x[16]
	s.x[17] = 0x50d5
	return next, nil
}

func runAccountBitByteLoopRoute(s *state) (uint64, error) {
	s.traceBlock(accountBitByteLoopRoute)
	s.x[2] = uint64(uint32(s.x[1]) | uint32(s.x[3]))
	s.x[10] = uint64(uint32(s.x[2]) + 0x76)
	s.x[10] = uint64(uint32(shiftLeft(s.x[7], int(s.x[10]), 32)))
	s.x[10] = uint64(uint32(s.x[10]) & uint32(s.x[4]))
	s.x[10] = uint64(uint32(s.x[7]) - uint32(s.x[10]))
	s.x[10] = uint64(uint32(s.x[10]) + 0x1b)
	if err := s.storeValue(s.x[14]+s.x[6], uint64(byte(s.x[10])), 1); err != nil {
		return 0, fmt.Errorf("account bit-byte output: %w", err)
	}
	value, err := s.load(s.sp+0x3ae8, 8)
	if err != nil {
		return 0, fmt.Errorf("account bit-byte counter: %w", err)
	}
	s.x[17] = value
	s.x[10] = s.x[17] + s.x[11]
	value, err = s.load(s.x[12]+s.x[10], 1)
	if err != nil {
		return 0, fmt.Errorf("account bit-byte input: %w", err)
	}
	s.x[10] = uint64(byte(value))
	s.compare32(uint32(s.x[10]), 0x2f)
	if s.condition("eq") {
		s.x[9] = s.x[17]
	}
	s.x[10] = s.x[17] + 1
	if err := s.storeValue(s.sp+0x3ae8, s.x[10], 8); err != nil {
		return 0, fmt.Errorf("account bit-byte counter output: %w", err)
	}
	s.x[17] += s.x[8]
	s.compare64(s.x[17], 0x400)
	s.x[17] = uint64(uint32(boolValue(s.condition("lo"))))
	s.x[17] = uint64(uint32(s.x[17]) * uint32(s.x[5]))
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[15]))
	value, err = s.load(s.sp+0xdf0, 8)
	if err != nil {
		return 0, fmt.Errorf("account bit-byte table: %w", err)
	}
	s.x[6] = value
	next, err := loadBootstrapPointer(s, s.x[6], uint32(s.x[17]))
	if err != nil {
		return 0, fmt.Errorf("account bit-byte target: %w", err)
	}
	s.x[6] = next
	s.x[17] = 0x50d5
	return next, nil
}
