package recovery

import (
	"encoding/binary"
	"math/rand"
	"testing"
)

func TestSharedEnrollmentWordCompareLoopMatchesSteps(t *testing.T) {
	random := rand.New(rand.NewSource(int64(sharedEnrollmentWordCompare)))
	const (
		words      = uint64(0x100000)
		dispatch   = uint64(0x200000)
		exitTarget = uint64(0x123456789)
	)
	for trial := 0; trial < 1024; trial++ {
		selector := random.Uint32()
		equalIndex := selector + 0x526c
		loopWord := (selector + 0x3d68) | 0x8e1
		loopMask := loopWord ^ 0x4aea
		loopBase := selector + 0xc153
		moreIndex := loopBase ^ loopMask
		pages := map[uint64][]byte{words: make([]byte, pageSize)}
		for _, index := range []uint32{equalIndex, equalIndex ^ 9, moreIndex, loopBase} {
			address := addOffset(dispatch, int64(int32(index))*8)
			if pages[address&^(pageSize-1)] == nil {
				pages[address&^(pageSize-1)] = make([]byte, pageSize)
			}
		}
		word := random.Uint32()
		for index := 0; index < 0x2c; index++ {
			binary.LittleEndian.PutUint32(pages[words][index*4:], word)
		}
		model := &model{pages: clonePages(pages)}
		base := newState(model)
		for register := range base.x {
			base.x[register] = random.Uint64()
		}
		base.x[8] = 0
		base.x[9] = words
		base.x[10] = 1
		base.x[21] = dispatch
		base.x[24] = uint64(selector)
		if err := base.storeValue(addOffset(dispatch, int64(int32(equalIndex))*8), sharedEnrollmentLoopSelect, 8); err != nil {
			t.Fatal(err)
		}
		if err := base.storeValue(addOffset(dispatch, int64(int32(equalIndex^9))*8), exitTarget, 8); err != nil {
			t.Fatal(err)
		}
		if err := base.storeValue(addOffset(dispatch, int64(int32(moreIndex))*8), sharedEnrollmentWordCompare, 8); err != nil {
			t.Fatal(err)
		}
		if err := base.storeValue(addOffset(dispatch, int64(int32(loopBase))*8), exitTarget, 8); err != nil {
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
		for repeatedNext == 0 || repeatedNext == sharedEnrollmentWordCompare {
			var err error
			repeatedNext, err = runSharedEnrollmentWordCompare(repeated)
			if err != nil {
				t.Fatal(err)
			}
			if repeatedNext == sharedEnrollmentLoopSelect {
				repeatedNext, err = runSharedEnrollmentLoopSelect(repeated)
				if err != nil {
					t.Fatal(err)
				}
			}
		}
		semanticNext, err := runSharedEnrollmentWordCompareLoop(semantic)
		if err != nil {
			t.Fatal(err)
		}
		if repeatedNext != semanticNext || repeated.x != semantic.x || repeated.n != semantic.n || repeated.z != semantic.z || repeated.c != semantic.c || repeated.v != semantic.v {
			t.Fatalf("trial %d state differs", trial)
		}
	}
}
