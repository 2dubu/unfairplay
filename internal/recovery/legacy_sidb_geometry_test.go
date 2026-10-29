package recovery

import (
	"encoding/binary"
	"fmt"
)

const (
	selectionSIDBTrailerBytes          = 16
	outerNaturalIterationCount         = 317
	selectionSIDBContainerBytes        = 993
	selectionSIDBExtendedBytes         = 1003
	selectionSIDBPayloadOffset         = 33
	selectionSIDBExtendedPayloadOffset = 43
	selectionSIDBPayloadBytes          = 944
	selectionSIDBOuterIterations       = outerNaturalIterationCount
	selectionSIDBExtendedIterations    = 427
)

// Historical fixed-size helpers serve only the original synthetic oracle tests.
func currentSIDBLayoutForContainer(input []byte) (currentSIDBLayout, error) {
	if len(input) < 4 || binary.BigEndian.Uint32(input[:4]) != currentSIDBContainerVersion {
		return currentSIDBLayout{}, fmt.Errorf("current SIDB container has an unsupported version")
	}
	switch len(input) {
	case selectionSIDBContainerBytes:
		return currentSIDBLayout{
			CodecBytes:      selectionSIDBPayloadOffset - 4,
			PayloadOffset:   selectionSIDBPayloadOffset,
			PayloadBytes:    selectionSIDBPayloadBytes,
			SourceBytes:     selectionSourceEnvelopeBytes,
			OuterIterations: selectionSIDBOuterIterations,
		}, nil
	case selectionSIDBExtendedBytes:
		return currentSIDBLayout{
			CodecBytes:      selectionSIDBExtendedPayloadOffset - 4,
			PayloadOffset:   selectionSIDBExtendedPayloadOffset,
			PayloadBytes:    selectionSIDBPayloadBytes,
			SourceBytes:     selectionSourceEnvelopeBytes,
			OuterIterations: selectionSIDBExtendedIterations,
		}, nil
	default:
		return currentSIDBLayout{}, fmt.Errorf("current SIDB container has unsupported size %d", len(input))
	}
}

func decodeSelectionSIDBContainerDynamic(input []byte) ([]byte, error) {
	layout, err := currentSIDBLayoutForContainer(input)
	if err != nil {
		return nil, err
	}
	return decodeSelectionSIDBContainerWithLayout(input, layout)
}

func decodeSelectionSIDBContainer(input []byte) ([selectionSIDBPayloadBytes]byte, error) {
	var output [selectionSIDBPayloadBytes]byte
	decoded, err := decodeSelectionSIDBContainerDynamic(input)
	if err != nil {
		return output, err
	}
	defer clear(decoded)
	if len(decoded) != len(output) {
		return output, fmt.Errorf("selection SIDB body has %d bytes, want %d", len(decoded), len(output))
	}
	copy(output[:], decoded)
	return output, nil
}

func decryptSelectionSIDBPayload(input [selectionSIDBPayloadBytes]byte, key [selectionSIDBBlockBytes]byte) ([selectionSIDBPayloadBytes]byte, error) {
	var output [selectionSIDBPayloadBytes]byte
	decoded, err := decryptSelectionSIDBPayloadDynamic(input[:], key)
	if err != nil {
		return output, err
	}
	copy(output[:], decoded)
	clear(decoded)
	return output, nil
}

func currentOuterInputs(container, hardware []byte) ([outerTableWords]uint32, outerTableRoundState, outerIterationControl, error) {
	table := currentOuterInitialTable
	state := outerTableRoundState{
		Words: currentOuterInitialWords, OptionalWord: currentOuterInitialOptionalWord,
		HardwareIndex: 0,
	}
	control := outerIterationControl{X1: currentOuterInitialX1, CounterSeed: currentOuterInitialCounterSeed}
	layout, err := currentSIDBLayoutForContainer(container)
	if err != nil {
		return table, state, control, err
	}
	if len(hardware) != len(state.Hardware) {
		return table, state, control, fmt.Errorf("current hardware identity has %d bytes, want %d", len(hardware), len(state.Hardware))
	}
	copy(state.Hardware[:], hardware)
	copy(state.Codec[4:4+layout.CodecBytes], container[4:4+layout.CodecBytes])
	return table, state, control, nil
}

func deriveCurrentSelectionSourceDynamic(container, hardware []byte) ([]byte, int, error) {
	layout, err := currentSIDBLayoutForContainer(container)
	if err != nil {
		return nil, 0, err
	}
	table, state, control, err := currentOuterInputs(container, hardware)
	if err != nil {
		return nil, 0, err
	}
	defer func() {
		clear(state.Source[:])
		clear(state.Writer[:])
		clear(state.Tail[:])
		clear(state.SelectionState[:])
		clear(state.Codec[:])
		clear(state.Hardware[:])
		clear(table[:])
	}()
	iterations, ok := applyOuterLoop(&table, &state, &control, layout.OuterIterations)
	if !ok || iterations != layout.OuterIterations {
		return nil, iterations, fmt.Errorf("typed current outer stage did not reach its natural exit")
	}
	source, err := selectionSourceFromOuterState(&state, layout.SourceBytes)
	return source, iterations, err
}

// Retained for callers that consume the original 224-byte source.
func deriveCurrentSelectionSource(container, hardware []byte) ([selectionSourceEnvelopeBytes]byte, int, error) {
	var source [selectionSourceEnvelopeBytes]byte
	dynamic, iterations, err := deriveCurrentSelectionSourceDynamic(container, hardware)
	if err != nil {
		return source, iterations, err
	}
	if len(dynamic) != len(source) {
		clear(dynamic)
		return source, iterations, fmt.Errorf("selection source has unsupported size for fixed wrapper")
	}
	copy(source[:], dynamic)
	clear(dynamic)
	return source, iterations, nil
}

// Historical fixed-iteration oracle, excluded from production recovery.
func applyOuterNaturalLoop(table *[outerTableWords]uint32, state *outerTableRoundState, control *outerIterationControl) (int, bool) {
	return applyOuterLoop(table, state, control, outerNaturalIterationCount)
}

func applyOuterLoop(table *[outerTableWords]uint32, state *outerTableRoundState, control *outerIterationControl, iterations int) (int, bool) {
	if table == nil || state == nil || control == nil || iterations < 1 {
		return 0, false
	}
	base := outerTableInitialBase
	cycleIndex := uint32(0)
	for iteration := 1; iteration <= iterations; iteration++ {
		next, ok := applyOuterFullIteration(table, base, state, control)
		if !ok {
			return iteration - 1, false
		}
		base = next
		if iteration >= 9 && (iteration-9)%11 == 0 {
			cycleIndex++
			if !applyOuterCycle(state, cycleIndex) {
				return iteration, false
			}
			control.X1 = 0x60319d87
		}
	}
	return iterations, true
}
