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

//go:embed profiles/current-frame-aes.bin.gz
var frameAESProfileCompressed []byte

var currentFrameAES = func() frameAESModel {
	stream, err := gzip.NewReader(bytes.NewReader(frameAESProfileCompressed))
	if err != nil {
		panic(fmt.Sprintf("decompress frame AES profile: %v", err))
	}
	encoded, err := io.ReadAll(stream)
	closeErr := stream.Close()
	if err != nil || closeErr != nil {
		panic(fmt.Sprintf("decompress frame AES profile: %v", err))
	}
	model, err := parseFrameAESProfile(encoded)
	clear(encoded)
	if err != nil {
		panic(fmt.Sprintf("load frame AES profile: %v", err))
	}
	return model
}()

func parseFrameAESProfile(encoded []byte) (frameAESModel, error) {
	var result frameAESModel
	if len(encoded) < sha256.Size {
		return result, fmt.Errorf("frame AES profile is truncated")
	}
	payload, checksum := encoded[:len(encoded)-sha256.Size], encoded[len(encoded)-sha256.Size:]
	digest := sha256.Sum256(payload)
	if !bytes.Equal(checksum, digest[:]) {
		return result, fmt.Errorf("frame AES profile checksum mismatch")
	}
	reader := bytes.NewReader(payload)
	var magic [8]byte
	if _, err := io.ReadFull(reader, magic[:]); err != nil || string(magic[:]) != "FPAES1\x00\x00" {
		return result, fmt.Errorf("invalid frame AES profile magic")
	}
	var version uint32
	if err := binary.Read(reader, binary.LittleEndian, &version); err != nil || version != 1 {
		return result, fmt.Errorf("unsupported frame AES profile version %d", version)
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
	readLinearBytes := func(values *[16]frameByteLinear) error {
		for index := range values {
			if _, err := io.ReadFull(reader, values[index][:]); err != nil {
				return err
			}
		}
		return nil
	}
	if err := readAffine(&result.first); err != nil {
		return result, err
	}
	if err := readLinearBytes(&result.firstPost); err != nil {
		return result, err
	}
	for mix := range result.mixes {
		for row := range result.mixes[mix] {
			if _, err := io.ReadFull(reader, result.mixes[mix][row][:]); err != nil {
				return result, err
			}
		}
	}
	for round := range result.roundKeys {
		if _, err := io.ReadFull(reader, result.roundKeys[round][:]); err != nil {
			return result, err
		}
	}
	for round := range result.roundMasks {
		if _, err := io.ReadFull(reader, result.roundMasks[round][:]); err != nil {
			return result, err
		}
	}
	if err := readLinearBytes(&result.lastPre); err != nil {
		return result, err
	}
	if _, err := io.ReadFull(reader, result.lastConstant[:]); err != nil {
		return result, err
	}
	if reader.Len() != 0 {
		return result, fmt.Errorf("frame AES profile has trailing data")
	}
	if rank := frameAffine128Rank(result.first); rank != 128 {
		return result, fmt.Errorf("frame AES initial affine map has rank %d", rank)
	}
	for position, value := range result.firstPost {
		if !value.invertible() {
			return result, fmt.Errorf("frame AES initial byte map %d is not invertible", position)
		}
	}
	for index, mix := range result.mixes {
		if !frameAESMixInvertible(mix) {
			return result, fmt.Errorf("frame AES mix %d is not invertible", index)
		}
	}
	for position, value := range result.lastPre {
		if !value.invertible() {
			return result, fmt.Errorf("frame AES final byte map %d is not invertible", position)
		}
	}
	return result, nil
}
