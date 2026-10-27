package recovery

import "fmt"

func runAccountWordRoundCircuit(s *state) (uint64, error) {
	s.traceBlock(accountWordRoundStart)
	s.x[2] = s.x[10]
	s.x[3] = s.x[8]
	s.x[10] = s.x[19]
	s.x[8] = uint64(uint32(s.x[25]) + 0x2cff82d0)
	s.x[11] = 0xd300ed7f
	s.x[0] = uint64(uint32(s.x[8]) & uint32(s.x[11]))
	s.x[8] = uint64(uint32(rotateRight(s.x[3], 27, 32)))
	s.x[11] = 0xbd84a2dd
	s.x[12] = uint64(uint32(s.x[8]) ^ uint32(s.x[11]))
	s.x[11] = 0x70158451
	s.x[12] = uint64(uint32(s.x[12]) + uint32(s.x[11]))
	s.x[11] = 0xc282af5a
	s.x[16] = uint64(uint32(s.x[8]) ^ uint32(s.x[11]))
	s.x[16] = uint64(uint32(s.x[0]) ^ uint32(s.x[16]))
	s.x[11] = 0x1cff1f5b
	s.x[8] = uint64(uint32(s.x[8]) ^ uint32(s.x[11]))
	s.x[11] = 0xd16e39d7
	s.x[8] = uint64(uint32(s.x[8]) + uint32(s.x[11]))
	s.x[8] = uint64(uint32(s.x[12]) ^ uint32(s.x[8]))
	s.x[11] = 0x0f132cab
	s.x[12] = uint64(uint32(s.x[16]) + uint32(s.x[11]))
	s.x[8] = uint64(uint32(s.x[8]) ^ uint32(s.x[12]))
	s.x[12] = uint64(uint32(s.x[2]) ^ uint32(s.x[19]))
	s.x[12] = uint64(uint32(s.x[12]) ^ uint32(s.x[22]))
	s.x[11] = 0x70b55e49
	s.x[16] = uint64(uint32(s.x[12]) ^ uint32(s.x[11]))
	s.x[11] = 0x940be0bf
	s.x[16] = uint64(uint32(s.x[16]) + uint32(s.x[11]))
	s.x[11] = 0xc76d7759
	s.x[17] = uint64(uint32(s.x[12]) ^ uint32(s.x[11]))
	s.x[11] = 0x23d3c9af
	s.x[17] = uint64(uint32(s.x[17]) + uint32(s.x[11]))
	s.x[16] = uint64(uint32(s.x[16]) ^ uint32(s.x[17]))
	s.x[11] = 0xbcc7e4ad
	s.x[12] = uint64(uint32(s.x[12]) ^ uint32(s.x[11]))
	s.x[11] = 0x58795a5b
	s.x[12] = uint64(uint32(s.x[12]) + uint32(s.x[11]))
	s.x[12] = uint64(uint32(s.x[16]) ^ uint32(s.x[12]))
	s.x[11] = 0x76886c70
	s.x[16] = uint64(uint32(s.x[15]) ^ uint32(s.x[11]))
	s.x[11] = 0x89779390
	s.x[16] = uint64(uint32(s.x[16]) + uint32(s.x[11]))
	s.x[11] = 0x47538645
	s.x[17] = uint64(uint32(s.x[15]) ^ uint32(s.x[11]))
	s.x[11] = 0xb8ac79bb
	s.x[17] = uint64(uint32(s.x[17]) + uint32(s.x[11]))
	s.x[16] = uint64(uint32(s.x[16]) ^ uint32(s.x[17]))
	s.x[11] = 0x3ac42788
	s.x[15] = uint64(uint32(s.x[15]) ^ uint32(s.x[11]))
	s.x[11] = 0xc53bd878
	s.x[15] = uint64(uint32(s.x[15]) + uint32(s.x[11]))
	s.x[15] = uint64(uint32(s.x[16]) ^ uint32(s.x[15]))
	s.x[16] = s.x[9] + (s.x[27] << 2)
	word, err := s.load(s.x[20]+s.x[16], 4)
	if err != nil {
		return 0, fmt.Errorf("account word-round table word: %w", err)
	}
	s.x[16] = uint64(uint32(word))
	s.x[11] = 0x667075aa
	s.x[17] = uint64(uint32(s.x[16]) ^ uint32(s.x[11]))
	s.x[11] = 0x998f8a56
	s.x[17] = uint64(uint32(s.x[17]) + uint32(s.x[11]))
	s.x[11] = 0x76bdad9a
	s.x[19] = uint64(uint32(s.x[16]) ^ uint32(s.x[11]))
	s.x[11] = 0x89425266
	s.x[19] = uint64(uint32(s.x[19]) + uint32(s.x[11]))
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[19]))
	s.x[11] = 0xdeb8a3f5
	s.x[16] = uint64(uint32(s.x[16]) ^ uint32(s.x[11]))
	s.x[11] = 0x21475c0b
	s.x[16] = uint64(uint32(s.x[16]) + uint32(s.x[11]))
	s.x[16] = uint64(uint32(s.x[17]) ^ uint32(s.x[16]))
	s.x[8] = uint64(uint32(s.x[8]) + uint32(s.x[15]))
	s.x[8] = uint64(uint32(s.x[8]) + uint32(s.x[12]))
	s.x[11] = 0x43bb4b41
	s.x[8] = uint64(uint32(s.x[8]) + uint32(s.x[11]))
	s.x[8] = uint64(uint32(s.x[8]) + uint32(s.x[16]))
	s.x[11] = 0xbe9e0c8d
	s.x[12] = uint64(uint32(s.x[8]) ^ uint32(s.x[11]))
	s.x[11] = 0xcec2562d
	s.x[15] = uint64(uint32(s.x[8]) ^ uint32(s.x[11]))
	s.x[11] = 0xa86bb345
	s.x[15] = uint64(uint32(s.x[15]) + uint32(s.x[11]))
	s.x[11] = 0x37f59fcb
	s.x[16] = uint64(uint32(s.x[8]) ^ uint32(s.x[11]))
	s.x[11] = 0xdeffdffd
	s.x[17] = uint64(uint32(s.x[8]) ^ uint32(s.x[11]))
	s.x[11] = 0xd837e9e5
	s.x[12] = uint64(uint32(s.x[12]) + uint32(s.x[11]))
	s.x[8] = uint64(uint32(s.x[8]) ^ uint32(s.x[12]))
	s.x[11] = 0x515c7aa3
	return accountWordRoundTail, nil
}
