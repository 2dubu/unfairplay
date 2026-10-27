package recovery

const (
	accountDispatchRoute020410     = uint64(0x100020410)
	accountCallSetup0650f0         = uint64(0x1000650f0)
	accountRecordSetup110c30       = uint64(0x100110c30)
	accountStateSetup150c5c        = uint64(0x100150c5c)
	accountVectorRoute15dcc0       = uint64(0x10015dcc0)
	accountFrameRoute178eb8        = uint64(0x100178eb8)
	accountLoopDispatchFront197fa4 = uint64(0x100197fa4)
	accountLoopDispatchTail198000  = uint64(0x100198000)
	accountStackReturn1b5db0       = uint64(0x1001b5db0)
	accountStateRoute2413b4        = uint64(0x1002413b4)
	accountTableRoute247b38        = uint64(0x100247b38)
	accountDispatchRoute24b5a8     = uint64(0x10024b5a8)
	accountNestedCallSetup24ed8c   = uint64(0x10024ed8c)
	accountNestedCallResume24eef0  = uint64(0x10024eef0)
	accountThresholdRoute30e674    = uint64(0x10030e674)
	accountDispatchSetup35ea64     = uint64(0x10035ea64)
	accountDispatchSetup386528     = uint64(0x100386528)
)

func runAccountDispatchRoute020410(s *state) (uint64, error) {
	s.traceBlock(0x100020410)
	var next uint64

	s.x[8] = uint64(uint32(uint64(0x6856)))

	s.x[8] = uint64(uint32(insertBits(uint64(uint32(s.x[8])), 0xebb8, 16, 16, 32)))

	s.x[8] = uint64(uint32((uint64(uint32(s.x[9])) + uint64(uint32(s.x[8])))))

	value3, err := s.load(addOffset(s.sp, 100), 4)
	if err != nil {
		return 0, err
	}
	s.x[27] = uint64(uint32(value3))

	value4, err := s.load((s.x[20] + shiftLeft(signExtend32(uint32(uint64(uint32(s.x[8])))), 3, 64)), 8)
	if err != nil {
		return 0, err
	}
	s.x[8] = uint64(value4)

	s.x[10] = uint64(uint32(uint64(0x7bc5)))

	s.x[10] = uint64(uint32(insertBits(uint64(uint32(s.x[10])), 0xce75, 16, 16, 32)))

	s.x[17] = uint64(uint64(0x165a))

	next = s.x[8]
	return next, nil
}

func runAccountCallSetup0650f0(s *state) (uint64, error) {
	s.traceBlock(0x1000650f0)
	var next uint64

	s.x[8] = uint64(uint32(uint64(0xfac)))

	s.x[8] = uint64(uint32(insertBits(uint64(uint32(s.x[8])), 0x1448, 16, 16, 32)))

	s.x[9] = uint64(uint32((uint64(uint32(s.x[10])) + uint64(uint32(s.x[8])))))

	s.x[8] = uint64(uint32(uint64(0x3054)))

	s.x[8] = uint64(uint32(insertBits(uint64(uint32(s.x[8])), 0x1448, 16, 16, 32)))

	s.x[23] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[8])))))

	value6, err := s.load(addOffset(s.x[12], 36), 4)
	if err != nil {
		return 0, err
	}
	s.x[11] = uint64(uint32(value6))

	s.x[8] = uint64(uint32(uint64(0xc04f)))

	s.x[8] = uint64(uint32(insertBits(uint64(uint32(s.x[8])), 0x7c79, 16, 16, 32)))

	s.x[8] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[8])))))

	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) + uint64(uint32(s.x[11])))))

	s.x[12] = uint64(uint32(uint64(0xe6a7)))

	s.x[12] = uint64(uint32(insertBits(uint64(uint32(s.x[12])), 0x7fff, 16, 16, 32)))

	s.x[12] = uint64(uint32((uint64(uint32(s.x[10])) + uint64(uint32(s.x[12])))))

	s.compare32(uint32(uint64(uint32(s.x[8]))), uint32(uint64(uint32(s.x[12]))))

	s.x[8] = uint64(uint32(boolValue(s.condition("lo"))))

	s.x[12] = uint64(uint32(uint64(0x219a)))

	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) * uint64(uint32(s.x[12])))))

	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) ^ uint64(uint32(s.x[10])))))

	value19, err := s.load((s.x[20] + shiftLeft(signExtend32(uint32(uint64(uint32(s.x[8])))), 3, 64)), 8)
	if err != nil {
		return 0, err
	}
	s.x[8] = uint64(value19)

	s.x[12] = uint64(uint32(uint64(0x201b)))

	s.x[12] = uint64(uint32(insertBits(uint64(uint32(s.x[12])), 0x318a, 16, 16, 32)))

	s.x[10] = uint64(uint32(uint64(0x7bc5)))

	s.x[10] = uint64(uint32(insertBits(uint64(uint32(s.x[10])), 0xce75, 16, 16, 32)))

	s.x[17] = uint64(uint64(0xdb39))

	next = s.x[8]
	return next, nil
}

func runAccountRecordSetup110c30(s *state) (uint64, error) {
	s.traceBlock(0x100110c30)
	var next uint64

	s.x[4] = uint64(s.x[25])

	s.x[8] = uint64(uint32(uint64(0x5371)))

	s.x[8] = uint64(uint32((uint64(uint32(s.x[21])) ^ uint64(uint32(s.x[8])))))

	s.x[9] = uint64(uint32(uint64(0xf23)))

	s.x[11] = uint64(uint32((uint64(uint32(s.x[8])) * uint64(uint32(s.x[9])))))

	value5, err := s.load(addOffset(s.sp, 152), 4)
	if err != nil {
		return 0, err
	}
	s.x[8] = uint64(uint32(value5))

	s.x[9] = uint64(uint32(uint64(0x7f05)))

	s.x[9] = uint64(uint32(insertBits(uint64(uint32(s.x[9])), 0x5df2, 16, 16, 32)))

	s.x[9] = uint64(uint32((uint64(uint32(s.x[8])) ^ uint64(uint32(s.x[9])))))

	s.x[10] = uint64(uint32(uint64(0x80fb)))

	s.x[10] = uint64(uint32(insertBits(uint64(uint32(s.x[10])), 0xa20d, 16, 16, 32)))

	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) + uint64(uint32(s.x[10])))))

	s.x[10] = uint64(uint32(uint64(0x296a)))

	s.x[10] = uint64(uint32(insertBits(uint64(uint32(s.x[10])), 0x2e65, 16, 16, 32)))

	s.x[10] = uint64(uint32((uint64(uint32(s.x[8])) ^ uint64(uint32(s.x[10])))))

	s.x[12] = uint64(uint32(uint64(0x5ead)))

	s.x[12] = uint64(uint32(insertBits(uint64(uint32(s.x[12])), 0xd19a, 16, 16, 32)))

	s.x[12] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[12])))))

	s.x[13] = uint64(uint32(uint64(0x15dc)))

	s.x[13] = uint64(uint32(insertBits(uint64(uint32(s.x[13])), 0x7777, 16, 16, 32)))

	s.x[13] = uint64(uint32((uint64(uint32(s.x[8])) ^ uint64(uint32(s.x[13])))))

	s.x[14] = uint64(uint32(uint64(0xea24)))

	s.x[14] = uint64(uint32(insertBits(uint64(uint32(s.x[14])), 0x8888, 16, 16, 32)))

	s.x[13] = uint64(uint32((uint64(uint32(s.x[13])) + uint64(uint32(s.x[14])))))

	s.x[9] = uint64(uint32((uint64(uint32(s.x[13])) ^ uint64(uint32(s.x[9])))))

	s.x[10] = uint64(uint32((uint64(uint32(s.x[12])) + uint64(uint32(s.x[10])))))

	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[10])))))

	s.x[10] = uint64(uint32(uint64(0xb5ee)))

	s.x[10] = uint64(uint32(insertBits(uint64(uint32(s.x[10])), 0xa980, 16, 16, 32)))

	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) + uint64(uint32(s.x[10])))))

	value30, err := s.load((s.x[24] + shiftLeft(signExtend32(uint32(uint64(uint32(s.x[21])))), 3, 64)), 8)
	if err != nil {
		return 0, err
	}
	s.x[10] = uint64(value30)

	s.x[19] = uint64(uint32(uint64(0xed9b)))

	s.x[19] = uint64(uint32(insertBits(uint64(uint32(s.x[19])), 0xd768, 16, 16, 32)))

	s.x[21] = uint64(uint32(uint64(0x80000003)))

	s.x[26] = uint64(uint32(uint64(0x1)))

	s.x[17] = uint64(uint64(0x24d4))

	next = s.x[10]
	return next, nil
}

