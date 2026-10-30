package recovery

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

// decodedSIDBRecord holds a candidate 16-byte payload and its surrounding
// bytes. Tree leaf roles and mask alignment require native selection or
// authentication; a structurally valid candidate is not authenticated.
type decodedSIDBRecord struct {
	Metadata      [decodedSIDBMetadataBytes]byte
	Payload       [decodedSIDBPayloadBytes]byte
	Tail          [decodedSIDBTailBytes]byte
	SourceOffset  int // Absolute offset of Payload in the decoded body.
	TailAvailable int // Number of bytes copied into Tail.
	MaskOffset    int // Mask phase at the start of Payload.
}

type decodedSIDBRegistry struct {
	Records     []decodedSIDBRecord
	Footer      [decodedSIDBFooterBytes]byte
	FooterBytes int
	Tree        bool
}

// parseDecodedSIDBRegistry exposes candidate payloads only from validated tree
// geometry. Structural validity does not authenticate a selected record.
func parseDecodedSIDBRegistry(data []byte, rootCursor int) (decodedSIDBRegistry, error) {
	nodes, err := parseDecodedSIDBTree(data, rootCursor)
	if err != nil {
		return decodedSIDBRegistry{}, err
	}
	result := decodedSIDBRegistry{Tree: true}
	for _, node := range nodes {
		if node.ChildCount != 0 || node.Length != decodedSIDBNodeHeaderBytes+decodedSIDBPayloadBytes {
			continue
		}
		var record decodedSIDBRecord
		record.SourceOffset = node.Offset + decodedSIDBNodeHeaderBytes
		record.MaskOffset = (record.SourceOffset + 12) % len(currentClassRecordMask)
		copy(record.Metadata[:], data[node.Offset:record.SourceOffset])
		copy(record.Payload[:], data[record.SourceOffset:record.SourceOffset+decodedSIDBPayloadBytes])

		result.Records = append(result.Records, record)
	}
	return result, nil
}
