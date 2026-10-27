package recovery

import "fmt"

const (
	accountByteLimit       = uint64(0x100700040)
	accountSixByteCopy     = uint64(0x10049d198)
	accountBitStore        = uint64(0x100251470)
	accountByteAdvance     = uint64(0x1001bbbe4)
	accountByteAdvanceTail = uint64(0x1001bbc00)
	accountSumSelect       = uint64(0x10014c8e8)
	accountOutputByte      = uint64(0x10011768c)
	accountPointerSelect   = uint64(0x1008fb824)
	accountCallReturn      = uint64(0x10091e620)
	accountConditionalPick = uint64(0x100949754)
)

func runAccountByteLimit(s *state) (uint64, error) {
	s.traceBlock(accountByteLimit)
	s.x[8] = uint64(uint32(s.x[12]) + 0x400)
	s.x[14] = s.x[19] + s.x[8]
	value, err := s.load(s.x[14], 1)
	if err != nil {
		return 0, fmt.Errorf("account byte-limit input: %w", err)
	}
	s.x[14] = uint64(byte(value))
	s.compare32(uint32(s.x[14]), uint32(byte(s.x[10])))
	gate := uint32(boolValue(s.condition("eq")))
	s.x[14] = uint64(gate * uint32(s.x[22]))
	s.x[14] = uint64(uint32(s.x[14]) ^ uint32(s.x[11]))
	table, err := s.load(s.sp+0xdf0, 8)
	if err != nil {
		return 0, fmt.Errorf("account byte-limit table: %w", err)
	}
	s.x[15] = table
	next, err := loadBootstrapPointer(s, table, uint32(s.x[14]))
	if err != nil {
		return 0, fmt.Errorf("account byte-limit target: %w", err)
	}
	s.x[14] = next
	s.x[12] = uint64(uint32(s.x[12]) + 1)
	s.x[17] = 0x279a
	return next, nil
}

func runAccountSixByteCopy(s *state) (uint64, error) {
	s.traceBlock(accountSixByteCopy)
	for offset := uint64(1); offset <= 6; offset++ {
		s.x[12] = s.x[19] - offset
		value, err := s.load(s.x[24]+s.x[12], 1)
		if err != nil {
			return 0, fmt.Errorf("account six-byte copy input %d: %w", offset, err)
		}
		s.x[13] = uint64(byte(value))
		if err := s.storeValue(s.x[20]+s.x[12], uint64(byte(value)), 1); err != nil {
			return 0, fmt.Errorf("account six-byte copy output %d: %w", offset, err)
		}
	}
	s.x[14] = uint64(uint32(s.x[10]) + uint32(s.x[7]))
	s.x[23] = uint64(uint32(s.x[14]) | uint32(s.x[21]))
	s.x[12] = s.x[19] + s.x[25]
	s.x[12] += s.x[23]
	value, err := s.load(s.x[24]+s.x[12], 1)
	if err != nil {
		return 0, fmt.Errorf("account six-byte copy seventh input: %w", err)
	}
	s.x[13] = uint64(byte(value))
	if err := s.storeValue(s.x[20]+s.x[12], uint64(byte(value)), 1); err != nil {
		return 0, fmt.Errorf("account six-byte copy seventh output: %w", err)
	}
	s.x[19] -= 8
	s.x[12] = s.x[19] - s.x[10]
	s.x[13] = s.x[10] - s.x[19]
	s.x[12] |= s.x[13]
	value, err = s.load(s.x[24]+s.x[19], 1)
	if err != nil {
		return 0, fmt.Errorf("account six-byte copy eighth input: %w", err)
	}
	s.x[13] = uint64(byte(value))
	s.x[12] >>= 63
	if err := s.storeValue(s.x[20]+s.x[19], uint64(byte(value)), 1); err != nil {
		return 0, fmt.Errorf("account six-byte copy eighth output: %w", err)
	}
	s.x[13] = s.x[12] << 2
	s.x[12] = uint64(uint32(s.x[13]) | uint32(s.x[12]<<4))
	s.x[13] = uint64(uint32(s.x[10]) ^ uint32(s.x[3]))
	s.x[12] = uint64(uint32(s.x[13]) ^ uint32(s.x[12]))
	table, err := s.load(s.sp+0xdf0, 8)
	if err != nil {
		return 0, fmt.Errorf("account six-byte copy table: %w", err)
	}
	s.x[13] = table
	next, err := loadBootstrapPointer(s, table, uint32(s.x[12]))
	if err != nil {
		return 0, fmt.Errorf("account six-byte copy target: %w", err)
	}
	s.x[12] = next
	s.x[17] = 0x3738
	return next, nil
}

