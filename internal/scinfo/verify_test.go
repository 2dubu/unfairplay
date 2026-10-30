package scinfo

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"testing"
	"time"
)

func TestVerifyEnvelopeVariantsAndTrust(t *testing.T) {
	ca, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		t.Fatal(err)
	}
	leaf, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		t.Fatal(err)
	}
	issuer := &x509.Certificate{SerialNumber: big.NewInt(1), Subject: pkix.Name{CommonName: "Apple FairPlay Certification Authority", Organization: []string{"Apple Inc."}, OrganizationalUnit: []string{"Apple Certification Authority"}}, NotBefore: time.Unix(1, 0), NotAfter: time.Unix(4, 0), IsCA: true, BasicConstraintsValid: true, KeyUsage: x509.KeyUsageCertSign}
	subject := &x509.Certificate{SerialNumber: big.NewInt(2), Subject: pkix.Name{CommonName: "AP.synthetic", Organization: []string{"Apple Inc."}, OrganizationalUnit: []string{"Apple FairPlay"}}, NotBefore: time.Unix(1, 0), NotAfter: time.Unix(3, 0), SignatureAlgorithm: x509.SHA1WithRSA}
	cert, err := x509.CreateCertificate(rand.Reader, subject, issuer, &leaf.PublicKey, ca)
	if err != nil {
		t.Fatal(err)
	}
	prefix := []byte("synthetic signed prefix")
	for _, variant := range []byte{1, 3, 4, 5} {
		var hash crypto.Hash
		var digest []byte
		if variant == 1 || variant == 3 {
			v := sha1.Sum(prefix)
			hash, digest = crypto.SHA1, v[:]
		} else {
			v := sha256.Sum256(prefix)
			hash, digest = crypto.SHA256, v[:]
		}
		sig, err := rsa.SignPKCS1v15(rand.Reader, leaf, hash, digest)
		if err != nil {
			t.Fatal(err)
		}
		if err = VerifyEnvelope(variant, prefix, cert, sig, &ca.PublicKey); err != nil {
			t.Fatal(err)
		}
		if VerifyEnvelope(variant, append(append([]byte(nil), prefix...), 0), cert, sig, &ca.PublicKey) == nil {
			t.Fatal("accepted altered prefix")
		}
		if VerifyEnvelope(variant, prefix, cert, sig, nil) == nil {
			t.Fatal("accepted a synthetic certificate under the pinned authority")
		}
		if VerifyEnvelope(variant, prefix, cert, sig[:len(sig)-1], &ca.PublicKey) == nil {
			t.Fatal("accepted truncated signature")
		}
		if VerifyEnvelope(2, prefix, cert, sig, &ca.PublicKey) == nil {
			t.Fatal("accepted unsupported variant")
		}
	}
	if VerifyEnvelope(1, prefix, []byte{1}, make([]byte, 128), &ca.PublicKey) == nil {
		t.Fatal("accepted malformed certificate")
	}
}
