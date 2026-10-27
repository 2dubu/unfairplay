package recovery

import (
	"math/rand"
	"testing"
)

func TestOuterOptionalIndexAndDecode(t *testing.T) {
	if got := outerOptionalFirstIndex(0x0a782fec); got != 20 {
		t.Fatalf("first index = %d, want 20", got)
	}
	if got := outerOptionalDecode(0x39f4c55b); got != 0x976deca9 {
		t.Fatalf("decode = %#x, want %#x", got, uint32(0x976deca9))
	}
}

func TestOuterOptionalNormalizedIndexExhaustive(t *testing.T) {
	for phase := uint32(0); phase < outerWriterSeedStates; phase++ {
		seed := outerWriterSeedBase + phase
		want := (phase + outerOptionalPhaseOffset) % outerWriterSeedStates
		if got := outerOptionalFirstIndexMBA(seed); got != want {
			t.Fatalf("phase %d: MBA=%d, semantic=%d", phase, got, want)
		}
		if got := outerOptionalFirstIndex(seed); got != want {
			t.Fatalf("phase %d: index=%d, want=%d", phase, got, want)
		}
	}
}

func TestOuterOptionalFirstIndexMatchesCircuit(t *testing.T) {
	for seed := uint32(0); seed < 1<<16; seed++ {
		if got, want := outerOptionalFirstIndex(seed), outerOptionalFirstIndexMBA(seed); got != want {
			t.Fatalf("seed %08x = %d, want %d", seed, got, want)
		}
	}
	random := rand.New(rand.NewSource(0x26a428))
	for sample := 0; sample < 1<<16; sample++ {
		seed := random.Uint32()
		if got, want := outerOptionalFirstIndex(seed), outerOptionalFirstIndexMBA(seed); got != want {
			t.Fatalf("seed %08x = %d, want %d", seed, got, want)
		}
	}
}

func TestOuterOptionalActive(t *testing.T) {
	for _, value := range []uint32{0x976dec80, 0x976decff, 0x976dec0d} {
		if !outerOptionalActive(value) {
			t.Fatalf("rejected active value %#x", value)
		}
	}
	for _, value := range []uint32{0x976dec00, 0x976dec7f, 0x976ded80, 0} {
		if outerOptionalActive(value) {
			t.Fatalf("accepted inactive value %#x", value)
		}
	}
}
