package recovery

import "fmt"

const (
	accountBooleanWriterA = uint64(0x1008f9a90)
	accountBooleanWriterB = uint64(0x1000b419c)
	accountBooleanWriterC = uint64(0x100136268)
	accountBooleanWriterD = uint64(0x10013e894)
	accountSmallSelectorA = uint64(0x1009e231c)
	accountSmallSelectorB = uint64(0x1009f7600)
	accountBoundsSelector = uint64(0x100907110)
)

func runAccountBooleanWriterA(s *state) (uint64, error) {
	s.traceBlock(accountBooleanWriterA)
	s.x[12] = uint64(uint32(s.x[11]) ^ uint32(s.x[1]))
	s.x[13] = uint64(uint32(s.x[11]) & uint32(s.x[2]))
	s.x[14] = uint64(uint32(s.x[11]) ^ uint32(s.x[3]))
	s.x[14] = uint64(uint32(s.x[14]) & uint32(s.x[13]<<1))
	s.x[13] = uint64(uint32(s.x[14]) ^ uint32(s.x[13]))
	s.x[14] = uint64(uint32(s.x[12]) << 1)
	s.x[15] = uint64(uint32(s.x[14]) ^ uint32(s.x[4]))
	s.x[15] = uint64(uint32(s.x[0]) ^ uint32(s.x[15]))
	s.x[12] = uint64(uint32(s.x[15]) & uint32(s.x[12]))
	s.x[14] = uint64(uint32(s.x[14]) & uint32(s.x[19]))
	s.x[12] = uint64(uint32(s.x[12]) ^ uint32(s.x[14]))
	s.x[14] = uint64(uint32(s.x[12]) ^ uint32(s.x[20]))
	s.x[12] = uint64(uint32(s.x[12]) ^ uint32(s.x[21]))
	s.x[12] = uint64(uint32(s.x[12]) & uint32(s.x[13]<<2))
	s.x[12] = uint64(uint32(s.x[12]) ^ uint32(s.x[13]))
	s.x[13] = uint64(uint32(s.x[14]) << 2)
	s.x[15] = uint64(uint32(s.x[13]) ^ uint32(s.x[22]))
	s.x[14] = uint64(uint32(s.x[15]) & uint32(s.x[14]))
	s.x[13] = uint64(uint32(s.x[13]) & uint32(s.x[25]))
	s.x[13] = uint64(uint32(s.x[13]) ^ uint32(s.x[27]))
	s.x[13] = uint64(uint32(s.x[13]) ^ uint32(s.x[14]))
	s.x[13] = uint64(uint32(s.x[13]) & uint32(s.x[12]<<4))
	s.x[12] = uint64(uint32(s.x[13]) ^ uint32(s.x[12]))
	s.x[12] = uint64(uint32(s.x[11]) ^ uint32(s.x[12]<<1))
	s.x[12] = uint64(uint32(s.x[12]) ^ uint32(s.x[28]))
	s.x[12] = uint64(uint32(s.x[12]) * uint32(s.x[10]))
	s.x[13] = s.x[23] + s.x[8]
	if err := s.storeValue(s.x[13], uint64(byte(s.x[12])), 1); err != nil {
		return 0, fmt.Errorf("account boolean-writer A output: %w", err)
	}
	s.x[11]++
	s.compare64(s.x[8], 0)
	gate := uint32(boolValue(s.condition("eq")))
	s.x[12] = uint64(gate * uint32(s.x[26]))
	s.x[12] = uint64(uint32(s.x[12]) ^ uint32(s.x[16]))
	next, err := loadBootstrapPointer(s, s.x[30], uint32(s.x[12]))
	if err != nil {
		return 0, fmt.Errorf("account boolean-writer A target: %w", err)
	}
	s.x[12] = next
	s.x[8]--
	s.x[10]++
	s.x[17] = 0x4f15
	return next, nil
}

