package recovery

import "fmt"

const (
	bootstrapBoundsDispatch      = uint64(0x1008eb978)
	bootstrapLimitCompareStart   = uint64(0x1008f6fe0)
	bootstrapLimitCompareTail    = uint64(0x1008f7000)
	bootstrapCounterAdvance      = uint64(0x100a024c8)
	bootstrapZeroByteDispatch    = uint64(0x100a10258)
	bootstrapPairReturn          = uint64(0x10004a38c)
	bootstrapIndexedCompare      = uint64(0x1000caaa0)
	bootstrapLimitCompare        = uint64(0x100115b44)
	bootstrapGuardReturnStart    = uint64(0x1002a8520)
	bootstrapGuardReturnTail     = uint64(0x1002a853c)
	bootstrapGuardReturnMismatch = uint64(0x1002a8550)
	bootstrapSignedChoice        = uint64(0x1002cb908)
)

func runBootstrapBoundsDispatch(s *state) (uint64, error) {
	s.traceBlock(bootstrapBoundsDispatch)
	s.compare32(uint32(s.x[14]), uint32(s.x[19]))
	first := uint32(boolValue(s.condition("lo")))
	s.x[17] = uint64(first)
	s.x[26] = uint64(uint32(s.x[12]) + uint32(s.x[23]))
	s.x[3] = uint64((uint32(s.x[12]) + uint32(s.x[22])) & uint32(s.x[25]))
	s.x[6] = uint64(uint32(s.x[7]) + uint32(s.x[3]))
	s.compare32(uint32(s.x[6]), uint32(s.x[19]))
	second := uint32(boolValue(s.condition("lo")))
	s.x[16] = uint64(first ^ second)
	s.compare32(uint32(s.x[14]), uint32(s.x[6]))
	ordered := uint32(boolValue(s.condition("hi")))
	s.x[6] = uint64(ordered)
	s.compare32(uint32(s.x[16]), 0)
	selected := ordered
	if s.condition("ne") {
		selected = first
	}
	s.x[16] = uint64((selected ^ 1) * uint32(s.x[27]))
	s.x[12] = uint64(uint32(s.x[12]) ^ uint32(s.x[16]))
	next, err := loadBootstrapPointer(s, s.x[13], uint32(s.x[12]))
	if err != nil {
		return 0, fmt.Errorf("bootstrap bounds dispatch: %w", err)
	}
	s.x[12] = next
	s.x[17] = 0x823b
	return next, nil
}

func runBootstrapLimitCompare(s *state) (uint64, error) {
	s.traceBlock(bootstrapLimitCompareStart)
	s.compare32(uint32(s.x[11]), uint32(s.x[28]))
	first := uint32(boolValue(s.condition("lo")))
	s.x[16] = uint64(first)
	s.x[12] = uint64(uint32(s.x[3]) + 0xffffc6bc)
	s.x[12] = uint64(uint32(s.x[12]) | 0x3022)
	s.x[17] = 0x9ece95f9

	s.traceBlock(bootstrapLimitCompareTail)
	s.x[17] = uint64(uint32(s.x[12]) ^ uint32(s.x[17]))
	s.x[17] = uint64(uint32(s.x[17]) + uint32(s.x[1]))
	s.compare32(uint32(s.x[17]), uint32(s.x[28]))
	second := uint32(boolValue(s.condition("lo")))
	s.x[6] = uint64(first ^ second)
	s.compare32(uint32(s.x[17]), uint32(s.x[11]))
	ordered := uint32(boolValue(s.condition("lo")))
	s.x[17] = uint64(ordered)
	s.compare32(uint32(s.x[6]), 0)
	selected := ordered
	if s.condition("ne") {
		selected = first
	}
	s.x[16] = uint64(selected * 0xc1)
	s.x[16] = uint64(uint32(s.x[16]) ^ uint32(s.x[3]))
	next, err := s.load(s.x[13]+uint64(uint32(s.x[16]))*8, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap limit-compare dispatch: %w", err)
	}
	s.x[6] = next
	s.x[17] = 0x93d3
	return next, nil
}

