package recovery

import "fmt"

const accountMixerSetup = uint64(0x100062280)

// runAccountMixerSetup applies the SHA-1 feed-forward for E. The remaining
// four words are applied by runAccountLargeBooleanTransform; the split is an
// artifact of the native dispatcher.
func runAccountMixerSetup(s *state) (uint64, error) {
	s.traceBlock(accountMixerSetup)
	encodedInitialE, err := s.load(s.sp+0x7c, 4)
	if err != nil {
		return 0, fmt.Errorf("account SHA-1 initial E: %w", err)
	}
	if uint32(s.x[7]) != 0x82835ef1 {
		return 0, fmt.Errorf("account SHA-1 feed-forward marker %#x is unsupported", uint32(s.x[7]))
	}
	initialE := uint32(encodedInitialE) ^ accountSHA1StateMask
	compressedE := uint32(s.x[9]) ^ accountSHA1StateMask
	encodedResult := (initialE + compressedE) ^ accountSHA1StateMask
	if err := s.storeValue(s.sp+0x7c, uint64(encodedResult), 4); err != nil {
		return 0, fmt.Errorf("account SHA-1 updated E: %w", err)
	}
	if err := s.storeValue(addOffset(s.x[29], -0x84), uint64(encodedResult), 4); err != nil {
		return 0, fmt.Errorf("account SHA-1 frame E: %w", err)
	}

	// Preserve the native dispatcher ABI until the whole compression loop is
	// replaced. These values are dispatcher scratch, not SHA-1 state.
	s.x[2] = 0x7bf7ff7f
	s.x[4] = 0x0554d5ca
	s.x[15] = 0x8283d7dd
	s.x[24] = uint64(uint32(s.x[7]) ^ uint32(s.x[15]))
	s.compare32(uint32(s.x[7]), 0x62e0114c)
	dispatchIndex := uint32(s.x[7]) ^ 0x82832578
	if s.condition("eq") {
		dispatchIndex ^= 0x8ca4
	}
	next, err := loadBootstrapPointer(s, s.x[21], dispatchIndex)
	if err != nil {
		return 0, fmt.Errorf("account mixer-setup target: %w", err)
	}
	s.x[9] = next
	s.x[17] = 0x169e
	return next, nil
}
