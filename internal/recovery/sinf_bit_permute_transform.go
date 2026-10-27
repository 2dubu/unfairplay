package recovery

import "fmt"

const (
	sinfBitPermuteStart = uint64(0x10002aba8)
	sinfBitPermuteTail  = uint64(0x10002ac00)
)

func runSINFBitPermute(s *state) (uint64, error) {
	if err := runSINFBitPermuteFront(s); err != nil {
		return 0, err
	}
	return runSINFBitPermuteFinish(s)
}

func runSINFBitPermuteFront(s *state) error {
	s.traceBlock(sinfBitPermuteStart)
	s.x[11] = uint64(uint32(s.x[1]) + 0x8c40)
	s.x[23] = uint64(uint32(s.x[11]) | 0x814)
	value, err := s.load(s.x[2], 8)
	if err != nil {
		return fmt.Errorf("SINF bit-permute first base: %w", err)
	}
	s.x[13] = value
	value, err = s.load(s.sp, 8)
	if err != nil {
		return fmt.Errorf("SINF bit-permute stack base: %w", err)
	}
	s.x[28] = value
	value, err = s.load(s.x[28]+2, 1)
	if err != nil {
		return fmt.Errorf("SINF bit-permute first byte: %w", err)
	}
	s.x[11] = uint64(byte(value))
	value, err = s.load(s.x[3], 8)
	if err != nil {
		return fmt.Errorf("SINF bit-permute second base: %w", err)
	}
	s.x[15] = value
	s.x[9] += s.x[15]
	value, err = s.load(s.x[9], 1)
	if err != nil {
		return fmt.Errorf("SINF bit-permute second byte: %w", err)
	}
	s.x[9] = uint64(byte(value))
	s.x[11] = uint64(uint32(s.x[9]) + uint32(s.x[11]))
	s.x[15] = 0x203
	s.x[17] = uint64(uint32(s.x[11]) * uint32(s.x[15]))
	s.x[17] = uint64(uint32(s.x[17]) >> 17)
	s.x[17] = uint64(uint32(s.x[17]) - uint32(s.x[17]<<8))
	s.x[11] = uint64(uint32(s.x[11]) + uint32(s.x[17]))
	s.x[11] = s.x[13] + uint64(uint16(s.x[11]))
	value, err = s.load(s.x[11], 1)
	if err != nil {
		return fmt.Errorf("SINF bit-permute third byte: %w", err)
	}
	s.x[11] = uint64(byte(value))
	s.x[17] = uint64(uint32(s.x[11]) & 0xffffffda)
	s.x[21] = 0xffffffa3
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[21]))
	return nil
}

func runSINFBitPermuteFinish(s *state) (uint64, error) {
	s.traceBlock(sinfBitPermuteTail)
	s.x[21] = uint64(uint32(s.x[23]) - 0x53)
	s.x[21] &= 0xff
	s.x[21] = uint64(uint32(s.x[11]) << (uint32(s.x[21]) & 31))
	s.x[22] = 0x28
	s.x[21] = uint64(uint32(s.x[21]) & uint32(s.x[22]))
	s.x[22] = 0x4e
	s.x[21] = uint64(uint32(s.x[21]) ^ uint32(s.x[22]))
	s.x[22] = 0x14
	s.x[26] = uint64(uint32(s.x[11]) ^ uint32(s.x[22]))
	s.x[21] = uint64(uint32(s.x[21]) & uint32(s.x[26]))
	s.x[21] = uint64(uint32(s.x[21]) ^ uint32(s.x[17]))
	for range 2 {
		s.x[21] = uint64(uint32(s.x[26]) & uint32(s.x[21]<<1))
		s.x[21] = uint64(uint32(s.x[21]) ^ uint32(s.x[17]))
	}
	s.x[21] = uint64(uint32(s.x[26]) & uint32(s.x[21]<<1))
	s.x[17] = uint64(uint32(s.x[21]) ^ uint32(s.x[17]))
	s.x[11] = uint64(uint32(s.x[11]) ^ uint32(s.x[17]<<1))
	s.x[11] &= 0xff
	s.x[17] = 0x635fa63a
	s.x[11] = uint64(uint32(s.x[11]) ^ uint32(s.x[17]))
	value, err := s.load(s.x[28]+3, 1)
	if err != nil {
		return 0, fmt.Errorf("SINF bit-permute fourth byte: %w", err)
	}
	s.x[17] = uint64(byte(value))
	s.x[9] = uint64(uint32(s.x[17]) + uint32(s.x[9]))
	s.x[15] = uint64(uint32(s.x[9]) * uint32(s.x[15]))
	s.x[15] = uint64(uint32(s.x[15]) >> 17)
	s.x[15] = uint64(uint32(s.x[15]) - uint32(s.x[15]<<8))
	s.x[9] = uint64(uint32(s.x[9]) + uint32(s.x[15]))
	s.x[9] = s.x[13] + uint64(uint16(s.x[9]))
	value, err = s.load(s.x[9], 1)
	if err != nil {
		return 0, fmt.Errorf("SINF bit-permute fifth byte: %w", err)
	}
	s.x[9] = uint64(byte(value))
	s.x[13] = uint64(uint32(s.x[9]) & 0xffffff94)
	s.x[13] = uint64(uint32(s.x[13]) | 0x60)
	s.x[15] = uint64(uint32(s.x[9]) & 0xffffffc0)
	s.x[15] = uint64(uint32(s.x[13]) ^ uint32(s.x[15]))
	s.x[17] = uint64(uint32(s.x[9]) ^ uint32(s.x[22]))
	for range 4 {
		s.x[13] = uint64(uint32(s.x[17]) & uint32(s.x[13]<<1))
		s.x[13] = uint64(uint32(s.x[13]) ^ uint32(s.x[15]))
	}
	s.x[13] = uint64(uint32(s.x[13] << 25))
	s.x[9] = uint64(uint32(s.x[13]) ^ uint32(s.x[9]<<24))
	s.x[13] = 0xb943636c
	s.x[26] = uint64(uint32(s.x[9]) ^ uint32(s.x[13]))
	next, err := loadBootstrapPointer(s, s.x[14], uint32(s.x[1]))
	if err != nil {
		return 0, fmt.Errorf("SINF bit-permute dispatch: %w", err)
	}
	s.x[9] = next
	s.x[17] = 0x1a30
	return next, nil
}
