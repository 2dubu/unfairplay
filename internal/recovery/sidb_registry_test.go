package recovery

import (
	"bytes"
	"testing"
)

func TestParseDecodedSIDBRegistry(t *testing.T) {
	data := syntheticSIDBRecordTree(2, 29, 7)
	registry, err := parseDecodedSIDBRegistry(data, 2)
	if err != nil {
		t.Fatal(err)
	}
	if !registry.Tree || len(registry.Records) != 2 {
		t.Fatalf("unexpected tree records: %+v", registry)
	}
	for i, record := range registry.Records {
		offset := 22 + i*36
		if record.SourceOffset != offset+20 || !bytes.Equal(record.Metadata[:], data[offset:offset+20]) || !bytes.Equal(record.Payload[:], data[offset+20:offset+36]) || record.TailAvailable != 0 {
			t.Fatalf("record %d has incorrect geometry or bytes", i)
		}
	}
}

func TestParseDecodedSIDBRegistryRejectsUnknownGeometry(t *testing.T) {
	for _, size := range []int{0, 24, 25, 63, 104, 576, 11680} {
		if _, err := parseDecodedSIDBRegistry(make([]byte, size), 2); err == nil {
			t.Fatalf("accepted %d-byte invalid tree", size)
		}
	}
}

// Flat-looking sizes previously selected a historical interpretation before
// tree geometry was understood. They are not evidence of a separate format.
func TestParseDecodedSIDBRegistryRejectsSyntheticFlatShapes(t *testing.T) {
	for _, size := range []int{2*40 + 24, 14*40 + 16, 292 * 40} {
		data := make([]byte, size)
		for i := range data {
			data[i] = byte(i*17 + 5)
		}
		if _, err := parseDecodedSIDBRegistry(data, 2); err == nil {
			t.Fatalf("accepted synthetic flat shape of %d bytes", size)
		}
	}
}
