package recovery

import "fmt"

const (
	sharedPermutationFront = uint64(0x100439e54)
	sharedPermutationTail  = uint64(0x10043a000)
)

func runSharedPermutationFront(s *state) (uint64, error) {
	s.traceBlock(sharedPermutationFront)
	s.x[17] = 0x7c16
	s.x[1] = uint64(uint32(s.x[2]) ^ uint32(s.x[17]))
	s.x[17] = uint64(uint32(s.x[11]) - 0x10)
	s.x[22] = 0x4224b4ad
	s.x[22] = uint64(uint32(s.x[1]) + uint32(s.x[22]))
	s.x[22] = uint64(uint32(s.x[22]) ^ uint32(s.x[17]))
	s.x[22] = uint64(uint32(s.x[22]) + 0xc1)
	s.x[23] = 0xd9cff57a
	s.x[23] = uint64(uint32(s.x[17]) ^ uint32(s.x[23]))
	s.x[23] = uint64(uint32(s.x[23]) + 0xd8)
	s.x[24] = uint64(uint32(s.x[17]) ^ uint32(s.x[6]))
	s.x[24] = uint64(uint32(s.x[24]) + 4)
	s.x[25] = 0xeb
	s.x[25] = uint64(uint32(s.x[17]) ^ uint32(s.x[25]))
	s.x[25] = uint64(uint32(s.x[25]) + 0x47)
	s.x[17] = uint64(uint32(s.x[23]) ^ uint32(s.x[17]))
	s.x[23] = uint64(uint32(s.x[24]) ^ uint32(s.x[25]))
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[23]))
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[22]))
	s.x[17] &= 0xff
	s.x[17] ^= s.x[7]
	value, err := s.load(s.sp+0x58, 8)
	if err != nil {
		return 0, fmt.Errorf("shared permutation table: %w", err)
	}
	s.x[28] = value
	value, err = s.load(s.x[28]+s.x[17], 1)
	if err != nil {
		return 0, fmt.Errorf("shared permutation selector byte: %w", err)
	}
	s.x[17] = uint64(byte(value)) & 7
	s.x[22] = 0x59787bf3
	s.x[22] = uint64(uint32(s.x[17]) ^ uint32(s.x[22]))
	s.x[23] = 0xa6878479
	s.x[22] = uint64(uint32(s.x[22]) + uint32(s.x[23]))
	s.x[22] = uint64(uint32(s.x[17]) ^ uint32(s.x[22]))
	s.x[23] = 0xe7a2820c
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[23]))
	s.x[23] = 0x185d7d88
	s.x[17] = uint64(uint32(s.x[17]) + uint32(s.x[23]))
	s.x[23] = 0x74
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[23]))
	s.x[17] = uint64(uint32(s.x[22]) ^ uint32(s.x[17]))
	s.x[22] = 0x569caf38
	s.x[22] = uint64(uint32(s.x[17]) + uint32(s.x[22]))
	s.x[23] = 0x69a8f24b
	s.x[23] = uint64(uint32(s.x[22]) ^ uint32(s.x[23]))
	s.x[24] = uint64(uint32(s.x[22]) & 0x3f)
	s.x[25] = 0x296a503f
	s.x[25] = uint64(uint32(s.x[22]) ^ uint32(s.x[25]))
	s.x[25] = uint64(uint32(s.x[25]) & uint32(s.x[24]<<1))
	s.x[24] = uint64(uint32(s.x[25]) ^ uint32(s.x[24]))
	s.x[25] = uint64(uint32(s.x[23]) << 1)
	s.x[27] = 0x819746e8
	s.x[27] = uint64(uint32(s.x[25]) ^ uint32(s.x[27]))
	s.x[23] = uint64(uint32(s.x[27]) & uint32(s.x[23]))
	s.x[27] = 0x40cba374
	s.x[25] = uint64(uint32(s.x[25]) & uint32(s.x[27]))
	s.x[23] = uint64(uint32(s.x[23]) ^ uint32(s.x[25]))
	s.x[25] = 0x4048a114
	s.x[25] = uint64(uint32(s.x[23]) ^ uint32(s.x[25]))
	s.x[23] = uint64(uint32(s.x[23]) ^ 0x60)
	s.x[23] = uint64(uint32(s.x[23]) & uint32(s.x[24]<<2))
	s.x[23] = uint64(uint32(s.x[23]) ^ uint32(s.x[24]))
	s.x[24] = uint64(uint32(s.x[25]) << 2)
	s.x[27] = 0x032e8dd0
	s.x[27] = uint64(uint32(s.x[24]) ^ uint32(s.x[27]))
	s.x[25] = uint64(uint32(s.x[27]) & uint32(s.x[25]))
	s.x[27] = 0xc0cba370
	s.x[24] = uint64(uint32(s.x[24]) & uint32(s.x[27]))
	s.x[24] = uint64(uint32(s.x[25]) ^ uint32(s.x[24]))
	s.x[25] = 0xc0c12224
	s.x[25] = uint64(uint32(s.x[24]) ^ uint32(s.x[25]))
	s.x[27] = 0x000a8150
	s.x[24] = uint64(uint32(s.x[24]) ^ uint32(s.x[27]))
	s.x[24] = uint64(uint32(s.x[24]) & uint32(s.x[23]<<4))
	s.x[23] = uint64(uint32(s.x[24]) ^ uint32(s.x[23]))
	s.x[24] = uint64(uint32(s.x[25]) << 4)
	s.x[27] = 0x0cba3740
	s.x[27] = uint64(uint32(s.x[24]) ^ uint32(s.x[27]))
	s.x[25] = uint64(uint32(s.x[27]) & uint32(s.x[25]))
	s.x[27] = 0xc0cba340
	s.x[24] = uint64(uint32(s.x[24]) & uint32(s.x[27]))
	s.x[24] = uint64(uint32(s.x[25]) ^ uint32(s.x[24]))
	s.x[25] = 0xc0418034
	s.x[25] = uint64(uint32(s.x[24]) ^ uint32(s.x[25]))
	s.x[27] = 0x008a2300
	s.x[24] = uint64(uint32(s.x[24]) ^ uint32(s.x[27]))
	s.x[24] = uint64(uint32(s.x[24]) & uint32(s.x[23]<<8))
	s.x[27] = 0xc0cba374
	s.x[23] = uint64(uint32(s.x[23]) ^ uint32(s.x[27]))
	s.x[23] = uint64(uint32(s.x[23]) ^ uint32(s.x[24]))
	s.x[24] = uint64(uint32(s.x[25]) << 8)
	return sharedPermutationTail, nil
}

