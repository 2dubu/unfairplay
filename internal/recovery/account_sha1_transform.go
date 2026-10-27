package recovery

import (
	"fmt"
	"math/bits"
)

const (
	accountHotMixFront      = uint64(0x100026bfc)
	accountHotMixTail       = uint64(0x100026c00)
	accountHotMixSecond     = uint64(0x10002dcd8)
	accountHotMixThirdFront = uint64(0x100065320)
	accountHotMixThirdTail  = uint64(0x100065400)
	accountHotMixFourth     = uint64(0x100088c14)
	accountSHA1StateMask    = uint32(0x0b1fcdbd)
	accountSHA1ScheduleMask = uint32(0xce757bc5)
)

type accountSHA1Phase uint8

const (
	accountSHA1Choose accountSHA1Phase = iota
	accountSHA1ParityEarly
	accountSHA1Majority
	accountSHA1ParityLate
)

type accountSHA1State struct {
	A, B, C, D, E uint32
}

// accountSHA1Round is one standard SHA-1 compression round. The native
// handlers split its 80 rounds into the usual four groups of twenty.
func accountSHA1Round(state accountSHA1State, word uint32, phase accountSHA1Phase) accountSHA1State {
	var mix, constant uint32
	switch phase {
	case accountSHA1Choose:
		mix = state.B&state.C | ^state.B&state.D
		constant = 0x5a827999
	case accountSHA1ParityEarly:
		mix = state.B ^ state.C ^ state.D
		constant = 0x6ed9eba1
	case accountSHA1Majority:
		mix = state.B&state.C | state.B&state.D | state.C&state.D
		constant = 0x8f1bbcdc
	case accountSHA1ParityLate:
		mix = state.B ^ state.C ^ state.D
		constant = 0xca62c1d6
	default:
		panic("invalid account SHA-1 phase")
	}
	return accountSHA1State{
		A: bits.RotateLeft32(state.A, 5) + mix + state.E + word + constant,
		B: state.A,
		C: bits.RotateLeft32(state.B, 30),
		D: state.C,
		E: state.D,
	}
}

func loadAccountSHA1Word(s *state, counter uint64) (uint32, error) {
	word, err := s.load(s.x[10]+s.x[6]+(counter<<2), 4)
	if err != nil {
		return 0, err
	}
	return uint32(word) ^ accountSHA1ScheduleMask, nil
}

func accountSHA1Decode(s *state, registers [5]int) accountSHA1State {
	return accountSHA1State{
		A: uint32(s.x[registers[0]]) ^ accountSHA1StateMask,
		B: uint32(s.x[registers[1]]) ^ accountSHA1StateMask,
		C: uint32(s.x[registers[2]]) ^ accountSHA1StateMask,
		D: uint32(s.x[registers[3]]) ^ accountSHA1StateMask,
		E: uint32(s.x[registers[4]]) ^ accountSHA1StateMask,
	}
}

// The four wrappers preserve the surrounding dispatcher ABI, but the state
// transition itself is the typed SHA-1 round above. State and schedule words
// use fixed XOR encodings at this boundary.
func runAccountHotMix(s *state) (uint64, error) {
	s.traceBlock(accountHotMixFront)
	s.traceBlock(accountHotMixTail)
	word, err := loadAccountSHA1Word(s, s.x[16])
	if err != nil {
		return 0, fmt.Errorf("account SHA-1 phase 0 word: %w", err)
	}
	oldA, oldC, oldD := s.x[12], s.x[13], s.x[11]
	value := accountSHA1Round(accountSHA1Decode(s, [5]int{12, 7, 13, 11, 19}), word, accountSHA1Choose)
	s.x[12] = uint64(value.A ^ accountSHA1StateMask)
	s.x[7], s.x[15] = oldA, oldA
	s.x[13] = uint64(value.C ^ accountSHA1StateMask)
	s.x[11] = oldC
	s.x[19], s.x[4] = oldD, oldD
	s.x[16]++
	gate := uint32(boolValue(s.x[16] != 0x2ba7c39f))
	index := gate*14 ^ uint32(s.x[9])
	next, err := loadBootstrapPointer(s, s.x[21], index)
	if err != nil {
		return 0, fmt.Errorf("account SHA-1 phase 0 target: %w", err)
	}
	s.x[23] = next
	s.x[2] = 0x6455814d
	s.x[20] = 0x2ba7c39f
	s.x[17] = 0x1d92
	return next, nil
}

