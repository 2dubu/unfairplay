package recovery

import (
	"bytes"
	"crypto/aes"
	"encoding/binary"
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
	container, hardware, layout := syntheticNativeSIDB(t, selectionSIDBPayloadBytes)
	for index := layout.PayloadOffset; index < layout.PayloadOffset+selectionSIDBPayloadBytes; index++ {
		container[index] = 0
	}
	source, _, iterations, err := deriveNativeSelectionSource(container, hardware)
	if err != nil || iterations != layout.OuterIterations {
		t.Fatalf("derive extended selection source = %d/%v", iterations, err)
	}
	keyCandidates, err := enrollmentKeyCandidates(source[selectionSourcePrefixBytes:])
	if err != nil {
		t.Fatal(err)
	}
	const selectedSlot = 8
	var plaintext [selectionSIDBPayloadBytes]byte
	copy(plaintext[:], syntheticSIDBTreeBuffer(len(plaintext)))
	putSyntheticSIDBNode(plaintext[:], 0, len(plaintext), 2)
	putSyntheticSIDBNode(plaintext[:], 20, 36, 0)
	putSyntheticSIDBNode(plaintext[:], 56, len(plaintext)-56, 0)
	fillSyntheticSelectionSIDBContainer(t, container, layout.PayloadOffset, plaintext[:], keyCandidates[selectedSlot].Key)
	candidates, _, gotIterations, err := deriveCurrentSelectionRegistryCandidates(
		map[string][]byte{"/private/var/db/fpsd/SC Info/SC Info.sidb": container}, hardware,
	)
	if err != nil {
		t.Fatal(err)
	}
	if gotIterations != layout.OuterIterations {
		t.Fatalf("extended selection iterations=%d", gotIterations)
	}
	treeMatches, flatMatches := 0, 0
	for _, candidate := range candidates {
		if bytes.Equal(candidate.Data, plaintext[:]) {
			if candidate.Registry.Tree {
				treeMatches++
			} else {
				flatMatches++
			}
		}
	}
	clearDecodedSIDBRegistryCandidates(candidates)
	if treeMatches != 1 || flatMatches != 0 {
		t.Fatalf("expected only the native tree interpretation of the plaintext")
	}
}

func TestRejectsFlatOnlyCurrentSelectionRegistryCandidates(t *testing.T) {
	container, hardware, layout := syntheticNativeSIDB(t, selectionSIDBPayloadBytes)
	for index := layout.PayloadOffset; index < layout.PayloadOffset+layout.PayloadBytes; index++ {
		container[index] = 0
	}
	source, _, iterations, err := deriveNativeSelectionSource(container, hardware)
	if err != nil || iterations != layout.OuterIterations {
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
	candidates, _, _, err := deriveCurrentSelectionRegistryCandidates(
		map[string][]byte{"/private/var/db/fpsd/SC Info/SC Info.sidb": container}, hardware,
	)
	clearDecodedSIDBRegistryCandidates(candidates)
	if err == nil {
		t.Fatal("accepted a length-only flat decoded body")
	}
}

func TestStructuralSIDBAndExtendedSINFProduceEnrollmentModel(t *testing.T) {
	container, hardware, layout := syntheticNativeSIDB(t, selectionSIDBPayloadBytes)
	source, _, iterations, err := deriveNativeSelectionSource(container, hardware)
	if err != nil || iterations != layout.OuterIterations {
		t.Fatalf("derive structural source = %d/%v", iterations, err)
	}
	keyCandidates, err := enrollmentKeyCandidates(source[selectionSourcePrefixBytes:])
	if err != nil {
		t.Fatal(err)
	}
	var plaintext [selectionSIDBPayloadBytes]byte
	copy(plaintext[:], syntheticSIDBTreeBuffer(len(plaintext)))
	copy(plaintext[:], syntheticSIDBRecordTree(25, 41, 17))
	registry, err := parseDecodedSIDBRegistry(plaintext[:])
	if err != nil {
		t.Fatal(err)
	}
	const recordIndex, windowOffset = 7, 3
	classKey, err := accountClassKey(hardware, 0x12345678, 2)
	if err != nil {
		t.Fatal(err)
	}
	precursor, err := deriveCurrentSINFPrecursorFromRecordWindow(
		registry.Records[recordIndex], windowOffset, classKey,
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
	if gotIterations != layout.OuterIterations {
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

func TestRejectsRecordAlignedCurrentSelectionRegistryCandidate(t *testing.T) {
	plaintext := make([]byte, 2*decodedSIDBRecordBytes)
	for index := range plaintext {
		plaintext[index] = byte(index*37 + 13)
	}
	container, hardware, layout := syntheticNativeSIDB(t, len(plaintext))
	source, _, _, err := deriveNativeSelectionSource(container, hardware)
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
	clearDecodedSIDBRegistryCandidates(candidates)
	if err == nil {
		t.Fatal("accepted a record-aligned flat decoded body")
	}
}

func TestCurrentSIDBContainerFilesAcceptsVersionBeforeNativeGeometry(t *testing.T) {
	data := make([]byte, 96)
	binary.BigEndian.PutUint32(data[:4], currentSIDBContainerVersion)
	_, err := currentSIDBContainerFiles(map[string][]byte{"state/fpsd-sc-info/SC Info.sidb": data})
	if err != nil {
		t.Fatalf("valid version was rejected before native layout discovery: %v", err)
	}
}

func TestSelectionRegistryDiagnosticsAreNonSecretAndComplete(t *testing.T) {
	diagnostics := selectionRegistryDiagnostics{
		Containers: 1, ContainerBytes: 11724, CodecBytes: 24, PayloadBytes: 11680,
		SourceBytes: 256, OuterIterations: 262, SourceLayouts: 2,
		KeyCandidates: 26, DecodedBodies: 26, TreeCandidates: 0, FlatCandidates: 0,
	}
	want := "containers=1 container_bytes=11724 codec_bytes=24 payload_bytes=11680 source_bytes=256 iterations=262 source_layouts=2 key_candidates=26 decoded_bodies=26 tree_candidates=0 flat_candidates=0"
	if got := diagnostics.String(); got != want {
		t.Fatalf("selection diagnostics = %q", got)
	}
}
