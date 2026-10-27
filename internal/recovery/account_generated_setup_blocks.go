package recovery

const (
	accountCallFrameSetup     = uint64(0x1001c3d10)
	accountRecordDecodeSetup  = uint64(0x10023733c)
	accountEntryWidePrologue  = uint64(0x100375f30)
	accountEntryWideSetup     = uint64(0x100376000)
	accountBooleanPermutation = uint64(0x10093d024)
)

func runAccountCallFrameSetup(s *state) (uint64, error) {
	s.traceBlock(0x1001c3d10)
	var next uint64

	s.x[11] = uint64(uint32(uint64(0x4ae)))

	s.x[11] = uint64(uint32(insertBits(uint64(uint32(s.x[11])), 0xe4ae, 16, 16, 32)))

	s.x[11] = uint64(uint32((uint64(uint32(s.x[28])) + uint64(uint32(s.x[11])))))

	s.x[12] = uint64((s.x[30] + uint64(0x48)))

	s.x[0] = uint64((s.x[12] + uint64(0x8)))

	value5, err := s.load((s.x[27] + shiftLeft(signExtend32(uint32(uint64(uint32(s.x[11])))), 3, 64)), 8)
	if err != nil {
		return 0, err
	}
	s.x[11] = uint64(value5)

	s.x[13] = uint64(uint32(uint64(0xffd)))

	s.x[13] = uint64(uint32(insertBits(uint64(uint32(s.x[13])), 0xe4ae, 16, 16, 32)))

	s.x[14] = uint64(uint32(uint64(0x2120)))

	s.x[15] = uint64(uint32(uint64(0xf5d7)))

	s.x[15] = uint64(uint32(insertBits(uint64(uint32(s.x[15])), 0xe4ad, 16, 16, 32)))

	s.x[2] = uint64(uint32(uint64(0xcfdb)))

	s.x[24] = uint64(uint32(uint64(0x1b)))

	s.x[25] = uint64(uint32(uint64(0xfaef)))

	s.x[25] = uint64(uint32(insertBits(uint64(uint32(s.x[25])), 0xffdf, 16, 16, 32)))

	s.x[23] = uint64(uint32(uint64(0xf5de)))

	s.x[23] = uint64(uint32(insertBits(uint64(uint32(s.x[23])), 0xffbf, 16, 16, 32)))

	s.x[22] = uint64(uint32(uint64(0xffffdec3)))

	s.x[20] = uint64(uint32(uint64(0x2140)))

	s.x[21] = uint64(uint32(uint64(0x7d30)))

	s.x[19] = uint64(uint32(uint64(0xffffdea1)))

	s.x[26] = uint64(uint32(uint64(0x84c6)))

	s.x[16] = uint64(uint32(uint64(0x26ef)))

	s.x[12] = uint64(uint32(uint64(0xcb8)))

	s.x[12] = uint64(uint32(insertBits(uint64(uint32(s.x[12])), 0x437f, 16, 16, 32)))

	s.x[12] = uint64(uint32((uint64(uint32(s.x[16])) + uint64(uint32(s.x[12])))))

	if err := s.storeValue(addOffset(s.sp, 132), uint64(uint32(s.x[12])), 4); err != nil {
		return 0, err
	}

	s.x[12] = uint64(uint32(uint64(0xb4a4)))

	s.x[12] = uint64(uint32(insertBits(uint64(uint32(s.x[12])), 0xc028, 16, 16, 32)))

	s.x[16] = uint64(uint32((uint64(uint32(s.x[16])) + uint64(uint32(s.x[12])))))

	s.x[17] = uint64(uint32(uint64(0x8237)))

	s.x[17] = uint64(uint32(insertBits(uint64(uint32(s.x[17])), 0x7f34, 16, 16, 32)))

	s.x[12] = uint64(uint32(uint64(0x81f5)))

	s.x[12] = uint64(uint32(insertBits(uint64(uint32(s.x[12])), 0x3f34, 16, 16, 32)))

	s.x[12] = uint64(uint32((uint64(uint32(s.x[17])) ^ uint64(uint32(s.x[12])))))

	if err := s.storeValue(addOffset(s.sp, 160), uint64(uint32(s.x[12])), 4); err != nil {
		return 0, err
	}

	s.x[12] = uint64(uint32(uint64(0xc231)))

	s.x[12] = uint64(uint32(insertBits(uint64(uint32(s.x[12])), 0x80cb, 16, 16, 32)))

	s.x[12] = uint64(uint32((uint64(uint32(s.x[17])) + uint64(uint32(s.x[12])))))

	if err := s.storeValue(addOffset(s.sp, 156), uint64(uint32(s.x[12])), 4); err != nil {
		return 0, err
	}

	s.x[12] = uint64((s.x[27] + shiftLeft(uint64(0x13), 12, 64)))

	s.x[12] = uint64((s.x[12] + uint64(0x778)))

	s.x[17] = uint64(uint64(0xc48a))

	next = s.x[11]
	return next, nil
}

