package recovery

import "fmt"

const (
	accountEntryByteGeneratorPass = uint64(0x100b1ddc0)
	accountEntryXORCopyPass       = uint64(0x1004bb284)
	accountEntryByteRemapPass     = uint64(0x100131cd8)
	accountEntryTableWriterPass   = uint64(0x10007e214)
	accountEntryWordMixerPass     = uint64(0x10007ba1c)
)

func runAccountEntrySemanticBlock(s *state, pc uint64) (uint64, bool, error) {
	switch pc {
	case accountEntryByteGeneratorPass:
		next, err := runAccountEntryByteGenerator(s)
		return next, true, err
	case accountEntryXORCopyPass:
		next, err := runAccountEntryXORCopy(s)
		return next, true, err
	case accountEntryByteRemapPass:
		next, err := runAccountEntryByteRemap(s)
		return next, true, err
	case accountEntryTableWriterPass:
		next, err := runAccountEntryTableWriter(s)
		return next, true, err
	case accountEntryWordMixerPass:
		next, err := runAccountEntryWordMixer(s)
		return next, true, err
	default:
		return pc, false, nil
	}
}

func runAccountEntryWordMixer(s *state) (uint64, error) {
	const maximumWords = 4096
	for index := 0; index < maximumWords; index++ {
		next, err := runAccountEntryWordMixerStep(s)
		if err != nil || next != accountEntryWordMixerPass {
			return next, err
		}
	}
	return 0, fmt.Errorf("account-entry word mixer exceeded %d words", maximumWords)
}

func runAccountEntryWordMixerStep(s *state) (uint64, error) {
	s.traceBlock(accountEntryWordMixerPass)
	s.x[9] = uint64(uint32(s.x[15]) + 0xd458908a)
	s.x[11] = s.x[20] << 2
	addresses := [...]uint64{
		0xffffffff5160f1c8,
		0xffffffff5160f1b4,
		0xffffffff5160f19c,
		0xffffffff5160f194,
	}
	registers := [...]int{12, 13, 16, 17}
	for index, offset := range addresses {
		value, err := s.load(s.x[10]+s.x[11]+offset, 4)
		if err != nil {
			return 0, fmt.Errorf("account-entry word mixer input %d: %w", index, err)
		}
		s.x[registers[index]] = uint64(uint32(value))
	}
	s.x[12] = uint64(uint32(s.x[13]) ^ uint32(s.x[12]))
	s.x[13] = uint64(uint32(s.x[16]) ^ uint32(s.x[17]))
	s.x[12] = uint64(uint32(s.x[12]) ^ uint32(s.x[13]))
	s.x[12] = uint64(uint32(s.x[9]) ^ uint32(rotateRight(s.x[12], 31, 32)))
	s.x[13] = 0xce752fa0
	s.x[12] = uint64(uint32(s.x[12]) ^ uint32(s.x[13]))
	s.x[11] += s.x[6]
	if err := s.storeValue(s.x[10]+s.x[11], uint64(uint32(s.x[12])), 4); err != nil {
		return 0, fmt.Errorf("account-entry word mixer output: %w", err)
	}
	s.x[20]++
	s.x[11] = s.x[20] - s.x[15]
	s.x[12] = s.x[15] - s.x[20]
	s.x[11] |= s.x[12]
	s.x[11] >>= 63
	s.x[11] = uint64(uint32(s.x[11]) * 0x72)
	s.x[12] = uint64(uint32(s.x[15]) + 0xd4589955)
	dispatchIndex := uint32(s.x[11]) ^ uint32(s.x[12])
	next, err := loadBootstrapPointer(s, s.x[21], dispatchIndex)
	if err != nil {
		return 0, fmt.Errorf("account-entry word mixer dispatch: %w", err)
	}
	limit, err := s.load(s.sp+0x7c, 4)
	if err != nil {
		return 0, fmt.Errorf("account-entry word mixer limit: %w", err)
	}
	s.x[4] = next
	s.x[7] = s.x[0]
	s.x[11] = s.x[3]
	s.x[12] = s.x[26]
	s.x[13] = s.x[22]
	s.x[16] = 0x2ba7c38b
	s.x[17] = 0x3601
	s.x[19] = uint64(uint32(limit))
	return next, nil
}

