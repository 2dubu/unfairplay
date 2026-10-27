package recovery

import "fmt"

const (
	accountSignedPairDispatch     = uint64(0x1008e5fd0)
	accountSignedPairDispatchTail = uint64(0x1008e6000)
	accountComputedByteDispatch   = uint64(0x1009ecac0)
)

func runAccountSignedPairDispatch(s *state) (uint64, error) {
	s.traceBlock(accountSignedPairDispatch)
	s.x[12] = uint64(uint32(s.x[9]) ^ uint32(s.x[26]))
	value, err := s.load(s.x[11], 1)
	if err != nil {
		return 0, fmt.Errorf("account signed-pair first input: %w", err)
	}
	s.x[13] = uint64(uint32(int32(int8(value))))
	value, err = s.load(s.x[10], 1)
	if err != nil {
		return 0, fmt.Errorf("account signed-pair second input: %w", err)
	}
	s.x[14] = uint64(uint32(int32(int8(value))))
	s.x[8] = uint64(uint32(s.x[14]) ^ uint32(s.x[13]))
	_ = s.addWithCarry(uint64(uint32(s.x[8])), 0x6b, false, 32)
	s.x[8] = uint64(uint32(boolValue(s.condition("eq"))))
	s.x[15] = 0xffffdc71
	s.x[15] = uint64(uint32(s.x[12]) + uint32(s.x[15]))
	s.x[8] = uint64(uint32(s.x[8]) * uint32(s.x[15]))
	s.x[9] = uint64(uint32(s.x[9]) + 0x12ed)
	s.x[8] = uint64(uint32(s.x[8]) ^ uint32(s.x[9]))
	s.traceBlock(accountSignedPairDispatchTail)
	next, err := loadBootstrapPointer(s, s.x[30], uint32(s.x[8]))
	if err != nil {
		return 0, fmt.Errorf("account signed-pair target: %w", err)
	}
	s.x[8] = next
	s.x[17] = 0xd666
	return next, nil
}

func runAccountComputedByteDispatch(s *state) (uint64, error) {
	s.traceBlock(accountComputedByteDispatch)
	s.x[12] = uint64(uint32(s.x[10]) ^ 0x4d73)
	s.x[13] = uint64(uint32(s.x[15]) ^ 0xf5de6358)
	s.x[14] = 0x0a219ca8
	s.x[13] = uint64(uint32(s.x[13]) + uint32(s.x[14]))
	s.x[14] = uint64(uint32(s.x[15]) ^ 0x97c5a20e)
	s.x[16] = 0x683a5df2
	s.x[14] = uint64(uint32(s.x[14]) + uint32(s.x[16]))
	s.x[13] = uint64(uint32(s.x[13]) ^ uint32(s.x[14]))
	s.x[14] = uint64(uint32(s.x[15]) ^ 0xac6eba93)
	s.x[16] = 0x5391288f
	s.x[14] = uint64(uint32(s.x[12]) + uint32(s.x[14]))
	s.x[14] = uint64(uint32(s.x[14]) + uint32(s.x[16]))
	s.x[13] = uint64(uint32(s.x[13]) ^ uint32(s.x[14]))
	s.x[13] = uint64(uint32(s.x[11]) + uint32(s.x[13]))
	value, err := s.load(s.x[1]+uint64(uint32(s.x[13])), 1)
	if err != nil {
		return 0, fmt.Errorf("account computed-byte input: %w", err)
	}
	s.x[13] = uint64(byte(value))
	s.compare32(uint32(s.x[13]), 0)
	s.x[13] = uint64(uint32(boolValue(s.condition("ne"))))
	s.x[14] = uint64(uint32(s.x[13]) << 5)
	s.x[13] = uint64(uint32(s.x[14]) - uint32(s.x[13]<<2))
	s.x[13] = uint64(uint32(s.x[13]) ^ uint32(s.x[10]))
	next, err := loadBootstrapPointer(s, s.x[22], uint32(s.x[13]))
	if err != nil {
		return 0, fmt.Errorf("account computed-byte target: %w", err)
	}
	s.x[13] = next
	s.x[17] = 0xc4a7
	return next, nil
}
