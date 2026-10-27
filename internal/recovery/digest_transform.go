package recovery

// digestLanes26A428 is the semantic form of the four lane stores at the front
// of each digest iteration. Results are ordered by their ascending stack slots
// (d78, d80, d90, d98).
func digestLanes26A428(a, b, c, d uint32) [4]uint32 {
	return [4]uint32{
		d - (0x34 & (d << 1)),
		a - (0xa0 & (a << 1)),
		c - (0x74 & (c << 1)),
		b - (0x184 & (b << 1)),
	}
}
