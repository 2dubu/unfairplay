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
		Words:     [5]uint32{0x01234567, 0x89abcdef, 0x0f1e2d3c, 0x4b5a6978, 0xdeadbeef},
		PendingS0: 0x10203040, PendingS1: 0x50607080, PendingS2: 0x90a0b0c0,
		PendingS3: 0xd0e0f000, PendingS4: 0x13579bdf,
		Mix16: 0x2468ace0, Mix17: 0xfedcba98, Lane19Mix: 0x55aa55aa, LateMix: 0xaa55aa55,
		Source: [3]uint32{0x12, 0x34, 0x56}, Codec: [3]uint32{0x78, 0x9a, 0xbc},
	}
	next, ok := applyOuterTableRound(&table, outerTableInitialBase, &state)
	if !ok || next != 11 {
		t.Fatalf("round advance = %d/%v, want 11/true", next, ok)
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
	if got, want := fmt.Sprintf("%x", hasher.Sum(nil)), "90ee6338eae739e24afd765adb8f45e7978156e00c1c2c83d74df32abf9a70a7"; got != want {
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
}
