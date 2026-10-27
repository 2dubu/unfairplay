package recovery

import "fmt"

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
