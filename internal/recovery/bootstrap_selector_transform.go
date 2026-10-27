package recovery

import "fmt"

const (
	bootstrapCounterSelector = uint64(0x1009b2174)
	bootstrapByteSelector    = uint64(0x1009f1804)
	bootstrapStackSelector   = uint64(0x1008a9500)
	bootstrapScaledSelector  = uint64(0x1008efca8)
)

func runBootstrapStackSelector(s *state) (uint64, error) {
	s.traceBlock(bootstrapStackSelector)
	first, err := s.load(s.sp+8, 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap stack-selector first input: %w", err)
	}
	second, err := s.load(s.sp+12, 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap stack-selector second input: %w", err)
	}
	s.x[10], s.x[9] = uint64(uint32(first)), uint64(uint32(second))
	s.x[9] = uint64(uint32(s.x[10]) + uint32(s.x[9]))
	s.x[10] = 0x7ff05135
	s.x[11] = 0xffffd156
	s.x[12] = 0x302
	s.x[11] = uint64(uint32(s.x[19]) + uint32(s.x[11]))
	s.x[15] = uint64(uint32(s.x[11]) | uint32(s.x[12]))
	s.x[11] = 0x800faef4
	s.x[10] = uint64(uint32(s.x[19]) + uint32(s.x[10]))
	s.x[10] = uint64(uint32(s.x[10]) & uint32(s.x[11]))
	s.x[11] = 0x32c7ace6
	s.x[12] = 0x47d041fc
	s.x[13] = uint64(uint32(s.x[15]) ^ 0xc7d016f2)
	s.x[11] = uint64(uint32(s.x[28]) + uint32(s.x[11]))
	s.compare32(uint32(s.x[11]), uint32(s.x[13]))
	firstGate := uint32(boolValue(s.condition("lt")))
	s.x[13] = uint64(firstGate)
	s.x[14] = uint64(uint32(s.x[10]) ^ 0x382f3ce3)
	if err := s.storeValue(s.sp+0x80, uint64(uint32(s.x[14])), 4); err != nil {
		return 0, fmt.Errorf("bootstrap stack-selector first state: %w", err)
	}
	s.compare32(uint32(s.x[8]), uint32(s.x[14]))
	secondGate := uint32(boolValue(s.condition("gt")))
	s.x[14] = uint64(secondGate)
	s.x[13] = uint64(firstGate ^ secondGate)
	s.x[12] = uint64(uint32(s.x[8]) + uint32(s.x[12]))
	if err := s.storeValue(s.sp+0x78, uint64(uint32(s.x[12])), 4); err != nil {
		return 0, fmt.Errorf("bootstrap stack-selector second state: %w", err)
	}
	s.compare32(uint32(s.x[11]), uint32(s.x[12]))
	thirdGate := uint32(boolValue(s.condition("lt")))
	s.x[11] = uint64(thirdGate)
	s.compare32(uint32(s.x[13]), 0)
	dispatchGate := thirdGate
	if s.condition("ne") {
		dispatchGate = secondGate
	}
	s.x[11] = uint64(dispatchGate * 0x344f)
	s.x[11] = uint64(uint32(s.x[11]) ^ uint32(s.x[19]))
	next, err := loadBootstrapPointer(s, s.x[25], uint32(s.x[11]))
	if err != nil {
		return 0, fmt.Errorf("bootstrap stack-selector target: %w", err)
	}
	s.x[11] = next
	s.x[17] = 0x7bf6
	return next, nil
}

