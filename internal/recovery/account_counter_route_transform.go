package recovery

import "fmt"

const (
	accountCounterRoute    = uint64(0x1008efe4c)
	accountSignedByteRoute = uint64(0x100a1ca10)
)

func runAccountCounterRoute(s *state) (uint64, error) {
	s.traceBlock(accountCounterRoute)
	s.x[11] = uint64(uint32(s.x[14]) + uint32(s.x[5]))
	s.x[9] = uint64(uint32(s.x[14]) + uint32(s.x[20]))
	s.x[9] = uint64(uint32(s.x[9]) & uint32(s.x[21]))
	s.x[8]++
	s.x[10]++
	s.compare32(uint32(s.x[12]), 0)
	s.x[12] = uint64(uint32(boolValue(s.condition("eq"))))
	s.x[13] = uint64(uint32(s.x[9]) - 0x606)
	s.x[12] = uint64(uint32(s.x[13]) * uint32(s.x[12]))
	s.x[12] = uint64(uint32(s.x[12]) ^ uint32(s.x[14]))
	next, err := loadBootstrapPointer(s, s.x[30], uint32(s.x[12]))
	if err != nil {
		return 0, fmt.Errorf("account counter-route target: %w", err)
	}
	s.x[13] = next
	s.x[12] = 0x7b5db20a
	s.x[17] = 0x829c
	return next, nil
}

func runAccountSignedByteRoute(s *state) (uint64, error) {
	s.traceBlock(accountSignedByteRoute)
	s.x[13] = uint64(uint32(s.x[11]) + uint32(s.x[2]))
	value, err := s.load(s.x[10], 1)
	if err != nil {
		return 0, fmt.Errorf("account signed-byte first input: %w", err)
	}
	s.x[12] = uint64(uint32(int32(int8(value))))
	value, err = s.load(s.x[8], 1)
	if err != nil {
		return 0, fmt.Errorf("account signed-byte second input: %w", err)
	}
	s.x[15] = uint64(uint32(int32(int8(value))))
	s.x[9] = uint64(uint32(s.x[15]) ^ uint32(s.x[12]))
	s.x[9] &= 0xff
	s.x[14] = uint64(uint32(s.x[11]) + uint32(s.x[3]))
	s.x[16] = uint64(uint32(s.x[14]) + uint32(s.x[13]))
	s.x[16] = uint64(uint32(s.x[16]) + 0x7e)
	s.compare32(uint32(s.x[9]), uint32(byte(s.x[16])))
	s.x[9] = uint64(uint32(boolValue(s.condition("eq"))))
	s.x[9] = uint64(uint32(s.x[9]) * uint32(s.x[4]))
	s.x[9] = uint64(uint32(s.x[9]) ^ uint32(s.x[11]))
	next, err := loadBootstrapPointer(s, s.x[30], uint32(s.x[9]))
	if err != nil {
		return 0, fmt.Errorf("account signed-byte target: %w", err)
	}
	s.x[9] = next
	s.x[17] = 0x1791
	return next, nil
}
