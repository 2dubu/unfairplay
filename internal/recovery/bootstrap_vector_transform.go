package recovery

import "fmt"

const (
	bootstrapVectorExpandStart = uint64(0x1005df840)
	bootstrapVectorExpandLoop  = uint64(0x1005df87c)
	bootstrapVectorExpandExit  = uint64(0x1005df900)

	bootstrapVectorSubstituteStart = uint64(0x10051b400)
	bootstrapVectorSubstituteLoop  = uint64(0x10051b470)
	bootstrapVectorSubstituteExit  = uint64(0x10051b4f0)

	bootstrapVectorBytes = uint64(0x880)
)

// runBootstrapVectorExpand is the direct form of the first 2,176-byte SIMD
// pass. The original code widens each input byte, evaluates four polynomial
// lanes, then maps the result through a four-vector byte table.
func runBootstrapVectorExpand(s *state) (uint64, error) {
	if s == nil {
		return 0, fmt.Errorf("bootstrap vector-expand state is unavailable")
	}
	s.traceBlock(bootstrapVectorExpandStart)
	s.x[11] = 0
	const tableBase = uint64(0x100d25000)
	for register, offset := range [...]uint64{0x860, 0x870, 0x880, 0x890, 0x8a0} {
		if err := s.loadVector(register, tableBase+offset, 16); err != nil {
			return 0, fmt.Errorf("bootstrap vector-expand table %d: %w", register, err)
		}
	}
	s.vectorDup(5, 4, 16, 0x9b)
	s.vectorDup(6, 8, 16, 0x9b)
	s.vectorDup(7, 4, 32, 0x10)
	s.vectorDup(16, 4, 32, uint64(uint32(s.x[10])))

	for iteration := 0; iteration < int(bootstrapVectorBytes/16); iteration++ {
		if iteration != 0 {
			s.traceBlock(bootstrapVectorExpandLoop)
		}
		if err := runBootstrapVectorExpandIteration(s); err != nil {
			return 0, err
		}
		if s.condition("eq") {
			return bootstrapVectorExpandExit, nil
		}
	}
	return 0, fmt.Errorf("bootstrap vector-expand did not reach its natural exit")
}

func runBootstrapVectorExpandIteration(s *state) error {
	s.x[12] = s.x[13] + s.x[11]
	if err := s.loadVector(17, s.x[12], 16); err != nil {
		return fmt.Errorf("bootstrap vector-expand input: %w", err)
	}
	s.vectorUnsignedShiftLong(18, 17, 8, 8, 0, true)
	s.vectorUnsignedShiftLong(17, 17, 8, 8, 0, false)
	s.vectorUnsignedMultiplyLong(19, 17, 5, 4, 16, false)
	s.vectorUnsignedMultiplyLong(17, 17, 6, 4, 16, true)
	s.vectorUnsignedMultiplyLong(20, 18, 5, 4, 16, false)
	s.vectorUnsignedMultiplyLong(18, 18, 6, 4, 16, true)
	for register := 0; register < 4; register++ {
		s.vectorBinary("add", 21+register, register, 16, 4, 32)
		s.vectorBinary("sub", 25+register, register, 16, 4, 32)
	}
	for register := 0; register < 4; register++ {
		s.vectorBinary("mul", 24-register, 24-register, 28-register, 4, 32)
	}
	s.vectorEor3(28, 21, 16, 18)
	s.vectorEor3(27, 22, 16, 20)
	s.vectorEor3(26, 23, 16, 17)
	s.vectorEor3(25, 24, 16, 19)
	s.vectorTable(17, []int{25, 26, 27, 28}, 4)
	if err := s.storeVector(17, s.x[12], 16); err != nil {
		return fmt.Errorf("bootstrap vector-expand output: %w", err)
	}
	s.x[11] += 16
	for register := 3; register >= 0; register-- {
		s.vectorBinary("add", register, register, 7, 4, 32)
	}
	s.compare64(s.x[11], bootstrapVectorBytes)
	return nil
}

