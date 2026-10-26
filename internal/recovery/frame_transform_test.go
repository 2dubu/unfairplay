package recovery

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"reflect"
	"testing"
)

func frameTestInput(index uint32) (result [16]byte) {
	seed := append([]byte(nil), []byte("macOS FairPlay Go oracle input")...)
	var encoded [4]byte
	binary.BigEndian.PutUint32(encoded[:], index)
	digest := sha256.Sum256(append(seed, encoded[:]...))
	copy(result[:], digest[:16])
	return result
}

func TestRecoveredFrameTransform(t *testing.T) {
	value, err := loadModelBytes(embeddedModelBytes())
	if err != nil {
		t.Fatal(err)
	}
	output, err := frameTransform(value, [16]byte{})
	if err != nil {
		t.Fatal(err)
	}
	want, _ := hex.DecodeString("2b7cdd00ede49ae3e3040f53c0046e65")
	if !reflect.DeepEqual(output[:], want) {
		t.Fatalf("zero frame output = %x, want %x", output, want)
	}
	if core := frameCore([16]byte{}); core != [16]byte{0x90, 0xf8, 0x3b, 0x63, 0xc1, 0xda, 0xf2, 0x88, 0x92, 0x72, 0xa7, 0xde, 0x8b, 0x4a, 0x93, 0x21} {
		t.Fatalf("zero frame core = %x", core)
	}
	for index, expected := range []string{
		"90cb4346ed4a09fc2cfb3f83680dd7f4",
		"d768891972be70c0b6cd6fb6c5219563",
		"92bf95b1688aadf7084c344aada2de32",
	} {
		output, err := frameTransform(value, frameTestInput(uint32(index)))
		if err != nil {
			t.Fatal(err)
		}
		if got := hex.EncodeToString(output[:]); got != expected {
			t.Fatalf("deterministic frame %d = %s, want %s", index, got, expected)
		}
	}
}
