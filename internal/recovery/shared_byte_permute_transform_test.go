package recovery

import (
	"math/rand"
	"reflect"
	"testing"
)

func sharedBytePermuteFirstCircuit(input byte) byte {
	value := uint32(input)
	left := (value ^ (0x19d6 ^ uint32(0xffffffe8))) - 0x3e
	middle := (value ^ 0x49) - 0x49
	right := (value ^ 0x6c) - 0x6c
	right ^= middle ^ left
	right *= 5
	right = 0xfffffff6 - right
	left = (right & 0x40) ^ 0x74
	middle = right ^ 0x18
	carry := (middle & 0x28) ^ left
	for range 2 {
		carry = (middle & (carry << 1)) ^ left
	}
	left = (middle & (carry << 1)) ^ left
	right ^= left << 1
	return byte(right ^ 0x5b)
}

func sharedBytePermuteSecondIndexCircuit(counter byte) byte {
	value := uint32(counter) - 6
	left := value ^ 0x3ff9651d
	middle := value & 0xba1c1cae
	right := (value ^ 0x3bd95c0d) & (middle << 1)
	middle = right ^ middle
	right = left << 1
	mask := right ^ 0x0bcaf366
	left = mask & left
	right &= 0x85e579b2
	left ^= right
	right = left ^ 0x84250891
	left ^= 0x00c07120
	left &= middle << 2
	left ^= middle
	middle = right << 2
	mask = middle ^ 0x1795e6c0
	right = mask & right
	middle &= 0x05e579a0
	middle = right ^ middle
	middle &= left << 4
	left = middle ^ left
	value ^= left << 1
	return byte(value) ^ 0xae
}

func sharedBytePermuteFinalCircuit(input byte) byte {
	value := uint32(input)
	left := (value ^ 0xffffffe0) + 0x32
	middle := (value ^ 0x52) - 0x7c
	left ^= middle
	value = ((value ^ 0xffffffe9) + 0x39) ^ left
	value -= 0x1d
	left = (value & 0xb) ^ 0x5c
	middle = value ^ 0x12
	carry := (middle & (value << 1) & 0x26) ^ left
	for range 4 {
		carry = (middle & (carry << 1)) ^ left
	}
	left = (middle & (carry << 1)) ^ left
	value ^= left << 1
	return byte(value ^ 0x10)
}

func TestSharedBytePermuteTablesMatchCircuits(t *testing.T) {
	for input := 0; input < 256; input++ {
		value := byte(input)
		if got, want := sharedBytePermuteFirst[value], sharedBytePermuteFirstCircuit(value); got != want {
			t.Fatalf("first input %#x: got %#x want %#x", input, got, want)
		}
		if got, want := byte(uint32(value)+0xa8), sharedBytePermuteSecondIndexCircuit(value); got != want {
			t.Fatalf("second input %#x: got %#x want %#x", input, got, want)
		}
		if got, want := sharedBytePermuteFinal[value], sharedBytePermuteFinalCircuit(value); got != want {
			t.Fatalf("final input %#x: got %#x want %#x", input, got, want)
		}
	}
}

func TestSharedBytePermuteTablesArePermutations(t *testing.T) {
	for name, table := range map[string]string{
		"first": sharedBytePermuteFirst,
		"final": sharedBytePermuteFinal,
	} {
		if len(table) != 256 {
			t.Fatalf("%s table has %d entries", name, len(table))
		}
		seen := [256]bool{}
		for index := range table {
			value := table[index]
			if seen[value] {
				t.Fatalf("%s table repeats %#x", name, value)
			}
			seen[value] = true
		}
	}
}

func TestSharedBytePermuteSemanticMatchesCircuit(t *testing.T) {
	random := rand.New(rand.NewSource(int64(sharedEnrollmentBytePermute)))
	const (
		stack    = uint64(0x100100)
		source   = uint64(0x200000)
		work     = uint64(0x300000)
		dispatch = uint64(0x400000)
		target   = uint64(0x123456789)
	)
	for trial := 0; trial < 2048; trial++ {
		selector := random.Uint32()
		dispatchAddress := addOffset(dispatch, int64(int32(selector))*8)
		pages := map[uint64][]byte{
			stack &^ (pageSize - 1): make([]byte, pageSize),
			source:                  make([]byte, pageSize), work: make([]byte, pageSize),
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
		base.sp = stack
		base.x[0] = uint64(byte(random.Uint32()))
		base.x[1] = uint64(selector)
		base.x[2] = 0x19d6
		base.x[5] = 0xba1c1cae
		base.x[6] = 0xae
		base.x[7] = 0x52
		base.x[8] = 0x1b
		base.x[11] = uint64(random.Uint32())
		base.x[20] = uint64(byte(random.Uint32()))
		base.x[21] = dispatch
		if err := base.storeValue(stack+0x50, source, 8); err != nil {
			t.Fatal(err)
		}
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
			item.dirty = map[uint64]bool{}
			for page := range item.writable {
				item.dirty[page] = true
			}
		}
		sourceByte, err := expected.load(source+uint64(byte(expected.x[20]^expected.x[8])), 1)
		if err != nil {
			t.Fatal(err)
		}
		first := sharedBytePermuteFirstCircuit(byte(sourceByte))
		if err := expected.storeValue(work+expected.x[0], uint64(first), 1); err != nil {
			t.Fatal(err)
		}
		secondIndex := sharedBytePermuteSecondIndexCircuit(byte(expected.x[11]))
		second, err := expected.load(work+uint64(secondIndex), 1)
		if err != nil {
			t.Fatal(err)
		}
		output := sharedBytePermuteFinalCircuit(first ^ 0x5b ^ byte(second))
		if err := expected.storeValue(work+expected.x[0], uint64(output), 1); err != nil {
			t.Fatal(err)
		}
		expected.x[1] = target
		expected.x[3] = 0x060be5d0
		expected.x[17] = 0xbf5f
		expected.x[20] = uint64(output)
		expected.x[27] = work
		expectedNext := target
		semanticNext, semanticErr := runSharedEnrollmentBytePermuteSemantic(semantic)
		if semanticErr != nil {
			t.Fatalf("trial %d semantic error: %v", trial, semanticErr)
		}
		if expectedNext != semanticNext || !reflect.DeepEqual(expected.writable, semantic.writable) {
			expectedByte, _ := expected.load(work+base.x[0], 1)
			semanticByte, _ := semantic.load(work+base.x[0], 1)
			t.Fatalf("trial %d output differs: next=%#x/%#x byte=%#x/%#x", trial,
				expectedNext, semanticNext, expectedByte, semanticByte)
		}
		for _, register := range []int{1, 2, 3, 17, 19, 20, 27, 28} {
			if expected.x[register] != semantic.x[register] {
				t.Fatalf("trial %d x%d: expected=%#x semantic=%#x", trial, register, expected.x[register], semantic.x[register])
			}
		}
	}
}
