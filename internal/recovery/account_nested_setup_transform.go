package recovery

import "fmt"

const (
	accountNestedSetupStart  = uint64(0x100854fac)
	accountNestedSetupCall   = uint64(0x100855000)
	accountNestedSetupReturn = uint64(0x1008550ec)
)

func runAccountNestedSetup(s *state) (uint64, error) {
	s.traceBlock(accountNestedSetupStart)
	s.x[8] = uint64(uint32(s.x[20]) + 0x6759)
	s.x[28] = uint64(uint32(s.x[8]) | 0x40)
	s.x[8] = uint64(uint32(s.x[23]) ^ 0xf94d2fa3)
	s.x[9] = 0x06b2d05d
	s.x[8] = uint64(uint32(s.x[8]) + uint32(s.x[9]))
	s.x[9] = uint64(uint32(s.x[23]) ^ 0x926551a4)
	s.x[10] = uint64(uint32(s.x[23]) ^ 0xa55d05c2)
	s.x[11] = 0x5aa2fa3e
	s.x[10] = uint64(uint32(s.x[10]) + uint32(s.x[11]))
	s.x[8] = uint64(uint32(s.x[10]) ^ uint32(s.x[8]))
	s.x[10] = 0x6d9a0e1a

	s.traceBlock(accountNestedSetupCall)
	s.x[10] = uint64(uint32(s.x[28]) ^ uint32(s.x[10]))
	s.x[9] = uint64(uint32(s.x[10]) + uint32(s.x[9]))
	s.x[8] = uint64(uint32(s.x[8]) ^ uint32(s.x[9]))
	s.x[10] = s.x[29] - 0xb0
	s.x[9] = 0x6eaad3b0
	s.x[9] = uint64(uint32(s.x[10]) & uint32(s.x[9]))
	s.x[9] = uint64(uint32(s.x[10]) - uint32(s.x[9]<<1))
	s.x[10] = 0xeeaad3b3
	s.x[9] = uint64(uint32(s.x[9]) + uint32(s.x[10]))
	s.x[10] = 0x099801d1
	s.x[9] = uint64(uint32(s.x[9]) ^ uint32(s.x[10]))
	s.x[10] = 0xba209ef7
	s.x[9] = uint64(uint32(s.x[9]) * uint32(s.x[10]))
	s.x[10] = uint64(uint32(s.x[20]) - uint32(s.x[9]))
	s.x[10] = uint64(uint32(s.x[10]) - 0xc0f)
	if err := s.storeValue(addOffset(s.x[29], -0xb0), uint64(uint32(s.x[10])), 4); err != nil {
		return 0, fmt.Errorf("account nested-setup frame word: %w", err)
	}
	s.x[10] = s.sp + 0x2c0
	if err := s.storeValue(addOffset(s.x[29], -0x90), s.x[10], 8); err != nil {
		return 0, fmt.Errorf("account nested-setup first frame pointer: %w", err)
	}
	if err := s.storeValue(addOffset(s.x[29], -0xa8), s.x[6], 8); err != nil {
		return 0, fmt.Errorf("account nested-setup saved register: %w", err)
	}
	s.x[10] = s.sp + 0x3ec
	if err := s.storeValue(addOffset(s.x[29], -0x98), s.x[10], 8); err != nil {
		return 0, fmt.Errorf("account nested-setup second frame pointer: %w", err)
	}
	s.x[8] = uint64(uint32(s.x[8]) - uint32(s.x[9]))
	s.x[9] = 0x2afa09fc
	s.x[8] = uint64(uint32(s.x[8]) + uint32(s.x[9]))
	if err := s.storeValue(addOffset(s.x[29], -0x88), uint64(uint32(s.x[8])), 4); err != nil {
		return 0, fmt.Errorf("account nested-setup selector: %w", err)
	}
	s.x[8] = s.sp + 0x324
	if err := s.storeValue(addOffset(s.x[29], -0x80), s.x[8], 8); err != nil {
		return 0, fmt.Errorf("account nested-setup third frame pointer: %w", err)
	}
	s.x[8] = s.sp + 0x3f4
	if err := s.storeValue(addOffset(s.x[29], -0x78), s.x[8], 8); err != nil {
		return 0, fmt.Errorf("account nested-setup fourth frame pointer: %w", err)
	}
	first, err := s.load(s.sp+0x1f8, 8)
	if err != nil {
		return 0, fmt.Errorf("account nested-setup first source: %w", err)
	}
	second, err := s.load(s.sp+0x200, 8)
	if err != nil {
		return 0, fmt.Errorf("account nested-setup second source: %w", err)
	}
	s.x[9], s.x[8] = first, second
	firstWord, err := s.load(s.x[8], 4)
	if err != nil {
		return 0, fmt.Errorf("account nested-setup first source word: %w", err)
	}
	s.x[8] = uint64(uint32(firstWord))
	secondWord, err := s.load(s.x[9], 4)
	if err != nil {
		return 0, fmt.Errorf("account nested-setup second source word: %w", err)
	}
	s.x[9] = uint64(uint32(secondWord))
	s.x[8] = uint64(uint32(s.x[9]) - uint32(s.x[8]))
	s.x[9] = 0xa6
	s.x[8] = uint64(uint32(s.x[8]) ^ uint32(s.x[9]))
	s.x[8] = uint64(uint32(s.x[8]) * uint32(s.x[7]))
	value, err := s.load(s.sp+0x138, 8)
	if err != nil {
		return 0, fmt.Errorf("account nested-setup first byte table: %w", err)
	}
	s.x[9] = value + uint64(byte(s.x[8]))
	value, err = s.load(s.x[9], 1)
	if err != nil {
		return 0, fmt.Errorf("account nested-setup first byte: %w", err)
	}
	s.x[9] = uint64(byte(value)) ^ 0xb5
	value, err = s.load(s.sp+0x130, 8)
	if err != nil {
		return 0, fmt.Errorf("account nested-setup second byte table: %w", err)
	}
	s.x[10] = value
	value, err = s.load(value+uint64(uint32(s.x[9])), 1)
	if err != nil {
		return 0, fmt.Errorf("account nested-setup second byte: %w", err)
	}
	s.x[9] = uint64(byte(value))
	s.x[8] = uint64(uint32(s.x[8]) ^ uint32(s.x[9]))
	s.x[9] = 0xd59bef3f
	s.x[22] = uint64(uint32(s.x[20]) + uint32(s.x[9]))
	s.x[9] = 0xd002
	s.x[8] = uint64(uint32(s.x[9]) + uint32(byte(s.x[8])))
	next, err := s.load(s.x[30]+uint64(uint32(s.x[8]))*8, 8)
	if err != nil {
		return 0, fmt.Errorf("account nested-setup call target: %w", err)
	}
	s.x[8] = next
	s.x[0] = s.x[29] - 0xb0
	s.x[26] = s.x[6]
	s.x[17] = 0x48c7
	s.x[30] = accountNestedSetupReturn
	return next, nil
}
