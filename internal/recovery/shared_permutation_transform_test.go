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

		expected := newState(base.model)
		semantic := newState(base.model)
		for _, item := range []*state{expected, semantic} {
			item.sp = base.sp
			item.x = base.x
			item.n, item.z, item.c, item.v = base.n, base.z, base.c, base.v
			item.writable = clonePages(base.writable)
			item.dirty = make(map[uint64]bool, len(item.writable))
			for page := range item.writable {
				item.dirty[page] = true
			}
		}
		selectorValue, err := expected.load(work+uint64(byte(uint32(expected.x[11])-0x62)), 1)
		if err != nil {
			t.Fatal(err)
		}
		rotation := int((byte(selectorValue) + 5) & 7)
		rotated := bits.RotateLeft8(byte(expected.x[20])^0x1b, rotation)
		if err := expected.storeValue(work+expected.x[0], uint64(rotated^0x1b), 1); err != nil {
			t.Fatal(err)
		}
		second, err := expected.load(work+uint64(byte(uint32(expected.x[11])+0xae)), 1)
		if err != nil {
			t.Fatal(err)
		}
		thirdIndex := byte(((uint32(expected.x[19]) & 0x7f) << 1) ^ 0x52)
		third, err := expected.load(work+uint64(thirdIndex), 1)
		if err != nil {
			t.Fatal(err)
		}
		fourthIndex := byte(uint32(expected.x[12])*uint32(expected.x[16]) + 0xa2)
		fourth, err := expected.load(work+uint64(fourthIndex), 1)
		if err != nil {
			t.Fatal(err)
		}
		word := uint32(rotated)<<24 |
			uint32(byte(second)^0x1b)<<16 |
			uint32(byte(third)^0x1b)<<8 |
			uint32(byte(fourth)^0x1b)
		output := sharedPermutationPolynomialCircuit(word)
		expected.x[1] = 0x65c0
		expected.x[2] = target
		expected.x[17] = 0xcd20
		expected.x[19] = 0x71
		expected.x[20] = uint64(output)
		expected.x[27] = bootstrapTargetTable
		expected.x[28] = work
		expectedNext := target
		semanticNext, err := runSharedPermutationSemantic(semantic)
		if err != nil {
			t.Fatalf("trial %d semantic: %v", trial, err)
		}
		if expectedNext != semanticNext || !reflect.DeepEqual(expected.writable, semantic.writable) {
			t.Fatalf("trial %d output differs", trial)
		}
		for _, register := range []int{1, 2, 3, 17, 19, 20, 27, 28} {
			if expected.x[register] != semantic.x[register] {
				t.Fatalf("trial %d x%d: expected=%#x semantic=%#x", trial, register, expected.x[register], semantic.x[register])
			}
		}
	}
}
