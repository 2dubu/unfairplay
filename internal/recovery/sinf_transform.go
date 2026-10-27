package recovery

import "fmt"

const (
	sinfPrecursorCodecPass = uint64(0x10015320c)
	sinfMaterialWriterPass = uint64(0x10005fde4)
	sinfParitySelector     = uint64(0x100143d6c)
	sinfModuloSelector     = uint64(0x10015a3c8)
	sinfWordSelector       = uint64(0x100803454)
	sinfByteSelector       = uint64(0x100605b18)
	sinfByteProjectionPass = uint64(0x10005bed4)
)

func runSINFEntrySemanticBlock(s *state, pc uint64) (uint64, bool, error) {
	switch pc {
	case sinfPrecursorCodecPass:
		next, err := runSINFPrecursorCodec(s)
		return next, true, err
	case sinfMaterialWriterPass:
		next, err := runSINFMaterialWriter(s)
		return next, true, err
	case sinfParitySelector:
		next, err := runSINFParitySelector(s)
		return next, true, err
	case sinfModuloSelector:
		next, err := runSINFModuloSelector(s)
		return next, true, err
	case sinfWordSelector:
		next, err := runSINFWordSelector(s)
		return next, true, err
	case sinfByteSelector:
		next, err := runSINFByteSelector(s)
		return next, true, err
	case sinfByteProjectionPass:
		next, err := runSINFByteProjection(s)
		return next, true, err
	case sinfModuloDispatch:
		next, err := runSINFModuloDispatch(s)
		return next, true, err
	case sinfByteWriter:
		next, err := runSINFByteWriter(s)
		return next, true, err
	case sinfReverseWriter:
		next, err := runSINFReverseWriter(s)
		return next, true, err
	case sinfStateSelector:
		next, err := runSINFStateSelector(s)
		return next, true, err
	case sinfBitSelector:
		next, err := runSINFBitSelector(s)
		return next, true, err
	case sinfBitPermuteStart:
		next, err := runSINFBitPermute(s)
		return next, true, err
	case sinfTableWriter:
		next, err := runSINFTableWriter(s)
		return next, true, err
	case sinfExpandedTableWriter:
		next, err := runSINFExpandedTableWriter(s)
		return next, true, err
	default:
		return pc, false, nil
	}
}

func runSINFByteProjection(s *state) (uint64, error) {
	const maximumBytes = 1024
	for index := 0; index < maximumBytes; index++ {
		next, err := runSINFByteProjectionStep(s)
		if err != nil || next != sinfByteProjectionPass {
			return next, err
		}
	}
	return 0, fmt.Errorf("SINF byte projection exceeded %d bytes", maximumBytes)
}

func runSINFByteProjectionStep(s *state) (uint64, error) {
	s.traceBlock(sinfByteProjectionPass)
	s.x[3] = uint64(uint32(s.x[13]) + 0x7ce)
	s.x[17] = s.x[3] ^ s.x[5]
	s.x[4] = s.addWithCarry(s.x[17], s.x[4], false, 64)
	value, err := s.load(s.x[16]+s.x[4], 1)
	if err != nil {
		return 0, fmt.Errorf("SINF byte-projection input: %w", err)
	}
	s.x[17] = uint64(byte(value))
	position := s.x[4] & 0xf
	s.x[20] = position
	first, err := s.load(s.x[1]+position, 1)
	if err != nil {
		return 0, fmt.Errorf("SINF byte-projection first lane: %w", err)
	}
	s.x[21] = uint64(byte(first))
	second, err := s.load(s.x[2]+position+3, 1)
	if err != nil {
		return 0, fmt.Errorf("SINF byte-projection second lane: %w", err)
	}
	s.x[22] = uint64(byte(second))
	s.x[23] = uint64(uint32(position) * uint32(s.x[6]))
	third, err := s.load(s.x[0]+position, 1)
	if err != nil {
		return 0, fmt.Errorf("SINF byte-projection third lane: %w", err)
	}
	s.x[20] = uint64(byte(third))
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[21]))
	s.x[20] = uint64(uint32(s.x[22]) ^ uint32(s.x[20]))
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[20]))
	s.x[20] = uint64(uint32(s.x[23]) ^ uint32(s.x[7]))
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[20]))
	if err := s.storeValue(s.x[12]+s.x[4], uint64(byte(s.x[17])), 1); err != nil {
		return 0, fmt.Errorf("SINF byte-projection output: %w", err)
	}
	dispatchIndex := uint32(boolValue(s.condition("eq")))*uint32(s.x[19]) ^ uint32(s.x[13])
	next, err := loadBootstrapPointer(s, s.x[8], dispatchIndex)
	if err != nil {
		return 0, fmt.Errorf("SINF byte-projection dispatch: %w", err)
	}
	s.x[17] = 0x1743
	s.x[20] = next
	return next, nil
}

