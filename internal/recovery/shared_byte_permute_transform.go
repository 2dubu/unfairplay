package recovery

import "fmt"

const sharedEnrollmentBytePermute = uint64(0x1003e9534)

func runSharedEnrollmentBytePermute(s *state) (uint64, error) {
	s.traceBlock(sharedEnrollmentBytePermute)
	s.x[17] = s.x[20] ^ s.x[8]
	source, err := s.load(s.sp+0x50, 8)
	if err != nil {
		return 0, fmt.Errorf("shared byte-permute source: %w", err)
	}
	s.x[20] = source
	value, err := s.load(source+s.x[17], 1)
	if err != nil {
		return 0, fmt.Errorf("shared byte-permute first input: %w", err)
	}
	s.x[17] = uint64(byte(value))
	s.x[3] = 0xffffffe8
	s.x[20] = uint64(uint32(s.x[2]) ^ uint32(s.x[3]))
	s.x[20] = uint64(uint32(s.x[17]) ^ uint32(s.x[20]))
	s.x[20] = uint64(uint32(s.x[20]) - 0x3e)
	s.x[22] = uint64(uint32(s.x[17]) ^ 0x49)
	s.x[22] = uint64(uint32(s.x[22]) - 0x49)
	s.x[17] = uint64(uint32(s.x[17]) ^ 0x6c)
	s.x[17] = uint64(uint32(s.x[17]) - 0x6c)
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[22]))
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[20]))
	s.x[17] = uint64(uint32(s.x[17]) + uint32(s.x[17]<<2))
	s.x[3] = 0xfffffff6
	s.x[17] = uint64(uint32(s.x[3]) - uint32(s.x[17]))
	s.x[20] = uint64((uint32(s.x[17]) & 0x40) ^ 0x74)
	s.x[22] = uint64(uint32(s.x[17]) ^ 0x18)
	s.x[23] = uint64((uint32(s.x[22]) & 0x28) ^ uint32(s.x[20]))
	for range 2 {
		s.x[23] = uint64(uint32(s.x[22]) & uint32(s.x[23]<<1))
		s.x[23] = uint64(uint32(s.x[23]) ^ uint32(s.x[20]))
	}
	s.x[22] = uint64(uint32(s.x[22]) & uint32(s.x[23]<<1))
	s.x[20] = uint64(uint32(s.x[22]) ^ uint32(s.x[20]))
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[20]<<1))
	s.x[20] = uint64(uint32(s.x[17]) ^ 0x5b)
	work, err := s.load(s.sp+0x58, 8)
	if err != nil {
		return 0, fmt.Errorf("shared byte-permute work buffer: %w", err)
	}
	s.x[27] = work
	if err := s.storeValue(work+s.x[0], uint64(byte(s.x[20])), 1); err != nil {
		return 0, fmt.Errorf("shared byte-permute first output: %w", err)
	}

	s.x[20] = uint64(uint32(s.x[11]) - 6)
	s.x[22] = uint64(uint32(s.x[20]) ^ 0x3ff9651d)
	s.x[23] = uint64(uint32(s.x[20]) & uint32(s.x[5]))
	s.x[24] = uint64(uint32(s.x[20]) ^ 0x3bd95c0d)
	s.x[24] = uint64(uint32(s.x[24]) & uint32(s.x[23]<<1))
	s.x[23] = uint64(uint32(s.x[24]) ^ uint32(s.x[23]))
	s.x[24] = uint64(uint32(s.x[22] << 1))
	s.x[25] = uint64(uint32(s.x[24]) ^ 0x0bcaf366)
	s.x[22] = uint64(uint32(s.x[25]) & uint32(s.x[22]))
	s.x[24] = uint64(uint32(s.x[24]) & 0x85e579b2)
	s.x[22] = uint64(uint32(s.x[22]) ^ uint32(s.x[24]))
	s.x[24] = uint64(uint32(s.x[22]) ^ 0x84250891)
	s.x[22] = uint64(uint32(s.x[22]) ^ 0x00c07120)
	s.x[22] = uint64(uint32(s.x[22]) & uint32(s.x[23]<<2))
	s.x[22] = uint64(uint32(s.x[22]) ^ uint32(s.x[23]))
	s.x[23] = uint64(uint32(s.x[24] << 2))
	s.x[25] = uint64(uint32(s.x[23]) ^ 0x1795e6c0)
	s.x[24] = uint64(uint32(s.x[25]) & uint32(s.x[24]))
	s.x[23] = uint64(uint32(s.x[23]) & 0x05e579a0)
	s.x[23] = uint64(uint32(s.x[24]) ^ uint32(s.x[23]))
	s.x[23] = uint64(uint32(s.x[23]) & uint32(s.x[22]<<4))
	s.x[22] = uint64(uint32(s.x[23]) ^ uint32(s.x[22]))
	s.x[20] = uint64(uint32(s.x[20]) ^ uint32(s.x[22]<<1))
	s.x[20] &= 0xff
	s.x[20] ^= s.x[6]
	value, err = s.load(work+s.x[20], 1)
	if err != nil {
		return 0, fmt.Errorf("shared byte-permute second input: %w", err)
	}
	s.x[20] = uint64(byte(value))
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[20]))
	s.x[20] = uint64(uint32(s.x[17]) ^ 0xffffffe0)
	s.x[20] = uint64(uint32(s.x[20]) + 0x32)
	s.x[22] = uint64(uint32(s.x[17]) ^ uint32(s.x[7]))
	s.x[22] = uint64(uint32(s.x[22]) - 0x7c)
	s.x[20] = uint64(uint32(s.x[20]) ^ uint32(s.x[22]))
	s.x[17] = uint64(uint32(s.x[17]) ^ 0xffffffe9)
	s.x[17] = uint64(uint32(s.x[17]) + 0x39)
	s.x[17] = uint64(uint32(s.x[20]) ^ uint32(s.x[17]))
	s.x[17] = uint64(uint32(s.x[17]) - 0x1d)
	s.x[20] = uint64(uint32(s.x[17]) & 0xb)
	s.x[3] = 0x060be5d0
	s.x[20] = uint64(uint32(s.x[20]) ^ 0x5c)
	s.x[22] = uint64(uint32(s.x[17]) ^ 0x12)
	s.x[23] = uint64(uint32(s.x[22]) & uint32(s.x[17]<<1))
	s.x[24] = 0x26
	s.x[23] = uint64(uint32(s.x[23]) & uint32(s.x[24]))
	s.x[23] = uint64(uint32(s.x[23]) ^ uint32(s.x[20]))
	for range 4 {
		s.x[23] = uint64(uint32(s.x[22]) & uint32(s.x[23]<<1))
		s.x[23] = uint64(uint32(s.x[23]) ^ uint32(s.x[20]))
	}
	s.x[22] = uint64(uint32(s.x[22]) & uint32(s.x[23]<<1))
	s.x[20] = uint64(uint32(s.x[22]) ^ uint32(s.x[20]))
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[20]<<1))
	s.x[17] = uint64(uint32(s.x[17]) ^ 0x10)
	if err := s.storeValue(work+s.x[0], uint64(byte(s.x[17])), 1); err != nil {
		return 0, fmt.Errorf("shared byte-permute final output: %w", err)
	}
	s.x[20] = uint64(uint32(s.x[17]) & 0xff)
	next, err := loadBootstrapPointer(s, s.x[21], uint32(s.x[1]))
	if err != nil {
		return 0, fmt.Errorf("shared byte-permute dispatch: %w", err)
	}
	s.x[1] = next
	s.x[17] = 0xbf5f
	return next, nil
}
