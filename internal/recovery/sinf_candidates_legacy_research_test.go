package recovery

func currentClassEnrollmentParameters() classEnrollmentParameters {
	return classEnrollmentParameters{ClassValue: currentClassSecondaryValue}
}

// deriveClassPrecursorCandidates evaluates every structurally valid decoded
// SIDB record. A later authenticated SINF check must accept exactly one; this
// function deliberately does not reproduce the native flattened selector.
func deriveClassPrecursorCandidates(registry decodedSIDBRegistry, parameters classEnrollmentParameters) ([]classPrecursorCandidate, error) {
	result := make([]classPrecursorCandidate, 0, len(registry.Records)*(maximumClassWindowOffset+1))
	key := parameters.classKey()
	for index, record := range registry.Records {
		for windowOffset := 0; windowOffset <= classRecordMaximumWindowOffset(record); windowOffset++ {
			precursor, err := deriveCurrentSINFPrecursorFromRecordWindow(record, windowOffset, key)
			if err != nil {
				return nil, err
			}
			result = append(result, classPrecursorCandidate{
				RecordIndex:  index,
				WindowOffset: windowOffset,
				Record:       record,
				Precursor:    precursor,
			})
		}
	}
	return result, nil
}

func deriveCurrentClassPrecursorCandidates(registry decodedSIDBRegistry) ([]classPrecursorCandidate, error) {
	return deriveClassPrecursorCandidates(registry, currentClassEnrollmentParameters())
}
