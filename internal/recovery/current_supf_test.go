package recovery

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/binary"
	"math/big"
	"testing"
	"time"
)

func syntheticSignedSUPF(t *testing.T, variant byte) ([]byte, *rsa.PublicKey) {
	t.Helper()
	authorityKey, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		t.Fatal(err)
	}
	authorityTemplate := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Country:            []string{"US"},
			Organization:       []string{"Apple Inc."},
			OrganizationalUnit: []string{"Apple Certification Authority"},
			CommonName:         "Apple FairPlay Certification Authority",
		},
		NotBefore: time.Unix(1, 0), NotAfter: time.Unix(4, 0),
		IsCA: true, BasicConstraintsValid: true, KeyUsage: x509.KeyUsageCertSign,
	}
	privateKey, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		t.Fatal(err)
	}
	certificateTemplate := &x509.Certificate{
		SerialNumber: big.NewInt(2),
		Subject: pkix.Name{
			Country:            []string{"US"},
			Organization:       []string{"Apple Inc."},
			OrganizationalUnit: []string{"Apple FairPlay"},
			CommonName:         "SINF.synthetic-test",
		},
		NotBefore: time.Unix(1, 0), NotAfter: time.Unix(3, 0),
		BasicConstraintsValid: true, KeyUsage: x509.KeyUsageDigitalSignature,
		SignatureAlgorithm: x509.SHA1WithRSA,
	}
	certificate, err := x509.CreateCertificate(rand.Reader, certificateTemplate, authorityTemplate, &privateKey.PublicKey, authorityKey)
	if err != nil {
		t.Fatal(err)
	}

	prefixBytes := 60
	segmentBytes := uint32(44)
	pageCountOffset := 24
	if variant == 3 {
		prefixBytes = 80
		segmentBytes = 64
		pageCountOffset = 44
	}
	prefix := make([]byte, prefixBytes)
	prefix[0] = variant
	copy(prefix[1:4], "612")
	binary.BigEndian.PutUint32(prefix[4:8], uint32(len(prefix)-8))
	binary.BigEndian.PutUint32(prefix[8:12], 1)
	binary.BigEndian.PutUint32(prefix[12:16], segmentBytes)
	binary.BigEndian.PutUint32(prefix[16:20], 0x0100000c)
	binary.BigEndian.PutUint32(prefix[20:24], 0)
	if variant == 3 {
		for index := 24; index < 44; index++ {
			prefix[index] = byte(index*11 + 3)
		}
	}
	binary.BigEndian.PutUint32(prefix[pageCountOffset:pageCountOffset+4], 1)
	for index := pageCountOffset + 4; index < len(prefix); index++ {
		prefix[index] = byte(index*29 + 7)
	}
	var hash crypto.Hash
	var digest []byte
	if variant == 3 {
		value := sha1.Sum(prefix)
		hash, digest = crypto.SHA1, value[:]
	} else {
		value := sha256.Sum256(prefix)
		hash, digest = crypto.SHA256, value[:]
	}
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, hash, digest)
	if err != nil {
		t.Fatal(err)
	}
	result := append([]byte(nil), prefix...)
	var word [4]byte
	binary.BigEndian.PutUint32(word[:], uint32(len(certificate)))
	result = append(result, word[:]...)
	result = append(result, certificate...)
	binary.BigEndian.PutUint32(word[:], uint32(len(signature)))
	result = append(result, word[:]...)
	result = append(result, signature...)
	return result, &authorityKey.PublicKey
}

func TestParseAndVerifyCurrentSUPFEnvelope(t *testing.T) {
	for _, variant := range []byte{3, 5} {
		data, authority := syntheticSignedSUPF(t, variant)
		envelope, err := parseCurrentSUPFEnvelope(data)
		if err != nil {
			t.Fatal(err)
		}
		if envelope.Variant != variant || len(envelope.Signature) != 128 {
			t.Fatalf("unexpected SUPF envelope geometry")
		}
		if err := verifyCurrentSUPFEnvelopeWithCA(envelope, authority); err != nil {
			t.Fatal(err)
		}
		wrongAuthority, err := rsa.GenerateKey(rand.Reader, 1024)
		if err != nil {
			t.Fatal(err)
		}
		if err := verifyCurrentSUPFEnvelopeWithCA(envelope, &wrongAuthority.PublicKey); err == nil {
			t.Fatal("SUPF certificate signed by another authority was accepted")
		}
		mutated := append([]byte(nil), data...)
		mutated[len(envelope.SignedPrefix)-1] ^= 1
		envelope, err = parseCurrentSUPFEnvelope(mutated)
		if err != nil {
			t.Fatal(err)
		}
		if err := verifyCurrentSUPFEnvelopeWithCA(envelope, authority); err == nil {
			t.Fatal("mutated SUPF signed prefix was accepted")
		}
	}
}

