package recovery

import "fmt"

const (
	accountByteBooleanStep = uint64(0x10030c11c)
	accountAccumulatorStep = uint64(0x100c3c104)
)

func runAccountByteBooleanStep(s *state) (uint64, error) {
	s.traceBlock(accountByteBooleanStep)
	s.x[14] = 0xffffe57c
	value, err := s.load(s.x[11]+uint64(uint32(s.x[4])), 1)
	if err != nil {
		return 0, fmt.Errorf("account byte-boolean input load: %w", err)
	}
	s.x[16] = uint64(uint32(value))
	s.x[4] = uint64(uint32(s.x[4]) + 1)
	s.x[17] = uint64(uint32(0x16) & uint32(shiftLeft(s.x[16], 1, 32)))
	s.x[5] = 0x8c1c1e0b
	s.x[16] = uint64(uint32(s.x[16]) - uint32(s.x[17]))
	s.x[17] = uint64(uint32(s.x[16]) + uint32(s.x[5]))
	s.x[16] = 0x2c0e98f2
	s.x[15] = uint64(uint32(s.x[16]) ^ uint32(shiftLeft(s.x[15], 8, 32)))
	s.x[5] = 0x73e3cd62
	s.x[16] = uint64(uint32(s.x[3]) + uint32(s.x[14]))
	s.x[14] = uint64(uint32(s.x[16]) ^ uint32(s.x[5]))
	s.x[14] = uint64(uint32(s.x[14]) ^ uint32(s.x[17]))
	s.x[14] = uint64(uint32(s.x[14]) & uint32(s.x[15]))
	s.x[15] = 0x8400020d
	s.x[15] = uint64(uint32(s.x[17]) & uint32(s.x[15]))
	s.x[17] = 0x1300243e
	s.x[15] = uint64(uint32(s.x[15]) ^ uint32(s.x[17]))
	s.x[15] = uint64(uint32(s.x[15]) ^ uint32(s.x[14]))
	s.x[2] = uint64(uint32(s.subtractWithFlags(uint64(uint32(s.x[2])), 1, 32)))
	s.x[14] = uint64(uint32(boolValue(s.condition("eq"))))
	s.x[17] = 0x2dc0
	s.x[14] = uint64(uint32(s.x[14]) * uint32(s.x[17]))
	s.x[14] = uint64(uint32(s.x[14]) ^ uint32(s.x[3]))
	value, err = s.load(s.x[10]+shiftLeft(signExtend32(uint32(s.x[14])), 3, 64), 8)
	if err != nil {
		return 0, fmt.Errorf("account byte-boolean target load: %w", err)
	}
	s.x[14] = value
	s.x[17] = 0xdd5a
	return s.x[14], nil
}

func runAccountAccumulatorStep(s *state) (uint64, error) {
	s.traceBlock(accountAccumulatorStep)
	s.x[17] = uint64(uint32(s.x[25]) ^ uint32(0x2bcd4aeb))
	s.x[3] = 0xe75b47b3
	s.x[17] = uint64(uint32(s.x[17]) + uint32(s.x[3]))
	s.x[3] = uint64(uint32(s.x[25]) ^ uint32(0xa9b68b56))
	s.x[4] = 0x65208610
	s.x[3] = uint64(uint32(s.x[3]) + uint32(s.x[4]))
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[3]))
	s.x[3] = 0x935c9020
	s.x[4] = uint64(uint32(s.x[25]) ^ uint32(s.x[3]))
	s.x[3] = uint64(uint32(s.x[9]) + 0x2d91)
	s.x[4] = uint64(uint32(s.x[4]) + uint32(s.x[3]))
	s.x[5] = 0x5fca12c0
	s.x[4] = uint64(uint32(s.x[4]) + uint32(s.x[5]))
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[4]))
	s.x[17] = uint64(uint32(s.x[10]) + uint32(s.x[17]))
	s.x[4] = uint64(uint32(s.x[17]) ^ uint32(0x97563477))
	s.x[25] = uint64(uint32(s.x[17]) ^ uint32(s.x[19]))
	s.x[25] = uint64(uint32(s.x[25]) + uint32(s.x[20]))
	s.x[26] = uint64(uint32(s.x[17]) ^ uint32(s.x[21]))
	s.x[5] = uint64(uint32(s.x[17]) ^ uint32(s.x[23]))
	s.x[4] = uint64(uint32(s.x[4]) + uint32(s.x[7]))
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[4]))
	s.x[4] = uint64(uint32(s.x[26]) + uint32(s.x[22]))
	s.x[4] = uint64(uint32(s.x[25]) ^ uint32(s.x[4]))
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[4]))
	s.x[4] = uint64(uint32(s.x[5]) + uint32(s.x[24]))
	s.x[4] = uint64(uint32(s.x[4]) ^ uint32(s.x[27]))
	s.x[4] = uint64(uint32(s.x[17]) ^ uint32(s.x[4]))
	s.x[17] = uint64(uint32(shiftLeft(1, int(uint32(s.x[2])), 32)))
	s.x[6] = uint64(uint32(s.x[17]) + uint32(s.x[6]))
	value, err := s.load(s.x[28]+shiftLeft(signExtend32(uint32(s.x[9])), 3, 64), 8)
	if err != nil {
		return 0, fmt.Errorf("account accumulator target load: %w", err)
	}
	s.x[9] = value
	s.x[17] = 0xc750
	return s.x[9], nil
}
