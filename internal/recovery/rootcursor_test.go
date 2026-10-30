package recovery

import "testing"

func TestNativeSIDBRootCursorRejectsInvalidHardware(t *testing.T) {
	for _, n := range []int{0, 5, 7} {
		if _, err := nativeSIDBRootCursor(make([]byte, n)); err == nil {
			t.Fatalf("accepted length %d", n)
		}
	}
}
