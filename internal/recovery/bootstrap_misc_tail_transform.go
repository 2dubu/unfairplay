package recovery

import "fmt"

const (
	bootstrapPointerTailA     = uint64(0x100022cb0)
	bootstrapGuardReturnA     = uint64(0x10006b4e0)
	bootstrapGuardReturnATail = uint64(0x10006b4f8)
	bootstrapGuardMismatchA   = uint64(0x10006b50c)
	bootstrapStoredDispatch   = uint64(0x100127c84)
	bootstrapPackedSelector   = uint64(0x1001ab594)
	bootstrapGuardReturnB     = uint64(0x1001bc008)
	bootstrapGuardReturnBTail = uint64(0x1001bc020)
	bootstrapGuardMismatchB   = uint64(0x1001bc038)
	bootstrapSignedSelector   = uint64(0x100362744)
	bootstrapStackPivot       = uint64(0x100c06e00)
	bootstrapStackPairSelect  = uint64(0x100caeaac)
	bootstrapZeroStore        = uint64(0x1003c7db4)
)

func runBootstrapPointerTailA(s *state) (uint64, error) {
	s.traceBlock(bootstrapPointerTailA)
	next, err := loadBootstrapPointer(s, s.x[19], uint32(s.x[20]))
	if err != nil {
		return 0, fmt.Errorf("bootstrap pointer tail A: %w", err)
	}
	s.x[8] = next
	s.x[17] = 0xb08d
	return next, nil
}

func runBootstrapGuardReturnA(s *state) (uint64, error) {
	return runBootstrapCompactGuardReturn(s, bootstrapGuardReturnA, bootstrapGuardReturnATail, bootstrapGuardMismatchA, 0x28, [][3]uint64{{29, 30, 0x50}, {20, 19, 0x40}, {22, 21, 0x30}})
}

func runBootstrapGuardReturnB(s *state) (uint64, error) {
	return runBootstrapCompactGuardReturn(s, bootstrapGuardReturnB, bootstrapGuardReturnBTail, bootstrapGuardMismatchB, 0x18, [][3]uint64{{29, 30, 0x50}, {20, 19, 0x40}, {22, 21, 0x30}, {24, 23, 0x20}})
}

func runBootstrapCompactGuardReturn(s *state, start, tail, mismatch, savedOffset uint64, pairs [][3]uint64) (uint64, error) {
	s.traceBlock(start)
	saved, err := s.load(s.sp+savedOffset, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap guard %#x saved value: %w", start, err)
	}
	s.x[8] = saved
	guardPointer, err := s.load(0x10106d000+0x408, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap guard %#x pointer: %w", start, err)
	}
	s.x[9] = guardPointer
	guard, err := s.load(guardPointer, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap guard %#x value: %w", start, err)
	}
	s.x[9] = guard
	s.compare64(guard, saved)
	if s.condition("ne") {
		return mismatch, nil
	}
	s.traceBlock(tail)
	for _, pair := range pairs {
		firstRegister, secondRegister, offset := int(pair[0]), int(pair[1]), pair[2]
		first, err := s.load(s.sp+offset, 8)
		if err != nil {
			return 0, fmt.Errorf("bootstrap guard %#x register %d: %w", start, firstRegister, err)
		}
		second, err := s.load(s.sp+offset+8, 8)
		if err != nil {
			return 0, fmt.Errorf("bootstrap guard %#x register %d: %w", start, secondRegister, err)
		}
		s.x[firstRegister], s.x[secondRegister] = first, second
	}
	s.sp += 0x60
	return s.x[30], nil
}

func runBootstrapStoredDispatch(s *state) (uint64, error) {
	s.traceBlock(bootstrapStoredDispatch)
	if err := s.storeValue(s.x[19]+4, uint64(uint32(s.x[8])), 4); err != nil {
		return 0, fmt.Errorf("bootstrap stored-dispatch output: %w", err)
	}
	next, err := loadBootstrapPointer(s, s.x[20], uint32(s.x[11]))
	if err != nil {
		return 0, fmt.Errorf("bootstrap stored-dispatch target: %w", err)
	}
	s.x[8] = next
	s.x[17] = 0x8d19
	return next, nil
}

