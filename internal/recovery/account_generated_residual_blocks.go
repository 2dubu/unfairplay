package recovery

const (
	accountExternalCallSetup     = uint64(0x100067950)
	accountExternalCallResume    = uint64(0x100067a74)
	accountRoundStateSetup       = uint64(0x10006ba18)
	accountConditionalTableSetup = uint64(0x1001d483c)
	accountIndexedByteStoreRoute = uint64(0x100240b58)
	accountSignedLimitRouteA     = uint64(0x1002ac5dc)
	accountSignedLimitRouteB     = uint64(0x1002c5070)
	accountSignedLimitRouteC     = uint64(0x1003be340)
)

func runAccountExternalCallSetup(s *state) (uint64, error) {
	s.traceBlock(0x100067950)
	var next uint64

	s.x[8] = uint64(uint32(uint64(0xfac)))

	s.x[8] = uint64(uint32(insertBits(uint64(uint32(s.x[8])), 0x1448, 16, 16, 32)))

	s.x[23] = uint64(uint32((uint64(uint32(s.x[21])) + uint64(uint32(s.x[8])))))

	s.x[8] = uint64(uint32(uint64(0xdda)))

	s.x[8] = uint64(uint32(insertBits(uint64(uint32(s.x[8])), 0x1448, 16, 16, 32)))

	s.x[22] = uint64(uint32((uint64(uint32(s.x[23])) ^ uint64(uint32(s.x[8])))))

	s.x[8] = uint64(uint32(uint64(0x38af)))

	s.x[8] = uint64(uint32(insertBits(uint64(uint32(s.x[8])), 0x233f, 16, 16, 32)))

	s.x[8] = uint64(uint32((uint64(uint32(s.x[24])) ^ uint64(uint32(s.x[8])))))

	s.x[9] = uint64(uint32(uint64(0xc751)))

	s.x[9] = uint64(uint32(insertBits(uint64(uint32(s.x[9])), 0xdcc0, 16, 16, 32)))

	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) + uint64(uint32(s.x[9])))))

	s.x[9] = uint64(uint32(uint64(0xd510)))

	s.x[9] = uint64(uint32(insertBits(uint64(uint32(s.x[9])), 0x9332, 16, 16, 32)))

	s.x[9] = uint64(uint32((uint64(uint32(s.x[24])) ^ uint64(uint32(s.x[9])))))

	s.x[10] = uint64(uint32(uint64(0x2af0)))

	s.x[10] = uint64(uint32(insertBits(uint64(uint32(s.x[10])), 0x6ccd, 16, 16, 32)))

	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) + uint64(uint32(s.x[10])))))

	s.x[8] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[8])))))

	s.x[9] = uint64(uint32(uint64(0xb2a9)))

	s.x[9] = uint64(uint32(insertBits(uint64(uint32(s.x[9])), 0x7e78, 16, 16, 32)))

	s.x[9] = uint64(uint32((uint64(uint32(s.x[24])) ^ uint64(uint32(s.x[9])))))

	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[22])))))

	s.x[10] = uint64(uint32(uint64(0x6986)))

	s.x[10] = uint64(uint32(insertBits(uint64(uint32(s.x[10])), 0x8187, 16, 16, 32)))

	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) + uint64(uint32(s.x[10])))))

	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) ^ uint64(uint32(s.x[9])))))

	s.x[9] = uint64((s.x[29] - uint64(0xa8)))

	s.x[10] = uint64(uint32(uint64(0xd262)))

	s.x[10] = uint64(uint32(insertBits(uint64(uint32(s.x[10])), 0xe732, 16, 16, 32)))

	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(uint32(s.x[10])))))

	s.x[10] = uint64(uint32(uint64(0x9ef7)))

	s.x[10] = uint64(uint32(insertBits(uint64(uint32(s.x[10])), 0xba20, 16, 16, 32)))

	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) * uint64(uint32(s.x[10])))))

	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) - uint64(uint32(s.x[9])))))

	s.x[10] = uint64(uint32(uint64(0x9fc)))

	s.x[10] = uint64(uint32(insertBits(uint64(uint32(s.x[10])), 0x2afa, 16, 16, 32)))

	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) + uint64(uint32(s.x[10])))))

	s.x[11] = uint64((s.sp + uint64(0x64)))

	if err := s.storeValue(addOffset(s.sp, 672), s.x[11], 8); err != nil {
		return 0, err
	}

	value40, err := s.load(addOffset(s.sp, 88), 8)
	if err != nil {
		return 0, err
	}
	s.x[11] = uint64(value40)

	if err := s.storeValue(addOffset(s.sp, 624), s.x[11], 8); err != nil {
		return 0, err
	}

	s.x[11] = uint64((s.sp + uint64(0x6c)))

	if err := s.storeValue(addOffset(s.sp, 664), s.x[11], 8); err != nil {
		return 0, err
	}

	s.x[11] = uint64((s.sp + uint64(0x70)))

	if err := s.storeValue(addOffset(s.sp, 648), s.x[11], 8); err != nil {
		return 0, err
	}

	s.x[9] = uint64(uint32((uint64(uint32(s.x[21])) - uint64(uint32(s.x[9])))))

	s.x[11] = uint64(uint32(uint64(0x1341)))

	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) + uint64(uint32(s.x[11])))))

	if err := s.storeValue(addOffset(s.x[29], -128), uint64(uint32(s.x[8])), 4); err != nil {
		return 0, err
	}

	if err := s.storeValue(addOffset(s.x[29], -168), uint64(uint32(s.x[9])), 4); err != nil {
		return 0, err
	}

	s.x[8] = uint64((s.sp + uint64(0x68)))

	if err := s.storeValue(addOffset(s.sp, 640), s.x[8], 8); err != nil {
		return 0, err
	}

	value53, err := s.load(s.x[6], 4)
	if err != nil {
		return 0, err
	}
	s.x[8] = uint64(uint32(value53))

	value54, err := s.load(addOffset(s.sp, 80), 8)
	if err != nil {
		return 0, err
	}
	s.x[9] = uint64(value54)

	value55, err := s.load(s.x[9], 4)
	if err != nil {
		return 0, err
	}
	s.x[9] = uint64(uint32(value55))

	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) ^ uint64(uint32(s.x[9])))))

	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) ^ uint64(uint32(s.x[3])))))

	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) * uint64(uint32(s.x[4])))))

	s.x[9] = uint64((s.x[7] + uint64(uint8(uint64(uint32(s.x[8]))))))

	value60, err := s.load(s.x[9], 1)
	if err != nil {
		return 0, err
	}
	s.x[9] = uint64(uint32(value60))

	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) ^ uint64(0x4))))

	value62, err := s.load((s.x[30] + uint64(uint32(uint64(uint32(s.x[9]))))), 1)
	if err != nil {
		return 0, err
	}
	s.x[9] = uint64(uint32(value62))

	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) ^ uint64(uint32(s.x[9])))))

	s.x[9] = uint64(uint32(uint64(0xd03b)))

	s.x[8] = uint64(uint32((uint64(uint32(s.x[9])) + uint64(uint8(uint64(uint32(s.x[8])))))))

	value66, err := s.load((s.x[20] + shiftLeft(uint64(uint32(uint64(uint32(s.x[8])))), 3, 64)), 8)
	if err != nil {
		return 0, err
	}
	s.x[8] = uint64(value66)

	s.x[0] = uint64((s.x[29] - uint64(0xa8)))

	s.x[19] = uint64(s.x[2])

	if err := s.storeValue(addOffset(s.sp, 40), s.x[30], 8); err != nil {
		return 0, err
	}
	if err := s.storeValue((addOffset(s.sp, 40))+8, s.x[6], 8); err != nil {
		return 0, err
	}

	if err := s.storeValue(addOffset(s.sp, 32), s.x[7], 8); err != nil {
		return 0, err
	}

	s.x[17] = uint64(uint64(0x48c7))

	s.x[30] = 0x100067a74
	next = s.x[8]
	return next, nil
}

