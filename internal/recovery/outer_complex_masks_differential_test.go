package recovery

import (
	"math/rand"
	"testing"
)

func TestOuterTableLane7MaskMatchesCircuit(t *testing.T) {
	random := rand.New(rand.NewSource(0x26a428))
	for sample := 0; sample < 100000; sample++ {
		var input [5]uint32
		for index := range input {
			input[index] = random.Uint32()
		}
		got := outerTableLane7Mask(input[0], input[1], input[2], input[3], input[4])
		want := outerTableLane7MaskCircuit(input[0], input[1], input[2], input[3], input[4])
		if got != want {
			t.Fatalf("sample %d = %08x, want %08x", sample, got, want)
		}
	}
}

func TestOuterTableLane2MaskMatchesCircuit(t *testing.T) {
	random := rand.New(rand.NewSource(0x26a428))
	for sample := 0; sample < 100000; sample++ {
		var input [5]uint32
		for index := range input {
			input[index] = random.Uint32()
		}
		got := outerTableLane2Mask(input[0], input[1], input[2], input[3], input[4])
		want := outerTableLane2MaskCircuit(input[0], input[1], input[2], input[3], input[4])
		if got != want {
			t.Fatalf("sample %d = %08x, want %08x", sample, got, want)
		}
	}
}

func TestOuterTableLane19MaskMatchesCircuit(t *testing.T) {
	random := rand.New(rand.NewSource(0x26a428))
	for sample := 0; sample < 100000; sample++ {
		var input [6]uint32
		for index := range input {
			input[index] = random.Uint32()
		}
		got := outerTableLane19Mask(input[0], input[1], input[2], input[3], input[4], input[5])
		want := outerTableLane19MaskCircuit(input[0], input[1], input[2], input[3], input[4], input[5])
		if got != want {
			t.Fatalf("sample %d = %08x, want %08x", sample, got, want)
		}
	}
}

func TestOuterTableLane18MaskMatchesCircuit(t *testing.T) {
	random := rand.New(rand.NewSource(0x26a428))
	for sample := 0; sample < 100000; sample++ {
		var input [5]uint32
		for index := range input {
			input[index] = random.Uint32()
		}
		got := outerTableLane18Mask(input[0], input[1], input[2], input[3], input[4])
		want := outerTableLane18MaskCircuit(input[0], input[1], input[2], input[3], input[4])
		if got != want {
			t.Fatalf("sample %d = %08x, want %08x", sample, got, want)
		}
	}
}