func runAccountHotMixSecond(s *state) (uint64, error) {
	s.traceBlock(accountHotMixSecond)
	selector := uint32(s.x[9])
	word, err := loadAccountSHA1Word(s, s.x[20])
	if err != nil {
		return 0, fmt.Errorf("account SHA-1 phase 1 word: %w", err)
	}
	oldA, oldC, oldD := s.x[12], s.x[13], s.x[11]
	value := accountSHA1Round(accountSHA1Decode(s, [5]int{12, 15, 13, 11, 4}), word, accountSHA1ParityEarly)
	s.x[12] = uint64(value.A ^ accountSHA1StateMask)
	s.x[15], s.x[19] = oldA, oldA
	s.x[13] = uint64(value.C ^ accountSHA1StateMask)
	s.x[11] = oldC
	s.x[4], s.x[7] = oldD, oldD
	s.x[20]++
	gate := uint32(boolValue(s.x[20] == 0x2ba7c3b3))
	index := gate*0xc7a1 ^ uint32(s.x[9])
	next, err := loadBootstrapPointer(s, s.x[21], index)
	if err != nil {
		return 0, fmt.Errorf("account SHA-1 phase 1 target: %w", err)
	}
	s.x[24] = next
	s.x[2] = 0x2a8dd0f2
	s.x[16] = uint64((selector + 0x0989137d) & 0xf676ffff)
	s.x[23] = 0x2ba7c3b3
	s.x[17] = 0x1d92
	return next, nil
}

func runAccountHotMixThird(s *state) (uint64, error) {
	s.traceBlock(accountHotMixThirdFront)
	s.traceBlock(accountHotMixThirdTail)
	selector := uint32(s.x[16])
	word, err := loadAccountSHA1Word(s, s.x[23])
	if err != nil {
		return 0, fmt.Errorf("account SHA-1 phase 2 word: %w", err)
	}
	oldA, oldC, oldD := s.x[12], s.x[13], s.x[11]
	value := accountSHA1Round(accountSHA1Decode(s, [5]int{12, 19, 13, 11, 7}), word, accountSHA1Majority)
	s.x[12] = uint64(value.A ^ accountSHA1StateMask)
	s.x[19], s.x[15] = oldA, oldA
	s.x[13] = uint64(value.C ^ accountSHA1StateMask)
	s.x[11] = oldC
	s.x[7], s.x[9] = oldD, oldD
	s.x[23]++
	gate := uint32(boolValue(s.x[23] != 0x2ba7c3c7))
	index := gate*0xffe ^ uint32(s.x[16])
	next, err := loadBootstrapPointer(s, s.x[21], index)
	if err != nil {
		return 0, fmt.Errorf("account SHA-1 phase 2 target: %w", err)
	}
	s.x[24] = next
	s.x[2] = 0x26c55028
	s.x[4] = uint64(selector + 0xffffd987)
	s.x[20] = 0x2ba7c3c7
	s.x[17] = 0xfab0
	return next, nil
}

func runAccountHotMixFourth(s *state) (uint64, error) {
	s.traceBlock(accountHotMixFourth)
	encodedWord, err := s.load(s.x[10]+s.x[6]+(s.x[20]<<2), 4)
	if err != nil {
		return 0, fmt.Errorf("account SHA-1 phase 3 word: %w", err)
	}
	word := uint32(encodedWord) ^ accountSHA1ScheduleMask
	oldA, oldC, oldD := s.x[12], s.x[13], s.x[11]
	value := accountSHA1Round(accountSHA1Decode(s, [5]int{12, 15, 13, 11, 9}), word, accountSHA1ParityLate)
	s.x[12] = uint64(value.A ^ accountSHA1StateMask)
	s.x[15], s.x[16] = oldA, oldA
	s.x[13] = uint64(value.C ^ accountSHA1StateMask)
	s.x[11] = oldC
	s.x[9] = oldD
	s.x[20]++
	gate := uint32(boolValue(s.x[20] == 0x2ba7c3db))
	index := gate*0x13 ^ uint32(s.x[4])
	next, err := loadBootstrapPointer(s, s.x[21], index)
	if err != nil {
		return 0, fmt.Errorf("account SHA-1 phase 3 target: %w", err)
	}
	s.x[19] = next
	s.x[2] = 0x3fffceff
	s.x[7] = uint64(uint32(s.x[4]) + 0x82831d88)
	s.x[23] = uint64(((uint32(encodedWord) ^ 0x4378c660) + 0xbc8739a0) ^ ((uint32(encodedWord) ^ 0xe0ddc8de) + 0x1f223722))
	s.x[24] = uint64((uint32(encodedWord) ^ 0xe0ddc8de) + 0x1f223722)
	s.x[17] = 0xfb38
	return next, nil
}
