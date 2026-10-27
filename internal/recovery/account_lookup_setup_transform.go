package recovery

import "fmt"

const (
	accountLookupSetupTarget     = uint64(0x100091e5c)
	accountMultiplyCompareTarget = uint64(0x1001efe24)
)

func runAccountLookupSetupTarget(s *state) (uint64, error) {
	s.traceBlock(accountLookupSetupTarget)
	s.x[8] = 0xa135
	s.x[9] = uint64(uint32(s.x[11]) ^ uint32(s.x[8]))
	s.x[8] = uint64(uint32(s.x[9]) << 4)
	s.x[12] = 0x038626ed
	s.compare32(uint32(s.x[28]), uint32(s.x[12]))
	s.x[12] = uint64(uint32(boolValue(s.condition("eq"))))
	s.x[13] = 0x144826f0
	s.x[9] = uint64(uint32(s.x[9]) + uint32(s.x[13]))
	s.x[12] = uint64(uint32(s.x[12]) ^ uint32(s.x[9]))
	s.x[12] &= 1
	s.x[13] = uint64(uint32(s.x[8]) ^ 0x21ca)
	s.x[12] = uint64(uint32(s.x[12]) * uint32(s.x[13]))
	s.x[11] = uint64(uint32(s.x[12]) ^ uint32(s.x[11]))
	next, err := loadBootstrapPointer(s, s.x[20], uint32(s.x[11]))
	if err != nil {
		return 0, fmt.Errorf("account lookup-setup first target: %w", err)
	}
	s.x[11] = next
	value, err := s.load(s.x[6], 4)
	if err != nil {
		return 0, fmt.Errorf("account lookup-setup first word: %w", err)
	}
	s.x[12] = uint64(uint32(value))
	value, err = s.load(s.x[2]+0xb40, 4)
	if err != nil {
		return 0, fmt.Errorf("account lookup-setup second word: %w", err)
	}
	s.x[13] = uint64(uint32(value))
	s.x[12] = uint64(uint32(s.x[13]) + uint32(s.x[12]))
	s.x[12] = uint64(uint32(s.x[12]) ^ uint32(s.x[3]))
	s.x[12] = uint64(uint32(s.x[12]) * uint32(s.x[4]))
	s.x[13] = s.x[7] + uint64(byte(s.x[12]))
	value, err = s.load(s.x[13], 1)
	if err != nil {
		return 0, fmt.Errorf("account lookup-setup first byte: %w", err)
	}
	s.x[13] = uint64(byte(value)) ^ 0x35
	value, err = s.load(s.x[30]+uint64(uint32(s.x[13])), 1)
	if err != nil {
		return 0, fmt.Errorf("account lookup-setup second byte: %w", err)
	}
	s.x[13] = uint64(byte(value))
	s.x[12] = uint64(uint32(s.x[12]) ^ uint32(s.x[13]))
	s.x[12] &= 0xff
	s.x[12] = uint64(uint32(s.x[12]) + 0x34f)
	value, err = s.load(s.x[5]+uint64(uint32(s.x[12]))*8, 8)
	if err != nil {
		return 0, fmt.Errorf("account lookup-setup value: %w", err)
	}
	s.x[12] = value
	s.x[13] = 0xffffffffcdd21c01
	s.x[26] = s.x[12] + s.x[13]
	s.x[17] = 0x9421
	return next, nil
}

func runAccountMultiplyCompareTarget(s *state) (uint64, error) {
	s.traceBlock(accountMultiplyCompareTarget)
	s.x[8] = 0x75302381
	s.x[8] = uint64(uint32(s.x[9]) ^ uint32(s.x[8]))
	s.x[12] = 0x232
	s.x[8] = uint64(uint32(s.x[8]) * uint32(s.x[12]))
	s.x[12] = uint64(uint32(s.x[8]) ^ 0xce756579)
	s.compare32(uint32(s.x[10]), uint32(s.x[12]))
	s.x[12] = uint64(uint32(boolValue(s.condition("eq"))))
	s.x[13] = 0x89ac
	s.x[12] = uint64(uint32(s.x[12]) * uint32(s.x[13]))
	s.x[13] = uint64(uint32(s.x[9]) + 0x0acffb2d)
	s.x[12] = uint64(uint32(s.x[12]) ^ uint32(s.x[13]))
	next, err := loadBootstrapPointer(s, s.x[24], uint32(s.x[12]))
	if err != nil {
		return 0, fmt.Errorf("account multiply-compare target: %w", err)
	}
	s.x[12] = next
	s.x[20] = s.x[11]
	s.x[19] = s.x[9]
	s.x[13] = s.x[10]
	s.x[17] = 0xa86b
	return next, nil
}
