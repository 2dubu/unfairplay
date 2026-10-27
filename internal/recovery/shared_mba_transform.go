package recovery

import "fmt"

const (
	sharedMBAFront = uint64(0x1005515e0)
	sharedMBATail  = uint64(0x100551800)
)

func runSharedMBAFront(s *state) (uint64, error) {
	s.traceBlock(sharedMBAFront)
	s.x[17] = 0x7c16
	s.x[1] = uint64(uint32(s.x[19]) ^ uint32(s.x[17]))
	s.x[17] = uint64(uint32(s.x[15]) ^ 0x2d)
	s.x[17] = uint64(uint32(s.x[17]) + 0x7f)
	s.x[2] = uint64(uint32(s.x[15]) ^ 0x23)
	s.x[2] = uint64(uint32(s.x[2]) + 0x71)
	s.x[20] = 0xffff9a9d
	s.x[20] = uint64(uint32(s.x[1]) + uint32(s.x[20]))
	s.x[20] = uint64(uint32(s.x[20]) ^ uint32(s.x[15]))
	s.x[20] = uint64(uint32(s.x[20]) + 0xf)
	s.x[23] = uint64(uint32(s.x[15]) ^ uint32(s.x[14]))
	s.x[2] = uint64(uint32(s.x[23]) ^ uint32(s.x[2]))
	s.x[17] = uint64(uint32(s.x[2]) ^ uint32(s.x[17]))
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[20]))
	s.x[17] &= 0xff
	s.x[2] = 0xac
	s.x[17] ^= s.x[2]
	value, err := s.load(s.sp+0x58, 8)
	if err != nil {
		return 0, fmt.Errorf("shared MBA table: %w", err)
	}
	s.x[28] = value
	value, err = s.load(s.x[28]+s.x[17], 1)
	if err != nil {
		return 0, fmt.Errorf("shared MBA table byte: %w", err)
	}
	s.x[17] = uint64(byte(value))
	s.x[2] = uint64(uint32(s.x[17]) ^ 0xcc667f2c)
	s.x[20] = 0x339980d4
	s.x[2] = uint64(uint32(s.x[2]) + uint32(s.x[20]))
	s.x[20] = uint64(uint32(s.x[17]) ^ 0x0aa63e3e)
	s.x[23] = 0xf559c1c2
	s.x[20] = uint64(uint32(s.x[20]) + uint32(s.x[23]))
	s.x[2] = uint64(uint32(s.x[2]) ^ uint32(s.x[20]))
	s.x[20] = 0xc6c04109
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[20]))
	s.x[20] = 0x393fbef7
	s.x[17] = uint64(uint32(s.x[17]) + uint32(s.x[20]))
	s.x[17] = uint64(uint32(s.x[2]) ^ uint32(s.x[17]))
	s.x[2] = uint64(uint32(s.x[17]) + 0x1b)
	s.x[3] = 0x75
	s.x[20] = 0x48575af0
	s.x[20] = uint64(uint32(s.x[22])*uint32(s.x[3]) + uint32(s.x[20]))
	s.x[3] = 0xb7a8b167
	s.x[23] = uint64(uint32(s.x[22]) * uint32(s.x[3]))
	s.x[24] = uint64(uint32(s.x[22]) + 0x1b)
	s.x[20] = uint64(uint32(s.x[20])*uint32(s.x[24]) + uint32(s.x[23]))
	s.x[3] = 0xc42b10eb
	s.x[23] = uint64(uint32(s.x[22]) + uint32(s.x[3]))
	s.x[3] = 0x3bd4ef30
	s.x[25] = uint64(uint32(s.x[20]) * uint32(s.x[3]))
	s.x[3] = 0x5ecab5dd
	s.x[20] = uint64(uint32(s.x[20]) + uint32(s.x[3]))
	s.x[3] = 0x9e9ca0fb
	s.x[25] = uint64(uint32(s.x[22])*uint32(s.x[3]) + uint32(s.x[25]))
	s.x[20] = uint64(uint32(s.x[20])*uint32(s.x[23]) + uint32(s.x[25]))
	s.x[3] = 0xb
	s.x[23] = 0xa40d49ed
	s.x[23] = uint64(uint32(s.x[17])*uint32(s.x[3]) + uint32(s.x[23]))
	s.x[23] = uint64(uint32(s.x[23]) * uint32(s.x[2]))
	s.x[3] = 0x5bf2b73c
	s.x[23] = uint64(uint32(s.x[17])*uint32(s.x[3]) + uint32(s.x[23]))
	s.x[3] = 0x4ed6e7d8
	s.x[23] = uint64(uint32(s.x[23]) + uint32(s.x[3]))
	s.x[23] = uint64(uint32(s.x[23]) * uint32(s.x[2]))
	s.x[3] = 0x63c26b7c
	s.x[23] = uint64(uint32(s.x[17])*uint32(s.x[3]) + uint32(s.x[23]))
	s.x[3] = 0xc961596a
	s.x[25] = uint64(uint32(s.x[17]) + uint32(s.x[3]))
	s.x[3] = 0x85815614
	s.x[27] = uint64(uint32(s.x[23]) + uint32(s.x[3]))
	s.x[25] = uint64(uint32(s.x[27]) * uint32(s.x[25]))
	s.x[27] = 0x9d449283
	s.x[27] = uint64(uint32(s.x[17]) + uint32(s.x[27]))
	s.x[27] = uint64(uint32(s.x[27]) * uint32(s.x[2]))
	s.x[3] = 0x62bb6d98
	s.x[27] = uint64(uint32(s.x[17])*uint32(s.x[3]) + uint32(s.x[27]))
	s.x[3] = 0x369ea6b1
	s.x[23] = uint64(uint32(s.x[23])*uint32(s.x[3]) + uint32(s.x[25]))
	s.x[3] = 0x6345d1c0
	s.x[25] = uint64(uint32(s.x[17]) + uint32(s.x[3]))
	s.x[3] = 0x9cba2e5b
	s.x[20] = uint64(uint32(s.x[27])*uint32(s.x[3]) + uint32(s.x[20]))
	s.x[3] = 0x69c48f08
	s.x[27] = uint64(uint32(s.x[27]) + uint32(s.x[3]))
	s.x[20] = uint64(uint32(s.x[27])*uint32(s.x[25]) + uint32(s.x[20]))
	s.x[3] = 0xe7
	s.x[25] = 0xd7876c42
	s.x[25] = uint64(uint32(s.x[17])*uint32(s.x[3]) + uint32(s.x[25]))
	s.x[3] = 0x2878ac1b
	s.x[27] = uint64(uint32(s.x[17]) * uint32(s.x[3]))
	s.x[3] = 0xd8f6830e
	s.x[27] = uint64(uint32(s.x[22])*uint32(s.x[3]) + uint32(s.x[27]))
	s.x[25] = uint64(uint32(s.x[25])*uint32(s.x[2]) + uint32(s.x[27]))
	s.x[27] = bootstrapTargetTable
	s.x[3] = 0x1ca7fc33
	s.x[23] = uint64(uint32(s.x[23]) + uint32(s.x[3]))
	s.x[2] = uint64(uint32(s.x[23]) * uint32(s.x[2]))
	s.x[3] = 0x67ed7fa1
	s.x[2] = uint64(uint32(s.x[17])*uint32(s.x[3]) + uint32(s.x[2]))
	s.x[3] = 0x576fd500
	s.x[2] = uint64(uint32(s.x[2]) + uint32(s.x[3]))
	s.x[2] = uint64(uint32(s.x[2])*uint32(s.x[24]) + uint32(s.x[20]))
	s.x[3] = 0x27c02bb3
	s.x[2] = uint64(uint32(s.x[2]) + uint32(s.x[3]))
	s.x[2] = uint64(uint32(s.x[2])*uint32(s.x[24]) + uint32(s.x[25]))
	s.x[3] = 0xcda8f9e7
	return sharedMBATail, nil
}

