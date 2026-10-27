package recovery

import "fmt"

const bootstrapByteCompareDispatch = uint64(0x1009ee0dc)

func runBootstrapByteCompareDispatch(s *state) (uint64, error) {
	s.traceBlock(bootstrapByteCompareDispatch)
	s.x[12] = 0
	s.x[8] = uint64(uint32(s.x[11]) + 0x0f6b5603)
	s.x[13] = uint64(uint32(s.x[8]) & 0xf094aecf)
	s.x[8] = uint64(uint32(s.x[11]) ^ 0x7b19)
	s.x[14] = 0x8f
	s.x[8] = uint64(uint32(s.x[8]) * uint32(s.x[14]))
	s.x[14] = uint64(uint32(shiftRight(s.x[9], int(s.x[10]), 32)))
	s.x[15] = s.x[5] - 1
	first, err := s.load(s.sp+0x150, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap byte-compare first source: %w", err)
	}
	second, err := s.load(s.sp+0x158, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap byte-compare second source: %w", err)
	}
	s.x[16], s.x[17] = first, second
	word, err := s.load(s.x[16], 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap byte-compare signed word: %w", err)
	}
	s.x[16] = uint64(int64(int32(word)))
	value, err := s.load(s.x[17], 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap byte-compare source base: %w", err)
	}
	s.x[17] = value
	s.x[0] = 0xfffffffffa6e3c78
	s.x[16] &= s.x[0]
	s.x[16] += s.x[17]
	word, err = s.load(s.x[16], 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap byte-compare source word: %w", err)
	}
	s.x[16] = uint64(uint32(word))
	s.x[15] += s.x[1]
	s.x[16] = uint64(uint32(s.x[16]) + uint32(s.x[15]))
	s.x[17] = uint64(uint32(s.x[8]) ^ 0xfa6e4244)
	s.x[0] = 0x101086000
	value, err = s.load(s.x[0]+0x138, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap byte-compare first mask: %w", err)
	}
	s.x[0] = value
	s.x[27] = 0x12af55048709d7d1
	s.x[13] += s.x[27]
	s.x[13] ^= s.x[0]
	s.x[0] = 0x101084000 + 0x6c0
	value, err = s.load(s.x[0], 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap byte-compare second mask: %w", err)
	}
	s.x[0] = value ^ 0xa7a730ec
	s.x[17] = uint64(uint32(s.x[16]) + uint32(s.x[17]))
	s.x[17] &= s.x[0]
	s.x[17] += s.x[13]
	value, err = s.load(s.x[17], 1)
	if err != nil {
		return 0, fmt.Errorf("bootstrap byte-compare first byte: %w", err)
	}
	s.x[17] = uint64(byte(value))
	s.x[27] = 0xfa6e3c7c
	s.x[16] = uint64(uint32(s.x[16]) + uint32(s.x[27]))
	s.x[16] &= s.x[0]
	s.x[0] = uint64(uint32(s.x[15]) * 0x23)
	s.x[13] += s.x[16]
	value, err = s.load(s.x[13], 1)
	if err != nil {
		return 0, fmt.Errorf("bootstrap byte-compare second byte: %w", err)
	}
	s.x[13] = uint64(byte(value))
	s.x[14] = uint64(uint32(s.x[14]) ^ uint32(s.x[0]))
	s.x[14] = uint64(uint32(s.x[14]) ^ uint32(s.x[17]))
	s.x[13] = uint64(uint32(s.x[14]) ^ uint32(s.x[13]))
	s.x[14] = uint64(uint32(s.x[13]) & 0xff)
	s.x[16] = 0x59
	s.x[13] = uint64(uint32(s.x[13]) ^ uint32(s.x[16]))
	if err := s.storeValue(s.x[15], uint64(byte(s.x[13])), 1); err != nil {
		return 0, fmt.Errorf("bootstrap byte-compare output: %w", err)
	}
	s.compare32(uint32(s.x[14]), 0x59)
	different := uint32(boolValue(s.condition("ne")))
	s.x[13] = uint64(different)
	s.x[5] = s.subtractWithFlags(uint64(uint32(s.x[5])), uint64(different), 32)
	s.x[10] = uint64(uint32(s.x[10]) + 8)
	done := uint32(boolValue(s.condition("eq")))
	s.x[13] = uint64(done * 0x7dc)
	s.x[11] = uint64(uint32(s.x[13]) ^ uint32(s.x[11]))
	next, err := loadBootstrapPointer(s, s.x[19], uint32(s.x[11]))
	if err != nil {
		return 0, fmt.Errorf("bootstrap byte-compare target: %w", err)
	}
	s.x[11] = next
	s.x[17] = 0x2b23
	return next, nil
}
