package recovery

import "testing"

func TestCurrentFairPlayByteMaskIsInvolutive(t *testing.T) {
	for value := 0; value < 256; value++ {
		input := byte(value)
		if got := applyCurrentFairPlayByteMask(applyCurrentFairPlayByteMask(input)); got != input {
			t.Fatalf("byte %#x round-tripped to %#x", input, got)
		}
	}
}

func TestCurrentFairPlayWordMaskUsesRecoveredLanes(t *testing.T) {
	const input = uint32(0x03020100)
	if got, want := applyCurrentFairPlayWordMask(input), uint32(0xd66c61de); got != want {
		t.Fatalf("word mask = %#x, want %#x", got, want)
	}
	if got := applyCurrentFairPlayWordMask(applyCurrentFairPlayWordMask(input)); got != input {
		t.Fatalf("word round trip = %#x, want %#x", got, input)
	}
}
