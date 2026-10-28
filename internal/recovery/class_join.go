package recovery

import (
	"bytes"
	"crypto/aes"
	"fmt"
)

// decodeClassSecondary removes the protected byte representation from the
// secondary class operand. Current enrollment state carries one eight-byte
// value twice; rejecting unequal halves keeps this boundary fail-closed.
func decodeClassSecondary(encoded [aes.BlockSize]byte) ([aes.BlockSize]byte, error) {
	if !bytes.Equal(encoded[:aes.BlockSize/2], encoded[aes.BlockSize/2:]) {
		return [aes.BlockSize]byte{}, fmt.Errorf("class secondary operand has unequal halves")
	}
	var key [aes.BlockSize]byte
	for index := range key {
		key[index] = encoded[index] ^ 0x1b
	}
	return key, nil
}

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
