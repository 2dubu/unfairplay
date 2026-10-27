package recovery

import "fmt"

const (
	bootstrapPointerCycleStart = uint64(0x1001dfbf4)
	bootstrapPointerCycleTail  = uint64(0x1001dfc00)
	bootstrapPointerCycleCheck = uint64(0x1000602d8)
	bootstrapPointerCycleStep  = uint64(0x100066ee0)
	bootstrapPointerCycleNext  = uint64(0x100087a4c)
	bootstrapByteClassifyPass  = uint64(0x100065d6c)
	bootstrapIndexedDispatch   = uint64(0x1005816cc)
	bootstrapCountdownStart    = uint64(0x1005807e0)
	bootstrapCountdownDispatch = uint64(0x100580800)
	bootstrapSmallFrameReturn  = uint64(0x10016e74c)
	bootstrapSavedFrameReturn  = uint64(0x100147c4c)
)

func runBootstrapIndexedDispatch(s *state) (uint64, error) {
	s.traceBlock(bootstrapIndexedDispatch)
	s.x[10] = uint64(uint32(s.x[8]) + 0x2a79)
	s.x[19] = uint64(uint32(s.x[10]) | 0xc00)
	s.x[10] = uint64((uint32(s.x[8]) + 0x3bf7a813) & 0xc40857cf)
	s.x[11] = uint64((uint32(s.x[10]) ^ 0xa9a335c6) + uint32(s.x[9]))
	s.x[12] = s.sp + 0x2d0
	word, err := s.load(s.x[12]+uint64(uint32(s.x[11]))*4, 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap indexed-dispatch lookup: %w", err)
	}
	s.x[11] = uint64(uint32(word))
	s.x[12] = uint64(uint32(s.x[19]) ^ 0x224e2e9f)
	s.compare32(uint32(s.x[11]), uint32(s.x[12]))
	s.x[11] = uint64(uint32(boolValue(s.condition("eq"))) * 0xc36e)
	s.x[8] = uint64(uint32(s.x[11]) ^ uint32(s.x[8]))
	next, err := loadBootstrapPointer(s, s.x[25], uint32(s.x[8]))
	if err != nil {
		return 0, fmt.Errorf("bootstrap indexed-dispatch target: %w", err)
	}
	s.x[8] = next
	s.x[17] = 0xe52e
	return next, nil
}

func runBootstrapCountdownDispatch(s *state) (uint64, error) {
	s.traceBlock(bootstrapCountdownStart)
	s.x[19] = uint64(uint32(s.x[10]) ^ 0xdaaa)
	s.x[8] = uint64(uint32(s.x[10]) ^ 0x25)
	s.x[12] = uint64((uint32(s.x[10]) ^ 0xa9a335c8) + uint32(s.x[9]))
	s.traceBlock(bootstrapCountdownDispatch)
	s.x[9] = uint64(uint32(s.x[9]) - 1)
	s.x[13] = uint64(uint32(s.x[8]) + 0x7fffa8d4)
	s.compare32(uint32(s.x[12]), uint32(s.x[13]))
	s.x[12] = uint64(uint32(boolValue(s.condition("hs"))))
	s.x[11] = uint64(uint32(s.x[12]) * uint32(s.x[11]))
	s.x[10] = uint64(uint32(s.x[11]) ^ uint32(s.x[10]))
	next, err := loadBootstrapPointer(s, s.x[25], uint32(s.x[10]))
	if err != nil {
		return 0, fmt.Errorf("bootstrap countdown-dispatch target: %w", err)
	}
	s.x[10] = next
	s.x[17] = 0xe52e
	return next, nil
}

func runBootstrapSmallFrameReturn(s *state) (uint64, error) {
	s.traceBlock(bootstrapSmallFrameReturn)
	s.sp += 0x10
	return s.x[30], nil
}

func runBootstrapSavedFrameReturn(s *state) (uint64, error) {
	s.traceBlock(bootstrapSavedFrameReturn)
	if err := s.storeValue(s.x[0]+8, uint64(uint32(s.x[11])), 4); err != nil {
		return 0, fmt.Errorf("bootstrap saved-frame result: %w", err)
	}
	for _, pair := range [...]struct {
		first, second int
		offset        uint64
	}{{29, 30, 0x60}, {20, 19, 0x50}, {22, 21, 0x40}, {24, 23, 0x30}, {26, 25, 0x20}, {28, 27, 0x10}} {
		first, err := s.load(s.sp+pair.offset, 8)
		if err != nil {
			return 0, fmt.Errorf("bootstrap saved-frame register %d: %w", pair.first, err)
		}
		second, err := s.load(s.sp+pair.offset+8, 8)
		if err != nil {
			return 0, fmt.Errorf("bootstrap saved-frame register %d: %w", pair.second, err)
		}
		s.x[pair.first] = first
		s.x[pair.second] = second
	}
	s.sp += 0x70
	return s.x[30], nil
}

func runBootstrapPointerCycle(s *state) (uint64, error) {
	pc := bootstrapPointerCycleStart
	const maximumBlocks = 8192
	for count := 0; count < maximumBlocks; count++ {
		var next uint64
		var err error
		switch pc {
		case bootstrapPointerCycleStart:
			next, err = runBootstrapPointerCycleStart(s)
		case bootstrapPointerCycleCheck:
			next, err = runBootstrapPointerCycleCheck(s)
		case bootstrapPointerCycleStep:
			next, err = runBootstrapPointerCycleStep(s)
		case bootstrapPointerCycleNext:
			next, err = runBootstrapPointerCycleNext(s)
		default:
			return pc, nil
		}
		if err != nil {
			return pc, err
		}
		pc = next
	}
	return 0, fmt.Errorf("bootstrap pointer cycle exceeded %d blocks", maximumBlocks)
}

