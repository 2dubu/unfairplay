package recovery

import "fmt"

const accountRecordWriterRoute = uint64(0x100a22a2c)

func runAccountRecordWriterRoute(s *state) (uint64, error) {
	s.traceBlock(accountRecordWriterRoute)
	value, err := s.load(s.sp+0x3f4, 4)
	if err != nil {
		return 0, fmt.Errorf("account record-writer value: %w", err)
	}
	s.x[12] = uint64(uint32(value))
	s.x[13] = 0x18
	s.x[11] = uint64(uint32(s.x[11])) * uint64(uint32(s.x[13]))
	s.x[13] = s.x[27] + 8
	s.x[14] = s.x[13] + s.x[11]
	if err := s.storeValue(s.x[14], uint64(uint32(s.x[25])), 4); err != nil {
		return 0, fmt.Errorf("account record-writer first word: %w", err)
	}
	if err := s.storeValue(s.x[14]+4, uint64(uint32(s.x[23])), 4); err != nil {
		return 0, fmt.Errorf("account record-writer second word: %w", err)
	}
	if err := s.storeValue(s.x[14]+8, s.x[6], 8); err != nil {
		return 0, fmt.Errorf("account record-writer pointer: %w", err)
	}
	if err := s.storeValue(s.x[14]+0x10, uint64(uint32(s.x[12])), 4); err != nil {
		return 0, fmt.Errorf("account record-writer trailing word: %w", err)
	}
	s.x[11] = uint64(uint32(s.x[9]) ^ 0x600e)
	s.x[10] = uint64(uint32(s.x[11]) + uint32(s.x[10]))
	if err := s.storeValue(s.x[27], uint64(uint32(s.x[10])), 4); err != nil {
		return 0, fmt.Errorf("account record-writer count: %w", err)
	}
	s.x[8] += 0x20
	value, err = s.load(s.x[8], 4)
	if err != nil {
		return 0, fmt.Errorf("account record-writer guard: %w", err)
	}
	s.x[25] = uint64(uint32(value))
	s.x[8] = 0xffffbfac
	s.x[8] = uint64(uint32(s.x[9]) + uint32(s.x[8]))
	s.x[10] = uint64(uint32(s.x[8]) ^ 0xd069c0b2)
	s.x[11] = uint64(uint32(s.x[25]) + 0x4bc21006)
	s.compare32(uint32(s.x[11]), uint32(s.x[10]))
	value, err = s.load(s.sp+0x1b0, 4)
	if err != nil {
		return 0, fmt.Errorf("account record-writer limit: %w", err)
	}
	s.x[10] = uint64(uint32(value))
	if s.condition("gt") {
		s.compare32(uint32(s.x[11]), uint32(s.x[10]))
	} else {
		s.setNZCV(0)
	}
	s.x[10] = uint64(uint32(boolValue(s.condition("lt"))))
	s.x[11] = 0xa3e3
	s.x[10] = uint64(uint32(s.x[10]) * uint32(s.x[11]))
	s.x[9] = uint64(uint32(s.x[10]) ^ uint32(s.x[9]))
	next, err := loadBootstrapPointer(s, s.x[30], uint32(s.x[9]))
	if err != nil {
		return 0, fmt.Errorf("account record-writer target: %w", err)
	}
	s.x[9] = next
	s.x[17] = 0xfd9d
	return next, nil
}
