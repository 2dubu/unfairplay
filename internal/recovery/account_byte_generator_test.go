package recovery

import (
	"math/rand"
	"reflect"
	"testing"
)

func TestAccountEntryByteGeneratorMatchesSteps(t *testing.T) {
	random := rand.New(rand.NewSource(int64(accountEntryByteGeneratorPass)))
	const stack, output, dispatch = uint64(0x100100), uint64(0x200000), uint64(0x300000)
	const exitTarget = uint64(0x123456789)
	for trial := 0; trial < 1024; trial++ {
		length := 1 + random.Intn(512)
		exitSelector := random.Uint32()
		continueAddress := dispatch
		exitAddress := addOffset(dispatch, int64(int32(exitSelector))*8)
		pages := map[uint64][]byte{stack &^ (pageSize - 1): make([]byte, pageSize), output: make([]byte, pageSize), continueAddress &^ (pageSize - 1): make([]byte, pageSize), exitAddress &^ (pageSize - 1): make([]byte, pageSize)}
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
		base.sp = stack
		base.x[1] = uint64(length)
		base.x[6] = output
		base.x[10] = 1
		base.x[13] = uint64(random.Uint32())
		base.x[14] = uint64(exitSelector)
		base.x[16] = 0
		base.x[19] = ^uint64(1)
		base.x[30] = dispatch
		if err := base.storeValue(continueAddress, accountEntryByteGeneratorPass, 8); err != nil {
			t.Fatal(err)
		}
		if err := base.storeValue(exitAddress, exitTarget, 8); err != nil {
			t.Fatal(err)
		}
		repeated := newState(model)
		semantic := newState(model)
		for _, item := range []*state{repeated, semantic} {
			item.sp, item.x = base.sp, base.x
			item.writable = clonePages(base.writable)
			item.dirty = map[uint64]bool{}
			for page := range item.writable {
				item.dirty[page] = true
			}
		}
		var repeatedNext uint64
		for index := 0; index < length; index++ {
			var err error
			repeatedNext, err = runAccountEntryByteGeneratorStep(repeated)
			if err != nil {
				t.Fatal(err)
			}
		}
		semanticNext, err := runAccountEntryByteGenerator(semantic)
		if err != nil {
			t.Fatal(err)
		}
		if repeatedNext != semanticNext || repeated.x != semantic.x || repeated.n != semantic.n || repeated.z != semantic.z || repeated.c != semantic.c || repeated.v != semantic.v || !reflect.DeepEqual(repeated.writable, semantic.writable) {
			t.Fatalf("trial %d state differs", trial)
		}
	}
}
