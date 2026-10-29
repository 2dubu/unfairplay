package recovery

import (
	"math/rand"
	"testing"
)

func TestOuterCycleMatchesProtectedCircuit(t *testing.T) {
	random := rand.New(rand.NewSource(0x6f75746572637963))
	for control := range 4096 {
		var input outerTableRoundState
		for index := range input.Words {
			input.Words[index] = random.Uint32()
		}
		for index := range input.Hardware {
			input.Hardware[index] = byte(random.Uint32())
		}
		for index := range input.Codec {
			input.Codec[index] = byte(random.Uint32())
		}
		input.HardwareIndex = byte(control % len(input.Hardware))
		cycleIndex := uint32(control%39 + 1)
		got, want := input, input
		if !applyOuterCycle(&got, cycleIndex) || !applyOuterCycleCircuit(&want, cycleIndex) {
			t.Fatalf("control %d rejected", control)
		}
		if got != want {
			t.Fatalf("control %d differs: got=%#v want=%#v", control, got, want)
		}
	}
}

func TestOuterOptionalDecodeMatchesCircuit(t *testing.T) {
	for word := uint32(0); word < 1<<16; word++ {
		if got, want := outerOptionalDecode(word), outerOptionalDecodeCircuit(word); got != want {
			t.Fatalf("word %08x = %08x, want %08x", word, got, want)
		}
	}
	random := rand.New(rand.NewSource(0x26a428))
	for sample := 0; sample < 1<<16; sample++ {
		word := random.Uint32()
		if got, want := outerOptionalDecode(word), outerOptionalDecodeCircuit(word); got != want {
			t.Fatalf("word %08x = %08x, want %08x", word, got, want)
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

func TestOuterWriterSeedStepMatchesCircuit(t *testing.T) {
	for seed := uint32(0); seed < 1<<16; seed++ {
		if got, want := outerWriterSeedStep(seed), outerWriterSeedStepMBA(seed); got != want {
			t.Fatalf("seed %08x = %08x, want %08x", seed, got, want)
		}
	}
	random := rand.New(rand.NewSource(0x26a428))
	for sample := 0; sample < 1<<16; sample++ {
		seed := random.Uint32()
		if got, want := outerWriterSeedStep(seed), outerWriterSeedStepMBA(seed); got != want {
			t.Fatalf("seed %08x = %08x, want %08x", seed, got, want)
		}
	}
}

func TestOuterWriterSeedNativeCircuitExhaustive(t *testing.T) {
	for phase := uint32(0); phase < outerWriterSeedStates; phase++ {
		seed := outerWriterSeedBase + phase
		want := outerWriterSeedBase + (phase+outerWriterSeedStride)%outerWriterSeedStates
		if got := outerWriterSeedStepMBA(seed); got != want {
			t.Fatalf("phase %d: MBA=%#x, recurrence=%#x", phase, got, want)
		}
	}
}

func TestOuterOptionalNativeCircuitExhaustive(t *testing.T) {
	for phase := uint32(0); phase < outerWriterSeedStates; phase++ {
		seed := outerWriterSeedBase + phase
		want := (phase + outerOptionalPhaseOffset) % outerWriterSeedStates
		if got := outerOptionalFirstIndexMBA(seed); got != want {
			t.Fatalf("phase %d: MBA=%d, semantic=%d", phase, got, want)
		}
	}
}
