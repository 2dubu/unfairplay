package recovery

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/DrWhax/unfairplay/internal/recovery/tagengine"
)

var errNoAuthenticatedClassPrecursor = errors.New("authenticated SINF accepted no class precursor")

type authenticatedClassPrecursor struct {
	classPrecursorCandidate
	RegistryIndex  int
	classKeyDigest [32]byte
}

func currentSINFPrivateSource(sinf []byte) ([]byte, error) {
	parsed, err := parseRecoverySINF(sinf)
	if err != nil {
		return nil, err
	}
	return parsed.private, nil
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
		typ  string
		size int
	}{
		{"key ", 24},
		{"iviv", 24},
		{"tmpr", 24},
		{"righ", 0},
		{"name", 264},
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
		if atom.typ == "righ" {
			if err := validateCurrentSINFRights(plaintext[cursor+8 : cursor+size]); err != nil {
				return false, nil
			}
		} else if size != atom.size {
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

func reportRecoveryProgress(progress func(Progress), event Progress) {
	if progress != nil {
		progress(event)
	}
}

// selectAuthenticatedEnrollmentClassPrecursor authenticates the selected tree route.
func selectAuthenticatedEnrollmentClassPrecursor(sinf []byte, registries []decodedSIDBRegistryCandidate, hardware []byte, user, keyID uint32) (authenticatedClassPrecursor, error) {
	return selectAuthenticatedEnrollmentClassPrecursorWithContext(context.Background(), sinf, registries, hardware, user, keyID, nil)
}

func selectAuthenticatedEnrollmentClassPrecursorWithContext(ctx context.Context, sinf []byte, registries []decodedSIDBRegistryCandidate, hardware []byte, user, keyID uint32, progress func(Progress)) (authenticatedClassPrecursor, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if err := ctx.Err(); err != nil {
		return authenticatedClassPrecursor{}, err
	}
	if len(registries) != 1 || !registries[0].Registry.Tree || len(registries[0].Data) == 0 {
		return authenticatedClassPrecursor{}, fmt.Errorf("tagged SIDB selection requires one decoded tree")
	}
	rootCursor, err := nativeSIDBRootCursor(hardware)
	if err != nil {
		return authenticatedClassPrecursor{}, err
	}
	book, err := tagengine.GenerateCodebookWithContext(ctx, hardware)
	if err != nil {
		return authenticatedClassPrecursor{}, fmt.Errorf("derive SIDB tag codebook: %w", err)
	}
	defer clear(book)
	codes := make(map[uint32]string, len(book))
	defer clear(codes)
	for _, entry := range book {
		if len(entry.Semantic) != 4 {
			return authenticatedClassPrecursor{}, fmt.Errorf("invalid SIDB tag codebook entry")
		}
		if _, exists := codes[entry.Raw]; exists {
			return authenticatedClassPrecursor{}, fmt.Errorf("duplicate SIDB tag codebook entry")
		}
		codes[entry.Raw] = entry.Semantic
	}
	decode := func(raw uint32) (string, bool) { tag, ok := codes[raw]; return tag, ok }
	return authenticateTaggedSelection(ctx, decode, sinf, registries, hardware, rootCursor, user, keyID, progress)
}
