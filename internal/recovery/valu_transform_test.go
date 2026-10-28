package recovery

import (
	"crypto/sha256"
	"encoding/binary"
	"testing"
)

func protectedClassValuReference(payload [16]byte, wordMasks [4]uint32) (valu [16]byte) {
	for word := range 4 {
		offset := word * 4
		x := binary.BigEndian.Uint32(payload[offset:]) ^ wordMasks[word]
		term := (x & uint32(0x4e757bdd)) ^ ((x & 0xff) & 0x18)
		y := x - 2*term + uint32(0xce757bc5)
		q := ((y ^ uint32(0xa5653f52)) + uint32(0x5a9ac0ae)) ^
			((y ^ uint32(0x0d404124)) + uint32(0xf2bfbedc)) ^
			((y ^ uint32(0x665005b3)) + uint32(0x99affa4d))
		z := q*uint32(0xa46fd0ff) -
			((q*uint32(0x08dfa1fe) + uint32(0x300f1cd8)) & uint32(0x377d4f36)) +
			uint32(0x93c63607)

		a := byte(z >> 24)
		encoded := a ^ 0x9b
		correction := (encoded & 0x1f) ^ (a & 4)
		valu[offset] = (encoded - 2*correction + 0x1b) ^ 0x1b
		encoded = byte(z>>16) ^ 0xbe
		valu[offset+1] = (encoded - (0x36 & (encoded << 1)) + 0x1b) ^ 0x1b
		encoded = byte(z>>8) ^ 0xa7
		valu[offset+2] = (encoded - (0x36 & (encoded << 1)) + 0x1b) ^ 0x1b
		valu[offset+3] = (byte(z) ^ 0x80) ^ 0x1b
	}
	return valu
}

func TestDeriveClassValuMatchesProtectedArithmetic(t *testing.T) {
	for control := range 4096 {
		digest := sha256.Sum256([]byte{0x56, 0x41, 0x4c, byte(control >> 8), byte(control)})
		var payload [16]byte
		copy(payload[:], digest[:16])
		masks := [4]uint32{
			binary.BigEndian.Uint32(digest[0:4]),
			binary.BigEndian.Uint32(digest[4:8]) ^ 0x13579bdf,
			binary.BigEndian.Uint32(digest[8:12]) ^ 0x2468ace0,
			binary.BigEndian.Uint32(digest[12:16]) ^ 0xa5a5a5a5,
		}
		if got, want := deriveClassValu(payload, masks), protectedClassValuReference(payload, masks); got != want {
			t.Fatalf("control %d = %x, want %x", control, got, want)
		}
	}
}

func TestDeriveCurrentSINFPrecursorFromRecord(t *testing.T) {
	payload := [16]byte{0x3d, 0x79, 0x22, 0x11, 0x48, 0x61, 0xa7, 0xc4, 0x6a, 0xb0, 0x52, 0x8f, 0x10, 0x39, 0xde, 0x72}
	masks := [4]uint32{0x10203040, 0x50607080, 0x90a0b0c0, 0xd0e0f000}
	key := [16]byte{0x01, 0x12, 0x23, 0x34, 0x45, 0x56, 0x67, 0x78, 0x89, 0x9a, 0xab, 0xbc, 0xcd, 0xde, 0xef, 0xf0}
	precursor, err := deriveCurrentSINFPrecursorFromRecord(payload, masks, key)
	if err != nil {
		t.Fatal(err)
	}
	valu := deriveClassValu(payload, masks)
	want, err := deriveCurrentSINFPrecursorFromClass(key, valu)
	if err != nil {
		t.Fatal(err)
	}
	if precursor != want {
		t.Fatalf("precursor = %x, want %x", precursor, want)
	}
}
