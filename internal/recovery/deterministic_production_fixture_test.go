package recovery

import (
	"context"
	"encoding/binary"
	"errors"
	"testing"

	"github.com/DrWhax/unfairplay/internal/recovery/tagengine"
)

// syntheticProductionRoute remasks a public synthetic hierarchy at the native
// root cursor and encodes its tags using this synthetic hardware's codebook.
func syntheticProductionRoute(t *testing.T, hardware []byte, bodySize int) ([]byte, decodedSIDBRegistry, []byte, uint32, uint32, [16]byte) {
	t.Helper()
	decode, _, old := syntheticRoute()
	nodes, err := parseDecodedSIDBTree(old, 2)
	if err != nil {
		t.Fatal(err)
	}
	root, err := nativeSIDBRootCursor(hardware)
	if err != nil {
		t.Fatal(err)
	}
	if bodySize == 0 {
		bodySize = root + 360
	}
	if bodySize < root+360 || (bodySize > root+360 && bodySize < root+380) {
		t.Fatal("invalid fixture body size")
	}
	body := syntheticSIDBTreeBuffer(bodySize)
	for i := 0; i < 360; i++ {
		src, dst := 2+i, root+i
		body[dst] = old[src] ^ currentClassRecordMask[(src+12)%16] ^ currentClassRecordMask[(dst+12)%16]
	}
	book, err := tagengine.GenerateCodebookWithContext(context.Background(), hardware)
	if err != nil {
		t.Fatal(err)
	}
	defer clear(book)
	byName := make(map[string]uint32, len(book))
	for _, entry := range book {
		byName[entry.Semantic] = entry.Raw
	}
	for _, node := range nodes {
		name, ok := decode(sidbTreeRawWord(old, node.Offset+4))
		if !ok {
			t.Fatal("unknown synthetic tag")
		}
		code, ok := byName[name]
		if !ok {
			t.Fatalf("missing generated tag %q", name)
		}
		putSyntheticSIDBRawWord(body, root+node.Offset-2+4, code)
	}
	if extra := bodySize - root - 360; extra > 0 {
		putSyntheticSIDBNode(body, root-2, bodySize-root, 3)
		putSyntheticSIDBNode(body, root+360-2, extra, 0)
		putSyntheticSIDBRawWord(body, root+360+4, byName["sign"])
	}
	registry, err := parseDecodedSIDBRegistry(body, root)
	if err != nil {
		t.Fatal(err)
	}
	index := -1
	for i, r := range registry.Records {
		if r.SourceOffset == root+172 {
			index = i
		}
	}
	if index < 0 {
		t.Fatal("missing selected value")
	}
	user := uint32(0x12345678)
	user = user*0xa46fd0ff + 0xfe063531
	const keyID uint32 = 7
	key, err := accountClassKey(hardware, user, keyID)
	if err != nil {
		t.Fatal(err)
	}
	precursor, err := deriveCurrentSINFPrecursorFromRecordWindow(registry.Records[index], 0, key)
	if err != nil {
		t.Fatal(err)
	}
	sinf := syntheticCurrentSINFForPrecursorRights(t, precursor, 104)
	start, end, err := enrollmentAtom(sinf, 0, len(sinf), "sinf")
	if err != nil {
		t.Fatal(err)
	}
	start, end, err = enrollmentAtom(sinf, start, end, "schi")
	if err != nil {
		t.Fatal(err)
	}
	u, _, err := enrollmentAtom(sinf, start, end, "user")
	if err != nil {
		t.Fatal(err)
	}
	binary.BigEndian.PutUint32(sinf[u:], user)
	k, _, err := enrollmentAtom(sinf, start, end, "key ")
	if err != nil {
		t.Fatal(err)
	}
	binary.BigEndian.PutUint32(sinf[k:], keyID)
	return body, registry, sinf, user, keyID, precursor
}

func TestProductionSelectionDoesNotFallbackToMetadata(t *testing.T) {
	hardware := []byte{0, 1, 2, 3, 4, 5}
	body, registry, sinf, user, keyID, selectedPrecursor := syntheticProductionRoute(t, hardware, 0)
	root, err := nativeSIDBRootCursor(hardware)
	if err != nil {
		t.Fatal(err)
	}
	index := -1
	for i, record := range registry.Records {
		if record.SourceOffset == root+40 {
			index = i
		}
	}
	if index < 0 {
		t.Fatal("missing root version record")
	}
	key, err := accountClassKey(hardware, user, keyID)
	if err != nil {
		t.Fatal(err)
	}
	metadataPrecursor, err := deriveCurrentSINFPrecursorFromRecordWindow(registry.Records[index], 0, key)
	if err != nil || metadataPrecursor == selectedPrecursor {
		t.Fatal("metadata precursor is not distinct", err)
	}
	trap := syntheticCurrentSINFForPrecursorRights(t, metadataPrecursor, 104)
	start, end, err := enrollmentAtom(trap, 0, len(trap), "sinf")
	if err != nil {
		t.Fatal(err)
	}
	start, end, err = enrollmentAtom(trap, start, end, "schi")
	if err != nil {
		t.Fatal(err)
	}
	u, _, err := enrollmentAtom(trap, start, end, "user")
	if err != nil {
		t.Fatal(err)
	}
	binary.BigEndian.PutUint32(trap[u:], user)
	k, _, err := enrollmentAtom(trap, start, end, "key ")
	if err != nil {
		t.Fatal(err)
	}
	binary.BigEndian.PutUint32(trap[k:], keyID)
	registries := []decodedSIDBRegistryCandidate{{Registry: registry, Data: body}}
	if accepted, err := validateCurrentSINFPrecursor(trap, metadataPrecursor); err != nil || !accepted {
		t.Fatalf("metadata precursor did not authenticate its fixture: %v", err)
	}

	if _, err := selectAuthenticatedEnrollmentClassPrecursor(sinf, registries, hardware, user, keyID); err != nil {
		t.Fatalf("selected fixture: %v", err)
	}
	if _, err := selectAuthenticatedEnrollmentClassPrecursor(trap, registries, hardware, user, keyID); !errors.Is(err, errNoAuthenticatedClassPrecursor) {
		t.Fatalf("fallback to metadata: %v", err)
	}
}
