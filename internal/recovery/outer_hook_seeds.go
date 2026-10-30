package recovery

import "math/bits"

type outerHookSeeds struct{ Even, Odd, Tail uint32 }

type outerHookTerm struct {
	mask    uint32
	product uint32
}

func outerHookPolynomial(inputs []uint32, terms []outerHookTerm) uint32 {
	var output uint32
	for _, term := range terms {
		value := term.mask
		for product := term.product; product != 0; product &= product - 1 {
			value &= inputs[bits.TrailingZeros32(product)]
		}
		output ^= value
	}
	return output
}

var outerHookOddTerms = [...]outerHookTerm{
	{0xde9fb102, 0x000}, {0xffffffff, 0x001}, {0xffffffff, 0x002}, {0xffffffff, 0x004},
	{0x07ffffff, 0x008}, {0x07ffffff, 0x010}, {0x07ffffff, 0x020}, {0x0402b030, 0x040},
	{0x00000001, 0x080}, {0x00b80200, 0x100}, {0x00280000, 0x140}, {0x00a80001, 0x180},
	{0x00280000, 0x1c0}, {0x00104006, 0x200}, {0x00023002, 0x240}, {0x00000005, 0x280},
	{0x001000c0, 0x300}, {0x00000040, 0x340}, {0x00000001, 0x380}, {0xffffffff, 0x400},
	{0x07ffffff, 0x800},
}

func deriveOuterHookSeedsWithCore(table *[outerTableWords]uint32, base uint32, digest *outerDigestState, word2Core uint32) (outerHookSeeds, bool) {
	if table == nil || digest == nil || base >= outerTableWords {
		return outerHookSeeds{}, false
	}

	even := digest.X28 ^ digest.X0>>15
	oddInputs := [12]uint32{
		table[(14*base+55)%outerTableWords],
		table[(14*base+13)%outerTableWords],
		table[(14*base+138)%outerTableWords],
		bits.RotateLeft32(table[(14*base+141)%outerTableWords], 27),
		bits.RotateLeft32(table[(14*base+99)%outerTableWords], 27),
		bits.RotateLeft32(table[(14*base+57)%outerTableWords], 27),
		bits.RotateLeft32(table[(14*base+134)%outerTableWords], 27),
		bits.RotateLeft32(table[(14*base+9)%outerTableWords], 27),
		bits.RotateLeft32(table[(14*base+93)%outerTableWords], 27),
		bits.RotateLeft32(table[(14*base+51)%outerTableWords], 27),
		digest.X27,
		bits.RotateLeft32(digest.X0, 27),
	}
	odd := outerHookPolynomial(oddInputs[:], outerHookOddTerms[:])

	// The tail circuit shares the front word-2 core, including its low
	// eleven bits folded from the top of the same word.
	tail := word2Core ^ word2Core>>21 ^ 0xea794ee0
	return outerHookSeeds{Even: even, Odd: odd, Tail: tail}, true
}
