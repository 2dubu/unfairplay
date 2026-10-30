package recovery

import "testing"

func TestRecoverRequiresInputs(t *testing.T) {
	if _, err := Recover(Options{}); err == nil {
		t.Fatal("empty recovery input was accepted")
	}
}
