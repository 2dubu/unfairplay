package recovery

import (
	"math/rand"
	"reflect"
	"testing"
)

func TestAccountEntryZeroFillMatchesSteps(t *testing.T) {
	random := rand.New(rand.NewSource(int64(accountEntryZeroFillPass)))
	const output, dispatch = uint64(0x200000), uint64(0x300000)
	const exitTarget = uint64(0x123456789)
	for trial := 0; trial < 1024; trial++ {
		length := 1 + random.Intn(512)
		start := random.Uint64() & 0xffff
		terminal := start + uint64(length-1)
		selector := uint32(1 + random.Intn(255))
		pages := map[uint64][]byte{
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
		base.x[8] = random.Uint64()
		base.x[9] = 0
		base.x[10] = start
		base.x[14] = base.x[8] ^ (terminal - 0x48f)
		base.x[15] = uint64(selector)
		base.x[16] = output - start
		base.x[30] = dispatch
		if err := base.storeValue(dispatch, accountEntryZeroFillPass, 8); err != nil {
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
			repeatedNext, err = runAccountEntryZeroFillStep(repeated)
			if err != nil {
				t.Fatal(err)
			}
		}
		semanticNext, err := runAccountEntryZeroFill(semantic)
		if err != nil {
			t.Fatal(err)
		}
		if repeatedNext != semanticNext || repeated.x != semantic.x || repeated.n != semantic.n || repeated.z != semantic.z || repeated.c != semantic.c || repeated.v != semantic.v || !reflect.DeepEqual(repeated.writable, semantic.writable) {
			t.Fatalf("trial %d state differs", trial)
		}
	}
}
