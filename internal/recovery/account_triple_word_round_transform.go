package recovery

import "fmt"

const (
	accountTripleWordRoundFront = uint64(0x1001c4ed0)
	accountTripleWordRoundTail  = uint64(0x1001c5000)
)

func runAccountTripleWordRound(s *state) (uint64, error) {
	s.traceBlock(accountTripleWordRoundFront)
	s.x[16] = s.x[10]
	s.x[12] = s.x[8]
	s.x[10] = s.x[19]
	s.x[8] = uint64(uint32(rotateRight(s.x[8], 27, 32)))
	s.x[11] = 0x15d0d56e
	s.x[15] = uint64(uint32(s.x[8]) ^ uint32(s.x[11]))
	s.x[11] = 0x4c717305
	s.x[15] = uint64(uint32(s.x[15]) + uint32(s.x[11]))
	s.x[11] = 0x45015859
	s.x[17] = uint64(uint32(s.x[8]) ^ uint32(s.x[11]))
	s.x[11] = 0x1ca0fe34
	s.x[17] = uint64(uint32(s.x[17]) + uint32(s.x[11]))
	s.x[15] = uint64(uint32(s.x[15]) ^ uint32(s.x[17]))
	s.x[11] = 0x33283a96
	s.x[8] = uint64(uint32(s.x[8]) ^ uint32(s.x[11]))
	s.x[11] = 0x6a899cfd
	s.x[8] = uint64(uint32(s.x[8]) + uint32(s.x[11]))
	s.x[8] = uint64(uint32(s.x[15]) ^ uint32(s.x[8]))
	s.x[11] = 0x03f3f4a1
	s.x[15] = uint64(uint32(s.x[19]) ^ uint32(s.x[11]))
	s.x[11] = 0xf4df8cc5
	s.x[17] = uint64(uint32(s.x[0]) + uint32(s.x[11]))
	s.x[17] = uint64(uint32(s.x[16]) ^ uint32(s.x[17]))
	s.x[15] = uint64(uint32(s.x[17]) & uint32(s.x[15]))
	s.x[11] = 0xf713c6e3
	s.x[17] = uint64(uint32(s.x[16]) & uint32(s.x[11]))
	s.x[15] = uint64(uint32(s.x[15]) ^ uint32(s.x[17]))
	s.x[11] = 0x5064c30f
	s.x[17] = uint64(uint32(s.x[3]) ^ uint32(s.x[11]))
	s.x[11] = 0x0bfffdbd
	s.x[19] = uint64(uint32(s.x[15]) ^ uint32(s.x[11]))
	s.x[17] = uint64(uint32(s.x[19]) & uint32(s.x[17]))
	s.x[11] = 0x5b7b0eb2
	s.x[15] = uint64(uint32(s.x[15]) & uint32(s.x[11]))
	s.x[11] = 0x619654a3
	s.x[19] = uint64(uint32(s.x[10]) ^ uint32(s.x[11]))
	s.x[22] = uint64(uint32(s.x[16]) ^ uint32(s.x[23]))
	s.x[19] = uint64(uint32(s.x[22]) & uint32(s.x[19]))
	s.x[11] = 0x6a89991e
	s.x[22] = uint64(uint32(s.x[16]) & uint32(s.x[11]))
	s.x[19] = uint64(uint32(s.x[19]) ^ uint32(s.x[22]))
	s.x[11] = 0xe239ae99
	s.x[15] = uint64(uint32(s.x[15]) ^ uint32(s.x[11]))
	s.x[15] = uint64(uint32(s.x[15]) ^ uint32(s.x[17]))
	s.x[11] = 0xf5f676e3
	s.x[17] = uint64(uint32(s.x[19]) ^ uint32(s.x[11]))
	s.x[15] = uint64(uint32(s.x[15]) & uint32(s.x[17]))
	s.x[11] = 0x16bd5dd6
	s.x[17] = uint64(uint32(s.x[19]) & uint32(s.x[11]))
	s.x[15] = uint64(uint32(s.x[15]) ^ uint32(s.x[17]))
	s.x[11] = 0x5a0aa5fb
	s.x[17] = uint64(uint32(s.x[15]) ^ uint32(s.x[11]))
	s.x[11] = 0xbdd69e58
	s.x[17] = uint64(uint32(s.x[17]) + uint32(s.x[11]))
	s.x[11] = 0xed8eceb0
	s.x[19] = uint64(uint32(s.x[15]) ^ uint32(s.x[11]))

	s.traceBlock(accountTripleWordRoundTail)
	s.x[11] = 0x0a52f51d
	s.x[19] = uint64(uint32(s.x[19]) + uint32(s.x[11]))
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[19]))
	s.x[11] = 0x5ccfc076
	s.x[15] = uint64(uint32(s.x[15]) ^ uint32(s.x[11]))
	s.x[11] = 0xbb13fbdb
	s.x[15] = uint64(uint32(s.x[15]) + uint32(s.x[11]))
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[15]))
	s.x[11] = 0xafb3d14c
	s.x[15] = uint64(uint32(s.x[2]) ^ uint32(s.x[11]))
	s.x[11] = 0x504c2eb4
	s.x[15] = uint64(uint32(s.x[15]) + uint32(s.x[11]))
	s.x[11] = 0x5ec1af1c
	s.x[19] = uint64(uint32(s.x[2]) ^ uint32(s.x[11]))
	s.x[11] = 0xa13e50e4
	s.x[19] = uint64(uint32(s.x[19]) + uint32(s.x[11]))
	s.x[15] = uint64(uint32(s.x[15]) ^ uint32(s.x[19]))
	s.x[11] = 0xfa6db3ed
	s.x[2] = uint64(uint32(s.x[2]) ^ uint32(s.x[11]))
	s.x[11] = 0x05924c13
	s.x[2] = uint64(uint32(s.x[2]) + uint32(s.x[11]))
	s.x[2] = uint64(uint32(s.x[15]) ^ uint32(s.x[2]))
	s.x[15] = s.x[9] + (s.x[26] << 2)
	word, err := s.load(s.x[20]+s.x[15], 4)
	if err != nil {
		return 0, fmt.Errorf("account triple-word round table word: %w", err)
	}
	s.x[19] = uint64(uint32(word))
	s.x[11] = 0x4bd2960a
	s.x[15] = uint64(uint32(s.x[19]) ^ uint32(s.x[11]))
	s.x[11] = 0xb42d69f6
	s.x[15] = uint64(uint32(s.x[15]) + uint32(s.x[11]))
	s.x[11] = 0xc43c4139
	s.x[22] = uint64(uint32(s.x[19]) ^ uint32(s.x[11]))
	s.x[11] = 0x3bc3bec7
	s.x[22] = uint64(uint32(s.x[22]) + uint32(s.x[11]))
	s.x[22] = uint64(uint32(s.x[15]) ^ uint32(s.x[22]))
	s.x[11] = 0xffffe106
	s.x[15] = uint64(uint32(s.x[0]) + uint32(s.x[11]))
	s.x[11] = 0x419bacf6
	s.x[19] = uint64(uint32(s.x[19]) ^ uint32(s.x[11]))
	s.x[11] = 0xbe64530a
	s.x[19] = uint64(uint32(s.x[19]) + uint32(s.x[11]))
	s.x[19] = uint64(uint32(s.x[22]) ^ uint32(s.x[19]))
	s.x[8] = uint64(uint32(s.x[8]) + uint32(s.x[2]))
	s.x[8] = uint64(uint32(s.x[8]) + uint32(s.x[19]))
	s.x[11] = 0xa1c502f7
	s.x[8] = uint64(uint32(s.x[8]) + uint32(s.x[11]))
	s.x[8] = uint64(uint32(s.x[8]) + uint32(s.x[17]))
	s.x[11] = 0x57752ccd
	s.x[17] = uint64(uint32(s.x[8]) ^ uint32(s.x[11]))
	s.x[11] = 0x6fa6c2c9
	s.x[2] = uint64(uint32(s.x[8]) ^ uint32(s.x[11]))
	s.x[11] = 0x105aa3c2
	s.x[2] = uint64(uint32(s.x[2]) + uint32(s.x[11]))
	s.x[11] = 0x4e2daf9c
	s.x[19] = uint64(uint32(s.x[8]) ^ uint32(s.x[11]))
	s.x[11] = 0xf6fddf6f
	s.x[22] = uint64(uint32(s.x[8]) ^ uint32(s.x[11]))
	s.x[11] = 0x28894dc6
	s.x[17] = uint64(uint32(s.x[17]) + uint32(s.x[11]))
	s.x[8] = uint64(uint32(s.x[8]) ^ uint32(s.x[17]))
	s.x[11] = 0x31d1ce95
	s.x[17] = uint64(uint32(s.x[19]) + uint32(s.x[11]))
	s.x[17] = uint64(uint32(s.x[2]) ^ uint32(s.x[17]))
	s.x[8] = uint64(uint32(s.x[8]) ^ uint32(s.x[17]))
	s.x[11] = 0x8901be68
	s.x[17] = uint64(uint32(s.x[22]) + uint32(s.x[11]))
	s.x[11] = 0x8b1c534a
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[11]))
	s.x[8] = uint64(uint32(s.x[8]) ^ uint32(s.x[17]))
	s.x[19] = uint64(uint32(s.x[21]) ^ uint32(rotateRight(s.x[3], 2, 32)))
	s.x[26]++
	s.x[11] = 0x481b9113
	s.compare64(s.x[26], s.x[11])
	gate := uint32(boolValue(s.condition("eq")))
	s.x[17] = uint64(gate * 0xf4)
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[0]))
	next, err := loadBootstrapPointer(s, s.x[24], uint32(s.x[17]))
	if err != nil {
		return 0, fmt.Errorf("account triple-word round target: %w", err)
	}
	s.x[25] = next
	s.x[3] = s.x[12]
	s.x[2] = s.x[16]
	s.x[22] = 0x481b9113
	s.x[17] = 0xcad9
	return next, nil
}
