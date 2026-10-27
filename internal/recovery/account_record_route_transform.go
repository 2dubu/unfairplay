package recovery

import "fmt"

const accountRecordRoute = uint64(0x1009f09d0)

func runAccountRecordRoute(s *state) (uint64, error) {
	s.traceBlock(accountRecordRoute)
	s.x[8] = uint64(uint32(s.x[9]) - 1)
	if err := s.storeValue(s.x[27], uint64(uint32(s.x[8])), 4); err != nil {
		return 0, fmt.Errorf("account record-route count: %w", err)
	}
	s.x[8] = 0x87a21b13
	s.x[8] = uint64(uint32(s.x[9]) + uint32(s.x[8]))
	s.x[14] = 0x18
	s.x[8] = uint64(uint32(s.x[8])) * uint64(uint32(s.x[14]))
	s.x[12] = s.x[27] + 8
	s.x[11] = s.x[12] + s.x[8]
	value, err := s.load(s.x[11], 4)
	if err != nil {
		return 0, fmt.Errorf("account record-route first pair word: %w", err)
	}
	s.x[8] = uint64(uint32(value))
	value, err = s.load(s.x[11]+4, 4)
	if err != nil {
		return 0, fmt.Errorf("account record-route second pair word: %w", err)
	}
	s.x[23] = uint64(uint32(value))
	value, err = s.load(s.x[11]+8, 8)
	if err != nil {
		return 0, fmt.Errorf("account record-route pair pointer: %w", err)
	}
	s.x[6] = value
	s.x[11] = uint64(uint32(s.x[10]) + 0xfb5891bd)
	s.x[13] = uint64(uint32(s.x[11]) ^ 0x87a27bd3)
	s.x[9] = uint64(uint32(s.x[13]) + uint32(s.x[9]))
	s.x[9] = uint64(uint32(s.x[9])) * uint64(uint32(s.x[14]))
	s.x[13] = s.x[12] + s.x[9]
	value, err = s.load(s.x[13]+0x10, 4)
	if err != nil {
		return 0, fmt.Errorf("account record-route selected word: %w", err)
	}
	s.x[9] = uint64(uint32(value))
	if err := s.storeValue(s.sp+0x3f4, uint64(uint32(s.x[9])), 4); err != nil {
		return 0, fmt.Errorf("account record-route selected output: %w", err)
	}
	s.x[9] = 0xfb5830fc
	s.x[8] = uint64(uint32(s.x[8]) + uint32(s.x[9]))
	s.x[8] = uint64(int64(int32(s.x[8])) * int64(int32(s.x[0])))
	s.x[9] = s.x[19] + s.x[8]
	value, err = s.load(s.x[9]+0x24, 4)
	if err != nil {
		return 0, fmt.Errorf("account record-route guard word: %w", err)
	}
	s.x[25] = uint64(uint32(value))
	s.x[8] = 0xfb588ada
	s.x[9] = uint64(uint32(s.x[10]) + uint32(s.x[8]))
	s.x[8] = 0x6f591c02
	s.x[8] = uint64(uint32(s.x[9]) + uint32(s.x[8]))
	s.x[8] = uint64(uint32(s.x[8]) + uint32(s.x[25]))
	value, err = s.load(s.sp+0x1c8, 4)
	if err != nil {
		return 0, fmt.Errorf("account record-route limit: %w", err)
	}
	s.x[12] = uint64(uint32(value))
	s.compare32(uint32(s.x[8]), uint32(s.x[12]))
	s.x[12] = 0xf40144e2
	if s.condition("lt") {
		s.compare32(uint32(s.x[8]), uint32(s.x[12]))
	} else {
		s.setNZCV(4)
	}
	s.x[8] = 0x318a201b
	if s.condition("gt") {
		s.x[8] = uint64(uint32(s.x[3]))
	}
	s.x[12] = uint64(uint32(s.x[25]) - uint32(s.x[10]))
	s.x[13] = uint64(uint32(s.x[10]) - uint32(s.x[25]))
	s.x[12] = uint64(uint32(s.x[12]) | uint32(s.x[13]))
	s.x[12] = uint64(^uint32(s.x[12]))
	s.x[12] = uint64(uint32(s.x[12]) >> 31)
	s.x[13] = 0xf69e
	s.x[12] = uint64(uint32(s.x[12]) * uint32(s.x[13]))
	s.x[13] = 0x04a7b64a
	s.x[10] = uint64(uint32(s.x[10]) ^ uint32(s.x[13]))
	s.x[10] = uint64(uint32(s.x[10]) ^ uint32(s.x[12]))
	next, err := loadBootstrapPointer(s, s.x[30], uint32(s.x[10]))
	if err != nil {
		return 0, fmt.Errorf("account record-route target: %w", err)
	}
	s.x[12] = next
	s.x[10] = s.x[25]
	s.x[17] = 0xfbad
	return next, nil
}
