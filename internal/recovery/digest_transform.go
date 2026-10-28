package recovery

// digestLanes is the semantic form of the four lane stores at the front of
// each digest iteration. Results use the fixed order consumed by the next
// typed stage.
func digestLanes(a, b, c, d uint32) [4]uint32 {
	return [4]uint32{
		d - (0x34 & (d << 1)),
		a - (0xa0 & (a << 1)),
		c - (0x74 & (c << 1)),
		b - (0x184 & (b << 1)),
	}
}
