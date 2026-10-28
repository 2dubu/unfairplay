package recovery

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
