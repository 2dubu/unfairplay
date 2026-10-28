package recovery

// enrollmentFrameCandidate is one real state slot and its portable frame
// transform. Native selection chooses one candidate; portable recovery can try
// each candidate against the authenticated SIDB/SINF structure instead.
type enrollmentFrameCandidate struct {
	Slot  int
	State [16]byte
	Frame [16]byte
	Key   [16]byte
}

func enrollmentFrameCandidates(source []byte) ([]enrollmentFrameCandidate, error) {
	states, realSlots, err := transformSelectionStateSlots(source)
	if err != nil {
		return nil, err
	}
	result := make([]enrollmentFrameCandidate, 0, realSlots)
	for slot := 0; slot < realSlots; slot++ {
		frame, err := frameTransform(states[slot])
		if err != nil {
			return nil, err
		}
		result = append(result, enrollmentFrameCandidate{
			Slot: slot, State: states[slot], Frame: frame, Key: deriveSelectionSIDBKey(states[slot]),
		})
	}
	return result, nil
}
