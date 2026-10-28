package recovery

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/binary"
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
	t.Helper()
	var plaintext []byte
	for index, atom := range []struct {
		name string
		size int
	}{{"key ", 24}, {"iviv", 24}, {"tmpr", 24}, {"righ", 88}, {"name", 264}} {
		payload := make([]byte, atom.size-8)
		for offset := range payload {
			payload[offset] = byte(index*47 + offset*19 + 3)
		}
		plaintext = append(plaintext, syntheticSINFAtom(atom.name, payload)...)
	}
	if len(plaintext) != 424 {
		t.Fatalf("synthetic plaintext bytes = %d", len(plaintext))
	}

	key := currentSINFModelKeyFromPrecursor(precursor)
	block, err := aes.NewCipher(key[:])
	clear(key[:])
	if err != nil {
		t.Fatal(err)
	}
	iv := [aes.BlockSize]byte{0x53, 0x49, 0x4e, 0x46, 0x2d, 0x49, 0x56, 0xa5, 0x5a, 0x13, 0x37, 0x24, 0x68, 0x9c, 0xde, 0xf0}
	logical := make([]byte, aes.BlockSize+416)
	copy(logical, iv[:])
	cipher.NewCBCEncrypter(block, iv[:]).CryptBlocks(logical[aes.BlockSize:], plaintext[:416])
	protected := make([]byte, 440)
	copy(protected, logical)
	fields := append(syntheticSINFAtom("user", []byte{0x12, 0x34, 0x56, 0x78}), syntheticSINFAtom("key ", []byte{0, 0, 0, 2})...)
	fields = append(fields, syntheticSINFAtom("priv", protected)...)
	return syntheticSINFAtom("sinf", syntheticSINFAtom("schi", fields))
}

func TestSelectAuthenticatedClassPrecursor(t *testing.T) {
	data := make([]byte, decodedSIDBRecordBytes+decodedSIDBFooterBytes)
	for index := range data {
		data[index] = byte(index*31 + 7)
	}
	registry, err := parseDecodedSIDBRegistry(data)
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
	falseRegistry, err := parseDecodedSIDBRegistry(falseData)
	if err != nil {
		t.Fatal(err)
	}
	registries := []decodedSIDBRegistryCandidate{
		{Registry: registry, Data: data},
		{Registry: falseRegistry, Data: falseData},
	}
	selected, err := selectAuthenticatedClassPrecursor(sinf, registries)
	if err != nil {
		t.Fatal(err)
	}
	if selected.RegistryIndex != 0 || selected.RecordIndex != 0 || selected.WindowOffset != selectedWindow || selected.Precursor != precursor {
		t.Fatalf("wrong authenticated candidate: %+v", selected)
	}
}

func TestSelectAuthenticatedWorkspaceClassPrecursor(t *testing.T) {
	data := make([]byte, 2*decodedSIDBRecordBytes+decodedSIDBFooterBytes)
	for index := range data {
		data[index] = byte(index*43 + 17)
	}
	registry, err := parseDecodedSIDBRegistry(data)
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
	selected, err := selectAuthenticatedClassPrecursor(sinf, []decodedSIDBRegistryCandidate{{Registry: registry, Data: data}})
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
