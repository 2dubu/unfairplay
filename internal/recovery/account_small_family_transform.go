package recovery

import "fmt"

const (
	accountSmallPairStart = uint64(0x1001453c0)
	accountSmallPairTail  = uint64(0x100145400)
	accountByteSelect     = uint64(0x10016e6f0)
	accountRangeDispatch  = uint64(0x100196998)
	accountStoreReturn    = uint64(0x100197310)
	accountOrderedSelect  = uint64(0x10019843c)
	accountDescriptorMove = uint64(0x10021eb98)
	accountLimitDispatch  = uint64(0x10030af80)
)

func runAccountSmallPair(s *state) (uint64, error) {
	s.traceBlock(accountSmallPairStart)
	s.x[14] = uint64(uint32(s.x[13]) + 0xd94b841a)
	s.x[15] = 0xe628
	s.x[2] = uint64(uint32(s.x[1]) ^ uint32(s.x[15]))
	s.x[16] = uint64(uint32(s.x[2]) + 0x4b2)
	s.x[15] = uint64(uint32(s.x[13]) + 0xd94b8419)
	value, err := s.load(s.x[11]+uint64(uint32(s.x[15])), 1)
	if err != nil {
		return 0, fmt.Errorf("account small-pair input: %w", err)
	}
	s.x[17] = uint64(byte(value))
	s.x[15] = 0x18a
	s.x[15] = uint64(uint32(s.x[15]) & uint32(s.x[17]<<1))
	s.x[3] = 0xce757bc5
	s.x[15] = uint64(uint32(s.x[17]) - uint32(s.x[15]))
	s.x[15] = uint64(uint32(s.x[15]) + uint32(s.x[3]))

	s.traceBlock(accountSmallPairTail)
	s.x[17] = uint64(uint32(s.x[17]) >> 7)
	s.x[17] = uint64(uint32(s.x[17]) ^ 1)
	s.x[3] = uint64(uint32(s.x[16]) ^ 0x2c85)
	s.x[17] = uint64(uint32(s.x[17]) * uint32(s.x[3]))
	s.x[17] = uint64(uint32(s.x[17]) | uint32(s.x[1]))
	next, err := loadBootstrapPointer(s, s.x[10], uint32(s.x[17]))
	if err != nil {
		return 0, fmt.Errorf("account small-pair target: %w", err)
	}
	s.x[1] = next
	s.x[17] = 0x68bc
	return next, nil
}

func runAccountByteSelect(s *state) (uint64, error) {
	s.traceBlock(accountByteSelect)
	base, err := s.load(s.x[0]+8, 8)
	if err != nil {
		return 0, fmt.Errorf("account byte-select base: %w", err)
	}
	s.x[11] = base
	s.x[12] = 0x6027
	s.x[16] = 0xd94b3011
	s.x[15] = uint64(uint32(s.x[14]) + uint32(s.x[12]))
	s.x[12] = uint64(uint32(s.x[15]) ^ uint32(s.x[16]))
	s.x[12] = uint64(uint32(s.x[12]) + uint32(s.x[13]))
	value, err := s.load(base+uint64(uint32(s.x[12])), 1)
	if err != nil {
		return 0, fmt.Errorf("account byte-select input: %w", err)
	}
	s.x[12] = uint64(byte(value))
	s.x[16] = 0x166
	s.x[16] = uint64(uint32(s.x[16]) & uint32(s.x[12]<<1))
	s.x[17] = 0x04e043b3
	s.x[12] = uint64(uint32(s.x[12]) - uint32(s.x[16]))
	s.x[12] = uint64(uint32(s.x[12]) + uint32(s.x[17]))
	s.x[16] = uint64(uint32(s.x[12]) & 0x1f)
	s.compare32(uint32(s.x[16]), 0xc)
	gate := uint32(boolValue(s.condition("eq")))
	s.x[16] = uint64(gate * 0x7eb)
	s.x[14] = uint64(uint32(s.x[16]) ^ uint32(s.x[14]))
	next, err := loadBootstrapPointer(s, s.x[10], uint32(s.x[14]))
	if err != nil {
		return 0, fmt.Errorf("account byte-select target: %w", err)
	}
	s.x[14] = next
	s.x[17] = 0xa133
	return next, nil
}