func runBootstrapCounterAdvance(s *state) (uint64, error) {
	s.traceBlock(bootstrapCounterAdvance)
	s.x[12] = uint64((uint32(s.x[26]) + uint32(s.x[4])) | uint32(s.x[30]))
	s.x[1] = uint64(uint32(s.x[1]) + 1)
	s.compare32(uint32(s.x[26]), uint32(s.x[5]))
	gate := uint32(boolValue(s.condition("eq")))
	s.x[16] = uint64(gate)
	s.x[17] = uint64(uint32(s.x[12]) + uint32(s.x[0]))
	s.x[16] = uint64(uint32(s.x[17]) * gate)
	s.x[17] = uint64(uint32(s.x[26]) ^ uint32(s.x[2]))
	s.x[16] = uint64(uint32(s.x[17]) ^ uint32(s.x[16]))
	next, err := loadBootstrapPointer(s, s.x[13], uint32(s.x[16]))
	if err != nil {
		return 0, fmt.Errorf("bootstrap counter-advance dispatch: %w", err)
	}
	s.x[3] = next
	s.x[7] = uint64(uint32(s.x[7]) + 1)
	s.x[10] = uint64(uint32(s.x[10]) + 0x25)
	s.x[17] = 0x34a3
	return next, nil
}

func runBootstrapZeroByteDispatch(s *state) (uint64, error) {
	s.traceBlock(bootstrapZeroByteDispatch)
	s.x[6] = uint64(uint32(s.x[13]) + uint32(s.x[3]))
	s.x[10] = uint64(uint32(s.x[11]) & 0xff)
	s.x[11] = uint64(uint32(s.x[11]) + 0x101)
	s.x[10] = uint64(uint32(s.x[10]) * uint32(s.x[12]))
	divisor := uint32(s.x[14])
	quotient := uint32(0)
	if divisor != 0 {
		quotient = uint32(s.x[10]) / divisor
	}
	s.x[16] = uint64(quotient)
	s.x[10] = uint64(uint32(s.x[10]) - quotient*divisor)
	s.x[16] = uint64(uint32(s.x[11]) & 0xfc)
	s.compare32(uint32(s.x[16]), 0x3c)
	gate := uint32(boolValue(s.condition("lo")))
	s.x[16] = uint64(gate)
	if err := s.storeValue(s.x[8]+uint64(uint32(s.x[10])), 0, 1); err != nil {
		return 0, fmt.Errorf("bootstrap zero-byte output: %w", err)
	}
	s.x[10] = uint64(gate * uint32(s.x[4]))
	s.x[10] = uint64(uint32(s.x[10]) ^ uint32(s.x[13]))
	next, err := loadBootstrapPointer(s, s.x[30], uint32(s.x[10]))
	if err != nil {
		return 0, fmt.Errorf("bootstrap zero-byte dispatch: %w", err)
	}
	s.x[0] = next
	s.x[10] = 1
	s.x[17] = 0xc004
	return next, nil
}

func runBootstrapPairReturn(s *state) (uint64, error) {
	s.traceBlock(bootstrapPairReturn)
	if err := s.storeValue(s.x[0]+0x10, uint64(uint32(s.x[20])), 4); err != nil {
		return 0, fmt.Errorf("bootstrap pair-return output: %w", err)
	}
	first, err := s.load(s.sp+0x10, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap pair-return register 20: %w", err)
	}
	second, err := s.load(s.sp+0x18, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap pair-return register 19: %w", err)
	}
	s.x[20], s.x[19] = first, second
	s.sp += 0x20
	return s.x[30], nil
}

func runBootstrapIndexedCompare(s *state) (uint64, error) {
	s.traceBlock(bootstrapIndexedCompare)
	s.x[9] = uint64(uint32(s.x[19]) + uint32(s.x[1]))
	s.x[17] = uint64(uint32(s.x[9]) | uint32(s.x[2]))
	s.x[9] = uint64((uint32(s.x[19]) ^ uint32(s.x[3])) * uint32(s.x[4]))
	base, err := s.load(s.x[8]+8, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap indexed-compare base: %w", err)
	}
	s.x[20] = base
	s.x[10] += s.x[5]
	s.x[10] += s.x[9]
	value, err := s.load(base+(s.x[10]<<2), 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap indexed-compare word: %w", err)
	}
	s.x[20] = uint64(uint32(value))
	s.compare32(uint32(s.x[20]), uint32(s.x[6]))
	gate := uint32(boolValue(s.condition("ne")))
	s.x[20] = uint64(gate)
	s.x[17] = uint64(uint32(s.x[17]) + uint32(s.x[7]))
	s.x[17] = uint64(gate * uint32(s.x[17]))
	s.x[17] = uint64(uint32(s.x[17]) | uint32(s.x[19]))
	next, err := loadBootstrapPointer(s, s.x[11], uint32(s.x[17]))
	if err != nil {
		return 0, fmt.Errorf("bootstrap indexed-compare dispatch: %w", err)
	}
	s.x[19] = next
	s.x[20] = 0x06d1a0bc6
	s.x[17] = 0x9e46
	return next, nil
}

