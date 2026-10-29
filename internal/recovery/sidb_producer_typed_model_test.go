package recovery

type sidbProducerTyped struct {
	firstTargets  [16]byte
	firstOutput   [16]frameByteLinear
	roundKeys     [8][16]byte
	lastTargets   [16]byte
	lastKeys      [16]byte
	finalXOR      [16]uint16
	finalOutput   [16]frameByteLinear
	finalConstant [16]byte
}

// apply evaluates the producer as ten direct AES inverse-S-box layers. Nine
// fixed InvShiftRows-plus-InvMixColumns operations separate the first nine
// layers; the last layer feeds four byte-wise XOR graphs and output maps.
func (model sidbProducerTyped) apply(input [16]byte) (output [16]byte) {
	var state [16]byte
	for source, target := range model.firstTargets {
		state[target] = model.firstOutput[target].apply(frameInverseSBox[input[source]])
	}
	state = applySIDBCanonicalInverseRound(state)
	for _, key := range model.roundKeys {
		for position := range state {
			state[position] = frameInverseSBox[state[position]^key[position]]
		}
		state = applySIDBCanonicalInverseRound(state)
	}
	var last [16]byte
	for source, target := range model.lastTargets {
		last[target] = frameInverseSBox[state[source]^model.lastKeys[source]]
	}
	for target, mask := range model.finalXOR {
		var value byte
		for source := range 16 {
			if mask&(1<<uint(source)) != 0 {
				value ^= last[source]
			}
		}
		output[target] = model.finalOutput[target].apply(value) ^ model.finalConstant[target]
	}
	return output
}