func runAccountEntryTableWriter(s *state) (uint64, error) {
	const maximumBytes = 4096
	for index := 0; index < maximumBytes; index++ {
		next, err := runAccountEntryTableWriterStep(s)
		if err != nil || next != accountEntryTableWriterPass {
			return next, err
		}
	}
	return 0, fmt.Errorf("account-entry table writer exceeded %d bytes", maximumBytes)
}

func runAccountEntryTableWriterStep(s *state) (uint64, error) {
	s.traceBlock(accountEntryTableWriterPass)
	s.x[11] = uint64(uint32(s.x[15]) + uint32(s.x[19]))
	s.x[13] = uint64(uint32(s.x[11]) | uint32(s.x[20]))
	value, err := s.load(s.x[2]+uint64(uint32(s.x[9])), 1)
	if err != nil {
		return 0, fmt.Errorf("account-entry table writer input: %w", err)
	}
	s.x[11] = uint64(byte(value))
	s.x[25] = uint64(uint32(s.x[9]) + 1)
	s.x[9] = s.x[16] + s.x[22]
	s.x[9] += s.x[13]
	s.x[12] = uint64(uint32(s.x[23]) & uint32(s.x[11]<<1))
	s.x[11] = uint64(uint32(s.x[11]) - uint32(s.x[12]))
	s.x[11] = uint64(uint32(s.x[11]) + 0x1b)
	if err := s.storeValue(s.x[24]+s.x[9], uint64(byte(s.x[11])), 1); err != nil {
		return 0, fmt.Errorf("account-entry table writer output: %w", err)
	}
	s.x[16]++
	s.compare64(s.x[16], s.x[26])
	dispatchIndex := uint32(boolValue(s.condition("eq")))*uint32(s.x[27]) ^ uint32(s.x[15])
	next, err := loadBootstrapPointer(s, s.x[21], dispatchIndex)
	if err != nil {
		return 0, fmt.Errorf("account-entry table writer dispatch: %w", err)
	}
	s.x[1] = s.x[24]
	s.x[4] = 0x12d6a942
	s.x[9] = s.x[25]
	s.x[11] = next
	s.x[17] = 0x71c7
	return next, nil
}

func runAccountEntryXORCopy(s *state) (uint64, error) {
	const maximumBytes = 4096
	for index := 0; index < maximumBytes; index++ {
		next, err := runAccountEntryXORCopyStep(s)
		if err != nil || next != accountEntryXORCopyPass {
			return next, err
		}
	}
	return 0, fmt.Errorf("account-entry XOR copy exceeded %d bytes", maximumBytes)
}

func runAccountEntryXORCopyStep(s *state) (uint64, error) {
	s.traceBlock(accountEntryXORCopyPass)
	s.x[19] = uint64(uint32(s.x[12]) | uint32(s.x[13]))
	s.x[8] = s.x[11] + s.x[14]
	s.x[11] = s.addWithCarry(s.x[8], s.x[19], false, 64)
	value, err := s.load(s.x[9]+s.x[11], 1)
	if err != nil {
		return 0, fmt.Errorf("account-entry XOR copy input: %w", err)
	}
	s.x[8] = uint64(uint32(value) ^ uint32(s.x[15]))
	if err := s.storeValue(s.x[0]+s.x[11], uint64(byte(s.x[8])), 1); err != nil {
		return 0, fmt.Errorf("account-entry XOR copy output: %w", err)
	}
	last := uint32(boolValue(s.condition("eq")))
	dispatchIndex := last*uint32(s.x[16]) ^ uint32(s.x[10])
	next, err := loadBootstrapPointer(s, s.x[20], dispatchIndex)
	if err != nil {
		return 0, fmt.Errorf("account-entry XOR copy dispatch: %w", err)
	}
	s.x[1] = next
	s.x[8] = 1
	s.x[17] = 0x56cc
	return next, nil
}

func runAccountEntryByteRemap(s *state) (uint64, error) {
	const maximumBytes = 4096
	for index := 0; index < maximumBytes; index++ {
		next, err := runAccountEntryByteRemapStep(s)
		if err != nil || next != accountEntryByteRemapPass {
			return next, err
		}
	}
	return 0, fmt.Errorf("account-entry byte remap exceeded %d bytes", maximumBytes)
}

