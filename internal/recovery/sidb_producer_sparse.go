package recovery

type sidbSparseTerm struct {
	target byte
	linear frameByteLinear
}

// sidbSparseAffine is an affine 128-bit layer whose byte-level dependency
// graph is sparse. It keeps the producer's actual byte transforms without the
// separate input and output encodings used by the protected implementation.
type sidbSparseAffine struct {
	constant [16]byte
	terms    [16][]sidbSparseTerm
}

func (value sidbSparseAffine) apply(input [16]byte) (output [16]byte) {
	output = value.constant
	for source, byteValue := range input {
		for _, term := range value.terms[source] {
			output[term.target] ^= term.linear.apply(byteValue)
		}
	}
	return output
}

type sidbProducerSparse struct {
	first  sidbSparseAffine
	rounds [9]sidbSparseAffine
	final  sidbSparseAffine
}

// apply evaluates the producer as ten ordinary AES S-box layers separated by
// the sparse byte-linear layers obtained after composing away every protected
// stage boundary.
func (model sidbProducerSparse) apply(input [16]byte) [16]byte {
	state := frameSBoxState(model.first.apply(input))
	for _, round := range model.rounds {
		state = frameSBoxState(round.apply(state))
	}
	return model.final.apply(state)
}

func deriveSelectionSIDBKey(state [16]byte) [16]byte {
	frameState := currentFrameAES.apply(frameCore(state))
	return currentSIDBProducerSparse.apply(frameState)
}
