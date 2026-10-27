package recovery

import "fmt"

const accountWordWriter = uint64(0x1001360b0)

func runAccountWordWriter(s *state) (uint64, error) {
	s.traceBlock(accountWordWriter)
	s.x[17] = s.x[0] + s.x[3]
	value, err := s.load(s.x[21]+s.x[17], 1)
	if err != nil {
		return 0, fmt.Errorf("account word-writer first byte: %w", err)
	}
	s.x[17] = uint64(byte(value))
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[4]))
	s.x[17] = uint64(uint32(s.x[17]) << 24)
	s.x[1] = s.x[0] + s.x[5]
	value, err = s.load(s.x[21]+s.x[1], 1)
	if err != nil {
		return 0, fmt.Errorf("account word-writer second byte: %w", err)
	}
	s.x[1] = uint64(byte(value))
	s.x[1] = uint64(uint32(s.x[1]) ^ uint32(s.x[4]))
	s.x[17] = uint64(uint32(s.x[17]) | uint32(s.x[1]<<16))
	s.x[1] = s.x[0] + s.x[24]
	value, err = s.load(s.x[21]+s.x[1], 1)
	if err != nil {
		return 0, fmt.Errorf("account word-writer third byte: %w", err)
	}
	s.x[1] = uint64(byte(value))
	s.x[2] = uint64(uint32(s.x[13]) ^ uint32(s.x[9]))
	s.x[1] = uint64(uint32(s.x[2]) ^ uint32(s.x[1]))
	s.x[17] = insertBits(uint64(uint32(s.x[17])), uint64(uint32(s.x[1])), 8, 8, 32)
	s.x[1] = s.x[0] + s.x[25]
	value, err = s.load(s.x[21]+s.x[1], 1)
	if err != nil {
		return 0, fmt.Errorf("account word-writer fourth byte: %w", err)
	}
	s.x[1] = uint64(byte(value))
	s.x[1] = uint64(uint32(s.x[1]) ^ uint32(s.x[4]))
	s.x[17] = uint64(uint32(s.x[17]) | uint32(s.x[1]))
	s.x[2] = uint64(uint32(s.x[17]) & uint32(s.x[23]))
	s.x[1] &= 0x10
	s.x[1] = uint64(uint32(s.x[2]) ^ uint32(s.x[1]))
	s.x[17] = uint64(uint32(s.x[17]) - uint32(s.x[1]<<1))
	s.x[17] = uint64(uint32(s.x[17]) + uint32(s.x[12]))
	s.x[1] = uint64(uint32(s.x[0]) ^ uint32(s.x[26]))
	s.x[2] = uint64(uint32(s.x[0]) ^ 0xbcc10a90)
	s.x[7] = 0x5ff6bdfd
	s.x[2] = uint64(uint32(s.x[2]) + uint32(s.x[7]))
	s.x[7] = uint64(uint32(s.x[0]) ^ 0xea1b2630)
	s.x[14] = 0x5bfbfbfc
	s.x[30] = uint64(uint32(s.x[0]) ^ uint32(s.x[14]))
	s.x[1] = uint64(uint32(s.x[1]) + uint32(s.x[27]))
	s.x[1] = uint64(uint32(s.x[1]) ^ uint32(s.x[0]))
	s.x[14] = 0x092c915e
	s.x[7] = uint64(uint32(s.x[7]) + uint32(s.x[14]))
	s.x[2] = uint64(uint32(s.x[2]) ^ uint32(s.x[7]))
	s.x[1] = uint64(uint32(s.x[1]) ^ uint32(s.x[2]))
	s.x[14] = 0xb8cc4c94
	s.x[2] = uint64(uint32(s.x[30]) + uint32(s.x[14]))
	s.x[30] = 0x49d83ed2
	s.x[1] = uint64(uint32(s.x[1]) ^ uint32(s.x[2]))
	s.x[1] = uint64(uint32(s.x[1]) >> 2)
	s.x[2] = 0x07321224
	s.x[1] ^= s.x[2]
	if err := s.storeValue(s.x[8]+(s.x[1]<<2), uint64(uint32(s.x[17])), 4); err != nil {
		return 0, fmt.Errorf("account word-writer output: %w", err)
	}
	s.x[17] = 0xffffffffe337b771
	s.x[17] += s.x[0]
	s.x[0] += 4
	s.compare64(s.x[17], 0x3f)
	gate := uint32(boolValue(s.condition("hi")))
	s.x[17] = uint64(gate * 0xc1)
	s.x[14] = 0xffffffff8cdeddb4
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[16]))
	next, err := loadBootstrapPointer(s, s.x[20], uint32(s.x[17]))
	if err != nil {
		return 0, fmt.Errorf("account word-writer target: %w", err)
	}
	s.x[2] = next
	s.x[1] = 0x1cc848a3
	s.x[17] = 0xe949
	return next, nil
}
