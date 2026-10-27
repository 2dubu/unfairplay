package recovery

import (
	"encoding/binary"
	"math/rand"
	"reflect"
	"testing"
)

func TestAccountFoldMatchesLiftedLoop(t *testing.T) {
	random := rand.New(rand.NewSource(int64(accountFoldRound)))
	profile := embedded26A428ExecutionProfile()
	for trial := 0; trial < 32; trial++ {
		const tableBase, dispatchBase = uint64(0x1000000), uint64(0x800000000)
		pages := map[uint64][]byte{tableBase: make([]byte, pageSize)}
		registers := [31]uint64{}
		for register := range registers {
			registers[register] = random.Uint64()
		}
		rounds := uint64(1)
		registers[7] = accountFoldLimit - rounds
		registers[8] = tableBase
		registers[14] = 0 - (registers[7] << 2)
		registers[20] = dispatchBase
		for _, index := range []uint32{uint32(registers[0]), uint32(registers[0]) ^ 0x83c3} {
			address := uint64(int64(dispatchBase) + int64(int32(index))*8)
			pages[address&^(pageSize-1)] = make([]byte, pageSize)
		}
		for _, page := range pages {
			if _, err := random.Read(page); err != nil {
				t.Fatal(err)
			}
		}
		for ordinal, index := range []uint32{uint32(registers[0]), uint32(registers[0]) ^ 0x83c3} {
			address := uint64(int64(dispatchBase) + int64(int32(index))*8)
			value := accountFoldRound
			if ordinal != 0 {
				value = 0x100000000
			}
			binary.LittleEndian.PutUint64(pages[address&^(pageSize-1)][address&(pageSize-1):], value)
		}
		lifted, direct := newState(&model{pages: clonePages(pages)}), newState(&model{pages: clonePages(pages)})
		lifted.x, direct.x = registers, registers
		next := accountFoldRound
		var stopped bool
		var err error
		for next == accountFoldRound {
			next, stopped, err = runCompactBlock(lifted, nil, profile.accountEntryBlocks[accountFoldRound])
			if err != nil || stopped || next != accountFoldTail {
				t.Fatalf("trial %d lifted front: next=%#x stopped=%v err=%v", trial, next, stopped, err)
			}
			next, stopped, err = runCompactBlock(lifted, nil, profile.accountEntryBlocks[accountFoldTail])
			if err != nil || stopped {
				t.Fatalf("trial %d lifted tail: next=%#x stopped=%v err=%v", trial, next, stopped, err)
			}
		}
		directNext, err := runAccountFoldCircuit(direct)
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
