package recovery

import "math/bits"

type outerFrontTerm struct {
	mask    uint32
	product uint64
}

func outerFrontPolynomial(inputs []uint32, terms []outerFrontTerm) uint32 {
	var output uint32
	for _, term := range terms {
		value := term.mask
		for product := term.product; product != 0; product &= product - 1 {
			value &= inputs[bits.TrailingZeros64(product)]
		}
		output ^= value
	}
	return output
}

var outerFrontWord0Terms = [...]outerFrontTerm{
	{0x6d53971b, 0x0000}, {0x60022140, 0x0001}, {0x224ea33b, 0x0002},
	{0x0b901081, 0x0004}, {0x08921080, 0x0005}, {0x7187de73, 0x0008},
	{0xffffffff, 0x000a}, {0x04048228, 0x0010}, {0x00480010, 0x0011},
	{0x08001080, 0x0014}, {0x08401090, 0x0015}, {0xddb15cc4, 0x0020},
	{0xffffffff, 0x0028}, {0x10200003, 0x0040}, {0x10202100, 0x0041},
	{0x00004001, 0x0044}, {0x00004000, 0x0045}, {0x00200008, 0x0050},
	{0x00280010, 0x0051}, {0x00004000, 0x0054}, {0x00004010, 0x0055},
	{0x80560608, 0x0080}, {0xffffffff, 0x0100}, {0x020800c0, 0x0200},
	{0x00460040, 0x0280}, {0xffffffff, 0x0400}, {0x00000025, 0x0800},
	{0x000000a0, 0x0a00}, {0x17004000, 0x1000}, {0x05000008, 0x1080},
	{0x02001800, 0x1200}, {0x00000800, 0x1280}, {0x15000020, 0x1800},
	{0x05000000, 0x1880}, {0x00000020, 0x1a00}, {0xffffffff, 0x2000},
	{0xffffffff, 0x4000},
}

var outerFrontWord1Terms = [...]outerFrontTerm{
	{0xf662b4bd, 0x0000000}, {0xddb15cc4, 0x0000001}, {0x224ea33b, 0x0000002},
	{0xffffffff, 0x0000004}, {0x21823101, 0x0000008}, {0x224ea33b, 0x0000010},
	{0xb1e73fd9, 0x0000020}, {0x09441402, 0x0000040}, {0x01001220, 0x0000048},
	{0xffffffff, 0x0000080}, {0x667f336b, 0x0000100}, {0xffffffff, 0x0000101},
	{0xffffffff, 0x0000110}, {0x00008814, 0x0000200}, {0x20002004, 0x0000208},
	{0x00000020, 0x0000248}, {0xddb15cc4, 0x0000400}, {0xffffffff, 0x0000402},
	{0xffffffff, 0x0000800}, {0xc2304000, 0x0001000}, {0x00120000, 0x0001008},
	{0x00200000, 0x0001040}, {0x00000020, 0x0001048}, {0x14000000, 0x0001200},
	{0x00000020, 0x0001248}, {0x80006090, 0x0002000}, {0xffffffff, 0x0004000},
	{0x00000001, 0x0008000}, {0x00000800, 0x000a000}, {0xffffffff, 0x0010000},
	{0x2246a33b, 0x0020000}, {0x00000410, 0x0040000}, {0x00000410, 0x0042000},
	{0x04010000, 0x0048000}, {0x00000800, 0x004a000}, {0xffffffff, 0x0080000},
	{0x15852c0f, 0x0100000}, {0xfff7ff7f, 0x0120000}, {0x08000014, 0x0200000},
	{0x00000012, 0x0202000}, {0x00400000, 0x0208000}, {0x00100000, 0x020a000},
	{0x00000010, 0x0240000}, {0x00000010, 0x0242000}, {0x04000000, 0x0248000},
	{0xddb15c44, 0x000400000}, {0xfff7ff7f, 0x000500000}, {0x80560608, 0x000800000},
	{0xffff3f7f, 0x001000000}, {0x02080040, 0x002000000}, {0x00460040, 0x002800000},
	{0xffff3f7f, 0x004000000}, {0x00000025, 0x008000000}, {0x00000020, 0x00a000000},
	{0x17000000, 0x010000000}, {0x05000008, 0x010800000}, {0x02001800, 0x012000000},
	{0x00000800, 0x012800000}, {0x15000020, 0x018000000}, {0x05000000, 0x018800000},
	{0x00000020, 0x01a000000}, {0xffff3f7f, 0x020000000}, {0xfff7ff7f, 0x040000000},
	{0xffffffff, 0x080000000}, {0xbfffffff, 0x100000000}, {0xffff3f7f, 0x200000000},
}

