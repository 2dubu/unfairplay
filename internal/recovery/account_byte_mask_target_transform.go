package recovery

import "fmt"

const accountByteMaskTarget = uint64(0x1008ebda8)

func runAccountByteMaskTarget(s *state) (uint64, error) {
	s.traceBlock(accountByteMaskTarget)
	s.x[13] &= 0xff
	s.x[8] = 0x6cddb847
	s.x[8] = uint64(uint32(s.x[12]) + uint32(s.x[8]))
	s.x[9] = 0x93227ddf
	s.x[8] = uint64(uint32(s.x[8]) & uint32(s.x[9]))
	s.x[9] = uint64(uint32(s.x[12]) ^ uint32(s.x[26]))
	s.x[10]++
	s.x[11]++
	s.x[14] = uint64(uint32(s.x[8]) + 0x23)
	s.compare32(uint32(s.x[13]), uint32(byte(s.x[14])))
	s.x[13] = uint64(uint32(boolValue(s.condition("eq"))))
	s.x[14] = uint64(uint32(s.x[9]) + 0x3af5)
	s.x[13] = uint64(uint32(s.x[13]) * uint32(s.x[14]))
	s.x[12] = uint64(uint32(s.x[13]) ^ uint32(s.x[12]))
	next, err := loadBootstrapPointer(s, s.x[30], uint32(s.x[12]))
	if err != nil {
		return 0, fmt.Errorf("account byte-mask target: %w", err)
	}
	s.x[13] = next
	s.x[12] = 0x7b5db20a
	s.x[17] = 0xa654
	return next, nil
}
