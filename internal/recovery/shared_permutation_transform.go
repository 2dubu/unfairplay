package recovery

import (
	"fmt"
	"math/bits"
)

const (
	sharedPermutationFront = uint64(0x100439e54)
	sharedPermutationTail  = uint64(0x10043a000)
)

// runSharedPermutationSemantic is the natural four-byte table round used by
// both record selection and SINF enrollment. The native implementation hides
// one byte rotation, three indexed table reads, and a uint32 polynomial behind
// two large Boolean circuits.
func runSharedPermutationSemantic(s *state) (uint64, error) {
	if s == nil {
		return 0, fmt.Errorf("shared permutation state is unavailable")
	}
	s.traceBlock(sharedPermutationFront)
	if uint32(s.x[2]) != 0x19d6 || uint32(s.x[5]) != 0xba1c1cae ||
		uint32(s.x[6]) != 0xae || uint32(s.x[7]) != 0x52 || uint32(s.x[8]) != 0x1b ||
		uint32(s.x[26]) != 0x40cb0000 || uint32(s.x[30]) != 0xb667b683 ||
		uint32(s.x[11])-uint32(s.x[12]) != 0x22f1f10c ||
		uint32(s.x[12])-uint32(s.x[16]) != 0xae {
		return 0, fmt.Errorf("unsupported shared permutation ABI")
	}
	work, err := s.load(s.sp+0x58, 8)
	if err != nil {
		return 0, fmt.Errorf("shared permutation table: %w", err)
	}
	s.x[28] = work
	selectorValue, err := s.load(work+uint64(byte(uint32(s.x[11])-0x62)), 1)
	if err != nil {
		return 0, fmt.Errorf("shared permutation selector: %w", err)
	}
	rotation := int((byte(selectorValue) + 5) & 7)
	rotated := bits.RotateLeft8(byte(s.x[20])^0x1b, rotation)
	if err := s.storeValue(work+s.x[0], uint64(rotated^0x1b), 1); err != nil {
		return 0, fmt.Errorf("shared permutation rotated output: %w", err)
	}

	second, err := s.load(work+uint64(byte(uint32(s.x[11])+0xae)), 1)
	if err != nil {
		return 0, fmt.Errorf("shared permutation second byte: %w", err)
	}
	thirdIndex := byte(((uint32(s.x[19]) & 0x7f) << 1) ^ 0x52)
	third, err := s.load(work+uint64(thirdIndex), 1)
	if err != nil {
		return 0, fmt.Errorf("shared permutation third byte: %w", err)
	}
	fourthIndex := byte(uint32(s.x[12])*uint32(s.x[16]) + 0xa2)
	fourth, err := s.load(work+uint64(fourthIndex), 1)
	if err != nil {
		return 0, fmt.Errorf("shared permutation fourth byte: %w", err)
	}
	word := uint32(rotated)<<24 |
		uint32(byte(second)^0x1b)<<16 |
		uint32(byte(third)^0x1b)<<8 |
		uint32(byte(fourth)^0x1b)
	output := sharedPermutationPolynomial(word)
	s.traceBlock(sharedPermutationTail)
	next, err := loadBootstrapPointer(s, s.x[21], uint32(s.x[2]))
	if err != nil {
		return 0, fmt.Errorf("shared permutation target: %w", err)
	}
	s.x[1] = 0x65c0
	s.x[2] = next
	s.x[17] = 0xcd20
	s.x[19] = 0x71
	s.x[20] = uint64(output)
	s.x[27] = bootstrapTargetTable
	return next, nil
}

func sharedPermutationPolynomial(word uint32) uint32 {
	mixed := (word << 3) ^ (word >> 5)
	word += mixed
	square := word * word
	highByte := word >> 8
	highWord := word >> 16
	result := square*word + highWord
	highSquare := highByte * highByte
	result += (highSquare * highByte) * highSquare
	wordSquare := highWord * highWord
	wordCube := wordSquare * highWord
	result += (wordCube * highWord) * wordCube
	result -= 0xd4 & (result << 1)
	return (result + 0x6a) ^ 0x71
}
