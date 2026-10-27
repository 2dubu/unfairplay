package recovery

import "fmt"

const (
	bootstrapWordTableRound = uint64(0x10082c698)
	bootstrapWordTableTail  = uint64(0x10082c800)
)

func runBootstrapWordTable(s *state) (uint64, error) {
	const maximumRounds = 256
	for round := 0; round < maximumRounds; round++ {
		next, err := runBootstrapWordTableRound(s)
		if err != nil || next != bootstrapWordTableRound {
			return next, err
		}
	}
	return 0, fmt.Errorf("bootstrap word-table pass exceeded %d rounds", maximumRounds)
}

func runBootstrapWordTableRound(s *state) (uint64, error) {
	s.traceBlock(bootstrapWordTableRound)
	s.x[11] = uint64(uint32(s.x[13]) + 0xffffb95f)
	s.x[23] = uint64(uint32(s.x[11]) | 0x4023)
	s.x[11] = uint64(uint32(s.x[16]) - 0xf)
	s.x[14] = uint64(uint32(s.x[11]) ^ 0x48)
	s.x[11] = uint64(uint32(0x90) & uint32(s.x[11]<<1))
	s.x[14] = uint64(uint32(s.x[11]) + uint32(s.x[14]))
	s.x[11] = uint64(uint32(s.x[16]) * uint32(s.x[19]))
	s.x[14] = uint64(uint32(s.x[14]) * uint32(s.x[19]))
	s.x[17] = uint64(uint32(s.x[23]) ^ 0x74f4)
	s.x[17] = uint64(uint32(s.x[17]) + uint32(s.x[11]))
	s.x[0] = uint64(uint32(s.x[17]) & 0xffff)
	s.x[0] = uint64(uint32(s.x[0]) * uint32(s.x[21]))
	s.x[0] = uint64(uint32(s.x[0]) >> 23)
	s.x[17] = uint64(uint32(s.x[17]) - uint32(s.x[0])*uint32(s.x[22]))
	s.x[12] = 0xe1fc780f
	s.x[0] = uint64(uint32(s.x[14])) * uint64(uint32(s.x[12]))
	s.x[0] >>= 39
	s.x[14] = uint64(uint32(s.x[14]) - uint32(s.x[0])*uint32(s.x[22]))
	s.x[17] = (s.x[17] & 0xffff) << 2
	value, err := s.load(s.x[9]+s.x[17], 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap word-table first lookup: %w", err)
	}
	s.x[17] = uint64(uint32(value))
	s.x[0] = uint64(uint32(s.x[11]) + 0x271a)
	s.x[3] = uint64(uint32(s.x[0]) & 0xffff)
	s.x[3] = uint64(uint32(s.x[3]) * uint32(s.x[21]))
	s.x[3] = uint64(uint32(s.x[3]) >> 23)
	value, err = s.load(s.x[9]+uint64(uint32(s.x[14]))*4, 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap word-table second lookup: %w", err)
	}
	s.x[14] = uint64(uint32(value))
	s.x[0] = uint64(uint32(s.x[0]) - uint32(s.x[3])*uint32(s.x[22]))
	s.x[0] = (s.x[0] & 0xffff) << 2
	value, err = s.load(s.x[9]+s.x[0], 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap word-table third lookup: %w", err)
	}
	s.x[0] = uint64(uint32(value))
	s.x[3] = uint64(uint32(s.x[11]) + uint32(s.x[20]))
	s.x[4] = uint64(uint32(s.x[17]) >> 3)
	s.x[6] = uint64(uint32(s.x[3]) & 0xffff)
	s.x[6] = uint64(uint32(s.x[6]) * uint32(s.x[21]))
	s.x[6] = uint64(uint32(s.x[6]) >> 23)
	s.x[3] = uint64(uint32(s.x[3]) - uint32(s.x[6])*uint32(s.x[22]))
	s.x[14] = uint64(uint32(extractPair(s.x[17], s.x[14], 7, 32)))
	s.x[12] = 0x333dea8c
	s.x[6] = uint64(uint32(s.x[12]) & (uint32(s.x[17]) >> 2))
	s.x[4] = uint64(uint32(s.x[4]) - uint32(s.x[6]))
	s.x[12] = 0x599ef546
	s.x[4] = uint64(uint32(s.x[4]) + uint32(s.x[12]))
	s.x[6] = uint64(uint32(rotateRight(s.x[0], 17, 32)))
	s.x[14] = uint64(uint32(s.x[14]) ^ uint32(rotateRight(s.x[17], 18, 32)))
	s.x[14] = uint64(uint32(s.x[14]) ^ uint32(s.x[4]))
	s.x[17] = uint64(uint32(s.x[0]) >> 10)
	s.x[16] = uint64(uint32(s.x[16]) + 1)
	s.x[4] = uint64(uint32(0x0018307a) & (uint32(s.x[0]) >> 9))
	s.x[17] = uint64(uint32(s.x[17]) - uint32(s.x[4]))
	s.x[14] = uint64(uint32(s.x[14]) ^ uint32(s.x[12]))
	s.x[12] = 0x46cc183d
	s.x[17] = uint64(uint32(s.x[17]) + uint32(s.x[12]))
	s.x[0] = uint64(uint32(s.x[6]) ^ uint32(rotateRight(s.x[0], 19, 32)))
	s.x[17] = uint64(uint32(s.x[0]) ^ uint32(s.x[17]))
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[12]))
	s.x[0] = (s.x[3] & 0xffff) << 2
	s.x[14] = uint64(uint32(s.x[17]) + uint32(s.x[14]))
	s.x[17] = uint64(uint32(s.x[11]) + 0x244f)
	s.x[3] = uint64(uint32(s.x[17]) & 0xffff)
	s.x[3] = uint64(uint32(s.x[3]) * uint32(s.x[21]))
	s.x[3] = uint64(uint32(s.x[3]) >> 23)
	if err := s.storeValue(s.x[9]+s.x[0], uint64(uint32(s.x[14])), 4); err != nil {
		return 0, fmt.Errorf("bootstrap word-table first output: %w", err)
	}
	s.x[17] = uint64(uint32(s.x[17]) - uint32(s.x[3])*uint32(s.x[22]))
	s.x[17] = (s.x[17] & 0xffff) << 2
	value, err = s.load(s.x[9]+s.x[17], 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap word-table fourth lookup: %w", err)
	}
	s.x[17] = uint64(uint32(value))
	s.x[11] = uint64(uint32(s.x[11]) + 0x1f48)
	s.x[14] = uint64(uint32(s.x[14]) + uint32(s.x[17]))
	s.x[17] = uint64(uint32(s.x[11]) & 0xffff)
	s.x[17] = uint64(uint32(s.x[17]) * uint32(s.x[21]))
	s.x[17] = uint64(uint32(s.x[17]) >> 23)
	s.x[11] = uint64(uint32(s.x[11]) - uint32(s.x[17])*uint32(s.x[22]))
	if err := s.storeValue(s.x[9]+s.x[0], uint64(uint32(s.x[14])), 4); err != nil {
		return 0, fmt.Errorf("bootstrap word-table second output: %w", err)
	}
	s.x[11] = (s.x[11] & 0xffff) << 2
	value, err = s.load(s.x[9]+s.x[11], 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap word-table fifth lookup: %w", err)
	}
	s.x[11] = uint64(uint32(value))
	s.x[11] = uint64(uint32(s.x[14]) + uint32(s.x[11]))
	if err := s.storeValue(s.x[9]+s.x[0], uint64(uint32(s.x[11])), 4); err != nil {
		return 0, fmt.Errorf("bootstrap word-table final output: %w", err)
	}
	s.compare32(uint32(s.x[16]), 0x40)
	s.x[11] = uint64(uint32(boolValue(s.condition("eq"))))

	s.traceBlock(bootstrapWordTableTail)
	s.x[12] = 0xfdd8
	s.x[11] = uint64(uint32(s.x[11]) * uint32(s.x[12]))
	s.x[11] = uint64(uint32(s.x[11]) ^ uint32(s.x[13]))
	next, err := loadBootstrapPointer(s, s.x[30], uint32(s.x[11]))
	if err != nil {
		return 0, fmt.Errorf("bootstrap word-table dispatch: %w", err)
	}
	s.x[11] = next
	s.x[17] = 0xf00c
	return next, nil
}
