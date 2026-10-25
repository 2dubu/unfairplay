// Code generated from the compact macOS recovery profile; DO NOT EDIT.
package recovery

func runRawSourceCycleBoundary(s *state, _ fault) error {
	if err := runCompactKernel(s, compact26A428Cycle); err != nil {
		return err
	}
	if s.x[9] == 0x10046119C {
		return runRawSourceCycleFinal(s, fault{})
	}
	return nil
}
