package recovery

import "testing"

// The standalone hook tests derive the shared core through the front path.
func deriveOuterHookSeeds(table *[outerTableWords]uint32, base uint32, digest *outerDigestState) (outerHookSeeds, bool) {
	_, core, ok := outerFrontWords(table, base, digest)
	if !ok {
		return outerHookSeeds{}, false
	}
	return deriveOuterHookSeedsWithCore(table, base, digest, core)
}

func TestOuterHookSeedsDeterministic(t *testing.T) {
	var table [outerTableWords]uint32
	for index := range table {
		table[index] = uint32(index)*0x9e3779b9 ^ 0x7f4a7c15
	}
	digest := outerDigestState{X0: 0x01234567, X4: 0x89abcdef, X10: 0x0f1e2d3c, X27: 0x4b5a6978, X28: 0xdeadbeef}
	got, ok := deriveOuterHookSeeds(&table, 49, &digest)
	if !ok {
		t.Fatal("hook seeds rejected valid input")
	}
	want := outerHookSeeds{Even: 0xdeadbca9, Odd: 0xb34f6db4, Tail: 0x8120781d}
	if got != want {
		t.Fatalf("hook seeds = %#v, want %#v", got, want)
	}
}

func TestOuterHookSeedsRejectInvalidInput(t *testing.T) {
	var table [outerTableWords]uint32
	var digest outerDigestState
	if _, ok := deriveOuterHookSeeds(nil, 0, &digest); ok {
		t.Fatal("accepted nil table")
	}
	if _, ok := deriveOuterHookSeeds(&table, 0, nil); ok {
		t.Fatal("accepted nil digest")
	}
	if _, ok := deriveOuterHookSeeds(&table, outerTableWords, &digest); ok {
		t.Fatal("accepted invalid base")
	}
}
