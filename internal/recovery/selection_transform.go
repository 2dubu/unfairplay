package recovery

import "fmt"

const (
	selectionSIDBProjectionPass = uint64(0x10006c78c)
	selectionPackedWordSeed     = uint64(0x10072c394)
	selectionPackedWordStep     = uint64(0x10072c400)
	selectionByteWriter         = uint64(0x10087c89c)
	selectionMixTail            = uint64(0x100450400)
	selectionReverseProjection  = uint64(0x10013df58)
)

func runSelectionSemanticBlock(s *state, pc uint64) (uint64, bool, error) {
	switch pc {
	case selectionSIDBProjectionPass:
		next, err := runSelectionSIDBProjection(s)
		return next, true, err
	case selectionPackedWordSeed:
		next, err := runSelectionPackedWordPass(s)
		return next, true, err
	case selectionByteWriter:
		next, err := runSelectionByteWriterPass(s)
		return next, true, err
	case selectionMixTail:
		next, err := runSelectionMixTail(s)
		return next, true, err
	case selectionReverseProjection:
		next, err := runSelectionReverseProjection(s)
		return next, true, err
	default:
		return pc, false, nil
	}
}

func runSelectionReverseProjection(s *state) (uint64, error) {
	const maximumBytes = 2048
	if s == nil || s.x[10] == 0 || s.x[10] > maximumBytes {
		return 0, fmt.Errorf("unsupported selection reverse-projection length %d", s.x[10])
	}
	targets := [2]uint64{}
	for more := uint32(0); more < 2; more++ {
		index := uint32(s.x[1]) ^ more*uint32(s.x[3])
		next, err := loadBootstrapPointer(s, s.x[22], index)
		if err != nil {
			return 0, fmt.Errorf("selection reverse-projection target %d: %w", more, err)
		}
		targets[more] = next
	}
	if s.x[10] > 1 && targets[1] != selectionReverseProjection {
		return 0, fmt.Errorf("unsupported selection reverse-projection continuation %#x", targets[1])
	}
	for s.x[10] != 0 {
		s.traceBlock(selectionReverseProjection)
		s.x[11] = uint64(uint32(s.x[16]) & uint32(s.x[0]))
		s.x[10] = s.subtractWithFlags(s.x[10], 1, 64)
		input, err := s.load(s.x[8]+s.x[10], 1)
		if err != nil {
			return 0, fmt.Errorf("selection reverse-projection input: %w", err)
		}
		lane := s.x[10] & 0xf
		first, err := s.load(s.x[15]+lane, 1)
		if err != nil {
			return 0, fmt.Errorf("selection reverse-projection first lane: %w", err)
		}
		second, err := s.load(s.x[12]+lane+1, 1)
		if err != nil {
			return 0, fmt.Errorf("selection reverse-projection second lane: %w", err)
		}
		third, err := s.load(s.x[13]+lane, 1)
		if err != nil {
			return 0, fmt.Errorf("selection reverse-projection third lane: %w", err)
		}
		s.x[4] = uint64(byte(third))
		s.x[5] = uint64(byte(first))
		s.x[6] = uint64(byte(second))
		s.x[7] = uint64(uint32(lane) * uint32(s.x[2]))
		s.x[19] = uint64(uint32(s.x[11]) - 0x10)
		output := byte(input) ^ byte(first) ^ byte(second) ^ byte(s.x[19]) ^ byte(s.x[7]) ^ byte(third)
		s.x[5] = uint64(uint32(byte(second)) ^ uint32(s.x[19]))
		if err := s.storeValue(s.x[14]+s.x[10], uint64(output), 1); err != nil {
			return 0, fmt.Errorf("selection reverse-projection output: %w", err)
		}
		more := uint32(boolValue(s.condition("ne")))
		next := targets[more]
		s.x[4] = next
		s.x[17] = 0x585b
		if next != selectionReverseProjection {
			return next, nil
		}
	}
	return 0, fmt.Errorf("selection reverse projection did not exit")
}