func runAccountRecordDecodeSetup(s *state) (uint64, error) {
	s.traceBlock(0x10023733c)
	var next uint64

	s.x[9] = uint64(uint32(uint64(0x6932)))

	s.x[9] = uint64(uint32((uint64(uint32(s.x[8])) ^ uint64(uint32(s.x[9])))))

	value2, err := s.load((s.x[27] + shiftLeft(signExtend32(uint32(uint64(uint32(s.x[9])))), 3, 64)), 8)
	if err != nil {
		return 0, err
	}
	s.x[11] = uint64(value2)

	value3, err := s.load(s.x[5], 4)
	if err != nil {
		return 0, err
	}
	s.x[9] = uint64(uint32(value3))

	value4, err := s.load(addOffset(s.sp, 96), 8)
	if err != nil {
		return 0, err
	}
	s.x[12] = uint64(value4)

	value5, err := s.load(s.x[12], 4)
	if err != nil {
		return 0, err
	}
	s.x[12] = uint64(uint32(value5))

	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[12])))))

	s.x[12] = uint64(uint32(uint64(0xa6)))

	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[12])))))

	s.x[12] = uint64(uint32(uint64(0xd4cf)))

	s.x[12] = uint64(uint32(insertBits(uint64(uint32(s.x[12])), 0xa5c, 16, 16, 32)))

	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) * uint64(uint32(s.x[12])))))

	value12, err := s.load(addOffset(s.sp, 80), 8)
	if err != nil {
		return 0, err
	}
	s.x[12] = uint64(value12)

	s.x[12] = uint64((s.x[12] + uint64(uint8(uint64(uint32(s.x[9]))))))

	value14, err := s.load(s.x[12], 1)
	if err != nil {
		return 0, err
	}
	s.x[12] = uint64(uint32(value14))

	s.x[13] = uint64(uint32(uint64(0x2a)))

	s.x[12] = uint64(uint32((uint64(uint32(s.x[12])) ^ uint64(uint32(s.x[13])))))

	value17, err := s.load(addOffset(s.sp, 72), 8)
	if err != nil {
		return 0, err
	}
	s.x[13] = uint64(value17)

	value18, err := s.load((s.x[13] + uint64(uint32(uint64(uint32(s.x[12]))))), 1)
	if err != nil {
		return 0, err
	}
	s.x[12] = uint64(uint32(value18))

	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[12])))))

	s.x[12] = uint64(uint32(uint64(0x13da)))

	s.x[9] = uint64(uint32((uint64(uint32(s.x[12])) + uint64(uint8(uint64(uint32(s.x[9])))))))

	value22, err := s.load((s.x[4] + shiftLeft(uint64(uint32(uint64(uint32(s.x[9])))), 3, 64)), 8)
	if err != nil {
		return 0, err
	}
	s.x[9] = uint64(value22)

	s.x[12] = uint64(uint64(0xffffffffffff6f5a))

	s.x[12] = uint64(insertBits(s.x[12], 0xc8cc, 16, 16, 64))

	s.x[9] = uint64((s.x[9] + s.x[12]))

	if err := s.storeValue(addOffset(s.sp, 48), uint64(uint32(s.x[10])), 4); err != nil {
		return 0, err
	}

	s.x[10] = uint64(uint32(uint64(0x3d3)))

	s.x[10] = uint64(uint32(insertBits(uint64(uint32(s.x[10])), 0x663d, 16, 16, 32)))

	s.x[13] = uint64(uint32(uint64(0xffd)))

	s.x[13] = uint64(uint32(insertBits(uint64(uint32(s.x[13])), 0xe4ae, 16, 16, 32)))

	s.x[14] = uint64(uint32(uint64(0x2120)))

	s.x[15] = uint64(uint32(uint64(0xf5d7)))

	s.x[15] = uint64(uint32(insertBits(uint64(uint32(s.x[15])), 0xe4ad, 16, 16, 32)))

	value34, err := s.load(addOffset(s.sp, 108), 4)
	if err != nil {
		return 0, err
	}
	s.x[16] = uint64(uint32(value34))

	value35, err := s.load(addOffset(s.sp, 112), 8)
	if err != nil {
		return 0, err
	}
	s.x[12] = uint64(value35)

	s.x[17] = uint64(uint64(0xc48a))

	next = s.x[11]
	return next, nil
}

