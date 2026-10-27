package recovery

import "fmt"

const (
	bootstrapPointerCycleStart = uint64(0x1001dfbf4)
	bootstrapPointerCycleTail  = uint64(0x1001dfc00)
	bootstrapPointerCycleCheck = uint64(0x1000602d8)
	bootstrapPointerCycleStep  = uint64(0x100066ee0)
	bootstrapPointerCycleNext  = uint64(0x100087a4c)
	bootstrapByteClassifyPass  = uint64(0x100065d6c)
	bootstrapIndexedDispatch   = uint64(0x1005816cc)
	bootstrapCountdownStart    = uint64(0x1005807e0)
	bootstrapCountdownDispatch = uint64(0x100580800)
	bootstrapSmallFrameReturn  = uint64(0x10016e74c)
	bootstrapSavedFrameReturn  = uint64(0x100147c4c)
	bootstrapProtectedCall     = uint64(0x10034d758)
	bootstrapProtectedDispatch = uint64(0x10034d800)
	bootstrapRangeDispatchA    = uint64(0x10015dd44)
	bootstrapRangeDispatchB    = uint64(0x100040340)
	bootstrapStateDispatch     = uint64(0x100187d9c)
	bootstrapWordDispatch      = uint64(0x1000f4834)
	bootstrapScalarRoundStart  = uint64(0x10092d608)
	bootstrapScalarRoundTail   = uint64(0x10092d800)
)

func runBootstrapScalarRound(s *state) (uint64, error) {
	roundDone, err := runBootstrapScalarRoundFront(s)
	if err != nil {
		return 0, err
	}
	return runBootstrapScalarRoundFinish(s, roundDone)
}

