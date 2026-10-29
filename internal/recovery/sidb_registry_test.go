package recovery

import (
	"bytes"
	"testing"
)

func TestParseDecodedSIDBRegistry(t *testing.T) {
	data := make([]byte, 2*decodedSIDBRecordBytes+decodedSIDBFooterBytes)
	for index := range data {
		data[index] = byte(index*29 + 7)
	}
	registry, err := parseDecodedSIDBRegistry(data)
	if err != nil {
		t.Fatal(err)
	}
	if len(registry.Records) != 2 {
		t.Fatalf("records = %d, want 2", len(registry.Records))
	}
	if registry.FooterBytes != decodedSIDBFooterBytes {
		t.Fatalf("footer bytes = %d, want %d", registry.FooterBytes, decodedSIDBFooterBytes)
	}
	for index, record := range registry.Records {
		offset := index * decodedSIDBRecordBytes
		if !bytes.Equal(record.Metadata[:], data[offset:offset+decodedSIDBMetadataBytes]) {
			t.Fatalf("record %d metadata mismatch", index)
		}
		if !bytes.Equal(record.Payload[:], data[offset+decodedSIDBMetadataBytes:offset+decodedSIDBMetadataBytes+decodedSIDBPayloadBytes]) {
			t.Fatalf("record %d payload mismatch", index)
		}
		if !bytes.Equal(record.Tail[:], data[offset+decodedSIDBMetadataBytes+decodedSIDBPayloadBytes:offset+decodedSIDBRecordBytes]) {
			t.Fatalf("record %d tail mismatch", index)
		}
	}
	if !bytes.Equal(registry.Footer[:], data[len(data)-decodedSIDBFooterBytes:]) {
		t.Fatal("footer mismatch")
	}
}

func TestParseDecodedSIDBRegistryRejectsUnknownGeometry(t *testing.T) {
	for _, size := range []int{0, decodedSIDBFooterBytes, decodedSIDBFooterBytes + 1, decodedSIDBFooterBytes + decodedSIDBRecordBytes - 1} {
		if _, err := parseDecodedSIDBRegistry(make([]byte, size)); err == nil {
			t.Fatalf("accepted %d-byte registry", size)
		}
	}
	tooMany := make([]byte, (maximumDecodedSIDBRecords+1)*decodedSIDBRecordBytes+decodedSIDBFooterBytes)
	if _, err := parseDecodedSIDBRegistry(tooMany); err == nil {
		t.Fatal("accepted excessive record count")
	}
}

func TestParseDecodedOlderSIDBRegistry(t *testing.T) {
	data := make([]byte, 14*decodedSIDBRecordBytes+decodedSIDBOlderFooterBytes)
	for index := range data {
		data[index] = byte(index*17 + 5)
	}
	registry, err := parseDecodedSIDBRegistry(data)
	if err != nil {
		t.Fatal(err)
	}
	if len(registry.Records) != 14 {
		t.Fatalf("older records = %d, want 14", len(registry.Records))
	}
	if registry.FooterBytes != decodedSIDBOlderFooterBytes {
		t.Fatalf("older footer bytes = %d, want %d", registry.FooterBytes, decodedSIDBOlderFooterBytes)
	}
	if !bytes.Equal(registry.Footer[:decodedSIDBOlderFooterBytes], data[len(data)-decodedSIDBOlderFooterBytes:]) || !bytes.Equal(registry.Footer[decodedSIDBOlderFooterBytes:], make([]byte, decodedSIDBFooterBytes-decodedSIDBOlderFooterBytes)) {
		t.Fatal("older footer mismatch")
	}
}

func TestParseDecodedFooterlessSIDBRegistry(t *testing.T) {
	data := make([]byte, 292*decodedSIDBRecordBytes)
	for index := range data {
		data[index] = byte(index*11 + 9)
	}
	registry, err := parseDecodedSIDBRegistry(data)
	if err != nil {
		t.Fatal(err)
	}
	if len(registry.Records) != 292 || registry.FooterBytes != decodedSIDBNoFooterBytes {
		t.Fatalf("footerless registry = %d records, %d footer bytes", len(registry.Records), registry.FooterBytes)
	}
	if registry.Footer != [decodedSIDBFooterBytes]byte{} {
		t.Fatal("footerless registry retained footer data")
	}
}
