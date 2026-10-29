package recovery

import (
	"crypto/sha256"
	"fmt"
	"testing"
)

func TestDecodeSelectionSIDBContainer(t *testing.T) {
	for _, payloadBytes := range []int{576, selectionSIDBPayloadBytes, 64 << 10} {
		input, _, layout := syntheticNativeSIDB(t, payloadBytes)
		output, err := decodeSelectionSIDBContainerWithLayout(input, layout)
		if err != nil || len(output) != payloadBytes {
			t.Fatalf("projection failed: %v", err)
		}
		for i, value := range output {
			if value != input[layout.PayloadOffset+i]^selectionSIDBInputMask[i%selectionSIDBBlockBytes] {
				t.Fatalf("projection mismatch at %d", i)
			}
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
	output, err := decryptSelectionSIDBPayloadDynamic(input[:], key)
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
	changed, err := decryptSelectionSIDBPayloadDynamic(mutated[:], key)
	if err != nil {
		t.Fatal(err)
	}
	if changed[selectionSIDBBlockBytes] != output[selectionSIDBBlockBytes]^1 {
		t.Fatal("selection SIDB external chaining was not preserved")
	}
}
