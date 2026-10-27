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
)

// runBootstrapMixer replaces the dominant flattened bootstrap path with its
// natural nested-loop structure. A round contains either 31 or 63 recurrence
// steps, followed by five small control blocks. If a different, valid path is
// selected, the caller resumes it through the verified CFG.
func runBootstrapMixer(s *state, runtime *rawSourceRuntime, profile *executionProfile) (uint64, error) {
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
		pc, err = runBootstrapMixerBlock(s, runtime, profile, pc)
		if err != nil || pc != bootstrapMixerPrepare {
			return pc, err
		}
		pc, err = runBootstrapMixerBlock(s, runtime, profile, pc)
		if err != nil || pc != bootstrapMixerDispatch {
			return pc, err
		}
		pc, err = runBootstrapMixerBlock(s, runtime, profile, pc)
		if err != nil || pc != bootstrapMixerRound {
			return pc, err
		}
		pc, err = runBootstrapMixerBlock(s, runtime, profile, pc)
		if err != nil || pc != bootstrapMixerContinue {
			return pc, err
		}
		pc, err = runBootstrapMixerBlock(s, runtime, profile, pc)
		if err != nil {
			return pc, err
		}
		if pc != profile.bootstrapStart {
			return pc, nil
		}
	}
	return 0, fmt.Errorf("bootstrap mixer exceeded %d rounds", maximumRounds)
}

func runBootstrapMixerBlock(s *state, runtime *rawSourceRuntime, profile *executionProfile, pc uint64) (uint64, error) {
	block, ok := recoveryBlock(s, pc, profile.bootstrapBlocks, profile.supplementalBlocks, profile.verifiedFallbackBlocks)
	if !ok {
		return pc, fmt.Errorf("bootstrap mixer is missing block %#x", pc)
	}
	next, stopped, err := runCompactBlock(s, runtime, block)
	if err != nil {
		return pc, fmt.Errorf("bootstrap mixer block %#x: %w", pc, err)
	}
	if stopped {
		return next, fmt.Errorf("bootstrap mixer block %#x stopped unexpectedly", pc)
	}
	return next, nil
}

// runBootstrapWordPass executes either of the two natural word-table passes
// used by the bootstrap finalizer. Their native handlers dispatch back to
// themselves once per word; this direct loop makes the real bounds explicit.
func runBootstrapWordPass(s *state, pc uint64) (uint64, error) {
	const maximumWords = 1024
	for word := 0; word < maximumWords; word++ {
		var next uint64
		var err error
		switch pc {
		case bootstrapShortWordPass:
			next, err = runBootstrapShortWordStep(s)
		case bootstrapLongWordPass:
			next, err = runBootstrapLongWordStep(s)
		default:
			return pc, fmt.Errorf("unknown bootstrap word pass %#x", pc)
		}
		if err != nil || next != pc {
			return next, err
		}
	}
	return 0, fmt.Errorf("bootstrap word pass %#x exceeded %d words", pc, maximumWords)
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
	mixed := (right ^ s.x[24]) + s.x[25]
	chain := (right ^ s.x[26]) + s.x[9]
	mixed ^= chain
	s.x[23] = 0xc0c02ce9
	combined := ((right ^ s.x[16]) + s.x[0]) ^ mixed
	combined += s.x[2]
	combined = combined*s.x[15] + carry
	combined -= s.x[21] & (combined << 1)
	result := combined + s.x[7]
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