// runBootstrapVectorSubstitute is the direct form of the second 2,176-byte
// SIMD pass. It applies four byte tables, combines their polynomial lanes,
// and XORs the result with the fixed per-build mask vector.
func runBootstrapVectorSubstitute(s *state) (uint64, error) {
	if s == nil {
		return 0, fmt.Errorf("bootstrap vector-substitute state is unavailable")
	}
	s.traceBlock(bootstrapVectorSubstituteStart)
	if err := s.loadVector(2, s.x[4]+0x40, 8); err != nil {
		return 0, fmt.Errorf("bootstrap vector-substitute seed: %w", err)
	}
	seed := s.vectorLane(0, 0, 64)
	s.clearVector(3)
	s.setVectorLane(3, 0, 64, seed)
	s.vectorMultiplyAccumulate(3, 2, 1, 2, 32, false)
	if err := s.storeVector(3, s.x[4]+0x40, 8); err != nil {
		return 0, fmt.Errorf("bootstrap vector-substitute seed output: %w", err)
	}
	value, err := s.load(s.x[4]+0x48, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap vector-substitute base: %w", err)
	}
	s.x[16] = value
	s.x[17] = uint64(uint32(s.x[16]) * uint32(s.x[17]))
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[0]))
	s.vectorDup(1, 4, 32, uint64(uint32(s.x[17])))

	const firstTableBase = uint64(0x100d25000)
	for register, offset := range [...]uint64{0x860, 0x870, 0x880, 0x890} {
		if err := s.loadVector(2+register, firstTableBase+offset, 16); err != nil {
			return 0, fmt.Errorf("bootstrap vector-substitute first table %d: %w", register, err)
		}
	}
	const secondTableBase = uint64(0x100f49000)
	for index, entry := range [...]struct {
		register int
		offset   uint64
	}{{6, 0x960}, {7, 0x970}, {16, 0x980}, {17, 0x990}} {
		if err := s.loadVector(entry.register, secondTableBase+entry.offset, 16); err != nil {
			return 0, fmt.Errorf("bootstrap vector-substitute second table %d: %w", index, err)
		}
	}
	if err := s.loadVector(18, firstTableBase+0x8a0, 16); err != nil {
		return 0, fmt.Errorf("bootstrap vector-substitute index table: %w", err)
	}
	s.vectorDup(19, 16, 8, 0x93)
	s.vectorDup(20, 4, 32, 0x10)

	for iteration := 0; iteration < int(bootstrapVectorBytes/16); iteration++ {
		if iteration != 0 {
			s.traceBlock(bootstrapVectorSubstituteLoop)
		}
		if err := runBootstrapVectorSubstituteIteration(s); err != nil {
			return 0, err
		}
		if s.condition("eq") {
			return bootstrapVectorSubstituteExit, nil
		}
	}
	return 0, fmt.Errorf("bootstrap vector-substitute did not reach its natural exit")
}

func runBootstrapVectorSubstituteIteration(s *state) error {
	s.x[0] = s.x[16] + s.x[11]
	if err := s.loadVector(21, s.x[0], 16); err != nil {
		return fmt.Errorf("bootstrap vector-substitute input: %w", err)
	}
	s.vectorTable(22, []int{21}, 6)
	s.vectorTable(23, []int{21}, 7)
	s.vectorTable(24, []int{21}, 16)
	s.vectorTable(21, []int{21}, 17)
	for index, source := range [...]int{4, 5, 3, 2} {
		addDestination := [...]int{25, 26, 27, 28}[index]
		subDestination := [...]int{29, 30, 31, 8}[index]
		s.vectorBinary("add", addDestination, source, 1, 4, 32)
		s.vectorBinary("sub", subDestination, source, 1, 4, 32)
	}
	for index, right := range [...]int{8, 31, 30, 29} {
		destination := 28 - index
		s.vectorBinary("mul", destination, destination, right, 4, 32)
	}
	s.vectorEor3(9, 21, 1, 25)
	s.vectorEor3(8, 24, 1, 26)
	s.vectorEor3(10, 23, 1, 27)
	s.vectorEor3(11, 22, 1, 28)
	s.vectorTable(21, []int{8, 9, 10, 11}, 18)
	s.vectorBinary("mul", 21, 21, 19, 16, 8)
	if err := s.storeVector(21, s.x[0], 16); err != nil {
		return fmt.Errorf("bootstrap vector-substitute output: %w", err)
	}
	s.x[11] += 16
	for register := 5; register >= 2; register-- {
		s.vectorBinary("add", register, register, 20, 4, 32)
	}
	s.compare64(s.x[11], bootstrapVectorBytes)
	return nil
}
