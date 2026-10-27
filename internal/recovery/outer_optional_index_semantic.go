package recovery

const outerOptionalPhaseOffset = uint32(70)

// outerOptionalFirstIndex selects the first optional-hook table word from the
// normalized 167-state writer seed. Only the first captured seed can be
// outside this domain; it retains the reference normalizer for now.
func outerOptionalFirstIndex(seed uint32) uint32 {
	if seed >= outerWriterSeedBase {
		phase := seed - outerWriterSeedBase
		if phase < outerWriterSeedStates {
			return (phase + outerOptionalPhaseOffset) % outerWriterSeedStates
		}
	}
	return outerOptionalFirstIndexMBA(seed)
}
