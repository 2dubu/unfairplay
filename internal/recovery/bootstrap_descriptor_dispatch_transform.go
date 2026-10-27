package recovery

import "fmt"

const bootstrapDescriptorDispatch = uint64(0x100382e30)

func runBootstrapDescriptorDispatch(s *state) (uint64, error) {
	s.traceBlock(bootstrapDescriptorDispatch)
	s.sp -= 0x10
	s.x[8] = 0x101074000
	value, err := s.load(s.x[8]+0x5e0, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap descriptor-dispatch first state: %w", err)
	}
	s.x[9] = value
	s.x[10] = bootstrapStateBase
	value, err = s.load(s.x[10]+0xb40, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap descriptor-dispatch second state: %w", err)
	}
	s.x[11] = value
	s.x[9] ^= value
	s.x[11] = s.sp + 8
	s.x[12] = 0x0a5cd4cf
	s.x[13] = 0x98c1c57ff8a0fd3a
	s.x[14] = 0x73d046208eaed1a6
	s.x[9] += s.x[11]
	s.x[11] = s.x[9] ^ s.x[14]
	s.x[13] = s.x[9]*s.x[12] + s.x[13]
	s.x[14] = s.x[11] * s.x[12]
	if err := s.storeValue(s.x[8]+0x5e0, s.x[13], 8); err != nil {
		return 0, fmt.Errorf("bootstrap descriptor-dispatch first state output: %w", err)
	}
	if err := s.storeValue(s.x[10]+0xb40, s.x[14], 8); err != nil {
		return 0, fmt.Errorf("bootstrap descriptor-dispatch second state output: %w", err)
	}
	s.x[8] = uint64(uint32(s.x[0]) & 0x6c3f6207)
	s.x[8] = uint64(uint32(s.x[8]) << 1)
	s.x[9] = 0x93c09df8
	s.x[8] = uint64(uint32(s.x[8]) - uint32(s.x[0]))
	s.x[8] = uint64(uint32(s.x[8]) + uint32(s.x[9]))
	s.x[9] = 0x6a14048c
	s.x[8] = uint64(uint32(s.x[8]) ^ uint32(s.x[9]))
	s.x[9] = 0xda51a827
	s.x[9] = uint64(uint32(s.x[8]) * uint32(s.x[9]))
	word, err := s.load(s.x[0]+4, 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap descriptor-dispatch first word: %w", err)
	}
	s.x[8] = uint64(uint32(word))
	s.x[15] = uint64(uint32(s.x[8]) - uint32(s.x[9]))
	value, err = s.load(s.x[0]+8, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap descriptor-dispatch nested descriptor: %w", err)
	}
	s.x[8] = value
	word, err = s.load(s.x[0]+0x10, 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap descriptor-dispatch second word: %w", err)
	}
	s.x[10] = uint64(uint32(word))
	s.x[9] = uint64(uint32(s.x[10]) ^ uint32(s.x[9]))
	s.x[10] = 0xf0bfee9a
	s.x[10] = s.addWithCarry(uint64(uint32(s.x[9])), uint64(uint32(s.x[10])), false, 32)
	s.x[11] = 0x0f401166
	s.x[9] = uint64(uint32(s.x[11]) - uint32(s.x[9]))
	if s.condition("pl") {
		s.x[9] = uint64(uint32(s.x[10]))
	}
	s.x[11] = uint64(uint32(s.x[9]) >> 5)
	word, err = s.load(s.x[8]+4, 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap descriptor-dispatch nested word: %w", err)
	}
	s.x[10] = uint64(uint32(word))
	s.x[16] = 0x842a8d93
	s.x[16] = s.addWithCarry(uint64(uint32(s.x[10])), uint64(uint32(s.x[16])), false, 32)
	s.x[17] = 0x7bd5726d
	s.x[10] = uint64(uint32(s.x[17]) - uint32(s.x[10]))
	if s.condition("pl") {
		s.x[10] = uint64(uint32(s.x[16]))
	}
	s.compare32(uint32(s.x[11]), uint32(s.x[10]))
	gate := uint32(boolValue(s.condition("lt")))
	s.x[16] = uint64(gate)
	s.x[10] = uint64(uint32(s.x[15]) + 0x32)
	s.x[16] = uint64(uint32(s.x[10]) ^ gate)
	s.x[17] = 1
	s.x[16] = uint64(uint32(s.x[17]) &^ uint32(s.x[16]))
	s.x[17] = 0x5e
	s.x[16] = uint64(uint32(s.x[16]) * uint32(s.x[17]))
	s.x[17] = 0xa6
	s.x[13] = uint64(uint32(s.x[14]) + uint32(s.x[13]))
	s.x[13] = uint64(uint32(s.x[13]) ^ uint32(s.x[17]))
	s.x[12] = uint64(uint32(s.x[13]) * uint32(s.x[12]))
	s.x[13] = 0x100fe9390 + uint64(byte(s.x[12]))
	value, err = s.load(s.x[13], 1)
	if err != nil {
		return 0, fmt.Errorf("bootstrap descriptor-dispatch first byte: %w", err)
	}
	s.x[13] = uint64(byte(value)) ^ 0x28
	s.x[14] = 0x100de2af0
	value, err = s.load(s.x[14]+uint64(uint32(s.x[13])), 1)
	if err != nil {
		return 0, fmt.Errorf("bootstrap descriptor-dispatch second byte: %w", err)
	}
	s.x[13] = uint64(byte(value))
	s.x[12] = uint64(uint32(s.x[12]) ^ uint32(s.x[13]))
	s.x[12] &= 0xff
	s.x[12] = uint64(uint32(s.x[12]) + 0x1b6)
	s.x[13] = bootstrapTargetTable
	value, err = s.load(s.x[13]+uint64(uint32(s.x[12]))*8, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap descriptor-dispatch table: %w", err)
	}
	s.x[12] = value
	s.x[14] = uint64(uint32(s.x[16]) ^ uint32(s.x[15]))
	s.x[13] = s.x[12] + 0xffffffff98245355
	next, err := loadBootstrapPointer(s, s.x[13], uint32(s.x[14]))
	if err != nil {
		return 0, fmt.Errorf("bootstrap descriptor-dispatch target: %w", err)
	}
	s.x[14] = next
	s.x[12] = 0x5d55faf3
	s.x[17] = 0xe952
	return next, nil
}
