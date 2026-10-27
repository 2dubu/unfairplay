package recovery

import (
	"math/big"
	"math/rand"
	"testing"
)

func TestMontgomeryCIOSSquare(t *testing.T) {
	random := rand.New(rand.NewSource(0x26a428))
	base := new(big.Int).Lsh(big.NewInt(1), 32)
	for _, limbs := range []int{1, 2, 32, 64} {
		for trial := 0; trial < 64; trial++ {
			modulus := make([]uint32, limbs)
			input := make([]uint32, limbs)
			for index := range modulus {
				modulus[index] = random.Uint32()
				input[index] = random.Uint32()
			}
			modulus[0] |= 1
			modulus[len(modulus)-1] |= 0x80000000
			modulusValue := littleEndianWordsBig(modulus)
			inputValue := littleEndianWordsBig(input)
			inputValue.Mod(inputValue, modulusValue)
			for index := range input {
				input[index] = uint32(new(big.Int).Rsh(new(big.Int).Set(inputValue), uint(32*index)).Uint64())
			}

			low := new(big.Int).SetUint64(uint64(modulus[0]))
			inverseValue := new(big.Int).ModInverse(low, base)
			if inverseValue == nil {
				t.Fatal("odd word has no radix inverse")
			}
			inverse := -uint32(inverseValue.Uint64())
			got, err := montgomeryCIOSSquare(input, modulus, inverse)
			if err != nil {
				t.Fatalf("limbs %d trial %d: %v", limbs, trial, err)
			}
			rawWords := append(append([]uint32(nil), got.Words...), got.Carry)
			raw := littleEndianWordsBig(rawWords)
			radix := new(big.Int).Lsh(big.NewInt(1), uint(32*limbs))
			radixInverse := new(big.Int).ModInverse(radix, modulusValue)
			want := new(big.Int).Mul(inputValue, inputValue)
			want.Mul(want, radixInverse)
			want.Mod(want, modulusValue)
			if new(big.Int).Mod(new(big.Int).Set(raw), modulusValue).Cmp(want) != 0 {
				t.Fatalf("limbs %d trial %d has wrong residue", limbs, trial)
			}
			if raw.Cmp(new(big.Int).Lsh(new(big.Int).Set(modulusValue), 1)) >= 0 {
				t.Fatalf("limbs %d trial %d result is not bounded by 2N", limbs, trial)
			}
		}
	}
}

func TestMontgomeryCIOSSquareRejectsInvalidParameters(t *testing.T) {
	if _, err := montgomeryCIOSSquare(nil, nil, 0); err == nil {
		t.Fatal("empty Montgomery input was accepted")
	}
	if _, err := montgomeryCIOSSquare([]uint32{1}, []uint32{2}, 0); err == nil {
		t.Fatal("even Montgomery modulus was accepted")
	}
}
