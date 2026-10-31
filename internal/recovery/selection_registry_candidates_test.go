package recovery

import (
	"bytes"
	"crypto/aes"
	"crypto/md5"
	"encoding/binary"
	"strings"
	"testing"
)

// sealSyntheticSelectionContainer independently constructs a public synthetic
// footer after the encrypted body has been written.
func sealSyntheticSelectionContainer(t *testing.T, container, source []byte) {
	t.Helper()
	var key [16]byte
	for i := range key {
		key[i] = source[i] ^ 0x1b
	}
	block, err := aes.NewCipher(key[:])
	if err != nil {
		t.Fatal(err)
	}
	digest := md5.Sum(container[:len(container)-16])
	block.Encrypt(container[len(container)-16:], digest[:])
}

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
	selectedKey, err := deriveSelectedEnrollmentKey(source[selectionSourcePrefixBytes:], layout.SelectedSlot)
	if err != nil {
		t.Fatal(err)
	}
	var plaintext [selectionSIDBPayloadBytes]byte
	copy(plaintext[:], syntheticSIDBTreeBuffer(len(plaintext)))
	root := layout.RootCursor - 2 // Synthetic writer uses the older two-byte-relative coordinate.
	rootLength := len(plaintext) - layout.RootCursor
	putSyntheticSIDBNode(plaintext[:], root, rootLength, 2)
	putSyntheticSIDBNode(plaintext[:], root+20, 36, 0)
	putSyntheticSIDBNode(plaintext[:], root+56, rootLength-56, 0)
	fillSyntheticSelectionSIDBContainer(t, container, layout.PayloadOffset, plaintext[:], selectedKey)
	sealSyntheticSelectionContainer(t, container, source)
	candidates, _, gotIterations, err := deriveCurrentSelectionRegistryCandidates(
		map[string][]byte{"/private/var/db/fpsd/SC Info/SC Info.sidb": container}, hardware,
	)
	if err != nil {
		t.Fatal(err)
	}
	if gotIterations != layout.OuterIterations {
		t.Fatalf("extended selection iterations=%d", gotIterations)
	}
	matches := 0
	for _, candidate := range candidates {
		if bytes.Equal(candidate.Data, plaintext[:]) {
			if !candidate.Registry.Tree {
				t.Fatal("unexpected non-tree candidate")
			}
			matches++
		}
	}
	clearDecodedSIDBRegistryCandidates(candidates)
	// The same byte length fits the former flat geometry; only the tree view
	// remains because the historical interpretation is not an established format.
	if matches != 1 {
		t.Fatalf("expected one tree interpretation, got %d", matches)
	}
	for _, offset := range []int{layout.PayloadOffset, len(container) - 1} {
		corrupt := append([]byte(nil), container...)
		corrupt[offset] ^= 1
		if got, _, _, err := deriveCurrentSelectionRegistryCandidates(map[string][]byte{"state/SC Info.sidb": corrupt}, hardware); err == nil || !strings.Contains(err.Error(), "SIDB container authentication failed") {
			clearDecodedSIDBRegistryCandidates(got)
			t.Fatalf("container mutation at offset %d: %v", offset, err)
		}
	}
	if layout.SelectorModulus > 1 {
		wrongKey, err := deriveSelectedEnrollmentKey(source[selectionSourcePrefixBytes:], (layout.SelectedSlot+1)%layout.SelectorModulus)
		if err != nil {
			t.Fatal(err)
		}
		wrongSlot := append([]byte(nil), container...)
		fillSyntheticSelectionSIDBContainer(t, wrongSlot, layout.PayloadOffset, plaintext[:], wrongKey)
		sealSyntheticSelectionContainer(t, wrongSlot, source)
		if got, _, _, err := deriveCurrentSelectionRegistryCandidates(map[string][]byte{"state/SC Info.sidb": wrongSlot}, hardware); err == nil {
			clearDecodedSIDBRegistryCandidates(got)
			t.Fatal("accepted body encrypted under a different slot")
		}
	}
}

