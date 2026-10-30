package recovery

import (
	"crypto/aes"
	"fmt"
)

const (
	selectionSIDBBlockBytes           = aes.BlockSize
	selectionSIDBExtendedTrailerBytes = 16
	selectionSIDBLargerSourceBytes    = 256
	maximumSelectionSIDBPayloadBytes  = 1 << 20
	maximumMappedSelectionModulus     = 18
)

type currentSIDBLayout struct {
	CodecBytes      int
	PayloadOffset   int
	PayloadBytes    int
	SourceBytes     int
	SelectorModulus int
	SelectedSlot    int
	RootCursor      int
	OuterIterations int
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

// decodeSelectionSIDBContainerWithLayout validates the payload and trailer.
func decodeSelectionSIDBContainerWithLayout(input []byte, layout currentSIDBLayout) ([]byte, error) {
	if layout.PayloadBytes <= 0 || layout.PayloadBytes > maximumSelectionSIDBPayloadBytes || layout.PayloadBytes%selectionSIDBBlockBytes != 0 || layout.PayloadOffset < 4 || layout.PayloadOffset > len(input)-selectionSIDBExtendedTrailerBytes || layout.PayloadBytes != len(input)-selectionSIDBExtendedTrailerBytes-layout.PayloadOffset {
		return nil, fmt.Errorf("selection SIDB has invalid body geometry")
	}
	output := make([]byte, layout.PayloadBytes)
	for index := range output {
		output[index] = input[layout.PayloadOffset+index] ^ selectionSIDBInputMask[index%selectionSIDBBlockBytes]
	}
	return output, nil
}

// decryptSelectionSIDBPayloadDynamic replaces the protected block loop once
// its per-enrollment AES key has been recovered by the upstream producer.
func decryptSelectionSIDBPayloadDynamic(input []byte, key [selectionSIDBBlockBytes]byte) ([]byte, error) {
	if len(input) == 0 || len(input) > maximumSelectionSIDBPayloadBytes || len(input)%selectionSIDBBlockBytes != 0 {
		return nil, fmt.Errorf("selection SIDB body has invalid length %d", len(input))
	}
	output := make([]byte, len(input))
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
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
