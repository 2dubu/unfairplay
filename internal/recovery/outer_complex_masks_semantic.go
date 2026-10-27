package recovery

func outerTableLane7Mask(s0, s1, s3, s4, _ uint32) uint32 {
	return s0 ^
		(s1 & 0x00010000) ^
		(s3 & 0x80000c00) ^
		(s4 & 0xffff6f7f) ^
		0xcac4ea7e
}
