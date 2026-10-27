package recovery

import "fmt"

const accountRoundSetupRoute = uint64(0x10005263c)

func runAccountRoundSetupRoute(s *state) (uint64, error) {
	s.traceBlock(accountRoundSetupRoute)
	s.x[2] = s.x[19]
	s.x[8] = uint64(uint32(s.x[22]) ^ 0x24ce)
	s.x[9] = 0x1c4
	s.x[12] = uint64(uint32(s.x[8]) * uint32(s.x[9]))
	value, err := s.load(s.sp+0x6c, 4)
	if err != nil {
		return 0, fmt.Errorf("account round-setup word: %w", err)
	}
	s.x[8] = uint64(uint32(value))
	s.x[9] = uint64(uint32(s.x[8]) ^ 0xdf80aef2)
	s.x[10] = 0x207f510e
	s.x[9] = uint64(uint32(s.x[9]) + uint32(s.x[10]))
	s.x[10] = uint64(uint32(s.x[8]) ^ 0x86909a80)
	s.x[11] = 0x796f6580
	s.x[10] = uint64(uint32(s.x[10]) + uint32(s.x[11]))
	s.x[9] = uint64(uint32(s.x[9]) ^ uint32(s.x[10]))
	s.x[10] = uint64(uint32(s.x[8]) ^ 0x5df044f5)
	s.x[10] = uint64(uint32(s.x[10]) ^ uint32(s.x[12]))
	s.x[11] = 0xa20f883f
	s.x[10] = uint64(uint32(s.x[10]) + uint32(s.x[11]))
	s.x[9] = uint64(uint32(s.x[9]) ^ uint32(s.x[10]))
	s.x[10] = 0x79921567
	s.x[9] = uint64(uint32(s.x[9]) + uint32(s.x[10]))
	next, err := loadBootstrapPointer(s, s.x[20], uint32(s.x[22]))
	if err != nil {
		return 0, fmt.Errorf("account round-setup target: %w", err)
	}
	s.x[10] = next
	s.x[22] = 0xfc79d912
	s.x[21] = 0x7c79d912
	s.x[23] = 0x80000004
	s.x[17] = 0x44e6
	return next, nil
}