func runAccountBitStore(s *state) (uint64, error) {
	s.traceBlock(accountBitStore)
	s.x[11] = uint64(uint32(s.x[1]) | uint32(s.x[2]))
	s.x[17] = s.x[11] ^ s.x[3]
	s.x[17] += s.x[10]
	s.x[22] = s.x[10] - 1
	base, err := s.load(s.x[19]+8, 8)
	if err != nil {
		return 0, fmt.Errorf("account bit-store base: %w", err)
	}
	s.x[23] = base
	word, err := s.load(base+(s.x[17]<<2), 4)
	if err != nil {
		return 0, fmt.Errorf("account bit-store word: %w", err)
	}
	s.x[17] = uint64(uint32(word))
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[12]))
	s.x[24] = uint64(uint32(shiftRight(s.x[17], int(s.x[13]), 32)))
	s.x[21] = uint64(uint32(s.x[21]) ^ uint32(s.x[4]))
	s.x[21] = uint64(uint32(s.x[21]) &^ uint32(s.x[24]))
	s.x[24] = uint64(uint32(s.x[24]) & uint32(s.x[5]))
	s.x[21] = uint64(uint32(s.x[24]) | uint32(s.x[21]))
	s.x[21] = uint64(uint32(s.x[21]) ^ uint32(s.x[4]))
	if err := s.storeValue(base+int32Offset(s.x[0], 2), uint64(uint32(s.x[21])), 4); err != nil {
		return 0, fmt.Errorf("account bit-store output: %w", err)
	}
	s.x[17] = uint64(uint32(shiftLeft(s.x[17], int(s.x[16]), 32)))
	s.x[17] = uint64(uint32(shiftLeft(s.x[17], int(s.x[15]), 32)))
	s.x[21] = uint64(uint32(s.x[17]) ^ uint32(s.x[12]))
	s.x[0]--
	s.x[10] += s.x[6]
	s.x[10] = uint64(^uint32(s.x[10]))
	s.x[10] = uint64(uint32(s.x[10]) >> 31)
	s.x[10] = uint64(uint32(s.x[10]) * uint32(s.x[7]))
	s.x[10] = uint64(uint32(s.x[10]) ^ uint32(s.x[14]))
	next, err := loadBootstrapPointer(s, s.x[20], uint32(s.x[10]))
	if err != nil {
		return 0, fmt.Errorf("account bit-store target: %w", err)
	}
	s.x[23] = next
	s.x[10] = s.x[22]
	s.x[17] = 0x390b
	return next, nil
}

func int32Offset(value uint64, shift uint) uint64 {
	return uint64(int64(int32(value)) << shift)
}

func runAccountByteAdvance(s *state) (uint64, error) {
	s.traceBlock(accountByteAdvance)
	s.x[17] = s.x[14] + 1
	s.x[9] = s.x[14] + s.x[11]
	s.x[9] += s.x[12]
	if err := s.storeValue(s.x[9], uint64(byte(s.x[13])), 1); err != nil {
		return 0, fmt.Errorf("account byte-advance output: %w", err)
	}
	s.x[9] = uint64(uint32(s.x[19]) + uint32(s.x[10]))
	s.x[1] = s.x[9] + s.x[16]
	s.x[14] += s.x[15]
	s.traceBlock(accountByteAdvanceTail)
	s.compare64(s.x[14], s.x[1])
	gate := uint32(boolValue(s.condition("lo")))
	s.x[14] = uint64(gate * uint32(s.x[0]))
	s.x[14] = uint64(uint32(s.x[14]) ^ uint32(s.x[8]))
	next, err := loadBootstrapPointer(s, s.x[24], uint32(s.x[14]))
	if err != nil {
		return 0, fmt.Errorf("account byte-advance target: %w", err)
	}
	s.x[1] = next
	s.x[14] = s.x[17]
	s.x[17] = 0xf36e
	return next, nil
}

