package recovery

import "math/bits"

type frameByteLinear [8]byte

func (value frameByteLinear) apply(input byte) (output byte) {
	for bit := range 8 {
		if input&(1<<uint(bit)) != 0 {
			output ^= value[bit]
		}
	}
	return output
}

func (value frameByteLinear) invertible() bool {
	var seen [256]bool
	for input := range 256 {
		output := value.apply(byte(input))
		if seen[output] {
			return false
		}
		seen[output] = true
	}
	return true
}

func frameAffine128Rank(value frameAffine128) int {
	var basis [128][2]uint64
	rank := 0
	for _, column := range value.columns {
		var vector [2]uint64
		for outputByte, byteValue := range column {
			word := outputByte / 8
			vector[word] |= uint64(byteValue) << uint((outputByte%8)*8)
		}
		for vector != [2]uint64{} {
			pivot := 0
			if vector[1] != 0 {
				pivot = 64 + bits.Len64(vector[1]) - 1
			} else {
				pivot = bits.Len64(vector[0]) - 1
			}
			if basis[pivot] == [2]uint64{} {
				basis[pivot] = vector
				rank++
				break
			}
			vector[0] ^= basis[pivot][0]
			vector[1] ^= basis[pivot][1]
		}
	}
	return rank
}

func frameAESMixInvertible(value [16][16]byte) bool {
	matrix := value
	for column := range 16 {
		pivot := column
		for pivot < 16 && matrix[pivot][column] == 0 {
			pivot++
		}
		if pivot == 16 {
			return false
		}
		matrix[column], matrix[pivot] = matrix[pivot], matrix[column]
		inverse := frameGFInverse(matrix[column][column])
		for index := range 16 {
			matrix[column][index] = frameGFMul(matrix[column][index], inverse)
		}
		for row := range 16 {
			if row == column || matrix[row][column] == 0 {
				continue
			}
			factor := matrix[row][column]
			for index := range 16 {
				matrix[row][index] ^= frameGFMul(matrix[column][index], factor)
			}
		}
	}
	return true
}

type frameAESModel struct {
	first        frameAffine128
	firstPost    [16]frameByteLinear
	mixes        [9][16][16]byte
	roundKeys    [8][16]byte
	roundMasks   [8][16]byte
	lastPre      [16]frameByteLinear
	lastConstant [16]byte
}

func applyFrameAESMix(matrix [16][16]byte, input [16]byte) (output [16]byte) {
	for target := range output {
		for source, coefficient := range matrix[target] {
			if coefficient != 0 {
				output[target] ^= frameGFMul(coefficient, input[source])
			}
		}
	}
	return output
}

// apply returns the semantic frame suffix immediately before the affine bridge
// into the SIDB-key producer. The former protected implementation used ten
// AES S-box layers separated by arbitrary binary encodings; those encodings
// are normalized here to nine GF(2^8) mixing layers and eight keyed S-box
// transitions.
func (model frameAESModel) apply(input [16]byte) [16]byte {
	state := frameSBoxState(model.first.apply(input))
	for position := range state {
		state[position] = model.firstPost[position].apply(state[position])
	}
	state = applyFrameAESMix(model.mixes[0], state)
	for round := range model.roundKeys {
		for position := range state {
			state[position] = frameSBox[state[position]^model.roundKeys[round][position]] ^ model.roundMasks[round][position]
		}
		state = applyFrameAESMix(model.mixes[round+1], state)
	}
	for position := range state {
		state[position] = model.lastPre[position].apply(state[position]) ^ model.lastConstant[position]
	}
	return frameSBoxState(state)
}
