package recovery

import (
	"crypto/aes"
	"testing"
)

func TestDecodeClassSecondary(t *testing.T) {
	key := [aes.BlockSize]byte{
		0x01, 0x12, 0x23, 0x34, 0x45, 0x56, 0x67, 0x78,
		0x01, 0x12, 0x23, 0x34, 0x45, 0x56, 0x67, 0x78,
	}
	var encoded [aes.BlockSize]byte
	for index := range encoded {
		encoded[index] = key[index] ^ 0x1b
	}
	got, err := decodeClassSecondary(encoded)
	if err != nil {
		t.Fatal(err)
	}
	if got != key {
		t.Fatalf("decoded key = %x, want %x", got, key)
	}
	encoded[aes.BlockSize-1] ^= 1
	if _, err := decodeClassSecondary(encoded); err == nil {
		t.Fatal("accepted unequal class secondary halves")
	}
}

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
