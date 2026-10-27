package recovery

import (
	"math/rand"
	"testing"
)

func sharedMBAFirstIndexCircuit(counter uint32) byte {
	otherCounter := uint32(0x24a2fbac - counter)
	x1 := uint32(0x19d6 ^ 0x7c16)
	left := (counter ^ 0x2d) + 0x7f
	right := (counter ^ 0x23) + 0x71
	mixed := ((x1 + 0xffff9a9d) ^ counter) + 0xf
	return byte((counter^otherCounter)^right^left^mixed) ^ 0xac
}

func sharedMBAFoldByteCircuit(value byte) uint32 {
	word := uint32(value)
	left := (word ^ 0xcc667f2c) + 0x339980d4
	right := (word ^ 0x0aa63e3e) + 0xf559c1c2
	tail := (word ^ 0xc6c04109) + 0x393fbef7
	return left ^ right ^ tail
}

func TestSharedMBAFirstIndexMatchesCircuit(t *testing.T) {
	random := rand.New(rand.NewSource(int64(sharedMBAFront)))
	for trial := 0; trial < 4096; trial++ {
		counter := random.Uint32()
		if got, want := sharedMBAFirstIndex(counter), sharedMBAFirstIndexCircuit(counter); got != want {
			t.Fatalf("counter %#x: got %#x want %#x", counter, got, want)
		}
	}
}

func TestSharedMBAFoldByteMatchesCircuit(t *testing.T) {
	for value := 0; value < 256; value++ {
		if got, want := sharedMBAFoldByte(byte(value)), sharedMBAFoldByteCircuit(byte(value)); got != want {
			t.Fatalf("value %#x: got %#x want %#x", value, got, want)
		}
	}
}

func sharedMBARemainderCircuit(folded, stateWord uint32) int32 {
	x17 := folded
	x2 := folded + 0x1b
	x20 := stateWord*0x75 + 0x48575af0
	x23 := stateWord * 0xb7a8b167
	x24 := stateWord + 0x1b
	x20 = x20*x24 + x23
	x23 = stateWord + 0xc42b10eb
	x25 := x20 * 0x3bd4ef30
	x20 += 0x5ecab5dd
	x25 = stateWord*0x9e9ca0fb + x25
	x20 = x20*x23 + x25
	x23 = x17*0xb + 0xa40d49ed
	x23 = x23*x2 + x17*0x5bf2b73c
	x23 += 0x4ed6e7d8
	x23 = x23*x2 + x17*0x63c26b7c
	x25 = x17 + 0xc961596a
	x27 := x23 + 0x85815614
	x25 = x27 * x25
	x27 = x17 + 0x9d449283
	x27 = x27*x2 + x17*0x62bb6d98
	x23 = x23*0x369ea6b1 + x25
	x25 = x17 + 0x6345d1c0
	x20 = x27*0x9cba2e5b + x20
	x27 += 0x69c48f08
	x20 = x27*x25 + x20
	x25 = x17*0xe7 + 0xd7876c42
	x27 = x17 * 0x2878ac1b
	x27 = stateWord*0xd8f6830e + x27
	x25 = x25*x2 + x27
	x23 += 0x1ca7fc33
	x2 = x23*x2 + x17*0x67ed7fa1
	x2 += 0x576fd500
	x2 = x2*x24 + x20
	x2 += 0x27c02bb3
	x2 = x2*x24 + x25
	x20 = stateWord + 0xcda8f9e7
	x17 = x2*0x32570634 + x17
	x2 += 0x26b9f953
	x17 = x2*x20 + x17
	x17 += 0x2ad396f7
	return int32(x17) % 599
}

func TestSharedMBARemainderMatchesCircuit(t *testing.T) {
	random := rand.New(rand.NewSource(0x257))
	for trial := 0; trial < 100000; trial++ {
		folded := random.Uint32()
		stateWord := random.Uint32()
		if got, want := sharedMBARemainder(folded, stateWord), sharedMBARemainderCircuit(folded, stateWord); got != want {
			t.Fatalf("trial %d: got %d want %d", trial, got, want)
		}
	}
}

func sharedMBATailIndexCircuit(remainder int32) byte {
	x17 := uint32(remainder)
	x2 := uint32(0x6e6cd6a0) & (x17 << 1)
	x17 -= x2
	x17 += 0x37366b50
	x2 = (x17 ^ 0x196dc519) + 0xe6923ae7
	x20 := (x17 ^ 0xb118082f) + 0x4ee7f7d1
	x2 ^= x20
	x17 = (x17 ^ 0x9f43a666) + 0x60bc599a
	x17 ^= x2
	x17 += x17 << 7
	x17 += 0x2ac71237
	x2 = x17 ^ 0xc5c184e0
	x20 = x17 & 0xa7a501e2
	x22 := (x17 ^ 0x878185e0) & (x20 << 1)
	x20 = x22 ^ x20
	x22 = 0xc4c90a04 ^ (x2 << 1)
	x2 = x22 & x2
	x22 = x2 & (x20 << 2)
	x20 ^= 0x62648502
	x20 ^= x22
	x22 = (x20 << 4) ^ 0x20
	x2 &= x2 << 2
	x2 &= x22
	x2 ^= x20
	x17 ^= x2 << 1
	return byte(x17) ^ 0xe6
}

func TestSharedMBATailIndexMatchesCircuit(t *testing.T) {
	for remainder := int32(-598); remainder <= 598; remainder++ {
		got := sharedMBATailIndex(byte(remainder))
		want := sharedMBATailIndexCircuit(remainder)
		if got != want {
			t.Fatalf("remainder %d: got %#x want %#x", remainder, got, want)
		}
	}
}

func TestSharedMBATailIndexIsPermutation(t *testing.T) {
	seen := [256]bool{}
	for value := 0; value < 256; value++ {
		output := sharedMBATailIndex(byte(value))
		if seen[output] {
			t.Fatalf("duplicate output %#x", output)
		}
		seen[output] = true
	}
}
