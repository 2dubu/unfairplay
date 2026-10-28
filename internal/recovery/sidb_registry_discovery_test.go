package recovery

import "testing"

func TestDiscoverDecodedSIDBRegistries(t *testing.T) {
	const first = uint64(0x710100000000)
	const duplicate = first + pageSize
	const unrelated = duplicate + pageSize
	data := make([]byte, 3*decodedSIDBRecordBytes+decodedSIDBFooterBytes)
	for index := range data {
		data[index] = byte(index*37 + 11)
	}
	pages := map[uint64][]byte{
		first:     make([]byte, pageSize),
		duplicate: make([]byte, pageSize),
		unrelated: make([]byte, pageSize),
	}
	copy(pages[first], data)
	copy(pages[duplicate], data)
	copy(pages[unrelated], data)
	pages[unrelated][len(data)+7] = 1
	s := newState(&model{})
	s.mapped = pages

	candidates, err := discoverDecodedSIDBRegistries(s)
	if err != nil {
		t.Fatal(err)
	}
	if len(candidates) != 1 {
		t.Fatalf("registry candidates = %d, want 1", len(candidates))
	}
	if len(candidates[0].Registry.Records) != 3 || len(candidates[0].Data) != len(data) {
		t.Fatalf("discovered registry has wrong geometry")
	}
	for index := range data {
		if candidates[0].Data[index] != data[index] {
			t.Fatalf("discovered registry differs at byte %d", index)
		}
	}
}

func TestDiscoverDecodedSIDBRegistriesRejectsUnboundedCandidates(t *testing.T) {
	s := newState(&model{})
	s.mapped = make(map[uint64][]byte)
	for candidate := 0; candidate <= maximumDecodedSIDBRegistryCandidates; candidate++ {
		page := make([]byte, pageSize)
		size := decodedSIDBRecordBytes + decodedSIDBFooterBytes
		for index := 0; index < size; index++ {
			page[index] = byte(index + candidate + 1)
		}
		s.mapped[uint64(0x710100000000)+uint64(candidate)*pageSize] = page
	}
	if _, err := discoverDecodedSIDBRegistries(s); err == nil {
		t.Fatal("unbounded decoded registry candidates were accepted")
	}
}
