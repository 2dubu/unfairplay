package recovery

import "testing"

func TestDeriveClassPrecursorCandidates(t *testing.T) {
	data := make([]byte, 3*decodedSIDBRecordBytes+decodedSIDBFooterBytes)
	for index := range data {
		data[index] = byte(index*17 + 3)
	}
	registry, err := parseDecodedSIDBRegistry(data)
	if err != nil {
		t.Fatal(err)
	}
	parameters := classEnrollmentParameters{
		ClassValue: [8]byte{0x01, 0x12, 0x23, 0x34, 0x45, 0x56, 0x67, 0x78},
	}
	candidates, err := deriveClassPrecursorCandidates(registry, parameters)
	if err != nil {
		t.Fatal(err)
	}
	wantCount := len(registry.Records) * (maximumClassWindowOffset + 1)
	if len(candidates) != wantCount {
		t.Fatalf("candidates = %d, want %d", len(candidates), wantCount)
	}
	for index, candidate := range candidates {
		recordIndex := index / (maximumClassWindowOffset + 1)
		windowOffset := index % (maximumClassWindowOffset + 1)
		if candidate.RecordIndex != recordIndex || candidate.WindowOffset != windowOffset || candidate.Record != registry.Records[recordIndex] {
			t.Fatalf("candidate %d lost record/window identity", index)
		}
		want, err := deriveCurrentSINFPrecursorFromRecordWindow(registry.Records[recordIndex], windowOffset, parameters.classKey())
		if err != nil {
			t.Fatal(err)
		}
		if candidate.Precursor != want {
			t.Fatalf("candidate %d precursor = %x, want %x", index, candidate.Precursor, want)
		}
	}
}