func runBootstrapScalarRoundFront(s *state) (bool, error) {
	s.traceBlock(bootstrapScalarRoundStart)
	s.x[17], s.x[6], s.x[13], s.x[24] = s.x[6], s.x[13], s.x[24], s.x[4]
	s.x[14] = 0x7a1df463
	s.x[3] = uint64(uint32(s.x[6]) ^ uint32(s.x[14]))
	s.x[14] = uint64(uint32(s.x[3]) ^ uint32(s.x[13]))
	s.x[4] = uint64(uint32(s.x[4]) ^ uint32(s.x[30]))
	s.x[26] = uint64((uint32(s.x[24]) ^ 0xdf956bdd) + 0x206a9423)
	s.x[27] = uint64((uint32(s.x[24]) ^ uint32(s.x[19])) + 0xea05010a)
	s.x[26] = uint64(uint32(s.x[26]) ^ uint32(s.x[27]))
	s.x[4] = uint64(uint32(s.x[4]) + 0xb0726149)
	s.x[4] = uint64(uint32(s.x[4]) ^ uint32(s.x[26]))
	s.x[1] = 0x1e877d19
	s.compare32(uint32(s.x[1]), uint32(s.x[4])>>2)
	if s.condition("eq") {
		s.x[1] = uint64(uint32(s.x[13]))
	} else {
		s.x[1] = uint64(uint32(s.x[14]))
	}
	s.x[4] = uint64(uint32(s.x[25]) * 0x8f)
	s.x[14] = uint64(uint32(s.x[4]) & 0xffff)
	s.x[14] = uint64(uint32(s.x[14]) * uint32(s.x[21]))
	s.x[14] = uint64(uint32(s.x[14]) >> 23)
	s.x[14] = uint64(uint32(s.x[4]) - uint32(s.x[14])*uint32(s.x[22]))
	s.x[26] = uint64(uint32(s.x[6]) ^ 0x9cfc962f)
	s.x[27] = uint64(uint32(s.x[1]) ^ 0x85e20b9c)
	s.x[2] = uint64(uint32(s.x[27]) & uint32(s.x[26]))
	s.x[4] = uint64(uint32(s.x[4]) + 0x2838)
	s.x[26] = uint64(uint32(s.x[4]) & 0xffff)
	s.x[26] = uint64(uint32(s.x[26]) * uint32(s.x[21]))
	s.x[26] = uint64(uint32(s.x[26]) >> 23)
	s.x[4] = uint64(uint32(s.x[4]) - uint32(s.x[26])*uint32(s.x[22]))
	s.x[26] = uint64(uint32(s.x[13]) ^ 0x3ff4b9cb)
	s.x[3] = uint64(uint32(s.x[3]) & uint32(s.x[26]))
	s.x[26] = uint64(uint32(s.x[6]) & 0x45e94da8)
	s.x[3] = uint64(uint32(s.x[3]) ^ uint32(s.x[26]))
	s.x[26] = uint64(uint32(s.x[24]) ^ uint32(s.x[12]))
	s.x[27] = uint64(uint32(s.x[3]) ^ 0xbff6bbdf)
	s.x[26] = uint64(uint32(s.x[27]) & uint32(s.x[26]))
	s.x[3] = uint64(uint32(s.x[3]) & 0xc811fec0)
	s.x[3] = uint64(uint32(s.x[26]) ^ uint32(s.x[3]))
	s.x[27], s.x[26], s.x[0] = s.x[0], s.x[16], s.x[11]
	s.x[11] = uint64(uint32(s.x[1]) & uint32(s.x[28]))
	s.x[11] = uint64(uint32(s.x[11]) ^ uint32(s.x[15]))
	s.x[11] = uint64(uint32(s.x[11]) ^ uint32(s.x[2]))
	s.x[16] = uint64(uint32(s.x[3]) ^ 0x77ef453f)
	s.x[11] = uint64(uint32(s.x[11]) & uint32(s.x[16]))
	s.x[16] = uint64(uint32(s.x[3]) & 0x133cf0b7)
	s.x[16] = uint64(uint32(s.x[16]) ^ uint32(s.x[5]))
	s.x[11] = uint64(uint32(s.x[16]) ^ uint32(s.x[11]))
	s.x[16] = uint64(uint32(s.x[0]) ^ uint32(s.x[27]))
	s.x[16] = uint64(uint32(s.x[16]) ^ uint32(s.x[8]))
	s.x[1] = uint64(uint32(s.x[26]) ^ uint32(s.x[20]))
	s.x[16] = uint64(uint32(s.x[16]) & uint32(s.x[1]))
	s.x[1] = uint64(uint32(s.x[26]) & uint32(s.x[8]))
	s.x[1] = uint64(uint32(s.x[1]) | 0xf0825106)
	s.x[1] = uint64(uint32(s.x[27]) ^ uint32(s.x[1]))
	s.x[16] = uint64(uint32(s.x[1]) ^ uint32(s.x[16]))
	s.x[1] = uint64(uint32(rotateRight(s.x[26], 6, 32)))
	s.x[1] = uint64(uint32(s.x[1]) ^ uint32(s.x[26]<<21))
	s.x[1] = uint64(uint32(s.x[1]) ^ uint32(s.x[26])>>11)
	s.x[1] = uint64(uint32(s.x[1]) ^ uint32(rotateRight(s.x[26], 25, 32)))
	s.x[1] = uint64(uint32(s.x[1]) ^ 0xfce9c4ea)
	s.x[2] = uint64(uint32(s.x[7]) ^ uint32(s.x[20]))
	s.x[16] = uint64(uint32(s.x[2]) + uint32(s.x[16]))
	s.x[2] = uint64(uint32(s.x[17]) ^ 0x0a8e3c74)
	s.x[3] = (s.x[4] & 0xffff) << 2
	value, err := s.load(s.x[9]+s.x[3], 4)
	if err != nil {
		return false, fmt.Errorf("bootstrap scalar-round first lookup: %w", err)
	}
	s.x[3] = uint64(uint32(value))
	s.x[16] = uint64(uint32(s.x[1]) + uint32(s.x[16]))
	s.x[1] = uint64(uint32(s.x[16]) & uint32(s.x[3]))
	s.x[16] = uint64(uint32(s.x[16]) ^ uint32(s.x[3]))
	s.x[3] = uint64(uint32(s.x[17]) ^ 0x40edb29d)
	s.x[4] = 0xcf818575
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[4]))
	s.x[14] = (s.x[14] & 0xffff) << 2
	value, err = s.load(s.x[9]+s.x[14], 4)
	if err != nil {
		return false, fmt.Errorf("bootstrap scalar-round second lookup: %w", err)
	}
	s.x[14] = uint64(uint32(value))
	s.x[14] = uint64(uint32(s.x[14]) + uint32(s.x[16]))
	s.x[16] = uint64(uint32(rotateRight(s.x[24], 2, 32)))
	s.x[16] = uint64(uint32(s.x[16]) ^ uint32(rotateRight(s.x[24], 13, 32)))
	s.x[16] = uint64(uint32(s.x[16]) ^ uint32(rotateRight(s.x[24], 22, 32)))
	s.x[16] = uint64(uint32(s.x[16]) ^ uint32(s.x[10]))
	s.compare32(uint32(s.x[25]), 0x3f)
	roundDone := s.condition("eq")
	s.x[25] = uint64(uint32(s.x[25]) + 1)
	s.x[14] = uint64(uint32(s.x[14]) + uint32(s.x[1]<<1))
	s.x[11] = uint64(uint32(s.x[16]) + uint32(s.x[11]))
	s.x[16] = 0x53a8b866
	s.x[11] = uint64(uint32(s.x[11]) + uint32(s.x[16]))
	return roundDone, nil
}

