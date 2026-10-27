package recovery

import "fmt"

const (
	sharedEnrollmentClassify    = uint64(0x1005419b4)
	sharedEnrollmentBranch      = uint64(0x1005649c4)
	sharedEnrollmentWrite       = uint64(0x1003f63b0)
	sharedEnrollmentAdvance     = uint64(0x1003f6400)
	sharedEnrollmentByteFill    = uint64(0x10051c16c)
	sharedEnrollmentLeftPath    = uint64(0x1005456ac)
	sharedEnrollmentRightPath   = uint64(0x1003d3294)
	sharedEnrollmentWordCompare = uint64(0x10055fb1c)
	sharedEnrollmentLoopSelect  = uint64(0x10027b7c8)
	sharedEnrollmentLoopTail    = uint64(0x10027b800)
	sharedEnrollmentSavedSelect = uint64(0x1003ea090)
)

// runSharedEnrollmentBlock replaces the small protected dispatcher core that
// is shared by natural record selection and current-SINF enrollment. Branch-
// specific transforms continue through their verified CFG blocks.
func runSharedEnrollmentBlock(s *state, pc uint64) (uint64, bool, error) {
	if s == nil {
		return 0, true, fmt.Errorf("shared enrollment state is unavailable")
	}
	switch pc {
	case sharedEnrollmentClassify:
		next, err := runSharedEnrollmentClassify(s)
		return next, true, err
	case sharedEnrollmentBranch:
		next, err := runSharedEnrollmentBranch(s)
		return next, true, err
	case sharedEnrollmentWrite:
		next, err := runSharedEnrollmentWrite(s)
		return next, true, err
	case sharedEnrollmentAdvance:
		next, err := runSharedEnrollmentAdvance(s)
		return next, true, err
	case sharedEnrollmentByteFill:
		next, err := runSharedEnrollmentByteFill(s)
		return next, true, err
	case sharedEnrollmentLeftPath:
		next, err := runSharedEnrollmentLeftPath(s)
		return next, true, err
	case sharedEnrollmentRightPath:
		next, err := runSharedEnrollmentRightPath(s)
		return next, true, err
	case sharedEnrollmentWordCompare:
		next, err := runSharedEnrollmentWordCompare(s)
		return next, true, err
	case sharedEnrollmentLoopSelect:
		next, err := runSharedEnrollmentLoopSelect(s)
		return next, true, err
	case sharedEnrollmentSavedSelect:
		next, err := runSharedEnrollmentSavedSelect(s)
		return next, true, err
	case sharedEnrollmentBytePermute:
		next, err := runSharedEnrollmentBytePermute(s)
		return next, true, err
	default:
		return pc, false, nil
	}
}

func runSharedEnrollmentSavedSelect(s *state) (uint64, error) {
	s.traceBlock(sharedEnrollmentSavedSelect)
	s.x[8] = 0x710bceae
	if err := s.storeValue(s.sp+0x3b3c, uint64(uint32(s.x[12])), 4); err != nil {
		return 0, fmt.Errorf("shared enrollment saved selector: %w", err)
	}
	s.x[9] = uint64(uint32(s.x[12]) + 0xf1fc0501)
	s.compare32(uint32(s.x[9]), 0x16c)
	selector := (uint32(s.x[19]) ^ uint32(boolValue(s.condition("lo")))) & 1
	s.x[9] = uint64(selector * 0x31)
	s.x[8] = uint64(uint32(s.x[19]) + uint32(s.x[8]))
	dispatchIndex := uint32(s.x[9]) ^ uint32(s.x[8])
	dispatchTable, err := s.load(s.sp+0xdf0, 8)
	if err != nil {
		return 0, fmt.Errorf("shared enrollment saved-select table: %w", err)
	}
	s.x[9] = dispatchTable
	next, err := loadBootstrapPointer(s, dispatchTable, dispatchIndex)
	if err != nil {
		return 0, fmt.Errorf("shared enrollment saved-select dispatch: %w", err)
	}
	s.x[8] = next
	s.x[17] = 0x58c6
	return next, nil
}

