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

//go:embed profiles/current-sidb-producer-typed.bin.gz
var sidbProducerTypedProfileCompressed []byte

func loadCurrentSIDBProducerTyped() sidbProducerTyped {
	stream, err := gzip.NewReader(bytes.NewReader(sidbProducerTypedProfileCompressed))
	if err != nil {
		panic(fmt.Sprintf("decompress SIDB producer typed profile: %v", err))
	}
	encoded, err := io.ReadAll(stream)
	closeErr := stream.Close()
	if err != nil || closeErr != nil {
		panic(fmt.Sprintf("decompress SIDB producer typed profile: %v", err))
	}
	model, err := parseSIDBProducerTypedProfile(encoded)
	clear(encoded)
	if err != nil {
		panic(fmt.Sprintf("load SIDB producer typed profile: %v", err))
	}
	return model
}

var currentSIDBProducerTyped = loadCurrentSIDBProducerTyped()

func sidbFinalXORInvertible(value [16]uint16) bool {
	matrix := value
	for column := range 16 {
		pivot := column
		for pivot < 16 && matrix[pivot]&(1<<uint(column)) == 0 {
			pivot++
		}
		if pivot == 16 {
			return false
		}
		matrix[column], matrix[pivot] = matrix[pivot], matrix[column]
		for row := range 16 {
			if row != column && matrix[row]&(1<<uint(column)) != 0 {
				matrix[row] ^= matrix[column]
			}
		}
	}
	return true
}

func readSIDBPermutation(reader *bytes.Reader, label string, targets *[16]byte) error {
	var seen [16]bool
	for source := range targets {
		target, err := reader.ReadByte()
		if err != nil || target >= 16 || seen[target] {
			return fmt.Errorf("%s source %d has invalid target %d", label, source, target)
		}
		targets[source], seen[target] = target, true
	}
	return nil
}

func parseSIDBProducerTypedProfile(encoded []byte) (sidbProducerTyped, error) {
	var result sidbProducerTyped
	if len(encoded) < 32 {
		return result, fmt.Errorf("SIDB producer typed profile is truncated")
	}
	payload, checksum := encoded[:len(encoded)-32], encoded[len(encoded)-32:]
	digest := sha256.Sum256(payload)
	if !bytes.Equal(checksum, digest[:]) {
		return result, fmt.Errorf("SIDB producer typed profile checksum mismatch")
	}
	reader := bytes.NewReader(payload)
	var magic [8]byte
	if _, err := io.ReadFull(reader, magic[:]); err != nil || string(magic[:]) != "FPSTP1\x00\x00" {
		return result, fmt.Errorf("invalid SIDB producer typed profile magic")
	}
	var version uint32
	if err := binary.Read(reader, binary.LittleEndian, &version); err != nil || version != 1 {
		return result, fmt.Errorf("unsupported SIDB producer typed profile version %d", version)
	}
	if err := readSIDBPermutation(reader, "first inverse-S-box", &result.firstTargets); err != nil {
		return result, err
	}
	for target := range result.firstOutput {
		if _, err := io.ReadFull(reader, result.firstOutput[target][:]); err != nil {
			return result, err
		}
		if !result.firstOutput[target].invertible() {
			return result, fmt.Errorf("first output map %d is singular", target)
		}
	}
	keyCount, err := reader.ReadByte()
	if err != nil || keyCount != byte(len(result.roundKeys)) {
		return result, fmt.Errorf("invalid SIDB producer typed key count %d", keyCount)
	}
	var reserved [3]byte
	if _, err := io.ReadFull(reader, reserved[:]); err != nil || reserved != [3]byte{} {
		return result, fmt.Errorf("invalid SIDB producer typed key header")
	}
	for round := range result.roundKeys {
		if _, err := io.ReadFull(reader, result.roundKeys[round][:]); err != nil {
			return result, err
		}
	}
	if err := readSIDBPermutation(reader, "last inverse-S-box", &result.lastTargets); err != nil {
		return result, err
	}
	if _, err := io.ReadFull(reader, result.lastKeys[:]); err != nil {
		return result, err
	}
	for target := range result.finalXOR {
		if err := binary.Read(reader, binary.LittleEndian, &result.finalXOR[target]); err != nil {
			return result, err
		}
		if result.finalXOR[target] == 0 {
			return result, fmt.Errorf("final XOR target %d has no inputs", target)
		}
	}
	if !sidbFinalXORInvertible(result.finalXOR) {
		return result, fmt.Errorf("final XOR graph is singular")
	}
	for target := range result.finalOutput {
		if _, err := io.ReadFull(reader, result.finalOutput[target][:]); err != nil {
			return result, err
		}
		if !result.finalOutput[target].invertible() {
			return result, fmt.Errorf("final output map %d is singular", target)
		}
	}
	if _, err := io.ReadFull(reader, result.finalConstant[:]); err != nil {
		return result, err
	}
	if reader.Len() != 0 {
		return result, fmt.Errorf("SIDB producer typed profile has trailing data")
	}
	return result, nil
}
