package recovery

import (
	"math/rand"
	"reflect"
	"testing"
)

func TestSelectionSIDBProjectionMatchesSteps(t *testing.T) {
	random := rand.New(rand.NewSource(int64(selectionSIDBProjectionPass)))
	const input, output = uint64(0x100000), uint64(0x200000)
	const laneA, laneB, laneC = uint64(0x300000), uint64(0x301000), uint64(0x302000)
	const dispatch, frame = uint64(0x400000), uint64(0x500100)
	const exitTarget = uint64(0x123456789)
	for trial := 0; trial < 1024; trial++ {
		length := 1 + random.Intn(512)
		selector := random.Uint32()
		continueAddress := addOffset(dispatch, int64(int32(selector^0x14))*8)
		exitAddress := addOffset(dispatch, int64(int32(selector))*8)
		pages := map[uint64][]byte{input: make([]byte, pageSize), output: make([]byte, pageSize), laneA: make([]byte, pageSize), laneB: make([]byte, pageSize), laneC: make([]byte, pageSize), frame &^ (pageSize - 1): make([]byte, pageSize), continueAddress &^ (pageSize - 1): make([]byte, pageSize), exitAddress &^ (pageSize - 1): make([]byte, pageSize)}
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
		base.x[0] = uint64(selector)
		base.x[9] = input
		base.x[10] = uint64(selector)
		base.x[11] = laneC
		base.x[12] = uint64(length)
		base.x[13] = 0
		base.x[14] = laneA
		base.x[15] = uint64(random.Uint32())
		base.x[16] = laneB
		base.x[20] = output
		base.x[29] = frame
		if err := base.storeValue(addOffset(frame, -0x90), dispatch, 8); err != nil {
			t.Fatal(err)
		}
		if err := base.storeValue(continueAddress, selectionSIDBProjectionPass, 8); err != nil {
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
			repeatedNext, err = runSelectionSIDBProjectionStep(repeated)
			if err != nil {
				t.Fatal(err)
			}
		}
		semanticNext, err := runSelectionSIDBProjection(semantic)
		if err != nil {
			t.Fatal(err)
		}
		if repeatedNext != semanticNext || repeated.x != semantic.x || repeated.n != semantic.n || repeated.z != semantic.z || repeated.c != semantic.c || repeated.v != semantic.v || !reflect.DeepEqual(repeated.writable, semantic.writable) {
			t.Fatalf("trial %d state differs", trial)
		}
	}
}
