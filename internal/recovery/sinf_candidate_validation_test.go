package recovery

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"encoding/binary"
	"errors"
	"strings"
	"testing"
)

func syntheticSINFAtom(name string, payload []byte) []byte {
	result := make([]byte, 8+len(payload))
	binary.BigEndian.PutUint32(result[:4], uint32(len(result)))
	copy(result[4:8], name)
	copy(result[8:], payload)
	return result
}

func TestCurrentSINFLogicalSourceIsIdentity(t *testing.T) {
	var precursor [16]byte
	for index := range precursor {
		precursor[index] = byte(index*13 + 7)
	}
	sinf := syntheticCurrentSINFForPrecursor(t, precursor)
	source, err := currentSINFPrivateSource(sinf)
	if err != nil {
		t.Fatal(err)
	}
	logical, _, _, err := currentSINFLogicalSource(sinf)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(logical, source[:len(logical)]) {
		t.Fatal("current-SINF logical source is not the authenticated private source")
	}
}

func syntheticCurrentSINFForPrecursor(t *testing.T, precursor [16]byte) []byte {
	return syntheticCurrentSINFForPrecursorRights(t, precursor, 88)
}

func syntheticCurrentSINFForPrecursorRights(t *testing.T, precursor [16]byte, rightsAtomBytes int) []byte {
	t.Helper()
	var plaintext []byte
	for index, atom := range []struct {
		name string
		size int
	}{{"key ", 24}, {"iviv", 24}, {"tmpr", 24}, {"righ", rightsAtomBytes}, {"name", 264}} {
		payload := make([]byte, atom.size-8)
		for offset := range payload {
			payload[offset] = byte(index*47 + offset*19 + 3)
		}
		plaintext = append(plaintext, syntheticSINFAtom(atom.name, payload)...)
	}
	if len(plaintext)%8 != 0 {
		t.Fatalf("synthetic plaintext has invalid geometry: %d bytes", len(plaintext))
	}

	key := currentSINFModelKeyFromPrecursor(precursor)
	block, err := aes.NewCipher(key[:])
	clear(key[:])
	if err != nil {
		t.Fatal(err)
	}
	iv := [aes.BlockSize]byte{0x53, 0x49, 0x4e, 0x46, 0x2d, 0x49, 0x56, 0xa5, 0x5a, 0x13, 0x37, 0x24, 0x68, 0x9c, 0xde, 0xf0}
	trailerBytes := len(plaintext) % aes.BlockSize
	ciphertextBytes := len(plaintext) - trailerBytes
	logical := make([]byte, aes.BlockSize+ciphertextBytes)
	copy(logical, iv[:])
	cipher.NewCBCEncrypter(block, iv[:]).CryptBlocks(logical[aes.BlockSize:], plaintext[:ciphertextBytes])
	protected := make([]byte, len(logical)+trailerBytes)
	copy(protected, logical)
	fields := append(syntheticSINFAtom("user", []byte{0x12, 0x34, 0x56, 0x78}), syntheticSINFAtom("key ", []byte{0, 0, 0, 2})...)
	fields = append(fields, syntheticSINFAtom("priv", protected)...)
	return syntheticSINFAtom("sinf", syntheticSINFAtom("schi", fields))
}

func TestValidateCurrentSINFPrecursorExtendedRights(t *testing.T) {
	var precursor [16]byte
	for index := range precursor {
		precursor[index] = byte(index*7 + 5)
	}
	sinf := syntheticCurrentSINFForPrecursorRights(t, precursor, 104)
	accepted, err := validateCurrentSINFPrecursor(sinf, precursor)
	if err != nil || !accepted {
		t.Fatalf("extended-rights precursor rejected: accepted=%v err=%v", accepted, err)
	}
}