func runBootstrapPointerCycleStart(s *state) (uint64, error) {
	s.traceBlock(bootstrapPointerCycleStart)
	s.x[1] = uint64(uint32(s.x[21]) + 0xbca)
	word, err := s.load(s.x[10], 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap pointer-cycle word: %w", err)
	}
	s.x[24] = uint64(uint32(word))
	next, err := loadBootstrapPointer(s, s.x[22], uint32(s.x[21]))
	if err != nil {
		return 0, fmt.Errorf("bootstrap pointer-cycle start dispatch: %w", err)
	}
	s.x[21] = next
	s.traceBlock(bootstrapPointerCycleTail)
	s.x[17] = 0x6dd3
	return next, nil
}

func runBootstrapPointerCycleCheck(s *state) (uint64, error) {
	s.traceBlock(bootstrapPointerCycleCheck)
	s.x[21] = uint64(uint32(s.x[1]) ^ uint32(s.x[12]))
	s.x[17] = uint64(uint32(s.x[1]) ^ uint32(s.x[13]))
	s.compare32(uint32(s.x[24]), uint32(s.x[17]))
	dispatchIndex := uint32(boolValue(s.condition("lt")))*uint32(s.x[14]) ^ uint32(s.x[1])
	next, err := loadBootstrapPointer(s, s.x[22], dispatchIndex)
	if err != nil {
		return 0, fmt.Errorf("bootstrap pointer-cycle check dispatch: %w", err)
	}
	s.x[1] = next
	s.x[17] = 0x99f1
	return next, nil
}

func runBootstrapPointerCycleStep(s *state) (uint64, error) {
	s.traceBlock(bootstrapPointerCycleStep)
	s.compare32(uint32(s.x[24]), 0x21)
	selected := uint32(boolValue(s.condition("eq")))
	s.x[1] = uint64(uint32(s.x[21]) + 0x9a48)
	s.x[2] = 0x3887
	s.x[24] = uint64(uint32(s.x[1]) ^ uint32(s.x[2]))
	dispatchIndex := uint32(s.x[24])*selected ^ uint32(s.x[21])
	next, err := loadBootstrapPointer(s, s.x[22], dispatchIndex)
	if err != nil {
		return 0, fmt.Errorf("bootstrap pointer-cycle step dispatch: %w", err)
	}
	s.x[17] = 0x1a09
	s.x[24] = next
	return next, nil
}

func runBootstrapPointerCycleNext(s *state) (uint64, error) {
	s.traceBlock(bootstrapPointerCycleNext)
	s.compare64(s.x[8], s.x[11])
	firstBound := uint32(boolValue(s.condition("lo")))
	s.x[24] = uint64(uint32(s.x[1]) + uint32(s.x[5]))
	s.x[2] = uint64(uint32(s.x[1]) + uint32(s.x[6]))
	s.x[21] = uint64(uint32(s.x[2]) & uint32(s.x[7]))
	word, err := s.load(s.x[10]+4, 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap pointer-cycle increment: %w", err)
	}
	s.x[2] = uint64(uint32(word))
	s.x[10] += s.x[2]
	nextCounter := s.x[9] + 1
	s.compare64(nextCounter, s.x[28])
	selector := firstBound ^ uint32(boolValue(s.condition("hi")))
	s.x[9] += s.x[27]
	s.compare64(s.x[9], s.x[8])
	secondBound := uint32(boolValue(s.condition("lo")))
	s.compare32(selector, 0)
	dispatchBit := secondBound
	if s.condition("ne") {
		dispatchBit = firstBound
	}
	s.x[17] = uint64(uint32(s.x[21]) + uint32(s.x[30]))
	dispatchIndex := dispatchBit*uint32(s.x[17]) ^ uint32(s.x[1])
	next, err := loadBootstrapPointer(s, s.x[22], dispatchIndex)
	if err != nil {
		return 0, fmt.Errorf("bootstrap pointer-cycle next dispatch: %w", err)
	}
	s.x[1] = next
	s.x[9] = nextCounter
	s.x[17] = 0x1329
	return next, nil
}

func runBootstrapByteClassify(s *state) (uint64, error) {
	const maximumBytes = 4096
	for index := 0; index < maximumBytes; index++ {
		next, err := runBootstrapByteClassifyStep(s)
		if err != nil || next != bootstrapByteClassifyPass {
			return next, err
		}
	}
	return 0, fmt.Errorf("bootstrap byte classify exceeded %d bytes", maximumBytes)
}

func runBootstrapByteClassifyStep(s *state) (uint64, error) {
	s.traceBlock(bootstrapByteClassifyPass)
	s.x[13] = uint64(uint32(s.x[12]))
	s.x[17] = shiftLeft(s.x[13], int(s.x[15]), 64)
	s.x[17] &= s.x[16]
	s.x[0] = s.x[13] + s.x[11]
	s.x[13] = s.x[0] - s.x[17]
	s.x[17] = s.x[13] ^ s.x[11]
	value, err := s.load(s.x[23]+s.x[17], 1)
	if err != nil {
		return 0, fmt.Errorf("bootstrap byte-classify input: %w", err)
	}
	s.x[17] = uint64(byte(value))
	s.compare32(uint32(s.x[17]), 0x1b)
	selected := uint32(boolValue(s.condition("ne")))
	s.x[12] = uint64(uint32(s.x[12]) + 1)
	dispatchIndex := ((selected << 2) | selected) ^ uint32(s.x[14])
	next, err := loadBootstrapPointer(s, s.x[21], dispatchIndex)
	if err != nil {
		return 0, fmt.Errorf("bootstrap byte-classify dispatch: %w", err)
	}
	s.x[1] = next
	s.x[17] = 0xd618
	return next, nil
}
