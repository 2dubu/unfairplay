package recovery

import "testing"

func TestOuterTableSchedule(t *testing.T) {
	wantBases := []uint32{5, 11, 155, 104, 49, 65, 115, 146}
	base := outerTableInitialBase
	for iteration, want := range wantBases {
		if base != want {
			t.Fatalf("iteration %d base = %d, want %d", iteration, base, want)
		}
		indexes, ok := outerTableIndexes(base)
		if !ok {
			t.Fatalf("iteration %d rejected base %d", iteration, base)
		}
		seen := make(map[uint32]bool, len(indexes))
		for lane, index := range indexes {
			if index >= outerTableWords || seen[index] {
				t.Fatalf("iteration %d lane %d produced invalid or duplicate index %d", iteration, lane, index)
			}
			seen[index] = true
		}
		base, ok = nextOuterTableBase(base)
		if !ok {
			t.Fatalf("iteration %d failed to advance", iteration)
		}
	}
}

func TestOuterTableScheduleRejectsInvalidBase(t *testing.T) {
	if _, ok := outerTableIndexes(outerTableWords); ok {
		t.Fatal("accepted an out-of-range table base")
	}
	if _, ok := nextOuterTableBase(outerTableWords); ok {
		t.Fatal("advanced an out-of-range table base")
	}
}
