package recovery

import "fmt"

const (
	accountIndexedByteRoute = uint64(0x10007c728)
	accountIndexedWordRoute = uint64(0x100176524)
)

func runAccountIndexedByteRoute(s *state) (uint64, error) {
	s.traceBlock(accountIndexedByteRoute)
	s.x[8] = uint64(uint32(s.x[12]) & uint32(s.x[13]))
	s.x[17] = s.x[9] + s.x[1]
	s.x[2] = s.x[9] + s.x[14]
	s.x[2] += s.x[8]
	s.x[9] += s.x[15]
	s.x[9] += s.x[16]
	if err := s.storeValue(s.x[9], uint64(byte(s.x[0])), 1); err != nil {
		return 0, fmt.Errorf("account indexed-byte output: %w", err)
	}
	s.compare64(s.x[17], 0x38)
	s.x[9] = uint64(uint32(boolValue(s.condition("lo"))))
	s.x[9] = uint64(uint32(s.x[9]) * uint32(s.x[10]))
	s.x[9] = uint64(uint32(s.x[9]) ^ uint32(s.x[11]))
	next, err := loadBootstrapPointer(s, s.x[21], uint32(s.x[9]))
	if err != nil {
		return 0, fmt.Errorf("account indexed-byte target: %w", err)
	}
	s.x[3] = next
	s.x[9] = s.x[2]
	s.x[17] = 0x5f4f
	return next, nil
}

func runAccountIndexedWordRoute(s *state) (uint64, error) {
	s.traceBlock(accountIndexedWordRoute)
	s.x[12] = uint64(uint32(s.x[11]) ^ 0xeb62)
	s.x[14] = 0xffff8181
	s.x[14] = uint64(uint32(s.x[13]) + uint32(s.x[14]))
	s.x[12] = uint64(uint32(s.x[12]) + uint32(s.x[14]))
	value, err := s.load(s.x[8]+8, 8)
	if err != nil {
		return 0, fmt.Errorf("account indexed-word first table: %w", err)
	}
	s.x[14] = value
	value, err = s.load(s.x[14]+int32Offset(s.x[12], 2), 4)
	if err != nil {
		return 0, fmt.Errorf("account indexed-word input: %w", err)
	}
	s.x[14] = uint64(uint32(value))
	value, err = s.load(s.x[9]+8, 8)
	if err != nil {
		return 0, fmt.Errorf("account indexed-word second table: %w", err)
	}
	s.x[15] = value
	if err := s.storeValue(s.x[15]+int32Offset(s.x[12], 2), uint64(uint32(s.x[14])), 4); err != nil {
		return 0, fmt.Errorf("account indexed-word first output: %w", err)
	}
	s.x[13] = s.subtractWithFlags(uint64(uint32(s.x[13])), 2, 32)
	value, err = s.load(s.x[8]+8, 8)
	if err != nil {
		return 0, fmt.Errorf("account indexed-word first table repeat: %w", err)
	}
	s.x[12] = value
	value, err = s.load(s.x[12]+int32Offset(s.x[13], 2), 4)
	if err != nil {
		return 0, fmt.Errorf("account indexed-word second input: %w", err)
	}
	s.x[12] = uint64(uint32(value))
	value, err = s.load(s.x[9]+8, 8)
	if err != nil {
		return 0, fmt.Errorf("account indexed-word second table repeat: %w", err)
	}
	s.x[14] = value
	if err := s.storeValue(s.x[14]+int32Offset(s.x[13], 2), uint64(uint32(s.x[12])), 4); err != nil {
		return 0, fmt.Errorf("account indexed-word second output: %w", err)
	}
	s.x[12] = uint64(uint32(boolValue(s.condition("ne"))))
	s.x[14] = uint64(uint32(s.x[12]) << 2)
	s.x[12] = uint64(uint32(s.x[14]) | uint32(s.x[12]<<5))
	s.x[12] = uint64(uint32(s.x[12]) ^ uint32(s.x[11]))
	next, err := loadBootstrapPointer(s, s.x[10], uint32(s.x[12]))
	if err != nil {
		return 0, fmt.Errorf("account indexed-word target: %w", err)
	}
	s.x[12] = next
	s.x[17] = 0xa668
	return next, nil
}
