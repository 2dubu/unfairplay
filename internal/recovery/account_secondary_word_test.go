package recovery

import (
	"math/rand"
	"reflect"
	"testing"
)

func TestAccountEntrySecondaryWordMatchesSteps(t *testing.T) {
	random := rand.New(rand.NewSource(int64(accountEntrySecondaryWordPass)))
	const schedule, dispatch = uint64(0x200000), uint64(0x300000)
	const exitTarget = uint64(0x123456789)
	const scheduleDisplacement = uint64(0xfffffffedf91bca4)
	for trial := 0; trial < 512; trial++ {
		length := 1 + random.Intn(128)
		selector := uint32(1 + random.Intn(255))
		terminal := uint64(selector+0x42) + 0x481b782a
		start := terminal - uint64(length)
		dispatchBase := addOffset(dispatch, -int64(int32(selector))*8)
		exitAddress := addOffset(dispatchBase, int64(int32(selector^0x46))*8)
		pages := map[uint64][]byte{
			schedule:                       make([]byte, pageSize),
			dispatch &^ (pageSize - 1):     make([]byte, pageSize),
			exitAddress &^ (pageSize - 1):  make([]byte, pageSize),
			dispatchBase &^ (pageSize - 1): make([]byte, pageSize),
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
		base.x[3] = start
		base.x[9] = scheduleDisplacement
		base.x[13] = uint64(random.Uint32())
		base.x[15] = uint64(selector)
		scheduleStart := uint64(schedule + 64)
		base.x[20] = scheduleStart - scheduleDisplacement - start*4
		base.x[24] = dispatchBase
		if err := base.storeValue(dispatch, accountEntrySecondaryWordPass, 8); err != nil {
			t.Fatal(err)
		}
		if err := base.storeValue(exitAddress, exitTarget, 8); err != nil {
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
			repeatedNext, err = runAccountEntrySecondaryWordStep(repeated)
			if err != nil {
				t.Fatal(err)
			}
		}
		semanticNext, err := runAccountEntrySecondaryWordPass(semantic)
		if err != nil {
			t.Fatal(err)
		}
		if repeatedNext != semanticNext || repeated.x != semantic.x || repeated.n != semantic.n || repeated.z != semantic.z || repeated.c != semantic.c || repeated.v != semantic.v || !reflect.DeepEqual(repeated.writable, semantic.writable) {
			t.Fatalf("trial %d state differs", trial)
		}
	}
}
