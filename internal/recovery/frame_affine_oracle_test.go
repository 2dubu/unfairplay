package recovery

type frameAffine128 struct {
	constant [16]byte
	columns  [128][16]byte
}

func (value frameAffine128) linear(input [16]byte) (output [16]byte) {
	for bit, column := range value.columns {
		if input[bit/8]&(1<<uint(bit%8)) == 0 {
			continue
		}
		for index := range output {
			output[index] ^= column[index]
		}
	}
	return output
}

func (value frameAffine128) apply(input [16]byte) (output [16]byte) {
	output = value.linear(input)
	for index := range output {
		output[index] ^= value.constant[index]
	}
	return output
}

type frameSPNStage struct {
	input  frameAffine128
	output frameAffine128
}

func frameSBoxState(input [16]byte) (output [16]byte) {
	for index, value := range input {
		output[index] = frameSBox[value]
	}
	return output
}