func runAccountRangeDispatch(s *state) (uint64, error) {
	s.traceBlock(accountRangeDispatch)
	s.x[17] = uint64(uint32(s.x[13]) ^ 0x27af)
	s.x[16] = uint64(uint32(s.x[17]) * 0xd5)
	s.x[2] = uint64(uint32(s.x[16]) + 0x410fa5df)
	s.x[2] = uint64(uint32(s.x[1]) + uint32(s.x[2]))
	s.x[3] = 0x1615efec
	s.x[4] = 0x729a7e23
	s.compare32(uint32(s.x[2]), uint32(s.x[4]))
	first := uint32(boolValue(s.condition("lo")))
	s.x[5] = uint64(first)
	s.x[9] = uint64(uint32(s.x[9]) + uint32(s.x[3]))
	s.compare32(uint32(s.x[9]), uint32(s.x[4]))
	second := uint32(boolValue(s.condition("lo")))
	s.x[3] = uint64(second)
	s.x[3] = uint64(second ^ first)
	s.compare32(uint32(s.x[2]), uint32(s.x[9]))
	ordered := uint32(boolValue(s.condition("hi")))
	s.x[9] = uint64(ordered)
	s.compare32(uint32(s.x[3]), 0)
	selected := ordered
	if s.condition("ne") {
		selected = first
	}
	s.x[9] = uint64(selected)
	s.x[17] = uint64(uint32(s.x[17]) - 0x3f)
	s.x[9] = uint64(uint32(s.x[9]) * uint32(s.x[17]))
	s.x[9] = uint64(uint32(s.x[9]) ^ uint32(s.x[13]))
	next, err := loadBootstrapPointer(s, s.x[10], uint32(s.x[9]))
	if err != nil {
		return 0, fmt.Errorf("account range-dispatch target: %w", err)
	}
	s.x[9] = next
	s.x[2] = 0x318a201b
	s.x[17] = 0x3cce
	return next, nil
}

func runAccountStoreReturn(s *state) (uint64, error) {
	s.traceBlock(accountStoreReturn)
	if err := s.storeValue(s.x[0]+0x10, uint64(uint32(s.x[2])), 4); err != nil {
		return 0, fmt.Errorf("account store-return output: %w", err)
	}
	s.sp += 0x10
	return s.x[30], nil
}

func runAccountOrderedSelect(s *state) (uint64, error) {
	s.traceBlock(accountOrderedSelect)
	s.x[13] = 0x1a9a8340
	s.x[17] = 0xe56577eb
	s.x[13] = uint64(uint32(s.x[16]) + uint32(s.x[13]))
	s.x[13] = uint64(uint32(s.x[13]) & uint32(s.x[17]))
	s.x[17] = uint64(uint32(s.x[15]) ^ 0x57b74ba1)
	s.x[1] = 0xa848b45f
	s.x[17] = uint64(uint32(s.x[17]) + uint32(s.x[1]))
	s.x[1] = uint64(uint32(s.x[15]) ^ 0xa9ac3165)
	s.x[2] = 0x5653ce9b
	s.x[1] = uint64(uint32(s.x[1]) + uint32(s.x[2]))
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[1]))
	s.x[1] = uint64(uint32(s.x[13]) ^ 0x306e26cb)
	s.x[1] = uint64(uint32(s.x[1]) ^ uint32(s.x[15]))
	s.x[2] = 0xcf91feff
	s.x[1] = uint64(uint32(s.x[1]) + uint32(s.x[2]))
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[1]))
	s.x[2] = 0x7ba7120e
	s.x[3] = 0x4a1c8dd3
	s.x[1] = uint64(uint32(s.x[17]) + uint32(s.x[14]))
	s.x[17] = uint64(uint32(s.x[1]) + uint32(s.x[3]))
	s.compare32(uint32(s.x[14]), uint32(s.x[3]))
	first := uint32(boolValue(s.condition("hi")))
	s.x[3] = uint64(first)
	s.compare32(uint32(s.x[17]), uint32(s.x[2]))
	second := uint32(boolValue(s.condition("lo")))
	s.x[4] = uint64(second)
	s.x[2] = uint64(uint32(s.x[14]) + uint32(s.x[2]))
	s.compare32(uint32(s.x[2]), uint32(s.x[17]))
	s.x[17] = 0x8d19
	s.x[2] = uint64(first ^ second)
	ordered := uint32(boolValue(s.condition("hi")))
	s.x[4] = uint64(ordered)
	s.compare32(uint32(s.x[2]), 0)
	selected := ordered
	if s.condition("ne") {
		selected = first
	}
	s.x[2] = uint64(selected)
	s.x[3] = uint64(selected << 1)
	s.x[2] = uint64(uint32(s.x[3]) | selected<<3)
	s.x[16] = uint64(uint32(s.x[16]) + uint32(s.x[17]))
	s.x[16] = uint64(uint32(s.x[2]) ^ uint32(s.x[16]))
	next, err := loadBootstrapPointer(s, s.x[10], uint32(s.x[16]))
	if err != nil {
		return 0, fmt.Errorf("account ordered-select target: %w", err)
	}
	s.x[3] = next
	s.x[2] = 0x318a201b
	s.x[17] = 0x564e
	return next, nil
}