func runAccountBooleanWriterB(s *state) (uint64, error) {
	s.traceBlock(accountBooleanWriterB)
	s.x[17] = uint64(uint32(s.x[9]) ^ uint32(s.x[14]))
	s.x[25] = uint64(uint32(s.x[9]) &^ uint32(s.x[10]))
	s.x[26] = uint64(uint32(s.x[10]) & uint32(s.x[15]))
	s.x[25] = uint64(uint32(s.x[26]) ^ uint32(s.x[25]))
	s.x[26] = uint64(uint32(s.x[17]) << 1)
	s.x[27] = uint64(uint32(s.x[26]) ^ uint32(s.x[16]))
	s.x[17] = uint64(uint32(s.x[27]) & uint32(s.x[17]))
	s.x[26] = uint64(uint32(s.x[26]) & uint32(s.x[0]))
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[26]))
	s.x[26] = uint64(uint32(s.x[17]) ^ uint32(s.x[1]))
	s.x[17] = uint64(uint32(s.x[17]) & uint32(s.x[25]<<2))
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[25]))
	s.x[25] = uint64(uint32(s.x[26]) << 2)
	s.x[27] = uint64(uint32(s.x[25]) ^ uint32(s.x[2]))
	s.x[26] = uint64(uint32(s.x[27]) & uint32(s.x[26]))
	s.x[25] = uint64(uint32(s.x[25]) & uint32(s.x[3]))
	s.x[25] = uint64(uint32(s.x[26]) ^ uint32(s.x[25]))
	s.x[25] = uint64(uint32(s.x[25]) & uint32(s.x[17]<<4))
	s.x[17] = uint64(uint32(s.x[25]) ^ uint32(s.x[17]))
	s.x[17] = uint64(uint32(s.x[9]) ^ uint32(s.x[17]<<1))
	s.x[25] = uint64(uint32(s.x[17]) ^ uint32(s.x[4]))
	s.x[25] = uint64(uint32(s.x[25]) + uint32(s.x[5]))
	s.x[26] = uint64(uint32(s.x[17]) ^ uint32(s.x[6]))
	s.x[26] = uint64(uint32(s.x[26]) + uint32(s.x[7]))
	s.x[25] = uint64(uint32(s.x[25]) ^ uint32(s.x[26]))
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[19]))
	s.x[17] = uint64(uint32(s.x[17]) + uint32(s.x[20]))
	s.x[17] = uint64(uint32(s.x[25]) ^ uint32(s.x[17]))
	s.x[17] = uint64(uint32(s.x[17]) + uint32(s.x[22]))
	s.x[17] = uint64(uint32(s.x[11])*uint32(s.x[17]) + uint32(s.x[12]))
	s.x[25] = s.x[23] + s.x[13]
	if err := s.storeValue(s.x[25], uint64(byte(s.x[17])), 1); err != nil {
		return 0, fmt.Errorf("account boolean-writer B output: %w", err)
	}
	s.x[9]++
	s.compare64(s.x[13], 0)
	gate := uint32(boolValue(s.condition("eq")))
	s.x[17] = uint64(gate * uint32(s.x[24]))
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[8]))
	next, err := loadBootstrapPointer(s, s.x[21], uint32(s.x[17]))
	if err != nil {
		return 0, fmt.Errorf("account boolean-writer B target: %w", err)
	}
	s.x[25] = next
	s.x[13]--
	s.x[10] += 2
	s.x[11]++
	s.x[12] += 0x8e
	s.x[17] = 0xe762
	return next, nil
}

func runAccountBooleanWriterC(s *state) (uint64, error) {
	s.traceBlock(accountBooleanWriterC)
	s.x[17] = uint64(uint32(s.x[11]) ^ uint32(s.x[14]))
	s.x[25] = uint64(uint32(s.x[11]) ^ uint32(s.x[16]))
	s.x[25] = uint64(uint32(s.x[25]) + uint32(s.x[0]))
	s.x[26] = uint64(uint32(s.x[11]) ^ uint32(s.x[1]))
	s.x[27] = uint64(uint32(s.x[11]) ^ uint32(s.x[3]))
	s.x[17] = uint64(uint32(s.x[17]) + uint32(s.x[15]))
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[11]))
	s.x[27] = uint64(uint32(s.x[27]) + uint32(s.x[4]))
	s.x[25] = uint64(uint32(s.x[25]) ^ uint32(s.x[27]))
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[25]))
	s.x[25] = uint64(uint32(s.x[2]) + uint32(s.x[26]))
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[25]))
	s.x[25] = uint64(uint32(s.x[17]) ^ uint32(s.x[5]))
	s.x[25] = uint64(uint32(s.x[25]) + uint32(s.x[6]))
	s.x[26] = uint64(uint32(s.x[17]) ^ uint32(s.x[7]))
	s.x[26] = uint64(uint32(s.x[26]) + uint32(s.x[19]))
	s.x[25] = uint64(uint32(s.x[25]) ^ uint32(s.x[26]))
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[20]))
	s.x[17] = uint64(uint32(s.x[17]) + uint32(s.x[21]))
	s.x[17] = uint64(uint32(s.x[25]) ^ uint32(s.x[17]))
	s.x[17] = uint64(uint32(s.x[17]) + uint32(s.x[22]))
	s.x[17] = uint64(uint32(s.x[8])*uint32(s.x[17]) + uint32(s.x[12]))
	s.x[25] = s.x[23] + s.x[13]
	if err := s.storeValue(s.x[25], uint64(byte(s.x[17])), 1); err != nil {
		return 0, fmt.Errorf("account boolean-writer C output: %w", err)
	}
	s.x[11]++
	s.compare64(s.x[13], 0)
	gate := uint32(boolValue(s.condition("ne")))
	s.x[17] = uint64(gate)
	s.x[25] = uint64(gate << 5)
	s.x[17] = uint64(uint32(s.x[25]) - gate<<1)
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[10]))
	next, err := loadBootstrapPointer(s, s.x[24], uint32(s.x[17]))
	if err != nil {
		return 0, fmt.Errorf("account boolean-writer C target: %w", err)
	}
	s.x[25] = next
	s.x[13]--
	s.x[8]++
	s.x[12] += 0xb6
	s.x[17] = 0x20ea
	return next, nil
}

