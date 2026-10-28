package recovery

import (
	"crypto/sha256"
	"fmt"
	"math/bits"
	"sort"
	"testing"
)

type linearByteMap struct {
	basis [8]byte
}

func (value linearByteMap) apply(input byte) (output byte) {
	for bit := range 8 {
		if input&(1<<bit) != 0 {
			output ^= value.basis[bit]
		}
	}
	return output
}

type partialLinearByteMap struct {
	inputs  [8]byte
	outputs [8]byte
	rank    int
}

func (value *partialLinearByteMap) add(input, output byte) bool {
	for input != 0 {
		pivot := bits.Len8(input) - 1
		if value.inputs[pivot] == 0 {
			value.inputs[pivot] = input
			value.outputs[pivot] = output
			value.rank++
			return true
		}
		input ^= value.inputs[pivot]
		output ^= value.outputs[pivot]
	}
	return output == 0
}

func (value partialLinearByteMap) complete() (linearByteMap, bool) {
	if value.rank != 8 {
		return linearByteMap{}, false
	}
	var result linearByteMap
	for bit := range 8 {
		input, output := byte(1<<bit), byte(0)
		for input != 0 {
			pivot := bits.Len8(input) - 1
			if value.inputs[pivot] == 0 {
				return linearByteMap{}, false
			}
			input ^= value.inputs[pivot]
			output ^= value.outputs[pivot]
		}
		result.basis[bit] = output
	}
	return result, true
}

func uniqueFourDifference(table [256]byte) ([256]byte, error) {
	var result [256]byte
	for inputDelta := 1; inputDelta < 256; inputDelta++ {
		var counts [256]uint16
		for input := range 256 {
			counts[table[input]^table[input^inputDelta]]++
		}
		found := -1
		for outputDelta, count := range counts {
			if count != 4 {
				continue
			}
			if found >= 0 {
				return result, fmt.Errorf("delta %#02x has multiple four-count outputs", inputDelta)
			}
			found = outputDelta
		}
		if found < 0 {
			return result, fmt.Errorf("delta %#02x has no four-count output", inputDelta)
		}
		result[inputDelta] = byte(found)
	}
	return result, nil
}

func solveLinearDifferenceEquivalence(reference, target [256]byte) (linearByteMap, linearByteMap, uint64, error) {
	referenceFour, err := uniqueFourDifference(reference)
	if err != nil {
		return linearByteMap{}, linearByteMap{}, 0, err
	}
	targetFour, err := uniqueFourDifference(target)
	if err != nil {
		return linearByteMap{}, linearByteMap{}, 0, err
	}
	var inputMap, outputMap linearByteMap
	var inputBasis [8]byte
	var inputSpan [256]byte
	var nodes uint64
	var search func(int, [256]byte, partialLinearByteMap) bool
	search = func(depth int, span [256]byte, partial partialLinearByteMap) bool {
		spanSize := 1 << depth
		var used [256]bool
		for _, item := range span[:spanSize] {
			used[item] = true
		}
		for candidate := 1; candidate < 256; candidate++ {
			if used[candidate] {
				continue
			}
			nodes++
			nextPartial := partial
			valid := true
			for suffix := 0; suffix < spanSize; suffix++ {
				referenceInput := span[suffix] ^ byte(candidate)
				targetInput := byte(spanSize + suffix)
				if !nextPartial.add(referenceFour[referenceInput], targetFour[targetInput]) {
					valid = false
					break
				}
			}
			if !valid {
				continue
			}
			inputBasis[depth] = byte(candidate)
			nextSpan := span
			for suffix := 0; suffix < spanSize; suffix++ {
				nextSpan[spanSize+suffix] = span[suffix] ^ byte(candidate)
			}
			if depth == 7 {
				completed, ok := nextPartial.complete()
				if !ok {
					continue
				}
				inputMap = linearByteMap{basis: inputBasis}
				outputMap = completed
				return true
			}
			if search(depth+1, nextSpan, nextPartial) {
				return true
			}
		}
		return false
	}
	if !search(0, inputSpan, partialLinearByteMap{}) {
		return inputMap, outputMap, nodes, fmt.Errorf("no linear difference equivalence")
	}
	return inputMap, outputMap, nodes, nil
}

