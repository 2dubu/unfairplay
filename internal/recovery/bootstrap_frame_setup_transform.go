package recovery

import "fmt"

const (
	bootstrapFrameSetup    = uint64(0x100579620)
	bootstrapDispatchSetup = uint64(0x10090d240)
)

func runBootstrapFrameSetup(s *state) (uint64, error) {
	s.traceBlock(bootstrapFrameSetup)
	s.x[10] = 0x64e21e1a
	word, err := s.load(s.sp+0x78, 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap frame-setup index: %w", err)
	}
	s.x[12] = uint64(uint32(word))
	s.x[10] = uint64(uint32(s.x[12]) + uint32(s.x[10]))
	if err := s.storeValue(s.sp+0x74, uint64(uint32(s.x[10])), 4); err != nil {
		return 0, fmt.Errorf("bootstrap frame-setup index output: %w", err)
	}
	s.x[10] = 0x2120
	s.x[20] = uint64(uint32(s.x[9]) + uint32(s.x[10]))
	value, err := s.load(s.sp+0x30, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap frame-setup descriptor: %w", err)
	}
	s.x[10] = value
	value, err = s.load(value+8, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap frame-setup descriptor base: %w", err)
	}
	s.x[10] = value
	if err := s.storeValue(s.sp+0x68, value, 8); err != nil {
		return 0, fmt.Errorf("bootstrap frame-setup saved base: %w", err)
	}
	word, err = s.load(value, 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap frame-setup descriptor word: %w", err)
	}
	s.x[10] = uint64(uint32(word))
	s.x[11] = 0x224e784e
	s.x[11] = uint64(uint32(s.x[9]) + uint32(s.x[11]))
	s.x[10] = uint64(uint32(s.x[10]) ^ uint32(s.x[11]))
	if err := s.storeValue(s.sp+0x60, s.x[10], 8); err != nil {
		return 0, fmt.Errorf("bootstrap frame-setup saved word: %w", err)
	}
	s.x[10] = 0xc1e1812b
	s.x[10] = uint64(uint32(s.x[12]) + uint32(s.x[10]))
	if err := s.storeValue(s.sp+0x3c, uint64(uint32(s.x[10])), 4); err != nil {
		return 0, fmt.Errorf("bootstrap frame-setup saved index: %w", err)
	}
	s.x[10] = 0x04c4d2ed
	s.x[22] = uint64(uint32(s.x[8]) + uint32(s.x[10]))
	s.x[10] = uint64(int64(int32(s.x[8])))
	s.x[11] = uint64(uint32(s.x[8]) - 1)
	s.x[12] = uint64(int64(int32(s.x[11])))
	value, err = s.load(s.sp+0x10, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap frame-setup first table: %w", err)
	}
	s.x[11] = value
	word, err = s.load(value+0x10, 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap frame-setup table word: %w", err)
	}
	s.x[11] = uint64(uint32(word))
	if err := s.storeValue(s.sp+0x5c, uint64(uint32(s.x[11])), 4); err != nil {
		return 0, fmt.Errorf("bootstrap frame-setup saved table word: %w", err)
	}
	value, err = s.load(s.sp+0x98, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap frame-setup second table: %w", err)
	}
	s.x[11] = value
	value, err = s.load(value+8, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap frame-setup dispatch base: %w", err)
	}
	s.x[11] = value
	if err := s.storeValue(s.sp+0x50, value, 8); err != nil {
		return 0, fmt.Errorf("bootstrap frame-setup saved dispatch base: %w", err)
	}
	next, err := s.load(s.x[25]+uint64(uint32(s.x[9]))*8, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap frame-setup dispatch: %w", err)
	}
	s.x[11] = next
	s.x[14] = 0x15089516
	s.x[24] = 0x224ea33b
	s.x[13] = 0x886438de
	s.x[4] = s.sp + 0x2d0
	s.x[6] = 0xb3563f76
	s.x[7] = 0x66ac7eed16cf0452
	s.x[26] = 0xb3563f768b678229
	s.x[27] = 0xa9292112
	s.x[9] = s.x[4] + (s.x[10] << 2)
	if err := s.storeValue(s.sp+0x48, s.x[9], 8); err != nil {
		return 0, fmt.Errorf("bootstrap frame-setup current slot: %w", err)
	}
	s.x[9] = s.x[4] + (s.x[12] << 2)
	if err := s.storeValue(s.sp+0x40, s.x[9], 8); err != nil {
		return 0, fmt.Errorf("bootstrap frame-setup previous slot: %w", err)
	}
	s.x[21] = 0x886438dd
	if err := s.storeValue(s.sp+0x38, uint64(uint32(s.x[22])), 4); err != nil {
		return 0, fmt.Errorf("bootstrap frame-setup saved selector: %w", err)
	}
	s.x[17] = 0xccdc
	return next, nil
}

