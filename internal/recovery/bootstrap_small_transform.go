package recovery

import "fmt"

const (
	bootstrapByteCompareSelect = uint64(0x100933ef8)
	bootstrapDistanceSelect    = uint64(0x100a28c28)
	bootstrapSwapPass          = uint64(0x1005df264)
	bootstrapSwapExit          = uint64(0x1005df2b0)
	bootstrapReverseWordWrite  = uint64(0x10011a6b0)
	bootstrapSparseZeroWrite   = uint64(0x10095e29c)
	bootstrapDualZeroWrite     = uint64(0x100988a34)
	bootstrapEqualitySelect    = uint64(0x100c1ee1c)
	bootstrapParitySelect      = uint64(0x100c94234)
)

func runBootstrapByteCompareSelect(s *state) (uint64, error) {
	s.traceBlock(bootstrapByteCompareSelect)
	s.x[16] = uint64(uint32(s.x[13]) ^ uint32(s.x[14]))
	s.x[21] = uint64(uint32(s.x[16]) * uint32(s.x[0]))
	s.x[16] = uint64(uint32(s.x[16]) + uint32(s.x[15]))
	s.x[17] = s.x[16] + s.x[1]
	s.x[17] ^= s.x[21]
	s.x[17] += s.x[12]
	left, err := s.load(s.x[9]+s.x[17], 1)
	if err != nil {
		return 0, fmt.Errorf("bootstrap byte-compare left: %w", err)
	}
	s.x[22] = uint64(byte(left))
	right, err := s.load(s.x[10]+s.x[17], 1)
	if err != nil {
		return 0, fmt.Errorf("bootstrap byte-compare right: %w", err)
	}
	s.x[20] = uint64(byte(right))
	s.compare32(uint32(s.x[22]), uint32(s.x[20]))
	gate := uint32(boolValue(s.condition("ne")))
	s.x[17] = uint64(gate)
	s.x[23] = uint64(gate << 6)
	s.x[17] = uint64(uint32(s.x[23]) - gate<<1)
	s.x[13] = uint64(uint32(s.x[17]) ^ uint32(s.x[13]))
	next, err := loadBootstrapPointer(s, s.x[30], uint32(s.x[13]))
	if err != nil {
		return 0, fmt.Errorf("bootstrap byte-compare dispatch: %w", err)
	}
	s.x[13] = next
	s.x[17] = 0xbfca
	return next, nil
}

func runBootstrapDistanceSelect(s *state) (uint64, error) {
	s.traceBlock(bootstrapDistanceSelect)
	s.x[13] = uint64(uint32(s.x[21]) + uint32(s.x[2]))
	s.x[16] = uint64(uint32(s.x[13]) ^ uint32(s.x[3]))
	s.x[13] = uint64(uint32(s.x[13]) & uint32(s.x[4]))
	s.x[17] = s.x[13] ^ s.x[5]
	s.x[20] = s.x[12] - s.x[11]
	s.x[12] = s.x[17] + s.x[12]
	s.x[17] = s.x[16] + s.x[6]
	s.compare64(s.x[20], s.x[17])
	gate := uint32(boolValue(s.condition("eq")))
	s.x[17] = uint64(gate * uint32(s.x[7]))
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[21]))
	next, err := loadBootstrapPointer(s, s.x[30], uint32(s.x[17]))
	if err != nil {
		return 0, fmt.Errorf("bootstrap distance dispatch: %w", err)
	}
	s.x[20] = next
	s.x[17] = 0xcb61
	return next, nil
}

func runBootstrapSwap(s *state) (uint64, error) {
	for iteration := 0; iteration < 0x45; iteration++ {
		s.traceBlock(bootstrapSwapPass)
		s.x[12] = uint64((uint32(s.x[8]) + uint32(s.x[10])) ^ uint32(s.x[8]))
		s.x[12] = uint64(uint32(s.x[12]) * uint32(s.x[9]))
		divisor := uint32(s.x[10])
		quotient := uint32(0)
		if divisor != 0 {
			quotient = uint32(s.x[12]) / divisor
		}
		s.x[13] = uint64(quotient)
		s.x[12] = uint64(uint32(s.x[12]) - quotient*divisor)
		linear := s.x[11] + (s.x[10] << 3)
		indexed := s.x[11] + uint64(uint32(s.x[12]))<<3
		first, err := s.load(linear, 8)
		if err != nil {
			return 0, fmt.Errorf("bootstrap swap linear input: %w", err)
		}
		second, err := s.load(indexed, 8)
		if err != nil {
			return 0, fmt.Errorf("bootstrap swap indexed input: %w", err)
		}
		if err := s.storeValue(linear, second, 8); err != nil {
			return 0, fmt.Errorf("bootstrap swap linear output: %w", err)
		}
		if err := s.storeValue(indexed, first, 8); err != nil {
			return 0, fmt.Errorf("bootstrap swap indexed output: %w", err)
		}
		s.x[14], s.x[15], s.x[12] = first, second, second
		s.x[10]++
		s.compare64(s.x[10], 0x45)
		if s.condition("eq") {
			return bootstrapSwapExit, nil
		}
	}
	return 0, fmt.Errorf("bootstrap swap did not reach its natural exit")
}

