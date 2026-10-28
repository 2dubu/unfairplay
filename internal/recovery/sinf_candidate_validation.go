package recovery

import (
	"bytes"
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
	rootStart, rootEnd, err := enrollmentAtom(sinf, 0, len(sinf), "sinf")
	if err != nil || rootStart != 8 || rootEnd != len(sinf) {
		return nil, fmt.Errorf("SINF root is invalid")
	}
	schiStart, schiEnd, err := enrollmentAtom(sinf, rootStart, rootEnd, "schi")
	if err != nil {
		return nil, err
	}
	privateStart, privateEnd, err := enrollmentAtom(sinf, schiStart, schiEnd, "priv")
	if err != nil {
		return nil, err
	}
	return sinf[privateStart:privateEnd], nil
}

func validateCurrentSINFPrecursor(sinf []byte, precursor [16]byte) (bool, error) {
	logical, ciphertextBytes, trailerBytes, err := currentSINFLogicalSource(sinf)
	if err != nil {
		return false, err
	}
	defer clear(logical)

	plaintext := make([]byte, ciphertextBytes+trailerBytes)
	decrypted := plaintext[:ciphertextBytes]
	defer clear(plaintext)
	key := currentSINFModelKeyFromPrecursor(precursor)
	block, err := aes.NewCipher(key[:])
	clear(key[:])
	if err != nil {
		return false, err
	}
	cipher.NewCBCDecrypter(block, logical[:aes.BlockSize]).CryptBlocks(decrypted, logical[aes.BlockSize:])

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

func currentSINFLogicalSource(sinf []byte) ([]byte, int, int, error) {
	source, err := currentSINFPrivateSource(sinf)
	if err != nil {
		return nil, 0, 0, err
	}
	if len(source) < 48 {
		return nil, 0, 0, fmt.Errorf("SINF private source is too short")
	}
	bodyBytes := len(source) - aes.BlockSize
	trailerBytes := bodyBytes % aes.BlockSize
	ciphertextBytes := bodyBytes - trailerBytes
	if ciphertextBytes < 2*aes.BlockSize || trailerBytes != 0 && trailerBytes != 8 {
		return nil, 0, 0, fmt.Errorf("SINF private block geometry is invalid")
	}
	encodedBytes := aes.BlockSize + ciphertextBytes
	for _, value := range source[encodedBytes:] {
		if value != 0 {
			return nil, 0, 0, fmt.Errorf("SINF private trailer is not zero")
		}
	}

	logical := append([]byte(nil), source[:encodedBytes]...)
	return logical, ciphertextBytes, trailerBytes, nil
}

type currentSINFModelKeyProbe struct {
	iv         [aes.BlockSize]byte
	ciphertext [aes.BlockSize]byte
}

func newCurrentSINFModelKeyProbe(sinf []byte) (currentSINFModelKeyProbe, error) {
	var result currentSINFModelKeyProbe
	logical, _, _, err := currentSINFLogicalSource(sinf)
	if err != nil {
		return result, err
	}
	defer clear(logical)
	copy(result.iv[:], logical[:aes.BlockSize])
	copy(result.ciphertext[:], logical[aes.BlockSize:2*aes.BlockSize])
	return result, nil
}

func (probe currentSINFModelKeyProbe) accepts(modelKey [aes.BlockSize]byte) (bool, error) {
	block, err := aes.NewCipher(modelKey[:])
	if err != nil {
		return false, err
	}
	var first [aes.BlockSize]byte
	block.Decrypt(first[:], probe.ciphertext[:])
	for index := range first {
		first[index] ^= probe.iv[index]
	}
	return bytes.Equal(first[:8], []byte{0, 0, 0, 24, 'k', 'e', 'y', ' '}), nil
}

func authenticateClassValu(probe currentSINFModelKeyProbe, sinf []byte, block cipher.Block, valu [aes.BlockSize]byte) ([aes.BlockSize]byte, bool, error) {
	var modelKey [aes.BlockSize]byte
	block.Decrypt(modelKey[:], valu[:])
	possible, err := probe.accepts(modelKey)
	if err != nil || !possible {
		clear(modelKey[:])
		return [aes.BlockSize]byte{}, false, err
	}
	var precursor [aes.BlockSize]byte
	for index := range precursor {
		precursor[index] = applyCurrentFairPlayByteMask(modelKey[index])
	}
	clear(modelKey[:])
	accepted, err := validateCurrentSINFPrecursor(sinf, precursor)
	if err != nil || !accepted {
		clear(precursor[:])
		return [aes.BlockSize]byte{}, false, err
	}
	return precursor, true, nil
}

func selectAuthenticatedClassPrecursor(sinf []byte, registries []decodedSIDBRegistryCandidate) (authenticatedClassPrecursor, error) {
	var selected authenticatedClassPrecursor
	probe, err := newCurrentSINFModelKeyProbe(sinf)
	if err != nil {
		return selected, err
	}
	type preparedClassCandidate struct {
		registryIndex int
		recordIndex   int
		windowOffset  int
		valu          [aes.BlockSize]byte
	}
	var candidates []preparedClassCandidate
	for registryIndex, registry := range registries {
		for recordIndex, record := range registry.Registry.Records {
			for windowOffset := 0; windowOffset <= maximumClassWindowOffset; windowOffset++ {
				valu, err := deriveClassValuWindow(record, windowOffset)
				if err != nil {
					return selected, err
				}
				candidates = append(candidates, preparedClassCandidate{
					registryIndex: registryIndex,
					recordIndex:   recordIndex,
					windowOffset:  windowOffset,
					valu:          valu,
				})
			}
		}
	}
	classValues := make([][8]byte, 0, classWorkspaceModulus+1)
	classValues = append(classValues, currentClassSecondaryValue)
	seenClassValues := map[[8]byte]bool{currentClassSecondaryValue: true}
	for seed := range uint32(classWorkspaceModulus) {
		value := classWorkspaceOutputFromSeed(uint16(seed))
		if seenClassValues[value] {
			continue
		}
		seenClassValues[value] = true
		classValues = append(classValues, value)
	}
	matches := 0
	for _, classValue := range classValues {
		parameters := classEnrollmentParameters{ClassValue: classValue}
		classKey := parameters.classKey()
		classBlock, err := aes.NewCipher(classKey[:])
		clear(classKey[:])
		if err != nil {
			return selected, err
		}
		for _, candidate := range candidates {
			precursor, accepted, err := authenticateClassValu(probe, sinf, classBlock, candidate.valu)
			if err != nil {
				return selected, err
			}
			if !accepted {
				continue
			}
			matches++
			selected = authenticatedClassPrecursor{
				classPrecursorCandidate: classPrecursorCandidate{
					RecordIndex:  candidate.recordIndex,
					WindowOffset: candidate.windowOffset,
					Record:       registries[candidate.registryIndex].Registry.Records[candidate.recordIndex],
					Precursor:    precursor,
				},
				RegistryIndex: candidate.registryIndex,
			}
			if matches > 1 {
				clear(selected.Precursor[:])
				return authenticatedClassPrecursor{}, fmt.Errorf("authenticated SINF accepted multiple class precursors")
			}
		}
	}
	if matches != 1 {
		return authenticatedClassPrecursor{}, fmt.Errorf("authenticated SINF accepted no class precursor")
	}
	return selected, nil
}

func selectAuthenticatedLegacyClassPrecursor(sinf []byte, registries []decodedSIDBRegistryCandidate, hardware []byte, user, keyID uint32) (authenticatedClassPrecursor, error) {
	var selected authenticatedClassPrecursor
	digest, err := legacyAccountDigest(hardware, user, keyID)
	if err != nil {
		return selected, err
	}
	parameters := classEnrollmentParameters{ClassValue: legacyClassWorkspaceValue(digest)}
	classKey := parameters.classKey()
	classBlock, err := aes.NewCipher(classKey[:])
	clear(classKey[:])
	clear(digest[:])
	if err != nil {
		return selected, err
	}
	probe, err := newCurrentSINFModelKeyProbe(sinf)
	if err != nil {
		return selected, err
	}
	matches := 0
	for registryIndex, registry := range registries {
		for recordIndex, record := range registry.Registry.Records {
			for maskPhase := range len(currentClassRecordMask) {
				valu, err := deriveLegacyClassValu(record, maskPhase)
				if err != nil {
					return selected, err
				}
				precursor, accepted, err := authenticateClassValu(probe, sinf, classBlock, valu)
				if err != nil {
					return selected, err
				}
				if !accepted {
					continue
				}
				matches++
				selected = authenticatedClassPrecursor{
					classPrecursorCandidate: classPrecursorCandidate{
						RecordIndex:  recordIndex,
						WindowOffset: maskPhase,
						Record:       record,
						Precursor:    precursor,
					},
					RegistryIndex: registryIndex,
				}
				if matches > 1 {
					clear(selected.Precursor[:])
					return authenticatedClassPrecursor{}, fmt.Errorf("authenticated SINF accepted multiple legacy class precursors")
				}
			}
		}
	}
	if matches != 1 {
		return authenticatedClassPrecursor{}, fmt.Errorf("authenticated SINF accepted no legacy class precursor")
	}
	return selected, nil
}

func selectAuthenticatedEnrollmentClassPrecursor(sinf []byte, registries []decodedSIDBRegistryCandidate, hardware []byte, user, keyID uint32) (authenticatedClassPrecursor, error) {
	if len(registries) == 0 {
		return authenticatedClassPrecursor{}, fmt.Errorf("class precursor selection requires one decoded registry")
	}
	older := len(registries[0].Data) == selectionSIDBOlderPayloadBytes
	for _, registry := range registries[1:] {
		if (len(registry.Data) == selectionSIDBOlderPayloadBytes) != older {
			return authenticatedClassPrecursor{}, fmt.Errorf("class precursor selection found mixed registry layouts")
		}
	}
	if older {
		return selectAuthenticatedLegacyClassPrecursor(sinf, registries, hardware, user, keyID)
	}
	return selectAuthenticatedClassPrecursor(sinf, registries)
}
