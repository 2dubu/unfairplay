package recovery

import "fmt"

const (
	accountEntryByteGeneratorPass = uint64(0x100b1ddc0)
	accountEntryXORCopyPass       = uint64(0x1004bb284)
	accountEntryByteRemapPass     = uint64(0x100131cd8)
	accountEntryTableWriterPass   = uint64(0x10007e214)
	accountEntryWordMixerPass     = uint64(0x10007ba1c)
	accountEntryByteCopyPass      = uint64(0x1001f39b8)
	accountEntryFourWordPass      = uint64(0x1002376d4)
	accountEntryZeroFillPass      = uint64(0x10092144c)
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
	case accountEntryByteCopyPass:
		next, err := runAccountEntryByteCopy(s)
		return next, true, err
	case accountEntryFourWordPass:
		next, err := runAccountEntryFourWordPass(s)
		return next, true, err
	case accountEntryZeroFillPass:
		next, err := runAccountEntryZeroFill(s)
		return next, true, err
	default:
		return pc, false, nil
	}
}

func runAccountEntryZeroFill(s *state) (uint64, error) {
	const maximumBytes = 4096
	for index := 0; index < maximumBytes; index++ {
		next, err := runAccountEntryZeroFillStep(s)
		if err != nil || next != accountEntryZeroFillPass {
			return next, err
		}
	}
	return 0, fmt.Errorf("account-entry zero fill exceeded %d bytes", maximumBytes)
}

func runAccountEntryZeroFillStep(s *state) (uint64, error) {
	s.traceBlock(accountEntryZeroFillPass)
	s.x[11] = s.x[16] + s.x[10]
	if err := s.storeValue(s.x[11], 0, 1); err != nil {
		return 0, fmt.Errorf("account-entry zero-fill output: %w", err)
	}
	s.x[11] = s.x[10] - 0x48f
	s.x[13] = s.x[8] ^ s.x[14]
	s.compare64(s.x[11], s.x[13])
	dispatchIndex := uint32(boolValue(s.condition("eq")))*uint32(s.x[15]) ^ uint32(s.x[9])
	next, err := loadBootstrapPointer(s, s.x[30], dispatchIndex)
	if err != nil {
		return 0, fmt.Errorf("account-entry zero-fill dispatch: %w", err)
	}
	s.x[10]++
	s.x[11] = next
	s.x[17] = 0x345f
	return next, nil
}

func runAccountEntryByteCopy(s *state) (uint64, error) {
	const maximumBytes = 4096
	for index := 0; index < maximumBytes; index++ {
		next, err := runAccountEntryByteCopyStep(s)
		if err != nil || next != accountEntryByteCopyPass {
			return next, err
		}
	}
	return 0, fmt.Errorf("account-entry byte copy exceeded %d bytes", maximumBytes)
}

func runAccountEntryByteCopyStep(s *state) (uint64, error) {
	s.traceBlock(accountEntryByteCopyPass)
	s.compare32(uint32(s.x[15]), uint32(s.x[14]))
	gate := uint32(boolValue(s.condition("lo")))
	s.x[10] = uint64(uint32(s.x[12]) + 0xb66ead0e)
	s.x[16] = 0x4990ff67
	s.x[10] = uint64(uint32(s.x[10]) & uint32(s.x[16]))
	s.x[16] = uint64(uint32(s.x[13]))
	value, err := s.load(s.x[9]+s.x[16], 1)
	if err != nil {
		return 0, fmt.Errorf("account-entry byte copy input: %w", err)
	}
	s.x[17] = uint64(byte(value))
	s.x[16] = s.x[8] + s.x[16]
	if err := s.storeValue(s.x[16], uint64(byte(s.x[17])), 1); err != nil {
		return 0, fmt.Errorf("account-entry byte copy output: %w", err)
	}
	s.x[16] = uint64(uint32(s.x[13]) + 0xffffb09b)
	s.x[16] = uint64(uint32(s.x[16]) + uint32(s.x[10]))
	s.x[17] = 0xe21ab501
	s.x[0] = 0x1de54aff
	s.compare32(uint32(s.x[16]), uint32(s.x[0]))
	selector := gate ^ uint32(boolValue(s.condition("hi")))
	s.x[13] = uint64(uint32(s.x[13]) + uint32(s.x[17]))
	s.compare32(uint32(s.x[13]), uint32(s.x[15]))
	secondBound := uint32(boolValue(s.condition("lo")))
	s.compare32(selector, 0)
	dispatchBit := secondBound
	if s.condition("ne") {
		dispatchBit = gate
	}
	dispatchIndex := dispatchBit*0x172c ^ uint32(s.x[12])
	next, err := loadBootstrapPointer(s, s.x[20], dispatchIndex)
	if err != nil {
		return 0, fmt.Errorf("account-entry byte copy dispatch: %w", err)
	}
	s.x[11] = next
	s.x[13] = s.x[16]
	s.x[17] = 0x63d5
	return next, nil
}