func runBootstrapScalarRoundFinish(s *state, roundDone bool) (uint64, error) {
	s.traceBlock(bootstrapScalarRoundTail)
	s.x[16] = uint64(uint32(s.x[14]) + 0x71c34584)
	s.x[11] = uint64(uint32(s.x[11]) + uint32(s.x[16]))
	s.x[14] = uint64((uint32(s.x[11]) ^ 0x0abde130) + 0x4d76a8a6)
	s.x[1] = uint64((uint32(s.x[11]) ^ 0x7da68b4d) + 0x3a6dc2d9)
	s.x[14] = uint64(uint32(s.x[14]) ^ uint32(s.x[1]))
	s.x[1] = uint64((uint32(s.x[11]) ^ 0x10f023e8) + 0x573b6a7e)
	s.x[1] = uint64(uint32(s.x[1]) ^ uint32(s.x[11]))
	s.x[4] = 0xdfdfffff
	s.x[11] = uint64(uint32(s.x[11]) ^ uint32(s.x[4]))
	s.x[14] = uint64(uint32(s.x[1]) ^ uint32(s.x[14]))
	s.x[11] = uint64(uint32(s.x[11]) + 0x9814b66b)
	s.x[14] = uint64(uint32(s.x[14]) ^ uint32(s.x[11]))
	s.x[11] = uint64(uint32(s.x[14]) ^ 0x3dd6bdf6)
	s.x[4] = s.x[11]
	s.x[11] = uint64(uint32(s.x[2]) + 0xf571c38c)
	s.x[1] = uint64(uint32(s.x[3]) + 0xbf124d63)
	s.x[11] = uint64(uint32(s.x[11]) ^ uint32(s.x[1]))
	s.x[17] = uint64(uint32(s.x[17]) + 0x307e7a8b)
	s.x[11] = uint64(uint32(s.x[11]) ^ uint32(s.x[17]))
	s.x[17] = 0x31b4e36b
	s.x[11] = uint64(uint32(s.x[11]) + uint32(s.x[17]))
	s.x[11] = uint64(uint32(s.x[11]) + uint32(s.x[16]))
	s.x[16] = uint64(uint32(s.x[11]) ^ 0x997804c8)
	s.x[17] = uint64(uint32(s.x[11]) ^ 0x62d57dbb)
	s.x[17] = uint64(uint32(s.x[17]) + 0x8d745798)
	s.x[1] = uint64(uint32(s.x[11]) ^ 0x940c935f)
	s.x[2] = uint64(uint32(s.x[11]) ^ 0x7fff3fff)
	s.x[16] = uint64(uint32(s.x[16]) + 0x76d92ee5)
	s.x[11] = uint64(uint32(s.x[16]) ^ uint32(s.x[11]))
	s.x[16] = uint64(uint32(s.x[1]) + 0x7badb974)
	s.x[16] = uint64(uint32(s.x[17]) ^ uint32(s.x[16]))
	s.x[11] = uint64(uint32(s.x[11]) ^ uint32(s.x[16]))
	s.x[16] = uint64(uint32(s.x[2]) + 0x905e15d4)
	s.x[3] = uint64(uint32(s.x[11]) ^ uint32(s.x[16]))
	s.x[16] = uint64(uint32(s.x[3]) ^ 0xe2dd2ee5)
	s.x[11] = uint64(uint32(boolValue(roundDone)) * 0x16)
	s.x[11] = uint64(uint32(s.x[11]) ^ uint32(s.x[23]))
	dispatchTable, err := s.load(s.sp+0x208, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap scalar-round dispatch table: %w", err)
	}
	s.x[17] = dispatchTable
	next, err := loadBootstrapPointer(s, dispatchTable, uint32(s.x[11]))
	if err != nil {
		return 0, fmt.Errorf("bootstrap scalar-round dispatch: %w", err)
	}
	s.x[1] = next
	s.x[11] = s.x[26]
	s.x[7] = s.x[27]
	s.x[17] = 0x5a72
	return next, nil
}

func runBootstrapStateDispatch(s *state) (uint64, error) {
	s.traceBlock(bootstrapStateDispatch)
	s.x[12] = uint64(uint32(s.x[14]) + 0xaf200e39)
	s.x[28] = uint64(uint32(s.x[12]) & 0x50df7fdf)
	s.x[12] = uint64(uint32(s.x[15]) + 0x042a8d93)
	s.x[13] = 0xffffc0b3
	s.x[16] = uint64(0xf7aae4da - uint32(s.x[15]))
	s.x[13] = uint64(uint32(s.x[28]) + uint32(s.x[13]))
	s.compare32(uint32(s.x[12]), uint32(s.x[13]))
	if s.condition("gt") {
		s.x[12] = uint64(uint32(s.x[16]))
	} else {
		s.x[12] = uint64(uint32(s.x[15]))
	}
	s.x[13] = 0x5eec043e
	if s.condition("gt") {
		s.x[16] = uint64(uint32(s.x[11]))
	} else {
		s.x[16] = uint64(uint32(s.x[13]))
	}
	next, err := loadBootstrapPointer(s, s.x[10], uint32(s.x[14]))
	if err != nil {
		return 0, fmt.Errorf("bootstrap state-dispatch target: %w", err)
	}
	s.x[11] = next
	s.x[30] = 0x43522004
	s.x[13] = 0x21037cc9
	s.x[1] = 0xd9b0b007
	s.x[7] = 0x264f10ad
	s.x[2] = 0x73b42a39
	s.x[19] = 0x8c4bd5c7
	s.x[20] = 0x1e35d90b
	s.x[21] = 0xe1ca26f5
	s.x[22] = 0xa9aa7c79
	s.x[23] = 0x56558387
	s.x[24] = 0x95d10649
	s.x[25] = 0x6a2ef9b7
	s.x[26] = 0x9d987e50
	s.x[27] = 0x7b49db15
	s.x[14] = 0xb
	s.x[17] = 0x1a6b
	return next, nil
}

