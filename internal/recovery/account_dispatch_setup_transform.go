package recovery

import "fmt"

const (
	accountDispatchSeedSetup = uint64(0x100358590)
	accountDispatchMixSetup  = uint64(0x100379ddc)
)

func runAccountDispatchSeedSetup(s *state) (uint64, error) {
	s.traceBlock(accountDispatchSeedSetup)
	s.sp -= 0x10

	s.x[8] = 0x101074000
	value, err := s.load(s.x[8]+0x568, 8)
	if err != nil {
		return 0, fmt.Errorf("account dispatch-seed state load: %w", err)
	}
	s.x[9] = value

	s.x[10] = 0x101075000
	value, err = s.load(s.x[10]+0xb40, 8)
	if err != nil {
		return 0, fmt.Errorf("account dispatch-seed shared load: %w", err)
	}
	s.x[11] = value

	s.x[12] = s.sp + 8
	s.x[9] = s.x[12] + s.x[9]
	s.x[9] += s.x[11]
	s.x[11] = 0x0a5cd4cf
	s.x[12] = 0x673e3a80075f02c6
	s.x[12] = s.x[9]*s.x[11] + s.x[12]
	s.x[13] = 0x73d046208eaed1a6
	s.x[9] ^= s.x[13]
	s.x[13] = s.x[9] * s.x[11]

	if err := s.storeValue(s.x[8]+0x568, s.x[12], 8); err != nil {
		return 0, fmt.Errorf("account dispatch-seed state store: %w", err)
	}
	if err := s.storeValue(s.x[10]+0xb40, s.x[13], 8); err != nil {
		return 0, fmt.Errorf("account dispatch-seed shared store: %w", err)
	}

	value, err = s.load(s.x[0], 4)
	if err != nil {
		return 0, fmt.Errorf("account dispatch-seed input load: %w", err)
	}
	s.x[9] = uint64(uint32(value))
	value, err = s.load(s.x[1]+4, 4)
	if err != nil {
		return 0, fmt.Errorf("account dispatch-seed selector load: %w", err)
	}
	s.x[16] = uint64(uint32(value))

	s.x[8] = uint64(uint32(s.x[16]) ^ uint32(0x25de65a8))
	s.x[10] = 0xda219a58
	s.x[8] = uint64(uint32(s.x[8]) + uint32(s.x[10]))
	s.x[10] = uint64(uint32(s.x[16]) ^ uint32(0x8bc9a7be))
	s.x[14] = 0x74365842
	s.x[10] = uint64(uint32(s.x[10]) + uint32(s.x[14]))
	s.x[8] = uint64(uint32(s.x[8]) ^ uint32(s.x[10]))
	s.x[10] = uint64(uint32(s.x[16]) ^ uint32(0x6062b9d3))
	s.x[14] = 0x9f9d462d
	s.x[10] = uint64(uint32(s.x[10]) + uint32(s.x[14]))
	s.x[8] = uint64(uint32(s.x[8]) ^ uint32(s.x[10]))
	s.x[10] = uint64(uint32(s.x[9]) + uint32(0xce757bc6))
	s.x[8] = uint64(uint32(s.x[10]) + uint32(s.x[8]))
	s.compare32(uint32(s.x[8]), 0xc9)
	s.x[8] = uint64(uint32(boolValue(s.condition("lo"))))
	s.x[10] = 0x2a
	s.x[8] = uint64(uint32(s.x[8]) * uint32(s.x[10]))
	s.x[10] = 0x5859
	s.x[14] = 0xa6
	s.x[12] = uint64(uint32(s.x[13]) + uint32(s.x[12]))
	s.x[12] = uint64(uint32(s.x[12]) ^ uint32(s.x[14]))
	s.x[11] = uint64(uint32(s.x[12]) * uint32(s.x[11]))

	s.x[12] = 0x100fe9390
	s.x[12] += uint64(uint8(s.x[11]))
	value, err = s.load(s.x[12], 1)
	if err != nil {
		return 0, fmt.Errorf("account dispatch-seed byte-table load: %w", err)
	}
	s.x[12] = uint64(uint32(value))
	s.x[13] = 0x84
	s.x[12] = uint64(uint32(s.x[12]) ^ uint32(s.x[13]))
	s.x[13] = 0x100de2af0
	value, err = s.load(s.x[13]+uint64(uint32(s.x[12])), 1)
	if err != nil {
		return 0, fmt.Errorf("account dispatch-seed remap load: %w", err)
	}
	s.x[12] = uint64(uint32(value))
	s.x[11] = uint64(uint32(s.x[11]) ^ uint32(s.x[12]))
	s.x[11] = uint64(uint32(s.x[11]) & 0xff)
	s.x[11] = uint64(uint32(s.x[11]) + 0x257)
	s.x[12] = 0x10105b6f0
	value, err = s.load(s.x[12]+uint64(uint32(s.x[11]))*8, 8)
	if err != nil {
		return 0, fmt.Errorf("account dispatch-seed first target load: %w", err)
	}
	s.x[11] = value
	s.x[8] = uint64(uint32(s.x[8]) ^ uint32(s.x[10]))
	s.x[10] = 0xffffffff98245355
	s.x[10] = s.x[11] + s.x[10]
	value, err = s.load(s.x[10]+uint64(uint32(s.x[8]))*8, 8)
	if err != nil {
		return 0, fmt.Errorf("account dispatch-seed final target load: %w", err)
	}
	s.x[11] = value
	s.x[8] = 0xffff5bde
	s.x[17] = 0x30d6
	return s.x[11], nil
}

