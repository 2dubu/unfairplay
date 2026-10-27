package recovery

import "fmt"

const (
	accountCompoundRouteA = uint64(0x100052750)
	accountCompoundRouteB = uint64(0x100059890)
	accountCompoundRouteC = uint64(0x100196e70)
)

func runAccountCompoundRouteA(s *state) (uint64, error) {
	s.traceBlock(accountCompoundRouteA)
	s.x[25] = 0x80000004
	s.x[21] = 0x7c79d912
	s.x[22] = 0xfc79d912
	s.x[10] = 0x1448107c
	s.x[13] = uint64(uint32(s.x[13]) ^ uint32(s.x[10]))
	s.x[10] = 0x663d03d0
	s.compare64(s.x[12], s.x[10])
	s.x[14] = uint64(uint32(boolValue(s.condition("lo"))))
	s.x[10] = 0xffffc78b
	s.x[10] = uint64(uint32(s.x[13]) + uint32(s.x[10]))
	s.x[15] = 0x2cce
	s.x[10] = uint64(uint32(s.x[10]) | uint32(s.x[15]))
	s.x[26] = uint64(uint32(s.x[13]) - 0xf9b)
	s.x[15] = s.x[11] + 1
	s.x[16] = 0x663cd9f7
	s.x[11] += s.x[16]
	s.x[11] += s.x[26]
	s.x[16] = 0xffffffff99c2fc2f
	s.compare64(s.x[15], s.x[16])
	s.x[16] = uint64(uint32(boolValue(s.condition("hi"))))
	s.x[16] = uint64(uint32(s.x[14]) ^ uint32(s.x[16]))
	s.compare64(s.x[11], s.x[12])
	s.x[11] = uint64(uint32(boolValue(s.condition("lo"))))
	s.compare32(uint32(s.x[16]), 0)
	if s.condition("ne") {
		s.x[11] = uint64(uint32(s.x[14]))
	}
	s.x[14] = uint64(uint32(s.x[10]) ^ 0x2d14)
	s.x[11] = uint64(uint32(s.x[11]) * uint32(s.x[14]))
	s.x[11] = uint64(uint32(s.x[11]) ^ uint32(s.x[13]))
	next, err := loadBootstrapPointer(s, s.x[20], uint32(s.x[11]))
	if err != nil {
		return 0, fmt.Errorf("account compound-route A target: %w", err)
	}
	s.x[13] = next
	s.x[11] = s.x[15]
	s.x[17] = 0xc702
	return next, nil
}

func runAccountCompoundRouteB(s *state) (uint64, error) {
	s.traceBlock(accountCompoundRouteB)
	value, err := s.load(s.x[8]+s.x[11], 1)
	if err != nil {
		return 0, fmt.Errorf("account compound-route B first input: %w", err)
	}
	s.x[13] = uint64(byte(value))
	value, err = s.load(s.x[9]+s.x[11], 1)
	if err != nil {
		return 0, fmt.Errorf("account compound-route B second input: %w", err)
	}
	s.x[14] = uint64(byte(value))
	s.x[15] = 0x1b
	s.x[14] = uint64(uint32(s.x[14]) ^ uint32(s.x[15]))
	s.x[13] = uint64(uint32(s.x[13]) - uint32(s.x[14]))
	s.x[15] = 0xffdffaef
	s.x[14] = uint64(uint32(s.x[13]) ^ uint32(s.x[15]))
	s.x[16] = 0xffbff5de
	s.x[13] = uint64(uint32(s.x[16]) & uint32(s.x[13]<<1))
	s.x[13] = uint64(uint32(s.x[13]) + uint32(s.x[14]))
	s.compare32(uint32(s.x[13]), uint32(s.x[15]))
	s.x[14] = uint64(uint32(boolValue(s.condition("eq"))))
	s.x[13] = uint64(uint32(s.x[10]) + 0x1447fb3b)
	s.x[15] = uint64(uint32(s.x[13]) ^ 0x14485a35)
	s.x[14] = uint64(uint32(s.x[14]) * uint32(s.x[15]))
	s.x[10] = uint64(uint32(s.x[14]) ^ uint32(s.x[10]))
	next, err := loadBootstrapPointer(s, s.x[20], uint32(s.x[10]))
	if err != nil {
		return 0, fmt.Errorf("account compound-route B target: %w", err)
	}
	s.x[10] = next
	s.x[0] = 0xffff5bde
	s.x[17] = 0x5613
	return next, nil
}

func runAccountCompoundRouteC(s *state) (uint64, error) {
	s.traceBlock(accountCompoundRouteC)
	s.x[16] = uint64(uint32(s.x[11]) - 0x3f6)
	s.compare64(s.x[15], s.x[14])
	s.x[17] = uint64(uint32(boolValue(s.condition("lo"))))
	s.x[11] = 0xffffe766
	s.x[13] = 0x1c90
	s.x[11] = uint64(uint32(s.x[16]) + uint32(s.x[11]))
	s.x[11] = uint64(uint32(s.x[11]) | uint32(s.x[13]))
	s.x[13] = uint64(uint32(s.x[11]) - 9)
	s.x[0] = 0xffffffffffffe376
	s.x[0] += s.x[12]
	s.x[0] += s.x[13]
	s.x[2] = 0x0c756d0b
	s.x[3] = 0xfffffffff38a92f5
	s.compare64(s.x[0], s.x[3])
	s.x[3] = uint64(uint32(boolValue(s.condition("hi"))))
	s.x[3] = uint64(uint32(s.x[17]) ^ uint32(s.x[3]))
	s.x[12] += s.x[2]
	s.compare64(s.x[12], s.x[15])
	s.x[12] = uint64(uint32(boolValue(s.condition("lo"))))
	s.compare32(uint32(s.x[3]), 0)
	if s.condition("ne") {
		s.x[12] = uint64(uint32(s.x[17]))
	}
	s.x[12] = uint64(uint32(s.x[12]) ^ uint32(s.x[16]))
	s.x[12] = uint64(uint32(s.x[12]) ^ 1)
	next, err := loadBootstrapPointer(s, s.x[20], uint32(s.x[12]))
	if err != nil {
		return 0, fmt.Errorf("account compound-route C target: %w", err)
	}
	s.x[2] = next
	s.x[12] = s.x[0]
	s.x[17] = 0x512c
	return next, nil
}
