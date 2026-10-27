package recovery

import "fmt"

const (
	accountScaledCompareTarget = uint64(0x100912ef8)
	accountShiftedBitTarget    = uint64(0x100aa1b9c)
	accountCCmpMaskTarget      = uint64(0x100aabd8c)
)

func runAccountScaledCompareTarget(s *state) (uint64, error) {
	s.traceBlock(accountScaledCompareTarget)
	s.x[11] = 0xa28cf2ac
	s.x[11] = uint64(uint32(s.x[12]) + uint32(s.x[11]))
	s.x[13] = 0x5d72dbae
	s.x[11] = uint64(uint32(s.x[11]) & uint32(s.x[13]))
	s.x[13] = 0xffffcb64
	s.x[13] = uint64(uint32(s.x[12]) + uint32(s.x[13]))
	s.x[8] += s.x[8] << 1
	s.x[14] = s.x[19] + (s.x[8] << 4)
	value, err := s.load(s.x[14]+0x28, 4)
	if err != nil {
		return 0, fmt.Errorf("account scaled-compare input: %w", err)
	}
	s.x[8] = uint64(uint32(value))
	s.x[8] &= 8
	s.x[15] = 0xffff72b7
	s.x[15] = uint64(uint32(s.x[12]) + uint32(s.x[15]))
	s.compare32(uint32(s.x[8]), uint32(s.x[15]))
	s.x[8] = uint64(uint32(boolValue(s.condition("eq"))))
	s.x[15] = uint64(uint32(s.x[11]) ^ 0x5b81)
	s.x[8] = uint64(uint32(s.x[8]) * uint32(s.x[15]))
	s.x[8] = uint64(uint32(s.x[8]) ^ uint32(s.x[12]))
	next, err := loadBootstrapPointer(s, s.x[30], uint32(s.x[8]))
	if err != nil {
		return 0, fmt.Errorf("account scaled-compare target: %w", err)
	}
	s.x[8] = next
	s.x[17] = 0xcb88
	return next, nil
}

func runAccountShiftedBitTarget(s *state) (uint64, error) {
	s.traceBlock(accountShiftedBitTarget)
	s.x[13] = uint64(uint32(s.x[11]) - 0x2f7)
	s.x[8] = 0xffffadc0
	s.x[8] = uint64(uint32(s.x[11]) + uint32(s.x[8]))
	s.x[12] = 0x5280
	s.x[8] = uint64(uint32(s.x[8]) | uint32(s.x[12]))
	value, err := s.load(s.x[14]+0x24, 4)
	if err != nil {
		return 0, fmt.Errorf("account shifted-bit input: %w", err)
	}
	s.x[25] = uint64(uint32(value))
	s.x[12] = uint64(uint32(s.x[8]) ^ 0x04a794e7)
	s.compare32(uint32(s.x[25]), uint32(s.x[12]))
	s.x[12] = uint64(uint32(boolValue(s.condition("eq"))))
	s.x[11] = uint64(uint32(s.x[11]) ^ uint32(s.x[12]<<6))
	value, err = s.load(s.x[30]+uint64(uint32(s.x[11]))*8, 8)
	if err != nil {
		return 0, fmt.Errorf("account shifted-bit target: %w", err)
	}
	s.x[11] = value
	s.x[17] = 0x9e46
	return value, nil
}

func runAccountCCmpMaskTarget(s *state) (uint64, error) {
	s.traceBlock(accountCCmpMaskTarget)
	s.x[11] = uint64(uint32(s.x[8]) + 0x65cd)
	s.x[12] = uint64(uint32(s.x[8]) + 0xd8e72009)
	s.x[13] = uint64(uint32(s.x[12]) & 0x2718dcbf)
	s.x[12] = uint64(uint32(s.x[25]) + 0x449ce05b)
	s.x[14] = 0xc944af5e
	s.compare32(uint32(s.x[12]), uint32(s.x[14]))
	if s.condition("gt") {
		s.compare32(uint32(s.x[12]), uint32(s.x[28]))
	} else {
		s.setNZCV(0)
	}
	s.x[12] = uint64(uint32(boolValue(s.condition("lt"))))
	s.x[14] = uint64(uint32(s.x[13]) + 0x4071)
	s.x[12] = uint64(uint32(s.x[14]) * uint32(s.x[12]))
	s.x[8] = uint64(uint32(s.x[12]) ^ uint32(s.x[8]))
	next, err := loadBootstrapPointer(s, s.x[30], uint32(s.x[8]))
	if err != nil {
		return 0, fmt.Errorf("account ccmp-mask target: %w", err)
	}
	s.x[8] = next
	s.x[17] = 0x9e46
	return next, nil
}
