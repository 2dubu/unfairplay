package recovery

const outerOptionalPhaseOffset = uint32(70)

// outerOptionalFirstIndex selects the first optional-hook table word from an
// arbitrary seed. Addition wraps at 32 bits before the modulo, matching the
// native recurrence.
func outerOptionalFirstIndex(seed uint32) uint32 {
	phase := (seed + 0xf587d4ee) % outerWriterSeedStates
	return (phase + 114) % outerWriterSeedStates
}
