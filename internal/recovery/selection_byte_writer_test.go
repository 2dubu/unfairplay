package recovery

import (
	"math/rand"
	"reflect"
	"testing"
)

func TestSelectionByteWriterPassMatchesSteps(t *testing.T) {
	random := rand.New(rand.NewSource(int64(selectionByteWriter)))
	const input, output, stack, dispatch = uint64(0x100000), uint64(0x200000), uint64(0x300100), uint64(0x400000)
	const exitTarget = uint64(0x123456789)
	for trial := 0; trial < 1024; trial++ {
		length := 1 + random.Intn(128)
		loopSelector := random.Uint32()
		exitSelector := random.Uint32()
		continueIndex := loopSelector ^ exitSelector
		continueAddress := addOffset(dispatch, int64(int32(continueIndex))*8)
		exitAddress := addOffset(dispatch, int64(int32(exitSelector))*8)
		pages := map[uint64][]byte{input: make([]byte, pageSize), output: make([]byte, pageSize), stack &^ (pageSize - 1): make([]byte, pageSize), continueAddress &^ (pageSize - 1): make([]byte, pageSize), exitAddress &^ (pageSize - 1): make([]byte, pageSize)}
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
		base.x[9] = 0
		base.x[11] = uint64(length - 1)
		base.x[12] = uint64(exitSelector)
		base.x[13] = 0x36
		base.x[14] = 0x6ff
		base.x[15] = uint64(loopSelector)
		base.x[16] = output
		base.x[20] = input
		if err := base.storeValue(stack+0xdf0, dispatch, 8); err != nil {
			t.Fatal(err)
		}
		if err := base.storeValue(continueAddress, selectionByteWriter, 8); err != nil {
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
			repeatedNext, err = runSelectionByteWriter(repeated)
			if err != nil {
				t.Fatal(err)
			}
		}
		semanticNext, err := runSelectionByteWriterPass(semantic)
		if err != nil {
			t.Fatal(err)
		}
		if repeatedNext != semanticNext || repeated.x != semantic.x || repeated.n != semantic.n || repeated.z != semantic.z || repeated.c != semantic.c || repeated.v != semantic.v || !reflect.DeepEqual(repeated.writable, semantic.writable) {
			t.Fatalf("trial %d state differs", trial)
		}
	}
}

func TestSelectionByteWriterEncodingIdentity(t *testing.T) {
	for input := uint32(0); input < 256; input++ {
		mixed := uint32(0x6ff) ^ (uint32(0x36) & (input << 1))
		old := byte(input + mixed + 0x1c)
		if got := byte(input) ^ 0x1b; got != old {
			t.Fatalf("input %#x encoded to %#x, want %#x", input, got, old)
		}
	}
}