func runAccountStateSetup150c5c(s *state) (uint64, error) {
	s.traceBlock(0x100150c5c)
	var next uint64

	s.x[8] = uint64(uint32(uint64(0x7372)))

	s.x[8] = uint64(uint32(insertBits(uint64(uint32(s.x[8])), 0x46e9, 16, 16, 32)))

	s.x[8] = uint64(uint32((uint64(uint32(s.x[27])) ^ uint64(uint32(s.x[8])))))

	s.x[11] = uint64(uint32(uint64(0x7603)))

	s.x[11] = uint64(uint32(insertBits(uint64(uint32(s.x[11])), 0xad0, 16, 16, 32)))

	s.x[28] = uint64(uint32((uint64(uint32(s.x[9])) + uint64(uint32(s.x[11])))))

	s.x[11] = uint64(uint32(uint64(0x151c)))

	s.x[11] = uint64(uint32(insertBits(uint64(uint32(s.x[11])), 0xb916, 16, 16, 32)))

	s.x[11] = uint64(uint32((uint64(uint32(s.x[28])) ^ uint64(uint32(s.x[11])))))

	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) + uint64(uint32(s.x[11])))))

	s.x[11] = uint64(uint32(uint64(0x6a12)))

	s.x[11] = uint64(uint32(insertBits(uint64(uint32(s.x[11])), 0x370e, 16, 16, 32)))

	s.x[11] = uint64(uint32((uint64(uint32(s.x[27])) ^ uint64(uint32(s.x[11])))))

	s.x[12] = uint64(uint32(uint64(0x95ee)))

	s.x[12] = uint64(uint32(insertBits(uint64(uint32(s.x[12])), 0xc8f1, 16, 16, 32)))

	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) + uint64(uint32(s.x[12])))))

	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) ^ uint64(uint32(s.x[11])))))

	s.x[11] = uint64(uint32(uint64(0x3f16)))

	s.x[11] = uint64(uint32(insertBits(uint64(uint32(s.x[11])), 0xca62, 16, 16, 32)))

	s.x[11] = uint64(uint32((uint64(uint32(s.x[9])) + uint64(uint32(s.x[11])))))

	s.x[11] = uint64(uint32((uint64(uint32(s.x[27])) ^ uint64(uint32(s.x[11])))))

	s.x[12] = uint64(uint32(uint64(0x9d5b)))

	s.x[12] = uint64(uint32(insertBits(uint64(uint32(s.x[12])), 0x406d, 16, 16, 32)))

	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) + uint64(uint32(s.x[12])))))

	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) ^ uint64(uint32(s.x[11])))))

	s.x[11] = uint64(uint32(uint64(0xba09)))

	s.x[11] = uint64(uint32(insertBits(uint64(uint32(s.x[11])), 0xd6fa, 16, 16, 32)))

	s.x[11] = uint64(uint32((uint64(uint32(s.x[22])) + uint64(uint32(s.x[11])))))

	s.x[14] = uint64(uint32(uint64(0x35f1)))

	s.x[14] = uint64(uint32(insertBits(uint64(uint32(s.x[14])), 0xfdaf, 16, 16, 32)))

	s.compare32(uint32(uint64(uint32(s.x[11]))), uint32(uint64(uint32(s.x[14]))))

	s.x[12] = uint64(uint32(boolValue(s.condition("lo"))))

	s.x[13] = uint64(uint32(uint64(0xb1b6)))

	s.x[13] = uint64(uint32(insertBits(uint64(uint32(s.x[13])), 0xcc24, 16, 16, 32)))

	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) + uint64(uint32(s.x[13])))))

	s.compare32(uint32(uint64(uint32(s.x[8]))), uint32(uint64(uint32(s.x[14]))))

	s.x[13] = uint64(uint32(boolValue(s.condition("lo"))))

	s.x[12] = uint64(uint32((uint64(uint32(s.x[12])) ^ uint64(uint32(s.x[13])))))

	s.compare32(uint32(uint64(uint32(s.x[11]))), uint32(uint64(uint32(s.x[8]))))

	s.x[8] = uint64(uint32(boolValue(s.condition("lo"))))

	s.compare32(uint32(uint64(uint32(s.x[12]))), uint32(uint64(0x0)))

	s.x[8] = uint64(uint32(selectValue(s.condition("ne"), uint64(uint32(s.x[13])), uint64(uint32(s.x[8])))))

	s.x[11] = uint64(uint32(uint64(0x29)))

	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) * uint64(uint32(s.x[11])))))

	s.x[11] = uint64(uint32(uint64(0x77ff)))

	s.x[11] = uint64(uint32(insertBits(uint64(uint32(s.x[11])), 0xad0, 16, 16, 32)))

	s.x[11] = uint64(uint32((uint64(uint32(s.x[9])) + uint64(uint32(s.x[11])))))

	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) ^ uint64(uint32(s.x[11])))))

	value48, err := s.load((s.x[24] + shiftLeft(signExtend32(uint32(uint64(uint32(s.x[8])))), 3, 64)), 8)
	if err != nil {
		return 0, err
	}
	s.x[8] = uint64(value48)

	s.x[17] = uint64(uint64(0xc2b))

	next = s.x[8]
	return next, nil
}

