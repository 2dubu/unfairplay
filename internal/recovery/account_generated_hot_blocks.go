package recovery

const (
	accountWideCallSetup           = uint64(0x100115bbc)
	accountWideCallTail            = uint64(0x100115c00)
	accountDigestStateSetup        = uint64(0x100160408)
	accountDigestRoute             = uint64(0x10018efb8)
	accountByteConditionSetup      = uint64(0x1001d3548)
	accountDigestMixFront          = uint64(0x100243704)
	accountDigestMixMiddle         = uint64(0x100243800)
	accountDigestMixTail           = uint64(0x100243c00)
	accountByteCompareFront        = uint64(0x10030d3e0)
	accountByteCompareDispatchTail = uint64(0x10030d400)
	accountLoopLimitRoute          = uint64(0x10097a170)
)

func runAccountWideCallSetup(s *state) (uint64, error) {
	s.traceBlock(0x100115bbc)

	s.x[11] = uint64(uint32(uint64(0x18fa)))

	s.x[11] = uint64(uint32(insertBits(uint64(uint32(s.x[11])), 0xad0, 16, 16, 32)))

	s.x[17] = uint64(uint32((uint64(uint32(s.x[15])) + uint64(uint32(s.x[11])))))

	value3, err := s.load(addOffset(s.sp, 756), 4)
	if err != nil {
		return 0, err
	}
	s.x[16] = uint64(uint32(value3))

	value4, err := s.load((s.x[24] + shiftLeft(signExtend32(uint32(uint64(uint32(s.x[17])))), 3, 64)), 8)
	if err != nil {
		return 0, err
	}
	s.x[3] = uint64(value4)

	s.x[11] = uint64(uint32(uint64(0xf52c)))

	s.x[11] = uint64(uint32(insertBits(uint64(uint32(s.x[11])), 0xacf, 16, 16, 32)))

	s.x[15] = uint64(uint32((uint64(uint32(s.x[15])) + uint64(uint32(s.x[11])))))

	s.x[11] = uint64(uint32(uint64(0x18ab)))

	s.x[0] = uint64(uint32((uint64(uint32(s.x[15])) ^ uint64(uint32(s.x[11])))))

	s.x[2] = uint64(uint32(uint64(0x90d7)))

	s.x[2] = uint64(uint32(insertBits(uint64(uint32(s.x[2])), 0x481b, 16, 16, 32)))

	s.x[5] = uint64(uint32(uint64(0x7bd5)))

	s.x[5] = uint64(uint32(insertBits(uint64(uint32(s.x[5])), 0x4e75, 16, 16, 32)))

	s.x[7] = uint64(uint32(uint64(0x1821)))

	s.x[7] = uint64(uint32(insertBits(uint64(uint32(s.x[7])), 0x5b5a, 16, 16, 32)))

	s.x[30] = uint64(uint32(uint64(0x129)))
	return 0x100115c00, nil
}

func runAccountWideCallTail(s *state) (uint64, error) {
	s.traceBlock(0x100115c00)
	var next uint64

	s.x[17] = uint64(uint64(0x426))

	next = s.x[3]
	return next, nil
}

