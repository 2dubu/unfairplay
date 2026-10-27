package recovery

import (
	"math/rand"
	"reflect"
	"testing"
)

func TestAccountWordWriterMatchesSteps(t *testing.T) {
	random := rand.New(rand.NewSource(int64(accountWordWriter)))
	const input, output, dispatch = uint64(0x200000), uint64(0x300000), uint64(0x400000)
	const exitTarget = uint64(0x123456789)
	const firstCounter = uint64(0x1cc84893)
	const selector = uint32(0xb375)
	for firstWord := 0; firstWord < 16; firstWord++ {
		for trial := 0; trial < 64; trial++ {
			length := 16 - firstWord
			dispatchBase := addOffset(dispatch, -int64(int32(selector))*8)
			exitAddress := addOffset(dispatchBase, int64(int32(selector^0xc1))*8)
			pages := map[uint64][]byte{
				input:                          make([]byte, pageSize),
				output:                         make([]byte, pageSize),
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
			negativeCounter := ^firstCounter + 1
			base.x[0] = firstCounter + uint64(firstWord)*4
			base.x[3] = negativeCounter
			base.x[4] = 0x1b
			base.x[5] = negativeCounter + 1
			base.x[8] = output
			base.x[9] = 0xaf
			base.x[12] = 0xce757bc5
			base.x[13] = 0xb3b4
			base.x[16] = uint64(selector)
			base.x[20] = dispatchBase
			base.x[21] = input
			base.x[23] = 0x4e757bd5
			base.x[24] = negativeCounter + 2
			base.x[25] = negativeCounter + 3
			base.x[26] = 0x11e99fcc
			base.x[27] = 0xf2de28a2
			if err := base.storeValue(dispatch, accountWordWriter, 8); err != nil {
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
				repeatedNext, err = runAccountWordWriterStep(repeated)
				if err != nil {
					t.Fatal(err)
				}
			}
			semanticNext, err := runAccountWordWriter(semantic)
			if err != nil {
				t.Fatal(err)
			}
			if repeatedNext != semanticNext || repeated.x != semantic.x || repeated.n != semantic.n || repeated.z != semantic.z || repeated.c != semantic.c || repeated.v != semantic.v || !reflect.DeepEqual(repeated.writable, semantic.writable) {
				t.Fatalf("word %d trial %d state differs", firstWord, trial)
			}
		}
	}
}
