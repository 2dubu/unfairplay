package recovery

import "fmt"

const (
	bootstrapNestedSetupStart  = uint64(0x10053db58)
	bootstrapNestedSetupCall   = uint64(0x10053dc00)
	bootstrapNestedSetupResume = uint64(0x10053dc84)
	bootstrapNestedSetupReturn = uint64(0x10053dd44)
)

func runBootstrapNestedSetupStart(s *state) (uint64, error) {
	s.traceBlock(bootstrapNestedSetupStart)
	s.x[8] = 0x262e
	s.x[28] = s.x[29] - 0x90
	s.x[9] = 0x164e578a
	s.x[9] = uint64(uint32(s.x[28]) | uint32(s.x[9]))
	s.x[9] = uint64(uint32(s.x[9]) ^ 0xfffffffe)
	s.x[10] = 0xe9b1a875
	s.x[10] = uint64(uint32(s.x[10]) | ^uint32(s.x[28]))
	s.x[11] = 0x9323ef83
	s.x[9] = uint64(uint32(s.x[9]) - uint32(s.x[10]))
	s.x[9] = uint64(uint32(s.x[9]) ^ uint32(s.x[11]))
	s.x[10] = 0x31b6b333
	first, err := s.load(s.sp+0x50, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap nested setup first saved value: %w", err)
	}
	second, err := s.load(s.sp+0x58, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap nested setup second saved value: %w", err)
	}
	s.x[11], s.x[24] = first, second
	if err := s.storeValue(addOffset(s.x[29], -0x90), s.x[11], 8); err != nil {
		return 0, fmt.Errorf("bootstrap nested setup first frame value: %w", err)
	}
	value, err := s.load(s.sp+0x48, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap nested setup third saved value: %w", err)
	}
	s.x[11] = value
	if err := s.storeValue(addOffset(s.x[29], -0x70), value, 8); err != nil {
		return 0, fmt.Errorf("bootstrap nested setup second frame value: %w", err)
	}
	s.x[11] = s.sp + 0x88
	if err := s.storeValue(addOffset(s.x[29], -0x80), s.x[11], 8); err != nil {
		return 0, fmt.Errorf("bootstrap nested setup first pair value: %w", err)
	}
	if err := s.storeValue(addOffset(s.x[29], -0x78), s.x[11], 8); err != nil {
		return 0, fmt.Errorf("bootstrap nested setup second pair value: %w", err)
	}
	s.x[9] = uint64(uint32(s.x[21]) - uint32(s.x[9])*uint32(s.x[10]))
	s.x[10] = 0xa433
	s.x[9] = uint64(uint32(s.x[9]) + uint32(s.x[10]))
	if err := s.storeValue(addOffset(s.x[29], -0x88), uint64(uint32(s.x[9])), 4); err != nil {
		return 0, fmt.Errorf("bootstrap nested setup frame word: %w", err)
	}
	word, err := s.load(s.x[24], 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap nested setup input word: %w", err)
	}
	s.x[9] = uint64(uint32(word))
	s.x[10] = bootstrapStateBase
	word, err = s.load(s.x[10]+0xb40, 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap nested setup state word: %w", err)
	}
	s.x[10] = uint64(uint32(word))
	s.x[11] = 0xa6
	s.x[9] = uint64(uint32(s.x[9]) ^ uint32(s.x[11]))
	s.x[12] = 0xa6
	s.x[9] = uint64(uint32(s.x[9]) ^ uint32(s.x[10]))
	s.x[25] = 0x0a5cd4cf
	s.x[9] = uint64(uint32(s.x[9]) * uint32(s.x[25]))
	value, err = s.load(s.sp+0x28, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap nested setup first byte table: %w", err)
	}
	s.x[10] = value + uint64(byte(s.x[9]))
	value, err = s.load(s.x[10], 1)
	if err != nil {
		return 0, fmt.Errorf("bootstrap nested setup first byte: %w", err)
	}
	s.x[10] = uint64(byte(value)) ^ 0x5d
	value, err = s.load(s.sp+0x20, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap nested setup second byte table: %w", err)
	}
	s.x[11] = value
	value, err = s.load(value+uint64(uint32(s.x[10])), 1)
	if err != nil {
		return 0, fmt.Errorf("bootstrap nested setup second byte: %w", err)
	}
	s.x[10] = uint64(byte(value))
	return bootstrapNestedSetupCall, nil
}

