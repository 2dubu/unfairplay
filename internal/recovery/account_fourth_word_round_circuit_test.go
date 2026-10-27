package recovery

import "fmt"

func runAccountFourthWordRoundCircuit(s *state) (uint64, error) {
	s.traceBlock(accountFourthWordRound)
	s.x[2] = s.x[16]
	s.x[17] = s.x[12]
	s.x[11] = 0x6eb9aea6
	s.x[3] = uint64(uint32(s.x[16]) ^ uint32(s.x[11]))
	s.x[16] = s.x[10]
	s.x[10] = s.x[19]
	s.x[12] = uint64(uint32(s.x[16]) ^ uint32(s.x[19]))
	s.x[19] = uint64(uint32(s.x[12]) ^ uint32(s.x[17]))
	s.x[12] = s.x[8]
	s.x[8] = uint64(uint32(s.x[15]) - 0xd33)
	s.x[11] = 0x422
	s.x[0] = uint64(uint32(s.x[8]) | uint32(s.x[11]))
	s.x[8] = uint64(uint32(rotateRight(s.x[12], 27, 32)))
	s.x[11] = 0xabdc3b00
	s.x[25] = uint64(uint32(s.x[8]) ^ uint32(s.x[11]))
	s.x[11] = 0x022127f4
	s.x[25] = uint64(uint32(s.x[25]) + uint32(s.x[11]))
	s.x[11] = 0x025ed594
	s.x[26] = uint64(uint32(s.x[8]) ^ uint32(s.x[11]))
	s.x[11] = 0xaba3c968
	s.x[26] = uint64(uint32(s.x[26]) + uint32(s.x[11]))
	s.x[25] = uint64(uint32(s.x[25]) ^ uint32(s.x[26]))
	s.x[11] = 0xca7b5935
	s.x[8] = uint64(uint32(s.x[8]) ^ uint32(s.x[11]))
	s.x[11] = 0x638645c7
	s.x[8] = uint64(uint32(s.x[8]) + uint32(s.x[11]))
	s.x[8] = uint64(uint32(s.x[25]) ^ uint32(s.x[8]))
	s.x[11] = 0x81e0b547
	s.x[25] = uint64(uint32(s.x[19]) ^ uint32(s.x[11]))
	s.x[11] = 0xa10b7eee
	s.x[26] = uint64(uint32(s.x[0]) + uint32(s.x[11]))
	s.x[25] = uint64(uint32(s.x[26]) + uint32(s.x[25]))
	s.x[11] = 0x71bd3438
	s.x[26] = uint64(uint32(s.x[19]) ^ uint32(s.x[11]))
	s.x[11] = 0x51567d21
	s.x[26] = uint64(uint32(s.x[26]) + uint32(s.x[11]))
	s.x[25] = uint64(uint32(s.x[25]) ^ uint32(s.x[26]))
	s.x[11] = 0x55310392
	s.x[26] = uint64(uint32(s.x[2]) ^ uint32(s.x[11]))
	s.x[11] = 0x30976089
	s.x[2] = uint64(uint32(s.x[2]) ^ uint32(s.x[11]))
	s.x[11] = 0xfb424cc2
	s.x[19] = uint64(uint32(s.x[19]) ^ uint32(s.x[11]))
	s.x[11] = 0xdba905db
	s.x[19] = uint64(uint32(s.x[19]) + uint32(s.x[11]))
	s.x[19] = uint64(uint32(s.x[25]) ^ uint32(s.x[19]))
	s.x[11] = 0x9146515a
	s.x[3] = uint64(uint32(s.x[3]) + uint32(s.x[11]))
	s.x[11] = 0xaacefc6e
	s.x[25] = uint64(uint32(s.x[26]) + uint32(s.x[11]))
	s.x[3] = uint64(uint32(s.x[3]) ^ uint32(s.x[25]))
	s.x[11] = 0xcf689f77
	s.x[2] = uint64(uint32(s.x[2]) + uint32(s.x[11]))
	s.x[2] = uint64(uint32(s.x[3]) ^ uint32(s.x[2]))
	s.x[3] = s.x[9] + (s.x[22] << 2)
	word, err := s.load(s.x[20]+s.x[3], 4)
	if err != nil {
		return 0, fmt.Errorf("account fourth-word round table word: %w", err)
	}
	s.x[3] = uint64(uint32(word))
	s.x[11] = 0x95d69bda
	s.x[25] = uint64(uint32(s.x[3]) ^ uint32(s.x[11]))
	s.x[11] = 0x6a296426
	s.x[25] = uint64(uint32(s.x[25]) + uint32(s.x[11]))
	s.x[11] = 0x59a1614a
	s.x[26] = uint64(uint32(s.x[3]) ^ uint32(s.x[11]))
	s.x[11] = 0xa65e9eb6
	s.x[26] = uint64(uint32(s.x[26]) + uint32(s.x[11]))
	s.x[25] = uint64(uint32(s.x[25]) ^ uint32(s.x[26]))
	s.x[11] = 0x02028155
	s.x[3] = uint64(uint32(s.x[3]) ^ uint32(s.x[11]))
	s.x[11] = 0xfdfd7eab
	s.x[3] = uint64(uint32(s.x[3]) + uint32(s.x[11]))
	s.x[3] = uint64(uint32(s.x[25]) ^ uint32(s.x[3]))
	s.x[8] = uint64(uint32(s.x[8]) + uint32(s.x[2]))
	s.x[8] = uint64(uint32(s.x[8]) + uint32(s.x[19]))
	s.x[11] = 0x74a9e021
	s.x[8] = uint64(uint32(s.x[8]) + uint32(s.x[11]))
	s.x[8] = uint64(uint32(s.x[8]) + uint32(s.x[3]))
	s.x[11] = 0x03003570
	s.x[2] = uint64(uint32(s.x[8]) ^ uint32(s.x[11]))
	s.x[11] = 0x3a54ce4e
	s.x[2] = uint64(uint32(s.x[2]) + uint32(s.x[11]))
	s.x[11] = 0xa623b9a4
	s.x[3] = uint64(uint32(s.x[8]) ^ uint32(s.x[11]))
	s.x[11] = 0x9f77429a
	s.x[3] = uint64(uint32(s.x[3]) + uint32(s.x[11]))
	s.x[2] = uint64(uint32(s.x[2]) ^ uint32(s.x[3]))
	s.x[11] = 0x1c7f676d
	s.x[3] = uint64(uint32(s.x[8]) ^ uint32(s.x[11]))
	s.x[11] = 0x252b9c51
	s.x[3] = uint64(uint32(s.x[3]) + uint32(s.x[11]))
	s.x[3] = uint64(uint32(s.x[8]) ^ uint32(s.x[3]))
	s.x[2] = uint64(uint32(s.x[3]) ^ uint32(s.x[2]))
	s.x[19] = uint64(uint32(s.x[21]) ^ uint32(rotateRight(s.x[17], 2, 32)))
	s.x[11] = 0x7ff7ef7b
	s.x[8] = uint64(uint32(s.x[8]) ^ uint32(s.x[11]))
	s.x[11] = 0x46a31447
	s.x[8] = uint64(uint32(s.x[8]) + uint32(s.x[11]))
	s.x[11] = 0xcdb4c97f
	s.x[8] = uint64(uint32(s.x[8]) ^ uint32(s.x[11]))
	s.x[8] = uint64(uint32(s.x[2]) ^ uint32(s.x[8]))
	s.x[22]++
	s.x[11] = 0x481b9127
	s.compare64(s.x[22], s.x[11])
	gate := uint32(boolValue(s.condition("eq")))
	s.x[17] = uint64(gate * 0x13)
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[15]))
	next, err := loadBootstrapPointer(s, s.x[24], uint32(s.x[17]))
	if err != nil {
		return 0, fmt.Errorf("account fourth-word round target: %w", err)
	}
	s.x[2] = next
	s.x[17] = 0x7ab5
	return next, nil
}