func runSelectionByteWriterPass(s *state) (uint64, error) {
	const maximumBytes = 4096
	if s == nil || s.x[11] >= maximumBytes {
		return 0, fmt.Errorf("unsupported selection byte-writer length %d", s.x[11]+1)
	}
	length := int(s.x[11] + 1)
	dispatchTable, err := s.load(s.sp+0xdf0, 8)
	if err != nil {
		return 0, fmt.Errorf("selection byte-writer table: %w", err)
	}
	continueTarget, err := loadBootstrapPointer(s, dispatchTable, uint32(s.x[15])^uint32(s.x[12]))
	if err != nil {
		return 0, fmt.Errorf("selection byte-writer continuation: %w", err)
	}
	if length > 1 && continueTarget != selectionByteWriter {
		return 0, fmt.Errorf("unsupported selection byte-writer continuation %#x", continueTarget)
	}
	if byte(s.x[13]) != 0x36 || byte(s.x[14]) != 0xff {
		return 0, fmt.Errorf("unsupported selection byte-writer encoding")
	}
	for index := 0; index < length; index++ {
		s.traceBlock(selectionByteWriter)
		value, err := s.load(s.x[20]+uint64(uint32(s.x[9])), 1)
		if err != nil {
			return 0, fmt.Errorf("selection byte-writer input: %w", err)
		}
		output := byte(value) ^ 0x1b
		if err := s.storeValue(s.x[16], uint64(output), 1); err != nil {
			return 0, fmt.Errorf("selection byte-writer output: %w", err)
		}
		s.x[9] = uint64(uint32(s.x[9]) + 1)
		s.x[16]++
	}
	next, err := loadBootstrapPointer(s, dispatchTable, uint32(s.x[12]))
	if err != nil {
		return 0, fmt.Errorf("selection byte-writer target: %w", err)
	}
	s.x[0] = next
	s.x[11] = ^uint64(0)
	s.x[17] = 0x18b5
	s.compare64(0, 0)
	return next, nil
}

func runSelectionMixTail(s *state) (uint64, error) {
	s.traceBlock(selectionMixTail)
	s.x[0] = uint64(uint32(s.x[13]) ^ uint32(s.x[0]))
	s.x[13] = uint64(uint32(s.x[13]) ^ uint32(s.x[22]))
	s.x[17] = uint64(uint32(s.x[0]) ^ uint32(s.x[17]))
	s.x[13] = uint64(uint32(s.x[13]) + uint32(s.x[23]))
	s.x[13] = uint64(uint32(s.x[13]) ^ uint32(s.x[24]))
	s.x[13] = uint64(uint32(s.x[17]) ^ uint32(s.x[13]))
	s.compare32(uint32(s.x[5]), 0x1f)
	gate := uint32(boolValue(s.condition("hi")))
	s.x[17] = uint64(gate)
	s.x[0] = uint64(gate << 7)
	s.x[17] = uint64(uint32(s.x[0]) - gate<<4)
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[1]))
	table, err := s.load(s.sp+0xdf0, 8)
	if err != nil {
		return 0, fmt.Errorf("selection mix-tail table: %w", err)
	}
	s.x[0] = table
	next, err := loadBootstrapPointer(s, table, uint32(s.x[17]))
	if err != nil {
		return 0, fmt.Errorf("selection mix-tail target: %w", err)
	}
	s.x[25] = next
	s.x[2]++
	s.x[5] = uint64(uint32(s.x[5]) + 1)
	s.x[0] = s.x[16]
	s.x[17] = 0x6202
	return next, nil
}

