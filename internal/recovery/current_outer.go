package recovery

import (
	"encoding/binary"
	"fmt"
)

// currentOuterInitialTable is semantic enrollment data, not an instruction or
// protected-register snapshot. Four independently signed v5 licenses with
// different page counts produced this same table exactly.
var currentOuterInitialTable = [outerTableWords]uint32{
	0xb1e42d3b, 0x7d21e1cc, 0x70eaae14, 0xfa0235dc, 0x5c8dfca6, 0x5f192caa,
	0x6050f318, 0x98f98ae2, 0xf5ea5c8b, 0x3da11bc7, 0xf2d6955d, 0xdf5cfb53,
	0xc6f54dfa, 0xbbf0cef6, 0x1bb7a270, 0xc9395d2c, 0xe22ec582, 0xd28bbb8c,
	0x85a72e4d, 0x1c209916, 0x333fab30, 0x9a65bf14, 0x01303b6d, 0xcf7bb0ae,
	0x8363e40a, 0x7f487da6, 0x86986651, 0xc13ca4e6, 0xf1322977, 0x25bb95b9,
	0x0fc06396, 0x9eb24fcb, 0x3e270430, 0x26bf4893, 0x3b2e9fba, 0xa4e85491,
	0x6f2fa819, 0x979fd626, 0x3ad52274, 0x82b47c83, 0x91bb520e, 0x1fcf1f2e,
	0x03d19907, 0x2a0ecec1, 0x8f4f201a, 0x7b795bef, 0x7c976db3, 0x6411d213,
	0xdf965397, 0x0b6c56a7, 0xc80bbd70, 0x9adb84e8, 0xabacece0, 0xd89659c1,
	0xd2a4ee6f, 0x2be7d16d, 0x4d2e0bc0, 0x9c339c92, 0x7dca0da4, 0xd7204d4e,
	0x39f4c55b, 0x6305d02f, 0x700a6f8f, 0x29254d6f, 0xffca57cc, 0x1f99a500,
	0x23e3a447, 0x07dd84db, 0x5d8916c6, 0x5a6ddbc3, 0xe91d9124, 0x2e368e41,
	0xbc7c5645, 0xb07c7380, 0xf84acfe6, 0x2ed5b31e, 0x46ae7bdd, 0x3bd1b83f,
	0x28f41103, 0x0e710c01, 0xc846da3f, 0xb932e1b9, 0xafeb5e3b, 0xa33da0cf,
	0x74e76ce5, 0x2a9d963f, 0xf4f64805, 0xcdb232ea, 0x6f4bbfd0, 0xe3d595db,
	0xfc6bf981, 0x29fecd6d, 0xde0f5314, 0x3b3ccd26, 0x19111329, 0xb40b0978,
	0x5182fd39, 0x080d8471, 0x6dfa885e, 0x9ab721b5, 0x076a33cf, 0xd75a2b6c,
	0x1c1987d9, 0xb91d8220, 0x18453a20, 0x190d2ef3, 0xbd7ebdfa, 0x5c725362,
	0xef8196bb, 0x327dc87f, 0xa4a999bd, 0xf17082cf, 0x631f3b8e, 0x93f073e4,
	0x17dc28ff, 0x66310380, 0xe677d0d6, 0xa4558308, 0x099de2da, 0x03340d8e,
	0x430fc8a5, 0xa7ab71df, 0x31aad011, 0x82716a32, 0xa2841402, 0x56dcbdc7,
	0xfb14d097, 0x72ac63a5, 0x325c7adc, 0xd478e88f, 0x0151af21, 0x803e9710,
	0x17fd43d6, 0xdc12c6d4, 0xba03a905, 0xb33d7a23, 0x91856d97, 0xcdf4ef6f,
	0x513c0c06, 0x5ba2583c, 0x57695b1f, 0x566679e5, 0x2e51cc3e, 0x0c908a0f,
	0x797b6f8e, 0x19862663, 0x086e61d3, 0x61f98df3, 0xcfe134d7, 0x0ae4bc81,
	0x9bb18921, 0xa1cfa419, 0x27b0afab, 0x621fc926, 0xa6624f77, 0x230ea81f,
	0x317d7b93, 0xe294cfa5, 0x0a16a5bc, 0x4d9eb290, 0x9965aef9, 0xdf83a499,
	0x869495f8, 0x1bf724ed, 0xf5ce2c38, 0x0de34be9, 0x966970e2,
}

var currentOuterInitialWords = [5]uint32{
	0x6dc9bf74, 0x31f0c2b6, 0x8310792a, 0xa2c14bc4, 0x3c55e264,
}

const (
	currentOuterInitialOptionalWord = uint32(0x976dec8d)
	currentOuterInitialX1           = uint32(0x60319d83)
	currentOuterInitialCounterSeed  = uint32(0x0a782fec)
)

func currentOuterInputs(container, hardware []byte) ([outerTableWords]uint32, outerTableRoundState, outerIterationControl, error) {
	table := currentOuterInitialTable
	state := outerTableRoundState{
		Words: currentOuterInitialWords, OptionalWord: currentOuterInitialOptionalWord,
		HardwareIndex: 0,
	}
	control := outerIterationControl{X1: currentOuterInitialX1, CounterSeed: currentOuterInitialCounterSeed}
	if len(container) != selectionSIDBContainerBytes {
		return table, state, control, fmt.Errorf("current SIDB container has %d bytes, want %d", len(container), selectionSIDBContainerBytes)
	}
	if binary.BigEndian.Uint32(container[:4]) != currentSIDBContainerVersion {
		return table, state, control, fmt.Errorf("current SIDB container has unsupported version %#08x", binary.BigEndian.Uint32(container[:4]))
	}
	if len(hardware) != len(state.Hardware) {
		return table, state, control, fmt.Errorf("current hardware identity has %d bytes, want %d", len(hardware), len(state.Hardware))
	}
	copy(state.Hardware[:], hardware)
	copy(state.Codec[4:], container[4:selectionSIDBPayloadOffset])
	return table, state, control, nil
}

func deriveCurrentSelectionSource(container, hardware []byte) ([selectionSourceEnvelopeBytes]byte, int, error) {
	var source [selectionSourceEnvelopeBytes]byte
	table, state, control, err := currentOuterInputs(container, hardware)
	if err != nil {
		return source, 0, err
	}
	iterations, ok := applyOuterNaturalLoop(&table, &state, &control)
	if !ok || iterations != outerNaturalIterationCount {
		return source, iterations, fmt.Errorf("typed current outer stage did not reach its natural exit")
	}
	copy(source[:], state.Source[:])
	clear(state.Source[:])
	clear(state.Writer[:])
	clear(state.Tail[:])
	clear(state.SelectionState[:])
	clear(table[:])
	return source, iterations, nil
}
