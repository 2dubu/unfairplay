package recovery

import "fmt"

const (
	decodedSIDBRecordBytes    = 40
	decodedSIDBMetadataBytes  = 20
	decodedSIDBPayloadBytes   = 16
	decodedSIDBTailBytes      = 4
	decodedSIDBFooterBytes    = 24
	maximumDecodedSIDBRecords = 4096
)

// decodedSIDBRecord is the build-independent geometry of one current SIDB
// registry entry. Metadata selects and authenticates the entry; Payload is the
// 16-byte input to the typed valu transform; Tail remains opaque until its
// semantic role is proved.
type decodedSIDBRecord struct {
	Metadata [decodedSIDBMetadataBytes]byte
	Payload  [decodedSIDBPayloadBytes]byte
	Tail     [decodedSIDBTailBytes]byte
}

type decodedSIDBRegistry struct {
	Records []decodedSIDBRecord
	Footer  [decodedSIDBFooterBytes]byte
}

// parseDecodedSIDBRegistry parses only authenticated, already-decoded SIDB
// bytes. It does not guess a native object address or select an app record.
func parseDecodedSIDBRegistry(data []byte) (decodedSIDBRegistry, error) {
	if len(data) <= decodedSIDBFooterBytes ||
		(len(data)-decodedSIDBFooterBytes)%decodedSIDBRecordBytes != 0 {
		return decodedSIDBRegistry{}, fmt.Errorf("decoded SIDB registry has invalid geometry: %d bytes", len(data))
	}
	count := (len(data) - decodedSIDBFooterBytes) / decodedSIDBRecordBytes
	if count == 0 || count > maximumDecodedSIDBRecords {
		return decodedSIDBRegistry{}, fmt.Errorf("decoded SIDB registry has invalid record count: %d", count)
	}
	result := decodedSIDBRegistry{Records: make([]decodedSIDBRecord, count)}
	for index := range result.Records {
		offset := index * decodedSIDBRecordBytes
		record := &result.Records[index]
		copy(record.Metadata[:], data[offset:offset+decodedSIDBMetadataBytes])
		copy(record.Payload[:], data[offset+decodedSIDBMetadataBytes:offset+decodedSIDBMetadataBytes+decodedSIDBPayloadBytes])
		copy(record.Tail[:], data[offset+decodedSIDBMetadataBytes+decodedSIDBPayloadBytes:offset+decodedSIDBRecordBytes])
	}
	copy(result.Footer[:], data[len(data)-decodedSIDBFooterBytes:])
	return result, nil
}
