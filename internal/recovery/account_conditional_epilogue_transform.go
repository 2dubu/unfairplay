package recovery

import "fmt"

const accountConditionalEpilogue = uint64(0x10011598c)

func runAccountConditionalEpilogue(s *state) (uint64, error) {
	s.traceBlock(accountConditionalEpilogue)
	value, err := s.load(addOffset(s.x[29], -0x60), 8)
	if err != nil {
		return 0, fmt.Errorf("account conditional epilogue saved value: %w", err)
	}
	s.x[8] = value
	s.x[9] = 0x10106d000
	value, err = s.load(s.x[9]+0x408, 8)
	if err != nil {
		return 0, fmt.Errorf("account conditional epilogue guard: %w", err)
	}
	s.x[9] = value
	value, err = s.load(s.x[9], 8)
	if err != nil {
		return 0, fmt.Errorf("account conditional epilogue guard value: %w", err)
	}
	s.x[9] = value
	s.compare64(s.x[9], s.x[8])
	if s.condition("ne") {
		return 0x1001159d4, nil
	}

	s.traceBlock(0x1001159a4)
	s.x[8] = 0xce757bc5
	value, err = s.load(s.sp+0x34, 4)
	if err != nil {
		return 0, fmt.Errorf("account conditional epilogue result: %w", err)
	}
	s.x[9] = uint64(uint32(value))
	s.x[0] = uint64(uint32(s.x[9]) ^ uint32(s.x[8]))
	for _, pair := range []struct {
		first, second int
		offset        uint64
	}{
		{29, 30, 0x160}, {20, 19, 0x150}, {22, 21, 0x140},
		{24, 23, 0x130}, {26, 25, 0x120}, {28, 27, 0x110},
	} {
		first, err := s.load(s.sp+pair.offset, 8)
		if err != nil {
			return 0, fmt.Errorf("account conditional epilogue x%d: %w", pair.first, err)
		}
		second, err := s.load(s.sp+pair.offset+8, 8)
		if err != nil {
			return 0, fmt.Errorf("account conditional epilogue x%d: %w", pair.second, err)
		}
		s.x[pair.first], s.x[pair.second] = first, second
	}
	s.sp += 0x170
	return s.x[30], nil
}
