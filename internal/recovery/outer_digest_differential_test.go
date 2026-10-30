package recovery

import (
	"math/rand"
	"testing"
)

func TestOuterDigestStepMatchesCircuit(t *testing.T) {
	random := rand.New(rand.NewSource(0x26a428))
	for sample := 0; sample < 100000; sample++ {
		state := outerDigestState{
			X0: random.Uint32(), X4: random.Uint32(), X10: random.Uint32(),
			X27: random.Uint32(), X28: random.Uint32(), Counter: random.Uint32(),
		}
		var hardware [6]byte
		var codec [33]byte
		for index := range hardware {
			hardware[index] = byte(random.Uint32())
		}
		for index := 4; index < 17; index++ {
			codec[index] = byte(random.Uint32())
		}
		hardwareBase := uint32(random.Intn(6))
		codecPhase := uint32(random.Intn(13))
		got, want := state, state
		gotExit := applyOuterDigestStep(&got, &hardware, &codec, hardwareBase, codecPhase)
		wantExit := applyOuterDigestStepCircuit(&want, &hardware, &codec, hardwareBase, codecPhase)
		if gotExit != wantExit || got != want {
			t.Fatalf("sample %d differs from circuit: got=%#v/%t want=%#v/%t", sample, got, gotExit, want, wantExit)
		}
	}
}
