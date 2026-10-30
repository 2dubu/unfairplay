package recovery

import "math/bits"

func outerTableLane2Mask(pendingS1, pendingS2, pendingS3, s0, s4 uint32) uint32 {
	combined := s0 ^
		bits.RotateLeft32(pendingS2, -1) ^
		bits.RotateLeft32(pendingS1, -9) ^
		bits.RotateLeft32(pendingS3, -18)
	parts := ((s4 & 0x0a213609) ^ 0x0a211601) |
		((combined & 0x3d838485) ^ 0x0d830085) |
		((pendingS2 & 0x1dc2642c) ^ 0x09006424) |
		((s4 & 0x111d4a54) ^ 0x00000210)
	return ^parts & 0x7fffffff
}

func outerTableLane7Mask(s0, s1, s3, s4, _ uint32) uint32 {
	return s0 ^
		(s1 & 0x00010000) ^
		(s3 & 0x80000c00) ^
		(s4 & 0xffff6f7f) ^
		0xcac4ea7e
}

func outerTableLane19Mask(mix, s1, s4, tap2, tap44, tap86 uint32) uint32 {
	combined := mix ^ tap2 ^ tap44 ^ tap86 ^ 0x0e0ca003
	parts := ((s1 & 0x0813cd04) ^ 0x08030900) |
		((s4 & 0x2b680826) ^ 0x2b600000) |
		((combined & 0x486a2a10) ^ 0x402a2800) |
		((combined & 0x39628612) ^ 0x28428610)
	return (^parts & 0x7fffffff) | (^s1 & 0x80000000)
}

func outerTableLane18Mask(s1, s3, s0, s2, _ uint32) uint32 {
	return s0 ^
		(s1 & 0x69fff7ff) ^
		(s3 & 0xee9cbd97) ^
		(s2 & 0xc4ff9cfa) ^
		(s3 & s2 & 0x3ffbf77f) ^
		0x3182dfbf
}
