package recovery

import "fmt"

const accountLargeBooleanTransform = uint64(0x10007a0b8)

// runAccountLargeBooleanTransform applies the SHA-1 feed-forward for A-D. E
// is handled immediately before this block by runAccountMixerSetup.
func runAccountLargeBooleanTransform(s *state) (uint64, error) {
	s.traceBlock(accountLargeBooleanTransform)
	if uint32(s.x[24]) != 0x892c {
		return 0, fmt.Errorf("account SHA-1 feed-forward selector %#x is unsupported", uint32(s.x[24]))
	}
	initial := accountSHA1State{
		A: uint32(s.x[26]) ^ accountSHA1StateMask,
		B: uint32(s.x[0]) ^ accountSHA1StateMask,
		C: uint32(s.x[22]) ^ accountSHA1StateMask,
		D: uint32(s.x[3]) ^ accountSHA1StateMask,
	}
	compressed := accountSHA1State{
		A: uint32(s.x[12]) ^ accountSHA1StateMask,
		B: uint32(s.x[16]) ^ accountSHA1StateMask,
		C: uint32(s.x[13]) ^ accountSHA1StateMask,
		D: uint32(s.x[11]) ^ accountSHA1StateMask,
	}
	result := accountSHA1FeedForward(initial, compressed)
	encodedA := result.A ^ accountSHA1StateMask
	encodedB := result.B ^ accountSHA1StateMask
	encodedC := result.C ^ accountSHA1StateMask
	encodedD := result.D ^ accountSHA1StateMask

	s.x[26], s.x[0], s.x[22], s.x[3] = uint64(encodedA), uint64(encodedB), uint64(encodedC), uint64(encodedD)
	if err := s.storeValue(addOffset(s.x[29], -0x94), uint64(encodedA), 4); err != nil {
		return 0, fmt.Errorf("account SHA-1 frame A: %w", err)
	}
	if err := s.storeValue(addOffset(s.x[29], -0x90), uint64(encodedB), 4); err != nil {
		return 0, fmt.Errorf("account SHA-1 frame B: %w", err)
	}
	if err := s.storeValue(addOffset(s.x[29], -0x8c), uint64(encodedC), 4); err != nil {
		return 0, fmt.Errorf("account SHA-1 frame C: %w", err)
	}
	if err := s.storeValue(addOffset(s.x[29], -0x88), uint64(encodedD), 4); err != nil {
		return 0, fmt.Errorf("account SHA-1 frame D: %w", err)
	}

	// Preserve the live native dispatcher outputs. Other circuit temporaries are
	// dead before the next semantic operation and deliberately are not modeled.
	s.x[9] = uint64(uint32(s.x[24]) + 0x52)
	s.x[7] = 0x7bb20000
	s.x[19] = 0x7bb2a446
	s.x[12] = 0x3d7f0910
	s.x[13] = uint64(encodedC ^ 0xef42ba1f)
	s.x[15] = 0xdd4e54af
	next, err := loadBootstrapPointer(s, s.x[21], uint32(s.x[24]))
	if err != nil {
		return 0, fmt.Errorf("account SHA-1 feed-forward target: %w", err)
	}
	s.x[11] = next
	s.x[17] = 0xaa2b
	return next, nil
}
