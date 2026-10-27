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
