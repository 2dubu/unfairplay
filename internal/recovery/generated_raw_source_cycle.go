// Code generated from the compact macOS recovery profile; DO NOT EDIT.
package recovery

func runRawSourceCycleBoundary(s *state, _ fault) error {
	profile, err := executionProfileForModel(s.model)
	if err != nil {
		return err
	}
	if err := withRecoveryStage(s, "outer.cycle", func() error { return profile.runCycle(s, fault{}) }); err != nil {
		return err
	}
	if s.x[9] == profile.cycleFinal {
		return withRecoveryStage(s, "outer.cycle-final", func() error { return profile.runCycleFinal(s, fault{}) })
	}
	return nil
}
