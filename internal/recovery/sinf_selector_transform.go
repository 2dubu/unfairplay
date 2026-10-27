package recovery

import "fmt"

const (
	sinfModuloDispatch = uint64(0x10027b13c)
	sinfByteWriter     = uint64(0x1001c9a34)
	sinfReverseWriter  = uint64(0x1004d1cd0)
)

func runSINFModuloDispatch(s *state) (uint64, error) {
	s.traceBlock(sinfModuloDispatch)
	s.x[10] = 0xffffb589
	s.x[11] = uint64(uint32(s.x[14]) + uint32(s.x[10]))
	s.x[11] = uint64(uint32(s.x[11]) | 0x4410)
	s.x[10] = uint64(uint32(s.x[14]) ^ 0x45ed)
	s.x[10] = uint64(uint32(s.x[10]) + uint32(s.x[10]<<2))
	s.x[12] = uint64(uint32(s.x[12]) * uint32(s.x[12]))
	s.x[15] = 0x00020005
	s.x[15] = uint64(uint32(s.x[12])) * uint64(uint32(s.x[15]))
	s.x[15] >>= 32
	s.x[16] = uint64(uint32(s.x[12]) - uint32(s.x[15]))
	s.x[15] = uint64(uint32(s.x[15]) + (uint32(s.x[16]) >> 1))
	s.x[15] = uint64(uint32(s.x[15]) >> 14)
	s.x[15] = uint64(uint32(s.x[15]) - uint32(s.x[15]<<15))
	s.x[12] = uint64(uint32(s.x[12]) + uint32(s.x[15]))
	s.x[15] = uint64(uint32(s.x[11]) ^ 0x541b)
	s.x[15] = uint64(uint32(s.x[9]) >> (uint32(s.x[15]) & 31))
	s.x[9] = uint64(uint32(s.x[10]) ^ 0x458abb43)
	s.x[9] = uint64(uint32(s.x[9]) ^ uint32(s.x[15]))
	s.x[16] = 0x3c861dcb
	s.compare32(uint32(s.x[15]), uint32(s.x[16]))
	gate := uint32(boolValue(s.condition("eq")))
	s.x[15] = uint64(gate)
	s.x[16] = 0x5a
	s.x[15] = uint64(gate * uint32(s.x[16]))
	s.x[14] = uint64(uint32(s.x[15]) ^ uint32(s.x[14]))
	next, err := loadBootstrapPointer(s, s.x[8], uint32(s.x[14]))
	if err != nil {
		return 0, fmt.Errorf("SINF modulo dispatch: %w", err)
	}
	s.x[14] = next
	s.x[17] = 0x36d7
	return next, nil
}

func runSINFByteWriter(s *state) (uint64, error) {
	s.traceBlock(sinfByteWriter)
	s.x[17] = uint64(uint32(s.x[9]) ^ uint32(s.x[11]))
	s.x[19] = uint64(uint32(s.x[11]) ^ uint32(s.x[1]))
	s.x[17] = uint64(uint32(s.x[17]) ^ 0x83918a14)
	s.x[26] = uint64(uint32(s.x[11]) ^ uint32(s.x[3]))
	s.x[26] = uint64(uint32(s.x[26]) + uint32(s.x[4]))
	s.x[16] = uint64(uint32(s.x[11]) ^ uint32(s.x[5]))
	s.x[19] = uint64(uint32(s.x[19]) + uint32(s.x[2]))
	s.x[19] = uint64(uint32(s.x[19]) ^ uint32(s.x[11]))
	s.x[16] = uint64(uint32(s.x[16]) + uint32(s.x[6]))
	s.x[16] = uint64(uint32(s.x[26]) ^ uint32(s.x[16]))
	s.x[16] = uint64(uint32(s.x[19]) ^ uint32(s.x[16]))
	s.x[19] = 0xa645f3d7
	s.x[17] = uint64(uint32(s.x[17]) + uint32(s.x[19]))
	s.x[16] = uint64(uint32(s.x[16]) ^ uint32(s.x[17]))
	s.x[17] = uint64(uint32(s.x[16]) ^ uint32(s.x[7]))
	s.x[10]--
	s.x[17] = uint64(uint32(s.x[17]) + uint32(s.x[28]))
	s.x[19] = uint64(uint32(s.x[16]) ^ uint32(s.x[30]))
	s.x[19] = uint64(uint32(s.x[19]) + uint32(s.x[21]))
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[19]))
	s.x[16] = uint64(uint32(s.x[16]) ^ uint32(s.x[22]))
	s.x[16] = uint64(uint32(s.x[16]) + uint32(s.x[23]))
	s.x[16] = uint64(uint32(s.x[17]) ^ uint32(s.x[16]))
	s.x[17] = uint64(uint32(s.x[16]) + uint32(s.x[20]))
	s.x[16] = uint64(uint32(s.x[16]) * uint32(s.x[12]))
	s.x[16] = uint64(uint32(s.x[14])*uint32(s.x[17]) + uint32(s.x[16]))
	s.x[16] = uint64(uint32(s.x[13]) + uint32(s.x[16]))
	if err := s.storeValue(s.sp+0x58, s.x[10], 8); err != nil {
		return 0, fmt.Errorf("SINF byte-writer counter: %w", err)
	}
	if err := s.storeValue(s.x[8]+s.x[10], uint64(byte(s.x[16])), 1); err != nil {
		return 0, fmt.Errorf("SINF byte-writer output: %w", err)
	}
	counter, err := s.load(s.sp+0x58, 8)
	if err != nil {
		return 0, fmt.Errorf("SINF byte-writer saved counter: %w", err)
	}
	s.x[10] = counter
	s.compare64(counter, 0)
	gate := uint32(boolValue(s.condition("eq")))
	s.x[16] = uint64(gate * uint32(s.x[15]))
	s.x[11]++
	s.x[16] = uint64(uint32(s.x[16]) ^ uint32(s.x[9]))
	next, err := s.load(s.x[24]+uint64(uint32(s.x[16]))*8, 8)
	if err != nil {
		return 0, fmt.Errorf("SINF byte-writer dispatch: %w", err)
	}
	s.x[19] = next
	s.x[14]++
	s.x[16] = 0x3a24
	s.x[13] += s.x[16]
	s.x[17] = 0x26a7
	return next, nil
}

