package recovery

import "fmt"

func runAccountFoldCircuit(s *state) (uint64, error) {
	const maximumRounds = 4096
	for round := 0; round < maximumRounds; round++ {
		next, err := runAccountFoldRoundCircuit(s)
		if err != nil || next != accountFoldRound {
			return next, err
		}
	}
	return 0, fmt.Errorf("account fold exceeded %d rounds", maximumRounds)
}

func runAccountFoldRoundCircuit(s *state) (uint64, error) {
	s.traceBlock(accountFoldRound)
	s.x[2], s.x[28], s.x[19] = s.x[28], s.x[15], s.x[11]
	s.x[11] = uint64(uint32(rotateRight(s.x[11], 27, 32)))
	s.x[13] = uint64(uint32(s.x[0]) ^ 0x9db2)
	s.x[15] = uint64((uint32(s.x[11]) ^ 0x8d9daa1c) + 0x2fbd1333)
	s.x[17] = uint64((uint32(s.x[11]) ^ 0x91ad4258) + 0x338dfb77)
	s.x[15] = uint64(uint32(s.x[15]) ^ uint32(s.x[17]))
	s.x[17] = uint64(uint32(s.x[13]) + 0x7fc8bda6)
	s.x[11] = uint64(uint32(s.x[11]) ^ uint32(s.x[17]))
	s.x[11] = uint64(uint32(s.x[11]) + 0xdde9e6cc)
	s.x[11] = uint64(uint32(s.x[15]) ^ uint32(s.x[11]))
	s.x[15] = uint64(uint32(s.x[28]) ^ uint32(s.x[1]))
	s.x[15] = uint64(uint32(s.x[15]) ^ uint32(s.x[2]))
	s.x[17] = uint64((uint32(s.x[15]) ^ 0x6305f35f) + 0x7e68333a)
	s.x[5] = uint64((uint32(s.x[15]) ^ 0xdae379fb) + 0xc78eb99e)
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[5]))
	s.x[15] = uint64((uint32(s.x[15]) ^ 0xb2f94719) + 0xaf948780)
	s.x[15] = uint64(uint32(s.x[17]) ^ uint32(s.x[15]))
	s.x[17] = uint64((uint32(s.x[16]) ^ 0xfba42f82) + 0x045bd07e)
	s.x[5] = uint64((uint32(s.x[16]) ^ 0x93e64b78) + 0x6c19b488)
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[5]))
	s.x[16] = uint64((uint32(s.x[16]) ^ 0x635da947) + 0x9ca256b9)
	address := s.x[14] + (s.x[7] << 2)
	value, err := s.load(s.x[8]+address, 4)
	if err != nil {
		return 0, fmt.Errorf("account fold input: %w", err)
	}
	s.x[5] = uint64(uint32(value))
	s.x[16] = uint64(uint32(s.x[17]) ^ uint32(s.x[16]))
	s.x[17] = uint64((uint32(s.x[5]) ^ 0x1cbd9e5e) + 0xe34261a2)
	s.x[6] = uint64((uint32(s.x[5]) ^ 0xeda321bb) + 0x125cde45)
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[6]))
	s.x[5] = uint64((uint32(s.x[5]) ^ 0x3f6bc420) + 0xc0943be0)
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[5]))
	s.x[11] = uint64(uint32(s.x[16]) + uint32(s.x[11]))
	s.x[11] = uint64(uint32(s.x[11]) + uint32(s.x[15]))
	s.x[11] = uint64(uint32(s.x[11]) + 0x41293bcd)
	s.x[11] = uint64(uint32(s.x[11]) + uint32(s.x[17]))
	s.x[15] = uint64(uint32(s.x[11]) ^ 0xb43f2b96)
	s.x[16] = uint64((uint32(s.x[11]) ^ 0x03261adf) + 0x89a5a9c2)
	s.x[17] = uint64(uint32(s.x[11]) ^ 0x3dba8053)
	s.x[15] = uint64(uint32(s.x[15]) + 0x3ebc9889)
	s.x[15] = uint64(uint32(s.x[11]) ^ uint32(s.x[15]))
	s.x[11] = uint64(uint32(s.x[11]) ^ 0xffdffdfb)
	s.x[17] = uint64(uint32(s.x[17]) + 0xb739334e)
	s.x[16] = uint64(uint32(s.x[16]) ^ uint32(s.x[17]))
	s.x[15] = uint64(uint32(s.x[15]) ^ uint32(s.x[16]))
	s.x[11] = uint64(uint32(s.x[11]) + 0x755c4ee6)
	s.x[11] = uint64(uint32(s.x[11]) ^ 0x7e63815c)
	s.x[11] = uint64(uint32(s.x[15]) ^ uint32(s.x[11]))
	s.x[15] = uint64(uint32(s.x[30]) ^ uint32(rotateRight(s.x[1], 2, 32)))
	s.x[7]++
	s.x[9] = accountFoldLimit
	s.compare64(s.x[7], s.x[9])
	gate := uint32(boolValue(s.condition("eq")))
	s.x[9] = 0x83c3
	s.x[16] = uint64(gate * uint32(s.x[9]))
	s.x[16] = uint64(uint32(s.x[16]) ^ uint32(s.x[0]))
	next, err := loadBootstrapPointer(s, s.x[20], uint32(s.x[16]))
	if err != nil {
		return 0, fmt.Errorf("account fold dispatch: %w", err)
	}
	s.traceBlock(accountFoldTail)
	s.x[5] = next
	s.x[1] = s.x[19]
	s.x[16] = s.x[2]
	s.x[17] = 0x2583
	return next, nil
}
