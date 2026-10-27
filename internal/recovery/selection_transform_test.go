package recovery

import (
	"math/rand"
	"reflect"
	"testing"
)

func TestSelectionPackedWordIterationMatchesLiftedBlocks(t *testing.T) {
	random := rand.New(rand.NewSource(int64(selectionPackedWordSeed)))
	for trial := 0; trial < 256; trial++ {
		pages := map[uint64][]byte{
			0x500000: make([]byte, pageSize),
			0x600000: make([]byte, pageSize),
			0x700000: make([]byte, pageSize),
		}
		const inputBase = uint64(0x100000000)
		for _, index := range []uint32{0x9b4e23ca, 0x9b4e23cb, 0x9b4e23cc, 0x9b4e23cd} {
			address := inputBase + uint64(index)
			if pages[address&^(pageSize-1)] == nil {
				pages[address&^(pageSize-1)] = make([]byte, pageSize)
			}
		}
		for _, page := range pages {
			if _, err := random.Read(page); err != nil {
				t.Fatal(err)
			}
		}
		lifted := newState(&model{pages: clonePages(pages)})
		direct := newState(&model{pages: clonePages(pages)})
		for register := range lifted.x {
			lifted.x[register] = random.Uint64()
		}
		lifted.sp = 0x700000
		lifted.x[3] = 0
		lifted.x[6] = 0
		lifted.x[8] = inputBase
		lifted.x[15] = 0x9b4e23cd
		lifted.x[16] = 0
		lifted.x[28] = 0
		lifted.x[30] = 0
		for register := 19; register <= 27; register++ {
			lifted.x[register] = 0
		}
		if err := lifted.storeValue(lifted.sp+0xd98, 0x500100, 8); err != nil {
			t.Fatal(err)
		}
		if err := lifted.storeValue(lifted.sp+0xdc0, 0, 4); err != nil {
			t.Fatal(err)
		}
		if err := lifted.storeValue(lifted.sp+0xdf0, 0x600800, 8); err != nil {
			t.Fatal(err)
		}
		direct.sp = lifted.sp
		direct.x = lifted.x
		direct.n, direct.z, direct.c, direct.v = lifted.n, lifted.z, lifted.c, lifted.v
		direct.writable = clonePages(lifted.writable)
		direct.dirty = map[uint64]bool{}
		for page := range direct.writable {
			direct.dirty[page] = true
		}

		seed := compact26A428SelectionBlocks[selectionPackedWordSeed]
		next, stopped, err := runCompactBlock(lifted, nil, seed)
		if err != nil || stopped || next != selectionPackedWordStep {
			t.Fatalf("trial %d lifted seed: next=%#x stopped=%v err=%v", trial, next, stopped, err)
		}
		step := compact26A428SelectionBlocks[selectionPackedWordStep]
		nextLifted, stopped, err := runCompactBlock(lifted, nil, step)
		if err != nil || stopped {
			t.Fatalf("trial %d lifted step: next=%#x stopped=%v err=%v", trial, nextLifted, stopped, err)
		}
		nextDirect, err := runSelectionPackedWordIteration(direct)
		if err != nil {
			t.Fatalf("trial %d direct iteration: %v", trial, err)
		}
		if nextLifted != nextDirect || lifted.x != direct.x {
			t.Fatalf("trial %d state differs: next lifted=%#x direct=%#x\nlifted=%#v\ndirect=%#v", trial, nextLifted, nextDirect, lifted.x, direct.x)
		}
		if lifted.n != direct.n || lifted.z != direct.z || lifted.c != direct.c || lifted.v != direct.v {
			t.Fatalf("trial %d flags differ", trial)
		}
		if !reflect.DeepEqual(lifted.writable, direct.writable) {
			liftedWord, _ := lifted.load(0x500100, 4)
			directWord, _ := direct.load(0x500100, 4)
			t.Fatalf("trial %d writes differ: lifted=%#x direct=%#x", trial, liftedWord, directWord)
		}
	}
}

func TestSelectionPackedWordPassMatchesIterations(t *testing.T) {
	random := rand.New(rand.NewSource(int64(selectionPackedWordSeed) ^ 0x400))
	const input, output, stack, dispatch = uint64(0x100000), uint64(0x200000), uint64(0x300000), uint64(0x400000)
	const exitTarget = uint64(0x123456789)
	for trial := 0; trial < 512; trial++ {
		length := 1 + random.Intn(16)
		selector := uint32(0)
		pages := map[uint64][]byte{
			input:    make([]byte, pageSize),
			output:   make([]byte, pageSize),
			stack:    make([]byte, pageSize),
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
		base.sp = stack
		base.x[0] = 0
		base.x[3] = 0
		base.x[6] = 0
		base.x[7] = 0
		inputBase := uint64(input)
		base.x[8] = inputBase - uint64(0x9b4e23ca)
		base.x[10] = 0
		base.x[15] = 0x9b4e23cd
		base.x[16] = 0
		base.x[19] = 0
		base.x[20] = 0
		base.x[21] = 0
		base.x[22] = 0
		base.x[23] = 0
		base.x[24] = 0
		base.x[25] = 0
		base.x[26] = 0
		base.x[27] = 0
		base.x[28] = uint64(64 - 4*(length-1))
		base.x[30] = 1
		if err := base.storeValue(stack+0xd98, output, 8); err != nil {
			t.Fatal(err)
		}
		if err := base.storeValue(stack+0xdc0, uint64(selector), 4); err != nil {
			t.Fatal(err)
		}
		if err := base.storeValue(stack+0xdf0, dispatch, 8); err != nil {
			t.Fatal(err)
		}
		if err := base.storeValue(dispatch, exitTarget, 8); err != nil {
			t.Fatal(err)
		}
		if err := base.storeValue(dispatch+8, selectionPackedWordSeed, 8); err != nil {
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
			repeatedNext, err = runSelectionPackedWordIteration(repeated)
			if err != nil {
				t.Fatal(err)
			}
		}
		semanticNext, err := runSelectionPackedWordPass(semantic)
		if err != nil {
			t.Fatal(err)
		}
		if repeatedNext != semanticNext || repeated.x != semantic.x || repeated.n != semantic.n || repeated.z != semantic.z || repeated.c != semantic.c || repeated.v != semantic.v || !reflect.DeepEqual(repeated.writable, semantic.writable) {
			t.Fatalf("trial %d state differs", trial)
		}
	}
}
