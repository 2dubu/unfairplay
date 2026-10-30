package recovery

import "fmt"

// classWorkspaceScheduleState holds the product and current residues for
// the nine-cell workspace. Each update replaces one cell in both planes.
// Initialization and class-input selection are separate stages.
type classWorkspaceScheduleState struct {
	Product [9]uint16
	Current [9]uint16
	Seed    uint16
}

func classWorkspaceScheduleCell(state *classWorkspaceScheduleState, row, column int) error {
	if state == nil || row < 0 || row >= 3 || column < 0 || column >= 3 || state.Seed >= uint16(classWorkspaceModulus) {
		return fmt.Errorf("invalid class workspace cell or seed")
	}
	for i := range 9 {
		if state.Product[i] >= uint16(classWorkspaceModulus) || state.Current[i] >= uint16(classWorkspaceModulus) {
			return fmt.Errorf("invalid class workspace residue")
		}
	}
	position := row*3 + column
	var force uint32
	for other := range 3 {
		if other == row {
			continue
		}
		var difference [3]uint16
		for coordinate := range 3 {
			left := uint32(state.Product[row*3+coordinate])
			right := uint32(state.Product[other*3+coordinate])
			difference[coordinate] = uint16((left + classWorkspaceModulus - right) % classWorkspaceModulus)
		}
		weight := uint32(classWorkspacePowerFactor(difference))
		delta := uint32(state.Product[other*3+column]) - uint32(state.Product[position])
		force += delta * weight // Native arithmetic wraps at 32 bits.
	}
	// The 5-power period is 150 modulo 32767.
	exponent := uint32(state.Product[position]) + uint32(state.Seed) + uint32(state.Current[position])
	power := classWorkspaceModularPower(5, uint16(exponent%150))
	// Force accumulation and this sum wrap at uint32 width before the
	// modulus; widening the arithmetic changes the recurrence.
	forward := (force + uint32(power) + uint32(state.Seed)) % classWorkspaceModulus
	current := (forward + uint32(state.Current[position])) % classWorkspaceModulus
	state.Current[position] = uint16(current)
	state.Product[position] = uint16((uint32(state.Product[position]) + current) % classWorkspaceModulus)
	return nil
}

// classWorkspaceSchedulePass follows the column-major visit order.
func classWorkspaceSchedulePass(state *classWorkspaceScheduleState) error {
	for column := range 3 {
		for row := range 3 {
			if err := classWorkspaceScheduleCell(state, row, column); err != nil {
				return err
			}
		}
	}
	return nil
}
