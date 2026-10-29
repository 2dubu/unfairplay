package recovery

import "testing"

func TestClassWorkspaceScheduleCellZeroControl(t *testing.T) {
	var state classWorkspaceScheduleState
	if err := classWorkspaceScheduleCell(&state, 0, 0); err != nil {
		t.Fatal(err)
	}
	if state.Product[0] != 1 || state.Current[0] != 1 || state.Forward[0] != 1 {
		t.Fatal("zero control did not produce expected unit residues")
	}
	for i := 1; i < 9; i++ {
		if state.Product[i] != 0 || state.Current[i] != 0 || state.Forward[i] != 0 {
			t.Fatal("cell updated a neighboring position")
		}
	}
}
func TestClassWorkspaceScheduleRejectsInvalidState(t *testing.T) {
	var state classWorkspaceScheduleState
	if err := classWorkspaceScheduleCell(nil, 0, 0); err == nil {
		t.Fatal("accepted nil state")
	}
	if err := classWorkspaceScheduleCell(&state, 3, 0); err == nil {
		t.Fatal("accepted invalid row")
	}
	state.Product[4] = uint16(classWorkspaceModulus)
	if err := classWorkspaceSchedulePass(&state); err == nil {
		t.Fatal("accepted invalid residue")
	}
}
