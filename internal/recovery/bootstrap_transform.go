package recovery

import "fmt"

const (
	bootstrapMixerAfterRecurrence = uint64(0x10058a898)
	bootstrapMixerPrepare         = uint64(0x1008553bc)
	bootstrapMixerDispatch        = uint64(0x100855400)
	bootstrapMixerRound           = uint64(0x10090d240)
	bootstrapMixerContinue        = uint64(0x1008c3b90)
	bootstrapShortWordPass        = uint64(0x100b4dea0)
	bootstrapLongWordPass         = uint64(0x100b4e878)
	bootstrapAccumulatorPassA     = uint64(0x1009729ac)
	bootstrapAccumulatorPassB     = uint64(0x100916db0)
	bootstrapShufflePass          = uint64(0x10051b558)
	bootstrapShuffleExit          = uint64(0x10051b5a8)
	bootstrapForwardShuffleSeed   = uint64(0x1005df7f4)
	bootstrapForwardShuffleStep   = uint64(0x1005df800)
	bootstrapForwardShuffleExit   = uint64(0x1005df840)
	bootstrapReverseBytePass      = uint64(0x10002d8fc)
	bootstrapLookupCopySeed       = uint64(0x1008197dc)
	bootstrapLookupCopyStep       = uint64(0x100819800)
	bootstrapPointerWalkSeed      = uint64(0x1000a3fa4)
	bootstrapPointerWalkDispatch  = uint64(0x1000a4000)
	bootstrapPairCopyPass         = uint64(0x100176524)
	bootstrapConstantFillSeed     = uint64(0x1008e6fa4)
	bootstrapConstantFillDispatch = uint64(0x1008e7000)
	bootstrapPackedBitPass        = uint64(0x100251470)
	bootstrapTableAccumulatorPass = uint64(0x10021aea4)
	bootstrapSparseFillSeed       = uint64(0x1008a97ac)
	bootstrapSparseFillDispatch   = uint64(0x1008a9800)
	bootstrapNonlinearBytePass    = uint64(0x10094ca9c)
	bootstrapByteStatePass        = uint64(0x1008eb0a8)
	bootstrapAccumulatorCycleSeed = uint64(0x1008d0e38)
	bootstrapAccumulatorCyclePrep = uint64(0x1009f458c)
	bootstrapAccumulatorCycleKey  = uint64(0x1009839e4)
	bootstrapAccumulatorCycleTail = uint64(0x10080f1a8)
	bootstrapByteWriterControl    = uint64(0x10024d828)
	bootstrapByteWriterMix        = uint64(0x100272950)
	bootstrapByteWriterLoad       = uint64(0x1003a5a8c)
	bootstrapByteWriterStore      = uint64(0x1003076e4)
)

// runBootstrapMixer replaces the dominant flattened bootstrap path with its
// natural nested-loop structure. A round contains either 31 or 63 recurrence
// steps, followed by five small control blocks. If a different, valid path is
// selected, the caller resumes it through the verified CFG.
func runBootstrapMixerCircuit(s *state, _ *rawSourceRuntime, profile *executionProfile) (uint64, error) {
	if s == nil || profile == nil || profile.runBootstrap == nil {
		return 0, fmt.Errorf("bootstrap mixer state is unavailable")
	}
	pc := profile.bootstrapStart
	const maximumRounds = 4096
	for round := 0; round < maximumRounds; round++ {
		for step := 0; pc == profile.bootstrapStart; step++ {
			if step >= 128 {
				return 0, fmt.Errorf("bootstrap mixer round %d exceeded 128 recurrence steps", round)
			}
			if err := profile.runBootstrap(s, fault{}); err != nil {
				return 0, err
			}
			pc = s.x[3]
		}
		if pc != bootstrapMixerAfterRecurrence {
			return pc, nil
		}

		var err error
		pc, err = runBootstrapMixerAfterRecurrence(s)
		if err != nil || pc != bootstrapMixerPrepare {
			return pc, err
		}
		pc, err = runBootstrapMixerPrepareBlock(s)
		if err != nil || pc != bootstrapMixerDispatch {
			return pc, err
		}
		pc, err = runBootstrapMixerDispatchBlock(s)
		if err != nil || pc != bootstrapMixerRound {
			return pc, err
		}
		pc, err = runBootstrapMixerRoundBlock(s)
		if err != nil || pc != bootstrapMixerContinue {
			return pc, err
		}
		pc, err = runBootstrapMixerContinueBlock(s)
		if err != nil {
			return pc, err
		}
		if pc != profile.bootstrapStart {
			return pc, nil
		}
	}
	return 0, fmt.Errorf("bootstrap mixer exceeded %d rounds", maximumRounds)
}

func runBootstrapMixerAfterRecurrence(s *state) (uint64, error) {
	if s == nil {
		return 0, fmt.Errorf("bootstrap mixer controller state is unavailable")
	}
	s.traceBlock(bootstrapMixerAfterRecurrence)
	s.x[10] = 0x7b3b2d12
	s.compare32(uint32(s.x[8]), uint32(s.x[10]))
	firstGate := uint32(boolValue(s.condition("gt")))
	s.x[12] = uint64((uint32(s.x[11]) + 0x5b58340c) & 0xa4a77ff9)
	value := uint32(s.x[11]) + 0x38e13ca5
	s.x[5] = uint64(value & 0xc71ebff7)
	value = (uint32(s.x[5]) ^ 0x8d29b67c) + uint32(s.x[2])
	s.compare32(value, 0x84c4d2ed)
	selector := firstGate ^ uint32(boolValue(s.condition("lt")))
	stackBound, err := s.load(s.sp+0x38, 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap mixer controller bound: %w", err)
	}
	s.x[22] = uint64(uint32(stackBound))
	s.compare32(value, uint32(stackBound))
	secondGate := uint32(boolValue(s.condition("lt")))
	s.compare32(selector, 0)
	dispatchGate := secondGate
	if s.condition("ne") {
		dispatchGate = firstGate
	}
	s.x[14] = uint64(uint32(s.x[12]) + 0xffff9647)
	dispatchIndex := dispatchGate*uint32(s.x[14]) ^ uint32(s.x[11])
	next, err := loadBootstrapPointer(s, s.x[25], dispatchIndex)
	if err != nil {
		return 0, fmt.Errorf("bootstrap mixer controller dispatch: %w", err)
	}
	s.x[0] = 0x53ceff38
	s.x[1] = 0x224ed61b
	s.x[6] = 0x886438de
	s.x[9] = next
	s.x[13] = 0x1eb
	s.x[15] = 0x886438dd
	s.x[16] = 0x84c4d2ed
	s.x[17] = 0x76d5
	s.x[21] = 0xa9292112
	s.x[23] = 0x8d290bcc
	s.x[27] = 0xb3563f76
	return next, nil
}

func runBootstrapMixerPrepareBlock(s *state) (uint64, error) {
	s.traceBlock(bootstrapMixerPrepare)
	s.x[9] = uint64(uint32(s.x[12]) + 0xb00)
	s.x[10] = 0xcc2d
	s.x[20] = uint64(uint32(s.x[9]) ^ uint32(s.x[10]))
	s.x[10] = 0x296
	s.x[19] = uint64(uint32(s.x[9]) | uint32(s.x[10]))
	s.x[9] = s.x[19] + 0xb355bcc0
	s.x[9] ^= s.x[30] >> 32
	s.x[9] += s.x[28] >> 32
	pointer, err := s.load(s.sp+0x48, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap mixer prepare pointer: %w", err)
	}
	s.x[10] = pointer
	word, err := s.load(pointer, 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap mixer prepare word: %w", err)
	}
	s.x[11] = word ^ s.x[24]
	s.x[13] = 0xf27a3d3c33e26bf7
	return bootstrapMixerDispatch, nil
}

func runBootstrapMixerDispatchBlock(s *state) (uint64, error) {
	s.traceBlock(bootstrapMixerDispatch)
	s.x[14] = s.x[20] ^ s.x[13]
	s.x[9] += s.x[11]
	s.x[11] = s.x[9] + s.x[14]
	s.x[13] = 0x67c44ff4
	s.x[9] = s.x[13] & (s.x[9] << 1)
	s.x[9] = s.x[11] - s.x[9]
	s.x[11] = uint64(uint32(s.x[9]) ^ 0x11ac84c1)
	outputPointer, err := s.load(s.sp+0x40, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap mixer dispatch output pointer: %w", err)
	}
	s.x[14] = outputPointer
	if err := s.storeValue(outputPointer, uint64(uint32(s.x[11])), 4); err != nil {
		return 0, fmt.Errorf("bootstrap mixer dispatch first output: %w", err)
	}
	s.x[9] >>= 32
	s.x[11] = 0xd0349e07
	s.x[9] = uint64(uint32(s.x[9]) ^ uint32(s.x[11]))
	if err := s.storeValue(s.x[10], uint64(uint32(s.x[9])), 4); err != nil {
		return 0, fmt.Errorf("bootstrap mixer dispatch second output: %w", err)
	}
	counter, err := s.load(s.sp+0x80, 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap mixer dispatch counter: %w", err)
	}
	s.x[9] = uint64(uint32(counter))
	s.x[28] = uint64(uint32(counter) + 1)
	s.x[10] = 0x4fd98905
	s.x[9] = uint64(uint32(s.x[9]) + uint32(s.x[10]))
	s.x[10] = 0xe4e21e1a
	s.compare32(uint32(s.x[9]), uint32(s.x[10]))
	bound, err := s.load(s.sp+0x74, 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap mixer dispatch bound: %w", err)
	}
	s.x[11] = uint64(uint32(bound))
	if s.condition("ge") {
		s.compare32(uint32(s.x[9]), uint32(bound))
	} else {
		s.setNZCV(0)
	}
	dispatchBit := uint32(boolValue(s.condition("ge")))
	dispatchIndex := uint32(s.x[12]) | dispatchBit*0x51
	next, err := loadBootstrapPointer(s, s.x[25], dispatchIndex)
	if err != nil {
		return 0, fmt.Errorf("bootstrap mixer dispatch continuation: %w", err)
	}
	s.x[6] = 0xb3563f76
	s.x[9] = next
	s.x[13] = 0x886438de
	s.x[14] = s.x[28]
	s.x[17] = 0x8766
	s.x[21] = 0x886438dd
	s.x[27] = 0xa9292112
	return next, nil
}

func runBootstrapMixerRoundBlock(s *state) (uint64, error) {
	s.traceBlock(bootstrapMixerRound)
	flag, err := s.load(s.sp+0x78, 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap mixer round flag: %w", err)
	}
	s.x[9] = uint64(uint32(flag))
	s.compare32(uint32(flag), 1)
	dispatchBit := uint32(boolValue(s.condition("ne")))
	s.x[11] = uint64(uint32(s.x[20]) + 0x5237)
	s.x[3] = uint64(uint32(s.x[11]) | 0x88)
	s.x[11] = uint64(uint32(s.x[20]) + 0x750d)
	s.x[12] = 0xeaf76aea
	if err := s.storeValue(s.sp+0x80, uint64(uint32(s.x[14])), 4); err != nil {
		return 0, fmt.Errorf("bootstrap mixer round counter: %w", err)
	}
	s.x[12] = uint64(uint32(s.x[14]) + uint32(s.x[12]))
	coefficient, err := s.load(s.sp+0x60, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap mixer round coefficient: %w", err)
	}
	table, err := s.load(s.sp+0x68, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap mixer round table: %w", err)
	}
	s.x[17], s.x[14] = coefficient, table
	tableWord, err := loadBootstrapWord(s, table, uint32(s.x[12]))
	if err != nil {
		return 0, fmt.Errorf("bootstrap mixer round table word: %w", err)
	}
	s.x[12] = tableWord
	accumulator, err := s.load(s.sp+0x2d0, 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap mixer round accumulator: %w", err)
	}
	s.x[14] = accumulator ^ s.x[24]
	s.x[12] ^= s.x[24]
	s.x[28] = uint64(uint32(s.x[12]))*uint64(uint32(coefficient)) + s.x[14]
	multiplier, err := s.load(s.sp+0x5c, 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap mixer round multiplier: %w", err)
	}
	s.x[14] = uint64(uint32(multiplier))
	s.x[19] = uint64(uint32(s.x[28]) * uint32(multiplier))
	s.x[14] = uint64(uint32(s.x[28]))
	wordPointer, err := s.load(s.sp+0x50, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap mixer round word pointer: %w", err)
	}
	word, err := s.load(wordPointer, 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap mixer round word: %w", err)
	}
	s.x[17] = uint64(uint32(word))
	s.x[2] = uint64(uint32(s.x[20]) + 0x224df1ea)
	s.x[2] = uint64(uint32(s.x[3]) ^ uint32(s.x[2]))
	s.x[17] = uint64(uint32(s.x[2]) ^ uint32(s.x[17]))
	s.x[14] = uint64(uint32(s.x[19]))*uint64(uint32(s.x[17])) + s.x[14]
	s.x[17] = uint64(0x99538112e930fbad) | ^(s.x[14] << 1)
	s.x[2] = 0xb3563f768b67822a
	s.x[14] += s.x[2]
	s.x[30] = s.x[14] + s.x[17]
	dispatchIndex := dispatchBit*0x17 ^ uint32(s.x[20])
	next, err := loadBootstrapPointer(s, s.x[25], dispatchIndex)
	if err != nil {
		return 0, fmt.Errorf("bootstrap mixer round dispatch: %w", err)
	}
	s.x[0] = 0xe663
	s.x[1], err = s.load(s.sp+0x3c, 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap mixer round bound: %w", err)
	}
	s.x[2] = 0x779bc723
	s.x[9] = next
	s.x[10] = 0xa40ddddf
	s.x[14] = 0x17
	s.x[15] = 0x41e1812a
	s.x[16] = 0x4a45ba0a
	s.x[17] = 0xcd95
	s.x[23] = 0xc13a
	return next, nil
}

