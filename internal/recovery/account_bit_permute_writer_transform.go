package recovery

import "fmt"

const accountBitPermuteWriter = uint64(0x100164d54)

func runAccountBitPermuteWriter(s *state) (uint64, error) {
	s.traceBlock(accountBitPermuteWriter)
	s.x[17] = s.x[2] + s.x[14]
	value, err := s.load(s.x[27]+s.x[17], 1)
	if err != nil {
		return 0, fmt.Errorf("account bit-permute first byte: %w", err)
	}
	s.x[17] = uint64(byte(value))
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[6]))
	s.x[3] = s.x[2] + s.x[1]
	value, err = s.load(s.x[27]+s.x[3], 1)
	if err != nil {
		return 0, fmt.Errorf("account bit-permute second byte: %w", err)
	}
	s.x[3] = uint64(byte(value))
	s.x[3] = uint64(uint32(s.x[3]) ^ uint32(s.x[6]))
	s.x[3] = uint64(uint32(shiftLeft(s.x[3], int(s.x[0]), 32)))
	s.x[17] = uint64(uint32(s.x[3]) + uint32(s.x[17]<<24))
	s.x[11] = 0xffffffffb7e46f2b
	s.x[3] = s.x[2] + s.x[11]
	value, err = s.load(s.x[27]+s.x[3], 1)
	if err != nil {
		return 0, fmt.Errorf("account bit-permute third byte: %w", err)
	}
	s.x[3] = uint64(byte(value))
	s.x[3] = uint64(uint32(s.x[3]) ^ uint32(s.x[6]))
	s.x[11] = 0xffffffffb7e46f2c
	s.x[22] = s.x[2] + s.x[11]
	value, err = s.load(s.x[27]+s.x[22], 1)
	if err != nil {
		return 0, fmt.Errorf("account bit-permute fourth byte: %w", err)
	}
	s.x[22] = uint64(byte(value))
	s.x[17] = uint64(uint32(s.x[17]) + uint32(s.x[3]<<8))
	s.x[3] = uint64(uint32(s.x[22]) ^ uint32(s.x[6]))
	s.x[17] = uint64(uint32(s.x[17]) + uint32(s.x[3]))
	s.x[22] = uint64(uint32(s.x[17]) & uint32(s.x[5]))
	s.x[3] &= 0x10
	s.x[3] = uint64(uint32(s.x[22]) ^ uint32(s.x[3]))
	s.x[17] = uint64(uint32(s.x[17]) - uint32(s.x[3]<<1))
	s.x[3] = uint64(uint32(s.x[2]) ^ uint32(s.x[7]))
	s.x[22] = uint64(uint32(s.x[2]) & uint32(s.x[30]))
	s.x[25] = uint64(uint32(s.x[22]) << 1)
	s.x[11] = 0x24800000
	s.x[26] = uint64(uint32(s.x[25]) | uint32(s.x[11]))
	s.x[26] = uint64(uint32(s.x[26]) & uint32(s.x[2]))
	s.x[25] &= 0x200
	s.x[11] = 0x00242200
	s.x[25] = uint64(uint32(s.x[25]) ^ uint32(s.x[11]))
	s.x[22] = uint64(uint32(s.x[25]) | uint32(s.x[22]))
	s.x[22] = uint64(uint32(s.x[22]) | uint32(s.x[26]))
	s.x[11] = 0xec9a5508
	s.x[22] = uint64(uint32(s.x[22]) ^ uint32(s.x[11]))
	s.x[25] = uint64(uint32(s.x[3]) << 1)
	s.x[11] = 0xfdd8ce10
	s.x[26] = uint64(uint32(s.x[25]) ^ uint32(s.x[11]))
	s.x[3] = uint64(uint32(s.x[26]) & uint32(s.x[3]))
	s.x[11] = 0x00021308
	s.x[25] = uint64(uint32(s.x[25]) & uint32(s.x[11]))
	s.x[25] = uint64(uint32(s.x[3]) ^ uint32(s.x[25]))
	s.x[11] = 0x951c74e0
	s.x[26] = uint64(uint32(s.x[11]) ^ uint32(s.x[22]<<2))
	s.x[11] = 0x034188c0
	s.x[3] = uint64(uint32(s.x[3]) & uint32(s.x[11]))
	s.x[11] = 0xecbc6600
	s.x[28] = uint64(uint32(s.x[25]) ^ uint32(s.x[11]))
	s.x[26] = uint64(uint32(s.x[26]) & uint32(s.x[28]))
	s.x[22] = uint64(uint32(s.x[22]) ^ uint32(s.x[3]))
	s.x[22] = uint64(uint32(s.x[22]) ^ uint32(s.x[26]))
	s.x[11] = 0x951430c0
	s.x[25] = uint64(uint32(s.x[11]) ^ uint32(s.x[25]<<2))
	s.x[25] = uint64(uint32(s.x[25]) & uint32(s.x[28]))
	s.x[3] = uint64(uint32(s.x[25]) ^ uint32(s.x[3]))
	s.x[11] = 0xa640d840
	s.x[25] = uint64(uint32(s.x[11]) ^ uint32(s.x[22]<<4))
	s.x[11] = 0x27e5a8c0
	s.x[26] = uint64(uint32(s.x[3]) & uint32(s.x[11]))
	s.x[21] = 0x24a42008
	s.x[28] = uint64(uint32(s.x[3]) ^ uint32(s.x[21]))
	s.x[25] = uint64(uint32(s.x[25]) & uint32(s.x[28]))
	s.x[22] = uint64(uint32(s.x[22]) ^ uint32(s.x[26]))
	s.x[22] = uint64(uint32(s.x[22]) ^ uint32(s.x[25]))
	s.x[25] = 0x6da7a8c0
	s.x[3] = uint64(uint32(s.x[25]) ^ uint32(s.x[3]<<4))
	s.x[3] = uint64(uint32(s.x[3]) & uint32(s.x[28]))
	s.x[3] = uint64(uint32(s.x[3]) ^ uint32(s.x[26]))
	s.x[25] = uint64(uint32(0x9992a0c0) ^ uint32(s.x[22]<<8))
	s.x[26] = uint64(uint32(s.x[3]) & uint32(s.x[11]))
	s.x[28] = uint64(uint32(s.x[3]) ^ uint32(s.x[21]))
	s.x[25] = uint64(uint32(s.x[25]) & uint32(s.x[28]))
	s.x[22] = uint64(uint32(s.x[22]) ^ uint32(s.x[26]))
	s.x[22] = uint64(uint32(s.x[22]) ^ uint32(s.x[25]))
	s.x[21] = 0x83c5a8c0
	s.x[3] = uint64(uint32(s.x[21]) ^ uint32(s.x[3]<<8))
	s.x[3] = uint64(uint32(s.x[3]) & uint32(s.x[28]))
	s.x[3] = uint64(uint32(s.x[3]) ^ uint32(s.x[26]))
	s.x[21] = 0x481a0000
	s.x[25] = uint64(uint32(s.x[3]) ^ uint32(s.x[21]))
	s.x[21] = 0x49d83ed2
	s.x[3] = uint64(uint32(s.x[3]) & uint32(s.x[11]))
	s.x[3] = uint64(uint32(s.x[22]) ^ uint32(s.x[3]))
	s.x[11] = 0x70eda8c1
	s.x[22] = uint64(uint32(s.x[11]) ^ uint32(s.x[22]<<16))
	s.x[25] = uint64(uint32(s.x[22]) & uint32(s.x[25]))
	s.x[3] = uint64(uint32(s.x[3]) ^ uint32(s.x[25]))
	s.x[11] = 0x6cbe2000
	s.x[22] = uint64(uint32(s.x[22]) & uint32(s.x[11]))
	s.x[3] = uint64(uint32(s.x[3]) ^ uint32(s.x[22]))
	s.x[17] = uint64(uint32(s.x[17]) + uint32(s.x[13]))
	s.x[3] = uint64(uint32(s.x[2]) ^ uint32(s.x[3]<<1))
	s.x[11] = 0x1ba6204e
	s.x[3] = uint64(uint32(s.x[11]) ^ uint32(s.x[3])>>2)
	if err := s.storeValue(s.x[20]+uint64(uint32(s.x[3]))*4, uint64(uint32(s.x[17])), 4); err != nil {
		return 0, fmt.Errorf("account bit-permute output: %w", err)
	}
	s.x[11] = 0xffffffffb7e46f2d
	s.x[17] = s.x[2] + s.x[11]
	s.x[2] += 4
	s.compare64(s.x[17], 0x40)
	gate := uint32(boolValue(s.condition("lo")))
	s.x[17] = uint64(gate * 0x2755)
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[15]))
	next, err := loadBootstrapPointer(s, s.x[24], uint32(s.x[17]))
	if err != nil {
		return 0, fmt.Errorf("account bit-permute target: %w", err)
	}
	s.x[22] = next
	s.x[3] = 0x481b90e7
	s.x[17] = 0xaaed
	return next, nil
}
