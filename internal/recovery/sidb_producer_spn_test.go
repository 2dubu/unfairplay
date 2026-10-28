package recovery

import (
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"testing"
)

func TestSIDBProducerSPNOracleProfile(t *testing.T) {
	digest := sha256.Sum256(sidbProducerSPNProfileCompressed)
	if got := hex.EncodeToString(digest[:]); got != "7f31d8e3b36bcd2981ff2edd8d742d83616c925a9a5c744c868af8e576edc904" {
		t.Fatalf("SIDB producer SPN profile digest = %s", got)
	}
	if len(currentSIDBProducerSPN.stages) != 10 {
		t.Fatalf("SIDB producer SPN has %d stages", len(currentSIDBProducerSPN.stages))
	}
}

func TestSIDBProducerCanonicalAESProfile(t *testing.T) {
	digest := sha256.Sum256(sidbProducerCanonicalAESProfileCompressed)
	if got := hex.EncodeToString(digest[:]); got != "b46470baad541cc627391f46e4c5039d7d3835b972408c183ceaf9a84b0571e9" {
		t.Fatalf("SIDB producer canonical-AES profile digest = %s", got)
	}
	hash := sha256.New()
	for index := range 4096 {
		digest := sha256.Sum256([]byte{0x53, 0x49, 0x44, 0x42, 0x4b, byte(index >> 8), byte(index)})
		var input [16]byte
		copy(input[:], digest[:16])
		key := deriveSelectionSIDBKey(input)
		_, _ = hash.Write(key[:])
	}
	if got := hex.EncodeToString(hash.Sum(nil)); got != "98b71a19908213bbd6b6499182d0d569c9554b5f0dc5c9e631adc0cb3df9c3ab" {
		t.Fatalf("SIDB producer key corpus digest = %s", got)
	}
	stream, err := gzip.NewReader(bytes.NewReader(sidbProducerCanonicalAESProfileCompressed))
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := io.ReadAll(stream)
	if err != nil {
		t.Fatal(err)
	}
	if err := stream.Close(); err != nil {
		t.Fatal(err)
	}
	if len(decoded) != 1064 {
		t.Fatalf("decoded SIDB producer canonical-AES profile size = %d", len(decoded))
	}
	decodedDigest := sha256.Sum256(decoded)
	if got := hex.EncodeToString(decodedDigest[:]); got != "a7895ad2c5e00b8650a71f932d45f86f1473781cdd5d412de756933c16816985" {
		t.Fatalf("decoded SIDB producer canonical-AES profile digest = %s", got)
	}
	corrupt := append([]byte(nil), decoded...)
	corrupt[len(corrupt)-1] ^= 1
	if _, err := parseSIDBProducerCanonicalAESProfile(corrupt); err == nil {
		t.Fatal("corrupt SIDB producer canonical-AES profile was accepted")
	}
}