func runSelectionByteWriter(s *state) (uint64, error) {
	s.traceBlock(selectionByteWriter)
	value, err := s.load(s.x[20]+uint64(uint32(s.x[9])), 1)
	if err != nil {
		return 0, fmt.Errorf("selection byte-writer input: %w", err)
	}
	s.x[17] = uint64(byte(value))
	s.x[9] = uint64(uint32(s.x[9]) + 1)
	s.x[0] = uint64(uint32(s.x[13]) & uint32(s.x[17]<<1))
	s.x[0] = uint64(uint32(s.x[14]) ^ uint32(s.x[0]))
	s.x[17] = uint64(uint32(s.x[17]) + uint32(s.x[0]))
	s.x[17] = uint64(uint32(s.x[17]) + 0x1c)
	if err := s.storeValue(s.x[16], uint64(byte(s.x[17])), 1); err != nil {
		return 0, fmt.Errorf("selection byte-writer output: %w", err)
	}
	s.x[16]++
	s.compare64(s.x[11], 0)
	gate := uint32(boolValue(s.condition("ne")))
	s.x[17] = uint64(gate * uint32(s.x[15]))
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[12]))
	table, err := s.load(s.sp+0xdf0, 8)
	if err != nil {
		return 0, fmt.Errorf("selection byte-writer table: %w", err)
	}
	s.x[0] = table
	next, err := loadBootstrapPointer(s, table, uint32(s.x[17]))
	if err != nil {
		return 0, fmt.Errorf("selection byte-writer target: %w", err)
	}
	s.x[0] = next
	s.x[11]--
	s.x[17] = 0x18b5
	return next, nil
}

func runSelectionPackedWordPass(s *state) (uint64, error) {
	const maximumWords = 2048
	if s == nil {
		return 0, fmt.Errorf("selection packed-word state is unavailable")
	}
	output, err := s.load(s.sp+0xd98, 8)
	if err != nil {
		return 0, fmt.Errorf("selection packed-word output table: %w", err)
	}
	selectorValue, err := s.load(s.sp+0xdc0, 4)
	if err != nil {
		return 0, fmt.Errorf("selection packed-word selector: %w", err)
	}
	selector := uint32(selectorValue)
	dispatchTable, err := s.load(s.sp+0xdf0, 8)
	if err != nil {
		return 0, fmt.Errorf("selection packed-word dispatch table: %w", err)
	}
	exitTarget, err := loadBootstrapPointer(s, dispatchTable, selector)
	if err != nil {
		return 0, fmt.Errorf("selection packed-word exit: %w", err)
	}
	continueTarget, err := loadBootstrapPointer(s, dispatchTable, selector^uint32(s.x[30]))
	if err != nil {
		return 0, fmt.Errorf("selection packed-word continuation: %w", err)
	}
	if continueTarget != selectionPackedWordSeed {
		return 0, fmt.Errorf("unsupported selection packed-word continuation %#x", continueTarget)
	}
	if uint32(s.x[0]) != 0x1b000000 || uint32(s.x[3]) != 0x4e757bdf ||
		uint32(s.x[6]) != 0x1a || uint32(s.x[7]) != 0xce757bc5 ||
		uint32(s.x[10]) != 0x1b || uint32(s.x[15]) != 0x9b4e23cd {
		return 0, fmt.Errorf("unsupported selection packed-word encoding")
	}

	for index := 0; index < maximumWords; index++ {
		s.traceBlock(selectionPackedWordSeed)
		var packed uint32
		for valueIndex, offset := range [...]uint32{0x9b4e23ca, 0x9b4e23cb, 0x9b4e23cc, 0x9b4e23cd} {
			value, err := s.load(s.x[8]+uint64(uint32(s.x[16])+offset), 1)
			if err != nil {
				return 0, fmt.Errorf("selection packed-word input %d: %w", valueIndex, err)
			}
			packed |= uint32(byte(value)) << (valueIndex * 8)
		}
		s.x[9] = uint64(packed ^ 0xd56e60de)

		s.traceBlock(selectionPackedWordStep)
		s.x[12] = uint64(uint32(s.x[16]) ^ uint32(s.x[19]))
		s.x[17] = uint64(uint32(s.x[16]) ^ uint32(s.x[21]))
		s.x[17] = uint64(uint32(s.x[17]) + uint32(s.x[22]))
		s.x[1] = uint64(uint32(s.x[16]) ^ uint32(s.x[23]))
		s.x[14] = uint64(uint32(s.x[16]) ^ uint32(s.x[25]))
		s.x[12] = uint64(uint32(s.x[12]) + uint32(s.x[20]))
		s.x[12] = uint64(uint32(s.x[12]) ^ uint32(s.x[16]))
		s.x[1] = uint64(uint32(s.x[1]) + uint32(s.x[24]))
		s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[1]))
		s.x[12] = uint64(uint32(s.x[12]) ^ uint32(s.x[17]))
		s.x[14] = uint64(uint32(s.x[14]) + uint32(s.x[26]))
		s.x[12] = uint64(uint32(s.x[12]) ^ uint32(s.x[14]))
		s.x[12] = uint64(uint32(s.x[27]) ^ uint32(s.x[12]>>2))
		s.x[14] = output
		if err := s.storeValue(output+uint64(uint32(s.x[12]))*4, uint64(uint32(s.x[9])), 4); err != nil {
			return 0, fmt.Errorf("selection packed-word output: %w", err)
		}
		s.x[9] = uint64(uint32(s.x[16]) + uint32(s.x[28]))
		s.x[16] = uint64(uint32(s.x[16]) + 4)
		s.compare32(uint32(s.x[9]), 0x40)
		more := s.condition("lo")
		s.x[12] = uint64(selector)
		s.x[12] = dispatchTable
		next := exitTarget
		if more {
			next = continueTarget
		}
		s.x[9] = next
		s.x[17] = 0xbb58
		if next != selectionPackedWordSeed {
			return next, nil
		}
	}
	return 0, fmt.Errorf("selection packed-word pass exceeded %d words", maximumWords)
}

