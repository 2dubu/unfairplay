package recovery

const (
	outerWriterSeedBase   = uint32(0x0a782f77)
	outerWriterSeedStates = uint32(167)
	outerWriterSeedStride = uint32(36)
)

// outerWriterSeedStep advances the normalized writer/table index state. The
// natural outer path stays inside this 167-value domain after its first step.
// An arbitrary captured seed still uses the reference normalizer until that
// one-time mapping is recovered independently.
func outerWriterSeedStep(seed uint32) uint32 {
	if seed >= outerWriterSeedBase {
		phase := seed - outerWriterSeedBase
		if phase < outerWriterSeedStates {
			return outerWriterSeedBase + (phase+outerWriterSeedStride)%outerWriterSeedStates
		}
	}
	return outerWriterSeedStepMBA(seed)
}
