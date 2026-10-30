package recovery

var frameSBox = func() (result [256]byte) {
	for input := range 256 {
		inverse := frameGFInverse(byte(input))
		rotated, value := inverse, inverse
		for range 4 {
			rotated = rotated<<1 | rotated>>7
			value ^= rotated
		}
		result[input] = value ^ 0x63
	}
	return result
}()

func frameGFMul(left, right byte) byte {
	var result byte
	for right != 0 {
		if right&1 != 0 {
			result ^= left
		}
		high := left & 0x80
		left <<= 1
		if high != 0 {
			left ^= 0x1b
		}
		right >>= 1
	}
	return result
}

func frameGFInverse(value byte) byte {
	if value == 0 {
		return 0
	}
	result, base, exponent := byte(1), value, 254
	for exponent != 0 {
		if exponent&1 != 0 {
			result = frameGFMul(result, base)
		}
		base = frameGFMul(base, base)
		exponent >>= 1
	}
	return result
}
