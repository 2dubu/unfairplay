package recovery

const pairTransformTableWords = 167

// pairTransform26A428 is the semantic form of the paired even/odd hooks for
// the proven 26A428 runtime domain. The hooks fill alternating bytes of one
// 100-byte buffer from a 167-word table.
func pairTransform26A428(table []uint32, evenSeed, oddSeed uint32, evenIndex, oddIndex uint32) ([100]byte, bool) {
	var output [100]byte
	if len(table) != pairTransformTableWords || evenIndex >= pairTransformTableWords || oddIndex >= pairTransformTableWords {
		return output, false
	}

	for index := range 50 {
		a := evenIndex
		b := (a + 42) % pairTransformTableWords
		c := (a + 41) % pairTransformTableWords
		output[index*2] = byte(evenSeed ^ (table[a] >> 9) ^ table[b] ^ (table[c] >> 23) ^ 0x6d)
		evenIndex = c

		a = oddIndex
		b = (a + 42) % pairTransformTableWords
		c = (a + 82) % pairTransformTableWords
		output[index*2+1] = byte(oddSeed ^ (table[a] >> 19) ^ table[b] ^ (table[c] >> 13) ^ 0x87)
		oddIndex = (2*a + 41) % pairTransformTableWords
	}
	return output, true
}