func runSharedMBATail(s *state) (uint64, error) {
	s.traceBlock(sharedMBATail)
	s.x[20] = uint64(uint32(s.x[22]) + uint32(s.x[3]))
	s.x[3] = 0x32570634
	s.x[17] = uint64(uint32(s.x[2])*uint32(s.x[3]) + uint32(s.x[17]))
	s.x[3] = 0x26b9f953
	s.x[2] = uint64(uint32(s.x[2]) + uint32(s.x[3]))
	s.x[17] = uint64(uint32(s.x[2])*uint32(s.x[20]) + uint32(s.x[17]))
	s.x[2] = 0x2ad396f7
	s.x[17] = uint64(uint32(s.x[17]) + uint32(s.x[2]))
	s.x[2] = 0x36b45a9b
	s.x[2] = uint64(int64(int32(s.x[17])) * int64(int32(s.x[2])))
	s.x[2] = shiftRightArithmetic(s.x[2], 39, 64)
	s.x[2] = uint64(uint32(s.x[2]) + uint32(shiftRight(s.x[2], 31, 32)))
	s.x[3] = 0x257
	s.x[17] = uint64(uint32(s.x[17]) - uint32(s.x[2])*uint32(s.x[3]))
	s.x[2] = 0x6e6cd6a0
	s.x[2] = uint64(uint32(s.x[2]) & uint32(s.x[17]<<1))
	s.x[17] = uint64(uint32(s.x[17]) - uint32(s.x[2]))
	s.x[2] = 0x37366b50
	s.x[17] = uint64(uint32(s.x[17]) + uint32(s.x[2]))
	s.x[2] = 0x196dc519
	s.x[2] = uint64(uint32(s.x[17]) ^ uint32(s.x[2]))
	s.x[3] = 0xe6923ae7
	s.x[2] = uint64(uint32(s.x[2]) + uint32(s.x[3]))
	s.x[3] = 0xb118082f
	s.x[20] = uint64(uint32(s.x[17]) ^ uint32(s.x[3]))
	s.x[3] = 0x4ee7f7d1
	s.x[20] = uint64(uint32(s.x[20]) + uint32(s.x[3]))
	s.x[2] = uint64(uint32(s.x[2]) ^ uint32(s.x[20]))
	s.x[3] = 0x9f43a666
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[3]))
	s.x[3] = 0x60bc599a
	s.x[17] = uint64(uint32(s.x[17]) + uint32(s.x[3]))
	s.x[17] = uint64(uint32(s.x[2]) ^ uint32(s.x[17]))
	s.x[17] = uint64(uint32(s.x[17]) + uint32(s.x[17]<<7))
	s.x[2] = 0x2ac71237
	s.x[17] = uint64(uint32(s.x[17]) + uint32(s.x[2]))
	s.x[2] = uint64(uint32(s.x[17]) ^ 0xc5c184e0)
	s.x[3] = 0xa7a501e2
	s.x[20] = uint64(uint32(s.x[17]) & uint32(s.x[3]))
	s.x[3] = 0x878185e0
	s.x[22] = uint64(uint32(s.x[17]) ^ uint32(s.x[3]))
	s.x[22] = uint64(uint32(s.x[22]) & uint32(s.x[20]<<1))
	s.x[20] = uint64(uint32(s.x[22]) ^ uint32(s.x[20]))
	s.x[3] = 0xc4c90a04
	s.x[22] = uint64(uint32(s.x[3]) ^ uint32(s.x[2]<<1))
	s.x[2] = uint64(uint32(s.x[22]) & uint32(s.x[2]))
	s.x[22] = uint64(uint32(s.x[2]) & uint32(s.x[20]<<2))
	s.x[3] = 0x62648502
	s.x[20] = uint64(uint32(s.x[20]) ^ uint32(s.x[3]))
	s.x[3] = 0x060be5d0
	s.x[20] = uint64(uint32(s.x[20]) ^ uint32(s.x[22]))
	s.x[22] = uint64(uint32(s.x[20]) << 4)
	s.x[22] = uint64(uint32(s.x[22]) ^ 0x20)
	s.x[2] = uint64(uint32(s.x[2]) & uint32(s.x[2]<<2))
	s.x[2] = uint64(uint32(s.x[2]) & uint32(s.x[22]))
	s.x[2] = uint64(uint32(s.x[2]) ^ uint32(s.x[20]))
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[2]<<1))
	s.x[17] &= 0xff
	s.x[2] = 0xe6
	s.x[17] ^= s.x[2]
	value, err := s.load(s.sp+0x50, 8)
	if err != nil {
		return 0, fmt.Errorf("shared MBA output table: %w", err)
	}
	s.x[2] = value
	value, err = s.load(value+s.x[17], 1)
	if err != nil {
		return 0, fmt.Errorf("shared MBA output byte: %w", err)
	}
	s.x[20] = uint64(byte(value))
	next, err := loadBootstrapPointer(s, s.x[21], uint32(s.x[19]))
	if err != nil {
		return 0, fmt.Errorf("shared MBA target: %w", err)
	}
	s.x[2] = next
	s.x[17] = 0xcd20
	return next, nil
}