func runAccountDigestStateSetupCircuit(s *state) (uint64, error) {
	s.traceBlock(0x100160408)
	var next uint64

	s.x[11] = uint64(uint32(uint64(0x723d)))

	s.x[11] = uint64(uint32(insertBits(uint64(uint32(s.x[11])), 0xf73d, 16, 16, 32)))

	s.x[15] = uint64(uint32((uint64(uint32(s.x[0])) + uint64(uint32(s.x[11])))))

	s.x[11] = uint64(uint32(uint64(0x33df)))

	s.x[11] = uint64(uint32(insertBits(uint64(uint32(s.x[11])), 0xfdf2, 16, 16, 32)))

	s.x[15] = uint64(uint32((uint64(uint32(s.x[15])) & uint64(uint32(s.x[11])))))

	value6, err := s.load(addOffset(s.sp, 756), 4)
	if err != nil {
		return 0, err
	}
	s.x[17] = uint64(uint32(value6))

	s.x[11] = uint64(uint32(uint64(0xfa19)))

	s.x[11] = uint64(uint32(insertBits(uint64(uint32(s.x[11])), 0x48c8, 16, 16, 32)))

	s.x[2] = uint64(uint32((uint64(uint32(s.x[17])) ^ uint64(uint32(s.x[11])))))

	s.x[11] = uint64(uint32(uint64(0x5e7)))

	s.x[11] = uint64(uint32(insertBits(uint64(uint32(s.x[11])), 0xb737, 16, 16, 32)))

	s.x[2] = uint64(uint32((uint64(uint32(s.x[2])) + uint64(uint32(s.x[11])))))

	s.x[11] = uint64(uint32(uint64(0x98)))

	s.x[11] = uint64(uint32(insertBits(uint64(uint32(s.x[11])), 0x526e, 16, 16, 32)))

	s.x[3] = uint64(uint32((uint64(uint32(s.x[17])) ^ uint64(uint32(s.x[11])))))

	s.x[11] = uint64(uint32(uint64(0xff68)))

	s.x[11] = uint64(uint32(insertBits(uint64(uint32(s.x[11])), 0xad91, 16, 16, 32)))

	s.x[3] = uint64(uint32((uint64(uint32(s.x[3])) + uint64(uint32(s.x[11])))))

	s.x[2] = uint64(uint32((uint64(uint32(s.x[2])) ^ uint64(uint32(s.x[3])))))

	s.x[11] = uint64(uint32(uint64(0x373c)))

	s.x[11] = uint64(uint32(insertBits(uint64(uint32(s.x[11])), 0x11b9, 16, 16, 32)))

	s.x[17] = uint64(uint32((uint64(uint32(s.x[17])) ^ uint64(uint32(s.x[11])))))

	s.x[11] = uint64(uint32(uint64(0xeb4b)))

	s.x[11] = uint64(uint32(insertBits(uint64(uint32(s.x[11])), 0x1b76, 16, 16, 32)))

	s.x[15] = uint64(uint32((uint64(uint32(s.x[15])) ^ uint64(uint32(s.x[11])))))

	s.x[15] = uint64(uint32((uint64(uint32(s.x[15])) + uint64(uint32(s.x[17])))))

	s.x[15] = uint64(uint32((uint64(uint32(s.x[2])) ^ uint64(uint32(s.x[15])))))

	s.x[11] = uint64(uint32(uint64(0x3811)))

	s.x[11] = uint64(uint32(insertBits(uint64(uint32(s.x[11])), 0x4a97, 16, 16, 32)))

	s.x[17] = uint64(uint32((uint64(uint32(s.x[16])) ^ uint64(uint32(s.x[11])))))

	s.x[11] = uint64(uint32(uint64(0xc7ef)))

	s.x[11] = uint64(uint32(insertBits(uint64(uint32(s.x[11])), 0xb568, 16, 16, 32)))

	s.x[17] = uint64(uint32((uint64(uint32(s.x[17])) + uint64(uint32(s.x[11])))))

	s.x[11] = uint64(uint32(uint64(0xfdea)))

	s.x[11] = uint64(uint32(insertBits(uint64(uint32(s.x[11])), 0xb190, 16, 16, 32)))

	s.x[2] = uint64(uint32((uint64(uint32(s.x[16])) ^ uint64(uint32(s.x[11])))))

	s.x[11] = uint64(uint32(uint64(0x216)))

	s.x[11] = uint64(uint32(insertBits(uint64(uint32(s.x[11])), 0x4e6f, 16, 16, 32)))

	s.x[2] = uint64(uint32((uint64(uint32(s.x[2])) + uint64(uint32(s.x[11])))))

	s.x[17] = uint64(uint32((uint64(uint32(s.x[17])) ^ uint64(uint32(s.x[2])))))

	s.x[11] = uint64(uint32(uint64(0x846)))

	s.x[11] = uint64(uint32(insertBits(uint64(uint32(s.x[11])), 0xf018, 16, 16, 32)))

	s.x[16] = uint64(uint32((uint64(uint32(s.x[16])) ^ uint64(uint32(s.x[11])))))

	s.x[11] = uint64(uint32(uint64(0xf7ba)))

	s.x[11] = uint64(uint32(insertBits(uint64(uint32(s.x[11])), 0xfe7, 16, 16, 32)))

	s.x[16] = uint64(uint32((uint64(uint32(s.x[16])) + uint64(uint32(s.x[11])))))

	s.x[16] = uint64(uint32((uint64(uint32(s.x[17])) ^ uint64(uint32(s.x[16])))))

	s.x[11] = uint64(uint32(uint64(0xcdca)))

	s.x[11] = uint64(uint32(insertBits(uint64(uint32(s.x[11])), 0x604c, 16, 16, 32)))

	s.x[16] = uint64(uint32((uint64(uint32(s.x[16])) + uint64(uint32(s.x[11])))))

	s.x[15] = uint64(uint32((uint64(uint32(s.x[16])) + uint64(uint32(s.x[15])))))

	s.x[11] = uint64(uint32(uint64(0x889e)))

	s.x[11] = uint64(uint32(insertBits(uint64(uint32(s.x[11])), 0xe72c, 16, 16, 32)))

	s.x[16] = uint64(uint32((uint64(uint32(s.x[15])) ^ uint64(uint32(s.x[11])))))

	s.x[11] = uint64(uint32(uint64(0xcdb0)))

	s.x[11] = uint64(uint32(insertBits(uint64(uint32(s.x[11])), 0xb5f2, 16, 16, 32)))

	s.x[17] = uint64(uint32((uint64(uint32(s.x[15])) & uint64(uint32(s.x[11])))))

	s.x[11] = uint64(uint32(uint64(0x89be)))

	s.x[11] = uint64(uint32(insertBits(uint64(uint32(s.x[11])), 0xa5e8, 16, 16, 32)))

	s.x[2] = uint64(uint32((uint64(uint32(s.x[15])) ^ uint64(uint32(s.x[11])))))

	s.x[2] = uint64(uint32((uint64(uint32(s.x[2])) & shiftLeft(uint64(uint32(s.x[17])), 1, 32))))

	s.x[17] = uint64(uint32((uint64(uint32(s.x[2])) ^ uint64(uint32(s.x[17])))))

	s.x[2] = uint64(uint32(shiftLeft(uint64(uint32(s.x[16])), int(uint64(0x1)), 32)))

	s.x[11] = uint64(uint32(uint64(0x8a5c)))

	s.x[11] = uint64(uint32(insertBits(uint64(uint32(s.x[11])), 0xa5bc, 16, 16, 32)))

	s.x[3] = uint64(uint32((uint64(uint32(s.x[2])) ^ uint64(uint32(s.x[11])))))

	s.x[16] = uint64(uint32((uint64(uint32(s.x[3])) & uint64(uint32(s.x[16])))))

	s.x[11] = uint64(uint32(uint64(0x452e)))

	s.x[11] = uint64(uint32(insertBits(uint64(uint32(s.x[11])), 0x52de, 16, 16, 32)))

	s.x[2] = uint64(uint32((uint64(uint32(s.x[2])) & uint64(uint32(s.x[11])))))

	s.x[16] = uint64(uint32((uint64(uint32(s.x[16])) ^ uint64(uint32(s.x[2])))))

	s.x[13] = uint64(uint32(uint64(0x4522)))

	s.x[13] = uint64(uint32(insertBits(uint64(uint32(s.x[13])), 0x5242, 16, 16, 32)))

	s.x[2] = uint64(uint32((uint64(uint32(s.x[16])) ^ uint64(uint32(s.x[13])))))

	s.x[13] = uint64(uint32(uint64(0x9c0000)))

	s.x[16] = uint64(uint32((uint64(uint32(s.x[16])) ^ uint64(uint32(s.x[13])))))

	s.x[16] = uint64(uint32((uint64(uint32(s.x[16])) & shiftLeft(uint64(uint32(s.x[17])), 2, 32))))

	s.x[16] = uint64(uint32((uint64(uint32(s.x[16])) ^ uint64(uint32(s.x[17])))))

	s.x[17] = uint64(uint32(shiftLeft(uint64(uint32(s.x[2])), int(uint64(0x2)), 32)))

	s.x[13] = uint64(uint32(uint64(0x14b8)))

	s.x[13] = uint64(uint32(insertBits(uint64(uint32(s.x[13])), 0x4b79, 16, 16, 32)))

	s.x[3] = uint64(uint32((uint64(uint32(s.x[17])) ^ uint64(uint32(s.x[13])))))

	s.x[2] = uint64(uint32((uint64(uint32(s.x[3])) & uint64(uint32(s.x[2])))))

	s.x[13] = uint64(uint32(uint64(0x4528)))

	s.x[13] = uint64(uint32(insertBits(uint64(uint32(s.x[13])), 0x52de, 16, 16, 32)))

	s.x[17] = uint64(uint32((uint64(uint32(s.x[17])) & uint64(uint32(s.x[13])))))

	s.x[17] = uint64(uint32((uint64(uint32(s.x[2])) ^ uint64(uint32(s.x[17])))))

	s.x[13] = uint64(uint32(uint64(0x4106)))

	s.x[13] = uint64(uint32(insertBits(uint64(uint32(s.x[13])), 0x1086, 16, 16, 32)))

	s.x[2] = uint64(uint32((uint64(uint32(s.x[17])) ^ uint64(uint32(s.x[13])))))

	s.x[13] = uint64(uint32(uint64(0x420)))

	s.x[13] = uint64(uint32(insertBits(uint64(uint32(s.x[13])), 0x4258, 16, 16, 32)))

	s.x[17] = uint64(uint32((uint64(uint32(s.x[17])) ^ uint64(uint32(s.x[13])))))

	s.x[17] = uint64(uint32((uint64(uint32(s.x[17])) & shiftLeft(uint64(uint32(s.x[16])), 4, 32))))

	s.x[16] = uint64(uint32((uint64(uint32(s.x[17])) ^ uint64(uint32(s.x[16])))))

	s.x[17] = uint64(uint32(shiftLeft(uint64(uint32(s.x[2])), int(uint64(0x4)), 32)))

	s.x[13] = uint64(uint32(uint64(0x52e0)))

	s.x[13] = uint64(uint32(insertBits(uint64(uint32(s.x[13])), 0x2de4, 16, 16, 32)))

	s.x[3] = uint64(uint32((uint64(uint32(s.x[17])) ^ uint64(uint32(s.x[13])))))

	s.x[2] = uint64(uint32((uint64(uint32(s.x[3])) & uint64(uint32(s.x[2])))))

	s.x[13] = uint64(uint32(uint64(0x4520)))

	s.x[13] = uint64(uint32(insertBits(uint64(uint32(s.x[13])), 0x52de, 16, 16, 32)))

	s.x[17] = uint64(uint32((uint64(uint32(s.x[17])) & uint64(uint32(s.x[13])))))

	s.x[17] = uint64(uint32((uint64(uint32(s.x[2])) ^ uint64(uint32(s.x[17])))))

	s.x[13] = uint64(uint32(uint64(0x50e)))

	s.x[13] = uint64(uint32(insertBits(uint64(uint32(s.x[13])), 0x521a, 16, 16, 32)))

	s.x[2] = uint64(uint32((uint64(uint32(s.x[17])) ^ uint64(uint32(s.x[13])))))

	s.x[13] = uint64(uint32(uint64(0x4000)))

	s.x[13] = uint64(uint32(insertBits(uint64(uint32(s.x[13])), 0xc4, 16, 16, 32)))

	s.x[17] = uint64(uint32((uint64(uint32(s.x[17])) ^ uint64(uint32(s.x[13])))))

	s.x[17] = uint64(uint32((uint64(uint32(s.x[17])) & shiftLeft(uint64(uint32(s.x[16])), 8, 32))))

	s.x[16] = uint64(uint32((uint64(uint32(s.x[16])) ^ uint64(uint32(s.x[11])))))

	s.x[16] = uint64(uint32((uint64(uint32(s.x[16])) ^ uint64(uint32(s.x[17])))))

	s.x[17] = uint64(uint32(shiftLeft(uint64(uint32(s.x[2])), int(uint64(0x8)), 32)))

	s.x[11] = uint64(uint32(uint64(0x5e450000)))

	s.x[3] = uint64(uint32((uint64(uint32(s.x[17])) ^ uint64(uint32(s.x[11])))))

	s.x[2] = uint64(uint32((uint64(uint32(s.x[3])) & uint64(uint32(s.x[2])))))

	s.x[11] = uint64(uint32(uint64(0x52de0000)))

	s.x[17] = uint64(uint32((uint64(uint32(s.x[17])) & uint64(uint32(s.x[11])))))

	s.x[13] = uint64(uint32(uint64(0x9a0000)))

	s.x[17] = uint64(uint32((uint64(uint32(s.x[17])) ^ uint64(uint32(s.x[13])))))

	s.x[17] = uint64(uint32((uint64(uint32(s.x[17])) ^ uint64(uint32(s.x[2])))))

	s.x[2] = uint64(uint32(shiftLeft(uint64(uint32(s.x[16])), int(uint64(0x10)), 32)))

	s.x[13] = uint64(uint32(uint64(0x452e0000)))

	s.x[3] = uint64(uint32((uint64(uint32(s.x[2])) ^ uint64(uint32(s.x[13])))))

	s.x[17] = uint64(uint32((uint64(uint32(s.x[3])) & uint64(uint32(s.x[17])))))

	s.x[2] = uint64(uint32((uint64(uint32(s.x[2])) & uint64(uint32(s.x[11])))))

	s.x[16] = uint64(uint32((uint64(uint32(s.x[2])) ^ uint64(uint32(s.x[16])))))

	s.x[16] = uint64(uint32((uint64(uint32(s.x[16])) ^ uint64(uint32(s.x[17])))))

	s.x[15] = uint64(uint32((uint64(uint32(s.x[15])) ^ shiftLeft(uint64(uint32(s.x[16])), 1, 32))))

	s.x[11] = uint64(uint32(uint64(0x8a51)))

	s.x[11] = uint64(uint32(insertBits(uint64(uint32(s.x[11])), 0x9b4d, 16, 16, 32)))

	s.x[15] = uint64(uint32((uint64(uint32(s.x[15])) ^ uint64(uint32(s.x[11])))))

	if err := s.storeValue(addOffset(s.sp, 756), uint64(uint32(s.x[15])), 4); err != nil {
		return 0, err
	}

	value135, err := s.load((s.x[24] + shiftLeft(signExtend32(uint32(uint64(uint32(s.x[0])))), 3, 64)), 8)
	if err != nil {
		return 0, err
	}
	s.x[15] = uint64(value135)

	value136, err := s.load(addOffset(s.sp, 56), 8)
	if err != nil {
		return 0, err
	}
	s.x[21] = uint64(value136)

	s.x[17] = uint64(uint64(0x9969))

	next = s.x[15]
	return next, nil
}