func runBootstrapMixerContinueBlock(s *state) (uint64, error) {
	s.traceBlock(bootstrapMixerContinue)
	leftPointer, err := s.load(s.sp+0x30, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap mixer continue left pointer: %w", err)
	}
	s.x[9] = leftPointer
	s.x[5], err = s.load(leftPointer+8, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap mixer continue left value: %w", err)
	}
	rightPointer, err := s.load(s.sp+0x98, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap mixer continue right pointer: %w", err)
	}
	s.x[9] = rightPointer
	s.x[22], err = s.load(rightPointer+8, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap mixer continue right value: %w", err)
	}
	s.x[20] = 0x779bc723
	s.x[9] = uint64(uint32(s.x[3]) + 0xffffad41)
	s.x[11] = 0x5bf2972d
	s.x[14] = uint64(uint32(s.x[9]) + uint32(s.x[11]))
	next, err := loadBootstrapPointer(s, s.x[25], uint32(s.x[3]))
	if err != nil {
		return 0, fmt.Errorf("bootstrap mixer continue dispatch: %w", err)
	}
	s.x[11] = next
	s.x[17] = 0xaa52
	return next, nil
}

// runBootstrapWordPass executes either of the two natural word-table passes
// used by the bootstrap finalizer. Their native handlers dispatch back to
// themselves once per word; this direct loop makes the real bounds explicit.
func runBootstrapWordPass(s *state, pc uint64) (uint64, error) {
	switch pc {
	case bootstrapShortWordPass:
		return runBootstrapShortWordDecode(s)
	case bootstrapLongWordPass:
		return runBootstrapLongWordDecode(s)
	default:
		return pc, fmt.Errorf("unknown bootstrap word pass %#x", pc)
	}
}

type bootstrapWordDecodeResult struct {
	Output, Encoded, Companion, Correction uint32
}

// decodeBootstrapWords removes the protected two-stream word representation
// used for RSA operands. One bit selects a correction word; the remaining bits
// form a shifted stream whose high bit is carried from the previous word.
func decodeBootstrapWords(companion, encoded []uint32, correction [2]uint32, previous uint32, shift uint, feedbackDistance int) ([]uint32, bootstrapWordDecodeResult) {
	output := make([]uint32, len(encoded))
	var result bootstrapWordDecodeResult
	for index, encodedWord := range encoded {
		folded := (encodedWord & 0x7ffffffe) | (previous & 0x80000000)
		correctionWord := correction[encodedWord&1]
		decoded := companion[index] ^ correctionWord ^ (folded >> shift)
		output[index] = decoded
		result = bootstrapWordDecodeResult{
			Output: decoded, Encoded: encodedWord,
			Companion: companion[index], Correction: correctionWord,
		}
		if feedbackDistance > 0 && index+feedbackDistance < len(companion) {
			companion[index+feedbackDistance] = decoded
		}
		previous = encodedWord
	}
	return output, result
}

func runBootstrapShortWordDecode(s *state) (uint64, error) {
	if s == nil {
		return 0, fmt.Errorf("bootstrap short-word state is unavailable")
	}
	start := uint32(s.x[8])
	if start > 0x388 || start&3 != 0 {
		return 0, fmt.Errorf("unsupported bootstrap short-word offset %#x", start)
	}
	count := int((0x388-start)/4 + 1)
	base := s.x[21] + uint64(start)
	companion, encoded, correction, err := loadBootstrapEncodedWordStreams(s, base+0x634, base+4, s.x[0], count)
	if err != nil {
		return 0, fmt.Errorf("bootstrap short-word streams: %w", err)
	}
	output, last := decodeBootstrapWords(companion, encoded, correction, uint32(s.x[10]), uint(s.x[12])&31, 0)
	for index, word := range output {
		if err := s.storeValue(base+uint64(index*4), uint64(word), 4); err != nil {
			return 0, fmt.Errorf("bootstrap short-word output %d: %w", index, err)
		}
	}
	next, err := loadBootstrapPointer(s, s.x[19], uint32(s.x[1])^uint32(s.x[9]))
	if err != nil {
		return 0, fmt.Errorf("bootstrap short-word dispatch: %w", err)
	}
	s.x[8] = 0x38c
	s.x[10] = uint64(last.Encoded)
	s.x[13] = next
	s.x[14] = uint64(last.Output)
	s.x[15] = uint64(last.Companion)
	s.x[16] = uint64(last.Correction)
	s.x[17] = 0x67ab
	s.compare64(0x388, 0x388)
	return next, nil
}

func runBootstrapLongWordDecode(s *state) (uint64, error) {
	if s == nil {
		return 0, fmt.Errorf("bootstrap long-word state is unavailable")
	}
	start := uint32(s.x[9])
	if start > 0x62c || start&3 != 0 {
		return 0, fmt.Errorf("unsupported bootstrap long-word offset %#x", start)
	}
	count := int((0x62c-start)/4 + 1)
	companion, encoded, correction, err := loadBootstrapEncodedWordStreams(s, s.x[21]+uint64(start), s.x[21]+uint64(start)+0x390, s.x[0], count)
	if err != nil {
		return 0, fmt.Errorf("bootstrap long-word streams: %w", err)
	}
	feedbackBytes := int64(s.x[11]) - int64(s.x[21])
	if feedbackBytes < 0 || feedbackBytes%4 != 0 || feedbackBytes/4 > int64(count) {
		return 0, fmt.Errorf("unsupported bootstrap long-word feedback distance %d", feedbackBytes)
	}
	output, last := decodeBootstrapWords(companion, encoded, correction, uint32(s.x[12]), 1, int(feedbackBytes/4))
	for index, word := range output {
		if err := s.storeValue(s.x[11]+uint64(start)+uint64(index*4), uint64(word), 4); err != nil {
			return 0, fmt.Errorf("bootstrap long-word output %d: %w", index, err)
		}
	}
	next, err := loadBootstrapPointer(s, s.x[19], uint32(s.x[1])^uint32(s.x[10]))
	if err != nil {
		return 0, fmt.Errorf("bootstrap long-word dispatch: %w", err)
	}
	s.x[9] = 0x630
	s.x[12] = uint64(last.Encoded)
	s.x[13] = next
	s.x[14] = uint64(last.Output)
	s.x[15] = uint64(last.Companion ^ last.Correction)
	s.x[16] = uint64(last.Correction)
	s.x[17] = 0x1028
	s.compare64(0x62c, 0x62c)
	return next, nil
}

func loadBootstrapEncodedWordStreams(s *state, companionBase, encodedBase, correctionBase uint64, count int) ([]uint32, []uint32, [2]uint32, error) {
	companion := make([]uint32, count)
	encoded := make([]uint32, count)
	for index := 0; index < count; index++ {
		value, err := s.load(companionBase+uint64(index*4), 4)
		if err != nil {
			return nil, nil, [2]uint32{}, fmt.Errorf("companion %d: %w", index, err)
		}
		companion[index] = uint32(value)
		value, err = s.load(encodedBase+uint64(index*4), 4)
		if err != nil {
			return nil, nil, [2]uint32{}, fmt.Errorf("encoded %d: %w", index, err)
		}
		encoded[index] = uint32(value)
	}
	var correction [2]uint32
	for index := range correction {
		value, err := s.load(correctionBase+uint64(index*4), 4)
		if err != nil {
			return nil, nil, [2]uint32{}, fmt.Errorf("correction %d: %w", index, err)
		}
		correction[index] = uint32(value)
	}
	return companion, encoded, correction, nil
}

// runBootstrapShortWordStep decodes one of 227 words. Bit zero selects a
// two-word correction table; the remaining bits address the protected word
// stream and are combined with its companion stream.
func runBootstrapShortWordStep(s *state) (uint64, error) {
	if s == nil {
		return 0, fmt.Errorf("bootstrap short word state is unavailable")
	}
	s.traceBlock(bootstrapShortWordPass)
	base := s.x[21] + s.x[8]
	previousHighBit := uint32(s.x[10]) & 0x80000000
	encoded, err := s.load(base+4, 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap short encoded word: %w", err)
	}
	companion, err := s.load(base+0x634, 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap short companion word: %w", err)
	}
	correction, err := s.load(s.x[0]+(encoded&1)*4, 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap short correction word: %w", err)
	}
	folded := (uint32(encoded) & 0x7ffffffe) | previousHighBit
	decoded := uint32(companion) ^ (folded >> (uint32(s.x[12]) & 31)) ^ uint32(correction)
	if err := s.storeValue(base, uint64(decoded), 4); err != nil {
		return 0, fmt.Errorf("bootstrap short output word: %w", err)
	}

	s.x[10] = uint64(uint32(encoded))
	s.x[13] = base
	s.x[14] = uint64(decoded)
	s.x[15] = uint64(uint32(companion))
	s.x[16] = uint64(uint32(correction))
	s.compare64(s.x[8], 0x388)
	dispatchIndex := uint32(boolValue(s.condition("eq")))*uint32(s.x[1]) ^ uint32(s.x[9])
	next, err := loadBootstrapPointer(s, s.x[19], dispatchIndex)
	if err != nil {
		return 0, fmt.Errorf("bootstrap short dispatch: %w", err)
	}
	s.x[8] += 4
	s.x[13] = next
	s.x[17] = 0x67ab
	return next, nil
}

// runBootstrapLongWordStep decodes one of 396 words. It uses the same
// correction-table format as the short pass with a fixed one-bit fold.
func runBootstrapLongWordStep(s *state) (uint64, error) {
	if s == nil {
		return 0, fmt.Errorf("bootstrap long word state is unavailable")
	}
	s.traceBlock(bootstrapLongWordPass)
	output := s.x[11] + s.x[9]
	input := s.x[21] + s.x[9]
	previousHighBit := uint32(s.x[12]) & 0x80000000
	encoded, err := s.load(input+0x390, 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap long encoded word: %w", err)
	}
	companion, err := s.load(input, 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap long companion word: %w", err)
	}
	correction, err := s.load(s.x[0]+(encoded&1)*4, 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap long correction word: %w", err)
	}
	folded := (uint32(encoded) & 0x7ffffffe) | previousHighBit
	decoded := uint32(companion) ^ uint32(correction) ^ (folded >> 1)
	if err := s.storeValue(output, uint64(decoded), 4); err != nil {
		return 0, fmt.Errorf("bootstrap long output word: %w", err)
	}

	s.x[12] = uint64(uint32(encoded))
	s.x[13] = output
	s.x[14] = uint64(decoded)
	s.x[15] = uint64(uint32(companion) ^ uint32(correction))
	s.x[16] = uint64(uint32(correction))
	s.compare64(s.x[9], 0x62c)
	dispatchIndex := uint32(boolValue(s.condition("eq")))*uint32(s.x[1]) ^ uint32(s.x[10])
	next, err := loadBootstrapPointer(s, s.x[19], dispatchIndex)
	if err != nil {
		return 0, fmt.Errorf("bootstrap long dispatch: %w", err)
	}
	s.x[9] += 4
	s.x[13] = next
	s.x[17] = 0x1028
	return next, nil
}

// runBootstrapAccumulatorPass executes either bounded accumulator loop. Both
// native handlers return their own address while a word range remains and a
// continuation address when the range is complete.
func runBootstrapAccumulatorPass(s *state, pc uint64) (uint64, error) {
	const maximumWords = 128
	for word := 0; word < maximumWords; word++ {
		var next uint64
		var err error
		switch pc {
		case bootstrapAccumulatorPassA:
			next, err = runBootstrapAccumulatorStepA(s)
		case bootstrapAccumulatorPassB:
			next, err = runBootstrapAccumulatorStepB(s)
		default:
			return pc, fmt.Errorf("unknown bootstrap accumulator pass %#x", pc)
		}
		if err != nil || next != pc {
			return next, err
		}
	}
	return 0, fmt.Errorf("bootstrap accumulator pass %#x exceeded %d words", pc, maximumWords)
}

