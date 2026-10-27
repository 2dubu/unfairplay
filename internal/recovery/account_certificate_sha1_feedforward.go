package recovery

import "fmt"

const (
	accountTripleBooleanTransform = uint64(0x1000912b0)
	accountStateSetupTransform    = uint64(0x1001a8524)
)

func runAccountStateSetupTransform(s *state) (uint64, error) {
	return runAccountCertificateSHA1FeedForward(s)
}

func runAccountTripleBooleanTransform(_ *state) (uint64, error) {
	return 0, fmt.Errorf("unexpected account certificate SHA-1 feed-forward continuation")
}

// runAccountCertificateSHA1FeedForward replaces the native feed-forward
// circuits for SHA-1 over signed certificate bodies.
func runAccountCertificateSHA1FeedForward(s *state) (uint64, error) {
	var encodedInitial [5]uint32
	for index := range encodedInitial {
		value, err := s.load(s.sp+0xe4+uint64(index*4), 4)
		if err != nil {
			return 0, fmt.Errorf("account certificate SHA-1 initial word %d: %w", index, err)
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
		A: uint32(s.x[11]) ^ accountSHA1StateMask,
		B: uint32(s.x[1]) ^ accountSHA1StateMask,
		C: uint32(s.x[15]) ^ accountSHA1StateMask,
		D: uint32(s.x[28]) ^ accountSHA1StateMask,
		E: uint32(s.x[16]) ^ accountSHA1StateMask,
	}
	result := accountSHA1FeedForward(initial, compressed)
	encodedResult := [...]uint32{
		result.A ^ accountSHA1StateMask, result.B ^ accountSHA1StateMask,
		result.C ^ accountSHA1StateMask, result.D ^ accountSHA1StateMask,
		result.E ^ accountSHA1StateMask,
	}
	for index, word := range encodedResult {
		if err := s.storeValue(s.sp+0xe4+uint64(index*4), uint64(word), 4); err != nil {
			return 0, fmt.Errorf("account certificate SHA-1 updated word %d: %w", index, err)
		}
	}

	framePointer, err := s.load(s.sp+0x80, 8)
	if err != nil {
		return 0, fmt.Errorf("account certificate SHA-1 frame pointer: %w", err)
	}
	s.x[26] = framePointer - 0x140
	if err := s.storeValue(s.sp+0x80, s.x[26], 8); err != nil {
		return 0, fmt.Errorf("account certificate SHA-1 updated frame pointer: %w", err)
	}
	first, err := s.load(s.sp+40, 8)
	if err != nil {
		return 0, fmt.Errorf("account certificate SHA-1 saved word 0: %w", err)
	}
	second, err := s.load(s.sp+48, 8)
	if err != nil {
		return 0, fmt.Errorf("account certificate SHA-1 saved word 1: %w", err)
	}
	next, err := s.load(s.sp+0x40, 8)
	if err != nil {
		return 0, fmt.Errorf("account certificate SHA-1 continuation: %w", err)
	}
	frameWord, err := s.load(s.sp+0x14, 4)
	if err != nil {
		return 0, fmt.Errorf("account certificate SHA-1 frame word: %w", err)
	}
	s.x[8] = next
	s.x[11] = second
	s.x[13] = 0x41225b35
	s.x[14] = uint64(uint32(frameWord))
	s.x[17] = 0x1c20
	s.x[21] = 0x56913c50
	s.x[27] = first
	return next, nil
}