func runBootstrapDispatchSetup(s *state) (uint64, error) {
	s.traceBlock(bootstrapDispatchSetup)
	word, err := s.load(s.sp+0x78, 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap dispatch-setup status: %w", err)
	}
	s.x[9] = uint64(uint32(word))
	s.compare32(uint32(s.x[9]), 1)
	selected := uint32(boolValue(s.condition("ne")))
	s.x[9] = uint64(selected)
	s.x[11] = 0x5237
	s.x[11] = uint64(uint32(s.x[20]) + uint32(s.x[11]))
	s.x[12] = 0x88
	s.x[3] = uint64(uint32(s.x[11]) | uint32(s.x[12]))
	s.x[11] = 0x750d
	s.x[11] = uint64(uint32(s.x[20]) + uint32(s.x[11]))
	s.x[12] = 0xeaf76aea
	if err := s.storeValue(s.sp+0x80, uint64(uint32(s.x[14])), 4); err != nil {
		return 0, fmt.Errorf("bootstrap dispatch-setup saved state: %w", err)
	}
	s.x[12] = uint64(uint32(s.x[14]) + uint32(s.x[12]))
	first, err := s.load(s.sp+0x60, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap dispatch-setup first pair: %w", err)
	}
	second, err := s.load(s.sp+0x68, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap dispatch-setup second pair: %w", err)
	}
	s.x[17], s.x[14] = first, second
	value, err := s.load(addOffset(s.x[14], int64(int32(s.x[12]))*4), 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap dispatch-setup indexed word: %w", err)
	}
	s.x[12] = uint64(uint32(value))
	word, err = s.load(s.sp+0x2d0, 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap dispatch-setup stack word: %w", err)
	}
	s.x[14] = uint64(uint32(word)) ^ s.x[24]
	s.x[12] ^= s.x[24]
	s.x[28] = uint64(uint32(s.x[12]))*uint64(uint32(s.x[17])) + s.x[14]
	word, err = s.load(s.sp+0x5c, 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap dispatch-setup multiplier: %w", err)
	}
	s.x[14] = uint64(uint32(word))
	s.x[19] = uint64(uint32(s.x[28]) * uint32(s.x[14]))
	s.x[14] = uint64(uint32(s.x[28]))
	value, err = s.load(s.sp+0x50, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap dispatch-setup word base: %w", err)
	}
	s.x[17] = value
	word, err = s.load(value, 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap dispatch-setup base word: %w", err)
	}
	s.x[17] = uint64(uint32(word))
	s.x[2] = 0x224df1ea
	s.x[2] = uint64(uint32(s.x[20]) + uint32(s.x[2]))
	s.x[2] = uint64(uint32(s.x[3]) ^ uint32(s.x[2]))
	s.x[17] = uint64(uint32(s.x[2]) ^ uint32(s.x[17]))
	s.x[14] = uint64(uint32(s.x[19]))*uint64(uint32(s.x[17])) + s.x[14]
	s.x[17] = 0x99538112e930fbad
	s.x[17] |= ^(s.x[14] << 1)
	s.x[2] = 0xb3563f768b67822a
	s.x[14] += s.x[2]
	s.x[30] = s.x[14] + s.x[17]
	s.x[14] = 0x17
	s.x[9] = uint64(selected * uint32(s.x[14]))
	s.x[9] = uint64(uint32(s.x[20]) ^ uint32(s.x[9]))
	next, err := loadBootstrapPointer(s, s.x[25], uint32(s.x[9]))
	if err != nil {
		return 0, fmt.Errorf("bootstrap dispatch-setup target: %w", err)
	}
	s.x[9] = next
	s.x[2] = 0x779bc723
	s.x[10] = 0xa40ddddf
	s.x[23] = 0xc13a
	s.x[16] = 0x4a45ba0a
	s.x[15] = 0x41e1812a
	s.x[0] = 0xe663
	word, err = s.load(s.sp+0x3c, 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap dispatch-setup saved selector: %w", err)
	}
	s.x[1] = uint64(uint32(word))
	s.x[17] = 0xcd95
	return next, nil
}