func runAccountDispatchMixSetup(s *state) (uint64, error) {
	s.traceBlock(accountDispatchMixSetup)
	s.sp -= 0x10

	s.x[8] = 0x101073000
	value, err := s.load(s.x[8]+0x50, 8)
	if err != nil {
		return 0, fmt.Errorf("account dispatch-mix state load: %w", err)
	}
	s.x[9] = value
	s.x[10] = 0x101075000
	value, err = s.load(s.x[10]+0xb40, 8)
	if err != nil {
		return 0, fmt.Errorf("account dispatch-mix shared load: %w", err)
	}
	s.x[11] = value
	s.x[9] = s.x[11] ^ s.x[9]
	s.x[11] = s.sp + 8
	s.x[12] = 0x73d046208eaed1a6
	s.x[9] += s.x[11]
	s.x[9] ^= s.x[12]
	s.x[11] = 0x0a5cd4cf
	s.x[9] *= s.x[11]
	if err := s.storeValue(s.x[8]+0x50, s.x[9], 8); err != nil {
		return 0, fmt.Errorf("account dispatch-mix state store: %w", err)
	}
	if err := s.storeValue(s.x[10]+0xb40, s.x[9], 8); err != nil {
		return 0, fmt.Errorf("account dispatch-mix shared store: %w", err)
	}

	value, err = s.load(s.x[0], 4)
	if err != nil {
		return 0, fmt.Errorf("account dispatch-mix input load: %w", err)
	}
	s.x[9] = uint64(uint32(value))
	s.x[11] = uint64(uint32(s.x[9]) + 1)
	s.x[8] = 0x100fe9000
	value, err = s.load(s.x[8]+0x1ca, 1)
	if err != nil {
		return 0, fmt.Errorf("account dispatch-mix byte load: %w", err)
	}
	s.x[8] = uint64(uint32(value))
	s.x[10] = 0xfa
	s.x[8] = uint64(uint32(s.x[8]) ^ uint32(s.x[10]))
	s.x[10] = 0x100de28f0
	value, err = s.load(s.x[10]+uint64(uint32(s.x[8])), 1)
	if err != nil {
		return 0, fmt.Errorf("account dispatch-mix remap load: %w", err)
	}
	s.x[8] = uint64(uint32(value))
	s.x[10] = 0x3a
	s.x[8] = uint64(uint32(s.x[8]) ^ uint32(s.x[10]))
	s.x[8] = uint64(uint32(s.x[8]) + 0x231)
	s.x[10] = 0x10105b6f0
	value, err = s.load(s.x[10]+uint64(uint32(s.x[8]))*8, 8)
	if err != nil {
		return 0, fmt.Errorf("account dispatch-mix first target load: %w", err)
	}
	s.x[12] = value
	value, err = s.load(s.x[1]+4, 4)
	if err != nil {
		return 0, fmt.Errorf("account dispatch-mix selector load: %w", err)
	}
	s.x[8] = uint64(uint32(value))
	s.x[10] = 0xffffffff98245355
	s.x[10] = s.x[12] + s.x[10]
	s.x[13] = 0xffffffff9829eb35
	s.x[12] += s.x[13]
	value, err = s.load(s.x[12], 8)
	if err != nil {
		return 0, fmt.Errorf("account dispatch-mix final target load: %w", err)
	}
	s.x[12] = value
	s.x[17] = 0x20d7
	return s.x[12], nil
}
