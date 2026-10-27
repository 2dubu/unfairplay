package recovery

const outerDigestCounterBias = uint32(0xf587d089)

// outerDigestInitialCounter derives the seven-state digest counter used by the
// natural outer path. The native implementation computes a 32-bit MBA hash,
// but only hash mod 7 is observable. A QF_BV proof over every 32-bit seed and
// each accepted control state reduces that value to the wrapped sum below.
func outerDigestInitialCounter(control, seed uint32) (uint32, bool) {
	if control != 0x60319d83 && control != 0x60319e4b && control != 0x60319d87 {
		return 0, false
	}
	base := uint32(0x62477516)
	if control == 0x60319d87 {
		base += 4
	}
	return base + uint32(seed+outerDigestCounterBias)%7, true
}
