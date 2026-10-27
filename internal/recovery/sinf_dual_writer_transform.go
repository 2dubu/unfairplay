package recovery

import "fmt"

const sinfDualWriter = uint64(0x1000a5c70)

func runSINFDualWriter(s *state) (uint64, error) {
	s.traceBlock(sinfDualWriter)
	s.x[5] = 0
	s.x[11] = uint64(uint32(s.x[30]) + 0xffffd672)
	s.x[11] = uint64(uint32(s.x[11]) | 0x660)
	s.x[13] = uint64(uint32(s.x[9] << 1))
	s.x[15] = uint64(uint32(s.x[13]) & 0xfe)
	s.x[17] = uint64(uint32(s.x[11]) + 0xcbdad46b)
	s.x[15] = uint64(uint32(s.x[17]) ^ uint32(s.x[15]))
	s.x[17] = uint64(0xcbd0341f ^ (uint32(s.x[9]) >> 7))
	s.x[15] = uint64(uint32(s.x[15]) & uint32(s.x[17]))
	s.x[15] = uint64(uint32(s.x[15]) ^ 0xcbdafc3f)
	s.x[17] = uint64(uint32(s.x[15] << 1))
	s.x[6] = uint64(uint32(s.x[17]) & 0xfe)
	s.x[6] = insertBits(uint64(uint32(s.x[6])), uint64(uint32(s.x[13])>>7), 0, 1, 32)
	s.x[13] = uint64(uint32(s.x[17]) & 0x1e)
	s.x[13] = uint64(uint32(s.x[6]) ^ uint32(s.x[13]))
	s.x[13] = uint64(uint32(s.x[6]) - uint32(s.x[13]<<1))
	s.x[9] = uint64(uint32(s.x[15]) ^ uint32(s.x[9]))
	s.x[13] = uint64(uint32(s.x[13]) + 0x1e9d9ae1)
	s.x[9] = uint64(uint32(s.x[9]) ^ uint32(s.x[13]))
	s.x[13] = uint64((uint32(s.x[13]) & 0x7f) << 1)
	s.x[13] = uint64(uint32(s.x[13]) ^ 0xc2)
	s.x[15] = uint64(uint32(s.x[13]) >> 7)
	s.x[6] = uint64(uint32(s.x[13] << 1))
	s.x[13] = insertBits(uint64(uint32(s.x[13])), uint64(uint32(s.x[17])>>7), 0, 1, 32)
	s.x[15] = insertBits(uint64(uint32(s.x[15])), uint64(uint32(s.x[13])), 1, 7, 32)
	s.x[17] = uint64(uint32(s.x[15]) & 0x3f)
	s.x[7] = uint64(uint32(0xa) & uint32(s.x[13]<<1))
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[7]))
	s.x[15] = uint64(uint32(s.x[15]) - uint32(s.x[17]<<1))
	s.x[17] = uint64(uint32(s.x[6]) & 0xffffff07)
	s.x[13] = uint64(uint32(s.x[13]) - uint32(s.x[17]))
	s.x[13] = uint64(uint32(s.x[13]) + 0x844d8782)
	s.x[9] = uint64(uint32(s.x[9]) ^ uint32(s.x[13]))
	s.x[13] = uint64(uint32(s.x[15]) + 0x34372535)
	s.x[9] = uint64(uint32(s.x[9]) ^ uint32(s.x[13]))
	s.x[9] &= 0xff
	s.x[28] = s.x[4]
	value, err := s.load(s.x[4]+0x210, 4)
	if err != nil {
		return 0, fmt.Errorf("SINF dual-writer first word: %w", err)
	}
	s.x[13] = uint64(uint32(value))
	const stateBase = uint64(0x101075000)
	s.x[17] = stateBase
	value, err = s.load(stateBase+0xb40, 4)
	if err != nil {
		return 0, fmt.Errorf("SINF dual-writer second word: %w", err)
	}
	s.x[15] = uint64(uint32(value))
	s.x[13] = uint64(uint32(s.x[13]) ^ uint32(s.x[15]))
	s.x[13] = uint64(uint32(s.x[13]) ^ uint32(s.x[22]))
	s.x[13] = uint64(uint32(s.x[13]) * uint32(s.x[23]))
	firstBase, err := s.load(s.sp+0x18, 8)
	if err != nil {
		return 0, fmt.Errorf("SINF dual-writer first byte base: %w", err)
	}
	secondBase, err := s.load(s.sp+0x20, 8)
	if err != nil {
		return 0, fmt.Errorf("SINF dual-writer second byte base: %w", err)
	}
	s.x[15], s.x[6] = firstBase, secondBase
	s.x[15] += uint64(byte(s.x[13]))
	value, err = s.load(s.x[15], 1)
	if err != nil {
		return 0, fmt.Errorf("SINF dual-writer first byte: %w", err)
	}
	s.x[15] = uint64(uint32(byte(value)) ^ uint32(s.x[25]))
	value, err = s.load(s.x[6]+uint64(uint32(s.x[15])), 1)
	if err != nil {
		return 0, fmt.Errorf("SINF dual-writer second byte: %w", err)
	}
	s.x[15] = uint64(byte(value))
	s.x[6] = 0xd
	s.x[9] ^= s.x[6]
	s.x[13] = uint64(uint32(s.x[13]) ^ uint32(s.x[15]))
	s.x[13] &= 0xff
	s.x[13] = uint64(uint32(s.x[13]) + 0x6bf)
	value, err = s.load(s.x[10]+uint64(uint32(s.x[13]))*8, 8)
	if err != nil {
		return 0, fmt.Errorf("SINF dual-writer first state: %w", err)
	}
	s.x[13] = value
	s.x[6] = s.x[13] + s.x[26]
	value, err = s.load(s.x[6], 8)
	if err != nil {
		return 0, fmt.Errorf("SINF dual-writer first output base: %w", err)
	}
	s.x[13] = value + s.x[21]
	if err := s.storeValue(s.x[13], uint64(byte(s.x[9])), 1); err != nil {
		return 0, fmt.Errorf("SINF dual-writer first output: %w", err)
	}
	value, err = s.load(s.x[4]+0x210, 4)
	if err != nil {
		return 0, fmt.Errorf("SINF dual-writer third word: %w", err)
	}
	s.x[13] = uint64(uint32(value))
	value, err = s.load(stateBase+0xb40, 4)
	if err != nil {
		return 0, fmt.Errorf("SINF dual-writer fourth word: %w", err)
	}
	s.x[15] = uint64(uint32(value))
	s.x[13] = uint64(uint32(s.x[15]) - uint32(s.x[13]))
	s.x[13] = uint64(uint32(s.x[13]) ^ uint32(s.x[22]))
	s.x[13] = uint64(uint32(s.x[13]) * uint32(s.x[23]))
	value, err = s.load(s.x[10]+0x55c0, 8)
	if err != nil {
		return 0, fmt.Errorf("SINF dual-writer third byte base: %w", err)
	}
	s.x[15] = value
	s.x[7] = s.x[15] + s.x[27]
	s.x[15] = s.x[7] + uint64(byte(s.x[13]))
	value, err = s.load(s.x[15], 1)
	if err != nil {
		return 0, fmt.Errorf("SINF dual-writer third byte: %w", err)
	}
	s.x[15] = uint64(uint32(byte(value)) ^ uint32(s.x[1]))
	value, err = s.load(s.x[10]+0x2df0, 8)
	if err != nil {
		return 0, fmt.Errorf("SINF dual-writer fourth byte base: %w", err)
	}
	s.x[17] = value
	s.x[19] = s.x[17] + s.x[12]
	value, err = s.load(s.x[19]+uint64(uint32(s.x[15])), 1)
	if err != nil {
		return 0, fmt.Errorf("SINF dual-writer fourth byte: %w", err)
	}
	s.x[15] = uint64(byte(value))
	s.x[13] = uint64(uint32(s.x[15]) ^ uint32(s.x[13]))
	s.x[13] &= 0xff
	s.x[13] = uint64(uint32(s.x[13]) + 0xc62)
	value, err = s.load(s.x[10]+uint64(uint32(s.x[13]))*8, 8)
	if err != nil {
		return 0, fmt.Errorf("SINF dual-writer second state: %w", err)
	}
	s.x[13] = value
	s.x[20] = s.x[13] + s.x[0]
	value, err = s.load(s.x[20], 8)
	if err != nil {
		return 0, fmt.Errorf("SINF dual-writer second output base: %w", err)
	}
	s.x[13] = value
	s.x[9] += s.x[13]
	if err := s.storeValue(s.x[9], uint64(byte(s.x[21])), 1); err != nil {
		return 0, fmt.Errorf("SINF dual-writer second output: %w", err)
	}
	s.x[21]++
	s.compare64(s.x[21], 0x100)
	gate := uint32(boolValue(s.condition("ne")))
	s.x[9] = uint64(gate | gate<<2)
	s.x[13] = uint64(uint32(s.x[30]) | uint32(s.x[9]))
	s.x[9] = uint64(uint32(s.x[30]) + 0x598a)
	next, err := loadBootstrapPointer(s, s.x[14], uint32(s.x[13]))
	if err != nil {
		return 0, fmt.Errorf("SINF dual-writer dispatch: %w", err)
	}
	s.x[13] = next
	s.x[30] = s.x[24]
	s.x[17] = 0x3cb1
	return next, nil
}
