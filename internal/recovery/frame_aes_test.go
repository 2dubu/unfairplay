package recovery

import (
	"crypto/sha256"
	"encoding/hex"
	"testing"
)

func TestFrameAESProfile(t *testing.T) {
	digest := sha256.Sum256(frameAESProfileCompressed)
	if got := hex.EncodeToString(digest[:]); got != "9d42fb362c39a1f5e1c152fb82fa371869afa0701bec3065b89fba10915af202" {
		t.Fatalf("frame AES profile digest = %s", got)
	}
	for control := range 4096 {
		digest := sha256.Sum256([]byte{0x46, 0x41, 0x45, 0x53, byte(control >> 8), byte(control)})
		var input [16]byte
		copy(input[:], digest[:16])
		got := currentFrameSPN.finalOutput.apply(currentFrameAES.apply(input))
		want, _ := currentFrameSPN.applyInternal(input)
		if got != want {
			t.Fatalf("control %d = %x, want %x", control, got, want)
		}
	}
}
