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

		// This single modular affine operation is the complete semantic
		// transform. It includes the native masked arithmetic, byte encoder,
		// and following XOR-0x1b input representation.
		binary.BigEndian.PutUint32(valu[offset:], x*uint32(0xa46fd0ff)+uint32(0xfe063531))
	}
	return valu
}

func deriveCurrentSINFPrecursorFromRecord(payload [16]byte, wordMasks [4]uint32, classKey [16]byte) ([16]byte, error) {
	return deriveCurrentSINFPrecursorFromClass(classKey, deriveClassValu(payload, wordMasks))
}
