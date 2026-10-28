package recovery

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/binary"
	"fmt"
)

type authenticatedClassPrecursor struct {
	classPrecursorCandidate
	RegistryIndex int
}

func currentSINFPrivateSource(sinf []byte) ([]byte, error) {
	rootStart, rootEnd, err := rawSourceAtom(sinf, 0, len(sinf), "sinf")
	if err != nil || rootStart != 8 || rootEnd != len(sinf) {
		return nil, fmt.Errorf("SINF root is invalid")
	}
	schiStart, schiEnd, err := rawSourceAtom(sinf, rootStart, rootEnd, "schi")
	if err != nil {
		return nil, err
	}
	privateStart, privateEnd, err := rawSourceAtom(sinf, schiStart, schiEnd, "priv")
	if err != nil {
		return nil, err
	}
	return sinf[privateStart:privateEnd], nil
}

func validateCurrentSINFPrecursor(sinf []byte, profile *semanticBuildProfile, precursor [16]byte) (bool, error) {
	if profile == nil {
		return false, fmt.Errorf("semantic enrollment profile is unavailable")
	}
	source, err := currentSINFPrivateSource(sinf)
	if err != nil {
		return false, err
	}
	if len(source) < 48 {
		return false, fmt.Errorf("SINF private source is too short")
	}
	bodyBytes := len(source) - aes.BlockSize
	trailerBytes := bodyBytes % aes.BlockSize
	ciphertextBytes := bodyBytes - trailerBytes
	if ciphertextBytes < 2*aes.BlockSize || trailerBytes != 0 && trailerBytes != 8 {
		return false, fmt.Errorf("SINF private block geometry is invalid")
	}
	encodedBytes := aes.BlockSize + ciphertextBytes
	for _, value := range source[encodedBytes:] {
		if value != 0 {
			return false, fmt.Errorf("SINF private trailer is not zero")
		}
	}

	logical := make([]byte, encodedBytes)
	plaintext := make([]byte, ciphertextBytes+trailerBytes)
	decrypted := plaintext[:ciphertextBytes]
	defer func() {
		clear(logical)
		clear(plaintext)
	}()
	for index, value := range source[:encodedBytes] {
		position := index % aes.BlockSize
		external := value ^ profile.currentSINFWhite[position]
		logical[index] = profile.currentSINFInput[position*256+int(external)]
	}
	key := currentSINFModelKeyFromPrecursor(precursor)
	block, err := aes.NewCipher(key[:])
	clear(key[:])
	if err != nil {
		return false, err
	}
	cipher.NewCBCDecrypter(block, logical[:aes.BlockSize]).CryptBlocks(decrypted, logical[aes.BlockSize:])
	copy(plaintext[ciphertextBytes:], source[encodedBytes:])

	expected := [...]struct {
		typ   string
		sizes []int
	}{
		{"key ", []int{24}},
		{"iviv", []int{24}},
		{"tmpr", []int{24}},
		{"righ", []int{88, 96}},
		{"name", []int{264}},
	}
	cursor := 0
	for _, atom := range expected {
		if len(plaintext)-cursor < 8 {
			return false, nil
		}
		size := int(binary.BigEndian.Uint32(plaintext[cursor : cursor+4]))
		if size < 8 || size > len(plaintext)-cursor || string(plaintext[cursor+4:cursor+8]) != atom.typ {
			return false, nil
		}
		validSize := false
		for _, expectedSize := range atom.sizes {
			validSize = validSize || size == expectedSize
		}
		if !validSize {
			return false, nil
		}
		cursor += size
	}
	return cursor == len(plaintext), nil
}

func selectAuthenticatedClassPrecursor(sinf []byte, profile *semanticBuildProfile, registries []decodedSIDBRegistryCandidate) (authenticatedClassPrecursor, error) {
	var selected authenticatedClassPrecursor
	matches := 0
	for registryIndex, registry := range registries {
		candidates, err := deriveCurrentClassPrecursorCandidates(registry.Registry)
		if err != nil {
			return selected, err
		}
		for _, candidate := range candidates {
			accepted, err := validateCurrentSINFPrecursor(sinf, profile, candidate.Precursor)
			if err != nil {
				return selected, err
			}
			if !accepted {
				continue
			}
			matches++
			selected = authenticatedClassPrecursor{
				classPrecursorCandidate: candidate,
				RegistryIndex:           registryIndex,
			}
			if matches > 1 {
				return authenticatedClassPrecursor{}, fmt.Errorf("authenticated SINF accepted multiple class precursors")
			}
		}
	}
	if matches != 1 {
		return authenticatedClassPrecursor{}, fmt.Errorf("authenticated SINF accepted no class precursor")
	}
	return selected, nil
}