func runBootstrapLimitCompareDispatch(s *state) (uint64, error) {
	s.traceBlock(bootstrapLimitCompare)
	s.x[17] = uint64(uint32(s.x[9]) + uint32(s.x[12]))
	s.x[19] = uint64(uint32(s.x[17]) | uint32(s.x[13]))
	s.x[17] = uint64(uint32(s.x[19]) + uint32(s.x[14]))
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[15]))
	s.compare32(uint32(s.x[10]), uint32(s.x[17]))
	gate := uint32(boolValue(s.condition("gt")))
	s.x[17] = uint64(gate * uint32(s.x[16]))
	s.x[9] = uint64(uint32(s.x[9]) ^ uint32(s.x[17]))
	next, err := loadBootstrapPointer(s, s.x[11], uint32(s.x[9]))
	if err != nil {
		return 0, fmt.Errorf("bootstrap limit dispatch: %w", err)
	}
	s.x[9] = next
	s.x[20] = 0x06d1a0bc7
	s.x[17] = 0x9e46
	return next, nil
}

func runBootstrapGuardReturn(s *state) (uint64, error) {
	s.traceBlock(bootstrapGuardReturnStart)
	if err := s.storeValue(s.x[19]+4, uint64(uint32(s.x[8])), 4); err != nil {
		return 0, fmt.Errorf("bootstrap guard-return output: %w", err)
	}
	saved, err := s.load(s.sp+0x28, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap guard-return saved value: %w", err)
	}
	s.x[8] = saved
	guardPointer, err := s.load(0x10106d000+0x408, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap guard-return guard pointer: %w", err)
	}
	s.x[9] = guardPointer
	guard, err := s.load(guardPointer, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap guard-return guard: %w", err)
	}
	s.x[9] = guard
	s.compare64(guard, saved)
	if s.condition("ne") {
		return bootstrapGuardReturnMismatch, nil
	}

	s.traceBlock(bootstrapGuardReturnTail)
	for _, pair := range [...]struct {
		first, second int
		offset        uint64
	}{{29, 30, 0x50}, {20, 19, 0x40}, {22, 21, 0x30}} {
		first, err := s.load(s.sp+pair.offset, 8)
		if err != nil {
			return 0, fmt.Errorf("bootstrap guard-return register %d: %w", pair.first, err)
		}
		second, err := s.load(s.sp+pair.offset+8, 8)
		if err != nil {
			return 0, fmt.Errorf("bootstrap guard-return register %d: %w", pair.second, err)
		}
		s.x[pair.first], s.x[pair.second] = first, second
	}
	s.sp += 0x60
	return s.x[30], nil
}

func runBootstrapSignedChoice(s *state) (uint64, error) {
	s.traceBlock(bootstrapSignedChoice)
	value, err := s.load(s.x[21]+4, 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap signed-choice input: %w", err)
	}
	s.x[8] = uint64(uint32(value))
	s.x[9] = 0x842a8d92
	s.x[10] = 0x7fffffff
	s.x[8] = uint64(uint32(s.x[8]) + uint32(s.x[9]))
	s.compare32(uint32(s.x[8]), uint32(s.x[10]))
	s.x[8] = 0x4ffc9143
	s.x[9] = 0x4ffc9145
	if s.condition("lo") {
		s.x[8] = s.x[9]
	}
	next, err := loadBootstrapPointer(s, s.x[22], uint32(s.x[20]))
	if err != nil {
		return 0, fmt.Errorf("bootstrap signed-choice dispatch: %w", err)
	}
	s.x[9] = next
	s.x[17] = 0xb89b
	return next, nil
}
