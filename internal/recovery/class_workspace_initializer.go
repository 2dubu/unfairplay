package recovery

// classWorkspaceInitialize builds the native product and current
// planes from one big-endian account-digest word. The Forward plane is
// scratch: the first scheduler pass replaces every cell before it is read.
// The exponent is in 1..32767; the value 32767 must not be reduced to zero.
func classWorkspaceInitialize(word uint16) classWorkspaceScheduleState {
	exponent := uint16(uint32(word)%classWorkspaceModulus + 1)
	state := classWorkspaceScheduleState{Seed: uint16(uint32(word) % classWorkspaceModulus)}
	for i := range state.Product {
		base := uint16((uint32(exponent) + uint32(i)) % classWorkspaceModulus)
		state.Product[i] = classWorkspaceModularPower(base, exponent)
		state.Current[i] = state.Seed
	}
	return state
}