func runAccountVectorRoute15dcc0(s *state) (uint64, error) {
	s.traceBlock(0x10015dcc0)
	var next uint64

	s.x[26] = uint64(uint32(uint64(0x80000003)))

	s.x[21] = uint64(uint32(uint64(0xed9b)))

	s.x[21] = uint64(uint32(insertBits(uint64(uint32(s.x[21])), 0xd768, 16, 16, 32)))

	s.x[9] = uint64(uint32(uint64(0x5b45)))

	s.x[9] = uint64(uint32((uint64(uint32(s.x[12])) ^ uint64(uint32(s.x[9])))))

	s.x[10] = uint64((s.x[10] + shiftLeft(s.x[10], 1, 64)))

	s.x[23] = uint64((s.x[23] + shiftLeft(s.x[10], 4, 64)))

	value7, err := s.load(addOffset(s.x[23], 40), 4)
	if err != nil {
		return 0, err
	}
	s.x[11] = uint64(uint32(value7))

	s.x[13] = uint64(uint32((uint64(uint32(s.x[11])) & uint64(0x1))))

	s.x[10] = uint64(uint32(uint64(0x4ec0)))

	s.x[10] = uint64(uint32((uint64(uint32(s.x[9])) + uint64(uint32(s.x[10])))))

	s.x[14] = uint64(uint32(uint64(0x6d7c)))

	s.x[14] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[14])))))

	s.compare32(uint32(uint64(uint32(s.x[13]))), uint32(uint64(uint32(s.x[14]))))

	s.x[13] = uint64(uint32(boolValue(s.condition("ne"))))

	s.x[14] = uint64(uint32(shiftLeft(uint64(uint32(s.x[13])), int(uint64(0xa)), 32)))

	s.x[13] = uint64(uint32((uint64(uint32(s.x[14])) - shiftLeft(uint64(uint32(s.x[13])), 2, 32))))

	s.x[12] = uint64(uint32((uint64(uint32(s.x[13])) ^ uint64(uint32(s.x[12])))))

	value18, err := s.load((s.x[24] + shiftLeft(signExtend32(uint32(uint64(uint32(s.x[12])))), 3, 64)), 8)
	if err != nil {
		return 0, err
	}
	s.x[12] = uint64(value18)

	s.x[17] = uint64(uint64(0xddc5))

	next = s.x[12]
	return next, nil
}

func runAccountFrameRoute178eb8(s *state) (uint64, error) {
	s.traceBlock(0x100178eb8)
	var next uint64

	s.x[8] = uint64(uint32((uint64(uint32(s.x[28])) - uint64(0x5bc))))

	s.x[9] = uint64(uint32(uint64(0xd57f)))

	s.x[9] = uint64(uint32(insertBits(uint64(uint32(s.x[9])), 0xfa08, 16, 16, 32)))

	s.x[9] = uint64(uint32((uint64(uint32(s.x[28])) + uint64(uint32(s.x[9])))))

	s.x[10] = uint64(uint32(uint64(0x9f9a)))

	s.x[10] = uint64(uint32(insertBits(uint64(uint32(s.x[10])), 0x5f7, 16, 16, 32)))

	s.x[10] = uint64(uint32((uint64(uint32(s.x[9])) & uint64(uint32(s.x[10])))))

	value7, err := s.load(addOffset(s.x[23], 32), 4)
	if err != nil {
		return 0, err
	}
	s.x[11] = uint64(uint32(value7))

	s.compare32(uint32(uint64(uint32(s.x[11]))), uint32(uint64(uint32(s.x[16]))))

	s.x[9] = uint64(uint32(boolValue(s.condition("ne"))))

	s.x[12] = uint64(uint32(uint64(0xffff78c9)))

	s.x[12] = uint64(uint32((uint64(uint32(s.x[8])) ^ uint64(uint32(s.x[12])))))

	s.x[12] = uint64(uint32((uint64(uint32(s.x[10])) + uint64(uint32(s.x[12])))))

	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) * uint64(uint32(s.x[12])))))

	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) | uint64(uint32(s.x[28])))))

	value15, err := s.load((s.x[24] + shiftLeft(signExtend32(uint32(uint64(uint32(s.x[9])))), 3, 64)), 8)
	if err != nil {
		return 0, err
	}
	s.x[9] = uint64(value15)

	s.x[17] = uint64(uint64(0x52ff))

	next = s.x[9]
	return next, nil
}

func runAccountLoopDispatchFront197fa4(s *state) (uint64, error) {
	s.traceBlock(0x100197fa4)

	s.x[10] = uint64(uint32(uint64(0x2774)))

	s.x[10] = uint64(uint32(insertBits(uint64(uint32(s.x[10])), 0xfd33, 16, 16, 32)))

	s.x[10] = uint64(uint32((uint64(uint32(s.x[11])) + uint64(uint32(s.x[10])))))

	s.x[12] = uint64(uint32(uint64(0x73df)))

	s.x[12] = uint64(uint32(insertBits(uint64(uint32(s.x[12])), 0xf7fc, 16, 16, 32)))

	s.x[12] = uint64(uint32((uint64(uint32(s.x[10])) & uint64(uint32(s.x[12])))))

	s.x[13] = uint64(uint32((uint64(uint32(s.x[20])) + uint64(uint32(s.x[19])))))

	s.x[10] = uint64(signExtend32(uint32(uint64(uint32(s.x[13])))))

	s.x[13] = uint64((s.x[10] + shiftLeft(signExtend32(uint32(uint64(uint32(s.x[13])))), 1, 64)))

	s.x[13] = uint64(shiftLeft(s.x[13], int(uint64(0x4)), 64))

	value10, err := s.load((s.x[23] + s.x[13]), 4)
	if err != nil {
		return 0, err
	}
	s.x[14] = uint64(uint32(value10))

	s.x[13] = uint64(uint32(uint64(0xf345)))

	s.x[13] = uint64(uint32(insertBits(uint64(uint32(s.x[13])), 0xacf, 16, 16, 32)))

	s.x[13] = uint64(uint32((uint64(uint32(s.x[12])) + uint64(uint32(s.x[13])))))

	s.x[15] = uint64(uint32(uint64(0xff04)))

	s.x[15] = uint64(uint32(insertBits(uint64(uint32(s.x[15])), 0x378d, 16, 16, 32)))

	s.x[15] = uint64(uint32((uint64(uint32(s.x[13])) ^ uint64(uint32(s.x[15])))))

	s.compare32(uint32(uint64(uint32(s.x[14]))), uint32(uint64(uint32(s.x[15]))))

	s.x[15] = uint64(uint32(boolValue(s.condition("eq"))))

	s.x[17] = uint64(uint32(uint64(0x9d)))

	s.x[15] = uint64(uint32((uint64(uint32(s.x[15])) * uint64(uint32(s.x[17])))))

	s.x[11] = uint64(uint32((uint64(uint32(s.x[15])) ^ uint64(uint32(s.x[11])))))

	value22, err := s.load((s.x[24] + shiftLeft(signExtend32(uint32(uint64(uint32(s.x[11])))), 3, 64)), 8)
	if err != nil {
		return 0, err
	}
	s.x[11] = uint64(value22)
	return 0x100198000, nil
}

func runAccountLoopDispatchTail198000(s *state) (uint64, error) {
	s.traceBlock(0x100198000)
	var next uint64

	s.x[17] = uint64(uint64(0xb16d))

	next = s.x[11]
	return next, nil
}

func runAccountStackReturn1b5db0(s *state) (uint64, error) {
	s.traceBlock(0x1001b5db0)
	var next uint64

	s.sp = uint64((s.sp + uint64(0x10)))

	next = s.x[30]
	return next, nil
}

