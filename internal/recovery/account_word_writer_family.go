package recovery

import "fmt"

const (
	accountFrameWordWriterStart = uint64(0x10009d35c)
	accountFrameWordWriterTail  = uint64(0x10009d400)
	accountStackWordWriterStart = uint64(0x1000fd328)
	accountStackWordWriterTail  = uint64(0x1000fd400)
	accountWideWordWriter       = uint64(0x1001e4b10)
)

func runAccountFrameWordWriter(s *state) (uint64, error) {
	s.traceBlock(accountFrameWordWriterStart)
	s.x[11] = s.x[29] - 0xdc
	s.x[11] += uint64(uint32(s.x[10]))
	value, err := s.load(s.x[11]+0x48, 4)
	if err != nil {
		return 0, fmt.Errorf("account frame-word input load: %w", err)
	}
	s.x[11] = uint64(uint32(value))
	s.x[12] = uint64(uint32(s.x[10]) ^ uint32(0xfdcf77f0))
	s.x[13] = uint64(uint32(0xfb9eefe0) & uint32(shiftLeft(s.x[10], 1, 32)))
	s.x[14] = uint64(uint32(shiftRight(s.x[11], 24, 32)))
	s.x[15] = uint64(uint32(s.x[14]) ^ 0x0b)
	s.x[16] = uint64(uint32(s.x[8]) ^ 0x4e)
	s.x[14] = uint64(uint32(s.x[14]) & uint32(s.x[16]))
	s.x[16] = uint64(uint32(s.x[15]) & 0x1f)
	s.x[14] = uint64(uint32(s.x[16]) ^ uint32(s.x[14]))
	s.x[14] = uint64(uint32(s.x[15]) - uint32(shiftLeft(s.x[14], 1, 32)))
	s.x[14] = uint64(uint32(s.x[14]) + 0x1b)
	s.x[15] = 0x02308810
	s.x[12] = uint64(uint32(s.x[13]) + uint32(s.x[12]))
	s.x[13] = uint64(uint32(s.x[12]) + uint32(s.x[15]))
	if err := s.storeValue(s.x[0]+uint64(uint32(s.x[13])), s.x[14], 1); err != nil {
		return 0, fmt.Errorf("account frame-word byte 3 store: %w", err)
	}

	s.x[13] = uint64(uint32(shiftRight(s.x[11], 16, 32)) ^ 0x1f)
	s.x[14] = 0x36
	s.x[15] = uint64(uint32(s.x[14]) & uint32(shiftLeft(s.x[13], 1, 32)))
	s.x[13] = uint64(uint32(s.x[13]) - uint32(s.x[15]))
	s.x[13] = uint64(uint32(s.x[13]) + 0x1b)
	s.x[15] = 0x02308811
	s.x[15] = uint64(uint32(s.x[12]) + uint32(s.x[15]))
	if err := s.storeValue(s.x[0]+uint64(uint32(s.x[15])), s.x[13], 1); err != nil {
		return 0, fmt.Errorf("account frame-word byte 2 store: %w", err)
	}

	s.x[13] = uint64(uint32(0xffffffcd) ^ uint32(shiftRight(s.x[11], 8, 32)))
	s.x[14] = uint64(uint32(s.x[14]) & uint32(shiftLeft(s.x[13], 1, 32)))
	s.x[13] = uint64(uint32(s.x[13]) - uint32(s.x[14]))
	s.x[13] = uint64(uint32(s.x[13]) + 0x1b)
	s.x[14] = 0x02308812
	return runAccountFrameWordWriterTail(s)
}

func runAccountFrameWordWriterTail(s *state) (uint64, error) {
	s.traceBlock(accountFrameWordWriterTail)
	s.x[14] = uint64(uint32(s.x[12]) + uint32(s.x[14]))
	if err := s.storeValue(s.x[0]+uint64(uint32(s.x[14])), s.x[13], 1); err != nil {
		return 0, fmt.Errorf("account frame-word byte 1 store: %w", err)
	}
	s.x[13] = 0xa6
	s.x[11] = uint64(uint32(s.x[11]) ^ uint32(s.x[13]))
	s.x[13] = 0x02308813
	s.x[12] = uint64(uint32(s.x[12]) + uint32(s.x[13]))
	if err := s.storeValue(s.x[0]+uint64(uint32(s.x[12])), s.x[11], 1); err != nil {
		return 0, fmt.Errorf("account frame-word byte 0 store: %w", err)
	}
	s.x[10] = uint64(uint32(s.x[10]) + 4)
	value, err := s.load(addOffset(s.x[29], -0x80), 4)
	if err != nil {
		return 0, fmt.Errorf("account frame-word limit load: %w", err)
	}
	s.x[11] = uint64(uint32(value))
	s.compare32(uint32(s.x[10]), uint32(s.x[11]))
	s.x[11] = uint64(uint32(boolValue(s.condition("lo"))))
	s.x[12] = 0x7a8
	s.x[11] = uint64(uint32(s.x[11]) * uint32(s.x[12]))
	s.x[11] = uint64(uint32(s.x[11]) ^ uint32(s.x[9]))
	value, err = s.load(s.x[21]+shiftLeft(signExtend32(uint32(s.x[11])), 3, 64), 8)
	if err != nil {
		return 0, fmt.Errorf("account frame-word target load: %w", err)
	}
	s.x[11] = value
	s.x[17] = 0x98cd
	return s.x[11], nil
}

