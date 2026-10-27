package recovery

import "fmt"

const bootstrapByteStreamPass = uint64(0x10087d0c4)

func runBootstrapByteStream(s *state) (uint64, error) {
	const maximumBytes = 4096
	for index := 0; index < maximumBytes; index++ {
		next, err := runBootstrapByteStreamStep(s)
		if err != nil || next != bootstrapByteStreamPass {
			return next, err
		}
	}
	return 0, fmt.Errorf("bootstrap byte stream exceeded %d bytes", maximumBytes)
}

func runBootstrapByteStreamStep(s *state) (uint64, error) {
	s.traceBlock(bootstrapByteStreamPass)
	s.x[17] = s.x[1] - 1
	if err := s.storeValue(s.sp+0x3a0, s.x[17], 8); err != nil {
		return 0, fmt.Errorf("bootstrap byte-stream counter: %w", err)
	}
	s.x[1] = uint64(uint32(s.x[10]) ^ uint32(s.x[13]))
	s.x[5] = uint64(uint32(s.x[10]) ^ uint32(s.x[15]))
	s.x[5] = uint64(uint32(s.x[11])*uint32(s.x[12]) + uint32(s.x[5]))
	s.x[6] = uint64(uint32(s.x[10]) ^ uint32(s.x[0]))
	s.x[7] = uint64(^uint32(s.x[10]) + 0x88)
	s.x[1] = uint64(uint32(s.x[1]) + uint32(s.x[14]))
	s.x[1] = uint64(uint32(s.x[10]) ^ uint32(s.x[1]))
	s.x[6] = uint64(uint32(s.x[6]) + uint32(s.x[2]))
	s.x[6] = uint64(uint32(s.x[6]) ^ uint32(s.x[7]))
	s.x[7] = s.x[25]
	s.x[1] = uint64(uint32(s.x[1]) ^ uint32(s.x[6]))
	s.x[1] = uint64(uint32(s.x[1]) ^ uint32(s.x[3]))
	s.x[5] = uint64(uint32(s.x[5]) + uint32(s.x[16]))
	s.x[1] = uint64(uint32(s.x[1]) ^ uint32(s.x[5]))
	s.x[5] = uint64(uint32(s.x[10]) + 0x8a)
	s.x[1] = uint64(uint32(s.x[1]) * uint32(s.x[5]))
	if err := s.storeValue(s.x[9]+s.x[17], uint64(byte(s.x[1])), 1); err != nil {
		return 0, fmt.Errorf("bootstrap byte-stream output: %w", err)
	}
	s.x[10]++
	counter, err := s.load(s.sp+0x3a0, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap byte-stream saved counter: %w", err)
	}
	s.x[1] = counter
	s.compare64(counter, 0)
	dispatchIndex := uint32(boolValue(s.condition("eq")))*uint32(s.x[4]) ^ uint32(s.x[8])
	next, err := loadBootstrapPointer(s, s.x[30], dispatchIndex)
	if err != nil {
		return 0, fmt.Errorf("bootstrap byte-stream dispatch: %w", err)
	}
	s.x[5] = next
	s.x[17] = 0xb22f
	return next, nil
}