func runBootstrapReverseWordWrite(s *state) (uint64, error) {
	s.traceBlock(bootstrapReverseWordWrite)
	s.x[15] = s.subtractWithFlags(s.x[15], 1, 64)
	base, err := s.load(s.x[19]+8, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap reverse-word base: %w", err)
	}
	s.x[10] = base
	if err := s.storeValue(base+(s.x[15]<<2), uint64(uint32(s.x[12])), 4); err != nil {
		return 0, fmt.Errorf("bootstrap reverse-word output: %w", err)
	}
	gate := uint32(boolValue(s.condition("hi")))
	s.x[16] = uint64(gate)
	s.x[10] = uint64(uint32(s.x[8]) + uint32(s.x[11]))
	s.x[16] = uint64(uint32(s.x[10]) ^ gate)
	s.x[16] = uint64(uint32(s.x[13]) &^ uint32(s.x[16]))
	s.x[16] = uint64(uint32(s.x[16]) * uint32(s.x[14]))
	s.x[16] = uint64(uint32(s.x[16]) ^ uint32(s.x[8]))
	next, err := loadBootstrapPointer(s, s.x[20], uint32(s.x[16]))
	if err != nil {
		return 0, fmt.Errorf("bootstrap reverse-word dispatch: %w", err)
	}
	s.x[0] = next
	s.x[17] = 0x9801
	return next, nil
}

func runBootstrapSparseZeroWrite(s *state) (uint64, error) {
	s.traceBlock(bootstrapSparseZeroWrite)
	s.x[17] = uint64(uint32(s.x[9]) - 1)
	s.x[17] &= 0xff
	s.x[3] = uint64(uint32(s.x[17]) * uint32(s.x[12]))
	s.x[4] = uint64(uint32(s.x[17]) * uint32(s.x[13]))
	s.x[4] = uint64(uint32(s.x[4]) >> 21)
	s.x[3] = uint64(uint32(s.x[3]) - uint32(s.x[4])*uint32(s.x[14]))
	s.x[3] = s.x[22] + uint64(uint16(s.x[3]))
	if err := s.storeValue(s.x[3], 0, 1); err != nil {
		return 0, fmt.Errorf("bootstrap sparse byte output: %w", err)
	}
	s.x[17] = uint64(uint32(s.x[17])*uint32(s.x[15]) + uint32(s.x[11]))
	s.x[17] = uint64(uint32(s.x[17]) + uint32(s.x[16]))
	s.x[3] = uint64(uint32(s.x[17]) & 0xfffe)
	s.x[3] = uint64(uint32(s.x[3]) * uint32(s.x[0]))
	s.x[3] = uint64(uint32(s.x[3]) >> 22)
	s.x[17] = uint64(uint32(s.x[17]) - uint32(s.x[3])*uint32(s.x[1]))
	s.x[17] = (s.x[17] & 0xffff) << 2
	if err := s.storeValue(s.x[8]+s.x[17], 0, 4); err != nil {
		return 0, fmt.Errorf("bootstrap sparse word output: %w", err)
	}
	s.test32(uint32(s.x[9]) & 0xc0)
	gate := uint32(boolValue(s.condition("ne")))
	s.x[17] = uint64(gate * uint32(s.x[2]))
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[10]))
	next, err := loadBootstrapPointer(s, s.x[30], uint32(s.x[17]))
	if err != nil {
		return 0, fmt.Errorf("bootstrap sparse dispatch: %w", err)
	}
	s.x[3] = next
	s.x[9] = uint64(uint32(s.x[9]) + 1)
	s.x[17] = 0x1e5e
	return next, nil
}

