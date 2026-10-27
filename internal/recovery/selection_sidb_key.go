package recovery

import (
	"fmt"
	"math/bits"
)

// selectionSIDBAffineByteDecoder is one external byte encoding around an
// AES-128 round-key byte. Bit zero of the 9-bit input is the affine constant;
// the remaining bits are the encoded byte.
type selectionSIDBAffineByteDecoder [8]uint16

func (decoder selectionSIDBAffineByteDecoder) apply(value byte) (result byte) {
	input := uint16(1) | uint16(value)<<1
	for bit, row := range decoder {
		result |= byte(bits.OnesCount16(input&row)&1) << bit
	}
	return result
}

type selectionSIDBRoundKeyDecoder [16]selectionSIDBAffineByteDecoder

// decoderKey converts the typed encoded round-10 state into the AES-128 key
// used by the SIDB consumer. The native path performs the same ordinary key
// schedule reversal under protected byte encodings.
func (decoder selectionSIDBRoundKeyDecoder) decoderKey(encoded [16]byte) ([16]byte, error) {
	var round10 [16]byte
	for position, byteDecoder := range decoder {
		var seen [256]bool
		for value := range 256 {
			decoded := byteDecoder.apply(byte(value))
			if seen[decoded] {
				return [16]byte{}, fmt.Errorf("SIDB round-key byte decoder %d is not bijective", position)
			}
			seen[decoded] = true
		}
		round10[position] = byteDecoder.apply(encoded[position])
	}
	return reverseSelectionSIDBAES128RoundKey(round10, 10)
}

func reverseSelectionSIDBAES128RoundKey(roundKey [16]byte, round int) (result [16]byte, err error) {
	if round < 0 || round > 10 {
		return result, fmt.Errorf("SIDB AES round %d is outside 0..10", round)
	}
	if round == 0 {
		return roundKey, nil
	}
	var words [44][4]byte
	for index := range 4 {
		copy(words[round*4+index][:], roundKey[index*4:index*4+4])
	}

	rcon := byte(1)
	for index := 1; index < round; index++ {
		rcon = frameGFMul(rcon, 2)
	}
	for word := round*4 + 3; word >= 4; word-- {
		if word%4 != 0 {
			for index := range 4 {
				words[word-4][index] = words[word][index] ^ words[word-1][index]
			}
			continue
		}
		previous := words[word-1]
		transformed := [4]byte{
			frameSBox[previous[1]] ^ rcon,
			frameSBox[previous[2]],
			frameSBox[previous[3]],
			frameSBox[previous[0]],
		}
		for index := range 4 {
			words[word-4][index] = words[word][index] ^ transformed[index]
		}
		if rcon&1 != 0 {
			rcon = rcon>>1 ^ 0x8d
		} else {
			rcon >>= 1
		}
	}
	for index := range 4 {
		copy(result[index*4:index*4+4], words[index][:])
	}
	return result, nil
}
