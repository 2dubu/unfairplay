package recovery

// writerTransform is the semantic form of the 140-byte current-enrollment
// writer stage. controls contains the per-byte words
// supplied by the enclosing outer stage.
//
// The protected implementation surrounds each XOR with three byte-local MBA
// networks. Over their complete 8-bit domains those networks are exactly XOR
// 0x9c, XOR 0x3b, and XOR 0x1b. Folding them with the surrounding XOR 0x49 and
// XOR 0x27 constants leaves the single 0xd2 constant below.
func writerTransform(table, controls []uint32, seed, tableIndex uint32) ([140]byte, bool) {
	var output [140]byte
	if len(table) != pairTransformTableWords || len(controls) != len(output) || tableIndex >= pairTransformTableWords {
		return output, false
	}

	for index, control := range controls {
		output[index] = byte(control>>11) ^ byte(seed) ^
			byte(table[tableIndex]>>19) ^
			byte(table[(tableIndex+42)%pairTransformTableWords]) ^
			byte(table[(tableIndex+82)%pairTransformTableWords]>>17) ^ 0xd2
		tableIndex = (tableIndex + 36) % pairTransformTableWords
	}
	return output, true
}
