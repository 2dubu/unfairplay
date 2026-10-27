package recovery

import "fmt"

// runAccountSHA1CompressionStage replaces the native flattened byte-pack,
// schedule expansion, 80 round handlers, and five feed-forward circuits with
// one typed SHA-1 compression.
func runAccountSHA1CompressionStage(s *state) (uint64, error) {
	if uint32(s.x[4]) != accountBytePackCounter {
		return 0, fmt.Errorf("account SHA-1 block counter %#x is unsupported", uint32(s.x[4]))
	}
	if s.x[8] > 0xff {
		return 0, fmt.Errorf("account SHA-1 input mask %#x is unsupported", s.x[8])
	}

	var block [64]byte
	for index := range block {
		value, err := s.load(s.x[1]+uint64(index), 1)
		if err != nil {
			return 0, fmt.Errorf("account SHA-1 block byte %d: %w", index, err)
		}
		block[index] = byte(value)
	}
	encodedInitialE, err := s.load(s.sp+0x7c, 4)
	if err != nil {
		return 0, fmt.Errorf("account SHA-1 initial E: %w", err)
	}
	initial := accountSHA1State{
		A: uint32(s.x[26]) ^ accountSHA1StateMask,
		B: uint32(s.x[0]) ^ accountSHA1StateMask,
		C: uint32(s.x[22]) ^ accountSHA1StateMask,
		D: uint32(s.x[3]) ^ accountSHA1StateMask,
		E: uint32(encodedInitialE) ^ accountSHA1StateMask,
	}
	compression := accountSHA1CompressBlock(initial, block, byte(s.x[8]))
	for index, word := range compression.Schedule {
		if err := s.storeValue(s.x[10]+uint64(index*4), uint64(word^accountSHA1ScheduleMask), 4); err != nil {
			return 0, fmt.Errorf("account SHA-1 schedule word %d: %w", index, err)
		}
	}

	encodedA := compression.State.A ^ accountSHA1StateMask
	encodedB := compression.State.B ^ accountSHA1StateMask
	encodedC := compression.State.C ^ accountSHA1StateMask
	encodedD := compression.State.D ^ accountSHA1StateMask
	encodedE := compression.State.E ^ accountSHA1StateMask
	for index, word := range [...]uint32{encodedA, encodedB, encodedC, encodedD, encodedE} {
		if err := s.storeValue(addOffset(s.x[29], int64(-0x94+index*4)), uint64(word), 4); err != nil {
			return 0, fmt.Errorf("account SHA-1 frame word %d: %w", index, err)
		}
	}
	if err := s.storeValue(s.sp+0x7c, uint64(encodedE), 4); err != nil {
		return 0, fmt.Errorf("account SHA-1 updated E: %w", err)
	}

	// Native live-outs required by the next dispatcher block. Values omitted
	// here are dead circuit temporaries and retain their entry values.
	s.x[0], s.x[3], s.x[22], s.x[26] = uint64(encodedB), uint64(encodedD), uint64(encodedC), uint64(encodedA)
	s.x[2] = 0x7bf7ff7f
	s.x[4] = 0x0554d5ca
	s.x[7] = 0x7bb20000
	s.x[9] = 0x897e
	s.x[12] = 0x3d7f0910
	s.x[13] = uint64(encodedC ^ 0xef42ba1f)
	s.x[15] = 0xdd4e54af
	s.x[16] = uint64(compression.Compressed.B ^ accountSHA1StateMask)
	s.x[17] = 0xaa2b
	s.x[19] = 0x7bb2a446
	s.x[20] = uint64(accountBytePackCounter + accountSHA1Rounds)
	lastEncodedWord := compression.Schedule[len(compression.Schedule)-1] ^ accountSHA1ScheduleMask
	s.x[23] = uint64(((lastEncodedWord ^ 0x4378c660) + 0xbc8739a0) ^ ((lastEncodedWord ^ 0xe0ddc8de) + 0x1f223722))
	s.x[24] = 0x892c
	s.compare32(0x82835ef1, 0x62e0114c)
	next, err := loadBootstrapPointer(s, s.x[21], uint32(s.x[24]))
	if err != nil {
		return 0, fmt.Errorf("account SHA-1 stage target: %w", err)
	}
	s.x[11] = next
	return next, nil
}
