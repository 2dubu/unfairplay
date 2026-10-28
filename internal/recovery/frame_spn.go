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

type frameSPNModel struct {
	initial           frameAffine128
	stages            []frameSPNStage
	finalInput        frameAffine128
	finalOutput       frameAffine128
	outputEncoding    [16][256]byte
	exitInputs        [4]byte
	exitBaseline      [4]byte
	exitContributions [4][4][256]byte
}

func frameSBoxState(input [16]byte) (output [16]byte) {
	for index, value := range input {
		output[index] = frameSBox[value]
	}
	return output
}

func (model frameSPNModel) applyInternal(input [16]byte) (state [16]byte, exitWord uint32) {
	state = model.initial.apply(input)
	for stageIndex, stage := range model.stages {
		if stageIndex == 6 {
			for outputByte := range 4 {
				value := model.exitBaseline[outputByte]
				for slot, position := range model.exitInputs {
					value ^= model.exitContributions[outputByte][slot][state[position]]
				}
				exitWord |= uint32(value) << uint(outputByte*8)
			}
		}
		state = stage.output.apply(frameSBoxState(stage.input.apply(state)))
	}
	state = model.finalOutput.apply(frameSBoxState(model.finalInput.apply(state)))
	return state, exitWord
}

func (model frameSPNModel) apply(input [16]byte) (output [16]byte, exitWord uint32) {
	state, exitWord := model.applyInternal(input)
	for position, value := range state {
		output[position] = model.outputEncoding[position][value]
	}
	return output, exitWord
}
