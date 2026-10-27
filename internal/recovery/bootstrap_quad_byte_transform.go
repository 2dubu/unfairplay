package recovery

import "fmt"

const bootstrapQuadByte = uint64(0x10001fd04)

func runBootstrapQuadByte(s *state) (uint64, error) {
	s.traceBlock(bootstrapQuadByte)
	s.x[14] = 0xdf091cbf
	s.compare64(s.x[13], s.x[14])
	s.x[9] = uint64(uint32(s.x[8]) ^ 0x9258)
	s.x[15] = 0x181
	s.x[9] = uint64(uint32(s.x[9]) * uint32(s.x[15]))
	s.x[15] = uint64(uint32(s.x[11]))
	s.x[16] = s.x[15] | 2
	value, err := s.load(s.x[23]+s.x[16], 1)
	if err != nil {
		return 0, fmt.Errorf("bootstrap quad-byte first input: %w", err)
	}
	s.x[16] = uint64(byte(value)) ^ 0xd1
	if err := s.storeValue(s.x[10], uint64(byte(s.x[16])), 1); err != nil {
		return 0, fmt.Errorf("bootstrap quad-byte first output: %w", err)
	}
	value, err = s.load(s.x[23]+s.x[12], 1)
	if err != nil {
		return 0, fmt.Errorf("bootstrap quad-byte second input: %w", err)
	}
	s.x[12] = uint64(byte(value)) ^ 0x93
	s.x[12] = uint64(uint32(s.x[9]) ^ uint32(s.x[12]))
	if err := s.storeValue(s.x[10]+1, uint64(byte(s.x[12])), 1); err != nil {
		return 0, fmt.Errorf("bootstrap quad-byte second output: %w", err)
	}
	s.x[12] = s.x[15] | 3
	value, err = s.load(s.x[23]+s.x[12], 1)
	if err != nil {
		return 0, fmt.Errorf("bootstrap quad-byte third input: %w", err)
	}
	s.x[12] = uint64(byte(value)) ^ 0x1a
	if err := s.storeValue(s.x[10]+2, uint64(byte(s.x[12])), 1); err != nil {
		return 0, fmt.Errorf("bootstrap quad-byte third output: %w", err)
	}
	s.x[12] = s.x[15] | 1
	value, err = s.load(s.x[23]+s.x[12], 1)
	if err != nil {
		return 0, fmt.Errorf("bootstrap quad-byte fourth input: %w", err)
	}
	s.x[12] = uint64(byte(value)) ^ 0xa1
	if err := s.storeValue(s.x[10]+3, uint64(byte(s.x[12])), 1); err != nil {
		return 0, fmt.Errorf("bootstrap quad-byte fourth output: %w", err)
	}
	s.x[10] += 4
	s.x[12] = uint64(uint32(s.x[11]) + 4)
	s.x[11] = s.x[12] + s.x[14]
	if s.condition("hs") {
		s.compare64(s.x[13], s.x[11])
	} else {
		s.setNZCV(2)
	}
	gate := uint32(boolValue(s.condition("hi")))
	s.x[11] = uint64(gate)
	s.x[14] = uint64(gate << 8)
	s.x[11] = uint64(uint32(s.x[14]) - gate)
	s.x[11] = uint64(uint32(s.x[11]) ^ uint32(s.x[8]))
	next, err := loadBootstrapPointer(s, s.x[21], uint32(s.x[11]))
	if err != nil {
		return 0, fmt.Errorf("bootstrap quad-byte target: %w", err)
	}
	s.x[14] = next
	s.x[11] = s.x[12]
	s.x[17] = 0xec7a
	return next, nil
}
