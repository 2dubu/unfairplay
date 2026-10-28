package recovery

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/x509"
	"encoding/binary"
	"fmt"
	"math/big"
	"strings"
)

const fairPlayCAModulusHex = "b2673c5d2ae78ff277cfef31bcd104fdea9fd08629900aa3ff644aef69f40ab9efc4c463ca3a2d213d528c238929aadb17dccb28c469cd7224c210c1f4781d3116c5e6ea191d2a33be80abdc1148062a1f0e6181cd799a5d12f48490c87cf6789b52f8faf1886e1072f85dcb88ecb6551497a44d0bc14f276799a6a4b3d7daad"

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

func fairPlayCAPublicKey() *rsa.PublicKey {
	modulus, ok := new(big.Int).SetString(fairPlayCAModulusHex, 16)
	if !ok || modulus.BitLen() != 1024 {
		panic("invalid embedded FairPlay CA public key")
	}
	return &rsa.PublicKey{N: modulus, E: 65537}
}

func verifyFairPlayLeafCertificate(certificate *x509.Certificate, authority *rsa.PublicKey) error {
	if certificate.SignatureAlgorithm != x509.SHA1WithRSA || len(certificate.Signature) != 128 {
		return fmt.Errorf("SUPF certificate has unsupported authority signature")
	}
	if certificate.Issuer.CommonName != "Apple FairPlay Certification Authority" ||
		len(certificate.Issuer.Organization) != 1 || certificate.Issuer.Organization[0] != "Apple Inc." ||
		len(certificate.Issuer.OrganizationalUnit) != 1 || certificate.Issuer.OrganizationalUnit[0] != "Apple Certification Authority" {
		return fmt.Errorf("SUPF certificate has an unexpected authority")
	}
	if len(certificate.Subject.Organization) != 1 || certificate.Subject.Organization[0] != "Apple Inc." ||
		len(certificate.Subject.OrganizationalUnit) != 1 || certificate.Subject.OrganizationalUnit[0] != "Apple FairPlay" ||
		(!strings.HasPrefix(certificate.Subject.CommonName, "SINF.") && !strings.HasPrefix(certificate.Subject.CommonName, "AP.")) || certificate.IsCA {
		return fmt.Errorf("SUPF certificate has an unexpected subject")
	}
	digest := sha1.Sum(certificate.RawTBSCertificate)
	if err := rsa.VerifyPKCS1v15(authority, crypto.SHA1, digest[:], certificate.Signature); err != nil {
		return fmt.Errorf("verify SUPF certificate authority signature: %w", err)
	}
	return nil
}

func verifyCurrentSUPFEnvelopeWithCA(envelope currentSUPFEnvelope, authority *rsa.PublicKey) error {
	certificate, err := x509.ParseCertificate(envelope.Certificate)
	if err != nil {
		return fmt.Errorf("parse SUPF certificate: %w", err)
	}
	if err := verifyFairPlayLeafCertificate(certificate, authority); err != nil {
		return err
	}
	publicKey, ok := certificate.PublicKey.(*rsa.PublicKey)
	if !ok || publicKey.N.BitLen() != 1024 || publicKey.E != 65537 || len(envelope.Signature) != 128 {
		return fmt.Errorf("SUPF certificate has unsupported RSA key geometry")
	}
	var hash crypto.Hash
	var digest []byte
	switch envelope.Variant {
	case 3:
		value := sha1.Sum(envelope.SignedPrefix)
		hash, digest = crypto.SHA1, value[:]
	case 5:
		value := sha256.Sum256(envelope.SignedPrefix)
		hash, digest = crypto.SHA256, value[:]
	default:
		return fmt.Errorf("SUPF variant %d is unsupported", envelope.Variant)
	}
	if err := rsa.VerifyPKCS1v15(publicKey, hash, digest, envelope.Signature); err != nil {
		return fmt.Errorf("verify SUPF signed prefix: %w", err)
	}
	return nil
}

func verifyCurrentSUPFEnvelope(envelope currentSUPFEnvelope) error {
	return verifyCurrentSUPFEnvelopeWithCA(envelope, fairPlayCAPublicKey())
}

func validateCurrentSUPFFilesWithCA(files map[string][]byte, authority *rsa.PublicKey) (int, error) {
	validated := 0
	for name, data := range files {
		if !strings.HasSuffix(name, ".supf") || len(data) == 0 {
			continue
		}
		envelope, err := parseCurrentSUPFEnvelope(data)
		if err != nil {
			return validated, fmt.Errorf("parse %s: %w", name, err)
		}
		if err := verifyCurrentSUPFEnvelopeWithCA(envelope, authority); err != nil {
			return validated, fmt.Errorf("validate %s: %w", name, err)
		}
		validated++
	}
	if validated == 0 {
		return 0, fmt.Errorf("recovery requires one signed SUPF")
	}
	return validated, nil
}

func validateCurrentSUPFFiles(files map[string][]byte) (int, error) {
	return validateCurrentSUPFFilesWithCA(files, fairPlayCAPublicKey())
}
