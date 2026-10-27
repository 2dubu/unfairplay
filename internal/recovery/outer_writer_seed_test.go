package recovery

import "testing"

func TestOuterWriterSeedStep(t *testing.T) {
	vectors := []struct {
		seed, want uint32
	}{
		{0x00000000, 0x0a782f7f},
		{0x00000001, 0x0a782f80},
		{0x12345678, 0x0a782fe5},
		{0x9abcdef0, 0x0a78301b},
		{0xffffffff, 0x0a782f7e},
		{0x0263747b, 0x0a782f7c},
	}
	for _, vector := range vectors {
		if got := outerWriterSeedStep(vector.seed); got != vector.want {
			t.Fatalf("writer seed step(%#x) = %#x, want %#x", vector.seed, got, vector.want)
		}
	}
}