func runAccountStackWordWriter(s *state) (uint64, error) {
	s.traceBlock(accountStackWordWriterStart)
	s.x[8] = uint64(uint32(s.x[9]) + uint32(0x54822a1f))
	s.x[11] = 0x0005fefa
	s.x[8] = uint64(uint32(s.x[8]) & uint32(s.x[11]))
	s.x[11] = s.sp + 0x9c
	s.x[11] += uint64(uint32(s.x[10]))
	value, err := s.load(s.x[11]+0x48, 4)
	if err != nil {
		return 0, fmt.Errorf("account stack-word input load: %w", err)
	}
	s.x[11] = uint64(uint32(value))
	s.x[12] = uint64(uint32(s.x[10]) ^ uint32(0xbeaf7ffd))
	s.x[13] = uint64(uint32(shiftRight(s.x[11], 24, 32)))
	s.x[14] = uint64(uint32(s.x[13]) ^ 0x0b)
	s.x[15] = uint64(uint32(s.x[14]) & 0x1f)
	s.x[13] = uint64(uint32(s.x[13]) & 4)
	s.x[13] = uint64(uint32(s.x[15]) ^ uint32(s.x[13]))
	s.x[15] = uint64(uint32(s.x[8]) + uint32(0x7d5e80fe))
	s.x[15] = uint64(uint32(s.x[15]) & uint32(shiftLeft(s.x[10], 1, 32)))
	s.x[13] = uint64(uint32(s.x[14]) - uint32(shiftLeft(s.x[13], 1, 32)))
	s.x[13] = uint64(uint32(s.x[13]) + 0x1b)
	s.x[14] = 0x41508003
	s.x[12] = uint64(uint32(s.x[15]) + uint32(s.x[12]))
	s.x[14] = uint64(uint32(s.x[12]) + uint32(s.x[14]))
	if err := s.storeValue(s.x[16]+uint64(uint32(s.x[14])), s.x[13], 1); err != nil {
		return 0, fmt.Errorf("account stack-word byte 3 store: %w", err)
	}

	s.x[13] = uint64(uint32(shiftRight(s.x[11], 16, 32)) ^ 0x1f)
	s.x[14] = 0x36
	s.x[15] = uint64(uint32(s.x[14]) & uint32(shiftLeft(s.x[13], 1, 32)))
	s.x[13] = uint64(uint32(s.x[13]) - uint32(s.x[15]))
	s.x[13] = uint64(uint32(s.x[13]) + 0x1b)
	s.x[15] = 0x41508004
	s.x[15] = uint64(uint32(s.x[12]) + uint32(s.x[15]))
	if err := s.storeValue(s.x[16]+uint64(uint32(s.x[15])), s.x[13], 1); err != nil {
		return 0, fmt.Errorf("account stack-word byte 2 store: %w", err)
	}

	s.x[13] = uint64(uint32(0xffffffcd) ^ uint32(shiftRight(s.x[11], 8, 32)))
	s.x[14] = uint64(uint32(s.x[14]) & uint32(shiftLeft(s.x[13], 1, 32)))
	s.x[13] = uint64(uint32(s.x[13]) - uint32(s.x[14]))
	s.x[13] = uint64(uint32(s.x[13]) + 0x1b)
	s.x[14] = 0x41508005
	s.x[14] = uint64(uint32(s.x[12]) + uint32(s.x[14]))
	if err := s.storeValue(s.x[16]+uint64(uint32(s.x[14])), s.x[13], 1); err != nil {
		return 0, fmt.Errorf("account stack-word byte 1 store: %w", err)
	}

	s.x[13] = 0xa6
	s.x[11] = uint64(uint32(s.x[11]) ^ uint32(s.x[13]))
	s.x[13] = 0x41508006
	s.x[12] = uint64(uint32(s.x[12]) + uint32(s.x[13]))
	if err := s.storeValue(s.x[16]+uint64(uint32(s.x[12])), s.x[11], 1); err != nil {
		return 0, fmt.Errorf("account stack-word byte 0 store: %w", err)
	}
	return runAccountStackWordWriterTail(s)
}

func runAccountStackWordWriterTail(s *state) (uint64, error) {
	s.traceBlock(accountStackWordWriterTail)
	s.x[10] = uint64(uint32(s.x[10]) + 4)
	value, err := s.load(s.sp+0xf8, 4)
	if err != nil {
		return 0, fmt.Errorf("account stack-word limit load: %w", err)
	}
	s.x[11] = uint64(uint32(value))
	s.compare32(uint32(s.x[10]), uint32(s.x[11]))
	s.x[11] = uint64(uint32(boolValue(s.condition("hs"))))
	s.x[12] = 0x15
	s.x[11] = uint64(uint32(s.x[11]) * uint32(s.x[12]))
	s.x[11] = uint64(uint32(s.x[11]) | uint32(s.x[9]))
	value, err = s.load(s.x[20]+uint64(uint32(s.x[11]))*8, 8)
	if err != nil {
		return 0, fmt.Errorf("account stack-word target load: %w", err)
	}
	s.x[11] = value
	s.x[17] = 0x9986
	return s.x[11], nil
}

