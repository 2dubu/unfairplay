package recovery

// enrollmentKeyCandidate is one real state slot and its derived SIDB key.
// Native selection chooses one candidate; portable recovery can try each
// candidate against the authenticated SIDB/SINF structure instead.
type enrollmentKeyCandidate struct {
	Slot  int
	State [16]byte
	Key   [16]byte
}

func enrollmentKeyCandidates(source []byte) ([]enrollmentKeyCandidate, error) {
	states, realSlots, err := transformSelectionStateSlots(source)
	if err != nil {
		return nil, err
	}
	result := make([]enrollmentKeyCandidate, 0, realSlots)
	for slot := 0; slot < realSlots; slot++ {
		result = append(result, enrollmentKeyCandidate{
			Slot: slot, State: states[slot], Key: deriveSelectionSIDBKey(states[slot]),
		})
	}
	return result, nil
}
