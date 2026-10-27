package recovery

import "fmt"

const bootstrapWordPack = uint64(0x100933428)

func runBootstrapWordPack(s *state) (uint64, error) {
	s.traceBlock(bootstrapWordPack)
	s.x[13] = uint64(uint32(s.x[6]) ^ uint32(s.x[12]))
	s.x[14] = uint64(uint32(s.x[11]) * uint32(s.x[7]))
	s.x[16] = uint64(uint32(s.x[14]) & 0xfffc)
	s.x[16] = uint64(uint32(s.x[16]) * uint32(s.x[5]))
	s.x[16] = uint64(uint32(s.x[16]) >> 23)
	s.x[16] = uint64(uint32(s.x[14]) - uint32(s.x[16])*uint32(s.x[23]))
	s.x[16] &= 0xfffe
	s.x[17] = uint64(uint32(s.x[13]) + uint32(s.x[25]))
	s.x[17] = uint64(uint32(s.x[17]) + uint32(s.x[16]))
	s.compare32(uint32(s.x[16]), 0x49)
	if !s.condition("lo") {
		s.x[16] = uint64(uint32(s.x[17]))
	}
	s.x[17] = uint64(uint32(s.x[11])*uint32(s.x[19]) + uint32(s.x[20]))
	s.x[11] = uint64(uint32(s.x[11]) + 1)
	s.x[0] = uint64(uint32(s.x[14]) + 0x15)
	s.x[3] = uint64(uint32(s.x[0]) & 0xffff)
	s.x[3] = uint64(uint32(s.x[3]) * uint32(s.x[5]))
	s.x[3] = uint64(uint32(s.x[3]) >> 22)
	s.x[0] = uint64(uint32(s.x[0]) - uint32(s.x[3])*uint32(s.x[15]))
	s.x[0] = s.x[8] + uint64(uint16(s.x[0]))
	value, err := s.load(s.x[0], 1)
	if err != nil {
		return 0, fmt.Errorf("bootstrap word-pack first byte: %w", err)
	}
	s.x[0] = uint64(byte(value))
	value, err = s.load(s.x[8]+uint64(uint32(s.x[16])), 1)
	if err != nil {
		return 0, fmt.Errorf("bootstrap word-pack second byte: %w", err)
	}
	s.x[16] = uint64(byte(value))
	s.x[3] = uint64(uint32(s.x[14]) + 0x2a)
	s.x[4] = uint64(uint32(s.x[3]) & 0xfffe)
	s.x[4] = uint64(uint32(s.x[4]) * uint32(s.x[5]))
	s.x[4] = uint64(uint32(s.x[4]) >> 22)
	s.x[3] = uint64(uint32(s.x[3]) - uint32(s.x[4])*uint32(s.x[15]))
	s.x[3] = s.x[8] + uint64(uint16(s.x[3]))
	value, err = s.load(s.x[3], 1)
	if err != nil {
		return 0, fmt.Errorf("bootstrap word-pack third byte: %w", err)
	}
	s.x[3] = uint64(byte(value))
	s.x[16] = uint64(uint32(s.x[16]) << 16)
	s.x[16] = uint64(uint32(s.x[16]) | uint32(s.x[0]<<8))
	s.x[14] = uint64(uint32(s.x[14]) + 0x3f)
	s.x[0] = uint64(uint32(s.x[14]) & 0xffff)
	s.x[0] = uint64(uint32(s.x[0]) * uint32(s.x[5]))
	s.x[0] = uint64(uint32(s.x[0]) >> 22)
	s.x[14] = uint64(uint32(s.x[14]) - uint32(s.x[0])*uint32(s.x[15]))
	s.x[16] = uint64(uint32(s.x[16]) | uint32(s.x[3]))
	s.x[14] = s.x[8] + uint64(uint16(s.x[14]))
	value, err = s.load(s.x[14], 1)
	if err != nil {
		return 0, fmt.Errorf("bootstrap word-pack fourth byte: %w", err)
	}
	s.x[14] = uint64(byte(value))
	s.x[14] = uint64(uint32(s.x[14]) | uint32(s.x[16]<<8))
	s.x[16] = uint64(uint32(s.x[17]) & 0xffff)
	s.x[16] = uint64(uint32(s.x[16]) * uint32(s.x[21]))
	s.x[16] = uint64(uint32(s.x[16]) >> 23)
	s.x[16] = uint64(uint32(s.x[17]) - uint32(s.x[16])*uint32(s.x[22]))
	s.x[16] = (s.x[16] & 0xffff) << 2
	if err := s.storeValue(s.x[9]+s.x[16], uint64(uint32(s.x[14])), 4); err != nil {
		return 0, fmt.Errorf("bootstrap word-pack output: %w", err)
	}
	s.compare32(uint32(s.x[11]), 0x10)
	gate := uint32(boolValue(s.condition("eq")))
	s.x[14] = uint64(gate * uint32(s.x[12]))
	s.x[14] = uint64(uint32(s.x[14]) ^ uint32(s.x[6]))
	next, err := loadBootstrapPointer(s, s.x[30], uint32(s.x[14]))
	if err != nil {
		return 0, fmt.Errorf("bootstrap word-pack target: %w", err)
	}
	s.x[14] = next
	s.x[16] = 0x10
	s.x[17] = 0xf00c
	return next, nil
}
