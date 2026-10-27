package recovery

import "fmt"

const accountMixerSetup = uint64(0x100062280)

func runAccountMixerSetup(s *state) (uint64, error) {
	s.traceBlock(accountMixerSetup)
	value, err := s.load(s.sp+0x7c, 4)
	if err != nil {
		return 0, fmt.Errorf("account mixer-setup word: %w", err)
	}
	w4 := uint32(value)
	w15 := w4 ^ uint32(0x9a9cb711)
	w17 := w4 ^ uint32(0x15662adf)
	w17 += 0xea99d521
	w2 := w4 ^ uint32(0x84e55073)
	w4 = 0x7b1aaf8d
	w2 += w4
	w17 ^= w2
	w2 = uint32(s.x[7]) ^ uint32(0xe7e0161e)
	w15 = w2 + w15
	w15 = w17 ^ w15

	w9 := uint32(s.x[9])
	w17 = w9 ^ uint32(0xe7cf2f5d)
	w17 += 0x1830d0a3
	w2 = w9 ^ uint32(0xfaab2a36)
	w4 = 0x0554d5ca
	w2 += w4
	w17 ^= w2
	w2 = 0x167bc8d6
	w9 ^= w2
	w2 = 0xe984372a
	w9 += w2
	w9 = w17 ^ w9
	w9 += w15
	w15 = 0x453d7a44
	w9 += w15
	w15 = w9 ^ uint32(0x3688b49f)
	w17 = 0xe78a95ab
	w15 += w17
	w17 = w9 ^ uint32(0x8ae06a35)
	w2 = 0x5be24b01
	w17 += w2
	w15 ^= w17
	w17 = w9 ^ uint32(0xe962ff1f)
	w2 = 0x3860de2b
	w17 += w2
	w17 = w9 ^ w17
	w2 = 0x7bf7ff7f
	w9 ^= w2
	w15 = w17 ^ w15
	w17 = 0xaaf5de4b
	w9 += w17
	w17 = 0x25e21377
	w9 ^= w17
	w17 = w15 ^ w9

	s.compare32(uint32(s.x[7]), 0x62e0114c)
	w9 = uint32(boolValue(s.condition("eq"))) * 0x8ca4
	w15 = uint32(s.x[7]) ^ uint32(0x82832578)
	w9 = w15 ^ w9
	w15 = 0x8283d7dd
	w24 := uint32(s.x[7]) ^ w15

	s.x[2] = uint64(w2)
	s.x[4] = uint64(w4)
	s.x[15] = uint64(w15)
	s.x[24] = uint64(w24)
	if err := s.storeValue(s.sp+0x7c, uint64(w17), 4); err != nil {
		return 0, fmt.Errorf("account mixer-setup saved word: %w", err)
	}
	if err := s.storeValue(addOffset(s.x[29], -0x84), uint64(w17), 4); err != nil {
		return 0, fmt.Errorf("account mixer-setup frame word: %w", err)
	}
	next, err := loadBootstrapPointer(s, s.x[21], w9)
	if err != nil {
		return 0, fmt.Errorf("account mixer-setup target: %w", err)
	}
	s.x[9] = next
	s.x[17] = 0x169e
	return next, nil
}
