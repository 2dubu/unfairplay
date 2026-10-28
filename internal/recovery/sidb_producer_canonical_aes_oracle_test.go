package recovery

import (
	"crypto/sha256"
	"testing"
)

func TestSIDBProducerCanonicalAESMatchesSPNOracle(t *testing.T) {
	for control := range 4096 {
		digest := sha256.Sum256([]byte{0x43, 0x41, 0x4e, 0x4f, byte(control >> 8), byte(control)})
		var input [16]byte
		copy(input[:], digest[:16])
		if got, want := currentSIDBProducerCanonicalAES.apply(input), currentSIDBProducerSPN.apply(input); got != want {
			t.Fatalf("control %d = %x, want %x", control, got, want)
		}
	}
}
