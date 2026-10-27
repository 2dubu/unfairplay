package recovery

// writerTransform26A428 is the semantic form of the 140-byte writer hook for
// the proven 26A428 runtime domain. controls contains the per-byte words
// supplied by the enclosing outer stage.
func writerTransform26A428(table, controls []uint32, seed, tableIndex uint32) ([140]byte, bool) {
	var output [140]byte
	if len(table) != pairTransformTableWords || len(controls) != len(output) || tableIndex >= pairTransformTableWords {
		return output, false
	}

	for index, control := range controls {
		value := control >> 11
		encoded := value ^ 0x00058840
		value = encoded - (((encoded & 0x000208de) ^ (value & 2)) << 1)
		value += 0xc7c208dc
		value ^= seed ^ 0x449 ^ table[tableIndex]>>19
		value -= 0x76 & (value << 1)
		value += 0x3b
		value ^= 0x27 ^ table[(tableIndex+42)%pairTransformTableWords] ^
			table[(tableIndex+82)%pairTransformTableWords]>>17
		value -= 0x36 & (value << 1)
		output[index] = byte(value + 0x1b)
		tableIndex = (tableIndex + 36) % pairTransformTableWords
	}
	return output, true
}