func TestSelectAuthenticatedClassPrecursor(t *testing.T) {
	data := make([]byte, decodedSIDBRecordBytes+decodedSIDBFooterBytes)
	for index := range data {
		data[index] = byte(index*31 + 7)
	}
	registry, err := parseLegacyDecodedSIDBRegistry(data)
	if err != nil {
		t.Fatal(err)
	}
	const selectedWindow = 3
	precursor, err := deriveCurrentSINFPrecursorFromRecordWindow(registry.Records[0], selectedWindow, currentClassEnrollmentParameters().classKey())
	if err != nil {
		t.Fatal(err)
	}
	sinf := syntheticCurrentSINFForPrecursor(t, precursor)
	accepted, err := validateCurrentSINFPrecursor(sinf, precursor)
	if err != nil || !accepted {
		t.Fatalf("correct precursor rejected: accepted=%v err=%v", accepted, err)
	}
	wrong := precursor
	wrong[0] ^= 1
	accepted, err = validateCurrentSINFPrecursor(sinf, wrong)
	if err != nil || accepted {
		t.Fatalf("wrong precursor accepted: accepted=%v err=%v", accepted, err)
	}

	falseData := append([]byte(nil), data...)
	falseData[decodedSIDBMetadataBytes+5] ^= 0x80
	falseRegistry, err := parseLegacyDecodedSIDBRegistry(falseData)
	if err != nil {
		t.Fatal(err)
	}
	registries := []decodedSIDBRegistryCandidate{
		{Registry: registry, Data: data},
		{Registry: falseRegistry, Data: falseData},
	}
	selected, err := selectAuthenticatedEnrollmentClassPrecursor(sinf, registries, []byte{0, 1, 2, 3, 4, 5}, 0x12345678, 2)
	if err != nil {
		t.Fatal(err)
	}
	if selected.RegistryIndex != 0 || selected.RecordIndex != 0 || selected.WindowOffset != selectedWindow || selected.Precursor != precursor {
		t.Fatalf("wrong authenticated candidate: %+v", selected)
	}
}

func TestSelectAuthenticatedFlatClassPrecursorReportsProgressAndCancels(t *testing.T) {
	data := make([]byte, decodedSIDBRecordBytes+decodedSIDBFooterBytes)
	for index := range data {
		data[index] = byte(index*23 + 5)
	}
	registry, err := parseLegacyDecodedSIDBRegistry(data)
	if err != nil {
		t.Fatal(err)
	}
	precursor, err := deriveCurrentSINFPrecursorFromRecordWindow(
		registry.Records[0], 1, currentClassEnrollmentParameters().classKey(),
	)
	if err != nil {
		t.Fatal(err)
	}
	sinf := syntheticCurrentSINFForPrecursor(t, precursor)
	registries := []decodedSIDBRegistryCandidate{{Registry: registry, Data: data}}
	var events []Progress
	hardware := []byte{0, 1, 2, 3, 4, 5}
	selected, err := selectAuthenticatedFlatClassPrecursorWithContext(
		context.Background(), sinf, registries, hardware, 0x12345678, 2,
		func(event Progress) { events = append(events, event) },
	)
	if err != nil {
		t.Fatal(err)
	}
	if selected.Precursor != precursor || len(events) < 2 {
		t.Fatalf("selected/progress = %x/%d", selected.Precursor, len(events))
	}
	previous := uint64(0)
	for _, event := range events {
		if event.Total == 0 || event.Completed < previous || event.Completed > event.Total ||
			event.RegistryCandidates != 1 || event.RecordCandidates != 1 || event.WindowCandidates != maximumClassWindowOffset+1 {
			t.Fatalf("invalid progress event: %+v", event)
		}
		previous = event.Completed
	}
	last := events[len(events)-1]
	if last.Phase != "search-sinf-complete" || last.Completed != last.Total {
		t.Fatalf("incomplete final progress: %+v", last)
	}

	canceled, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := selectAuthenticatedFlatClassPrecursorWithContext(canceled, sinf, registries, hardware, 0x12345678, 2, nil); !errors.Is(err, context.Canceled) {
		t.Fatalf("canceled class search = %v", err)
	}
}

func TestSelectAuthenticatedFlatClassPrecursorRejectsExcessiveWork(t *testing.T) {
	var precursor [16]byte
	sinf := syntheticCurrentSINFForPrecursor(t, precursor)
	registries := make([]decodedSIDBRegistryCandidate, 128)
	for index := range registries {
		registries[index].Registry.Records = make([]decodedSIDBRecord, 512)
	}
	_, err := selectAuthenticatedFlatClassPrecursorWithContext(context.Background(), sinf, registries, []byte{0, 1, 2, 3, 4, 5}, 0x12345678, 2, nil)
	if err == nil || !strings.Contains(err.Error(), "unsupported size") {
		t.Fatalf("excessive class search = %v", err)
	}
}

func TestSelectAuthenticatedWorkspaceClassPrecursor(t *testing.T) {
	data := make([]byte, 2*decodedSIDBRecordBytes+decodedSIDBFooterBytes)
	for index := range data {
		data[index] = byte(index*43 + 17)
	}
	registry, err := parseLegacyDecodedSIDBRegistry(data)
	if err != nil {
		t.Fatal(err)
	}
	const selectedRecord = 1
	const selectedWindow = 4
	workspaceValue := classWorkspaceOutputFromSeed(0)
	parameters := classEnrollmentParameters{ClassValue: workspaceValue}
	precursor, err := deriveCurrentSINFPrecursorFromRecordWindow(registry.Records[selectedRecord], selectedWindow, parameters.classKey())
	if err != nil {
		t.Fatal(err)
	}
	sinf := syntheticCurrentSINFForPrecursor(t, precursor)
	selected, err := selectAuthenticatedEnrollmentClassPrecursor(sinf, []decodedSIDBRegistryCandidate{{Registry: registry, Data: data}}, []byte{0, 1, 2, 3, 4, 5}, 0x12345678, 2)
	if err != nil {
		t.Fatal(err)
	}
	if selected.RegistryIndex != 0 || selected.RecordIndex != selectedRecord || selected.WindowOffset != selectedWindow || selected.Precursor != precursor {
		t.Fatalf("wrong workspace candidate: %+v", selected)
	}
}