func runAccountDigestRoute(s *state) (uint64, error) {
	s.traceBlock(0x10018efb8)
	var next uint64

	s.x[13] = uint64(uint32(uint64(0x7bc5)))

	s.x[13] = uint64(uint32(insertBits(uint64(uint32(s.x[13])), 0xce75, 16, 16, 32)))

	s.x[11] = uint64(uint32(uint64(0xffff8931)))

	s.x[17] = uint64(uint32((uint64(uint32(s.x[15])) + uint64(uint32(s.x[11])))))

	s.x[11] = uint64(uint32(uint64(0x2509)))

	s.x[0] = uint64(uint32((uint64(uint32(s.x[17])) | uint64(uint32(s.x[11])))))

	s.x[2] = uint64(uint32((uint64(uint32(s.x[0])) + uint64(0x59))))

	s.x[11] = uint64(uint32(uint64(0x94)))

	s.x[17] = uint64(uint32((uint64(uint32(s.x[2])) ^ uint64(uint32(s.x[11])))))

	s.compare32(uint32(uint64(uint32(s.x[16]))), uint32(uint64(int64(int8(uint64(uint32(s.x[17])))))))

	s.x[17] = uint64(uint32(boolValue(s.condition("lt"))))

	s.x[11] = uint64(uint32(uint64(0x35)))

	s.x[17] = uint64(uint32((uint64(uint32(s.x[17])) * uint64(uint32(s.x[11])))))

	s.x[15] = uint64(uint32((uint64(uint32(s.x[17])) ^ uint64(uint32(s.x[15])))))

	value14, err := s.load((s.x[24] + shiftLeft(signExtend32(uint32(uint64(uint32(s.x[15])))), 3, 64)), 8)
	if err != nil {
		return 0, err
	}
	s.x[15] = uint64(value14)

	s.x[17] = uint64(uint64(0xfd8a))

	next = s.x[15]
	return next, nil
}

func runAccountByteConditionSetup(s *state) (uint64, error) {
	s.traceBlock(0x1001d3548)
	var next uint64

	s.x[8] = uint64(uint32(uint64(0x6932)))

	s.x[8] = uint64(uint32(insertBits(uint64(uint32(s.x[8])), 0xad0, 16, 16, 32)))

	s.x[16] = uint64(uint32((uint64(uint32(s.x[20])) + uint64(uint32(s.x[8])))))

	s.x[8] = uint64(uint32(uint64(0x421e)))

	s.x[15] = uint64(uint32((uint64(uint32(s.x[16])) | uint64(uint32(s.x[8])))))

	value5, err := s.load((s.x[24] + shiftLeft(signExtend32(uint32(uint64(uint32(s.x[2])))), 3, 64)), 8)
	if err != nil {
		return 0, err
	}
	s.x[8] = uint64(value5)

	if err := s.storeValue(addOffset(s.sp, 96), s.x[8], 8); err != nil {
		return 0, err
	}

	s.x[8] = uint64(uint32(uint64(0x5d46)))

	s.x[8] = uint64(uint32(insertBits(uint64(uint32(s.x[8])), 0xad0, 16, 16, 32)))

	s.x[17] = uint64(uint32((uint64(uint32(s.x[20])) + uint64(uint32(s.x[8])))))

	s.x[8] = uint64((s.sp + uint64(0xc0)))

	s.x[20] = uint64((s.x[8] + s.x[25]))

	s.x[8] = uint64((s.x[25] + uint64(0x140)))

	if err := s.storeValue(addOffset(s.sp, 144), s.x[8], 8); err != nil {
		return 0, err
	}

	value14, err := s.load(addOffset(s.sp, 740), 4)
	if err != nil {
		return 0, err
	}
	s.x[8] = uint64(uint32(value14))

	value15, err := s.load(addOffset(s.sp, 744), 4)
	if err != nil {
		return 0, err
	}
	s.x[12] = uint64(uint32(value15))

	value16, err := s.load(addOffset(s.sp, 748), 4)
	if err != nil {
		return 0, err
	}
	s.x[19] = uint64(uint32(value16))

	value17, err := s.load(addOffset(s.sp, 752), 4)
	if err != nil {
		return 0, err
	}
	s.x[10] = uint64(uint32(value17))

	value18, err := s.load(addOffset(s.sp, 764), 1)
	if err != nil {
		return 0, err
	}
	s.x[0] = uint64(uint32(value18))

	s.x[16] = uint64(uint32((uint64(uint32(s.x[16])) ^ uint64(0xc0))))

	s.x[11] = uint64(uint32(uint64(0xe1)))

	s.x[16] = uint64((s.x[16] & s.x[11]))

	s.x[16] = uint64(uint32(shiftLeft(uint64(uint32(s.x[0])), int(uint64(uint32(s.x[16]))), 32)))

	s.x[16] = uint64(uint32((uint64(uint32(s.x[16])) & uint64(0xc0))))

	s.x[16] = uint64(uint32((uint64(uint32(s.x[0])) - uint64(uint32(s.x[16])))))

	s.x[16] = uint64(uint32((uint64(uint32(s.x[16])) - uint64(0x1c))))

	s.x[16] = uint64(uint32(signExtendBits(uint64(uint32(s.x[16])), 8, 32)))

	value27, err := s.load((s.x[24] + shiftLeft(signExtend32(uint32(uint64(uint32(s.x[17])))), 3, 64)), 8)
	if err != nil {
		return 0, err
	}
	s.x[0] = uint64(value27)

	s.x[17] = uint64(uint64(0x9151))

	next = s.x[0]
	return next, nil
}

