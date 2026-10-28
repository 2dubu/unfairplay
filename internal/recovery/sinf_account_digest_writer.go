package recovery

import "fmt"

// runSINFAccountDigestWriter serializes the semantic account SHA-1 digest used
// by the class-value stage. The native loop writes four protected bytes per
// visit; this typed stage writes the complete digest and retains its semantic
// value for structural candidate evaluation.
func runSINFAccountDigestWriter(s *state) (uint64, error) {
	if s == nil {
		return 0, fmt.Errorf("SINF account-digest state is unavailable")
	}
	if uint32(s.x[10]) != 0 {
		return 0, fmt.Errorf("SINF account-digest writer starts at %#x", uint32(s.x[10]))
	}
	limit, err := s.load(s.sp+0x1b0, 4)
	if err != nil {
		return 0, fmt.Errorf("SINF account-digest length: %w", err)
	}
	if uint32(limit) != 20 {
		return 0, fmt.Errorf("SINF account-digest length is %d, want 20", uint32(limit))
	}
	if err := writeAccountSHA1Digest(s, s.sp+0x19c, s.x[3]); err != nil {
		return 0, err
	}

	// At the completed loop boundary the protected dispatch predicate selects
	// x9 directly. Preserve only the live counter, target, and branch modifier;
	// the byte-lane temporaries are dead after this stage.
	s.x[10] = 20
	target, err := s.load(s.x[30]+shiftLeft(signExtend32(uint32(s.x[9])), 3, 64), 8)
	if err != nil {
		return 0, fmt.Errorf("SINF account-digest target: %w", err)
	}
	s.x[12] = target
	s.x[17] = 0x195a
	return target, nil
}