func runAccountExternalCallResume(s *state) (uint64, error) {
	s.traceBlock(0x100067a74)
	var next uint64

	s.x[1] = uint64(uint32(uint64(0x7bc5)))

	s.x[1] = uint64(uint32(insertBits(uint64(uint32(s.x[1])), 0xce75, 16, 16, 32)))

	value2a, err := s.load(addOffset(s.sp, 32), 8)
	if err != nil {
		return 0, err
	}
	value2b, err := s.load((addOffset(s.sp, 32))+8, 8)
	if err != nil {
		return 0, err
	}
	s.x[7] = uint64(value2a)
	s.x[30] = uint64(value2b)

	value3, err := s.load(addOffset(s.sp, 48), 8)
	if err != nil {
		return 0, err
	}
	s.x[6] = uint64(value3)

	s.x[5] = uint64(0x10105b000)

	s.x[5] = uint64((s.x[5] + uint64(0x6f0)))

	s.x[4] = uint64(uint32(uint64(0xd4cf)))

	s.x[4] = uint64(uint32(insertBits(uint64(uint32(s.x[4])), 0xa5c, 16, 16, 32)))

	s.x[3] = uint64(uint32(uint64(0xa6)))

	value9, err := s.load(addOffset(s.x[29], -152), 4)
	if err != nil {
		return 0, err
	}
	s.x[12] = uint64(uint32(value9))

	s.compare32(uint32(uint64(uint32(s.x[12]))), uint32(uint64(uint32(s.x[1]))))

	s.x[8] = uint64(uint32(boolValue(s.condition("eq"))))

	s.x[9] = uint64(uint32(uint64(0x4862)))

	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) * uint64(uint32(s.x[9])))))

	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) ^ uint64(uint32(s.x[21])))))

	value15, err := s.load((s.x[20] + shiftLeft(signExtend32(uint32(uint64(uint32(s.x[8])))), 3, 64)), 8)
	if err != nil {
		return 0, err
	}
	s.x[8] = uint64(value15)

	s.x[17] = uint64(uint64(0xdb39))

	next = s.x[8]
	return next, nil
}

