package recovery

import "fmt"

const (
	accountPublicKeySHA1Counter = uint32(0x481b90d7)
	accountPublicKeySHA1Done    = uint32(0x481b9127)
	accountWordRoundStart       = uint64(0x10010de18)
	accountDualWordRoundFront   = uint64(0x1001c2f34)
	accountDualWordRoundTail    = uint64(0x1001c3000)
	accountTripleWordRoundFront = uint64(0x1001c4ed0)
	accountTripleWordRoundTail  = uint64(0x1001c5000)
	accountFourthWordRound      = uint64(0x1002a759c)
)

func runAccountDualWordRound(s *state) (uint64, error) {
	return runAccountPublicKeySHA1Rounds(s)
}

func runAccountWordRound(_ *state) (uint64, error) {
	return 0, fmt.Errorf("unexpected account public-key SHA-1 phase-1 entry")
}

func runAccountTripleWordRound(_ *state) (uint64, error) {
	return 0, fmt.Errorf("unexpected account public-key SHA-1 phase-2 entry")
}

func runAccountFourthWordRound(_ *state) (uint64, error) {
	return 0, fmt.Errorf("unexpected account public-key SHA-1 phase-3 entry")
}

// runAccountPublicKeySHA1Rounds replaces the third flattened SHA-1 family,
// used for the DER-encoded RSA public keys in the certificate chain.
func runAccountPublicKeySHA1Rounds(s *state) (uint64, error) {
	if uint32(s.x[2]) != accountPublicKeySHA1Counter || uint32(s.x[23]) != accountSHA1StateMask {
		return 0, fmt.Errorf("account public-key SHA-1 state is unsupported")
	}
	scheduleBase := s.x[20] + s.x[9] + (s.x[2] << 2)
	var schedule [80]uint32
	for index := range schedule {
		word, err := s.load(scheduleBase+uint64(index*4), 4)
		if err != nil {
			return 0, fmt.Errorf("account public-key SHA-1 schedule word %d: %w", index, err)
		}
		schedule[index] = uint32(word) ^ accountSHA1ScheduleMask
	}
	initial := accountSHA1State{
		A: uint32(s.x[8]) ^ accountSHA1StateMask,
		B: uint32(s.x[12]) ^ accountSHA1StateMask,
		C: uint32(s.x[19]) ^ accountSHA1StateMask,
		D: uint32(s.x[10]) ^ accountSHA1StateMask,
		E: uint32(s.x[16]) ^ accountSHA1StateMask,
	}
	compressed := accountSHA1CompressSchedule(initial, schedule)
	encodedA := compressed.A ^ accountSHA1StateMask
	encodedB := compressed.B ^ accountSHA1StateMask
	encodedC := compressed.C ^ accountSHA1StateMask
	encodedD := compressed.D ^ accountSHA1StateMask
	encodedE := compressed.E ^ accountSHA1StateMask

	// Preserve the live ABI at the old final-round boundary. The following
	// native blocks apply feed-forward and continue certificate processing.
	s.x[0] = 0x7d72
	s.x[8] = uint64(encodedA)
	s.x[10] = uint64(encodedD)
	s.x[11] = uint64(accountPublicKeySHA1Done)
	s.x[12] = uint64(encodedB)
	s.x[15] = 0x8683
	s.x[16] = uint64(encodedE)
	s.x[17] = 0x7ab5
	s.x[19] = uint64(encodedC)
	s.x[22] = uint64(accountPublicKeySHA1Done)
	s.x[27] = 0x481b90ff
	s.compare64(s.x[22], uint64(accountPublicKeySHA1Done))
	next, err := loadBootstrapPointer(s, s.x[24], 0x8690)
	if err != nil {
		return 0, fmt.Errorf("account public-key SHA-1 target: %w", err)
	}
	s.x[2] = next
	return next, nil
}