func TestSelectAuthenticatedLegacyClassPrecursor(t *testing.T) {
	data := make([]byte, 14*decodedSIDBRecordBytes+decodedSIDBOlderFooterBytes)
	for index := range data {
		data[index] = byte(index*37 + 11)
	}
	registry, err := parseDecodedSIDBRegistry(data)
	if err != nil {
		t.Fatal(err)
	}
	hardware := []byte{1, 2, 3, 4, 5, 6}
	const user = uint32(0x12345678)
	const keyID = uint32(5)
	digest, err := legacyAccountDigest(hardware, user, keyID)
	if err != nil {
		t.Fatal(err)
	}
	parameters := classEnrollmentParameters{ClassValue: legacyClassWorkspaceValue(digest)}
	clear(digest[:])
	const selectedRecord = 7
	const selectedPhase = 8
	valu, err := deriveLegacyClassValu(registry.Records[selectedRecord], selectedPhase)
	if err != nil {
		t.Fatal(err)
	}
	precursor, err := deriveCurrentSINFPrecursorFromClass(parameters.classKey(), valu)
	if err != nil {
		t.Fatal(err)
	}
	sinf := syntheticCurrentSINFForPrecursor(t, precursor)
	rootStart, rootEnd, err := enrollmentAtom(sinf, 0, len(sinf), "sinf")
	if err != nil {
		t.Fatal(err)
	}
	schiStart, schiEnd, err := enrollmentAtom(sinf, rootStart, rootEnd, "schi")
	if err != nil {
		t.Fatal(err)
	}
	keyStart, keyEnd, err := enrollmentAtom(sinf, schiStart, schiEnd, "key ")
	if err != nil || keyEnd-keyStart != 4 {
		t.Fatalf("synthetic key atom: %v", err)
	}
	binary.BigEndian.PutUint32(sinf[keyStart:keyEnd], keyID)
	selected, err := selectAuthenticatedEnrollmentClassPrecursor(sinf, []decodedSIDBRegistryCandidate{{Registry: registry, Data: data}}, hardware, user, keyID)
	if err != nil {
		t.Fatal(err)
	}
	if selected.RegistryIndex != 0 || selected.RecordIndex != selectedRecord || selected.WindowOffset != selectedPhase || selected.Precursor != precursor {
		t.Fatalf("wrong legacy candidate: %+v", selected)
	}
}

func TestSelectAuthenticatedDerivedClassPrecursor(t *testing.T) {
	hardware := []byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55}
	data := syntheticSIDBRecordTree(1, 13, 5)
	registry, err := parseDecodedSIDBRegistry(data)
	if err != nil {
		t.Fatal(err)
	}
	key, err := accountClassKey(hardware, 0x12345678, 2)
	if err != nil {
		t.Fatal(err)
	}
	const selectedWindow = 2
	precursor, err := deriveCurrentSINFPrecursorFromRecordWindow(registry.Records[0], selectedWindow, key)
	if err != nil {
		t.Fatal(err)
	}
	sinf := syntheticCurrentSINFForPrecursor(t, precursor)
	selected, err := selectAuthenticatedClassPrecursor(sinf, []decodedSIDBRegistryCandidate{{Registry: registry, Data: data}}, hardware)
	if err != nil {
		t.Fatal(err)
	}
	if selected.RegistryIndex != 0 || selected.RecordIndex != 0 || selected.WindowOffset != selectedWindow || selected.Precursor != precursor {
		t.Fatal("derived class path selected the wrong synthetic candidate")
	}
	wrongHardware := append([]byte(nil), hardware...)
	wrongHardware[0] ^= 1
	if _, err := selectAuthenticatedClassPrecursor(sinf, []decodedSIDBRegistryCandidate{{Registry: registry, Data: data}}, wrongHardware); err == nil {
		t.Fatal("accepted SINF with wrong hardware")
	}
}