func runAccountDigestMixFront(s *state) (uint64, error) {
	s.traceBlock(0x100243704)

	value0, err := s.load(addOffset(s.sp, 740), 4)
	if err != nil {
		return 0, err
	}
	s.x[11] = uint64(uint32(value0))

	value1, err := s.load(addOffset(s.sp, 744), 4)
	if err != nil {
		return 0, err
	}
	s.x[13] = uint64(uint32(value1))

	s.x[15] = uint64(uint32(uint64(0x5cf7)))

	s.x[15] = uint64(uint32(insertBits(uint64(uint32(s.x[15])), 0x4d37, 16, 16, 32)))

	s.x[15] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[15])))))

	s.x[16] = uint64(uint32(uint64(0xa309)))

	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0xb2c8, 16, 16, 32)))

	s.x[15] = uint64(uint32((uint64(uint32(s.x[15])) + uint64(uint32(s.x[16])))))

	s.x[16] = uint64(uint32(uint64(0x5e28)))

	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0xe5d7, 16, 16, 32)))

	s.x[16] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[16])))))

	s.x[17] = uint64(uint32(uint64(0xa1d8)))

	s.x[17] = uint64(uint32(insertBits(uint64(uint32(s.x[17])), 0x1a28, 16, 16, 32)))

	s.x[16] = uint64(uint32((uint64(uint32(s.x[16])) + uint64(uint32(s.x[17])))))

	s.x[15] = uint64(uint32((uint64(uint32(s.x[15])) ^ uint64(uint32(s.x[16])))))

	s.x[16] = uint64(uint32(uint64(0xcf62)))

	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0xa3ff, 16, 16, 32)))

	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[16])))))

	s.x[16] = uint64(uint32(uint64(0x309e)))

	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0x5c00, 16, 16, 32)))

	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) + uint64(uint32(s.x[16])))))

	s.x[11] = uint64(uint32((uint64(uint32(s.x[15])) ^ uint64(uint32(s.x[11])))))

	s.x[15] = uint64(uint32(uint64(0x94ac)))

	s.x[15] = uint64(uint32(insertBits(uint64(uint32(s.x[15])), 0x9dd5, 16, 16, 32)))

	s.x[15] = uint64(uint32((uint64(uint32(s.x[8])) ^ uint64(uint32(s.x[15])))))

	s.x[16] = uint64(uint32(uint64(0x6b54)))

	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0x622a, 16, 16, 32)))

	s.x[15] = uint64(uint32((uint64(uint32(s.x[15])) + uint64(uint32(s.x[16])))))

	s.x[16] = uint64(uint32(uint64(0xbb30)))

	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0x324c, 16, 16, 32)))

	s.x[16] = uint64(uint32((uint64(uint32(s.x[8])) ^ uint64(uint32(s.x[16])))))

	s.x[17] = uint64(uint32(uint64(0x44d0)))

	s.x[17] = uint64(uint32(insertBits(uint64(uint32(s.x[17])), 0xcdb3, 16, 16, 32)))

	s.x[16] = uint64(uint32((uint64(uint32(s.x[16])) + uint64(uint32(s.x[17])))))

	s.x[15] = uint64(uint32((uint64(uint32(s.x[15])) ^ uint64(uint32(s.x[16])))))

	s.x[16] = uint64(uint32(uint64(0xe221)))

	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0xa486, 16, 16, 32)))

	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) ^ uint64(uint32(s.x[16])))))

	s.x[16] = uint64(uint32(uint64(0x1ddf)))

	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0x5b79, 16, 16, 32)))

	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) + uint64(uint32(s.x[16])))))

	s.x[8] = uint64(uint32((uint64(uint32(s.x[15])) ^ uint64(uint32(s.x[8])))))

	s.x[15] = uint64(uint32(uint64(0x82fd)))

	s.x[15] = uint64(uint32(insertBits(uint64(uint32(s.x[15])), 0x1b25, 16, 16, 32)))

	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) + uint64(uint32(s.x[15])))))

	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) + uint64(uint32(s.x[11])))))

	s.x[11] = uint64(uint32(uint64(0xb399)))

	s.x[11] = uint64(uint32(insertBits(uint64(uint32(s.x[11])), 0xcc91, 16, 16, 32)))

	s.x[11] = uint64(uint32((uint64(uint32(s.x[8])) ^ uint64(uint32(s.x[11])))))

	s.x[15] = uint64(uint32(uint64(0xabe6)))

	s.x[15] = uint64(uint32(insertBits(uint64(uint32(s.x[15])), 0x378b, 16, 16, 32)))

	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) + uint64(uint32(s.x[15])))))

	s.x[15] = uint64(uint32(uint64(0x4c30)))

	s.x[15] = uint64(uint32(insertBits(uint64(uint32(s.x[15])), 0x27ab, 16, 16, 32)))

	s.x[15] = uint64(uint32((uint64(uint32(s.x[8])) ^ uint64(uint32(s.x[15])))))

	s.x[16] = uint64(uint32(uint64(0x544d)))

	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0xdcb1, 16, 16, 32)))

	s.x[15] = uint64(uint32((uint64(uint32(s.x[15])) + uint64(uint32(s.x[16])))))

	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[15])))))

	s.x[15] = uint64(uint32(uint64(0xe3d5)))

	s.x[15] = uint64(uint32(insertBits(uint64(uint32(s.x[15])), 0x1120, 16, 16, 32)))

	s.x[15] = uint64(uint32((uint64(uint32(s.x[8])) ^ uint64(uint32(s.x[15])))))

	s.x[16] = uint64(uint32(uint64(0xfbaa)))
	return 0x100243800, nil
}