func runAccountStateRoute2413b4(s *state) (uint64, error) {
	s.traceBlock(0x1002413b4)
	var next uint64

	s.x[9] = uint64(uint32(uint64(0xffffb6fc)))

	s.x[28] = uint64(uint32((uint64(uint32(s.x[10])) + uint64(uint32(s.x[9])))))

	s.x[9] = uint64(uint32(uint64(0xffffacb4)))

	s.x[9] = uint64(uint32((uint64(uint32(s.x[10])) + uint64(uint32(s.x[9])))))

	s.x[12] = uint64(uint32(uint64(0x48c)))

	s.x[12] = uint64(uint32((uint64(uint32(s.x[9])) | uint64(uint32(s.x[12])))))

	s.x[9] = uint64(uint32(uint64(0xffffe146)))

	s.x[9] = uint64(uint32((uint64(uint32(s.x[12])) + uint64(uint32(s.x[9])))))

	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) & uint64(uint32(s.x[11])))))

	s.x[11] = uint64(uint32(uint64(0xffff9284)))

	s.x[11] = uint64(uint32((uint64(uint32(s.x[10])) + uint64(uint32(s.x[11])))))

	s.compare32(uint32(uint64(uint32(s.x[9]))), uint32(uint64(uint32(s.x[11]))))

	s.x[9] = uint64(uint32(boolValue(s.condition("eq"))))

	s.x[11] = uint64(uint32(uint64(0xfd97)))

	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) * uint64(uint32(s.x[11])))))

	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[10])))))

	value16, err := s.load((s.x[24] + shiftLeft(signExtend32(uint32(uint64(uint32(s.x[9])))), 3, 64)), 8)
	if err != nil {
		return 0, err
	}
	s.x[9] = uint64(value16)

	s.x[17] = uint64(uint64(0xa86b))

	next = s.x[9]
	return next, nil
}

func runAccountTableRoute247b38(s *state) (uint64, error) {
	s.traceBlock(0x100247b38)
	var next uint64

	s.x[11] = uint64(uint32(uint64(0x1188)))

	s.x[11] = uint64(uint32(insertBits(uint64(uint32(s.x[11])), 0x6499, 16, 16, 32)))

	s.x[15] = uint64(uint32((uint64(uint32(s.x[13])) + uint64(uint32(s.x[11])))))

	s.x[11] = uint64(uint32(uint64(0xfff3)))

	s.x[11] = uint64(uint32(insertBits(uint64(uint32(s.x[11])), 0x9b66, 16, 16, 32)))

	s.x[11] = uint64(uint32((uint64(uint32(s.x[15])) & uint64(uint32(s.x[11])))))

	s.x[12] = uint64(uint32(uint64(0xd12b)))

	s.x[12] = uint64(uint32(insertBits(uint64(uint32(s.x[12])), 0xe6f7, 16, 16, 32)))

	s.x[12] = uint64(uint32((uint64(uint32(s.x[13])) + uint64(uint32(s.x[12])))))

	s.x[17] = uint64(uint32(uint64(0x5df9)))

	s.x[17] = uint64(uint32(insertBits(uint64(uint32(s.x[17])), 0x1908, 16, 16, 32)))

	s.x[12] = uint64(uint32((uint64(uint32(s.x[12])) & uint64(uint32(s.x[17])))))

	s.x[14] = uint64(uint32((uint64(uint32(s.x[9])) - uint64(uint32(s.x[14])))))

	s.x[17] = uint64(uint32(uint64(0x886b)))

	s.x[17] = uint64(uint32(insertBits(uint64(uint32(s.x[17])), 0x6c12, 16, 16, 32)))

	s.compare32(uint32(uint64(uint32(s.x[14]))), uint32(uint64(uint32(s.x[17]))))

	s.x[14] = uint64(uint32(boolValue(s.condition("eq"))))

	s.x[14] = uint64(uint32((uint64(uint32(s.x[15])) ^ uint64(uint32(s.x[14])))))

	s.x[14] = uint64(uint32((uint64(uint32(s.x[26])) &^ uint64(uint32(s.x[14])))))

	s.x[15] = uint64(uint32(uint64(0xffffba09)))

	s.x[15] = uint64(uint32((uint64(uint32(s.x[12])) + uint64(uint32(s.x[15])))))

	s.x[14] = uint64(uint32((uint64(uint32(s.x[14])) * uint64(uint32(s.x[15])))))

	s.x[13] = uint64(uint32((uint64(uint32(s.x[14])) ^ uint64(uint32(s.x[13])))))

	value23, err := s.load((s.x[24] + shiftLeft(signExtend32(uint32(uint64(uint32(s.x[13])))), 3, 64)), 8)
	if err != nil {
		return 0, err
	}
	s.x[13] = uint64(value23)

	s.x[17] = uint64(uint64(0xcc6))

	next = s.x[13]
	return next, nil
}

func runAccountDispatchRoute24b5a8(s *state) (uint64, error) {
	s.traceBlock(0x10024b5a8)
	var next uint64

	s.x[8] = uint64(uint32(uint64(0x2c04)))

	s.x[8] = uint64(uint32(insertBits(uint64(uint32(s.x[8])), 0xad0, 16, 16, 32)))

	s.x[8] = uint64(uint32((uint64(uint32(s.x[9])) + uint64(uint32(s.x[8])))))

	value3, err := s.load(addOffset(s.sp, 172), 4)
	if err != nil {
		return 0, err
	}
	s.x[22] = uint64(uint32(value3))

	value4, err := s.load((s.x[24] + shiftLeft(signExtend32(uint32(uint64(uint32(s.x[8])))), 3, 64)), 8)
	if err != nil {
		return 0, err
	}
	s.x[8] = uint64(value4)

	s.x[10] = uint64(uint32(uint64(0x7bc5)))

	s.x[10] = uint64(uint32(insertBits(uint64(uint32(s.x[10])), 0xce75, 16, 16, 32)))

	s.x[17] = uint64(uint64(0x3b05))

	next = s.x[8]
	return next, nil
}

