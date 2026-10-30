package recovery

import (
	"context"
	"errors"
	"testing"
)

func TestAuthenticateSelectedClassWindow(t *testing.T) {
	data := syntheticSIDBRecordTree(2, 31, 7)
	registry, err := parseDecodedSIDBRegistry(data, 2)
	if err != nil {
		t.Fatal(err)
	}
	key := currentClassEnrollmentParameters().classKey()
	const window = 0
	precursor, err := deriveCurrentSINFPrecursorFromRecordWindow(registry.Records[1], window, key)
	if err != nil {
		t.Fatal(err)
	}
	sinf := syntheticCurrentSINFForPrecursor(t, precursor)
	registries := []decodedSIDBRegistryCandidate{{}, {Registry: registry, Data: data}}
	var events []Progress
	selected, err := authenticateSelectedClassWindow(context.Background(), sinf, registries, 1, 1, window, key, func(p Progress) { events = append(events, p) })
	if err != nil {
		t.Fatal(err)
	}
	if selected.Precursor != precursor || selected.RegistryIndex != 1 || selected.RecordIndex != 1 || selected.WindowOffset != window {
		t.Fatalf("wrong selected tuple: %+v", selected)
	}
	if len(events) < 2 {
		t.Fatalf("missing progress events: %+v", events)
	}
	var previous uint64
	for _, p := range events {
		if p.Total != 1 || p.Completed < previous || p.Completed > 1 || p.RegistryCandidates != 1 || p.RecordCandidates != 1 || p.WindowCandidates != 1 {
			t.Fatalf("invalid progress: %+v", p)
		}
		previous = p.Completed
	}
	last := events[len(events)-1]
	if last.Phase != "search-sinf-complete" || last.Completed != 1 {
		t.Fatalf("missing completion: %+v", last)
	}

	// A correct tuple exists in the same registry, but a selected wrong tuple
	// must fail instead of searching for it.
	events = nil
	_, err = authenticateSelectedClassWindow(context.Background(), sinf, registries, 1, 0, window, key, func(p Progress) { events = append(events, p) })
	if !errors.Is(err, errNoAuthenticatedClassPrecursor) {
		t.Fatalf("wrong tuple: %v", err)
	}
	for _, p := range events {
		if p.Phase == "search-sinf-complete" {
			t.Fatal("failure reported success")
		}
	}

	wrongKey := key
	wrongKey[0] ^= 0x80
	if _, err := authenticateSelectedClassWindow(context.Background(), sinf, registries, 1, 1, window, wrongKey, nil); !errors.Is(err, errNoAuthenticatedClassPrecursor) {
		t.Fatalf("wrong key: %v", err)
	}
	if _, err := authenticateSelectedClassWindow(context.Background(), sinf, registries, 1, 1, window+1, key, nil); err == nil {
		t.Fatalf("wrong window: %v", err)
	}
	// Keep the first encrypted block (the fast probe) intact, but corrupt a
	// later block so full SINF structure authentication must reject it.
	corrupt := append([]byte(nil), sinf...)
	private, err := currentSINFPrivateSource(corrupt)
	if err != nil {
		t.Fatal(err)
	}
	private[32] ^= 0x80
	if _, err := authenticateSelectedClassWindow(context.Background(), corrupt, registries, 1, 1, window, key, nil); !errors.Is(err, errNoAuthenticatedClassPrecursor) {
		t.Fatalf("corrupt SINF: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := authenticateSelectedClassWindow(ctx, sinf, registries, 1, 1, window, key, nil); !errors.Is(err, context.Canceled) {
		t.Fatalf("pre-canceled: %v", err)
	}
	ctx, cancel = context.WithCancel(context.Background())
	defer cancel()
	events = nil
	_, err = authenticateSelectedClassWindow(ctx, sinf, registries, 1, 1, window, key, func(p Progress) { events = append(events, p); cancel() })
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("callback cancellation: %v", err)
	}
	for _, p := range events {
		if p.Phase == "search-sinf-complete" {
			t.Fatal("cancellation reported success")
		}
	}
	ctx, cancel = context.WithCancel(context.Background())
	defer cancel()
	_, err = authenticateSelectedClassWindow(ctx, sinf, registries, 1, 1, window, key, func(p Progress) {
		if p.Phase == "search-sinf" && p.Completed == 1 {
			cancel()
		}
	})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("completed-work cancellation: %v", err)
	}
}

func TestAuthenticateSelectedClassWindowRejectsInvalidSelection(t *testing.T) {
	registries := []decodedSIDBRegistryCandidate{{Registry: decodedSIDBRegistry{Records: []decodedSIDBRecord{{}}}}}
	for _, indices := range [][3]int{{-1, 0, 0}, {1, 0, 0}, {0, -1, 0}, {0, 1, 0}, {0, 0, -1}, {0, 0, 1}} {
		if _, err := authenticateSelectedClassWindow(context.Background(), nil, registries, indices[0], indices[1], indices[2], [16]byte{}, nil); err == nil {
			t.Fatalf("accepted invalid selection %v", indices)
		}
	}
	for _, phase := range []int{-1, len(currentClassRecordMask)} {
		registries[0].Registry.Records[0].MaskOffset = phase
		if _, err := authenticateSelectedClassWindow(context.Background(), nil, registries, 0, 0, 0, [16]byte{}, nil); err == nil {
			t.Fatalf("accepted invalid mask offset %d", phase)
		}
	}
	registries[0].Registry.Records[0].MaskOffset = 0
	for _, tail := range []int{-1, decodedSIDBTailBytes + 1} {
		registries[0].Registry.Records[0].TailAvailable = tail
		if _, err := authenticateSelectedClassWindow(context.Background(), nil, registries, 0, 0, 0, [16]byte{}, nil); err == nil {
			t.Fatalf("accepted invalid tail length %d", tail)
		}
	}
}
