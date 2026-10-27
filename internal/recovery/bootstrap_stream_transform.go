package recovery

import "fmt"

const (
	bootstrapStreamMixStart = uint64(0x100a13b30)
	bootstrapStreamMixTail  = uint64(0x100a13c00)
	bootstrapStreamState    = uint64(0x1009f6df0)
	bootstrapStreamAdvance  = uint64(0x1009f3154)
)

// runBootstrapStreamMix is the arithmetic front of the 126-step bootstrap
// stream loop. It advances the 64-bit mixer and resolves the protected next
// state from the loop's signed dispatch table.
func runBootstrapStreamMix(s *state) (uint64, error) {
	s.traceBlock(bootstrapStreamMixStart)
	if err := s.storeValue(s.sp+0x8c, uint64(uint32(s.x[20])), 4); err != nil {
		return 0, fmt.Errorf("bootstrap stream counter: %w", err)
	}
	s.x[8] = uint64(uint32(s.x[22]) + 0x84549631)
	s.x[11] = uint64(uint32(s.x[8]) & 0x7babbd79)
	value, err := s.load(s.sp+0xd8, 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap stream seed: %w", err)
	}
	s.x[8] = uint64(uint32(value)) ^ s.x[2]
	value, err = s.load(s.sp+0x7c, 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap stream multiplier: %w", err)
	}
	s.x[10] = uint64(uint32(value))
	s.x[15] = uint64(uint32(s.x[8]) * uint32(s.x[10]))
	base, err := s.load(s.sp+0x70, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap stream source base: %w", err)
	}
	s.x[10] = base
	value, err = s.load(base, 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap stream source word: %w", err)
	}
	s.x[17] = uint64(uint32(value))
	s.x[1] = (s.x[17] ^ 0x7643e7bb50f3522b) + 0x89bc1844af0cadd5
	s.x[30] = s.x[17] ^ 0x0f18ce1a03bfc4a0
	s.x[14] = uint64(uint32(s.x[11]) + 0xffffc9a9)
	if err := s.storeValue(s.sp+0x98, s.x[14], 8); err != nil {
		return 0, fmt.Errorf("bootstrap stream saved state: %w", err)
	}
	s.x[14] ^= 0xf0e731e5fc405292
	s.x[17] ^= 0x795b29a1710235b0
	s.x[17] += 0x86a4d65e8efdca50
	s.x[17] = s.x[1] ^ s.x[17]
	s.x[14] = s.x[30] + s.x[14]
	s.x[14] = s.x[17] ^ s.x[14]
	s.x[14] += s.x[2]
	s.x[8] += s.x[14] * s.x[15]
	s.x[14] = s.x[21] & (s.x[8] << 1)

	s.traceBlock(bootstrapStreamMixTail)
	s.x[8] -= s.x[14]
	s.x[1] = s.x[8] + s.x[7]
	value, err = s.load(s.sp+0x6c, 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap stream selector: %w", err)
	}
	s.x[8] = uint64(uint32(s.x[22]) ^ uint32(value))
	dispatchTable, err := s.load(s.sp+0xa8, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap stream dispatch table: %w", err)
	}
	s.x[10] = dispatchTable
	next, err := loadBootstrapPointer(s, dispatchTable, uint32(s.x[8]))
	if err != nil {
		return 0, fmt.Errorf("bootstrap stream dispatch: %w", err)
	}
	s.x[8] = next
	s.x[22] = 0x4c9c24d6
	s.x[17] = 0x207f
	return next, nil
}

func runBootstrapStreamState(s *state) (uint64, error) {
	s.traceBlock(bootstrapStreamState)
	s.x[8] = uint64(uint32(s.x[11]) - 0x95d)
	dispatchTable, err := s.load(s.sp+0xa8, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap stream-state dispatch table: %w", err)
	}
	s.x[10] = dispatchTable
	next, err := loadBootstrapPointer(s, dispatchTable, uint32(s.x[8]))
	if err != nil {
		return 0, fmt.Errorf("bootstrap stream-state dispatch: %w", err)
	}
	s.x[20] = next
	base, err := s.load(s.sp+0xa0, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap stream-state base: %w", err)
	}
	s.x[8] = base
	value, err := s.load(base+8, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap stream-state value: %w", err)
	}
	s.x[14] = value
	s.x[8] = 0x3fd1a9c9
	s.x[17] = uint64(uint32(s.x[11]) ^ 0xa0f6)
	s.x[17] = uint64(uint32(s.x[17]) * 0x8e)
	s.x[10] = 0xffffac56
	s.x[11] = uint64(uint32(s.x[11]) + uint32(s.x[10]))
	s.x[10] = 0x224ecac9
	if err := s.storeValue(s.sp+0x98, s.x[17], 8); err != nil {
		return 0, fmt.Errorf("bootstrap stream-state output: %w", err)
	}
	s.x[30] = uint64(uint32(s.x[17]) ^ uint32(s.x[10]))
	s.x[17] = 0x1b85
	return next, nil
}

