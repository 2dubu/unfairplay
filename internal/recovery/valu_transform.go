package recovery

import "encoding/binary"

// deriveClassValu converts the protected 16-byte valu payload into the
// semantic ciphertext consumed by the typed class join. The four masks are
// enrollment data recovered from the enclosing registry record; all other
// constants belong to the shared 32-bit transform.
func deriveClassValu(payload [16]byte, wordMasks [4]uint32) (valu [16]byte) {
	for word := range 4 {
		offset := word * 4
		x := binary.BigEndian.Uint32(payload[offset:]) ^ wordMasks[word]

		// The native masked subtract/add stage is exactly this XOR.
		y := x ^ uint32(0xce757bc5)
		q := ((y ^ uint32(0xa5653f52)) + uint32(0x5a9ac0ae)) ^
			((y ^ uint32(0x0d404124)) + uint32(0xf2bfbedc)) ^
			((y ^ uint32(0x665005b3)) + uint32(0x99affa4d))
		z := q*uint32(0xa46fd0ff) -
			((q*uint32(0x08dfa1fe) + uint32(0x300f1cd8)) & uint32(0x377d4f36)) +
			uint32(0x93c63607)

		// This combines the native byte encoder with the following XOR-0x1b
		// input representation. The result is the semantic AES ciphertext.
		binary.BigEndian.PutUint32(valu[offset:], z^uint32(0x9bbea79b))
	}
	return valu
}

func deriveCurrentSINFPrecursorFromRecord(payload [16]byte, wordMasks [4]uint32, classKey [16]byte) ([16]byte, error) {
	return deriveCurrentSINFPrecursorFromClass(classKey, deriveClassValu(payload, wordMasks))
}