func runAccountNestedCallSetup24ed8c(s *state) (uint64, error) {
	s.traceBlock(0x10024ed8c)
	var next uint64

	s.x[8] = uint64(uint32(uint64(0x667d)))

	s.x[8] = uint64(uint32(insertBits(uint64(uint32(s.x[8])), 0xd40f, 16, 16, 32)))

	s.x[8] = uint64(uint32((uint64(uint32(s.x[28])) + uint64(uint32(s.x[8])))))

	s.x[9] = uint64(uint32(uint64(0x2380)))

	s.x[9] = uint64(uint32(insertBits(uint64(uint32(s.x[9])), 0x2120, 16, 16, 32)))

	s.x[19] = uint64(uint32((uint64(uint32(s.x[8])) | uint64(uint32(s.x[9])))))

	s.x[8] = uint64(uint32(uint64(0x2fe9)))

	s.x[8] = uint64(uint32(insertBits(uint64(uint32(s.x[8])), 0xad0, 16, 16, 32)))

	s.x[21] = uint64(uint32((uint64(uint32(s.x[19])) + uint64(uint32(s.x[8])))))

	s.x[8] = uint64(uint32(uint64(0x7667)))

	s.x[8] = uint64(uint32(insertBits(uint64(uint32(s.x[8])), 0xac1d, 16, 16, 32)))

	s.x[8] = uint64(uint32((uint64(uint32(s.x[27])) ^ uint64(uint32(s.x[8])))))

	s.x[9] = uint64(uint32(uint64(0x8999)))

	s.x[9] = uint64(uint32(insertBits(uint64(uint32(s.x[9])), 0x53e2, 16, 16, 32)))

	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) + uint64(uint32(s.x[9])))))

	s.x[9] = uint64(uint32(uint64(0x5f58)))

	s.x[9] = uint64(uint32(insertBits(uint64(uint32(s.x[9])), 0x3b6a, 16, 16, 32)))

	s.x[9] = uint64(uint32((uint64(uint32(s.x[27])) ^ uint64(uint32(s.x[9])))))

	s.x[10] = uint64(uint32(uint64(0xa0a8)))

	s.x[10] = uint64(uint32(insertBits(uint64(uint32(s.x[10])), 0xc495, 16, 16, 32)))

	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) + uint64(uint32(s.x[10])))))

	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) ^ uint64(uint32(s.x[9])))))

	s.x[9] = uint64(uint32(uint64(0x52fa)))

	s.x[9] = uint64(uint32(insertBits(uint64(uint32(s.x[9])), 0x5902, 16, 16, 32)))

	s.x[9] = uint64(uint32((uint64(uint32(s.x[27])) ^ uint64(uint32(s.x[9])))))

	s.x[10] = uint64(uint32(uint64(0x598e)))

	s.x[10] = uint64(uint32(insertBits(uint64(uint32(s.x[10])), 0xa6fd, 16, 16, 32)))

	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) + uint64(uint32(s.x[10])))))

	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) + uint64(uint32(s.x[21])))))

	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) ^ uint64(uint32(s.x[9])))))

	s.x[10] = uint64((s.x[29] - uint64(0xa0)))

	s.x[9] = uint64(uint32(uint64(0xf400)))

	s.x[9] = uint64(uint32(insertBits(uint64(uint32(s.x[9])), 0x6221, 16, 16, 32)))

	s.x[9] = uint64(uint32((uint64(uint32(s.x[10])) & uint64(uint32(s.x[9])))))

	s.x[9] = uint64(uint32(shiftLeft(uint64(uint32(s.x[9])), int(uint64(0x1)), 32)))

	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) - uint64(uint32(s.x[10])))))

	s.x[10] = uint64(uint32(uint64(0xbfc)))

	s.x[10] = uint64(uint32(insertBits(uint64(uint32(s.x[10])), 0x1dde, 16, 16, 32)))

	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) + uint64(uint32(s.x[10])))))

	s.x[10] = uint64(uint32(uint64(0xd99e)))

	s.x[10] = uint64(uint32(insertBits(uint64(uint32(s.x[10])), 0xfaec, 16, 16, 32)))

	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[10])))))

	s.x[10] = uint64(uint32(uint64(0x9ef7)))

	s.x[10] = uint64(uint32(insertBits(uint64(uint32(s.x[10])), 0xba20, 16, 16, 32)))

	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) * uint64(uint32(s.x[10])))))

	s.x[10] = uint64(uint32((uint64(uint32(s.x[28])) - uint64(uint32(s.x[9])))))

	s.x[11] = uint64(uint32(uint64(0xffff930c)))

	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) + uint64(uint32(s.x[11])))))

	if err := s.storeValue(addOffset(s.x[29], -160), uint64(uint32(s.x[10])), 4); err != nil {
		return 0, err
	}

	s.x[10] = uint64((s.sp + uint64(0x98)))

	if err := s.storeValue(addOffset(s.x[29], -112), s.x[10], 8); err != nil {
		return 0, err
	}

	s.x[10] = uint64((s.sp + uint64(0xac)))

	if err := s.storeValue(addOffset(s.x[29], -104), s.x[10], 8); err != nil {
		return 0, err
	}

	s.x[10] = uint64((s.sp + uint64(0xb0)))

	if err := s.storeValue(addOffset(s.x[29], -128), s.x[10], 8); err != nil {
		return 0, err
	}

	s.x[10] = uint64((s.sp + uint64(0x9c)))

	if err := s.storeValue(addOffset(s.x[29], -136), s.x[10], 8); err != nil {
		return 0, err
	}

	value57, err := s.load(addOffset(s.sp, 80), 8)
	if err != nil {
		return 0, err
	}
	s.x[10] = uint64(value57)

	if err := s.storeValue(addOffset(s.x[29], -152), s.x[10], 8); err != nil {
		return 0, err
	}

	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) - uint64(uint32(s.x[9])))))

	s.x[9] = uint64(uint32(uint64(0x9fc)))

	s.x[9] = uint64(uint32(insertBits(uint64(uint32(s.x[9])), 0x2afa, 16, 16, 32)))

	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) + uint64(uint32(s.x[9])))))

	if err := s.storeValue(addOffset(s.x[29], -120), uint64(uint32(s.x[8])), 4); err != nil {
		return 0, err
	}

	value64a, err := s.load(addOffset(s.sp, 128), 8)
	if err != nil {
		return 0, err
	}
	value64b, err := s.load((addOffset(s.sp, 128))+8, 8)
	if err != nil {
		return 0, err
	}
	s.x[9] = uint64(value64a)
	s.x[8] = uint64(value64b)

	value65, err := s.load(s.x[8], 4)
	if err != nil {
		return 0, err
	}
	s.x[8] = uint64(uint32(value65))

	value66, err := s.load(s.x[9], 4)
	if err != nil {
		return 0, err
	}
	s.x[9] = uint64(uint32(value66))

	s.x[8] = uint64(uint32((uint64(uint32(s.x[9])) - uint64(uint32(s.x[8])))))

	s.x[9] = uint64(uint32(uint64(0xa6)))

	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) ^ uint64(uint32(s.x[9])))))

	s.x[9] = uint64(uint32(uint64(0xd4cf)))

	s.x[9] = uint64(uint32(insertBits(uint64(uint32(s.x[9])), 0xa5c, 16, 16, 32)))

	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) * uint64(uint32(s.x[9])))))

	value73, err := s.load(addOffset(s.sp, 120), 8)
	if err != nil {
		return 0, err
	}
	s.x[9] = uint64(value73)

	s.x[9] = uint64((s.x[9] + uint64(uint8(uint64(uint32(s.x[8]))))))

	value75, err := s.load(s.x[9], 1)
	if err != nil {
		return 0, err
	}
	s.x[9] = uint64(uint32(value75))

	s.x[10] = uint64(uint32(uint64(0x35)))

	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[10])))))

	value78, err := s.load(addOffset(s.sp, 112), 8)
	if err != nil {
		return 0, err
	}
	s.x[10] = uint64(value78)

	value79, err := s.load((s.x[10] + uint64(uint32(uint64(uint32(s.x[9]))))), 1)
	if err != nil {
		return 0, err
	}
	s.x[9] = uint64(uint32(value79))

	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) ^ uint64(uint32(s.x[9])))))

	s.x[9] = uint64(uint32(uint64(0xcf7d)))

	s.x[8] = uint64(uint32((uint64(uint32(s.x[9])) + uint64(uint8(uint64(uint32(s.x[8])))))))

	value83, err := s.load((s.x[24] + shiftLeft(uint64(uint32(uint64(uint32(s.x[8])))), 3, 64)), 8)
	if err != nil {
		return 0, err
	}
	s.x[8] = uint64(value83)

	s.x[0] = uint64((s.x[29] - uint64(0xa0)))

	if err := s.storeValue(addOffset(s.sp, 104), s.x[5], 8); err != nil {
		return 0, err
	}

	s.x[25] = uint64(s.x[4])

	s.x[17] = uint64(uint64(0x48c7))

	s.x[30] = 0x10024eef0
	next = s.x[8]
	return next, nil
}