func runBootstrapAccumulatorStepA(s *state) (uint64, error) {
	if s == nil {
		return 0, fmt.Errorf("bootstrap accumulator A state is unavailable")
	}
	s.traceBlock(bootstrapAccumulatorPassA)
	index := uint32(s.x[11]) + uint32(s.x[14])
	left, err := loadBootstrapWord(s, s.x[9], index)
	if err != nil {
		return 0, fmt.Errorf("bootstrap accumulator A left table: %w", err)
	}
	masked := uint32(s.x[4]) & uint32(s.x[20])
	left = (left ^ s.x[16]) + (s.x[8] >> 32)
	right, err := loadBootstrapWord(s, s.x[27], index)
	if err != nil {
		return 0, fmt.Errorf("bootstrap accumulator A right table: %w", err)
	}
	coefficient := masked + uint32(s.x[19])
	right = uint64(uint32(right) ^ coefficient)
	s.x[8] = uint64(uint32(s.x[7]))*uint64(uint32(right)) + left

	accumulator := s.x[13] >> 32
	tableWord, err := loadBootstrapWord(s, s.x[26], index)
	if err != nil {
		return 0, fmt.Errorf("bootstrap accumulator A coefficient table: %w", err)
	}
	accumulator += uint64(uint32(tableWord^s.x[16]))*uint64(uint32(s.x[0])) + uint64(uint32(s.x[8]))
	fold := uint32(s.x[24]) & uint32(accumulator<<1)
	output := uint32(accumulator) - fold + uint32(s.x[16])
	storeIndex := uint32(s.x[11]) + uint32(s.x[12])
	if err := storeBootstrapWord(s, s.x[9], storeIndex, output); err != nil {
		return 0, fmt.Errorf("bootstrap accumulator A output table: %w", err)
	}

	nextCounter := uint32(s.x[11]) + 1
	branchCounter := uint32(s.x[11]) + uint32(s.x[28])
	s.compare32(branchCounter, uint32(s.x[25]))
	if s.condition("gt") {
		s.compare32(branchCounter, uint32(s.x[1]))
	} else {
		s.setNZCV(0)
	}
	dispatchIndex := uint32(boolValue(s.condition("lt")))*uint32(s.x[21]) ^ uint32(s.x[2])
	dispatchBase, err := s.load(s.sp+0xa8, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap accumulator A dispatch base: %w", err)
	}
	next, err := loadBootstrapPointer(s, dispatchBase, dispatchIndex)
	if err != nil {
		return 0, fmt.Errorf("bootstrap accumulator A dispatch: %w", err)
	}

	s.x[5] = uint64(nextCounter)
	s.x[6] = next
	s.x[11] = uint64(nextCounter)
	s.x[13] = accumulator
	s.x[15] = dispatchBase
	s.x[17] = 0x5bb3
	s.x[22] = uint64(masked)
	return next, nil
}

func runBootstrapAccumulatorStepB(s *state) (uint64, error) {
	if s == nil {
		return 0, fmt.Errorf("bootstrap accumulator B state is unavailable")
	}
	s.traceBlock(bootstrapAccumulatorPassB)
	s.compare32(uint32(s.x[27]), uint32(s.x[22]))
	gate := uint32(boolValue(s.condition("gt")))
	index := uint32(s.x[8]) + uint32(s.x[19])
	left, err := loadBootstrapWord(s, s.x[28], index)
	if err != nil {
		return 0, fmt.Errorf("bootstrap accumulator B left table: %w", err)
	}
	carry := s.x[23] ^ (s.x[1] >> 32)
	left = uint64(uint32(s.x[30]) ^ uint32(left))
	carry += left
	right, err := loadBootstrapWord(s, s.x[14], index)
	if err != nil {
		return 0, fmt.Errorf("bootstrap accumulator B right table: %w", err)
	}
	if s.x[24] != 0x19c0683c464360c5 || s.x[25] != 0xe63f97c3b9bc9f3b ||
		s.x[26] != 0x4042601c6f58ed3f || s.x[9] != 0xbfbd9fe390a712c1 ||
		s.x[16] != 0x598208200b552ec1 || s.x[0] != 0xa67df7dff4aad13f ||
		s.x[2] != 0x224ea33b || s.x[7] != 0xc0c02ce911b20ca0 ||
		s.x[21] != 0x818059d223641940 || s.x[3] != 0x33fcaf9b || s.x[15] > 0xffffffff {
		return 0, fmt.Errorf("unsupported bootstrap accumulator encoding")
	}
	decodedRight := uint32(right) ^ uint32(s.x[2])
	s.x[23] = 0xc0c02ce9
	combined := uint64(decodedRight)*uint64(uint32(s.x[15])) + carry
	result := combined ^ s.x[7]
	output := uint32(result) ^ uint32(s.x[3])
	storeIndex := uint32(s.x[8]) + uint32(s.x[4])
	if err := storeBootstrapWord(s, s.x[28], storeIndex, output); err != nil {
		return 0, fmt.Errorf("bootstrap accumulator B output table: %w", err)
	}

	nextCounter := uint32(s.x[8]) + 1
	branchCounter := uint32(s.x[8]) + uint32(s.x[5])
	s.compare32(branchCounter, uint32(s.x[6]))
	lessThanFirstBound := uint32(boolValue(s.condition("lt")))
	selector := gate ^ lessThanFirstBound
	s.compare32(branchCounter, uint32(s.x[13]))
	lessThanSecondBound := uint32(boolValue(s.condition("lt")))
	s.compare32(selector, 0)
	dispatchBit := lessThanSecondBound
	if s.condition("ne") {
		dispatchBit = gate
	}
	dispatchIndex := dispatchBit*uint32(s.x[12]) ^ uint32(s.x[11])
	dispatchBase, err := s.load(s.sp+0xa8, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap accumulator B dispatch base: %w", err)
	}
	next, err := loadBootstrapPointer(s, dispatchBase, dispatchIndex)
	if err != nil {
		return 0, fmt.Errorf("bootstrap accumulator B dispatch: %w", err)
	}

	s.x[1] = result
	s.x[8] = uint64(nextCounter)
	s.x[10] = uint64(nextCounter)
	s.x[17] = 0x207f
	s.x[20] = next
	return next, nil
}

// runBootstrapShufflePass is a descending Fisher-Yates-style permutation of
// 64-bit table entries. The native code expresses each XOR swap as a handler
// that dispatches back to itself.
func runBootstrapShufflePass(s *state) (uint64, error) {
	const maximumEntries = 4096
	for entry := 0; entry < maximumEntries; entry++ {
		next, err := runBootstrapShuffleStep(s)
		if err != nil || next != bootstrapShufflePass {
			return next, err
		}
	}
	return 0, fmt.Errorf("bootstrap shuffle exceeded %d entries", maximumEntries)
}

func runBootstrapShuffleStep(s *state) (uint64, error) {
	if s == nil {
		return 0, fmt.Errorf("bootstrap shuffle state is unavailable")
	}
	s.traceBlock(bootstrapShufflePass)
	limit := uint32(s.x[11])
	seed := uint32(s.x[17])
	mixed := (seed + limit) ^ seed
	mixed *= uint32(s.x[0])
	selected := uint32(0)
	if limit != 0 {
		selected = mixed % limit
	}
	leftAddress := s.x[16] + uint64(limit)*8
	rightAddress := s.x[16] + uint64(selected)*8
	left, err := s.load(leftAddress, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap shuffle left entry: %w", err)
	}
	right, err := s.load(rightAddress, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap shuffle right entry: %w", err)
	}
	left ^= right
	if err := s.storeValue(leftAddress, left, 8); err != nil {
		return 0, fmt.Errorf("bootstrap shuffle left xor: %w", err)
	}
	rightLoaded, err := s.load(rightAddress, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap shuffle updated right entry: %w", err)
	}
	right = rightLoaded ^ left
	if err := s.storeValue(rightAddress, right, 8); err != nil {
		return 0, fmt.Errorf("bootstrap shuffle right xor: %w", err)
	}
	left, err = s.load(leftAddress, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap shuffle updated left entry: %w", err)
	}
	left ^= right
	if err := s.storeValue(leftAddress, left, 8); err != nil {
		return 0, fmt.Errorf("bootstrap shuffle final left entry: %w", err)
	}

	s.x[1] = uint64(limit) + 1
	s.x[2] = leftAddress
	s.x[3] = right
	s.x[4] = rightLoaded
	s.x[11] = uint64(limit) - 1
	s.compare64(s.x[1], 2)
	if s.condition("gt") {
		return bootstrapShufflePass, nil
	}
	return bootstrapShuffleExit, nil
}

// runBootstrapForwardShufflePass performs the ascending companion to the
// descending shuffle above. The first entry is prepared by the preceding CFG;
// this loop owns entries 2 through 272.
func runBootstrapForwardShufflePass(s *state) (uint64, error) {
	const maximumEntries = 4096
	for entry := 0; entry < maximumEntries; entry++ {
		next, err := runBootstrapForwardShuffleIteration(s)
		if err != nil || next != bootstrapForwardShuffleSeed {
			return next, err
		}
	}
	return 0, fmt.Errorf("bootstrap forward shuffle exceeded %d entries", maximumEntries)
}

func runBootstrapForwardShuffleIteration(s *state) (uint64, error) {
	if s == nil {
		return 0, fmt.Errorf("bootstrap forward shuffle state is unavailable")
	}
	s.traceBlock(bootstrapForwardShuffleSeed)
	limit := uint32(s.x[12])
	seed := uint32(s.x[10])
	mixed := ((seed + limit) ^ seed) * uint32(s.x[11])

	s.traceBlock(bootstrapForwardShuffleStep)
	selected := uint32(0)
	if limit != 0 {
		selected = mixed % limit
	}
	leftAddress := s.x[13] + uint64(limit)*8
	rightAddress := s.x[13] + uint64(selected)*8
	left, err := s.load(leftAddress, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap forward shuffle left entry: %w", err)
	}
	rightLoaded, err := s.load(rightAddress, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap forward shuffle right entry: %w", err)
	}
	xor := left ^ rightLoaded
	if err := s.storeValue(leftAddress, xor, 8); err != nil {
		return 0, fmt.Errorf("bootstrap forward shuffle left xor: %w", err)
	}
	rightLoaded, err = s.load(rightAddress, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap forward shuffle updated right entry: %w", err)
	}
	right := rightLoaded ^ xor
	if err := s.storeValue(rightAddress, right, 8); err != nil {
		return 0, fmt.Errorf("bootstrap forward shuffle right xor: %w", err)
	}
	left, err = s.load(leftAddress, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap forward shuffle updated left entry: %w", err)
	}
	left ^= right
	if err := s.storeValue(leftAddress, left, 8); err != nil {
		return 0, fmt.Errorf("bootstrap forward shuffle final left entry: %w", err)
	}

	s.x[12] = uint64(limit) + 1
	s.x[14] = left
	s.x[15] = leftAddress
	s.x[16] = right
	s.x[17] = rightLoaded
	s.compare64(s.x[12], 0x111)
	if s.condition("ne") {
		return bootstrapForwardShuffleSeed, nil
	}
	return bootstrapForwardShuffleExit, nil
}

// runBootstrapReverseBytePass copies a byte range backwards while applying a
// fixed XOR byte. The protected handler uses its zero flag to select the exit
// continuation after the last byte.
func runBootstrapReverseBytePass(s *state) (uint64, error) {
	const maximumBytes = 4096
	for index := 0; index < maximumBytes; index++ {
		next, err := runBootstrapReverseByteStep(s)
		if err != nil || next != bootstrapReverseBytePass {
			return next, err
		}
	}
	return 0, fmt.Errorf("bootstrap reverse byte pass exceeded %d bytes", maximumBytes)
}

func runBootstrapReverseByteStep(s *state) (uint64, error) {
	if s == nil {
		return 0, fmt.Errorf("bootstrap reverse byte state is unavailable")
	}
	s.traceBlock(bootstrapReverseBytePass)
	s.x[8] = uint64(uint32(s.x[9]) ^ uint32(s.x[11]))
	counter := s.x[10] - 1
	s.x[10] = counter
	s.compare64(counter+1, 1)
	value, err := s.load(s.x[24]+counter, 1)
	if err != nil {
		return 0, fmt.Errorf("bootstrap reverse byte input: %w", err)
	}
	value = uint64(uint32(value) ^ uint32(s.x[12]))
	if err := s.storeValue(s.x[13]+counter, uint64(byte(value)), 1); err != nil {
		return 0, fmt.Errorf("bootstrap reverse byte output: %w", err)
	}
	last := uint32(boolValue(s.condition("eq")))
	s.x[16] = uint64(uint32(s.x[8]) + uint32(s.x[14]))
	dispatchIndex := last*uint32(s.x[16]) ^ uint32(s.x[9])
	next, err := loadBootstrapPointer(s, s.x[21], dispatchIndex)
	if err != nil {
		return 0, fmt.Errorf("bootstrap reverse byte dispatch: %w", err)
	}
	s.x[15] = next
	s.x[17] = 0x864c
	return next, nil
}

// runBootstrapLookupCopyPass resolves a source table through two byte lookup
// layers, then copies one 32-bit word per iteration to the output table.
func runBootstrapLookupCopyPass(s *state) (uint64, error) {
	const maximumWords = 4096
	for word := 0; word < maximumWords; word++ {
		next, err := runBootstrapLookupCopyIteration(s)
		if err != nil || next != bootstrapLookupCopySeed {
			return next, err
		}
	}
	return 0, fmt.Errorf("bootstrap lookup-copy pass exceeded %d words", maximumWords)
}

