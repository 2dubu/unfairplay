package recovery

import "fmt"

const (
	accountByteRouteWriter = uint64(0x100186f4c)
	accountBytePairRoute   = uint64(0x10025f088)
)

func runAccountByteRouteWriter(s *state) (uint64, error) {
	s.traceBlock(accountByteRouteWriter)
	s.compare32(uint32(s.x[12]), uint32(s.x[9]))
	value, err := s.load(s.x[27]+uint64(uint32(s.x[11])), 1)
	if err != nil {
		return 0, fmt.Errorf("account byte-route input: %w", err)
	}
	s.x[13] = uint64(byte(value))
	s.x[14] = uint64(uint32(boolValue(s.condition("lo"))))
	s.x[15] = uint64(uint32(s.x[11]) + 1)
	s.x[16] = 0x835dd153
	s.x[16] = uint64(uint32(s.x[8]) + uint32(s.x[16]))
	s.x[8] = uint64(uint32(s.x[8]) + 1)
	s.x[17] = uint64(uint32(0x36) & uint32(s.x[13]<<1))
	s.x[0] = uint64(uint32(s.x[19]) ^ 0x9a)
	s.x[13] = uint64(uint32(s.x[13]) + uint32(s.x[0]))
	s.x[13] = uint64(uint32(s.x[13]) - uint32(s.x[17]))
	s.x[17] = s.sp + 0x9c
	if err := s.storeValue(s.x[17]+uint64(uint32(s.x[16])), uint64(byte(s.x[13])), 1); err != nil {
		return 0, fmt.Errorf("account byte-route output: %w", err)
	}
	s.x[13] = 0x649d920c
	s.compare32(uint32(s.x[15]), uint32(s.x[10]))
	s.x[16] = uint64(uint32(boolValue(s.condition("hi"))))
	s.x[16] = uint64(uint32(s.x[14]) ^ uint32(s.x[16]))
	s.x[11] = uint64(uint32(s.x[11]) + uint32(s.x[13]))
	s.compare32(uint32(s.x[11]), uint32(s.x[12]))
	s.x[11] = uint64(uint32(boolValue(s.condition("lo"))))
	s.compare32(uint32(s.x[16]), 0)
	if s.condition("ne") {
		s.x[11] = uint64(uint32(s.x[14]))
	}
	s.x[13] = uint64(uint32(s.x[11]) << 1)
	s.x[11] = uint64(uint32(s.x[13]) | uint32(s.x[11]<<4))
	s.x[11] = uint64(uint32(s.x[11]) ^ uint32(s.x[21]))
	next, err := loadBootstrapPointer(s, s.x[20], uint32(s.x[11]))
	if err != nil {
		return 0, fmt.Errorf("account byte-route target: %w", err)
	}
	s.x[13] = next
	s.x[11] = s.x[15]
	s.x[17] = 0xfc8d
	return next, nil
}

func runAccountBytePairRoute(s *state) (uint64, error) {
	s.traceBlock(accountBytePairRoute)
	s.x[13] = uint64(uint32(s.x[0]) - 3)
	s.x[16] = 0xffffffffffff616d
	s.x[16] = s.x[14] + s.x[16]
	s.x[16] += s.x[15]
	value, err := s.load(s.x[12]+s.x[16], 1)
	if err != nil {
		return 0, fmt.Errorf("account byte-pair first input: %w", err)
	}
	s.x[17] = uint64(byte(value))
	value, err = s.load(s.x[9]+s.x[16], 1)
	if err != nil {
		return 0, fmt.Errorf("account byte-pair second input: %w", err)
	}
	s.x[2] = uint64(byte(value))
	s.x[2] = uint64(uint32(s.x[2]) ^ uint32(s.x[24]))
	s.x[17] = uint64(uint32(s.x[17]) - uint32(s.x[2]))
	s.x[2] = uint64(uint32(s.x[17]) ^ uint32(s.x[25]))
	s.x[17] = uint64(uint32(s.x[23]) & uint32(s.x[17]<<1))
	s.x[17] = uint64(uint32(s.x[17]) + uint32(s.x[2]))
	s.compare32(uint32(s.x[17]), uint32(s.x[25]))
	s.x[17] = uint64(uint32(boolValue(s.condition("eq"))))
	s.x[2] = uint64(uint32(s.x[0]) + uint32(s.x[26]))
	s.x[17] = uint64(uint32(s.x[17]) * uint32(s.x[2]))
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[0]))
	next, err := loadBootstrapPointer(s, s.x[27], uint32(s.x[17]))
	if err != nil {
		return 0, fmt.Errorf("account byte-pair target: %w", err)
	}
	s.x[0] = next
	s.x[17] = 0x7db5
	return next, nil
}
