package recovery

import "fmt"

const (
	sinfStateSelector           = uint64(0x10008ecd0)
	sinfBitSelector             = uint64(0x10000f8fc)
	sharedEnrollmentCodecWriter = uint64(0x100119cd8)
)

func runSINFStateSelector(s *state) (uint64, error) {
	s.traceBlock(sinfStateSelector)
	s.x[4], s.x[24] = s.x[28], s.x[30]
	s.x[11] = uint64(uint32(s.x[9]) + 0xcbcafa21)
	s.x[30] = uint64(uint32(s.x[11]) ^ 0xcbcbd5af)
	s.x[11] = uint64(uint32(s.x[11]) & 0x3434fb5e)
	s.compare64(s.x[21], 0)
	gate := uint32(boolValue(s.condition("eq")))
	s.x[13] = uint64(uint32(s.x[30]) ^ gate)
	s.x[15] = 1
	s.x[13] = uint64(uint32(s.x[15]) &^ uint32(s.x[13]))
	s.x[15] = uint64(uint32(s.x[11]) ^ 0x9bc9)
	s.x[13] = uint64(uint32(s.x[13]) * uint32(s.x[15]))
	s.x[9] = uint64(uint32(s.x[13]) ^ uint32(s.x[9]))
	next, err := loadBootstrapPointer(s, s.x[14], uint32(s.x[9]))
	if err != nil {
		return 0, fmt.Errorf("SINF state-selector target: %w", err)
	}
	s.x[13] = next
	s.x[9] = 0x07647038
	s.x[17] = 0x5551
	return next, nil
}

func runSINFBitSelector(s *state) (uint64, error) {
	s.traceBlock(sinfBitSelector)
	s.x[9] = uint64(uint32(s.x[11]) ^ 0xbf95)
	s.x[30] = uint64(uint32(s.x[9] << 1))
	value, err := s.load(s.x[2], 8)
	if err != nil {
		return 0, fmt.Errorf("SINF bit-selector first base: %w", err)
	}
	s.x[9] = value
	value, err = s.load(s.x[3], 8)
	if err != nil {
		return 0, fmt.Errorf("SINF bit-selector second base: %w", err)
	}
	s.x[13] = value + s.x[21]
	value, err = s.load(s.x[13], 1)
	if err != nil {
		return 0, fmt.Errorf("SINF bit-selector first byte: %w", err)
	}
	s.x[13] = uint64(^uint32(byte(value)))
	s.x[9] += uint64(byte(s.x[13]))
	value, err = s.load(s.x[9], 1)
	if err != nil {
		return 0, fmt.Errorf("SINF bit-selector second byte: %w", err)
	}
	s.x[9] = uint64(byte(value))
	s.x[13] = uint64(uint32(s.x[9]) & 0xffffff90)
	s.x[13] = uint64(uint32(s.x[13]) ^ 0xffffff86)
	s.x[15] = uint64(uint32(0x28) & uint32(s.x[9]<<1))
	s.x[15] = uint64(uint32(s.x[15]) | 4)
	s.x[17] = uint64(uint32(s.x[30]) + 2)
	s.x[17] = uint64(uint32(s.x[9]) ^ uint32(s.x[17]))
	s.x[15] = uint64(uint32(s.x[15]) & uint32(s.x[17]))
	s.x[15] = uint64(uint32(s.x[15]) ^ uint32(s.x[13]))
	for range 2 {
		s.x[15] = uint64(uint32(s.x[17]) & uint32(s.x[15]<<1))
		s.x[15] = uint64(uint32(s.x[15]) ^ uint32(s.x[13]))
	}
	s.x[15] = uint64(uint32(s.x[17]) & uint32(s.x[15]<<1))
	s.x[13] = uint64(uint32(s.x[15]) ^ uint32(s.x[13]))
	s.x[9] = uint64(uint32(s.x[9]) ^ uint32(s.x[13]<<1))
	s.x[9] &= 0xff
	s.x[13] = 0x076470a8
	s.x[9] = uint64(uint32(s.x[9]) ^ uint32(s.x[13]))
	value, err = s.load(s.x[14]+uint64(uint32(s.x[11]))*8, 8)
	if err != nil {
		return 0, fmt.Errorf("SINF bit-selector dispatch: %w", err)
	}
	s.x[11] = value
	s.x[17] = 0x970e
	return value, nil
}