func runAccountDigestMixMiddle(s *state) (uint64, error) {
	s.traceBlock(0x100243800)

	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0xea3a, 16, 16, 32)))

	s.x[15] = uint64(uint32((uint64(uint32(s.x[15])) + uint64(uint32(s.x[16])))))

	s.x[15] = uint64(uint32((uint64(uint32(s.x[8])) ^ uint64(uint32(s.x[15])))))

	s.x[16] = uint64(uint32(uint64(0xfbff)))

	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0xfeff, 16, 16, 32)))

	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) ^ uint64(uint32(s.x[16])))))

	s.x[11] = uint64(uint32((uint64(uint32(s.x[15])) ^ uint64(uint32(s.x[11])))))

	s.x[15] = uint64(uint32(uint64(0xe384)))

	s.x[15] = uint64(uint32(insertBits(uint64(uint32(s.x[15])), 0x5e5, 16, 16, 32)))

	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) + uint64(uint32(s.x[15])))))

	s.x[15] = uint64(uint32(uint64(0x2a3e)))

	s.x[15] = uint64(uint32(insertBits(uint64(uint32(s.x[15])), 0xffa, 16, 16, 32)))

	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) ^ uint64(uint32(s.x[15])))))

	s.x[8] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[8])))))

	s.x[11] = uint64(uint32(uint64(0x9ecb)))

	s.x[11] = uint64(uint32(insertBits(uint64(uint32(s.x[11])), 0x9247, 16, 16, 32)))

	s.x[11] = uint64(uint32((uint64(uint32(s.x[13])) ^ uint64(uint32(s.x[11])))))

	s.x[15] = uint64(uint32(uint64(0x6135)))

	s.x[15] = uint64(uint32(insertBits(uint64(uint32(s.x[15])), 0x6db8, 16, 16, 32)))

	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) + uint64(uint32(s.x[15])))))

	s.x[15] = uint64(uint32(uint64(0x1945)))

	s.x[15] = uint64(uint32(insertBits(uint64(uint32(s.x[15])), 0xa3c2, 16, 16, 32)))

	s.x[15] = uint64(uint32((uint64(uint32(s.x[13])) ^ uint64(uint32(s.x[15])))))

	s.x[16] = uint64(uint32(uint64(0xe6bb)))

	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0x5c3d, 16, 16, 32)))

	s.x[15] = uint64(uint32((uint64(uint32(s.x[15])) + uint64(uint32(s.x[16])))))

	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[15])))))

	s.x[15] = uint64(uint32(uint64(0x4a33)))

	s.x[15] = uint64(uint32(insertBits(uint64(uint32(s.x[15])), 0x3a9a, 16, 16, 32)))

	s.x[13] = uint64(uint32((uint64(uint32(s.x[13])) ^ uint64(uint32(s.x[15])))))

	s.x[15] = uint64(uint32(uint64(0xb5cd)))

	s.x[15] = uint64(uint32(insertBits(uint64(uint32(s.x[15])), 0xc565, 16, 16, 32)))

	s.x[13] = uint64(uint32((uint64(uint32(s.x[13])) + uint64(uint32(s.x[15])))))

	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[13])))))

	s.x[13] = uint64(uint32(uint64(0x6369)))

	s.x[13] = uint64(uint32(insertBits(uint64(uint32(s.x[13])), 0x73aa, 16, 16, 32)))

	s.x[13] = uint64(uint32((uint64(uint32(s.x[12])) ^ uint64(uint32(s.x[13])))))

	s.x[15] = uint64(uint32(uint64(0x9c97)))

	s.x[15] = uint64(uint32(insertBits(uint64(uint32(s.x[15])), 0x8c55, 16, 16, 32)))

	s.x[13] = uint64(uint32((uint64(uint32(s.x[13])) + uint64(uint32(s.x[15])))))

	s.x[15] = uint64(uint32(uint64(0xd2f9)))

	s.x[15] = uint64(uint32(insertBits(uint64(uint32(s.x[15])), 0x3fdf, 16, 16, 32)))

	s.x[15] = uint64(uint32((uint64(uint32(s.x[12])) ^ uint64(uint32(s.x[15])))))

	s.x[16] = uint64(uint32(uint64(0x2d07)))

	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0xc020, 16, 16, 32)))

	s.x[15] = uint64(uint32((uint64(uint32(s.x[15])) + uint64(uint32(s.x[16])))))

	s.x[13] = uint64(uint32((uint64(uint32(s.x[13])) ^ uint64(uint32(s.x[15])))))

	s.x[15] = uint64(uint32(uint64(0x7c2d)))

	s.x[15] = uint64(uint32(insertBits(uint64(uint32(s.x[15])), 0x476a, 16, 16, 32)))

	s.x[12] = uint64(uint32((uint64(uint32(s.x[12])) ^ uint64(uint32(s.x[15])))))

	s.x[15] = uint64(uint32(uint64(0x83d3)))

	s.x[15] = uint64(uint32(insertBits(uint64(uint32(s.x[15])), 0xb895, 16, 16, 32)))

	s.x[12] = uint64(uint32((uint64(uint32(s.x[12])) + uint64(uint32(s.x[15])))))

	s.x[12] = uint64(uint32((uint64(uint32(s.x[13])) ^ uint64(uint32(s.x[12])))))

	s.x[13] = uint64(uint32(uint64(0x9dfa)))

	s.x[13] = uint64(uint32(insertBits(uint64(uint32(s.x[13])), 0x7d6c, 16, 16, 32)))

	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) + uint64(uint32(s.x[13])))))

	s.x[11] = uint64(uint32((uint64(uint32(s.x[12])) + uint64(uint32(s.x[11])))))

	s.x[12] = uint64(uint32(uint64(0x541)))

	s.x[12] = uint64(uint32(insertBits(uint64(uint32(s.x[12])), 0xda5a, 16, 16, 32)))

	s.x[12] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[12])))))

	s.x[13] = uint64(uint32(uint64(0xaed)))

	s.x[13] = uint64(uint32(insertBits(uint64(uint32(s.x[13])), 0x121f, 16, 16, 32)))

	s.x[13] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[13])))))

	s.x[15] = uint64(uint32(uint64(0xf793)))

	s.x[15] = uint64(uint32(insertBits(uint64(uint32(s.x[15])), 0x8acd, 16, 16, 32)))

	s.x[13] = uint64(uint32((uint64(uint32(s.x[13])) + uint64(uint32(s.x[15])))))

	s.x[15] = uint64(uint32(uint64(0xf2db)))

	s.x[15] = uint64(uint32(insertBits(uint64(uint32(s.x[15])), 0x5097, 16, 16, 32)))

	s.x[15] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[15])))))

	s.x[16] = uint64(uint32(uint64(0xf83f)))

	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0x4288, 16, 16, 32)))

	s.x[12] = uint64(uint32((uint64(uint32(s.x[12])) + uint64(uint32(s.x[16])))))

	s.x[12] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[12])))))

	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(0xfffffff7))))

	s.x[16] = uint64(uint32(uint64(0xfa5)))

	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0xc845, 16, 16, 32)))

	s.x[15] = uint64(uint32((uint64(uint32(s.x[15])) + uint64(uint32(s.x[16])))))

	s.x[13] = uint64(uint32((uint64(uint32(s.x[13])) ^ uint64(uint32(s.x[15])))))

	s.x[12] = uint64(uint32((uint64(uint32(s.x[12])) ^ uint64(uint32(s.x[13])))))

	s.x[13] = uint64(uint32(uint64(0x289)))

	s.x[13] = uint64(uint32(insertBits(uint64(uint32(s.x[13])), 0x672d, 16, 16, 32)))

	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) + uint64(uint32(s.x[13])))))

	s.x[13] = uint64(uint32(uint64(0xcf3d)))

	s.x[13] = uint64(uint32(insertBits(uint64(uint32(s.x[13])), 0x6c32, 16, 16, 32)))

	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[13])))))

	s.x[11] = uint64(uint32((uint64(uint32(s.x[12])) ^ uint64(uint32(s.x[11])))))

	if err := s.storeValue(addOffset(s.sp, 740), uint64(uint32(s.x[8])), 4); err != nil {
		return 0, err
	}

	if err := s.storeValue(addOffset(s.sp, 744), uint64(uint32(s.x[11])), 4); err != nil {
		return 0, err
	}

	value89, err := s.load(addOffset(s.sp, 748), 4)
	if err != nil {
		return 0, err
	}
	s.x[8] = uint64(uint32(value89))

	value90, err := s.load(addOffset(s.sp, 752), 4)
	if err != nil {
		return 0, err
	}
	s.x[12] = uint64(uint32(value90))

	s.x[11] = uint64(uint32(uint64(0xe68d)))

	s.x[11] = uint64(uint32(insertBits(uint64(uint32(s.x[11])), 0x8165, 16, 16, 32)))

	s.x[11] = uint64(uint32((uint64(uint32(s.x[8])) ^ uint64(uint32(s.x[11])))))

	s.x[13] = uint64(uint32(uint64(0x1973)))

	s.x[13] = uint64(uint32(insertBits(uint64(uint32(s.x[13])), 0x7e9a, 16, 16, 32)))

	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) + uint64(uint32(s.x[13])))))

	s.x[13] = uint64(uint32(uint64(0xbc58)))

	s.x[13] = uint64(uint32(insertBits(uint64(uint32(s.x[13])), 0xdccc, 16, 16, 32)))

	s.x[13] = uint64(uint32((uint64(uint32(s.x[8])) ^ uint64(uint32(s.x[13])))))

	s.x[15] = uint64(uint32(uint64(0x43a8)))

	s.x[15] = uint64(uint32(insertBits(uint64(uint32(s.x[15])), 0x2333, 16, 16, 32)))

	s.x[13] = uint64(uint32((uint64(uint32(s.x[13])) + uint64(uint32(s.x[15])))))

	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[13])))))

	s.x[13] = uint64(uint32(uint64(0x9768)))

	s.x[13] = uint64(uint32(insertBits(uint64(uint32(s.x[13])), 0x56b6, 16, 16, 32)))

	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) ^ uint64(uint32(s.x[13])))))

	s.x[13] = uint64(uint32(uint64(0x6898)))

	s.x[13] = uint64(uint32(insertBits(uint64(uint32(s.x[13])), 0xa949, 16, 16, 32)))

	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) + uint64(uint32(s.x[13])))))

	s.x[8] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[8])))))

	s.x[11] = uint64(uint32(uint64(0xca45)))

	s.x[11] = uint64(uint32(insertBits(uint64(uint32(s.x[11])), 0xde24, 16, 16, 32)))

	s.x[11] = uint64(uint32((uint64(uint32(s.x[19])) ^ uint64(uint32(s.x[11])))))

	s.x[13] = uint64(uint32(uint64(0x35bb)))

	s.x[13] = uint64(uint32(insertBits(uint64(uint32(s.x[13])), 0x21db, 16, 16, 32)))

	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) + uint64(uint32(s.x[13])))))

	s.x[13] = uint64(uint32(uint64(0x1f8d)))

	s.x[13] = uint64(uint32(insertBits(uint64(uint32(s.x[13])), 0x99c1, 16, 16, 32)))

	s.x[13] = uint64(uint32((uint64(uint32(s.x[19])) ^ uint64(uint32(s.x[13])))))

	s.x[15] = uint64(uint32(uint64(0xe073)))

	s.x[15] = uint64(uint32(insertBits(uint64(uint32(s.x[15])), 0x663e, 16, 16, 32)))

	s.x[13] = uint64(uint32((uint64(uint32(s.x[13])) + uint64(uint32(s.x[15])))))

	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[13])))))

	s.x[13] = uint64(uint32(uint64(0x1875)))

	s.x[13] = uint64(uint32(insertBits(uint64(uint32(s.x[13])), 0x4cfa, 16, 16, 32)))

	s.x[13] = uint64(uint32((uint64(uint32(s.x[19])) ^ uint64(uint32(s.x[13])))))

	s.x[15] = uint64(uint32(uint64(0xe78b)))

	s.x[15] = uint64(uint32(insertBits(uint64(uint32(s.x[15])), 0xb305, 16, 16, 32)))

	s.x[13] = uint64(uint32((uint64(uint32(s.x[13])) + uint64(uint32(s.x[15])))))

	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[13])))))

	s.x[13] = uint64(uint32(uint64(0xa863)))

	s.x[13] = uint64(uint32(insertBits(uint64(uint32(s.x[13])), 0x88a2, 16, 16, 32)))

	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) + uint64(uint32(s.x[13])))))

	s.x[8] = uint64(uint32((uint64(uint32(s.x[11])) + uint64(uint32(s.x[8])))))

	s.x[11] = uint64(uint32(uint64(0x65aa)))

	s.x[11] = uint64(uint32(insertBits(uint64(uint32(s.x[11])), 0xb798, 16, 16, 32)))

	s.x[11] = uint64(uint32((uint64(uint32(s.x[8])) ^ uint64(uint32(s.x[11])))))

	s.x[13] = uint64(uint32(uint64(0xf317)))

	s.x[13] = uint64(uint32(insertBits(uint64(uint32(s.x[13])), 0x8d9c, 16, 16, 32)))

	s.x[13] = uint64(uint32((uint64(uint32(s.x[8])) & uint64(uint32(s.x[13])))))

	s.x[15] = uint64(uint32(uint64(0xe386)))

	s.x[15] = uint64(uint32(insertBits(uint64(uint32(s.x[15])), 0xad98, 16, 16, 32)))

	s.x[15] = uint64(uint32((uint64(uint32(s.x[8])) ^ uint64(uint32(s.x[15])))))

	s.x[15] = uint64(uint32((uint64(uint32(s.x[15])) & shiftLeft(uint64(uint32(s.x[13])), 1, 32))))

	s.x[13] = uint64(uint32((uint64(uint32(s.x[15])) ^ uint64(uint32(s.x[13])))))

	s.x[15] = uint64(uint32(shiftLeft(uint64(uint32(s.x[11])), int(uint64(0x1)), 32)))

	s.x[16] = uint64(uint32(uint64(0x2d7a)))

	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0x7409, 16, 16, 32)))

	s.x[16] = uint64(uint32((uint64(uint32(s.x[15])) ^ uint64(uint32(s.x[16])))))

	s.x[11] = uint64(uint32((uint64(uint32(s.x[16])) & uint64(uint32(s.x[11])))))

	s.x[17] = uint64(uint32(uint64(0x96bc)))

	s.x[17] = uint64(uint32(insertBits(uint64(uint32(s.x[17])), 0x3a04, 16, 16, 32)))

	s.x[15] = uint64(uint32((uint64(uint32(s.x[15])) & uint64(uint32(s.x[17])))))

	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[15])))))

	s.x[15] = uint64(uint32(uint64(0x9285)))

	s.x[15] = uint64(uint32(insertBits(uint64(uint32(s.x[15])), 0xa04, 16, 16, 32)))

	s.x[15] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[15])))))

	s.x[16] = uint64(uint32(uint64(0x438)))

	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0x3000, 16, 16, 32)))

	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[16])))))

	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) & shiftLeft(uint64(uint32(s.x[13])), 2, 32))))

	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[13])))))

	s.x[13] = uint64(uint32(shiftLeft(uint64(uint32(s.x[15])), int(uint64(0x2)), 32)))

	s.x[16] = uint64(uint32(uint64(0x5af4)))

	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0xe812, 16, 16, 32)))

	s.x[16] = uint64(uint32((uint64(uint32(s.x[13])) ^ uint64(uint32(s.x[16])))))

	s.x[15] = uint64(uint32((uint64(uint32(s.x[16])) & uint64(uint32(s.x[15])))))

	s.x[13] = uint64(uint32((uint64(uint32(s.x[13])) & uint64(uint32(s.x[17])))))

	s.x[13] = uint64(uint32((uint64(uint32(s.x[15])) ^ uint64(uint32(s.x[13])))))

	s.x[15] = uint64(uint32(uint64(0x8409)))

	s.x[15] = uint64(uint32(insertBits(uint64(uint32(s.x[15])), 0x1204, 16, 16, 32)))

	s.x[15] = uint64(uint32((uint64(uint32(s.x[13])) ^ uint64(uint32(s.x[15])))))

	s.x[16] = uint64(uint32(uint64(0x12b0)))

	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0x2800, 16, 16, 32)))

	s.x[13] = uint64(uint32((uint64(uint32(s.x[13])) ^ uint64(uint32(s.x[16])))))

	s.x[13] = uint64(uint32((uint64(uint32(s.x[13])) & shiftLeft(uint64(uint32(s.x[11])), 4, 32))))

	s.x[11] = uint64(uint32((uint64(uint32(s.x[13])) ^ uint64(uint32(s.x[11])))))

	s.x[13] = uint64(uint32(shiftLeft(uint64(uint32(s.x[15])), int(uint64(0x4)), 32)))

	s.x[16] = uint64(uint32(uint64(0x6bd0)))

	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0xa049, 16, 16, 32)))

	s.x[16] = uint64(uint32((uint64(uint32(s.x[13])) ^ uint64(uint32(s.x[16])))))

	s.x[15] = uint64(uint32((uint64(uint32(s.x[16])) & uint64(uint32(s.x[15])))))

	s.x[16] = uint64(uint32(uint64(0x9690)))

	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0x3a04, 16, 16, 32)))

	s.x[13] = uint64(uint32((uint64(uint32(s.x[13])) & uint64(uint32(s.x[16])))))

	s.x[13] = uint64(uint32((uint64(uint32(s.x[15])) ^ uint64(uint32(s.x[13])))))

	s.x[15] = uint64(uint32(uint64(0x942d)))

	s.x[15] = uint64(uint32(insertBits(uint64(uint32(s.x[15])), 0x1a04, 16, 16, 32)))

	s.x[15] = uint64(uint32((uint64(uint32(s.x[13])) ^ uint64(uint32(s.x[15])))))

	s.x[16] = uint64(uint32(uint64(0x200)))

	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0x2000, 16, 16, 32)))

	s.x[13] = uint64(uint32((uint64(uint32(s.x[13])) ^ uint64(uint32(s.x[16])))))

	s.x[13] = uint64(uint32((uint64(uint32(s.x[13])) & shiftLeft(uint64(uint32(s.x[11])), 8, 32))))

	s.x[16] = uint64(uint32(uint64(0x96bd)))

	s.x[16] = uint64(uint32(insertBits(uint64(uint32(s.x[16])), 0x3a04, 16, 16, 32)))

	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[16])))))

	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[13])))))

	s.x[13] = uint64(uint32(shiftLeft(uint64(uint32(s.x[15])), int(uint64(0x8)), 32)))

	s.x[16] = uint64(uint32(uint64(0x4960000)))

	s.x[16] = uint64(uint32((uint64(uint32(s.x[13])) ^ uint64(uint32(s.x[16])))))

	s.x[15] = uint64(uint32((uint64(uint32(s.x[16])) & uint64(uint32(s.x[15])))))

	s.x[17] = uint64(uint32(uint64(0x3a040000)))

	s.x[13] = uint64(uint32((uint64(uint32(s.x[13])) & uint64(uint32(s.x[17])))))

	s.x[16] = uint64(uint32(uint64(0x3a000000)))

	s.x[13] = uint64(uint32((uint64(uint32(s.x[13])) ^ uint64(uint32(s.x[16])))))

	s.x[13] = uint64(uint32((uint64(uint32(s.x[13])) ^ uint64(uint32(s.x[15])))))

	s.x[15] = uint64(uint32(shiftLeft(uint64(uint32(s.x[11])), int(uint64(0x10)), 32)))

	s.x[16] = uint64(uint32(uint64(0x16bd0000)))

	s.x[16] = uint64(uint32((uint64(uint32(s.x[15])) ^ uint64(uint32(s.x[16])))))

	s.x[13] = uint64(uint32((uint64(uint32(s.x[16])) & uint64(uint32(s.x[13])))))

	s.x[15] = uint64(uint32((uint64(uint32(s.x[15])) & uint64(uint32(s.x[17])))))

	s.x[11] = uint64(uint32((uint64(uint32(s.x[15])) ^ uint64(uint32(s.x[11])))))

	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[13])))))

	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) ^ shiftLeft(uint64(uint32(s.x[11])), 1, 32))))

	s.x[11] = uint64(uint32(uint64(0xa072)))

	s.x[11] = uint64(uint32(insertBits(uint64(uint32(s.x[11])), 0xf731, 16, 16, 32)))

	s.x[11] = uint64(uint32((uint64(uint32(s.x[12])) ^ uint64(uint32(s.x[11])))))

	s.x[13] = uint64(uint32(uint64(0x5f8e)))

	s.x[13] = uint64(uint32(insertBits(uint64(uint32(s.x[13])), 0x8ce, 16, 16, 32)))

	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) + uint64(uint32(s.x[13])))))

	s.x[13] = uint64(uint32(uint64(0x6413)))

	s.x[13] = uint64(uint32(insertBits(uint64(uint32(s.x[13])), 0xa740, 16, 16, 32)))

	s.x[13] = uint64(uint32((uint64(uint32(s.x[12])) ^ uint64(uint32(s.x[13])))))

	s.x[15] = uint64(uint32(uint64(0x9bed)))

	s.x[15] = uint64(uint32(insertBits(uint64(uint32(s.x[15])), 0x58bf, 16, 16, 32)))

	s.x[13] = uint64(uint32((uint64(uint32(s.x[13])) + uint64(uint32(s.x[15])))))

	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[13])))))

	s.x[13] = uint64(uint32(uint64(0x9dc)))

	s.x[13] = uint64(uint32(insertBits(uint64(uint32(s.x[13])), 0x5b6e, 16, 16, 32)))

	s.x[12] = uint64(uint32((uint64(uint32(s.x[12])) ^ uint64(uint32(s.x[13])))))

	s.x[13] = uint64(uint32(uint64(0xf624)))

	s.x[13] = uint64(uint32(insertBits(uint64(uint32(s.x[13])), 0xa491, 16, 16, 32)))

	s.x[12] = uint64(uint32((uint64(uint32(s.x[12])) + uint64(uint32(s.x[13])))))

	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[12])))))

	s.x[12] = uint64(uint32(uint64(0x7d3c)))

	s.x[12] = uint64(uint32(insertBits(uint64(uint32(s.x[12])), 0xe3b3, 16, 16, 32)))

	s.x[12] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[12])))))

	s.x[13] = uint64(uint32(uint64(0x82c4)))

	s.x[13] = uint64(uint32(insertBits(uint64(uint32(s.x[13])), 0x1c4c, 16, 16, 32)))

	s.x[12] = uint64(uint32((uint64(uint32(s.x[12])) + uint64(uint32(s.x[13])))))

	s.x[13] = uint64(uint32(uint64(0xaf00)))

	s.x[13] = uint64(uint32(insertBits(uint64(uint32(s.x[13])), 0xf6b1, 16, 16, 32)))

	s.x[13] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[13])))))

	s.x[15] = uint64(uint32(uint64(0x5100)))

	s.x[15] = uint64(uint32(insertBits(uint64(uint32(s.x[15])), 0x94e, 16, 16, 32)))

	s.x[13] = uint64(uint32((uint64(uint32(s.x[13])) + uint64(uint32(s.x[15])))))

	s.x[12] = uint64(uint32((uint64(uint32(s.x[12])) ^ uint64(uint32(s.x[13])))))

	s.x[13] = uint64(uint32(uint64(0x13d0)))

	s.x[13] = uint64(uint32(insertBits(uint64(uint32(s.x[13])), 0xd682, 16, 16, 32)))

	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) ^ uint64(uint32(s.x[13])))))

	s.x[13] = uint64(uint32(uint64(0x1f81)))

	s.x[13] = uint64(uint32(insertBits(uint64(uint32(s.x[13])), 0x1e1d, 16, 16, 32)))

	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[13])))))

	s.x[13] = uint64(uint32(uint64(0xe07f)))

	s.x[13] = uint64(uint32(insertBits(uint64(uint32(s.x[13])), 0xe1e2, 16, 16, 32)))
	return 0x100243c00, nil
}

