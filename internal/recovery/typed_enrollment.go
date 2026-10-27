package recovery

import (
	"encoding/binary"
	"fmt"
)

const typedEnrollmentModelBytes = 8 + 4 + 4 + 32 + 4 + 4 + 16 + 16 + 16*256

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