func runSelectionPackedWordIteration(s *state) (uint64, error) {
	s.traceBlock(selectionPackedWordSeed)
	indices := [...]uint32{0x9b4e23ca, 0x9b4e23cb, 0x9b4e23cc, uint32(s.x[15])}
	var values [4]uint32
	for index, offset := range indices {
		value, err := s.load(s.x[8]+uint64(uint32(s.x[16])+offset), 1)
		if err != nil {
			return 0, fmt.Errorf("selection packed-word input %d: %w", index, err)
		}
		values[index] = uint32(byte(value))
		if index != len(values)-1 {
			values[index] ^= uint32(s.x[10])
		}
	}
	s.x[17] = uint64(values[0])
	s.x[1] = uint64(values[0] | values[1]<<8)
	s.x[9] = uint64(values[2]<<16 | values[3]<<24 | uint32(s.x[1]))
	s.x[9] = uint64(uint32(s.x[9]) ^ uint32(s.x[0]))
	s.x[12] = uint64(uint32(s.x[9]) & uint32(s.x[3]))
	s.x[17] = uint64(uint32(s.x[17]) & uint32(s.x[6]))
	s.x[12] = uint64(uint32(s.x[12]) ^ uint32(s.x[17]))
	s.x[9] = uint64(uint32(s.x[9]) - uint32(s.x[12]<<1))
	s.x[9] = uint64(uint32(s.x[9]) + uint32(s.x[7]))

	s.traceBlock(selectionPackedWordStep)
	s.x[12] = uint64(uint32(s.x[16]) ^ uint32(s.x[19]))
	s.x[17] = uint64(uint32(s.x[16]) ^ uint32(s.x[21]))
	s.x[17] = uint64(uint32(s.x[17]) + uint32(s.x[22]))
	s.x[1] = uint64(uint32(s.x[16]) ^ uint32(s.x[23]))
	s.x[14] = uint64(uint32(s.x[16]) ^ uint32(s.x[25]))
	s.x[12] = uint64(uint32(s.x[12]) + uint32(s.x[20]))
	s.x[12] = uint64(uint32(s.x[12]) ^ uint32(s.x[16]))
	s.x[1] = uint64(uint32(s.x[1]) + uint32(s.x[24]))
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[1]))
	s.x[12] = uint64(uint32(s.x[12]) ^ uint32(s.x[17]))
	s.x[14] = uint64(uint32(s.x[14]) + uint32(s.x[26]))
	s.x[12] = uint64(uint32(s.x[12]) ^ uint32(s.x[14]))
	s.x[12] = uint64(uint32(s.x[27]) ^ uint32(s.x[12]>>2))
	output, err := s.load(s.sp+0xd98, 8)
	if err != nil {
		return 0, fmt.Errorf("selection packed-word output table: %w", err)
	}
	s.x[14] = output
	if err := s.storeValue(output+uint64(uint32(s.x[12]))*4, uint64(uint32(s.x[9])), 4); err != nil {
		return 0, fmt.Errorf("selection packed-word output: %w", err)
	}
	s.x[9] = uint64(uint32(s.x[16]) + uint32(s.x[28]))
	s.x[16] = uint64(uint32(s.x[16]) + 4)
	s.compare32(uint32(s.x[9]), 0x40)
	dispatchIndex := uint32(boolValue(s.condition("lo"))) * uint32(s.x[30])
	selector, err := s.load(s.sp+0xdc0, 4)
	if err != nil {
		return 0, fmt.Errorf("selection packed-word selector: %w", err)
	}
	s.x[12] = uint64(uint32(selector))
	dispatchIndex ^= uint32(selector)
	dispatchTable, err := s.load(s.sp+0xdf0, 8)
	if err != nil {
		return 0, fmt.Errorf("selection packed-word dispatch table: %w", err)
	}
	s.x[12] = dispatchTable
	next, err := loadBootstrapPointer(s, dispatchTable, dispatchIndex)
	if err != nil {
		return 0, fmt.Errorf("selection packed-word dispatch: %w", err)
	}
	s.x[9] = next
	s.x[17] = 0xbb58
	return next, nil
}

