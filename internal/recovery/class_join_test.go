package recovery

import (
	"crypto/aes"
	"testing"
)

func TestDeriveCurrentSINFPrecursorFromClass(t *testing.T) {
	classKey := [aes.BlockSize]byte{
		0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
		0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
	}
	plaintext := [aes.BlockSize]byte{
		0x6b, 0xc1, 0xbe, 0xe2, 0x2e, 0x40, 0x9f, 0x96,
		0xe9, 0x3d, 0x7e, 0x11, 0x73, 0x93, 0x17, 0x2a,
	}
	block, err := aes.NewCipher(classKey[:])
	if err != nil {
		t.Fatal(err)
	}
	var valu [aes.BlockSize]byte
	block.Encrypt(valu[:], plaintext[:])
	precursor, err := deriveCurrentSINFPrecursorFromClass(classKey, valu)
	if err != nil {
		t.Fatal(err)
	}
	for index, value := range plaintext {
		if want := applyCurrentFairPlayByteMask(value); precursor[index] != want {
			t.Fatalf("precursor byte %d = %#02x, want %#02x", index, precursor[index], want)
		}
	}
	if got, want := currentSINFModelKeyFromPrecursor(precursor), plaintext; got != want {
		t.Fatalf("model key = %x, want %x", got, want)
	}
}
