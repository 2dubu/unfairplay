package recovery

type sidbBytePermutationTerm struct {
	target byte
	linear frameByteLinear
}

type sidbBytePermutationAffine struct {
	constant [16]byte
	terms    [16]sidbBytePermutationTerm
}

func (value sidbBytePermutationAffine) apply(input [16]byte) (output [16]byte) {
	output = value.constant
	for source, term := range value.terms {
		output[term.target] ^= term.linear.apply(input[source])
	}
	return output
}

type sidbProducerCanonicalAES struct {
	first     sidbSparseAffine
	firstPost sidbBytePermutationAffine
	roundKeys [8][16]byte
	lastPre   sidbBytePermutationAffine
	final     sidbSparseAffine
}

var frameInverseSBox = func() (result [256]byte) {
	for input, output := range frameSBox {
		result[output] = byte(input)
	}
	return result
}()

func applySIDBCanonicalInverseRound(input [16]byte) (output [16]byte) {
	coefficients := [4][4]byte{
		{0x0e, 0x0b, 0x0d, 0x09},
		{0x09, 0x0e, 0x0b, 0x0d},
		{0x0d, 0x09, 0x0e, 0x0b},
		{0x0b, 0x0d, 0x09, 0x0e},
	}
	for column := range 4 {
		for row := range 4 {
			for sourceRow := range 4 {
				sourceColumn := (column + sourceRow) % 4
				source := sourceColumn*4 + sourceRow
				output[column*4+row] ^= frameGFMul(coefficients[row][sourceRow], input[source])
			}
		}
	}
	return output
}

// apply evaluates the producer in one fixed AES state order. All nine middle
// layers are the standard InvShiftRows-plus-InvMixColumns transform and all
// eight nonlinear layers are keyed AES inverse S-boxes.
func (model sidbProducerCanonicalAES) apply(input [16]byte) [16]byte {
	state := model.firstPost.apply(frameSBoxState(model.first.apply(input)))
	state = applySIDBCanonicalInverseRound(state)
	for _, key := range model.roundKeys {
		for position := range state {
			state[position] = frameInverseSBox[state[position]^key[position]]
		}
		state = applySIDBCanonicalInverseRound(state)
	}
	state = model.lastPre.apply(state)
	return model.final.apply(frameSBoxState(state))
}

func deriveSelectionSIDBKey(state [16]byte) [16]byte {
	frameState := currentFrameAES.apply(frameCore(state))
	return currentSIDBProducerCanonicalAES.apply(frameState)
}
