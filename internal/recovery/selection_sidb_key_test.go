package recovery

import (
	"crypto/sha256"
	"testing"
)

var selectionSIDBTestRoundKeyDecoder = selectionSIDBRoundKeyDecoder{
	{0x4f, 0x14, 0x23, 0xd7, 0x1f5, 0x19b, 0xe9, 0x133},
	{0x80, 0x9b, 0x182, 0x1e6, 0x2e, 0xf1, 0xb3, 0x12f},
	{0x149, 0x133, 0x157, 0x1f6, 0x15b, 0xf9, 0x16c, 0x1d7},
	{0x171, 0xe8, 0xb7, 0x1b, 0xd2, 0xa6, 0x12d, 0x1d0},
	{0x4e, 0x15, 0x22, 0xd6, 0x1f4, 0x19a, 0xe8, 0x132},
	{0x81, 0x9b, 0x182, 0x1e6, 0x2f, 0xf1, 0xb3, 0x12f},
	{0x148, 0x133, 0x157, 0x1f7, 0x15a, 0xf9, 0x16c, 0x1d7},
	{0x170, 0xe8, 0xb7, 0x1b, 0xd3, 0xa6, 0x12d, 0x1d1},
	{0x4f, 0x15, 0x23, 0xd7, 0x1f4, 0x19a, 0xe8, 0x133},
	{0x80, 0x9a, 0x183, 0x1e6, 0x2e, 0xf1, 0xb2, 0x12f},
	{0x149, 0x133, 0x157, 0x1f7, 0x15b, 0xf8, 0x16d, 0x1d6},
	{0x170, 0xe8, 0xb7, 0x1a, 0xd3, 0xa7, 0x12c, 0x1d0},
	{0x4f, 0x15, 0x22, 0xd7, 0x1f5, 0x19a, 0xe9, 0x133},
	{0x81, 0x9a, 0x182, 0x1e6, 0x2e, 0xf0, 0xb3, 0x12e},
	{0x148, 0x133, 0x156, 0x1f7, 0x15a, 0xf9, 0x16c, 0x1d7},
	{0x170, 0xe9, 0xb6, 0x1b, 0xd3, 0xa7, 0x12d, 0x1d1},
}

func expandSelectionSIDBAES128Key(key [16]byte) (result [176]byte) {
	copy(result[:16], key[:])
	rcon := byte(1)
	for offset := 16; offset < len(result); offset += 16 {
		previous := result[offset-4 : offset]
		word := [4]byte{
			frameSBox[previous[1]] ^ rcon,
			frameSBox[previous[2]],
			frameSBox[previous[3]],
			frameSBox[previous[0]],
		}
		for index := range 4 {
			result[offset+index] = result[offset-16+index] ^ word[index]
		}
		for index := 4; index < 16; index++ {
			result[offset+index] = result[offset-16+index] ^ result[offset+index-4]
		}
		rcon = frameGFMul(rcon, 2)
	}
	return result
}

func encodeSelectionSIDBRoundKeyByte(t *testing.T, decoder selectionSIDBAffineByteDecoder, want byte) byte {
	t.Helper()
	for value := range 256 {
		if decoder.apply(byte(value)) == want {
			return byte(value)
		}
	}
	t.Fatalf("round-key decoder has no encoding for %#x", want)
	return 0
}

func TestSelectionSIDBRoundKeyDecoder(t *testing.T) {
	for test := range 256 {
		key := sha256.Sum256([]byte{byte(test), byte(test >> 8), 0x53, 0x49, 0x44, 0x42})
		var want [16]byte
		copy(want[:], key[:16])
		expanded := expandSelectionSIDBAES128Key(want)
		var encoded [16]byte
		for position := range encoded {
			encoded[position] = encodeSelectionSIDBRoundKeyByte(t, selectionSIDBTestRoundKeyDecoder[position], expanded[160+position])
		}
		got, err := selectionSIDBTestRoundKeyDecoder.decoderKey(encoded)
		if err != nil {
			t.Fatal(err)
		}
		if got != want {
			t.Fatalf("test %d decoder key = %x, want %x", test, got, want)
		}
	}
}

func TestSelectionSIDBRoundKeyDecoderRejectsNonBijectiveEncoding(t *testing.T) {
	decoder := selectionSIDBTestRoundKeyDecoder
	decoder[7] = selectionSIDBAffineByteDecoder{}
	if _, err := decoder.decoderKey([16]byte{}); err == nil {
		t.Fatal("non-bijective round-key decoder was accepted")
	}
}
