package recovery

import (
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	_ "embed"
	"encoding/binary"
	"fmt"
	"io"
	"sort"
)

//go:embed profiles/current-sidb-producer-inverse-aes.bin.gz
var sidbProducerInverseAESProfileCompressed []byte

func loadCurrentSIDBProducerInverseAES() sidbProducerInverseAES {
	stream, err := gzip.NewReader(bytes.NewReader(sidbProducerInverseAESProfileCompressed))
	if err != nil {
		panic(fmt.Sprintf("decompress SIDB producer inverse-AES profile: %v", err))
	}
	encoded, err := io.ReadAll(stream)
	closeErr := stream.Close()
	if err != nil || closeErr != nil {
		panic(fmt.Sprintf("decompress SIDB producer inverse-AES profile: %v", err))
	}
	model, err := parseSIDBProducerInverseAESProfile(encoded)
	clear(encoded)
	if err != nil {
		panic(fmt.Sprintf("load SIDB producer inverse-AES profile: %v", err))
	}
	return model
}

var currentSIDBProducerInverseAES = loadCurrentSIDBProducerInverseAES()

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

func parseSIDBProducerInverseAESProfile(encoded []byte) (sidbProducerInverseAES, error) {
	var result sidbProducerInverseAES
	if len(encoded) < 32 {
		return result, fmt.Errorf("SIDB producer inverse-AES profile is truncated")
	}
	payload, checksum := encoded[:len(encoded)-32], encoded[len(encoded)-32:]
	digest := sha256.Sum256(payload)
	if !bytes.Equal(checksum, digest[:]) {
		return result, fmt.Errorf("SIDB producer inverse-AES profile checksum mismatch")
	}
	reader := bytes.NewReader(payload)
	var magic [8]byte
	if _, err := io.ReadFull(reader, magic[:]); err != nil || string(magic[:]) != "FPSIA1\x00\x00" {
		return result, fmt.Errorf("invalid SIDB producer inverse-AES profile magic")
	}
	var version uint32
	if err := binary.Read(reader, binary.LittleEndian, &version); err != nil || version != 1 {
		return result, fmt.Errorf("unsupported SIDB producer inverse-AES profile version %d", version)
	}
	var err error
	result.first, err = readSIDBSparseAffine(reader, "first layer", 16, func(count int) bool { return count == 1 })
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
	for position := range result.firstPost {
		if _, err := io.ReadFull(reader, result.firstPost[position][:]); err != nil {
			return result, err
		}
		if !result.firstPost[position].invertible() {
			return result, fmt.Errorf("first post-map %d is singular", position)
		}
	}
	mixCount, err := reader.ReadByte()
	if err != nil || mixCount != byte(len(result.mixes)) {
		return result, fmt.Errorf("invalid SIDB producer inverse-AES mix count %d", mixCount)
	}
	var reserved [3]byte
	if _, err := io.ReadFull(reader, reserved[:]); err != nil || reserved != [3]byte{} {
		return result, fmt.Errorf("invalid SIDB producer inverse-AES mix header")
	}
	for mixIndex := range result.mixes {
		var full [16][16]byte
		for target := range result.mixes[mixIndex] {
			var seen [16]bool
			coefficients := make([]int, 0, 4)
			for termIndex := range result.mixes[mixIndex][target] {
				term := &result.mixes[mixIndex][target][termIndex]
				source, err := reader.ReadByte()
				if err != nil || source >= 16 || seen[source] {
					return result, fmt.Errorf("mix %d target %d has invalid source %d", mixIndex, target, source)
				}
				coefficient, err := reader.ReadByte()
				if err != nil || coefficient == 0 {
					return result, fmt.Errorf("mix %d target %d has invalid coefficient %d", mixIndex, target, coefficient)
				}
				term.source, term.coefficient, seen[source] = source, coefficient, true
				full[target][source] = coefficient
				coefficients = append(coefficients, int(coefficient))
			}
			sort.Ints(coefficients)
			if coefficients[0] != 0x09 || coefficients[1] != 0x0b || coefficients[2] != 0x0d || coefficients[3] != 0x0e {
				return result, fmt.Errorf("mix %d target %d is not InvMixColumns", mixIndex, target)
			}
		}
		if !frameAESMixInvertible(full) {
			return result, fmt.Errorf("mix %d is singular", mixIndex)
		}
	}
	keyCount, err := reader.ReadByte()
	if err != nil || keyCount != byte(len(result.roundKeys)) {
		return result, fmt.Errorf("invalid SIDB producer inverse-AES key count %d", keyCount)
	}
	if _, err := io.ReadFull(reader, reserved[:]); err != nil || reserved != [3]byte{} {
		return result, fmt.Errorf("invalid SIDB producer inverse-AES key header")
	}
	for round := range result.roundKeys {
		if _, err := io.ReadFull(reader, result.roundKeys[round][:]); err != nil {
			return result, err
		}
	}
	for position := range result.lastPre {
		if _, err := io.ReadFull(reader, result.lastPre[position][:]); err != nil {
			return result, err
		}
		if !result.lastPre[position].invertible() {
			return result, fmt.Errorf("last pre-map %d is singular", position)
		}
	}
	if _, err := io.ReadFull(reader, result.lastConstant[:]); err != nil {
		return result, err
	}
	result.final, err = readSIDBSparseAffine(reader, "final layer", 40, func(count int) bool { return count == 2 || count == 3 })
	if err != nil {
		return result, err
	}
	if reader.Len() != 0 {
		return result, fmt.Errorf("SIDB producer inverse-AES profile has trailing data")
	}
	return result, nil
}
