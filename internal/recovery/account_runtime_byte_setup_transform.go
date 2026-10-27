package recovery

import "fmt"

const accountRuntimeByteSetup = uint64(0x100103374)

func runAccountRuntimeByteSetup(s *state) (uint64, error) {
	s.traceBlock(accountRuntimeByteSetup)
	s.x[8] = 0x1163
	s.x[12] = 0xd0c0
	s.x[12] = s.x[27] + s.x[12]
	value, err := s.load(s.x[12], 8)
	if err != nil {
		return 0, fmt.Errorf("account runtime-byte dispatch: %w", err)
	}
	s.x[12] = value
	value, err = s.load(s.x[5], 4)
	if err != nil {
		return 0, fmt.Errorf("account runtime-byte first word: %w", err)
	}
	s.x[13] = uint64(uint32(value))
	value, err = s.load(s.x[3]+0xb40, 4)
	if err != nil {
		return 0, fmt.Errorf("account runtime-byte second word: %w", err)
	}
	s.x[14] = uint64(uint32(value))
	s.x[13] = uint64(uint32(s.x[14]) + uint32(s.x[13]))
	s.x[9] = uint64(uint32(s.x[13]) ^ uint32(s.x[9]))
	s.x[9] = uint64(uint32(s.x[9]) * uint32(s.x[11]))
	value, err = s.load(s.x[4]+0x3a08, 8)
	if err != nil {
		return 0, fmt.Errorf("account runtime-byte first base: %w", err)
	}
	s.x[11] = value + uint64(byte(s.x[9]))
	s.x[13] = 0xfffffffffdc52ec2
	s.x[11] += s.x[13]
	value, err = s.load(s.x[11], 1)
	if err != nil {
		return 0, fmt.Errorf("account runtime-byte first byte: %w", err)
	}
	s.x[11] = uint64(byte(value))
	s.x[10] = uint64(uint32(s.x[11]) ^ uint32(s.x[10]))
	value, err = s.load(s.x[4]+0x828, 8)
	if err != nil {
		return 0, fmt.Errorf("account runtime-byte second base: %w", err)
	}
	s.x[11] = value
	s.x[13] = 0xffffffffba636d86
	s.x[11] += s.x[13]
	s.x[10] += s.x[11]
	value, err = s.load(s.x[10], 1)
	if err != nil {
		return 0, fmt.Errorf("account runtime-byte second byte: %w", err)
	}
	s.x[10] = uint64(byte(value))
	s.x[9] = uint64(uint32(s.x[10]) ^ uint32(s.x[9]))
	s.x[9] &= 0xff
	s.x[9] = uint64(uint32(s.x[9]) + 0xb1)
	value, err = s.load(s.x[4]+uint64(uint32(s.x[9]))*8, 8)
	if err != nil {
		return 0, fmt.Errorf("account runtime-byte target value: %w", err)
	}
	s.x[9] = value
	s.x[10] = 0xffffffff8b9e3eb6
	s.x[9] += s.x[10]
	s.x[28] = 0x1b527405
	s.x[10] = 0x663d03d3
	s.x[17] = 0xd3f7
	return s.x[12], nil
}