func runBootstrapNestedSetupCall(s *state) (uint64, error) {
	s.traceBlock(bootstrapNestedSetupCall)
	s.x[9] = uint64(uint32(s.x[9]) ^ uint32(s.x[10]))
	s.x[10] = 0x1448
	s.x[9] = uint64(uint32(s.x[10]) + uint32(byte(s.x[9])))
	s.x[10] = bootstrapTargetTable
	value, err := s.load(s.x[10]+uint64(uint32(s.x[9]))*8, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap nested setup first target table: %w", err)
	}
	s.x[9] = value
	s.x[10] = 0xffffffffe2130ede
	s.x[9] += s.x[10]
	value, err = s.load(s.x[9], 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap nested setup first target: %w", err)
	}
	s.x[9] = value
	if err := s.storeValue(s.sp+0x98, value, 8); err != nil {
		return 0, fmt.Errorf("bootstrap nested setup saved target: %w", err)
	}
	word, err := s.load(s.x[24], 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap nested setup refreshed input word: %w", err)
	}
	s.x[9] = uint64(uint32(word))
	value, err = s.load(s.sp+0x60, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap nested setup state pointer: %w", err)
	}
	s.x[22] = value
	word, err = s.load(value, 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap nested setup refreshed state word: %w", err)
	}
	s.x[10] = uint64(uint32(word))
	s.x[9] = uint64(uint32(s.x[10]) + uint32(s.x[9]))
	s.x[9] = uint64(uint32(s.x[9]) ^ uint32(s.x[12]))
	s.x[9] = uint64(uint32(s.x[9]) * uint32(s.x[25]))
	first, err := s.load(s.sp+0x38, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap nested setup third byte table: %w", err)
	}
	second, err := s.load(s.sp+0x40, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap nested setup fourth byte table: %w", err)
	}
	s.x[23], s.x[26] = first, second
	s.x[10] = s.x[26] + uint64(byte(s.x[9]))
	value, err = s.load(s.x[10], 1)
	if err != nil {
		return 0, fmt.Errorf("bootstrap nested setup third byte: %w", err)
	}
	s.x[10] = uint64(byte(value)) ^ 0xb5
	value, err = s.load(s.x[23]+uint64(uint32(s.x[10])), 1)
	if err != nil {
		return 0, fmt.Errorf("bootstrap nested setup fourth byte: %w", err)
	}
	s.x[10] = uint64(byte(value))
	s.x[19] = uint64(uint32(s.x[21]) + uint32(s.x[8]))
	s.x[8] = uint64(uint32(s.x[9]) ^ uint32(s.x[10]))
	s.x[9] = 0xcfa6
	s.x[8] = uint64(uint32(s.x[9]) + uint32(byte(s.x[8])))
	value, err = s.load(s.sp+0x70, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap nested setup dispatch table: %w", err)
	}
	s.x[27] = value
	next, err := s.load(value+uint64(uint32(s.x[8]))*8, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap nested setup call target: %w", err)
	}
	s.x[8] = next
	s.x[0] = s.x[29] - 0x90
	s.x[1] = s.sp + 0x98
	s.x[17] = 0x3442
	s.x[30] = bootstrapNestedSetupResume
	return next, nil
}