func TestDeriveStructuralCurrentSelectionRegistryCandidates(t *testing.T) {
	container, hardware, layout := syntheticNativeSIDB(t, selectionSIDBPayloadBytes)
	for index := layout.PayloadOffset; index < layout.PayloadOffset+layout.PayloadBytes; index++ {
		container[index] = 0
	}
	source, _, iterations, err := deriveNativeSelectionSource(container, hardware)
	if err != nil || iterations != layout.OuterIterations {
		t.Fatalf("derive structural selection source = %d/%v", iterations, err)
	}
	selectedKey, err := deriveSelectedEnrollmentKey(source[selectionSourcePrefixBytes:], layout.SelectedSlot)
	if err != nil {
		t.Fatal(err)
	}
	var plaintext [selectionSIDBPayloadBytes]byte
	copy(plaintext[:], syntheticSIDBTreeBuffer(len(plaintext)))
	root := layout.RootCursor - 2 // Synthetic writer uses the older two-byte-relative coordinate.
	rootLength := len(plaintext) - layout.RootCursor
	putSyntheticSIDBNode(plaintext[:], root, rootLength, 2)
	putSyntheticSIDBNode(plaintext[:], root+20, 36, 0)
	putSyntheticSIDBNode(plaintext[:], root+56, rootLength-56, 0)
	fillSyntheticSelectionSIDBContainer(t, container, layout.PayloadOffset, plaintext[:], selectedKey)
	sealSyntheticSelectionContainer(t, container, source)
	candidates, _, gotIterations, err := deriveCurrentSelectionRegistryCandidates(
		map[string][]byte{"/private/var/db/fpsd/SC Info/SC Info.sidb": container}, hardware,
	)
	if err != nil {
		t.Fatal(err)
	}
	if gotIterations != layout.OuterIterations {
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
	container, hardware, layout := syntheticNativeSIDB(t, selectionSIDBPayloadBytes)
	source, _, iterations, err := deriveNativeSelectionSource(container, hardware)
	if err != nil || iterations != layout.OuterIterations {
		t.Fatalf("derive structural source = %d/%v", iterations, err)
	}
	selectedKey, err := deriveSelectedEnrollmentKey(source[selectionSourcePrefixBytes:], layout.SelectedSlot)
	if err != nil {
		t.Fatal(err)
	}
	plaintext, _, sinf, user, keyID, precursor := syntheticProductionRoute(t, hardware, selectionSIDBPayloadBytes)
	fillSyntheticSelectionSIDBContainer(t, container, layout.PayloadOffset, plaintext, selectedKey)
	sealSyntheticSelectionContainer(t, container, source)
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
	if len(pack) != typedEnrollmentModelBytes || binary.LittleEndian.Uint32(pack[48:52]) != user || binary.LittleEndian.Uint32(pack[52:56]) != keyID {
		t.Fatal("structural SIDB/SINF model envelope is invalid")
	}
	wantKey := currentSINFModelKeyFromPrecursor(precursor)
	if !bytes.Equal(pack[56:72], wantKey[:]) {
		t.Fatal("structural SIDB/SINF model key differs")
	}
}

func TestDeriveRejectsSyntheticFlatSelectionRegistry(t *testing.T) {
	plaintext := make([]byte, 2*decodedSIDBRecordBytes)
	for i := range plaintext {
		plaintext[i] = byte(i*37 + 13)
	}
	container, hardware, layout := syntheticNativeSIDB(t, len(plaintext))
	source, _, _, err := deriveNativeSelectionSource(container, hardware)
	if err != nil {
		t.Fatal(err)
	}
	selectedKey, err := deriveSelectedEnrollmentKey(source[selectionSourcePrefixBytes:], layout.SelectedSlot)
	if err != nil {
		t.Fatal(err)
	}
	fillSyntheticSelectionSIDBContainer(t, container, layout.PayloadOffset, plaintext, selectedKey)
	sealSyntheticSelectionContainer(t, container, source)
	candidates, _, _, err := deriveCurrentSelectionRegistryCandidates(map[string][]byte{"state/SC Info.sidb": container}, hardware)
	if err == nil {
		clearDecodedSIDBRegistryCandidates(candidates)
		t.Fatal("accepted synthetic flat-only body")
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

func TestCurrentSIDBContainerFilesRequireOneDistinctContainer(t *testing.T) {
	first := make([]byte, 16)
	binary.BigEndian.PutUint32(first, currentSIDBContainerVersion)
	second := append([]byte(nil), first...)
	second[4] = 1
	otherVersion := append([]byte(nil), first...)
	otherVersion[3] ^= 1
	aliases := map[string][]byte{"z/SC Info.sidb": first, "a/SC Info.sidb": append([]byte(nil), first...), "old/SC Info.sidb": otherVersion}
	got, err := currentSIDBContainerFiles(aliases)
	if err != nil || len(got) != 1 || got[0].name != "a/SC Info.sidb" || !bytes.Equal(got[0].data, first) {
		t.Fatalf("identical aliases: %v, count %d", err, len(got))
	}
	aliases["b/SC Info.sidb"] = second
	if _, err := currentSIDBContainerFiles(aliases); err == nil || !strings.Contains(err.Error(), "one distinct") {
		t.Fatalf("distinct containers: %v", err)
	}
	if _, err := currentSIDBContainerFiles(map[string][]byte{"old/SC Info.sidb": otherVersion}); err == nil {
		t.Fatal("accepted an unknown version as current")
	}
}

func TestSelectionRegistryDiagnosticsAreNonSecretAndComplete(t *testing.T) {
	diagnostics := selectionRegistryDiagnostics{
		Containers: 1, ContainerBytes: 11724, CodecBytes: 24, PayloadBytes: 11680,
		SourceBytes: 320, OuterIterations: 262, SourceLayouts: 1,
		KeyCandidates: 1, DecodedBodies: 1, TreeCandidates: 0,
	}
	want := "containers=1 container_bytes=11724 codec_bytes=24 payload_bytes=11680 source_bytes=320 iterations=262 source_layouts=1 key_candidates=1 decoded_bodies=1 tree_candidates=0"
	if got := diagnostics.String(); got != want {
		t.Fatalf("selection diagnostics = %q", got)
	}
}

func TestSelectionMappedSlot14UsesCBCRecurrence(t *testing.T) {
	container, hardware, layout := syntheticNativeSIDBMatching(t, selectionSIDBPayloadBytes, func(layout currentSIDBLayout) bool {
		return layout.SelectedSlot == 14
	})
	source, _, _, err := deriveNativeSelectionSource(container, hardware)
	if err != nil {
		t.Fatal(err)
	}
	selectedKey, err := deriveSelectedEnrollmentKey(source[selectionSourcePrefixBytes:], layout.SelectedSlot)
	if err != nil {
		t.Fatal(err)
	}
	plaintext := syntheticSIDBTreeBuffer(layout.PayloadBytes)
	root := layout.RootCursor - 2
	rootLength := len(plaintext) - layout.RootCursor
	putSyntheticSIDBNode(plaintext, root, rootLength, 2)
	putSyntheticSIDBNode(plaintext, root+20, 36, 0)
	putSyntheticSIDBNode(plaintext, root+56, rootLength-56, 0)
	fillSyntheticSelectionSIDBContainer(t, container, layout.PayloadOffset, plaintext, selectedKey)
	sealSyntheticSelectionContainer(t, container, source)
	candidates, _, _, err := deriveCurrentSelectionRegistryCandidates(map[string][]byte{"state/SC Info.sidb": container}, hardware)
	if err != nil {
		t.Fatal(err)
	}
	defer clearDecodedSIDBRegistryCandidates(candidates)
	if len(candidates) != 1 || !bytes.Equal(candidates[0].Data, plaintext) {
		t.Fatalf("selected slot %d did not produce one authenticated registry", layout.SelectedSlot)
	}
}
