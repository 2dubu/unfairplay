package recovery

import "testing"

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
