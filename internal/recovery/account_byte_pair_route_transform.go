package recovery

import "fmt"

const (
	accountByteDifferenceRoute     = uint64(0x1001bbfac)
	accountByteDifferenceRouteTail = uint64(0x1001bc000)
	accountByteSubtractRoute       = uint64(0x10026f664)
	accountByteCCmpRoute           = uint64(0x100273308)
)

func runAccountByteDifferenceRoute(s *state) (uint64, error) {
	s.traceBlock(accountByteDifferenceRoute)
	value, err := s.load(s.x[10]+s.x[12], 1)
	if err != nil {
		return 0, fmt.Errorf("account byte-difference first input: %w", err)
	}
	s.x[13] = uint64(byte(value))
	value, err = s.load(s.x[9]+s.x[12], 1)
	if err != nil {
		return 0, fmt.Errorf("account byte-difference second input: %w", err)
	}
	s.x[16] = uint64(byte(value))
	s.x[13] = uint64(uint32(s.x[13]) - uint32(s.x[16]))
	s.x[16] = uint64(uint32(s.x[11]) + 0x9ceadaf6)
	s.x[16] = uint64(uint32(s.x[16]) & uint32(s.x[13]<<1))
	s.x[13] = uint64(uint32(s.x[13]) - uint32(s.x[16]))
	s.x[16] = 0xce757be4
	s.x[13] = uint64(uint32(s.x[13]) + uint32(s.x[16]))
	s.x[16] = uint64(uint32(s.x[11]) + 0xce755f50)
	s.compare32(uint32(s.x[13]), uint32(s.x[16]))
	s.x[13] = uint64(uint32(boolValue(s.condition("eq"))))
	s.x[16] = 0xa905
	s.x[13] = uint64(uint32(s.x[13]) * uint32(s.x[16]))
	s.x[13] = uint64(uint32(s.x[13]) ^ uint32(s.x[11]))
	next, err := loadBootstrapPointer(s, s.x[20], uint32(s.x[13]))
	if err != nil {
		return 0, fmt.Errorf("account byte-difference target: %w", err)
	}
	s.x[13] = next
	s.x[0] = 0xffff5bde
	s.traceBlock(accountByteDifferenceRouteTail)
	s.x[17] = 0xe7a6
	return next, nil
}

func runAccountByteSubtractRoute(s *state) (uint64, error) {
	s.traceBlock(accountByteSubtractRoute)
	s.x[17] = uint64(uint32(s.x[16]) + 0xd68)
	s.x[15] = 0x6a62
	s.x[16] = uint64(uint32(s.x[17]) ^ uint32(s.x[15]))
	s.x[3] = 0x1a84
	s.x[15] = 0xd94b841b
	s.x[4] = uint64(uint32(s.x[13]) + uint32(s.x[15]))
	value, err := s.load(s.x[11]+uint64(uint32(s.x[14])), 1)
	if err != nil {
		return 0, fmt.Errorf("account byte-subtract input: %w", err)
	}
	s.x[14] = uint64(byte(value))
	s.x[15] = 0x18a
	s.x[15] = uint64(uint32(s.x[15]) & uint32(s.x[14]<<1))
	s.x[2] = 0xce757bc5
	s.x[14] = uint64(uint32(s.x[14]) - uint32(s.x[15]))
	s.x[15] = uint64(uint32(s.x[14]) + uint32(s.x[2]))
	s.x[2] = s.subtractWithFlags(uint64(uint32(s.x[1])), 1, 32)
	s.x[14] = uint64(uint32(boolValue(s.condition("ne"))))
	s.x[5] = 0x46fc
	s.x[3] = uint64(uint32(s.x[16]) + uint32(s.x[3]))
	s.x[5] = uint64(uint32(s.x[3]) ^ uint32(s.x[5]))
	s.x[14] = uint64(uint32(s.x[5]) * uint32(s.x[14]))
	s.x[14] = uint64(uint32(s.x[14]) ^ uint32(s.x[17]))
	next, err := loadBootstrapPointer(s, s.x[10], uint32(s.x[14]))
	if err != nil {
		return 0, fmt.Errorf("account byte-subtract target: %w", err)
	}
	s.x[5] = next
	s.x[14] = s.x[4]
	s.x[17] = 0xdd5a
	return next, nil
}

func runAccountByteCCmpRoute(s *state) (uint64, error) {
	s.traceBlock(accountByteCCmpRoute)
	s.x[16] = uint64(uint32(s.x[2]) ^ 0x31a3)
	s.x[17] = 0xffffe3e2
	s.x[16] = uint64(uint32(s.x[17]) + uint32(s.x[16]<<1))
	s.x[17] = 0x2914
	s.x[16] = uint64(uint32(s.x[16]) | uint32(s.x[17]))
	s.x[15] &= 0x7f
	s.x[17] = 0x45
	s.x[1] = uint64(uint32(s.x[15]) ^ uint32(s.x[17]))
	s.compare32(uint32(s.x[15]), 0x45)
	if s.condition("ne") {
		s.compare32(uint32(s.x[1]), 5)
	} else {
		s.setNZCV(2)
	}
	s.x[15] = uint64(uint32(boolValue(s.condition("lo"))))
	s.x[17] = uint64(uint32(s.x[16]) + 0x7041)
	s.x[15] = uint64(uint32(s.x[17]) * uint32(s.x[15]))
	s.x[15] = uint64(uint32(s.x[15]) ^ uint32(s.x[2]))
	next, err := loadBootstrapPointer(s, s.x[10], uint32(s.x[15]))
	if err != nil {
		return 0, fmt.Errorf("account byte-ccmp target: %w", err)
	}
	s.x[15] = next
	s.x[2] = 0x318a201b
	s.x[17] = 0x3cce
	return next, nil
}