func TestValidateCurrentSUPFFiles(t *testing.T) {
	data, authority := syntheticSignedSUPF(t, 5)
	if count, err := validateCurrentSUPFFilesWithCA(map[string][]byte{"App.v5.supf": data, "App.sinf": {1, 2, 3}}, authority); err != nil || count != 1 {
		t.Fatalf("validated SUPFs = %d, err=%v", count, err)
	}
	if _, err := validateCurrentSUPFFilesWithCA(map[string][]byte{"App.sinf": {1, 2, 3}}, authority); err == nil {
		t.Fatal("missing signed SUPF was accepted")
	}
}

func TestParseCurrentSUPFEnvelopeRejectsGeometry(t *testing.T) {
	data, _ := syntheticSignedSUPF(t, 5)
	for _, mutated := range [][]byte{data[:11], data[:len(data)-1]} {
		if _, err := parseCurrentSUPFEnvelope(mutated); err == nil {
			t.Fatal("invalid SUPF geometry was accepted")
		}
	}
}

func syntheticSignedSUPP(t *testing.T, variant byte) ([]byte, *rsa.PublicKey) {
	t.Helper()
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
	prefix := make([]byte, 4)
	prefix[0] = variant
	if variant == 1 {
		prefix = append(prefix, make([]byte, 20)...)
	}
	var word [4]byte
	binary.BigEndian.PutUint32(word[:], 2)
	prefix = append(prefix, word[:]...)
	prefix = append(prefix, make([]byte, 64)...)
	var h crypto.Hash
	var digest []byte
	if variant == 1 {
		v := sha1.Sum(prefix)
		h, digest = crypto.SHA1, v[:]
	} else {
		v := sha256.Sum256(prefix)
		h, digest = crypto.SHA256, v[:]
	}
	signature, err := rsa.SignPKCS1v15(rand.Reader, leaf, h, digest)
	if err != nil {
		t.Fatal(err)
	}
	result := append([]byte(nil), prefix...)
	binary.BigEndian.PutUint32(word[:], uint32(len(cert)))
	result = append(result, word[:]...)
	result = append(result, cert...)
	result = append(result, signature...)
	return result, &ca.PublicKey
}

func TestCurrentRecoverySUPPEnvelope(t *testing.T) {
	for _, variant := range []byte{1, 4} {
		t.Run(string(rune('0'+variant)), func(t *testing.T) {
			data, authority := syntheticSignedSUPP(t, variant)
			envelope, err := parseCurrentSUPFEnvelope(data)
			if err != nil {
				t.Fatal(err)
			}
			if envelope.Variant != variant || len(envelope.Signature) != 128 {
				t.Fatal("invalid envelope")
			}
			if err := verifyCurrentSUPFEnvelopeWithCA(envelope, authority); err != nil {
				t.Fatal(err)
			}
			if count, err := validateCurrentSUPFFilesWithCA(map[string][]byte{"App.supp": data}, authority); err != nil || count != 1 {
				t.Fatalf("SUPP-only validation failed: %d, %v", count, err)
			}
			if _, err := validateCurrentSUPFFilesWithCA(map[string][]byte{"App.supf": data}, authority); err == nil {
				t.Fatal("accepted p-class envelope as an f-class file")
			}
			changed := append([]byte(nil), data...)
			changed[len(envelope.SignedPrefix)-1] ^= 1
			env, err := parseCurrentSUPFEnvelope(changed)
			if err != nil {
				t.Fatal(err)
			}
			if verifyCurrentSUPFEnvelopeWithCA(env, authority) == nil {
				t.Fatal("accepted changed signed payload")
			}
			for _, bad := range [][]byte{data[:len(data)-1], append(append([]byte(nil), data...), 0), data[:12]} {
				if _, err := parseCurrentSUPFEnvelope(bad); err == nil {
					t.Fatal("accepted malformed SUPP")
				}
			}
			malformed := append([]byte(nil), data...)
			pos := 4
			if variant == 1 {
				pos += 20
			}
			binary.BigEndian.PutUint32(malformed[pos:], ^uint32(0))
			if _, err := parseCurrentSUPFEnvelope(malformed); err == nil {
				t.Fatal("accepted invalid page count")
			}
		})
	}
	if !isCurrentRecoveryFile("App.supp") || !isCurrentRecoveryFile("App.v4.supp") {
		t.Fatal("SUPP omitted from recovery inputs")
	}
}

func TestCurrentRecoveryRejectsFClassInSUPP(t *testing.T) {
	data, authority := syntheticSignedSUPF(t, 3)
	if _, err := validateCurrentSUPFFilesWithCA(map[string][]byte{"App.supp": data}, authority); err == nil {
		t.Fatal("accepted f-class envelope as p-class file")
	}
}