func runBootstrapLookupCopyIteration(s *state) (uint64, error) {
	if s == nil {
		return 0, fmt.Errorf("bootstrap lookup-copy state is unavailable")
	}
	s.traceBlock(bootstrapLookupCopySeed)
	leftPointer, err := s.load(s.sp+0xde8, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap lookup-copy left pointer: %w", err)
	}
	left, err := s.load(leftPointer, 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap lookup-copy left value: %w", err)
	}
	rightPointer, err := s.load(s.sp+0xde0, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap lookup-copy right pointer: %w", err)
	}
	right, err := s.load(rightPointer, 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap lookup-copy right value: %w", err)
	}
	selector := (uint32(right) - uint32(left)) ^ uint32(s.x[14])
	selector *= uint32(s.x[15])
	first, err := s.load(s.x[4]+uint64(byte(selector)), 1)
	if err != nil {
		return 0, fmt.Errorf("bootstrap lookup-copy first table: %w", err)
	}

	s.traceBlock(bootstrapLookupCopyStep)
	secondIndex := uint32(first) ^ uint32(s.x[16])
	second, err := s.load(s.x[3]+uint64(secondIndex), 1)
	if err != nil {
		return 0, fmt.Errorf("bootstrap lookup-copy second table: %w", err)
	}
	sourceIndex := (selector ^ uint32(second)) & 0xff
	sourceIndex += 0x1b9
	sourceTable, err := s.load(s.x[5]+uint64(sourceIndex)*8, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap lookup-copy source table: %w", err)
	}
	position := s.addWithCarry(s.x[13], s.x[12], false, 64)
	source := sourceTable + s.x[0]
	value, err := s.load(source+shiftLeft(position, 2, 64), 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap lookup-copy source word: %w", err)
	}
	if err := s.storeValue(s.x[11]+shiftLeft(position, 2, 64), uint64(uint32(value)), 4); err != nil {
		return 0, fmt.Errorf("bootstrap lookup-copy output word: %w", err)
	}
	last := uint32(boolValue(s.condition("eq")))
	dispatchIndex := last*uint32(s.x[1]) ^ uint32(s.x[9])
	dispatchTable, err := s.load(s.sp+0xdf0, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap lookup-copy dispatch table: %w", err)
	}
	next, err := loadBootstrapPointer(s, dispatchTable, dispatchIndex)
	if err != nil {
		return 0, fmt.Errorf("bootstrap lookup-copy dispatch: %w", err)
	}
	s.x[2] = next
	s.x[12] = position
	s.x[17] = 0x8085
	return next, nil
}

// runBootstrapPointerWalkPass advances one encoded counter and uses a byte
// classifier to select the next protected continuation.
func runBootstrapPointerWalkPass(s *state) (uint64, error) {
	const maximumSteps = 4096
	for step := 0; step < maximumSteps; step++ {
		next, err := runBootstrapPointerWalkIteration(s)
		if err != nil || next != bootstrapPointerWalkSeed {
			return next, err
		}
	}
	return 0, fmt.Errorf("bootstrap pointer walk exceeded %d steps", maximumSteps)
}

func runBootstrapPointerWalkIteration(s *state) (uint64, error) {
	if s == nil {
		return 0, fmt.Errorf("bootstrap pointer-walk state is unavailable")
	}
	s.traceBlock(bootstrapPointerWalkSeed)
	const (
		counterBias = uint32(0xffff65f4)
		counterStep = uint32(0xffffe33a)
		mask        = uint64(0x6f164ec4)
		pointerSalt = uint64(0xf22d1348378b2762)
	)
	counter := uint32(s.x[10])
	s.x[12] = uint64(uint32(s.x[16]) + counterBias)
	s.x[13] = uint64(counter)
	removed := mask & (uint64(counter) << 1)
	s.x[11] = pointerSalt
	s.x[13] = s.x[13] - removed + pointerSalt
	lookupOffset := s.x[13] ^ pointerSalt
	classification, err := s.load(s.x[24]+lookupOffset, 1)
	if err != nil {
		return 0, fmt.Errorf("bootstrap pointer-walk classifier: %w", err)
	}
	s.compare32(uint32(classification), 0x1b)
	classified := uint32(boolValue(s.condition("ne")))
	s.x[15] = uint64(counterStep)
	s.x[10] = uint64(uint32(s.x[12]) + uint32(counter+counterStep))
	classified |= classified << 2
	dispatchIndex := classified ^ uint32(s.x[16])
	next, err := loadBootstrapPointer(s, s.x[21], dispatchIndex)
	if err != nil {
		return 0, fmt.Errorf("bootstrap pointer-walk dispatch: %w", err)
	}
	s.x[14] = next
	s.x[17] = 0x5724
	s.traceBlock(bootstrapPointerWalkDispatch)
	return next, nil
}

// runBootstrapPairCopyPass copies two 32-bit words per iteration from the
// source object at [x8+8] to the destination object at [x9+8].
func runBootstrapPairCopyPass(s *state) (uint64, error) {
	if s == nil {
		return 0, fmt.Errorf("bootstrap pair-copy state is unavailable")
	}
	const maximumPairs = 8192
	counter := uint32(s.x[13])
	if counter < 2 || counter&1 != 0 || counter/2 > maximumPairs {
		return 0, fmt.Errorf("unsupported bootstrap pair-copy length %d", counter)
	}
	source, err := s.load(s.x[8]+8, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap pair-copy source pointer: %w", err)
	}
	destination, err := s.load(s.x[9]+8, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap pair-copy destination pointer: %w", err)
	}
	indexBias := (uint32(s.x[11]) ^ 0xeb62) + 0xffff8181
	for current := counter; current != 0; current -= 2 {
		firstIndex := indexBias + current
		first, err := loadBootstrapWord(s, source, firstIndex)
		if err != nil {
			return 0, fmt.Errorf("bootstrap pair-copy first source: %w", err)
		}
		if err := storeBootstrapWord(s, destination, firstIndex, uint32(first)); err != nil {
			return 0, fmt.Errorf("bootstrap pair-copy first destination: %w", err)
		}
		secondIndex := current - 2
		second, err := loadBootstrapWord(s, source, secondIndex)
		if err != nil {
			return 0, fmt.Errorf("bootstrap pair-copy second source: %w", err)
		}
		if err := storeBootstrapWord(s, destination, secondIndex, uint32(second)); err != nil {
			return 0, fmt.Errorf("bootstrap pair-copy second destination: %w", err)
		}
	}
	next, err := loadBootstrapPointer(s, s.x[10], uint32(s.x[11]))
	if err != nil {
		return 0, fmt.Errorf("bootstrap pair-copy dispatch: %w", err)
	}
	s.x[12] = next
	s.x[13] = 0
	s.x[14] = 0
	s.x[15] = destination
	s.x[17] = 0xa668
	s.compare32(2, 2)
	return next, nil
}

func runBootstrapPairCopyStep(s *state) (uint64, error) {
	if s == nil {
		return 0, fmt.Errorf("bootstrap pair-copy state is unavailable")
	}
	s.traceBlock(bootstrapPairCopyPass)
	firstIndex := uint32(s.x[11]) ^ 0xeb62
	firstIndex += uint32(s.x[13]) + 0xffff8181
	source, err := s.load(s.x[8]+8, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap pair-copy source pointer: %w", err)
	}
	first, err := loadBootstrapWord(s, source, firstIndex)
	if err != nil {
		return 0, fmt.Errorf("bootstrap pair-copy first source: %w", err)
	}
	destination, err := s.load(s.x[9]+8, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap pair-copy destination pointer: %w", err)
	}
	if err := storeBootstrapWord(s, destination, firstIndex, uint32(first)); err != nil {
		return 0, fmt.Errorf("bootstrap pair-copy first destination: %w", err)
	}

	s.x[13] = uint64(uint32(s.x[13]) - 2)
	s.compare32(uint32(s.x[13])+2, 2)
	second, err := loadBootstrapWord(s, source, uint32(s.x[13]))
	if err != nil {
		return 0, fmt.Errorf("bootstrap pair-copy second source: %w", err)
	}
	if err := storeBootstrapWord(s, destination, uint32(s.x[13]), uint32(second)); err != nil {
		return 0, fmt.Errorf("bootstrap pair-copy second destination: %w", err)
	}
	more := uint32(boolValue(s.condition("ne")))
	dispatchIndex := ((more << 2) | (more << 5)) ^ uint32(s.x[11])
	next, err := loadBootstrapPointer(s, s.x[10], dispatchIndex)
	if err != nil {
		return 0, fmt.Errorf("bootstrap pair-copy dispatch: %w", err)
	}
	s.x[12] = next
	s.x[14] = uint64(more << 2)
	s.x[15] = destination
	s.x[17] = 0xa668
	return next, nil
}

// runBootstrapConstantFillPass writes one encoded constant per iteration to a
// sparse stack table, then selects the next continuation from two bounds.
func runBootstrapConstantFillPass(s *state) (uint64, error) {
	const maximumWords = 8192
	if s == nil {
		return 0, fmt.Errorf("bootstrap constant-fill state is unavailable")
	}
	targets := [2]uint64{}
	for bit := uint32(0); bit < 2; bit++ {
		next, err := loadBootstrapPointer(s, s.x[25], uint32(s.x[9])^bit*0x1d4)
		if err != nil {
			return 0, fmt.Errorf("bootstrap constant-fill target %d: %w", bit, err)
		}
		targets[bit] = next
	}
	if targets[0] == bootstrapConstantFillSeed && targets[1] == bootstrapConstantFillSeed {
		return 0, fmt.Errorf("bootstrap constant-fill has no exit target")
	}

	for word := 0; word < maximumWords; word++ {
		s.traceBlock(bootstrapConstantFillSeed)
		s.x[11] = 0x77b16a0c
		s.compare32(uint32(s.x[8]), uint32(s.x[11]))
		firstBound := uint32(boolValue(s.condition("gt")))
		s.x[11] = uint64((uint32(s.x[9]) + 0xffffe51b) | 0x4561)
		counter := uint32(s.x[15])
		s.x[16] = uint64(counter + 0xde3d1312 + uint32(s.x[11]))
		s.x[17] = s.sp + 0x2d0
		s.x[1] = 0x224ea33b
		if err := storeBootstrapWord(s, s.x[17], uint32(s.x[16]), uint32(s.x[1])); err != nil {
			return 0, fmt.Errorf("bootstrap constant-fill output: %w", err)
		}
		s.x[16] = uint64(firstBound)
		s.x[17] = uint64(counter + 1)
		s.x[1] = 0xe68c1e77
		s.x[2] = 0x884e95f3
		s.x[15] = uint64(counter + uint32(s.x[1]))
		s.compare32(uint32(s.x[15]), uint32(s.x[2]))

		s.traceBlock(bootstrapConstantFillDispatch)
		secondBound := uint32(boolValue(s.condition("lt")))
		s.x[16] = uint64(firstBound ^ secondBound)
		s.compare32(uint32(s.x[15]), uint32(s.x[14]))
		thirdBound := uint32(boolValue(s.condition("gt")))
		s.compare32(uint32(s.x[16]), 0)
		dispatchBit := thirdBound
		if s.condition("ne") {
			dispatchBit = secondBound
		}
		s.x[16] = 0x1d4
		next := targets[dispatchBit]
		s.x[1] = next
		s.x[15] = s.x[17]
		s.x[17] = 0x88c4
		if next != bootstrapConstantFillSeed {
			return next, nil
		}
	}
	return 0, fmt.Errorf("bootstrap constant-fill pass exceeded %d words", maximumWords)
}

func runBootstrapConstantFillIteration(s *state) (uint64, error) {
	if s == nil {
		return 0, fmt.Errorf("bootstrap constant-fill state is unavailable")
	}
	s.traceBlock(bootstrapConstantFillSeed)
	s.x[11] = 0x77b16a0c
	s.compare32(uint32(s.x[8]), uint32(s.x[11]))
	firstBound := uint32(boolValue(s.condition("gt")))
	s.x[11] = uint64(uint32(s.x[9]) + 0xffffe51b)
	s.x[16] = 0x4561
	s.x[11] = uint64(uint32(s.x[11]) | uint32(s.x[16]))
	s.x[16] = uint64(uint32(s.x[15]) + 0xde3d1312)
	s.x[16] = uint64(uint32(s.x[16]) + uint32(s.x[11]))
	s.x[17] = s.sp + 0x2d0
	s.x[1] = 0x224ea33b
	if err := storeBootstrapWord(s, s.x[17], uint32(s.x[16]), uint32(s.x[1])); err != nil {
		return 0, fmt.Errorf("bootstrap constant-fill output: %w", err)
	}
	s.x[16] = uint64(firstBound)
	s.x[17] = uint64(uint32(s.x[15]) + 1)
	s.x[1] = 0xe68c1e77
	s.x[2] = 0x884e95f3
	s.x[15] = uint64(uint32(s.x[15]) + uint32(s.x[1]))
	s.compare32(uint32(s.x[15]), uint32(s.x[2]))

	s.traceBlock(bootstrapConstantFillDispatch)
	secondBound := uint32(boolValue(s.condition("lt")))
	s.x[16] = uint64(uint32(s.x[16]) ^ secondBound)
	s.compare32(uint32(s.x[15]), uint32(s.x[14]))
	thirdBound := uint32(boolValue(s.condition("gt")))
	s.compare32(uint32(s.x[16]), 0)
	dispatchBit := thirdBound
	if s.condition("ne") {
		dispatchBit = secondBound
	}
	s.x[16] = 0x1d4
	dispatchIndex := dispatchBit*uint32(s.x[16]) ^ uint32(s.x[9])
	next, err := loadBootstrapPointer(s, s.x[25], dispatchIndex)
	if err != nil {
		return 0, fmt.Errorf("bootstrap constant-fill dispatch: %w", err)
	}
	s.x[1] = next
	s.x[15] = s.x[17]
	s.x[17] = 0x88c4
	return next, nil
}