// The core is shared with the tail hook seed. The remaining terms depend
// only on four additional table taps.
var outerFrontWord2CoreTerms = [...]outerFrontTerm{
	{0x1a83c5e4, 0x0000}, {0xffffffff, 0x0001}, {0xffffffff, 0x0004},
	{0x218a3101, 0x0008}, {0x00008894, 0x0080}, {0x20002004, 0x0088},
	{0xffffffff, 0x0100}, {0xc2304000, 0x0200}, {0x001a0000, 0x0208},
	{0x14000000, 0x0280}, {0x224ea33b, 0x0400}, {0x158d2c0f, 0x0800},
	{0xffffffff, 0x0c00}, {0xddb15cc4, 0x1000}, {0xffffffff, 0x1800},
	{0x09441402, 0x2000}, {0x01001220, 0x2008}, {0x00000020, 0x2088},
	{0x00200000, 0x2200}, {0x00000020, 0x2208}, {0x00000020, 0x2288},
	{0xffffffff, 0x4000},
}

var outerFrontWord2ExtraTerms = [...]outerFrontTerm{
	{0x05401555, 0x0002}, {0x10c80002, 0x0010}, {0x00400000, 0x0012},
	{0x8014c800, 0x0020}, {0x05000004, 0x0022}, {0x10000000, 0x0030},
	{0x40006008, 0x0040}, {0x04002000, 0x0042}, {0x10000008, 0x0050},
	{0x00024200, 0x0060}, {0x04000200, 0x0062}, {0x10020000, 0x0070},
}

var outerFrontWord3Terms = [...]outerFrontTerm{
	{0x0b400509, 0x000}, {0xc2000000, 0x001}, {0x0a429600, 0x002},
	{0x80000002, 0x004}, {0x80010100, 0x005}, {0x00008008, 0x008},
	{0x0004c000, 0x009}, {0x00100000, 0x00c}, {0x00100000, 0x00d},
	{0x00480080, 0x010}, {0x0a400088, 0x012}, {0x00200240, 0x020},
	{0x00001250, 0x022}, {0x00000080, 0x030}, {0x00000080, 0x032},
	{0x01020000, 0x040}, {0x40020040, 0x041}, {0x08000000, 0x044},
	{0x00802000, 0x048}, {0x00001000, 0x049}, {0x08000000, 0x04c},
	{0x11020100, 0x080}, {0x01029000, 0x082}, {0x01000000, 0x090},
	{0x01000000, 0x092}, {0x00001000, 0x0a2}, {0x80006090, 0x0100},
	{0x00000001, 0x0200}, {0x00000800, 0x0300}, {0x00000410, 0x0400},
	{0x00000410, 0x0500}, {0x04010000, 0x0600}, {0x00000800, 0x0700},
	{0x08000014, 0x0800}, {0x00000012, 0x0900}, {0x00400000, 0x0a00},
	{0x00100000, 0x0b00}, {0x00000010, 0x0c00}, {0x00000010, 0x0d00},
	{0x04000000, 0x0e00}, {0xffffffff, 0x1000},
}

var outerFrontWord4Terms = [...]outerFrontTerm{
	{0x9a13abe8, 0x000}, {0x224ea33b, 0x001}, {0xffffffff, 0x002},
	{0xef7becfb, 0x004}, {0xffffffff, 0x008}, {0xddb15cc4, 0x010},
	{0xffffffff, 0x011}, {0xffffffff, 0x020}, {0x4a856940, 0x040},
	{0x35081104, 0x080}, {0x00000100, 0x0c0}, {0x80220080, 0x100},
	{0x08804800, 0x140}, {0x10001000, 0x180}, {0x0020001a, 0x200},
	{0x48800018, 0x240}, {0x00000004, 0x280}, {0x00200000, 0x300},
	{0x08800000, 0x340}, {0xffffffff, 0x400},
}

var outerFrontSourceExtras = [12][3]uint32{
	{103, 60, 143}, {58, 15, 98}, {13, 137, 53}, {135, 92, 8},
	{90, 47, 130}, {45, 2, 85}, {0, 124, 40}, {122, 79, 162},
	{77, 34, 117}, {32, 156, 72}, {154, 111, 27}, {109, 66, 149},
}

