package recovery

import (
	"math/rand"
	"testing"
)

func TestOuterFrontMatchesCircuit(t *testing.T) {
	random := rand.New(rand.NewSource(0x26a428))
	controls := [...]uint32{0x60319d83, 0x60319e4b, 0x60319d87}
	for sample := 0; sample < 10000; sample++ {
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
		x1 := controls[random.Intn(len(controls))]
		counterSeed := random.Uint32()
		gotState, wantState := state, state
		var gotHooks, wantHooks outerHookSeeds
		gotOK := applyOuterFront(&table, base, &gotState, x1, counterSeed, &gotHooks)
		wantOK := applyOuterFrontCircuit(&table, base, &wantState, x1, counterSeed, &wantHooks)
		if gotOK != wantOK || gotState != wantState || gotHooks != wantHooks {
			t.Fatalf("sample %d differs from circuit", sample)
		}
	}
}