func runAccountRoundStateSetup(s *state) (uint64, error) {
	s.traceBlock(0x10006ba18)
	var next uint64

	s.x[8] = uint64(uint32(uint64(0xe707)))

	s.x[8] = uint64(uint32(insertBits(uint64(uint32(s.x[8])), 0xebb7, 16, 16, 32)))

	s.x[8] = uint64(uint32((uint64(uint32(s.x[9])) + uint64(uint32(s.x[8])))))

	s.x[11] = uint64(uint32(uint64(0x94d)))

	s.x[21] = uint64(uint32((uint64(uint32(s.x[8])) | uint64(uint32(s.x[11])))))

	s.x[8] = uint64(uint32(uint64(0x5afb)))

	s.x[8] = uint64(uint32(insertBits(uint64(uint32(s.x[8])), 0xebb8, 16, 16, 32)))

	s.x[8] = uint64(uint32((uint64(uint32(s.x[9])) + uint64(uint32(s.x[8])))))

	s.x[11] = uint64(uint32(uint64(0x2248)))

	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) | uint64(uint32(s.x[11])))))

	s.x[11] = uint64(uint32(uint64(0x11cc)))

	s.x[11] = uint64(uint32(insertBits(uint64(uint32(s.x[11])), 0xd834, 16, 16, 32)))

	s.x[11] = uint64(uint32((uint64(uint32(s.x[24])) ^ uint64(uint32(s.x[11])))))

	s.x[11] = uint64(uint32((uint64(uint32(s.x[21])) ^ uint64(uint32(s.x[11])))))

	s.x[12] = uint64(uint32(uint64(0x7037)))

	s.x[12] = uint64(uint32(insertBits(uint64(uint32(s.x[12])), 0x7323, 16, 16, 32)))

	s.x[12] = uint64(uint32((uint64(uint32(s.x[24])) ^ uint64(uint32(s.x[12])))))

	s.x[13] = uint64(uint32(uint64(0x8fc9)))

	s.x[13] = uint64(uint32(insertBits(uint64(uint32(s.x[13])), 0x8cdc, 16, 16, 32)))

	s.x[12] = uint64(uint32((uint64(uint32(s.x[12])) + uint64(uint32(s.x[13])))))

	s.x[13] = uint64(uint32(uint64(0x363)))

	s.x[13] = uint64(uint32(insertBits(uint64(uint32(s.x[13])), 0x6562, 16, 16, 32)))

	s.x[13] = uint64(uint32((uint64(uint32(s.x[24])) ^ uint64(uint32(s.x[13])))))

	s.x[14] = uint64(uint32(uint64(0xfc9d)))

	s.x[14] = uint64(uint32(insertBits(uint64(uint32(s.x[14])), 0x9a9d, 16, 16, 32)))

	s.x[13] = uint64(uint32((uint64(uint32(s.x[13])) + uint64(uint32(s.x[14])))))

	s.x[12] = uint64(uint32((uint64(uint32(s.x[12])) ^ uint64(uint32(s.x[13])))))

	s.x[13] = uint64(uint32(uint64(0xf76f)))

	s.x[13] = uint64(uint32(insertBits(uint64(uint32(s.x[13])), 0x27cb, 16, 16, 32)))

	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) + uint64(uint32(s.x[13])))))

	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[12])))))

	s.x[12] = uint64(uint32(uint64(0xfc75)))

	s.x[12] = uint64(uint32(insertBits(uint64(uint32(s.x[12])), 0xac13, 16, 16, 32)))

	s.x[12] = uint64(uint32((uint64(uint32(s.x[27])) + uint64(uint32(s.x[12])))))

	s.x[12] = uint64(uint32((uint64(uint32(s.x[12])) + uint64(uint32(s.x[8])))))

	s.x[15] = uint64(uint32(uint64(0x1ea9)))

	s.x[15] = uint64(uint32(insertBits(uint64(uint32(s.x[15])), 0xd2c9, 16, 16, 32)))

	s.compare32(uint32(uint64(uint32(s.x[12]))), uint32(uint64(uint32(s.x[15]))))

	s.x[13] = uint64(uint32(boolValue(s.condition("lo"))))

	s.x[14] = uint64(uint32(uint64(0x9a6e)))

	s.x[14] = uint64(uint32(insertBits(uint64(uint32(s.x[14])), 0xa13e, 16, 16, 32)))

	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) + uint64(uint32(s.x[14])))))

	s.compare32(uint32(uint64(uint32(s.x[11]))), uint32(uint64(uint32(s.x[15]))))

	s.x[14] = uint64(uint32(boolValue(s.condition("lo"))))

	s.x[13] = uint64(uint32((uint64(uint32(s.x[13])) ^ uint64(uint32(s.x[14])))))

	s.compare32(uint32(uint64(uint32(s.x[12]))), uint32(uint64(uint32(s.x[11]))))

	s.x[11] = uint64(uint32(boolValue(s.condition("lo"))))

	s.compare32(uint32(uint64(uint32(s.x[13]))), uint32(uint64(0x0)))

	s.x[11] = uint64(uint32(selectValue(s.condition("ne"), uint64(uint32(s.x[14])), uint64(uint32(s.x[11])))))

	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(0x1))))

	s.x[12] = uint64(uint32(uint64(0x7d)))

	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) * uint64(uint32(s.x[12])))))

	s.x[12] = uint64(uint32(uint64(0x90a7)))

	s.x[12] = uint64(uint32(insertBits(uint64(uint32(s.x[12])), 0xebb8, 16, 16, 32)))

	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) + uint64(uint32(s.x[12])))))

	s.x[9] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[9])))))

	value56, err := s.load((s.x[20] + shiftLeft(signExtend32(uint32(uint64(uint32(s.x[9])))), 3, 64)), 8)
	if err != nil {
		return 0, err
	}
	s.x[9] = uint64(value56)

	s.x[17] = uint64(uint64(0x33b0))

	next = s.x[9]
	return next, nil
}

