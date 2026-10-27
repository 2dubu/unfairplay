package recovery

import "fmt"

const sharedEnrollmentBytePermute = uint64(0x1003e9534)

const sharedBytePermuteFirst = "\x95\x88\x83\x86\x61\x64\x9f\x92\xad\xa0\xdb\xde\xb9\xbc\xb7\xaa\xc5\xf8\xf3\xf6\xd1\xd4\xcf\xc2\x1d\x10\x0b\x0e\xe9\xec\xe7\x1a\xf5\xe8\xe3\xe6\xc1\xc4\xff\xf2\x0d\x00\x3b\x3e\x19\x1c\x17\x0a\x25\x58\x53\x56\x31\x34\x2f\x22\x7d\x70\x6b\x6e\x49\x4c\x47\x7a\x55\x48\x43\x46\x21\x24\x5f\x52\x6d\x60\x9b\x9e\x79\x7c\x77\x6a\x85\xb8\xb3\xb6\x91\x94\x8f\x82\xdd\xd0\xcb\xce\xa9\xac\xa7\xda\xb5\xa8\xa3\xa6\x81\x84\xbf\xb2\xcd\xc0\xfb\xfe\xd9\xdc\xd7\xca\xe5\x18\x13\x16\xf1\xf4\xef\xe2\x3d\x30\x2b\x2e\x09\x0c\x07\x3a\x15\x08\x03\x06\xe1\xe4\x1f\x12\x2d\x20\x5b\x5e\x39\x3c\x37\x2a\x45\x78\x73\x76\x51\x54\x4f\x42\x9d\x90\x8b\x8e\x69\x6c\x67\x9a\x75\x68\x63\x66\x41\x44\x7f\x72\x8d\x80\xbb\xbe\x99\x9c\x97\x8a\xa5\xd8\xd3\xd6\xb1\xb4\xaf\xa2\xfd\xf0\xeb\xee\xc9\xcc\xc7\xfa\xd5\xc8\xc3\xc6\xa1\xa4\xdf\xd2\xed\xe0\x1b\x1e\xf9\xfc\xf7\xea\x05\x38\x33\x36\x11\x14\x0f\x02\x5d\x50\x4b\x4e\x29\x2c\x27\x5a\x35\x28\x23\x26\x01\x04\x3f\x32\x4d\x40\x7b\x7e\x59\x5c\x57\x4a\x65\x98\x93\x96\x71\x74\x6f\x62\xbd\xb0\xab\xae\x89\x8c\x87\xba"
const sharedBytePermuteSecondIndex = "\xa8\xa9\xaa\xab\xac\xad\xae\xaf\xb0\xb1\x92\x93\x94\x95\xf6\xf7\xf8\xf9\xfa\xfb\xfc\xfd\xbe\xbf\xc0\xc1\xe2\xe3\xe4\xe5\xc6\xc7\xc8\xc9\xca\xcb\xcc\xcd\xce\xcf\xd0\xd1\xf2\xf3\xf4\xf5\x96\x97\x98\x99\x9a\x9b\x9c\x9d\xde\xdf\xe0\xe1\xc2\xc3\xc4\xc5\xe6\xe7\xe8\xe9\xea\xeb\xec\xed\xee\xef\xf0\xf1\x52\x53\x54\x55\xb6\xb7\xb8\xb9\xba\xbb\xbc\xbd\xfe\xff\x00\x01\xa2\xa3\xa4\xa5\x06\x07\x08\x09\x0a\x0b\x0c\x0d\x0e\x0f\x10\x11\x32\x33\x34\x35\x56\x57\x58\x59\x5a\x5b\x5c\x5d\x1e\x1f\x20\x21\x02\x03\x04\x05\x26\x27\x28\x29\x2a\x2b\x2c\x2d\x2e\x2f\x30\x31\x12\x13\x14\x15\x76\x77\x78\x79\x7a\x7b\x7c\x7d\x3e\x3f\x40\x41\x62\x63\x64\x65\x46\x47\x48\x49\x4a\x4b\x4c\x4d\x4e\x4f\x50\x51\x72\x73\x74\x75\x16\x17\x18\x19\x1a\x1b\x1c\x1d\x5e\x5f\x60\x61\x42\x43\x44\x45\x66\x67\x68\x69\x6a\x6b\x6c\x6d\x6e\x6f\x70\x71\xd2\xd3\xd4\xd5\x36\x37\x38\x39\x3a\x3b\x3c\x3d\x7e\x7f\x80\x81\x22\x23\x24\x25\x86\x87\x88\x89\x8a\x8b\x8c\x8d\x8e\x8f\x90\x91\xb2\xb3\xb4\xb5\xd6\xd7\xd8\xd9\xda\xdb\xdc\xdd\x9e\x9f\xa0\xa1\x82\x83\x84\x85\xa6\xa7"
const sharedBytePermuteFinal = "\x47\x40\x41\x42\x7b\x44\x45\x46\x4f\x48\x49\x4a\x43\x4c\x4d\x4e\x57\x50\x51\x52\x4b\x54\x55\x56\x5f\x58\x59\x5a\x53\x5c\x5d\x5e\x67\x60\x61\x62\x9b\x64\x65\x66\x6f\x68\x69\x6a\x63\x6c\x6d\x6e\x77\x70\x71\x72\x6b\x74\x75\x76\x7f\x78\x79\x7a\x73\x7c\x7d\x7e\x07\x00\x01\x02\x3b\x04\x05\x06\x0f\x08\x09\x0a\x03\x0c\x0d\x0e\x17\x10\x11\x12\x0b\x14\x15\x16\x1f\x18\x19\x1a\x13\x1c\x1d\x1e\x27\x20\x21\x22\x5b\x24\x25\x26\x2f\x28\x29\x2a\x23\x2c\x2d\x2e\x37\x30\x31\x32\x2b\x34\x35\x36\x3f\x38\x39\x3a\x33\x3c\x3d\x3e\xc7\xc0\xc1\xc2\xfb\xc4\xc5\xc6\xcf\xc8\xc9\xca\xc3\xcc\xcd\xce\xd7\xd0\xd1\xd2\xcb\xd4\xd5\xd6\xdf\xd8\xd9\xda\xd3\xdc\xdd\xde\xe7\xe0\xe1\xe2\x1b\xe4\xe5\xe6\xef\xe8\xe9\xea\xe3\xec\xed\xee\xf7\xf0\xf1\xf2\xeb\xf4\xf5\xf6\xff\xf8\xf9\xfa\xf3\xfc\xfd\xfe\x87\x80\x81\x82\xbb\x84\x85\x86\x8f\x88\x89\x8a\x83\x8c\x8d\x8e\x97\x90\x91\x92\x8b\x94\x95\x96\x9f\x98\x99\x9a\x93\x9c\x9d\x9e\xa7\xa0\xa1\xa2\xdb\xa4\xa5\xa6\xaf\xa8\xa9\xaa\xa3\xac\xad\xae\xb7\xb0\xb1\xb2\xab\xb4\xb5\xb6\xbf\xb8\xb9\xba\xb3\xbc\xbd\xbe"