func runBootstrapStreamAdvance(s *state) (uint64, error) {
	s.traceBlock(bootstrapStreamAdvance)
	s.compare32(uint32(s.x[27]), 0x2e0eb4b5)
	firstGate := uint32(boolValue(s.condition("gt")))
	s.x[8] = uint64(firstGate)
	state, err := s.load(s.sp+0x98, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap stream-advance state: %w", err)
	}
	s.x[20] = state
	s.x[10] = uint64(uint32(s.x[20]) + 0xffffa214)
	s.x[22] = uint64(uint32(s.x[10]) | 0x4098)
	s.x[10] = s.x[23] ^ (s.x[1] >> 32)
	output, err := s.load(s.sp+0x60, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap stream-advance output pointer: %w", err)
	}
	s.x[14] = output
	value, err := s.load(output, 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap stream-advance output word: %w", err)
	}
	s.x[11] = uint64(uint32(value)) ^ s.x[2]
	s.x[10] += s.x[11]
	s.x[11] = 0x05a522f0
	s.x[11] &= s.x[10] << 1
	s.x[10] -= s.x[11]
	s.x[15] = uint64(uint32(s.x[20]) + 0xffffc2e2)
	s.x[11] = uint64(uint32(s.x[15]) ^ uint32(s.x[22]))
	s.x[11] = uint64(uint32(s.x[11]) ^ 0x209c5209)
	s.x[10] += 0x84a96e6a02d29178
	s.x[11] = uint64(uint32(s.x[11]) ^ uint32(s.x[10]))
	wordOutput, err := s.load(s.sp+0x58, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap stream-advance word pointer: %w", err)
	}
	s.x[17] = wordOutput
	if err := s.storeValue(wordOutput, uint64(uint32(s.x[11])), 4); err != nil {
		return 0, fmt.Errorf("bootstrap stream-advance first output: %w", err)
	}
	s.x[10] >>= 32
	s.x[10] = uint64(uint32(s.x[10]) ^ 0xa6e7cd51)
	if err := s.storeValue(output, uint64(uint32(s.x[10])), 4); err != nil {
		return 0, fmt.Errorf("bootstrap stream-advance second output: %w", err)
	}
	counter, err := s.load(s.sp+0x8c, 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap stream-advance counter: %w", err)
	}
	s.x[11] = uint64(uint32(counter))
	s.x[10] = uint64(uint32(s.x[11]) + 1)
	s.x[11] = uint64(uint32(s.x[11]) + 0xfc08e30c)
	s.compare32(uint32(s.x[11]), 0xd1f14b4a)
	secondGate := uint32(boolValue(s.condition("lt")))
	s.x[14] = uint64(firstGate ^ secondGate)
	limit, err := s.load(s.sp+0x80, 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap stream-advance limit: %w", err)
	}
	s.x[17] = uint64(uint32(limit))
	s.compare32(uint32(s.x[11]), uint32(s.x[17]))
	thirdGate := uint32(boolValue(s.condition("lt")))
	s.x[11] = uint64(thirdGate)
	s.compare32(uint32(s.x[14]), 0)
	dispatchGate := thirdGate
	if s.condition("ne") {
		dispatchGate = firstGate
	}
	s.x[8] = uint64(dispatchGate * 0x17cf)
	s.x[8] = uint64(uint32(s.x[20]) ^ uint32(s.x[8]))
	dispatchTable, err := s.load(s.sp+0xa8, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap stream-advance dispatch table: %w", err)
	}
	s.x[11] = dispatchTable
	next, err := loadBootstrapPointer(s, dispatchTable, uint32(s.x[8]))
	if err != nil {
		return 0, fmt.Errorf("bootstrap stream-advance dispatch: %w", err)
	}
	s.x[8] = next
	s.x[20] = s.x[10]
	s.x[17] = 0xc98e
	return next, nil
}
