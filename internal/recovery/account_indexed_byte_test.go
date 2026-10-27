package recovery

import (
	"math/rand"
	"reflect"
	"testing"
)

func TestAccountEntryIndexedByteMatchesSteps(t *testing.T) {
	random := rand.New(rand.NewSource(int64(accountEntryIndexedBytePass)))
	const input, output, dispatch = uint64(0x200000), uint64(0x300000), uint64(0x400000)
	const exitTarget = uint64(0x123456789)
	for trial := 0; trial < 1024; trial++ {
		length := 1 + random.Intn(512)
		selector := uint32(1 + random.Intn(255))
		pages := map[uint64][]byte{
			input:    make([]byte, pageSize),
			output:   make([]byte, pageSize),
			dispatch: make([]byte, pageSize),
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
		base.x[0] = random.Uint64()
		base.x[2] = 1
		base.x[3] = uint64(length)
		base.x[4] = 0
		base.x[5] = uint64(selector)
		base.x[8] = 0
		base.x[11] = 0
		base.x[12] = 0
		base.x[13] = 0
		base.x[21] = output
		base.x[24] = dispatch
		base.x[26] = input
		if err := base.storeValue(dispatch, accountEntryIndexedBytePass, 8); err != nil {
			t.Fatal(err)
		}
		if err := base.storeValue(dispatch+uint64(selector)*8, exitTarget, 8); err != nil {
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
			repeatedNext, err = runAccountEntryIndexedByteStep(repeated)
			if err != nil {
				t.Fatal(err)
			}
		}
		semanticNext, err := runAccountEntryIndexedBytePass(semantic)
		if err != nil {
			t.Fatal(err)
		}
		if repeatedNext != semanticNext || repeated.x != semantic.x || repeated.n != semantic.n || repeated.z != semantic.z || repeated.c != semantic.c || repeated.v != semantic.v || !reflect.DeepEqual(repeated.writable, semantic.writable) {
			t.Fatalf("trial %d state differs", trial)
		}
	}
}
