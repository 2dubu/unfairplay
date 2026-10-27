package recovery

import (
	"encoding/binary"
	"math/rand"
	"reflect"
	"testing"
)

func TestAccountSecondFoldMatchesLiftedLoop(t *testing.T) {
	random := rand.New(rand.NewSource(int64(accountSecondFoldRound)))
	profile := embedded26A428ExecutionProfile()
	for trial := 0; trial < 32; trial++ {
		const tableBase, dispatchBase = uint64(0x1000000), uint64(0x800000000)
		pages := map[uint64][]byte{tableBase: make([]byte, pageSize)}
		registers := [31]uint64{}
		for register := range registers {
			registers[register] = random.Uint64()
		}
		registers[8], registers[20] = tableBase, dispatchBase
		registers[21] = accountSecondFoldLimit - 1
		registers[14] = 0 - (registers[21] << 2)
		for _, index := range []uint32{uint32(registers[0]), uint32(registers[0]) ^ 0x1f} {
			address := uint64(int64(dispatchBase) + int64(int32(index))*8)
			pages[address&^(pageSize-1)] = make([]byte, pageSize)
		}
		for _, page := range pages {
			if _, err := random.Read(page); err != nil {
				t.Fatal(err)
			}
		}
		for ordinal, index := range []uint32{uint32(registers[0]), uint32(registers[0]) ^ 0x1f} {
			address := uint64(int64(dispatchBase) + int64(int32(index))*8)
			value := uint64(0x100000000)
			if ordinal != 0 {
				value = accountSecondFoldRound
			}
			binary.LittleEndian.PutUint64(pages[address&^(pageSize-1)][address&(pageSize-1):], value)
		}
		lifted, direct := newState(&model{pages: clonePages(pages)}), newState(&model{pages: clonePages(pages)})
		lifted.x, direct.x = registers, registers
		next, stopped, err := runCompactBlock(lifted, nil, profile.accountEntryBlocks[accountSecondFoldRound])
		if err != nil || stopped || next != accountSecondFoldTail {
			t.Fatalf("trial %d lifted front: next=%#x stopped=%v err=%v", trial, next, stopped, err)
		}
		next, stopped, err = runCompactBlock(lifted, nil, profile.accountEntryBlocks[accountSecondFoldTail])
		if err != nil || stopped {
			t.Fatalf("trial %d lifted tail: next=%#x stopped=%v err=%v", trial, next, stopped, err)
		}
		directNext, err := runAccountSecondFoldCircuit(direct)
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
