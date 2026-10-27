package recovery

const (
	accountBlock1003aed78 = uint64(0x1003aed78)
	accountBlock1004bec84 = uint64(0x1004bec84)
	accountBlock1004df3b0 = uint64(0x1004df3b0)
	accountBlock10055579c = uint64(0x10055579c)
	accountBlock10092147c = uint64(0x10092147c)
	accountBlock1009e3a28 = uint64(0x1009e3a28)
	accountBlock100a0b4d4 = uint64(0x100a0b4d4)
	accountBlock100aba4f8 = uint64(0x100aba4f8)
	accountBlock100b5f734 = uint64(0x100b5f734)
	accountBlock1000f50f4 = uint64(0x1000f50f4)
	accountBlock1008d5268 = uint64(0x1008d5268)
	accountBlock1008e64c0 = uint64(0x1008e64c0)
	accountBlock1008f5e1c = uint64(0x1008f5e1c)
	accountBlock100900484 = uint64(0x100900484)
	accountBlock1009ecc88 = uint64(0x1009ecc88)
	accountBlock100a070f8 = uint64(0x100a070f8)
	accountBlock100a4aa04 = uint64(0x100a4aa04)
	accountBlock100a5bf1c = uint64(0x100a5bf1c)
	accountBlock100bed25c = uint64(0x100bed25c)
)

func runAccountBlock1003aed78(s *state) (uint64, error) {
	s.traceBlock(0x1003aed78)
	var next uint64

	s.x[8] = uint64(uint32(uint64(0x1e56)))

	s.x[8] = uint64(uint32(insertBits(uint64(uint32(s.x[8])), 0x1ecb, 16, 16, 32)))

	if err := s.storeValue(addOffset(s.x[0], 4), uint64(uint32(s.x[8])), 4); err != nil {
		return 0, err
	}

	s.sp = uint64((s.sp + uint64(0x10)))

	next = s.x[30]
	return next, nil
}

func runAccountBlock1004bec84(s *state) (uint64, error) {
	s.traceBlock(0x1004bec84)
	var next uint64

	s.x[11] = uint64(uint32(uint64(0xffffe1e3)))

	s.x[12] = uint64(uint32(uint64(0x6227)))

	s.x[11] = uint64(uint32((uint64(uint32(s.x[10])) + uint64(uint32(s.x[11])))))

	s.x[10] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[12])))))

	s.x[12] = uint64(uint32(uint64(0x2c4d)))

	s.x[12] = uint64(uint32(insertBits(uint64(uint32(s.x[12])), 0xc8cd, 16, 16, 32)))

	if err := s.storeValue(addOffset(s.x[8], 88), uint64(uint32(s.x[12])), 4); err != nil {
		return 0, err
	}

	value7, err := s.load((s.x[9] + shiftLeft(signExtend32(uint32(uint64(uint32(s.x[11])))), 3, 64)), 8)
	if err != nil {
		return 0, err
	}
	s.x[12] = uint64(value7)

	s.x[11] = uint64(uint32(uint64(0x14)))

	s.x[17] = uint64(uint64(0xbb44))

	next = s.x[12]
	return next, nil
}

func runAccountBlock1004df3b0(s *state) (uint64, error) {
	s.traceBlock(0x1004df3b0)
	var next uint64

	s.x[10] = uint64(uint32(uint64(0x7e3b)))

	s.x[10] = uint64(uint32(insertBits(uint64(uint32(s.x[10])), 0x5bfe, 16, 16, 32)))

	s.x[13] = uint64(uint32(uint64(0x7d7a)))

	s.x[13] = uint64(uint32(insertBits(uint64(uint32(s.x[13])), 0xa401, 16, 16, 32)))

	s.x[10] = uint64(uint32((uint64(uint32(s.x[11])) + uint64(uint32(s.x[10])))))

	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) & uint64(uint32(s.x[13])))))

	s.x[13] = uint64(uint32(signExtendBits(uint64(uint32(s.x[12])), 8, 32)))

	s.compare32(uint32(uint64(uint32(s.x[13]))), uint32(uint64(0x1)))

	s.x[13] = uint64(uint32(boolValue(s.condition("gt"))))

	s.x[14] = uint64(uint32(shiftLeft(uint64(uint32(s.x[13])), int(uint64(0x2)), 32)))

	s.x[13] = uint64(uint32((uint64(uint32(s.x[14])) | shiftLeft(uint64(uint32(s.x[13])), 3, 32))))

	s.x[11] = uint64(uint32((uint64(uint32(s.x[13])) ^ uint64(uint32(s.x[11])))))

	value12, err := s.load((s.x[9] + shiftLeft(signExtend32(uint32(uint64(uint32(s.x[11])))), 3, 64)), 8)
	if err != nil {
		return 0, err
	}
	s.x[11] = uint64(value12)

	s.x[17] = uint64(uint64(0x1c2a))

	next = s.x[11]
	return next, nil
}

func runAccountBlock10055579c(s *state) (uint64, error) {
	s.traceBlock(0x10055579c)
	var next uint64

	s.x[11] = uint64(uint32(uint64(0xffffe797)))

	s.x[13] = uint64(uint32(uint64(0x5a)))

	s.x[11] = uint64(uint32((uint64(uint32(s.x[10])) + uint64(uint32(s.x[11])))))

	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) | uint64(uint32(s.x[13])))))

	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) & uint64(0xff))))

	s.x[13] = uint64(uint32(uint64(0x59)))

	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[13])))))

	s.compare32(uint32(uint64(uint32(s.x[12]))), uint32(uint64(uint32(s.x[11]))))

	s.x[11] = uint64(uint32(boolValue(s.condition("eq"))))

	s.x[12] = uint64(uint32(uint64(0xf7d2)))

	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) * uint64(uint32(s.x[12])))))

	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[10])))))

	value12, err := s.load((s.x[9] + shiftLeft(signExtend32(uint32(uint64(uint32(s.x[11])))), 3, 64)), 8)
	if err != nil {
		return 0, err
	}
	s.x[11] = uint64(value12)

	s.x[17] = uint64(uint64(0x8c87))

	next = s.x[11]
	return next, nil
}

