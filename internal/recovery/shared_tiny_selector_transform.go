package recovery

import "fmt"

const (
	sharedEnrollmentRegisterReturn = uint64(0x10030723c)
	sharedEnrollmentBoundsSelect   = uint64(0x10028723c)
	sharedEnrollmentProductSelect  = uint64(0x1002da2f4)
)

func runSharedEnrollmentRegisterReturn(s *state) (uint64, error) {
	s.traceBlock(sharedEnrollmentRegisterReturn)
	s.x[17] = 0x8b3c
	return s.x[6], nil
}

func runSharedEnrollmentBoundsSelect(s *state) (uint64, error) {
	s.traceBlock(sharedEnrollmentBoundsSelect)
	s.x[9] = 0x180b67ab
	s.x[12] = uint64(uint32(s.x[11]) + uint32(s.x[9]))
	s.x[9] = 0xe7f46d79
	s.x[3] = uint64(uint32(s.x[12]) & uint32(s.x[9]))
	s.x[9] = 0x180b9fdb
	s.x[12] = uint64(uint32(s.x[12]) ^ uint32(s.x[9]))
	s.x[9] = 0xa2bd6725
	s.x[16] = uint64(uint32(s.x[1]) + uint32(s.x[9]))
	s.compare32(uint32(s.x[16]), uint32(s.x[12]))
	if s.condition("hs") {
		s.compare32(uint32(s.x[16]), uint32(s.x[23]))
	} else {
		s.setNZCV(2)
	}
	gate := uint32(boolValue(s.condition("hi")))
	s.x[12] = uint64(gate)
	s.x[9] = 0xd43
	s.x[16] = uint64(uint32(s.x[3]) ^ uint32(s.x[9]))
	s.x[12] = uint64(gate * uint32(s.x[16]))
	s.x[11] = uint64(uint32(s.x[12]) ^ uint32(s.x[11]))
	next, err := loadBootstrapPointer(s, s.x[14], uint32(s.x[11]))
	if err != nil {
		return 0, fmt.Errorf("shared bounds-selector target: %w", err)
	}
	s.x[11] = next
	s.x[2] = 0x5d4298ef
	s.x[17] = 0x4b0e
	return next, nil
}

func runSharedEnrollmentProductSelect(s *state) (uint64, error) {
	s.traceBlock(sharedEnrollmentProductSelect)
	s.x[9] = 0xd4b
	s.x[11] = uint64(uint32(s.x[3]) ^ uint32(s.x[9]))
	s.x[9] = 0x377d
	s.x[16] = uint64(uint32(s.x[11]) * uint32(s.x[9]))
	s.x[9] = 0x653
	s.x[11] = uint64(uint32(s.x[3]) ^ uint32(s.x[9]))
	s.x[11] = uint64(uint32(s.x[11]) ^ uint32(s.x[16]))
	next, err := loadBootstrapPointer(s, s.x[14], uint32(s.x[11]))
	if err != nil {
		return 0, fmt.Errorf("shared product-selector target: %w", err)
	}
	s.x[11] = next
	s.x[17] = 0xbd21
	return next, nil
}
