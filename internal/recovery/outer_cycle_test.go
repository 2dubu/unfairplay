package recovery

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"testing"
)

func TestOuterCycleDeterministic(t *testing.T) {
	state := outerTableRoundState{
		Words:         [5]uint32{0x243f6a88, 0x85a308d3, 0x13198a2e, 0x03707344, 0xa4093822},
		Hardware:      [6]byte{0x12, 0x34, 0x56, 0x78, 0x9a, 0xbc},
		HardwareIndex: 3,
		Codec: [33]byte{
			4: 0x10, 5: 0x21, 6: 0x32, 7: 0x43, 8: 0x54, 9: 0x65,
			10: 0x76, 11: 0x87, 12: 0x98, 13: 0xa9, 14: 0xba, 15: 0xcb, 16: 0xdc,
		},
	}
	if !applyOuterCycle(&state, 1) {
		t.Fatal("outer cycle rejected valid input")
	}
	var encoded [20]byte
	for index, value := range state.Words {
		binary.LittleEndian.PutUint32(encoded[index*4:], value)
	}
	hasher := sha256.New()
	_, _ = hasher.Write(encoded[:])
	var wide [8]byte
	for _, value := range state.CycleRegisters {
		binary.LittleEndian.PutUint64(wide[:], value)
		_, _ = hasher.Write(wide[:])
	}
	if got, want := fmt.Sprintf("%x", hasher.Sum(nil)), "209d67e44059e6e3201b84db7ebd4efd3a7306825fe50e6222d8a32dd68bd46a"; got != want {
		t.Fatalf("cycle digest = %s, want %s", got, want)
	}
}

func TestOuterCycleRejectsInvalidInput(t *testing.T) {
	var state outerTableRoundState
	if applyOuterCycle(nil, 1) || applyOuterCycle(&state, 0) {
		t.Fatal("accepted invalid cycle input")
	}
	state.HardwareIndex = 6
	if applyOuterCycle(&state, 1) {
		t.Fatal("accepted invalid cycle hardware index")
	}
}
