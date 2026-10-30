package tagengine

import (
	"context"
	"errors"
	"testing"
)

func TestConstructorCipherCancellation(t *testing.T) {
	var state [256]byte
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := transformConstructorState(ctx, &state); !errors.Is(err, context.Canceled) {
		t.Fatalf("canceled transform: %v", err)
	}
	if err := transformConstructorState(nil, &state); err == nil {
		t.Fatal("nil context accepted")
	}
	if err := transformConstructorState(context.Background(), nil); err == nil {
		t.Fatal("nil state accepted")
	}
}