func runAccountSumSelect(s *state) (uint64, error) {
	s.traceBlock(accountSumSelect)
	s.x[14] = s.x[16]
	s.x[15] = s.x[13] + s.x[22]
	s.x[0] = uint64(uint32(s.x[15]) | uint32(s.x[20]))
	s.x[15] = s.x[0] + s.x[21]
	s.compare64(s.x[16], s.x[11])
	gate := uint32(boolValue(s.condition("eq")))
	s.x[16] = uint64(gate)
	s.x[17] = s.x[0] + s.x[19]
	s.x[16] = uint64(gate * uint32(s.x[17]))
	s.x[13] = uint64(uint32(s.x[16]) ^ uint32(s.x[13]))
	next, err := loadBootstrapPointer(s, s.x[27], uint32(s.x[13]))
	if err != nil {
		return 0, fmt.Errorf("account sum-select target: %w", err)
	}
	s.x[13] = next
	s.x[17] = 0x7db5
	return next, nil
}

func runAccountOutputByte(s *state) (uint64, error) {
	s.traceBlock(accountOutputByte)
	s.compare32(uint32(s.x[14]), uint32(s.x[10]))
	first := uint32(boolValue(s.condition("lo")))
	s.x[12] = uint64(first)
	s.x[19] = uint64(uint32(s.x[8]) ^ 0xc65f)
	s.x[16] = uint64(uint32(s.x[8]) ^ 0x63cd)
	value, err := s.load(s.x[2]+uint64(uint32(s.x[9])), 1)
	if err != nil {
		return 0, fmt.Errorf("account output-byte input: %w", err)
	}
	s.x[17] = uint64(byte(value))
	s.x[0] = uint64(uint32(s.x[11]) + 0xbb949c9f)
	s.x[1] = uint64(0x36 & uint32(s.x[17]<<1))
	s.x[17] = uint64(uint32(s.x[17]) - uint32(s.x[1]))
	s.x[17] = uint64(uint32(s.x[17]) + 0x1b)
	s.x[1] = s.x[29] - 0xdc
	if err := s.storeValue(s.x[1]+uint64(uint32(s.x[0])), uint64(byte(s.x[17])), 1); err != nil {
		return 0, fmt.Errorf("account output-byte output: %w", err)
	}
	s.x[11] = uint64(uint32(s.x[11]) + 1)
	s.x[17] = 0x575d8837
	s.x[16] = uint64(uint32(s.x[9]) + uint32(s.x[16]))
	s.compare32(uint32(s.x[16]), uint32(s.x[13]))
	second := uint32(boolValue(s.condition("hi")))
	s.x[0] = uint64(first ^ second)
	s.x[9] = uint64(uint32(s.x[9]) + uint32(s.x[17]))
	s.compare32(uint32(s.x[9]), uint32(s.x[14]))
	ordered := uint32(boolValue(s.condition("lo")))
	s.x[9] = uint64(ordered)
	s.compare32(uint32(s.x[0]), 0)
	selected := ordered
	if s.condition("ne") {
		selected = first
	}
	s.x[9] = uint64(selected ^ 1)
	s.x[9] = uint64(uint32(s.x[9]) * uint32(s.x[15]))
	s.x[9] = uint64(uint32(s.x[9]) ^ uint32(s.x[8]))
	next, err := loadBootstrapPointer(s, s.x[21], uint32(s.x[9]))
	if err != nil {
		return 0, fmt.Errorf("account output-byte target: %w", err)
	}
	s.x[12] = next
	s.x[9] = s.x[16]
	s.x[17] = 0xc4ce
	return next, nil
}

