package recovery

import (
	"math/rand"
	"reflect"
	"testing"
)

func TestSelectionReverseProjectionMatchesSteps(t *testing.T) {
	random := rand.New(rand.NewSource(int64(selectionReverseProjection)))
	const input, output, laneA, laneB, laneC, dispatch = uint64(0x100000), uint64(0x200000), uint64(0x300000), uint64(0x400000), uint64(0x500000), uint64(0x600000)
	const exitTarget = uint64(0x123456789)
	for trial := 0; trial < 512; trial++ {
		length := 1 + random.Intn(512)
		pages := map[uint64][]byte{
			input: make([]byte, pageSize), output: make([]byte, pageSize),
			laneA: make([]byte, pageSize), laneB: make([]byte, pageSize), laneC: make([]byte, pageSize),
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
		base.x[1] = 0
		base.x[3] = 1
		base.x[8] = input
		base.x[10] = uint64(length)
		base.x[12] = laneB
		base.x[13] = laneC
		base.x[14] = output
		base.x[15] = laneA
		base.x[22] = dispatch
		if err := base.storeValue(dispatch, exitTarget, 8); err != nil {
			t.Fatal(err)
		}
		if err := base.storeValue(dispatch+8, selectionReverseProjection, 8); err != nil {
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
			repeatedNext, err = runSelectionDirectBlock10013df58(repeated)
			if err != nil {
				t.Fatal(err)
			}
		}
		semanticNext, err := runSelectionReverseProjection(semantic)
		if err != nil {
			t.Fatal(err)
		}
		if repeatedNext != semanticNext || repeated.x != semantic.x || repeated.n != semantic.n || repeated.z != semantic.z || repeated.c != semantic.c || repeated.v != semantic.v || !reflect.DeepEqual(repeated.writable, semantic.writable) {
			t.Fatalf("trial %d state differs", trial)
		}
	}
}