func runSharedEnrollmentWordCompare(s *state) (uint64, error) {
	s.traceBlock(sharedEnrollmentWordCompare)
	s.x[12] = 0x526c
	left, err := s.load(s.x[9]+uint64(uint32(s.x[8]))*4, 4)
	if err != nil {
		return 0, fmt.Errorf("shared enrollment word-compare left: %w", err)
	}
	s.x[11] = uint64(uint32(left))
	right, err := s.load(s.x[9]+uint64(uint32(s.x[10]))*4, 4)
	if err != nil {
		return 0, fmt.Errorf("shared enrollment word-compare right: %w", err)
	}
	s.x[13] = uint64(uint32(right))
	s.compare32(uint32(s.x[11]), uint32(s.x[13]))
	different := uint32(boolValue(s.condition("ne")))
	s.x[13] = uint64(different | different<<3)
	s.x[12] = uint64(uint32(s.x[24]) + uint32(s.x[12]))
	dispatchIndex := uint32(s.x[13]) ^ uint32(s.x[12])
	next, err := loadBootstrapPointer(s, s.x[21], dispatchIndex)
	if err != nil {
		return 0, fmt.Errorf("shared enrollment word-compare dispatch: %w", err)
	}
	s.x[12] = next
	s.x[17] = 0x7cc
	return next, nil
}

func runSharedEnrollmentLoopSelect(s *state) (uint64, error) {
	s.traceBlock(sharedEnrollmentLoopSelect)
	s.x[11] = uint64(uint32(s.x[24]) + 0x3d68)
	s.x[12] = 0x8e1
	s.x[11] = uint64(uint32(s.x[11]) | uint32(s.x[12]))
	s.x[12] = 0xc153
	s.x[10] = uint64(uint32(s.x[10]) + 1)
	s.compare32(uint32(s.x[10]), 0x2c)
	more := uint32(boolValue(s.condition("lo")))
	s.x[14] = uint64(uint32(s.x[11]) ^ 0x4aea)
	s.x[13] = uint64(more * uint32(s.x[14]))
	s.x[12] = uint64(uint32(s.x[24]) + uint32(s.x[12]))
	dispatchIndex := uint32(s.x[13]) ^ uint32(s.x[12])
	next, err := loadBootstrapPointer(s, s.x[21], dispatchIndex)
	if err != nil {
		return 0, fmt.Errorf("shared enrollment loop-select dispatch: %w", err)
	}
	s.x[12] = next
	s.traceBlock(sharedEnrollmentLoopTail)
	s.x[17] = 0xc87e
	return next, nil
}

func runSharedEnrollmentLeftPath(s *state) (uint64, error) {
	s.traceBlock(sharedEnrollmentLeftPath)
	s.x[17] = uint64(uint32(s.x[24]) + 0xdd472f56)
	s.x[1] = 0x22b83dd7
	s.x[2] = uint64(uint32(s.x[17]) & uint32(s.x[1]))
	s.x[17] = 0xffff6d2e
	s.x[19] = uint64(uint32(s.x[24]) + uint32(s.x[17]))
	s.x[17] = uint64(uint32(s.x[19]) ^ uint32(s.x[2]))
	s.x[17] = uint64(uint32(s.x[17]) ^ 0xfffffff7)
	s.compare32(uint32(s.x[23]), uint32(int32(int8(s.x[17]))))
	dispatchBit := uint32(boolValue(s.condition("lt")))
	dispatchIndex := dispatchBit*0x49 ^ uint32(s.x[24])
	next, err := loadBootstrapPointer(s, s.x[21], dispatchIndex)
	if err != nil {
		return 0, fmt.Errorf("shared enrollment left-path dispatch: %w", err)
	}
	s.x[1] = next
	s.x[17] = 0x60e8
	return next, nil
}

func runSharedEnrollmentRightPath(s *state) (uint64, error) {
	s.traceBlock(sharedEnrollmentRightPath)
	s.x[17] = uint64(uint32(s.x[25]) + 0xffffce15)
	s.x[1] = uint64(uint32(s.x[17]) | 0x3200)
	s.x[17] = 0xffffa626
	s.x[2] = uint64(uint32(s.x[1]) + uint32(s.x[17]))
	s.x[17] = uint64(uint32(s.x[25]) + 0xffff745a)
	s.x[22] = uint64(uint32(s.x[23]) & 0xff)
	s.compare32(uint32(s.x[22]), uint32(byte(s.x[17])))
	dispatchBit := uint32(boolValue(s.condition("ne")))
	dispatchIndex := dispatchBit*0x2b ^ uint32(s.x[25])
	next, err := loadBootstrapPointer(s, s.x[21], dispatchIndex)
	if err != nil {
		return 0, fmt.Errorf("shared enrollment right-path dispatch: %w", err)
	}
	s.x[17] = 0xbf5f
	s.x[22] = next
	return next, nil
}

