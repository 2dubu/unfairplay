package recovery

import (
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"testing"
)

func readSIDBProducerResearchProfile(path string) ([]generalizedByteRound, error) {
	compressed, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	stream, err := gzip.NewReader(bytes.NewReader(compressed))
	if err != nil {
		return nil, err
	}
	encoded, err := io.ReadAll(stream)
	closeErr := stream.Close()
	if err != nil || closeErr != nil {
		return nil, fmt.Errorf("decompress producer profile: %v", err)
	}
	if len(encoded) < 32 {
		return nil, fmt.Errorf("producer profile is truncated")
	}
	payload, checksum := encoded[:len(encoded)-32], encoded[len(encoded)-32:]
	digest := sha256.Sum256(payload)
	if !bytes.Equal(checksum, digest[:]) {
		return nil, fmt.Errorf("producer profile checksum mismatch")
	}
	reader := bytes.NewReader(payload)
	var magic [8]byte
	if _, err := io.ReadFull(reader, magic[:]); err != nil || string(magic[:]) != "FPSIDP1\x00" {
		return nil, fmt.Errorf("invalid producer profile magic")
	}
	var version uint32
	if err := binary.Read(reader, binary.LittleEndian, &version); err != nil || version != 1 {
		return nil, fmt.Errorf("unsupported producer profile version %d", version)
	}
	roundCount, err := reader.ReadByte()
	if err != nil || roundCount != 10 {
		return nil, fmt.Errorf("invalid producer round count %d", roundCount)
	}
	var reserved [3]byte
	if _, err := io.ReadFull(reader, reserved[:]); err != nil || reserved != [3]byte{} {
		return nil, fmt.Errorf("invalid producer profile header")
	}
	rounds := make([]generalizedByteRound, roundCount)
	for roundIndex := range rounds {
		round := &rounds[roundIndex]
		if _, err := io.ReadFull(reader, round.baseline[:]); err != nil {
			return nil, err
		}
		for position := range 16 {
			if _, err := io.ReadFull(reader, round.targets[position][:]); err != nil {
				return nil, err
			}
			var seen [16]bool
			for _, target := range round.targets[position] {
				if target >= 16 || seen[target] {
					return nil, fmt.Errorf("producer round %d input %d has invalid target %d", roundIndex, position, target)
				}
				seen[target] = true
			}
		}
		for position := range 16 {
			for value := range 256 {
				if _, err := io.ReadFull(reader, round.tables[position][value][:]); err != nil {
					return nil, err
				}
			}
		}
	}
	if reader.Len() != 0 {
		return nil, fmt.Errorf("producer profile has trailing data")
	}
	return rounds, nil
}

type sidbKeyLayerResearch struct {
	baseline      [16]byte
	contributions [16][256][16]byte
}

func readSIDBKeyLayerResearch(path string) (sidbKeyLayerResearch, error) {
	var result sidbKeyLayerResearch
	encoded, err := os.ReadFile(path)
	if err != nil {
		return result, err
	}
	if len(encoded) < 32 {
		return result, fmt.Errorf("key-layer profile is truncated")
	}
	payload, checksum := encoded[:len(encoded)-32], encoded[len(encoded)-32:]
	digest := sha256.Sum256(payload)
	if !bytes.Equal(checksum, digest[:]) {
		return result, fmt.Errorf("key-layer profile checksum mismatch")
	}
	reader := bytes.NewReader(payload)
	var magic [8]byte
	if _, err := io.ReadFull(reader, magic[:]); err != nil || string(magic[:]) != "FPSIDK1\x00" {
		return result, fmt.Errorf("invalid key-layer profile magic")
	}
	var version uint32
	if err := binary.Read(reader, binary.LittleEndian, &version); err != nil || version != 1 {
		return result, fmt.Errorf("unsupported key-layer profile version %d", version)
	}
	if _, err := io.ReadFull(reader, result.baseline[:]); err != nil {
		return result, err
	}
	for position := range 16 {
		for value := range 256 {
			if _, err := io.ReadFull(reader, result.contributions[position][value][:]); err != nil {
				return result, err
			}
		}
	}
	if reader.Len() != 0 {
		return result, fmt.Errorf("key-layer profile has trailing data")
	}
	return result, nil
}