func runAccountDigestMixTail(s *state) (uint64, error) {
	s.traceBlock(0x100243c00)
	var next uint64

	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) + uint64(uint32(s.x[13])))))

	s.x[10] = uint64(uint32((uint64(uint32(s.x[12])) ^ uint64(uint32(s.x[10])))))

	s.x[12] = uint64(uint32(uint64(0x7451)))

	s.x[12] = uint64(uint32(insertBits(uint64(uint32(s.x[12])), 0x2f9c, 16, 16, 32)))

	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) + uint64(uint32(s.x[12])))))

	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) + uint64(uint32(s.x[11])))))

	s.x[11] = uint64(uint32(uint64(0x9526)))

	s.x[11] = uint64(uint32(insertBits(uint64(uint32(s.x[11])), 0xa73f, 16, 16, 32)))

	s.x[11] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[11])))))

	s.x[12] = uint64(uint32(uint64(0xb20f)))

	s.x[12] = uint64(uint32(insertBits(uint64(uint32(s.x[12])), 0x419c, 16, 16, 32)))

	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) + uint64(uint32(s.x[12])))))

	s.x[12] = uint64(uint32(uint64(0x77a0)))

	s.x[12] = uint64(uint32(insertBits(uint64(uint32(s.x[12])), 0xb156, 16, 16, 32)))

	s.x[12] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[12])))))

	s.x[13] = uint64(uint32(uint64(0x5089)))

	s.x[13] = uint64(uint32(insertBits(uint64(uint32(s.x[13])), 0x57f5, 16, 16, 32)))

	s.x[12] = uint64(uint32((uint64(uint32(s.x[12])) + uint64(uint32(s.x[13])))))

	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[12])))))

	s.x[12] = uint64(uint32(uint64(0xc3ae)))

	s.x[12] = uint64(uint32(insertBits(uint64(uint32(s.x[12])), 0xf0da, 16, 16, 32)))

	s.x[12] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[12])))))

	s.x[13] = uint64(uint32(uint64(0xe487)))

	s.x[13] = uint64(uint32(insertBits(uint64(uint32(s.x[13])), 0x1679, 16, 16, 32)))

	s.x[12] = uint64(uint32((uint64(uint32(s.x[12])) + uint64(uint32(s.x[13])))))

	s.x[12] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[12])))))

	s.x[13] = uint64(uint32(uint64(0xf9ff)))

	s.x[13] = uint64(uint32(insertBits(uint64(uint32(s.x[13])), 0xffef, 16, 16, 32)))

	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[13])))))

	s.x[11] = uint64(uint32((uint64(uint32(s.x[12])) ^ uint64(uint32(s.x[11])))))

	s.x[12] = uint64(uint32(uint64(0xded8)))

	s.x[12] = uint64(uint32(insertBits(uint64(uint32(s.x[12])), 0x194c, 16, 16, 32)))

	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) + uint64(uint32(s.x[12])))))

	s.x[12] = uint64(uint32(uint64(0x156a)))

	s.x[12] = uint64(uint32(insertBits(uint64(uint32(s.x[12])), 0x1243, 16, 16, 32)))

	s.x[10] = uint64(uint32((uint64(uint32(s.x[10])) ^ uint64(uint32(s.x[12])))))

	s.x[10] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[10])))))

	if err := s.storeValue(addOffset(s.sp, 748), uint64(uint32(s.x[8])), 4); err != nil {
		return 0, err
	}

	if err := s.storeValue(addOffset(s.sp, 752), uint64(uint32(s.x[10])), 4); err != nil {
		return 0, err
	}

	value39, err := s.load(addOffset(s.sp, 144), 8)
	if err != nil {
		return 0, err
	}
	s.x[8] = uint64(value39)

	s.x[25] = uint64((s.x[8] - uint64(0x140)))

	if err := s.storeValue(addOffset(s.sp, 144), s.x[25], 8); err != nil {
		return 0, err
	}

	s.x[10] = uint64(uint32(uint64(0xe82b)))

	s.x[10] = uint64(uint32(insertBits(uint64(uint32(s.x[10])), 0x4696, 16, 16, 32)))

	s.x[19] = uint64(uint32(uint64(0x238f)))

	s.x[19] = uint64(uint32(insertBits(uint64(uint32(s.x[19])), 0xf530, 16, 16, 32)))

	value46a, err := s.load(addOffset(s.sp, 72), 8)
	if err != nil {
		return 0, err
	}
	value46b, err := s.load((addOffset(s.sp, 72))+8, 8)
	if err != nil {
		return 0, err
	}
	s.x[8] = uint64(value46a)
	s.x[26] = uint64(value46b)

	value47, err := s.load(addOffset(s.sp, 96), 8)
	if err != nil {
		return 0, err
	}
	s.x[11] = uint64(value47)

	s.x[17] = uint64(uint64(0x377c))

	next = s.x[11]
	return next, nil
}