func runSharedPermutationTail(s *state) (uint64, error) {
	s.traceBlock(sharedPermutationTail)
	s.x[27] = 0x4ba30000
	s.x[27] = uint64(uint32(s.x[24]) ^ uint32(s.x[27]))
	s.x[25] = uint64(uint32(s.x[27]) & uint32(s.x[25]))
	s.x[24] = uint64(uint32(s.x[24]) & uint32(s.x[26]))
	s.x[27] = 0x00480000
	s.x[24] = uint64(uint32(s.x[24]) ^ uint32(s.x[27]))
	s.x[24] = uint64(uint32(s.x[24]) ^ uint32(s.x[25]))
	s.x[25] = uint64(uint32(s.x[23]) << 16)
	s.x[27] = 0x23740000
	s.x[27] = uint64(uint32(s.x[25]) ^ uint32(s.x[27]))
	s.x[24] = uint64(uint32(s.x[27]) & uint32(s.x[24]))
	s.x[25] = uint64(uint32(s.x[25]) & uint32(s.x[26]))
	s.x[23] = uint64(uint32(s.x[25]) ^ uint32(s.x[23]))
	s.x[23] = uint64(uint32(s.x[23]) ^ uint32(s.x[24]))
	s.x[22] = uint64(uint32(s.x[22]) ^ uint32(s.x[23]<<1))
	s.x[23] = 0x5854c196
	s.x[22] = uint64(uint32(s.x[22]) ^ uint32(s.x[23]))
	s.x[23] = 0xbb1a751e
	s.x[23] = uint64(uint32(s.x[22]) & uint32(s.x[23]))
	s.x[24] = 0x30005400
	s.x[23] = uint64(uint32(s.x[23]) ^ uint32(s.x[24]))
	s.x[24] = 0x44e58ae1
	s.x[22] = uint64(uint32(s.x[22]) & uint32(s.x[24]))
	s.x[24] = 0x40208241
	s.x[22] = uint64(uint32(s.x[22]) ^ uint32(s.x[24]))
	s.x[20] = uint64(uint32(s.x[20]) ^ uint32(s.x[8]))
	s.x[22] = uint64(uint32(shiftRight(s.x[20], int(s.x[22]), 32)))
	s.x[22] = uint64(uint32(shiftRight(s.x[22], int(s.x[23]), 32)))
	s.x[23] = 0x6a53f896
	s.x[17] = uint64(uint32(s.x[23]) - uint32(s.x[17]))
	s.x[23] = 0x721158a0
	s.x[23] = uint64(uint32(s.x[17]) ^ uint32(s.x[23]))
	s.x[24] = 0x213a5b37
	s.x[24] = uint64(uint32(s.x[17]) ^ uint32(s.x[24]))
	s.x[25] = 0xb4965dce
	s.x[24] = uint64(uint32(s.x[24]) + uint32(s.x[25]))
	s.x[25] = 0xcf97057d
	s.x[25] = uint64(uint32(s.x[17]) ^ uint32(s.x[25]))
	s.x[27] = 0x5a3b0388
	s.x[25] = uint64(uint32(s.x[25]) + uint32(s.x[27]))
	s.x[24] = uint64(uint32(s.x[24]) ^ uint32(s.x[25]))
	s.x[25] = 0xe7bd5e5b
	s.x[23] = uint64(uint32(s.x[23]) + uint32(s.x[25]))
	s.x[23] = uint64(uint32(s.x[17]) ^ uint32(s.x[23]))
	s.x[23] = uint64(uint32(s.x[23]) ^ uint32(s.x[24]))
	s.x[24] = 0xf6efffef
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[24]))
	s.x[24] = 0x6343f916
	s.x[17] = uint64(uint32(s.x[17]) + uint32(s.x[24]))
	s.x[24] = 0xb3d9cd29
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[24]))
	s.x[17] = uint64(uint32(s.x[23]) ^ uint32(s.x[17]))
	s.x[23] = 0x1152ddeb
	s.x[23] = uint64(uint32(s.x[17]) & uint32(s.x[23]))
	s.x[24] = 0x11021428
	s.x[23] = uint64(uint32(s.x[23]) ^ uint32(s.x[24]))
	s.x[24] = 0xeead2214
	s.x[17] = uint64(uint32(s.x[17]) & uint32(s.x[24]))
	s.x[24] = 0xc8882004
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[24]))
	s.x[17] = uint64(uint32(shiftLeft(s.x[20], int(s.x[17]), 32)))
	s.x[17] = uint64(uint32(shiftLeft(s.x[17], int(s.x[23]), 32)))
	s.x[17] = uint64(uint32(s.x[22]) ^ uint32(s.x[17]))
	s.x[20] = uint64(uint32(s.x[17]) ^ uint32(s.x[8]))
	if err := s.storeValue(s.x[28]+s.x[0], uint64(byte(s.x[20])), 1); err != nil {
		return 0, fmt.Errorf("shared permutation first output: %w", err)
	}
	s.x[20] = uint64(uint32(s.x[11]) ^ 0x65bd27ab)
	s.x[22] = uint64(uint32(s.x[11]) & uint32(s.x[5]))
	s.x[23] = uint64(uint32(s.x[11]) ^ 0x319d1eaf)
	s.x[23] = uint64(uint32(s.x[23]) & uint32(s.x[22]<<1))
	s.x[22] = uint64(uint32(s.x[23]) ^ uint32(s.x[22]))
	s.x[23] = uint64(0xbf42760a ^ uint32(s.x[20]<<1))
	s.x[20] = uint64(uint32(s.x[23]) & uint32(s.x[20]))
	s.x[23] = uint64(uint32(s.x[20]) & uint32(s.x[22]<<2))
	s.x[22] = uint64(uint32(s.x[23]) ^ uint32(s.x[22]))
	s.x[20] = uint64(uint32(s.x[20]) & uint32(s.x[20]<<2))
	s.x[20] = uint64(uint32(s.x[20]) & uint32(s.x[22]<<4))
	s.x[20] = uint64(uint32(s.x[20]) ^ uint32(s.x[22]))
	s.x[20] = uint64(uint32(s.x[11]) ^ uint32(s.x[20]<<1))
	s.x[20] &= 0xff
	s.x[20] ^= s.x[6]
	value, err := s.load(s.x[28]+s.x[20], 1)
	if err != nil {
		return 0, fmt.Errorf("shared permutation first table byte: %w", err)
	}
	s.x[20] = uint64(byte(value))
	s.x[20] = uint64(uint32(s.x[20]) ^ uint32(s.x[8]))
	s.x[20] = uint64(uint32(s.x[20]) << 16)
	s.x[17] = uint64(uint32(s.x[30]) ^ uint32(s.x[17]<<24))
	s.x[22] = uint64(uint32(s.x[20]) ^ 0xff77b7a7)
	s.x[17] = uint64(uint32(s.x[17]) & uint32(s.x[22]))
	s.x[22] = 0x00980000
	s.x[20] = uint64(uint32(s.x[20]) & uint32(s.x[22]))
	s.x[17] = uint64(uint32(s.x[17]) | uint32(s.x[20]))
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[30]))
	s.x[19] = uint64((uint32(s.x[19]) & 0x7f) << 1)
	s.x[19] ^= s.x[7]
	value, err = s.load(s.x[28]+s.x[19], 1)
	if err != nil {
		return 0, fmt.Errorf("shared permutation second table byte: %w", err)
	}
	s.x[19] = uint64(byte(value))
	s.x[19] = uint64(uint32(s.x[19]) ^ uint32(s.x[8]))
	s.x[17] = uint64(uint32(s.x[17]) | uint32(s.x[19]<<8))
	s.x[19] = uint64(uint32(s.x[12])*uint32(s.x[16]) + 0xa2)
	s.x[19] &= 0xff
	value, err = s.load(s.x[28]+s.x[19], 1)
	if err != nil {
		return 0, fmt.Errorf("shared permutation third table byte: %w", err)
	}
	s.x[19] = uint64(byte(value))
	s.x[20] = 0x7db1fa90
	s.x[19] = uint64(uint32(s.x[19]) ^ uint32(s.x[20]))
	s.x[19] = uint64(uint32(s.x[19]) &^ uint32(s.x[17]))
	s.x[20] = 0x7db1fa8b
	s.x[17] = uint64(uint32(s.x[17]) | uint32(s.x[20]))
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[19]))
	s.x[19] = uint64(uint32(s.x[17]) << 3)
	s.x[19] = uint64(uint32(s.x[19]) ^ uint32(s.x[17])>>5)
	s.x[17] = uint64(uint32(s.x[19]) + uint32(s.x[17]))
	s.x[19] = uint64(uint32(s.x[17]) * uint32(s.x[17]))
	s.x[20] = uint64(uint32(s.x[17]) >> 8)
	s.x[22] = uint64(uint32(s.x[17]) >> 16)
	s.x[17] = uint64(uint32(s.x[19])*uint32(s.x[17]) + uint32(s.x[22]))
	s.x[19] = uint64(uint32(s.x[20]) * uint32(s.x[20]))
	s.x[20] = uint64(uint32(s.x[19]) * uint32(s.x[20]))
	s.x[17] = uint64(uint32(s.x[20])*uint32(s.x[19]) + uint32(s.x[17]))
	s.x[19] = uint64(uint32(s.x[22]) * uint32(s.x[22]))
	s.x[19] = uint64(uint32(s.x[19]) * uint32(s.x[22]))
	s.x[20] = uint64(uint32(s.x[19]) * uint32(s.x[22]))
	s.x[17] = uint64(uint32(s.x[20])*uint32(s.x[19]) + uint32(s.x[17]))
	s.x[19] = 0xd4
	s.x[19] = uint64(uint32(s.x[19]) & uint32(s.x[17]<<1))
	s.x[17] = uint64(uint32(s.x[17]) - uint32(s.x[19]))
	s.x[17] = uint64(uint32(s.x[17]) + 0x6a)
	s.x[19] = 0x71
	s.x[20] = uint64(uint32(s.x[17]) ^ uint32(s.x[19]))
	next, err := loadBootstrapPointer(s, s.x[21], uint32(s.x[2]))
	if err != nil {
		return 0, fmt.Errorf("shared permutation target: %w", err)
	}
	s.x[2] = next
	s.x[27] = bootstrapTargetTable
	s.x[17] = 0xcd20
	return next, nil
}