func normalizeSIDBKeyLayer(layer sidbKeyLayerResearch) (frameSPNStage, error) {
	var result frameSPNStage
	result.output.constant = layer.baseline
	for position := range 16 {
		chosen := -1
		var chosenTable [256]byte
		for target := range 16 {
			var table [256]byte
			for value := range 256 {
				table[value] = layer.contributions[position][value][target]
			}
			if byteTablePermutation(table) {
				chosen, chosenTable = target, table
				break
			}
		}
		if chosen < 0 {
			return result, fmt.Errorf("key-layer input %d has no permutation component", position)
		}
		inputMap, _, inputOffset, _, _, err := solveAffineSBox(frameSBox, chosenTable)
		if err != nil {
			return result, fmt.Errorf("key-layer input %d: %w", position, err)
		}
		var canonical [256]byte
		for value := range 256 {
			canonical[value] = frameSBox[inputMap.apply(byte(value))^inputOffset]
		}
		result.input.constant[position] = inputOffset
		for bit := range 8 {
			result.input.columns[position*8+bit][position] = inputMap.apply(1 << bit)
		}
		for target := range 16 {
			var table [256]byte
			for value := range 256 {
				table[value] = layer.contributions[position][value][target]
			}
			linear, offset, err := affineRelation(canonical, table)
			if err != nil {
				return result, fmt.Errorf("key-layer input %d target %d: %w", position, target, err)
			}
			result.output.constant[target] ^= offset
			for bit := range 8 {
				result.output.columns[position*8+bit][target] ^= linear.apply(1 << bit)
			}
		}
	}
	return result, nil
}

func sidbPrefixToSemantic(input [16]byte) (output [16]byte) {
	copy(output[0:4], input[12:16])
	copy(output[4:8], input[4:8])
	copy(output[12:16], input[0:4])
	for index := range 4 {
		output[8+index] = input[8+index] ^ input[index]
	}
	return output
}

type sidbPrefixNormalization struct {
	substitutions [16][256]byte
	output        frameAffine128
}

func normalizeSIDBPrefix(round generalizedByteRound) (sidbPrefixNormalization, error) {
	var result sidbPrefixNormalization
	result.output.constant = round.baseline
	for position := range 16 {
		chosen := -1
		for slot := range 4 {
			var table [256]byte
			for value := range 256 {
				table[value] = round.tables[position][value][slot]
			}
			if byteTablePermutation(table) {
				chosen = slot
				result.substitutions[position] = table
				break
			}
		}
		if chosen < 0 {
			return result, fmt.Errorf("prefix input %d has no permutation component", position)
		}
		for slot, target := range round.targets[position] {
			var table [256]byte
			for value := range 256 {
				table[value] = round.tables[position][value][slot]
			}
			linear, offset, err := affineRelation(result.substitutions[position], table)
			if err != nil {
				return result, fmt.Errorf("prefix input %d target %d: %w", position, target, err)
			}
			result.output.constant[target] ^= offset
			for bit := range 8 {
				result.output.columns[position*8+bit][target] ^= linear.apply(1 << bit)
			}
		}
	}
	return result, nil
}

func composeFrameAffine(outer, inner frameAffine128) (result frameAffine128) {
	result.constant = outer.apply(inner.constant)
	for index, column := range inner.columns {
		result.columns[index] = outer.linear(column)
	}
	return result
}

func affineFromStateFunction(function func([16]byte) [16]byte) frameAffine128 {
	var result frameAffine128
	result.constant = function([16]byte{})
	for bit := range 128 {
		var input [16]byte
		input[bit/8] = 1 << uint(bit%8)
		output := function(input)
		for index := range output {
			result.columns[bit][index] = output[index] ^ result.constant[index]
		}
	}
	return result
}

