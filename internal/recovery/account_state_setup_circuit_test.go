package recovery

import "fmt"

func runAccountStateSetupTransformCircuit(s *state) (uint64, error) {
	s.traceBlock(0x1001a8524)
	s.x[8] = 0x5abc
	s.x[8] = uint64(uint32(insertBits(uint64(uint32(s.x[8])), 0x77f7, 16, 16, 32)))
	s.x[8] = uint64(uint32(s.x[13]) + uint32(s.x[8]))
	s.x[9] = 0x3f54
	s.x[9] = uint64(uint32(insertBits(uint64(uint32(s.x[9])), 0xde99, 16, 16, 32)))
	s.x[8] = uint64(uint32(s.x[8]) & uint32(s.x[9]))
	value6, err := s.load(s.sp+0xf4, 4)
	if err != nil {
		return 0, fmt.Errorf("account state-setup load 0x1001a853c: %w", err)
	}
	s.x[16] = uint64(uint32(value6))
	s.x[9] = 0xd62d
	s.x[9] = uint64(uint32(insertBits(uint64(uint32(s.x[9])), 0x1e52, 16, 16, 32)))
	s.x[17] = uint64(uint32(s.x[16]) ^ uint32(s.x[9]))
	s.x[9] = 0x29d3
	s.x[9] = uint64(uint32(insertBits(uint64(uint32(s.x[9])), 0xe1ad, 16, 16, 32)))
	s.x[17] = uint64(uint32(s.x[17]) + uint32(s.x[9]))
	s.x[9] = 0xede
	s.x[9] = uint64(uint32(insertBits(uint64(uint32(s.x[9])), 0x4ab5, 16, 16, 32)))
	s.x[0] = uint64(uint32(s.x[16]) ^ uint32(s.x[9]))
	s.x[9] = 0xf122
	s.x[9] = uint64(uint32(insertBits(uint64(uint32(s.x[9])), 0xb54a, 16, 16, 32)))
	s.x[0] = uint64(uint32(s.x[0]) + uint32(s.x[9]))
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[0]))
	s.x[8] = uint64(uint32(s.x[16]) ^ uint32(s.x[8]))
	s.x[9] = 0x291e
	s.x[9] = uint64(uint32(insertBits(uint64(uint32(s.x[9])), 0x969, 16, 16, 32)))
	s.x[8] = uint64(uint32(s.x[8]) ^ uint32(s.x[9]))
	s.x[9] = 0xeab2
	s.x[9] = uint64(uint32(insertBits(uint64(uint32(s.x[9])), 0xa007, 16, 16, 32)))
	s.x[8] = uint64(uint32(s.x[8]) + uint32(s.x[9]))
	s.x[8] = uint64(uint32(s.x[17]) ^ uint32(s.x[8]))
	s.x[9] = 0xe51b
	s.x[9] = uint64(uint32(insertBits(uint64(uint32(s.x[9])), 0x5cca, 16, 16, 32)))
	s.x[16] = uint64(uint32(s.x[2]) ^ uint32(s.x[9]))
	s.x[9] = 0x1ae5
	s.x[9] = uint64(uint32(insertBits(uint64(uint32(s.x[9])), 0xa335, 16, 16, 32)))
	s.x[16] = uint64(uint32(s.x[16]) + uint32(s.x[9]))
	s.x[9] = 0x5ca3
	s.x[9] = uint64(uint32(insertBits(uint64(uint32(s.x[9])), 0x4a7e, 16, 16, 32)))
	s.x[17] = uint64(uint32(s.x[2]) ^ uint32(s.x[9]))
	s.x[9] = 0xa35d
	s.x[9] = uint64(uint32(insertBits(uint64(uint32(s.x[9])), 0xb581, 16, 16, 32)))
	s.x[17] = uint64(uint32(s.x[17]) + uint32(s.x[9]))
	s.x[16] = uint64(uint32(s.x[16]) ^ uint32(s.x[17]))
	s.x[9] = 0x7405
	s.x[9] = uint64(uint32(insertBits(uint64(uint32(s.x[9])), 0x1dab, 16, 16, 32)))
	s.x[17] = uint64(uint32(s.x[2]) ^ uint32(s.x[9]))
	s.x[9] = 0x8bfb
	s.x[9] = uint64(uint32(insertBits(uint64(uint32(s.x[9])), 0xe254, 16, 16, 32)))
	s.x[17] = uint64(uint32(s.x[17]) + uint32(s.x[9]))
	s.x[16] = uint64(uint32(s.x[16]) ^ uint32(s.x[17]))
	s.x[9] = 0xf2c2
	s.x[9] = uint64(uint32(insertBits(uint64(uint32(s.x[9])), 0x196a, 16, 16, 32)))
	s.x[16] = uint64(uint32(s.x[16]) + uint32(s.x[9]))
	s.x[8] = uint64(uint32(s.x[16]) + uint32(s.x[8]))
	s.x[9] = 0xbb96
	s.x[9] = uint64(uint32(insertBits(uint64(uint32(s.x[9])), 0x6a99, 16, 16, 32)))
	s.x[16] = uint64(uint32(s.x[8]) ^ uint32(s.x[9]))
	s.x[9] = 0xcc71
	s.x[9] = uint64(uint32(insertBits(uint64(uint32(s.x[9])), 0x28ce, 16, 16, 32)))
	s.x[17] = uint64(uint32(s.x[8]) ^ uint32(s.x[9]))
	s.x[9] = 0x64c7
	s.x[9] = uint64(uint32(insertBits(uint64(uint32(s.x[9])), 0xd41a, 16, 16, 32)))
	s.x[17] = uint64(uint32(s.x[17]) + uint32(s.x[9]))
	s.x[9] = 0xdf90
	s.x[9] = uint64(uint32(insertBits(uint64(uint32(s.x[9])), 0x3e83, 16, 16, 32)))
	s.x[0] = uint64(uint32(s.x[8]) ^ uint32(s.x[9]))
	s.x[9] = 0xff3f
	s.x[9] = uint64(uint32(insertBits(uint64(uint32(s.x[9])), 0x7fff, 16, 16, 32)))
	s.x[1] = uint64(uint32(s.x[8]) ^ uint32(s.x[9]))
	s.x[9] = 0x1322
	s.x[9] = uint64(uint32(insertBits(uint64(uint32(s.x[9])), 0x964d, 16, 16, 32)))
	s.x[16] = uint64(uint32(s.x[16]) + uint32(s.x[9]))
	s.x[8] = uint64(uint32(s.x[8]) ^ uint32(s.x[16]))
	s.x[9] = 0x7728
	s.x[9] = uint64(uint32(insertBits(uint64(uint32(s.x[9])), 0xc257, 16, 16, 32)))
	s.x[16] = uint64(uint32(s.x[0]) + uint32(s.x[9]))
	s.x[16] = uint64(uint32(s.x[17]) ^ uint32(s.x[16]))
	s.x[8] = uint64(uint32(s.x[8]) ^ uint32(s.x[16]))
	s.x[9] = 0x5789
	s.x[9] = uint64(uint32(insertBits(uint64(uint32(s.x[9])), 0x832b, 16, 16, 32)))
	s.x[16] = uint64(uint32(s.x[1]) + uint32(s.x[9]))
	s.x[9] = 0x9af5
	s.x[9] = uint64(uint32(insertBits(uint64(uint32(s.x[9])), 0x834, 16, 16, 32)))
	s.x[16] = uint64(uint32(s.x[16]) ^ uint32(s.x[9]))
	s.x[8] = uint64(uint32(s.x[8]) ^ uint32(s.x[16]))
	if err := s.storeValue(s.sp+0xf4, uint64(uint32(s.x[8])), 4); err != nil {
		return 0, fmt.Errorf("account state-setup store 0x1001a8670: %w", err)
	}
	value84, err := s.load(s.x[20]+shiftLeft(signExtend32(uint32(s.x[13])), 3, 64), 8)
	if err != nil {
		return 0, fmt.Errorf("account state-setup target 0x1001a8674: %w", err)
	}
	s.x[8] = value84
	value85, err := s.load(s.sp+0x14, 4)
	if err != nil {
		return 0, fmt.Errorf("account state-setup load 0x1001a8678: %w", err)
	}
	s.x[14] = uint64(uint32(value85))
	s.x[17] = 0x4f50
	return s.x[8], nil
}
