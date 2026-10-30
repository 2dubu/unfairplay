package recovery

import "fmt"

// sidbSchemaVersion is the root schema gate's semantic version. Recognition
// here does not establish acceptance of the complete enrollment or recovery.
type sidbSchemaVersion uint32

const (
	sidbSchemaVersion10001 sidbSchemaVersion = 0x10001
	sidbSchemaVersion10002 sidbSchemaVersion = 0x10002
)

// readSIDBSchemaVersion reads a version leaf using native length-field cursor
// coordinates. The caller must establish that cursor names the intended vers
// leaf and length is its decoded declared length. Four encoded words start at
// cursor+20. This reader independently bounds [cursor, cursor+length) to the
// body as a safety restriction; native extent enforcement is not established.
func readSIDBSchemaVersion(data []byte, cursor, length int) (sidbSchemaVersion, error) {
	const payloadOffset, payloadBytes = 20, 16
	if cursor < 0 || cursor > len(data) || length < payloadOffset+payloadBytes || length > len(data)-cursor {
		return 0, fmt.Errorf("SIDB schema version leaf has invalid extent")
	}
	var version uint32
	for i := range 4 {
		word := sidbTreeWord(data, cursor+payloadOffset+4*i, 4, 5)
		version = version<<8 | word&0xff
	}
	switch sidbSchemaVersion(version) {
	case sidbSchemaVersion10001, sidbSchemaVersion10002:
		return sidbSchemaVersion(version), nil
	default:
		return 0, fmt.Errorf("unsupported SIDB schema version %#x", version)
	}
}
