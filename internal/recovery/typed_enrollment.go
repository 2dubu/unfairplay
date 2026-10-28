package recovery

import (
	"encoding/binary"
	"fmt"
	"strings"
)

const typedEnrollmentModelBytes = 8 + 4 + 4 + 32 + 4 + 4 + 16 + 16 + 16*256

// typedEnrollmentAlgorithmIdentity is the SHA-256 of
// "unfairplay/typed-current-enrollment/v1". It identifies the portable
// algorithm carried by the model, not any Apple executable or OS build.
var typedEnrollmentAlgorithmIdentity = [32]byte{
	0xa3, 0xaa, 0x88, 0x17, 0x49, 0x84, 0x86, 0x36,
	0xd8, 0x8a, 0x8b, 0x0e, 0xc0, 0x2a, 0xdb, 0x35,
	0xb7, 0x71, 0x01, 0xd0, 0xa0, 0x01, 0xde, 0xa6,
	0xa0, 0x67, 0x16, 0x96, 0xd8, 0x59, 0x54, 0x02,
}

func enrollmentAtom(blob []byte, start, limit int, name string) (int, int, error) {
	for cursor := start; cursor < limit; {
		if limit-cursor < 8 {
			return 0, 0, fmt.Errorf("SINF atom list is truncated")
		}
		size := int(binary.BigEndian.Uint32(blob[cursor : cursor+4]))
		if size < 8 || size > limit-cursor {
			return 0, 0, fmt.Errorf("SINF atom has invalid bounds")
		}
		if string(blob[cursor+4:cursor+8]) == name {
			return cursor + 8, cursor + size, nil
		}
		cursor += size
	}
	return 0, 0, fmt.Errorf("SINF has no %q atom", name)
}

func enrollmentSINFSelector(blob []byte) (uint32, uint32, error) {
	rootStart, rootEnd, err := enrollmentAtom(blob, 0, len(blob), "sinf")
	if err != nil || rootStart != 8 || rootEnd != len(blob) {
		return 0, 0, fmt.Errorf("SINF root is invalid")
	}
	schiStart, schiEnd, err := enrollmentAtom(blob, rootStart, rootEnd, "schi")
	if err != nil {
		return 0, 0, err
	}
	values := [2]uint32{}
	for index, name := range []string{"user", "key "} {
		start, end, atomErr := enrollmentAtom(blob, schiStart, schiEnd, name)
		if atomErr != nil || end-start != 4 {
			return 0, 0, fmt.Errorf("SINF selector atom %q is invalid", name)
		}
		values[index] = binary.BigEndian.Uint32(blob[start:end])
	}
	return values[0], values[1], nil
}

func currentEnrollmentPack(files map[string][]byte, registries []decodedSIDBRegistryCandidate, hardware []byte) ([]byte, error) {
	var sinf []byte
	for suffix, data := range files {
		if strings.HasSuffix(suffix, ".sinf") {
			if sinf != nil {
				return nil, fmt.Errorf("enrollment model requires exactly one installed SINF")
			}
			sinf = data
		}
	}
	if sinf == nil {
		return nil, fmt.Errorf("enrollment model requires one installed SINF mapping")
	}
	user, keyID, err := enrollmentSINFSelector(sinf)
	if err != nil {
		return nil, err
	}
	selected, err := selectAuthenticatedEnrollmentClassPrecursor(sinf, registries, hardware, user, keyID)
	if err != nil {
		return nil, err
	}
	pack, err := enrollmentPackForPrecursor(typedEnrollmentAlgorithmIdentity, user, keyID, selected.Precursor)
	clear(selected.Precursor[:])
	return pack, err
}

func enrollmentPackForPrecursor(algorithmIdentity [32]byte, user, keyID uint32, precursor [16]byte) ([]byte, error) {
	pack := make([]byte, 0, typedEnrollmentModelBytes)
	pack = append(pack, []byte("FPSENR2\x00")...)
	var word [4]byte
	binary.LittleEndian.PutUint32(word[:], 2)
	pack = append(pack, word[:]...)
	binary.LittleEndian.PutUint32(word[:], typedEnrollmentModelBytes)
	pack = append(pack, word[:]...)
	pack = append(pack, algorithmIdentity[:]...)
	binary.LittleEndian.PutUint32(word[:], user)
	pack = append(pack, word[:]...)
	binary.LittleEndian.PutUint32(word[:], keyID)
	pack = append(pack, word[:]...)
	modelKey := currentSINFModelKeyFromPrecursor(precursor)
	pack = append(pack, modelKey[:]...)
	pack = append(pack, make([]byte, 16)...)
	for range 16 {
		for value := range 256 {
			pack = append(pack, byte(value))
		}
	}
	clear(modelKey[:])
	if len(pack) != typedEnrollmentModelBytes {
		clear(pack)
		return nil, fmt.Errorf("typed enrollment model has invalid size %d", len(pack))
	}
	return pack, nil
}
