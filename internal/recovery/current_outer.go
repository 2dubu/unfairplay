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

// selectionSourceFromOuterState reconstructs the source bytes consumed by
// selection from the outer state.
func selectionSourceFromOuterState(state *outerTableRoundState, sourceBytes int) ([]byte, error) {
	if state == nil {
		return nil, fmt.Errorf("current outer state is unavailable")
	}
	if sourceBytes < selectionSourcePrefixBytes+16 || sourceBytes > selectionSourcePrefixBytes+16*maximumMappedSelectionModulus || (sourceBytes-selectionSourcePrefixBytes)%16 != 0 {
		return nil, fmt.Errorf("unsupported current selection source size %d", sourceBytes)
	}
	// The mapped continuation follows Source with eight Writer bytes and the
	// tail beginning at offset 40. It is bounded to the currently mapped profile.
	var full [selectionSourcePrefixBytes + 16*maximumMappedSelectionModulus]byte
	copy(full[:], state.Source[:])
	copy(full[224:232], state.Writer[132:140])
	copy(full[232:], state.Tail[40:128])
	return append([]byte(nil), full[:sourceBytes]...), nil
}

// nativeOuterStopStep evaluates the cycle-final dispatch at 0x10046119c.
// The accumulator is initialized at 0x100845fc8 and updated only when the
// cycle-final path runs, beginning with cycle 14.
func nativeOuterStopStep(words [5]uint32, accumulator *uint32) bool {
	const target uint32 = 0x05010000
	w11 := target ^ uint32(0x050183f6)
	w9 := words[3] ^ uint32(0x189c788b) ^ words[2]
	w10 := words[1]
	w9 |= w10 ^ uint32(0xdc6e37ea)
	w12 := words[4] ^ uint32(0xc7c208dc)
	w14 := words[2] ^ uint32(0x00000c09)
	w13 := words[2] ^ uint32(0x399b2000)
	w13 = uint32((uint64(w14)<<32 | uint64(w13)) >> 13)
	w14 = w11 + uint32(0xfa8e2a8c)
	w10 = (w14 ^ w10) & w12 & w13
	w9 |= w10
	w10 = uint32(0x0fc20444) & (w9 << 1)
	w9 = w9 - w10 + uint32(0x87e10222)
	w10 = (*accumulator ^ uint32(0x5db37cfe)) & (w9 ^ uint32(0x87e10222))
	w9 = (w9 & uint32(0x4d191d94)) ^ w10
	*accumulator = w9 ^ uint32(0x15ab616a)
	return w9 == target
}

// deriveNativeSelectionSource follows the native codec cycle and stop
// dispatch. The codec boundary is determined by the stateful native exit,
// then the remaining bytes are checked as aligned body and trailer.
func deriveNativeSelectionSource(container, hardware []byte) ([]byte, currentSIDBLayout, int, error) {
	var layout currentSIDBLayout
	if len(container) < 4+13+selectionSIDBExtendedTrailerBytes || binary.BigEndian.Uint32(container[:4]) != currentSIDBContainerVersion {
		return nil, layout, 0, fmt.Errorf("current SIDB container has an unsupported version or truncated header")
	}
	if len(hardware) != 6 {
		return nil, layout, 0, fmt.Errorf("current hardware identity has invalid length")
	}
	table := currentOuterInitialTable
	state := outerTableRoundState{Words: currentOuterInitialWords, OptionalWord: currentOuterInitialOptionalWord}
	copy(state.Hardware[:], hardware)
	const maximumOuterIterations = 4096
	const maximumCodecBytes = (maximumOuterIterations + 2) / 11
	available := min(maximumCodecBytes, len(container)-4-selectionSIDBExtendedTrailerBytes)
	copy(state.Codec[4:17], container[4:17])
	defer func() {
		clear(table[:])
		clear(state.Source[:])
		clear(state.Writer[:])
		clear(state.Tail[:])
		clear(state.SelectionState[:])
		clear(state.Codec[:])
		clear(state.Hardware[:])
	}()
	control := outerIterationControl{X1: currentOuterInitialX1, CounterSeed: currentOuterInitialCounterSeed}
	base := outerTableInitialBase
	accumulator := uint32(0xaf549e95)
	cycles := 0
	for iteration := 1; iteration <= 11*available-2; iteration++ {
		next, ok := applyOuterFullIteration(&table, base, &state, &control)
		if !ok {
			return nil, layout, iteration - 1, fmt.Errorf("typed current outer stage failed")
		}
		base = next
		if iteration >= 9 && (iteration-9)%11 == 0 {
			cycles++
			if !applyOuterCycleWithByte(&state, container[4+cycles-1]) {
				return nil, layout, iteration, fmt.Errorf("typed current outer cycle failed")
			}
			control.X1 = 0x60319d87
			if cycles >= 14 && nativeOuterStopStep(state.Words, &accumulator) {
				bodyOffset := 4 + cycles
				bodyBytes := len(container) - bodyOffset - selectionSIDBExtendedTrailerBytes
				layout = currentSIDBLayout{CodecBytes: cycles, PayloadOffset: bodyOffset, PayloadBytes: bodyBytes, OuterIterations: iteration}
				if bodyBytes <= 0 || bodyBytes%selectionSIDBBlockBytes != 0 || bodyBytes > maximumSelectionSIDBPayloadBytes {
					return nil, layout, iteration, fmt.Errorf("current SIDB has invalid body geometry")
				}
				count := uint32(0x0976dec8) ^ (state.OptionalWord >> 4)
				if count == 0 || count > maximumMappedSelectionModulus {
					return nil, layout, iteration, fmt.Errorf("current SIDB selection modulus is outside mapped profile")
				}
				rootCursor, err := nativeSIDBRootCursor(hardware)
				if err != nil {
					return nil, layout, iteration, err
				}
				layout.RootCursor = rootCursor
				layout.SelectorModulus = int(count)
				layout.SourceBytes = selectionSourcePrefixBytes + 16*layout.SelectorModulus
				layout.SelectedSlot = int(container[bodyOffset-1]) % layout.SelectorModulus
				source, err := selectionSourceFromOuterState(&state, layout.SourceBytes)
				return source, layout, iteration, err
			}
		}
	}
	return nil, layout, 11*available - 2, fmt.Errorf("current SIDB codec did not reach the native exit within its bounded input")
}