func runAccountConditionalTableSetup(s *state) (uint64, error) {
	s.traceBlock(0x1001d483c)
	var next uint64

	s.x[8] = uint64(uint32((uint64(uint32(s.x[8])) + uint64(uint32(s.x[1])))))

	s.x[9] = uint64(uint32(uint64(0xed1f)))

	s.x[9] = uint64(uint32(insertBits(uint64(uint32(s.x[9])), 0xf772, 16, 16, 32)))

	s.x[9] = uint64(uint32((uint64(uint32(s.x[8])) + uint64(uint32(s.x[9])))))

	s.x[11] = uint64(uint32(uint64(0xe7bf)))

	s.x[11] = uint64(uint32(insertBits(uint64(uint32(s.x[11])), 0xfdbc, 16, 16, 32)))

	s.x[9] = uint64(uint32((uint64(uint32(s.x[9])) & uint64(uint32(s.x[11])))))

	s.compare32(uint32(uint64(uint32(s.x[20]))), uint32(uint64(uint32(s.x[16]))))

	s.x[11] = uint64(uint32(boolValue(s.condition("eq"))))

	s.x[12] = uint64(uint32(uint64(0xdcd4)))

	s.x[12] = uint64(uint32(insertBits(uint64(uint32(s.x[12])), 0xacf, 16, 16, 32)))

	s.x[12] = uint64(uint32((uint64(uint32(s.x[9])) + uint64(uint32(s.x[12])))))

	s.x[11] = uint64(uint32((uint64(uint32(s.x[12])) * uint64(uint32(s.x[11])))))

	s.x[11] = uint64(uint32((uint64(uint32(s.x[11])) ^ uint64(uint32(s.x[8])))))

	value14, err := s.load((s.x[24] + shiftLeft(signExtend32(uint32(uint64(uint32(s.x[11])))), 3, 64)), 8)
	if err != nil {
		return 0, err
	}
	s.x[11] = uint64(value14)

	value15a, err := s.load(addOffset(s.sp, 128), 8)
	if err != nil {
		return 0, err
	}
	value15b, err := s.load((addOffset(s.sp, 128))+8, 8)
	if err != nil {
		return 0, err
	}
	s.x[13] = uint64(value15a)
	s.x[12] = uint64(value15b)

	value16, err := s.load(s.x[12], 4)
	if err != nil {
		return 0, err
	}
	s.x[12] = uint64(uint32(value16))

	value17, err := s.load(s.x[13], 4)
	if err != nil {
		return 0, err
	}
	s.x[13] = uint64(uint32(value17))

	s.x[12] = uint64(uint32((uint64(uint32(s.x[13])) + uint64(uint32(s.x[12])))))

	s.x[13] = uint64(uint32(uint64(0xa6)))

	s.x[12] = uint64(uint32((uint64(uint32(s.x[12])) ^ uint64(uint32(s.x[13])))))

	s.x[13] = uint64(uint32(uint64(0xd4cf)))

	s.x[13] = uint64(uint32(insertBits(uint64(uint32(s.x[13])), 0xa5c, 16, 16, 32)))

	s.x[12] = uint64(uint32((uint64(uint32(s.x[12])) * uint64(uint32(s.x[13])))))

	value24, err := s.load(addOffset(s.sp, 32), 8)
	if err != nil {
		return 0, err
	}
	s.x[13] = uint64(value24)

	s.x[13] = uint64((s.x[13] + uint64(uint8(uint64(uint32(s.x[12]))))))

	value26, err := s.load(s.x[13], 1)
	if err != nil {
		return 0, err
	}
	s.x[13] = uint64(uint32(value26))

	s.x[14] = uint64(uint32(uint64(0xb1)))

	s.x[13] = uint64(uint32((uint64(uint32(s.x[13])) ^ uint64(uint32(s.x[14])))))

	value29, err := s.load(addOffset(s.sp, 24), 8)
	if err != nil {
		return 0, err
	}
	s.x[14] = uint64(value29)

	value30, err := s.load((s.x[14] + uint64(uint32(uint64(uint32(s.x[13]))))), 1)
	if err != nil {
		return 0, err
	}
	s.x[13] = uint64(uint32(value30))

	s.x[12] = uint64(uint32((uint64(uint32(s.x[12])) ^ uint64(uint32(s.x[13])))))

	s.x[13] = uint64(uint32(uint64(0x140d)))

	s.x[12] = uint64(uint32((uint64(uint32(s.x[13])) + uint64(uint8(uint64(uint32(s.x[12])))))))

	s.x[13] = uint64(0x10105b000)

	s.x[13] = uint64((s.x[13] + uint64(0x6f0)))

	value36, err := s.load((s.x[13] + shiftLeft(uint64(uint32(uint64(uint32(s.x[12])))), 3, 64)), 8)
	if err != nil {
		return 0, err
	}
	s.x[12] = uint64(value36)

	s.x[13] = uint64(uint64(0xffffffffffff3cb6))

	s.x[13] = uint64(insertBits(s.x[13], 0xe43c, 16, 16, 64))

	s.x[23] = uint64((s.x[12] + s.x[13]))

	s.x[17] = uint64(uint64(0xd396))

	next = s.x[11]
	return next, nil
}

