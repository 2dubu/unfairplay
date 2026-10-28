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