func runAccountNestedCallResume24eef0(s *state) (uint64, error) {
	s.traceBlock(0x10024eef0)
	var next uint64

	s.x[30] = uint64(uint32(uint64(0xa401)))

	s.x[7] = uint64(uint32(uint64(0xf3a8)))

	s.x[7] = uint64(uint32(insertBits(uint64(uint32(s.x[7])), 0x2896, 16, 16, 32)))

	s.x[6] = uint64(uint32(uint64(0x3d33)))

	s.x[6] = uint64(uint32(insertBits(uint64(uint32(s.x[6])), 0xf530, 16, 16, 32)))

	s.x[1] = uint64(uint32(uint64(0x2ff4)))

	s.x[3] = uint64(uint32(uint64(0x8d41)))

	s.x[3] = uint64(uint32(insertBits(uint64(uint32(s.x[3])), 0x567f, 16, 16, 32)))

	s.x[2] = uint64(uint32(uint64(0x7bc5)))

	s.x[2] = uint64(uint32(insertBits(uint64(uint32(s.x[2])), 0xce75, 16, 16, 32)))

	value10, err := s.load(addOffset(s.sp, 104), 8)
	if err != nil {
		return 0, err
	}
	s.x[5] = uint64(value10)

	value11, err := s.load(addOffset(s.x[29], -144), 4)
	if err != nil {
		return 0, err
	}
	s.x[13] = uint64(uint32(value11))

	s.compare32(uint32(uint64(uint32(s.x[13]))), uint32(uint64(uint32(s.x[2]))))

	s.x[8] = uint64(uint32(boolValue(s.condition("eq"))))

	s.x[9] = uint64(uint32(uint64(0x53e2)))

	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) * uint64(uint32(s.x[9])))))

	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) ^ uint64(uint32(s.x[28])))))

	value17, err := s.load((s.x[24] + shiftLeft(signExtend32(uint32(uint64(uint32(s.x[8])))), 3, 64)), 8)
	if err != nil {
		return 0, err
	}
	s.x[8] = uint64(value17)

	s.x[16] = uint64(uint32(uint64(0x1264)))

	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0x2897, 16, 16, 32)))

	s.x[17] = uint64(uint64(0x8960))

	next = s.x[8]
	return next, nil
}

func runAccountThresholdRoute30e674(s *state) (uint64, error) {
	s.traceBlock(0x10030e674)
	var next uint64

	s.x[12] = uint64(uint32(uint64(0xeebc)))

	s.x[12] = uint64(uint32(insertBits(uint64(uint32(s.x[12])), 0x6c5a, 16, 16, 32)))

	s.x[13] = uint64(uint32(uint64(0x6634)))

	s.x[13] = uint64(uint32(insertBits(uint64(uint32(s.x[13])), 0xe4d2, 16, 16, 32)))

	if err := s.storeValue(addOffset(s.x[8], 72), uint64(uint32(s.x[12])), 4); err != nil {
		return 0, err
	}
	if err := s.storeValue((addOffset(s.x[8], 72))+4, uint64(uint32(s.x[13])), 4); err != nil {
		return 0, err
	}

	s.x[12] = uint64(uint32(uint64(0x1143)))

	s.x[12] = uint64(uint32(insertBits(uint64(uint32(s.x[12])), 0x93a5, 16, 16, 32)))

	s.x[13] = uint64(uint32(uint64(0x99cb)))

	s.x[13] = uint64(uint32(insertBits(uint64(uint32(s.x[13])), 0x1b2d, 16, 16, 32)))

	if err := s.storeValue(addOffset(s.x[8], 80), uint64(uint32(s.x[12])), 4); err != nil {
		return 0, err
	}
	if err := s.storeValue((addOffset(s.x[8], 80))+4, uint64(uint32(s.x[13])), 4); err != nil {
		return 0, err
	}

	if err := s.storeValue(addOffset(s.x[8], 92), uint64(uint32(s.x[11])), 4); err != nil {
		return 0, err
	}

	value11, err := s.load((s.x[9] + shiftLeft(signExtend32(uint32(uint64(uint32(s.x[10])))), 3, 64)), 8)
	if err != nil {
		return 0, err
	}
	s.x[8] = uint64(value11)

	s.x[17] = uint64(uint64(0x8c87))

	next = s.x[8]
	return next, nil
}

