// Package scinfo verifies signed FairPlay metadata envelopes.
package scinfo

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/x509"
	"fmt"
	"math/big"
	"strings"
)

const fairPlayCAModulusHex = "b2673c5d2ae78ff277cfef31bcd104fdea9fd08629900aa3ff644aef69f40ab9efc4c463ca3a2d213d528c238929aadb17dccb28c469cd7224c210c1f4781d3116c5e6ea191d2a33be80abdc1148062a1f0e6181cd799a5d12f48490c87cf6789b52f8faf1886e1072f85dcb88ecb6551497a44d0bc14f276799a6a4b3d7daad"

func AuthorityPublicKey() *rsa.PublicKey {
	modulus, ok := new(big.Int).SetString(fairPlayCAModulusHex, 16)
	if !ok || modulus.BitLen() != 1024 {
		panic("invalid embedded FairPlay CA public key")
	}
	return &rsa.PublicKey{N: modulus, E: 65537}
}

func verifyFairPlayLeafCertificate(certificate *x509.Certificate, authority *rsa.PublicKey) error {
	if certificate.SignatureAlgorithm != x509.SHA1WithRSA || len(certificate.Signature) != 128 {
		return fmt.Errorf("FairPlay certificate has unsupported authority signature")
	}
	if certificate.Issuer.CommonName != "Apple FairPlay Certification Authority" ||
		len(certificate.Issuer.Organization) != 1 || certificate.Issuer.Organization[0] != "Apple Inc." ||
		len(certificate.Issuer.OrganizationalUnit) != 1 || certificate.Issuer.OrganizationalUnit[0] != "Apple Certification Authority" {
		return fmt.Errorf("FairPlay certificate has an unexpected authority")
	}
	if len(certificate.Subject.Organization) != 1 || certificate.Subject.Organization[0] != "Apple Inc." ||
		len(certificate.Subject.OrganizationalUnit) != 1 || certificate.Subject.OrganizationalUnit[0] != "Apple FairPlay" ||
		(!strings.HasPrefix(certificate.Subject.CommonName, "SINF.") && !strings.HasPrefix(certificate.Subject.CommonName, "AP.")) || certificate.IsCA {
		return fmt.Errorf("FairPlay certificate has an unexpected subject")
	}
	digest := sha1.Sum(certificate.RawTBSCertificate)
	if err := rsa.VerifyPKCS1v15(authority, crypto.SHA1, digest[:], certificate.Signature); err != nil {
		return fmt.Errorf("verify FairPlay certificate authority signature: %w", err)
	}
	return nil
}

// VerifyEnvelope authenticates a certificate and the signed envelope prefix.
// A nil authority selects the pinned FairPlay certification authority.
func VerifyEnvelope(variant byte, signedPrefix, certificateBytes, signature []byte, authority *rsa.PublicKey) error {
	if authority == nil {
		authority = AuthorityPublicKey()
	}
	if authority.N == nil || authority.N.Sign() <= 0 || authority.E <= 1 {
		return fmt.Errorf("invalid FairPlay authority key")
	}
	certificate, err := x509.ParseCertificate(certificateBytes)
	if err != nil {
		return fmt.Errorf("parse FairPlay certificate: %w", err)
	}
	if err := verifyFairPlayLeafCertificate(certificate, authority); err != nil {
		return err
	}
	publicKey, ok := certificate.PublicKey.(*rsa.PublicKey)
	if !ok || publicKey.N == nil || publicKey.N.BitLen() != 1024 || publicKey.E != 65537 || len(signature) != 128 {
		return fmt.Errorf("FairPlay certificate has unsupported RSA key geometry")
	}
	var hash crypto.Hash
	var digest []byte
	switch variant {
	case 1, 3:
		value := sha1.Sum(signedPrefix)
		hash, digest = crypto.SHA1, value[:]
	case 4, 5:
		value := sha256.Sum256(signedPrefix)
		hash, digest = crypto.SHA256, value[:]
	default:
		return fmt.Errorf("FairPlay envelope variant %d is unsupported", variant)
	}
	if err := rsa.VerifyPKCS1v15(publicKey, hash, digest, signature); err != nil {
		return fmt.Errorf("verify FairPlay signed prefix: %w", err)
	}
	return nil
}