func runAccountIndexedByteStoreRoute(s *state) (uint64, error) {
	s.traceBlock(0x100240b58)
	var next uint64

	s.x[8] = uint64(uint32(uint64(0x666c)))

	s.x[8] = uint64(uint32(insertBits(uint64(uint32(s.x[8])), 0x80ff, 16, 16, 32)))

	s.x[8] = uint64(uint32((uint64(uint32(s.x[13])) ^ uint64(uint32(s.x[8])))))

	s.x[14] = uint64(uint32(uint64(0x16a9)))

	s.x[14] = uint64(uint32((uint64(uint32(s.x[8])) * uint64(uint32(s.x[14])))))

	s.x[8] = uint64(uint32(uint64(0xf816)))

	s.x[8] = uint64(uint32(insertBits(uint64(uint32(s.x[8])), 0x7f00, 16, 16, 32)))

	s.x[8] = uint64(uint32((uint64(uint32(s.x[13])) + uint64(uint32(s.x[8])))))

	s.x[13] = uint64((s.x[0] + uint64(0x4)))

	if err := s.storeValue((s.x[13] + uint64(uint32(uint64(uint32(s.x[11]))))), uint64(uint32(s.x[16])), 1); err != nil {
		return 0, err
	}

	s.x[11] = uint64(uint32((uint64(uint32(s.x[1])) + uint64(uint32(s.x[9])))))

	value11, err := s.load((s.x[10] + shiftLeft(signExtend32(uint32(uint64(uint32(s.x[8])))), 3, 64)), 8)
	if err != nil {
		return 0, err
	}
	s.x[8] = uint64(value11)

	s.x[17] = uint64(uint64(0x9478))

	next = s.x[8]
	return next, nil
}

