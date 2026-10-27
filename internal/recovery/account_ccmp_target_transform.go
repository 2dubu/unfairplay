package recovery

import "fmt"

const accountCCmpTargetRoute = uint64(0x1002798e4)

func runAccountCCmpTargetRoute(s *state) (uint64, error) {
	s.traceBlock(accountCCmpTargetRoute)
	s.x[19] = uint64(uint32(s.x[8]) ^ uint32(s.x[6]))
	value, err := s.load(s.x[5], 4)
	if err != nil {
		return 0, fmt.Errorf("account ccmp-target input: %w", err)
	}
	s.x[10] = uint64(uint32(value))
	s.compare32(uint32(s.x[10]), uint32(s.x[3]))
	s.x[9] = uint64(uint32(s.x[8]) + uint32(s.x[7]))
	if s.condition("eq") {
		s.compare32(uint32(s.x[20]), uint32(s.x[9]))
	} else {
		s.setNZCV(0)
	}
	s.x[9] = uint64(uint32(boolValue(s.condition("ne"))))
	s.x[11] = uint64(uint32(s.x[8]) ^ uint32(s.x[30]))
	s.x[9] = uint64(uint32(s.x[9]) * uint32(s.x[11]))
	s.x[9] = uint64(uint32(s.x[9]) ^ uint32(s.x[8]))
	next, err := loadBootstrapPointer(s, s.x[24], uint32(s.x[9]))
	if err != nil {
		return 0, fmt.Errorf("account ccmp-target target: %w", err)
	}
	s.x[9] = next
	s.x[13] = 0xce757bc5
	s.x[17] = 0xa86b
	return next, nil
}