// runSelectionSIDBProjection reduces the native SIDB container to its 944-byte
// encoded payload. The four lookup lanes are build-global byte encodings.
func runSelectionSIDBProjection(s *state) (uint64, error) {
	const maximumBytes = 2048
	if s == nil {
		return 0, fmt.Errorf("selection SIDB projection state is unavailable")
	}
	boundary := uint64(uint32(s.x[10])+uint32(s.x[13])) ^ s.x[0]
	if boundary >= s.x[12] || s.x[12]-boundary > maximumBytes {
		return 0, fmt.Errorf("unsupported selection SIDB projection range %#x..%#x", boundary, s.x[12])
	}
	dispatchTable, err := s.load(addOffset(s.x[29], -0x90), 8)
	if err != nil {
		return 0, fmt.Errorf("selection SIDB projection dispatch table: %w", err)
	}
	continueTarget, err := loadBootstrapPointer(s, dispatchTable, uint32(s.x[10])^0x14)
	if err != nil {
		return 0, fmt.Errorf("selection SIDB projection continuation: %w", err)
	}
	if s.x[12]-boundary > 1 && continueTarget != selectionSIDBProjectionPass {
		return 0, fmt.Errorf("unsupported selection SIDB projection continuation %#x", continueTarget)
	}
	for position := s.x[12] - 1; ; position-- {
		s.traceBlock(selectionSIDBProjectionPass)
		input, err := s.load(s.x[9]+position, 1)
		if err != nil {
			return 0, fmt.Errorf("selection SIDB projection input: %w", err)
		}
		lane := position & 0xf
		first, err := s.load(s.x[14]+lane, 1)
		if err != nil {
			return 0, fmt.Errorf("selection SIDB projection first lane: %w", err)
		}
		second, err := s.load(s.x[16]+lane+3, 1)
		if err != nil {
			return 0, fmt.Errorf("selection SIDB projection second lane: %w", err)
		}
		third, err := s.load(s.x[11]+lane+2, 1)
		if err != nil {
			return 0, fmt.Errorf("selection SIDB projection third lane: %w", err)
		}
		output := byte(input) ^ byte(first) ^ byte(second) ^ byte(third) ^ byte(uint32(lane)*uint32(s.x[15])) ^ 0x99
		if err := s.storeValue(s.x[20]+position, uint64(output), 1); err != nil {
			return 0, fmt.Errorf("selection SIDB projection output: %w", err)
		}
		if position == boundary {
			s.x[2] = uint64(byte(first))
			s.x[3] = uint64(uint32(s.x[10]) + uint32(s.x[13]))
			s.x[4] = uint64(byte(third))
			break
		}
	}
	next, err := loadBootstrapPointer(s, dispatchTable, uint32(s.x[10]))
	if err != nil {
		return 0, fmt.Errorf("selection SIDB projection dispatch: %w", err)
	}
	s.x[1] = next
	s.x[12] = boundary
	s.x[17] = 0x4b2b
	s.compare64(boundary, boundary)
	return next, nil
}

