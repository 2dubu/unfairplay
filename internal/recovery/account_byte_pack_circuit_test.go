package recovery

import "fmt"

func runAccountBytePackCircuit(s *state) (uint64, error) {
	for index := 0; index < 4096; index++ {
		next, err := runAccountBytePackStepCircuit(s)
		if err != nil || next != accountBytePackPass {
			return next, err
		}
	}
	return 0, fmt.Errorf("account byte-pack exceeded 4096 words")
}

func runAccountBytePackStepCircuit(s *state) (uint64, error) {
	s.traceBlock(accountBytePackPass)
	value, err := s.load(s.x[1]+s.x[4]+s.x[14], 1)
	if err != nil {
		return 0, fmt.Errorf("account byte-pack first byte: %w", err)
	}
	s.x[11] = uint64(uint32(byte(value)) ^ uint32(s.x[8]))
	s.x[11] = uint64(uint32(s.x[11] << 24))
	s.x[15] = uint64(uint32(s.x[9]) + 0x2ba76699)
	s.x[12] = s.x[15] ^ 0xffffffffffffffad
	s.x[12] += s.x[4]
	value, err = s.load(s.x[1]+s.x[12], 1)
	if err != nil {
		return 0, fmt.Errorf("account byte-pack second byte: %w", err)
	}
	s.x[12] = uint64(uint32(byte(value)) ^ uint32(s.x[8]))
	value, err = s.load(s.x[1]+s.x[4]+s.x[28], 1)
	if err != nil {
		return 0, fmt.Errorf("account byte-pack third byte: %w", err)
	}
	s.x[13] = uint64(byte(value))
	s.x[11] = uint64(uint32(s.x[11]) | uint32(s.x[12]<<16))
	s.x[12] = uint64(uint32(s.x[13]) ^ uint32(s.x[8]))
	s.x[11] = uint64(uint32(s.x[11]) | uint32(s.x[12]<<8))
	s.x[12] = 0xffffffffd4583c78 + s.x[4]
	value, err = s.load(s.x[1]+s.x[12], 1)
	if err != nil {
		return 0, fmt.Errorf("account byte-pack fourth byte: %w", err)
	}
	s.x[12] = uint64(uint32(byte(value)) ^ uint32(s.x[8]))
	s.x[11] = uint64(uint32(s.x[11]) | uint32(s.x[12]))
	s.x[13] = uint64(uint32(s.x[11]) & 0x4e757bd5)
	s.x[12] = uint64(uint32(s.x[12]) & 0x10)
	s.x[12] = uint64(uint32(s.x[13]) ^ uint32(s.x[12]))
	s.x[11] = uint64(uint32(s.x[11]) - uint32(s.x[12]<<1))
	s.x[12] = uint64(uint32(s.x[4]) ^ 0xc1c7a744)
	s.x[13] = uint64(uint32(s.x[4]) & 0x75)
	s.x[16] = uint64(uint32(s.x[13] << 1))
	s.x[2] = 0x140e1300
	s.x[17] = uint64(uint32(s.x[16]) | uint32(s.x[2]))
	s.x[17] = uint64(uint32(s.x[17]) & uint32(s.x[12]))
	s.x[16] = uint64(uint32(s.x[16]) & 0x20)
	s.x[13] = uint64(uint32(s.x[16]) ^ uint32(s.x[13]))
	s.x[13] = uint64(uint32(s.x[13]) ^ uint32(s.x[17]))
	s.x[16] = uint64(uint32(s.x[12] << 1))
	s.x[17] = uint64(uint32(s.x[16]) ^ 0x3f312562)
	s.x[12] = uint64(uint32(s.x[17]) & uint32(s.x[12]))
	s.x[16] = uint64(uint32(s.x[16]) & 0x00808830)
	s.x[16] = uint64(uint32(s.x[16]) ^ 0x00808811)
	s.x[12] = uint64(uint32(s.x[16]) ^ uint32(s.x[12]))
	s.x[16] = uint64(uint32(s.x[13] << 2))
	s.x[17] = uint64(uint32(s.x[16]) ^ uint32(s.x[2]))
	s.x[17] = uint64(uint32(s.x[17]) & uint32(s.x[12]))
	s.x[16] = uint64(uint32(s.x[16]) & 0x10180b30)
	s.x[13] = uint64(uint32(s.x[16]) ^ uint32(s.x[13]))
	s.x[13] = uint64(uint32(s.x[13]) ^ uint32(s.x[17]))
	s.x[16] = uint64(uint32(s.x[12] << 2))
	s.x[17] = uint64(uint32(s.x[16]) ^ 0x42707fc4)
	s.x[12] = uint64(uint32(s.x[17]) & uint32(s.x[12]))
	s.x[16] = uint64(uint32(s.x[16]) & 0x159f9b30)
	s.x[12] = uint64(uint32(s.x[12]) ^ uint32(s.x[16]))
	s.x[16] = uint64(uint32(s.x[12]) ^ 0x01819331)
	s.x[12] = uint64(uint32(s.x[12]) ^ 0x141e0801)
	s.x[12] = uint64(uint32(s.x[12]) & uint32(s.x[13]<<4))
	s.x[12] = uint64(uint32(s.x[12]) ^ uint32(s.x[13]))
	s.x[13] = uint64(uint32(s.x[16] << 4))
	s.x[17] = uint64(uint32(s.x[13]) ^ 0x59f9b310)
	s.x[16] = uint64(uint32(s.x[17]) & uint32(s.x[16]))
	s.x[13] = uint64(uint32(s.x[13]) & 0x159f9b10)
	s.x[13] = uint64(uint32(s.x[16]) ^ uint32(s.x[13]))
	s.x[16] = uint64(uint32(s.x[13]) ^ 0x04060821)
	s.x[13] = uint64(uint32(s.x[13]) ^ 0x11999321)
	s.x[13] = uint64(uint32(s.x[13]) & uint32(s.x[12]<<8))
	s.x[12] = uint64(uint32(s.x[12]) ^ 0x159f9b31)
	s.x[12] = uint64(uint32(s.x[12]) ^ uint32(s.x[13]))
	s.x[13] = uint64(uint32(s.x[16] << 8))
	s.x[17] = uint64(uint32(s.x[13]) ^ 0x1f9b0000)
	s.x[16] = uint64(uint32(s.x[17]) & uint32(s.x[16]))
	s.x[2] = 0x159f0000
	s.x[13] = uint64(uint32(s.x[13]) & uint32(s.x[2]))
	s.x[13] = uint64(uint32(s.x[13]) ^ uint32(s.x[16]))
	s.x[13] = uint64(uint32(s.x[13]) ^ 0x00040000)
	s.x[16] = uint64(uint32(s.x[12] << 16))
	s.x[17] = uint64(uint32(s.x[16]) ^ 0x1b310000)
	s.x[13] = uint64(uint32(s.x[17]) & uint32(s.x[13]))
	s.x[16] = uint64(uint32(s.x[16]) & uint32(s.x[2]))
	s.x[12] = uint64(uint32(s.x[16]) ^ uint32(s.x[12]))
	s.x[12] = uint64(uint32(s.x[12]) ^ uint32(s.x[13]))
	s.x[13] = 0xce757bc5
	s.x[11] = uint64(uint32(s.x[11]) + uint32(s.x[13]))
	s.x[12] = uint64(uint32(s.x[4]) ^ uint32(s.x[12]<<1))
	s.x[13] = 0x37514285
	s.x[12] = uint64(uint32(s.x[13]) ^ (uint32(s.x[12]) >> 2))
	if err := s.storeValue(s.x[10]+uint64(uint32(s.x[12]))*4, uint64(uint32(s.x[11])), 4); err != nil {
		return 0, fmt.Errorf("account byte-pack output: %w", err)
	}
	s.x[11] = 0xffffffffd4583c79 + s.x[4]
	s.x[4] += 4
	s.compare64(s.x[11], 0x40)
	gate := uint32(boolValue(s.condition("lo")))
	s.x[12] = 0xc285
	s.x[11] = uint64(gate * uint32(s.x[12]))
	s.x[11] = uint64(uint32(s.x[11]) ^ uint32(s.x[9]))
	next, err := loadBootstrapPointer(s, s.x[21], uint32(s.x[11]))
	if err != nil {
		return 0, fmt.Errorf("account byte-pack dispatch: %w", err)
	}
	s.x[11] = next
	s.x[20] = 0x2ba7c39b
	s.x[17] = 0x3601
	return next, nil
}