func runAccountBooleanWriterD(s *state) (uint64, error) {
	s.traceBlock(accountBooleanWriterD)
	s.x[0] = 0
	s.x[10] = s.subtractWithFlags(s.x[10], 1, 64)
	s.x[17] = uint64(uint32(s.x[9]) ^ uint32(s.x[11]))
	s.x[27] = uint64(uint32(s.x[9]) ^ uint32(s.x[13]))
	s.x[27] = uint64(uint32(s.x[27]) + uint32(s.x[14]))
	s.x[28] = uint64(uint32(s.x[9]) ^ uint32(s.x[15]))
	s.x[28] = uint64(uint32(s.x[28]) + uint32(s.x[16]))
	s.x[27] = uint64(uint32(s.x[27]) ^ uint32(s.x[28]))
	s.x[28] = uint64(uint32(s.x[9]) ^ uint32(s.x[1]))
	s.x[17] = uint64(uint32(s.x[17]) + uint32(s.x[12]))
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[9]))
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[27]))
	s.x[27] = uint64(uint32(s.x[28]) + uint32(s.x[2]))
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[27]))
	s.x[27] = uint64(uint32(s.x[17]) ^ uint32(s.x[4]))
	s.x[27] = uint64(uint32(s.x[27]) + uint32(s.x[5]))
	s.x[28] = uint64(uint32(s.x[17]) ^ uint32(s.x[6]))
	s.x[28] = uint64(uint32(s.x[28]) + uint32(s.x[7]))
	s.x[27] = uint64(uint32(s.x[27]) ^ uint32(s.x[28]))
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[19]))
	s.x[17] = uint64(uint32(s.x[17]) + uint32(s.x[21]))
	s.x[17] = uint64(uint32(s.x[27]) ^ uint32(s.x[17]))
	s.x[17] = uint64(uint32(s.x[17]) + uint32(s.x[22]))
	s.x[27] = uint64(uint32(s.x[9]) + uint32(s.x[23]))
	s.x[17] = uint64(uint32(s.x[17]) * uint32(s.x[27]))
	s.x[17] = uint64(uint32(s.x[9])*uint32(s.x[3]) + uint32(s.x[17]))
	s.x[17] = uint64(uint32(s.x[17]) + uint32(s.x[24]))
	s.x[27] = s.x[25] + s.x[10]
	if err := s.storeValue(s.x[27], uint64(byte(s.x[17])), 1); err != nil {
		return 0, fmt.Errorf("account boolean-writer D output: %w", err)
	}
	s.x[9]++
	gate := uint32(boolValue(s.condition("eq")))
	s.x[17] = uint64(gate * uint32(s.x[26]))
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[8]))
	next, err := s.load(s.x[20]+uint64(uint32(s.x[17]))*8, 8)
	if err != nil {
		return 0, fmt.Errorf("account boolean-writer D target: %w", err)
	}
	s.x[27] = next
	s.x[17] = 0xc94a
	return next, nil
}

