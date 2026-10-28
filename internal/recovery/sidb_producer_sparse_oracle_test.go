package recovery

import (
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/binary"
	"io"
	"os"
	"testing"
)

func sparseSIDBProducerAffine(value frameAffine128) sidbSparseAffine {
	result := sidbSparseAffine{constant: value.constant}
	for source := range 16 {
		for target := range 16 {
			var linear frameByteLinear
			for bit := range 8 {
				linear[bit] = value.columns[source*8+bit][target]
			}
			if linear != (frameByteLinear{}) {
				result.terms[source] = append(result.terms[source], sidbSparseTerm{target: byte(target), linear: linear})
			}
		}
	}
	return result
}

func sparseSIDBProducerFromSPN(value sidbProducerSPN) sidbProducerSparse {
	result := sidbProducerSparse{
		first: sparseSIDBProducerAffine(composeFrameAffine(value.stages[0].input, value.bridge)),
		final: sparseSIDBProducerAffine(value.stages[len(value.stages)-1].output),
	}
	for index := range result.rounds {
		result.rounds[index] = sparseSIDBProducerAffine(composeFrameAffine(value.stages[index+1].input, value.stages[index].output))
	}
	return result
}

func writeSIDBProducerSparseProfile(writer io.Writer, value sidbProducerSparse) error {
	var payload bytes.Buffer
	payload.WriteString("FPSPA1\x00\x00")
	if err := binary.Write(&payload, binary.LittleEndian, uint32(1)); err != nil {
		return err
	}
	writeAffine := func(affine sidbSparseAffine) {
		payload.Write(affine.constant[:])
		for _, terms := range affine.terms {
			payload.WriteByte(byte(len(terms)))
			for _, term := range terms {
				payload.WriteByte(term.target)
				payload.Write(term.linear[:])
			}
		}
	}
	writeAffine(value.first)
	payload.WriteByte(byte(len(value.rounds)))
	payload.Write([]byte{0, 0, 0})
	for _, round := range value.rounds {
		writeAffine(round)
	}
	writeAffine(value.final)
	checksum := sha256.Sum256(payload.Bytes())
	payload.Write(checksum[:])
	stream, err := gzip.NewWriterLevel(writer, gzip.BestCompression)
	if err != nil {
		return err
	}
	if _, err := stream.Write(payload.Bytes()); err != nil {
		return err
	}
	return stream.Close()
}

func TestSIDBProducerSparseMatchesSPNOracle(t *testing.T) {
	model := sparseSIDBProducerFromSPN(currentSIDBProducerSPN)
	for control := range 4096 {
		digest := sha256.Sum256([]byte{0x53, 0x50, 0x41, 0x52, byte(control >> 8), byte(control)})
		var input [16]byte
		copy(input[:], digest[:16])
		if got, want := model.apply(input), currentSIDBProducerSPN.apply(input); got != want {
			t.Fatalf("control %d = %x, want %x", control, got, want)
		}
	}
	if path := os.Getenv("UNFAIRPLAY_SIDB_PRODUCER_SPARSE_OUTPUT"); path != "" {
		file, err := os.Create(path)
		if err != nil {
			t.Fatal(err)
		}
		writeErr := writeSIDBProducerSparseProfile(file, model)
		closeErr := file.Close()
		if writeErr != nil {
			t.Fatal(writeErr)
		}
		if closeErr != nil {
			t.Fatal(closeErr)
		}
	}
}