func runBootstrapPackedBitPass(s *state) (uint64, error) {
	const maximumWords = 8192
	if s == nil {
		return 0, fmt.Errorf("bootstrap packed-bit state is unavailable")
	}
	targets := [2]uint64{}
	for bit := uint32(0); bit < 2; bit++ {
		next, err := loadBootstrapPointer(s, s.x[20], uint32(s.x[14])^bit*uint32(s.x[7]))
		if err != nil {
			return 0, fmt.Errorf("bootstrap packed-bit target %d: %w", bit, err)
		}
		targets[bit] = next
	}
	if targets[0] == bootstrapPackedBitPass && targets[1] == bootstrapPackedBitPass {
		return 0, fmt.Errorf("bootstrap packed-bit pass has no exit target")
	}
	table, err := s.load(s.x[19]+8, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap packed-bit table: %w", err)
	}

	for word := 0; word < maximumWords; word++ {
		s.traceBlock(bootstrapPackedBitPass)
		s.x[11] = uint64(uint32(s.x[1]) | uint32(s.x[2]))
		index := (s.x[11] ^ s.x[3]) + s.x[10]
		s.x[22] = s.x[10] - 1
		s.x[23] = table
		value, err := s.load(table+shiftLeft(index, 2, 64), 4)
		if err != nil {
			return 0, fmt.Errorf("bootstrap packed-bit input: %w", err)
		}
		s.x[17] = uint64(uint32(value) ^ uint32(s.x[12]))
		s.x[24] = uint64(uint32(shiftRight(s.x[17], int(s.x[13]), 32)))
		fragment := uint32(s.x[24])
		current := uint32(s.x[21])
		s.x[21] = uint64(current ^ (fragment & (current ^ uint32(s.x[4]) ^ uint32(s.x[5]))))
		s.x[24] = uint64(fragment & uint32(s.x[5]))
		if err := storeBootstrapWord(s, table, uint32(s.x[0]), uint32(s.x[21])); err != nil {
			return 0, fmt.Errorf("bootstrap packed-bit output: %w", err)
		}
		s.x[17] = uint64(uint32(shiftLeft(s.x[17], int(s.x[16]), 32)))
		s.x[17] = uint64(uint32(shiftLeft(s.x[17], int(s.x[15]), 32)))
		s.x[21] = uint64(uint32(s.x[17]) ^ uint32(s.x[12]))
		s.x[0] = uint64(uint32(s.x[0]) - 1)
		dispatchBit := (^uint32(uint32(s.x[10]) + uint32(s.x[6]))) >> 31
		next := targets[dispatchBit]
		s.x[10] = s.x[22]
		s.x[17] = 0x390b
		s.x[23] = next
		if next != bootstrapPackedBitPass {
			return next, nil
		}
	}
	return 0, fmt.Errorf("bootstrap packed-bit pass exceeded %d words", maximumWords)
}

func runBootstrapPackedBitStep(s *state) (uint64, error) {
	if s == nil {
		return 0, fmt.Errorf("bootstrap packed-bit state is unavailable")
	}
	s.traceBlock(bootstrapPackedBitPass)
	s.x[11] = uint64(uint32(s.x[1]) | uint32(s.x[2]))
	index := s.x[11] ^ s.x[3]
	index += s.x[10]
	s.x[22] = s.x[10] - 1
	table, err := s.load(s.x[19]+8, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap packed-bit table: %w", err)
	}
	s.x[23] = table
	word, err := s.load(table+shiftLeft(index, 2, 64), 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap packed-bit input: %w", err)
	}
	s.x[17] = uint64(uint32(word) ^ uint32(s.x[12]))
	s.x[24] = uint64(uint32(shiftRight(s.x[17], int(s.x[13]), 32)))
	s.x[21] = uint64(uint32(s.x[21]) ^ uint32(s.x[4]))
	s.x[21] = uint64(uint32(s.x[21]) &^ uint32(s.x[24]))
	s.x[24] = uint64(uint32(s.x[24]) & uint32(s.x[5]))
	s.x[21] = uint64(uint32(s.x[24]) | uint32(s.x[21]))
	s.x[21] = uint64(uint32(s.x[21]) ^ uint32(s.x[4]))
	if err := storeBootstrapWord(s, table, uint32(s.x[0]), uint32(s.x[21])); err != nil {
		return 0, fmt.Errorf("bootstrap packed-bit output: %w", err)
	}
	s.x[17] = uint64(uint32(shiftLeft(s.x[17], int(s.x[16]), 32)))
	s.x[17] = uint64(uint32(shiftLeft(s.x[17], int(s.x[15]), 32)))
	s.x[21] = uint64(uint32(s.x[17]) ^ uint32(s.x[12]))
	s.x[0] = uint64(uint32(s.x[0]) - 1)
	s.x[10] = uint64(uint32(s.x[10]) + uint32(s.x[6]))
	s.x[10] = uint64(^uint32(s.x[10]))
	s.x[10] = uint64(uint32(s.x[10]) >> 31)
	dispatchIndex := uint32(s.x[10])*uint32(s.x[7]) ^ uint32(s.x[14])
	next, err := loadBootstrapPointer(s, s.x[20], dispatchIndex)
	if err != nil {
		return 0, fmt.Errorf("bootstrap packed-bit dispatch: %w", err)
	}
	s.x[10] = s.x[22]
	s.x[17] = 0x390b
	s.x[23] = next
	return next, nil
}

func runBootstrapTableAccumulatorPass(s *state) (uint64, error) {
	const maximumWords = 8192
	if s == nil {
		return 0, fmt.Errorf("bootstrap table-accumulator state is unavailable")
	}
	targets := [2]uint64{}
	for bit := uint32(0); bit < 2; bit++ {
		index := uint32(s.x[13]) ^ bit*uint32(s.x[6])
		next, err := s.load(s.x[20]+uint64(index)*8, 8)
		if err != nil {
			return 0, fmt.Errorf("bootstrap table-accumulator target %d: %w", bit, err)
		}
		targets[bit] = next
	}
	if targets[0] == bootstrapTableAccumulatorPass && targets[1] == bootstrapTableAccumulatorPass {
		return 0, fmt.Errorf("bootstrap table accumulator has no exit target")
	}
	leftTable, err := s.load(s.x[23]+8, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap table accumulator left pointer: %w", err)
	}
	rightTable, err := s.load(s.x[21]+8, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap table accumulator right pointer: %w", err)
	}
	destination, err := s.load(s.x[9]+8, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap table accumulator destination: %w", err)
	}

	for word := 0; word < maximumWords; word++ {
		s.traceBlock(bootstrapTableAccumulatorPass)
		s.compare32(uint32(s.x[12]), uint32(s.x[15]))
		gate := uint32(boolValue(s.condition("lt")))
		index := uint32(s.x[7]) + uint32(s.x[0])
		left, err := loadBootstrapWord(s, leftTable, index)
		if err != nil {
			return 0, fmt.Errorf("bootstrap table accumulator left word: %w", err)
		}
		right, err := loadBootstrapWord(s, rightTable, index)
		if err != nil {
			return 0, fmt.Errorf("bootstrap table accumulator right word: %w", err)
		}
		accumulator := s.x[11] + (left ^ s.x[1]) - (right ^ s.x[1]) + s.x[2]
		s.x[14] = uint64(uint32(s.x[13]) + uint32(s.x[16]))
		s.x[11] = uint64(uint32(accumulator) + (uint32(s.x[14]) ^ uint32(s.x[4])))
		s.x[25] = uint64(uint32(s.x[3]) & uint32(accumulator<<1))
		s.x[11] = uint64(uint32(s.x[11]) - uint32(s.x[25]))
		s.x[25] = destination
		if err := storeBootstrapWord(s, destination, index, uint32(s.x[11])); err != nil {
			return 0, fmt.Errorf("bootstrap table accumulator output: %w", err)
		}
		s.x[11] = accumulator >> 32
		nextCounter := uint32(s.x[7]) + 1
		branchCounter := uint32(s.x[7]) + uint32(s.x[5])
		s.compare32(branchCounter, uint32(s.x[15]))
		firstBound := uint32(boolValue(s.condition("lt")))
		selector := gate ^ firstBound
		s.compare32(branchCounter, uint32(s.x[12]))
		secondBound := uint32(boolValue(s.condition("lt")))
		s.compare32(selector, 0)
		dispatchBit := secondBound
		if s.condition("ne") {
			dispatchBit = gate
		}
		dispatchBit ^= 1
		next := targets[dispatchBit]
		s.x[7] = uint64(nextCounter)
		s.x[10] = uint64(nextCounter)
		s.x[17] = 0x4c2
		s.x[19] = next
		if next != bootstrapTableAccumulatorPass {
			return next, nil
		}
	}
	return 0, fmt.Errorf("bootstrap table accumulator exceeded %d words", maximumWords)
}

func runBootstrapTableAccumulatorStep(s *state) (uint64, error) {
	if s == nil {
		return 0, fmt.Errorf("bootstrap table-accumulator state is unavailable")
	}
	s.traceBlock(bootstrapTableAccumulatorPass)
	s.compare32(uint32(s.x[12]), uint32(s.x[15]))
	gate := uint32(boolValue(s.condition("lt")))
	leftTable, err := s.load(s.x[23]+8, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap table accumulator left pointer: %w", err)
	}
	index := uint32(s.x[7]) + uint32(s.x[0])
	left, err := loadBootstrapWord(s, leftTable, index)
	if err != nil {
		return 0, fmt.Errorf("bootstrap table accumulator left word: %w", err)
	}
	left ^= s.x[1]
	rightTable, err := s.load(s.x[21]+8, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap table accumulator right pointer: %w", err)
	}
	right, err := loadBootstrapWord(s, rightTable, index)
	if err != nil {
		return 0, fmt.Errorf("bootstrap table accumulator right word: %w", err)
	}
	right ^= s.x[1]
	accumulator := s.x[11] + left - right
	s.x[14] = uint64(uint32(s.x[13]) + uint32(s.x[16]))
	s.x[11] = uint64(uint32(s.x[14]) ^ uint32(s.x[4]))
	accumulator += s.x[2]
	s.x[11] = uint64(uint32(accumulator) + uint32(s.x[11]))
	s.x[25] = uint64(uint32(s.x[3]) & uint32(accumulator<<1))
	s.x[11] = uint64(uint32(s.x[11]) - uint32(s.x[25]))
	destination, err := s.load(s.x[9]+8, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap table accumulator destination: %w", err)
	}
	s.x[25] = destination
	if err := storeBootstrapWord(s, destination, index, uint32(s.x[11])); err != nil {
		return 0, fmt.Errorf("bootstrap table accumulator output: %w", err)
	}
	s.x[11] = accumulator >> 32
	nextCounter := uint32(s.x[7]) + 1
	branchCounter := uint32(s.x[7]) + uint32(s.x[5])
	s.compare32(branchCounter, uint32(s.x[15]))
	firstBound := uint32(boolValue(s.condition("lt")))
	selector := gate ^ firstBound
	s.compare32(branchCounter, uint32(s.x[12]))
	secondBound := uint32(boolValue(s.condition("lt")))
	s.compare32(selector, 0)
	dispatchBit := secondBound
	if s.condition("ne") {
		dispatchBit = gate
	}
	dispatchBit ^= 1
	dispatchIndex := dispatchBit*uint32(s.x[6]) ^ uint32(s.x[13])
	next, err := s.load(s.x[20]+uint64(dispatchIndex)*8, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap table accumulator dispatch: %w", err)
	}
	s.x[7] = uint64(nextCounter)
	s.x[10] = uint64(nextCounter)
	s.x[17] = 0x4c2
	s.x[19] = next
	return next, nil
}

