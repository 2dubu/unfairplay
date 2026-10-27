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
