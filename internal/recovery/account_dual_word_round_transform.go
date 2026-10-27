package recovery

import "fmt"

const (
	accountDualWordRoundFront = uint64(0x1001c2f34)
	accountDualWordRoundTail  = uint64(0x1001c3000)
)

func runAccountDualWordRound(s *state) (uint64, error) {
	s.traceBlock(accountDualWordRoundFront)
	s.x[15] = s.x[10]
	s.x[22] = s.x[8]
	s.x[10] = s.x[19]
	s.x[8] = uint64(uint32(rotateRight(s.x[8], 27, 32)))
	s.x[11] = 0x02b671db
	s.x[17] = uint64(uint32(s.x[0]) + uint32(s.x[11]))
	s.x[17] = uint64(uint32(s.x[8]) ^ uint32(s.x[17]))
	s.x[11] = 0x302c8b99
	s.x[17] = uint64(uint32(s.x[17]) + uint32(s.x[11]))
	s.x[11] = 0x456f7aa5
	s.x[3] = uint64(uint32(s.x[8]) ^ uint32(s.x[11]))
	s.x[11] = 0x77f57be6
	s.x[3] = uint64(uint32(s.x[3]) + uint32(s.x[11]))
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[3]))
	s.x[11] = 0x242047dc
	s.x[8] = uint64(uint32(s.x[8]) ^ uint32(s.x[11]))
	s.x[11] = 0x16ba469d
	s.x[8] = uint64(uint32(s.x[8]) + uint32(s.x[11]))
	s.x[8] = uint64(uint32(s.x[17]) ^ uint32(s.x[8]))
	s.x[11] = 0x8eec14d0
	s.x[17] = uint64(uint32(s.x[12]) ^ uint32(s.x[11]))
	s.x[3] = uint64(uint32(s.x[19]) ^ uint32(s.x[23]))
	s.x[17] = uint64(uint32(s.x[17]) & uint32(s.x[3]))
	s.x[11] = 0x85f3d96d
	s.x[3] = uint64(uint32(s.x[19]) & uint32(s.x[11]))
	s.x[11] = 0x389c3d76
	s.x[19] = uint64(uint32(s.x[15]) ^ uint32(s.x[11]))
	s.x[11] = 0xf4e03242
	s.x[25] = uint64(uint32(s.x[12]) ^ uint32(s.x[11]))
	s.x[19] = uint64(uint32(s.x[25]) & uint32(s.x[19]))
	s.x[11] = 0x3383f0cb
	s.x[25] = uint64(uint32(s.x[12]) & uint32(s.x[11]))
	s.x[19] = uint64(uint32(s.x[19]) ^ uint32(s.x[25]))
	s.x[11] = 0x09061b2d
	s.x[3] = uint64(uint32(s.x[3]) ^ uint32(s.x[11]))
	s.x[17] = uint64(uint32(s.x[3]) ^ uint32(s.x[17]))
	s.x[11] = 0xcf7fcfbd
	s.x[3] = uint64(uint32(s.x[19]) ^ uint32(s.x[11]))

	s.traceBlock(accountDualWordRoundTail)
	s.x[17] = uint64(uint32(s.x[3]) & uint32(s.x[17]))
	s.x[11] = 0xf7ea2dff
	s.x[3] = uint64(uint32(s.x[19]) & uint32(s.x[11]))
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[3]))
	s.x[11] = 0xbe568618
	s.x[3] = uint64(uint32(s.x[17]) ^ uint32(s.x[11]))
	s.x[11] = 0xfe0cd165
	s.x[3] = uint64(uint32(s.x[3]) + uint32(s.x[11]))
	s.x[11] = 0xfc43bc19
	s.x[19] = uint64(uint32(s.x[17]) ^ uint32(s.x[11]))
	s.x[11] = 0xbc19eb66
	s.x[19] = uint64(uint32(s.x[19]) + uint32(s.x[11]))
	s.x[3] = uint64(uint32(s.x[3]) ^ uint32(s.x[19]))
	s.x[11] = 0x7a80c843
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[11]))
	s.x[11] = 0x3ada9f40
	s.x[17] = uint64(uint32(s.x[17]) + uint32(s.x[11]))
	s.x[17] = uint64(uint32(s.x[3]) ^ uint32(s.x[17]))
	s.x[11] = 0x89008eba
	s.x[3] = uint64(uint32(s.x[16]) ^ uint32(s.x[11]))
	s.x[11] = 0x76ff7146
	s.x[3] = uint64(uint32(s.x[3]) + uint32(s.x[11]))
	s.x[11] = 0x5d56faa7
	s.x[19] = uint64(uint32(s.x[16]) ^ uint32(s.x[11]))
	s.x[11] = 0xa2a90559
	s.x[19] = uint64(uint32(s.x[19]) + uint32(s.x[11]))
	s.x[3] = uint64(uint32(s.x[3]) ^ uint32(s.x[19]))
	s.x[11] = 0xdf49b9a0
	s.x[16] = uint64(uint32(s.x[16]) ^ uint32(s.x[11]))
	s.x[11] = 0x20b64660
	s.x[16] = uint64(uint32(s.x[16]) + uint32(s.x[11]))
	s.x[16] = uint64(uint32(s.x[3]) ^ uint32(s.x[16]))
	s.x[3] = s.x[9] + (s.x[2] << 2)
	word, err := s.load(s.x[20]+s.x[3], 4)
	if err != nil {
		return 0, fmt.Errorf("account dual-word round table word: %w", err)
	}
	s.x[3] = uint64(uint32(word))
	s.x[11] = 0x230cc73d
	s.x[19] = uint64(uint32(s.x[3]) ^ uint32(s.x[11]))
	s.x[11] = 0xdcf338c3
	s.x[19] = uint64(uint32(s.x[19]) + uint32(s.x[11]))
	s.x[11] = 0xd3ffec52
	s.x[25] = uint64(uint32(s.x[3]) ^ uint32(s.x[11]))
	s.x[11] = 0x2c0013ae
	s.x[25] = uint64(uint32(s.x[25]) + uint32(s.x[11]))
	s.x[19] = uint64(uint32(s.x[19]) ^ uint32(s.x[25]))
	s.x[11] = 0x1c30
	s.x[25] = uint64(uint32(s.x[0]) + uint32(s.x[11]))
	s.x[11] = 0x3e8650aa
	s.x[3] = uint64(uint32(s.x[3]) ^ uint32(s.x[11]))
	s.x[11] = 0xc179af56
	s.x[3] = uint64(uint32(s.x[3]) + uint32(s.x[11]))
	s.x[3] = uint64(uint32(s.x[19]) ^ uint32(s.x[3]))
	s.x[8] = uint64(uint32(s.x[8]) + uint32(s.x[16]))
	s.x[8] = uint64(uint32(s.x[8]) + uint32(s.x[17]))
	s.x[8] = uint64(uint32(s.x[8]) + uint32(s.x[3]))
	s.x[11] = 0xe317c21f
	s.x[8] = uint64(uint32(s.x[8]) + uint32(s.x[11]))
	s.x[11] = 0x751d6dda
	s.x[16] = uint64(uint32(s.x[8]) ^ uint32(s.x[11]))
	s.x[11] = 0x606e5658
	s.x[17] = uint64(uint32(s.x[8]) ^ uint32(s.x[11]))
	s.x[11] = 0xe6a2f283
	s.x[17] = uint64(uint32(s.x[17]) + uint32(s.x[11]))
	s.x[11] = 0x03b7dd59
	s.x[3] = uint64(uint32(s.x[8]) ^ uint32(s.x[11]))
	s.x[11] = 0x6ff7bdfe
	s.x[19] = uint64(uint32(s.x[8]) ^ uint32(s.x[11]))
	s.x[11] = 0xf3d1c901
	s.x[16] = uint64(uint32(s.x[16]) + uint32(s.x[11]))
	s.x[8] = uint64(uint32(s.x[8]) ^ uint32(s.x[16]))
	s.x[11] = 0x857b7984
	s.x[16] = uint64(uint32(s.x[3]) + uint32(s.x[11]))
	s.x[16] = uint64(uint32(s.x[17]) ^ uint32(s.x[16]))
	s.x[8] = uint64(uint32(s.x[8]) ^ uint32(s.x[16]))
	s.x[11] = 0xe93b1925
	s.x[16] = uint64(uint32(s.x[19]) + uint32(s.x[11]))
	s.x[11] = 0x722c9698
	s.x[16] = uint64(uint32(s.x[16]) ^ uint32(s.x[11]))
	s.x[8] = uint64(uint32(s.x[8]) ^ uint32(s.x[16]))
	s.x[19] = uint64(uint32(s.x[21]) ^ uint32(rotateRight(s.x[12], 2, 32)))
	s.x[2]++
	s.x[11] = 0x481b90eb
	s.compare64(s.x[2], s.x[11])
	gate := uint32(boolValue(s.condition("eq")))
	s.x[12] = uint64(gate * 0x2dd0)
	s.x[12] = uint64(uint32(s.x[12]) ^ uint32(s.x[0]))
	next, err := loadBootstrapPointer(s, s.x[24], uint32(s.x[12]))
	if err != nil {
		return 0, fmt.Errorf("account dual-word round target: %w", err)
	}
	s.x[3] = next
	s.x[12] = s.x[22]
	s.x[16] = s.x[15]
	s.x[27] = 0x481b90eb
	s.x[17] = 0xaaed
	return next, nil
}
