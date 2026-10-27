package recovery

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"testing"
)

func TestOuterTableRoundDeterministic(t *testing.T) {
	var table [outerTableWords]uint32
	for index := range table {
		table[index] = uint32(index)*0x9e3779b9 ^ 0x7f4a7c15
	}
	state := outerTableRoundState{
		Words:    [5]uint32{0x01234567, 0x89abcdef, 0x0f1e2d3c, 0x4b5a6978, 0xdeadbeef},
		Hardware: [6]byte{0x12, 0x34, 0x56, 0x78, 0x9a, 0xbc}, Codec: [3]uint32{0x78, 0x9a, 0xbc},
	}
	next, ok := applyOuterTableRound(&table, outerTableInitialBase, &state)
	if !ok || next != 11 {
		t.Fatalf("round advance = %d/%v, want 11/true", next, ok)
	}
	if state.HardwareIndex != 5 {
		t.Fatalf("hardware index = %d, want 5", state.HardwareIndex)
	}
	hasher := sha256.New()
	var word [4]byte
	for _, value := range table {
		binary.LittleEndian.PutUint32(word[:], value)
		_, _ = hasher.Write(word[:])
	}
	for _, value := range state.Words {
		binary.LittleEndian.PutUint32(word[:], value)
		_, _ = hasher.Write(word[:])
	}
	if got, want := fmt.Sprintf("%x", hasher.Sum(nil)), "692e78d0430d5ac551b71f6352aaa29cdb10420b4e08677ec867d6fc7e62bcf3"; got != want {
		t.Fatalf("semantic round digest = %s, want %s", got, want)
	}
}

func TestOuterTableRoundRejectsInvalidInput(t *testing.T) {
	var table [outerTableWords]uint32
	var state outerTableRoundState
	if _, ok := applyOuterTableRound(nil, 0, &state); ok {
		t.Fatal("accepted a nil table")
	}
	if _, ok := applyOuterTableRound(&table, 0, nil); ok {
		t.Fatal("accepted nil state")
	}
	if _, ok := applyOuterTableRound(&table, outerTableWords, &state); ok {
		t.Fatal("accepted an invalid base")
	}
	state.HardwareIndex = uint8(len(state.Hardware))
	if _, ok := applyOuterTableRound(&table, 0, &state); ok {
		t.Fatal("accepted an invalid hardware index")
	}
}
