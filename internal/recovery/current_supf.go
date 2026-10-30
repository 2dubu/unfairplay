package recovery

import (
	"crypto/rsa"
	"encoding/binary"
	"fmt"
	"strings"

	"github.com/DrWhax/unfairplay/internal/scinfo"
)

type currentSUPFEnvelope struct {
	Variant      byte
	SignedPrefix []byte
	Certificate  []byte
	Signature    []byte
}

func parseCurrentSUPFEnvelope(data []byte) (currentSUPFEnvelope, error) {
	var result currentSUPFEnvelope
	if len(data) < 12 {
		return result, fmt.Errorf("SUPF is truncated")
	}
	result.Variant = data[0]
	if result.Variant == 1 || result.Variant == 4 {
		cursor := 4
		if result.Variant == 1 {
			cursor += 20
		}
		if len(data)-cursor < 4 {
			return result, fmt.Errorf("SUPP header is truncated")
		}
		pages := uint64(binary.BigEndian.Uint32(data[cursor : cursor+4]))
		cursor += 4
		if pages == 0 || len(data)-cursor < 4+128 || pages > uint64((len(data)-cursor-4-128)/32) {
			return result, fmt.Errorf("SUPP page records have invalid bounds")
		}
		end := cursor + int(pages)*32
		certificateBytes := uint64(binary.BigEndian.Uint32(data[end : end+4]))
		certificateStart := end + 4
		if certificateBytes == 0 || certificateBytes != uint64(len(data)-certificateStart-128) {
			return result, fmt.Errorf("SUPP certificate and signature have invalid bounds")
		}
		certificateEnd := certificateStart + int(certificateBytes)
		result.SignedPrefix = data[:end]
		result.Certificate = data[certificateStart:certificateEnd]
		result.Signature = data[certificateEnd:]
		return result, nil
	}
	if result.Variant != 3 && result.Variant != 5 {
		return result, fmt.Errorf("SUPF variant %d is unsupported", result.Variant)
	}
	outerEnd := 8 + int(binary.BigEndian.Uint32(data[4:8]))
	if outerEnd < 12 || outerEnd > len(data) {
		return result, fmt.Errorf("SUPF outer payload has invalid bounds")
	}
	architectureCount := int(binary.BigEndian.Uint32(data[8:12]))
	if architectureCount < 1 || architectureCount > 32 {
		return result, fmt.Errorf("SUPF architecture count is invalid")
	}
	cursor := 12
	for architecture := 0; architecture < architectureCount; architecture++ {
		if outerEnd-cursor < 4 {
			return result, fmt.Errorf("SUPF architecture %d length is truncated", architecture)
		}
		length := int(binary.BigEndian.Uint32(data[cursor : cursor+4]))
		start, end := cursor+4, cursor+4+length
		minimum := 12
		if result.Variant == 3 {
			minimum += 20
		}
		if length < minimum || end > outerEnd {
			return result, fmt.Errorf("SUPF architecture %d has invalid bounds", architecture)
		}
		pageCountOffset := start + 8
		if result.Variant == 3 {
			pageCountOffset += 20
		}
		pageCount := uint64(binary.BigEndian.Uint32(data[pageCountOffset : pageCountOffset+4]))
		keyBytes := pageCount * 32
		if keyBytes > uint64(end-(pageCountOffset+4)) || pageCountOffset+4+int(keyBytes) != end {
			return result, fmt.Errorf("SUPF architecture %d page-key geometry is invalid", architecture)
		}
		cursor = end
	}
	if cursor != outerEnd || len(data)-cursor < 4 {
		return result, fmt.Errorf("SUPF architectures do not fill the signed prefix")
	}
	certificateBytes := int(binary.BigEndian.Uint32(data[cursor : cursor+4]))
	certificateStart, certificateEnd := cursor+4, cursor+4+certificateBytes
	if certificateBytes < 1 || certificateEnd > len(data) || len(data)-certificateEnd < 4 {
		return result, fmt.Errorf("SUPF certificate has invalid bounds")
	}
	signatureBytes := int(binary.BigEndian.Uint32(data[certificateEnd : certificateEnd+4]))
	signatureStart, signatureEnd := certificateEnd+4, certificateEnd+4+signatureBytes
	if signatureBytes < 1 || signatureEnd != len(data) {
		return result, fmt.Errorf("SUPF signature has invalid bounds")
	}
	result.SignedPrefix = data[:outerEnd]
	result.Certificate = data[certificateStart:certificateEnd]
	result.Signature = data[signatureStart:signatureEnd]
	return result, nil
}

func fairPlayCAPublicKey() *rsa.PublicKey { return scinfo.AuthorityPublicKey() }

func verifyCurrentSUPFEnvelopeWithCA(envelope currentSUPFEnvelope, authority *rsa.PublicKey) error {
	if envelope.Variant != 1 && envelope.Variant != 3 && envelope.Variant != 4 && envelope.Variant != 5 {
		return fmt.Errorf("SUPF variant %d is unsupported", envelope.Variant)
	}
	return scinfo.VerifyEnvelope(envelope.Variant, envelope.SignedPrefix, envelope.Certificate, envelope.Signature, authority)
}

func verifyCurrentSUPFEnvelope(envelope currentSUPFEnvelope) error {
	return verifyCurrentSUPFEnvelopeWithCA(envelope, fairPlayCAPublicKey())
}

func validateCurrentSUPFFilesWithCA(files map[string][]byte, authority *rsa.PublicKey) (int, error) {
	validated := 0
	for name, data := range files {
		if (!strings.HasSuffix(name, ".supf") && !strings.HasSuffix(name, ".supp")) || len(data) == 0 {
			continue
		}
		envelope, err := parseCurrentSUPFEnvelope(data)
		if err != nil {
			return validated, fmt.Errorf("parse %s: %w", name, err)
		}
		pClass := strings.HasSuffix(name, ".supp")
		if pClass != (envelope.Variant == 1 || envelope.Variant == 4) {
			return validated, fmt.Errorf("keyfile %s has a variant inconsistent with its class", name)
		}
		if err := verifyCurrentSUPFEnvelopeWithCA(envelope, authority); err != nil {
			return validated, fmt.Errorf("validate %s: %w", name, err)
		}
		validated++
	}
	if validated == 0 {
		return 0, fmt.Errorf("recovery requires one signed keyfile")
	}
	return validated, nil
}

func validateCurrentSUPFFiles(files map[string][]byte) (int, error) {
	return validateCurrentSUPFFilesWithCA(files, fairPlayCAPublicKey())
}
