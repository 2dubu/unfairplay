package recovery

import (
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

func syntheticCurrentSINFForPrecursor(t *testing.T, profile *semanticBuildProfile, precursor [16]byte) []byte {
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
	for index, value := range logical {
		position := index % aes.BlockSize
		encoded := -1
		for candidate, decoded := range profile.currentSINFInput[position*256 : (position+1)*256] {
			if decoded == value {
				encoded = candidate
				break
			}
		}
		if encoded < 0 {
			t.Fatalf("no current-SINF encoding for position %d value %#02x", position, value)
		}
		protected[index] = byte(encoded) ^ profile.currentSINFWhite[position]
	}
	fields := append(syntheticSINFAtom("user", []byte{0x12, 0x34, 0x56, 0x78}), syntheticSINFAtom("key ", []byte{0, 0, 0, 2})...)
	fields = append(fields, syntheticSINFAtom("priv", protected)...)
	return syntheticSINFAtom("sinf", syntheticSINFAtom("schi", fields))
}

func TestSelectAuthenticatedClassPrecursor(t *testing.T) {
	profile := embedded26A428ExecutionProfile().semantic
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
	sinf := syntheticCurrentSINFForPrecursor(t, profile, precursor)
	accepted, err := validateCurrentSINFPrecursor(sinf, profile, precursor)
	if err != nil || !accepted {
		t.Fatalf("correct precursor rejected: accepted=%v err=%v", accepted, err)
	}
	wrong := precursor
	wrong[0] ^= 1
	accepted, err = validateCurrentSINFPrecursor(sinf, profile, wrong)
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
	selected, err := selectAuthenticatedClassPrecursor(sinf, profile, registries)
	if err != nil {
		t.Fatal(err)
	}
	if selected.RegistryIndex != 0 || selected.RecordIndex != 0 || selected.WindowOffset != selectedWindow || selected.Precursor != precursor {
		t.Fatalf("wrong authenticated candidate: %+v", selected)
	}
}