func runBootstrapWordDispatch(s *state) (uint64, error) {
	s.traceBlock(bootstrapWordDispatch)
	firstBase, err := s.load(s.x[8]+8, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap word-dispatch first base: %w", err)
	}
	s.x[11] = firstBase
	s.x[15] = 0xfffffffe10aa3648 + (s.x[12] << 2)
	first, err := s.load(s.x[11]+s.x[15], 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap word-dispatch first word: %w", err)
	}
	s.x[3] = uint64(uint32(first))
	secondBase, err := s.load(s.x[9]+8, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap word-dispatch second base: %w", err)
	}
	s.x[11] = secondBase
	second, err := s.load(s.x[11]+s.x[15], 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap word-dispatch second word: %w", err)
	}
	s.x[15] = uint64(uint32(second))
	s.x[4] = uint64(uint32(s.x[5]) ^ 0x421)
	s.x[11] = uint64((uint32(s.x[3]) ^ 0x037eabc3) + 0xfc81543d)
	s.x[17] = uint64(uint32(s.x[4]) + 0xffff8dfb)
	s.x[17] = uint64(uint32(s.x[17]) ^ 0x7d619862)
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[3]))
	s.x[28] = uint64((uint32(s.x[3]) ^ 0x5c51afb4) + 0xa3ae504c)
	s.x[11] = uint64(uint32(s.x[28]) ^ uint32(s.x[11]))
	s.x[17] = uint64(uint32(s.x[17]) + 0x829e58b4)
	s.x[11] = uint64(uint32(s.x[11]) ^ uint32(s.x[17]))
	s.x[17] = uint64((uint32(s.x[15]) ^ 0xbdf50467) + 0x420afb99)
	s.x[28] = uint64((uint32(s.x[15]) ^ 0xdb3d9f3a) + 0x24c260c6)
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[28]))
	s.x[28] = uint64((uint32(s.x[15]) ^ 0x44863866) + 0xbb79c79a)
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[28]))
	s.x[11] = uint64(uint32(s.x[11]) + uint32(s.x[30]))
	s.compare32(uint32(s.x[11]), uint32(s.x[13]))
	firstBound := uint32(boolValue(s.condition("lo")))
	s.x[28] = uint64(firstBound)
	s.x[17] = uint64(uint32(s.x[17]) + uint32(s.x[30]))
	s.compare32(uint32(s.x[17]), uint32(s.x[13]))
	secondBound := uint32(boolValue(s.condition("lo")))
	s.x[6] = uint64(firstBound ^ secondBound)
	s.compare32(uint32(s.x[11]), uint32(s.x[17]))
	higher := uint32(boolValue(s.condition("hi")))
	s.x[11] = uint64(higher)
	s.compare32(uint32(s.x[6]), 0)
	if s.condition("ne") {
		s.x[11] = uint64(firstBound)
	}
	s.x[11] = uint64(uint32(s.x[11]) ^ 1)
	s.x[11] = uint64(uint32(s.x[11]) * 0x51)
	s.x[11] = uint64(uint32(s.x[11]) ^ uint32(s.x[5]))
	next, err := loadBootstrapPointer(s, s.x[10], uint32(s.x[11]))
	if err != nil {
		return 0, fmt.Errorf("bootstrap word-dispatch target: %w", err)
	}
	s.x[5] = next
	s.x[11] = s.x[16]
	s.x[17] = 0x578f
	return next, nil
}