func runBootstrapSparseFillPass(s *state) (uint64, error) {
	const maximumWords = 4096
	if s == nil {
		return 0, fmt.Errorf("bootstrap sparse-fill state is unavailable")
	}
	dispatchTable, err := s.load(s.sp+0xa8, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap sparse-fill dispatch table: %w", err)
	}
	targets := [2]uint64{}
	for bit := uint32(0); bit < 2; bit++ {
		next, err := loadBootstrapPointer(s, dispatchTable, uint32(s.x[16])^bit*uint32(s.x[2]))
		if err != nil {
			return 0, fmt.Errorf("bootstrap sparse-fill target %d: %w", bit, err)
		}
		targets[bit] = next
	}
	if targets[0] == bootstrapSparseFillSeed && targets[1] == bootstrapSparseFillSeed {
		return 0, fmt.Errorf("bootstrap sparse-fill has no exit target")
	}

	for word := 0; word < maximumWords; word++ {
		s.traceBlock(bootstrapSparseFillSeed)
		s.x[8] = uint64((uint32(s.x[16]) ^ 0x2662) * 0x3d)
		s.x[17] = uint64((uint32(s.x[8]) ^ 0xa2c3baf5) + uint32(s.x[3]))
		s.x[4] = s.sp + 0xd8
		s.x[5] = 0x224ea33b
		if err := storeBootstrapWord(s, s.x[4], uint32(s.x[17]), uint32(s.x[5])); err != nil {
			return 0, fmt.Errorf("bootstrap sparse-fill output: %w", err)
		}
		s.x[17] = uint64(uint32(s.x[3]) + 1)
		s.x[4] = 0x7b932eb7
		s.x[5] = 0x58cf47f3
		s.x[3] = uint64(uint32(s.x[3]) + uint32(s.x[4]))
		s.compare32(uint32(s.x[3]), uint32(s.x[5]))
		if s.condition("le") {
			s.compare32(uint32(s.x[3]), uint32(s.x[1]))
		} else {
			s.setNZCV(4)
		}
		dispatchBit := uint32(boolValue(s.condition("le")))

		s.traceBlock(bootstrapSparseFillDispatch)
		next := targets[dispatchBit]
		s.x[3] = s.x[17]
		s.x[4] = next
		s.x[17] = 0x6d1a
		if next != bootstrapSparseFillSeed {
			return next, nil
		}
	}
	return 0, fmt.Errorf("bootstrap sparse-fill pass exceeded %d words", maximumWords)
}

func runBootstrapSparseFillIteration(s *state) (uint64, error) {
	if s == nil {
		return 0, fmt.Errorf("bootstrap sparse-fill state is unavailable")
	}
	s.traceBlock(bootstrapSparseFillSeed)
	s.x[8] = uint64((uint32(s.x[16]) ^ 0x2662) * 0x3d)
	s.x[17] = uint64(uint32(s.x[8]) ^ 0xa2c3baf5)
	s.x[17] = uint64(uint32(s.x[17]) + uint32(s.x[3]))
	s.x[4] = s.sp + 0xd8
	s.x[5] = 0x224ea33b
	if err := storeBootstrapWord(s, s.x[4], uint32(s.x[17]), uint32(s.x[5])); err != nil {
		return 0, fmt.Errorf("bootstrap sparse-fill output: %w", err)
	}
	s.x[17] = uint64(uint32(s.x[3]) + 1)
	s.x[4] = 0x7b932eb7
	s.x[5] = 0x58cf47f3
	s.x[3] = uint64(uint32(s.x[3]) + uint32(s.x[4]))
	s.compare32(uint32(s.x[3]), uint32(s.x[5]))
	if s.condition("le") {
		s.compare32(uint32(s.x[3]), uint32(s.x[1]))
	} else {
		s.setNZCV(4)
	}
	dispatchBit := uint32(boolValue(s.condition("le")))

	s.traceBlock(bootstrapSparseFillDispatch)
	dispatchIndex := dispatchBit*uint32(s.x[2]) ^ uint32(s.x[16])
	dispatchTable, err := s.load(s.sp+0xa8, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap sparse-fill dispatch table: %w", err)
	}
	next, err := loadBootstrapPointer(s, dispatchTable, dispatchIndex)
	if err != nil {
		return 0, fmt.Errorf("bootstrap sparse-fill dispatch: %w", err)
	}
	s.x[3] = s.x[17]
	s.x[4] = next
	s.x[17] = 0x6d1a
	return next, nil
}

func runBootstrapNonlinearBytePass(s *state) (uint64, error) {
	const maximumBytes = 4096
	if s == nil {
		return 0, fmt.Errorf("bootstrap nonlinear-byte state is unavailable")
	}
	targets := [2]uint64{}
	for bit := uint32(0); bit < 2; bit++ {
		next, err := loadBootstrapPointer(s, s.x[30], uint32(s.x[10])^bit*uint32(s.x[19]))
		if err != nil {
			return 0, fmt.Errorf("bootstrap nonlinear-byte target %d: %w", bit, err)
		}
		targets[bit] = next
	}
	if targets[0] == bootstrapNonlinearBytePass && targets[1] == bootstrapNonlinearBytePass {
		return 0, fmt.Errorf("bootstrap nonlinear-byte pass has no exit target")
	}

	for index := 0; index < maximumBytes; index++ {
		s.traceBlock(bootstrapNonlinearBytePass)
		position := s.x[20] - 1
		if err := s.storeValue(s.sp+0x280, position, 8); err != nil {
			return 0, fmt.Errorf("bootstrap nonlinear-byte position: %w", err)
		}
		s.x[20] = uint64(uint32(s.x[12]) ^ uint32(s.x[15]))
		s.x[21] = uint64(uint32(s.x[16]) & uint32(s.x[12]<<1))
		s.x[11] = uint64(uint32(s.x[10]) + uint32(s.x[13]))
		s.x[22] = uint64(uint32(s.x[11]) ^ uint32(s.x[0]))
		s.x[21] = uint64(uint32(s.x[22]) ^ uint32(s.x[21]))
		s.x[21] = uint64(uint32(s.x[21]) & uint32(s.x[12]))
		s.x[22] = uint64(uint32(s.x[20]) << 1)
		s.x[23] = uint64(uint32(s.x[22]) ^ uint32(s.x[1]))
		s.x[20] = uint64(uint32(s.x[23]) & uint32(s.x[20]))
		s.x[22] = uint64(uint32(s.x[22]) & uint32(s.x[2]))
		s.x[20] = uint64(uint32(s.x[20]) ^ uint32(s.x[22]))
		s.x[22] = uint64(uint32(s.x[20]) ^ uint32(s.x[3]))
		s.x[20] = uint64(uint32(s.x[20]) & uint32(s.x[21]<<2))
		s.x[20] = uint64(uint32(s.x[20]) ^ uint32(s.x[21]))
		s.x[21] = uint64(uint32(s.x[22]) << 2)
		s.x[23] = uint64(uint32(s.x[21]) ^ uint32(s.x[4]))
		s.x[22] = uint64(uint32(s.x[23]) & uint32(s.x[22]))
		s.x[21] = uint64(uint32(s.x[21]) & uint32(s.x[5]))
		s.x[21] = uint64(uint32(s.x[21]) ^ uint32(s.x[6]))
		s.x[21] = uint64(uint32(s.x[21]) ^ uint32(s.x[22]))
		s.x[21] = uint64(uint32(s.x[21]) & uint32(s.x[20]<<4))
		s.x[20] = uint64(uint32(s.x[21]) ^ uint32(s.x[20]))
		s.x[20] = uint64(uint32(s.x[12]) ^ uint32(s.x[20]<<1))
		s.x[20] = uint64(uint32(s.x[20]) ^ uint32(s.x[7]))
		s.x[21] = uint64(uint32(s.x[12]) + 0xe2)
		s.x[20] = uint64(uint32(s.x[20]) * uint32(s.x[21]))
		if err := s.storeValue(s.x[9]+position, uint64(byte(s.x[20])), 1); err != nil {
			return 0, fmt.Errorf("bootstrap nonlinear-byte output: %w", err)
		}
		s.x[12]++
		s.x[20] = position
		s.compare64(position, s.x[14])
		dispatchBit := uint32(boolValue(s.condition("hi")))
		next := targets[dispatchBit]
		s.x[17] = 0x9465
		s.x[21] = next
		if next != bootstrapNonlinearBytePass {
			return next, nil
		}
	}
	return 0, fmt.Errorf("bootstrap nonlinear byte pass exceeded %d bytes", maximumBytes)
}

func runBootstrapNonlinearByteStep(s *state) (uint64, error) {
	if s == nil {
		return 0, fmt.Errorf("bootstrap nonlinear-byte state is unavailable")
	}
	s.traceBlock(bootstrapNonlinearBytePass)
	position := s.x[20] - 1
	if err := s.storeValue(s.sp+0x280, position, 8); err != nil {
		return 0, fmt.Errorf("bootstrap nonlinear-byte position: %w", err)
	}
	s.x[20] = uint64(uint32(s.x[12]) ^ uint32(s.x[15]))
	s.x[21] = uint64(uint32(s.x[16]) & uint32(s.x[12]<<1))
	s.x[11] = uint64(uint32(s.x[10]) + uint32(s.x[13]))
	s.x[22] = uint64(uint32(s.x[11]) ^ uint32(s.x[0]))
	s.x[21] = uint64(uint32(s.x[22]) ^ uint32(s.x[21]))
	s.x[21] = uint64(uint32(s.x[21]) & uint32(s.x[12]))
	s.x[22] = uint64(uint32(s.x[20]) << 1)
	s.x[23] = uint64(uint32(s.x[22]) ^ uint32(s.x[1]))
	s.x[20] = uint64(uint32(s.x[23]) & uint32(s.x[20]))
	s.x[22] = uint64(uint32(s.x[22]) & uint32(s.x[2]))
	s.x[20] = uint64(uint32(s.x[20]) ^ uint32(s.x[22]))
	s.x[22] = uint64(uint32(s.x[20]) ^ uint32(s.x[3]))
	s.x[20] = uint64(uint32(s.x[20]) & uint32(s.x[21]<<2))
	s.x[20] = uint64(uint32(s.x[20]) ^ uint32(s.x[21]))
	s.x[21] = uint64(uint32(s.x[22]) << 2)
	s.x[23] = uint64(uint32(s.x[21]) ^ uint32(s.x[4]))
	s.x[22] = uint64(uint32(s.x[23]) & uint32(s.x[22]))
	s.x[21] = uint64(uint32(s.x[21]) & uint32(s.x[5]))
	s.x[21] = uint64(uint32(s.x[21]) ^ uint32(s.x[6]))
	s.x[21] = uint64(uint32(s.x[21]) ^ uint32(s.x[22]))
	s.x[21] = uint64(uint32(s.x[21]) & uint32(s.x[20]<<4))
	s.x[20] = uint64(uint32(s.x[21]) ^ uint32(s.x[20]))
	s.x[20] = uint64(uint32(s.x[12]) ^ uint32(s.x[20]<<1))
	s.x[20] = uint64(uint32(s.x[20]) ^ uint32(s.x[7]))
	s.x[21] = uint64(uint32(s.x[12]) + 0xe2)
	s.x[20] = uint64(uint32(s.x[20]) * uint32(s.x[21]))
	if err := s.storeValue(s.x[9]+position, uint64(byte(s.x[20])), 1); err != nil {
		return 0, fmt.Errorf("bootstrap nonlinear-byte output: %w", err)
	}
	s.x[12]++
	savedPosition, err := s.load(s.sp+0x280, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap nonlinear-byte saved position: %w", err)
	}
	s.x[20] = savedPosition
	s.compare64(s.x[20], s.x[14])
	dispatchIndex := uint32(boolValue(s.condition("hi")))*uint32(s.x[19]) ^ uint32(s.x[10])
	next, err := loadBootstrapPointer(s, s.x[30], dispatchIndex)
	if err != nil {
		return 0, fmt.Errorf("bootstrap nonlinear-byte dispatch: %w", err)
	}
	s.x[17] = 0x9465
	s.x[21] = next
	return next, nil
}

func runBootstrapByteStatePass(s *state) (uint64, error) {
	const maximumBytes = 4096
	if s == nil {
		return 0, fmt.Errorf("bootstrap byte-state state is unavailable")
	}
	targets := [2]uint64{}
	for bit := uint32(0); bit < 2; bit++ {
		index := uint32(s.x[24]) ^ bit*0x22
		next, err := loadBootstrapPointer(s, s.x[30], index)
		if err != nil {
			return 0, fmt.Errorf("bootstrap byte-state target %d: %w", bit, err)
		}
		targets[bit] = next
	}
	if targets[0] == bootstrapByteStatePass && targets[1] == bootstrapByteStatePass {
		return 0, fmt.Errorf("bootstrap byte-state pass has no exit target")
	}

	for index := 0; index < maximumBytes; index++ {
		s.traceBlock(bootstrapByteStatePass)
		s.x[10] = 0xadf89487
		s.compare32(uint32(s.x[11]), uint32(s.x[10]))
		firstBound := uint32(boolValue(s.condition("lo")))
		s.x[12] = uint64(firstBound)
		s.x[10] = 0x1eda
		value, err := s.load(s.x[20]+uint64(uint32(s.x[9])), 1)
		if err != nil {
			return 0, fmt.Errorf("bootstrap byte-state input: %w", err)
		}
		s.x[13] = uint64(byte(value))
		s.x[14] = uint64(uint32(s.x[9]) + 1)
		s.x[15] = 0x21dd12bf
		s.x[8] = uint64(uint32(s.x[8]) ^ uint32(s.x[15]))
		s.x[15] = 0xffffff1b
		s.x[15] = uint64(uint32(s.x[13]) ^ uint32(s.x[15]))
		s.x[8] = uint64(uint32(s.x[15]) & uint32(s.x[8]))
		s.x[13] = uint64(uint32(s.x[13]) & 0xffffffbf)
		s.x[15] = 0x21dd121f
		s.x[13] = uint64(uint32(s.x[13]) | uint32(s.x[15]))
		s.x[8] = uint64(uint32(s.x[13]) ^ uint32(s.x[8]))
		s.x[13] = 0xadf80939
		s.x[10] = uint64(uint32(s.x[24]) + uint32(s.x[10]))
		s.x[13] = uint64(uint32(s.x[10]) ^ uint32(s.x[13]))
		s.x[15] = 0x52076b78
		s.compare32(uint32(s.x[14]), uint32(s.x[15]))
		selector := firstBound ^ uint32(boolValue(s.condition("hi")))
		s.x[9] = uint64(uint32(s.x[9]) + uint32(s.x[13]))
		s.compare32(uint32(s.x[9]), uint32(s.x[11]))
		secondBound := uint32(boolValue(s.condition("lo")))
		s.compare32(selector, 0)
		dispatchBit := secondBound
		if s.condition("ne") {
			dispatchBit = firstBound
		}
		s.x[12] = uint64(dispatchBit << 1)
		next := targets[dispatchBit]
		s.x[9] = s.x[14]
		s.x[12] = next
		s.x[17] = 0x37fb
		if next != bootstrapByteStatePass {
			return next, nil
		}
	}
	return 0, fmt.Errorf("bootstrap byte-state pass exceeded %d bytes", maximumBytes)
}

