package recovery

import "fmt"

func runAccountHotMixSecondCircuit(s *state) (uint64, error) {
	s.traceBlock(accountHotMixSecond)
	s.x[7] = s.x[11]
	s.x[11] = s.x[13]
	s.x[19] = s.x[12]
	s.x[12] = uint64(uint32(s.x[9]) + 0x0989137d)
	s.x[13] = 0xf676ffff
	s.x[16] = uint64(uint32(s.x[12]) & uint32(s.x[13]))
	s.x[12] = uint64(uint32(rotateRight(s.x[19], 27, 32)))
	s.x[13] = uint64(uint32(s.x[12]) ^ 0x05b55cde)
	s.x[17] = uint64(uint32(s.x[16]) ^ 0xd447c6c8)
	s.x[2] = 0xc1946e24
	s.x[23] = uint64(uint32(s.x[12]) ^ uint32(s.x[2]))
	s.x[2] = 0x106693d4
	s.x[23] = uint64(uint32(s.x[23]) + uint32(s.x[2]))
	s.x[2] = 0xa7d8855b
	s.x[12] = uint64(uint32(s.x[12]) ^ uint32(s.x[2]))
	s.x[2] = 0x762a78ad
	s.x[12] = uint64(uint32(s.x[12]) + uint32(s.x[2]))
	s.x[12] = uint64(uint32(s.x[23]) ^ uint32(s.x[12]))
	s.x[13] = uint64(uint32(s.x[13]) + uint32(s.x[17]))
	s.x[12] = uint64(uint32(s.x[13]) ^ uint32(s.x[12]))
	s.x[13] = uint64(uint32(s.x[11]) ^ uint32(s.x[15]))
	s.x[13] = uint64(uint32(s.x[13]) ^ uint32(s.x[7]))
	s.x[17] = uint64(uint32(s.x[13]) ^ 0x934ce200)
	s.x[2] = 0xfff36668
	s.x[17] = uint64(uint32(s.x[17]) + uint32(s.x[2]))
	s.x[2] = 0x01f12ca6
	s.x[23] = uint64(uint32(s.x[13]) ^ uint32(s.x[2]))
	s.x[2] = 0x6d4ea8c2
	s.x[23] = uint64(uint32(s.x[23]) + uint32(s.x[2]))
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[23]))
	s.x[2] = 0x99a2031b
	s.x[13] = uint64(uint32(s.x[13]) ^ uint32(s.x[2]))
	s.x[2] = 0xf51d877d
	s.x[13] = uint64(uint32(s.x[13]) + uint32(s.x[2]))
	s.x[13] = uint64(uint32(s.x[17]) ^ uint32(s.x[13]))
	s.x[17] = uint64(uint32(s.x[4]) ^ 0xe1dd6fa5)
	s.x[2] = 0x1e22905b
	s.x[17] = uint64(uint32(s.x[17]) + uint32(s.x[2]))
	s.x[2] = 0x2118de94
	s.x[23] = uint64(uint32(s.x[4]) ^ uint32(s.x[2]))
	s.x[2] = 0xdee7216c
	s.x[23] = uint64(uint32(s.x[23]) + uint32(s.x[2]))
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[23]))
	s.x[2] = 0xcbda7c8c
	s.x[4] = uint64(uint32(s.x[4]) ^ uint32(s.x[2]))
	s.x[2] = 0x34258374
	s.x[4] = uint64(uint32(s.x[4]) + uint32(s.x[2]))
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[4]))
	s.x[4] = s.x[6] + (s.x[20] << 2)
	word, err := s.load(s.x[10]+s.x[4], 4)
	if err != nil {
		return 0, fmt.Errorf("account second hot-mix table word: %w", err)
	}
	s.x[4] = uint64(uint32(word))
	s.x[23] = uint64(uint32(s.x[4]) ^ 0xaa31e45f)
	s.x[2] = 0x55ce1ba1
	s.x[23] = uint64(uint32(s.x[23]) + uint32(s.x[2]))
	s.x[24] = uint64(uint32(s.x[4]) ^ 0xfaf6ab14)
	s.x[2] = 0x050954ec
	s.x[24] = uint64(uint32(s.x[24]) + uint32(s.x[2]))
	s.x[23] = uint64(uint32(s.x[23]) ^ uint32(s.x[24]))
	s.x[2] = 0x9eb2348e
	s.x[4] = uint64(uint32(s.x[4]) ^ uint32(s.x[2]))
	s.x[2] = 0x614dcb72
	s.x[4] = uint64(uint32(s.x[4]) + uint32(s.x[2]))
	s.x[4] = uint64(uint32(s.x[23]) ^ uint32(s.x[4]))
	s.x[12] = uint64(uint32(s.x[12]) + uint32(s.x[17]))
	s.x[12] = uint64(uint32(s.x[12]) + uint32(s.x[13]))
	s.x[13] = 0xe37ef5e7
	s.x[12] = uint64(uint32(s.x[12]) + uint32(s.x[13]))
	s.x[12] = uint64(uint32(s.x[12]) + uint32(s.x[4]))
	s.x[13] = uint64(uint32(s.x[12]) ^ 0x61c97bf8)
	s.x[17] = uint64(uint32(s.x[12]) ^ 0xf828ef83)
	s.x[2] = 0xb36c448b
	s.x[17] = uint64(uint32(s.x[17]) + uint32(s.x[2]))
	s.x[2] = 0xf2a13ff3
	s.x[4] = uint64(uint32(s.x[12]) ^ uint32(s.x[2]))
	s.x[2] = 0xdffbff7e
	s.x[23] = uint64(uint32(s.x[12]) ^ uint32(s.x[2]))
	s.x[2] = 0x2a8dd0f2
	s.x[13] = uint64(uint32(s.x[13]) + uint32(s.x[2]))
	s.x[12] = uint64(uint32(s.x[12]) ^ uint32(s.x[13]))
	s.x[13] = uint64(uint32(s.x[4]) + 0xb9e594fb)
	s.x[13] = uint64(uint32(s.x[17]) ^ uint32(s.x[13]))
	s.x[12] = uint64(uint32(s.x[12]) ^ uint32(s.x[13]))
	s.x[13] = uint64(uint32(s.x[23]) + 0x94bf5478)
	s.x[17] = 0xbfa4994b
	s.x[13] = uint64(uint32(s.x[13]) ^ uint32(s.x[17]))
	s.x[12] = uint64(uint32(s.x[12]) ^ uint32(s.x[13]))
	s.x[13] = uint64(uint32(s.x[30]) ^ uint32(rotateRight(s.x[15], 2, 32)))
	s.x[20]++
	s.x[15] = 0x2ba7c3b3
	s.compare64(s.x[20], s.x[15])
	gate := uint32(boolValue(s.condition("eq")))
	s.x[15] = uint64(gate * 0xc7a1)
	s.x[15] = uint64(uint32(s.x[15]) ^ uint32(s.x[9]))
	next, err := loadBootstrapPointer(s, s.x[21], uint32(s.x[15]))
	if err != nil {
		return 0, fmt.Errorf("account second hot-mix target: %w", err)
	}
	s.x[24] = next
	s.x[15] = s.x[19]
	s.x[4] = s.x[7]
	s.x[23] = 0x2ba7c3b3
	s.x[17] = 0x1d92
	return next, nil
}