func runAccountSignedLimitRouteA(s *state) (uint64, error) {
	s.traceBlock(0x1002ac5dc)
	var next uint64

	s.x[13] = uint64(uint32(uint64(0xf31)))

	s.x[14] = uint64(uint32((uint64(uint32(s.x[15])) ^ uint64(uint32(s.x[13])))))

	s.x[13] = uint64(uint32(uint64(0x2471)))

	s.compare32(uint32(uint64(uint32(s.x[12]))), uint32(uint64(0xa)))

	s.x[16] = uint64(uint32(boolValue(s.condition("lt"))))

	s.x[13] = uint64(uint32((uint64(uint32(s.x[15])) + uint64(uint32(s.x[13])))))

	s.x[17] = uint64(uint32((uint64(uint32(s.x[14])) + uint64(uint32(s.x[13])))))

	s.x[1] = uint64(uint32(uint64(0xffff8cc5)))

	s.x[17] = uint64(uint32((uint64(uint32(s.x[17])) + uint64(uint32(s.x[1])))))

	s.x[16] = uint64(uint32((uint64(uint32(s.x[17])) * uint64(uint32(s.x[16])))))

	s.x[15] = uint64(uint32((uint64(uint32(s.x[16])) ^ uint64(uint32(s.x[15])))))

	value11, err := s.load((s.x[10] + shiftLeft(signExtend32(uint32(uint64(uint32(s.x[15])))), 3, 64)), 8)
	if err != nil {
		return 0, err
	}
	s.x[15] = uint64(value11)

	s.x[17] = uint64(uint64(0x3137))

	next = s.x[15]
	return next, nil
}

