package recovery

import "fmt"

const (
	accountByteValueRoute = uint64(0x1009b72bc)
	accountWideGuardRoute = uint64(0x1009b2174)
	accountByteGuardRoute = uint64(0x1009f1804)
)

func runAccountByteValueRoute(s *state) (uint64, error) {
	s.traceBlock(accountByteValueRoute)
	s.x[13] = uint64(uint32(s.x[11]) & 0xff)
	s.x[8] = uint64(uint32(s.x[12]) ^ 0x2d60)
	s.x[11] = 0x6b
	s.x[8] = uint64(uint32(s.x[8]) * uint32(s.x[11]))
	s.x[11] = uint64(uint32(s.x[12]) ^ 0x2e9a)
	s.x[11] = uint64(uint32(s.x[11]) + uint32(s.x[11]<<1))
	s.x[11] = uint64(uint32(s.x[11]) << 1)
	s.x[9]++
	s.x[10]++
	s.x[14] = uint64(uint32(s.x[8]) ^ 0xdddddddd)
	s.compare32(uint32(s.x[13]), uint32(byte(s.x[14])))
	s.x[13] = uint64(uint32(boolValue(s.condition("eq"))))
	s.x[14] = uint64(uint32(s.x[11]) ^ 0xb18c)
	s.x[13] = uint64(uint32(s.x[13]) * uint32(s.x[14]))
	s.x[12] = uint64(uint32(s.x[13]) ^ uint32(s.x[12]))
	next, err := loadBootstrapPointer(s, s.x[30], uint32(s.x[12]))
	if err != nil {
		return 0, fmt.Errorf("account byte-value target: %w", err)
	}
	s.x[13] = next
	s.x[12] = 0x7b5db20a
	s.x[17] = 0xe762
	return next, nil
}

func runAccountWideGuardRoute(s *state) (uint64, error) {
	s.traceBlock(accountWideGuardRoute)
	s.x[8] = uint64(uint32(s.x[13]) + uint32(s.x[4]))
	s.x[8] = uint64(uint32(s.x[8]) | uint32(s.x[5]))
	s.x[14] = uint64(uint32(s.x[13]) + uint32(s.x[6]))
	s.x[14] = uint64(uint32(s.x[14]) | uint32(s.x[7]))
	s.compare32(uint32(s.x[22]), 0)
	s.x[10]++
	s.x[17] = s.x[8] ^ s.x[19]
	if s.condition("ne") {
		s.compare64(s.x[10], s.x[17])
	} else {
		s.setNZCV(4)
	}
	s.x[17] = uint64(uint32(boolValue(s.condition("ne"))))
	s.x[21] = uint64(uint32(s.x[14]) ^ uint32(s.x[20]))
	s.x[17] = uint64(uint32(s.x[17]) * uint32(s.x[21]))
	s.x[13] = uint64(uint32(s.x[17]) ^ uint32(s.x[13]))
	next, err := loadBootstrapPointer(s, s.x[30], uint32(s.x[13]))
	if err != nil {
		return 0, fmt.Errorf("account wide-guard target: %w", err)
	}
	s.x[21] = next
	s.x[13] = 0x30afef76
	s.x[17] = 0xbabc
	return next, nil
}

func runAccountByteGuardRoute(s *state) (uint64, error) {
	s.traceBlock(accountByteGuardRoute)
	s.x[8] = uint64(uint32(s.x[14]) ^ uint32(s.x[16]))
	s.x[13] = uint64(uint32(s.x[8]) * uint32(s.x[0]))
	s.x[8] = uint64(uint32(s.x[14]) ^ uint32(s.x[1]))
	s.x[21] = uint64(uint32(s.x[8]) * uint32(s.x[2]))
	value, err := s.load(s.x[12]+s.x[10], 1)
	if err != nil {
		return 0, fmt.Errorf("account byte-guard first input: %w", err)
	}
	s.x[17] = uint64(byte(value))
	s.x[22] = uint64(uint32(int32(int8(s.x[17]))))
	value, err = s.load(s.x[15]+s.x[10], 1)
	if err != nil {
		return 0, fmt.Errorf("account byte-guard second input: %w", err)
	}
	s.x[8] = uint64(byte(value))
	s.compare32(uint32(s.x[17]), uint32(s.x[8]))
	s.x[17] = uint64(uint32(boolValue(s.condition("eq"))))
	s.x[17] = uint64(uint32(s.x[13]) ^ uint32(s.x[17]))
	s.x[17] &= 1
	s.x[23] = uint64(uint32(s.x[21]) ^ uint32(s.x[3]))
	s.x[17] = uint64(uint32(s.x[17]) * uint32(s.x[23]))
	s.x[14] = uint64(uint32(s.x[17]) ^ uint32(s.x[14]))
	next, err := loadBootstrapPointer(s, s.x[30], uint32(s.x[14]))
	if err != nil {
		return 0, fmt.Errorf("account byte-guard target: %w", err)
	}
	s.x[14] = next
	s.x[17] = 0xac7c
	return next, nil
}
