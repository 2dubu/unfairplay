package recovery

import (
	"crypto/aes"
	"crypto/cipher"
)

// These are public algorithm constants, not enrollment credentials. The
// recovered round keys 1..10 are the AES-128 schedule for selectionSIDBAESKey.
// The original initial key differs by selectionSIDBInputWhitening; its XOR
// is absorbed into the input. The suffix/producer cancellation is proved
// independently by TestSIDBKeySuffixCancellation.
var selectionSIDBAESKey = [16]byte{
	0xa3, 0x22, 0xec, 0xaf, 0xa8, 0xc2, 0x68, 0xbf,
	0x49, 0x3c, 0x21, 0xc5, 0xd2, 0x57, 0x65, 0xfc,
}

var selectionSIDBInputWhitening = [16]byte{
	0xe6, 0xf8, 0xca, 0x44, 0x9d, 0x25, 0x58, 0xad,
	0x57, 0x3c, 0x2d, 0x76, 0x5d, 0x9c, 0xa5, 0xcc,
}

var selectionSIDBCipher = func() cipher.Block {
	block, err := aes.NewCipher(selectionSIDBAESKey[:])
	if err != nil {
		panic(err)
	}
	return block
}()

func deriveSelectionSIDBKey(input [16]byte) (output [16]byte) {
	for index := range input {
		input[index] ^= selectionSIDBInputWhitening[index]
	}
	selectionSIDBCipher.Encrypt(output[:], input[:])
	return output
}
