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
		Hardware: [6]byte{0x12, 0x34, 0x56, 0x78, 0x9a, 0xbc},
		Codec:    [17]byte{4: 0x78, 6: 0x9a, 8: 0xbc},
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

func TestOuterSemanticIterationRejectsInvalidInput(t *testing.T) {
	var table [outerTableWords]uint32
	var state outerTableRoundState
	if _, ok := applyOuterSemanticIteration(nil, 0, &state); ok {
		t.Fatal("accepted a nil semantic table")
	}
	if _, ok := applyOuterSemanticIteration(&table, 0, nil); ok {
		t.Fatal("accepted nil semantic state")
	}
	if _, ok := applyOuterFullIteration(&table, 0, &state, nil); ok {
		t.Fatal("accepted nil full-iteration control state")
	}
}

func TestOuterSemanticIterationDeterministic(t *testing.T) {
	var table [outerTableWords]uint32
	for index := range table {
		table[index] = uint32(index)*0x85ebca6b ^ 0xc2b2ae35
	}
	state := outerTableRoundState{
		Words:         [5]uint32{0x243f6a88, 0x85a308d3, 0x13198a2e, 0x03707344, 0xa4093822},
		Hardware:      [6]byte{0x12, 0x34, 0x56, 0x78, 0x9a, 0xbc},
		HardwareIndex: 4,
		Codec: [17]byte{
			4: 0xde, 6: 0xad, 7: 0xbe, 8: 0xef, 9: 0x13,
			10: 0x57, 11: 0x9b, 12: 0xdf, 14: 0x24, 15: 0x68, 16: 0xac,
		},
	}
	next, ok := applyOuterSemanticIteration(&table, 49, &state)
	if !ok || next != 65 || state.HardwareIndex != 3 {
		t.Fatalf("semantic iteration advance = %d/%d/%v, want 65/3/true", next, state.HardwareIndex, ok)
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
	if got, want := fmt.Sprintf("%x", hasher.Sum(nil)), "cb32356c29d5e9c9514288305b49dd3d5b39ffb4cf426d5dd2efd64688355e13"; got != want {
		t.Fatalf("semantic iteration digest = %s, want %s", got, want)
	}
}

func TestOuterFullIterationDeterministic(t *testing.T) {
	var table [outerTableWords]uint32
	for index := range table {
		table[index] = uint32(index)*0x27d4eb2d ^ 0x165667b1
	}
	state := outerTableRoundState{
		Words:         [5]uint32{0x6a09e667, 0xbb67ae85, 0x3c6ef372, 0xa54ff53a, 0x510e527f},
		Hardware:      [6]byte{0x12, 0x34, 0x56, 0x78, 0x9a, 0xbc},
		HardwareIndex: 0,
		Codec: [17]byte{
			4: 0x10, 5: 0x21, 6: 0x32, 7: 0x43, 8: 0x54, 9: 0x65,
			10: 0x76, 11: 0x87, 12: 0x98, 13: 0xa9, 14: 0xba, 15: 0xcb, 16: 0xdc,
		},
	}
	control := outerIterationControl{X1: 0x60319d83, CounterSeed: 0x12345678}
	next, ok := applyOuterFullIteration(&table, outerTableInitialBase, &state, &control)
	if !ok || next != 11 || state.HardwareIndex != 5 {
		t.Fatalf("full iteration advance = %d/%d/%v, want 11/5/true", next, state.HardwareIndex, ok)
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
	for _, value := range [...]uint32{control.X1, control.CounterSeed} {
		binary.LittleEndian.PutUint32(word[:], value)
		_, _ = hasher.Write(word[:])
	}
	if got, want := fmt.Sprintf("%x", hasher.Sum(nil)), "7d75c9b0f553425e8a979606b7eaa48af3415c7ca396aa89f23571890417a8e0"; got != want {
		t.Fatalf("full iteration digest = %s, want %s", got, want)
	}
}