func runSINFParitySelector(s *state) (uint64, error) {
	s.traceBlock(sinfParitySelector)
	s.x[10] = 0xe9b
	s.x[14] = uint64(uint32(s.x[11]) ^ uint32(s.x[10]))
	s.test32(uint32(s.x[9]) & 1)
	selected := uint32(boolValue(s.condition("eq")))
	s.x[15] = uint64(uint32(s.x[14]) + 0xffffa5d7)
	dispatchIndex := selected*uint32(s.x[15]) ^ uint32(s.x[11])
	next, err := loadBootstrapPointer(s, s.x[8], dispatchIndex)
	if err != nil {
		return 0, fmt.Errorf("SINF parity-selector dispatch: %w", err)
	}
	s.x[10] = next
	s.x[17] = 0xaee1
	return next, nil
}

func runSINFModuloSelector(s *state) (uint64, error) {
	s.traceBlock(sinfModuloSelector)
	s.x[10] = uint64(uint32(s.x[11]) - 0x199)
	s.x[14] = uint64(uint32(s.x[10]) | 0x800)
	s.x[10] = uint64(uint32(s.x[12]) * uint32(s.x[13]))
	s.x[13] = uint64(uint32(s.x[14]) ^ 0x257e)
	divisor := uint32(s.x[13])
	quotient := uint32(0)
	if divisor != 0 {
		quotient = uint32(s.x[10]) / divisor
	}
	s.x[15] = uint64(quotient)
	s.x[13] = uint64(uint32(s.x[10]) - quotient*divisor)
	next, err := loadBootstrapPointer(s, s.x[8], uint32(s.x[11]))
	if err != nil {
		return 0, fmt.Errorf("SINF modulo-selector dispatch: %w", err)
	}
	s.x[10] = next
	s.x[17] = 0xaee1
	return next, nil
}

func runSINFWordSelector(s *state) (uint64, error) {
	s.traceBlock(sinfWordSelector)
	s.x[0] = s.x[11] & 0xff
	word, err := s.load(s.x[10]+shiftLeft(s.x[0], 2, 64), 4)
	if err != nil {
		return 0, fmt.Errorf("SINF word-selector input: %w", err)
	}
	s.x[11] = uint64(uint32(s.x[8]) - uint32(word))
	s.compare32(uint32(s.x[11]), uint32(s.x[14]))
	matched := uint32(boolValue(s.condition("eq")))
	s.x[11] = uint64(uint32(s.x[9]) ^ matched)
	s.x[17] = s.x[11] & 1
	s.x[11] = uint64(uint32(s.x[9]) + uint32(s.x[12]))
	s.x[1] = uint64(uint32(s.x[11]) ^ uint32(s.x[15]))
	dispatchIndex := uint32(s.x[17])*uint32(s.x[1]) ^ uint32(s.x[9])
	dispatchTable, err := s.load(s.sp+0xdf0, 8)
	if err != nil {
		return 0, fmt.Errorf("SINF word-selector table: %w", err)
	}
	next, err := loadBootstrapPointer(s, dispatchTable, dispatchIndex)
	if err != nil {
		return 0, fmt.Errorf("SINF word-selector dispatch: %w", err)
	}
	s.x[1] = next
	s.x[17] = 0xe3da
	return next, nil
}

