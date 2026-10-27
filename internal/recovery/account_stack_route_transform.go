package recovery

import "fmt"

const (
	accountStackRouteSetup = uint64(0x100139734)
	accountStackRouteSave  = uint64(0x1001b7670)
	accountStackRouteLoad  = uint64(0x1001d1498)
)

func runAccountStackRouteSetup(s *state) (uint64, error) {
	s.traceBlock(accountStackRouteSetup)
	s.x[8] = 0xa96ed437
	s.x[8] = uint64(uint32(s.x[19]) + uint32(s.x[8]))
	value, err := s.load(s.sp+0x80, 8)
	if err != nil {
		return 0, fmt.Errorf("account stack-route setup pointer: %w", err)
	}
	s.x[26] = value
	s.x[9] = 0xa96f085f
	s.x[13] = uint64(uint32(s.x[19]) + uint32(s.x[9]))
	next, err := loadBootstrapPointer(s, s.x[20], uint32(s.x[13]))
	if err != nil {
		return 0, fmt.Errorf("account stack-route setup target: %w", err)
	}
	s.x[13] = next
	if err := s.storeValue(s.sp+0x3c, uint64(uint32(s.x[11])), 4); err != nil {
		return 0, fmt.Errorf("account stack-route setup word: %w", err)
	}
	s.x[21] = s.x[12]
	s.x[12] = 0xaf
	s.x[17] = 0x8b46
	return next, nil
}

func runAccountStackRouteSave(s *state) (uint64, error) {
	s.traceBlock(accountStackRouteSave)
	if err := s.storeValue(s.sp+0x70, s.x[12], 8); err != nil {
		return 0, fmt.Errorf("account stack-route saved value: %w", err)
	}
	if err := s.storeValue(s.sp+0x6c, uint64(uint32(s.x[16])), 4); err != nil {
		return 0, fmt.Errorf("account stack-route saved word: %w", err)
	}
	s.x[11] = uint64(uint32(s.x[28]) + uint32(s.x[13]))
	s.x[12] = uint64(uint32(s.x[11]) | uint32(s.x[14]))
	s.x[11] = uint64(uint32(s.x[28]) + uint32(s.x[15]))
	next, err := loadBootstrapPointer(s, s.x[27], uint32(s.x[11]))
	if err != nil {
		return 0, fmt.Errorf("account stack-route save target: %w", err)
	}
	s.x[13] = next
	s.x[11] = s.x[0]
	s.x[17] = 0x5406
	return next, nil
}

func runAccountStackRouteLoad(s *state) (uint64, error) {
	s.traceBlock(accountStackRouteLoad)
	value, err := s.load(s.sp+0x18, 8)
	if err != nil {
		return 0, fmt.Errorf("account stack-route load pointer: %w", err)
	}
	s.x[8] = value
	value, err = s.load(s.x[8], 8)
	if err != nil {
		return 0, fmt.Errorf("account stack-route load target: %w", err)
	}
	s.x[8] = value
	value, err = s.load(s.sp+0x3c, 4)
	if err != nil {
		return 0, fmt.Errorf("account stack-route load word: %w", err)
	}
	s.x[11] = uint64(uint32(value))
	s.x[17] = 0x1c20
	return s.x[8], nil
}
