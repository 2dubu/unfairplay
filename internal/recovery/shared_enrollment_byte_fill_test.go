package recovery

import (
	"math/rand"
	"reflect"
	"testing"
)

func TestSharedEnrollmentByteFillMatchesSteps(t *testing.T) {
	random := rand.New(rand.NewSource(int64(sharedEnrollmentByteFill)))
	const (
		stack      = uint64(0x100100)
		output     = uint64(0x200000)
		dispatch   = uint64(0x300000)
		exitTarget = uint64(0x123456789)
	)
	for trial := 0; trial < 1024; trial++ {
		selector := random.Uint32()
		continueAddress := addOffset(dispatch, int64(int32(selector))*8)
		exitAddress := addOffset(dispatch, int64(int32(selector^8))*8)
		pages := map[uint64][]byte{
			stack &^ (pageSize - 1):            make([]byte, pageSize),
			(stack + 0x3af8) &^ (pageSize - 1): make([]byte, pageSize),
			output:                             make([]byte, pageSize),
			continueAddress &^ (pageSize - 1):  make([]byte, pageSize),
			exitAddress &^ (pageSize - 1):      make([]byte, pageSize),
		}
		for _, page := range pages {
			if _, err := random.Read(page); err != nil {
				t.Fatal(err)
			}
		}
		base := newState(&model{pages: clonePages(pages)})
		for register := range base.x {
			base.x[register] = random.Uint64()
		}
		length := 1 + random.Intn(512)
		base.sp = stack
		base.x[1] = uint64(random.Uint32())
		base.x[3] = random.Uint64()
		base.x[8] = uint64(length)
		base.x[9] = uint64(selector)
		base.x[10] = random.Uint64()
		base.x[20] = output
		base.x[30] = uint64(random.Uint32())
		if err := base.storeValue(stack+0xdf0, dispatch, 8); err != nil {
			t.Fatal(err)
		}
		if err := base.storeValue(continueAddress, sharedEnrollmentByteFill, 8); err != nil {
			t.Fatal(err)
		}
		if err := base.storeValue(exitAddress, exitTarget, 8); err != nil {
			t.Fatal(err)
		}

		repeated := newState(base.model)
		semantic := newState(base.model)
		for _, item := range []*state{repeated, semantic} {
			item.sp, item.x = base.sp, base.x
			item.n, item.z, item.c, item.v = base.n, base.z, base.c, base.v
			item.writable = clonePages(base.writable)
			item.dirty = map[uint64]bool{}
			for page := range item.writable {
				item.dirty[page] = true
			}
		}
		var repeatedNext uint64
		for index := 0; index < length; index++ {
			var err error
			repeatedNext, err = runSharedEnrollmentByteFillStep(repeated)
			if err != nil {
				t.Fatalf("trial %d step %d: %v", trial, index, err)
			}
		}
		semanticNext, err := runSharedEnrollmentByteFill(semantic)
		if err != nil {
			t.Fatalf("trial %d semantic: %v", trial, err)
		}
		if repeatedNext != semanticNext || repeated.x != semantic.x ||
			repeated.n != semantic.n || repeated.z != semantic.z || repeated.c != semantic.c || repeated.v != semantic.v ||
			!reflect.DeepEqual(repeated.writable, semantic.writable) {
			t.Fatalf("trial %d state differs", trial)
		}
	}
}

func TestSharedEnrollmentFillByte(t *testing.T) {
	random := rand.New(rand.NewSource(0x51c16c))
	for trial := 0; trial < 100000; trial++ {
		stateValue, key := random.Uint32(), random.Uint32()
		want := byte((stateValue ^ key) * (stateValue + 0x11))
		if got := sharedEnrollmentFillByte(stateValue, key); got != want {
			t.Fatalf("trial %d: got %#x want %#x", trial, got, want)
		}
	}
}
