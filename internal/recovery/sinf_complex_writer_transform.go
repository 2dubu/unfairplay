package recovery

import "fmt"

const (
	sinfComplexWriterStart = uint64(0x100011eb4)
	sinfComplexWriterTail  = uint64(0x100012000)
)

func runSINFComplexWriter(s *state) (uint64, error) {
	if _, err := runSINFComplexWriterFront(s); err != nil {
		return 0, err
	}
	return runSINFComplexWriterFinish(s)
}

func runSINFComplexWriterFront(s *state) (uint64, error) {
	s.traceBlock(sinfComplexWriterStart)
	value, err := s.load(s.x[28]+0x210, 4)
	if err != nil {
		return 0, fmt.Errorf("SINF complex-writer first word: %w", err)
	}
	s.x[17] = uint64(uint32(value))
	const stateBase = uint64(0x101075000)
	s.x[8] = stateBase
	value, err = s.load(stateBase+0xb40, 4)
	if err != nil {
		return 0, fmt.Errorf("SINF complex-writer second word: %w", err)
	}
	s.x[2] = uint64(uint32(value))
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[2]))
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[5]))
	s.x[17] = uint64(uint32(s.x[17]) * uint32(s.x[6]))
	s.x[2] = s.x[19] + uint64(byte(s.x[17]))
	value, err = s.load(s.x[2], 1)
	if err != nil {
		return 0, fmt.Errorf("SINF complex-writer first byte: %w", err)
	}
	s.x[2] = uint64(uint32(byte(value)) ^ 0xcd)
	value, err = s.load(s.x[22]+uint64(uint32(s.x[2])), 1)
	if err != nil {
		return 0, fmt.Errorf("SINF complex-writer second byte: %w", err)
	}
	s.x[2] = uint64(byte(value))
	s.x[17] = uint64(uint32(s.x[2]) ^ uint32(s.x[17]))
	s.x[17] &= 0xff
	s.x[17] = uint64(uint32(s.x[17]) + 0x1a8)
	value, err = s.load(s.x[10]+uint64(uint32(s.x[17]))*8, 8)
	if err != nil {
		return 0, fmt.Errorf("SINF complex-writer first state: %w", err)
	}
	s.x[17] = value
	s.x[8] = 0xffffffffc2609e4e
	s.x[2] = s.x[17] + s.x[8]
	value, err = s.load(s.x[2], 8)
	if err != nil {
		return 0, fmt.Errorf("SINF complex-writer first output base: %w", err)
	}
	s.x[17] = value + s.x[4]
	if err := s.storeValue(s.x[17], uint64(byte(s.x[7])), 1); err != nil {
		return 0, fmt.Errorf("SINF complex-writer first output: %w", err)
	}

	value, err = s.load(s.x[28]+0x210, 4)
	if err != nil {
		return 0, fmt.Errorf("SINF complex-writer third word: %w", err)
	}
	s.x[17] = uint64(uint32(value))
	value, err = s.load(s.x[30], 4)
	if err != nil {
		return 0, fmt.Errorf("SINF complex-writer fourth word: %w", err)
	}
	s.x[3] = uint64(uint32(value))
	s.x[3] = uint64(uint32(s.x[3]) ^ uint32(s.x[5]))
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[3]))
	s.x[17] = uint64(uint32(s.x[17]) * uint32(s.x[6]))
	s.x[3] = s.x[15] + uint64(byte(s.x[17]))
	value, err = s.load(s.x[3], 1)
	if err != nil {
		return 0, fmt.Errorf("SINF complex-writer third byte: %w", err)
	}
	s.x[3] = uint64(uint32(byte(value)) ^ 0xfa)
	value, err = s.load(s.x[1]+uint64(uint32(s.x[3])), 1)
	if err != nil {
		return 0, fmt.Errorf("SINF complex-writer fourth byte: %w", err)
	}
	s.x[3] = uint64(byte(value))
	s.x[17] = uint64(uint32(s.x[3]) ^ uint32(s.x[17]))
	s.x[17] &= 0xff
	s.x[17] = uint64(uint32(s.x[17]) + 0xfab)
	value, err = s.load(s.x[10]+uint64(uint32(s.x[17]))*8, 8)
	if err != nil {
		return 0, fmt.Errorf("SINF complex-writer second state: %w", err)
	}
	s.x[17] = value
	s.x[8] = 0xffffffffddc4a14d
	s.x[3] = s.x[17] + s.x[8]
	value, err = s.load(s.x[3], 8)
	if err != nil {
		return 0, fmt.Errorf("SINF complex-writer second output base: %w", err)
	}
	s.x[17] = value
	s.x[13] = uint64(uint32(s.x[7]) - 0x6c)
	s.x[13] = s.x[17] + uint64(byte(s.x[13]))
	if err := s.storeValue(s.x[13], uint64(byte(s.x[4])), 1); err != nil {
		return 0, fmt.Errorf("SINF complex-writer second output: %w", err)
	}

	runSINFComplexWriterNetwork(s)
	return 0, nil
}

