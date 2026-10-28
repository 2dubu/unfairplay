package recovery

import "math/bits"

// applyOuterCycle folds one codec byte and the six-byte hardware identity into
// the five-word outer state. The protected implementation expressed these XOR
// masks as hundreds of bit slices and opaque arithmetic identities.
func applyOuterCycle(state *outerTableRoundState, cycleIndex uint32) bool {
	if state == nil || state.HardwareIndex >= 6 || cycleIndex == 0 || int(cycleIndex)+3 >= len(state.Codec) {
		return false
	}
	byte32 := func(value byte) uint32 { return uint32(value) }
	hardware := func(offset int) uint32 {
		return byte32(state.Hardware[(int(state.HardwareIndex)+offset)%len(state.Hardware)])
	}
	h0, h1, h2 := hardware(0), hardware(1), hardware(2)
	h3, h4, h5 := hardware(3), hardware(4), hardware(5)
	current := byte32(state.Codec[int(cycleIndex)+3])
	c4, c5, c7 := byte32(state.Codec[4]), byte32(state.Codec[5]), byte32(state.Codec[7])
	c8, c9, c10 := byte32(state.Codec[8]), byte32(state.Codec[9]), byte32(state.Codec[10])
	c12, c13 := byte32(state.Codec[12]), byte32(state.Codec[13])
	c14, c15, c16 := byte32(state.Codec[14]), byte32(state.Codec[15]), byte32(state.Codec[16])

	state.Words[0] ^= h0<<5 ^ h1<<8 ^ bits.RotateLeft32(h2, 30) ^
		current<<3 ^ current<<9 ^ current<<20
	state.Words[1] ^= h3<<4 ^ h4<<14 ^ h5<<18 ^ current<<3 ^ current<<10 ^
		bits.RotateLeft32(current, 30) ^ bits.RotateLeft32(c8, 26) ^ c10<<14 ^ c12<<10
	state.Words[2] ^= h1<<7 ^ h2<<2 ^ h2<<15 ^ h3<<12 ^ h3<<18 ^
		bits.RotateLeft32(h4, 28) ^ current<<5 ^ current<<6 ^ current<<13 ^
		current<<16 ^ bits.RotateLeft32(current, 28) ^ c5<<16 ^ c7<<5 ^
		bits.RotateLeft32(c10, 25) ^ c12<<17 ^ c14<<7 ^ c16<<23
	state.Words[3] ^= h0<<12 ^ h1<<22 ^ h5<<1 ^ current<<2 ^ current<<3 ^
		current<<15 ^ current<<21 ^ current<<23 ^ bits.RotateLeft32(c5, 28) ^
		c7<<19 ^ c9<<4
	state.Words[4] ^= bits.RotateLeft32(h0, 26) ^ h4<<1 ^ h5<<11 ^ current<<7 ^
		current<<11 ^ current<<16 ^ current<<18 ^ bits.RotateLeft32(current, 28) ^
		c4<<3 ^ bits.RotateLeft32(c13, 28) ^ c15<<14

	index := uint64(state.HardwareIndex)
	state.CycleRegisters = [5]uint64{
		index + 1,
		0xffe67ddb + index,
		uint64(h2),
		index >> 2,
		uint64(h5),
	}
	return true
}
