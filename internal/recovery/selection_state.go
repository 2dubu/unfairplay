package recovery

import (
	"crypto/aes"
	"crypto/cipher"
)

// The recovered state update is one slot-independent CBC recurrence. The
// native source mapper bounds the sequence to this many complete blocks.
const selectionStateSlotCount = maximumMappedSelectionModulus

// Public algorithm constants, not enrollment credentials. The mapped shared
// round keys form this AES-128 schedule. Slot-specific initial keys reduce to
// one CBC IV and a repeating input XOR; the output XOR cancels the SIDB key
// producer's input whitening.
var selectionStateAESKey = [16]byte{
	0x36, 0xc0, 0x11, 0x2f, 0xc8, 0x4c, 0x50, 0x17,
	0x62, 0x7f, 0x53, 0x53, 0x35, 0xe5, 0xe7, 0xc9,
}

var selectionStateIV = [16]byte{
	0x1c, 0x64, 0x35, 0xdf, 0x5c, 0xe3, 0x16, 0xce,
	0x0f, 0x8f, 0xd0, 0xa4, 0x17, 0x9d, 0x12, 0x56,
}

const selectionStateInputWhitening byte = 0x1b

var selectionStateCipher = func() cipher.Block {
	block, err := aes.NewCipher(selectionStateAESKey[:])
	if err != nil {
		panic(err)
	}
	return block
}()

// advanceSelectionState keeps the decoded CBC state between blocks.
func advanceSelectionState(previous *[16]byte, input []byte) {
	var block [16]byte
	defer clear(block[:])
	for i := range block {
		block[i] = input[i] ^ selectionStateInputWhitening ^ previous[i]
	}
	selectionStateCipher.Encrypt(previous[:], block[:])
}
