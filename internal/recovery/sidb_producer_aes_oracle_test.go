package recovery

var frameInverseSBox = func() (result [256]byte) {
	for input, output := range frameSBox {
		result[output] = byte(input)
	}
	return result
}()

func applySIDBCanonicalInverseRound(input [16]byte) (output [16]byte) {
	coefficients := [4][4]byte{
		{0x0e, 0x0b, 0x0d, 0x09},
		{0x09, 0x0e, 0x0b, 0x0d},
		{0x0d, 0x09, 0x0e, 0x0b},
		{0x0b, 0x0d, 0x09, 0x0e},
	}
	for column := range 4 {
		for row := range 4 {
			for sourceRow := range 4 {
				sourceColumn := (column + sourceRow) % 4
				source := sourceColumn*4 + sourceRow
				output[column*4+row] ^= frameGFMul(coefficients[row][sourceRow], input[source])
			}
		}
	}
	return output
}
