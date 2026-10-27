package recovery

import (
	"math/rand"
	"testing"
)

func TestOuterPreRoundWord1MatchesCircuit(t *testing.T) {
	random := rand.New(rand.NewSource(0x26a428))
	for sample := 0; sample < 100000; sample++ {
		var table [outerTableWords]uint32
		for index := range table {
			table[index] = random.Uint32()
		}
		state := outerTableRoundState{
			Words:         [5]uint32{random.Uint32(), random.Uint32(), random.Uint32(), random.Uint32(), random.Uint32()},
			HardwareIndex: uint8(random.Intn(6)),
		}
		for index := range state.Hardware {
			state.Hardware[index] = byte(random.Uint32())
		}
		for index := range state.Codec {
			state.Codec[index] = byte(random.Uint32())
		}
		base := uint32(random.Intn(int(outerTableWords)))
		wantState := state
		if !applyOuterPreRound(&table, base, &wantState) {
			t.Fatal("circuit rejected input")
		}
		got, ok := outerPreRoundWord1(&table, base, &state)
		if !ok || got != wantState.Words[1] {
			t.Fatalf("sample %d = %08x/%t, want %08x", sample, got, ok, wantState.Words[1])
		}
	}
}