func runAccountEntryFourWordPass(s *state) (uint64, error) {
	const maximumWords = 4096
	for index := 0; index < maximumWords; index++ {
		next, err := runAccountEntryFourWordStep(s)
		if err != nil || next != accountEntryFourWordPass {
			return next, err
		}
	}
	return 0, fmt.Errorf("account-entry four-word pass exceeded %d words", maximumWords)
}

func runAccountEntryFourWordStep(s *state) (uint64, error) {
	s.traceBlock(accountEntryFourWordPass)
	s.x[9] = 0xffff5fad
	s.x[16] = uint64(uint32(s.x[13]) + uint32(s.x[9]))
	s.x[9] = 0x410
	s.x[0] = uint64(uint32(s.x[16]) | uint32(s.x[9]))
	s.x[16] = s.x[1] << 2
	offsets := [...]uint64{
		0xffffffff8cdedda8,
		0xffffffff8cdedd94,
		0xffffffff8cdedd7c,
	}
	first, err := s.load(s.x[8]+s.x[16]+offsets[0], 4)
	if err != nil {
		return 0, fmt.Errorf("account-entry four-word input 0: %w", err)
	}
	s.x[17] = uint64(uint32(first))
	second, err := s.load(s.x[8]+s.x[16]+offsets[1], 4)
	if err != nil {
		return 0, fmt.Errorf("account-entry four-word input 1: %w", err)
	}
	s.x[2] = uint64(uint32(second))
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[2]))
	third, err := s.load(s.x[8]+s.x[16]+offsets[2], 4)
	if err != nil {
		return 0, fmt.Errorf("account-entry four-word input 2: %w", err)
	}
	s.x[2] = uint64(uint32(third))
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[2]))
	s.x[2] = s.x[1] + s.x[0]
	s.x[9] = 0xffffffff8cde7fb0
	s.x[2] = s.x[9] + (s.x[2] << 2)
	fourth, err := s.load(s.x[8]+s.x[2], 4)
	if err != nil {
		return 0, fmt.Errorf("account-entry four-word input 3: %w", err)
	}
	s.x[2] = uint64(uint32(fourth))
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[2]))
	s.x[17] = uint64(uint32(s.x[12]) ^ uint32(rotateRight(s.x[17], 31, 32)))
	s.x[16] += s.x[14]
	if err := s.storeValue(s.x[8]+s.x[16], uint64(uint32(s.x[17])), 4); err != nil {
		return 0, fmt.Errorf("account-entry four-word output: %w", err)
	}
	s.x[1]++
	s.x[9] = 0x1cc848e3
	s.compare64(s.x[1], s.x[9])
	dispatchIndex := uint32(boolValue(s.condition("eq")))*0x7615 ^ uint32(s.x[13])
	next, err := loadBootstrapPointer(s, s.x[20], dispatchIndex)
	if err != nil {
		return 0, fmt.Errorf("account-entry four-word dispatch: %w", err)
	}
	s.x[2] = next
	s.x[16] = 0x1cc84893
	s.x[17] = 0xe949
	return next, nil
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
