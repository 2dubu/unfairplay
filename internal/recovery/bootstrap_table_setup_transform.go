package recovery

import "fmt"

const bootstrapTableSetup = uint64(0x100a09ea0)

func runBootstrapTableSetup(s *state) (uint64, error) {
	s.traceBlock(bootstrapTableSetup)
	s.x[12] = uint64(uint32(s.x[8]) + uint32(s.x[7]))
	s.x[11] = uint64(uint32(s.x[12]) ^ uint32(s.x[21]))
	s.x[12] = uint64(uint32(s.x[12]) & uint32(s.x[22]))
	s.compare32(uint32(s.x[10]), 0x20)
	gate := uint32(boolValue(s.condition("eq")))
	s.x[13] = uint64(gate)
	s.x[13] = uint64(gate ^ uint32(s.x[11]))
	s.x[13] &= 1
	s.x[14] = uint64(uint32(s.x[12]) + uint32(s.x[23]))
	s.x[13] = uint64(uint32(s.x[13]) * uint32(s.x[14]))
	s.x[8] = uint64(uint32(s.x[13]) ^ uint32(s.x[8]))
	next, err := loadBootstrapPointer(s, s.x[19], uint32(s.x[8]))
	if err != nil {
		return 0, fmt.Errorf("bootstrap table-setup target: %w", err)
	}
	s.x[8] = next
	word, err := s.load(s.x[3], 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap table-setup input word: %w", err)
	}
	s.x[13] = uint64(uint32(word))
	word, err = s.load(s.x[25]+0xb40, 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap table-setup state word: %w", err)
	}
	s.x[14] = uint64(uint32(word))
	s.x[15] = uint64(uint32(s.x[14]) + uint32(s.x[13]))
	s.x[15] = uint64(uint32(s.x[15]) ^ uint32(s.x[26]))
	s.x[15] = uint64(uint32(s.x[15]) * uint32(s.x[28]))
	value, err := s.load(s.x[2]+0x1af8, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap table-setup first table: %w", err)
	}
	s.x[16] = value + s.x[24]
	if err := s.storeValue(s.sp+0x190, s.x[16], 8); err != nil {
		return 0, fmt.Errorf("bootstrap table-setup saved first table: %w", err)
	}
	s.x[16] += uint64(byte(s.x[15]))
	value, err = s.load(s.x[16], 1)
	if err != nil {
		return 0, fmt.Errorf("bootstrap table-setup first byte: %w", err)
	}
	s.x[16] = uint64(byte(value)) ^ 0x17
	value, err = s.load(s.x[2]+0x2a58, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap table-setup second table: %w", err)
	}
	s.x[17] = value + s.x[30]
	if err := s.storeValue(s.sp+0x188, s.x[17], 8); err != nil {
		return 0, fmt.Errorf("bootstrap table-setup saved second table: %w", err)
	}
	value, err = s.load(s.x[17]+uint64(uint32(s.x[16])), 1)
	if err != nil {
		return 0, fmt.Errorf("bootstrap table-setup second byte: %w", err)
	}
	s.x[16] = uint64(byte(value))
	s.x[15] = uint64(uint32(s.x[16]) ^ uint32(s.x[15]))
	s.x[16] = 0x10d9
	s.x[15] = uint64(uint32(s.x[16]) + uint32(byte(s.x[15])))
	value, err = s.load(s.x[2]+uint64(uint32(s.x[15]))*8, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap table-setup third table: %w", err)
	}
	s.x[15] = value + 0xfffffffff1508aea
	if err := s.storeValue(s.sp+0x158, s.x[15], 8); err != nil {
		return 0, fmt.Errorf("bootstrap table-setup saved third table: %w", err)
	}
	s.x[13] = uint64(uint32(s.x[13]) ^ uint32(s.x[26]))
	s.x[13] = uint64(uint32(s.x[13]) ^ uint32(s.x[14]))
	s.x[13] = uint64(uint32(s.x[13]) * uint32(s.x[28]))
	value, err = s.load(s.x[2]+0x3a08, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap table-setup fourth table: %w", err)
	}
	s.x[14] = value + 0xfffffffffdc52ec2
	if err := s.storeValue(s.sp+0x1a0, s.x[14], 8); err != nil {
		return 0, fmt.Errorf("bootstrap table-setup saved fourth table: %w", err)
	}
	s.x[14] += uint64(byte(s.x[13]))
	value, err = s.load(s.x[14], 1)
	if err != nil {
		return 0, fmt.Errorf("bootstrap table-setup third byte: %w", err)
	}
	s.x[14] = uint64(byte(value)) ^ 0xb5
	value, err = s.load(s.x[2]+0x828, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap table-setup fifth table: %w", err)
	}
	s.x[15] = value + 0xffffffffba636d86
	if err := s.storeValue(s.sp+0x198, s.x[15], 8); err != nil {
		return 0, fmt.Errorf("bootstrap table-setup saved fifth table: %w", err)
	}
	value, err = s.load(s.x[15]+uint64(uint32(s.x[14])), 1)
	if err != nil {
		return 0, fmt.Errorf("bootstrap table-setup fourth byte: %w", err)
	}
	s.x[14] = uint64(byte(value))
	s.x[13] = uint64(uint32(s.x[14]) ^ uint32(s.x[13]))
	s.x[13] &= 0xff
	s.x[13] = uint64(uint32(s.x[13]) + 0xb98)
	value, err = s.load(s.x[2]+uint64(uint32(s.x[13]))*8, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap table-setup sixth table: %w", err)
	}
	s.x[13] = value + 0xffffffffaf39e8c5
	if err := s.storeValue(s.sp+0x150, s.x[13], 8); err != nil {
		return 0, fmt.Errorf("bootstrap table-setup saved sixth table: %w", err)
	}
	s.x[17] = 0x9227
	return next, nil
}