func runAccountBlock10092147c(s *state) (uint64, error) {
	s.traceBlock(0x10092147c)
	var next uint64

	s.x[8] = uint64(uint32(uint64(0x5369)))

	s.x[10] = uint64(uint32(uint64(0xa709)))

	s.x[8] = uint64(uint32((uint64(uint32(s.x[9])) + uint64(uint32(s.x[8])))))

	s.x[11] = uint64(uint32((uint64(uint32(s.x[8])) ^ uint64(uint32(s.x[10])))))

	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) - uint64(0xf40))))

	value5, err := s.load(addOffset(s.sp, 1396), 4)
	if err != nil {
		return 0, err
	}
	s.x[10] = uint64(uint32(value5))

	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) + uint64(0x8))))

	if err := s.storeValue(addOffset(s.sp, 1396), uint64(uint32(s.x[10])), 4); err != nil {
		return 0, err
	}

	if err := s.storeValue(addOffset(s.sp, 1392), uint64(uint32(s.x[19])), 4); err != nil {
		return 0, err
	}

	s.x[12] = uint64(uint32(uint64(0x41a7)))

	s.x[12] = uint64(uint32(insertBits(uint64(uint32(s.x[12])), 0xce75, 16, 16, 32)))

	s.x[12] = uint64(uint32((uint64(uint32(s.x[11])) + uint64(uint32(s.x[12])))))

	s.compare32(uint32(uint64(uint32(s.x[19]))), uint32(uint64(uint32(s.x[12]))))

	s.x[12] = uint64(uint32(boolValue(s.condition("eq"))))

	s.x[13] = uint64(uint32(uint64(0x86b5)))

	s.x[12] = uint64(uint32((uint64(uint32(s.x[12])) * uint64(uint32(s.x[13])))))

	s.x[9] = uint64(uint32((uint64(uint32(s.x[12])) ^ uint64(uint32(s.x[9])))))

	value17, err := s.load((s.x[22] + shiftLeft(signExtend32(uint32(uint64(uint32(s.x[9])))), 3, 64)), 8)
	if err != nil {
		return 0, err
	}
	s.x[9] = uint64(value17)

	s.x[17] = uint64(uint64(0xe25e))

	next = s.x[9]
	return next, nil
}

func runAccountBlock1009e3a28(s *state) (uint64, error) {
	s.traceBlock(0x1009e3a28)
	var next uint64

	s.x[8] = uint64(uint32(uint64(0xaaf9)))

	s.x[8] = uint64(uint32(insertBits(uint64(uint32(s.x[8])), 0x129, 16, 16, 32)))

	s.x[9] = uint64(uint32(uint64(0x7f37)))

	s.x[9] = uint64(uint32(insertBits(uint64(uint32(s.x[9])), 0xfed6, 16, 16, 32)))

	s.x[8] = uint64(uint32((uint64(uint32(s.x[11])) + uint64(uint32(s.x[8])))))

	s.x[9] = uint64(uint32((uint64(uint32(s.x[8])) & uint64(uint32(s.x[9])))))

	value6, err := s.load(addOffset(s.sp, 2088), 4)
	if err != nil {
		return 0, err
	}
	s.x[8] = uint64(uint32(value6))

	s.x[12] = uint64(uint32(uint64(0x33f1)))

	s.x[12] = uint64(uint32(insertBits(uint64(uint32(s.x[12])), 0x4431, 16, 16, 32)))

	s.x[13] = uint64(uint32(uint64(0x7d67)))

	s.x[13] = uint64(uint32(insertBits(uint64(uint32(s.x[13])), 0x86f9, 16, 16, 32)))

	s.x[14] = uint64(uint32(uint64(0x7049)))

	s.x[14] = uint64(uint32(insertBits(uint64(uint32(s.x[14])), 0xb028, 16, 16, 32)))

	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) + uint64(uint32(s.x[13])))))

	s.x[12] = uint64(uint32((uint64(uint32(s.x[10])) + uint64(uint32(s.x[12])))))

	s.x[13] = uint64(uint32(uint64(0x1418)))

	s.x[13] = uint64(uint32(insertBits(uint64(uint32(s.x[13])), 0xb028, 16, 16, 32)))

	s.x[13] = uint64(uint32((uint64(uint32(s.x[9])) + uint64(uint32(s.x[13])))))

	s.compare32(uint32(uint64(uint32(s.x[12]))), uint32(uint64(uint32(s.x[13]))))

	s.x[13] = uint64(uint32(boolValue(s.condition("lo"))))

	s.compare32(uint32(uint64(uint32(s.x[8]))), uint32(uint64(uint32(s.x[14]))))

	s.x[14] = uint64(uint32(boolValue(s.condition("lo"))))

	s.x[13] = uint64(uint32((uint64(uint32(s.x[14])) ^ uint64(uint32(s.x[13])))))

	s.compare32(uint32(uint64(uint32(s.x[12]))), uint32(uint64(uint32(s.x[8]))))

	s.x[8] = uint64(uint32(boolValue(s.condition("ls"))))

	s.compare32(uint32(uint64(uint32(s.x[13]))), uint32(uint64(0x0)))

	s.x[8] = uint64(uint32(selectValue(s.condition("ne"), uint64(uint32(s.x[14])), uint64(uint32(s.x[8])))))

	s.x[12] = uint64(uint32(uint64(0xe17)))

	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) * uint64(uint32(s.x[12])))))

	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) ^ uint64(uint32(s.x[11])))))

	value30, err := s.load((s.x[22] + shiftLeft(signExtend32(uint32(uint64(uint32(s.x[8])))), 3, 64)), 8)
	if err != nil {
		return 0, err
	}
	s.x[11] = uint64(value30)

	s.x[12] = uint64(uint32(uint64(0x1)))

	s.x[8] = uint64(uint32(uint64(0x7bc5)))

	s.x[8] = uint64(uint32(insertBits(uint64(uint32(s.x[8])), 0xce75, 16, 16, 32)))

	s.x[17] = uint64(uint64(0x7705))

	next = s.x[11]
	return next, nil
}

