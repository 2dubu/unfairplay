package recovery

var currentClassSecondaryValue = [8]byte{0x4a, 0xc6, 0xfc, 0xf2, 0xd3, 0x92, 0x5e, 0x05}

type classPrecursorCandidate struct {
	RecordIndex  int
	WindowOffset int
	Record       decodedSIDBRecord
	Precursor    [16]byte
}

// classEnrollmentParameters holds the repeated-half class value.
type classEnrollmentParameters struct {
	ClassValue [8]byte
}

func (parameters classEnrollmentParameters) classKey() (key [16]byte) {
	copy(key[:8], parameters.ClassValue[:])
	copy(key[8:], parameters.ClassValue[:])
	return key
}
