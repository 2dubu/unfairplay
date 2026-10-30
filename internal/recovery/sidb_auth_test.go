package recovery

import (
	"encoding/hex"
	"testing"
)

func TestSIDBContainerAuthentication(t *testing.T) {
	// Expected footer is AES-128-ECB over the MD5 of this public wire prefix,
	// independently computed with OpenSSL using key 000102...0f.
	prefix := append([]byte{0, 7, 0, 11}, []byte("public synthetic authentication fixture")...)
	footer, err := hex.DecodeString("75d682d73ca1534716a02fa2479b3fae")
	if err != nil {
		t.Fatal(err)
	}
	container := append(prefix, footer...)
	source := make([]byte, selectionSIDBLargerSourceBytes)
	for i := range 16 {
		source[i] = byte(i) ^ 0x1b
	}
	if err := validateSIDBContainerAuthentication(container, source); err != nil {
		t.Fatal(err)
	}
	for _, offset := range []int{0, 4, len(container) - 1} {
		changed := append([]byte(nil), container...)
		changed[offset] ^= 1
		if err := validateSIDBContainerAuthentication(changed, source); err == nil {
			t.Fatalf("accepted changed wire byte %d", offset)
		}
	}
	changedSource := append([]byte(nil), source...)
	changedSource[0] ^= 1
	if err := validateSIDBContainerAuthentication(container, changedSource); err == nil {
		t.Fatal("accepted changed key source")
	}
	for _, size := range []int{0, 15, 32, 255, 257, 321, 336} {
		if err := validateSIDBContainerAuthentication(container, make([]byte, size)); err == nil {
			t.Fatalf("accepted %d-byte source", size)
		}
	}
	if err := validateSIDBContainerAuthentication(make([]byte, 15), source); err == nil {
		t.Fatal("accepted truncated footer")
	}
}
