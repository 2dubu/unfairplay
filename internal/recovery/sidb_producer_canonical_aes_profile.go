package recovery

import (
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	_ "embed"
	"encoding/binary"
	"fmt"
	"io"
)

//go:embed profiles/current-sidb-producer-canonical-aes.bin.gz
var sidbProducerCanonicalAESProfileCompressed []byte

func loadCurrentSIDBProducerCanonicalAES() sidbProducerCanonicalAES {
	stream, err := gzip.NewReader(bytes.NewReader(sidbProducerCanonicalAESProfileCompressed))
	if err != nil {
		panic(fmt.Sprintf("decompress SIDB producer canonical-AES profile: %v", err))
	}
	encoded, err := io.ReadAll(stream)
	closeErr := stream.Close()
	if err != nil || closeErr != nil {
		panic(fmt.Sprintf("decompress SIDB producer canonical-AES profile: %v", err))
	}
	model, err := parseSIDBProducerCanonicalAESProfile(encoded)
	clear(encoded)
	if err != nil {
		panic(fmt.Sprintf("load SIDB producer canonical-AES profile: %v", err))
	}
	return model
}

var currentSIDBProducerCanonicalAES = loadCurrentSIDBProducerCanonicalAES()

func readSIDBSparseAffine(reader *bytes.Reader, label string, expectedTerms int, perSource func(int) bool) (sidbSparseAffine, error) {
	var value sidbSparseAffine
	if _, err := io.ReadFull(reader, value.constant[:]); err != nil {
		return value, err
	}
	total := 0
	for source := range value.terms {
		count, err := reader.ReadByte()
		if err != nil || !perSource(int(count)) {
			return value, fmt.Errorf("%s source %d has invalid term count %d", label, source, count)
		}
		var seen [16]bool
		value.terms[source] = make([]sidbSparseTerm, count)
		for index := range value.terms[source] {
			term := &value.terms[source][index]
			target, err := reader.ReadByte()
			if err != nil || target >= 16 || seen[target] {
				return value, fmt.Errorf("%s source %d has invalid target %d", label, source, target)
			}
			term.target, seen[target] = target, true
			if _, err := io.ReadFull(reader, term.linear[:]); err != nil {
				return value, err
			}
			if term.linear == (frameByteLinear{}) {
				return value, fmt.Errorf("%s source %d target %d has zero map", label, source, target)
			}
		}
		total += int(count)
	}
	if total != expectedTerms {
		return value, fmt.Errorf("%s has %d terms, want %d", label, total, expectedTerms)
	}
	if frameAffine128Rank(sidbSparseAffineToFull(value)) != 128 {
		return value, fmt.Errorf("%s is singular", label)
	}
	return value, nil
}

func readSIDBBytePermutationAffine(reader *bytes.Reader, label string) (sidbBytePermutationAffine, error) {
	var result sidbBytePermutationAffine
	if _, err := io.ReadFull(reader, result.constant[:]); err != nil {
		return result, err
	}
	var seen [16]bool
	for source := range result.terms {
		term := &result.terms[source]
		target, err := reader.ReadByte()
		if err != nil || target >= 16 || seen[target] {
			return result, fmt.Errorf("%s source %d has invalid target %d", label, source, target)
		}
		term.target, seen[target] = target, true
		if _, err := io.ReadFull(reader, term.linear[:]); err != nil {
			return result, err
		}
		if !term.linear.invertible() {
			return result, fmt.Errorf("%s source %d is singular", label, source)
		}
	}
	return result, nil
}

func parseSIDBProducerCanonicalAESProfile(encoded []byte) (sidbProducerCanonicalAES, error) {
	var result sidbProducerCanonicalAES
	if len(encoded) < 32 {
		return result, fmt.Errorf("SIDB producer canonical-AES profile is truncated")
	}
	payload, checksum := encoded[:len(encoded)-32], encoded[len(encoded)-32:]
	digest := sha256.Sum256(payload)
	if !bytes.Equal(checksum, digest[:]) {
		return result, fmt.Errorf("SIDB producer canonical-AES profile checksum mismatch")
	}
	reader := bytes.NewReader(payload)
	var magic [8]byte
	if _, err := io.ReadFull(reader, magic[:]); err != nil || string(magic[:]) != "FPSCA1\x00\x00" {
		return result, fmt.Errorf("invalid SIDB producer canonical-AES profile magic")
	}
	var version uint32
	if err := binary.Read(reader, binary.LittleEndian, &version); err != nil || version != 1 {
		return result, fmt.Errorf("unsupported SIDB producer canonical-AES profile version %d", version)
	}
	var err error
	result.first, err = readSIDBSparseAffine(reader, "first layer", 16, func(count int) bool { return count == 1 })
	if err != nil {
		return result, err
	}
	result.firstPost, err = readSIDBBytePermutationAffine(reader, "first post-map")
	if err != nil {
		return result, err
	}
	keyCount, err := reader.ReadByte()
	if err != nil || keyCount != byte(len(result.roundKeys)) {
		return result, fmt.Errorf("invalid SIDB producer canonical-AES key count %d", keyCount)
	}
	var reserved [3]byte
	if _, err := io.ReadFull(reader, reserved[:]); err != nil || reserved != [3]byte{} {
		return result, fmt.Errorf("invalid SIDB producer canonical-AES key header")
	}
	for round := range result.roundKeys {
		if _, err := io.ReadFull(reader, result.roundKeys[round][:]); err != nil {
			return result, err
		}
	}
	result.lastPre, err = readSIDBBytePermutationAffine(reader, "last pre-map")
	if err != nil {
		return result, err
	}
	result.final, err = readSIDBSparseAffine(reader, "final layer", 40, func(count int) bool { return count == 2 || count == 3 })
	if err != nil {
		return result, err
	}
	if reader.Len() != 0 {
		return result, fmt.Errorf("SIDB producer canonical-AES profile has trailing data")
	}
	return result, nil
}
