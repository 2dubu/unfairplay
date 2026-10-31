package recovery

import "fmt"

// taggedSIDBRoute is a structural selection only; it neither authenticates the
// record nor establishes that a real enrollment follows this route.
type taggedSIDBRoute struct {
	Version        sidbSchemaVersion
	AccountVersion sidbSchemaVersion
	TypeCursor     int
	TypeLength     int
	SignCursor     int
	SignLength     int
	ValueCursor    int
	ValueLength    int
	RecordIndex    int
	Record         decodedSIDBRecord
	Window         [decodedSIDBPayloadBytes]byte
}

// selectTaggedSIDBRoute selects a bounded decoded tree path. The caller
// supplies the hardware-derived tag decoder and decoded body.
func selectTaggedSIDBRoute(decode sidbTagDecoder, sinf, body []byte, rootCursor int) (taggedSIDBRoute, error) {
	var result taggedSIDBRoute
	parsed, err := parseRecoverySINF(sinf)
	if err != nil {
		return result, fmt.Errorf("invalid SINF: %w", err)
	}
	tree, err := parseDecodedSIDBTaggedTree(body, rootCursor, decode)
	if err != nil {
		return result, err
	}
	if tree.Nodes[0].Tag != "sean" {
		return result, fmt.Errorf("unsupported SIDB root tag")
	}
	unsupported := func(parent int, tags ...string) error {
		for _, i := range tree.children[parent] {
			for _, tag := range tags {
				if tree.Nodes[i].Tag == tag {
					return fmt.Errorf("unsupported SIDB route tag %q", tag)
				}
			}
		}
		return nil
	}
	if err := unsupported(0, "divi"); err != nil {
		return result, err
	}
	unique := func(parent int, tag string, id *uint32) (*decodedSIDBTaggedNode, error) {
		var match *decodedSIDBTaggedNode
		for _, i := range tree.children[parent] {
			n := &tree.Nodes[i]
			if n.Tag != tag || (id != nil && (!n.HasIdentifier || n.Identifier != *id)) {
				continue
			}
			if match != nil {
				return nil, fmt.Errorf("duplicate SIDB route tag %q", tag)
			}
			match = n
		}
		if match == nil {
			return nil, fmt.Errorf("missing SIDB route tag %q", tag)
		}
		return match, nil
	}
	one := uint32(1)
	version, err := unique(0, "vers", &one)
	if err != nil {
		return result, err
	}
	if version.ChildCount != 0 {
		return result, fmt.Errorf("SIDB version is not a leaf")
	}
	result.Version, err = readSIDBSchemaVersion(body, version.Offset, version.Length)
	if err != nil {
		return taggedSIDBRoute{}, err
	}
	group, err := unique(0, "grup", nil)
	if err != nil {
		return taggedSIDBRoute{}, err
	}
	user, err := unique(tree.byCursor[group.Offset], "user", &parsed.selector.User)
	if err != nil {
		return taggedSIDBRoute{}, err
	}
	accountVersion, err := unique(tree.byCursor[user.Offset], "vers", &one)
	if err != nil {
		return taggedSIDBRoute{}, err
	}
	if accountVersion.ChildCount != 0 {
		return taggedSIDBRoute{}, fmt.Errorf("SIDB account version is not a leaf")
	}
	result.AccountVersion, err = readSIDBSchemaVersion(body, accountVersion.Offset, accountVersion.Length)
	if err != nil {
		return taggedSIDBRoute{}, err
	}
	if result.AccountVersion != sidbSchemaVersion10002 {
		return taggedSIDBRoute{}, fmt.Errorf("unsupported SIDB account version")
	}
	if err := unsupported(tree.byCursor[user.Offset], "gval"); err != nil {
		return taggedSIDBRoute{}, err
	}
	key, err := unique(tree.byCursor[user.Offset], "key ", &parsed.selector.Key)
	if err != nil {
		return taggedSIDBRoute{}, err
	}
	keyIndex := tree.byCursor[key.Offset]
	typeNode, err := unique(keyIndex, "type", &one)
	if err != nil {
		return taggedSIDBRoute{}, err
	}
	if typeNode.ChildCount != 0 || typeNode.Length != 24 {
		return taggedSIDBRoute{}, fmt.Errorf("unsupported SIDB type leaf geometry")
	}
	result.TypeCursor, result.TypeLength = typeNode.Offset, typeNode.Length
	sign, err := unique(keyIndex, "sign", &one)
	if err != nil {
		return taggedSIDBRoute{}, err
	}
	if sign.ChildCount != 0 || sign.Length != 148 {
		return taggedSIDBRoute{}, fmt.Errorf("unsupported SIDB sign leaf geometry")
	}
	result.SignCursor, result.SignLength = sign.Offset, sign.Length
	if err := unsupported(keyIndex, "salt", "iviv", "dong"); err != nil {
		return taggedSIDBRoute{}, err
	}
	value, err := unique(keyIndex, "valu", nil)
	if err != nil {
		return taggedSIDBRoute{}, err
	}
	first, err := tree.firstChild(key.Offset)
	if err != nil || first == nil || first.Offset != value.Offset {
		return taggedSIDBRoute{}, fmt.Errorf("SIDB first key child is not valu")
	}
	if value.ChildCount != 0 || value.Length != decodedSIDBNodeHeaderBytes+decodedSIDBPayloadBytes {
		return taggedSIDBRoute{}, fmt.Errorf("invalid SIDB value leaf")
	}
	result.ValueCursor, result.ValueLength = value.Offset, value.Length
	start := value.Offset + decodedSIDBNodeHeaderBytes
	window := start
	if window > len(body) || len(body)-window < len(result.Window) {
		return taggedSIDBRoute{}, fmt.Errorf("SIDB value window exceeds body")
	}
	registry, err := parseDecodedSIDBRegistry(body, rootCursor)
	if err != nil {
		return taggedSIDBRoute{}, err
	}
	found := false
	for i, record := range registry.Records {
		if record.SourceOffset == start {
			if found {
				return taggedSIDBRoute{}, fmt.Errorf("duplicate SIDB record offset")
			}
			result.RecordIndex, result.Record, found = i, record, true
		}
	}
	if !found {
		return taggedSIDBRoute{}, fmt.Errorf("SIDB value is not a registry record")
	}
	copy(result.Window[:], body[window:window+len(result.Window)])
	return result, nil
}