func solveAffineSBox(reference, target [256]byte) (linearByteMap, linearByteMap, byte, byte, uint64, error) {
	inputMap, outputMap, nodes, err := solveLinearDifferenceEquivalence(reference, target)
	if err != nil {
		return inputMap, outputMap, 0, 0, nodes, err
	}
	for inputOffset := range 256 {
		outputOffset := target[0] ^ outputMap.apply(reference[inputOffset])
		matched := true
		for input := range 256 {
			want := outputMap.apply(reference[inputMap.apply(byte(input))^byte(inputOffset)]) ^ outputOffset
			if target[input] != want {
				matched = false
				break
			}
		}
		if matched {
			return inputMap, outputMap, byte(inputOffset), outputOffset, nodes, nil
		}
	}
	return inputMap, outputMap, 0, 0, nodes, fmt.Errorf("difference equivalence has no affine S-box lift")
}

func affineRelation(source, target [256]byte) (linearByteMap, byte, error) {
	var inverse [256]byte
	var seen [256]bool
	for input, value := range source {
		if seen[value] {
			return linearByteMap{}, 0, fmt.Errorf("source is not a permutation")
		}
		seen[value] = true
		inverse[value] = byte(input)
	}
	var relation [256]byte
	for value := range 256 {
		relation[value] = target[inverse[value]]
	}
	offset := relation[0]
	var linear linearByteMap
	for bit := range 8 {
		linear.basis[bit] = relation[1<<bit] ^ offset
	}
	for value := range 256 {
		if relation[value] != linear.apply(byte(value))^offset {
			return linearByteMap{}, 0, fmt.Errorf("relation is not affine")
		}
	}
	return linear, offset, nil
}

type linear32Inverse struct {
	rows [32]uint32
}

func invertLinear32(columns [32]uint32) (linear32Inverse, bool) {
	var rows [32]uint64
	for row := range 32 {
		for column := range 32 {
			if columns[column]&(1<<row) != 0 {
				rows[row] |= 1 << column
			}
		}
		rows[row] |= uint64(1) << (32 + row)
	}
	for column := range 32 {
		pivot := -1
		for row := column; row < 32; row++ {
			if rows[row]&(1<<column) != 0 {
				pivot = row
				break
			}
		}
		if pivot < 0 {
			return linear32Inverse{}, false
		}
		rows[column], rows[pivot] = rows[pivot], rows[column]
		for row := range 32 {
			if row != column && rows[row]&(1<<column) != 0 {
				rows[row] ^= rows[column]
			}
		}
	}
	var result linear32Inverse
	for row := range 32 {
		result.rows[row] = uint32(rows[row] >> 32)
	}
	return result, true
}

func (value linear32Inverse) apply(input uint32) (output uint32) {
	for bit, row := range value.rows {
		if bits.OnesCount32(row&input)&1 != 0 {
			output |= 1 << bit
		}
	}
	return output
}

type frameFinalGroup struct {
	inputs    []int
	outputs   []int
	decode    [16][256]byte
	constant  uint32
	forward   [32]uint32
	inverse   linear32Inverse
	canonical [16][256]byte
}

