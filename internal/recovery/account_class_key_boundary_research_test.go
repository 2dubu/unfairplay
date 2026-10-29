package recovery

import (
	"encoding/binary"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

// The caller retains all traces outside the repository; no digest or key is
// emitted by this test. Both SHA-1 suffix words are perturbed together.
func TestAccountClassKeyNativeBoundaries(t *testing.T) {
	dir, digestPath := os.Getenv("UNFAIRPLAY_CLASS_BOUNDARY_DIRECTORY"), os.Getenv("UNFAIRPLAY_CLASS_DIGEST_ORACLE")
	if dir == "" || digestPath == "" {
		t.Skip("native key boundary controls not configured")
	}
	original, err := os.ReadFile(digestPath)
	if err != nil {
		t.Fatal(err)
	}
	defer clear(original)
	if len(original) != 20 {
		t.Fatal("invalid digest size")
	}
	controls := []struct {
		name string
		word uint16
	}{{"0", 0}, {"1", 1}, {"0x3fff", 0x3fff}, {"0x4000", 0x4000}, {"0x7ffe", 0x7ffe}, {"0x7fff", 0x7fff}, {"0x8000", 0x8000}, {"0xffff", 0xffff}}
	for _, tc := range controls {
		t.Run(tc.name, func(t *testing.T) {
			blob, err := os.ReadFile(filepath.Join(dir, "sidb-join-boundary-"+tc.name+".json"))
			if err != nil {
				t.Fatal(err)
			}
			defer clear(blob)
			var rows []struct {
				PC  uint64
				Key []byte
			}
			if err := json.Unmarshal(blob, &rows); err != nil {
				t.Fatal(err)
			}
			var digest [20]byte
			copy(digest[:], original)
			binary.BigEndian.PutUint16(digest[16:18], tc.word)
			binary.BigEndian.PutUint16(digest[18:20], tc.word)
			got, err := accountClassKeyFromDigest(digest)
			clear(digest[:])
			if err != nil {
				t.Fatal(err)
			}
			defer clear(got[:])
			matches := 0
			for _, row := range rows {
				if row.PC != 0x1008d7634 {
					continue
				}
				if len(row.Key) != 16 {
					t.Fatal("invalid native join key length")
				}
				same := true
				for i := range got {
					if got[i] != row.Key[i]^0x1b {
						same = false
					}
				}
				if same {
					matches++
				}
				clear(row.Key)
			}
			if matches != 1 {
				t.Fatalf("native join matched %d keys, want one", matches)
			}
		})
	}
}
