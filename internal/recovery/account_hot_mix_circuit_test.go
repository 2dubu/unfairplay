package recovery

import "fmt"

func runAccountHotMixCircuit(s *state) (uint64, error) {
	s.traceBlock(accountHotMixFront)
	s.x[4] = s.x[11]
	s.traceBlock(accountHotMixTail)
	s.x[11] = s.x[13]
	s.x[15] = s.x[12]
	s.x[12] = uint64(uint32(rotateRight(s.x[12], 27, 32)))
	s.x[13] = uint64(uint32(s.x[12]) ^ 0x7d1595f1)
	s.x[17] = 0xcc7a4f97
	s.x[13] = uint64(uint32(s.x[13]) + uint32(s.x[17]))
	s.x[17] = uint64(uint32(s.x[12]) ^ 0x0c015332)
	s.x[2] = 0xbd6e8956
	s.x[17] = uint64(uint32(s.x[17]) + uint32(s.x[2]))
	s.x[13] = uint64(uint32(s.x[13]) ^ uint32(s.x[17]))
	s.x[17] = 0x12ed7162
	s.x[12] = uint64(uint32(s.x[12]) ^ uint32(s.x[17]))
	s.x[17] = 0xa382ab06
	s.x[12] = uint64(uint32(s.x[12]) + uint32(s.x[17]))
	s.x[12] = uint64(uint32(s.x[13]) ^ uint32(s.x[12]))
	s.x[13] = uint64(uint32(s.x[7]) ^ 0x92596cde)
	s.x[17] = uint64(uint32(s.x[11]) ^ uint32(s.x[5]))
	s.x[13] = uint64(uint32(s.x[17]) & uint32(s.x[13]))
	s.x[17] = uint64(uint32(s.x[11]) & 0x9946a163)
	s.x[2] = 0xa3f799e2
	s.x[20] = uint64(uint32(s.x[4]) ^ uint32(s.x[2]))
	s.x[2] = 0xf4e03242
	s.x[23] = uint64(uint32(s.x[7]) ^ uint32(s.x[2]))
	s.x[20] = uint64(uint32(s.x[20]) & uint32(s.x[23]))
	s.x[2] = 0xa8e8545f
	s.x[23] = uint64(uint32(s.x[7]) & uint32(s.x[2]))
	s.x[20] = uint64(uint32(s.x[20]) ^ uint32(s.x[23]))
	s.x[2] = 0x1fc92292
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[2]))
	s.x[13] = uint64(uint32(s.x[17]) ^ uint32(s.x[13]))
	s.x[17] = uint64(uint32(s.x[20]) ^ 0x5f1fefbd)
	s.x[13] = uint64(uint32(s.x[17]) & uint32(s.x[13]))
	s.x[17] = uint64(uint32(s.x[20]) & 0xe9305c4c)
	s.x[13] = uint64(uint32(s.x[13]) ^ uint32(s.x[17]))
	s.x[17] = uint64(uint32(s.x[13]) ^ 0x53e72a39)
	s.x[2] = 0x56f1ca0e
	s.x[17] = uint64(uint32(s.x[17]) + uint32(s.x[2]))
	s.x[2] = 0x3673fed6
	s.x[20] = uint64(uint32(s.x[13]) ^ uint32(s.x[2]))
	s.x[2] = 0x33651ee3
	s.x[20] = uint64(uint32(s.x[20]) + uint32(s.x[2]))
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[20]))
	s.x[2] = 0xd37b671c
	s.x[13] = uint64(uint32(s.x[13]) ^ uint32(s.x[2]))
	s.x[2] = 0xd66d8729
	s.x[13] = uint64(uint32(s.x[13]) + uint32(s.x[2]))
	s.x[13] = uint64(uint32(s.x[17]) ^ uint32(s.x[13]))
	s.x[17] = uint64(uint32(s.x[19]) ^ 0x38b0f973)
	s.x[2] = 0xc74f068d
	s.x[17] = uint64(uint32(s.x[17]) + uint32(s.x[2]))
	s.x[2] = 0xc5ec8b23
	s.x[20] = uint64(uint32(s.x[19]) ^ uint32(s.x[2]))
	s.x[2] = 0x3a1374dd
	s.x[20] = uint64(uint32(s.x[20]) + uint32(s.x[2]))
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[20]))
	s.x[2] = 0xf643bfed
	s.x[19] = uint64(uint32(s.x[19]) ^ uint32(s.x[2]))
	s.x[2] = 0x09bc4013
	s.x[19] = uint64(uint32(s.x[19]) + uint32(s.x[2]))
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[19]))
	s.x[19] = s.x[6] + (s.x[16] << 2)
	word, err := s.load(s.x[10]+s.x[19], 4)
	if err != nil {
		return 0, fmt.Errorf("account hot-mix table word: %w", err)
	}
	s.x[19] = uint64(uint32(word))
	s.x[20] = uint64(uint32(s.x[19]) ^ 0xd3dddaa9)
	s.x[2] = 0x2c222557
	s.x[20] = uint64(uint32(s.x[20]) + uint32(s.x[2]))
	s.x[23] = uint64(uint32(s.x[19]) ^ 0x45602a6c)
	s.x[2] = 0xba9fd594
	s.x[23] = uint64(uint32(s.x[23]) + uint32(s.x[2]))
	s.x[20] = uint64(uint32(s.x[20]) ^ uint32(s.x[23]))
	s.x[19] = uint64(uint32(s.x[19]) ^ 0x58c88b00)
	s.x[2] = 0xa7377500
	s.x[19] = uint64(uint32(s.x[19]) + uint32(s.x[2]))
	s.x[19] = uint64(uint32(s.x[20]) ^ uint32(s.x[19]))
	s.x[12] = uint64(uint32(s.x[12]) + uint32(s.x[17]))
	s.x[12] = uint64(uint32(s.x[12]) + uint32(s.x[13]))
	s.x[12] = uint64(uint32(s.x[12]) + uint32(s.x[19]))
	s.x[13] = 0x11cb71c7
	s.x[12] = uint64(uint32(s.x[12]) + uint32(s.x[13]))
	s.x[13] = uint64(uint32(s.x[12]) ^ 0xffe90e88)
	s.x[17] = uint64(uint32(s.x[12]) ^ 0x45d023ce)
	s.x[2] = 0xde6cac0b
	s.x[17] = uint64(uint32(s.x[17]) + uint32(s.x[2]))
	s.x[2] = 0x35c0a282
	s.x[19] = uint64(uint32(s.x[12]) ^ uint32(s.x[2]))
	s.x[2] = 0xebbaffff
	s.x[20] = uint64(uint32(s.x[12]) ^ uint32(s.x[2]))
	s.x[2] = 0x6455814d
	s.x[13] = uint64(uint32(s.x[13]) + uint32(s.x[2]))
	s.x[12] = uint64(uint32(s.x[12]) ^ uint32(s.x[13]))
	s.x[13] = uint64(uint32(s.x[19]) + 0xae7c2d47)
	s.x[13] = uint64(uint32(s.x[17]) ^ uint32(s.x[13]))
	s.x[12] = uint64(uint32(s.x[12]) ^ uint32(s.x[13]))
	s.x[13] = uint64(uint32(s.x[20]) + 0x7006703c)
	s.x[17] = 0x6f5cbd86
	s.x[13] = uint64(uint32(s.x[13]) ^ uint32(s.x[17]))
	s.x[12] = uint64(uint32(s.x[12]) ^ uint32(s.x[13]))
	s.x[13] = uint64(uint32(s.x[30]) ^ uint32(rotateRight(s.x[7], 2, 32)))
	s.x[16]++
	s.x[17] = 0x2ba7c39f
	s.compare64(s.x[16], s.x[17])
	gate := uint32(boolValue(s.condition("ne")))
	s.x[17] = uint64(gate)
	s.x[7] = uint64(gate << 4)
	s.x[17] = uint64(uint32(s.x[7]) - gate<<1)
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[9]))
	next, err := loadBootstrapPointer(s, s.x[21], uint32(s.x[17]))
	if err != nil {
		return 0, fmt.Errorf("account hot-mix target: %w", err)
	}
	s.x[23] = next
	s.x[7] = s.x[15]
	s.x[19] = s.x[4]
	s.x[20] = 0x2ba7c39f
	s.x[17] = 0x1d92
	return next, nil
}
