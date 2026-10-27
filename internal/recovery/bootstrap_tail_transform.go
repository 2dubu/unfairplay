package recovery

import "fmt"

const (
	bootstrapStoreWordReturn   = uint64(0x100285604)
	bootstrapBitLookupDispatch = uint64(0x1003cece8)
	bootstrapSignedRange       = uint64(0x1005332bc)
	bootstrapFrameStatus       = uint64(0x10053dd44)
	bootstrapByteMerge         = uint64(0x100a17200)
	bootstrapBitSelect         = uint64(0x10006c768)
	bootstrapWordLookup        = uint64(0x100120a5c)
	bootstrapByteUpdate        = uint64(0x1004f10c4)
	bootstrapByteStore         = uint64(0x10051c16c)
	bootstrapPreparedDispatch  = uint64(0x100344c00)
	bootstrapStackChoice       = uint64(0x10037fc4c)
)

func runBootstrapStoreWordReturn(s *state) (uint64, error) {
	s.traceBlock(bootstrapStoreWordReturn)
	if err := s.storeValue(s.x[0], uint64(uint32(s.x[12])), 4); err != nil {
		return 0, fmt.Errorf("bootstrap store-word return: %w", err)
	}
	s.sp += 0x10
	return s.x[30], nil
}

func runBootstrapBitLookupDispatch(s *state) (uint64, error) {
	s.traceBlock(bootstrapBitLookupDispatch)
	base, err := s.load(s.x[8]+8, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap bit-lookup base: %w", err)
	}
	s.x[8] = base
	word, err := s.load(base+uint64(uint32(s.x[11]))*4, 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap bit-lookup word: %w", err)
	}
	s.x[8] = uint64(uint32(word))
	s.x[11] = 0x18d7
	s.x[11] = uint64(uint32(s.x[10]) ^ uint32(s.x[11]))
	s.x[9] = uint64(uint32(s.x[11]) & uint32(s.x[9]))
	s.x[11] = 0x224ea33b
	s.x[8] = uint64(uint32(s.x[8]) ^ uint32(s.x[11]))
	s.x[8] = uint64(uint32(shiftRight(s.x[8], int(s.x[9]), 32)))
	s.x[8] &= 1
	s.x[12] = uint64(uint32(s.x[8]) ^ uint32(s.x[12]))
	next, err := loadBootstrapPointer(s, s.x[13], uint32(s.x[10]))
	if err != nil {
		return 0, fmt.Errorf("bootstrap bit-lookup dispatch: %w", err)
	}
	s.x[8] = next
	s.x[17] = 0xe952
	return next, nil
}

func runBootstrapSignedRange(s *state) (uint64, error) {
	s.traceBlock(bootstrapSignedRange)
	s.x[8] = uint64(uint32(s.x[20]) - 1)
	s.x[9] = 0x36db5c1a
	s.x[10] = 0xc9247db7
	word, err := s.load(s.sp+0x6c, 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap signed-range stack word: %w", err)
	}
	s.x[26] = uint64(uint32(word))
	s.x[9] = uint64(uint32(s.x[26]) + uint32(s.x[9]))
	s.x[21] = uint64(uint32(s.x[9]) & uint32(s.x[10]))
	s.x[9] = 0x871d9549
	s.x[10] = 0x7fffffff
	s.x[9] = uint64(uint32(s.x[20]) + uint32(s.x[9]))
	s.compare32(uint32(s.x[9]), uint32(s.x[10]))
	gate := uint32(boolValue(s.condition("lo")))
	s.x[9] = uint64(gate)
	s.x[10] = 0xffffe6cd
	s.x[10] = uint64(uint32(s.x[21]) + uint32(s.x[10]))
	s.x[9] = uint64(uint32(s.x[10]) * gate)
	s.x[9] = uint64(uint32(s.x[9]) ^ uint32(s.x[26]))
	table, err := s.load(s.sp+0x70, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap signed-range table: %w", err)
	}
	s.x[10] = table
	next, err := loadBootstrapPointer(s, table, uint32(s.x[9]))
	if err != nil {
		return 0, fmt.Errorf("bootstrap signed-range dispatch: %w", err)
	}
	s.x[9] = next
	s.x[20] = s.x[8]
	s.x[17] = 0x5ca6
	return next, nil
}

