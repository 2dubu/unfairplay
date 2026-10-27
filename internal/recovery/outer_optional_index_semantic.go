package recovery

const outerOptionalPhaseOffset = uint32(70)

// outerOptionalFirstIndex selects the first optional-hook table word from an
// arbitrary seed. Addition wraps at 32 bits before the modulo, matching the
// native recurrence.
func outerOptionalFirstIndex(seed uint32) uint32 {
	phase := (seed + 0xf587d4ee) % outerWriterSeedStates
	return (phase + 114) % outerWriterSeedStates
}

func outerOptionalDecode(word uint32) uint32 {
	positive := int64(word & 0xddb15cc4)
	negative := int64(word &^ 0xddb15cc4)
	remainder := -(positive - negative) % 300
	if remainder < 0 {
		remainder += 300
	}

	centered := remainder
	if centered >= 140 {
		centered -= 300
	}
	value := byte(centered + 4)
	low := (value - 2*(value&0xf0)) ^ 0x02

	high := uint32(0)
	if remainder >= 140 && remainder < 184 {
		high = 1 << 8
	}
	return 0x976dec00 | high | uint32(low)
}
