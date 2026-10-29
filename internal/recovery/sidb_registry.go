package recovery

import "fmt"

const (
	decodedSIDBRecordBytes           = 40
	decodedSIDBMetadataBytes         = 20
	decodedSIDBPayloadBytes          = 16
	decodedSIDBTailBytes             = 4
	decodedSIDBNoSeparateFooterBytes = 0
	decodedSIDBOlderFooterBytes      = 16
	decodedSIDBFooterBytes           = 24
	maximumDecodedSIDBRecords        = 4096
)

// decodedSIDBRecord is the build-independent geometry of one current SIDB
// registry entry. Metadata is opaque native selection/validation state;
// Payload is the 16-byte input to the typed valu transform; Tail remains
// opaque until its semantic role is proved.
type decodedSIDBRecord struct {
	Metadata [decodedSIDBMetadataBytes]byte
	Payload  [decodedSIDBPayloadBytes]byte
	Tail     [decodedSIDBTailBytes]byte
}

type decodedSIDBRegistry struct {
	Records     []decodedSIDBRecord
	Footer      [decodedSIDBFooterBytes]byte
	FooterBytes int
}

// parseDecodedSIDBRegistry parses only authenticated, already-decoded SIDB
// bytes. It does not guess a native object address or select an app record.
func parseDecodedSIDBRegistry(data []byte) (decodedSIDBRegistry, error) {
	footerBytes := -1
	for _, candidate := range [...]int{decodedSIDBNoSeparateFooterBytes, decodedSIDBOlderFooterBytes, decodedSIDBFooterBytes} {
		if len(data) > candidate && (len(data)-candidate)%decodedSIDBRecordBytes == 0 {
			if footerBytes != -1 {
				return decodedSIDBRegistry{}, fmt.Errorf("decoded SIDB registry has ambiguous geometry: %d bytes", len(data))
			}
			footerBytes = candidate
		}
	}
	if footerBytes < 0 {
		return decodedSIDBRegistry{}, fmt.Errorf("decoded SIDB registry has invalid geometry: %d bytes", len(data))
	}
	count := (len(data) - footerBytes) / decodedSIDBRecordBytes
	if count == 0 || count > maximumDecodedSIDBRecords {
		return decodedSIDBRegistry{}, fmt.Errorf("decoded SIDB registry has invalid record count: %d", count)
	}
	result := decodedSIDBRegistry{Records: make([]decodedSIDBRecord, count), FooterBytes: footerBytes}
	for index := range result.Records {
		offset := index * decodedSIDBRecordBytes
		record := &result.Records[index]
		copy(record.Metadata[:], data[offset:offset+decodedSIDBMetadataBytes])
		copy(record.Payload[:], data[offset+decodedSIDBMetadataBytes:offset+decodedSIDBMetadataBytes+decodedSIDBPayloadBytes])
		copy(record.Tail[:], data[offset+decodedSIDBMetadataBytes+decodedSIDBPayloadBytes:offset+decodedSIDBRecordBytes])
	}
	if footerBytes != 0 {
		copy(result.Footer[:], data[len(data)-footerBytes:])
	}
	return result, nil
}
