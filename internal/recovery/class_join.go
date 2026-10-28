package recovery

import (
	"crypto/aes"
	"fmt"
)

// deriveCurrentSINFPrecursorFromClass replaces the protected two-input class
// consumer. The secondary class value is an AES-128 key and the primary valu
// field is one ciphertext block. The native consumer stores its plaintext in
// the protected byte representation used by the following material writer.
func deriveCurrentSINFPrecursorFromClass(classKey, valu [aes.BlockSize]byte) ([aes.BlockSize]byte, error) {
	var plaintext, precursor [aes.BlockSize]byte
	block, err := aes.NewCipher(classKey[:])
	if err != nil {
		return precursor, fmt.Errorf("create class AES: %w", err)
	}
	block.Decrypt(plaintext[:], valu[:])
	for index, value := range plaintext {
		precursor[index] = applyCurrentFairPlayByteMask(value)
	}
	clear(plaintext[:])
	return precursor, nil
}
