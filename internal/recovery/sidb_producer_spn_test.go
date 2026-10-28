package recovery

import (
	"crypto/sha256"
	"encoding/hex"
	"testing"
)

func TestSIDBProducerSPNProfile(t *testing.T) {
	digest := sha256.Sum256(sidbProducerSPNProfileCompressed)
	if got := hex.EncodeToString(digest[:]); got != "0b1ce0c897216723802c38dcffd81875cf69afd1d2a164e80cb1391d619733ed" {
		t.Fatalf("SIDB producer SPN profile digest = %s", got)
	}
	if len(sidbProducerSPN26A428.stages) != 10 {
		t.Fatalf("SIDB producer SPN has %d stages", len(sidbProducerSPN26A428.stages))
	}
	hash := sha256.New()
	for index := range 4096 {
		digest := sha256.Sum256([]byte{0x53, 0x49, 0x44, 0x42, 0x4b, byte(index >> 8), byte(index)})
		var input [16]byte
		copy(input[:], digest[:16])
		key := deriveSelectionSIDBKey(input)
		_, _ = hash.Write(key[:])
	}
	if got := hex.EncodeToString(hash.Sum(nil)); got != "98b71a19908213bbd6b6499182d0d569c9554b5f0dc5c9e631adc0cb3df9c3ab" {
		t.Fatalf("SIDB producer key corpus digest = %s", got)
	}
}
