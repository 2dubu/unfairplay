package recovery

import "fmt"

const (
	sharedEnrollmentClassify = uint64(0x1005419b4)
	sharedEnrollmentBranch   = uint64(0x1005649c4)
	sharedEnrollmentWrite    = uint64(0x1003f63b0)
	sharedEnrollmentAdvance  = uint64(0x1003f6400)
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
	default:
		return pc, false, nil
	}
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