func runAccountEntryByteRemapStep(s *state) (uint64, error) {
	s.traceBlock(accountEntryByteRemapPass)
	s.x[19] = uint64(uint32(s.x[14]) & uint32(s.x[0]))
	value, err := s.load(s.x[27]+uint64(uint32(s.x[11])), 1)
	if err != nil {
		return 0, fmt.Errorf("account-entry byte remap input: %w", err)
	}
	s.x[8] = uint64(byte(value))
	s.x[11] = uint64(uint32(s.x[11]) + 1)
	s.x[15] = s.x[13] + s.x[1]
	s.x[16] = uint64(uint32(s.x[19]) - 0x1a)
	s.x[16] = uint64(uint32(s.x[16]) & uint32(s.x[8]<<1))
	s.x[8] = uint64(uint32(s.x[8]) - uint32(s.x[16]))
	s.x[8] = uint64(uint32(s.x[8]) + 0x1b)
	if err := s.storeValue(s.x[12]+s.x[15], uint64(byte(s.x[8])), 1); err != nil {
		return 0, fmt.Errorf("account-entry byte remap output: %w", err)
	}
	s.x[13]++
	s.compare64(s.x[13], s.x[2])
	dispatchIndex := uint32(boolValue(s.condition("eq")))*uint32(s.x[5]) ^ uint32(s.x[9])
	next, err := s.load(s.x[20]+uint64(dispatchIndex)*8, 8)
	if err != nil {
		return 0, fmt.Errorf("account-entry byte remap dispatch: %w", err)
	}
	s.x[8] = next
	s.x[17] = 0xd90f
	return next, nil
}

func runAccountEntryByteGenerator(s *state) (uint64, error) {
	const maximumBytes = 4096
	for index := 0; index < maximumBytes; index++ {
		next, err := runAccountEntryByteGeneratorStep(s)
		if err != nil || next != accountEntryByteGeneratorPass {
			return next, err
		}
	}
	return 0, fmt.Errorf("account-entry byte generator exceeded %d bytes", maximumBytes)
}

func runAccountEntryByteGeneratorStep(s *state) (uint64, error) {
	if s == nil {
		return 0, fmt.Errorf("account-entry byte-generator state is unavailable")
	}
	s.traceBlock(accountEntryByteGeneratorPass)
	s.x[0] = uint64(uint32(s.x[16]) ^ uint32(s.x[10]))
	s.x[15] = (s.x[0] ^ s.x[19]) + s.x[1]
	if err := s.storeValue(s.sp+0x318, s.x[15], 8); err != nil {
		return 0, fmt.Errorf("account-entry byte-generator position: %w", err)
	}
	s.x[17] = uint64(uint32(s.x[13]) ^ uint32(s.x[20]))
	s.x[1] = uint64(uint32(s.x[13]) ^ uint32(s.x[23]))
	s.x[1] = uint64(uint32(s.x[1]) + uint32(s.x[26]))
	s.x[2] = uint64(uint32(s.x[13]) ^ uint32(s.x[27]))
	s.x[3] = uint64((uint32(s.x[13]) ^ 0xfffffffb) + 0x8f)
	s.x[17] = uint64(uint32(s.x[17]) + uint32(s.x[21]))
	s.x[17] = uint64(uint32(s.x[13]) ^ uint32(s.x[17]))
	s.x[2] = uint64(uint32(s.x[2]) + uint32(s.x[28]))
	s.x[1] = uint64(uint32(s.x[1]) ^ uint32(s.x[2]))
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[1]))
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[3]))
	s.x[17] = uint64(uint32(s.x[17]) ^ 0x30)
	s.x[1] = uint64(uint32(s.x[13]) + 0x87)
	s.x[17] = uint64(uint32(s.x[17]) * uint32(s.x[1]))
	if err := s.storeValue(s.x[6]+s.x[15], uint64(byte(s.x[17])), 1); err != nil {
		return 0, fmt.Errorf("account-entry byte-generator output: %w", err)
	}
	s.x[13]++
	position, err := s.load(s.sp+0x318, 8)
	if err != nil {
		return 0, fmt.Errorf("account-entry byte-generator saved position: %w", err)
	}
	s.x[1] = position
	s.compare64(position, 0)
	dispatchIndex := uint32(boolValue(s.condition("eq")))*uint32(s.x[14]) ^ uint32(s.x[16])
	next, err := loadBootstrapPointer(s, s.x[30], dispatchIndex)
	if err != nil {
		return 0, fmt.Errorf("account-entry byte-generator dispatch: %w", err)
	}
	s.x[15] = next
	s.x[17] = 0x623c
	return next, nil
}