func runAccountByteCompareFront(s *state) (uint64, error) {
	s.traceBlock(0x10030d3e0)

	s.x[11] = uint64(uint32(uint64(0xa61d)))

	s.x[11] = uint64(uint32(insertBits(uint64(uint32(s.x[11])), 0xf52f, 16, 16, 32)))

	s.x[15] = uint64(uint32((uint64(uint32(s.x[2])) + uint64(uint32(s.x[11])))))

	s.x[11] = uint64(uint32(uint64(0x69)))

	s.x[17] = uint64(uint32((uint64(uint32(s.x[15])) ^ uint64(uint32(s.x[11])))))

	s.x[16] = uint64(uint32((uint64(uint32(s.x[16])) & uint64(0xff))))

	s.compare32(uint32(uint64(uint32(s.x[16]))), uint32(uint64(uint8(uint64(uint32(s.x[17]))))))

	s.x[16] = uint64(uint32(boolValue(s.condition("eq"))))
	return 0x10030d400, nil
}

func runAccountByteCompareDispatchTail(s *state) (uint64, error) {
	s.traceBlock(0x10030d400)
	var next uint64

	s.x[11] = uint64(uint32(uint64(0xe562)))

	s.x[16] = uint64(uint32((uint64(uint32(s.x[16])) * uint64(uint32(s.x[11])))))

	s.x[16] = uint64(uint32((uint64(uint32(s.x[16])) ^ uint64(uint32(s.x[2])))))

	value3, err := s.load((s.x[24] + shiftLeft(signExtend32(uint32(uint64(uint32(s.x[16])))), 3, 64)), 8)
	if err != nil {
		return 0, err
	}
	s.x[0] = uint64(value3)

	s.x[17] = uint64(uint64(0x9969))

	next = s.x[0]
	return next, nil
}