func runSINFByteSelector(s *state) (uint64, error) {
	s.traceBlock(sinfByteSelector)
	s.x[17] = uint64(uint32(s.x[9]) ^ uint32(s.x[13]))
	s.x[9] = uint64(uint32(s.x[17]) ^ uint32(s.x[13]))
	s.x[11] = uint64(uint32(s.x[0]) + 1)
	s.x[1] = uint64(uint32(s.x[9]) - 0x16)
	s.compare32(uint32(s.x[0]), uint32(byte(s.x[1])))
	dispatchIndex := uint32(boolValue(s.condition("lo")))*uint32(s.x[16]) ^ uint32(s.x[17])
	dispatchTable, err := s.load(s.sp+0xdf0, 8)
	if err != nil {
		return 0, fmt.Errorf("SINF byte-selector table: %w", err)
	}
	next, err := loadBootstrapPointer(s, dispatchTable, dispatchIndex)
	if err != nil {
		return 0, fmt.Errorf("SINF byte-selector dispatch: %w", err)
	}
	s.x[0] = next
	s.x[17] = 0x38bd
	return next, nil
}

// runSINFPrecursorCodec applies the position-dependent byte codec used for the
// 432-byte current-SINF private source. Bytes are independent; only the reverse
// loop counter and protected continuation are carried between iterations.
func runSINFPrecursorCodec(s *state) (uint64, error) {
	const maximumBytes = 1024
	for index := 0; index < maximumBytes; index++ {
		next, err := runSINFPrecursorCodecStep(s)
		if err != nil || next != sinfPrecursorCodecPass {
			return next, err
		}
	}
	return 0, fmt.Errorf("SINF precursor codec exceeded %d bytes", maximumBytes)
}

