package recovery

import (
	"math/rand"
	"reflect"
	"testing"
)

func TestAccountEntryByteCopyMatchesSteps(t *testing.T) {
	random := rand.New(rand.NewSource(int64(accountEntryByteCopyPass)))
	const input, output, dispatch = uint64(0x200000), uint64(0x300000), uint64(0x400000)
	const exitTarget = uint64(0x123456789)
	const selector = uint32(0x4991a258)
	const threshold = uint32(0x1de54aff)
	for trial := 0; trial < 512; trial++ {
		length := 1 + random.Intn(512)
		start := threshold - uint32(length-1)
		dispatchBase := addOffset(dispatch, -int64(int32(selector))*8)
		continueAddress := addOffset(dispatchBase, int64(int32(selector^0x172c))*8)
		pages := map[uint64][]byte{
			input:                             make([]byte, pageSize),
			output:                            make([]byte, pageSize),
			dispatch &^ (pageSize - 1):        make([]byte, pageSize),
			continueAddress &^ (pageSize - 1): make([]byte, pageSize),
			dispatchBase &^ (pageSize - 1):    make([]byte, pageSize),
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
		base.x[8] = output - uint64(start)
		base.x[9] = input - uint64(start)
		base.x[12] = uint64(selector)
		base.x[13] = uint64(start)
		base.x[14] = 1
		base.x[15] = 0
		base.x[20] = dispatchBase
		if err := base.storeValue(dispatch, exitTarget, 8); err != nil {
			t.Fatal(err)
		}
		if err := base.storeValue(continueAddress, accountEntryByteCopyPass, 8); err != nil {
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
			repeatedNext, err = runAccountEntryByteCopyStep(repeated)
			if err != nil {
				t.Fatal(err)
			}
		}
		semanticNext, err := runAccountEntryByteCopy(semantic)
		if err != nil {
			t.Fatal(err)
		}
		if repeatedNext != semanticNext || repeated.x != semantic.x || repeated.n != semantic.n || repeated.z != semantic.z || repeated.c != semantic.c || repeated.v != semantic.v || !reflect.DeepEqual(repeated.writable, semantic.writable) {
			t.Fatalf("trial %d state differs", trial)
		}
	}
}