func runAccountSmallSelectorA(s *state) (uint64, error) {
	s.traceBlock(accountSmallSelectorA)
	s.x[11] = uint64(uint32(s.x[9]) ^ 0x59d6)
	s.x[11] = uint64(uint32(s.x[11]) * 0x527)
	s.x[22] = uint64(uint32(s.x[11]) ^ 0xd59c1f41)
	s.x[9] = uint64(uint32(s.x[9]) + 0x632)
	s.compare32(uint32(s.x[8]), uint32(s.x[3]))
	gate := uint32(boolValue(s.condition("eq")))
	s.x[12] = uint64(gate)
	s.x[13] = uint64(uint32(s.x[22]) + 0x2a647c71)
	s.x[12] = uint64(gate * uint32(s.x[13]))
	s.x[9] = uint64(uint32(s.x[12]) ^ uint32(s.x[9]))
	next, err := loadBootstrapPointer(s, s.x[30], uint32(s.x[9]))
	if err != nil {
		return 0, fmt.Errorf("account small-selector A target: %w", err)
	}
	s.x[9] = next
	s.x[25] = s.x[10]
	s.x[5] = s.x[8]
	s.x[17] = 0xfd9d
	return next, nil
}

func runAccountSmallSelectorB(s *state) (uint64, error) {
	s.traceBlock(accountSmallSelectorB)
	s.x[8] = uint64(uint32(s.x[10]) + 0xd47fecbb)
	s.x[11] = uint64(uint32(s.x[8]) & uint32(s.x[4]))
	s.compare32(uint32(s.x[25]), uint32(s.x[1]))
	gate := uint32(boolValue(s.condition("ne")))
	s.x[8] = uint64(gate * uint32(s.x[2]))
	s.x[8] = uint64(uint32(s.x[8]) ^ uint32(s.x[10]))
	next, err := loadBootstrapPointer(s, s.x[30], uint32(s.x[8]))
	if err != nil {
		return 0, fmt.Errorf("account small-selector B target: %w", err)
	}
	s.x[12] = next
	s.x[8] = 0xce757bc5
	s.x[17] = 0xfbad
	return next, nil
}

func runAccountBoundsSelector(s *state) (uint64, error) {
	s.traceBlock(accountBoundsSelector)
	s.x[9] = 0xffffd7ed
	s.x[20] = uint64(uint32(s.x[11]) + uint32(s.x[9]))
	word, err := s.load(s.sp+0x3f4, 4)
	if err != nil {
		return 0, fmt.Errorf("account bounds-selector stack word: %w", err)
	}
	s.x[21] = uint64(uint32(word))
	s.x[9] = uint64(uint32(s.x[23]) ^ 0xda2e7ab5)
	s.x[10] = 0x25d14c9e
	s.x[9] = uint64(uint32(s.x[9]) + uint32(s.x[10]))
	s.x[9] = uint64(uint32(s.x[20]) + uint32(s.x[9]))
	s.x[10] = uint64(uint32(s.x[23]) ^ 0x019c852d)
	s.x[12] = 0xfe637ad3
	s.x[10] = uint64(uint32(s.x[10]) + uint32(s.x[12]))
	s.x[12] = uint64(uint32(s.x[23]) ^ 0x15c7845d)
	s.x[13] = 0xea387ba3
	s.x[12] = uint64(uint32(s.x[12]) + uint32(s.x[13]))
	s.x[10] = uint64(uint32(s.x[10]) ^ uint32(s.x[12]))
	s.x[9] = uint64(uint32(s.x[9]) ^ uint32(s.x[10]))
	s.x[10] = uint64(uint32(s.x[11]) + 0x3563)
	s.x[12] = uint64(uint32(s.x[10]) + 0x2b7b2092)
	s.x[13] = uint64(uint32(s.x[21]) + 0x3651bf08)
	s.x[15] = 0x5d063af0
	s.compare32(uint32(s.x[13]), uint32(s.x[15]))
	first := uint32(boolValue(s.condition("lo")))
	s.x[14] = uint64(first)
	s.x[9] = uint64(uint32(s.x[12]) + uint32(s.x[9]))
	s.compare32(uint32(s.x[9]), uint32(s.x[15]))
	second := uint32(boolValue(s.condition("lo")))
	s.x[12] = uint64(second)
	s.x[14] = uint64(first ^ second)
	s.compare32(uint32(s.x[13]), uint32(s.x[9]))
	ordered := uint32(boolValue(s.condition("lo")))
	s.x[9] = uint64(ordered)
	s.compare32(uint32(s.x[14]), 0)
	selected := ordered
	if s.condition("ne") {
		selected = second
	}
	s.x[9] = uint64(selected * 0x17)
	s.x[9] = uint64(uint32(s.x[9]) | uint32(s.x[11]))
	next, err := loadBootstrapPointer(s, s.x[30], uint32(s.x[9]))
	if err != nil {
		return 0, fmt.Errorf("account bounds-selector target: %w", err)
	}
	s.x[9] = next
	s.x[17] = 0xaa48
	return next, nil
}