func runAccountLoopLimitRoute(s *state) (uint64, error) {
	s.traceBlock(0x10097a170)
	var next uint64

	s.x[8] = uint64(uint32(uint64(0x62f9)))

	s.x[9] = uint64(uint32(uint64(0x338b)))

	s.x[23] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[9])))))

	value3, err := s.load(addOffset(s.sp, 2148), 4)
	if err != nil {
		return 0, err
	}
	s.x[9] = uint64(uint32(value3))

	s.compare32(uint32(uint64(uint32(s.x[10]))), uint32(uint64(uint32(s.x[9]))))

	s.x[9] = uint64(uint32(boolValue(s.condition("lo"))))

	s.x[8] = uint64(uint32((uint64(uint32(s.x[11])) + uint64(uint32(s.x[8])))))

	s.x[9] = uint64(uint32((uint64(uint32(s.x[8])) ^ uint64(uint32(s.x[9])))))

	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) & uint64(0x1))))

	s.x[12] = uint64(uint32((uint64(uint32(s.x[23])) - uint64(0x925))))

	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) * uint64(uint32(s.x[12])))))

	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[11])))))

	value12, err := s.load((s.x[22] + shiftLeft(signExtend32(uint32(uint64(uint32(s.x[9])))), 3, 64)), 8)
	if err != nil {
		return 0, err
	}
	s.x[9] = uint64(value12)

	s.x[19] = uint64(uint32(uint64(0x7bc5)))

	s.x[19] = uint64(uint32(insertBits(uint64(uint32(s.x[19])), 0xce75, 16, 16, 32)))

	s.x[17] = uint64(uint64(0xe25e))

	next = s.x[9]
	return next, nil
}
