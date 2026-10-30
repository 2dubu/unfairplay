package recovery

import (
	"bytes"
	"testing"
)

func TestSelectionSourceLayouts(t *testing.T) {
	for count := 1; count <= maximumMappedSelectionModulus; count++ {
		source := make([]byte, selectionSourcePrefixBytes+16*count)
		for index := range source {
			source[index] = byte(index*17 + 5)
		}
		layouts, err := selectionSourceLayouts(source)
		if err != nil || len(layouts) != 1 || len(layouts[0]) != 16*count || !bytes.Equal(layouts[0], source[selectionSourcePrefixBytes:]) {
			t.Fatalf("invalid state source for %d slots", count)
		}
	}
	for _, size := range []int{0, 32, 47, 255, 321, 336} {
		if _, err := selectionSourceLayouts(make([]byte, size)); err == nil {
			t.Fatalf("accepted unsupported source length %d", size)
		}
	}
}