func runAccountSignedLimitRouteB(s *state) (uint64, error) {
	s.traceBlock(0x1002c5070)
	var next uint64

	s.x[15] = uint64(uint32(uint64(0x220f)))

	s.x[14] = uint64(uint32(uint64(0xa6b)))

	s.x[14] = uint64(uint32((uint64(uint32(s.x[13])) ^ uint64(uint32(s.x[14])))))

	s.x[15] = uint64(uint32((uint64(uint32(s.x[13])) + uint64(uint32(s.x[15])))))

	s.x[16] = uint64(uint32((uint64(uint32(s.x[15])) + uint64(uint32(s.x[14])))))

	s.x[16] = uint64(uint32((uint64(uint32(s.x[16])) + uint64(0x6c))))

	s.compare32(uint32(uint64(uint32(s.x[12]))), uint32(uint64(int64(int8(uint64(uint32(s.x[16])))))))

	s.x[16] = uint64(uint32(boolValue(s.condition("lt"))))

	s.x[17] = uint64(uint32(uint64(0xf4)))

	s.x[16] = uint64(uint32((uint64(uint32(s.x[16])) * uint64(uint32(s.x[17])))))

	s.x[13] = uint64(uint32((uint64(uint32(s.x[16])) ^ uint64(uint32(s.x[13])))))

	value11, err := s.load((s.x[10] + shiftLeft(signExtend32(uint32(uint64(uint32(s.x[13])))), 3, 64)), 8)
	if err != nil {
		return 0, err
	}
	s.x[13] = uint64(value11)

	s.x[17] = uint64(uint64(0xa4a8))

	next = s.x[13]
	return next, nil
}

func runAccountSignedLimitRouteC(s *state) (uint64, error) {
	s.traceBlock(0x1003be340)
	var next uint64

	s.x[14] = uint64(uint32(uint64(0xffffc16c)))

	s.x[15] = uint64(uint32(uint64(0x2512)))

	s.x[14] = uint64(uint32((uint64(uint32(s.x[13])) + uint64(uint32(s.x[14])))))

	s.x[14] = uint64(uint32((uint64(uint32(s.x[14])) | uint64(uint32(s.x[15])))))

	s.x[15] = uint64(uint32(uint64(0x59)))

	s.x[15] = uint64(uint32((uint64(uint32(s.x[14])) ^ uint64(uint32(s.x[15])))))

	s.compare32(uint32(uint64(uint32(s.x[12]))), uint32(uint64(int64(int8(uint64(uint32(s.x[15])))))))

	s.x[12] = uint64(uint32(boolValue(s.condition("lt"))))

	s.x[15] = uint64(uint32(uint64(0x34)))

	s.x[12] = uint64(uint32((uint64(uint32(s.x[12])) * uint64(uint32(s.x[15])))))

	s.x[12] = uint64(uint32((uint64(uint32(s.x[12])) ^ uint64(uint32(s.x[13])))))

	value11, err := s.load((s.x[10] + shiftLeft(signExtend32(uint32(uint64(uint32(s.x[12])))), 3, 64)), 8)
	if err != nil {
		return 0, err
	}
	s.x[12] = uint64(value11)

	s.x[17] = uint64(uint64(0xf921))

	next = s.x[12]
	return next, nil
}