func runBootstrapFrameStatus(s *state) (uint64, error) {
	s.traceBlock(bootstrapFrameStatus)
	value, err := s.load(addOffset(s.x[29], -0x90), 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap frame-status input: %w", err)
	}
	s.x[8] = uint64(uint32(value))
	s.x[9] = 0x5d55faf3
	s.compare32(uint32(s.x[8]), uint32(s.x[9]))
	gate := uint32(boolValue(s.condition("ne")))
	s.x[8] = uint64(gate)
	s.x[9] = 0x4a
	s.x[8] = uint64(gate * uint32(s.x[9]))
	s.x[8] = uint64(uint32(s.x[8]) ^ uint32(s.x[21]))
	next, err := loadBootstrapPointer(s, s.x[27], uint32(s.x[8]))
	if err != nil {
		return 0, fmt.Errorf("bootstrap frame-status dispatch: %w", err)
	}
	s.x[8] = next
	s.x[17] = 0x125c
	return next, nil
}

func runBootstrapByteMerge(s *state) (uint64, error) {
	s.traceBlock(bootstrapByteMerge)
	left, err := s.load(s.x[20]+uint64(uint32(s.x[11])), 1)
	if err != nil {
		return 0, fmt.Errorf("bootstrap byte-merge left: %w", err)
	}
	s.x[16] = uint64(byte(left))
	right, err := s.load(s.x[22], 1)
	if err != nil {
		return 0, fmt.Errorf("bootstrap byte-merge right: %w", err)
	}
	s.x[17] = uint64(byte(right))
	s.x[22]++
	s.x[16] = uint64(uint32(s.x[17]) ^ uint32(s.x[16]))
	s.x[8] = uint64(uint32(s.x[8]) &^ uint32(s.x[16]))
	s.x[16] = uint64(uint32(s.x[16]) & uint32(s.x[14]))
	s.x[8] = uint64(uint32(s.x[16]) | uint32(s.x[8]))
	s.x[11] = uint64(uint32(s.x[11]) + 1)
	s.compare64(s.x[13], 0)
	gate := uint32(boolValue(s.condition("eq")))
	s.x[16] = uint64(gate * uint32(s.x[15]))
	s.x[16] = uint64(uint32(s.x[16]) ^ uint32(s.x[12]))
	next, err := loadBootstrapPointer(s, s.x[30], uint32(s.x[16]))
	if err != nil {
		return 0, fmt.Errorf("bootstrap byte-merge dispatch: %w", err)
	}
	s.x[0] = next
	s.x[13]--
	s.x[17] = 0xd87d
	return next, nil
}

func runBootstrapBitSelect(s *state) (uint64, error) {
	s.traceBlock(bootstrapBitSelect)
	s.x[16] = uint64(uint32(shiftRight(s.x[13], int(s.x[14]), 32)))
	s.x[14] = uint64(uint32(s.x[14]) - 1)
	s.x[16] &= 1
	s.x[16] = uint64(uint32(s.x[16]) * uint32(s.x[15]))
	s.x[16] = uint64(uint32(s.x[16]) | uint32(s.x[10]))
	next, err := s.load(s.x[9]+uint64(uint32(s.x[16]))*8, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap bit-select dispatch: %w", err)
	}
	s.x[1] = next
	s.x[11] = uint64(uint32(s.x[11]) - 1)
	s.x[17] = 0x3edc
	return next, nil
}

func runBootstrapWordLookup(s *state) (uint64, error) {
	s.traceBlock(bootstrapWordLookup)
	s.x[17] = uint64(uint32(s.x[11]) + uint32(s.x[10]))
	s.x[17] = uint64(uint32(s.x[17]) + uint32(s.x[14]))
	word, err := s.load(s.x[13]+uint64(uint32(s.x[17]))*4, 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap word-lookup input: %w", err)
	}
	s.x[17] = uint64(uint32(word))
	s.compare32(uint32(s.x[17]), uint32(s.x[15]))
	gate := uint32(boolValue(s.condition("eq")))
	s.x[17] = uint64(gate * uint32(s.x[16]))
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[11]))
	next, err := loadBootstrapPointer(s, s.x[20], uint32(s.x[17]))
	if err != nil {
		return 0, fmt.Errorf("bootstrap word-lookup dispatch: %w", err)
	}
	s.x[2] = next
	s.x[17] = 0x5da3
	return next, nil
}

