package recovery

import (
	"math/rand"
	"testing"
)

func TestOuterHookSeedsMatchCircuit(t *testing.T) {
	random := rand.New(rand.NewSource(0x26a428))
	for sample := 0; sample < 100000; sample++ {
		var table [outerTableWords]uint32
		for index := range table {
			table[index] = random.Uint32()
		}
		digest := outerDigestState{
			X0: random.Uint32(), X4: random.Uint32(), X10: random.Uint32(),
			X27: random.Uint32(), X28: random.Uint32(),
		}
		base := uint32(random.Intn(int(outerTableWords)))
		got, gotOK := deriveOuterHookSeeds(&table, base, &digest)
		want, wantOK := deriveOuterHookSeedsCircuit(&table, base, &digest)
		if gotOK != wantOK || got != want {
			t.Fatalf("sample %d differs from circuit: got=%#v/%t want=%#v/%t", sample, got, gotOK, want, wantOK)
		}
	}
}
