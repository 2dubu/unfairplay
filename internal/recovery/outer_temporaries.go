package recovery

func outerTablePendingS2(s2, tap70, tap112, tap154 uint32) uint32 {
	return s2 ^
		(tap70 & 0x224ea33b) ^
		(tap112 & 0x0e054c0f) ^
		(tap70 & tap112) ^
		(tap154 & 0xddb15cc4) ^
		(tap112 & tap154) ^
		0x739060ae
}

func outerTableLane0PendingS1(pendingS2, s0, s1, s3, tap27, tap69, tap111, tap152 uint32) uint32 {
	result := uint32(0x569c3294) ^
		(pendingS2 & 0xf7fedfef) ^
		(s0 & 0x50e8c30c) ^
		s1 ^
		(s3 & 0xffbffcd7) ^
		(s0 & s3 & 0x51e9e7af) ^
		(tap27 & 0x10612130) ^
		(tap69 & 0x0a048000) ^
		(tap27 & tap69 & 0x08402000) ^
		(tap111 & 0x20a00002) ^
		(tap27 & tap111 & 0x00200000) ^
		(tap152 & 0xc0081008) ^
		(tap27 & tap152 & 0x00001000)
	result ^= tap69 & tap152 &^ tap27 &^ tap111 & 0x00000004
	return result
}

func outerTableLane0PendingS3(pendingS2, s0, s3, s4 uint32) uint32 {
	return (pendingS2 & 0xefffefef) ^
		(s0 & 0xf87ff72b) ^
		s3 ^
		(s4 & 0x0e09cd90) ^
		(s0 & s4 & 0x1fbdfdf7) ^
		0x2e6f04e6
}

func outerTableLane5PendingS4(s4, s0, pendingS1, s3 uint32) uint32 {
	return s4 ^
		(s0 & 0x6dedf7dd) ^
		(pendingS1 & 0xf3ffefff) ^
		(s3 & 0xdfdfeffb) ^
		0xd6113160
}

func outerTableLane18Mix(s0, s2, s1, s3 uint32) uint32 {
	return s0 ^
		(s2 & 0xc4ff9cfa) ^
		(s1 & 0x69fff7ff) ^
		(s3 & 0x11634268) ^
		(s2 & s3 & 0x3ffbf77f) ^
		0x3f98f797
}

func outerTableLane19LateMix(mix, tap2, tap44, tap86 uint32) uint32 {
	return mix ^ tap2 ^ tap44 ^ tap86 ^ 0xf0820e81
}
