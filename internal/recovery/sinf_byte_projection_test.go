package recovery

import (
	"math/rand"
	"reflect"
	"testing"
)

func TestSINFByteProjectionMatchesSteps(t *testing.T) {
	random := rand.New(rand.NewSource(int64(sinfByteProjectionPass)))
	const input, output, laneA, laneB, laneC, dispatch = uint64(0x100000), uint64(0x200000), uint64(0x300000), uint64(0x400000), uint64(0x500000), uint64(0x600000)
	const exitTarget = uint64(0x123456789)
	for trial := 0; trial < 512; trial++ {
		length := 1 + random.Intn(512)
		pages := map[uint64][]byte{
			input:    make([]byte, pageSize),
			output:   make([]byte, pageSize),
			laneA:    make([]byte, pageSize),
			laneB:    make([]byte, pageSize),
			laneC:    make([]byte, pageSize),
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
		base.x[0] = laneC
		base.x[1] = laneA
		base.x[2] = laneB
		base.x[4] = uint64(length)
		base.x[6] = uint64(random.Uint32())
		base.x[7] = uint64(random.Uint32())
		base.x[8] = dispatch
		base.x[12] = output
		base.x[13] = 0
		deltaBase := uint64(0x7ce)
		base.x[5] = deltaBase ^ ^uint64(0)
		base.x[16] = input
		base.x[19] = 1
		if err := base.storeValue(dispatch, sinfByteProjectionPass, 8); err != nil {
			t.Fatal(err)
		}
		if err := base.storeValue(dispatch+8, exitTarget, 8); err != nil {
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
			repeatedNext, err = runSINFByteProjectionStep(repeated)
			if err != nil {
				t.Fatal(err)
			}
		}
		semanticNext, err := runSINFByteProjection(semantic)
		if err != nil {
			t.Fatal(err)
		}
		if repeatedNext != semanticNext || repeated.x != semantic.x || repeated.n != semantic.n || repeated.z != semantic.z || repeated.c != semantic.c || repeated.v != semantic.v || !reflect.DeepEqual(repeated.writable, semantic.writable) {
			t.Fatalf("trial %d state differs", trial)
		}
	}
}