func runBootstrapRangeDispatchA(s *state) (uint64, error) {
	s.traceBlock(bootstrapRangeDispatchA)
	s.x[11] = uint64((uint32(s.x[28]) + 0xffda80b3) & 0x00257f2e)
	s.x[5] = uint64(uint32(s.x[28]) ^ 0x8a5e)
	s.x[15] = uint64(uint32(s.x[12]) + 0x842a8d92)
	s.x[17] = uint64(uint32(s.x[5]) + 0xffff4aed)
	s.compare32(uint32(s.x[15]), uint32(s.x[17]))
	s.x[15] = uint64(uint32(boolValue(s.condition("gt"))))
	s.x[11] = uint64(uint32(s.x[11]) ^ 0x4060)
	s.x[11] = uint64(uint32(s.x[15]) * uint32(s.x[11]))
	s.x[11] = uint64(uint32(s.x[11]) ^ uint32(s.x[28]))
	next, err := loadBootstrapPointer(s, s.x[10], uint32(s.x[11]))
	if err != nil {
		return 0, fmt.Errorf("bootstrap range-dispatch A: %w", err)
	}
	s.x[15] = next
	s.x[11] = 0x5eec043d
	s.x[17] = 0x63e8
	return next, nil
}

func runBootstrapRangeDispatchB(s *state) (uint64, error) {
	s.traceBlock(bootstrapRangeDispatchB)
	s.x[11] = 0xffffe558
	s.x[16] = 0xac2aa0cd
	s.x[17] = 0xac2ad0ac
	s.x[14] = uint64(uint32(s.x[13]) + uint32(s.x[11]))
	s.x[11] = uint64(uint32(s.x[14]) ^ uint32(s.x[17]))
	s.x[17] = 0xffff7353
	s.x[12] = uint64(uint32(s.x[12]) + uint32(s.x[13]))
	s.x[12] = uint64(uint32(s.x[12]) + uint32(s.x[17]))
	s.x[17] = 0xa800133a
	s.x[16] = uint64(uint32(s.x[15]) + uint32(s.x[16]))
	s.compare32(uint32(s.x[16]), uint32(s.x[17]))
	s.x[1] = uint64(uint32(boolValue(s.condition("lt"))))
	s.x[11] = uint64(uint32(s.x[12]) + uint32(s.x[11]))
	s.compare32(uint32(s.x[11]), uint32(s.x[17]))
	s.x[12] = uint64(uint32(boolValue(s.condition("lt"))))
	s.x[17] = uint64(uint32(s.x[1]) ^ uint32(s.x[12]))
	s.compare32(uint32(s.x[16]), uint32(s.x[11]))
	s.x[11] = uint64(uint32(boolValue(s.condition("lt"))))
	s.compare32(uint32(s.x[17]), 0)
	if s.condition("ne") {
		s.x[11] = uint64(uint32(s.x[12]))
	}
	s.x[11] = uint64(uint32(s.x[11]) * 0x3d)
	s.x[11] = uint64(uint32(s.x[11]) ^ uint32(s.x[13]))
	next, err := loadBootstrapPointer(s, s.x[10], uint32(s.x[11]))
	if err != nil {
		return 0, fmt.Errorf("bootstrap range-dispatch B: %w", err)
	}
	s.x[12] = next
	s.x[11] = 0x5eec043c
	s.x[17] = 0xfe60
	return next, nil
}