// runSharedEnrollmentBytePermuteSemantic replaces two native MBA byte
// circuits with their exhaustive 8-bit permutations.
func runSharedEnrollmentBytePermuteSemantic(s *state) (uint64, error) {
	if s == nil {
		return 0, fmt.Errorf("shared byte-permute state is unavailable")
	}
	s.traceBlock(sharedEnrollmentBytePermute)
	if uint32(s.x[2]) != 0x19d6 || uint32(s.x[5]) != 0xba1c1cae ||
		uint32(s.x[6]) != 0xae || uint32(s.x[7]) != 0x52 || uint32(s.x[8]) != 0x1b {
		return 0, fmt.Errorf("unsupported shared byte-permute ABI")
	}
	source, err := s.load(s.sp+0x50, 8)
	if err != nil {
		return 0, fmt.Errorf("shared byte-permute source: %w", err)
	}
	input, err := s.load(source+uint64(byte(s.x[20]^s.x[8])), 1)
	if err != nil {
		return 0, fmt.Errorf("shared byte-permute first input: %w", err)
	}
	first := sharedBytePermuteFirst[byte(input)]
	work, err := s.load(s.sp+0x58, 8)
	if err != nil {
		return 0, fmt.Errorf("shared byte-permute work buffer: %w", err)
	}
	if err := s.storeValue(work+s.x[0], uint64(first), 1); err != nil {
		return 0, fmt.Errorf("shared byte-permute first output: %w", err)
	}
	secondIndex := byte(uint32(s.x[11]) + 0xa8)
	second, err := s.load(work+uint64(secondIndex), 1)
	if err != nil {
		return 0, fmt.Errorf("shared byte-permute second input: %w", err)
	}
	output := sharedBytePermuteFinal[first^0x5b^byte(second)]
	if err := s.storeValue(work+s.x[0], uint64(output), 1); err != nil {
		return 0, fmt.Errorf("shared byte-permute final output: %w", err)
	}
	next, err := loadBootstrapPointer(s, s.x[21], uint32(s.x[1]))
	if err != nil {
		return 0, fmt.Errorf("shared byte-permute dispatch: %w", err)
	}
	s.x[1] = next
	s.x[3] = 0x060be5d0
	s.x[17] = 0xbf5f
	s.x[20] = uint64(output)
	s.x[27] = work
	return next, nil
}

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