func runBootstrapNestedSetupResume(s *state) (uint64, error) {
	s.traceBlock(bootstrapNestedSetupResume)
	s.x[8] = s.sp
	value, err := s.load(s.sp+0x98, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap nested resume stack adjustment: %w", err)
	}
	s.x[9] = value
	s.x[8] += value
	s.sp = s.x[8]
	s.x[8] = 0x965d9902
	if err := s.storeValue(s.sp+0x6c, uint64(uint32(s.x[19])), 4); err != nil {
		return 0, fmt.Errorf("bootstrap nested resume saved selector: %w", err)
	}
	s.x[8] = uint64(uint32(s.x[19]) ^ uint32(s.x[8]))
	s.x[9] = 0x4e8c5470
	s.x[9] = uint64(uint32(s.x[28]) & uint32(s.x[9]))
	s.x[9] = uint64(uint32(s.x[28]) - uint32(s.x[9]<<1))
	s.x[10] = 0x4e8c5472
	s.x[11] = 0xb758cd06
	s.x[9] = uint64(uint32(s.x[9]) + uint32(s.x[10]))
	s.x[9] = uint64(uint32(s.x[9]) ^ uint32(s.x[11]))
	s.x[10] = 0xda51a827
	s.x[9] = uint64(uint32(s.x[9]) * uint32(s.x[10]))
	s.x[10] = uint64(uint32(s.x[9]) + uint32(s.x[21]))
	s.x[10] = uint64(uint32(s.x[10]) - 0xee)
	if err := s.storeValue(addOffset(s.x[29], -0x8c), uint64(uint32(s.x[10])), 4); err != nil {
		return 0, fmt.Errorf("bootstrap nested resume frame word: %w", err)
	}
	value, err = s.load(s.sp+0x30, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap nested resume saved pointer: %w", err)
	}
	s.x[10] = value
	if err := s.storeValue(addOffset(s.x[29], -0x88), value, 8); err != nil {
		return 0, fmt.Errorf("bootstrap nested resume frame pointer: %w", err)
	}
	s.x[8] = uint64(uint32(s.x[8]) + uint32(s.x[20]))
	s.x[8] = uint64(uint32(s.x[8]) ^ uint32(s.x[9]))
	if err := s.storeValue(addOffset(s.x[29], -0x80), uint64(uint32(s.x[8])), 4); err != nil {
		return 0, fmt.Errorf("bootstrap nested resume frame selector: %w", err)
	}
	word, err := s.load(s.x[24], 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap nested resume input word: %w", err)
	}
	s.x[8] = uint64(uint32(word))
	word, err = s.load(s.x[22], 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap nested resume state word: %w", err)
	}
	s.x[9] = uint64(uint32(word))
	s.x[8] = uint64(uint32(s.x[9]) + uint32(s.x[8]))
	s.x[9] = 0xa6
	s.x[8] = uint64(uint32(s.x[8]) ^ uint32(s.x[9]))
	s.x[8] = uint64(uint32(s.x[8]) * uint32(s.x[25]))
	s.x[9] = s.x[26] + uint64(byte(s.x[8]))
	value, err = s.load(s.x[9], 1)
	if err != nil {
		return 0, fmt.Errorf("bootstrap nested resume first byte: %w", err)
	}
	s.x[9] = uint64(byte(value))
	s.x[25] = 0xb5
	s.x[10] = 0xb5
	s.x[9] = uint64(uint32(s.x[9]) ^ uint32(s.x[10]))
	value, err = s.load(s.x[23]+uint64(uint32(s.x[9])), 1)
	if err != nil {
		return 0, fmt.Errorf("bootstrap nested resume second byte: %w", err)
	}
	s.x[9] = uint64(byte(value))
	s.x[8] = uint64(uint32(s.x[8]) ^ uint32(s.x[9]))
	s.x[9] = 0xd075
	s.x[8] = uint64(uint32(s.x[9]) + uint32(byte(s.x[8])))
	next, err := s.load(s.x[27]+uint64(uint32(s.x[8]))*8, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap nested resume call target: %w", err)
	}
	s.x[8] = next
	s.x[0] = s.x[29] - 0x90
	s.x[17] = 0xa288
	s.x[30] = bootstrapNestedSetupReturn
	return next, nil
}
