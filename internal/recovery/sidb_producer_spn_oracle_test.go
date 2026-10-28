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

type sidbProducerSPN struct {
	bridge frameAffine128
	stages []frameSPNStage
}

func (model sidbProducerSPN) apply(input [16]byte) [16]byte {
	state := model.bridge.apply(input)
	for _, stage := range model.stages {
		state = stage.output.apply(frameSBoxState(stage.input.apply(state)))
	}
	return state
}

//go:embed profiles/macos26a428-sidb-producer-spn-oracle.bin.gz
var sidbProducerSPNProfileCompressed []byte

var currentSIDBProducerSPN = func() sidbProducerSPN {
	stream, err := gzip.NewReader(bytes.NewReader(sidbProducerSPNProfileCompressed))
	if err != nil {
		panic(fmt.Sprintf("decompress SIDB producer SPN profile: %v", err))
	}
	encoded, err := io.ReadAll(stream)
	closeErr := stream.Close()
	if err != nil || closeErr != nil {
		panic(fmt.Sprintf("decompress SIDB producer SPN profile: %v", err))
	}
	model, err := parseSIDBProducerSPNProfile(encoded)
	clear(encoded)
	if err != nil {
		panic(fmt.Sprintf("load SIDB producer SPN profile: %v", err))
	}
	return model
}()

func parseSIDBProducerSPNProfile(encoded []byte) (sidbProducerSPN, error) {
	var result sidbProducerSPN
	if len(encoded) < 32 {
		return result, fmt.Errorf("SIDB producer SPN profile is truncated")
	}
	payload, checksum := encoded[:len(encoded)-32], encoded[len(encoded)-32:]
	digest := sha256.Sum256(payload)
	if !bytes.Equal(checksum, digest[:]) {
		return result, fmt.Errorf("SIDB producer SPN profile checksum mismatch")
	}
	reader := bytes.NewReader(payload)
	var magic [8]byte
	if _, err := io.ReadFull(reader, magic[:]); err != nil || string(magic[:]) != "FPSPK1\x00\x00" {
		return result, fmt.Errorf("invalid SIDB producer SPN profile magic")
	}
	var version uint32
	if err := binary.Read(reader, binary.LittleEndian, &version); err != nil || version != 1 {
		return result, fmt.Errorf("unsupported SIDB producer SPN profile version %d", version)
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
	if err := readAffine(&result.bridge); err != nil {
		return result, err
	}
	stageCount, err := reader.ReadByte()
	if err != nil || stageCount != 10 {
		return result, fmt.Errorf("invalid SIDB producer SPN stage count %d", stageCount)
	}
	var reserved [3]byte
	if _, err := io.ReadFull(reader, reserved[:]); err != nil || reserved != [3]byte{} {
		return result, fmt.Errorf("invalid SIDB producer SPN header")
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
	if reader.Len() != 0 {
		return result, fmt.Errorf("SIDB producer SPN profile has trailing data")
	}
	return result, nil
}