func runAccountBlock100a0b4d4(s *state) (uint64, error) {
	s.traceBlock(0x100a0b4d4)
	var next uint64

	value0, err := s.load(addOffset(s.sp, 856), 8)
	if err != nil {
		return 0, err
	}
	s.x[9] = uint64(value0)

	s.x[10] = uint64((s.x[9] - uint64(0x1b0)))

	if err := s.storeValue(addOffset(s.sp, 856), s.x[10], 8); err != nil {
		return 0, err
	}

	value3, err := s.load(addOffset(s.sp, 688), 8)
	if err != nil {
		return 0, err
	}
	s.x[8] = uint64(value3)

	value4, err := s.load(addOffset(s.sp, 368), 8)
	if err != nil {
		return 0, err
	}
	s.x[12] = uint64(value4)

	value5, err := s.load(addOffset(s.sp, 480), 8)
	if err != nil {
		return 0, err
	}
	s.x[6] = uint64(value5)

	value6, err := s.load(addOffset(s.sp, 284), 4)
	if err != nil {
		return 0, err
	}
	s.x[16] = uint64(uint32(value6))

	value7, err := s.load(addOffset(s.sp, 280), 4)
	if err != nil {
		return 0, err
	}
	s.x[0] = uint64(uint32(value7))

	s.x[19] = uint64(uint32(uint64(0x4e2e)))

	s.x[19] = uint64(uint32(insertBits(uint64(uint32(s.x[19])), 0xb6a1, 16, 16, 32)))

	s.x[20] = uint64(uint32(uint64(0x4221)))

	s.x[20] = uint64(uint32(insertBits(uint64(uint32(s.x[20])), 0x92a1, 16, 16, 32)))

	s.x[21] = uint64(uint32(uint64(0xc0c)))

	s.x[21] = uint64(uint32(insertBits(uint64(uint32(s.x[21])), 0x2400, 16, 16, 32)))

	s.x[22] = uint64(uint32(uint64(0x38b0)))

	s.x[22] = uint64(uint32(insertBits(uint64(uint32(s.x[22])), 0x5a85, 16, 16, 32)))

	s.x[25] = uint64(uint32(uint64(0x4e20)))

	s.x[25] = uint64(uint32(insertBits(uint64(uint32(s.x[25])), 0x36a1, 16, 16, 32)))

	s.x[27] = uint64(uint32(uint64(0x820)))

	s.x[27] = uint64(uint32(insertBits(uint64(uint32(s.x[27])), 0x1281, 16, 16, 32)))

	s.x[17] = uint64(uint64(0x12f))

	next = s.x[8]
	return next, nil
}

func runAccountBlock100aba4f8(s *state) (uint64, error) {
	s.traceBlock(0x100aba4f8)
	var next uint64

	s.x[11] = uint64(uint32(uint64(0xffff4f51)))

	s.x[11] = uint64(uint32((uint64(uint32(s.x[9])) + uint64(uint32(s.x[11])))))

	s.x[13] = uint64(uint32(uint64(0x3021)))

	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) | uint64(uint32(s.x[13])))))

	value4, err := s.load((s.x[30] + shiftLeft(signExtend32(uint32(uint64(uint32(s.x[8])))), 3, 64)), 8)
	if err != nil {
		return 0, err
	}
	s.x[8] = uint64(value4)

	if err := s.storeValue(addOffset(s.sp, 688), s.x[8], 8); err != nil {
		return 0, err
	}

	s.x[8] = uint64((s.sp + uint64(0x490)))

	s.x[27] = uint64((s.x[8] + s.x[10]))

	if err := s.storeValue(addOffset(s.sp, 728), s.x[27], 8); err != nil {
		return 0, err
	}

	s.x[8] = uint64((s.x[10] + uint64(0x1b0)))

	if err := s.storeValue(addOffset(s.sp, 856), s.x[8], 8); err != nil {
		return 0, err
	}

	if err := s.storeValue(addOffset(s.sp, 972), uint64(uint32(s.x[20])), 4); err != nil {
		return 0, err
	}

	if err := s.storeValue(addOffset(s.sp, 872), s.x[6], 8); err != nil {
		return 0, err
	}

	s.x[8] = uint64(uint32(uint64(0x43b3)))

	s.x[8] = uint64(uint32(insertBits(uint64(uint32(s.x[8])), 0x4e0, 16, 16, 32)))

	if err := s.storeValue(addOffset(s.sp, 804), uint64(uint32(s.x[8])), 4); err != nil {
		return 0, err
	}

	if err := s.storeValue(addOffset(s.sp, 1004), uint64(uint32(s.x[3])), 4); err != nil {
		return 0, err
	}

	if err := s.storeValue(addOffset(s.sp, 704), uint64(0), 8); err != nil {
		return 0, err
	}

	s.x[8] = uint64(uint32(uint64(0x7be8)))

	s.x[8] = uint64(uint32(insertBits(uint64(uint32(s.x[8])), 0x26b4, 16, 16, 32)))

	if err := s.storeValue(addOffset(s.sp, 1012), uint64(uint32(s.x[8])), 4); err != nil {
		return 0, err
	}

	s.x[25] = uint64(uint32(uint64(0xcf04)))

	s.x[25] = uint64(uint32(insertBits(uint64(uint32(s.x[25])), 0x4a7, 16, 16, 32)))

	if err := s.storeValue(addOffset(s.sp, 1076), uint64(uint32(s.x[25])), 4); err != nil {
		return 0, err
	}

	s.x[8] = uint64(uint32(uint64(0xe4ec)))

	s.x[8] = uint64(uint32(insertBits(uint64(uint32(s.x[8])), 0x785d, 16, 16, 32)))

	if err := s.storeValue(s.x[27], uint64(uint32(s.x[8])), 4); err != nil {
		return 0, err
	}

	s.x[8] = uint64(uint32(uint64(0x9b44)))

	s.x[8] = uint64(uint32(insertBits(uint64(uint32(s.x[8])), 0xfcac, 16, 16, 32)))

	s.x[8] = uint64(uint32((uint64(uint32(s.x[12])) + uint64(uint32(s.x[8])))))

	if err := s.storeValue(addOffset(s.sp, 456), uint64(uint32(s.x[8])), 4); err != nil {
		return 0, err
	}

	s.x[8] = uint64(uint32(uint64(0x688)))

	s.x[8] = uint64(uint32(insertBits(uint64(uint32(s.x[8])), 0x6d2f, 16, 16, 32)))

	s.x[8] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[8])))))

	s.x[8] = uint64(uint32((uint64(uint32(s.x[12])) + uint64(uint32(s.x[8])))))

	if err := s.storeValue(addOffset(s.sp, 496), uint64(uint32(s.x[8])), 4); err != nil {
		return 0, err
	}

	s.x[8] = uint64(uint32(uint64(0x5c0)))

	s.x[8] = uint64(uint32(insertBits(uint64(uint32(s.x[8])), 0xd1f0, 16, 16, 32)))

	s.x[8] = uint64(uint32((uint64(uint32(s.x[12])) + uint64(uint32(s.x[8])))))

	if err := s.storeValue(addOffset(s.sp, 448), uint64(uint32(s.x[8])), 4); err != nil {
		return 0, err
	}

	s.x[8] = uint64(uint32(uint64(0x356b)))

	s.x[8] = uint64(uint32(insertBits(uint64(uint32(s.x[8])), 0xd915, 16, 16, 32)))

	s.x[8] = uint64(uint32((uint64(uint32(s.x[12])) + uint64(uint32(s.x[8])))))

	if err := s.storeValue(addOffset(s.sp, 432), uint64(uint32(s.x[8])), 4); err != nil {
		return 0, err
	}

	s.x[8] = uint64(uint32(uint64(0x943b)))

	s.x[8] = uint64(uint32(insertBits(uint64(uint32(s.x[8])), 0x47c7, 16, 16, 32)))

	s.x[8] = uint64(uint32((uint64(uint32(s.x[12])) + uint64(uint32(s.x[8])))))

	if err := s.storeValue(addOffset(s.sp, 440), uint64(uint32(s.x[8])), 4); err != nil {
		return 0, err
	}

	s.x[8] = uint64((s.x[27] + uint64(0x188)))

	if err := s.storeValue(addOffset(s.sp, 344), s.x[8], 8); err != nil {
		return 0, err
	}

	if err := s.storeValue(addOffset(s.sp, 472), uint64(uint32(s.x[20])), 4); err != nil {
		return 0, err
	}

	s.x[23] = uint64(s.x[20])

	value52, err := s.load((s.x[30] + shiftLeft(signExtend32(uint32(uint64(uint32(s.x[9])))), 3, 64)), 8)
	if err != nil {
		return 0, err
	}
	s.x[8] = uint64(value52)

	s.x[17] = uint64(uint64(0x364))

	next = s.x[8]
	return next, nil
}