func runBootstrapByteStateStep(s *state) (uint64, error) {
	if s == nil {
		return 0, fmt.Errorf("bootstrap byte-state state is unavailable")
	}
	s.traceBlock(bootstrapByteStatePass)
	s.x[10] = 0xadf89487
	s.compare32(uint32(s.x[11]), uint32(s.x[10]))
	firstBound := uint32(boolValue(s.condition("lo")))
	s.x[12] = uint64(firstBound)
	s.x[10] = 0x1eda
	value, err := s.load(s.x[20]+uint64(uint32(s.x[9])), 1)
	if err != nil {
		return 0, fmt.Errorf("bootstrap byte-state input: %w", err)
	}
	s.x[13] = uint64(byte(value))
	s.x[14] = uint64(uint32(s.x[9]) + 1)
	s.x[15] = 0x21dd12bf
	s.x[8] = uint64(uint32(s.x[8]) ^ uint32(s.x[15]))
	s.x[15] = 0xffffff1b
	s.x[15] = uint64(uint32(s.x[13]) ^ uint32(s.x[15]))
	s.x[8] = uint64(uint32(s.x[15]) & uint32(s.x[8]))
	s.x[13] = uint64(uint32(s.x[13]) & 0xffffffbf)
	s.x[15] = 0x21dd121f
	s.x[13] = uint64(uint32(s.x[13]) | uint32(s.x[15]))
	s.x[8] = uint64(uint32(s.x[13]) ^ uint32(s.x[8]))
	s.x[13] = 0xadf80939
	s.x[10] = uint64(uint32(s.x[24]) + uint32(s.x[10]))
	s.x[13] = uint64(uint32(s.x[10]) ^ uint32(s.x[13]))
	s.x[15] = 0x52076b78
	s.compare32(uint32(s.x[14]), uint32(s.x[15]))
	selector := firstBound ^ uint32(boolValue(s.condition("hi")))
	s.x[9] = uint64(uint32(s.x[9]) + uint32(s.x[13]))
	s.compare32(uint32(s.x[9]), uint32(s.x[11]))
	secondBound := uint32(boolValue(s.condition("lo")))
	s.compare32(selector, 0)
	dispatchBit := secondBound
	if s.condition("ne") {
		dispatchBit = firstBound
	}
	s.x[12] = uint64(dispatchBit << 1)
	dispatchIndex := ((dispatchBit << 1) | (dispatchBit << 5)) ^ uint32(s.x[24])
	next, err := loadBootstrapPointer(s, s.x[30], dispatchIndex)
	if err != nil {
		return 0, fmt.Errorf("bootstrap byte-state dispatch: %w", err)
	}
	s.x[9] = s.x[14]
	s.x[12] = next
	s.x[17] = 0x37fb
	return next, nil
}

func runBootstrapAccumulatorCycleCircuit(s *state) (uint64, error) {
	const maximumRounds = 4096
	for round := 0; round < maximumRounds; round++ {
		next, err := runBootstrapAccumulatorCycleSeed(s)
		if err != nil || next != bootstrapAccumulatorCyclePrep {
			return next, err
		}
		next, err = runBootstrapAccumulatorCyclePrepBlock(s)
		if err != nil || next != bootstrapAccumulatorPassA {
			return next, err
		}
		next, err = runBootstrapAccumulatorPass(s, bootstrapAccumulatorPassA)
		if err != nil || next != bootstrapAccumulatorCycleKey {
			return next, err
		}
		next, err = runBootstrapAccumulatorCycleKeyBlock(s)
		if err != nil || next != bootstrapAccumulatorCycleTail {
			return next, err
		}
		next, err = runBootstrapAccumulatorCycleTailBlock(s)
		if err != nil || next != bootstrapAccumulatorCycleSeed {
			return next, err
		}
	}
	return 0, fmt.Errorf("bootstrap accumulator cycle exceeded %d rounds", maximumRounds)
}

func runBootstrapAccumulatorCycleSeed(s *state) (uint64, error) {
	if s == nil {
		return 0, fmt.Errorf("bootstrap accumulator-cycle state is unavailable")
	}
	s.traceBlock(bootstrapAccumulatorCycleSeed)
	bound, err := s.load(s.sp+0x7c, 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap accumulator-cycle bound: %w", err)
	}
	s.x[8] = uint64(uint32(bound))
	s.x[13] = 0x7ffffffe
	s.compare32(uint32(s.x[8]), uint32(s.x[13]))
	gate := uint32(boolValue(s.condition("lo")))
	s.x[17] = uint64(gate)
	s.x[8] = uint64(uint32(s.x[11]) ^ 0x6797)
	s.x[13] = 0x8b1
	s.x[22] = uint64(uint32(s.x[8]) * uint32(s.x[13]))
	s.x[8] = 0xaa1797c1
	if err := s.storeValue(s.sp+0x98, uint64(uint32(s.x[4])), 4); err != nil {
		return 0, fmt.Errorf("bootstrap accumulator-cycle counter: %w", err)
	}
	s.x[8] = uint64(uint32(s.x[4]) + uint32(s.x[8]))
	table, err := s.load(s.sp+0x70, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap accumulator-cycle table: %w", err)
	}
	s.x[13] = table
	word, err := loadBootstrapWord(s, table, uint32(s.x[8]))
	if err != nil {
		return 0, fmt.Errorf("bootstrap accumulator-cycle table word: %w", err)
	}
	s.x[8] = word
	s.x[13] = uint64(uint32(s.x[22]) + uint32(s.x[2]))
	s.x[7] = uint64(uint32(s.x[13]) ^ uint32(s.x[8]))
	seed, err := s.load(s.sp+0xd8, 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap accumulator-cycle seed: %w", err)
	}
	s.x[8] = seed ^ s.x[16]
	coefficient, err := s.load(s.sp+0x80, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap accumulator-cycle coefficient: %w", err)
	}
	s.x[13] = coefficient
	s.x[8] = uint64(uint32(coefficient))*uint64(uint32(s.x[7])) + s.x[8]
	multiplier, err := s.load(s.sp+0x6c, 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap accumulator-cycle multiplier: %w", err)
	}
	s.x[13] = uint64(uint32(multiplier))
	s.x[0] = uint64(uint32(s.x[8]) * uint32(multiplier))
	s.x[4] = uint64(uint32(s.x[11]) + 0x32ce)
	s.x[13] = s.x[4] ^ uint64(uint32(0xffff65a0))
	s.x[13] = s.x[8] & s.x[13]
	s.x[19] = 0x224e6f15
	wordPointer, err := s.load(s.sp+0x60, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap accumulator-cycle word pointer: %w", err)
	}
	word, err = s.load(wordPointer, 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap accumulator-cycle word: %w", err)
	}
	s.x[2] = word ^ s.x[16]
	s.x[13] += uint64(uint32(s.x[2])) * uint64(uint32(s.x[0]))
	s.x[2] = uint64(gate << 6)
	s.x[17] = uint64(uint32(s.x[2]) - (gate << 3))
	s.x[11] = uint64(uint32(s.x[11]) ^ uint32(s.x[17]))
	dispatchTable, err := s.load(s.sp+0xa8, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap accumulator-cycle dispatch table: %w", err)
	}
	s.x[15] = dispatchTable
	next, err := loadBootstrapPointer(s, dispatchTable, uint32(s.x[11]))
	if err != nil {
		return 0, fmt.Errorf("bootstrap accumulator-cycle dispatch: %w", err)
	}
	s.x[5] = 0x4bef61bf
	s.x[11] = next
	s.x[17] = 0x5bb3
	s.x[20] = 0x360cf6ff
	s.x[21] = 0xdca6
	s.x[25] = 0xe9e5d8f9
	s.x[28] = 0x1df6773d
	return next, nil
}

func runBootstrapAccumulatorCyclePrepBlock(s *state) (uint64, error) {
	s.traceBlock(bootstrapAccumulatorCyclePrep)
	s.x[10] = 0xfdce
	s.x[2] = uint64(uint32(s.x[4]) ^ uint32(s.x[10]))
	object, err := s.load(s.sp+0x40, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap accumulator-cycle object: %w", err)
	}
	s.x[10] = object
	s.x[27], err = s.load(object+8, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap accumulator-cycle right table: %w", err)
	}
	s.x[10], err = s.load(s.sp+0xa0, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap accumulator-cycle left object: %w", err)
	}
	s.x[11], err = s.load(s.sp+0xa8, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap accumulator-cycle dispatch: %w", err)
	}
	s.x[26], err = s.load(s.x[10]+8, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap accumulator-cycle left table: %w", err)
	}
	s.x[10], err = s.load(s.sp+0x90, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap accumulator-cycle source table: %w", err)
	}
	s.x[5], err = s.load(s.x[11]+shiftLeft(signExtend32(uint32(s.x[4])), 3, 64), 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap accumulator-cycle source: %w", err)
	}
	s.x[11] = 0x4bef61bf
	s.x[17] = 0xc9f2d595
	s.x[4] = uint64(uint32(s.x[2]) + uint32(s.x[17]))
	s.x[17] = 0x58a9
	return bootstrapAccumulatorPassA, nil
}

func runBootstrapAccumulatorCycleKeyBlock(s *state) (uint64, error) {
	s.traceBlock(bootstrapAccumulatorCycleKey)
	s.x[4] = 0xf1c88bf3
	s.compare32(uint32(s.x[10]), uint32(s.x[4]))
	gate := uint32(boolValue(s.condition("gt")))
	s.x[17] = uint64(uint32(s.x[22]) ^ 0x34a7)
	s.x[2] = 0xb7
	s.x[7] = uint64(uint32(s.x[17]) * uint32(s.x[2]))
	s.x[17] = uint64(uint32(s.x[22]) + 0xee2e6dd3)
	s.x[2] = 0x11d1fe4f
	s.x[27] = uint64(uint32(s.x[17]) & uint32(s.x[2]))
	s.x[17] = uint64(uint32(s.x[7]) ^ 0x42484e79)
	s.x[17] = uint64(uint32(s.x[17]) + uint32(s.x[5]))
	s.x[12] = 0x0e37740c
	s.compare32(uint32(s.x[17]), uint32(s.x[12]))
	selector := gate ^ uint32(boolValue(s.condition("lt")))
	s.compare32(uint32(s.x[17]), uint32(s.x[3]))
	secondBound := uint32(boolValue(s.condition("lt")))
	s.compare32(selector, 0)
	dispatchBit := secondBound
	if s.condition("ne") {
		dispatchBit = gate
	}
	dispatchBit ^= 1
	dispatchMultiplier := uint32(s.x[27]) ^ 0xa046
	dispatchIndex := dispatchBit*dispatchMultiplier ^ uint32(s.x[22])
	dispatchTable, err := s.load(s.sp+0xa8, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap accumulator-cycle rekey table: %w", err)
	}
	s.x[15], s.x[23] = dispatchTable, dispatchTable
	next, err := loadBootstrapPointer(s, dispatchTable, dispatchIndex)
	if err != nil {
		return 0, fmt.Errorf("bootstrap accumulator-cycle rekey dispatch: %w", err)
	}
	s.x[1] = 0x2cfc0ae8a98f89c6
	s.x[6] = 0xb4109e41
	s.x[11] = next
	s.x[12] = 0x449d4676
	s.x[14] = 0x4248124f
	s.x[15] = 0x1430e33e4efdacf8
	s.x[17] = 0xe4e0
	s.x[19] = 0x446c
	s.x[20] = 0xfc7e
	s.x[21] = ^uint64(0x5c16)
	s.x[24] = 0xd303f5175670763a
	s.x[25] = 0x38cce9d6c53c8609
	s.x[26] = 0xc73316293ac379f7
	s.x[28] = 0xebcf1cc1b1025308
	return next, nil
}

