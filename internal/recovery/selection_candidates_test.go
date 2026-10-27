package recovery

import (
	"crypto/sha256"
	"fmt"
	"testing"
)

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
		}
		want := map[int]string{
			192: "31203fc8c2853b0a675952004acc455e11e23c410df8bb76ca895ce191100340",
			224: "194ef8ef3f9ea21dfab7471ef441087f7ac9d82fa3cdd6a9f53615790eca6b83",
		}[size]
		if got := fmt.Sprintf("%x", hash.Sum(nil)); got != want {
			t.Fatalf("%d-byte candidate digest=%s, want %s", size, got, want)
		}
	}
}