func runAccountBlock100b5f734(s *state) (uint64, error) {
	s.traceBlock(0x100b5f734)
	var next uint64

	s.x[8] = uint64(uint32(uint64(0x7e31)))

	s.x[8] = uint64(uint32(insertBits(uint64(uint32(s.x[8])), 0xd59c, 16, 16, 32)))

	s.x[8] = uint64(uint32((uint64(uint32(s.x[22])) ^ uint64(uint32(s.x[8])))))

	if err := s.storeValue(addOffset(s.sp, 1076), uint64(uint32(s.x[25])), 4); err != nil {
		return 0, err
	}

	if err := s.storeValue(addOffset(s.sp, 972), uint64(uint32(s.x[23])), 4); err != nil {
		return 0, err
	}

	if err := s.storeValue(addOffset(s.sp, 872), s.x[6], 8); err != nil {
		return 0, err
	}

	s.x[9] = uint64(uint32(uint64(0xea67)))

	s.x[9] = uint64(uint32(insertBits(uint64(uint32(s.x[9])), 0x6ea0, 16, 16, 32)))

	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) + uint64(uint32(s.x[9])))))

	s.compare32(uint32(uint64(uint32(s.x[22]))), uint32(uint64(uint32(s.x[8]))))

	s.x[8] = uint64(uint32(boolValue(s.condition("hi"))))

	s.x[9] = uint64(uint32(uint64(0x40b2)))

	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) * uint64(uint32(s.x[9])))))

	s.x[9] = uint64(uint32(uint64(0xf033)))

	s.x[9] = uint64(uint32(insertBits(uint64(uint32(s.x[9])), 0x2a63, 16, 16, 32)))

	s.x[9] = uint64(uint32((uint64(uint32(s.x[22])) + uint64(uint32(s.x[9])))))

	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) ^ uint64(uint32(s.x[9])))))

	value17, err := s.load((s.x[30] + shiftLeft(signExtend32(uint32(uint64(uint32(s.x[8])))), 3, 64)), 8)
	if err != nil {
		return 0, err
	}
	s.x[8] = uint64(value17)

	value18, err := s.load(addOffset(s.sp, 512), 8)
	if err != nil {
		return 0, err
	}
	s.x[7] = uint64(value18)

	s.x[17] = uint64(uint64(0x83c0))

	next = s.x[8]
	return next, nil
}

func runAccountBlock1000f50f4(s *state) (uint64, error) {
	s.traceBlock(0x1000f50f4)
	var next uint64

	value0, err := s.load(addOffset(s.sp, 88), 4)
	if err != nil {
		return 0, err
	}
	s.x[11] = uint64(uint32(value0))

	s.x[13] = uint64(uint32(uint64(0x53e8)))

	s.x[13] = uint64(uint32(insertBits(uint64(uint32(s.x[13])), 0x8a41, 16, 16, 32)))

	s.compare32(uint32(uint64(uint32(s.x[11]))), uint32(uint64(uint32(s.x[13]))))

	s.x[10] = uint64(uint32(boolValue(s.condition("lo"))))

	s.x[19] = uint64(uint32(uint64(0x3238)))

	s.x[12] = uint64(uint32(uint64(0x5428)))

	s.x[12] = uint64(uint32(insertBits(uint64(uint32(s.x[12])), 0x8a41, 16, 16, 32)))

	s.x[12] = uint64(uint32((uint64(uint32(s.x[8])) + uint64(uint32(s.x[12])))))

	s.compare32(uint32(uint64(uint32(s.x[12]))), uint32(uint64(uint32(s.x[13]))))

	s.x[15] = uint64(uint32(boolValue(s.condition("lo"))))

	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[15])))))

	s.compare32(uint32(uint64(uint32(s.x[12]))), uint32(uint64(uint32(s.x[11]))))

	s.x[12] = uint64(uint32(boolValue(s.condition("hi"))))

	s.compare32(uint32(uint64(uint32(s.x[10]))), uint32(uint64(0x0)))

	s.x[10] = uint64(uint32(selectValue(s.condition("ne"), uint64(uint32(s.x[15])), uint64(uint32(s.x[12])))))

	s.x[11] = uint64(uint32(uint64(0x58)))

	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) * uint64(uint32(s.x[11])))))

	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[4])))))

	value19, err := s.load((s.x[24] + shiftLeft(uint64(uint32(uint64(uint32(s.x[10])))), 3, 64)), 8)
	if err != nil {
		return 0, err
	}
	s.x[15] = uint64(value19)

	s.x[10] = uint64(uint32(uint64(0x9a4f)))

	s.x[10] = uint64(uint32(insertBits(uint64(uint32(s.x[10])), 0x67d2, 16, 16, 32)))

	s.x[12] = uint64(uint32(uint64(0xe2a)))

	s.x[12] = uint64(uint32(insertBits(uint64(uint32(s.x[12])), 0x4a9e, 16, 16, 32)))

	s.x[11] = uint64(uint32(uint64(0xcff5)))

	s.x[11] = uint64(uint32(insertBits(uint64(uint32(s.x[11])), 0xf52f, 16, 16, 32)))

	s.x[13] = uint64(uint64(0xfffffffffffff1d6))

	s.x[13] = uint64(insertBits(s.x[13], 0xb561, 16, 16, 64))

	s.x[0] = uint64(uint32(uint64(0x36)))

	s.x[2] = uint64(uint32(uint64(0x238e)))

	s.x[2] = uint64(uint32(insertBits(uint64(uint32(s.x[2])), 0xf530, 16, 16, 32)))

	s.x[3] = uint64(uint32(uint64(0xe6a)))

	s.x[3] = uint64(uint32(insertBits(uint64(uint32(s.x[3])), 0x4a9e, 16, 16, 32)))

	s.x[5] = uint64(uint32(uint64(0xe15a)))

	s.x[17] = uint64(uint64(0xe66f))

	next = s.x[15]
	return next, nil
}

