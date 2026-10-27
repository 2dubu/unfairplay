package recovery

import "fmt"

const accountHotMixFourth = uint64(0x100088c14)

func runAccountHotMixFourth(s *state) (uint64, error) {
	s.traceBlock(accountHotMixFourth)
	s.x[17] = s.x[9]
	s.x[9] = 0x50c9d63e
	s.x[19] = uint64(uint32(s.x[17]) ^ uint32(s.x[9]))
	s.x[9] = s.x[11]
	s.x[11] = s.x[13]
	s.x[16] = s.x[12]
	s.x[12] = uint64(uint32(rotateRight(s.x[12], 27, 32)))
	s.x[13] = uint64(uint32(s.x[12]) ^ 0x832b95db)
	s.x[2] = 0x85fb187e
	s.x[13] = uint64(uint32(s.x[13]) + uint32(s.x[2]))
	s.x[7] = uint64(uint32(s.x[12]) ^ 0xffe58a6a)
	s.x[2] = 0xf93507cd
	s.x[7] = uint64(uint32(s.x[7]) + uint32(s.x[2]))
	s.x[13] = uint64(uint32(s.x[13]) ^ uint32(s.x[7]))
	s.x[12] = uint64(uint32(s.x[12]) ^ 0x1f37a810)
	s.x[7] = uint64(uint32(s.x[4]) + 0x82831d88)
	s.x[23] = uint64(uint32(s.x[7]) + 0x9763c6c6)
	s.x[12] = uint64(uint32(s.x[23]) + uint32(s.x[12]))
	s.x[12] = uint64(uint32(s.x[13]) ^ uint32(s.x[12]))
	s.x[13] = uint64(uint32(s.x[11]) ^ uint32(s.x[15]))
	s.x[13] = uint64(uint32(s.x[13]) ^ uint32(s.x[9]))
	s.x[23] = uint64(uint32(s.x[13]) ^ 0x7ff86437)
	s.x[2] = 0xcedca140
	s.x[23] = uint64(uint32(s.x[23]) + uint32(s.x[2]))
	s.x[24] = uint64(uint32(s.x[13]) ^ 0x00a6528a)
	s.x[2] = 0xb1829783
	s.x[24] = uint64(uint32(s.x[24]) + uint32(s.x[2]))
	s.x[23] = uint64(uint32(s.x[23]) ^ uint32(s.x[24]))
	s.x[24] = uint64(uint32(s.x[17]) ^ 0x78b58bc7)
	s.x[17] = uint64(uint32(s.x[17]) ^ 0x23639044)
	s.x[13] = uint64(uint32(s.x[13]) ^ 0x7441fb00)
	s.x[2] = 0xc5653e09
	s.x[13] = uint64(uint32(s.x[13]) + uint32(s.x[2]))
	s.x[13] = uint64(uint32(s.x[23]) ^ uint32(s.x[13]))
	s.x[19] = uint64(uint32(s.x[19]) + 0xaf3629c2)
	s.x[23] = uint64(uint32(s.x[24]) + 0x874a7439)
	s.x[19] = uint64(uint32(s.x[19]) ^ uint32(s.x[23]))
	s.x[17] = uint64(uint32(s.x[17]) + 0xdc9c6fbc)
	s.x[17] = uint64(uint32(s.x[19]) ^ uint32(s.x[17]))
	s.x[19] = s.x[6] + (s.x[20] << 2)
	word, err := s.load(s.x[10]+s.x[19], 4)
	if err != nil {
		return 0, fmt.Errorf("account fourth hot-mix table word: %w", err)
	}
	s.x[19] = uint64(uint32(word))
	s.x[23] = uint64(uint32(s.x[19]) ^ 0x4378c660)
	s.x[2] = 0xbc8739a0
	s.x[23] = uint64(uint32(s.x[23]) + uint32(s.x[2]))
	s.x[24] = uint64(uint32(s.x[19]) ^ 0xe0ddc8de)
	s.x[2] = 0x1f223722
	s.x[24] = uint64(uint32(s.x[24]) + uint32(s.x[2]))
	s.x[23] = uint64(uint32(s.x[23]) ^ uint32(s.x[24]))
	s.x[19] = uint64(uint32(s.x[19]) ^ 0x6dd0757b)
	s.x[2] = 0x922f8a85
	s.x[19] = uint64(uint32(s.x[19]) + uint32(s.x[2]))
	s.x[19] = uint64(uint32(s.x[23]) ^ uint32(s.x[19]))
	s.x[12] = uint64(uint32(s.x[12]) + uint32(s.x[17]))
	s.x[12] = uint64(uint32(s.x[12]) + uint32(s.x[13]))
	s.x[13] = 0x43498dda
	s.x[12] = uint64(uint32(s.x[12]) + uint32(s.x[13]))
	s.x[12] = uint64(uint32(s.x[12]) + uint32(s.x[19]))
	s.x[13] = uint64(uint32(s.x[12]) ^ 0x205d0af7)
	s.x[13] = uint64(uint32(s.x[13]) + 0x61173049)
	s.x[17] = uint64(uint32(s.x[12]) ^ 0xeec07b79)
	s.x[2] = 0xaf8a41c7
	s.x[17] = uint64(uint32(s.x[17]) + uint32(s.x[2]))
	s.x[13] = uint64(uint32(s.x[13]) ^ uint32(s.x[17]))
	s.x[17] = uint64(uint32(s.x[12]) ^ 0x4fd77a31)
	s.x[2] = 0x0e9d408f
	s.x[17] = uint64(uint32(s.x[17]) + uint32(s.x[2]))
	s.x[17] = uint64(uint32(s.x[12]) ^ uint32(s.x[17]))
	s.x[2] = 0x3fffceff
	s.x[12] = uint64(uint32(s.x[12]) ^ uint32(s.x[2]))
	s.x[13] = uint64(uint32(s.x[17]) ^ uint32(s.x[13]))
	s.x[17] = 0x7eb5f441
	s.x[12] = uint64(uint32(s.x[12]) + uint32(s.x[17]))
	s.x[17] = 0xb5aa08fd
	s.x[12] = uint64(uint32(s.x[12]) ^ uint32(s.x[17]))
	s.x[12] = uint64(uint32(s.x[13]) ^ uint32(s.x[12]))
	s.x[13] = uint64(uint32(s.x[30]) ^ uint32(rotateRight(s.x[15], 2, 32)))
	s.x[20]++
	s.x[15] = 0x2ba7c3db
	s.compare64(s.x[20], s.x[15])
	gate := uint32(boolValue(s.condition("eq")))
	s.x[15] = uint64(gate * 0x13)
	s.x[15] = uint64(uint32(s.x[15]) ^ uint32(s.x[4]))
	next, err := loadBootstrapPointer(s, s.x[21], uint32(s.x[15]))
	if err != nil {
		return 0, fmt.Errorf("account fourth hot-mix target: %w", err)
	}
	s.x[19] = next
	s.x[15] = s.x[16]
	s.x[17] = 0xfb38
	return next, nil
}