func runBootstrapScaledSelector(s *state) (uint64, error) {
	s.traceBlock(bootstrapScaledSelector)
	s.x[8] = uint64(uint32(s.x[15]) ^ 0x5924)
	s.x[8] = uint64(uint32(s.x[8]) + uint32(s.x[8]<<2))
	s.x[19] = uint64(uint32(s.x[8] << 1))
	s.x[8] = uint64(uint32(s.x[15]) + 0xe462081d)
	s.x[8] = uint64(uint32(s.x[8]) & 0x1b9df7ff)
	s.x[10] = uint64(uint32(s.x[9]) + 0xa9a2d523)
	s.x[10] = uint64(uint32(s.x[10]) + uint32(s.x[19]))
	s.x[11] = 0x7ffffffe
	s.compare32(uint32(s.x[10]), uint32(s.x[11]))
	dispatchGate := uint32(boolValue(s.condition("hi")))
	s.x[10] = uint64(dispatchGate)
	s.x[11] = uint64(uint32(s.x[8]) ^ 0x570e)
	s.x[10] = uint64(dispatchGate * uint32(s.x[11]))
	s.x[10] = uint64(uint32(s.x[10]) ^ uint32(s.x[15]))
	next, err := loadBootstrapPointer(s, s.x[25], uint32(s.x[10]))
	if err != nil {
		return 0, fmt.Errorf("bootstrap scaled-selector target: %w", err)
	}
	s.x[10] = next
	s.x[17] = 0xe52e
	return next, nil
}

func runBootstrapCounterSelector(s *state) (uint64, error) {
	s.traceBlock(bootstrapCounterSelector)
	s.x[8] = uint64(uint32(s.x[13]) + uint32(s.x[4]))
	s.x[8] = uint64(uint32(s.x[8]) | uint32(s.x[5]))
	s.x[14] = uint64(uint32(s.x[13]) + uint32(s.x[6]))
	s.x[14] = uint64(uint32(s.x[14]) | uint32(s.x[7]))
	s.compare32(uint32(s.x[22]), 0)
	s.x[10]++
	s.x[17] = s.x[8] ^ s.x[19]
	if s.condition("ne") {
		s.compare64(s.x[10], s.x[17])
	} else {
		s.setNZCV(4)
	}
	s.x[17] = uint64(uint32(boolValue(s.condition("ne"))))
	s.x[21] = uint64(uint32(s.x[14]) ^ uint32(s.x[20]))
	s.x[17] = uint64(uint32(s.x[17]) * uint32(s.x[21]))
	s.x[13] = uint64(uint32(s.x[17]) ^ uint32(s.x[13]))
	next, err := loadBootstrapPointer(s, s.x[30], uint32(s.x[13]))
	if err != nil {
		return 0, fmt.Errorf("bootstrap counter-selector target: %w", err)
	}
	s.x[21] = next
	s.x[13] = 0x30afef76
	s.x[17] = 0xbabc
	return next, nil
}

func runBootstrapByteSelector(s *state) (uint64, error) {
	s.traceBlock(bootstrapByteSelector)
	s.x[8] = uint64(uint32(s.x[14]) ^ uint32(s.x[16]))
	s.x[13] = uint64(uint32(s.x[8]) * uint32(s.x[0]))
	s.x[8] = uint64(uint32(s.x[14]) ^ uint32(s.x[1]))
	s.x[21] = uint64(uint32(s.x[8]) * uint32(s.x[2]))
	first, err := s.load(s.x[12]+uint64(uint32(s.x[10])), 1)
	if err != nil {
		return 0, fmt.Errorf("bootstrap byte-selector first input: %w", err)
	}
	s.x[17] = uint64(byte(first))
	s.x[22] = uint64(uint32(int32(int8(s.x[17]))))
	second, err := s.load(s.x[15]+uint64(uint32(s.x[10])), 1)
	if err != nil {
		return 0, fmt.Errorf("bootstrap byte-selector second input: %w", err)
	}
	s.x[8] = uint64(byte(second))
	s.compare32(uint32(s.x[17]), uint32(s.x[8]))
	s.x[17] = uint64(uint32(boolValue(s.condition("eq"))))
	s.x[17] = uint64(uint32(s.x[13]) ^ uint32(s.x[17]))
	s.x[17] &= 1
	s.x[23] = uint64(uint32(s.x[21]) ^ uint32(s.x[3]))
	s.x[17] = uint64(uint32(s.x[17]) * uint32(s.x[23]))
	s.x[14] = uint64(uint32(s.x[17]) ^ uint32(s.x[14]))
	next, err := loadBootstrapPointer(s, s.x[30], uint32(s.x[14]))
	if err != nil {
		return 0, fmt.Errorf("bootstrap byte-selector target: %w", err)
	}
	s.x[14] = next
	s.x[17] = 0xac7c
	return next, nil
}
