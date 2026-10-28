package recovery

import (
	"crypto/aes"
	"encoding/binary"
	"fmt"
)

const (
	selectionSIDBContainerBytes        = 993
	selectionSIDBExtendedBytes         = 1003
	selectionSIDBPayloadOffset         = 33
	selectionSIDBExtendedPayloadOffset = 43
	selectionSIDBPayloadBytes          = 944
	selectionSIDBBlockBytes            = aes.BlockSize
	selectionSIDBExtendedTrailerBytes  = 16
	selectionSIDBOuterIterations       = outerNaturalIterationCount
	selectionSIDBExtendedIterations    = 427
)

type currentSIDBLayout struct {
	CodecBytes      int
	PayloadOffset   int
	OuterIterations int
}

func currentSIDBLayoutForContainer(input []byte) (currentSIDBLayout, error) {
	if len(input) < 4 || binary.BigEndian.Uint32(input[:4]) != currentSIDBContainerVersion {
		return currentSIDBLayout{}, fmt.Errorf("current SIDB container has an unsupported version")
	}
	switch len(input) {
	case selectionSIDBContainerBytes:
		return currentSIDBLayout{
			CodecBytes:      selectionSIDBPayloadOffset - 4,
			PayloadOffset:   selectionSIDBPayloadOffset,
			OuterIterations: selectionSIDBOuterIterations,
		}, nil
	case selectionSIDBExtendedBytes:
		return currentSIDBLayout{
			CodecBytes:      selectionSIDBExtendedPayloadOffset - 4,
			PayloadOffset:   selectionSIDBExtendedPayloadOffset,
			OuterIterations: selectionSIDBExtendedIterations,
		}, nil
	default:
		return currentSIDBLayout{}, fmt.Errorf("current SIDB container has unsupported size %d", len(input))
	}
}

// These are fixed semantic external encodings, not account or build material.
// The protected decoder applies each encoding as a byte-wise XOR around an
// ordinary AES-128 decryption.
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

// The selection sweep reduces either structurally supported current container
// to the same 59 encoded AES blocks. The extended layout has ten additional
// codec bytes before the body and a 16-byte trailer after it. Neither belongs
// to the encrypted registry body.
func decodeSelectionSIDBContainer(input []byte) ([selectionSIDBPayloadBytes]byte, error) {
	var output [selectionSIDBPayloadBytes]byte
	layout, err := currentSIDBLayoutForContainer(input)
	if err != nil {
		return output, err
	}
	end := layout.PayloadOffset + len(output)
	if end > len(input) {
		return output, fmt.Errorf("selection SIDB body exceeds its container")
	}
	for index := range output {
		output[index] = input[layout.PayloadOffset+index] ^ selectionSIDBInputMask[index%selectionSIDBBlockBytes]
	}
	return output, nil
}

// decryptSelectionSIDBPayload replaces the protected 59-block white-box loop
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