func runAccountDispatchSetup35ea64(s *state) (uint64, error) {
	s.traceBlock(0x10035ea64)
	var next uint64

	s.sp = uint64((s.sp - uint64(0x10)))

	s.x[8] = uint64(0x101075000)

	value2, err := s.load(addOffset(s.x[8], 1664), 8)
	if err != nil {
		return 0, err
	}
	s.x[9] = uint64(value2)

	s.x[10] = uint64(0x101075000)

	value4, err := s.load(addOffset(s.x[10], 2880), 8)
	if err != nil {
		return 0, err
	}
	s.x[11] = uint64(value4)

	s.x[9] = uint64((s.x[11] ^ s.x[9]))

	s.x[11] = uint64((s.sp + uint64(0x8)))

	s.x[12] = uint64(uint64(0xd1a6))

	s.x[12] = uint64(insertBits(s.x[12], 0x8eae, 16, 16, 64))

	s.x[12] = uint64(insertBits(s.x[12], 0x4620, 32, 16, 64))

	s.x[12] = uint64(insertBits(s.x[12], 0x73d0, 48, 16, 64))

	s.x[9] = uint64((s.x[9] + s.x[11]))

	s.x[11] = uint64((s.x[9] ^ s.x[12]))

	s.x[9] = uint64(uint32(uint64(0xd4cf)))

	s.x[9] = uint64(uint32(insertBits(uint64(uint32(s.x[9])), 0xa5c, 16, 16, 32)))

	s.x[11] = uint64((s.x[11] * s.x[9]))

	if err := s.storeValue(addOffset(s.x[8], 1664), s.x[11], 8); err != nil {
		return 0, err
	}

	if err := s.storeValue(addOffset(s.x[10], 2880), s.x[11], 8); err != nil {
		return 0, err
	}

	s.x[8] = uint64(uint32(uint64(0xd102)))

	s.x[8] = uint64(uint32(insertBits(uint64(uint32(s.x[8])), 0x3715, 16, 16, 32)))

	s.x[8] = uint64(uint32((uint64(uint32(s.x[0])) & uint64(uint32(s.x[8])))))

	s.x[8] = uint64(uint32(shiftLeft(uint64(uint32(s.x[8])), int(uint64(0x1)), 32)))

	s.x[10] = uint64(uint32(uint64(0x2efd)))

	s.x[10] = uint64(uint32(insertBits(uint64(uint32(s.x[10])), 0x48ea, 16, 16, 32)))

	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) - uint64(uint32(s.x[0])))))

	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) + uint64(uint32(s.x[10])))))

	s.x[10] = uint64(uint32(uint64(0xc9bf)))

	s.x[10] = uint64(uint32(insertBits(uint64(uint32(s.x[10])), 0xfdcb, 16, 16, 32)))

	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) ^ uint64(uint32(s.x[10])))))

	s.x[10] = uint64(uint32(uint64(0xe98d)))

	s.x[10] = uint64(uint32(insertBits(uint64(uint32(s.x[10])), 0x74ca, 16, 16, 32)))

	s.x[10] = uint64(uint32((uint64(uint32(s.x[8])) * uint64(uint32(s.x[10])))))

	value32, err := s.load(addOffset(s.x[0], 8), 8)
	if err != nil {
		return 0, err
	}
	s.x[8] = uint64(value32)

	value33, err := s.load(s.x[0], 4)
	if err != nil {
		return 0, err
	}
	s.x[12] = uint64(uint32(value33))

	value34, err := s.load(addOffset(s.x[8], 68), 4)
	if err != nil {
		return 0, err
	}
	s.x[13] = uint64(uint32(value34))

	s.x[12] = uint64(uint32((uint64(uint32(s.x[12])) + uint64(uint32(s.x[10])))))

	s.x[13] = uint64(uint32((uint64(uint32(s.x[12])) + uint64(uint32(s.x[13])))))

	s.x[14] = uint64(uint32(uint64(0xc81f)))

	s.x[14] = uint64(uint32(insertBits(uint64(uint32(s.x[14])), 0xc663, 16, 16, 32)))

	s.x[15] = uint64(uint32(uint64(0xbfc3)))

	s.x[15] = uint64(uint32(insertBits(uint64(uint32(s.x[15])), 0xacd9, 16, 16, 32)))

	s.x[16] = uint64(uint32(uint64(0x85b)))

	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0x198a, 16, 16, 32)))

	s.x[14] = uint64(uint32((uint64(uint32(s.x[13])) + uint64(uint32(s.x[14])))))

	s.compare32(uint32(uint64(uint32(s.x[14]))), uint32(uint64(uint32(s.x[16]))))

	s.x[16] = uint64(uint32(boolValue(s.condition("hi"))))

	s.x[17] = uint64(uint32(uint64(0xf7a4)))

	s.x[17] = uint64(uint32(insertBits(uint64(uint32(s.x[17])), 0xe675, 16, 16, 32)))

	s.x[12] = uint64(uint32((uint64(uint32(s.x[12])) + uint64(uint32(s.x[15])))))

	s.compare32(uint32(uint64(uint32(s.x[12]))), uint32(uint64(uint32(s.x[17]))))

	s.x[17] = uint64(uint32(boolValue(s.condition("lo"))))

	s.x[13] = uint64(uint32((uint64(uint32(s.x[13])) + uint64(uint32(s.x[15])))))

	s.compare32(uint32(uint64(uint32(s.x[13]))), uint32(uint64(uint32(s.x[12]))))

	value53, err := s.load(addOffset(s.x[0], 16), 4)
	if err != nil {
		return 0, err
	}
	s.x[12] = uint64(uint32(value53))

	if err := s.storeValue(addOffset(s.x[8], 68), uint64(uint32(s.x[14])), 4); err != nil {
		return 0, err
	}

	s.x[13] = uint64(uint32((uint64(uint32(s.x[17])) ^ uint64(uint32(s.x[16])))))

	s.x[14] = uint64(uint32(boolValue(s.condition("lo"))))

	s.compare32(uint32(uint64(uint32(s.x[13]))), uint32(uint64(0x0)))

	s.x[13] = uint64(uint32(selectValue(s.condition("ne"), uint64(uint32(s.x[17])), uint64(uint32(s.x[14])))))

	s.x[14] = uint64(uint32(uint64(0xcbda)))

	s.x[14] = uint64(uint32(insertBits(uint64(uint32(s.x[14])), 0xd35a, 16, 16, 32)))

	s.x[12] = uint64(uint32((uint64(uint32(s.x[12])) - uint64(uint32(s.x[10])))))

	s.x[10] = uint64(uint32((uint64(uint32(s.x[12])) + uint64(uint32(s.x[14])))))

	s.x[14] = uint64(uint32(uint64(0xf67f)))

	s.x[14] = uint64(uint32(insertBits(uint64(uint32(s.x[14])), 0x6fbf, 16, 16, 32)))

	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) & uint64(uint32(s.x[14])))))

	s.x[14] = uint64(uint32(uint64(0xfbf2)))

	s.x[14] = uint64(uint32(insertBits(uint64(uint32(s.x[14])), 0xbce4, 16, 16, 32)))

	s.x[14] = uint64(uint32((uint64(uint32(s.x[10])) + uint64(uint32(s.x[14])))))

	s.x[13] = uint64(uint32((uint64(uint32(s.x[13])) * uint64(uint32(s.x[14])))))

	s.x[14] = uint64(uint32(uint64(0xa6)))

	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) + uint64(uint32(s.x[11])))))

	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[14])))))

	s.x[9] = uint64(uint32((uint64(uint32(s.x[11])) * uint64(uint32(s.x[9])))))

	s.x[11] = uint64(0x100fe9000)

	s.x[11] = uint64((s.x[11] + uint64(0x190)))

	s.x[11] = uint64((s.x[11] + uint64(uint8(uint64(uint32(s.x[9]))))))

	value77, err := s.load(s.x[11], 1)
	if err != nil {
		return 0, err
	}
	s.x[11] = uint64(uint32(value77))

	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(0x10))))

	s.x[14] = uint64(0x100de2000)

	s.x[14] = uint64((s.x[14] + uint64(0x8f0)))

	value81, err := s.load((s.x[14] + uint64(uint32(uint64(uint32(s.x[11]))))), 1)
	if err != nil {
		return 0, err
	}
	s.x[11] = uint64(uint32(value81))

	s.x[12] = uint64(uint32((uint64(uint32(s.x[13])) ^ uint64(uint32(s.x[12])))))

	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[11])))))

	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) & uint64(0xff))))

	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) + uint64(0x24e))))

	s.x[11] = uint64(0x10105b000)

	s.x[11] = uint64((s.x[11] + uint64(0x6f0)))

	value88, err := s.load((s.x[11] + shiftLeft(uint64(uint32(uint64(uint32(s.x[9])))), 3, 64)), 8)
	if err != nil {
		return 0, err
	}
	s.x[9] = uint64(value88)

	s.x[11] = uint64(uint64(0xffffffffffff5355))

	s.x[11] = uint64(insertBits(s.x[11], 0x9824, 16, 16, 64))

	s.x[9] = uint64((s.x[9] + s.x[11]))

	value92, err := s.load((s.x[9] + shiftLeft(signExtend32(uint32(uint64(uint32(s.x[12])))), 3, 64)), 8)
	if err != nil {
		return 0, err
	}
	s.x[11] = uint64(value92)

	s.x[17] = uint64(uint64(0x331e))

	next = s.x[11]
	return next, nil
}