func runSINFPrecursorCodecStep(s *state) (uint64, error) {
	if s == nil {
		return 0, fmt.Errorf("SINF precursor codec state is unavailable")
	}
	s.traceBlock(sinfPrecursorCodecPass)
	s.x[9] = uint64(uint32(s.x[12]) | 0x02600022)
	oldCounter := s.x[11]
	s.x[11] = s.subtractWithFlags(oldCounter, 1, 64)
	input, err := s.load(s.x[23]+s.x[11], 1)
	if err != nil {
		return 0, fmt.Errorf("SINF precursor input: %w", err)
	}
	s.x[10] = uint64(byte(input))
	firstWord, err := s.load(0x101075000+0x588, 4)
	if err != nil {
		return 0, fmt.Errorf("SINF precursor first static word: %w", err)
	}
	s.x[17] = uint64(uint32(firstWord))
	secondWord, err := s.load(s.x[25]+0xb40, 4)
	if err != nil {
		return 0, fmt.Errorf("SINF precursor second static word: %w", err)
	}
	s.x[20] = uint64(uint32(secondWord))
	s.x[21] = uint64(uint32(s.x[17]) ^ uint32(s.x[20]))
	s.x[21] = uint64(uint32(s.x[21]) ^ uint32(s.x[15]))
	s.x[21] = uint64(uint32(s.x[21]) * uint32(s.x[16]))
	firstBase, err := s.load(addOffset(s.x[29], -0xc8), 8)
	if err != nil {
		return 0, fmt.Errorf("SINF precursor first dynamic table: %w", err)
	}
	s.x[1] = firstBase
	s.x[22] = firstBase + uint64(byte(s.x[21]))
	first, err := s.load(s.x[22], 1)
	if err != nil {
		return 0, fmt.Errorf("SINF precursor first dynamic byte: %w", err)
	}
	s.x[22] = uint64(byte(first))
	s.x[22] = uint64(uint32(s.x[22]) ^ uint32(s.x[0]))
	first, err = s.load(s.x[14]+uint64(uint32(s.x[22])), 1)
	if err != nil {
		return 0, fmt.Errorf("SINF precursor first substitution: %w", err)
	}
	s.x[22] = uint64(byte(first))
	s.x[21] = uint64(uint32(s.x[21]) ^ uint32(s.x[22]))
	s.x[17] = uint64(uint32(s.x[20]) - uint32(s.x[17]))
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[15]))
	s.x[17] = uint64(uint32(s.x[17]) * uint32(s.x[16]))
	s.x[20] = s.x[17] & 0xff
	s.x[21] &= 0xff
	s.x[22] = s.x[26] + s.x[20]
	second, err := s.load(s.x[22], 1)
	if err != nil {
		return 0, fmt.Errorf("SINF precursor second dynamic byte: %w", err)
	}
	s.x[22] = uint64(byte(second))
	s.x[22] = uint64(uint32(s.x[22]) ^ uint32(s.x[3]))
	second, err = s.load(s.x[28]+uint64(uint32(s.x[22])), 1)
	if err != nil {
		return 0, fmt.Errorf("SINF precursor second substitution: %w", err)
	}
	s.x[22] = uint64(byte(second))
	s.x[21] = uint64(uint32(s.x[21]) + 0x577)
	s.x[22] = uint64(uint32(s.x[17]) ^ uint32(s.x[22]))
	s.x[22] &= 0xff
	s.x[22] = uint64(uint32(s.x[22]) + 0xd6f)
	secondBase, err := s.load(addOffset(s.x[29], -0xf8), 8)
	if err != nil {
		return 0, fmt.Errorf("SINF precursor second dynamic table: %w", err)
	}
	s.x[20] = secondBase + s.x[20]
	s.x[1] = 0x10105b000 + 0x6f0
	s.x[22], err = s.load(s.x[1]+uint64(uint32(s.x[22]))*8, 8)
	if err != nil {
		return 0, fmt.Errorf("SINF precursor second pointer: %w", err)
	}
	second, err = s.load(s.x[20], 1)
	if err != nil {
		return 0, fmt.Errorf("SINF precursor second table byte: %w", err)
	}
	s.x[20] = uint64(byte(second))
	s.x[20] = uint64(uint32(s.x[20]) ^ uint32(s.x[6]))
	s.x[30] = s.x[23]
	positionTable, err := s.load(s.sp+0x4c8, 8)
	if err != nil {
		return 0, fmt.Errorf("SINF precursor position table: %w", err)
	}
	s.x[23] = positionTable
	second, err = s.load(positionTable+uint64(uint32(s.x[20])), 1)
	if err != nil {
		return 0, fmt.Errorf("SINF precursor position byte: %w", err)
	}
	s.x[20] = uint64(byte(second))
	s.x[21], err = s.load(s.x[1]+uint64(uint32(s.x[21]))*8, 8)
	if err != nil {
		return 0, fmt.Errorf("SINF precursor first pointer: %w", err)
	}
	position := s.x[11] & 0xf
	s.x[23] = position
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[20]))
	s.x[17] &= 0xff
	s.x[17] = uint64(uint32(s.x[17]) + 0x946)
	s.x[20] = position + s.x[4]
	s.x[20] = s.x[22] + s.x[20]
	s.x[17], err = s.load(s.x[1]+uint64(uint32(s.x[17]))*8, 8)
	if err != nil {
		return 0, fmt.Errorf("SINF precursor third pointer: %w", err)
	}
	s.x[22] = position + s.x[7]
	s.x[17] += s.x[22]
	s.x[22] = position + s.x[2]
	s.x[21] += s.x[22]
	s.x[22] = uint64(uint32(s.x[9]) - 0x1e)
	s.x[10] = uint64(uint32(s.x[22]) ^ uint32(s.x[10]))
	s.x[22] = uint64(uint32(position) * uint32(s.x[5]))
	s.x[23] = s.x[30]
	s.x[30], err = s.load(addOffset(s.x[29], -0x100), 8)
	if err != nil {
		return 0, fmt.Errorf("SINF precursor final table: %w", err)
	}
	s.x[10] = uint64(uint32(s.x[10]) ^ uint32(s.x[22]))
	first, err = s.load(s.x[21], 1)
	if err != nil {
		return 0, fmt.Errorf("SINF precursor first result byte: %w", err)
	}
	s.x[21] = uint64(byte(first))
	s.x[10] = uint64(uint32(s.x[10]) ^ uint32(s.x[21]))
	second, err = s.load(s.x[20], 1)
	if err != nil {
		return 0, fmt.Errorf("SINF precursor second result byte: %w", err)
	}
	s.x[20] = uint64(byte(second))
	third, err := s.load(s.x[17], 1)
	if err != nil {
		return 0, fmt.Errorf("SINF precursor third result byte: %w", err)
	}
	s.x[17] = uint64(uint32(s.x[20]) ^ uint32(byte(third)))
	s.x[10] = uint64(uint32(s.x[10]) ^ uint32(s.x[17]))
	if err := s.storeValue(s.x[8]+s.x[11], uint64(byte(s.x[10])), 1); err != nil {
		return 0, fmt.Errorf("SINF precursor output: %w", err)
	}
	last := uint32(boolValue(s.condition("eq")))
	dispatchIndex := last*uint32(s.x[19]) ^ uint32(s.x[13])
	next, err := loadBootstrapPointer(s, s.x[27], dispatchIndex)
	if err != nil {
		return 0, fmt.Errorf("SINF precursor dispatch: %w", err)
	}
	s.x[10] = next
	s.x[17] = 0x1840
	return next, nil
}

