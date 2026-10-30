package tagengine

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"fmt"
)

// These are public algorithm constants from the mapped constructor, not
// enrollment credentials. The decoded round keys form this AES-128 schedule;
// its input/output byte encodings and transport whitening cancel.
var constructorAESKey = [aes.BlockSize]byte{0x8c, 0x92, 0x0e, 0x14, 0x8b, 0x4f, 0x48, 0xa6, 0x2f, 0x8a, 0x54, 0x97, 0x95, 0x23, 0x38, 0xcd}
var constructorIV = [aes.BlockSize]byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f}
var constructorCipher = func() cipher.Block {
	block, err := aes.NewCipher(constructorAESKey[:])
	if err != nil {
		panic(err)
	}
	return block
}()

// transformConstructorState encrypts the decoded mixer state with AES-CBC.
// Explicit chaining retains block-boundary cancellation and state wiping.
func transformConstructorState(ctx context.Context, state *[256]byte) error {
	if ctx == nil || state == nil {
		return fmt.Errorf("invalid constructor cipher input")
	}
	chain := constructorIV
	var input [aes.BlockSize]byte
	defer clear(chain[:])
	defer clear(input[:])
	for start := 0; start < len(state); start += aes.BlockSize {
		if err := ctx.Err(); err != nil {
			return err
		}
		for i := range input {
			input[i] = state[start+i] ^ chain[i]
		}
		constructorCipher.Encrypt(state[start:start+aes.BlockSize], input[:])
		copy(chain[:], state[start:start+aes.BlockSize])
	}
	return nil
}
