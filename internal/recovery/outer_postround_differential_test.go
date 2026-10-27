package recovery

import (
	"math/rand"
	"testing"
)

func TestOuterPostRoundStepMatchesCircuit(t *testing.T) {
	random := rand.New(rand.NewSource(0x26a428))
	for sample := 0; sample < 100000; sample++ {
		var input [6]uint32
		for index := range input {
			input[index] = random.Uint32()
		}
		got0, got1 := outerPostRoundStep(input[0], input[1], input[2], input[3], input[4], input[5])
		want0, want1 := outerPostRoundStepCircuit(input[0], input[1], input[2], input[3], input[4], input[5])
		if got0 != want0 || got1 != want1 {
			t.Fatalf("sample %d = %08x/%08x, want %08x/%08x", sample, got0, got1, want0, want1)
		}
	}
}