func outerFrontWords(table *[outerTableWords]uint32, base uint32, digest *outerDigestState) ([5]uint32, uint32, bool) {
	if table == nil || digest == nil || base >= outerTableWords {
		return [5]uint32{}, 0, false
	}
	word0Inputs := [15]uint32{
		table[(base+106)%outerTableWords], table[(base+127)%outerTableWords],
		table[(base+148)%outerTableWords], table[(base+2)%outerTableWords],
		table[(base+23)%outerTableWords], table[(base+44)%outerTableWords],
		table[(base+64)%outerTableWords], table[(14*base+134)%outerTableWords],
		table[(14*base+141)%outerTableWords], table[(14*base+51)%outerTableWords],
		table[(14*base+57)%outerTableWords], table[(14*base+9)%outerTableWords],
		table[(14*base+93)%outerTableWords], table[(14*base+99)%outerTableWords], digest.X0,
	}
	word1Inputs := [34]uint32{
		table[(base+104)%outerTableWords], table[(base+125)%outerTableWords],
		table[(base+131)%outerTableWords], table[(base+132)%outerTableWords],
		table[(base+20)%outerTableWords], table[(base+41)%outerTableWords],
		table[(base+49)%outerTableWords], table[(base+6)%outerTableWords],
		table[(base+62)%outerTableWords], table[(base+7)%outerTableWords],
		table[(base+83)%outerTableWords], table[(base+89)%outerTableWords],
		table[(base+90)%outerTableWords], table[(14*base+10)%outerTableWords],
		table[(14*base+13)%outerTableWords], table[(14*base+136)%outerTableWords],
		table[(14*base+138)%outerTableWords], table[(14*base+14)%outerTableWords],
		table[(14*base+52)%outerTableWords], table[(14*base+55)%outerTableWords],
		table[(14*base+56)%outerTableWords], table[(14*base+94)%outerTableWords],
		table[(14*base+98)%outerTableWords], table[(14*base+134)%outerTableWords],
		table[(14*base+141)%outerTableWords], table[(14*base+51)%outerTableWords],
		table[(14*base+57)%outerTableWords], table[(14*base+9)%outerTableWords],
		table[(14*base+93)%outerTableWords], table[(14*base+99)%outerTableWords],
		digest.X10, digest.X27, digest.X4, digest.X0,
	}
	word2Inputs := [15]uint32{
		table[(base+1)%outerTableWords], table[(base+105)%outerTableWords],
		table[(base+126)%outerTableWords], table[(base+132)%outerTableWords],
		table[(base+147)%outerTableWords], table[(base+21)%outerTableWords],
		table[(base+63)%outerTableWords], table[(base+7)%outerTableWords],
		table[(base+84)%outerTableWords], table[(base+90)%outerTableWords],
		table[(14*base+14)%outerTableWords], table[(14*base+56)%outerTableWords],
		table[(14*base+98)%outerTableWords], table[(base+49)%outerTableWords], digest.X10,
	}
	word3Inputs := [13]uint32{
		table[(base+107)%outerTableWords], table[(base+128)%outerTableWords],
		table[(base+149)%outerTableWords], table[(base+24)%outerTableWords],
		table[(base+3)%outerTableWords], table[(base+45)%outerTableWords],
		table[(base+66)%outerTableWords], table[(base+87)%outerTableWords],
		table[(14*base+10)%outerTableWords], table[(14*base+136)%outerTableWords],
		table[(14*base+52)%outerTableWords], table[(14*base+94)%outerTableWords], digest.X4,
	}
	word4Inputs := [11]uint32{
		table[(base+109)%outerTableWords], table[(base+130)%outerTableWords],
		table[(base+25)%outerTableWords], table[(base+46)%outerTableWords],
		table[(base+67)%outerTableWords], table[(base+88)%outerTableWords],
		table[(14*base+12)%outerTableWords], table[(14*base+137)%outerTableWords],
		table[(14*base+53)%outerTableWords], table[(14*base+95)%outerTableWords], digest.X28,
	}
	word2Core := outerFrontPolynomial(word2Inputs[:], outerFrontWord2CoreTerms[:])
	return [5]uint32{
		outerFrontPolynomial(word0Inputs[:], outerFrontWord0Terms[:]),
		outerFrontPolynomial(word1Inputs[:], outerFrontWord1Terms[:]),
		word2Core ^ outerFrontPolynomial(word2Inputs[:], outerFrontWord2ExtraTerms[:]),
		outerFrontPolynomial(word3Inputs[:], outerFrontWord3Terms[:]),
		outerFrontPolynomial(word4Inputs[:], outerFrontWord4Terms[:]),
	}, word2Core, true
}

// outerFrontSourcePrefix shares the front words already computed from the
// digest. The common source circuit cancels to word 1 XOR word 4 shifted by
// eleven; only three table taps vary between output bytes.
func outerFrontSourcePrefix(table *[outerTableWords]uint32, base uint32, words [5]uint32) ([12]byte, bool) {
	var output [12]byte
	if table == nil || base >= outerTableWords {
		return output, false
	}
	common := words[1] ^ words[4]>>11 ^ 0x3b
	for index, offsets := range outerFrontSourceExtras {
		output[index] = byte(common ^ table[(base+offsets[0])%outerTableWords] ^
			table[(base+offsets[1])%outerTableWords]>>9 ^
			table[(base+offsets[2])%outerTableWords]>>22)
	}
	return output, true
}

func outerFrontSourceTail(table *[outerTableWords]uint32, base, oldWord1, oldWord4 uint32) ([4]byte, bool) {
	var output [4]byte
	if table == nil || base >= outerTableWords {
		return output, false
	}
	pairs := [4][3]uint32{
		{36, 160, 118}, {159, 116, 74}, {115, 72, 30}, {71, 28, 153},
	}
	for index, offsets := range pairs {
		output[index] = byte(table[(7*base+offsets[0])%outerTableWords]) ^
			byte(oldWord1) ^ byte(oldWord4>>13) ^
			byte(table[(7*base+offsets[1])%outerTableWords]>>17) ^
			byte(table[(7*base+offsets[2])%outerTableWords]>>7) ^ 0xd3
	}
	return output, true
}