func runAccountBlock1008d5268(s *state) (uint64, error) {
	s.traceBlock(0x1008d5268)
	var next uint64

	s.x[8] = uint64(uint32(uint64(0xffffbfee)))

	if err := s.storeValue(addOffset(s.sp, 1264), uint64(uint32(s.x[20])), 4); err != nil {
		return 0, err
	}

	value2, err := s.load(addOffset(s.sp, 1256), 8)
	if err != nil {
		return 0, err
	}
	s.x[9] = uint64(value2)

	value3, err := s.load(s.x[9], 4)
	if err != nil {
		return 0, err
	}
	s.x[9] = uint64(uint32(value3))

	value4, err := s.load(addOffset(s.sp, 1388), 4)
	if err != nil {
		return 0, err
	}
	s.x[10] = uint64(uint32(value4))

	s.x[11] = uint64(uint32(uint64(0x1fb9)))

	s.x[11] = uint64(uint32(insertBits(uint64(uint32(s.x[11])), 0x9e92, 16, 16, 32)))

	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) + uint64(uint32(s.x[11])))))

	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) + uint64(uint32(s.x[10])))))

	s.x[10] = uint64(uint32(uint64(0x1175)))

	s.x[23] = uint64(uint32((uint64(uint32(s.x[3])) + uint64(uint32(s.x[8])))))

	s.x[8] = uint64(uint32((uint64(uint32(s.x[23])) ^ uint64(uint32(s.x[10])))))

	value12, err := s.load((s.x[22] + shiftLeft(signExtend32(uint32(uint64(uint32(s.x[8])))), 3, 64)), 8)
	if err != nil {
		return 0, err
	}
	s.x[10] = uint64(value12)

	s.x[24] = uint64(s.x[3])

	s.x[8] = uint64(s.x[1])

	s.x[17] = uint64(uint64(0x8178))

	next = s.x[10]
	return next, nil
}

func runAccountBlock1008e64c0(s *state) (uint64, error) {
	s.traceBlock(0x1008e64c0)
	var next uint64

	s.x[8] = uint64(uint32(uint64(0x2cbb)))

	s.x[8] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[8])))))

	s.x[9] = uint64(uint32((uint64(uint32(s.x[8])) + shiftLeft(uint64(uint32(s.x[8])), 1, 32))))

	s.x[8] = uint64(uint32(uint64(0xc3a0)))

	s.x[8] = uint64(uint32(insertBits(uint64(uint32(s.x[8])), 0x9408, 16, 16, 32)))

	s.x[12] = uint64(uint32(uint64(0xffffa3c9)))

	s.x[12] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[12])))))

	s.x[8] = uint64(uint32((uint64(uint32(s.x[10])) + uint64(uint32(s.x[8])))))

	s.compare32(uint32(uint64(uint32(s.x[8]))), uint32(uint64(uint32(s.x[12]))))

	s.x[8] = uint64(uint32(boolValue(s.condition("lo"))))

	s.x[12] = uint64(uint32(uint64(0x9998)))

	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) * uint64(uint32(s.x[12])))))

	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) ^ uint64(uint32(s.x[11])))))

	value13, err := s.load((s.x[22] + shiftLeft(signExtend32(uint32(uint64(uint32(s.x[8])))), 3, 64)), 8)
	if err != nil {
		return 0, err
	}
	s.x[13] = uint64(value13)

	s.x[12] = uint64(uint32(uint64(0x1)))

	s.x[8] = uint64(uint32(uint64(0x7bc5)))

	s.x[8] = uint64(uint32(insertBits(uint64(uint32(s.x[8])), 0xce75, 16, 16, 32)))

	s.x[17] = uint64(uint64(0x7705))

	next = s.x[13]
	return next, nil
}

