package recovery

import (
	"crypto/aes"
	"fmt"
)

const (
	selectionSIDBContainerBytes = 993
	selectionSIDBPayloadOffset  = 33
	selectionSIDBPayloadBytes   = 944
	selectionSIDBBlockBytes     = aes.BlockSize
)

// These are build-global external encodings, not account material. The native
// decoder applies each encoding as a byte-wise XOR around an ordinary AES-128
// decryption.
var selectionSIDBInputMask = [selectionSIDBBlockBytes]byte{
	0x0c, 0xc2, 0xa2, 0xea, 0xe7, 0xb1, 0x7b, 0x94,
	0xb7, 0xd8, 0x82, 0x58, 0x10, 0x4b, 0xc5, 0xad,
}

var selectionSIDBOutputMask = [selectionSIDBBlockBytes]byte{
	0xe4, 0x34, 0x93, 0x3c, 0x0b, 0x85, 0x61, 0x5e,
	0x68, 0x18, 0xfa, 0x61, 0x4f, 0x2f, 0x29, 0x9c,
}

// Blocks after the first are chained in the external byte domain. This is not
// standard CBC: the previous encoded input and this fixed bias are XORed into
// the current encoded output after AES decryption.
var selectionSIDBChainMask = [selectionSIDBBlockBytes]byte{
	0x31, 0x7a, 0xa2, 0xc1, 0xa5, 0xbf, 0x38, 0xdb,
	0x5d, 0x73, 0x61, 0xbc, 0x2f, 0xb9, 0x7c, 0x7e,
}

// The selection sweep reduces the 993-byte container to 59 encoded AES
// blocks. Its only data transform is this repeating XOR mask.
func decodeSelectionSIDBContainer(input []byte) ([selectionSIDBPayloadBytes]byte, error) {
	var output [selectionSIDBPayloadBytes]byte
	if len(input) != selectionSIDBContainerBytes {
		return output, fmt.Errorf("selection SIDB container has %d bytes, want %d", len(input), selectionSIDBContainerBytes)
	}
	for index := range output {
		output[index] = input[selectionSIDBPayloadOffset+index] ^ selectionSIDBInputMask[index%selectionSIDBBlockBytes]
	}
	return output, nil
}

// decryptSelectionSIDBPayload replaces the native 59-block white-box loop
// once its per-enrollment AES key has been recovered by the upstream producer.
func decryptSelectionSIDBPayload(input [selectionSIDBPayloadBytes]byte, key [selectionSIDBBlockBytes]byte) ([selectionSIDBPayloadBytes]byte, error) {
	var output [selectionSIDBPayloadBytes]byte
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return output, err
	}
	for offset := 0; offset < len(input); offset += selectionSIDBBlockBytes {
		var encodedInput, decodedOutput [selectionSIDBBlockBytes]byte
		for index := range encodedInput {
			encodedInput[index] = input[offset+index] ^ selectionSIDBInputMask[index]
		}
		block.Decrypt(decodedOutput[:], encodedInput[:])
		for index := range decodedOutput {
			value := decodedOutput[index] ^ selectionSIDBOutputMask[index]
			if offset != 0 {
				value ^= input[offset-selectionSIDBBlockBytes+index] ^ selectionSIDBChainMask[index]
			}
			output[offset+index] = value
		}
	}
	return output, nil
}
