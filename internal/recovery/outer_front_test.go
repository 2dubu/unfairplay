package recovery

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"testing"
)

func TestOuterFrontDeterministic(t *testing.T) {
	var table [outerTableWords]uint32
	for index := range table {
		table[index] = uint32(index)*0x9e3779b9 ^ 0x85ebca6b
	}
	state := outerTableRoundState{
		Words:         [5]uint32{0x243f6a88, 0x85a308d3, 0x13198a2e, 0x03707344, 0xa4093822},
		Hardware:      [6]byte{0x12, 0x34, 0x56, 0x78, 0x9a, 0xbc},
		HardwareIndex: 4,
		Codec: [33]byte{
			4: 0x10, 5: 0x21, 6: 0x32, 7: 0x43, 8: 0x54, 9: 0x65,
			10: 0x76, 11: 0x87, 12: 0x98, 13: 0xa9, 14: 0xba, 15: 0xcb, 16: 0xdc,
		},
	}
	if !applyOuterFront(&table, 49, &state, 0x60319e4b, 0x12345678) {
		t.Fatal("semantic outer front rejected valid input")
	}
	var encoded [20]byte
	for index, value := range state.Words {
		binary.LittleEndian.PutUint32(encoded[index*4:], value)
	}
	if got, want := fmt.Sprintf("%x", sha256.Sum256(encoded[:])), "4a61d08eb6000a16b07da60fce983d91c3cce0b5e0828a477a16f8a770429233"; got != want {
		t.Fatalf("outer front digest = %s, want %s", got, want)
	}
}

func TestOuterFrontRejectsInvalidInput(t *testing.T) {
	var table [outerTableWords]uint32
	var state outerTableRoundState
	if applyOuterFront(nil, 0, &state, 0x60319d83, 0) || applyOuterFront(&table, 0, nil, 0x60319d83, 0) {
		t.Fatal("accepted a nil outer front input")
	}
	if applyOuterFront(&table, outerTableWords, &state, 0x60319d83, 0) {
		t.Fatal("accepted an invalid outer front base")
	}
	state.HardwareIndex = 6
	if applyOuterFront(&table, 0, &state, 0x60319d83, 0) {
		t.Fatal("accepted an invalid hardware index")
	}
	state.HardwareIndex = 0
	if applyOuterFront(&table, 0, &state, 0, 0) {
		t.Fatal("accepted a counter state outside the proved natural path")
	}
}