func TestEnrollmentClassSelectionAuthenticatesMixedRegistryCandidates(t *testing.T) {
	hardware := []byte{0, 1, 2, 3, 4, 5}
	data := syntheticSIDBRecordTree(1, 13, 5)
	tree, err := parseDecodedSIDBTree(data)
	if err != nil || len(tree) == 0 {
		t.Fatal("invalid synthetic tree")
	}
	registry, err := parseDecodedSIDBRegistry(data)
	if err != nil {
		t.Fatal(err)
	}
	key, err := accountClassKey(hardware, 0x12345678, 2)
	if err != nil {
		t.Fatal(err)
	}
	precursor, err := deriveCurrentSINFPrecursorFromRecordWindow(registry.Records[0], 2, key)
	if err != nil {
		t.Fatal(err)
	}
	sinf := syntheticCurrentSINFForPrecursor(t, precursor)
	flat, err := parseLegacyDecodedSIDBRegistry(make([]byte, 2*decodedSIDBRecordBytes+decodedSIDBOlderFooterBytes))
	if err != nil {
		t.Fatal(err)
	}
	candidates := []decodedSIDBRegistryCandidate{{Registry: flat}, {Registry: registry}}
	selected, err := selectAuthenticatedEnrollmentClassPrecursor(sinf, candidates, hardware, 0x12345678, 2)
	if err != nil || selected.RegistryIndex != 1 || selected.Precursor != precursor {
		t.Fatalf("mixed candidate authentication failed: %v", err)
	}
	candidates = append(candidates, decodedSIDBRegistryCandidate{Registry: registry})
	if _, err := selectAuthenticatedEnrollmentClassPrecursor(sinf, candidates, hardware, 0x12345678, 2); err == nil {
		t.Fatal("accepted multiple authenticated candidates")
	}
}

func TestTreeSelectionDoesNotTryFlatClassFamilies(t *testing.T) {
	hardware := []byte{0, 1, 2, 3, 4, 5}
	data := syntheticSIDBRecordTree(1, 13, 5)
	registry, err := parseDecodedSIDBRegistry(data)
	if err != nil || !registry.Tree {
		t.Fatal("synthetic tree was not parsed")
	}
	precursor, err := deriveCurrentSINFPrecursorFromRecordWindow(registry.Records[0], 2, currentClassEnrollmentParameters().classKey())
	if err != nil {
		t.Fatal(err)
	}
	sinf := syntheticCurrentSINFForPrecursor(t, precursor)
	if _, err := selectAuthenticatedEnrollmentClassPrecursor(sinf, []decodedSIDBRegistryCandidate{{Registry: registry}}, hardware, 0x12345678, 2); !errors.Is(err, errNoAuthenticatedClassPrecursor) ||
		!strings.Contains(err.Error(), "registries=1 tree=1 flat=0 older=0 records=1 windows=5") {
		t.Fatalf("tree flat-family rejection diagnostic = %v", err)
	}
}

func TestEnrollmentClassSelectionRejectsMixedFlatTreeAmbiguity(t *testing.T) {
	hardware := []byte{0, 1, 2, 3, 4, 5}
	data := syntheticSIDBRecordTree(1, 13, 5)
	registry, err := parseDecodedSIDBRegistry(data)
	if err != nil || !registry.Tree {
		t.Fatal("synthetic tree was not parsed")
	}
	key, err := accountClassKey(hardware, 0x12345678, 2)
	if err != nil {
		t.Fatal(err)
	}
	precursor, err := deriveCurrentSINFPrecursorFromRecordWindow(registry.Records[0], 2, key)
	if err != nil {
		t.Fatal(err)
	}
	sinf := syntheticCurrentSINFForPrecursor(t, precursor)
	flat := registry
	flat.Tree = false
	flat.FooterBytes = decodedSIDBFooterBytes
	candidates := []decodedSIDBRegistryCandidate{{Registry: flat}, {Registry: registry}}
	if _, err := selectAuthenticatedEnrollmentClassPrecursor(sinf, candidates, hardware, 0x12345678, 2); err == nil {
		t.Fatal("accepted distinct authenticated flat and tree candidates")
	}
}

func TestNativeTreeSuppressesOverlappingFlatInterpretation(t *testing.T) {
	data := syntheticSIDBTreeBuffer(224)
	for i := range data {
		data[i] = byte(i*37 + 11)
	}
	putSyntheticSIDBNode(data, 0, len(data), 3)
	putSyntheticSIDBNode(data, 20, 60, 0)
	putSyntheticSIDBNode(data, 80, 36, 0)
	putSyntheticSIDBNode(data, 116, len(data)-116, 0)
	views := decodedSIDBRegistryInterpretations(data)
	if len(views) != 1 || !views[0].Tree || len(views[0].Records) != 1 {
		t.Fatal("native tree was not the unique interpretation")
	}
}