func TestSIDBProducerSPNResearch(t *testing.T) {
	profilePath := os.Getenv("UNFAIRPLAY_SIDB_PRODUCER_PROFILE")
	keyLayerPath := os.Getenv("UNFAIRPLAY_SIDB_KEY_LAYER_PROFILE")
	if profilePath == "" || keyLayerPath == "" {
		t.Skip("private SIDB producer profile is disabled")
	}
	rounds, err := readSIDBProducerResearchProfile(profilePath)
	if err != nil {
		t.Fatal(err)
	}
	prefix := rounds[0]
	permutations, affineRelations := 0, 0
	for position := range 16 {
		var source [256]byte
		for value := range 256 {
			source[value] = prefix.tables[position][value][0]
		}
		positionPermutation := byteTablePermutation(source)
		if positionPermutation {
			permutations++
		}
		positionAffine := 0
		for slot := range 4 {
			var target [256]byte
			for value := range 256 {
				target[value] = prefix.tables[position][value][slot]
			}
			if _, _, err := affineRelation(source, target); err == nil {
				affineRelations++
				positionAffine++
			}
		}
	}
	if permutations != 12 || affineRelations != 48 {
		t.Fatalf("prefix geometry = %d permutations and %d affine relations", permutations, affineRelations)
	}
	prefixNormalized, err := normalizeSIDBPrefix(prefix)
	if err != nil {
		t.Fatal(err)
	}
	var bridgeInput frameAffine128
	for position := range 16 {
		var combined [256]byte
		for value := range 256 {
			combined[value] = prefixNormalized.substitutions[position][currentFrameSPN.outputEncoding[position][value]]
		}
		if !affineByteMap(combined) {
			t.Fatalf("frame-to-producer byte %d does not reduce to an affine map", position)
		}
		bridgeInput.constant[position] = combined[0]
		for bit := range 8 {
			bridgeInput.columns[position*8+bit][position] = combined[1<<bit] ^ combined[0]
		}
	}
	prefixToSemanticAffine := affineFromStateFunction(sidbPrefixToSemantic)
	bridge := composeFrameAffine(prefixToSemanticAffine, composeFrameAffine(prefixNormalized.output, bridgeInput))
	for index := range 4096 {
		digest := sha256.Sum256([]byte{0x42, 0x52, 0x49, 0x44, byte(index >> 8), byte(index)})
		var input [16]byte
		copy(input[:], digest[:16])
		var encoded [16]byte
		for position, value := range input {
			encoded[position] = currentFrameSPN.outputEncoding[position][value]
		}
		want := sidbPrefixToSemantic(prefix.apply(encoded))
		if got := bridge.apply(input); got != want {
			t.Fatalf("bridge control %d = %x, want %x", index, got, want)
		}
	}
	stages := make([]frameSPNStage, len(rounds)-1)
	for index, round := range rounds[1:] {
		stage, err := normalizeFrameRound(round)
		if err != nil {
			t.Fatalf("round %d: %v", index+1, err)
		}
		stages[index] = stage
	}
	for index := range 4096 {
		digest := sha256.Sum256([]byte{0x53, 0x49, 0x44, 0x50, byte(index >> 8), byte(index)})
		var input [16]byte
		copy(input[:], digest[:16])
		want := sidbPrefixToSemantic(rounds[0].apply(input))
		got := want
		for roundIndex := 1; roundIndex < len(rounds); roundIndex++ {
			want = rounds[roundIndex].apply(want)
			stage := stages[roundIndex-1]
			got = stage.output.apply(frameSBoxState(stage.input.apply(got)))
		}
		if got != want {
			t.Fatalf("control %d = %x, want %x", index, got, want)
		}
	}
	keyLayer, err := readSIDBKeyLayerResearch(keyLayerPath)
	if err != nil {
		t.Fatal(err)
	}
	keyStage, err := normalizeSIDBKeyLayer(keyLayer)
	if err != nil {
		t.Fatal(err)
	}
	for index := range 4096 {
		digest := sha256.Sum256([]byte{0x4b, 0x45, 0x59, 0x4c, byte(index >> 8), byte(index)})
		var input [16]byte
		copy(input[:], digest[:16])
		want := keyLayer.baseline
		for position, value := range input {
			for target, contribution := range keyLayer.contributions[position][value] {
				want[target] ^= contribution
			}
		}
		got := keyStage.output.apply(frameSBoxState(keyStage.input.apply(input)))
		if got != want {
			t.Fatalf("key-layer control %d = %x, want %x", index, got, want)
		}
	}
	stages = append(stages, keyStage)
	applyProducer := func(input [16]byte) [16]byte {
		state := bridge.apply(input)
		for _, stage := range stages {
			state = stage.output.apply(frameSBoxState(stage.input.apply(state)))
		}
		return state
	}
	for index := range 4096 {
		digest := sha256.Sum256([]byte{0x50, 0x52, 0x4f, 0x44, byte(index >> 8), byte(index)})
		var input [16]byte
		copy(input[:], digest[:16])
		if got, want := currentSIDBProducerSPN.apply(input), applyProducer(input); got != want {
			t.Fatalf("production producer control %d = %x, want %x", index, got, want)
		}
	}
	if reportPath := os.Getenv("UNFAIRPLAY_SIDB_KEY_REPORT_INPUT"); reportPath != "" {
		type reportItem struct {
			Frame string `json:"frame"`
			Key   string `json:"key"`
		}
		raw, err := os.ReadFile(reportPath)
		if err != nil {
			t.Fatal(err)
		}
		var report []reportItem
		if err := json.Unmarshal(raw, &report); err != nil {
			t.Fatal(err)
		}
		var inverseOutput [16][256]byte
		for position, table := range currentFrameSPN.outputEncoding {
			for input, output := range table {
				inverseOutput[position][output] = byte(input)
			}
		}
		for index, item := range report {
			frame, err := hex.DecodeString(item.Frame)
			if err != nil || len(frame) != 16 {
				t.Fatalf("invalid frame %q", item.Frame)
			}
			key, err := hex.DecodeString(item.Key)
			if err != nil || len(key) != 16 {
				t.Fatalf("invalid key %q", item.Key)
			}
			var internal [16]byte
			for position, value := range frame {
				internal[position] = inverseOutput[position][value]
			}
			if got := applyProducer(internal); !bytes.Equal(got[:], key) {
				t.Fatalf("key report %d = %x, want %x", index, got, key)
			}
		}
	}
	if outputPath := os.Getenv("UNFAIRPLAY_SIDB_PRODUCER_SPN_OUTPUT"); outputPath != "" {
		var payload bytes.Buffer
		payload.WriteString("FPSPK1\x00\x00")
		if err := binary.Write(&payload, binary.LittleEndian, uint32(1)); err != nil {
			t.Fatal(err)
		}
		writeAffine := func(value frameAffine128) {
			payload.Write(value.constant[:])
			for _, column := range value.columns {
				payload.Write(column[:])
			}
		}
		writeAffine(bridge)
		payload.WriteByte(byte(len(stages)))
		payload.Write([]byte{0, 0, 0})
		for _, stage := range stages {
			writeAffine(stage.input)
			writeAffine(stage.output)
		}
		checksum := sha256.Sum256(payload.Bytes())
		payload.Write(checksum[:])
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
}

func TestEncodeSIDBProducerKeyReport(t *testing.T) {
	inputPath := os.Getenv("UNFAIRPLAY_SIDB_KEY_REPORT_INPUT")
	outputPath := os.Getenv("UNFAIRPLAY_SIDB_KEY_REPORT_OUTPUT")
	if inputPath == "" || outputPath == "" {
		t.Skip("SIDB key-report encoding is disabled")
	}
	type inputItem struct {
		State string  `json:"state"`
		Frame string  `json:"frame"`
		Key   string  `json:"key"`
		MS    float64 `json:"milliseconds"`
	}
	type outputItem struct {
		inputItem
		Encoded string `json:"encoded_round10"`
	}
	raw, err := os.ReadFile(inputPath)
	if err != nil {
		t.Fatal(err)
	}
	var input []inputItem
	if err := json.Unmarshal(raw, &input); err != nil {
		t.Fatal(err)
	}
	output := make([]outputItem, 0, len(input))
	for _, item := range input {
		keyBytes, err := hex.DecodeString(item.Key)
		if err != nil || len(keyBytes) != 16 {
			t.Fatalf("invalid key %q", item.Key)
		}
		var key [16]byte
		copy(key[:], keyBytes)
		expanded := expandSelectionSIDBAES128Key(key)
		var encoded [16]byte
		for position := range encoded {
			encoded[position] = encodeSelectionSIDBRoundKeyByte(t, selectionSIDBTestRoundKeyDecoder[position], expanded[160+position])
		}
		output = append(output, outputItem{inputItem: item, Encoded: hex.EncodeToString(encoded[:])})
	}
	encoded, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	encoded = append(encoded, '\n')
	if err := os.WriteFile(outputPath, encoded, 0o600); err != nil {
		t.Fatal(err)
	}
}
