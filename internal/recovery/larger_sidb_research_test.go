package recovery

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/DrWhax/unfairplay/internal/fairplay"
)

// TestSIDBNativeSourceAndBodyOracle compares the production native-exit path
// with external, private native-reference outputs. It neither imports an existing
// enrollment model nor exports keys or decoded input.
func TestSIDBNativeSourceAndBodyOracle(t *testing.T) {
	capture := os.Getenv("UNFAIRPLAY_TEST_CAPTURE")
	sourceOracle := os.Getenv("UNFAIRPLAY_SIDB_SOURCE_ORACLE")
	bodyOracle := os.Getenv("UNFAIRPLAY_SIDB_BODY_ORACLE")
	if capture == "" || sourceOracle == "" || bodyOracle == "" {
		t.Skip("private capture and native oracle files not configured")
	}
	manifest, err := fairplay.LoadEnrollmentCapture(capture, true)
	if err != nil {
		t.Fatal(err)
	}
	canonical, err := fairplay.NormalizeEnrollmentCapture(capture, manifest)
	if err != nil {
		t.Fatal(err)
	}
	if len(canonical.SIDBs) != 1 {
		t.Fatal("oracle comparison requires one SIDB")
	}
	container, err := readPrivateInput(filepath.Join(capture, canonical.SIDBs[0].File.CapturePath))
	if err != nil {
		t.Fatal(err)
	}
	defer clear(container)
	hardware, err := readPrivateInput(canonical.PrivateHardwarePath())
	if err != nil {
		t.Fatal(err)
	}
	defer clear(hardware)
	expectedSource, err := readPrivateInput(sourceOracle)
	if err != nil {
		t.Fatal(err)
	}
	defer clear(expectedSource)
	expectedBody, err := readPrivateInput(bodyOracle)
	if err != nil {
		t.Fatal(err)
	}
	defer clear(expectedBody)

	source, layout, _, err := deriveNativeSelectionSource(container, hardware)
	if err != nil {
		t.Fatal(err)
	}
	defer clear(source)
	if len(expectedSource) < len(source) || !bytes.Equal(source, expectedSource[:len(source)]) {
		t.Fatal("complete selection source differs from native oracle")
	}
	projected, err := decodeSelectionSIDBContainerWithLayout(container, layout)
	if err != nil {
		t.Fatal(err)
	}
	defer clear(projected)
	candidates, err := enrollmentKeyCandidates(source[selectionSourcePrefixBytes:])
	if err != nil {
		t.Fatal(err)
	}
	defer clear(candidates)
	matches := 0
	for _, candidate := range candidates {
		body, err := decryptSelectionSIDBPayloadDynamic(projected, candidate.Key)
		if err != nil {
			t.Fatal(err)
		}
		if bytes.Equal(body, expectedBody) {
			matches++
		}
		clear(body)
	}
	if matches != 1 {
		t.Fatalf("native body matched %d candidate keys, want one", matches)
	}
	if _, err := parseDecodedSIDBRegistry(expectedBody); err != nil {
		t.Fatalf("native decoded body fails production grammar: %v", err)
	}
	t.Log("native source and body match")
}
