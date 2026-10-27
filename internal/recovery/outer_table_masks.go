package recovery

type outerBooleanTerm struct {
	mask    uint32
	product uint8
}

func outerBooleanMask(inputs []uint32, terms []outerBooleanTerm) uint32 {
	var output uint32
	for _, term := range terms {
		value := term.mask
		for input := uint8(0); input < uint8(len(inputs)); input++ {
			if term.product&(1<<input) != 0 {
				value &= inputs[input]
			}
		}
		output ^= value
	}
	return output
}

var outerTableLane1Terms = [...]outerBooleanTerm{
	{0xe42514dd, 0}, {0x02810004, 1}, {0x19000010, 2}, {0x02080800, 3},
	{0x80520000, 4}, {0x00100000, 5}, {0x01000000, 6}, {0xc4006130, 8},
	{0x00008222, 9}, {0x01004010, 10}, {0x00000200, 11}, {0x80400000, 12},
	{0x00000200, 13}, {0x01000000, 14}, {0x00000200, 15},
}

var outerTableLane2Terms = [...]outerBooleanTerm{
	{0x705c894a, 0}, {0x34828080, 1}, {0x10406028, 2}, {0x15820004, 3},
	{0x123c1a58, 4}, {0x10000001, 5}, {0x10006408, 6}, {0x19000404, 7},
}

var outerTableLane5Terms = [...]outerBooleanTerm{
	{0x868aca50, 0}, {0x854064c8, 1}, {0x20300016, 2}, {0x20002002, 3},
	{0x520d1200, 4}, {0x00401100, 5}, {0x00310000, 6},
}

var outerTableLane6Terms = [...]outerBooleanTerm{
	{0x530a702f, 0}, {0x00824441, 1}, {0x08035830, 2}, {0x08824190, 3},
	{0xa5544001, 4}, {0x00204001, 5}, {0x8c014000, 6}, {0x08004100, 7},
}

var outerTableLane9Terms = [...]outerBooleanTerm{
	{0x17a2bee8, 0}, {0x383df723, 1}, {0xc7c208dc, 2}, {0x30cd3dc4, 4},
	{0xffffffff, 5}, {0xffffffff, 6},
}

var outerTableLane11Terms = [...]outerBooleanTerm{
	{0xd0bb2ebd, 0}, {0x0000a502, 1}, {0x00200005, 2}, {0x00900000, 3},
	{0x010100b8, 4}, {0x00082020, 5}, {0x00000024, 6}, {0x00000020, 7},
	{0x00004200, 8}, {0xd0000800, 10}, {0x80100000, 11}, {0x00020010, 12},
	{0x40000000, 14},
}

var outerTableLane12Terms = [...]outerBooleanTerm{
	{0x7c3ed75e, 0}, {0x80001044, 1}, {0x81830909, 2}, {0x84020010, 3},
	{0x40400800, 4}, {0x40001040, 5}, {0x40020a00, 6}, {0x40020010, 7},
	{0x201ccc82, 8}, {0x20240402, 9}, {0x00000802, 10}, {0x04200002, 11},
	{0x18048800, 12}, {0x10040000, 13}, {0x00000800, 14},
}

var outerTableLane13Terms = [...]outerBooleanTerm{
	{0xc19056fd, 0}, {0x01542000, 1}, {0x32b0086c, 2}, {0x00120000, 3},
	{0x410001a2, 4}, {0x05000100, 5}, {0x20088020, 6}, {0x04000000, 7},
}

var outerTableLane14Terms = [...]outerBooleanTerm{
	{0x3cfd3b53, 0}, {0x04008404, 1}, {0x9a092020, 2}, {0x20092000, 3},
	{0x42040852, 4}, {0x00400002, 5}, {0x02800002, 6}, {0x00400002, 7},
	{0x02120180, 8}, {0x00000200, 9}, {0x02000001, 10}, {0x02201810, 12},
	{0x00001000, 13}, {0x02a00000, 14},
}

var outerTableLane20Terms = [...]outerBooleanTerm{
	{0x09b1c801, 0}, {0x54719560, 1}, {0x3c38c757, 2}, {0x00342202, 3},
	{0x5100c41d, 4}, {0xffcbddfd, 6}, {0xaa8a8a82, 8}, {0x5004a602, 9},
	{0xafabcbdf, 10}, {0x5475b762, 11}, {0x5004a602, 12}, {0x5475b762, 14},
}

func outerTableLane1Mask(pendingS1, pendingS2, pendingS3, s4 uint32) uint32 {
	return outerBooleanMask([]uint32{pendingS1, pendingS2, pendingS3, s4}, outerTableLane1Terms[:])
}

func outerTableLane2Mask(s0, s2, s4 uint32) uint32 {
	return outerBooleanMask([]uint32{s0, s2, s4}, outerTableLane2Terms[:])
}

func outerTableLane3Mask(s0, pendingS3 uint32) uint32 {
	return s0 ^ pendingS3 ^ 0x5d0bebe6
}

func outerTableLane4Mask(a, b, c, d uint32) uint32 {
	return ^(a ^ b ^ c ^ d)
}

// outerTableLane5Mask is the semantic Boolean polynomial for writer lane 5.
// pendingS1 is the next value staged for state word 1; s3 and s0 are the
// current state words at this point in the outer iteration.
func outerTableLane5Mask(pendingS1, s3, s0 uint32) uint32 {
	return outerBooleanMask([]uint32{pendingS1, s3, s0}, outerTableLane5Terms[:])
}

func outerTableLane6Mask(s1, pendingS4, s0 uint32) uint32 {
	return outerBooleanMask([]uint32{s1, pendingS4, s0}, outerTableLane6Terms[:])
}

func outerTableLane8Mask(a, b, c, d uint32) uint32 {
	return ^(a ^ b ^ c ^ d)
}

func outerTableLane9Mask(s2, s0, s4 uint32) uint32 {
	return outerBooleanMask([]uint32{s2, s0, s4}, outerTableLane9Terms[:])
}

func outerTableLane10Mask(s2, s3, s4 uint32) uint32 {
	return s2 ^ s3 ^ s4 ^ 0x389cbd17
}

func outerTableLane11Mask(s2, s0, s3, s1 uint32) uint32 {
	return outerBooleanMask([]uint32{s2, s0, s3, s1}, outerTableLane11Terms[:])
}

func outerTableLane12Mask(s1, s0, s3, s4 uint32) uint32 {
	return outerBooleanMask([]uint32{s1, s0, s3, s4}, outerTableLane12Terms[:])
}

func outerTableLane13Mask(s3, s2, s4 uint32) uint32 {
	return outerBooleanMask([]uint32{s3, s2, s4}, outerTableLane13Terms[:])
}

func outerTableLane14Mask(s0, s2, s3, s4 uint32) uint32 {
	return outerBooleanMask([]uint32{s0, s2, s3, s4}, outerTableLane14Terms[:])
}

func outerTableLane15Mask(a, b, c, d uint32) uint32 {
	return ^(a ^ b ^ c ^ d)
}

func outerTableLane16Mask(s1, s2, s3 uint32) uint32 {
	return s1 ^ s2 ^ s3 ^ 0xeb1c4c71
}

func outerTableLane17Mask(s0, s1, s3 uint32) uint32 {
	return s0 ^ s1 ^ s3 ^ 0x5d1651a0
}

func outerTableLane20Mask(pendingS0, s1, lateMix, s3 uint32) uint32 {
	return outerBooleanMask([]uint32{pendingS0, s1, lateMix, s3}, outerTableLane20Terms[:])
}