func runAccountDescriptorMove(s *state) (uint64, error) {
	s.traceBlock(accountDescriptorMove)
	s.x[9] = uint64(uint32(s.x[1]) + 0xf529f7ad)
	first, err := s.load(s.x[0]+0x18, 8)
	if err != nil {
		return 0, fmt.Errorf("account descriptor-move first pointer: %w", err)
	}
	second, err := s.load(s.x[0]+0x20, 8)
	if err != nil {
		return 0, fmt.Errorf("account descriptor-move second pointer: %w", err)
	}
	s.x[17], s.x[13] = first, second
	value, err := s.load(s.x[0]+0x30, 8)
	if err != nil {
		return 0, fmt.Errorf("account descriptor-move third pointer: %w", err)
	}
	s.x[1] = value
	if err := s.storeValue(s.x[8], uint64(uint32(s.x[9])), 4); err != nil {
		return 0, fmt.Errorf("account descriptor-move first output: %w", err)
	}
	if err := s.storeValue(s.x[1], uint64(uint32(s.x[12])), 4); err != nil {
		return 0, fmt.Errorf("account descriptor-move second output: %w", err)
	}
	if err := s.storeValue(s.x[17], uint64(uint32(s.x[15])), 4); err != nil {
		return 0, fmt.Errorf("account descriptor-move third output: %w", err)
	}
	s.x[8] = s.x[11] + uint64(uint32(s.x[14]))
	if err := s.storeValue(s.x[13], s.x[8], 8); err != nil {
		return 0, fmt.Errorf("account descriptor-move pointer output: %w", err)
	}
	next, err := loadBootstrapPointer(s, s.x[10], uint32(s.x[16]))
	if err != nil {
		return 0, fmt.Errorf("account descriptor-move target: %w", err)
	}
	s.x[8] = next
	s.x[2] = 0xce757bc5
	s.x[17] = 0xa133
	return next, nil
}

func runAccountLimitDispatch(s *state) (uint64, error) {
	s.traceBlock(accountLimitDispatch)
	s.x[14] = 0x75cd
	s.x[1] = uint64(uint32(s.x[15]) ^ uint32(s.x[14]))
	s.x[14] = 0xffff6ada
	s.x[16] = 0x72410ca4
	s.x[17] = 0x3c70d6cd
	s.x[14] = uint64(uint32(s.x[1]) + uint32(s.x[14]))
	s.x[14] = uint64(uint32(s.x[14]) ^ uint32(s.x[17]))
	s.x[17] = 0x98f5888a
	s.x[16] = uint64(uint32(s.x[13]) + uint32(s.x[16]))
	s.compare32(uint32(s.x[16]), uint32(s.x[17]))
	first := uint32(boolValue(s.condition("lo")))
	s.x[2] = uint64(first)
	s.x[14] = uint64(uint32(s.x[14]) + uint32(s.x[9]))
	s.compare32(uint32(s.x[14]), uint32(s.x[17]))
	second := uint32(boolValue(s.condition("lo")))
	s.x[17] = uint64(second)
	s.x[17] = uint64(first ^ second)
	s.compare32(uint32(s.x[16]), uint32(s.x[14]))
	ordered := uint32(boolValue(s.condition("hi")))
	s.x[14] = uint64(ordered)
	s.compare32(uint32(s.x[17]), 0)
	selected := ordered
	if s.condition("ne") {
		selected = first
	}
	s.x[14] = uint64(selected ^ 1)
	s.x[16] = 0x1d
	s.x[14] = uint64(uint32(s.x[14]) * uint32(s.x[16]))
	s.x[14] = uint64(uint32(s.x[14]) ^ uint32(s.x[15]))
	next, err := loadBootstrapPointer(s, s.x[10], uint32(s.x[14]))
	if err != nil {
		return 0, fmt.Errorf("account limit-dispatch target: %w", err)
	}
	s.x[14] = next
	s.x[2] = 0x318a201b
	s.x[17] = 0xaae4
	return next, nil
}
