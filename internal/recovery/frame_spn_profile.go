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

//go:embed profiles/current-frame-spn.bin.gz
var frameSPNProfileCompressed []byte

var currentFrameSPN = func() frameSPNModel {
	stream, err := gzip.NewReader(bytes.NewReader(frameSPNProfileCompressed))
	if err != nil {
		panic(fmt.Sprintf("decompress frame SPN profile: %v", err))
	}
	encoded, err := io.ReadAll(stream)
	closeErr := stream.Close()
	if err != nil || closeErr != nil {
		panic(fmt.Sprintf("decompress frame SPN profile: %v", err))
	}
	model, err := parseFrameSPNProfile(encoded)
	clear(encoded)
	if err != nil {
		panic(fmt.Sprintf("load frame SPN profile: %v", err))
	}
	return model
}()

func parseFrameSPNProfile(encoded []byte) (frameSPNModel, error) {
	var result frameSPNModel
	if len(encoded) < 32 {
		return result, fmt.Errorf("frame SPN profile is truncated")
	}
	payload, checksum := encoded[:len(encoded)-32], encoded[len(encoded)-32:]
	digest := sha256.Sum256(payload)
	if !bytes.Equal(checksum, digest[:]) {
		return result, fmt.Errorf("frame SPN profile checksum mismatch")
	}
	reader := bytes.NewReader(payload)
	var magic [8]byte
	if _, err := io.ReadFull(reader, magic[:]); err != nil || string(magic[:]) != "FPSPN1\x00\x00" {
		return result, fmt.Errorf("invalid frame SPN profile magic")
	}
	var version uint32
	if err := binary.Read(reader, binary.LittleEndian, &version); err != nil || version != 1 {
		return result, fmt.Errorf("unsupported frame SPN profile version %d", version)
	}
	stageCount, err := reader.ReadByte()
	if err != nil || stageCount != 9 {
		return result, fmt.Errorf("invalid frame SPN stage count %d", stageCount)
	}
	var reserved [3]byte
	if _, err := io.ReadFull(reader, reserved[:]); err != nil || reserved != [3]byte{} {
		return result, fmt.Errorf("invalid frame SPN header")
	}
	readAffine := func(value *frameAffine128) error {
		if _, err := io.ReadFull(reader, value.constant[:]); err != nil {
			return err
		}
		for index := range value.columns {
			if _, err := io.ReadFull(reader, value.columns[index][:]); err != nil {
				return err
			}
		}
		return nil
	}
	if err := readAffine(&result.initial); err != nil {
		return result, err
	}
	result.stages = make([]frameSPNStage, stageCount)
	for index := range result.stages {
		if err := readAffine(&result.stages[index].input); err != nil {
			return result, err
		}
		if err := readAffine(&result.stages[index].output); err != nil {
			return result, err
		}
	}
	if err := readAffine(&result.finalInput); err != nil {
		return result, err
	}
	if err := readAffine(&result.finalOutput); err != nil {
		return result, err
	}
	for position := range result.outputEncoding {
		if _, err := io.ReadFull(reader, result.outputEncoding[position][:]); err != nil {
			return result, err
		}
		var seen [256]bool
		for _, value := range result.outputEncoding[position] {
			if seen[value] {
				return result, fmt.Errorf("frame SPN output %d is not a permutation", position)
			}
			seen[value] = true
		}
	}
	if _, err := io.ReadFull(reader, result.exitInputs[:]); err != nil {
		return result, err
	}
	for _, input := range result.exitInputs {
		if input >= 16 {
			return result, fmt.Errorf("invalid frame SPN exit input %d", input)
		}
	}
	if _, err := io.ReadFull(reader, result.exitBaseline[:]); err != nil {
		return result, err
	}
	for outputByte := range 4 {
		for slot := range 4 {
			if _, err := io.ReadFull(reader, result.exitContributions[outputByte][slot][:]); err != nil {
				return result, err
			}
		}
	}
	if reader.Len() != 0 {
		return result, fmt.Errorf("frame SPN profile has trailing data")
	}
	return result, nil
}
