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
		WordMasks:  [4]uint32{0x10203040, 0x50607080, 0x90a0b0c0, 0xd0e0f000},
		ClassValue: [8]byte{0x01, 0x12, 0x23, 0x34, 0x45, 0x56, 0x67, 0x78},
	}
	candidates, err := deriveClassPrecursorCandidates(registry, parameters)
	if err != nil {
		t.Fatal(err)
	}
	if len(candidates) != len(registry.Records) {
		t.Fatalf("candidates = %d, want %d", len(candidates), len(registry.Records))
	}
	for index, candidate := range candidates {
		if candidate.RecordIndex != index || candidate.Record != registry.Records[index] {
			t.Fatalf("candidate %d lost record identity", index)
		}
		want, err := deriveCurrentSINFPrecursorFromRecord(registry.Records[index].Payload, parameters.WordMasks, parameters.classKey())
		if err != nil {
			t.Fatal(err)
		}
		if candidate.Precursor != want {
			t.Fatalf("candidate %d precursor = %x, want %x", index, candidate.Precursor, want)
		}
	}
}
