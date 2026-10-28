package recovery

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"testing"
)

func TestSelectionSourceLayouts(t *testing.T) {
	envelope := make([]byte, selectionSourceEnvelopeBytes)
	for index := range envelope {
		envelope[index] = byte(index*17 + 5)
	}
	layouts, err := selectionSourceLayouts(envelope)
	if err != nil {
		t.Fatal(err)
	}
	if len(layouts) != 2 || len(layouts[0]) != selectionSourceBodyBytes || len(layouts[1]) != selectionSourceEnvelopeBytes ||
		!bytes.Equal(layouts[0], envelope[selectionSourcePrefixBytes:]) || !bytes.Equal(layouts[1], envelope) {
		t.Fatal("224-byte selection envelope was not split structurally")
	}
	body := envelope[:selectionSourceBodyBytes]
	layouts, err = selectionSourceLayouts(body)
	if err != nil || len(layouts) != 1 || !bytes.Equal(layouts[0], body) {
		t.Fatal("192-byte selection body was not retained")
	}
	if _, err := selectionSourceLayouts(envelope[:208]); err == nil {
		t.Fatal("unsupported selection source geometry was accepted")
	}
}

func TestEnrollmentFrameCandidates(t *testing.T) {
	for _, size := range []int{192, 224} {
		source := make([]byte, size)
		for index := range source {
			source[index] = byte(index*17 + size/16)
		}
		candidates, err := enrollmentFrameCandidates(source)
		if err != nil {
			t.Fatal(err)
		}
		if len(candidates) != size/16 {
			t.Fatalf("%d-byte source produced %d candidates", size, len(candidates))
		}
		hash := sha256.New()
		for slot, candidate := range candidates {
			if candidate.Slot != slot {
				t.Fatalf("candidate %d has slot %d", slot, candidate.Slot)
			}
			_, _ = hash.Write(candidate.State[:])
			_, _ = hash.Write(candidate.Frame[:])
			_, _ = hash.Write(candidate.Key[:])
		}
		want := map[int]string{
			192: "b7fa7cd0bfea718156a026078df095b7913966b2eaa1c7362fb8d4b3738a4b0e",
			224: "0fa30ba4791845e36c5603ba1ab4d18594dd4750ac5cc8d3ae9b7e55d667135a",
		}[size]
		if got := fmt.Sprintf("%x", hash.Sum(nil)); got != want {
			t.Fatalf("%d-byte candidate digest=%s, want %s", size, got, want)
		}
	}
}