func runBootstrapPackedSelector(s *state) (uint64, error) {
	s.traceBlock(bootstrapPackedSelector)
	s.x[8] = uint64(uint32(s.x[23]) + uint32(s.x[22]))
	s.x[13] = uint64(uint32(s.x[8]) & 0x1f)
	s.x[8] = uint64(uint32(s.x[8]) >> 5)
	s.x[9] = 0x00a2d3e2
	s.x[9] = uint64(uint32(s.x[8]) ^ uint32(s.x[9]))
	s.x[10] = uint64(uint32(s.x[21]) << 1)
	s.x[12] = 0x095ae2a9
	s.x[14] = 0xf6a4f5ff
	s.x[12] = uint64(uint32(s.x[11]) + uint32(s.x[12]))
	s.x[12] = uint64(uint32(s.x[12]) & uint32(s.x[14]))
	s.compare32(uint32(s.x[13]), 0x14)
	gate := uint32(boolValue(s.condition("eq")))
	s.x[14] = uint64(gate)
	s.x[15] = 0xa7
	s.x[14] = uint64(gate * uint32(s.x[15]))
	s.x[11] = uint64(uint32(s.x[11]) ^ uint32(s.x[14]))
	next, err := loadBootstrapPointer(s, s.x[20], uint32(s.x[11]))
	if err != nil {
		return 0, fmt.Errorf("bootstrap packed-selector target: %w", err)
	}
	s.x[11] = next
	s.x[17] = 0x150f
	return next, nil
}

func runBootstrapSignedSelector(s *state) (uint64, error) {
	s.traceBlock(bootstrapSignedSelector)
	s.compare32(uint32(s.x[21]), 0)
	gate := uint32(boolValue(s.condition("gt")))
	s.x[8] = uint64(gate)
	s.x[9] = 0xffffb623
	s.x[10] = 0x5020
	s.x[9] = uint64(uint32(s.x[24]) + uint32(s.x[9]))
	s.x[11] = uint64(uint32(s.x[9]) | uint32(s.x[10]))
	s.x[9] = 0xffffa7b7
	s.x[9] = uint64(uint32(s.x[11]) + uint32(s.x[9]))
	s.x[8] = uint64(gate * uint32(s.x[9]))
	s.x[8] = uint64(uint32(s.x[8]) ^ uint32(s.x[24]))
	next, err := loadBootstrapPointer(s, s.x[20], uint32(s.x[8]))
	if err != nil {
		return 0, fmt.Errorf("bootstrap signed-selector target: %w", err)
	}
	s.x[8] = next
	s.x[17] = 0x8d19
	return next, nil
}

func runBootstrapStackPivot(s *state) (uint64, error) {
	s.traceBlock(bootstrapStackPivot)
	s.x[8] = s.sp
	value, err := s.load(s.sp+0xc0, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap stack-pivot offset: %w", err)
	}
	s.x[9] = value
	s.x[8] += value
	s.sp = s.x[8]
	s.x[9] = uint64(uint32(s.x[24]) >> 1)
	word, err := s.load(s.sp+0x2c, 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap stack-pivot word: %w", err)
	}
	s.x[8] = uint64(uint32(word))
	if err := s.storeValue(s.sp+0x4c, uint64(uint32(s.x[9])), 4); err != nil {
		return 0, fmt.Errorf("bootstrap stack-pivot first saved word: %w", err)
	}
	if err := s.storeValue(s.sp+0x50, uint64(uint32(s.x[24])), 4); err != nil {
		return 0, fmt.Errorf("bootstrap stack-pivot second saved word: %w", err)
	}
	s.x[8] = uint64(uint32(s.x[9]) & uint32(s.x[8]))
	s.x[9] = uint64(uint32(s.x[8]) - uint32(s.x[21]))
	s.x[8] = uint64(uint32(s.x[21]) - uint32(s.x[8]))
	s.x[8] = uint64(uint32(s.x[9]) | uint32(s.x[8]))
	s.x[8] = uint64(^uint32(s.x[8]))
	s.x[8] = uint64(uint32(s.x[8]) >> 31)
	s.x[9] = 0xc980b09e
	s.x[9] = uint64(uint32(s.x[26]) ^ uint32(s.x[9]))
	s.x[9] = uint64(uint32(s.x[9]) + uint32(s.x[27]))
	s.x[8] = uint64(uint32(s.x[8]) * uint32(s.x[9]))
	s.x[9] = 0xc981006e
	s.x[9] = uint64(uint32(s.x[26]) + uint32(s.x[9]))
	s.x[8] = uint64(uint32(s.x[8]) ^ uint32(s.x[9]))
	next, err := loadBootstrapPointer(s, s.x[28], uint32(s.x[8]))
	if err != nil {
		return 0, fmt.Errorf("bootstrap stack-pivot target: %w", err)
	}
	s.x[8] = next
	s.x[28] = 0xba209ef7
	s.x[17] = 0x7fb9
	return next, nil
}

