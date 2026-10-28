package recovery

import (
	"bytes"
	"crypto/aes"
	"encoding/binary"
	"testing"
)

func syntheticSelectionSIDBContainer(t *testing.T, plaintext [selectionSIDBPayloadBytes]byte, key [16]byte) []byte {
	t.Helper()
	container := make([]byte, selectionSIDBContainerBytes)
	binary.BigEndian.PutUint32(container[:4], currentSIDBContainerVersion)
	fillSyntheticSelectionSIDBContainer(t, container, selectionSIDBPayloadOffset, plaintext, key)
	return container
}

func fillSyntheticSelectionSIDBContainer(t *testing.T, container []byte, payloadOffset int, plaintext [selectionSIDBPayloadBytes]byte, key [16]byte) {
	t.Helper()
	block, err := aes.NewCipher(key[:])
	if err != nil {
		t.Fatal(err)
	}
	var previousInput [selectionSIDBBlockBytes]byte
	for offset := 0; offset < len(plaintext); offset += selectionSIDBBlockBytes {
		var decoded, ciphertext [selectionSIDBBlockBytes]byte
		for index := range decoded {
			decoded[index] = plaintext[offset+index] ^ selectionSIDBOutputMask[index]
			if offset != 0 {
				decoded[index] ^= previousInput[index] ^ selectionSIDBChainMask[index]
			}
		}
		block.Encrypt(ciphertext[:], decoded[:])
		copy(container[payloadOffset+offset:], ciphertext[:])
		for index := range previousInput {
			previousInput[index] = ciphertext[index] ^ selectionSIDBInputMask[index]
		}
	}
}

func TestDeriveSelectionRegistryCandidates(t *testing.T) {
	profile := embedded26A428ExecutionProfile()
	source := make([]byte, selectionSourceEnvelopeBytes)
	for index := range source {
		source[index] = byte(index*23 + 9)
	}
	frames, err := enrollmentFrameCandidates(source[selectionSourcePrefixBytes:])
	if err != nil {
		t.Fatal(err)
	}
	const selectedSlot = 8
	var plaintext [selectionSIDBPayloadBytes]byte
	for index := range plaintext {
		plaintext[index] = byte(index*41 + 17)
	}
	container := syntheticSelectionSIDBContainer(t, plaintext, frames[selectedSlot].Key)

	pageBase := profile.sourceAddress &^ (pageSize - 1)
	page := make([]byte, pageSize)
	copy(page[profile.sourceAddress-pageBase:], source)
	s := newState(&model{binaryHash: profile.binaryHash, pages: map[uint64][]byte{pageBase: page}})
	runtime := &rawSourceRuntime{files: map[string][]byte{"/private/var/db/fpsd/SC Info/SC Info.sidb": container}}
	candidates, err := deriveSelectionRegistryCandidates(s, runtime)
	if err != nil {
		t.Fatal(err)
	}
	matches := 0
	for _, candidate := range candidates {
		if bytes.Equal(candidate.Data, plaintext[:]) {
			matches++
		}
	}
	if matches != 1 {
		t.Fatalf("typed SIDB selection found plaintext %d times, want 1", matches)
	}
}

func TestDeriveExtendedCurrentSelectionRegistryCandidates(t *testing.T) {
	container, hardware := syntheticExtendedCurrentOuterInputs()
	for index := selectionSIDBExtendedPayloadOffset; index < selectionSIDBExtendedPayloadOffset+selectionSIDBPayloadBytes; index++ {
		container[index] = 0
	}
	source, iterations, err := deriveCurrentSelectionSource(container, hardware)
	if err != nil || iterations != selectionSIDBExtendedIterations {
		t.Fatalf("derive extended selection source = %d/%v", iterations, err)
	}
	frames, err := enrollmentFrameCandidates(source[selectionSourcePrefixBytes:])
	if err != nil {
		t.Fatal(err)
	}
	const selectedSlot = 8
	var plaintext [selectionSIDBPayloadBytes]byte
	for index := range plaintext {
		plaintext[index] = byte(index*31 + 19)
	}
	fillSyntheticSelectionSIDBContainer(t, container, selectionSIDBExtendedPayloadOffset, plaintext, frames[selectedSlot].Key)
	candidates, _, gotIterations, err := deriveCurrentSelectionRegistryCandidates(
		map[string][]byte{"/private/var/db/fpsd/SC Info/SC Info.sidb": container}, hardware,
	)
	if err != nil {
		t.Fatal(err)
	}
	if gotIterations != selectionSIDBExtendedIterations {
		t.Fatalf("extended selection iterations=%d", gotIterations)
	}
	matches := 0
	for _, candidate := range candidates {
		if bytes.Equal(candidate.Data, plaintext[:]) {
			matches++
		}
	}
	clearDecodedSIDBRegistryCandidates(candidates)
	if matches != 1 {
		t.Fatalf("extended typed SIDB selection found plaintext %d times, want 1", matches)
	}
}
