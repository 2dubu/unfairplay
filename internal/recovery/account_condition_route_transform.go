package recovery

import "fmt"

const (
	accountConditionRoute     = uint64(0x100b073e0)
	accountConditionRouteTail = uint64(0x100b07400)
	accountSignedIndexRoute   = uint64(0x100016c90)
	accountStackCompareRoute  = uint64(0x10002da1c)
)

func runAccountConditionRoute(s *state) (uint64, error) {
	s.traceBlock(accountConditionRoute)
	s.x[20] = s.x[13]
	s.compare32(uint32(s.x[8]), uint32(s.x[9]))
	s.x[13] = uint64(uint32(boolValue(s.condition("lt"))))
	s.compare32(uint32(s.x[20]), uint32(s.x[11]))
	s.x[14] = uint64(uint32(boolValue(s.condition("gt"))))
	s.x[13] = uint64(uint32(s.x[13]) ^ uint32(s.x[14]))
	s.x[15] = uint64(uint32(s.x[20]) + uint32(s.x[10]))
	s.compare32(uint32(s.x[15]), uint32(s.x[8]))
	s.traceBlock(accountConditionRouteTail)
	s.x[15] = uint64(uint32(boolValue(s.condition("gt"))))
	s.compare32(uint32(s.x[13]), 0)
	if s.condition("ne") {
		s.x[14] = uint64(uint32(s.x[14]))
	} else {
		s.x[14] = uint64(uint32(s.x[15]))
	}
	s.x[13] = uint64(uint32(s.x[20]) << 1)
	s.x[14] = uint64(uint32(s.x[14]) * uint32(s.x[12]))
	s.x[14] = uint64(uint32(s.x[14]) ^ uint32(s.x[23]))
	next, err := loadBootstrapPointer(s, s.x[28], uint32(s.x[14]))
	if err != nil {
		return 0, fmt.Errorf("account condition-route target: %w", err)
	}
	s.x[14] = next
	s.x[17] = 0x3bbe
	return next, nil
}

func runAccountSignedIndexRoute(s *state) (uint64, error) {
	s.traceBlock(accountSignedIndexRoute)
	s.x[10] = 0x6c95
	s.x[13] = uint64(uint32(s.x[12]) ^ uint32(s.x[10]))
	s.x[10] = uint64(uint32(s.x[12]) + 0x0047d4cc)
	s.x[11] = uint64(uint32(s.x[10]) | 0x14002109)
	s.x[14] = uint64(uint32(s.x[28]) + uint32(s.x[22]))
	s.x[10] = uint64(int64(int32(s.x[14])))
	s.x[14] = s.x[10] + (uint64(int64(int32(s.x[14]))) << 1)
	s.x[14] <<= 4
	value, err := s.load(s.x[26]+s.x[14], 4)
	if err != nil {
		return 0, fmt.Errorf("account signed-index input: %w", err)
	}
	s.x[14] = uint64(uint32(value))
	s.x[15] = uint64(uint32(s.x[11]) ^ 0x03a428ae)
	s.compare32(uint32(s.x[14]), uint32(s.x[15]))
	s.x[15] = uint64(uint32(boolValue(s.condition("eq"))))
	s.x[16] = 0xffffa0b4
	s.x[16] = uint64(uint32(s.x[13]) + uint32(s.x[16]))
	s.x[15] = uint64(uint32(s.x[15]) * uint32(s.x[16]))
	s.x[12] = uint64(uint32(s.x[15]) ^ uint32(s.x[12]))
	next, err := loadBootstrapPointer(s, s.x[20], uint32(s.x[12]))
	if err != nil {
		return 0, fmt.Errorf("account signed-index target: %w", err)
	}
	s.x[12] = next
	s.x[17] = 0x8448
	return next, nil
}

func runAccountStackCompareRoute(s *state) (uint64, error) {
	s.traceBlock(accountStackCompareRoute)
	s.x[8] = 0xffffee7b
	s.x[8] = uint64(uint32(s.x[26]) + uint32(s.x[8]))
	s.x[9] = 0x108
	s.x[10] = uint64(uint32(s.x[8]) | uint32(s.x[9]))
	value, err := s.load(s.sp+0x10, 8)
	if err != nil {
		return 0, fmt.Errorf("account stack-compare pointer: %w", err)
	}
	s.x[12] = value
	value, err = s.load(s.x[12]+0x20, 4)
	if err != nil {
		return 0, fmt.Errorf("account stack-compare input: %w", err)
	}
	s.x[11] = uint64(uint32(value))
	s.x[8] = uint64(uint32(s.x[10]) ^ 0x03863fb0)
	s.compare32(uint32(s.x[11]), uint32(s.x[8]))
	s.x[8] = uint64(uint32(boolValue(s.condition("eq"))))
	s.x[9] = uint64(uint32(s.x[8]) << 1)
	s.x[8] = uint64(uint32(s.x[9]) | uint32(s.x[8]<<2))
	s.x[8] = uint64(uint32(s.x[8]) ^ uint32(s.x[26]))
	next, err := loadBootstrapPointer(s, s.x[20], uint32(s.x[8]))
	if err != nil {
		return 0, fmt.Errorf("account stack-compare target: %w", err)
	}
	s.x[8] = next
	s.x[17] = 0x1283
	return next, nil
}