func runSINFReverseWriter(s *state) (uint64, error) {
	s.traceBlock(sinfReverseWriter)
	s.x[13] = s.x[16] - 1
	s.x[15] = uint64(uint32(s.x[12]) ^ 0x2b899e06)
	s.x[16] = uint64(uint32(s.x[9]) ^ 0x6156c0e7)
	s.x[17] = uint64(uint32(s.x[12]) ^ 0xcabc01c8)
	s.x[15] = uint64(uint32(s.x[15]) + uint32(s.x[16]))
	s.x[16] = uint64(uint32(s.x[12]) ^ uint32(s.x[23]))
	s.x[0] = uint64(uint32(s.x[12]) ^ uint32(s.x[27]))
	s.x[17] = uint64(uint32(s.x[17]) + 0x80633b1b)
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[12]))
	s.x[16] = uint64(uint32(s.x[16]) + uint32(s.x[24]))
	s.x[15] = uint64(uint32(s.x[15]) ^ uint32(s.x[16]))
	s.x[15] = uint64(uint32(s.x[17]) ^ uint32(s.x[15]))
	s.x[16] = uint64(uint32(s.x[0]) + uint32(s.x[28]))
	s.x[15] = uint64(uint32(s.x[15]) ^ uint32(s.x[16]))
	s.x[16] = uint64(uint32(s.x[15]) ^ uint32(s.x[5]))
	s.x[17] = uint64(uint32(s.x[12]) + uint32(s.x[7]))
	s.x[16] = uint64(uint32(s.x[16]) + uint32(s.x[6]))
	s.x[0] = uint64(uint32(s.x[15]) ^ uint32(s.x[11]))
	s.x[0] = uint64(uint32(s.x[0]) + uint32(s.x[19]))
	s.x[16] = uint64(uint32(s.x[16]) ^ uint32(s.x[0]))
	s.x[15] = uint64(uint32(s.x[15]) ^ uint32(s.x[20]))
	s.x[15] = uint64(uint32(s.x[15]) + uint32(s.x[3]))
	s.x[15] = uint64(uint32(s.x[16]) ^ uint32(s.x[15]))
	s.x[15] = uint64(uint32(s.x[15]) + uint32(s.x[4]))
	s.x[16] = uint64(uint32(s.x[15]) * uint32(s.x[17]))
	s.x[15] = uint64(uint32(s.x[15]) * uint32(s.x[30]))
	s.x[0] = 0x50
	s.x[15] = uint64(uint32(s.x[17])*uint32(s.x[0]) + uint32(s.x[15]))
	if err := s.storeValue(s.sp+0x120, s.x[13], 8); err != nil {
		return 0, fmt.Errorf("SINF reverse-writer counter: %w", err)
	}
	s.x[15] = uint64(uint32(s.x[15])*uint32(s.x[1]) + uint32(s.x[16]))
	s.x[15] = uint64(uint32(s.x[15]) + uint32(s.x[8]))
	if err := s.storeValue(s.x[14]+s.x[13], uint64(byte(s.x[15])), 1); err != nil {
		return 0, fmt.Errorf("SINF reverse-writer output: %w", err)
	}
	counter, err := s.load(s.sp+0x120, 8)
	if err != nil {
		return 0, fmt.Errorf("SINF reverse-writer saved counter: %w", err)
	}
	s.x[16] = counter
	s.compare64(counter, 0)
	gate := uint32(boolValue(s.condition("eq")))
	s.x[13] = uint64(gate * uint32(s.x[2]))
	s.x[13] = uint64(uint32(s.x[13]) ^ uint32(s.x[9]))
	next, err := loadBootstrapPointer(s, s.x[21], uint32(s.x[13]))
	if err != nil {
		return 0, fmt.Errorf("SINF reverse-writer dispatch: %w", err)
	}
	s.x[13] = next
	s.x[12]++
	s.x[17] = 0x3a56
	return next, nil
}