func runBootstrapProtectedCall(s *state) (uint64, error) {
	s.traceBlock(bootstrapProtectedCall)
	s.sp -= 0x70
	for _, pair := range [...]struct {
		first, second int
		offset        uint64
	}{{28, 27, 0x10}, {26, 25, 0x20}, {24, 23, 0x30}, {22, 21, 0x40}, {20, 19, 0x50}, {29, 30, 0x60}} {
		if err := s.storeValue(s.sp+pair.offset, s.x[pair.first], 8); err != nil {
			return 0, fmt.Errorf("bootstrap protected-call register %d: %w", pair.first, err)
		}
		if err := s.storeValue(s.sp+pair.offset+8, s.x[pair.second], 8); err != nil {
			return 0, fmt.Errorf("bootstrap protected-call register %d: %w", pair.second, err)
		}
	}

	const firstState = uint64(0x101081000)
	value, err := s.load(firstState+0x290, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap protected-call first state: %w", err)
	}
	s.x[8] = firstState
	s.x[9] = value
	const secondState = uint64(0x101075000)
	value, err = s.load(secondState+0xb40, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap protected-call second state: %w", err)
	}
	s.x[10] = secondState
	s.x[11] = value
	s.x[12] = s.sp + 8
	s.x[9] = s.x[9] - s.x[12] + s.x[11]
	s.x[11] = 0x0a5cd4cf
	s.x[12] = 0x8c2fb9df71512e5a
	s.x[14] = s.x[9]*s.x[11] + s.x[12]
	s.x[12] = 0x73d046208eaed1a6
	s.x[9] ^= s.x[12]
	s.x[16] = s.x[9] * s.x[11]
	if err := s.storeValue(firstState+0x290, s.x[14], 8); err != nil {
		return 0, fmt.Errorf("bootstrap protected-call first state output: %w", err)
	}
	if err := s.storeValue(secondState+0xb40, s.x[16], 8); err != nil {
		return 0, fmt.Errorf("bootstrap protected-call second state output: %w", err)
	}

	s.x[8] = uint64(uint32(s.x[0]) ^ 0xe732d262)
	s.x[9] = 0xba209ef7
	value, err = s.load(s.x[0]+0xc, 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap protected-call descriptor word: %w", err)
	}
	s.x[10] = uint64(uint32(value) - uint32(s.x[8])*uint32(s.x[9]))
	value, err = s.load(s.x[0], 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap protected-call first descriptor: %w", err)
	}
	s.x[8] = value
	value, err = s.load(s.x[0]+0x10, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap protected-call second descriptor: %w", err)
	}
	s.x[9] = value
	value, err = s.load(s.x[8]+4, 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap protected-call first length: %w", err)
	}
	s.x[15] = uint64(uint32(value))
	value, err = s.load(s.x[9]+4, 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap protected-call second length: %w", err)
	}
	s.x[12] = uint64(uint32(value))
	s.x[13] = 0x710846d9

	s.traceBlock(bootstrapProtectedDispatch)
	s.x[17] = 0x6cddb946
	s.x[1] = uint64(uint32(s.x[15]) + uint32(s.x[13]))
	s.compare32(uint32(s.x[1]), uint32(s.x[17]))
	s.x[2] = uint64(uint32(boolValue(s.condition("lt"))))
	s.x[13] = uint64(uint32(s.x[12]) + 0x710846d9)
	s.compare32(uint32(s.x[13]), uint32(s.x[17]))
	s.x[17] = uint64(uint32(boolValue(s.condition("lt"))))
	s.x[17] = uint64(uint32(s.x[2]) ^ uint32(s.x[17]))
	s.compare32(uint32(s.x[1]), uint32(s.x[13]))
	s.x[13] = uint64(uint32(boolValue(s.condition("gt"))))
	s.compare32(uint32(s.x[17]), 0)
	if s.condition("ne") {
		s.x[17] = uint64(uint32(s.x[2]))
	} else {
		s.x[17] = uint64(uint32(s.x[13]))
	}
	s.x[13] = uint64((uint32(s.x[10]) ^ 0x3f39) * 0x8dd)
	s.x[1] = uint64(uint32(s.x[13]) ^ 0xcbb9)
	s.x[17] = uint64(uint32(s.x[17]) * uint32(s.x[1]))
	s.x[14] = uint64(uint32(s.x[14]) ^ 0xa6)
	s.x[14] = uint64(uint32(s.x[14]) ^ uint32(s.x[16]))
	s.x[11] = uint64(uint32(s.x[14]) * 0x0a5cd4cf)
	const firstByteTable = uint64(0x100fe9390)
	value, err = s.load(firstByteTable+uint64(byte(s.x[11])), 1)
	if err != nil {
		return 0, fmt.Errorf("bootstrap protected-call first byte table: %w", err)
	}
	s.x[14] = uint64(byte(value) ^ 0x28)
	const secondByteTable = uint64(0x100de2af0)
	value, err = s.load(secondByteTable+uint64(uint32(s.x[14])), 1)
	if err != nil {
		return 0, fmt.Errorf("bootstrap protected-call second byte table: %w", err)
	}
	s.x[14] = uint64(byte(value))
	s.x[11] = uint64(((uint32(s.x[11]) ^ uint32(s.x[14])) & 0xff) + 0x1b6)
	const firstDispatchTable = uint64(0x10105b6f0)
	value, err = s.load(firstDispatchTable+uint64(uint32(s.x[11]))*8, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap protected-call first dispatch: %w", err)
	}
	s.x[11] = value
	s.x[14] = uint64(uint32(s.x[17]) ^ uint32(s.x[10]))
	s.x[10] = 0xffffffff98245355
	s.x[10] += s.x[11]
	next, err := loadBootstrapPointer(s, s.x[10], uint32(s.x[14]))
	if err != nil {
		return 0, fmt.Errorf("bootstrap protected-call target: %w", err)
	}
	s.x[14] = next
	s.x[11] = 0x5eec043e
	s.x[17] = 0x63e8
	return next, nil
}

func runBootstrapIndexedDispatch(s *state) (uint64, error) {
	s.traceBlock(bootstrapIndexedDispatch)
	s.x[10] = uint64(uint32(s.x[8]) + 0x2a79)
	s.x[19] = uint64(uint32(s.x[10]) | 0xc00)
	s.x[10] = uint64((uint32(s.x[8]) + 0x3bf7a813) & 0xc40857cf)
	s.x[11] = uint64((uint32(s.x[10]) ^ 0xa9a335c6) + uint32(s.x[9]))
	s.x[12] = s.sp + 0x2d0
	word, err := s.load(s.x[12]+uint64(uint32(s.x[11]))*4, 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap indexed-dispatch lookup: %w", err)
	}
	s.x[11] = uint64(uint32(word))
	s.x[12] = uint64(uint32(s.x[19]) ^ 0x224e2e9f)
	s.compare32(uint32(s.x[11]), uint32(s.x[12]))
	s.x[11] = uint64(uint32(boolValue(s.condition("eq"))) * 0xc36e)
	s.x[8] = uint64(uint32(s.x[11]) ^ uint32(s.x[8]))
	next, err := loadBootstrapPointer(s, s.x[25], uint32(s.x[8]))
	if err != nil {
		return 0, fmt.Errorf("bootstrap indexed-dispatch target: %w", err)
	}
	s.x[8] = next
	s.x[17] = 0xe52e
	return next, nil
}

