package recovery

import (
	"crypto/sha256"
	"encoding/hex"
	"testing"
)

// Public synthetic controls were evaluated through the original frame-core,
// suffix, and producer composition. Expected values do not use the AES
// implementation under test. Complete derivation proofs run with -tags research.
func TestSIDBKeyKnownAnswers(t *testing.T) {
	controls := [][2]string{
		{"00000000000000000000000000000000", "21b5c58d38b5be627bb53a26bf8b80a6"},
		{"ffffffffffffffffffffffffffffffff", "c7887d14ff3361ad63e8792b99099563"},
		{"000102030405060708090a0b0c0d0e0f", "bcd67ece0b6b78c8d306b8ba2b519a30"},
		{"80000000000000000000000000000000", "90b912d34e78dd99264e0975e9b21652"},
		{"00000000000000000000000000000001", "367aba9c2df9d44413c32e4fa7aa14a0"},
		{"0123456789abcdef0123456789abcdef", "6a5c0996656510bf9fa4ce7896a33e1e"},
	}
	for _, control := range controls {
		t.Run(control[0], func(t *testing.T) {
			raw, err := hex.DecodeString(control[0])
			if err != nil {
				t.Fatal(err)
			}
			var input [16]byte
			copy(input[:], raw)
			got := deriveSelectionSIDBKey(input)
			if hex.EncodeToString(got[:]) != control[1] {
				t.Fatalf("key=%x, want %s", got, control[1])
			}
		})
	}
}

// Preserve the original 4,096-input corpus digest from the upstream producer
// test without loading that historical producer or either profile.
func TestSIDBKeyCorpus(t *testing.T) {
	hash := sha256.New()
	for index := range 4096 {
		digest := sha256.Sum256([]byte{0x53, 0x49, 0x44, 0x42, 0x4b, byte(index >> 8), byte(index)})
		var input [16]byte
		copy(input[:], digest[:16])
		key := deriveSelectionSIDBKey(input)
		_, _ = hash.Write(key[:])
	}
	const want = "98b71a19908213bbd6b6499182d0d569c9554b5f0dc5c9e631adc0cb3df9c3ab"
	if got := hex.EncodeToString(hash.Sum(nil)); got != want {
		t.Fatalf("key corpus digest=%s, want %s", got, want)
	}
}
