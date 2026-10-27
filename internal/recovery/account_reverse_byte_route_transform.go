package recovery

import "fmt"

const (
	accountReverseByteRouteA = uint64(0x100177eac)
	accountReverseByteRouteB = uint64(0x1002182dc)
	accountByteCompareRoute  = uint64(0x10016e398)
	accountByteCompareTail   = uint64(0x10016e400)
)

func runAccountReverseByteRoute(s *state, pc, marker uint64) (uint64, error) {
	s.traceBlock(pc)
	s.x[8] = 0
	s.x[12] = s.subtractWithFlags(s.x[12], 1, 64)
	value, err := s.load(s.x[11]+s.x[12], 1)
	if err != nil {
		return 0, fmt.Errorf("account reverse-byte input: %w", err)
	}
	s.x[16] = uint64(byte(value))
	s.x[17] = uint64(uint32(s.x[13]) & uint32(s.x[16]<<1))
	s.x[16] = uint64(uint32(s.x[16]) - uint32(s.x[17]))
	s.x[16] = uint64(uint32(s.x[16]) + 0x1b)
	if err := s.storeValue(s.x[10]+s.x[12], uint64(byte(s.x[16])), 1); err != nil {
		return 0, fmt.Errorf("account reverse-byte output: %w", err)
	}
	s.x[16] = uint64(uint32(boolValue(s.condition("eq"))))
	s.x[16] = uint64(uint32(s.x[16]) * uint32(s.x[14]))
	s.x[16] = uint64(uint32(s.x[16]) ^ uint32(s.x[15]))
	value, err = s.load(s.x[9]+uint64(uint32(s.x[16]))*8, 8)
	if err != nil {
		return 0, fmt.Errorf("account reverse-byte target: %w", err)
	}
	s.x[0] = value
	s.x[17] = marker
	return value, nil
}

func runAccountByteCompareRoute(s *state) (uint64, error) {
	s.traceBlock(accountByteCompareRoute)
	s.compare32(uint32(s.x[12]), uint32(s.x[9]))
	s.x[14] = uint64(uint32(boolValue(s.condition("lo"))))
	value, err := s.load(s.x[26]+uint64(uint32(s.x[8])), 1)
	if err != nil {
		return 0, fmt.Errorf("account byte-compare input: %w", err)
	}
	s.x[15] = uint64(byte(value))
	s.x[16] = uint64(uint32(s.x[8]) + 1)
	s.x[17] = uint64(uint32(s.x[10]) + 0x982d65b1)
	s.x[10] = uint64(uint32(s.x[10]) + 1)
	s.x[0] = 0xffffe0fe
	s.x[0] = uint64(uint32(s.x[19]) + uint32(s.x[0]))
	s.x[0] = uint64(uint32(s.x[0]) & uint32(s.x[15]<<1))
	s.x[15] = uint64(uint32(s.x[15]) - uint32(s.x[0]))
	s.x[15] = uint64(uint32(s.x[15]) + 0x1b)
	s.x[0] = s.sp + 0x29c
	if err := s.storeValue(s.x[0]+uint64(uint32(s.x[17])), uint64(byte(s.x[15])), 1); err != nil {
		return 0, fmt.Errorf("account byte-compare output: %w", err)
	}
	s.x[15] = 0x0d2a90ec
	s.compare32(uint32(s.x[16]), uint32(s.x[11]))
	s.x[17] = uint64(uint32(boolValue(s.condition("hi"))))
	s.x[17] = uint64(uint32(s.x[14]) ^ uint32(s.x[17]))
	s.x[8] = uint64(uint32(s.x[8]) + uint32(s.x[15]))
	s.compare32(uint32(s.x[8]), uint32(s.x[12]))
	s.x[8] = uint64(uint32(boolValue(s.condition("lo"))))
	s.compare32(uint32(s.x[17]), 0)
	if s.condition("ne") {
		s.x[8] = uint64(uint32(s.x[14]))
	}
	s.x[8] = uint64(uint32(s.x[8]) * uint32(s.x[13]))

	s.traceBlock(accountByteCompareTail)
	s.x[8] = uint64(uint32(s.x[8]) ^ uint32(s.x[19]))
	value, err = s.load(s.x[24]+uint64(uint32(s.x[8]))*8, 8)
	if err != nil {
		return 0, fmt.Errorf("account byte-compare target: %w", err)
	}
	s.x[14] = value
	s.x[8] = s.x[16]
	s.x[17] = 0x5b3e
	return value, nil
}
