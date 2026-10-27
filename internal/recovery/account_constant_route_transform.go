package recovery

import "fmt"

const (
	accountConstantRouteA = uint64(0x10000f828)
	accountConstantRouteB = uint64(0x1000134f0)
	accountConstantRouteC = uint64(0x100035e94)
)

func runAccountConstantRouteA(s *state) (uint64, error) {
	s.traceBlock(accountConstantRouteA)
	s.x[9] = 0xffffc335
	s.x[9] = uint64(uint32(s.x[12]) + uint32(s.x[9]))
	s.x[11] = uint64(uint32(s.x[9]) ^ 0xf3)
	_ = s.addWithCarry(uint64(uint32(s.x[16])), 0x18, false, 32)
	s.x[13] = uint64(uint32(boolValue(s.condition("gt"))))
	s.x[15] = 0xffff84bc
	s.x[15] = uint64(uint32(s.x[11]) + uint32(s.x[15]))
	s.x[13] = uint64(uint32(s.x[15]) * uint32(s.x[13]))
	s.x[12] = uint64(uint32(s.x[13]) ^ uint32(s.x[12]))
	next, err := loadBootstrapPointer(s, s.x[21], uint32(s.x[12]))
	if err != nil {
		return 0, fmt.Errorf("account constant-route A target: %w", err)
	}
	s.x[12] = next
	s.x[17] = 0x117d
	return next, nil
}

func runAccountConstantRouteB(s *state) (uint64, error) {
	s.traceBlock(accountConstantRouteB)
	s.x[9] = 0x527fd2d7
	s.x[9] = uint64(uint32(s.x[11]) ^ uint32(s.x[9]))
	s.x[12] = 0x1f6
	s.x[12] = uint64(uint32(s.x[9]) * uint32(s.x[12]))
	s.x[9] = 0xbc00b302
	s.compare32(uint32(s.x[11]), uint32(s.x[9]))
	s.x[9] = uint64(uint32(boolValue(s.condition("hi"))))
	s.x[13] = uint64(uint32(s.x[12]) ^ 0x6255)
	s.x[9] = uint64(uint32(s.x[13]) * uint32(s.x[9]))
	s.x[13] = 0xd27f9f9a
	s.x[11] = uint64(uint32(s.x[11]) ^ uint32(s.x[13]))
	s.x[9] = uint64(uint32(s.x[11]) ^ uint32(s.x[9]))
	next, err := loadBootstrapPointer(s, s.x[21], uint32(s.x[9]))
	if err != nil {
		return 0, fmt.Errorf("account constant-route B target: %w", err)
	}
	s.x[9] = next
	s.x[17] = 0xb3ab
	return next, nil
}

func runAccountConstantRouteC(s *state) (uint64, error) {
	s.traceBlock(accountConstantRouteC)
	s.x[9] = 0xe28a7dc2
	s.x[9] = uint64(uint32(s.x[11]) + uint32(s.x[9]))
	s.x[12] = 0x1d758fef
	s.x[24] = uint64(uint32(s.x[9]) & uint32(s.x[12]))
	_ = s.addWithCarry(uint64(uint32(s.x[16])), 0x18, false, 32)
	s.x[12] = uint64(uint32(boolValue(s.condition("ne"))))
	s.x[9] = 0xffffe1c8
	s.x[9] = uint64(uint32(s.x[11]) + uint32(s.x[9]))
	s.x[13] = 0xffff1a85
	s.x[13] = uint64(uint32(s.x[9]) + uint32(s.x[13]))
	s.x[13] = uint64(uint32(s.x[13]) + uint32(s.x[24]))
	s.x[12] = uint64(uint32(s.x[13]) * uint32(s.x[12]))
	s.x[11] = uint64(uint32(s.x[12]) ^ uint32(s.x[11]))
	next, err := loadBootstrapPointer(s, s.x[21], uint32(s.x[11]))
	if err != nil {
		return 0, fmt.Errorf("account constant-route C target: %w", err)
	}
	s.x[15] = next
	s.x[4] = 0x2ba7c38b
	s.x[11] = s.x[3]
	s.x[13] = s.x[22]
	s.x[12] = s.x[26]
	s.x[16] = s.x[0]
	s.x[17] = 0x169e
	return next, nil
}
