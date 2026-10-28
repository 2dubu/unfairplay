package recovery

import (
	"encoding/hex"
	"testing"
)

func TestLegacyAccountDigestSyntheticControl(t *testing.T) {
	digest, err := legacyAccountDigest([]byte{1, 2, 3, 4, 5, 6}, 0x12345678, 5)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := hex.EncodeToString(digest[:]), "b28371549cfd9d0c04ba7e50da504968cb8f0c4b"; got != want {
		t.Fatalf("legacy account digest = %s, want %s", got, want)
	}
}

func TestLegacyAccountDigestRejectsInvalidInputs(t *testing.T) {
	if _, err := legacyAccountDigest(make([]byte, 5), 1, 2); err == nil {
		t.Fatal("short hardware identity was accepted")
	}
	if _, err := legacyAccountDigest(make([]byte, 6), 1, 0x100); err == nil {
		t.Fatal("wide key selector was accepted")
	}
}