func runBootstrapAccumulatorCycleTailBlock(s *state) (uint64, error) {
	s.traceBlock(bootstrapAccumulatorCycleTail)
	bound, err := s.load(s.sp+0x8c, 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap accumulator-cycle tail bound: %w", err)
	}
	s.x[2] = uint64(uint32(bound))
	s.x[6] = 0x68ed69f7
	s.compare32(uint32(s.x[2]), uint32(s.x[6]))
	gate := uint32(boolValue(s.condition("lt")))
	s.x[11] = uint64(uint32(s.x[7]) + 0xb5a)
	s.x[0] = uint64(uint32(s.x[7]) + 0x303c9dc4)
	s.x[4] = 0xcfc36fbe
	s.x[0] = uint64(uint32(s.x[0]) & uint32(s.x[4]))
	s.x[13] >>= 32
	s.x[8] = s.x[13] + (s.x[8] >> 32)
	s.x[5], err = s.load(s.sp+0x50, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap accumulator-cycle tail output: %w", err)
	}
	s.x[4], err = s.load(s.sp+0x58, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap accumulator-cycle tail word pointer: %w", err)
	}
	word, err := s.load(s.x[4], 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap accumulator-cycle tail word: %w", err)
	}
	s.x[13] = word ^ s.x[16]
	s.x[8] += s.x[13]
	s.x[13] = s.x[0] + 0x00000003bc522814
	s.x[13] &= s.x[8] << 1
	s.x[8] -= s.x[13]
	s.x[8] += 0x64938749de2948e7
	s.x[13] = uint64(uint32(s.x[8]) ^ 0xfc67ebdc)
	if err := s.storeValue(s.x[5], uint64(uint32(s.x[13])), 4); err != nil {
		return 0, fmt.Errorf("bootstrap accumulator-cycle tail first output: %w", err)
	}
	s.x[8] >>= 32
	s.x[13] = 0x46dd2472
	s.x[8] = uint64(uint32(s.x[8]) ^ uint32(s.x[13]))
	if err := s.storeValue(s.x[4], uint64(uint32(s.x[8])), 4); err != nil {
		return 0, fmt.Errorf("bootstrap accumulator-cycle tail second output: %w", err)
	}
	counter, err := s.load(s.sp+0x98, 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap accumulator-cycle tail counter: %w", err)
	}
	s.x[8] = uint64(uint32(counter))
	s.x[20] = uint64(uint32(counter) + 1)
	s.x[13] = 0x930501b9
	s.x[8] = uint64(uint32(s.x[8]) + uint32(s.x[13]))
	s.compare32(uint32(s.x[8]), uint32(s.x[6]))
	firstBound := uint32(boolValue(s.condition("lt")))
	selector := gate ^ firstBound
	s.compare32(uint32(s.x[8]), uint32(s.x[2]))
	secondBound := uint32(boolValue(s.condition("lt")))
	s.compare32(selector, 0)
	dispatchBit := secondBound
	if s.condition("ne") {
		dispatchBit = gate
	}
	dispatchIndex := dispatchBit*0x56 ^ uint32(s.x[7])
	next, err := loadBootstrapPointer(s, s.x[23], dispatchIndex)
	if err != nil {
		return 0, fmt.Errorf("bootstrap accumulator-cycle tail dispatch: %w", err)
	}
	s.x[1], err = s.load(s.sp+0x4c, 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap accumulator-cycle tail limit: %w", err)
	}
	s.x[2] = 0x224e6f15
	s.x[4] = s.x[20]
	s.x[8] = next
	s.x[12] = 0xb4109e41
	s.x[14] = 0xb4109e42
	s.x[17] = 0x902d
	s.x[24] = 0x449d4676
	return next, nil
}

func runBootstrapByteWriterCycle(s *state) (uint64, error) {
	if s == nil {
		return 0, fmt.Errorf("bootstrap byte-writer state is unavailable")
	}
	const maximumBytes = 4096
	for count := 0; count < maximumBytes; count++ {
		next, err := runBootstrapByteWriterControl(s)
		if err != nil {
			return 0, err
		}
		switch next {
		case bootstrapByteWriterMix:
			next, err = runBootstrapByteWriterMix(s)
			if err != nil {
				return 0, err
			}
			if next != bootstrapByteWriterLoad {
				return next, nil
			}
			next, err = runBootstrapByteWriterLoad(s)
			if err != nil {
				return 0, err
			}
			if next != bootstrapByteWriterStore {
				return next, nil
			}
		case bootstrapByteWriterStore:
		default:
			return next, nil
		}

		next, err = runBootstrapByteWriterStore(s)
		if err != nil {
			return 0, err
		}
		if next != bootstrapByteWriterControl {
			return next, nil
		}
	}
	return 0, fmt.Errorf("bootstrap byte-writer cycle exceeded %d bytes", maximumBytes)
}

func runBootstrapByteWriterControl(s *state) (uint64, error) {
	s.traceBlock(bootstrapByteWriterControl)
	s.x[17] = uint64(uint32(s.x[6]) + uint32(s.x[14]))
	s.x[26] = uint64(uint32(s.x[17]) & uint32(s.x[16]))
	s.x[17] = uint64(uint32(s.x[26]) + uint32(s.x[1]))
	s.compare32(uint32(s.x[15]), uint32(s.x[17]))
	dispatchIndex := uint32(boolValue(s.condition("eq")))*uint32(s.x[2]) ^ uint32(s.x[6])
	next, err := loadBootstrapPointer(s, s.x[11], dispatchIndex)
	if err != nil {
		return 0, fmt.Errorf("bootstrap byte-writer control dispatch: %w", err)
	}
	s.x[17] = 0x4989
	s.x[28] = next
	return next, nil
}

func runBootstrapByteWriterMix(s *state) (uint64, error) {
	s.traceBlock(bootstrapByteWriterMix)
	s.x[15] = 0
	s.x[17] = uint64(uint32(s.x[6]) + uint32(s.x[3]))
	s.x[26] = uint64(uint32(s.x[17]) & uint32(s.x[4]))
	less := uint32(0)
	if uint32(s.x[12]) < uint32(s.x[9]) {
		less = 1
	}
	s.x[27] = uint64(uint32(s.x[26]) ^ uint32(s.x[5]))
	s.x[17] = uint64(uint32(s.x[27]) * less)
	dispatchIndex := uint32(s.x[17]) ^ uint32(s.x[6])
	next, err := loadBootstrapPointer(s, s.x[11], dispatchIndex)
	if err != nil {
		return 0, fmt.Errorf("bootstrap byte-writer mix dispatch: %w", err)
	}
	s.x[17] = 0x4989
	s.x[27] = 0x224ea33b
	s.x[28] = next
	return next, nil
}

func runBootstrapByteWriterLoad(s *state) (uint64, error) {
	s.traceBlock(bootstrapByteWriterLoad)
	s.x[15] = 0
	s.x[17] = uint64(uint32(s.x[6]) + uint32(s.x[7]))
	s.x[26] = uint64(uint32(s.x[17]) & uint32(s.x[19]))
	table, err := s.load(s.x[8]+8, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap byte-writer input table: %w", err)
	}
	s.x[17] = table
	s.x[28] = uint64(uint32(s.x[26]) ^ uint32(s.x[20]))
	word, err := loadBootstrapWord(s, table, uint32(s.x[12]))
	if err != nil {
		return 0, fmt.Errorf("bootstrap byte-writer input: %w", err)
	}
	s.x[27] = word
	s.x[12] = uint64(uint32(s.x[28]) + uint32(s.x[12]))
	next, err := loadBootstrapPointer(s, s.x[11], uint32(s.x[6]))
	if err != nil {
		return 0, fmt.Errorf("bootstrap byte-writer load dispatch: %w", err)
	}
	s.x[6] = next
	s.x[17] = 0x4989
	return next, nil
}

func runBootstrapByteWriterStore(s *state) (uint64, error) {
	s.traceBlock(bootstrapByteWriterStore)
	output, err := s.load(s.x[0]+0x18, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap byte-writer output: %w", err)
	}
	s.x[17] = output
	s.x[6] = uint64(uint32(s.x[26]) + uint32(s.x[21]))
	s.x[6] = uint64(uint32(s.x[6]) & uint32(s.x[22]))
	if uint32(s.x[23]) != 0xca {
		return 0, fmt.Errorf("unsupported bootstrap byte-writer mask %#x", uint32(s.x[23]))
	}
	encoded := uint32(s.x[27]) ^ uint32(s.x[13])
	s.x[28] = uint64(byte(encoded>>uint32(s.x[15])) ^ 0x1b)
	s.x[30] = uint64(uint32(s.x[10]) - 1)
	s.x[10] = uint64(uint32(s.x[10]) + uint32(s.x[24]))
	s.x[10] = uint64(uint32(s.x[10]) + uint32(s.x[6]))
	if err := s.storeValue(output+uint64(uint32(s.x[10])), uint64(byte(s.x[28])), 1); err != nil {
		return 0, fmt.Errorf("bootstrap byte-writer output byte: %w", err)
	}
	s.x[15] = uint64(uint32(s.x[15]) + 8)
	s.compare32(uint32(s.x[30]), uint32(s.x[25]))
	more := uint32(boolValue(s.condition("ne")))
	s.x[17] = uint64(more << 8)
	s.x[10] = uint64(uint32(s.x[17]) - (more << 3))
	dispatchIndex := uint32(s.x[26]) ^ uint32(s.x[10])
	next, err := loadBootstrapPointer(s, s.x[11], dispatchIndex)
	if err != nil {
		return 0, fmt.Errorf("bootstrap byte-writer store dispatch: %w", err)
	}
	s.x[10] = s.x[30]
	s.x[17] = 0x320d
	s.x[26] = 0xce757bc5
	s.x[28] = next
	return next, nil
}

// runBootstrapRecurrenceStep is the direct form of the hot recurrence used by
// the 26A428 enrollment bootstrap. The surrounding generated CFG still sets up
// its register ABI, but the repeated arithmetic no longer needs the compact
// instruction interpreter.
func runBootstrapRecurrenceStep(s *state, _ fault) error {
	if s == nil {
		return fmt.Errorf("bootstrap recurrence state is unavailable")
	}

	previousTracePC := s.tracePC
	defer func() { s.tracePC = previousTracePC }()
	s.tracePC = compact26A428BootstrapKernel.start
	s.traceBlock(s.tracePC)

	maskedSelector := uint32(s.x[14]) & uint32(s.x[10])
	tableIndex := uint32(s.x[20]) + uint32(s.x[13])

	left, err := loadBootstrapWord(s, s.x[4], tableIndex)
	if err != nil {
		return fmt.Errorf("bootstrap recurrence left table: %w", err)
	}
	left = (left ^ s.x[24]) + (s.x[28] >> 32)

	right, err := loadBootstrapWord(s, s.x[5], tableIndex)
	if err != nil {
		return fmt.Errorf("bootstrap recurrence right table: %w", err)
	}
	right ^= s.x[24]
	s.x[28] = uint64(uint32(right))*uint64(uint32(s.x[12])) + left

	shift := uint64(maskedSelector) ^ s.x[23]
	mixed := shiftRight(s.x[30], int(shift), 64) ^ s.x[6]

	coefficient, err := loadBootstrapWord(s, s.x[22], tableIndex)
	if err != nil {
		return fmt.Errorf("bootstrap recurrence coefficient table: %w", err)
	}
	coefficient ^= s.x[24]
	mixed += uint64(uint32(coefficient))*uint64(uint32(s.x[19])) + uint64(uint32(s.x[28]))
	mixed -= s.x[7] & (mixed << 1)
	s.x[30] = mixed + s.x[26]

	stored := uint32(s.x[30]) ^ uint32(s.x[27])
	storeIndex := uint32(s.x[20]) + uint32(s.x[21])
	if err := storeBootstrapWord(s, s.x[4], storeIndex, stored); err != nil {
		return fmt.Errorf("bootstrap recurrence output table: %w", err)
	}

	next := uint32(s.x[20]) + 1
	branchIndex := uint32(s.x[20]) + uint32(s.x[16])
	// Preserve the flags left by CMP/CCMP. The following dispatcher only uses
	// the boolean result, but exact flags make this a drop-in replacement for
	// the lifted sequence and keep differential tests honest.
	s.compare32(branchIndex, uint32(s.x[15]))
	if s.condition("le") {
		s.compare32(branchIndex, uint32(s.x[1]))
	} else {
		s.setNZCV(8)
	}
	dispatchIndex := uint32(boolValue(s.condition("lt")))
	dispatchIndex = dispatchIndex*uint32(s.x[0]) ^ uint32(s.x[9])
	dispatch, err := loadBootstrapPointer(s, s.x[25], dispatchIndex)
	if err != nil {
		return fmt.Errorf("bootstrap recurrence dispatch table: %w", err)
	}

	s.x[2] = uint64(next)
	s.x[3] = dispatch
	s.x[11] = uint64(maskedSelector)
	s.x[17] = 0xcd95
	s.x[20] = uint64(next)
	return nil
}

func loadBootstrapWord(s *state, base uint64, index uint32) (uint64, error) {
	return s.load(addOffset(base, int64(int32(index))*4), 4)
}

func storeBootstrapWord(s *state, base uint64, index, value uint32) error {
	return s.storeValue(addOffset(base, int64(int32(index))*4), uint64(value), 4)
}

func loadBootstrapPointer(s *state, base uint64, index uint32) (uint64, error) {
	return s.load(addOffset(base, int64(int32(index))*8), 8)
}
