package recovery

import "fmt"

const (
	bootstrapRoundFront      = uint64(0x10093ffd0)
	bootstrapRoundTail       = uint64(0x100940000)
	bootstrapPreparedStart   = uint64(0x100344b1c)
	bootstrapPreparedTail    = uint64(0x100344c00)
	bootstrapPreparedTable   = uint64(0x101075000)
	bootstrapPreparedBytesA  = uint64(0x100f49000)
	bootstrapPreparedBytesB  = uint64(0x100fe9000)
	bootstrapPreparedTargets = uint64(0x10105b000)
)

func runBootstrapRound(s *state) (uint64, error) {
	s.traceBlock(bootstrapRoundFront)
	s.x[9] = s.subtractWithFlags(s.x[9], 1, 64)
	word, err := s.load(s.x[7], 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap round input word: %w", err)
	}
	s.x[10] = uint64(uint32(word))
	word, err = s.load(s.x[6]+0xb40, 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap round state word: %w", err)
	}
	s.x[17] = uint64(uint32(word))
	s.x[5] = uint64(uint32(s.x[17]) + uint32(s.x[10]))
	s.x[5] = uint64(uint32(s.x[5]) ^ uint32(s.x[13]))
	s.x[5] = uint64(uint32(s.x[5]) * uint32(s.x[14]))
	s.x[21] = s.x[6]
	s.x[6] = s.x[5] & 0xff
	s.x[7] = s.x[28] + s.x[6]
	value, err := s.load(s.x[7], 1)
	if err != nil {
		return 0, fmt.Errorf("bootstrap round first byte: %w", err)
	}
	s.x[7] = uint64(byte(value)) ^ 4
	value, err = s.load(s.x[26]+uint64(uint32(s.x[7])), 1)
	if err != nil {
		return 0, fmt.Errorf("bootstrap round first substitution: %w", err)
	}
	s.x[7] = uint64(byte(value))

	s.traceBlock(bootstrapRoundTail)
	s.x[7] = uint64(uint32(s.x[5]) ^ uint32(s.x[7]))
	s.x[10] = uint64(uint32(s.x[10]) ^ uint32(s.x[13]))
	s.x[10] = uint64(uint32(s.x[10]) ^ uint32(s.x[17]))
	s.x[10] = uint64(uint32(s.x[10]) * uint32(s.x[14]))
	s.x[17] = s.x[24] + uint64(byte(s.x[10]))
	s.x[7] = uint64(uint32(s.x[15]) + uint32(byte(s.x[7])))
	byteValue, err := s.load(s.x[17], 1)
	if err != nil {
		return 0, fmt.Errorf("bootstrap round second state byte: %w", err)
	}
	s.x[17] = uint64(byte(byteValue))
	s.x[17] = uint64(uint32(s.x[17]) ^ 0x10)
	value, err = s.load(s.x[22]+uint64(uint32(s.x[17])), 1)
	if err != nil {
		return 0, fmt.Errorf("bootstrap round second substitution: %w", err)
	}
	s.x[17] = uint64(byte(value))
	value, err = s.load(s.x[19]+uint64(uint32(s.x[7]))*8, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap round fourth dispatch: %w", err)
	}
	s.x[7] = value
	s.x[12] = s.x[19]
	s.x[19] = s.x[9] & 0xf
	s.x[10] = uint64(uint32(s.x[10]) ^ uint32(s.x[17]))
	s.x[17] = s.x[19] + s.x[16]
	s.x[10] &= 0xff
	s.x[6] = s.x[25] + s.x[6]
	s.x[10] = uint64(uint32(s.x[10]) + 0x6b2)
	byteValue, err = s.load(s.x[6], 1)
	if err != nil {
		return 0, fmt.Errorf("bootstrap round first byte: %w", err)
	}
	s.x[6] = uint64(byte(byteValue))
	s.x[6] = uint64(uint32(s.x[6]) ^ uint32(s.x[1]))
	value, err = s.load(s.x[27]+uint64(uint32(s.x[6])), 1)
	if err != nil {
		return 0, fmt.Errorf("bootstrap round substituted byte: %w", err)
	}
	s.x[6] = uint64(byte(value))
	value, err = s.load(s.x[12]+uint64(uint32(s.x[10]))*8, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap round fifth dispatch: %w", err)
	}
	s.x[10] = value
	s.x[17] += s.x[7]
	s.x[5] = uint64(uint32(s.x[5]) ^ uint32(s.x[6]))
	s.x[5] &= 0xff
	s.x[5] = uint64(uint32(s.x[5]) + 0x75e)
	s.x[6] = s.x[19] + s.x[0]
	s.x[10] += s.x[6]
	value, err = s.load(s.x[12]+uint64(uint32(s.x[5]))*8, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap round sixth dispatch: %w", err)
	}
	s.x[5] = value
	s.x[6] = s.x[19] + s.x[2]
	s.x[5] += s.x[6]
	byteValue, err = s.load(s.x[23]+s.x[9], 1)
	if err != nil {
		return 0, fmt.Errorf("bootstrap round stream byte: %w", err)
	}
	s.x[6] = uint64(byte(byteValue))
	s.x[7] = uint64(uint32(s.x[19]) * uint32(s.x[3]))
	s.x[19] = s.x[12]
	s.x[6] = uint64(uint32(s.x[6]) ^ uint32(s.x[7]))
	value, err = s.load(s.sp+0x200, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap round output: %w", err)
	}
	s.x[7] = value
	byteValue, err = s.load(s.x[10], 1)
	if err != nil {
		return 0, fmt.Errorf("bootstrap round fifth byte: %w", err)
	}
	s.x[10] = uint64(byte(byteValue))
	s.x[10] = uint64(uint32(s.x[6]) ^ uint32(s.x[10]))
	s.x[6] = s.x[21]
	byteValue, err = s.load(s.x[17], 1)
	if err != nil {
		return 0, fmt.Errorf("bootstrap round fourth byte: %w", err)
	}
	s.x[17] = uint64(byte(byteValue))
	s.x[17] = uint64(uint32(s.x[10]) ^ uint32(s.x[17]))
	s.x[10] = 0x79377ddd
	s.x[10] = uint64(uint32(s.x[11]) & uint32(s.x[10]))
	byteValue, err = s.load(s.x[5], 1)
	if err != nil {
		return 0, fmt.Errorf("bootstrap round sixth byte: %w", err)
	}
	s.x[5] = uint64(byte(byteValue))
	s.x[5] = uint64(uint32(s.x[5]) ^ uint32(s.x[4]))
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[5]))
	if err := s.storeValue(s.x[20]+s.x[9], uint64(byte(s.x[17])), 1); err != nil {
		return 0, fmt.Errorf("bootstrap round output byte: %w", err)
	}
	gate := uint32(boolValue(s.condition("eq")))
	s.x[17] = uint64(gate)
	s.x[5] = uint64(uint32(s.x[10]) - 0xc1)
	s.x[17] = uint64(gate * uint32(s.x[5]))
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[8]))
	next, err := loadBootstrapPointer(s, s.x[30], uint32(s.x[17]))
	if err != nil {
		return 0, fmt.Errorf("bootstrap round target: %w", err)
	}
	s.x[5] = next
	s.x[17] = 0x5e79
	return next, nil
}