// runSINFMaterialWriter applies the recovered direct precursor-to-material
// projection. The native protected handler performs the same byte-wise mask
// while dispatching once per byte.
func runSINFMaterialWriter(s *state) (uint64, error) {
	if s == nil {
		return 0, fmt.Errorf("SINF material writer state is unavailable")
	}
	source, err := s.load(s.x[24]+0x10, 8)
	if err != nil {
		return 0, fmt.Errorf("SINF material precursor: %w", err)
	}
	const maximumBytes = 16
	for count := 0; count < maximumBytes; count++ {
		s.traceBlock(sinfMaterialWriterPass)
		oldCounter := s.x[8]
		s.x[8] = s.subtractWithFlags(oldCounter, 1, 64)
		if s.x[8] >= uint64(len(currentSINFMaterialMask)) {
			return 0, fmt.Errorf("SINF material index %d is outside the 16-byte block", s.x[8])
		}
		input, err := s.load(source+s.x[8], 1)
		if err != nil {
			return 0, fmt.Errorf("SINF material precursor byte: %w", err)
		}
		material := byte(input) ^ currentSINFMaterialMask[s.x[8]]
		destination := s.x[29] - 0x70 + s.x[8]
		if err := s.storeValue(destination, uint64(material), 1); err != nil {
			return 0, fmt.Errorf("SINF material output byte: %w", err)
		}
		more := uint32(boolValue(s.condition("ne")))
		dispatchIndex := uint32(s.x[13]) ^ (more << 6)
		next, err := loadBootstrapPointer(s, s.x[27], dispatchIndex)
		if err != nil {
			return 0, fmt.Errorf("SINF material writer dispatch: %w", err)
		}
		s.x[9] = destination
		s.x[10] = uint64(material)
		s.x[11] = next
		s.x[17] = 0x39ba
		if next != sinfMaterialWriterPass {
			return next, nil
		}
	}
	return 0, fmt.Errorf("SINF material writer did not exit after %d bytes", maximumBytes)
}
