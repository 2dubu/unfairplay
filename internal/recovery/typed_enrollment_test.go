package recovery

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"testing"
)

func TestTypedEnrollmentAlgorithmIdentity(t *testing.T) {
	want := sha256.Sum256([]byte("unfairplay/typed-current-enrollment/v1"))
	if typedEnrollmentAlgorithmIdentity != want {
		t.Fatalf("typed enrollment algorithm identity = %x, want %x", typedEnrollmentAlgorithmIdentity, want)
	}
}

func TestEnrollmentPackForPrecursor(t *testing.T) {
	profile := portableCurrentSemanticProfile()
	var precursor [16]byte
	for index := range precursor {
		precursor[index] = byte(index*17 + 3)
	}
	pack, err := enrollmentPackForPrecursor(profile, typedEnrollmentAlgorithmIdentity, 0x12345678, 2, precursor)
	if err != nil {
		t.Fatal(err)
	}
	defer clear(pack)
	if len(pack) != typedEnrollmentModelBytes || string(pack[:8]) != "FPSENR2\x00" ||
		binary.LittleEndian.Uint32(pack[8:12]) != 2 || binary.LittleEndian.Uint32(pack[12:16]) != uint32(len(pack)) ||
		!bytes.Equal(pack[16:48], typedEnrollmentAlgorithmIdentity[:]) || binary.LittleEndian.Uint32(pack[48:52]) != 0x12345678 ||
		binary.LittleEndian.Uint32(pack[52:56]) != 2 {
		t.Fatal("typed enrollment pack envelope is invalid")
	}
	wantKey := currentSINFModelKeyFromPrecursor(precursor)
	if !bytes.Equal(pack[56:72], wantKey[:]) || !bytes.Equal(pack[72:88], profile.currentSINFWhite[:]) ||
		!bytes.Equal(pack[88:], profile.currentSINFInput[:]) {
		t.Fatal("typed enrollment pack payload is invalid")
	}
}