func runSINFComplexWriterNetwork(s *state) {
	runSINFComplexWriterNetworkFirst(s)
	runSINFComplexWriterNetworkSecond(s)
}

func runSINFComplexWriterNetworkFirst(s *state) {
	s.x[8] = s.x[1]
	s.x[13] = uint64((uint32(s.x[7]) & 0xa) ^ 0xffffffdb)
	s.x[17] = uint64(uint32(s.x[7]) << (uint32(s.x[24]) & 31))
	s.x[17] = uint64(uint32(s.x[17]) & uint32(s.x[25]))
	s.x[1] = uint64(uint32(s.x[17]) ^ 0x1e)
	s.x[16] = s.x[15]
	s.x[15] = uint64(uint32(s.x[7]) ^ uint32(s.x[26]))
	s.x[1] = uint64(uint32(s.x[1]) & uint32(s.x[15]))
	s.x[1] = uint64(uint32(s.x[1]) ^ uint32(s.x[13]))
	for range 2 {
		s.x[1] = uint64(uint32(s.x[15]) & uint32(s.x[1]<<1))
		s.x[1] = uint64(uint32(s.x[1]) ^ uint32(s.x[13]))
	}
	s.x[1] = uint64(uint32(s.x[15]) & uint32(s.x[1]<<1))
	s.x[13] = uint64(uint32(s.x[1]) ^ uint32(s.x[13]))
}

func runSINFComplexWriterNetworkSecond(s *state) {
	runSINFComplexWriterNetworkSecondFront(s)
	runSINFComplexWriterNetworkSecondFinish(s)
}

func runSINFComplexWriterNetworkSecondFront(s *state) {
	s.x[1] = uint64(uint32(s.x[7]) + 0x14)
	s.x[1] = uint64(uint32(int32(int8(s.x[1]))))
	_ = s.addWithCarry(uint64(uint32(s.x[1])), 1, false, 32)
	s.x[1] = uint64((uint32(s.x[7]) & 0x3e) ^ uint32(s.x[21]))
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[23]))
	s.x[17] = uint64(uint32(s.x[17]) & uint32(s.x[7]))
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[1]))
	for range 2 {
		s.x[17] = uint64(uint32(s.x[15]) & uint32(s.x[17]<<1))
		s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[1]))
	}
	s.x[15] = uint64(uint32(s.x[15]) & uint32(s.x[17]<<1))
	s.x[17] = uint64(uint32(s.x[1]) ^ uint32(s.x[7]))
	s.x[1] = s.x[8]
	s.x[15] = uint64(uint32(s.x[17]) ^ uint32(s.x[15]))
}

func runSINFComplexWriterNetworkSecondFinish(s *state) {
	if s.condition("gt") {
		s.x[17] = uint64(uint32(s.x[20]))
	} else {
		s.x[17] = uint64(uint32(s.x[27]))
	}
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[7]))
	s.x[13] = uint64(uint32(s.x[17]) ^ uint32(s.x[13]<<2))
	s.x[13] = uint64(uint32(s.x[13]) ^ uint32(s.x[15]<<1))
	s.x[15] = uint64(uint32(s.x[13]) ^ uint32(s.x[12]))
	s.x[15] = uint64(uint32(s.x[15]) - 0x40)
	s.x[17] = uint64(uint32(s.x[13]) ^ uint32(s.x[11]))
}

func runSINFComplexWriterFinish(s *state) (uint64, error) {
	s.traceBlock(sinfComplexWriterTail)
	s.x[17] = uint64(uint32(s.x[17]) - 0x6d)
	s.x[15] = uint64(uint32(s.x[15]) ^ uint32(s.x[17]))
	s.x[13] = uint64(uint32(s.x[13]) ^ 3)
	s.x[13] = uint64(uint32(s.x[13]) - 0x31)
	s.x[13] = uint64(uint32(s.x[15]) ^ uint32(s.x[13]))
	s.x[15] = s.x[16]
	s.x[7] = uint64(uint32(s.x[13]) - 0x78)
	s.compare64(s.x[4], 0xff)
	s.x[4]++
	gate := uint32(boolValue(s.condition("ne")))
	s.x[13] = uint64(gate * uint32(s.x[9]))
	s.x[13] = uint64(uint32(s.x[13]) ^ uint32(s.x[0]))
	next, err := loadBootstrapPointer(s, s.x[14], uint32(s.x[13]))
	if err != nil {
		return 0, fmt.Errorf("SINF complex-writer dispatch: %w", err)
	}
	s.x[13] = next
	s.x[17] = 0xf807
	return next, nil
}
