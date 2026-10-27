package recovery

import "fmt"

const (
	accountEqualityDispatch = uint64(0x1007f651c)
	accountByteCCmpDispatch = uint64(0x10098451c)
	accountOffsetDispatch   = uint64(0x10029397c)
)

func runAccountEqualityDispatch(s *state) (uint64, error) {
	s.traceBlock(accountEqualityDispatch)
	s.x[8] = uint64(uint32(s.x[5]) ^ uint32(s.x[16]))
	s.x[6] = uint64(uint32(s.x[8]) * uint32(s.x[1]))
	s.x[17] = s.x[6] ^ s.x[2]
	s.x[4] = uint64(uint32(s.x[8]) + uint32(s.x[0]))
	s.x[8] = s.x[4] ^ s.x[3]
	s.x[10] += s.x[17]
	s.compare64(s.x[10], s.x[8])
	s.x[8] = uint64(uint32(boolValue(s.condition("eq"))))
	s.x[17] = uint64(uint32(s.x[8]) << 5)
	s.x[8] = uint64(uint32(s.x[17]) - uint32(s.x[8]<<1))
	s.x[8] = uint64(uint32(s.x[8]) ^ uint32(s.x[5]))
	value, err := s.load(s.sp+0xdf0, 8)
	if err != nil {
		return 0, fmt.Errorf("account equality-dispatch table: %w", err)
	}
	s.x[17] = value
	next, err := loadBootstrapPointer(s, s.x[17], uint32(s.x[8]))
	if err != nil {
		return 0, fmt.Errorf("account equality-dispatch target: %w", err)
	}
	s.x[8] = next
	s.x[17] = 0xeb12
	return next, nil
}

func runAccountByteCCmpDispatch(s *state) (uint64, error) {
	s.traceBlock(accountByteCCmpDispatch)
	s.x[17] = s.x[11] + s.x[14]
	value, err := s.load(s.x[17]+0x468, 1)
	if err != nil {
		return 0, fmt.Errorf("account byte-ccmp first input: %w", err)
	}
	s.x[17] = uint64(byte(value))
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[16]))
	s.x[2] = s.x[10] + s.x[14]
	value, err = s.load(s.x[2], 1)
	if err != nil {
		return 0, fmt.Errorf("account byte-ccmp second input: %w", err)
	}
	s.x[2] = uint64(byte(value))
	s.x[2] = uint64(uint32(s.x[12]) ^ uint32(s.x[2]))
	s.x[2] = uint64(uint32(s.x[2]) ^ 0x7c)
	s.x[17] = uint64(uint32(s.x[17]) - uint32(byte(s.x[2])))
	s.x[2] = uint64(uint32(s.x[0]) & uint32(s.x[17]<<1))
	s.x[2] = uint64(uint32(s.x[17]) - uint32(s.x[2]))
	s.compare64(s.x[14], 0x13)
	s.x[14]++
	if s.condition("ne") {
		s.compare32(uint32(s.x[2]), 0)
	} else {
		s.setNZCV(0)
	}
	s.x[17] = uint64(uint32(boolValue(s.condition("ne"))))
	s.x[17] = uint64(uint32(s.x[17]) * uint32(s.x[1]))
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[8]))
	next, err := loadBootstrapPointer(s, s.x[30], uint32(s.x[17]))
	if err != nil {
		return 0, fmt.Errorf("account byte-ccmp dispatch target: %w", err)
	}
	s.x[3] = next
	s.x[17] = 0x8e77
	return next, nil
}

func runAccountOffsetDispatch(s *state) (uint64, error) {
	s.traceBlock(accountOffsetDispatch)
	s.x[14] = 0x1a5e
	s.x[17] = uint64(uint32(s.x[16]) + uint32(s.x[14]))
	s.x[14] = 0xd94b841a
	s.x[13] = uint64(uint32(s.x[13]) + uint32(s.x[1]))
	s.x[14] = uint64(uint32(s.x[13]) + uint32(s.x[14]))
	next, err := loadBootstrapPointer(s, s.x[10], uint32(s.x[17]))
	if err != nil {
		return 0, fmt.Errorf("account offset-dispatch target: %w", err)
	}
	s.x[13] = next
	s.x[17] = 0xdd5a
	return next, nil
}