func runBootstrapCountdownDispatch(s *state) (uint64, error) {
	s.traceBlock(bootstrapCountdownStart)
	s.x[19] = uint64(uint32(s.x[10]) ^ 0xdaaa)
	s.x[8] = uint64(uint32(s.x[10]) ^ 0x25)
	s.x[12] = uint64((uint32(s.x[10]) ^ 0xa9a335c8) + uint32(s.x[9]))
	s.traceBlock(bootstrapCountdownDispatch)
	s.x[9] = uint64(uint32(s.x[9]) - 1)
	s.x[13] = uint64(uint32(s.x[8]) + 0x7fffa8d4)
	s.compare32(uint32(s.x[12]), uint32(s.x[13]))
	s.x[12] = uint64(uint32(boolValue(s.condition("hs"))))
	s.x[11] = uint64(uint32(s.x[12]) * uint32(s.x[11]))
	s.x[10] = uint64(uint32(s.x[11]) ^ uint32(s.x[10]))
	next, err := loadBootstrapPointer(s, s.x[25], uint32(s.x[10]))
	if err != nil {
		return 0, fmt.Errorf("bootstrap countdown-dispatch target: %w", err)
	}
	s.x[10] = next
	s.x[17] = 0xe52e
	return next, nil
}

func runBootstrapSmallFrameReturn(s *state) (uint64, error) {
	s.traceBlock(bootstrapSmallFrameReturn)
	s.sp += 0x10
	return s.x[30], nil
}

func runBootstrapSavedFrameReturn(s *state) (uint64, error) {
	s.traceBlock(bootstrapSavedFrameReturn)
	if err := s.storeValue(s.x[0]+8, uint64(uint32(s.x[11])), 4); err != nil {
		return 0, fmt.Errorf("bootstrap saved-frame result: %w", err)
	}
	for _, pair := range [...]struct {
		first, second int
		offset        uint64
	}{{29, 30, 0x60}, {20, 19, 0x50}, {22, 21, 0x40}, {24, 23, 0x30}, {26, 25, 0x20}, {28, 27, 0x10}} {
		first, err := s.load(s.sp+pair.offset, 8)
		if err != nil {
			return 0, fmt.Errorf("bootstrap saved-frame register %d: %w", pair.first, err)
		}
		second, err := s.load(s.sp+pair.offset+8, 8)
		if err != nil {
			return 0, fmt.Errorf("bootstrap saved-frame register %d: %w", pair.second, err)
		}
		s.x[pair.first] = first
		s.x[pair.second] = second
	}
	s.sp += 0x70
	return s.x[30], nil
}

func runBootstrapPointerCycle(s *state) (uint64, error) {
	pc := bootstrapPointerCycleStart
	const maximumBlocks = 8192
	for count := 0; count < maximumBlocks; count++ {
		var next uint64
		var err error
		switch pc {
		case bootstrapPointerCycleStart:
			next, err = runBootstrapPointerCycleStart(s)
		case bootstrapPointerCycleCheck:
			next, err = runBootstrapPointerCycleCheck(s)
		case bootstrapPointerCycleStep:
			next, err = runBootstrapPointerCycleStep(s)
		case bootstrapPointerCycleNext:
			next, err = runBootstrapPointerCycleNext(s)
		default:
			return pc, nil
		}
		if err != nil {
			return pc, err
		}
		pc = next
	}
	return 0, fmt.Errorf("bootstrap pointer cycle exceeded %d blocks", maximumBlocks)
}

func runBootstrapPointerCycleStart(s *state) (uint64, error) {
	s.traceBlock(bootstrapPointerCycleStart)
	s.x[1] = uint64(uint32(s.x[21]) + 0xbca)
	word, err := s.load(s.x[10], 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap pointer-cycle word: %w", err)
	}
	s.x[24] = uint64(uint32(word))
	next, err := loadBootstrapPointer(s, s.x[22], uint32(s.x[21]))
	if err != nil {
		return 0, fmt.Errorf("bootstrap pointer-cycle start dispatch: %w", err)
	}
	s.x[21] = next
	s.traceBlock(bootstrapPointerCycleTail)
	s.x[17] = 0x6dd3
	return next, nil
}

