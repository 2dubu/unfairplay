package recovery

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"encoding/binary"
	"errors"
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

func TestTreeAuthenticationProgressCancellationAndWrongSINF(t *testing.T) {
	hardware := []byte{0, 1, 2, 3, 4, 5}
	data, registry, sinf, user, keyID, precursor := syntheticProductionRoute(t, hardware, 0)
	registries := []decodedSIDBRegistryCandidate{{Registry: registry, Data: data}}
	var events []Progress
	selected, err := selectAuthenticatedEnrollmentClassPrecursorWithContext(context.Background(), sinf, registries, hardware, user, keyID, func(event Progress) { events = append(events, event) })
	if err != nil || selected.Precursor != precursor || selected.WindowOffset != 0 {
		t.Fatalf("tree selection: %v", err)
	}
	if len(events) < 2 {
		t.Fatal("missing progress")
	}
	var previous uint64
	for _, event := range events {
		if event.Total == 0 || event.Completed < previous || event.Completed > event.Total || event.RegistryCandidates != 1 || event.RecordCandidates != 1 || event.WindowCandidates != 1 {
			t.Fatalf("invalid progress: %+v", event)
		}
		previous = event.Completed
	}
	last := events[len(events)-1]
	if last.Phase != "search-sinf-complete" || last.Completed != last.Total {
		t.Fatalf("incomplete progress: %+v", last)
	}
	canceled, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := selectAuthenticatedEnrollmentClassPrecursorWithContext(canceled, sinf, registries, hardware, user, keyID, nil); !errors.Is(err, context.Canceled) {
		t.Fatalf("canceled: %v", err)
	}
	wrong := precursor
	wrong[0] ^= 1
	if accepted, err := validateCurrentSINFPrecursor(sinf, wrong); err != nil || accepted {
		t.Fatalf("wrong precursor accepted: %v", err)
	}
}

func TestEnrollmentClassSelectionRejectsUnsupportedGrammarAndAmbiguity(t *testing.T) {
	hardware := []byte{0, 1, 2, 3, 4, 5}
	data, registry, sinf, user, keyID, _ := syntheticProductionRoute(t, hardware, 0)
	unsupported := registry
	unsupported.Tree = false
	for _, candidates := range [][]decodedSIDBRegistryCandidate{
		{{Registry: unsupported, Data: data}},
		{{Registry: registry, Data: data}, {Registry: registry, Data: data}},
		{{Registry: registry, Data: data}, {Registry: unsupported, Data: data}},
	} {
		if _, err := selectAuthenticatedEnrollmentClassPrecursor(sinf, candidates, hardware, user, keyID); err == nil {
			t.Fatal("accepted unsupported or multiple registries")
		}
	}
}

func TestTreeSelectionDoesNotTryAlternativeClassFamilies(t *testing.T) {
	hardware := []byte{0, 1, 2, 3, 4, 5}
	data, registry, sinf, user, keyID, _ := syntheticProductionRoute(t, hardware, 0)
	wrong := append([]byte(nil), sinf...)
	priv, err := currentSINFPrivateSource(wrong)
	if err != nil {
		t.Fatal(err)
	}
	priv[32] ^= 1
	if _, err := selectAuthenticatedEnrollmentClassPrecursor(wrong, []decodedSIDBRegistryCandidate{{Registry: registry, Data: data}}, hardware, user, keyID); !errors.Is(err, errNoAuthenticatedClassPrecursor) {
		t.Fatalf("wrong SINF: %v", err)
	}
}
