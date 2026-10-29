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
			if control.name == "extended" && control.payloadOffset+len(output)+selectionSIDBTrailerBytes != len(input) {
				t.Fatal("extended SIDB trailer geometry changed")
			}
		})
	}

	structural := make([]byte, 4+28+selectionSIDBPayloadBytes+selectionSIDBTrailerBytes)
	binary.BigEndian.PutUint32(structural[:4], currentSIDBContainerVersion)
	for index := 4; index < len(structural); index++ {
		structural[index] = byte(index*13 + 5)
	}
	layout, err := currentSIDBLayoutForContainer(structural)
	if err != nil {
		t.Fatal(err)
	}
	if layout.CodecBytes != 28 || layout.PayloadOffset != 32 || layout.PayloadBytes != selectionSIDBPayloadBytes || layout.OuterIterations != 306 {
		t.Fatalf("structural SIDB layout = %+v", layout)
	}
	decoded, err := decodeSelectionSIDBContainer(structural)
	if err != nil || len(decoded) != selectionSIDBPayloadBytes {
		t.Fatalf("decode structural SIDB = %d, %v", len(decoded), err)
	}

	large := make([]byte, 11724)
	binary.BigEndian.PutUint32(large[:4], currentSIDBContainerVersion)
	largeLayout, err := currentSIDBLayoutForContainer(large)
	if err != nil {
		t.Fatal(err)
	}
	if largeLayout.CodecBytes != 24 || largeLayout.PayloadOffset != 28 || largeLayout.PayloadBytes != 11680 || largeLayout.OuterIterations != 262 {
		t.Fatalf("large structural SIDB layout = %+v", largeLayout)
	}

	invalidVersion := make([]byte, selectionSIDBContainerBytes)
	if _, err := decodeSelectionSIDBContainer(invalidVersion); err == nil {
		t.Fatal("accepted a selection SIDB container with an invalid version")
	}
	for _, size := range []int{4 + 12 + 64 + selectionSIDBTrailerBytes} {
		unsupported := make([]byte, size)
		binary.BigEndian.PutUint32(unsupported[:4], currentSIDBContainerVersion)
		if _, err := decodeSelectionSIDBContainer(unsupported); err == nil {
			t.Fatalf("accepted unsupported %d-byte selection SIDB container", size)
		}
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