func runAccountWideWordWriter(s *state) (uint64, error) {
	s.traceBlock(accountWideWordWriter)
	s.x[11] = s.sp + 0x29c
	s.x[11] += uint64(uint32(s.x[10]))
	value, err := s.load(s.x[11]+0x48, 4)
	if err != nil {
		return 0, fmt.Errorf("account wide-word input load: %w", err)
	}
	s.x[11] = uint64(uint32(value))
	s.x[12] = uint64(uint32(s.x[10]) ^ uint32(0xbc97f6f0))
	s.x[13] = uint64(uint32(shiftRight(s.x[11], 24, 32)))
	s.x[14] = uint64(uint32(s.x[13]) ^ 0x0b)
	s.x[15] = uint64(uint32(s.x[8]) ^ 0x2f)
	s.x[15] = uint64(uint32(s.x[14]) & uint32(s.x[15]))
	s.x[13] = uint64(uint32(s.x[13]) & 4)
	s.x[13] = uint64(uint32(s.x[15]) ^ uint32(s.x[13]))
	s.x[15] = uint64(uint32(0x792fede0) & uint32(shiftLeft(s.x[10], 1, 32)))
	s.x[13] = uint64(uint32(s.x[14]) - uint32(shiftLeft(s.x[13], 1, 32)))
	s.x[13] = uint64(uint32(s.x[13]) + 0x1b)
	s.x[14] = 0x43680910
	s.x[12] = uint64(uint32(s.x[15]) + uint32(s.x[12]))
	s.x[14] = uint64(uint32(s.x[12]) + uint32(s.x[14]))
	if err := s.storeValue(s.x[16]+uint64(uint32(s.x[14])), s.x[13], 1); err != nil {
		return 0, fmt.Errorf("account wide-word byte 3 store: %w", err)
	}

	s.x[13] = uint64(uint32(shiftRight(s.x[11], 16, 32)) ^ 0x1f)
	s.x[14] = 0x36
	s.x[15] = uint64(uint32(s.x[14]) & uint32(shiftLeft(s.x[13], 1, 32)))
	s.x[13] = uint64(uint32(s.x[13]) - uint32(s.x[15]))
	s.x[13] = uint64(uint32(s.x[13]) + 0x1b)
	s.x[15] = 0x43680911
	s.x[15] = uint64(uint32(s.x[12]) + uint32(s.x[15]))
	if err := s.storeValue(s.x[16]+uint64(uint32(s.x[15])), s.x[13], 1); err != nil {
		return 0, fmt.Errorf("account wide-word byte 2 store: %w", err)
	}

	s.x[13] = uint64(uint32(0xffffffcd) ^ uint32(shiftRight(s.x[11], 8, 32)))
	s.x[14] = uint64(uint32(s.x[14]) & uint32(shiftLeft(s.x[13], 1, 32)))
	s.x[13] = uint64(uint32(s.x[13]) - uint32(s.x[14]))
	s.x[13] = uint64(uint32(s.x[13]) + 0x1b)
	s.x[14] = 0x43680912
	s.x[14] = uint64(uint32(s.x[12]) + uint32(s.x[14]))
	if err := s.storeValue(s.x[16]+uint64(uint32(s.x[14])), s.x[13], 1); err != nil {
		return 0, fmt.Errorf("account wide-word byte 1 store: %w", err)
	}

	s.x[13] = 0xa6
	s.x[11] = uint64(uint32(s.x[11]) ^ uint32(s.x[13]))
	s.x[13] = 0x43680913
	s.x[12] = uint64(uint32(s.x[12]) + uint32(s.x[13]))
	if err := s.storeValue(s.x[16]+uint64(uint32(s.x[12])), s.x[11], 1); err != nil {
		return 0, fmt.Errorf("account wide-word byte 0 store: %w", err)
	}
	s.x[10] = uint64(uint32(s.x[10]) + 4)
	value, err = s.load(s.sp+0x2f8, 4)
	if err != nil {
		return 0, fmt.Errorf("account wide-word limit load: %w", err)
	}
	s.x[11] = uint64(uint32(value))
	s.compare32(uint32(s.x[10]), uint32(s.x[11]))
	s.x[11] = uint64(uint32(boolValue(s.condition("hs"))))
	s.x[12] = 0x1f6
	s.x[11] = uint64(uint32(s.x[11]) * uint32(s.x[12]))
	s.x[11] = uint64(uint32(s.x[11]) ^ uint32(s.x[9]))
	value, err = s.load(s.x[24]+shiftLeft(signExtend32(uint32(s.x[11])), 3, 64), 8)
	if err != nil {
		return 0, fmt.Errorf("account wide-word target load: %w", err)
	}
	s.x[11] = value
	s.x[17] = 0x5b8c
	return s.x[11], nil
}
