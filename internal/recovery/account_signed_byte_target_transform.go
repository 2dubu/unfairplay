package recovery

import "fmt"

const (
	accountSignedByteTarget     = uint64(0x1002917f4)
	accountSignedByteTargetTail = uint64(0x100291800)
	accountPointerByteTarget    = uint64(0x1002e50a0)
)

func runAccountSignedByteTarget(s *state) (uint64, error) {
	s.traceBlock(accountSignedByteTarget)
	s.x[12] = 0x80ff440d
	s.x[15] = uint64(uint32(s.x[14]) ^ uint32(s.x[12]))
	s.traceBlock(accountSignedByteTargetTail)
	s.x[13] = uint64(uint32(s.x[15]) + 0x4d6)
	s.x[16] = 0x7f00bec4
	s.x[12] = uint64(uint32(int32(int8(s.x[8]))))
	s.compare32(uint32(s.x[12]), 8)
	s.x[17] = uint64(uint32(boolValue(s.condition("lt"))))
	s.x[1] = uint64(uint32(s.x[13]) ^ 0x27aa)
	s.x[17] = uint64(uint32(s.x[1]) * uint32(s.x[17]))
	s.x[14] = uint64(uint32(s.x[14]) + uint32(s.x[16]))
	s.x[14] = uint64(uint32(s.x[17]) ^ uint32(s.x[14]))
	next, err := loadBootstrapPointer(s, s.x[10], uint32(s.x[14]))
	if err != nil {
		return 0, fmt.Errorf("account signed-byte target: %w", err)
	}
	s.x[14] = next
	s.x[17] = 0xeaba
	return next, nil
}

func runAccountPointerByteTarget(s *state) (uint64, error) {
	s.traceBlock(accountPointerByteTarget)
	value, err := s.load(s.x[1]+8, 8)
	if err != nil {
		return 0, fmt.Errorf("account pointer-byte pointer: %w", err)
	}
	s.x[13] = value
	value, err = s.load(s.x[13], 1)
	if err != nil {
		return 0, fmt.Errorf("account pointer-byte input: %w", err)
	}
	s.x[8] = uint64(byte(value))
	s.x[14] = uint64(uint32(s.x[12]) ^ 0x3c)
	s.x[14] = uint64(uint32(s.x[14]) - 0x19)
	s.compare32(uint32(s.x[8]), uint32(byte(s.x[14])))
	s.x[8] = uint64(uint32(boolValue(s.condition("eq"))))
	s.x[14] = 0x41ed
	s.x[8] = uint64(uint32(s.x[8]) * uint32(s.x[14]))
	s.x[8] = uint64(uint32(s.x[8]) ^ uint32(s.x[12]))
	value, err = s.load(s.x[10]+uint64(uint32(s.x[8]))*8, 8)
	if err != nil {
		return 0, fmt.Errorf("account pointer-byte target: %w", err)
	}
	s.x[14] = value
	s.x[8] = 0xffff5bde
	s.x[17] = 0x9478
	return value, nil
}