func runSharedEnrollmentByteFill(s *state) (uint64, error) {
	const maximumBytes = 16384
	for index := 0; index < maximumBytes; index++ {
		next, err := runSharedEnrollmentByteFillStep(s)
		if err != nil || next != sharedEnrollmentByteFill {
			return next, err
		}
	}
	return 0, fmt.Errorf("shared enrollment byte fill exceeded %d bytes", maximumBytes)
}

func runSharedEnrollmentByteFillStep(s *state) (uint64, error) {
	s.traceBlock(sharedEnrollmentByteFill)
	s.x[11] = uint64(uint32(s.x[9]) + uint32(s.x[30]))
	s.x[8]--
	if err := s.storeValue(s.sp+0x3af8, s.x[8], 8); err != nil {
		return 0, fmt.Errorf("shared enrollment byte-fill counter: %w", err)
	}
	s.x[12] = uint64(uint32(s.x[10]) ^ uint32(s.x[1]))
	s.x[13] = uint64(uint32(s.x[10]) + 0x11)
	s.x[12] = uint64(uint32(s.x[12]) * uint32(s.x[13]))
	if err := s.storeValue(s.x[20]+s.x[8], uint64(byte(s.x[12])), 1); err != nil {
		return 0, fmt.Errorf("shared enrollment byte-fill output: %w", err)
	}
	s.x[8] = s.x[10] + s.x[3]
	s.x[10] = s.x[8] + s.x[11]
	counter, err := s.load(s.sp+0x3af8, 8)
	if err != nil {
		return 0, fmt.Errorf("shared enrollment byte-fill saved counter: %w", err)
	}
	s.x[8] = counter
	s.compare64(counter, 0)
	last := uint32(boolValue(s.condition("eq")))
	dispatchIndex := uint32(s.x[9]) ^ (last << 3)
	dispatchTable, err := s.load(s.sp+0xdf0, 8)
	if err != nil {
		return 0, fmt.Errorf("shared enrollment byte-fill dispatch table: %w", err)
	}
	s.x[12] = dispatchTable
	next, err := loadBootstrapPointer(s, dispatchTable, dispatchIndex)
	if err != nil {
		return 0, fmt.Errorf("shared enrollment byte-fill dispatch: %w", err)
	}
	s.x[11] = next
	s.x[17] = 0x3548
	return next, nil
}

func runSharedEnrollmentClassify(s *state) (uint64, error) {
	s.traceBlock(sharedEnrollmentClassify)
	s.x[16] = uint64(uint32(s.x[12]) - 0xae)
	s.x[17] = uint64((uint32(s.x[16]) ^ 0x3a) + 0x6d)
	s.x[0] = uint64((uint32(s.x[16]) ^ 0x584ee074) + 0x23)
	s.x[1] = uint64((uint32(s.x[16]) ^ 0x18) + 0x4f)
	s.x[19] = uint64(uint32(s.x[13]) ^ uint32(s.x[16]))
	s.x[0] = uint64(uint32(s.x[19]) ^ uint32(s.x[0]))
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[1]))
	s.x[17] = uint64(uint32(s.x[0]) ^ uint32(s.x[17]))
	s.x[0] = 0x380f
	s.x[1] = uint64(uint32(s.x[2]) + uint32(s.x[0]))
	s.x[0] = 0x5f43
	s.x[0] = uint64(uint32(s.x[1]) ^ uint32(s.x[0]))
	s.x[19] = uint64(uint32(s.x[17]) & uint32(s.x[0]))
	s.x[17] = 0xa9
	s.x[0] = s.x[19] ^ s.x[17]
	value, err := s.load(s.x[28]+s.x[0], 1)
	if err != nil {
		return 0, fmt.Errorf("shared enrollment classifier: %w", err)
	}
	s.x[20] = uint64(byte(value))
	s.x[17] = uint64(uint32(s.x[20]) ^ 0x22e6ee17)
	s.x[22] = 0xdd1911e9
	s.x[17] = uint64(uint32(s.x[17]) + uint32(s.x[22]))
	s.x[22] = uint64(uint32(s.x[20]) ^ 0x6482ee5e)
	s.x[23] = 0x9b7d11a2
	s.x[22] = uint64(uint32(s.x[22]) + uint32(s.x[23]))
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[22]))
	s.x[22] = uint64(uint32(s.x[20]) ^ 0x46640052)
	s.x[23] = 0xb99bffae
	s.x[22] = uint64(uint32(s.x[22]) + uint32(s.x[23]))
	s.x[22] = uint64(uint32(s.x[17]) ^ uint32(s.x[22]))
	s.x[23] = uint64(0xfffffff4 | (uint32(s.x[20]) & 3))
	dispatchIndex := uint32(s.x[2])
	next, err := loadBootstrapPointer(s, s.x[21], dispatchIndex)
	if err != nil {
		return 0, fmt.Errorf("shared enrollment classifier dispatch: %w", err)
	}
	s.x[2] = next
	s.x[17] = 0x1359
	return next, nil
}

