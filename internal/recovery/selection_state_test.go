package recovery

import (
	"crypto/sha256"
	"fmt"
	"testing"
)

func TestTransformSelectionStateSlots(t *testing.T) {
	encodedZeroSource := make([]byte, 224)
	inputMask := [...]byte{0xa8, 0xd6, 0x01, 0x92, 0xfe, 0x37, 0x94, 0x0f, 0xe3, 0x33, 0x04, 0xfd, 0x34, 0xe0, 0xad, 0x43}
	for index := range encodedZeroSource {
		encodedZeroSource[index] = inputMask[index%len(inputMask)]
	}
	zero, _, err := transformSelectionStateSlots(encodedZeroSource)
	if err != nil {
		t.Fatal(err)
	}
	zeroHash := sha256.New()
	for _, slot := range zero {
		_, _ = zeroHash.Write(slot[:])
	}
	if got, want := fmt.Sprintf("%x", zeroHash.Sum(nil)), "9e9f0073f8732cd41f7afa3106b5d25c068864e4974814944186664865b3aa93"; got != want {
		t.Fatalf("encoded-zero selection state digest=%s, want native %s", got, want)
	}

	source := make([]byte, 224)
	for index := range source {
		source[index] = byte(index*41 + index/16*13 + 5)
	}
	output, realSlots, err := transformSelectionStateSlots(source)
	if err != nil {
		t.Fatal(err)
	}
	if realSlots != selectionStateSlotCount {
		t.Fatalf("real slots=%d", realSlots)
	}
	hash := sha256.New()
	for _, slot := range output {
		_, _ = hash.Write(slot[:])
	}
	if got, want := fmt.Sprintf("%x", hash.Sum(nil)), "0b923115d70c4b26b594b1f5833f3b916ff9abb1f7fbc021f5190e0793573f83"; got != want {
		t.Fatalf("selection state digest=%s, want %s", got, want)
	}

	short, realSlots, err := transformSelectionStateSlots(source[:192])
	if err != nil {
		t.Fatal(err)
	}
	if realSlots != 12 {
		t.Fatalf("short real slots=%d", realSlots)
	}
	for slot := range 12 {
		if short[slot] != output[slot] {
			t.Fatalf("short source changed real slot %d", slot)
		}
	}
	if short[12] == [16]byte{} || short[13] == [16]byte{} {
		t.Fatal("short source did not preserve native padded-slot transforms")
	}
	if _, _, err := transformSelectionStateSlots(source[:208]); err == nil {
		t.Fatal("accepted unsupported state source length")
	}
}

func TestTransformSelectionStateSlotRejectsIndex(t *testing.T) {
	if _, err := transformSelectionStateSlot(-1, [16]byte{}); err == nil {
		t.Fatal("accepted negative state slot")
	}
	if _, err := transformSelectionStateSlot(selectionStateSlotCount, [16]byte{}); err == nil {
		t.Fatal("accepted out-of-range state slot")
	}
}
