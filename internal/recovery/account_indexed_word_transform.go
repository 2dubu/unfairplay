package recovery

import "fmt"

const (
	accountIndexedWordStore = uint64(0x10011a6b0)
	accountWordCopyRoute    = uint64(0x10021aea4)
	accountRangeRoute       = uint64(0x1004325e4)
)

func runAccountIndexedWordStore(s *state) (uint64, error) {
	s.traceBlock(accountIndexedWordStore)
	s.x[15] = s.subtractWithFlags(s.x[15], 1, 64)
	value, err := s.load(s.x[19]+8, 8)
	if err != nil {
		return 0, fmt.Errorf("account indexed-word store table: %w", err)
	}
	s.x[10] = value
	if err := s.storeValue(s.x[10]+(s.x[15]<<2), uint64(uint32(s.x[12])), 4); err != nil {
		return 0, fmt.Errorf("account indexed-word store output: %w", err)
	}
	s.x[16] = uint64(uint32(boolValue(s.condition("hi"))))
	s.x[10] = uint64(uint32(s.x[8]) + uint32(s.x[11]))
	s.x[16] = uint64(uint32(s.x[10]) ^ uint32(s.x[16]))
	s.x[16] = uint64(uint32(s.x[13]) &^ uint32(s.x[16]))
	s.x[16] = uint64(uint32(s.x[16]) * uint32(s.x[14]))
	s.x[16] = uint64(uint32(s.x[16]) ^ uint32(s.x[8]))
	next, err := loadBootstrapPointer(s, s.x[20], uint32(s.x[16]))
	if err != nil {
		return 0, fmt.Errorf("account indexed-word store target: %w", err)
	}
	s.x[0] = next
	s.x[17] = 0x9801
	return next, nil
}

func runAccountWordCopyRoute(s *state) (uint64, error) {
	s.traceBlock(accountWordCopyRoute)
	s.compare32(uint32(s.x[12]), uint32(s.x[15]))
	s.x[17] = uint64(uint32(boolValue(s.condition("lt"))))
	value, err := s.load(s.x[23]+8, 8)
	if err != nil {
		return 0, fmt.Errorf("account word-copy first table: %w", err)
	}
	s.x[10] = value
	s.x[19] = uint64(uint32(s.x[7]) + uint32(s.x[0]))
	value, err = s.load(s.x[10]+int32Offset(s.x[19], 2), 4)
	if err != nil {
		return 0, fmt.Errorf("account word-copy first input: %w", err)
	}
	s.x[10] = uint64(uint32(value)) ^ s.x[1]
	value, err = s.load(s.x[21]+8, 8)
	if err != nil {
		return 0, fmt.Errorf("account word-copy second table: %w", err)
	}
	s.x[14] = value
	value, err = s.load(s.x[14]+int32Offset(s.x[19], 2), 4)
	if err != nil {
		return 0, fmt.Errorf("account word-copy second input: %w", err)
	}
	s.x[14] = uint64(uint32(value)) ^ s.x[1]
	s.x[10] = s.x[11] + s.x[10]
	s.x[10] -= s.x[14]
	s.x[14] = uint64(uint32(s.x[13]) + uint32(s.x[16]))
	s.x[11] = uint64(uint32(s.x[14]) ^ uint32(s.x[4]))
	s.x[10] += s.x[2]
	s.x[11] = uint64(uint32(s.x[10]) + uint32(s.x[11]))
	s.x[25] = uint64(uint32(s.x[3]) & uint32(s.x[10]<<1))
	s.x[11] = uint64(uint32(s.x[11]) - uint32(s.x[25]))
	value, err = s.load(s.x[9]+8, 8)
	if err != nil {
		return 0, fmt.Errorf("account word-copy output table: %w", err)
	}
	s.x[25] = value
	if err := s.storeValue(s.x[25]+int32Offset(s.x[19], 2), uint64(uint32(s.x[11])), 4); err != nil {
		return 0, fmt.Errorf("account word-copy output: %w", err)
	}
	s.x[11] = s.x[10] >> 32
	s.x[10] = uint64(uint32(s.x[7]) + 1)
	s.x[7] = uint64(uint32(s.x[7]) + uint32(s.x[5]))
	s.compare32(uint32(s.x[7]), uint32(s.x[15]))
	s.x[19] = uint64(uint32(boolValue(s.condition("lt"))))
	s.x[19] = uint64(uint32(s.x[17]) ^ uint32(s.x[19]))
	s.compare32(uint32(s.x[7]), uint32(s.x[12]))
	s.x[7] = uint64(uint32(boolValue(s.condition("lt"))))
	s.compare32(uint32(s.x[19]), 0)
	if s.condition("ne") {
		s.x[17] = uint64(uint32(s.x[17]))
	} else {
		s.x[17] = uint64(uint32(s.x[7]))
	}
	s.x[17] = uint64(uint32(s.x[17]) ^ 1)
	s.x[17] = uint64(uint32(s.x[17]) * uint32(s.x[6]))
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[13]))
	value, err = s.load(s.x[20]+uint64(uint32(s.x[17]))*8, 8)
	if err != nil {
		return 0, fmt.Errorf("account word-copy target: %w", err)
	}
	s.x[19] = value
	s.x[7] = s.x[10]
	s.x[17] = 0x4c2
	return value, nil
}

func runAccountRangeRoute(s *state) (uint64, error) {
	s.traceBlock(accountRangeRoute)
	s.x[10] = uint64(uint32(s.x[23]) + uint32(s.x[3]))
	s.x[13] = uint64(uint32(s.x[10]) & uint32(s.x[4]))
	s.x[17] = uint64(uint32(s.x[19]) + 1)
	s.x[10] = uint64(uint32(s.x[11]) + uint32(s.x[19]))
	s.x[19] = s.x[10] + s.x[6]
	s.x[22] = s.x[12] ^ s.x[19]
	s.x[19] -= s.x[12]
	s.x[19] &= ^s.x[22]
	s.x[22] = uint64(uint32(s.x[13]) + uint32(s.x[5]))
	s.x[25] = s.x[22] + s.x[7]
	s.x[25] &= s.x[12]
	s.x[19] = s.x[25] | s.x[19]
	s.x[19] >>= 63
	s.x[19] = uint64(uint32(s.x[19]) * uint32(s.x[21]))
	s.x[19] = uint64(uint32(s.x[19]) ^ uint32(s.x[23]))
	value, err := s.load(s.sp+0xdf0, 8)
	if err != nil {
		return 0, fmt.Errorf("account range-route table: %w", err)
	}
	s.x[23] = value
	next, err := loadBootstrapPointer(s, s.x[23], uint32(s.x[19]))
	if err != nil {
		return 0, fmt.Errorf("account range-route target: %w", err)
	}
	s.x[23] = next
	s.x[19] = s.x[17]
	s.x[17] = 0xf8c0
	return next, nil
}
