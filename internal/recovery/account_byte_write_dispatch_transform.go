package recovery

import "fmt"

const (
	accountByteWriteDispatchA = uint64(0x10053137c)
	accountByteWriteDispatchB = uint64(0x100641d5c)
	accountGreaterEqualRoute  = uint64(0x1008f8160)
)

func runAccountByteWriteDispatchA(s *state) (uint64, error) {
	s.traceBlock(accountByteWriteDispatchA)
	value, err := s.load(s.x[14]+uint64(uint32(s.x[9])), 1)
	if err != nil {
		return 0, fmt.Errorf("account byte-write dispatch A input: %w", err)
	}
	s.x[17] = uint64(byte(value))
	s.x[8] = uint64(uint32(s.x[10]) + uint32(s.x[11]))
	s.x[8] = uint64(uint32(s.x[8]) & uint32(s.x[12]))
	s.x[2] = uint64(uint32(s.x[8]) + 0x1b)
	s.x[2] = uint64(uint32(s.x[2]) ^ uint32(s.x[17]))
	if err := s.storeValue(s.x[19]+s.x[9], uint64(byte(s.x[2])), 1); err != nil {
		return 0, fmt.Errorf("account byte-write dispatch A output: %w", err)
	}
	s.compare32(uint32(s.x[17]), 0x1b)
	s.x[17] = uint64(uint32(boolValue(s.condition("eq"))))
	s.x[2] = uint64(uint32(s.x[17]) << 6)
	s.x[17] = uint64(uint32(s.x[2]) - uint32(s.x[17]<<3))
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[10]))
	value, err = s.load(s.sp+0xdf0, 8)
	if err != nil {
		return 0, fmt.Errorf("account byte-write dispatch A table: %w", err)
	}
	s.x[2] = value
	next, err := loadBootstrapPointer(s, s.x[2], uint32(s.x[17]))
	if err != nil {
		return 0, fmt.Errorf("account byte-write dispatch A target: %w", err)
	}
	s.x[2] = next
	s.x[17] = 0xd936
	return next, nil
}

func runAccountByteWriteDispatchB(s *state) (uint64, error) {
	s.traceBlock(accountByteWriteDispatchB)
	s.x[8] = uint64(uint32(s.x[4]) ^ uint32(s.x[9]))
	s.x[8] = uint64(uint32(s.x[8]) * uint32(s.x[11]))
	value, err := s.load(s.x[12]+uint64(uint32(s.x[10])), 1)
	if err != nil {
		return 0, fmt.Errorf("account byte-write dispatch B input: %w", err)
	}
	s.x[17] = uint64(byte(value))
	s.x[5] = uint64(uint32(s.x[4]) ^ uint32(s.x[13]))
	s.x[6] = uint64(uint32(s.x[17]) ^ uint32(s.x[14]))
	if err := s.storeValue(s.x[19]+s.x[10], uint64(byte(s.x[6])), 1); err != nil {
		return 0, fmt.Errorf("account byte-write dispatch B output: %w", err)
	}
	s.x[6] = uint64(uint32(s.x[8]) - 0x53)
	s.x[6] = uint64(uint32(s.x[6]) ^ uint32(s.x[5]))
	s.compare32(uint32(s.x[17]), uint32(byte(s.x[6])))
	s.x[17] = uint64(uint32(boolValue(s.condition("ne"))))
	s.x[17] = uint64(uint32(s.x[17]) * uint32(s.x[15]))
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[4]))
	value, err = s.load(s.sp+0xdf0, 8)
	if err != nil {
		return 0, fmt.Errorf("account byte-write dispatch B table: %w", err)
	}
	s.x[4] = value
	next, err := loadBootstrapPointer(s, s.x[4], uint32(s.x[17]))
	if err != nil {
		return 0, fmt.Errorf("account byte-write dispatch B target: %w", err)
	}
	s.x[4] = next
	s.x[17] = 0x6008
	return next, nil
}

func runAccountGreaterEqualRoute(s *state) (uint64, error) {
	s.traceBlock(accountGreaterEqualRoute)
	s.x[11] = 0xffff5c30
	s.x[13] = uint64(uint32(s.x[14]) + uint32(s.x[11]))
	s.x[11] = 0x04df8802
	s.x[11] = uint64(uint32(s.x[14]) + uint32(s.x[11]))
	s.compare32(uint32(s.x[9]), uint32(s.x[11]))
	s.x[11] = uint64(uint32(boolValue(s.condition("ge"))))
	s.x[12] = 0xb
	s.x[11] = uint64(uint32(s.x[11]) * uint32(s.x[12]))
	s.x[11] = uint64(uint32(s.x[11]) ^ uint32(s.x[14]))
	next, err := loadBootstrapPointer(s, s.x[30], uint32(s.x[11]))
	if err != nil {
		return 0, fmt.Errorf("account greater-equal target: %w", err)
	}
	s.x[11] = next
	s.x[17] = 0xd242
	return next, nil
}
