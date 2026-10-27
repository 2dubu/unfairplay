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

var outerFrontWord2Terms = [...]outerFrontTerm{
	{0x1a83c5e4, 0x0000}, {0xffffffff, 0x0001}, {0x05401555, 0x0002},
	{0xffffffff, 0x0004}, {0x218a3101, 0x0008}, {0x10c80002, 0x0010},
	{0x00400000, 0x0012}, {0x8014c800, 0x0020}, {0x05000004, 0x0022},
	{0x10000000, 0x0030}, {0x40006008, 0x0040}, {0x04002000, 0x0042},
	{0x10000008, 0x0050}, {0x00024200, 0x0060}, {0x04000200, 0x0062},
	{0x10020000, 0x0070}, {0x00008894, 0x0080}, {0x20002004, 0x0088},
	{0xffffffff, 0x0100}, {0xc2304000, 0x0200}, {0x001a0000, 0x0208},
	{0x14000000, 0x0280}, {0x224ea33b, 0x0400}, {0x158d2c0f, 0x0800},
	{0xffffffff, 0x0c00}, {0xddb15cc4, 0x1000}, {0xffffffff, 0x1800},
	{0x09441402, 0x2000}, {0x01001220, 0x2008}, {0x00000020, 0x2088},
	{0x00200000, 0x2200}, {0x00000020, 0x2208}, {0x00000020, 0x2288},
	{0xffffffff, 0x4000},
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

var outerFrontSourceTerms = [...]outerFrontTerm{
	{0x000000f3, 0x000000000000}, {0x000000c4, 0x000000000001},
	{0x000000d4, 0x000000000002}, {0x0000003b, 0x000000000004},
	{0x000000ff, 0x000000000008}, {0x000000ff, 0x000000000010},
	{0x00000001, 0x000000000020}, {0x0000003b, 0x000000000040},
	{0x0000007d, 0x000000000080}, {0x000000d9, 0x000000000100},
	{0x000000ff, 0x000000000200}, {0x00000002, 0x000000000400},
	{0x00000020, 0x000000000420}, {0x000000ff, 0x000000000800},
	{0x0000006b, 0x000000001000}, {0x000000ff, 0x000000001001},
	{0x000000ff, 0x000000001040}, {0x0000002b, 0x000000002000},
	{0x000000ff, 0x000000002002}, {0x00000014, 0x000000004000},
	{0x00000004, 0x000000004020}, {0x00000020, 0x000000004420},
	{0x000000c4, 0x000000008000}, {0x000000ff, 0x000000008004},
	{0x000000ff, 0x000000010000}, {0x000000ff, 0x000000020000},
	{0x00000020, 0x000000040420}, {0x00000020, 0x000000044420},
	{0x00000090, 0x000000080000}, {0x000000ad, 0x000000100000},
	{0x000000ff, 0x000000200000}, {0x00000008, 0x000000400000},
	{0x00000001, 0x000000800000}, {0x00000002, 0x000001000000},
	{0x000000ff, 0x000002000000}, {0x0000003b, 0x000004000000},
	{0x0000007f, 0x000008000000}, {0x00000040, 0x000010000000},
	{0x00000040, 0x000010400000}, {0x00000010, 0x000020000000},
	{0x00000010, 0x000020080000}, {0x00000040, 0x000040000000},
	{0x00000009, 0x000040100000}, {0x00000002, 0x000041000000},
	{0x000000ff, 0x000080000000}, {0x0000000f, 0x000100000000},
	{0x0000007f, 0x000104000000}, {0x0000007f, 0x000200000000},
	{0x00000025, 0x000400000000}, {0x00000020, 0x000410000000},
	{0x00000008, 0x000800400000}, {0x00000020, 0x000c00000000},
	{0x00000020, 0x000c10000000}, {0x00000014, 0x001000000000},
	{0x00000012, 0x001000080000}, {0x00000010, 0x001020000000},
	{0x00000010, 0x001020080000}, {0x00000044, 0x002000000000},
	{0x0000007f, 0x002100000000}, {0x0000007f, 0x004000000000},
	{0x0000007f, 0x008000000000}, {0x000000ff, 0x010000000000},
	{0x0000007f, 0x020000000000}, {0x000000ff, 0x040000000000},
	{0x000000ff, 0x080000000000}, {0x000000ff, 0x100000000000},
	{0x000000ff, 0x200000000000}, {0x000000ff, 0x400000000000},
}

type outerFrontTableRoute struct {
	multiplier uint32
	offset     uint32
}

var outerFrontSourceCoreRoutes = [...]outerFrontTableRoute{
	{1, 104}, {1, 109}, {1, 125}, {1, 130}, {1, 131}, {1, 132},
	{1, 20}, {1, 25}, {1, 41}, {1, 46}, {1, 49}, {1, 6}, {1, 62},
	{1, 67}, {1, 7}, {1, 83}, {1, 88}, {1, 89}, {1, 90},
	{14, 10}, {14, 12}, {14, 13}, {14, 134}, {14, 136}, {14, 137},
	{14, 138}, {14, 14}, {14, 141}, {14, 51}, {14, 52}, {14, 53},
	{14, 55}, {14, 56}, {14, 57}, {14, 9}, {14, 93}, {14, 94},
	{14, 98}, {14, 99},
}

var outerFrontSourceCoreShifts = [...]int{
	0, 21, 0, 21, 0, 0, 0, 21, 0, 21, 0, 0, 0, 21, 0, 0,
	21, 0, 0, 0, 21, 0, 0, 0, 21, 0, 0, 0, 0, 0, 21, 0,
	0, 0, 0, 0, 0, 0, 0,
}

var outerFrontSourceExtras = [12][3]uint32{
	{103, 60, 143}, {58, 15, 98}, {13, 137, 53}, {135, 92, 8},
	{90, 47, 130}, {45, 2, 85}, {0, 124, 40}, {122, 79, 162},
	{77, 34, 117}, {32, 156, 72}, {154, 111, 27}, {109, 66, 149},
}

func outerFrontWords(table *[outerTableWords]uint32, base uint32, digest *outerDigestState) ([5]uint32, bool) {
	if table == nil || digest == nil || base >= outerTableWords {
		return [5]uint32{}, false
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
	return [5]uint32{
		outerFrontPolynomial(word0Inputs[:], outerFrontWord0Terms[:]),
		outerFrontPolynomial(word1Inputs[:], outerFrontWord1Terms[:]),
		outerFrontPolynomial(word2Inputs[:], outerFrontWord2Terms[:]),
		outerFrontPolynomial(word3Inputs[:], outerFrontWord3Terms[:]),
		outerFrontPolynomial(word4Inputs[:], outerFrontWord4Terms[:]),
	}, true
}

func outerFrontSourcePrefix(table *[outerTableWords]uint32, base uint32, digest *outerDigestState) ([12]byte, bool) {
	var output [12]byte
	if table == nil || digest == nil || base >= outerTableWords {
		return output, false
	}
	var common [44]uint32
	for index, route := range outerFrontSourceCoreRoutes {
		word := table[(route.multiplier*base+route.offset)%outerTableWords]
		common[index] = bits.RotateLeft32(word, outerFrontSourceCoreShifts[index])
	}
	digestWords := [...]uint32{digest.X0, digest.X4, digest.X10, digest.X27, digest.X28}
	digestShifts := [...]int{0, 0, 0, 0, 21}
	for index, word := range digestWords {
		common[len(outerFrontSourceCoreRoutes)+index] = bits.RotateLeft32(word, digestShifts[index])
	}
	for outputIndex, offsets := range outerFrontSourceExtras {
		var inputs [47]uint32
		copy(inputs[:], common[:])
		inputs[44] = table[(base+offsets[0])%outerTableWords]
		inputs[45] = bits.RotateLeft32(table[(base+offsets[1])%outerTableWords], 23)
		inputs[46] = bits.RotateLeft32(table[(base+offsets[2])%outerTableWords], 10)
		output[outputIndex] = byte(outerFrontPolynomial(inputs[:], outerFrontSourceTerms[:]))
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