func runBootstrapByteUpdate(s *state) (uint64, error) {
	s.traceBlock(bootstrapByteUpdate)
	s.x[17] = s.x[21] + s.x[19]
	value, err := s.load(s.x[17], 1)
	if err != nil {
		return 0, fmt.Errorf("bootstrap byte-update input: %w", err)
	}
	s.x[17] = uint64(byte(value))
	s.x[24] = s.x[20] + s.x[19]
	s.x[25] = uint64(uint32(s.x[4]) & uint32(s.x[17]<<1))
	s.x[17] = uint64(uint32(s.x[17]) - uint32(s.x[25]))
	s.x[17] = uint64(uint32(s.x[17]) + 0x1b)
	if err := s.storeValue(s.x[24], uint64(byte(s.x[17])), 1); err != nil {
		return 0, fmt.Errorf("bootstrap byte-update output: %w", err)
	}
	s.compare64(s.x[19], s.x[23])
	gate := uint32(boolValue(s.condition("eq")))
	s.x[17] = uint64(gate * uint32(s.x[6]))
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[14]))
	table, err := s.load(s.sp+0xdf0, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap byte-update table: %w", err)
	}
	s.x[24] = table
	next, err := loadBootstrapPointer(s, table, uint32(s.x[17]))
	if err != nil {
		return 0, fmt.Errorf("bootstrap byte-update dispatch: %w", err)
	}
	s.x[24] = next
	s.x[19]--
	s.x[17] = 0x5cd7
	return next, nil
}

func runBootstrapByteStore(s *state) (uint64, error) {
	s.traceBlock(bootstrapByteStore)
	s.x[11] = uint64(uint32(s.x[9]) + uint32(s.x[30]))
	s.x[8]--
	if err := s.storeValue(s.sp+0x3af8, s.x[8], 8); err != nil {
		return 0, fmt.Errorf("bootstrap byte-store counter: %w", err)
	}
	s.x[12] = uint64(uint32(s.x[10]) ^ uint32(s.x[1]))
	s.x[13] = uint64(uint32(s.x[10]) + 0x11)
	s.x[12] = uint64(uint32(s.x[12]) * uint32(s.x[13]))
	if err := s.storeValue(s.x[20]+s.x[8], uint64(byte(s.x[12])), 1); err != nil {
		return 0, fmt.Errorf("bootstrap byte-store output: %w", err)
	}
	s.x[8] = s.x[10] + s.x[3]
	s.x[10] = s.x[8] + s.x[11]
	counter, err := s.load(s.sp+0x3af8, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap byte-store saved counter: %w", err)
	}
	s.x[8] = counter
	s.compare64(counter, 0)
	gate := uint32(boolValue(s.condition("eq")))
	s.x[11] = uint64(uint32(s.x[9]) ^ gate<<3)
	table, err := s.load(s.sp+0xdf0, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap byte-store table: %w", err)
	}
	s.x[12] = table
	next, err := loadBootstrapPointer(s, table, uint32(s.x[11]))
	if err != nil {
		return 0, fmt.Errorf("bootstrap byte-store dispatch: %w", err)
	}
	s.x[11] = next
	s.x[17] = 0x3548
	return next, nil
}

func runBootstrapPreparedDispatch(s *state) (uint64, error) {
	s.traceBlock(bootstrapPreparedDispatch)
	s.x[13] = 0xffffffff98245355
	s.x[11] += s.x[13]
	next, err := loadBootstrapPointer(s, s.x[11], uint32(s.x[12]))
	if err != nil {
		return 0, fmt.Errorf("bootstrap prepared dispatch: %w", err)
	}
	s.x[19] = next
	s.x[12] = 0xffffef9b
	s.x[13] = 0x1040
	s.x[14] = 0xffffed0b
	s.x[15] = 0x48cd
	s.x[16] = 0x94cd
	s.x[1] = 0xffffe43f
	s.x[2] = 0x8cc
	s.x[3] = 0x5bc5
	s.x[4] = 0xd21
	s.x[5] = 0xffffffffffffa418
	s.x[6] = 0x224ea33b
	s.x[7] = 0xffffb758
	s.x[17] = 0xced6
	return next, nil
}

func runBootstrapStackChoice(s *state) (uint64, error) {
	s.traceBlock(bootstrapStackChoice)
	value, err := s.load(s.sp+0x20, 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap stack-choice input: %w", err)
	}
	s.x[8] = uint64(uint32(value))
	s.x[9] = 0x06d1a0bc6
	s.compare32(uint32(s.x[8]), uint32(s.x[9]))
	gate := uint32(boolValue(s.condition("eq")))
	s.x[8] = uint64(gate * 0x6044)
	s.x[8] = uint64(uint32(s.x[8]) ^ uint32(s.x[20]))
	next, err := loadBootstrapPointer(s, s.x[22], uint32(s.x[8]))
	if err != nil {
		return 0, fmt.Errorf("bootstrap stack-choice dispatch: %w", err)
	}
	s.x[9] = next
	s.x[8] = 0x4ffc9144
	s.x[17] = 0xb89b
	return next, nil
}