func runSharedEnrollmentCodecWriter(s *state) (uint64, error) {
	s.traceBlock(sharedEnrollmentCodecWriter)
	s.x[9] = s.subtractWithFlags(s.x[9], 1, 64)
	value, err := s.load(s.x[14]+s.x[9], 1)
	if err != nil {
		return 0, fmt.Errorf("shared codec-writer input: %w", err)
	}
	s.x[15] = uint64(byte(value))
	value, err = s.load(s.x[27]+0xfa0, 4)
	if err != nil {
		return 0, fmt.Errorf("shared codec-writer first word: %w", err)
	}
	s.x[17] = uint64(uint32(value))
	value, err = s.load(s.x[28]+0xb40, 4)
	if err != nil {
		return 0, fmt.Errorf("shared codec-writer second word: %w", err)
	}
	s.x[21] = uint64(uint32(value))
	s.x[21] = uint64(uint32(s.x[21]) + uint32(s.x[17]))
	s.x[21] = uint64(uint32(s.x[21]) ^ uint32(s.x[0]))
	s.x[21] = uint64(uint32(s.x[21]) * uint32(s.x[1]))
	s.x[22] = s.x[26] + uint64(byte(s.x[21]))
	value, err = s.load(s.x[22], 1)
	if err != nil {
		return 0, fmt.Errorf("shared codec-writer first table byte: %w", err)
	}
	s.x[22] = uint64(uint32(byte(value)) ^ uint32(s.x[2]))
	value, err = s.load(s.x[12]+uint64(uint32(s.x[22])), 1)
	if err != nil {
		return 0, fmt.Errorf("shared codec-writer second table byte: %w", err)
	}
	s.x[22] = uint64(byte(value))
	s.x[23] = s.x[16] & s.x[9]
	s.x[21] = uint64(uint32(s.x[21]) ^ uint32(s.x[22]))
	value, err = s.load(s.x[24], 4)
	if err != nil {
		return 0, fmt.Errorf("shared codec-writer dynamic word: %w", err)
	}
	s.x[22] = uint64(uint32(value))
	s.x[17] = uint64(uint32(s.x[22]) - uint32(s.x[17]))
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[0]))
	s.x[21] &= 0xff
	s.x[17] = uint64(uint32(s.x[17]) * uint32(s.x[1]))
	s.x[22] = s.x[17] & 0xff
	s.x[8] = s.x[24]
	s.x[24] = s.x[10] + s.x[22]
	value, err = s.load(s.x[24], 1)
	if err != nil {
		return 0, fmt.Errorf("shared codec-writer third table byte: %w", err)
	}
	s.x[24] = uint64(uint32(byte(value)) ^ 0x10)
	s.x[21] = uint64(uint32(s.x[21]) + 0xf9b)
	value, err = s.load(s.x[11]+uint64(uint32(s.x[24])), 1)
	if err != nil {
		return 0, fmt.Errorf("shared codec-writer fourth table byte: %w", err)
	}
	s.x[24] = uint64(byte(value))
	s.x[24] = uint64(uint32(s.x[17]) ^ uint32(s.x[24]))
	s.x[24] &= 0xff
	firstPointer, err := s.load(s.x[30]+uint64(uint32(s.x[21]))*8, 8)
	if err != nil {
		return 0, fmt.Errorf("shared codec-writer first pointer: %w", err)
	}
	s.x[21] = firstPointer
	s.x[24] = uint64(uint32(s.x[24]) + 0xde4)
	secondPointer, err := s.load(s.x[30]+uint64(uint32(s.x[24]))*8, 8)
	if err != nil {
		return 0, fmt.Errorf("shared codec-writer second pointer: %w", err)
	}
	s.x[24] = secondPointer
	s.x[25] = s.x[23] + s.x[4]
	s.x[26] = s.x[23] + s.x[3]
	s.x[21] += s.x[26]
	value, err = s.load(addOffset(s.x[29], -0x90), 8)
	if err != nil {
		return 0, fmt.Errorf("shared codec-writer saved base: %w", err)
	}
	s.x[26] = value
	s.x[24] += s.x[25]
	s.x[25] = uint64(uint32(s.x[23]) * uint32(s.x[5]))
	s.x[22] = s.x[26] + s.x[22]
	value, err = s.load(s.x[22], 1)
	if err != nil {
		return 0, fmt.Errorf("shared codec-writer fifth table byte: %w", err)
	}
	s.x[22] = uint64(uint32(byte(value)) ^ uint32(s.x[6]))
	value, err = s.load(s.x[21], 1)
	if err != nil {
		return 0, fmt.Errorf("shared codec-writer fifth lane byte: %w", err)
	}
	s.x[21] = uint64(byte(value))
	value, err = s.load(s.x[12]+uint64(uint32(s.x[22])), 1)
	if err != nil {
		return 0, fmt.Errorf("shared codec-writer sixth table byte: %w", err)
	}
	s.x[22] = uint64(byte(value))
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[22]))
	s.x[17] &= 0xff
	s.x[17] = uint64(uint32(s.x[17]) + 0xa0b)
	value, err = s.load(s.x[24], 1)
	if err != nil {
		return 0, fmt.Errorf("shared codec-writer seventh lane byte: %w", err)
	}
	s.x[22] = uint64(byte(value))
	s.x[24] = s.x[8]
	thirdPointer, err := s.load(s.x[30]+uint64(uint32(s.x[17]))*8, 8)
	if err != nil {
		return 0, fmt.Errorf("shared codec-writer third pointer: %w", err)
	}
	s.x[17] = thirdPointer
	s.x[23] += s.x[7]
	s.x[17] += s.x[23]
	value, err = s.load(s.x[17], 1)
	if err != nil {
		return 0, fmt.Errorf("shared codec-writer seventh table byte: %w", err)
	}
	s.x[17] = uint64(byte(value))
	s.x[15] = uint64(uint32(s.x[15]) ^ uint32(s.x[25]))
	s.x[15] = uint64(uint32(s.x[15]) ^ uint32(s.x[21]))
	s.x[17] = uint64(uint32(s.x[22]) ^ uint32(s.x[17]))
	s.x[15] = uint64(uint32(s.x[15]) ^ uint32(s.x[17]))
	s.x[15] = uint64(uint32(s.x[15]) ^ uint32(s.x[19]))
	if err := s.storeValue(s.x[14]+s.x[9], uint64(byte(s.x[15])), 1); err != nil {
		return 0, fmt.Errorf("shared codec-writer output: %w", err)
	}
	gate := uint32(boolValue(s.condition("eq")))
	s.x[15] = uint64(gate * uint32(s.x[20]))
	s.x[15] = uint64(uint32(s.x[15]) ^ uint32(s.x[13]))
	dispatchTable, err := s.load(addOffset(s.x[29], -0x80), 8)
	if err != nil {
		return 0, fmt.Errorf("shared codec-writer dispatch table: %w", err)
	}
	s.x[17] = dispatchTable
	next, err := loadBootstrapPointer(s, dispatchTable, uint32(s.x[15]))
	if err != nil {
		return 0, fmt.Errorf("shared codec-writer dispatch: %w", err)
	}
	s.x[15] = next
	s.x[17] = 0x8182
	return next, nil
}
