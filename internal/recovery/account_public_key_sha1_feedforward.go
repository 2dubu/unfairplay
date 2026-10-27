package recovery

import "fmt"

func runAccountDigestStateSetup(s *state) (uint64, error) {
	return runAccountPublicKeySHA1FeedForward(s)
}

// runAccountPublicKeySHA1FeedForward replaces the five native feed-forward
// circuits following runAccountPublicKeySHA1Rounds.
func runAccountPublicKeySHA1FeedForward(s *state) (uint64, error) {
	if uint32(s.x[0]) != 0x7d72 || uint32(s.x[23]) != accountSHA1StateMask {
		return 0, fmt.Errorf("account public-key SHA-1 feed-forward state is unsupported")
	}
	var encodedInitial [5]uint32
	for index := range encodedInitial {
		value, err := s.load(s.sp+740+uint64(index*4), 4)
		if err != nil {
			return 0, fmt.Errorf("account public-key SHA-1 initial word %d: %w", index, err)
		}
		encodedInitial[index] = uint32(value)
	}
	initial := accountSHA1State{
		A: encodedInitial[0] ^ accountSHA1StateMask,
		B: encodedInitial[1] ^ accountSHA1StateMask,
		C: encodedInitial[2] ^ accountSHA1StateMask,
		D: encodedInitial[3] ^ accountSHA1StateMask,
		E: encodedInitial[4] ^ accountSHA1StateMask,
	}
	compressed := accountSHA1State{
		A: uint32(s.x[8]) ^ accountSHA1StateMask,
		B: uint32(s.x[12]) ^ accountSHA1StateMask,
		C: uint32(s.x[19]) ^ accountSHA1StateMask,
		D: uint32(s.x[10]) ^ accountSHA1StateMask,
		E: uint32(s.x[16]) ^ accountSHA1StateMask,
	}
	result := accountSHA1FeedForward(initial, compressed)
	encodedResult := [...]uint32{
		result.A ^ accountSHA1StateMask, result.B ^ accountSHA1StateMask,
		result.C ^ accountSHA1StateMask, result.D ^ accountSHA1StateMask,
		result.E ^ accountSHA1StateMask,
	}
	for index, word := range encodedResult {
		if err := s.storeValue(s.sp+740+uint64(index*4), uint64(word), 4); err != nil {
			return 0, fmt.Errorf("account public-key SHA-1 updated word %d: %w", index, err)
		}
	}

	framePointer, err := s.load(s.sp+144, 8)
	if err != nil {
		return 0, fmt.Errorf("account public-key SHA-1 frame pointer: %w", err)
	}
	s.x[25] = framePointer - 0x140
	if err := s.storeValue(s.sp+144, s.x[25], 8); err != nil {
		return 0, fmt.Errorf("account public-key SHA-1 updated frame pointer: %w", err)
	}
	first, err := s.load(s.sp+72, 8)
	if err != nil {
		return 0, fmt.Errorf("account public-key SHA-1 saved word 0: %w", err)
	}
	second, err := s.load(s.sp+80, 8)
	if err != nil {
		return 0, fmt.Errorf("account public-key SHA-1 saved word 1: %w", err)
	}
	next, err := s.load(s.sp+96, 8)
	if err != nil {
		return 0, fmt.Errorf("account public-key SHA-1 continuation: %w", err)
	}
	saved, err := s.load(s.sp+56, 8)
	if err != nil {
		return 0, fmt.Errorf("account public-key SHA-1 saved continuation: %w", err)
	}
	s.x[8], s.x[26], s.x[21] = first, second, saved
	s.x[10] = 0x4696e82b
	s.x[11] = next
	s.x[12] = 0x1243156a
	s.x[13] = 0xffeff9ff
	s.x[15] = 0x094e5100
	s.x[17] = 0x377c
	s.x[19] = 0xf530238f
	return next, nil
}
