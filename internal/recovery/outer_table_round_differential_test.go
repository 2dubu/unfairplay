package recovery

import (
	"math/rand"
	"testing"
)

func TestOuterTableRoundReducedMatchesCircuit(t *testing.T) {
	random := rand.New(rand.NewSource(0x26a428))
	for base := uint32(0); base < outerTableWords; base++ {
		for sample := 0; sample < 16; sample++ {
			var table [outerTableWords]uint32
			for index := range table {
				table[index] = random.Uint32()
			}
			state := outerTableRoundState{
				Words: [5]uint32{
					random.Uint32(), random.Uint32(), random.Uint32(),
					random.Uint32(), random.Uint32(),
				},
				HardwareIndex: uint8(random.Intn(6)),
			}
			for index := range state.Hardware {
				state.Hardware[index] = byte(random.Uint32())
			}
			for index := range state.Codec {
				state.Codec[index] = byte(random.Uint32())
			}

			circuitTable, recoveredTable := table, table
			circuitState, recoveredState := state, state
			circuitNext, circuitOK := applyOuterTableRoundCircuit(&circuitTable, base, &circuitState)
			recoveredNext, recoveredOK := applyOuterTableRoundReduced(&recoveredTable, base, &recoveredState)
			if circuitOK != recoveredOK || circuitNext != recoveredNext ||
				circuitTable != recoveredTable || circuitState.Words != recoveredState.Words ||
				circuitState.HardwareIndex != recoveredState.HardwareIndex {
				t.Fatalf("base %d sample %d differs from circuit oracle", base, sample)
			}
		}
	}
}
