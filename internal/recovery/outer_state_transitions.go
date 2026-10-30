// Recovered outer-table state transitions.
package recovery

import "math/bits"

func outerTableLane2NextS0(pendingS2, pendingS1, pendingS3, s0 uint32) uint32 {
	return s0 ^
		bits.RotateLeft32(pendingS2, -1) ^
		bits.RotateLeft32(pendingS1, -9) ^
		bits.RotateLeft32(pendingS3, -18) ^
		0x5de927cd
}

func outerTableLane7NextS0(s0, s1, s3, s4 uint32) uint32 {
	return s0 ^
		(s1 & 0xfffeffff) ^
		(s3 & 0x7ffff3ff) ^
		(s4 & 0xffff6f7f) ^
		0xb9b2aaa4
}

func outerTableLane8NextS1(mix17, s2, mix16, s4 uint32) uint32 {
	left := ((s2 & 0xfd9fffff) ^ 0xc404d3f6) |
		(mix16 ^ 0x48d4e698)
	right := ((mix16 & 0x7ffc7654) ^ 0x7b006204) |
		(s4 ^ 0x71aae24b)
	return mix17 ^ ^(left & right)
}

func outerTableLane15NextS4(s1, s4, s2, s0 uint32) uint32 {
	mixMask := uint32(0x9c2cbf46)
	fromS2 := bits.RotateLeft32(s2, -4)&mixMask ^ 0x9008b240
	fromS0 := bits.RotateLeft32(s0, 21)&mixMask ^ 0x98283706
	return s4 ^ bits.RotateLeft32(s1, -17) ^ 0x57417d47 ^ (fromS2 | fromS0)
}

func outerTableLane19NextS0(mix uint32) uint32 {
	return mix ^ 0xfe8eae82
}
func outerTableLane5NextS2(s2, source0, source1, source2, codec0, codec1, codec2 uint32) uint32 {
	// The low byte of each source is packed into consecutive, overlapping
	// bit lanes. Overlap is intentional: each contribution is XORed.
	return s2 ^
		(codec0>>3)&0x1f ^
		(codec0&0x07)<<29 ^
		(source0&0xff)<<6 ^
		(codec2&0xff)<<10 ^
		(source1&0xff)<<14 ^
		(codec1&0xff)<<19 ^
		(source2&0xff)<<22
}
