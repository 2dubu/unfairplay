package recovery

import (
	"bytes"
	"crypto/aes"
	"encoding/binary"
	"strings"
	"testing"
)

func fillSyntheticSelectionSIDBContainer(t *testing.T, container []byte, payloadOffset int, plaintext []byte, key [16]byte) {
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

func TestDeriveExtendedCurrentSelectionRegistryCandidates(t *testing.T) {
	container, hardware := syntheticExtendedCurrentOuterInputs()
	for index := selectionSIDBExtendedPayloadOffset; index < selectionSIDBExtendedPayloadOffset+selectionSIDBPayloadBytes; index++ {
		container[index] = 0
	}
	source, iterations, err := deriveCurrentSelectionSource(container, hardware)
	if err != nil || iterations != selectionSIDBExtendedIterations {
		t.Fatalf("derive extended selection source = %d/%v", iterations, err)
	}
	keyCandidates, err := enrollmentKeyCandidates(source[selectionSourcePrefixBytes:])
	if err != nil {
		t.Fatal(err)
	}
	const selectedSlot = 8
	var plaintext [selectionSIDBPayloadBytes]byte
	for index := range plaintext {
		plaintext[index] = byte(index*31 + 19)
	}
	fillSyntheticSelectionSIDBContainer(t, container, selectionSIDBExtendedPayloadOffset, plaintext[:], keyCandidates[selectedSlot].Key)
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

func TestDeriveStructuralCurrentSelectionRegistryCandidates(t *testing.T) {
	container, hardware := syntheticStructuralCurrentOuterInputs(28, selectionSIDBPayloadBytes)
	layout, err := currentSIDBLayoutForContainer(container)
	if err != nil {
		t.Fatal(err)
	}
	for index := layout.PayloadOffset; index < layout.PayloadOffset+layout.PayloadBytes; index++ {
		container[index] = 0
	}
	source, iterations, err := deriveCurrentSelectionSource(container, hardware)
	if err != nil || iterations != 306 {
		t.Fatalf("derive structural selection source = %d/%v", iterations, err)
	}
	keyCandidates, err := enrollmentKeyCandidates(source[selectionSourcePrefixBytes:])
	if err != nil {
		t.Fatal(err)
	}
	const selectedSlot = 8
	var plaintext [selectionSIDBPayloadBytes]byte
	for index := range plaintext {
		plaintext[index] = byte(index*23 + 11)
	}
	fillSyntheticSelectionSIDBContainer(t, container, layout.PayloadOffset, plaintext[:], keyCandidates[selectedSlot].Key)
	candidates, _, gotIterations, err := deriveCurrentSelectionRegistryCandidates(
		map[string][]byte{"/private/var/db/fpsd/SC Info/SC Info.sidb": container}, hardware,
	)
	if err != nil {
		t.Fatal(err)
	}
	if gotIterations != 306 {
		t.Fatalf("structural selection iterations=%d", gotIterations)
	}
	matches := 0
	for _, candidate := range candidates {
		if bytes.Equal(candidate.Data, plaintext[:]) {
			matches++
		}
	}
	clearDecodedSIDBRegistryCandidates(candidates)
	if matches != 1 {
		t.Fatalf("structural typed SIDB selection found plaintext %d times, want 1", matches)
	}
}

func TestStructuralSIDBAndExtendedSINFProduceEnrollmentModel(t *testing.T) {
	container, hardware := syntheticStructuralCurrentOuterInputs(28, selectionSIDBPayloadBytes)
	layout, err := currentSIDBLayoutForContainer(container)
	if err != nil {
		t.Fatal(err)
	}
	source, iterations, err := deriveCurrentSelectionSource(container, hardware)
	if err != nil || iterations != 306 {
		t.Fatalf("derive structural source = %d/%v", iterations, err)
	}
	keyCandidates, err := enrollmentKeyCandidates(source[selectionSourcePrefixBytes:])
	if err != nil {
		t.Fatal(err)
	}
	var plaintext [selectionSIDBPayloadBytes]byte
	for index := range plaintext {
		plaintext[index] = byte(index*41 + 17)
	}
	registry, err := parseDecodedSIDBRegistry(plaintext[:])
	if err != nil {
		t.Fatal(err)
	}
	const recordIndex, windowOffset = 7, 3
	precursor, err := deriveCurrentSINFPrecursorFromRecordWindow(
		registry.Records[recordIndex], windowOffset, currentClassEnrollmentParameters().classKey(),
	)
	if err != nil {
		t.Fatal(err)
	}
	sinf := syntheticCurrentSINFForPrecursorRights(t, precursor, 104)
	fillSyntheticSelectionSIDBContainer(t, container, layout.PayloadOffset, plaintext[:], keyCandidates[8].Key)
	registries, _, gotIterations, err := deriveCurrentSelectionRegistryCandidates(
		map[string][]byte{"/private/var/db/fpsd/SC Info/SC Info.sidb": container}, hardware,
	)
	if err != nil {
		t.Fatal(err)
	}
	defer clearDecodedSIDBRegistryCandidates(registries)
	if gotIterations != 306 {
		t.Fatalf("structural iteration count = %d", gotIterations)
	}
	pack, err := currentEnrollmentPack(map[string][]byte{"App.sinf": sinf}, registries, hardware)
	if err != nil {
		t.Fatal(err)
	}
	defer clear(pack)
	if len(pack) != typedEnrollmentModelBytes || binary.LittleEndian.Uint32(pack[48:52]) != 0x12345678 || binary.LittleEndian.Uint32(pack[52:56]) != 2 {
		t.Fatal("structural SIDB/SINF model envelope is invalid")
	}
	wantKey := currentSINFModelKeyFromPrecursor(precursor)
	if !bytes.Equal(pack[56:72], wantKey[:]) {
		t.Fatal("structural SIDB/SINF model key differs")
	}
}

func TestDeriveRecordAlignedCurrentSelectionRegistryCandidate(t *testing.T) {
	const codecBytes = 24
	plaintext := make([]byte, 2*decodedSIDBRecordBytes)
	for index := range plaintext {
		plaintext[index] = byte(index*37 + 13)
	}
	container, hardware := syntheticStructuralCurrentOuterInputs(codecBytes, len(plaintext))
	layout, err := currentSIDBLayoutForContainer(container)
	if err != nil {
		t.Fatal(err)
	}
	if layout.CodecBytes != codecBytes || layout.PayloadBytes != len(plaintext) || layout.OuterIterations != 262 {
		t.Fatalf("record-aligned structural layout = %+v", layout)
	}
	source, _, err := deriveCurrentSelectionSource(container, hardware)
	if err != nil {
		t.Fatal(err)
	}
	keys, err := enrollmentKeyCandidates(source[selectionSourcePrefixBytes:])
	if err != nil {
		t.Fatal(err)
	}
	fillSyntheticSelectionSIDBContainer(t, container, layout.PayloadOffset, plaintext, keys[8].Key)
	candidates, _, _, err := deriveCurrentSelectionRegistryCandidates(
		map[string][]byte{"/private/var/db/fpsd/SC Info/SC Info.sidb": container}, hardware,
	)
	if err != nil {
		t.Fatal(err)
	}
	defer clearDecodedSIDBRegistryCandidates(candidates)
	var selected []decodedSIDBRegistryCandidate
	for _, candidate := range candidates {
		if bytes.Equal(candidate.Data, plaintext) {
			selected = append(selected, candidate)
		}
	}
	if len(selected) != 1 || selected[0].Registry.FooterBytes != decodedSIDBNoSeparateFooterBytes {
		t.Fatalf("record-aligned decoded matches = %d", len(selected))
	}
	const recordIndex, windowOffset = 1, 2
	precursor, err := deriveCurrentSINFPrecursorFromRecordWindow(
		selected[0].Registry.Records[recordIndex], windowOffset, currentClassEnrollmentParameters().classKey(),
	)
	if err != nil {
		t.Fatal(err)
	}
	sinf := syntheticCurrentSINFForPrecursor(t, precursor)
	authenticated, err := selectAuthenticatedEnrollmentClassPrecursor(sinf, selected, hardware, 0x12345678, 2)
	if err != nil {
		t.Fatal(err)
	}
	if authenticated.RecordIndex != recordIndex || authenticated.WindowOffset != windowOffset || authenticated.Precursor != precursor {
		t.Fatalf("record-aligned authenticated precursor = %+v", authenticated)
	}
}

func TestCurrentSIDBContainerFilesReportsRejectedGeometry(t *testing.T) {
	data := make([]byte, 96)
	binary.BigEndian.PutUint32(data[:4], currentSIDBContainerVersion)
	_, err := currentSIDBContainerFiles(map[string][]byte{"state/fpsd-sc-info/SC Info.sidb": data})
	if err == nil || !strings.Contains(err.Error(), "96 bytes") || !strings.Contains(err.Error(), "unsupported or ambiguous size") {
		t.Fatalf("rejected SIDB diagnostic = %v", err)
	}
}
