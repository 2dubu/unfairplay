package recovery

import "fmt"

func runAccountSecondFoldCircuit(s *state) (uint64, error) {
	const maximumRounds = 4096
	for round := 0; round < maximumRounds; round++ {
		next, err := runAccountSecondFoldRoundCircuit(s)
		if err != nil || next != accountSecondFoldRound {
			return next, err
		}
	}
	return 0, fmt.Errorf("account second fold exceeded %d rounds", maximumRounds)
}

func runAccountSecondFoldRoundCircuit(s *state) (uint64, error) {
	s.traceBlock(accountSecondFoldRound)
	s.x[16], s.x[28], s.x[1] = s.x[28], s.x[15], s.x[11]
	s.x[11] = uint64(uint32(rotateRight(s.x[11], 27, 32)))
	s.x[13] = uint64((uint32(s.x[11]) ^ 0x4e0b5915) + 0x8a2c71a6)
	s.x[15] = uint64((uint32(s.x[11]) ^ 0x22a6723d) + 0xe6815a8e)
	s.x[13] = uint64(uint32(s.x[13]) ^ uint32(s.x[15]))
	s.x[11] = uint64((uint32(s.x[11]) ^ 0x0f549c89) + 0xcb73b43a)
	s.x[11] = uint64(uint32(s.x[13]) ^ uint32(s.x[11]))
	s.x[13] = uint64(uint32(s.x[28]) ^ 0xad98f639)
	s.x[15] = uint64(uint32(s.x[16]) ^ uint32(s.x[22]))
	s.x[13] = uint64(uint32(s.x[15]) & uint32(s.x[13]))
	s.x[15] = uint64(uint32(s.x[16]) & 0x5978c47b)
	s.x[13] = uint64(uint32(s.x[13]) ^ uint32(s.x[15]))
	s.x[15] = uint64(uint32(s.x[6]) ^ 0x641d5cfa)
	s.x[17] = uint64(uint32(s.x[13]) ^ 0xaf9fffbd)
	s.x[15] = uint64(uint32(s.x[17]) & uint32(s.x[15]))
	s.x[13] = uint64(uint32(s.x[13]) & 0x6f029147)
	s.x[17] = uint64(uint32(s.x[28]) ^ 0x8537880a)
	s.x[2] = uint64(uint32(s.x[16]) ^ uint32(s.x[10]))
	s.x[17] = uint64(uint32(s.x[2]) & uint32(s.x[17]))
	s.x[2] = uint64(uint32(s.x[16]) & 0x8e2845b7)
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[2]))
	s.x[13] = uint64(uint32(s.x[13]) ^ 0xfb5b407c)
	s.x[13] = uint64(uint32(s.x[13]) ^ uint32(s.x[15]))
	s.x[15] = uint64(uint32(s.x[17]) ^ 0xf5f7ba4a)
	s.x[13] = uint64(uint32(s.x[13]) & uint32(s.x[15]))
	s.x[15] = uint64(uint32(s.x[17]) & 0x2ba62e86)
	s.x[13] = uint64(uint32(s.x[13]) ^ uint32(s.x[15]))
	s.x[15] = uint64((uint32(s.x[13]) ^ 0x8ac8bda8) + 0x8ad78b4c)
	s.x[17] = uint64((uint32(s.x[13]) ^ 0xf48d333a) + 0xf49205da)
	s.x[15] = uint64(uint32(s.x[15]) ^ uint32(s.x[17]))
	s.x[13] = uint64((uint32(s.x[13]) ^ 0xa01c5b6f) + 0xa0036d8d)
	s.x[13] = uint64(uint32(s.x[15]) ^ uint32(s.x[13]))
	s.x[15] = uint64((uint32(s.x[19]) ^ 0x2c4d2cc1) + 0xd3b2d33f)
	s.x[17] = uint64((uint32(s.x[19]) ^ 0xd51d1e6a) + 0x2ae2e196)
	s.x[15] = uint64(uint32(s.x[15]) ^ uint32(s.x[17]))
	s.x[17] = uint64((uint32(s.x[19]) ^ 0xf24fff16) + 0x0db000ea)
	s.x[15] = uint64(uint32(s.x[15]) ^ uint32(s.x[17]))
	address := s.x[14] + (s.x[21] << 2)
	value, err := s.load(s.x[8]+address, 4)
	if err != nil {
		return 0, fmt.Errorf("account second fold input: %w", err)
	}
	s.x[17] = uint64(uint32(value))

	s.traceBlock(accountSecondFoldTail)
	s.x[9] = 0xdd0677b5
	s.x[2] = uint64((uint32(s.x[17]) ^ uint32(s.x[9])) + 0x22f9884b)
	s.x[5] = uint64((uint32(s.x[17]) ^ 0x4efd1e4b) + 0xb102e1b5)
	s.x[2] = uint64(uint32(s.x[2]) ^ uint32(s.x[5]))
	s.x[17] = uint64((uint32(s.x[17]) ^ 0x5d8e123b) + 0xa271edc5)
	s.x[17] = uint64(uint32(s.x[2]) ^ uint32(s.x[17]))
	s.x[11] = uint64(uint32(s.x[15]) + uint32(s.x[11]))
	s.x[11] = uint64(uint32(s.x[11]) + uint32(s.x[17]))
	s.x[11] = uint64(uint32(s.x[11]) + uint32(s.x[13]))
	s.x[11] = uint64(uint32(s.x[11]) + 0x9d99a70c)
	s.x[13] = uint64(uint32(s.x[11]) ^ 0x4575db32)
	s.x[15] = uint64((uint32(s.x[11]) ^ 0x491a2872) + 0x0debf553)
	s.x[17] = uint64(uint32(s.x[11]) ^ 0x8a9e5ee0)
	s.x[2] = uint64(uint32(s.x[11]) ^ 0x3dff8f7f)
	s.x[13] = uint64(uint32(s.x[13]) + 0x01840613)
	s.x[11] = uint64(uint32(s.x[11]) ^ uint32(s.x[13]))
	s.x[13] = uint64(uint32(s.x[17]) + 0xce6f83c1)
	s.x[13] = uint64(uint32(s.x[15]) ^ uint32(s.x[13]))
	s.x[11] = uint64(uint32(s.x[11]) ^ uint32(s.x[13]))
	s.x[13] = uint64(uint32(s.x[2]) + 0x790e5260)
	s.x[13] = uint64(uint32(s.x[13]) ^ 0xb011ef62)
	s.x[11] = uint64(uint32(s.x[11]) ^ uint32(s.x[13]))
	s.x[15] = uint64(uint32(s.x[30]) ^ uint32(rotateRight(s.x[6], 2, 32)))
	s.x[21]++
	s.x[9] = accountSecondFoldLimit
	s.compare64(s.x[21], s.x[9])
	gate := uint32(boolValue(s.condition("ne")))
	s.x[13] = uint64(gate << 5)
	s.x[13] = uint64(uint32(s.x[13]) - gate)
	s.x[13] = uint64(uint32(s.x[13]) ^ uint32(s.x[0]))
	next, err := loadBootstrapPointer(s, s.x[20], uint32(s.x[13]))
	if err != nil {
		return 0, fmt.Errorf("account second fold dispatch: %w", err)
	}
	s.x[6], s.x[19] = s.x[1], s.x[16]
	s.x[7] = accountSecondFoldLimit
	s.x[13] = next
	s.x[17] = 0x2583
	return next, nil
}
