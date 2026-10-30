package recovery

import (
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"os"
	"testing"
)

func normalizeFrameRound(round generalizedByteRound) (frameSPNStage, error) {
	var result frameSPNStage
	result.output.constant = round.baseline
	for position := range 16 {
		var table [256]byte
		for input := range 256 {
			table[input] = round.tables[position][input][0]
		}
		inputMap, _, inputOffset, _, _, err := solveAffineSBox(frameSBox, table)
		if err != nil {
			return result, fmt.Errorf("input %d: %w", position, err)
		}
		var canonical [256]byte
		for input := range 256 {
			canonical[input] = frameSBox[inputMap.apply(byte(input))^inputOffset]
		}
		result.input.constant[position] = inputOffset
		for bit := range 8 {
			result.input.columns[position*8+bit][position] = inputMap.apply(1 << bit)
		}
		for slot, target := range round.targets[position] {
			var contribution [256]byte
			for input := range 256 {
				contribution[input] = round.tables[position][input][slot]
			}
			linear, offset, err := affineRelation(canonical, contribution)
			if err != nil {
				return result, fmt.Errorf("input %d target %d: %w", position, target, err)
			}
			result.output.constant[target] ^= offset
			for bit := range 8 {
				result.output.columns[position*8+bit][target] ^= linear.apply(1 << bit)
			}
		}
	}
	return result, nil
}

func normalizeFrameSPN(model frameSemanticModel) (frameSPNModel, error) {
	var result frameSPNModel
	if len(model.rounds) != 9 {
		return result, fmt.Errorf("frame model has %d generalized rounds", len(model.rounds))
	}
	rounds := make([]frameSPNStage, len(model.rounds))
	for index, round := range model.rounds {
		normalized, err := normalizeFrameRound(round)
		if err != nil {
			return result, fmt.Errorf("round %d: %w", index, err)
		}
		rounds[index] = normalized
	}
	result.initial = frameAffine128{constant: model.inputConstant, columns: model.inputColumns}
	result.stages = rounds
	groups, err := frameFinalComponents(model)
	if err != nil {
		return result, err
	}
	var finalInput frameAffine128
	var inverseSBox [256]byte
	for input, output := range frameSBox {
		inverseSBox[output] = byte(input)
	}
	var assigned [16]bool
	for _, group := range groups {
		for _, position := range group.inputs {
			if assigned[position] {
				return result, fmt.Errorf("duplicate final input %d", position)
			}
			assigned[position] = true
			var preimage [256]byte
			for input := range 256 {
				preimage[input] = inverseSBox[group.canonical[position][input]]
			}
			offset := preimage[0]
			var linear linearByteMap
			for bit := range 8 {
				linear.basis[bit] = preimage[1<<bit] ^ offset
			}
			for input := range 256 {
				if preimage[input] != linear.apply(byte(input))^offset {
					return result, fmt.Errorf("final input %d is not affine", position)
				}
			}
			finalInput.constant[position] = offset
			for bit := range 8 {
				finalInput.columns[position*8+bit][position] = linear.apply(1 << bit)
			}
		}
		for outputSlot, output := range group.outputs {
			result.finalOutput.constant[output] = byte(group.constant >> (outputSlot * 8))
		}
		for inputSlot, input := range group.inputs {
			for bit := range 8 {
				column := group.forward[inputSlot*8+bit]
				for outputSlot, output := range group.outputs {
					result.finalOutput.columns[input*8+bit][output] = byte(column >> (outputSlot * 8))
				}
			}
		}
		for _, output := range group.outputs {
			entry := model.final[output]
			if entry.inputCount == 1 {
				var seen [256]bool
				for raw := range 256 {
					canonical := group.decode[output][raw]
					if seen[canonical] {
						return result, fmt.Errorf("final direct output %d decoder is not invertible", output)
					}
					seen[canonical] = true
					result.outputEncoding[output][canonical] = byte(raw)
				}
				continue
			}
			for value, encoded := range entry.encode {
				result.outputEncoding[output][value] = entry.baseline ^ encoded
			}
		}
	}
	result.finalInput = finalInput
	result.exitInputs = model.exitInputs
	result.exitBaseline = model.exitBaseline
	result.exitContributions = model.exitContributions
	return result, nil
}

func TestFrameSPNNormalization(t *testing.T) {
	normalized, err := normalizeFrameSPN(frameSemantic26A428)
	if err != nil {
		t.Fatal(err)
	}
	for index := range 4096 {
		digest := sha256.Sum256([]byte{0x53, 0x50, 0x4e, byte(index >> 8), byte(index)})
		var input [16]byte
		copy(input[:], digest[:16])
		got, gotExit := normalized.apply(input)
		want, wantExit := frameSemantic26A428.apply(input)
		if got != want {
			t.Fatalf("control %d = %x, want %x", index, got, want)
		}
		if gotExit != wantExit {
			t.Fatalf("control %d exit = %#x, want %#x", index, gotExit, wantExit)
		}
	}
}

func TestWriteFrameSPNProfile(t *testing.T) {
	outputPath := os.Getenv("UNFAIRPLAY_FRAME_SPN_OUTPUT")
	if outputPath == "" {
		t.Skip("frame SPN profile output is disabled")
	}
	model, err := normalizeFrameSPN(frameSemantic26A428)
	if err != nil {
		t.Fatal(err)
	}
	var payload bytes.Buffer
	payload.WriteString("FPSPN1\x00\x00")
	if err := binary.Write(&payload, binary.LittleEndian, uint32(1)); err != nil {
		t.Fatal(err)
	}
	payload.WriteByte(byte(len(model.stages)))
	payload.Write([]byte{0, 0, 0})
	writeAffine := func(value frameAffine128) {
		payload.Write(value.constant[:])
		for _, column := range value.columns {
			payload.Write(column[:])
		}
	}
	writeAffine(model.initial)
	for _, stage := range model.stages {
		writeAffine(stage.input)
		writeAffine(stage.output)
	}
	writeAffine(model.finalInput)
	writeAffine(model.finalOutput)
	for _, table := range model.outputEncoding {
		payload.Write(table[:])
	}
	payload.Write(model.exitInputs[:])
	payload.Write(model.exitBaseline[:])
	for outputByte := range 4 {
		for slot := range 4 {
			payload.Write(model.exitContributions[outputByte][slot][:])
		}
	}
	digest := sha256.Sum256(payload.Bytes())
	payload.Write(digest[:])
	var compressed bytes.Buffer
	stream, err := gzip.NewWriterLevel(&compressed, gzip.BestCompression)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := stream.Write(payload.Bytes()); err != nil {
		t.Fatal(err)
	}
	if err := stream.Close(); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(outputPath, compressed.Bytes(), 0o644); err != nil {
		t.Fatal(err)
	}
}
