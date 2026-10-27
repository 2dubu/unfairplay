package recovery

import "fmt"

const sharedEnrollmentByteGenerate = uint64(0x1004b1cf0)

func runSharedEnrollmentByteGenerate(s *state) (uint64, error) {
	s.traceBlock(sharedEnrollmentByteGenerate)
	s.x[15] = 0xd06c4952
	s.compare64(s.x[12], s.x[15])
	firstBound := uint32(boolValue(s.condition("lo")))
	s.x[16] = uint64(firstBound)
	s.x[10] = uint64(uint32(s.x[8]) + 0x8ef060c1)
	s.x[10] = uint64(uint32(s.x[10]) | 0x00040169)
	s.x[14]--
	if err := s.storeValue(s.sp+0x5000, s.x[14], 8); err != nil {
		return 0, fmt.Errorf("shared byte-generate counter: %w", err)
	}
	s.x[17] = uint64(uint32(s.x[13]) ^ 0xc384deba)
	s.x[0] = 0xa49596e9
	s.x[1] = uint64(uint32(s.x[13]) ^ 0x6076ccc0)
	s.x[1] = uint64(uint32(s.x[1]) + 0x07678493)
	s.x[2] = uint64(uint32(s.x[13]) ^ 0x70)
	s.x[3] = 0xa1a2d223
	s.x[4] = uint64(uint32(s.x[10]) + 0x6eba9e3c)
	s.x[4] = uint64(uint32(s.x[4]) ^ uint32(s.x[13]))
	s.x[4] = uint64(uint32(s.x[4]) + 0xf6)
	s.x[5] = uint64(uint32(s.x[13]) + 0x64)
	s.x[17] = uint64(uint32(s.x[17]) + uint32(s.x[0]))
	s.x[17] = uint64(uint32(s.x[13]) ^ uint32(s.x[17]))
	s.x[0] = uint64(uint32(s.x[2]) + uint32(s.x[3]))
	s.x[3] = 0xffffffff710b5e96
	s.x[0] = uint64(uint32(s.x[1]) ^ uint32(s.x[0]))
	s.x[1] = 0xba
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[0]))
	s.x[0] = 0x17
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[0]))
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[4]))
	s.x[17] = uint64(uint32(s.x[17]) * uint32(s.x[5]))
	if err := s.storeValue(s.x[9]+s.x[14], uint64(byte(s.x[17])), 1); err != nil {
		return 0, fmt.Errorf("shared byte-generate output: %w", err)
	}
	s.x[13]++
	counter, err := s.load(s.sp+0x5000, 8)
	if err != nil {
		return 0, fmt.Errorf("shared byte-generate saved counter: %w", err)
	}
	s.x[14] = counter
	s.compare64(s.x[14], s.x[11])
	secondBound := uint32(boolValue(s.condition("hi")))
	s.x[17] = uint64(secondBound)
	s.x[16] = uint64(firstBound ^ secondBound)
	s.x[15] = s.x[14] + 0xd06c4952
	s.compare64(s.x[15], s.x[12])
	thirdBound := uint32(boolValue(s.condition("hi")))
	s.x[15] = uint64(thirdBound)
	s.compare32(uint32(s.x[16]), 0)
	dispatchGate := thirdBound
	if s.condition("ne") {
		dispatchGate = secondBound
	}
	s.x[16] = 0x89c3
	s.x[15] = uint64(dispatchGate * uint32(s.x[16]))
	s.x[15] = uint64(uint32(s.x[15]) ^ uint32(s.x[8]))
	dispatchTable, err := s.load(s.sp+0xdf0, 8)
	if err != nil {
		return 0, fmt.Errorf("shared byte-generate dispatch table: %w", err)
	}
	s.x[16] = dispatchTable
	next, err := loadBootstrapPointer(s, dispatchTable, uint32(s.x[15]))
	if err != nil {
		return 0, fmt.Errorf("shared byte-generate dispatch: %w", err)
	}
	s.x[15] = next
	s.x[17] = 0x8686
	return next, nil
}
