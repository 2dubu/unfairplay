package recovery

import (
	"encoding/binary"
	"fmt"
	"strings"
)

const typedEnrollmentModelBytes = 8 + 4 + 4 + 32 + 4 + 4 + 16 + 16 + 16*256

var typedEnrollmentSourceHash = [32]byte{0x1e, 0x58, 0xd3, 0x7e, 0x96, 0x02, 0xdd, 0x0c, 0xa5, 0xe2, 0x4b, 0x14, 0x76, 0x20, 0x63, 0xb6, 0xc3, 0x10, 0x24, 0x8b, 0x4d, 0x3a, 0x94, 0x27, 0x27, 0x71, 0xd8, 0x15, 0xe4, 0xea, 0x17, 0x84}

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

func currentEnrollmentPack(files map[string][]byte, profile *semanticBuildProfile, registries []decodedSIDBRegistryCandidate) ([]byte, error) {
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
	if profile == nil {
		return nil, fmt.Errorf("semantic enrollment profile is unavailable")
	}
	selected, err := selectAuthenticatedClassPrecursor(sinf, profile, registries)
	if err != nil {
		return nil, err
	}
	pack, err := enrollmentPackForPrecursor(profile, typedEnrollmentSourceHash, user, keyID, selected.Precursor)
	clear(selected.Precursor[:])
	return pack, err
}

func enrollmentPackForPrecursor(profile *semanticBuildProfile, binaryHash [32]byte, user, keyID uint32, precursor [16]byte) ([]byte, error) {
	if profile == nil {
		return nil, fmt.Errorf("semantic enrollment profile is unavailable")
	}
	pack := make([]byte, 0, typedEnrollmentModelBytes)
	pack = append(pack, []byte("FPSENR2\x00")...)
	var word [4]byte
	binary.LittleEndian.PutUint32(word[:], 2)
	pack = append(pack, word[:]...)
	binary.LittleEndian.PutUint32(word[:], typedEnrollmentModelBytes)
	pack = append(pack, word[:]...)
	pack = append(pack, binaryHash[:]...)
	binary.LittleEndian.PutUint32(word[:], user)
	pack = append(pack, word[:]...)
	binary.LittleEndian.PutUint32(word[:], keyID)
	pack = append(pack, word[:]...)
	modelKey := currentSINFModelKeyFromPrecursor(precursor)
	pack = append(pack, modelKey[:]...)
	pack = append(pack, profile.currentSINFWhite[:]...)
	pack = append(pack, profile.currentSINFInput[:]...)
	clear(modelKey[:])
	if len(pack) != typedEnrollmentModelBytes {
		clear(pack)
		return nil, fmt.Errorf("typed enrollment model has invalid size %d", len(pack))
	}
	return pack, nil
}
