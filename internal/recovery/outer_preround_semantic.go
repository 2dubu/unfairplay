package recovery

import "math/bits"

var outerPreRoundWord1Terms = [...]outerFrontTerm{
	{0x1003f369, 0x0000000}, {0x44200408, 0x0000001}, {0x0049a400, 0x0000002},
	{0x10201490, 0x0000004}, {0xf7fffb9f, 0x0000008}, {0x00061892, 0x0000010},
	{0x04060010, 0x0000011}, {0xffffffff, 0x0000020}, {0x00340244, 0x0000040},
	{0x00708040, 0x0000042}, {0xc481407a, 0x0000080}, {0x08000010, 0x0000084},
	{0xddb15cc4, 0x0000100}, {0xffffffff, 0x0000200}, {0x21001039, 0x0000400},
	{0x00000002, 0x0000402}, {0x20200045, 0x0000440}, {0x00200042, 0x0000442},
	{0x0008aa00, 0x0000800}, {0x10200000, 0x0000804}, {0xc0000040, 0x0000880},
	{0x08000000, 0x0000884}, {0x2ff6ffcb, 0x0001000}, {0x224ea33b, 0x0002000},
	{0xffffffff, 0x0002100}, {0xddb15cc4, 0x0004000}, {0x2ad1a081, 0x0008000},
	{0x00202000, 0x0008001}, {0x000000a0, 0x0008010}, {0x00000020, 0x0008011},
	{0xffffffff, 0x0010000}, {0x10800800, 0x0020000}, {0x00008800, 0x0020002},
	{0x08004000, 0x0020040}, {0x0000c000, 0x0020042}, {0x00000018, 0x0020400},
	{0x00000002, 0x0020402}, {0x00000002, 0x0020442}, {0x02020100, 0x0040000},
	{0x02020180, 0x0040004}, {0xc2830040, 0x0040080}, {0x02020000, 0x0040084},
	{0x20040000, 0x0040800}, {0xc0000040, 0x0040880}, {0x224ea33b, 0x0080000},
	{0xffffffff, 0x0084000}, {0x01084100, 0x0100000}, {0x00004002, 0x0100010},
	{0x00400000, 0x0108000}, {0xffffffff, 0x0200000}, {0x00003fc0, 0x0400000},
	{0x0007f800, 0x0800000}, {0xff000000, 0x1000000}, {0xc000003f, 0x2000000},
	{0x01fe0000, 0x4000000}, {0x000007f8, 0x8000000},
}

func outerPreRoundWord1(table *[outerTableWords]uint32, base uint32, state *outerTableRoundState) (uint32, bool) {
	if table == nil || state == nil || base >= outerTableWords || state.HardwareIndex >= 6 {
		return 0, false
	}
	offsets := [...]uint32{10, 105, 117, 123, 136, 143, 147, 159, 165, 18, 22, 33, 4, 40, 46, 52, 60, 63, 75, 88, 94}
	var inputs [28]uint32
	for index, offset := range offsets {
		inputs[index] = table[(base+offset)%outerTableWords]
	}
	inputs[21] = state.Words[4]
	inputs[22] = bits.RotateLeft32(uint32(state.Hardware[(int(state.HardwareIndex)+3)%6]), 6)
	inputs[23] = bits.RotateLeft32(uint32(state.Hardware[(int(state.HardwareIndex)+4)%6]), 11)
	inputs[24] = bits.RotateLeft32(uint32(state.Hardware[(int(state.HardwareIndex)+5)%6]), 24)
	inputs[25] = bits.RotateLeft32(uint32(state.Codec[15]), 30)
	inputs[26] = bits.RotateLeft32(uint32(state.Codec[4]), 17)
	inputs[27] = bits.RotateLeft32(uint32(state.Codec[6]), 3)
	return outerFrontPolynomial(inputs[:], outerPreRoundWord1Terms[:]), true
}
