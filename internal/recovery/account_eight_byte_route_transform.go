package recovery

import "fmt"

const accountEightByteRoute = uint64(0x1000dbd70)

func runAccountEightByteRoute(s *state) (uint64, error) {
	s.traceBlock(accountEightByteRoute)
	s.x[14] = s.x[12] - 1
	s.x[15] = s.x[10] + s.x[12]
	value, err := s.load(s.x[15], 1)
	if err != nil {
		return 0, fmt.Errorf("account eight-byte first input: %w", err)
	}
	s.x[15] = uint64(byte(value))
	if err := s.storeValue(s.x[8]+s.x[14], uint64(byte(s.x[15])), 1); err != nil {
		return 0, fmt.Errorf("account eight-byte first output: %w", err)
	}
	s.x[14] = 0xffffd4b5
	s.x[15] = s.x[12] - 2
	value, err = s.load(s.x[13]+s.x[15], 1)
	if err != nil {
		return 0, fmt.Errorf("account eight-byte second input: %w", err)
	}
	s.x[16] = uint64(byte(value))
	if err := s.storeValue(s.x[8]+s.x[15], uint64(byte(s.x[16])), 1); err != nil {
		return 0, fmt.Errorf("account eight-byte second output: %w", err)
	}
	s.x[15] = 0x4448
	s.x[16] = s.x[12] - 3
	value, err = s.load(s.x[13]+s.x[16], 1)
	if err != nil {
		return 0, fmt.Errorf("account eight-byte third input: %w", err)
	}
	s.x[17] = uint64(byte(value))
	if err := s.storeValue(s.x[8]+s.x[16], uint64(byte(s.x[17])), 1); err != nil {
		return 0, fmt.Errorf("account eight-byte third output: %w", err)
	}
	s.x[14] = uint64(uint32(s.x[11]) + uint32(s.x[14]))
	s.x[16] = s.x[12] - 4
	value, err = s.load(s.x[13]+s.x[16], 1)
	if err != nil {
		return 0, fmt.Errorf("account eight-byte fourth input: %w", err)
	}
	s.x[17] = uint64(byte(value))
	if err := s.storeValue(s.x[8]+s.x[16], uint64(byte(s.x[17])), 1); err != nil {
		return 0, fmt.Errorf("account eight-byte fourth output: %w", err)
	}
	s.x[14] = uint64(uint32(s.x[14]) | uint32(s.x[15]))
	s.x[15] = s.x[12] - 5
	value, err = s.load(s.x[13]+s.x[15], 1)
	if err != nil {
		return 0, fmt.Errorf("account eight-byte fifth input: %w", err)
	}
	s.x[16] = uint64(byte(value))
	if err := s.storeValue(s.x[8]+s.x[15], uint64(byte(s.x[16])), 1); err != nil {
		return 0, fmt.Errorf("account eight-byte fifth output: %w", err)
	}
	s.x[15] = s.x[12] - 6
	value, err = s.load(s.x[13]+s.x[15], 1)
	if err != nil {
		return 0, fmt.Errorf("account eight-byte sixth input: %w", err)
	}
	s.x[16] = uint64(byte(value))
	s.x[17] = 0xffffffffffffbbb1
	if err := s.storeValue(s.x[8]+s.x[15], uint64(byte(s.x[16])), 1); err != nil {
		return 0, fmt.Errorf("account eight-byte sixth output: %w", err)
	}
	s.x[15] = s.x[14] ^ s.x[17]
	s.x[15] += s.x[12]
	value, err = s.load(s.x[13]+s.x[15], 1)
	if err != nil {
		return 0, fmt.Errorf("account eight-byte seventh input: %w", err)
	}
	s.x[16] = uint64(byte(value))
	if err := s.storeValue(s.x[8]+s.x[15], uint64(byte(s.x[16])), 1); err != nil {
		return 0, fmt.Errorf("account eight-byte seventh output: %w", err)
	}
	s.x[12] = s.subtractWithFlags(s.x[12], 8, 64)
	value, err = s.load(s.x[13]+s.x[12], 1)
	if err != nil {
		return 0, fmt.Errorf("account eight-byte eighth input: %w", err)
	}
	s.x[15] = uint64(byte(value))
	if err := s.storeValue(s.x[8]+s.x[12], uint64(byte(s.x[15])), 1); err != nil {
		return 0, fmt.Errorf("account eight-byte eighth output: %w", err)
	}
	s.x[15] = uint64(uint32(boolValue(s.condition("eq"))))
	s.x[16] = 0xdc
	s.x[15] = uint64(uint32(s.x[15]) * uint32(s.x[16]))
	s.x[15] = uint64(uint32(s.x[15]) ^ uint32(s.x[11]))
	next, err := loadBootstrapPointer(s, s.x[9], uint32(s.x[15]))
	if err != nil {
		return 0, fmt.Errorf("account eight-byte target: %w", err)
	}
	s.x[15] = next
	s.x[17] = 0x5012
	return next, nil
}
