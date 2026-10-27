package recovery

import (
	"math/bits"
	"math/rand"
	"reflect"
	"testing"
)

func sharedPermutationPolynomialCircuit(word uint32) uint32 {
	mixed := (word << 3) ^ (word >> 5)
	word = mixed + word
	square := word * word
	highByte := word >> 8
	highWord := word >> 16
	result := square*word + highWord
	square = highByte * highByte
	highByte = square * highByte
	result = highByte*square + result
	square = highWord * highWord
	square *= highWord
	highByte = square * highWord
	result = highByte*square + result
	result -= 0xd4 & (result << 1)
	result += 0x6a
	return result ^ 0x71
}

func TestSharedPermutationPolynomialMatchesCircuit(t *testing.T) {
	random := rand.New(rand.NewSource(int64(sharedPermutationFront)))
	for trial := 0; trial < 100000; trial++ {
		word := random.Uint32()
		if got, want := sharedPermutationPolynomial(word), sharedPermutationPolynomialCircuit(word); got != want {
			t.Fatalf("trial %d word %#x: got %#x want %#x", trial, word, got, want)
		}
	}
}

func TestSharedPermutationRotation(t *testing.T) {
	for selector := 0; selector < 8; selector++ {
		for current := 0; current < 256; current++ {
			got := bits.RotateLeft8(byte(current)^0x1b, (selector+5)&7) ^ 0x1b
			want := bits.RotateLeft8(byte(current)^0x1b, selector-3) ^ 0x1b
			if got != want {
				t.Fatalf("selector %d current %#x: got %#x want %#x", selector, current, got, want)
			}
		}
	}
}

func TestSharedPermutationSemanticMatchesCircuit(t *testing.T) {
	random := rand.New(rand.NewSource(0x43a000))
	const (
		stack    = uint64(0x100100)
		work     = uint64(0x200000)
		dispatch = uint64(0x300000)
		target   = uint64(0x123456789)
	)
	dispatchAddress := dispatch + 0x19d6*8
	for trial := 0; trial < 512; trial++ {
		pages := map[uint64][]byte{
			stack &^ (pageSize - 1):           make([]byte, pageSize),
			work:                              make([]byte, pageSize),
			dispatchAddress &^ (pageSize - 1): make([]byte, pageSize),
		}
		for _, page := range pages {
			if _, err := random.Read(page); err != nil {
				t.Fatal(err)
			}
		}
		base := newState(&model{pages: clonePages(pages)})
		for register := range base.x {
			base.x[register] = random.Uint64()
		}
		iteration := random.Uint32()
		base.sp = stack
		base.x[0] = uint64(byte(random.Uint32()))
		base.x[2] = 0x19d6
		base.x[5] = 0xba1c1cae
		base.x[6] = 0xae
		base.x[7] = 0x52
		base.x[8] = 0x1b
		base.x[11] = uint64(0x45e3e363 + iteration)
		base.x[12] = uint64(0x22f1f257 + iteration)
		base.x[16] = uint64(0x22f1f1a9 + iteration)
		base.x[19] = uint64(random.Uint32())
		base.x[20] = uint64(byte(random.Uint32()))
		base.x[21] = dispatch
		base.x[26] = 0x40cb0000
		base.x[30] = 0xb667b683
		if err := base.storeValue(stack+0x58, work, 8); err != nil {
			t.Fatal(err)
		}
		if err := base.storeValue(dispatchAddress, target, 8); err != nil {
			t.Fatal(err)
		}

		circuit := newState(base.model)
		semantic := newState(base.model)
		for _, item := range []*state{circuit, semantic} {
			item.sp = base.sp
			item.x = base.x
			item.n, item.z, item.c, item.v = base.n, base.z, base.c, base.v
			item.writable = clonePages(base.writable)
			item.dirty = make(map[uint64]bool, len(item.writable))
			for page := range item.writable {
				item.dirty[page] = true
			}
		}
		frontNext, err := runSharedPermutationFront(circuit)
		if err != nil || frontNext != sharedPermutationTail {
			t.Fatalf("trial %d front: next=%#x err=%v", trial, frontNext, err)
		}
		circuitNext, err := runSharedPermutationTail(circuit)
		if err != nil {
			t.Fatalf("trial %d tail: %v", trial, err)
		}
		semanticNext, err := runSharedPermutationSemantic(semantic)
		if err != nil {
			t.Fatalf("trial %d semantic: %v", trial, err)
		}
		if circuitNext != semanticNext || !reflect.DeepEqual(circuit.writable, semantic.writable) {
			t.Fatalf("trial %d output differs", trial)
		}
		for _, register := range []int{1, 2, 3, 17, 19, 20, 27, 28} {
			if circuit.x[register] != semantic.x[register] {
				t.Fatalf("trial %d x%d: circuit=%#x semantic=%#x", trial, register, circuit.x[register], semantic.x[register])
			}
		}
	}
}
