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

// parseDecodedSIDBRegistry prefers the tree interpretation and otherwise
// returns a flat compatibility interpretation. Neither is authenticated.
func parseDecodedSIDBRegistry(data []byte) (decodedSIDBRegistry, error) {
	nodes, err := parseDecodedSIDBTree(data)
	if err != nil {
		return parseLegacyDecodedSIDBRegistry(data)
	}
	result := decodedSIDBRegistry{Tree: true}
	rootEnd := nodes[0].Length
	for _, node := range nodes {
		if node.ChildCount != 0 || node.Length != decodedSIDBNodeHeaderBytes+decodedSIDBPayloadBytes {
			continue
		}
		var record decodedSIDBRecord
		record.SourceOffset = node.Offset + decodedSIDBNodeHeaderBytes
		record.MaskOffset = (record.SourceOffset + 12) % len(currentClassRecordMask)
		copy(record.Metadata[:], data[node.Offset:record.SourceOffset])
		copy(record.Payload[:], data[record.SourceOffset:record.SourceOffset+decodedSIDBPayloadBytes])
		tailStart := record.SourceOffset + decodedSIDBPayloadBytes
		tailEnd := min(tailStart+decodedSIDBTailBytes, rootEnd)
		record.TailAvailable = copy(record.Tail[:], data[tailStart:tailEnd])
		result.Records = append(result.Records, record)
	}
	return result, nil
}

// decodedSIDBRegistryInterpretations retains each structurally valid view for
// subsequent SINF authentication; neither grammar proves the selected format.
func decodedSIDBRegistryInterpretations(data []byte) []decodedSIDBRegistry {
	var registries []decodedSIDBRegistry
	if registry, err := parseDecodedSIDBRegistry(data); err == nil && registry.Tree {
		registries = append(registries, registry)
	}
	if registry, err := parseLegacyDecodedSIDBRegistry(data); err == nil {
		registries = append(registries, registry)
	}
	return registries
}

// The upstream flat interpretation is a compatibility candidate. Its
// geometry alone does not authenticate a record or identify a native branch.
func parseLegacyDecodedSIDBRegistry(data []byte) (decodedSIDBRegistry, error) {
	footerBytes := -1
	for _, candidate := range [...]int{decodedSIDBNoSeparateFooterBytes, decodedSIDBOlderFooterBytes, decodedSIDBFooterBytes} {
		if len(data) > candidate && (len(data)-candidate)%decodedSIDBRecordBytes == 0 {
			if footerBytes != -1 {
				return decodedSIDBRegistry{}, fmt.Errorf("ambiguous SIDB registry geometry")
			}
			footerBytes = candidate
		}
	}
	if footerBytes < 0 {
		return decodedSIDBRegistry{}, fmt.Errorf("invalid SIDB registry geometry")
	}
	count := (len(data) - footerBytes) / decodedSIDBRecordBytes
	if count == 0 || count > maximumDecodedSIDBRecords {
		return decodedSIDBRegistry{}, fmt.Errorf("invalid SIDB registry count")
	}
	result := decodedSIDBRegistry{Records: make([]decodedSIDBRecord, count), FooterBytes: footerBytes}
	for i := range result.Records {
		offset := i * decodedSIDBRecordBytes
		r := &result.Records[i]
		r.SourceOffset = offset + decodedSIDBMetadataBytes
		r.TailAvailable = decodedSIDBTailBytes
		copy(r.Metadata[:], data[offset:r.SourceOffset])
		copy(r.Payload[:], data[r.SourceOffset:r.SourceOffset+decodedSIDBPayloadBytes])
		copy(r.Tail[:], data[r.SourceOffset+decodedSIDBPayloadBytes:offset+decodedSIDBRecordBytes])
	}
	copy(result.Footer[:], data[len(data)-footerBytes:])
	return result, nil
}