func runAccountBlock1008f5e1c(s *state) (uint64, error) {
	s.traceBlock(0x1008f5e1c)
	var next uint64

	s.x[17] = uint64(uint32((uint64(uint32(s.x[0])) ^ uint64(uint32(s.x[5])))))

	s.x[1] = uint64(uint32((uint64(uint32(s.x[15])) & uint64(uint32(s.x[0])))))

	s.x[2] = uint64(uint32((uint64(uint32(s.x[0])) ^ uint64(uint32(s.x[20])))))

	s.x[2] = uint64(uint32((uint64(uint32(s.x[2])) & shiftLeft(uint64(uint32(s.x[1])), 1, 32))))

	s.x[1] = uint64(uint32((uint64(uint32(s.x[2])) ^ uint64(uint32(s.x[1])))))

	s.x[2] = uint64(uint32(shiftLeft(uint64(uint32(s.x[17])), int(uint64(0x1)), 32)))

	s.x[3] = uint64(uint32((uint64(uint32(s.x[2])) ^ uint64(uint32(s.x[21])))))

	s.x[17] = uint64(uint32((uint64(uint32(s.x[3])) & uint64(uint32(s.x[17])))))

	s.x[2] = uint64(uint32((uint64(uint32(s.x[2])) & uint64(uint32(s.x[22])))))

	s.x[17] = uint64(uint32((uint64(uint32(s.x[17])) ^ uint64(uint32(s.x[2])))))

	s.x[2] = uint64(uint32((uint64(uint32(s.x[17])) ^ uint64(uint32(s.x[28])))))

	s.x[17] = uint64(uint32((uint64(uint32(s.x[17])) ^ uint64(uint32(s.x[10])))))

	s.x[17] = uint64(uint32((uint64(uint32(s.x[17])) & shiftLeft(uint64(uint32(s.x[1])), 2, 32))))

	s.x[17] = uint64(uint32((uint64(uint32(s.x[17])) ^ uint64(uint32(s.x[1])))))

	s.x[1] = uint64(uint32(shiftLeft(uint64(uint32(s.x[2])), int(uint64(0x2)), 32)))

	s.x[3] = uint64(uint32((uint64(uint32(s.x[1])) ^ uint64(uint32(s.x[19])))))

	s.x[2] = uint64(uint32((uint64(uint32(s.x[3])) & uint64(uint32(s.x[2])))))

	s.x[1] = uint64(uint32((uint64(uint32(s.x[1])) & uint64(uint32(s.x[23])))))

	s.x[1] = uint64(uint32((uint64(uint32(s.x[1])) ^ uint64(uint32(s.x[7])))))

	s.x[1] = uint64(uint32((uint64(uint32(s.x[1])) ^ uint64(uint32(s.x[2])))))

	s.x[1] = uint64(uint32((uint64(uint32(s.x[1])) & shiftLeft(uint64(uint32(s.x[17])), 4, 32))))

	s.x[17] = uint64(uint32((uint64(uint32(s.x[1])) ^ uint64(uint32(s.x[17])))))

	s.x[17] = uint64(uint32((uint64(uint32(s.x[0])) ^ shiftLeft(uint64(uint32(s.x[17])), 1, 32))))

	s.x[17] = uint64(uint32((uint64(uint32(s.x[17])) ^ uint64(0xffffffcf))))

	s.x[17] = uint64(uint32((uint64(uint32(s.x[17])) * uint64(uint32(s.x[16])))))

	s.x[1] = uint64((s.x[24] + s.x[14]))

	if err := s.storeValue(s.x[1], uint64(uint32(s.x[17])), 1); err != nil {
		return 0, err
	}

	s.x[0] = uint64((s.x[0] + uint64(0x1)))

	s.compare64(uint64(s.x[14]), uint64(uint64(0x0)))

	s.x[17] = uint64(uint32(boolValue(s.condition("eq"))))

	s.x[1] = uint64(uint32(shiftLeft(uint64(uint32(s.x[17])), int(uint64(0x8)), 32)))

	s.x[17] = uint64(uint32((uint64(uint32(s.x[1])) - shiftLeft(uint64(uint32(s.x[17])), 2, 32))))

	s.x[17] = uint64(uint32((uint64(uint32(s.x[17])) ^ uint64(uint32(s.x[12])))))

	value33, err := s.load((s.x[30] + shiftLeft(signExtend32(uint32(uint64(uint32(s.x[17])))), 3, 64)), 8)
	if err != nil {
		return 0, err
	}
	s.x[1] = uint64(value33)

	s.x[14] = uint64((s.x[14] - uint64(0x1)))

	s.x[16] = uint64((s.x[16] + uint64(0x1)))

	s.x[17] = uint64(uint64(0x60de))

	next = s.x[1]
	return next, nil
}

func runAccountBlock100900484(s *state) (uint64, error) {
	s.traceBlock(0x100900484)
	var next uint64

	s.x[8] = uint64(uint32(uint64(0xffffad3b)))

	s.x[11] = uint64(uint32(uint64(0x4042)))

	s.x[23] = uint64(uint32((uint64(uint32(s.x[9])) + uint64(uint32(s.x[8])))))

	s.x[24] = uint64(uint32((uint64(uint32(s.x[23])) + uint64(uint32(s.x[11])))))

	value4, err := s.load(addOffset(s.sp, 512), 8)
	if err != nil {
		return 0, err
	}
	s.x[8] = uint64(value4)

	if err := s.storeValue(addOffset(s.sp, 1256), s.x[8], 8); err != nil {
		return 0, err
	}

	value6, err := s.load(addOffset(s.sp, 2092), 4)
	if err != nil {
		return 0, err
	}
	s.x[8] = uint64(uint32(value6))

	s.x[9] = uint64(uint32(uint64(0x1fb5)))

	s.x[9] = uint64(uint32(insertBits(uint64(uint32(s.x[9])), 0x9e92, 16, 16, 32)))

	s.x[8] = uint64(uint32((uint64(uint32(s.x[10])) + uint64(uint32(s.x[8])))))

	s.x[9] = uint64(uint32((uint64(uint32(s.x[8])) + uint64(uint32(s.x[9])))))

	s.x[8] = uint64(uint32(uint64(0x582b)))

	s.x[8] = uint64(uint32((uint64(uint32(s.x[23])) ^ uint64(uint32(s.x[8])))))

	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) ^ uint64(uint32(s.x[24])))))

	value14, err := s.load((s.x[22] + shiftLeft(signExtend32(uint32(uint64(uint32(s.x[8])))), 3, 64)), 8)
	if err != nil {
		return 0, err
	}
	s.x[10] = uint64(value14)

	s.x[8] = uint64(s.x[1])

	s.x[17] = uint64(uint64(0x8178))

	next = s.x[10]
	return next, nil
}

func runAccountBlock1009ecc88(s *state) (uint64, error) {
	s.traceBlock(0x1009ecc88)
	var next uint64

	s.x[10] = uint64(uint32(uint64(0x2196)))

	s.x[11] = uint64(uint32(uint64(0x6b46)))

	s.x[10] = uint64(uint32((uint64(uint32(s.x[3])) + uint64(uint32(s.x[10])))))

	s.x[11] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[11])))))

	s.x[12] = uint64(uint32(uint64(0x24d7)))

	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) * uint64(uint32(s.x[12])))))

	value6, err := s.load(addOffset(s.sp, 2088), 4)
	if err != nil {
		return 0, err
	}
	s.x[12] = uint64(uint32(value6))

	s.x[13] = uint64(uint32(uint64(0xedac)))

	s.x[13] = uint64(uint32(insertBits(uint64(uint32(s.x[13])), 0x7453, 16, 16, 32)))

	s.x[14] = uint64(uint32(uint64(0x372a)))

	s.x[14] = uint64(uint32(insertBits(uint64(uint32(s.x[14])), 0xb71c, 16, 16, 32)))

	s.x[15] = uint64(uint32(uint64(0x2a0c)))

	s.x[15] = uint64(uint32(insertBits(uint64(uint32(s.x[15])), 0xe04b, 16, 16, 32)))

	s.x[13] = uint64(uint32((uint64(uint32(s.x[8])) + uint64(uint32(s.x[13])))))

	s.compare32(uint32(uint64(uint32(s.x[13]))), uint32(uint64(uint32(s.x[15]))))

	s.x[16] = uint64(uint32(boolValue(s.condition("lo"))))

	s.x[12] = uint64(uint32((uint64(uint32(s.x[12])) + uint64(uint32(s.x[14])))))

	s.compare32(uint32(uint64(uint32(s.x[12]))), uint32(uint64(uint32(s.x[15]))))

	s.x[14] = uint64(uint32(boolValue(s.condition("lo"))))

	s.x[14] = uint64(uint32((uint64(uint32(s.x[16])) ^ uint64(uint32(s.x[14])))))

	s.compare32(uint32(uint64(uint32(s.x[13]))), uint32(uint64(uint32(s.x[12]))))

	s.x[12] = uint64(uint32(boolValue(s.condition("hi"))))

	s.compare32(uint32(uint64(uint32(s.x[14]))), uint32(uint64(0x0)))

	s.x[12] = uint64(uint32(selectValue(s.condition("ne"), uint64(uint32(s.x[16])), uint64(uint32(s.x[12])))))

	s.x[12] = uint64(uint32((uint64(uint32(s.x[12])) ^ uint64(0x1))))

	s.x[13] = uint64(uint32(uint64(0xffffb661)))

	s.x[13] = uint64(uint32((uint64(uint32(s.x[11])) + uint64(uint32(s.x[13])))))

	s.x[12] = uint64(uint32((uint64(uint32(s.x[12])) * uint64(uint32(s.x[13])))))

	s.x[10] = uint64(uint32((uint64(uint32(s.x[12])) ^ uint64(uint32(s.x[10])))))

	value29, err := s.load((s.x[22] + shiftLeft(signExtend32(uint32(uint64(uint32(s.x[10])))), 3, 64)), 8)
	if err != nil {
		return 0, err
	}
	s.x[12] = uint64(value29)

	s.x[10] = uint64(s.x[9])

	s.x[17] = uint64(uint64(0x83e7))

	next = s.x[12]
	return next, nil
}

