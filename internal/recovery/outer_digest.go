package recovery

const (
	outerDigestStop              = uint32(0x6247750e)
	outerDigestHardwareBias      = uint32(0xb03992b8)
	outerDigestCodecBias         = uint32(0x9fce6288)
	outerDigestMaximumIterations = 1024
)

// outerDigestState is the live semantic state carried by the protected digest
// recurrence. The selectors are the original x5 and x8 words; their opaque
// division sequences are simply modulo-six and modulo-thirteen rotations.
type outerDigestState struct {
	X0               uint32
	X4               uint32
	X10              uint32
	X27              uint32
	X28              uint32
	Counter          uint32
	HardwareSelector uint32
	CodecSelector    uint32
}

// applyOuterDigest evaluates the complete protected digest recurrence. The
// generated step below is a direct bit-vector circuit: it does not interpret
// lifted instructions or depend on a captured address space.
func applyOuterDigest(state *outerDigestState, hardware *[6]byte, codec *[17]byte) bool {
	if state == nil || hardware == nil || codec == nil {
		return false
	}
	hardwareBase := (state.HardwareSelector + outerDigestHardwareBias) % 6
	for range outerDigestMaximumIterations {
		codecPhase := (state.CodecSelector + outerDigestCodecBias) % 13
		exit := applyOuterDigestStep(state, hardware, codec, hardwareBase, codecPhase)
		state.CodecSelector--
		if exit {
			return true
		}
	}
	return false
}