func runBootstrapStackPairSelect(s *state) (uint64, error) {
	s.traceBlock(bootstrapStackPairSelect)
	s.x[8] = 0x767fc75d
	s.x[8] = uint64(uint32(s.x[27]) + uint32(s.x[8]))
	s.x[9] = 0xbfff777f
	s.x[26] = uint64(uint32(s.x[8]) & uint32(s.x[9]))
	s.x[8] = 0xc98088d0
	s.x[8] = uint64(uint32(s.x[26]) + uint32(s.x[8]))
	first, err := s.load(s.sp+0x4c, 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap stack-pair first word: %w", err)
	}
	second, err := s.load(s.sp+0x50, 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap stack-pair second word: %w", err)
	}
	s.x[24], s.x[9] = uint64(uint32(first)), uint64(uint32(second))
	s.compare32(uint32(s.x[9]), uint32(s.x[8]))
	gate := uint32(boolValue(s.condition("hi")))
	s.x[8] = uint64(gate)
	s.x[9] = 0x9e30
	s.x[8] = uint64(gate * uint32(s.x[9]))
	s.x[8] = uint64(uint32(s.x[8]) ^ uint32(s.x[27]))
	value, err := s.load(s.sp+0x68, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap stack-pair table: %w", err)
	}
	s.x[28] = value
	next, err := loadBootstrapPointer(s, value, uint32(s.x[8]))
	if err != nil {
		return 0, fmt.Errorf("bootstrap stack-pair target: %w", err)
	}
	s.x[8] = next
	s.x[12] = 0xa6
	s.x[13] = 0x0a5cd4cf
	s.x[17] = 0xe334
	return next, nil
}

func runBootstrapZeroStore(s *state) (uint64, error) {
	s.traceBlock(bootstrapZeroStore)
	value, err := s.load(s.x[8], 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap zero-store input: %w", err)
	}
	s.x[17] = value
	s.x[0] = s.x[11]*s.x[12] + s.x[13]
	s.x[17] += 0x18
	if err := s.storeValue(s.x[17]+s.x[0], 0, 4); err != nil {
		return 0, fmt.Errorf("bootstrap zero-store output: %w", err)
	}
	s.x[11] = s.x[14] + s.x[11]
	s.compare64(s.x[11], s.x[15])
	gate := uint32(boolValue(s.condition("eq")))
	s.x[17] = uint64(gate * uint32(s.x[16]))
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[10]))
	table, err := s.load(s.sp+0xdf0, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap zero-store table: %w", err)
	}
	s.x[0] = table
	next, err := loadBootstrapPointer(s, table, uint32(s.x[17]))
	if err != nil {
		return 0, fmt.Errorf("bootstrap zero-store target: %w", err)
	}
	s.x[0] = next
	s.x[17] = 0x6e8c
	return next, nil
}