func runAccountBlock100a070f8(s *state) (uint64, error) {
	s.traceBlock(0x100a070f8)
	var next uint64

	s.x[10] = uint64(uint32(uint64(0x3b1f)))

	value1, err := s.load(addOffset(s.sp, 1264), 4)
	if err != nil {
		return 0, err
	}
	s.x[11] = uint64(uint32(value1))

	s.x[8] = uint64(uint32(uint64(0x7c64)))

	s.x[8] = uint64(uint32(insertBits(uint64(uint32(s.x[8])), 0x22e6, 16, 16, 32)))

	s.x[8] = uint64(uint32((uint64(uint32(s.x[3])) + uint64(uint32(s.x[8])))))

	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) ^ uint64(uint32(s.x[11])))))

	s.x[9] = uint64(uint32(uint64(0x39ee)))

	s.x[9] = uint64(uint32(insertBits(uint64(uint32(s.x[9])), 0xdd19, 16, 16, 32)))

	s.x[12] = uint64(uint32(uint64(0xc612)))

	s.x[12] = uint64(uint32(insertBits(uint64(uint32(s.x[12])), 0xb278, 16, 16, 32)))

	s.x[12] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[12])))))

	s.x[13] = uint64(uint32(uint64(0x39ee)))

	s.x[13] = uint64(uint32(insertBits(uint64(uint32(s.x[13])), 0x4d87, 16, 16, 32)))

	s.x[12] = uint64(uint32((uint64(uint32(s.x[12])) + uint64(uint32(s.x[13])))))

	s.x[13] = uint64(uint32(uint64(0x7bc5)))

	s.x[13] = uint64(uint32(insertBits(uint64(uint32(s.x[13])), 0x5eeb, 16, 16, 32)))

	s.x[13] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[13])))))

	s.x[14] = uint64(uint32(uint64(0x843b)))

	s.x[14] = uint64(uint32(insertBits(uint64(uint32(s.x[14])), 0xa114, 16, 16, 32)))

	s.x[13] = uint64(uint32((uint64(uint32(s.x[13])) + uint64(uint32(s.x[14])))))

	s.x[12] = uint64(uint32((uint64(uint32(s.x[13])) ^ uint64(uint32(s.x[12])))))

	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) + uint64(uint32(s.x[9])))))

	s.x[8] = uint64(uint32((uint64(uint32(s.x[12])) ^ uint64(uint32(s.x[8])))))

	value23, err := s.load(addOffset(s.sp, 1388), 4)
	if err != nil {
		return 0, err
	}
	s.x[9] = uint64(uint32(value23))

	s.x[12] = uint64(uint32(uint64(0x7bc5)))

	s.x[12] = uint64(uint32(insertBits(uint64(uint32(s.x[12])), 0xce75, 16, 16, 32)))

	s.x[13] = uint64(uint32((uint64(uint32(s.x[9])) + uint64(uint32(s.x[12])))))

	s.x[8] = uint64(uint32((uint64(uint32(s.x[13])) + uint64(uint32(s.x[8])))))

	if err := s.storeValue(addOffset(s.sp, 1268), uint64(uint32(s.x[8])), 4); err != nil {
		return 0, err
	}

	s.x[13] = uint64(uint32(uint64(0xfec)))

	s.x[13] = uint64(uint32(insertBits(uint64(uint32(s.x[13])), 0xf1be, 16, 16, 32)))

	s.x[13] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[13])))))

	s.x[14] = uint64(uint32(uint64(0xf014)))

	s.x[14] = uint64(uint32(insertBits(uint64(uint32(s.x[14])), 0xe41, 16, 16, 32)))

	s.x[13] = uint64(uint32((uint64(uint32(s.x[13])) + uint64(uint32(s.x[14])))))

	s.x[14] = uint64(uint32(uint64(0x85b0)))

	s.x[14] = uint64(uint32(insertBits(uint64(uint32(s.x[14])), 0xbdc, 16, 16, 32)))

	s.x[14] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[14])))))

	s.x[15] = uint64(uint32(uint64(0x7a50)))

	s.x[15] = uint64(uint32(insertBits(uint64(uint32(s.x[15])), 0xf423, 16, 16, 32)))

	s.x[14] = uint64(uint32((uint64(uint32(s.x[14])) + uint64(uint32(s.x[15])))))

	s.x[13] = uint64(uint32((uint64(uint32(s.x[13])) ^ uint64(uint32(s.x[14])))))

	s.x[14] = uint64(uint32(uint64(0xf199)))

	s.x[14] = uint64(uint32(insertBits(uint64(uint32(s.x[14])), 0x3417, 16, 16, 32)))

	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[14])))))

	s.x[14] = uint64(uint32(uint64(0xe67)))

	s.x[14] = uint64(uint32(insertBits(uint64(uint32(s.x[14])), 0xcbe8, 16, 16, 32)))

	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) + uint64(uint32(s.x[14])))))

	s.x[11] = uint64(uint32((uint64(uint32(s.x[13])) ^ uint64(uint32(s.x[11])))))

	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) + uint64(uint32(s.x[12])))))

	s.compare32(uint32(uint64(uint32(s.x[11]))), uint32(uint64(0x8)))

	s.x[11] = uint64(uint32(boolValue(s.condition("lo"))))

	s.x[12] = uint64(uint32(uint64(0x1a)))

	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) * uint64(uint32(s.x[12])))))

	s.x[10] = uint64(uint32((uint64(uint32(s.x[3])) + uint64(uint32(s.x[10])))))

	s.x[10] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[10])))))

	value56, err := s.load((s.x[22] + shiftLeft(signExtend32(uint32(uint64(uint32(s.x[10])))), 3, 64)), 8)
	if err != nil {
		return 0, err
	}
	s.x[12] = uint64(value56)

	s.x[11] = uint64(s.x[3])

	s.x[10] = uint64(s.x[9])

	s.x[17] = uint64(uint64(0xb2ae))

	next = s.x[12]
	return next, nil
}