func frameFinalComponents(model frameSemanticModel) ([]frameFinalGroup, error) {
	var assigned [16]bool
	var groups []frameFinalGroup
	for start := range model.final {
		if assigned[start] {
			continue
		}
		outputSet := map[int]bool{start: true}
		inputSet := make(map[int]bool)
		changed := true
		for changed {
			changed = false
			for output := range outputSet {
				entry := model.final[output]
				for slot := 0; slot < int(entry.inputCount); slot++ {
					input := int(entry.inputs[slot])
					if !inputSet[input] {
						inputSet[input] = true
						changed = true
					}
				}
			}
			for output, entry := range model.final {
				if outputSet[output] {
					continue
				}
				for slot := 0; slot < int(entry.inputCount); slot++ {
					if inputSet[int(entry.inputs[slot])] {
						outputSet[output] = true
						changed = true
						break
					}
				}
			}
		}
		group := frameFinalGroup{}
		for input := range inputSet {
			group.inputs = append(group.inputs, input)
		}
		for output := range outputSet {
			group.outputs = append(group.outputs, output)
			assigned[output] = true
		}
		sort.Ints(group.inputs)
		sort.Ints(group.outputs)
		if len(group.inputs) != 4 || len(group.outputs) != 4 {
			return nil, fmt.Errorf("final component has %d inputs and %d outputs", len(group.inputs), len(group.outputs))
		}
		for _, input := range group.inputs {
			found := false
			var table [256]byte
			for _, output := range group.outputs {
				entry := model.final[output]
				if entry.inputCount != 3 {
					continue
				}
				for slot := 0; slot < 3; slot++ {
					if int(entry.inputs[slot]) == input {
						table = entry.contributions[slot]
						found = true
						break
					}
				}
				if found {
					break
				}
			}
			if !found {
				return nil, fmt.Errorf("final input %d has no shared contribution", input)
			}
			inputMap, _, inputOffset, _, _, err := solveAffineSBox(frameSBox, table)
			if err != nil {
				return nil, fmt.Errorf("final input %d: %w", input, err)
			}
			for value := range 256 {
				group.canonical[input][value] = frameSBox[inputMap.apply(byte(value))^inputOffset]
			}
		}
		var columns [32]uint32
		for outputSlot, output := range group.outputs {
			entry := model.final[output]
			if entry.inputCount == 1 {
				input := int(entry.inputs[0])
				inputSlot := sort.SearchInts(group.inputs, input)
				if inputSlot >= 4 || group.inputs[inputSlot] != input {
					return nil, fmt.Errorf("final direct input %d is outside its group", input)
				}
				var seen [256]bool
				for value, encoded := range entry.direct {
					if seen[encoded] {
						return nil, fmt.Errorf("final direct output %d is not invertible", output)
					}
					seen[encoded] = true
					group.decode[output][encoded] = group.canonical[input][value]
				}
				for bit := range 8 {
					columns[inputSlot*8+bit] |= uint32(1) << (outputSlot*8 + bit)
				}
				continue
			}
			var encodeInverse [256]byte
			var seen [256]bool
			for value, encoded := range entry.encode {
				if seen[encoded] {
					return nil, fmt.Errorf("final outer output %d is not invertible", output)
				}
				seen[encoded] = true
				encodeInverse[encoded] = byte(value)
			}
			for raw := range 256 {
				group.decode[output][raw] = encodeInverse[byte(raw)^entry.baseline]
			}
			for slot := 0; slot < 3; slot++ {
				input := int(entry.inputs[slot])
				inputSlot := sort.SearchInts(group.inputs, input)
				linear, offset, err := affineRelation(group.canonical[input], entry.contributions[slot])
				if err != nil {
					return nil, fmt.Errorf("final output %d input %d: %w", output, input, err)
				}
				group.constant ^= uint32(offset) << (outputSlot * 8)
				for bit := range 8 {
					columns[inputSlot*8+bit] |= uint32(linear.apply(1<<bit)) << (outputSlot * 8)
				}
			}
		}
		inverse, ok := invertLinear32(columns)
		if !ok {
			return nil, fmt.Errorf("final component is not linearly invertible")
		}
		group.inverse = inverse
		group.forward = columns
		groups = append(groups, group)
	}
	if len(groups) != 4 {
		return nil, fmt.Errorf("final layer has %d components", len(groups))
	}
	return groups, nil
}

func applyFrameFinal(model frameSemanticModel, state [16]byte) (output [16]byte) {
	for target, entry := range model.final {
		if entry.inputCount == 1 {
			output[target] = entry.direct[state[entry.inputs[0]]]
			continue
		}
		decoded := byte(0)
		for slot := 0; slot < 3; slot++ {
			decoded ^= entry.contributions[slot][state[entry.inputs[slot]]]
		}
		output[target] = entry.baseline ^ entry.encode[decoded]
	}
	return output
}

func TestFrameFinalGroupDecoder(t *testing.T) {
	groups, err := frameFinalComponents(frameSemantic26A428)
	if err != nil {
		t.Fatal(err)
	}
	for index := range 4096 {
		digest := sha256.Sum256([]byte{0x46, 0x49, 0x4e, byte(index >> 8), byte(index)})
		var state [16]byte
		copy(state[:], digest[:16])
		raw := applyFrameFinal(frameSemantic26A428, state)
		for _, group := range groups {
			var encoded uint32
			for outputSlot, output := range group.outputs {
				encoded |= uint32(group.decode[output][raw[output]]) << (outputSlot * 8)
			}
			decoded := group.inverse.apply(encoded ^ group.constant)
			for inputSlot, input := range group.inputs {
				got := byte(decoded >> (inputSlot * 8))
				want := group.canonical[input][state[input]]
				if got != want {
					t.Fatalf("control %d final input %d = %#02x, want %#02x", index, input, got, want)
				}
			}
		}
	}
}
