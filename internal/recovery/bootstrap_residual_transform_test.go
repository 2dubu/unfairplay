package recovery

import (
	"encoding/binary"
	"math/rand"
	"reflect"
	"testing"
)

func TestBootstrapScalarRoundMatchesLiftedBlocks(t *testing.T) {
	random := rand.New(rand.NewSource(int64(bootstrapScalarRoundStart)))
	profile := embedded26A428ExecutionProfile()
	for trial := 0; trial < 32; trial++ {
		const (
			lookupBase   = uint64(0x1000000)
			stackBase    = uint64(0x700000)
			dispatchBase = uint64(0x800000000)
		)
		pages := make(map[uint64][]byte)
		for offset := uint64(0); offset < 0x40000; offset += pageSize {
			pages[lookupBase+offset] = make([]byte, pageSize)
		}
		pages[stackBase] = make([]byte, pageSize)
		registers := [31]uint64{}
		for register := range registers {
			registers[register] = random.Uint64()
		}
		registers[9] = lookupBase
		for _, index := range []uint32{uint32(registers[23]), uint32(registers[23]) ^ 0x16} {
			address := uint64(int64(dispatchBase) + int64(int32(index))*8)
			page := address &^ (pageSize - 1)
			if pages[page] == nil {
				pages[page] = make([]byte, pageSize)
			}
		}
		for _, page := range pages {
			if _, err := random.Read(page); err != nil {
				t.Fatal(err)
			}
		}
		binary.LittleEndian.PutUint64(pages[stackBase][0x208:], dispatchBase)
		for ordinal, index := range []uint32{uint32(registers[23]), uint32(registers[23]) ^ 0x16} {
			address := uint64(int64(dispatchBase) + int64(int32(index))*8)
			binary.LittleEndian.PutUint64(pages[address&^(pageSize-1)][address&(pageSize-1):], uint64(0x100000000+ordinal))
		}

		lifted := newState(&model{pages: clonePages(pages)})
		direct := newState(&model{pages: clonePages(pages)})
		lifted.sp = stackBase
		lifted.x = registers
		direct.sp = lifted.sp
		direct.x = lifted.x
		direct.n, direct.z, direct.c, direct.v = lifted.n, lifted.z, lifted.c, lifted.v

		next, stopped, err := runCompactBlock(lifted, nil, profile.bootstrapBlocks[bootstrapScalarRoundStart])
		if err != nil || stopped || next != bootstrapScalarRoundTail {
			t.Fatalf("trial %d lifted start: next=%#x stopped=%v err=%v", trial, next, stopped, err)
		}
		roundDone, err := runBootstrapScalarRoundFront(direct)
		if err != nil {
			t.Fatalf("trial %d direct front: %v", trial, err)
		}
		if lifted.x != direct.x {
			t.Fatalf("trial %d front scalar state differs\nlifted=%#v\ndirect=%#v", trial, lifted.x, direct.x)
		}
		if lifted.n != direct.n || lifted.z != direct.z || lifted.c != direct.c || lifted.v != direct.v {
			t.Fatalf("trial %d front flags differ", trial)
		}
		next, stopped, err = runCompactBlock(lifted, nil, profile.bootstrapBlocks[bootstrapScalarRoundTail])
		if err != nil || stopped {
			t.Fatalf("trial %d lifted tail: next=%#x stopped=%v err=%v", trial, next, stopped, err)
		}
		directNext, err := runBootstrapScalarRoundFinish(direct, roundDone)
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