func runAccountBlock100a4aa04(s *state) (uint64, error) {
	s.traceBlock(0x100a4aa04)
	var next uint64

	s.x[13] = uint64(uint32(uint64(0xd509)))

	s.x[13] = uint64(uint32(insertBits(uint64(uint32(s.x[13])), 0xbb05, 16, 16, 32)))

	s.x[12] = uint64(uint32(uint64(0xac0f)))

	s.x[12] = uint64(uint32(insertBits(uint64(uint32(s.x[12])), 0xc247, 16, 16, 32)))

	s.x[14] = uint64(uint32(uint64(0x6efb)))

	s.x[14] = uint64(uint32(insertBits(uint64(uint32(s.x[14])), 0x3db8, 16, 16, 32)))

	s.x[12] = uint64(uint32((uint64(uint32(s.x[11])) + uint64(uint32(s.x[12])))))

	s.x[12] = uint64(uint32((uint64(uint32(s.x[12])) & uint64(uint32(s.x[14])))))

	s.x[14] = uint64(uint32(uint64(0xae47)))

	s.x[14] = uint64(uint32(insertBits(uint64(uint32(s.x[14])), 0x1107, 16, 16, 32)))

	s.x[13] = uint64(uint32((uint64(uint32(s.x[11])) + uint64(uint32(s.x[13])))))

	s.x[14] = uint64(uint32((uint64(uint32(s.x[13])) ^ uint64(uint32(s.x[14])))))

	s.x[14] = uint64(uint32((uint64(uint32(s.x[12])) + uint64(uint32(s.x[14])))))

	s.compare32(uint32(uint64(uint32(s.x[20]))), uint32(uint64(uint32(s.x[14]))))

	s.x[14] = uint64(uint32(boolValue(s.condition("lt"))))

	s.x[15] = uint64(uint32(uint64(0x27)))

	s.x[14] = uint64(uint32((uint64(uint32(s.x[14])) * uint64(uint32(s.x[15])))))

	s.x[15] = uint64(uint32(uint64(0x3a9c)))

	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) + uint64(uint32(s.x[15])))))

	s.x[11] = uint64(uint32((uint64(uint32(s.x[14])) ^ uint64(uint32(s.x[11])))))

	value20, err := s.load((s.x[22] + shiftLeft(signExtend32(uint32(uint64(uint32(s.x[11])))), 3, 64)), 8)
	if err != nil {
		return 0, err
	}
	s.x[11] = uint64(value20)

	s.x[17] = uint64(uint64(0x1a88))

	next = s.x[11]
	return next, nil
}

func runAccountBlock100a5bf1c(s *state) (uint64, error) {
	s.traceBlock(0x100a5bf1c)
	var next uint64

	s.x[10] = uint64(s.x[8])

	s.x[8] = uint64(uint32(uint64(0x63cf)))

	s.x[12] = uint64(uint32((uint64(uint32(s.x[11])) + uint64(uint32(s.x[8])))))

	value3, err := s.load(addOffset(s.sp, 1384), 4)
	if err != nil {
		return 0, err
	}
	s.x[8] = uint64(uint32(value3))

	value4, err := s.load((s.x[22] + shiftLeft(signExtend32(uint32(uint64(uint32(s.x[12])))), 3, 64)), 8)
	if err != nil {
		return 0, err
	}
	s.x[12] = uint64(value4)

	s.x[17] = uint64(uint64(0x4bf8))

	next = s.x[12]
	return next, nil
}

func runAccountBlock100bed25c(s *state) (uint64, error) {
	s.traceBlock(0x100bed25c)
	var next uint64

	s.x[12] = uint64(uint32(uint64(0x0)))

	s.x[9] = uint64(uint32(uint64(0x4b35)))

	s.x[9] = uint64(uint32(insertBits(uint64(uint32(s.x[9])), 0x1f0d, 16, 16, 32)))

	s.x[10] = uint64(uint32(uint64(0xdefb)))

	s.x[10] = uint64(uint32(insertBits(uint64(uint32(s.x[10])), 0xe0f2, 16, 16, 32)))

	s.x[9] = uint64(uint32((uint64(uint32(s.x[11])) + uint64(uint32(s.x[9])))))

	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) & uint64(uint32(s.x[10])))))

	s.x[10] = uint64(uint32(uint64(0x27f4)))

	s.x[10] = uint64(uint32(insertBits(uint64(uint32(s.x[10])), 0xce75, 16, 16, 32)))

	s.x[10] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[10])))))

	s.compare32(uint32(uint64(uint32(s.x[8]))), uint32(uint64(uint32(s.x[10]))))

	s.x[10] = uint64(uint32(boolValue(s.condition("eq"))))

	s.x[13] = uint64(uint32(uint64(0xa56a)))

	s.x[13] = uint64(uint32((uint64(uint32(s.x[10])) * uint64(uint32(s.x[13])))))

	value14, err := s.load(addOffset(s.sp, 1268), 4)
	if err != nil {
		return 0, err
	}
	s.x[10] = uint64(uint32(value14))

	if err := s.storeValue(addOffset(s.sp, 1388), uint64(uint32(s.x[10])), 4); err != nil {
		return 0, err
	}

	if err := s.storeValue(addOffset(s.sp, 1384), uint64(uint32(s.x[8])), 4); err != nil {
		return 0, err
	}

	s.x[13] = uint64(uint32((uint64(uint32(s.x[13])) ^ uint64(uint32(s.x[11])))))

	value18, err := s.load((s.x[22] + shiftLeft(signExtend32(uint32(uint64(uint32(s.x[13])))), 3, 64)), 8)
	if err != nil {
		return 0, err
	}
	s.x[13] = uint64(value18)

	s.x[17] = uint64(uint64(0x7705))

	next = s.x[13]
	return next, nil
}