func runAccountPointerSelect(s *state) (uint64, error) {
	s.traceBlock(accountPointerSelect)
	s.x[8] = uint64(uint32(s.x[20]) + 0xcfefd0ea)
	s.x[9] = 0x14b7cfbf
	s.x[13] = uint64(uint32(s.x[8]) & uint32(s.x[9]))
	s.x[8] = uint64(uint32(s.x[25]) + 0xf6b061f9)
	s.x[8] = uint64(uint32(s.x[8]) + uint32(s.x[13]))
	s.x[8] = uint64(int64(int32(s.x[8])))
	next, err := s.load(s.x[30]+uint64(uint32(s.x[20]))*8, 8)
	if err != nil {
		return 0, fmt.Errorf("account pointer-select target: %w", err)
	}
	s.x[9] = next
	s.x[17] = 0xb6ab
	return next, nil
}

func runAccountCallReturn(s *state) (uint64, error) {
	s.traceBlock(accountCallReturn)
	s.x[6] = s.x[26]
	s.x[4] = 0x2b7ffafa
	s.x[3] = 0xce757bc5
	s.x[2] = 0x43
	value, err := s.load(s.sp+0x208, 8)
	if err != nil {
		return 0, fmt.Errorf("account call-return link: %w", err)
	}
	s.x[30] = value
	s.x[8] = 0x9ceaf78a
	s.x[8] = uint64(uint32(s.x[8]) & uint32(s.x[0]<<1))
	s.x[8] = s.subtractWithFlags(uint64(uint32(s.x[0])), uint64(uint32(s.x[8])), 32)
	s.x[5] = uint64(uint32(s.x[8]) + uint32(s.x[3]))
	gate := uint32(boolValue(s.condition("eq")))
	s.x[8] = uint64(gate * 0x6d)
	s.x[8] = uint64(uint32(s.x[8]) ^ uint32(s.x[28]))
	next, err := loadBootstrapPointer(s, s.x[30], uint32(s.x[8]))
	if err != nil {
		return 0, fmt.Errorf("account call-return target: %w", err)
	}
	s.x[8] = next
	s.x[1] = 0x04a7cf03
	s.x[17] = 0x8231
	return next, nil
}

func runAccountConditionalPick(s *state) (uint64, error) {
	s.traceBlock(accountConditionalPick)
	s.x[9] = uint64(uint32(s.x[11]) + 0x410be7ff)
	s.x[10] = 0x94900740
	s.x[22] = uint64(uint32(s.x[9]) | uint32(s.x[10]))
	s.x[9] = uint64(uint32(s.x[11]) + 0x8217a752)
	s.x[10] = 0x7de879dd
	s.x[9] = uint64(uint32(s.x[9]) & uint32(s.x[10]))
	word, err := s.load(s.x[8]+0x24, 4)
	if err != nil {
		return 0, fmt.Errorf("account conditional-pick word: %w", err)
	}
	s.x[10] = uint64(uint32(word))
	s.x[8] = uint64(uint32(s.x[10]) + 0xba746ed6)
	s.x[12] = 0x3f1c3dd9
	s.compare32(uint32(s.x[8]), uint32(s.x[12]))
	stackWord, err := s.load(s.sp+0x1b8, 4)
	if err != nil {
		return 0, fmt.Errorf("account conditional-pick stack word: %w", err)
	}
	s.x[12] = uint64(uint32(stackWord))
	if s.condition("le") {
		s.compare32(uint32(s.x[8]), uint32(s.x[12]))
	} else {
		s.setNZCV(8)
	}
	gate := uint32(boolValue(s.condition("lt")))
	s.x[8] = uint64(gate)
	s.x[12] = uint64(uint32(s.x[9]) + uint32(s.x[22]))
	s.x[13] = 0x2a6472b4
	s.x[12] = uint64(uint32(s.x[12]) + uint32(s.x[13]))
	s.x[8] = uint64(gate * uint32(s.x[12]))
	s.x[8] = uint64(uint32(s.x[8]) ^ uint32(s.x[11]))
	next, err := loadBootstrapPointer(s, s.x[30], uint32(s.x[8]))
	if err != nil {
		return 0, fmt.Errorf("account conditional-pick target: %w", err)
	}
	s.x[11] = next
	s.x[5] = 0x318a201b
	s.x[8] = 0xce757bc5
	s.x[25] = s.x[10]
	s.x[1] = 0x04a7cf03
	s.x[0] = 0x30
	s.x[7] = 0x0a5cd4cf
	s.x[17] = 0x8231
	return next, nil
}
