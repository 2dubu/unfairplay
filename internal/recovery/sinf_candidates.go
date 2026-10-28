package recovery

// classEnrollmentParameters are the enrollment-global inputs shared by every
// record in one decoded SIDB registry. Their producer remains a separate typed
// stage; candidate evaluation does not depend on how they were discovered.
type classEnrollmentParameters struct {
	WordMasks  [4]uint32
	ClassValue [8]byte
}

func (parameters classEnrollmentParameters) classKey() (key [16]byte) {
	copy(key[:8], parameters.ClassValue[:])
	copy(key[8:], parameters.ClassValue[:])
	return key
}

type classPrecursorCandidate struct {
	RecordIndex int
	Record      decodedSIDBRecord
	Precursor   [16]byte
}

// deriveClassPrecursorCandidates evaluates every structurally valid decoded
// SIDB record. A later authenticated SINF check must accept exactly one; this
// function deliberately does not reproduce the native flattened selector.
func deriveClassPrecursorCandidates(registry decodedSIDBRegistry, parameters classEnrollmentParameters) ([]classPrecursorCandidate, error) {
	result := make([]classPrecursorCandidate, 0, len(registry.Records))
	key := parameters.classKey()
	for index, record := range registry.Records {
		precursor, err := deriveCurrentSINFPrecursorFromRecord(record.Payload, parameters.WordMasks, key)
		if err != nil {
			return nil, err
		}
		result = append(result, classPrecursorCandidate{
			RecordIndex: index,
			Record:      record,
			Precursor:   precursor,
		})
	}
	return result, nil
}
