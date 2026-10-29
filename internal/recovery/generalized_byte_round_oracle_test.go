package recovery

// generalizedByteRound is a typed 16-byte substitution/diffusion stage. Each
// input byte contributes four bytes to fixed output positions; the sixteen
// contributions are XORed with one baseline. This is the semantic shape used
// by the recovered FairPlay frame and SIDB producer rounds after register,
// stack, and MBA encodings are removed.
type generalizedByteRound struct {
	baseline [16]byte
	targets  [16][4]byte
	tables   [16][256][4]byte
}

func (round generalizedByteRound) apply(input [16]byte) (output [16]byte) {
	output = round.baseline
	for position, value := range input {
		for slot, target := range round.targets[position] {
			output[target] ^= round.tables[position][value][slot]
		}
	}
	return output
}