func runBootstrapPreparedStart(s *state) (uint64, error) {
	s.traceBlock(bootstrapPreparedStart)
	s.sp -= 0x20
	if err := s.storeValue(s.sp+0x10, s.x[20], 8); err != nil {
		return 0, fmt.Errorf("bootstrap prepared frame register 20: %w", err)
	}
	if err := s.storeValue(s.sp+0x18, s.x[19], 8); err != nil {
		return 0, fmt.Errorf("bootstrap prepared frame register 19: %w", err)
	}
	s.x[8] = bootstrapPreparedTable
	first, err := s.load(s.x[8]+0x630, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap prepared first state: %w", err)
	}
	s.x[9] = first
	s.x[10] = bootstrapPreparedTable
	second, err := s.load(s.x[10]+0xb40, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap prepared second state: %w", err)
	}
	s.x[11] = second
	s.x[12] = s.sp + 8
	s.x[9] = s.x[11] - s.x[9]
	s.x[9] += s.x[12]
	s.x[11] = 0x73d046208eaed1a6
	s.x[9] ^= s.x[11]
	s.x[11] = 0x0a5cd4cf
	s.x[9] *= s.x[11]
	if err := s.storeValue(s.x[8]+0x630, s.x[9], 8); err != nil {
		return 0, fmt.Errorf("bootstrap prepared first state output: %w", err)
	}
	if err := s.storeValue(s.x[10]+0xb40, s.x[9], 8); err != nil {
		return 0, fmt.Errorf("bootstrap prepared second state output: %w", err)
	}
	s.x[8] = uint64((uint32(0xe1f2b71e) | uint32(s.x[0]<<1)))
	s.x[9] = 0x8f06a471
	s.x[8] = uint64(uint32(s.x[8]) - uint32(s.x[0]))
	s.x[8] = uint64(uint32(s.x[8]) + uint32(s.x[9]))
	s.x[9] = 0x97cb89ed
	s.x[8] = uint64(uint32(s.x[8]) ^ uint32(s.x[9]))
	s.x[9] = 0xba209ef7
	word, err := s.load(s.x[0], 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap prepared input word: %w", err)
	}
	s.x[10] = uint64(uint32(word))
	s.x[12] = uint64(uint32(s.x[10]) - uint32(s.x[8])*uint32(s.x[9]))
	value, err := s.load(s.x[0]+8, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap prepared input base: %w", err)
	}
	s.x[8] = value
	s.x[9] = 0x11b7
	s.x[10] = 0x163
	s.x[9] = uint64(uint32(s.x[12]) + uint32(s.x[9]))
	s.x[9] = uint64(uint32(s.x[9]) | uint32(s.x[10]))
	word, err = s.load(s.x[8]+4, 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap prepared comparison word: %w", err)
	}
	s.x[10] = uint64(uint32(word))
	s.x[11] = 0x842ad674
	s.x[11] = uint64(uint32(s.x[9]) ^ uint32(s.x[11]))
	s.x[10] = s.addWithCarry(uint64(uint32(s.x[11])), uint64(uint32(s.x[10])), false, 32)
	s.x[11] = bootstrapPreparedBytesA
	byteValue, err := s.load(s.x[11]+0x76a, 1)
	if err != nil {
		return 0, fmt.Errorf("bootstrap prepared sign byte: %w", err)
	}
	s.x[11] = uint64(byte(byteValue)) ^ 0x5d
	s.x[13] = bootstrapPreparedBytesB + 0x590
	byteValue, err = s.load(s.x[13]+uint64(uint32(s.x[11])), 1)
	if err != nil {
		return 0, fmt.Errorf("bootstrap prepared substituted byte: %w", err)
	}
	s.x[11] = uint64(byte(byteValue))
	if s.condition("mi") {
		s.x[10] = uint64(uint32(0) - uint32(s.x[10]))
	}
	s.x[11] = uint64(uint32(s.x[11]) ^ 0x3a)
	s.x[11] = uint64(uint32(s.x[11]) + 0x27d)
	s.x[13] = bootstrapPreparedTargets + 0x6f0
	value, err = s.load(s.x[13]+uint64(uint32(s.x[11]))*8, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap prepared target: %w", err)
	}
	s.x[11] = value
	return bootstrapPreparedTail, nil
}
