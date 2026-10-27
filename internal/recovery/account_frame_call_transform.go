package recovery

import "fmt"

const accountFrameCall = uint64(0x10091e464)

func runAccountFrameCall(s *state) (uint64, error) {
	s.traceBlock(accountFrameCall)
	s.x[8] = uint64(uint32(s.x[28]) + 0xf51e372e)
	s.x[11] = 0x0ae1be7d
	s.x[20] = uint64(uint32(s.x[8]) & uint32(s.x[11]))

	value, err := s.load(s.sp+0x158, 8)
	if err != nil {
		return 0, fmt.Errorf("account frame-call first output pointer: %w", err)
	}
	s.x[1] = value
	if err := s.storeValue(s.x[1], uint64(uint32(s.x[9])), 4); err != nil {
		return 0, fmt.Errorf("account frame-call first output: %w", err)
	}
	value, err = s.load(s.sp+0x3ec, 4)
	if err != nil {
		return 0, fmt.Errorf("account frame-call first word: %w", err)
	}
	s.x[8] = uint64(uint32(value))
	if err := s.storeValue(s.x[27]+0x18c, s.x[8], 4); err != nil {
		return 0, fmt.Errorf("account frame-call second output: %w", err)
	}
	value, err = s.load(s.sp+0x2c0, 8)
	if err != nil {
		return 0, fmt.Errorf("account frame-call first pointer: %w", err)
	}
	s.x[8] = value
	if err := s.storeValue(s.x[27]+0x190, s.x[8], 8); err != nil {
		return 0, fmt.Errorf("account frame-call third output: %w", err)
	}
	value, err = s.load(s.sp+0x3f4, 4)
	if err != nil {
		return 0, fmt.Errorf("account frame-call second word: %w", err)
	}
	s.x[8] = uint64(uint32(value))

	s.x[8] = uint64(uint32(s.x[8]) - uint32(s.x[21]))
	s.x[9] = 0x290f0c70
	s.x[8] = uint64(uint32(s.x[8]) + uint32(s.x[9]))
	s.x[9] = uint64(uint32(s.x[8]) ^ 0xb86a5731)
	s.x[11] = uint64(uint32(s.x[8]) & 0xd6f0f390)
	s.x[12] = uint64(uint32(s.x[8]) ^ 0x94eaf311)
	s.x[12] = uint64(uint32(s.x[12]) & uint32(s.x[11]<<1))
	s.x[11] = uint64(uint32(s.x[12]) ^ uint32(s.x[11]))
	s.x[12] = uint64(uint32(s.x[9]) << 1)
	s.x[13] = uint64(uint32(s.x[20]) + 0xdd351ae9)
	s.x[13] = uint64(uint32(s.x[12]) ^ uint32(s.x[13]))
	s.x[9] = uint64(uint32(s.x[13]) & uint32(s.x[9]))
	s.x[14] = 0x6e9aa4a0
	s.x[12] = uint64(uint32(s.x[12]) & uint32(s.x[14]))
	s.x[9] = uint64(uint32(s.x[9]) ^ uint32(s.x[12]))
	s.x[12] = uint64(uint32(s.x[9]) ^ 0x228aa4a1)
	s.x[13] = 0x4c000000
	s.x[9] = uint64(uint32(s.x[9]) ^ uint32(s.x[13]))
	s.x[9] = uint64(uint32(s.x[9]) & uint32(s.x[11]<<2))
	s.x[9] = uint64(uint32(s.x[9]) ^ uint32(s.x[11]))
	s.x[11] = uint64(uint32(s.x[12]) << 2)
	s.x[13] = uint64(uint32(s.x[11]) ^ 0xba6a9284)
	s.x[12] = uint64(uint32(s.x[13]) & uint32(s.x[12]))
	s.x[11] = uint64(uint32(s.x[11]) & uint32(s.x[14]))
	s.x[11] = uint64(uint32(s.x[12]) ^ uint32(s.x[11]))
	s.x[12] = uint64(uint32(s.x[11]) ^ 0x44902421)
	s.x[13] = 0x2a0a8080
	s.x[11] = uint64(uint32(s.x[11]) ^ uint32(s.x[13]))
	s.x[11] = uint64(uint32(s.x[11]) & uint32(s.x[9]<<4))
	s.x[9] = uint64(uint32(s.x[11]) ^ uint32(s.x[9]))
	s.x[11] = uint64(uint32(s.x[12]) << 4)
	s.x[13] = uint64(uint32(s.x[11]) ^ 0xe9aa4a10)
	s.x[12] = uint64(uint32(s.x[13]) & uint32(s.x[12]))
	s.x[13] = 0x6e9aa480
	s.x[11] = uint64(uint32(s.x[11]) & uint32(s.x[13]))
	s.x[11] = uint64(uint32(s.x[12]) ^ uint32(s.x[11]))
	s.x[12] = uint64(uint32(s.x[11]) ^ 0x0610a4a1)
	s.x[13] = 0x688a0000
	s.x[11] = uint64(uint32(s.x[11]) ^ uint32(s.x[13]))
	s.x[11] = uint64(uint32(s.x[11]) & uint32(s.x[9]<<8))
	s.x[13] = 0x6e9aa4a1
	s.x[9] = uint64(uint32(s.x[9]) ^ uint32(s.x[13]))
	s.x[9] = uint64(uint32(s.x[9]) ^ uint32(s.x[11]))
	s.x[11] = uint64(uint32(s.x[12]) << 8)
	s.x[13] = 0x1aa40000
	s.x[13] = uint64(uint32(s.x[11]) ^ uint32(s.x[13]))
	s.x[12] = uint64(uint32(s.x[13]) & uint32(s.x[12]))
	s.x[14] = 0x6e9a0000
	s.x[11] = uint64(uint32(s.x[11]) & uint32(s.x[14]))
	s.x[13] = 0x641a0000
	s.x[11] = uint64(uint32(s.x[11]) ^ uint32(s.x[13]))
	s.x[11] = uint64(uint32(s.x[11]) ^ uint32(s.x[12]))
	s.x[12] = uint64(uint32(s.x[9]) << 16)
	s.x[13] = 0x24a10000
	s.x[13] = uint64(uint32(s.x[12]) ^ uint32(s.x[13]))
	s.x[11] = uint64(uint32(s.x[13]) & uint32(s.x[11]))
	s.x[12] = uint64(uint32(s.x[12]) & uint32(s.x[14]))
	s.x[9] = uint64(uint32(s.x[12]) ^ uint32(s.x[9]))
	s.x[9] = uint64(uint32(s.x[9]) ^ uint32(s.x[11]))
	s.x[8] = uint64(uint32(s.x[8]) ^ uint32(s.x[9]<<1))
	s.x[8] = uint64(uint32(s.x[8]) ^ 0x8cb0c117)
	if err := s.storeValue(s.x[27]+0x198, s.x[8], 4); err != nil {
		return 0, fmt.Errorf("account frame-call fourth output: %w", err)
	}

	s.x[8] = uint64(uint32(s.x[21]) + 0xd94b8418)
	s.x[8] = s.x[6] + s.x[8]
	if err := s.storeValue(s.x[27]+0x1a0, s.x[8], 8); err != nil {
		return 0, fmt.Errorf("account frame-call fifth output: %w", err)
	}
	s.x[8] = 0xd59bf993
	s.x[22] = uint64(uint32(s.x[20]) + uint32(s.x[8]))
	next, err := s.load(s.x[10]+8, 8)
	if err != nil {
		return 0, fmt.Errorf("account frame-call target: %w", err)
	}
	s.x[8] = next
	value, err = s.load(s.sp+0x1a8, 8)
	if err != nil {
		return 0, fmt.Errorf("account frame-call argument: %w", err)
	}
	s.x[0] = value
	s.x[30] = accountCallReturn
	return next, nil
}
