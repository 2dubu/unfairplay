package recovery

import "testing"

func TestSelectedEnrollmentKeyIgnoresLaterBlocks(t *testing.T) {
	source := make([]byte, 18*16)
	for i := range source {
		source[i] = byte(i*41 + i/16*13 + 5)
	}
	for nativeCount := 1; nativeCount <= 18; nativeCount++ {
		for slot := 0; slot < nativeCount && slot < selectionStateSlotCount; slot++ {
			key, err := deriveSelectedEnrollmentKey(source[:nativeCount*16], slot)
			if err != nil {
				t.Fatal(err)
			}
			changed := append([]byte(nil), source[:nativeCount*16]...)
			for i := (slot + 1) * 16; i < len(changed); i++ {
				changed[i] ^= 0xff
			}
			again, err := deriveSelectedEnrollmentKey(changed, slot)
			if err != nil || again != key {
				t.Fatalf("slot %d depends on later inputs: %v", slot, err)
			}
		}
	}
}
func TestSelectedEnrollmentKeyRejectsInvalidInput(t *testing.T) {
	for _, length := range []int{0, 15, 17, 19 * 16} {
		if _, err := deriveSelectedEnrollmentKey(make([]byte, length), 0); err == nil {
			t.Fatalf("accepted length %d", length)
		}
	}
	source := make([]byte, 16)
	for _, slot := range []int{-1, 1, selectionStateSlotCount} {
		if _, err := deriveSelectedEnrollmentKey(source, slot); err == nil {
			t.Fatalf("accepted slot %d", slot)
		}
	}
}

func TestSelectedEnrollmentKeyRejectsUnverifiedSlot(t *testing.T) {
	source := make([]byte, 18*16)
	if _, err := deriveSelectedEnrollmentKey(source, 13); err != nil {
		t.Fatalf("verified slot in larger source: %v", err)
	}
	for _, slot := range []int{14, 15, 16, 17} {
		if _, err := deriveSelectedEnrollmentKey(source, slot); err == nil {
			t.Fatalf("accepted unverified slot %d", slot)
		}
	}
}
