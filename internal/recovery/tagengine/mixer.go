package tagengine

import "math/bits"

// mixStep operates on decoded bytes. All indices wrap modulo 256 and all
// polynomial arithmetic wraps modulo 2^32, matching the mapped constructor.
// Stores precede subsequent reads deliberately: indexed operands can alias.
func mixStep(state *[256]byte, iteration int) {
	i := byte(iteration)
	switch state[i] & 3 {
	case 0:
		x, y := uint32(state[i]), uint32(state[byte(5*int(i)+8)])
		x2, y2 := x*x, y*y
		p := y + 231*x*y2 + x2*y2*y + 11*x2*x*y2*y2*y + 117*x2*x2*x
		// Preserve signed remainder semantics after the polynomial wraps.
		r := uint32(int32(p) % 599)
		state[i] = constructorSubstitution[byte(129*r+201)] ^ 0x1b
	case 1:
		x := uint32(state[i])<<8 | uint32(state[i+3])
		y := uint32(state[i+5])<<8 | uint32(state[i+9])
		x2, y2 := x*x, y*y
		n := uint32(0x7e0d) + 0xf707*x*y2*y2*y + 0xa14b*x2*x2*x*y
		state[i] = byte(int32(n) % 0xf6eb)
	case 2:
		state[i] = 21 - 5*(constructorSubstitution[state[i]]^0x1b)
		state[i] = (state[i] ^ state[i+11]) + 1
		fallthrough
	case 3:
		state[i] = bits.RotateLeft8(state[i], -int(state[i+1]))
		p := uint32(state[i])<<24 | uint32(state[i+17])<<16 |
			uint32(state[byte(2*int(i))])<<8 | uint32(state[byte(int(i)*int(i)+17)])
		x := p + ((p << 3) ^ (p >> 5))
		y, z := x>>8, x>>16
		y2, z2 := y*y, z*z
		state[i] = byte(x*x*x + y2*y2*y + z2*z2*z2*z + z)
	}
}
