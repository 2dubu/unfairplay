package recovery

import "fmt"

const (
	crossStageBitCycleStart = uint64(0x100035c1c)
	crossStageBitCyclePlace = uint64(0x10007d8c8)
	crossStageBitCycleTest  = uint64(0x100080290)
	crossStageBitCycleWrite = uint64(0x1000534a0)
)

func runCrossStageSemanticBlock(s *state, pc uint64) (uint64, bool, error) {
	if pc != crossStageBitCycleStart {
		return pc, false, nil
	}
	next, err := runCrossStageBitCycle(s)
	return next, true, err
}

func runCrossStageBitCycle(s *state) (uint64, error) {
	pc := crossStageBitCycleStart
	const maximumBlocks = 16384
	for count := 0; count < maximumBlocks; count++ {
		var next uint64
		var err error
		switch pc {
		case crossStageBitCycleStart:
			next, err = runCrossStageBitCycleStart(s)
		case crossStageBitCyclePlace:
			next, err = runCrossStageBitCyclePlace(s)
		case crossStageBitCycleTest:
			next, err = runCrossStageBitCycleTest(s)
		case crossStageBitCycleWrite:
			next, err = runCrossStageBitCycleWrite(s)
		default:
			return pc, nil
		}
		if err != nil {
			return pc, err
		}
		pc = next
	}
	return 0, fmt.Errorf("cross-stage bit cycle exceeded %d blocks", maximumBlocks)
}

func runCrossStageBitCycleStart(s *state) (uint64, error) {
	s.traceBlock(crossStageBitCycleStart)
	s.x[13] = 0
	s.x[25] = uint64(uint32(s.x[30]) + uint32(s.x[14]))
	next, err := loadBootstrapPointer(s, s.x[8], uint32(s.x[30]))
	if err != nil {
		return 0, fmt.Errorf("cross-stage bit-cycle start dispatch: %w", err)
	}
	s.x[11] = next
	s.x[17] = 0x558b
	s.x[28] = 0xec365926
	return next, nil
}

func runCrossStageBitCyclePlace(s *state) (uint64, error) {
	s.traceBlock(crossStageBitCyclePlace)
	output, err := s.load(s.x[0]+0x18, 8)
	if err != nil {
		return 0, fmt.Errorf("cross-stage bit-cycle source: %w", err)
	}
	s.x[17] = output
	s.x[30] = uint64(uint32(s.x[25]) ^ uint32(s.x[15]))
	s.x[5] = uint64(uint32(s.x[30]) + 0x91a)
	s.x[11] = uint64(uint32(s.x[4]) - 1)
	s.x[6] = uint64(uint32(s.x[4]) + uint32(s.x[16]))
	value, err := s.load(output+uint64(uint32(s.x[6])), 1)
	if err != nil {
		return 0, fmt.Errorf("cross-stage bit-cycle source byte: %w", err)
	}
	s.x[17] = uint64(uint32(shiftLeft(uint64(byte(value)), int(s.x[13]), 32)))
	s.x[6] = uint64(uint32(s.x[28]) &^ uint32(s.x[17]))
	s.x[17] = uint64(uint32(s.x[17]) & uint32(s.x[1]))
	s.x[28] = uint64(uint32(s.x[17]) | uint32(s.x[6]))
	s.x[13] = uint64(uint32(s.x[13]) + uint32(s.x[2]))
	s.x[13] = uint64(uint32(s.x[13]) + uint32(s.x[5]))
	s.compare32(uint32(s.x[13]), 0x20)
	dispatchIndex := uint32(boolValue(s.condition("eq")))*uint32(s.x[3]) ^ uint32(s.x[25])
	next, err := loadBootstrapPointer(s, s.x[8], dispatchIndex)
	if err != nil {
		return 0, fmt.Errorf("cross-stage bit-cycle place dispatch: %w", err)
	}
	s.x[6] = next
	s.x[17] = 0xc6db
	return next, nil
}

func runCrossStageBitCycleTest(s *state) (uint64, error) {
	s.traceBlock(crossStageBitCycleTest)
	s.x[17] = uint64(uint32(s.x[30]) + uint32(s.x[21]))
	s.x[5] = uint64(uint32(s.x[17]) & uint32(s.x[22]))
	s.x[17] = uint64(uint32(s.x[30]) ^ uint32(s.x[23]))
	s.x[25] = uint64(uint32(s.x[17]) * uint32(s.x[24]))
	s.x[17] = uint64(uint32(s.x[25]) + uint32(s.x[26]))
	s.compare32(uint32(s.x[11]), uint32(s.x[17]))
	s.x[4] = uint64(uint32(s.x[5]) ^ uint32(s.x[27]))
	dispatchIndex := uint32(boolValue(s.condition("eq")))*uint32(s.x[4]) ^ uint32(s.x[30])
	next, err := loadBootstrapPointer(s, s.x[8], dispatchIndex)
	if err != nil {
		return 0, fmt.Errorf("cross-stage bit-cycle test dispatch: %w", err)
	}
	s.x[4] = s.x[11]
	s.x[6] = next
	s.x[17] = 0x558b
	return next, nil
}

func runCrossStageBitCycleWrite(s *state) (uint64, error) {
	s.traceBlock(crossStageBitCycleWrite)
	s.x[25] = uint64(uint32(s.x[5]) + 0x4b)
	s.x[13] = uint64(uint32(s.x[5]) + 0xffffb4c2)
	s.x[17] = 0x4224
	s.x[30] = uint64(uint32(s.x[13]) | uint32(s.x[17]))
	s.x[13] = 0xce78fa1d
	s.x[17] = uint64(uint32(s.x[28]) ^ uint32(s.x[13]))
	table, err := s.load(s.x[9]+8, 8)
	if err != nil {
		return 0, fmt.Errorf("cross-stage bit-cycle output table: %w", err)
	}
	s.x[6] = table
	s.x[13] = uint64(uint32(s.x[12]) + 1)
	s.x[12] = uint64(uint32(s.x[12]) + uint32(s.x[7]))
	if err := storeBootstrapWord(s, table, uint32(s.x[12]), uint32(s.x[17])); err != nil {
		return 0, fmt.Errorf("cross-stage bit-cycle output: %w", err)
	}
	s.compare32(uint32(s.x[4]), uint32(s.x[19]))
	s.x[17] = uint64(uint32(s.x[30]) + uint32(s.x[20]))
	dispatchIndex := uint32(s.x[17])*uint32(boolValue(s.condition("eq"))) ^ uint32(s.x[5])
	next, err := loadBootstrapPointer(s, s.x[8], dispatchIndex)
	if err != nil {
		return 0, fmt.Errorf("cross-stage bit-cycle write dispatch: %w", err)
	}
	s.x[4] = s.x[11]
	s.x[5] = next
	s.x[12] = s.x[13]
	s.x[17] = 0x5233
	return next, nil
}
