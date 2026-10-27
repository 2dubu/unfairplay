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

var outerHookTailTerms = [...]outerHookTerm{
	{0xf0fa8bd0, 0x000000}, {0xffffffff, 0x000001}, {0x000007ff, 0x000002},
	{0xffffffff, 0x000004}, {0x000007ff, 0x000008}, {0xffffffff, 0x000010},
	{0x000007ff, 0x000020}, {0xc2304000, 0x000040}, {0x00000611, 0x000080},
	{0x218a3101, 0x000100}, {0x001a0000, 0x000140}, {0x0000010c, 0x000200},
	{0x00008894, 0x000400}, {0x14000000, 0x000440}, {0x20002004, 0x000500},
	{0x000000a0, 0x000880}, {0x00000100, 0x000a00}, {0x09441402, 0x001000},
	{0x00200000, 0x001040}, {0x01001220, 0x001100}, {0x00000020, 0x001140},
	{0x00000020, 0x001500}, {0x00000020, 0x001540}, {0x0000004a, 0x002000},
	{0x00000001, 0x002080}, {0x00000008, 0x002200}, {0x158d2c0f, 0x004000},
	{0x000000ac, 0x008000}, {0x224ea33b, 0x010000}, {0xffffffff, 0x014000},
	{0x00000112, 0x020000}, {0x000007ff, 0x028000}, {0xddb15cc4, 0x040000},
	{0xffffffff, 0x044000}, {0x000006ed, 0x080000}, {0x000007ff, 0x088000},
	{0xffffffff, 0x100000}, {0x000007ff, 0x200000},
}

func deriveOuterHookSeeds(table *[outerTableWords]uint32, base uint32, digest *outerDigestState) (outerHookSeeds, bool) {
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

	tailWords := [11]uint32{
		table[(base+1)%outerTableWords], table[(base+126)%outerTableWords],
		table[(base+84)%outerTableWords], table[(base+90)%outerTableWords],
		table[(base+132)%outerTableWords], table[(base+7)%outerTableWords],
		table[(base+49)%outerTableWords], table[(14*base+56)%outerTableWords],
		table[(14*base+14)%outerTableWords], table[(14*base+98)%outerTableWords],
		digest.X10,
	}
	var tailInputs [22]uint32
	for index, word := range tailWords {
		tailInputs[index*2] = word
		tailInputs[index*2+1] = bits.RotateLeft32(word, 11)
	}
	tail := outerHookPolynomial(tailInputs[:], outerHookTailTerms[:])
	return outerHookSeeds{Even: even, Odd: odd, Tail: tail}, true
}