func runSharedEnrollmentBranch(s *state) (uint64, error) {
	s.traceBlock(sharedEnrollmentBranch)
	s.x[17] = uint64(uint32(s.x[1]) + 0x4064)
	s.x[2] = 0xc88
	s.x[24] = uint64(uint32(s.x[17]) | uint32(s.x[2]))
	s.x[17] = 0x13df
	s.x[25] = uint64(uint32(s.x[1]) + uint32(s.x[17]))
	s.x[17] = uint64(uint32(s.x[25]) ^ 0x6d)
	right := uint32(int32(int8(s.x[17])))
	s.compare32(uint32(s.x[23]), right)
	selected := uint32(boolValue(s.condition("ge")))
	s.x[2] = uint64(uint32(s.x[24]) ^ 0xace3)
	dispatchIndex := selected*uint32(s.x[2]) ^ uint32(s.x[1])
	next, err := loadBootstrapPointer(s, s.x[21], dispatchIndex)
	if err != nil {
		return 0, fmt.Errorf("shared enrollment branch dispatch: %w", err)
	}
	s.x[1] = next
	s.x[17] = 0x8eec
	return next, nil
}

func runSharedEnrollmentWrite(s *state) (uint64, error) {
	s.traceBlock(sharedEnrollmentWrite)
	s.compare32(uint32(s.x[9]), uint32(s.x[3]))
	gate := uint32(boolValue(s.condition("lo")))
	s.x[17] = uint64(gate)
	s.x[2] = uint64(uint32(s.x[1]) + 0xffffbe48)
	s.x[19] = 0x3a5
	s.x[2] = uint64(uint32(s.x[2]) | uint32(s.x[19]))
	s.x[19] = uint64(uint32(s.x[1]) + 0xffff9c4d)
	s.x[22] = uint64(uint32(s.x[19]) | 0x48e2)
	if err := s.storeValue(s.x[28]+s.x[0], uint64(byte(s.x[20])), 1); err != nil {
		return 0, fmt.Errorf("shared enrollment output byte: %w", err)
	}
	s.x[0] = uint64(uint32(s.x[22]) ^ 0xcad53862)
	s.x[19] = 0xa7e380e5
	s.x[16] = uint64(uint32(s.x[16]) + uint32(s.x[19]))
	s.compare32(uint32(s.x[16]), uint32(s.x[0]))
	s.x[0] = uint64(gate ^ uint32(boolValue(s.condition("lo"))))
	return sharedEnrollmentAdvance, nil
}

func runSharedEnrollmentAdvance(s *state) (uint64, error) {
	s.traceBlock(sharedEnrollmentAdvance)
	s.compare32(uint32(s.x[16]), uint32(s.x[10]))
	secondBound := uint32(boolValue(s.condition("lo")))
	s.compare32(uint32(s.x[0]), 0)
	dispatchBit := secondBound
	if s.condition("ne") {
		dispatchBit = uint32(s.x[17])
	}
	dispatchBit ^= 1
	dispatchMultiplier := uint32(s.x[2]) ^ uint32(s.x[4])
	dispatchIndex := dispatchBit*dispatchMultiplier ^ uint32(s.x[1])
	s.x[16] = uint64(dispatchIndex)
	next, err := loadBootstrapPointer(s, s.x[21], dispatchIndex)
	if err != nil {
		return 0, fmt.Errorf("shared enrollment advance dispatch: %w", err)
	}
	s.x[0] = next
	s.x[11] = uint64(uint32(s.x[11]) + 1)
	s.x[12] = uint64(uint32(s.x[12]) + 1)
	s.x[13] = uint64(uint32(s.x[13]) - 1)
	s.x[14] = uint64(uint32(s.x[14]) - 5)
	s.x[15] = uint64(uint32(s.x[15]) + 5)
	s.x[17] = 0x4799
	return next, nil
}
