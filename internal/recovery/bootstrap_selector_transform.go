package recovery

import "fmt"

const (
	bootstrapCounterSelector = uint64(0x1009b2174)
	bootstrapByteSelector    = uint64(0x1009f1804)
)

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
