package recovery

import "testing"

func TestCurrentSINFMaterialProjection(t *testing.T) {
	precursor := [16]byte{
		0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77,
		0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff,
	}
	material := currentSINFMaterialFromPrecursor(precursor)
	key := currentSINFModelKeyFromPrecursor(precursor)
	for index := range key {
		if got, want := material[index]^generatedCurrentSINFKeyMask[index], key[index]; got != want {
			t.Fatalf("key byte %d=%#x, want %#x", index, got, want)
		}
		if got, want := currentSINFMaterialMask[index]^generatedCurrentSINFKeyMask[index], byte(0x1b); got != want {
			t.Fatalf("combined mask byte %d=%#x, want %#x", index, got, want)
		}
	}
	if key != [16]byte{
		0x1b, 0x0a, 0x39, 0x28, 0x5f, 0x4e, 0x7d, 0x6c,
		0x93, 0x82, 0xb1, 0xa0, 0xd7, 0xc6, 0xf5, 0xe4,
	} {
		t.Fatalf("direct current SINF key=%x", key)
	}
}