func runAccountEntryWidePrologue(s *state) (uint64, error) {
	s.traceBlock(0x100375f30)

	s.sp = uint64((s.sp - uint64(0x170)))

	if err := s.storeValue(addOffset(s.sp, 272), s.x[28], 8); err != nil {
		return 0, err
	}
	if err := s.storeValue((addOffset(s.sp, 272))+8, s.x[27], 8); err != nil {
		return 0, err
	}

	if err := s.storeValue(addOffset(s.sp, 288), s.x[26], 8); err != nil {
		return 0, err
	}
	if err := s.storeValue((addOffset(s.sp, 288))+8, s.x[25], 8); err != nil {
		return 0, err
	}

	if err := s.storeValue(addOffset(s.sp, 304), s.x[24], 8); err != nil {
		return 0, err
	}
	if err := s.storeValue((addOffset(s.sp, 304))+8, s.x[23], 8); err != nil {
		return 0, err
	}

	if err := s.storeValue(addOffset(s.sp, 320), s.x[22], 8); err != nil {
		return 0, err
	}
	if err := s.storeValue((addOffset(s.sp, 320))+8, s.x[21], 8); err != nil {
		return 0, err
	}

	if err := s.storeValue(addOffset(s.sp, 336), s.x[20], 8); err != nil {
		return 0, err
	}
	if err := s.storeValue((addOffset(s.sp, 336))+8, s.x[19], 8); err != nil {
		return 0, err
	}

	if err := s.storeValue(addOffset(s.sp, 352), s.x[29], 8); err != nil {
		return 0, err
	}
	if err := s.storeValue((addOffset(s.sp, 352))+8, s.x[30], 8); err != nil {
		return 0, err
	}

	s.x[29] = uint64((s.sp + uint64(0x160)))

	s.x[30] = uint64(s.x[0])

	s.x[8] = uint64(0x10106d000)

	value11, err := s.load(addOffset(s.x[8], 1032), 8)
	if err != nil {
		return 0, err
	}
	s.x[8] = uint64(value11)

	value12, err := s.load(s.x[8], 8)
	if err != nil {
		return 0, err
	}
	s.x[8] = uint64(value12)

	if err := s.storeValue(addOffset(s.x[29], -96), s.x[8], 8); err != nil {
		return 0, err
	}

	s.x[12] = uint64(0x101080000)

	value15, err := s.load(addOffset(s.x[12], 3776), 4)
	if err != nil {
		return 0, err
	}
	s.x[8] = uint64(uint32(value15))

	s.x[3] = uint64(0x101075000)

	value17, err := s.load(addOffset(s.x[3], 2880), 4)
	if err != nil {
		return 0, err
	}
	s.x[13] = uint64(uint32(value17))

	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) ^ uint64(uint32(s.x[13])))))

	s.x[9] = uint64(uint32(uint64(0xa6)))

	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) ^ uint64(uint32(s.x[9])))))

	s.x[11] = uint64(uint32(uint64(0xd4cf)))

	s.x[11] = uint64(uint32(insertBits(uint64(uint32(s.x[11])), 0xa5c, 16, 16, 32)))

	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) * uint64(uint32(s.x[11])))))

	s.x[10] = uint64(0x100f49000)

	s.x[10] = uint64((s.x[10] + uint64(0x830)))

	s.x[10] = uint64((s.x[10] + uint64(uint8(uint64(uint32(s.x[8]))))))

	value27, err := s.load(s.x[10], 1)
	if err != nil {
		return 0, err
	}
	s.x[14] = uint64(uint32(value27))

	s.x[10] = uint64(uint32(uint64(0x9a)))

	s.x[14] = uint64(uint32((uint64(uint32(s.x[14])) ^ uint64(uint32(s.x[10])))))

	s.x[15] = uint64(0x100fe9000)

	s.x[15] = uint64((s.x[15] + uint64(0x690)))

	value32, err := s.load((s.x[15] + uint64(uint32(uint64(uint32(s.x[14]))))), 1)
	if err != nil {
		return 0, err
	}
	s.x[14] = uint64(uint32(value32))

	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) ^ uint64(uint32(s.x[14])))))

	s.x[14] = uint64(uint32(uint64(0x14fa)))

	s.x[4] = uint64(0x10105b000)

	s.x[4] = uint64((s.x[4] + uint64(0x6f0)))

	s.x[8] = uint64(uint32((uint64(uint32(s.x[14])) + uint64(uint8(uint64(uint32(s.x[8])))))))

	value38, err := s.load((s.x[4] + shiftLeft(uint64(uint32(uint64(uint32(s.x[8])))), 3, 64)), 8)
	if err != nil {
		return 0, err
	}
	s.x[8] = uint64(value38)

	s.x[8] = uint64((s.x[8] - shiftLeft(uint64(0x1d6), 12, 64)))

	s.x[5] = uint64((s.x[8] - uint64(0xdc3)))

	value41, err := s.load(s.x[5], 8)
	if err != nil {
		return 0, err
	}
	s.x[8] = uint64(value41)

	s.x[13] = uint64(uint32((uint64(uint32(s.x[13])) ^ uint64(uint32(s.x[9])))))

	s.x[13] = uint64(uint32((uint64(uint32(s.x[8])) ^ uint64(uint32(s.x[13])))))

	s.x[13] = uint64(uint32((uint64(uint32(s.x[13])) * uint64(uint32(s.x[11])))))

	value45, err := s.load(addOffset(s.x[4], 30544), 8)
	if err != nil {
		return 0, err
	}
	s.x[14] = uint64(value45)

	s.x[15] = uint64(uint64(0xffffffffffff2dad))

	s.x[15] = uint64(insertBits(s.x[15], 0xed6f, 16, 16, 64))

	s.x[6] = uint64((s.x[14] + s.x[15]))

	s.x[14] = uint64((s.x[6] + uint64(uint8(uint64(uint32(s.x[13]))))))

	value50, err := s.load(s.x[14], 1)
	if err != nil {
		return 0, err
	}
	s.x[14] = uint64(uint32(value50))

	s.x[14] = uint64(uint32((uint64(uint32(s.x[14])) ^ uint64(0x10))))
	return 0x100376000, nil
}

