package recovery

import "fmt"

const accountRuntimeSetup = uint64(0x100a39d3c)

func runAccountRuntimeSetup(s *state) (uint64, error) {
	s.traceBlock(accountRuntimeSetup)
	next, err := loadBootstrapPointer(s, s.x[22], uint32(s.x[23]))
	if err != nil {
		return 0, fmt.Errorf("account runtime-setup target: %w", err)
	}
	s.x[13] = next

	s.x[10] = uint64(uint32(s.x[9]) + 0xf576a3eb)
	value, err := s.load(s.x[8]+uint64(uint32(s.x[10])), 1)
	if err != nil {
		return 0, fmt.Errorf("account runtime-setup first byte: %w", err)
	}
	s.x[10] = uint64(uint32(value) << 24)
	s.x[11] = uint64(uint32(s.x[9]) + 0xf576a3ec)
	value, err = s.load(s.x[8]+uint64(uint32(s.x[11])), 1)
	if err != nil {
		return 0, fmt.Errorf("account runtime-setup second byte: %w", err)
	}
	s.x[11] = uint64(byte(value))
	s.x[10] = uint64(uint32(s.x[10]) | uint32(s.x[11]<<16))
	s.x[11] = uint64(uint32(s.x[9]) + 0xf576a3ed)
	value, err = s.load(s.x[8]+uint64(uint32(s.x[11])), 1)
	if err != nil {
		return 0, fmt.Errorf("account runtime-setup third byte: %w", err)
	}
	s.x[11] = uint64(byte(value))
	s.x[10] = uint64(uint32(s.x[10]) | uint32(s.x[11]<<8))
	s.x[11] = 0x5b5e1000
	s.x[11] = uint64(uint32(s.x[11]) & uint32(s.x[10]<<1))
	s.x[12] = 0xf576a3ee
	s.x[9] = uint64(uint32(s.x[9]) + uint32(s.x[12]))
	value, err = s.load(s.x[8]+uint64(uint32(s.x[9])), 1)
	if err != nil {
		return 0, fmt.Errorf("account runtime-setup fourth byte: %w", err)
	}
	s.x[8] = uint64(byte(value))
	s.x[9] = uint64(uint32(s.x[10]) - uint32(s.x[11]))
	s.x[10] = 0xadaf0894
	s.x[9] = uint64(uint32(s.x[9]) + uint32(s.x[10]))
	s.x[10] = uint64((uint32(s.x[8]) & 0x1f) << 1)
	s.x[8] = uint64(uint32(s.x[8]) - uint32(s.x[10]))
	s.x[10] = 0x5aa17d1f
	s.x[8] = uint64(uint32(s.x[8]) + uint32(s.x[10]))
	s.x[10] = 0xdb967a69
	s.x[11] = uint64(uint32(s.x[9]) ^ uint32(s.x[10]))
	s.x[9] = 0xa55e82e0
	s.x[9] = uint64(uint32(s.x[8]) ^ uint32(s.x[9]))
	s.x[12] = uint64(uint32(s.x[11]) & uint32(s.x[9]))
	s.x[9] = 0x08000d02
	s.x[8] = uint64(uint32(s.x[8]) & uint32(s.x[9]))
	s.x[10] = uint64(uint32(s.x[12]) ^ uint32(s.x[8]))
	s.x[8] = 0xb04c043a
	s.x[20] = uint64(uint32(s.x[10]) ^ uint32(s.x[8]))

	s.x[9] = 0x10105b000
	s.x[9] += 0xf18
	value, err = s.load(s.x[9]+0x66e8, 8)
	if err != nil {
		return 0, fmt.Errorf("account runtime-setup first base: %w", err)
	}
	s.x[8] = value
	s.x[14] = 0xffffffff954bc2a5
	s.x[6] = s.x[8] + s.x[14]
	value, err = s.load(s.x[9]+0x12d0, 8)
	if err != nil {
		return 0, fmt.Errorf("account runtime-setup second base: %w", err)
	}
	s.x[8] = value
	s.x[14] = 0xffffffffcd114862
	s.x[27] = s.x[8] + s.x[14]
	value, err = s.load(s.x[9]+0x2230, 8)
	if err != nil {
		return 0, fmt.Errorf("account runtime-setup third base: %w", err)
	}
	s.x[8] = value
	s.x[14] = 0xffffffffe5b0614e
	s.x[24] = s.x[8] + s.x[14]
	value, err = s.load(s.x[9]+0x4d98, 8)
	if err != nil {
		return 0, fmt.Errorf("account runtime-setup fourth base: %w", err)
	}
	s.x[8] = value
	s.x[14] = 0xffffffffc10da2f9
	s.x[5] = s.x[8] + s.x[14]
	value, err = s.load(s.x[9]+0x25c8, 8)
	if err != nil {
		return 0, fmt.Errorf("account runtime-setup fifth base: %w", err)
	}
	s.x[8] = value
	s.x[14] = 0xffffffff9f79fd81
	s.x[4] = s.x[8] + s.x[14]
	value, err = s.load(s.x[9]+0x31e0, 8)
	if err != nil {
		return 0, fmt.Errorf("account runtime-setup sixth base: %w", err)
	}
	s.x[8] = value
	s.x[14] = 0xfffffffffdc52ec2
	s.x[2] = s.x[8] + s.x[14]
	value, err = s.load(s.x[9], 8)
	if err != nil {
		return 0, fmt.Errorf("account runtime-setup seventh base: %w", err)
	}
	s.x[8] = value
	s.x[9] = 0xffffffffba636d86
	s.x[0] = s.x[8] + s.x[9]
	s.x[17] = 0x927e
	return next, nil
}
