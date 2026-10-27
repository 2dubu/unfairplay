package recovery

import "fmt"

const (
	accountByteShiftRoute     = uint64(0x1000e73e0)
	accountByteShiftRouteTail = uint64(0x1000e7400)
	accountByteWriteRoute     = uint64(0x1001f9898)
)

func runAccountByteShiftRoute(s *state) (uint64, error) {
	s.traceBlock(accountByteShiftRoute)
	s.x[11] = s.subtractWithFlags(s.x[11], 1, 64)
	value, err := s.load(s.x[10]+s.x[11], 1)
	if err != nil {
		return 0, fmt.Errorf("account byte-shift input: %w", err)
	}
	s.x[14] = uint64(byte(value))
	s.x[15] = uint64(uint32(shiftLeft(s.x[14], int(s.x[12]), 32)))
	s.x[15] = uint64(uint32(s.x[15]) & uint32(s.x[13]))
	s.x[14] = uint64(uint32(s.x[14]) - uint32(s.x[15]))
	s.x[14] = uint64(uint32(s.x[14]) + 0x1b)
	if err := s.storeValue(s.x[9]+s.x[11], uint64(byte(s.x[14])), 1); err != nil {
		return 0, fmt.Errorf("account byte-shift output: %w", err)
	}
	s.x[14] = uint64(uint32(boolValue(s.condition("ne"))))

	s.traceBlock(accountByteShiftRouteTail)
	s.x[15] = uint64(uint32(s.x[14]) << 3)
	s.x[14] = uint64(uint32(s.x[15]) | uint32(s.x[14]<<4))
	s.x[14] = uint64(uint32(s.x[14]) ^ uint32(s.x[8]))
	next, err := loadBootstrapPointer(s, s.x[27], uint32(s.x[14]))
	if err != nil {
		return 0, fmt.Errorf("account byte-shift target: %w", err)
	}
	s.x[14] = next
	s.x[15] = 0xce757bc5
	s.x[17] = 0x30d6
	return next, nil
}

func runAccountByteWriteRoute(s *state) (uint64, error) {
	s.traceBlock(accountByteWriteRoute)
	s.x[9] = uint64(uint32(s.x[21]) + uint32(s.x[8]))
	s.x[17] = s.x[13] + s.x[16]
	s.x[1] = s.x[13] + s.x[11]
	s.x[1] += s.x[9]
	s.x[13] += s.x[12]
	s.x[13] += s.x[14]
	if err := s.storeValue(s.x[13], uint64(byte(s.x[15])), 1); err != nil {
		return 0, fmt.Errorf("account byte-write output: %w", err)
	}
	s.compare64(s.x[17], 0x38)
	s.x[13] = uint64(uint32(boolValue(s.condition("lo"))))
	s.x[13] = uint64(uint32(s.x[13]) * uint32(s.x[0]))
	s.x[13] = uint64(uint32(s.x[13]) ^ uint32(s.x[10]))
	next, err := loadBootstrapPointer(s, s.x[20], uint32(s.x[13]))
	if err != nil {
		return 0, fmt.Errorf("account byte-write target: %w", err)
	}
	s.x[2] = next
	s.x[13] = s.x[1]
	s.x[17] = 0xbfc0
	return next, nil
}
