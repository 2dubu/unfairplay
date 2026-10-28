package recovery

type sidbInverseAESTerm struct {
	source      byte
	coefficient byte
}

type sidbInverseAESMix [16][4]sidbInverseAESTerm

func (mix sidbInverseAESMix) apply(input [16]byte) (output [16]byte) {
	for target, terms := range mix {
		for _, term := range terms {
			output[target] ^= frameGFMul(term.coefficient, input[term.source])
		}
	}
	return output
}

type sidbProducerInverseAES struct {
	first        sidbSparseAffine
	firstPost    [16]frameByteLinear
	mixes        [9]sidbInverseAESMix
	roundKeys    [8][16]byte
	lastPre      [16]frameByteLinear
	lastConstant [16]byte
	final        sidbSparseAffine
}

// apply evaluates the producer after removing its protected byte encodings.
// The middle is eight keyed AES inverse-S-box layers separated by ordinary
// InvMixColumns coefficient matrices. The endpoint maps retain the external
// frame and SIDB-key coordinate systems.
func (model sidbProducerInverseAES) apply(input [16]byte) [16]byte {
	state := frameSBoxState(model.first.apply(input))
	for position := range state {
		state[position] = model.firstPost[position].apply(state[position])
	}
	state = model.mixes[0].apply(state)
	for round, key := range model.roundKeys {
		for position := range state {
			state[position] = frameInverseSBox[state[position]^key[position]]
		}
		state = model.mixes[round+1].apply(state)
	}
	for position := range state {
		state[position] = model.lastPre[position].apply(state[position]) ^ model.lastConstant[position]
	}
	return model.final.apply(frameSBoxState(state))
}

var frameInverseSBox = func() (result [256]byte) {
	for input, output := range frameSBox {
		result[output] = byte(input)
	}
	return result
}()

func deriveSelectionSIDBKey(state [16]byte) [16]byte {
	frameState := currentFrameAES.apply(frameCore(state))
	return currentSIDBProducerInverseAES.apply(frameState)
}
