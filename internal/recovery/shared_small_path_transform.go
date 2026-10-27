package recovery

import "fmt"

const (
	sharedIndexMixStart = uint64(0x1000c2f90)
	sharedIndexMixTail  = uint64(0x1000c3000)
	sharedRangeSelect   = uint64(0x100218310)
)

func runSharedIndexMix(s *state) (uint64, error) {
	s.traceBlock(sharedIndexMixStart)
	s.x[8] = 0xffffc66a
	s.x[17] = uint64(uint32(s.x[19]) + uint32(s.x[8]))
	s.x[0] = uint64(uint32(s.x[17]) + 0x104)
	s.x[8] = uint64(uint32(s.x[17]) - 0x360)
	s.x[1] = 0x214
	s.x[8] = uint64(uint32(s.x[8]) | uint32(s.x[1]))
	s.x[1] = 0x3892
	s.x[19] = uint64(uint32(s.x[0]) | uint32(s.x[1]))
	base, err := s.load(s.x[22]+0x20, 8)
	if err != nil {
		return 0, fmt.Errorf("shared index-mix table: %w", err)
	}
	s.x[1] = base
	s.x[0] = s.x[6] & 0xff
	word, err := s.load(base+(s.x[0]<<2), 4)
	if err != nil {
		return 0, fmt.Errorf("shared index-mix word: %w", err)
	}
	s.x[1] = uint64(uint32(word))
	s.x[2] = uint64(uint32(s.x[4]) ^ 0x3a95dabc)
	s.x[5] = 0xc569e84d
	s.x[2] = uint64(uint32(s.x[2]) + uint32(s.x[5]))
	s.x[2] = uint64(uint32(s.x[2]) + uint32(s.x[19]))
	s.x[5] = uint64(uint32(s.x[4]) ^ 0x7ebe4825)
	s.x[7] = uint64(uint32(s.x[8]) ^ 0x8141b5ce)
	s.x[23] = uint64(uint32(s.x[4]) ^ 0x8a5ee95c)
	s.x[25] = 0x16a4

	s.traceBlock(sharedIndexMixTail)
	s.x[25] = 0x75a116a4
	s.x[23] = uint64(uint32(s.x[23]) + uint32(s.x[25]))
	s.x[2] = uint64(uint32(s.x[2]) ^ uint32(s.x[23]))
	s.x[5] = uint64(uint32(s.x[7]) + uint32(s.x[5]))
	s.x[2] = uint64(uint32(s.x[2]) ^ uint32(s.x[5]))
	s.x[5] = 0x002230e3
	s.x[1] = uint64(uint32(s.x[2]) - uint32(s.x[1]))
	s.compare32(uint32(s.x[1]), uint32(s.x[5]))
	gate := uint32(boolValue(s.condition("eq")))
	s.x[1] = uint64(gate)
	s.x[6] = uint64(uint32(s.x[6]) + 1)
	s.x[2] = 0x2cfd
	s.x[1] = uint64(gate * uint32(s.x[2]))
	s.x[17] = uint64(uint32(s.x[1]) ^ uint32(s.x[17]))
	next, err := loadBootstrapPointer(s, s.x[21], uint32(s.x[17]))
	if err != nil {
		return 0, fmt.Errorf("shared index-mix target: %w", err)
	}
	s.x[1] = next
	s.x[17] = 0xca77
	return next, nil
}

func runSharedRangeSelect(s *state) (uint64, error) {
	s.traceBlock(sharedRangeSelect)
	s.x[8] = 0xffffc51e
	s.x[17] = 0xffffc635
	s.x[17] = uint64(uint32(s.x[19]) + uint32(s.x[17]))
	s.x[0] = uint64(uint32(s.x[6]) & 0xff)
	s.compare32(uint32(s.x[0]), uint32(byte(s.x[17])))
	gate := uint32(boolValue(s.condition("lo")))
	s.x[17] = uint64(gate)
	s.x[0] = 0x7471
	s.x[17] = uint64(gate * uint32(s.x[0]))
	s.x[8] = uint64(uint32(s.x[19]) + uint32(s.x[8]))
	s.x[8] = uint64(uint32(s.x[17]) ^ uint32(s.x[8]))
	next, err := loadBootstrapPointer(s, s.x[21], uint32(s.x[8]))
	if err != nil {
		return 0, fmt.Errorf("shared range-select target: %w", err)
	}
	s.x[8] = next
	s.x[5] = 0xce757bc5
	s.x[17] = 0x9354
	return next, nil
}