func runSelectionSIDBProjectionStep(s *state) (uint64, error) {
	if s == nil {
		return 0, fmt.Errorf("selection SIDB projection state is unavailable")
	}
	s.traceBlock(selectionSIDBProjectionPass)
	s.x[12]--
	input, err := s.load(s.x[9]+s.x[12], 1)
	if err != nil {
		return 0, fmt.Errorf("selection SIDB projection input: %w", err)
	}
	s.x[17] = uint64(byte(input))
	lane := s.x[12] & 0xf
	s.x[1] = lane
	first, err := s.load(s.x[14]+lane, 1)
	if err != nil {
		return 0, fmt.Errorf("selection SIDB projection first lane: %w", err)
	}
	s.x[2] = uint64(byte(first))
	second, err := s.load(s.x[16]+lane+3, 1)
	if err != nil {
		return 0, fmt.Errorf("selection SIDB projection second lane: %w", err)
	}
	s.x[3] = uint64(byte(second))
	third, err := s.load(s.x[11]+lane+2, 1)
	if err != nil {
		return 0, fmt.Errorf("selection SIDB projection third lane: %w", err)
	}
	s.x[4] = uint64(byte(third))
	s.x[1] = uint64(uint32(lane) * uint32(s.x[15]))
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[2]))
	s.x[1] = uint64(uint32(s.x[3]) ^ uint32(s.x[1]))
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[1]))
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[4]))
	s.x[17] = uint64(uint32(s.x[17]) ^ 0x99999999)
	if err := s.storeValue(s.x[20]+s.x[12], uint64(byte(s.x[17])), 1); err != nil {
		return 0, fmt.Errorf("selection SIDB projection output: %w", err)
	}
	s.x[3] = uint64(uint32(s.x[10]) + uint32(s.x[13]))
	s.x[17] = s.x[3] ^ s.x[0]
	s.compare64(s.x[12], s.x[17])
	last := uint32(boolValue(s.condition("ne")))
	s.x[1] = uint64(last << 2)
	dispatchIndex := ((last << 2) | (last << 4)) ^ uint32(s.x[10])
	dispatchTable, err := s.load(addOffset(s.x[29], -0x90), 8)
	if err != nil {
		return 0, fmt.Errorf("selection SIDB projection dispatch table: %w", err)
	}
	next, err := loadBootstrapPointer(s, dispatchTable, dispatchIndex)
	if err != nil {
		return 0, fmt.Errorf("selection SIDB projection dispatch: %w", err)
	}
	s.x[1] = next
	s.x[17] = 0x4b2b
	return next, nil
}