func runBootstrapPointerCycleCheck(s *state) (uint64, error) {
	s.traceBlock(bootstrapPointerCycleCheck)
	s.x[21] = uint64(uint32(s.x[1]) ^ uint32(s.x[12]))
	s.x[17] = uint64(uint32(s.x[1]) ^ uint32(s.x[13]))
	s.compare32(uint32(s.x[24]), uint32(s.x[17]))
	dispatchIndex := uint32(boolValue(s.condition("lt")))*uint32(s.x[14]) ^ uint32(s.x[1])
	next, err := loadBootstrapPointer(s, s.x[22], dispatchIndex)
	if err != nil {
		return 0, fmt.Errorf("bootstrap pointer-cycle check dispatch: %w", err)
	}
	s.x[1] = next
	s.x[17] = 0x99f1
	return next, nil
}

func runBootstrapPointerCycleStep(s *state) (uint64, error) {
	s.traceBlock(bootstrapPointerCycleStep)
	s.compare32(uint32(s.x[24]), 0x21)
	selected := uint32(boolValue(s.condition("eq")))
	s.x[1] = uint64(uint32(s.x[21]) + 0x9a48)
	s.x[2] = 0x3887
	s.x[24] = uint64(uint32(s.x[1]) ^ uint32(s.x[2]))
	dispatchIndex := uint32(s.x[24])*selected ^ uint32(s.x[21])
	next, err := loadBootstrapPointer(s, s.x[22], dispatchIndex)
	if err != nil {
		return 0, fmt.Errorf("bootstrap pointer-cycle step dispatch: %w", err)
	}
	s.x[17] = 0x1a09
	s.x[24] = next
	return next, nil
}

func runBootstrapPointerCycleNext(s *state) (uint64, error) {
	s.traceBlock(bootstrapPointerCycleNext)
	s.compare64(s.x[8], s.x[11])
	firstBound := uint32(boolValue(s.condition("lo")))
	s.x[24] = uint64(uint32(s.x[1]) + uint32(s.x[5]))
	s.x[2] = uint64(uint32(s.x[1]) + uint32(s.x[6]))
	s.x[21] = uint64(uint32(s.x[2]) & uint32(s.x[7]))
	word, err := s.load(s.x[10]+4, 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap pointer-cycle increment: %w", err)
	}
	s.x[2] = uint64(uint32(word))
	s.x[10] += s.x[2]
	nextCounter := s.x[9] + 1
	s.compare64(nextCounter, s.x[28])
	selector := firstBound ^ uint32(boolValue(s.condition("hi")))
	s.x[9] += s.x[27]
	s.compare64(s.x[9], s.x[8])
	secondBound := uint32(boolValue(s.condition("lo")))
	s.compare32(selector, 0)
	dispatchBit := secondBound
	if s.condition("ne") {
		dispatchBit = firstBound
	}
	s.x[17] = uint64(uint32(s.x[21]) + uint32(s.x[30]))
	dispatchIndex := dispatchBit*uint32(s.x[17]) ^ uint32(s.x[1])
	next, err := loadBootstrapPointer(s, s.x[22], dispatchIndex)
	if err != nil {
		return 0, fmt.Errorf("bootstrap pointer-cycle next dispatch: %w", err)
	}
	s.x[1] = next
	s.x[9] = nextCounter
	s.x[17] = 0x1329
	return next, nil
}

func runBootstrapByteClassify(s *state) (uint64, error) {
	const maximumBytes = 4096
	for index := 0; index < maximumBytes; index++ {
		next, err := runBootstrapByteClassifyStep(s)
		if err != nil || next != bootstrapByteClassifyPass {
			return next, err
		}
	}
	return 0, fmt.Errorf("bootstrap byte classify exceeded %d bytes", maximumBytes)
}

func runBootstrapByteClassifyStep(s *state) (uint64, error) {
	s.traceBlock(bootstrapByteClassifyPass)
	s.x[13] = uint64(uint32(s.x[12]))
	s.x[17] = shiftLeft(s.x[13], int(s.x[15]), 64)
	s.x[17] &= s.x[16]
	s.x[0] = s.x[13] + s.x[11]
	s.x[13] = s.x[0] - s.x[17]
	s.x[17] = s.x[13] ^ s.x[11]
	value, err := s.load(s.x[23]+s.x[17], 1)
	if err != nil {
		return 0, fmt.Errorf("bootstrap byte-classify input: %w", err)
	}
	s.x[17] = uint64(byte(value))
	s.compare32(uint32(s.x[17]), 0x1b)
	selected := uint32(boolValue(s.condition("ne")))
	s.x[12] = uint64(uint32(s.x[12]) + 1)
	dispatchIndex := ((selected << 2) | selected) ^ uint32(s.x[14])
	next, err := loadBootstrapPointer(s, s.x[21], dispatchIndex)
	if err != nil {
		return 0, fmt.Errorf("bootstrap byte-classify dispatch: %w", err)
	}
	s.x[1] = next
	s.x[17] = 0xd618
	return next, nil
}
