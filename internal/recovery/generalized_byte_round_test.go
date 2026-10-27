package recovery

import (
	"crypto/sha256"
	"math/rand"
	"testing"
)

func referenceGeneralizedByteRound(round generalizedByteRound, input [16]byte) (output [16]byte) {
	copy(output[:], round.baseline[:])
	for position := range input {
		value := input[position]
		for slot := range 4 {
			target := round.targets[position][slot]
			output[target] = output[target] ^ round.tables[position][value][slot]
		}
	}
	return output
}

func TestGeneralizedByteRound(t *testing.T) {
	random := rand.New(rand.NewSource(0x53494442))
	var round generalizedByteRound
	if _, err := random.Read(round.baseline[:]); err != nil {
		t.Fatal(err)
	}
	for position := range 16 {
		column := position % 4
		round.targets[position] = [4]byte{byte(column * 4), byte(column*4 + 1), byte(column*4 + 2), byte(column*4 + 3)}
		for value := range 256 {
			if _, err := random.Read(round.tables[position][value][:]); err != nil {
				t.Fatal(err)
			}
		}
	}
	var corpus []byte
	for index := range 4096 {
		digest := sha256.Sum256([]byte{byte(index), byte(index >> 8), 0x47, 0x52})
		var input [16]byte
		copy(input[:], digest[:16])
		got := round.apply(input)
		want := referenceGeneralizedByteRound(round, input)
		if got != want {
			t.Fatalf("input %d output = %x, want %x", index, got, want)
		}
		corpus = append(corpus, got[:]...)
	}
	if got := sha256.Sum256(corpus); got == [32]byte{} {
		t.Fatal("generalized round corpus has an empty digest")
	}
}