func runBootstrapDualZeroWrite(s *state) (uint64, error) {
	s.traceBlock(bootstrapDualZeroWrite)
	s.x[17] = uint64(uint32(s.x[14]) & 0xff)
	s.x[14] = uint64(uint32(s.x[14]) + 0x88)
	s.x[4] = uint64(uint32(s.x[17]) * uint32(s.x[12]))
	divisor := uint32(s.x[10])
	quotient := uint32(0)
	if divisor != 0 {
		quotient = uint32(s.x[4]) / divisor
	}
	s.x[5] = uint64(quotient)
	s.x[4] = uint64(uint32(s.x[4]) - quotient*divisor)
	if err := s.storeValue(s.x[8]+uint64(uint32(s.x[4])), 0, 1); err != nil {
		return 0, fmt.Errorf("bootstrap dual byte output: %w", err)
	}
	s.x[4] = uint64(uint32(s.x[11]) & uint32(s.x[14]<<1))
	s.x[14] = uint64(uint32(s.x[14]) ^ uint32(s.x[15]))
	s.x[14] = uint64(uint32(s.x[4]) + uint32(s.x[14]))
	s.test32(uint32(s.x[14]) & 0xc0)
	s.x[14] = uint64(uint32(s.x[17])*uint32(s.x[16]) + uint32(s.x[14]))
	s.x[4] = uint64(uint32(s.x[14]) & 0xffff)
	s.x[4] = uint64(uint32(s.x[4]) * uint32(s.x[1]))
	s.x[4] = uint64(uint32(s.x[4]) >> 23)
	s.x[14] = uint64(uint32(s.x[14]) - uint32(s.x[4])*uint32(s.x[2]))
	s.x[14] = (s.x[14] & 0xffff) << 2
	if err := s.storeValue(s.x[9]+s.x[14], 0, 4); err != nil {
		return 0, fmt.Errorf("bootstrap dual word output: %w", err)
	}
	gate := uint32(boolValue(s.condition("ne")))
	s.x[4] = uint64(gate)
	s.x[14] = uint64(uint32(s.x[17]) + 1)
	s.x[17] = uint64(gate * uint32(s.x[3]))
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[24]))
	next, err := loadBootstrapPointer(s, s.x[13], uint32(s.x[17]))
	if err != nil {
		return 0, fmt.Errorf("bootstrap dual dispatch: %w", err)
	}
	s.x[4] = next
	s.x[17] = 0x1597
	return next, nil
}

func runBootstrapEqualitySelect(s *state) (uint64, error) {
	s.traceBlock(bootstrapEqualitySelect)
	s.x[9] = uint64(uint32(s.x[3]) ^ uint32(s.x[30]))
	s.x[25] = uint64(uint32(s.x[9]) * uint32(s.x[8]))
	s.x[9] = uint64(uint32(s.x[3]) + uint32(s.x[11]))
	s.x[9] = uint64(uint32(s.x[9]) | uint32(s.x[15]))
	s.x[17] = uint64(uint32(s.x[2]) + uint32(s.x[16]))
	s.x[2] = uint64(uint32(s.x[17]) + uint32(s.x[9]))
	s.x[17] = uint64(uint32(s.x[25]) + uint32(s.x[0]))
	s.compare32(uint32(s.x[2]), uint32(s.x[17]))
	gate := uint32(boolValue(s.condition("eq")))
	s.x[17] = uint64(gate * uint32(s.x[1]))
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[3]))
	next, err := loadBootstrapPointer(s, s.x[28], uint32(s.x[17]))
	if err != nil {
		return 0, fmt.Errorf("bootstrap equality dispatch: %w", err)
	}
	s.x[3] = next
	s.x[26] = s.x[4]
	s.x[17] = 0x119a
	return next, nil
}

func runBootstrapParitySelect(s *state) (uint64, error) {
	s.traceBlock(bootstrapParitySelect)
	s.x[17] = uint64(uint32(s.x[9]) ^ uint32(s.x[12]))
	s.x[17] = uint64(uint32(s.x[17]) + uint32(s.x[17]<<3))
	s.x[3] = uint64(uint32(s.x[17] << 1))
	s.x[25] = uint64(uint32(s.x[26]) >> 1)
	s.x[4] = uint64(uint32(s.x[25]) ^ uint32(s.x[13]))
	s.test32(uint32(s.x[26]) & 2)
	gate := uint32(boolValue(s.condition("eq")))
	s.x[17] = uint64(gate)
	s.x[26] = uint64(uint32(s.x[3]) ^ uint32(s.x[14]))
	s.x[17] = uint64(gate * uint32(s.x[26]))
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[9]))
	next, err := loadBootstrapPointer(s, s.x[28], uint32(s.x[17]))
	if err != nil {
		return 0, fmt.Errorf("bootstrap parity dispatch: %w", err)
	}
	s.x[26] = next
	s.x[17] = 0xc750
	return next, nil
}
