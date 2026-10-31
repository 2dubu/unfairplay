package recovery

import (
	"encoding/binary"
	"fmt"
)

var currentClassRecordMask = [16]byte{
	0x49, 0x8b, 0x22, 0x11, 0x82, 0xb3, 0x19, 0x85,
	0x70, 0xdd, 0x90, 0x4f, 0xd9, 0x8c, 0x93, 0x17,
}

// deriveClassValuFromMaskedWords converts the already-masked 16-byte record
// window into the semantic ciphertext consumed by the typed class join.
func deriveClassValuFromMaskedWords(masked [16]byte) (valu [16]byte) {
	for word := range 4 {
		offset := word * 4
		x := binary.BigEndian.Uint32(masked[offset:])

		// This single modular affine operation is the complete semantic
		// transform. It includes the native masked arithmetic, byte encoder,
		// and following XOR-0x1b input representation.
		binary.BigEndian.PutUint32(valu[offset:], x*uint32(0xa46fd0ff)+uint32(0xfe063531))
	}
	return valu
}

func classRecordMaximumWindowOffset(record decodedSIDBRecord) int {
	if record.TailAvailable < 0 || record.TailAvailable > decodedSIDBTailBytes {
		return -1
	}
	return record.TailAvailable
}

func deriveClassValuWindow(record decodedSIDBRecord, windowOffset int) ([16]byte, error) {
	if windowOffset < 0 || windowOffset > classRecordMaximumWindowOffset(record) {
		return [16]byte{}, fmt.Errorf("class record window offset is invalid: %d", windowOffset)
	}
	var material [decodedSIDBPayloadBytes + decodedSIDBTailBytes]byte
	copy(material[:decodedSIDBPayloadBytes], record.Payload[:])
	copy(material[decodedSIDBPayloadBytes:], record.Tail[:])
	var masked [16]byte
	for index := range masked {
		position := windowOffset + index
		maskPosition := position + record.MaskOffset
		masked[index] = material[position] ^ currentClassRecordMask[maskPosition%len(currentClassRecordMask)]
	}
	return deriveClassValuFromMaskedWords(masked), nil
}

func deriveCurrentSINFPrecursorFromRecordWindow(record decodedSIDBRecord, windowOffset int, classKey [16]byte) ([16]byte, error) {
	valu, err := deriveClassValuWindow(record, windowOffset)
	if err != nil {
		return [16]byte{}, err
	}
	return deriveCurrentSINFPrecursorFromClass(classKey, valu)
}
