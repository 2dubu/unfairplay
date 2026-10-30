package tagengine

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestConstructorFailsClosed(t *testing.T) {
	for _, size := range []int{0, 5, 7} {
		if _, err := GenerateCodebook(make([]byte, size)); err == nil {
			t.Fatalf("accepted hardware length %d", size)
		}
	}
}

func TestPublicSyntheticCodebookShape(t *testing.T) {
	entries, err := GenerateCodebook([]byte{0, 1, 2, 3, 4, 5})
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 44 {
		t.Fatalf("entry count=%d", len(entries))
	}
	seen := map[uint32]bool{}
	for _, entry := range entries {
		if len(entry.Semantic) != 4 || seen[entry.Raw] {
			t.Fatal("invalid or ambiguous codebook")
		}
		seen[entry.Raw] = true
	}
}

type countingContext struct{ calls, cancelAt int }

func (c *countingContext) Deadline() (time.Time, bool) { return time.Time{}, false }
func (c *countingContext) Done() <-chan struct{}       { return nil }
func (c *countingContext) Value(any) any               { return nil }
func (c *countingContext) Err() error {
	c.calls++
	if c.calls >= c.cancelAt {
		return context.Canceled
	}
	return nil
}
func TestConstructorCancellationAndInputOwnership(t *testing.T) {
	hw := []byte{0, 1, 2, 3, 4, 5}
	original := append([]byte(nil), hw...)
	for _, n := range []int{1, 3, 40} {
		ctx := &countingContext{cancelAt: n}
		out, err := GenerateCodebookWithContext(ctx, hw)
		if !errors.Is(err, context.Canceled) || out != nil {
			t.Fatalf("cancellation %d: %v", n, err)
		}
	}
	for i := range hw {
		if hw[i] != original[i] {
			t.Fatal("caller hardware modified")
		}
	}
	if _, err := GenerateCodebookWithContext(nil, hw); err == nil {
		t.Fatal("accepted nil context")
	}
}
