package recovery

import "fmt"

func runBootstrapAccumulatorStepBCircuit(s *state) (uint64, error) {
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
