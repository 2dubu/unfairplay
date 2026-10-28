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

func TestDeriveClassPrecursorCandidatesFromAccountDigest(t *testing.T) {
	data := make([]byte, decodedSIDBRecordBytes+decodedSIDBFooterBytes)
	for index := range data {
		data[index] = byte(index*29 + 7)
	}
	registry, err := parseDecodedSIDBRegistry(data)
	if err != nil {
		t.Fatal(err)
	}
	digest := [20]byte{0x41, 0x43, 0x43, 0x4f, 0x55, 0x4e, 0x54, 0x2d, 0x44, 0x49, 0x47, 0x45, 0x53, 0x54, 0xa5, 0x5a, 0x13, 0x37, 0x24, 0x68}
	got, err := deriveClassPrecursorCandidatesFromAccountDigest(registry, digest)
	if err != nil {
		t.Fatal(err)
	}
	want, err := deriveClassPrecursorCandidates(registry, classEnrollmentParametersFromAccountDigest(digest))
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != len(want) {
		t.Fatalf("candidate count = %d, want %d", len(got), len(want))
	}
	for index := range got {
		if got[index] != want[index] {
			t.Fatalf("candidate %d differs", index)
		}
	}
}
