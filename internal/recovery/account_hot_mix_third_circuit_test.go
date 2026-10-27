package recovery

import "fmt"

func runAccountHotMixThirdCircuit(s *state) (uint64, error) {
	s.traceBlock(accountHotMixThirdFront)
	s.x[9] = s.x[11]
	s.x[11] = s.x[13]
	s.x[15] = s.x[12]
	s.x[12] = uint64(uint32(rotateRight(s.x[12], 27, 32)))
	s.x[13] = uint64(uint32(s.x[12]) ^ 0x8d0af128)
	s.x[17] = 0x512f646e
	s.x[13] = uint64(uint32(s.x[13]) + uint32(s.x[17]))
	s.x[17] = uint64(uint32(s.x[12]) ^ 0x40b7bf96)
	s.x[2] = 0x9c922ad4
	s.x[17] = uint64(uint32(s.x[17]) + uint32(s.x[2]))
	s.x[13] = uint64(uint32(s.x[13]) ^ uint32(s.x[17]))
	s.x[17] = 0xae44f91f
	s.x[12] = uint64(uint32(s.x[12]) ^ uint32(s.x[17]))
	s.x[17] = 0x72616c5b
	s.x[12] = uint64(uint32(s.x[12]) + uint32(s.x[17]))
	s.x[12] = uint64(uint32(s.x[13]) ^ uint32(s.x[12]))
	s.x[13] = uint64(uint32(s.x[11]) ^ 0x4c15405f)
	s.x[17] = uint64(uint32(s.x[9]) ^ 0xf4e0732b)
	s.x[2] = 0xffffd987
	s.x[4] = uint64(uint32(s.x[16]) + uint32(s.x[2]))
	s.x[17] = uint64(uint32(s.x[4]) ^ uint32(s.x[17]))
	s.x[13] = uint64(uint32(s.x[17]) & uint32(s.x[13]))
	s.x[17] = uint64(uint32(s.x[9]) & 0xb8f5721d)
	s.x[13] = uint64(uint32(s.x[13]) ^ uint32(s.x[17]))
	s.x[17] = uint64(uint32(s.x[19]) ^ 0xe8ee4639)
	s.x[20] = uint64(uint32(s.x[13]) ^ 0x4f1fcdff)
	s.x[17] = uint64(uint32(s.x[20]) & uint32(s.x[17]))
	s.x[13] = uint64(uint32(s.x[13]) & 0xe3f18b84)
	s.x[20] = uint64(uint32(s.x[11]) ^ 0x98b33544)
	s.x[24] = uint64(uint32(s.x[9]) ^ uint32(s.x[5]))
	s.x[20] = uint64(uint32(s.x[24]) & uint32(s.x[20]))
	s.x[24] = uint64(uint32(s.x[9]) & 0x93acf8f9)

	s.traceBlock(accountHotMixThirdTail)
	s.x[20] = uint64(uint32(s.x[20]) ^ uint32(s.x[24]))
	s.x[13] = uint64(uint32(s.x[13]) ^ 0x8e182efc)
	s.x[13] = uint64(uint32(s.x[13]) ^ uint32(s.x[17]))
	s.x[17] = uint64(uint32(s.x[20]) ^ 0xfcf33746)
	s.x[13] = uint64(uint32(s.x[13]) & uint32(s.x[17]))
	s.x[17] = uint64(uint32(s.x[20]) & 0x32f65887)
	s.x[13] = uint64(uint32(s.x[13]) ^ uint32(s.x[17]))
	s.x[17] = uint64(uint32(s.x[13]) ^ 0x98cd4834)
	s.x[2] = 0x3cf8fcc2
	s.x[17] = uint64(uint32(s.x[17]) + uint32(s.x[2]))
	s.x[2] = 0x74497714
	s.x[20] = uint64(uint32(s.x[13]) ^ uint32(s.x[2]))
	s.x[2] = 0xd07cc3e2
	s.x[20] = uint64(uint32(s.x[20]) + uint32(s.x[2]))
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[20]))
	s.x[2] = 0x2389d0d9
	s.x[13] = uint64(uint32(s.x[13]) ^ uint32(s.x[2]))
	s.x[2] = 0x87bc642d
	s.x[13] = uint64(uint32(s.x[13]) + uint32(s.x[2]))
	s.x[13] = uint64(uint32(s.x[17]) ^ uint32(s.x[13]))
	s.x[17] = uint64(uint32(s.x[7]) ^ 0x31aef698)
	s.x[2] = 0xce510968
	s.x[17] = uint64(uint32(s.x[17]) + uint32(s.x[2]))
	s.x[2] = 0x791fa8be
	s.x[20] = uint64(uint32(s.x[7]) ^ uint32(s.x[2]))
	s.x[2] = 0x86e05742
	s.x[20] = uint64(uint32(s.x[20]) + uint32(s.x[2]))
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[20]))
	s.x[2] = 0x43ae939b
	s.x[7] = uint64(uint32(s.x[7]) ^ uint32(s.x[2]))
	s.x[2] = 0xbc516c65
	s.x[7] = uint64(uint32(s.x[7]) + uint32(s.x[2]))
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[7]))
	s.x[7] = s.x[6] + (s.x[23] << 2)
	word, err := s.load(s.x[10]+s.x[7], 4)
	if err != nil {
		return 0, fmt.Errorf("account third hot-mix table word: %w", err)
	}
	s.x[7] = uint64(uint32(word))
	s.x[20] = uint64(uint32(s.x[7]) ^ 0x18f555d6)
	s.x[2] = 0xe70aaa2a
	s.x[20] = uint64(uint32(s.x[20]) + uint32(s.x[2]))
	s.x[24] = uint64(uint32(s.x[7]) ^ 0xa420d2d0)
	s.x[2] = 0x5bdf2d30
	s.x[24] = uint64(uint32(s.x[24]) + uint32(s.x[2]))
	s.x[20] = uint64(uint32(s.x[20]) ^ uint32(s.x[24]))
	s.x[2] = 0x72a0fcc3
	s.x[7] = uint64(uint32(s.x[7]) ^ uint32(s.x[2]))
	s.x[2] = 0x8d5f033d
	s.x[7] = uint64(uint32(s.x[7]) + uint32(s.x[2]))
	s.x[7] = uint64(uint32(s.x[20]) ^ uint32(s.x[7]))
	s.x[12] = uint64(uint32(s.x[12]) + uint32(s.x[17]))
	s.x[12] = uint64(uint32(s.x[12]) + uint32(s.x[7]))
	s.x[17] = 0xfcc1e3f7
	s.x[12] = uint64(uint32(s.x[12]) + uint32(s.x[17]))
	s.x[12] = uint64(uint32(s.x[12]) + uint32(s.x[13]))
	s.x[13] = uint64(uint32(s.x[12]) ^ 0x661ff453)
	s.x[17] = uint64(uint32(s.x[12]) ^ 0x235eb2c7)
	s.x[2] = 0x638416b4
	s.x[17] = uint64(uint32(s.x[17]) + uint32(s.x[2]))
	s.x[2] = 0x059bc6e8
	s.x[7] = uint64(uint32(s.x[12]) ^ uint32(s.x[2]))
	s.x[2] = 0xffffdbf7
	s.x[20] = uint64(uint32(s.x[12]) ^ uint32(s.x[2]))
	s.x[2] = 0x26c55028
	s.x[13] = uint64(uint32(s.x[13]) + uint32(s.x[2]))
	s.x[12] = uint64(uint32(s.x[12]) ^ uint32(s.x[13]))
	s.x[13] = uint64(uint32(s.x[7]) + 0x4541629d)
	s.x[13] = uint64(uint32(s.x[17]) ^ uint32(s.x[13]))
	s.x[12] = uint64(uint32(s.x[12]) ^ uint32(s.x[13]))
	s.x[13] = uint64(uint32(s.x[20]) + 0xbf257f84)
	s.x[17] = 0xb43a9636
	s.x[13] = uint64(uint32(s.x[13]) ^ uint32(s.x[17]))
	s.x[12] = uint64(uint32(s.x[12]) ^ uint32(s.x[13]))
	s.x[13] = uint64(uint32(s.x[30]) ^ uint32(rotateRight(s.x[19], 2, 32)))
	s.x[23]++
	s.x[17] = 0x2ba7c3c7
	s.compare64(s.x[23], s.x[17])
	gate := uint32(boolValue(s.condition("ne")))
	s.x[17] = uint64(gate)
	s.x[7] = uint64(gate << 12)
	s.x[17] = uint64(uint32(s.x[7]) - gate<<1)
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[16]))
	next, err := loadBootstrapPointer(s, s.x[21], uint32(s.x[17]))
	if err != nil {
		return 0, fmt.Errorf("account third hot-mix target: %w", err)
	}
	s.x[24] = next
	s.x[19] = s.x[15]
	s.x[7] = s.x[9]
	s.x[20] = 0x2ba7c3c7
	s.x[17] = 0xfab0
	return next, nil
}
