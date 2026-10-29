package recovery

import (
	"crypto/aes"
	"encoding/binary"
	"fmt"
)

const (
	selectionSIDBOlderContainerBytes   = 630
	selectionSIDBOlderPayloadOffset    = 38
	selectionSIDBOlderPayloadBytes     = 576
	selectionSIDBOlderIterations       = 372
	selectionSIDBContainerBytes        = 993
	selectionSIDBExtendedBytes         = 1003
	selectionSIDBPayloadOffset         = 33
	selectionSIDBExtendedPayloadOffset = 43
	selectionSIDBPayloadBytes          = 944
	selectionSIDBBlockBytes            = aes.BlockSize
	selectionSIDBTrailerBytes          = 16
	selectionSIDBOuterIterations       = outerNaturalIterationCount
	selectionSIDBExtendedIterations    = 427
	selectionSIDBMinimumCodecBytes     = 13
	selectionSIDBMaximumCodecBytes     = 39
	selectionSIDBIterationsPerCodec    = 11
	selectionSIDBIterationBias         = -2
	maximumStructuralSIDBRecords       = 512
)

type currentSIDBLayout struct {
	CodecBytes      int
	PayloadOffset   int
	PayloadBytes    int
	OuterIterations int
}

func currentSIDBLayoutForContainer(input []byte) (currentSIDBLayout, error) {
	if len(input) < 4 || binary.BigEndian.Uint32(input[:4]) != currentSIDBContainerVersion {
		return currentSIDBLayout{}, fmt.Errorf("current SIDB container has an unsupported version")
	}
	switch len(input) {
	case selectionSIDBOlderContainerBytes:
		return currentSIDBLayout{
			CodecBytes:      selectionSIDBOlderPayloadOffset - 4,
			PayloadOffset:   selectionSIDBOlderPayloadOffset,
			PayloadBytes:    selectionSIDBOlderPayloadBytes,
			OuterIterations: selectionSIDBOlderIterations,
		}, nil
	case selectionSIDBContainerBytes:
		return currentSIDBLayout{
			CodecBytes:      selectionSIDBPayloadOffset - 4,
			PayloadOffset:   selectionSIDBPayloadOffset,
			PayloadBytes:    selectionSIDBPayloadBytes,
			OuterIterations: selectionSIDBOuterIterations,
		}, nil
	case selectionSIDBExtendedBytes:
		return currentSIDBLayout{
			CodecBytes:      selectionSIDBExtendedPayloadOffset - 4,
			PayloadOffset:   selectionSIDBExtendedPayloadOffset,
			PayloadBytes:    selectionSIDBPayloadBytes,
			OuterIterations: selectionSIDBExtendedIterations,
		}, nil
	}

	// Every observed current container is version || codec || encrypted
	// registry || 16-byte trailer. The decoded registry is a bounded sequence
	// of 40-byte records, optionally followed by a 16- or 24-byte footer. The
	// outer loop consumes one codec byte every eleven iterations, with the same
	// two-step bias in all three independently validated layouts.
	var layouts []currentSIDBLayout
	for codecBytes := selectionSIDBMinimumCodecBytes; codecBytes <= selectionSIDBMaximumCodecBytes; codecBytes++ {
		payloadBytes := len(input) - 4 - codecBytes - selectionSIDBTrailerBytes
		if payloadBytes <= 0 || payloadBytes%selectionSIDBBlockBytes != 0 || !validDecodedSIDBRegistrySize(payloadBytes) {
			continue
		}
		layouts = append(layouts, currentSIDBLayout{
			CodecBytes:      codecBytes,
			PayloadOffset:   4 + codecBytes,
			PayloadBytes:    payloadBytes,
			OuterIterations: selectionSIDBIterationsPerCodec*codecBytes + selectionSIDBIterationBias,
		})
	}
	if len(layouts) != 1 {
		return currentSIDBLayout{}, fmt.Errorf("current SIDB container has unsupported or ambiguous size %d", len(input))
	}
	return layouts[0], nil
}

func validDecodedSIDBRegistrySize(size int) bool {
	for _, footerBytes := range [...]int{decodedSIDBNoSeparateFooterBytes, decodedSIDBOlderFooterBytes, decodedSIDBFooterBytes} {
		if size > footerBytes && (size-footerBytes)%decodedSIDBRecordBytes == 0 {
			count := (size - footerBytes) / decodedSIDBRecordBytes
			return count > 0 && count <= maximumStructuralSIDBRecords
		}
	}
	return false
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

// The selection sweep reduces each structurally supported container to its
// encoded AES blocks. The two current layouts contain 59 blocks; the older
// 630-byte layout contains 36. Codec prefixes and optional trailers do not
// belong to the encrypted registry body.
func decodeSelectionSIDBContainer(input []byte) ([]byte, error) {
	layout, err := currentSIDBLayoutForContainer(input)
	if err != nil {
		return nil, err
	}
	output := make([]byte, layout.PayloadBytes)
	end := layout.PayloadOffset + len(output)
	if end > len(input) {
		return nil, fmt.Errorf("selection SIDB body exceeds its container")
	}
	for index := range output {
		output[index] = input[layout.PayloadOffset+index] ^ selectionSIDBInputMask[index%selectionSIDBBlockBytes]
	}
	return output, nil
}

// decryptSelectionSIDBPayload replaces the protected block loop once its
// per-enrollment AES key has been recovered by the upstream producer.
func decryptSelectionSIDBPayload(input []byte, key [selectionSIDBBlockBytes]byte) ([]byte, error) {
	if len(input) == 0 || len(input)%selectionSIDBBlockBytes != 0 {
		return nil, fmt.Errorf("selection SIDB payload has invalid size %d", len(input))
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
