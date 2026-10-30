package recovery

import (
	"math/rand"
	"testing"
)

func TestOuterPreRoundMatchesCircuit(t *testing.T) {
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
		wantState, gotState := state, state
		if !applyOuterPreRoundCircuit(&table, base, &wantState) {
			t.Fatal("circuit rejected input")
		}
		if !applyOuterPreRound(&table, base, &gotState) {
			t.Fatal("reduced implementation rejected input")
		}
		if gotState.Words != wantState.Words {
			t.Fatalf("sample %d = %08x, want %08x", sample, gotState.Words, wantState.Words)
		}
	}
}
