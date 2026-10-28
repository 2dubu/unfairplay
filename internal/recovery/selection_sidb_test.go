package recovery

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"testing"
)

func TestDecodeSelectionSIDBContainer(t *testing.T) {
	for _, control := range []struct {
		name          string
		containerSize int
		payloadOffset int
		payloadBytes  int
	}{
		{name: "older", containerSize: selectionSIDBOlderContainerBytes, payloadOffset: selectionSIDBOlderPayloadOffset, payloadBytes: selectionSIDBOlderPayloadBytes},
		{name: "compact", containerSize: selectionSIDBContainerBytes, payloadOffset: selectionSIDBPayloadOffset, payloadBytes: selectionSIDBPayloadBytes},
		{name: "extended", containerSize: selectionSIDBExtendedBytes, payloadOffset: selectionSIDBExtendedPayloadOffset, payloadBytes: selectionSIDBPayloadBytes},
	} {
		t.Run(control.name, func(t *testing.T) {
			input := make([]byte, control.containerSize)
			for index := range input {
				input[index] = byte(index*37 + 11)
			}
			binary.BigEndian.PutUint32(input[:4], currentSIDBContainerVersion)
			output, err := decodeSelectionSIDBContainer(input)
			if err != nil {
				t.Fatal(err)
			}
			if len(output) != control.payloadBytes {
				t.Fatalf("decoded payload bytes = %d, want %d", len(output), control.payloadBytes)
			}
			for index, value := range output {
				want := input[control.payloadOffset+index] ^ selectionSIDBInputMask[index%selectionSIDBBlockBytes]
				if value != want {
					t.Fatalf("decoded byte %d=%#x, want %#x", index, value, want)
				}
			}
			if control.name == "extended" && control.payloadOffset+len(output)+selectionSIDBExtendedTrailerBytes != len(input) {
				t.Fatal("extended SIDB trailer geometry changed")
			}
		})
	}

	invalidVersion := make([]byte, selectionSIDBContainerBytes)
	if _, err := decodeSelectionSIDBContainer(invalidVersion); err == nil {
		t.Fatal("accepted a selection SIDB container with an invalid version")
	}
	truncated := make([]byte, selectionSIDBContainerBytes-1)
	binary.BigEndian.PutUint32(truncated[:4], currentSIDBContainerVersion)
	if _, err := decodeSelectionSIDBContainer(truncated); err == nil {
		t.Fatal("accepted a truncated selection SIDB container")
	}
	oversized := make([]byte, selectionSIDBExtendedBytes+1)
	binary.BigEndian.PutUint32(oversized[:4], currentSIDBContainerVersion)
	if _, err := decodeSelectionSIDBContainer(oversized); err == nil {
		t.Fatal("accepted an oversized selection SIDB container")
	}
	for _, size := range []int{0, 1, 2, 3} {
		if _, err := decodeSelectionSIDBContainer(make([]byte, size)); err == nil {
			t.Fatalf("accepted a %d-byte selection SIDB container", size)
		}
	}
}

func TestDecryptSelectionSIDBPayload(t *testing.T) {
	var input [selectionSIDBPayloadBytes]byte
	for index := range input {
		input[index] = byte(index*29 + index/16*7 + 3)
	}
	key := [selectionSIDBBlockBytes]byte{
		0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77,
		0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff,
	}
	output, err := decryptSelectionSIDBPayload(input[:], key)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := fmt.Sprintf("%x", sha256.Sum256(output[:])), "36ae1110e54116748e1599b6354ca0026f94c17b93ba13c8f6ff146e934e0ec8"; got != want {
		t.Fatalf("selection SIDB payload digest=%s, want %s", got, want)
	}
	// The chaining rule is external-domain feedback, not CBC. A change in the
	// first encoded block therefore reaches the same byte of block two after
	// its independent AES decryption.
	mutated := input
	mutated[0] ^= 1
	changed, err := decryptSelectionSIDBPayload(mutated[:], key)
	if err != nil {
		t.Fatal(err)
	}
	if changed[selectionSIDBBlockBytes] != output[selectionSIDBBlockBytes]^1 {
		t.Fatal("selection SIDB external chaining was not preserved")
	}
}
