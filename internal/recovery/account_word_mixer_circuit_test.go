package recovery

import "fmt"

// runAccountEntryWordMixerStepCircuit is the mechanically simplified native
// handler retained only as a differential oracle for the semantic SHA-1
// schedule implementation.
func runAccountEntryWordMixerStepCircuit(s *state) (uint64, error) {
	s.traceBlock(accountEntryWordMixerPass)
	s.x[9] = uint64(uint32(s.x[15]) + 0xd458908a)
	s.x[11] = s.x[20] << 2
	addresses := [...]uint64{
		0xffffffff5160f1c8,
		0xffffffff5160f1b4,
		0xffffffff5160f19c,
		0xffffffff5160f194,
	}
	registers := [...]int{12, 13, 16, 17}
	for index, offset := range addresses {
		value, err := s.load(s.x[10]+s.x[11]+offset, 4)
		if err != nil {
			return 0, fmt.Errorf("account-entry word mixer input %d: %w", index, err)
		}
		s.x[registers[index]] = uint64(uint32(value))
	}
	s.x[12] = uint64(uint32(s.x[13]) ^ uint32(s.x[12]))
	s.x[13] = uint64(uint32(s.x[16]) ^ uint32(s.x[17]))
	s.x[12] = uint64(uint32(s.x[12]) ^ uint32(s.x[13]))
	s.x[12] = uint64(uint32(s.x[9]) ^ uint32(rotateRight(s.x[12], 31, 32)))
	s.x[13] = 0xce752fa0
	s.x[12] = uint64(uint32(s.x[12]) ^ uint32(s.x[13]))
	s.x[11] += s.x[6]
	if err := s.storeValue(s.x[10]+s.x[11], uint64(uint32(s.x[12])), 4); err != nil {
		return 0, fmt.Errorf("account-entry word mixer output: %w", err)
	}
	s.x[20]++
	s.x[11] = s.x[20] - s.x[15]
	s.x[12] = s.x[15] - s.x[20]
	s.x[11] |= s.x[12]
	s.x[11] >>= 63
	s.x[11] = uint64(uint32(s.x[11]) * 0x72)
	s.x[12] = uint64(uint32(s.x[15]) + 0xd4589955)
	dispatchIndex := uint32(s.x[11]) ^ uint32(s.x[12])
	next, err := loadBootstrapPointer(s, s.x[21], dispatchIndex)
	if err != nil {
		return 0, fmt.Errorf("account-entry word mixer dispatch: %w", err)
	}
	limit, err := s.load(s.sp+0x7c, 4)
	if err != nil {
		return 0, fmt.Errorf("account-entry word mixer limit: %w", err)
	}
	s.x[4] = next
	s.x[7] = s.x[0]
	s.x[11] = s.x[3]
	s.x[12] = s.x[26]
	s.x[13] = s.x[22]
	s.x[16] = 0x2ba7c38b
	s.x[17] = 0x3601
	s.x[19] = uint64(uint32(limit))
	return next, nil
}
