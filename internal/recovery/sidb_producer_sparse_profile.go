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

//go:embed profiles/current-sidb-producer-sparse.bin.gz
var sidbProducerSparseProfileCompressed []byte

func loadCurrentSIDBProducerSparse() sidbProducerSparse {
	stream, err := gzip.NewReader(bytes.NewReader(sidbProducerSparseProfileCompressed))
	if err != nil {
		panic(fmt.Sprintf("decompress SIDB producer sparse profile: %v", err))
	}
	encoded, err := io.ReadAll(stream)
	closeErr := stream.Close()
	if err != nil || closeErr != nil {
		panic(fmt.Sprintf("decompress SIDB producer sparse profile: %v", err))
	}
	model, err := parseSIDBProducerSparseProfile(encoded)
	clear(encoded)
	if err != nil {
		panic(fmt.Sprintf("load SIDB producer sparse profile: %v", err))
	}
	return model
}

var currentSIDBProducerSparse = loadCurrentSIDBProducerSparse()

func parseSIDBProducerSparseProfile(encoded []byte) (sidbProducerSparse, error) {
	var result sidbProducerSparse
	if len(encoded) < 32 {
		return result, fmt.Errorf("SIDB producer sparse profile is truncated")
	}
	payload, checksum := encoded[:len(encoded)-32], encoded[len(encoded)-32:]
	digest := sha256.Sum256(payload)
	if !bytes.Equal(checksum, digest[:]) {
		return result, fmt.Errorf("SIDB producer sparse profile checksum mismatch")
	}
	reader := bytes.NewReader(payload)
	var magic [8]byte
	if _, err := io.ReadFull(reader, magic[:]); err != nil || string(magic[:]) != "FPSPA1\x00\x00" {
		return result, fmt.Errorf("invalid SIDB producer sparse profile magic")
	}
	var version uint32
	if err := binary.Read(reader, binary.LittleEndian, &version); err != nil || version != 1 {
		return result, fmt.Errorf("unsupported SIDB producer sparse profile version %d", version)
	}
	readAffine := func(label string, expectedTerms int, perSource func(int) bool) (sidbSparseAffine, error) {
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
	var err error
	result.first, err = readAffine("first layer", 16, func(count int) bool { return count == 1 })
	if err != nil {
		return result, err
	}
	var firstTargets [16]bool
	for source := range result.first.terms {
		target := result.first.terms[source][0].target
		if firstTargets[target] {
			return result, fmt.Errorf("first layer target %d is repeated", target)
		}
		firstTargets[target] = true
	}
	roundCount, err := reader.ReadByte()
	if err != nil || roundCount != byte(len(result.rounds)) {
		return result, fmt.Errorf("invalid SIDB producer sparse round count %d", roundCount)
	}
	var reserved [3]byte
	if _, err := io.ReadFull(reader, reserved[:]); err != nil || reserved != [3]byte{} {
		return result, fmt.Errorf("invalid SIDB producer sparse header")
	}
	for index := range result.rounds {
		result.rounds[index], err = readAffine(fmt.Sprintf("round %d", index), 64, func(count int) bool { return count == 4 })
		if err != nil {
			return result, err
		}
	}
	result.final, err = readAffine("final layer", 40, func(count int) bool { return count == 2 || count == 3 })
	if err != nil {
		return result, err
	}
	if reader.Len() != 0 {
		return result, fmt.Errorf("SIDB producer sparse profile has trailing data")
	}
	return result, nil
}

func sidbSparseAffineToFull(value sidbSparseAffine) (result frameAffine128) {
	result.constant = value.constant
	for source, terms := range value.terms {
		for _, term := range terms {
			for bit, mapped := range term.linear {
				result.columns[source*8+bit][term.target] = mapped
			}
		}
	}
	return result
}