func runAccountEntryWideSetup(s *state) (uint64, error) {
	s.traceBlock(0x100376000)
	var next uint64

	s.x[15] = uint64(uint32(uint64(0xa078)))

	s.x[15] = uint64((s.x[4] + s.x[15]))

	value2, err := s.load(s.x[15], 8)
	if err != nil {
		return 0, err
	}
	s.x[15] = uint64(value2)

	s.x[16] = uint64(uint64(0xfffffffffffff512))

	s.x[16] = uint64(insertBits(s.x[16], 0x997c, 16, 16, 64))

	s.x[7] = uint64((s.x[15] + s.x[16]))

	value6, err := s.load((s.x[7] + uint64(uint32(uint64(uint32(s.x[14]))))), 1)
	if err != nil {
		return 0, err
	}
	s.x[14] = uint64(uint32(value6))

	s.x[13] = uint64(uint32((uint64(uint32(s.x[14])) ^ uint64(uint32(s.x[13])))))

	s.x[14] = uint64(uint32(uint64(0x131e)))

	s.x[13] = uint64(uint32((uint64(uint32(s.x[14])) + uint64(uint8(uint64(uint32(s.x[13])))))))

	value10, err := s.load((s.x[4] + shiftLeft(uint64(uint32(uint64(uint32(s.x[13])))), 3, 64)), 8)
	if err != nil {
		return 0, err
	}
	s.x[13] = uint64(value10)

	s.x[14] = uint64(uint64(0xffffffffffffa825))

	s.x[14] = uint64(insertBits(s.x[14], 0x9608, 16, 16, 64))

	s.x[15] = uint64((s.x[13] + s.x[14]))

	value14, err := s.load(s.x[15], 8)
	if err != nil {
		return 0, err
	}
	s.x[13] = uint64(value14)

	s.x[14] = uint64((s.x[29] - uint64(0xa8)))

	s.x[8] = uint64((s.x[13] + s.x[8]))

	s.x[8] = uint64((s.x[8] ^ s.x[14]))

	s.x[13] = uint64(uint64(0xd1a6))

	s.x[13] = uint64(insertBits(s.x[13], 0x8eae, 16, 16, 64))

	s.x[13] = uint64(insertBits(s.x[13], 0x4620, 32, 16, 64))

	s.x[13] = uint64(insertBits(s.x[13], 0x73d0, 48, 16, 64))

	s.x[14] = uint64((s.x[8]*s.x[11] + s.x[13]))

	s.x[8] = uint64((s.x[8] ^ s.x[13]))

	s.x[8] = uint64((s.x[8] * s.x[11]))

	if err := s.storeValue(s.x[5], s.x[14], 8); err != nil {
		return 0, err
	}

	if err := s.storeValue(addOffset(s.sp, 96), s.x[15], 8); err != nil {
		return 0, err
	}

	if err := s.storeValue(s.x[15], s.x[8], 8); err != nil {
		return 0, err
	}

	value28, err := s.load(addOffset(s.x[0], 72), 4)
	if err != nil {
		return 0, err
	}
	s.x[8] = uint64(uint32(value28))

	s.x[13] = uint64(uint32(uint64(0x7bc6)))

	s.x[13] = uint64(uint32(insertBits(uint64(uint32(s.x[13])), 0xce75, 16, 16, 32)))

	s.compare32(uint32(uint64(uint32(s.x[8]))), uint32(uint64(uint32(s.x[13]))))

	s.x[13] = uint64(uint32(boolValue(s.condition("eq"))))

	s.x[14] = uint64(uint32(uint64(0x39)))

	s.x[13] = uint64(uint32((uint64(uint32(s.x[13])) * uint64(uint32(s.x[14])))))

	value35, err := s.load(addOffset(s.x[12], 3776), 4)
	if err != nil {
		return 0, err
	}
	s.x[12] = uint64(uint32(value35))

	value36, err := s.load(addOffset(s.x[3], 2880), 4)
	if err != nil {
		return 0, err
	}
	s.x[14] = uint64(uint32(value36))

	s.x[12] = uint64(uint32((uint64(uint32(s.x[12])) ^ uint64(uint32(s.x[14])))))

	s.x[12] = uint64(uint32((uint64(uint32(s.x[12])) ^ uint64(uint32(s.x[9])))))

	s.x[12] = uint64(uint32((uint64(uint32(s.x[12])) * uint64(uint32(s.x[11])))))

	s.x[14] = uint64(0x100fe9000)

	s.x[14] = uint64((s.x[14] + uint64(0x190)))

	s.x[14] = uint64((s.x[14] + uint64(uint8(uint64(uint32(s.x[12]))))))

	value43, err := s.load(s.x[14], 1)
	if err != nil {
		return 0, err
	}
	s.x[14] = uint64(uint32(value43))

	s.x[15] = uint64(uint32(uint64(0x29)))

	s.x[14] = uint64(uint32((uint64(uint32(s.x[14])) ^ uint64(uint32(s.x[15])))))

	s.x[15] = uint64(0x100de2000)

	s.x[15] = uint64((s.x[15] + uint64(0x8f0)))

	value48, err := s.load((s.x[15] + uint64(uint32(uint64(uint32(s.x[14]))))), 1)
	if err != nil {
		return 0, err
	}
	s.x[14] = uint64(uint32(value48))

	s.x[15] = uint64(uint32(uint64(0x271b)))

	s.x[12] = uint64(uint32((uint64(uint32(s.x[12])) ^ uint64(uint32(s.x[14])))))

	s.x[12] = uint64(uint32((uint64(uint32(s.x[12])) & uint64(0xff))))

	s.x[12] = uint64(uint32((uint64(uint32(s.x[12])) + uint64(0x1c9))))

	value53, err := s.load((s.x[4] + shiftLeft(uint64(uint32(uint64(uint32(s.x[12])))), 3, 64)), 8)
	if err != nil {
		return 0, err
	}
	s.x[12] = uint64(value53)

	s.x[13] = uint64(uint32((uint64(uint32(s.x[13])) ^ uint64(uint32(s.x[15])))))

	s.x[14] = uint64(uint64(0xffffffffffff5355))

	s.x[14] = uint64(insertBits(s.x[14], 0x9824, 16, 16, 64))

	s.x[27] = uint64((s.x[12] + s.x[14]))

	value58, err := s.load((s.x[27] + shiftLeft(uint64(uint32(uint64(uint32(s.x[13])))), 3, 64)), 8)
	if err != nil {
		return 0, err
	}
	s.x[12] = uint64(value58)

	value59, err := s.load(addOffset(s.x[4], 24096), 8)
	if err != nil {
		return 0, err
	}
	s.x[13] = uint64(value59)

	s.x[13] = uint64((s.x[13] - shiftLeft(uint64(0x9fd), 12, 64)))

	s.x[15] = uint64((s.x[13] - uint64(0xcfb)))

	s.x[13] = uint64(uint32(uint64(0x9898)))

	s.x[13] = uint64((s.x[4] + s.x[13]))

	value64, err := s.load(s.x[13], 8)
	if err != nil {
		return 0, err
	}
	s.x[13] = uint64(value64)

	s.x[14] = uint64(uint64(0xffffffffffffc916))

	s.x[14] = uint64(insertBits(s.x[14], 0xf3e6, 16, 16, 64))

	s.x[13] = uint64((s.x[13] + s.x[14]))

	if err := s.storeValue(addOffset(s.sp, 136), s.x[13], 8); err != nil {
		return 0, err
	}
	if err := s.storeValue((addOffset(s.sp, 136))+8, s.x[15], 8); err != nil {
		return 0, err
	}

	s.x[13] = uint64(uint32(uint64(0x8d70)))

	s.x[13] = uint64((s.x[4] + s.x[13]))

	value71, err := s.load(s.x[13], 8)
	if err != nil {
		return 0, err
	}
	s.x[13] = uint64(value71)

	s.x[14] = uint64(uint64(0xffffffffffffb571))

	s.x[14] = uint64(insertBits(s.x[14], 0x93fc, 16, 16, 64))

	s.x[13] = uint64((s.x[13] + s.x[14]))

	if err := s.storeValue(addOffset(s.sp, 80), s.x[13], 8); err != nil {
		return 0, err
	}

	value76, err := s.load(addOffset(s.x[4], 28432), 8)
	if err != nil {
		return 0, err
	}
	s.x[13] = uint64(value76)

	s.x[14] = uint64(uint64(0xffffffffffffc2a5))

	s.x[14] = uint64(insertBits(s.x[14], 0x954b, 16, 16, 64))

	s.x[13] = uint64((s.x[13] + s.x[14]))

	if err := s.storeValue(addOffset(s.sp, 72), s.x[13], 8); err != nil {
		return 0, err
	}

	value81, err := s.load(addOffset(s.x[4], 6904), 8)
	if err != nil {
		return 0, err
	}
	s.x[13] = uint64(value81)

	s.x[14] = uint64(uint64(0xffffffffffff4862))

	s.x[14] = uint64(insertBits(s.x[14], 0xcd11, 16, 16, 64))

	s.x[13] = uint64((s.x[13] + s.x[14]))

	if err := s.storeValue(addOffset(s.sp, 64), s.x[13], 8); err != nil {
		return 0, err
	}

	value86, err := s.load(addOffset(s.x[4], 10840), 8)
	if err != nil {
		return 0, err
	}
	s.x[13] = uint64(value86)

	s.x[14] = uint64(uint64(0xffffffffffff614e))

	s.x[14] = uint64(insertBits(s.x[14], 0xe5b0, 16, 16, 64))

	s.x[13] = uint64((s.x[13] + s.x[14]))

	if err := s.storeValue(addOffset(s.sp, 56), s.x[13], 8); err != nil {
		return 0, err
	}

	s.x[25] = uint64(uint32(uint64(0x7bc5)))

	s.x[25] = uint64(uint32(insertBits(uint64(uint32(s.x[25])), 0xce75, 16, 16, 32)))

	s.x[13] = uint64(uint32(uint64(0x1177)))

	s.x[15] = uint64(uint32(uint64(0x7bc5)))

	s.x[15] = uint64(uint32(insertBits(uint64(uint32(s.x[15])), 0xce75, 16, 16, 32)))

	if err := s.storeValue(addOffset(s.sp, 24), s.x[5], 8); err != nil {
		return 0, err
	}
	if err := s.storeValue((addOffset(s.sp, 24))+8, s.x[1], 8); err != nil {
		return 0, err
	}

	if err := s.storeValue(addOffset(s.sp, 16), s.x[6], 8); err != nil {
		return 0, err
	}

	if err := s.storeValue(addOffset(s.sp, 120), s.x[7], 8); err != nil {
		return 0, err
	}

	s.x[17] = uint64(uint64(0xfed4))

	next = s.x[12]
	return next, nil
}

