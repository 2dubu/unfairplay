package recovery

import (
	"math/rand"
	"reflect"
	"testing"
)

func TestAccountEntryByteStoreMatchesSteps(t *testing.T) {
	random := rand.New(rand.NewSource(int64(accountEntryByteStoreStart)))
	const descriptor, input, output, dispatch = uint64(0x100000), uint64(0x200000), uint64(0x300000), uint64(0x400000)
	const exitTarget = uint64(0x123456789)
	const rangeStart = uint32(0x4539237b)
	const positionBias = uint32(0x3489f9a3)
	for trial := 0; trial < 512; trial++ {
		length := 1 + random.Intn(256)
		position := rangeStart + 1
		counter := position - positionBias
		rangeEnd := position + uint32(length-1)
		encodedOffset := uint32(0x545dd584) & 0x00025fdf
		inputBase := input - uint64(counter+0xef50d627)
		pages := map[uint64][]byte{
			descriptor: make([]byte, pageSize),
			input:      make([]byte, pageSize),
			output:     make([]byte, pageSize),
			dispatch:   make([]byte, pageSize),
		}
		for _, page := range pages {
			if _, err := random.Read(page); err != nil {
				t.Fatal(err)
			}
		}
		model := &model{pages: clonePages(pages)}
		base := newState(model)
		for register := range base.x {
			base.x[register] = random.Uint64()
		}
		base.x[0] = output - 4
		base.x[1] = descriptor
		base.x[2] = 1
		base.x[3] = uint64(counter)
		base.x[8] = 0
		base.x[9] = 0
		base.x[10] = dispatch
		base.x[12] = 0
		base.x[13] = 0
		base.x[14] = 0
		base.x[15] = 0
		base.x[16] = uint64(rangeEnd - 0x13ae9f40 - (encodedOffset + 0xcd463667))
		if err := base.storeValue(descriptor+4, 0, 4); err != nil {
			t.Fatal(err)
		}
		if err := base.storeValue(descriptor+8, inputBase, 8); err != nil {
			t.Fatal(err)
		}
		if err := base.storeValue(dispatch, exitTarget, 8); err != nil {
			t.Fatal(err)
		}
		if err := base.storeValue(dispatch+8, accountEntryByteStoreStart, 8); err != nil {
			t.Fatal(err)
		}

		repeated := newState(model)
		semantic := newState(model)
		for _, item := range []*state{repeated, semantic} {
			item.x = base.x
			item.writable = clonePages(base.writable)
			item.dirty = map[uint64]bool{}
			for page := range item.writable {
				item.dirty[page] = true
			}
		}
		var repeatedNext uint64
		for index := 0; index < length; index++ {
			var err error
			repeatedNext, err = runAccountEntryByteStoreStep(repeated)
			if err != nil {
				t.Fatal(err)
			}
		}
		semanticNext, err := runAccountEntryByteStore(semantic)
		if err != nil {
			t.Fatal(err)
		}
		if repeatedNext != semanticNext || repeated.x != semantic.x || repeated.n != semantic.n || repeated.z != semantic.z || repeated.c != semantic.c || repeated.v != semantic.v || !reflect.DeepEqual(repeated.writable, semantic.writable) {
			t.Fatalf("trial %d state differs", trial)
		}
	}
}
