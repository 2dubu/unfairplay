package recovery

import (
	"math/rand"
	"testing"
)

func TestOuterFrontWordsAndSourceTailMatchCircuit(t *testing.T) {
	random := rand.New(rand.NewSource(0x26a428))
	for sample := 0; sample < 100000; sample++ {
		var table [outerTableWords]uint32
		for index := range table {
			table[index] = random.Uint32()
		}
		digest := outerDigestState{X0: random.Uint32(), X4: random.Uint32(), X10: random.Uint32(), X27: random.Uint32(), X28: random.Uint32()}
		base := uint32(random.Intn(int(outerTableWords)))
		oldWord1, oldWord4 := random.Uint32(), random.Uint32()
		state := outerTableRoundState{Words: [5]uint32{1: oldWord1, 4: oldWord4}}
		if !applyOuterFrontTailCircuit(&table, base, &state, &digest) {
			t.Fatal("circuit rejected input")
		}
		words, ok := outerFrontWords(&table, base, &digest)
		if !ok || words != state.Words {
			t.Fatalf("sample %d base %d words=%08x/%08x/%08x/%08x/%08x want=%08x/%08x/%08x/%08x/%08x", sample, base,
				words[0], words[1], words[2], words[3], words[4], state.Words[0], state.Words[1], state.Words[2], state.Words[3], state.Words[4])
		}
		tail, ok := outerFrontSourceTail(&table, base, oldWord1, oldWord4)
		if !ok || tail != [4]byte(state.Source[12:16]) {
			t.Fatalf("sample %d source tail differs", sample)
		}
		prefix, ok := outerFrontSourcePrefix(&table, base, &digest)
		if !ok || prefix != [12]byte(state.Source[:12]) {
			t.Fatalf("sample %d source prefix differs", sample)
		}
	}
}