func runAccountBooleanPermutation(s *state) (uint64, error) {
	s.traceBlock(0x10093d024)
	var next uint64

	s.x[14] = uint64(uint32(uint64(0xba26)))

	s.x[14] = uint64(uint32(insertBits(uint64(uint32(s.x[14])), 0x4187, 16, 16, 32)))

	s.compare32(uint32(uint64(uint32(s.x[9]))), uint32(uint64(uint32(s.x[14]))))

	s.x[13] = uint64(uint32(boolValue(s.condition("lo"))))

	s.x[12] = uint64(uint32(uint64(0xffffcb31)))

	s.x[16] = uint64(uint32(uint64(0x97b6)))

	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0x29d0, 16, 16, 32)))

	s.x[16] = uint64(uint32((uint64(uint32(s.x[15])) ^ uint64(uint32(s.x[16])))))

	s.x[17] = uint64(uint32(uint64(0x684a)))

	s.x[17] = uint64(uint32(insertBits(uint64(uint32(s.x[17])), 0xd62f, 16, 16, 32)))

	s.x[16] = uint64(uint32((uint64(uint32(s.x[16])) + uint64(uint32(s.x[17])))))

	s.x[17] = uint64(uint32(uint64(0xc27d)))

	s.x[17] = uint64(uint32(insertBits(uint64(uint32(s.x[17])), 0xf2ab, 16, 16, 32)))

	s.x[17] = uint64(uint32((uint64(uint32(s.x[15])) ^ uint64(uint32(s.x[17])))))

	s.x[0] = uint64(uint32(uint64(0x3d83)))

	s.x[0] = uint64(uint32(insertBits(uint64(uint32(s.x[0])), 0xd54, 16, 16, 32)))

	s.x[17] = uint64(uint32((uint64(uint32(s.x[17])) + uint64(uint32(s.x[0])))))

	s.x[16] = uint64(uint32((uint64(uint32(s.x[16])) ^ uint64(uint32(s.x[17])))))

	s.x[17] = uint64(uint32(uint64(0x32d0)))

	s.x[17] = uint64(uint32(insertBits(uint64(uint32(s.x[17])), 0x150e, 16, 16, 32)))

	s.x[15] = uint64(uint32((uint64(uint32(s.x[15])) ^ uint64(uint32(s.x[17])))))

	s.x[12] = uint64(uint32((uint64(uint32(s.x[10])) + uint64(uint32(s.x[12])))))

	s.x[15] = uint64(uint32((uint64(uint32(s.x[15])) ^ uint64(uint32(s.x[12])))))

	s.x[17] = uint64(uint32(uint64(0xd1f2)))

	s.x[17] = uint64(uint32(insertBits(uint64(uint32(s.x[17])), 0xeaf1, 16, 16, 32)))

	s.x[15] = uint64(uint32((uint64(uint32(s.x[15])) + uint64(uint32(s.x[17])))))

	s.x[15] = uint64(uint32((uint64(uint32(s.x[16])) ^ uint64(uint32(s.x[15])))))

	s.x[16] = uint64(uint32(uint64(0x5a10)))

	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0xf3cf, 16, 16, 32)))

	s.x[17] = uint64(uint32(uint64(0x8941)))

	s.x[17] = uint64(uint32(insertBits(uint64(uint32(s.x[17])), 0x7263, 16, 16, 32)))

	s.x[15] = uint64(uint32((uint64(uint32(s.x[15])) + uint64(uint32(s.x[16])))))

	s.x[16] = uint64(uint32((uint64(uint32(s.x[15])) ^ uint64(uint32(s.x[17])))))

	s.x[17] = uint64(uint32(uint64(0x21b6)))

	s.x[17] = uint64(uint32(insertBits(uint64(uint32(s.x[17])), 0xdaa6, 16, 16, 32)))

	s.x[17] = uint64(uint32((uint64(uint32(s.x[15])) & uint64(uint32(s.x[17])))))

	s.x[0] = uint64(uint32(uint64(0x8925)))

	s.x[0] = uint64(uint32(insertBits(uint64(uint32(s.x[0])), 0xd227, 16, 16, 32)))

	s.x[0] = uint64(uint32((uint64(uint32(s.x[15])) ^ uint64(uint32(s.x[0])))))

	s.x[0] = uint64(uint32((uint64(uint32(s.x[0])) & shiftLeft(uint64(uint32(s.x[17])), 1, 32))))

	s.x[17] = uint64(uint32((uint64(uint32(s.x[0])) ^ uint64(uint32(s.x[17])))))

	s.x[0] = uint64(uint32(shiftLeft(uint64(uint32(s.x[16])), int(uint64(0x1)), 32)))

	s.x[1] = uint64(uint32(uint64(0x51ee)))

	s.x[1] = uint64(uint32(insertBits(uint64(uint32(s.x[1])), 0x518b, 16, 16, 32)))

	s.x[1] = uint64(uint32((uint64(uint32(s.x[0])) ^ uint64(uint32(s.x[1])))))

	s.x[16] = uint64(uint32((uint64(uint32(s.x[1])) & uint64(uint32(s.x[16])))))

	s.x[1] = uint64(uint32(uint64(0xa8f6)))

	s.x[1] = uint64(uint32(insertBits(uint64(uint32(s.x[1])), 0xa8c5, 16, 16, 32)))

	s.x[0] = uint64(uint32((uint64(uint32(s.x[0])) & uint64(uint32(s.x[1])))))

	s.x[16] = uint64(uint32((uint64(uint32(s.x[16])) ^ uint64(uint32(s.x[0])))))

	s.x[0] = uint64(uint32(uint64(0xa811)))

	s.x[0] = uint64(uint32(insertBits(uint64(uint32(s.x[0])), 0xa844, 16, 16, 32)))

	s.x[0] = uint64(uint32((uint64(uint32(s.x[16])) ^ uint64(uint32(s.x[0])))))

	s.x[1] = uint64(uint32(uint64(0xe0)))

	s.x[1] = uint64(uint32(insertBits(uint64(uint32(s.x[1])), 0x81, 16, 16, 32)))

	s.x[16] = uint64(uint32((uint64(uint32(s.x[16])) ^ uint64(uint32(s.x[1])))))

	s.x[16] = uint64(uint32((uint64(uint32(s.x[16])) & shiftLeft(uint64(uint32(s.x[17])), 2, 32))))

	s.x[16] = uint64(uint32((uint64(uint32(s.x[16])) ^ uint64(uint32(s.x[17])))))

	s.x[17] = uint64(uint32(shiftLeft(uint64(uint32(s.x[0])), int(uint64(0x2)), 32)))

	s.x[1] = uint64(uint32(uint64(0xa3dc)))

	s.x[1] = uint64(uint32(insertBits(uint64(uint32(s.x[1])), 0xa316, 16, 16, 32)))

	s.x[1] = uint64(uint32((uint64(uint32(s.x[17])) ^ uint64(uint32(s.x[1])))))

	s.x[0] = uint64(uint32((uint64(uint32(s.x[1])) & uint64(uint32(s.x[0])))))

	s.x[1] = uint64(uint32(uint64(0xa8f4)))

	s.x[1] = uint64(uint32(insertBits(uint64(uint32(s.x[1])), 0xa8c5, 16, 16, 32)))

	s.x[17] = uint64(uint32((uint64(uint32(s.x[17])) & uint64(uint32(s.x[1])))))

	s.x[17] = uint64(uint32((uint64(uint32(s.x[0])) ^ uint64(uint32(s.x[17])))))

	s.x[0] = uint64(uint32(uint64(0x823)))

	s.x[0] = uint64(uint32(insertBits(uint64(uint32(s.x[0])), 0x8c1, 16, 16, 32)))

	s.x[0] = uint64(uint32((uint64(uint32(s.x[17])) ^ uint64(uint32(s.x[0])))))

	s.x[1] = uint64(uint32(uint64(0xa0c0)))

	s.x[1] = uint64(uint32(insertBits(uint64(uint32(s.x[1])), 0xa004, 16, 16, 32)))

	s.x[17] = uint64(uint32((uint64(uint32(s.x[17])) ^ uint64(uint32(s.x[1])))))

	s.x[17] = uint64(uint32((uint64(uint32(s.x[17])) & shiftLeft(uint64(uint32(s.x[16])), 4, 32))))

	s.x[16] = uint64(uint32((uint64(uint32(s.x[17])) ^ uint64(uint32(s.x[16])))))

	s.x[17] = uint64(uint32(shiftLeft(uint64(uint32(s.x[0])), int(uint64(0x4)), 32)))

	s.x[1] = uint64(uint32(uint64(0x8f70)))

	s.x[1] = uint64(uint32(insertBits(uint64(uint32(s.x[1])), 0x8c5a, 16, 16, 32)))

	s.x[1] = uint64(uint32((uint64(uint32(s.x[17])) ^ uint64(uint32(s.x[1])))))

	s.x[0] = uint64(uint32((uint64(uint32(s.x[1])) & uint64(uint32(s.x[0])))))

	s.x[1] = uint64(uint32(uint64(0xa8f0)))

	s.x[1] = uint64(uint32(insertBits(uint64(uint32(s.x[1])), 0xa8c5, 16, 16, 32)))

	s.x[17] = uint64(uint32((uint64(uint32(s.x[17])) & uint64(uint32(s.x[1])))))

	s.x[17] = uint64(uint32((uint64(uint32(s.x[0])) ^ uint64(uint32(s.x[17])))))

	s.x[0] = uint64(uint32(uint64(0x2087)))

	s.x[0] = uint64(uint32(insertBits(uint64(uint32(s.x[0])), 0x2085, 16, 16, 32)))

	s.x[0] = uint64(uint32((uint64(uint32(s.x[17])) ^ uint64(uint32(s.x[0])))))

	s.x[1] = uint64(uint32(uint64(0x8800)))

	s.x[1] = uint64(uint32(insertBits(uint64(uint32(s.x[1])), 0x8840, 16, 16, 32)))

	s.x[17] = uint64(uint32((uint64(uint32(s.x[17])) ^ uint64(uint32(s.x[1])))))

	s.x[17] = uint64(uint32((uint64(uint32(s.x[17])) & shiftLeft(uint64(uint32(s.x[16])), 8, 32))))

	s.x[1] = uint64(uint32(uint64(0xa8f7)))

	s.x[1] = uint64(uint32(insertBits(uint64(uint32(s.x[1])), 0xa8c5, 16, 16, 32)))

	s.x[16] = uint64(uint32((uint64(uint32(s.x[16])) ^ uint64(uint32(s.x[1])))))

	s.x[16] = uint64(uint32((uint64(uint32(s.x[16])) ^ uint64(uint32(s.x[17])))))

	s.x[17] = uint64(uint32(shiftLeft(uint64(uint32(s.x[0])), int(uint64(0x8)), 32)))

	s.x[1] = uint64(uint32(uint64(0x45a80000)))

	s.x[1] = uint64(uint32((uint64(uint32(s.x[17])) ^ uint64(uint32(s.x[1])))))

	s.x[0] = uint64(uint32((uint64(uint32(s.x[1])) & uint64(uint32(s.x[0])))))

	s.x[1] = uint64(uint32(uint64(0x28c50000)))

	s.x[17] = uint64(uint32((uint64(uint32(s.x[17])) & uint64(uint32(s.x[1])))))

	s.x[2] = uint64(uint32(uint64(0x28450000)))

	s.x[17] = uint64(uint32((uint64(uint32(s.x[17])) ^ uint64(uint32(s.x[2])))))

	s.x[17] = uint64(uint32((uint64(uint32(s.x[17])) ^ uint64(uint32(s.x[0])))))

	s.x[0] = uint64(uint32(shiftLeft(uint64(uint32(s.x[16])), int(uint64(0x10)), 32)))

	s.x[2] = uint64(uint32(uint64(0x28f70000)))

	s.x[2] = uint64(uint32((uint64(uint32(s.x[0])) ^ uint64(uint32(s.x[2])))))

	s.x[17] = uint64(uint32((uint64(uint32(s.x[2])) & uint64(uint32(s.x[17])))))

	s.x[0] = uint64(uint32((uint64(uint32(s.x[0])) & uint64(uint32(s.x[1])))))

	s.x[16] = uint64(uint32((uint64(uint32(s.x[0])) ^ uint64(uint32(s.x[16])))))

	s.x[16] = uint64(uint32((uint64(uint32(s.x[16])) ^ uint64(uint32(s.x[17])))))

	s.x[16] = uint64(uint32((uint64(uint32(s.x[15])) ^ shiftLeft(uint64(uint32(s.x[16])), 1, 32))))

	s.x[15] = uint64(uint32(uint64(0xb9d)))

	s.x[15] = uint64(uint32(insertBits(uint64(uint32(s.x[15])), 0x14d2, 16, 16, 32)))

	s.x[15] = uint64(uint32((uint64(uint32(s.x[16])) ^ uint64(uint32(s.x[15])))))

	s.x[17] = uint64(uint32(uint64(0xa61d)))

	s.x[17] = uint64(uint32(insertBits(uint64(uint32(s.x[17])), 0x6b12, 16, 16, 32)))

	s.x[17] = uint64(uint32((uint64(uint32(s.x[16])) ^ uint64(uint32(s.x[17])))))

	s.x[0] = uint64(uint32(uint64(0x5280)))

	s.x[0] = uint64(uint32(insertBits(uint64(uint32(s.x[0])), 0x803f, 16, 16, 32)))

	s.x[17] = uint64(uint32((uint64(uint32(s.x[17])) + uint64(uint32(s.x[0])))))

	s.x[0] = uint64(uint32(uint64(0x998d)))

	s.x[0] = uint64(uint32(insertBits(uint64(uint32(s.x[0])), 0x6760, 16, 16, 32)))

	s.x[0] = uint64(uint32((uint64(uint32(s.x[16])) ^ uint64(uint32(s.x[0])))))

	s.x[1] = uint64(uint32(uint64(0x6df0)))

	s.x[1] = uint64(uint32(insertBits(uint64(uint32(s.x[1])), 0x8c4d, 16, 16, 32)))

	s.x[0] = uint64(uint32((uint64(uint32(s.x[0])) + uint64(uint32(s.x[1])))))

	value127, err := s.load(addOffset(s.sp, 608), 8)
	if err != nil {
		return 0, err
	}
	s.x[1] = uint64(value127)

	s.x[17] = uint64(uint32((uint64(uint32(s.x[17])) ^ uint64(uint32(s.x[0])))))

	s.x[0] = uint64(uint32(uint64(0x4fc8)))

	s.x[0] = uint64(uint32(insertBits(uint64(uint32(s.x[0])), 0xd6d5, 16, 16, 32)))

	s.x[16] = uint64(uint32((uint64(uint32(s.x[16])) ^ uint64(uint32(s.x[0])))))

	s.x[0] = uint64(uint32(uint64(0xbbab)))

	s.x[0] = uint64(uint32(insertBits(uint64(uint32(s.x[0])), 0x3df8, 16, 16, 32)))

	s.x[16] = uint64(uint32((uint64(uint32(s.x[16])) + uint64(uint32(s.x[0])))))

	s.x[16] = uint64(uint32((uint64(uint32(s.x[17])) ^ uint64(uint32(s.x[16])))))

	s.x[17] = uint64(uint32(uint64(0x35f3)))

	s.x[17] = uint64(uint32(insertBits(uint64(uint32(s.x[17])), 0xffd, 16, 16, 32)))

	s.x[16] = uint64(uint32((uint64(uint32(s.x[16])) + uint64(uint32(s.x[17])))))

	s.compare32(uint32(uint64(uint32(s.x[16]))), uint32(uint64(uint32(s.x[14]))))

	s.x[14] = uint64(uint32(boolValue(s.condition("lo"))))

	s.x[14] = uint64(uint32((uint64(uint32(s.x[13])) ^ uint64(uint32(s.x[14])))))

	s.compare32(uint32(uint64(uint32(s.x[16]))), uint32(uint64(uint32(s.x[9]))))

	s.x[16] = uint64(uint32(boolValue(s.condition("lo"))))

	s.compare32(uint32(uint64(uint32(s.x[14]))), uint32(uint64(0x0)))

	s.x[13] = uint64(uint32(selectValue(s.condition("ne"), uint64(uint32(s.x[13])), uint64(uint32(s.x[16])))))

	s.x[14] = uint64(uint32(uint64(0xc117)))

	s.x[13] = uint64(uint32((uint64(uint32(s.x[13])) * uint64(uint32(s.x[14])))))

	s.x[13] = uint64(uint32((uint64(uint32(s.x[13])) ^ uint64(uint32(s.x[10])))))

	value149, err := s.load((s.x[22] + shiftLeft(signExtend32(uint32(uint64(uint32(s.x[13])))), 3, 64)), 8)
	if err != nil {
		return 0, err
	}
	s.x[13] = uint64(value149)

	s.x[17] = uint64(uint64(0xc4a7))

	next = s.x[13]
	return next, nil
}
