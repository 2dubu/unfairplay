package recovery

import "fmt"

const (
	accountCertificateSHA1Counter = uint32(0x1cc84893)
	accountCertificateSHA1Done    = uint32(0x1cc848e3)
	accountRoundPermuteA          = uint64(0x100138d30)
	accountRoundPermuteB          = uint64(0x1001478d0)
	accountFoldRound              = uint64(0x10019f718)
	accountFoldTail               = uint64(0x10019f800)
	accountFoldLimit              = uint64(0x1cc848e3)
	accountSecondFoldRound        = uint64(0x10021d654)
	accountSecondFoldTail         = uint64(0x10021d800)
	accountSecondFoldLimit        = uint64(0x1cc848cf)
)

func runAccountRoundPermuteB(s *state) (uint64, error) {
	return runAccountCertificateSHA1Rounds(s)
}

func runAccountRoundPermuteA(_ *state) (uint64, error) {
	return 0, fmt.Errorf("unexpected account certificate SHA-1 phase-1 entry")
}

func runAccountSecondFold(_ *state) (uint64, error) {
	return 0, fmt.Errorf("unexpected account certificate SHA-1 phase-2 entry")
}

func runAccountFold(_ *state) (uint64, error) {
	return 0, fmt.Errorf("unexpected account certificate SHA-1 phase-3 entry")
}

// runAccountCertificateSHA1Rounds replaces a second flattened SHA-1 round
// family used while hashing the signed certificate bodies. Its schedule has
// already been prepared by the surrounding certificate parser.
func runAccountCertificateSHA1Rounds(s *state) (uint64, error) {
	if uint32(s.x[16]) != accountCertificateSHA1Counter || uint32(s.x[10]) != accountSHA1StateMask {
		return 0, fmt.Errorf("account certificate SHA-1 state is unsupported")
	}
	scheduleBase := s.x[8] + s.x[14] + (s.x[16] << 2)
	var schedule [80]uint32
	for index := range schedule {
		word, err := s.load(scheduleBase+uint64(index*4), 4)
		if err != nil {
			return 0, fmt.Errorf("account certificate SHA-1 schedule word %d: %w", index, err)
		}
		schedule[index] = uint32(word) ^ accountSHA1ScheduleMask
	}
	initial := accountSHA1State{
		A: uint32(s.x[11]) ^ accountSHA1StateMask,
		B: uint32(s.x[19]) ^ accountSHA1StateMask,
		C: uint32(s.x[15]) ^ accountSHA1StateMask,
		D: uint32(s.x[28]) ^ accountSHA1StateMask,
		E: uint32(s.x[6]) ^ accountSHA1StateMask,
	}
	compressed := accountSHA1CompressSchedule(initial, schedule)
	encodedA := compressed.A ^ accountSHA1StateMask
	encodedB := compressed.B ^ accountSHA1StateMask
	encodedC := compressed.C ^ accountSHA1StateMask
	encodedD := compressed.D ^ accountSHA1StateMask
	encodedE := compressed.E ^ accountSHA1StateMask

	// Preserve the live ABI at the old final-round boundary. The following
	// native blocks apply SHA-1 feed-forward and update the certificate parser.
	s.x[0] = 0x3f8d
	s.x[1], s.x[19] = uint64(encodedB), uint64(encodedB)
	s.x[2], s.x[16] = uint64(encodedE), uint64(encodedE)
	s.x[7] = uint64(accountCertificateSHA1Done)
	s.x[9] = 0x83c3
	s.x[11] = uint64(encodedA)
	s.x[12] = 0x0d990000
	s.x[13] = 0xa23f
	s.x[15] = uint64(encodedC)
	s.x[17] = 0x2583
	s.x[21] = uint64(accountSecondFoldLimit)
	s.x[28] = uint64(encodedD)
	s.compare64(s.x[7], uint64(accountCertificateSHA1Done))
	next, err := loadBootstrapPointer(s, s.x[20], 0xbc4e)
	if err != nil {
		return 0, fmt.Errorf("account certificate SHA-1 target: %w", err)
	}
	s.x[5] = next
	return next, nil
}
