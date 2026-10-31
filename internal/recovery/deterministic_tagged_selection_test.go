package recovery

import (
	"context"
	"encoding/binary"
	"errors"
	"testing"
)

func TestAuthenticateTaggedSelection(t *testing.T) {
	decode, _, body := syntheticRoute()
	user := uint32(0x12345678)
	user = user*0xa46fd0ff + 0xfe063531
	const keyID = uint32(7)
	hardware := []byte{0, 1, 2, 3, 4, 5}
	registry, err := parseDecodedSIDBRegistry(body, 2)
	if err != nil {
		t.Fatal(err)
	}
	key, err := accountClassKey(hardware, user, keyID)
	if err != nil {
		t.Fatal(err)
	}
	recordIndex := -1
	for i, record := range registry.Records {
		if record.SourceOffset == 174 {
			recordIndex = i
		}
	}
	if recordIndex < 0 {
		t.Fatal("missing value record")
	}
	precursor, err := deriveCurrentSINFPrecursorFromRecordWindow(registry.Records[recordIndex], 0, key)
	if err != nil {
		t.Fatal(err)
	}
	sinf := syntheticCurrentSINFForPrecursor(t, precursor)
	start, end, err := enrollmentAtom(sinf, 0, len(sinf), "sinf")
	if err != nil {
		t.Fatal(err)
	}
	start, end, err = enrollmentAtom(sinf, start, end, "schi")
	if err != nil {
		t.Fatal(err)
	}
	p, _, err := enrollmentAtom(sinf, start, end, "user")
	if err != nil {
		t.Fatal(err)
	}
	userOffset := p
	binary.BigEndian.PutUint32(sinf[p:], user)
	p, _, err = enrollmentAtom(sinf, start, end, "key ")
	if err != nil {
		t.Fatal(err)
	}
	keyOffset := p
	binary.BigEndian.PutUint32(sinf[p:], keyID)
	regs := []decodedSIDBRegistryCandidate{{Registry: registry, Data: body}}
	var events []Progress
	got, err := authenticateTaggedSelection(context.Background(), decode, sinf, regs, hardware, 2, user, keyID, func(p Progress) { events = append(events, p) })
	if err != nil || got.Precursor != precursor || got.RecordIndex != recordIndex {
		t.Fatalf("selected: %+v, %v", got, err)
	}
	if got.routeEvidence != (selectedRouteEvidence{
		RootVersion: 0x10001, AccountVersion: 0x10002,
		TypeCursor: 190, TypeLength: 24, SignCursor: 214, SignLength: 148,
		ValueCursor: 154, ValueLength: 36, SelectedRecordIndex: recordIndex,
		SelectedPayloadOffset: 174, SelectedWindowOffset: 0,
	}) {
		t.Fatalf("route evidence: %+v", got.routeEvidence)
	}
	if len(events) < 2 || events[len(events)-1].Phase != "search-sinf-complete" || events[len(events)-1].Total != 1 {
		t.Fatal("invalid progress")
	}
	// The root version leaf is also a 36-byte record. Construct a SINF that
	// authenticates it; that must not authorize a different selected route.
	metadataIndex := -1
	for i, record := range registry.Records {
		if record.SourceOffset == 42 {
			metadataIndex = i
		}
	}
	if metadataIndex < 0 || metadataIndex == recordIndex {
		t.Fatal("missing distinct metadata record")
	}
	metadataPrecursor, err := deriveCurrentSINFPrecursorFromRecordWindow(registry.Records[metadataIndex], 0, key)
	if err != nil || metadataPrecursor == precursor {
		t.Fatal("metadata precursor must differ from selected value", err)
	}
	trap := syntheticCurrentSINFForPrecursor(t, metadataPrecursor)
	binary.BigEndian.PutUint32(trap[userOffset:], user)
	binary.BigEndian.PutUint32(trap[keyOffset:], keyID)
	if accepted, err := validateCurrentSINFPrecursor(trap, metadataPrecursor); err != nil || !accepted {
		t.Fatalf("metadata precursor did not authenticate its fixture: %v", err)
	}

	if _, err := authenticateTaggedSelection(context.Background(), decode, trap, regs, hardware, 2, user, keyID, nil); !errors.Is(err, errNoAuthenticatedClassPrecursor) {
		t.Fatalf("selected value unexpectedly fell back to metadata: %v", err)
	}
	if _, err := authenticateTaggedSelection(context.Background(), decode, sinf, regs, hardware, 2, user+1, keyID, nil); err == nil {
		t.Fatal("accepted mismatched selector")
	}
	if _, err := authenticateTaggedSelection(context.Background(), decode, sinf, append(regs, regs...), hardware, 2, user, keyID, nil); err == nil {
		t.Fatal("accepted multiple registries")
	}
	stale := registry
	stale.Records = append([]decodedSIDBRecord(nil), registry.Records...)
	stale.Records[recordIndex].Payload[0] ^= 1
	if _, err := authenticateTaggedSelection(context.Background(), decode, sinf, []decodedSIDBRegistryCandidate{{Registry: stale, Data: body}}, hardware, 2, user, keyID, nil); err == nil {
		t.Fatal("accepted stale selected record")
	}
	wrong := append([]byte(nil), sinf...)
	priv, err := currentSINFPrivateSource(wrong)
	if err != nil {
		t.Fatal(err)
	}
	priv[32] ^= 0x80
	if _, err := authenticateTaggedSelection(context.Background(), decode, wrong, regs, hardware, 2, user, keyID, nil); !errors.Is(err, errNoAuthenticatedClassPrecursor) {
		t.Fatalf("corrupt SINF: %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := authenticateTaggedSelection(ctx, decode, sinf, regs, hardware, 2, user, keyID, nil); !errors.Is(err, context.Canceled) {
		t.Fatalf("cancel: %v", err)
	}
	ctx, cancel = context.WithCancel(context.Background())
	defer cancel()
	if _, err := authenticateTaggedSelection(ctx, decode, sinf, regs, hardware, 2, user, keyID, func(Progress) { cancel() }); !errors.Is(err, context.Canceled) {
		t.Fatalf("progress cancellation: %v", err)
	}
}
