package recovery

// classEnrollmentParameters are the enrollment-global inputs shared by every
// record in one decoded SIDB registry. Candidate evaluation does not depend on
// the native path that produced them.
type classEnrollmentParameters struct {
	ClassValue [8]byte
}

var currentClassSecondaryValue = [8]byte{0x4a, 0xc6, 0xfc, 0xf2, 0xd3, 0x92, 0x5e, 0x05}

func currentClassEnrollmentParameters() classEnrollmentParameters {
	return classEnrollmentParameters{ClassValue: currentClassSecondaryValue}
}

func (parameters classEnrollmentParameters) classKey() (key [16]byte) {
	copy(key[:8], parameters.ClassValue[:])
	copy(key[8:], parameters.ClassValue[:])
	return key
}

type classPrecursorCandidate struct {
	RecordIndex  int
	WindowOffset int
	Record       decodedSIDBRecord
	Precursor    [16]byte
}

// deriveClassPrecursorCandidates evaluates every structurally valid decoded
// SIDB record. A later authenticated SINF check must accept exactly one; this
// function deliberately does not reproduce the native flattened selector.
func deriveClassPrecursorCandidates(registry decodedSIDBRegistry, parameters classEnrollmentParameters) ([]classPrecursorCandidate, error) {
	result := make([]classPrecursorCandidate, 0, len(registry.Records)*(maximumClassWindowOffset+1))
	key := parameters.classKey()
	for index, record := range registry.Records {
		for windowOffset := 0; windowOffset <= maximumClassWindowOffset; windowOffset++ {
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
