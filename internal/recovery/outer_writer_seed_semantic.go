package recovery

const (
	outerWriterSeedBase   = uint32(0x0a782f77)
	outerWriterSeedStates = uint32(167)
	outerWriterSeedStride = uint32(36)
)

// outerWriterSeedStep normalizes an arbitrary seed into the 167-state writer
// domain and advances it. Addition wraps at 32 bits before the modulo, matching
// the native recurrence.
func outerWriterSeedStep(seed uint32) uint32 {
	phase := (seed + 0xf587d953) % outerWriterSeedStates
	return outerWriterSeedBase + (phase+124)%outerWriterSeedStates
}
