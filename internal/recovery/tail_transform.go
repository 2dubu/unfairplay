package recovery

// tailTransform is the semantic form of the 140-byte current-enrollment tail
// stage.
func tailTransform(table []uint32, seed, tableIndex uint32) ([140]byte, bool) {
	var output [140]byte
	if len(table) != pairTransformTableWords || tableIndex >= pairTransformTableWords {
		return output, false
	}
	for index := range output {
		a := tableIndex
		b := (a + 42) % pairTransformTableWords
		c := (a + 83) % pairTransformTableWords
		output[index] = byte(seed ^ (table[a] >> 18) ^ table[b] ^ (table[c] >> 2) ^ 0xb8)
		tableIndex = (tableIndex + pairTransformTableWords - 1) % pairTransformTableWords
	}
	return output, true
}