func runAccountDispatchSetup386528(s *state) (uint64, error) {
	s.traceBlock(0x100386528)
	var next uint64

	s.sp = uint64((s.sp - uint64(0x10)))

	s.x[8] = uint64(0x101081000)

	value2, err := s.load(addOffset(s.x[8], 680), 8)
	if err != nil {
		return 0, err
	}
	s.x[9] = uint64(value2)

	s.x[10] = uint64(0x101075000)

	value4, err := s.load(addOffset(s.x[10], 2880), 8)
	if err != nil {
		return 0, err
	}
	s.x[11] = uint64(value4)

	s.x[12] = uint64((s.sp + uint64(0x8)))

	s.x[9] = uint64((s.x[11] - s.x[9]))

	s.x[9] = uint64((s.x[9] ^ s.x[12]))

	s.x[12] = uint64(uint32(uint64(0xd4cf)))

	s.x[12] = uint64(uint32(insertBits(uint64(uint32(s.x[12])), 0xa5c, 16, 16, 32)))

	s.x[11] = uint64(uint64(0x2c6))

	s.x[11] = uint64(insertBits(s.x[11], 0x75f, 16, 16, 64))

	s.x[11] = uint64(insertBits(s.x[11], 0x3a80, 32, 16, 64))

	s.x[11] = uint64(insertBits(s.x[11], 0x673e, 48, 16, 64))

	s.x[13] = uint64((s.x[9]*s.x[12] + s.x[11]))

	s.x[11] = uint64(uint64(0xd1a6))

	s.x[11] = uint64(insertBits(s.x[11], 0x8eae, 16, 16, 64))

	s.x[11] = uint64(insertBits(s.x[11], 0x4620, 32, 16, 64))

	s.x[11] = uint64(insertBits(s.x[11], 0x73d0, 48, 16, 64))

	s.x[9] = uint64((s.x[9] ^ s.x[11]))

	s.x[9] = uint64((s.x[9] * s.x[12]))

	if err := s.storeValue(addOffset(s.x[8], 680), s.x[13], 8); err != nil {
		return 0, err
	}

	if err := s.storeValue(addOffset(s.x[10], 2880), s.x[9], 8); err != nil {
		return 0, err
	}

	s.x[8] = uint64(uint32(uint64(0x2e6e)))

	s.x[8] = uint64(uint32(insertBits(uint64(uint32(s.x[8])), 0x6155, 16, 16, 32)))

	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) | shiftLeft(uint64(uint32(s.x[0])), 1, 32))))

	s.x[10] = uint64(uint32(uint64(0x68c9)))

	s.x[10] = uint64(uint32(insertBits(uint64(uint32(s.x[10])), 0xcf55, 16, 16, 32)))

	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) - uint64(uint32(s.x[0])))))

	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) + uint64(uint32(s.x[10])))))

	s.x[10] = uint64(uint32(uint64(0xe43)))

	s.x[10] = uint64(uint32(insertBits(uint64(uint32(s.x[10])), 0xc97e, 16, 16, 32)))

	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) ^ uint64(uint32(s.x[10])))))

	s.x[10] = uint64(uint32(uint64(0xa827)))

	s.x[10] = uint64(uint32(insertBits(uint64(uint32(s.x[10])), 0xda51, 16, 16, 32)))

	s.x[10] = uint64(uint32((uint64(uint32(s.x[8])) * uint64(uint32(s.x[10])))))

	value36, err := s.load(s.x[0], 4)
	if err != nil {
		return 0, err
	}
	s.x[8] = uint64(uint32(value36))

	s.x[14] = uint64(uint32((uint64(uint32(s.x[8])) ^ uint64(uint32(s.x[10])))))

	value38, err := s.load(addOffset(s.x[0], 16), 8)
	if err != nil {
		return 0, err
	}
	s.x[8] = uint64(value38)

	value39, err := s.load(addOffset(s.x[0], 8), 4)
	if err != nil {
		return 0, err
	}
	s.x[11] = uint64(uint32(value39))

	s.x[10] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[10])))))

	s.x[11] = uint64(uint32((uint64(uint32(s.x[14])) + uint64(0x735))))

	s.x[15] = uint64(uint32(uint64(0x124)))

	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) | uint64(uint32(s.x[15])))))

	s.x[15] = uint64(uint32((uint64(uint32(s.x[10])) & uint64(0x1ff))))

	s.x[16] = uint64(uint32(uint64(0x45d1)))

	s.x[15] = uint64(uint32((uint64(uint32(s.x[15])) ^ uint64(uint32(s.x[16])))))

	s.x[15] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[15])))))

	if err := s.storeValue(addOffset(s.x[8], 64), uint64(0), 8); err != nil {
		return 0, err
	}

	if err := s.storeValue(addOffset(s.x[8], 92), uint64(0), 4); err != nil {
		return 0, err
	}
	if err := s.storeValue((addOffset(s.x[8], 92))+4, uint64(uint32(s.x[15])), 4); err != nil {
		return 0, err
	}

	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) & uint64(0xff))))

	s.x[15] = uint64(uint32(uint64(0x65)))

	s.x[16] = uint64(uint32(uint64(0xa6)))

	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) - uint64(uint32(s.x[13])))))

	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[16])))))

	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) * uint64(uint32(s.x[12])))))

	s.x[12] = uint64(0x100fe9000)

	s.x[12] = uint64((s.x[12] + uint64(0x390)))

	s.x[12] = uint64((s.x[12] + uint64(uint8(uint64(uint32(s.x[9]))))))

	value59, err := s.load(s.x[12], 1)
	if err != nil {
		return 0, err
	}
	s.x[12] = uint64(uint32(value59))

	s.x[13] = uint64(uint32(uint64(0x2a)))

	s.x[12] = uint64(uint32((uint64(uint32(s.x[12])) ^ uint64(uint32(s.x[13])))))

	s.x[13] = uint64(0x100de2000)

	s.x[13] = uint64((s.x[13] + uint64(0xaf0)))

	value64, err := s.load((s.x[13] + uint64(uint32(uint64(uint32(s.x[12]))))), 1)
	if err != nil {
		return 0, err
	}
	s.x[13] = uint64(uint32(value64))

	s.x[12] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[15])))))

	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[13])))))

	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) & uint64(0xff))))

	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) + uint64(0x215))))

	s.x[10] = uint64(0x10105b000)

	s.x[10] = uint64((s.x[10] + uint64(0x6f0)))

	value71, err := s.load((s.x[10] + shiftLeft(uint64(uint32(uint64(uint32(s.x[9])))), 3, 64)), 8)
	if err != nil {
		return 0, err
	}
	s.x[9] = uint64(value71)

	s.x[10] = uint64(uint64(0xffffffffffff5355))

	s.x[10] = uint64(insertBits(s.x[10], 0x9824, 16, 16, 64))

	s.x[9] = uint64((s.x[9] + s.x[10]))

	value75, err := s.load((s.x[9] + shiftLeft(signExtend32(uint32(uint64(uint32(s.x[14])))), 3, 64)), 8)
	if err != nil {
		return 0, err
	}
	s.x[10] = uint64(value75)

	s.x[17] = uint64(uint64(0x6b95))

	next = s.x[10]
	return next, nil
}
