package recovery

import (
	"encoding/binary"
	"math/rand"
	"reflect"
	"testing"
)

func TestSharedEnrollmentBytePermuteMatchesLiftedBlock(t *testing.T) {
	random := rand.New(rand.NewSource(int64(sharedEnrollmentBytePermute)))
	profile := embedded26A428ExecutionProfile()
	block, ok := profile.selectionBlocks[sharedEnrollmentBytePermute]
	if !ok {
		block, ok = profile.supplementalBlocks[sharedEnrollmentBytePermute]
	}
	if !ok {
		t.Fatal("shared byte-permute block is unavailable")
	}
	for trial := 0; trial < 128; trial++ {
		const stackBase, sourceBase, workBase, dispatchBase = uint64(0x700000), uint64(0x710000000000), uint64(0x720000000000), uint64(0x800000000)
		registers := [31]uint64{}
		for register := range registers {
			registers[register] = random.Uint64()
		}
		registers[0] = uint64(random.Intn(256))
		registers[6] = 0
		registers[21] = dispatchBase
		inputAddress := sourceBase + (registers[20] ^ registers[8])
		pages := map[uint64][]byte{
			stackBase:                      make([]byte, pageSize),
			workBase:                       make([]byte, pageSize),
			inputAddress &^ (pageSize - 1): make([]byte, pageSize),
		}
		dispatchAddress := uint64(int64(dispatchBase) + int64(int32(uint32(registers[1])))*8)
		pages[dispatchAddress&^(pageSize-1)] = make([]byte, pageSize)
		for _, page := range pages {
			if _, err := random.Read(page); err != nil {
				t.Fatal(err)
			}
		}
		binary.LittleEndian.PutUint64(pages[stackBase][0x50:], sourceBase)
		binary.LittleEndian.PutUint64(pages[stackBase][0x58:], workBase)
		binary.LittleEndian.PutUint64(pages[dispatchAddress&^(pageSize-1)][dispatchAddress&(pageSize-1):], 0x100000000)
		lifted, direct := newState(&model{pages: clonePages(pages)}), newState(&model{pages: clonePages(pages)})
		lifted.sp, lifted.x = stackBase, registers
		direct.sp, direct.x = lifted.sp, lifted.x
		next, stopped, err := runCompactBlock(lifted, nil, block)
		if err != nil || stopped {
			t.Fatalf("trial %d lifted: next=%#x stopped=%v err=%v", trial, next, stopped, err)
		}
		directNext, err := runSharedEnrollmentBytePermute(direct)
		if err != nil {
			t.Fatalf("trial %d direct: %v", trial, err)
		}
		if next != directNext || lifted.x != direct.x {
			t.Fatalf("trial %d scalar state differs: next lifted=%#x direct=%#x\nlifted=%#v\ndirect=%#v", trial, next, directNext, lifted.x, direct.x)
		}
		if lifted.n != direct.n || lifted.z != direct.z || lifted.c != direct.c || lifted.v != direct.v {
			t.Fatalf("trial %d flags differ", trial)
		}
		if !reflect.DeepEqual(lifted.writable, direct.writable) {
			t.Fatalf("trial %d memory writes differ", trial)
		}
	}
}
