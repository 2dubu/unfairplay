package recovery

import (
	"crypto/sha256"
	"testing"
)

func TestSIDBProducerTypedMatchesSPNOracle(t *testing.T) {
	for control := range 4096 {
		digest := sha256.Sum256([]byte{0x54, 0x59, 0x50, 0x45, byte(control >> 8), byte(control)})
		var input [16]byte
		copy(input[:], digest[:16])
		if got, want := currentSIDBProducerTyped.apply(input), currentSIDBProducerSPN.apply(input); got != want {
			t.Fatalf("control %d = %x, want %x", control, got, want)
		}
	}
}
