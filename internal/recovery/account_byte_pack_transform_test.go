package recovery

import (
	"encoding/binary"
	"math/rand"
	"reflect"
	"testing"
)

func TestAccountBytePackStepMatchesLiftedBlock(t *testing.T) {
	random := rand.New(rand.NewSource(int64(accountBytePackPass)))
	profile := embedded26A428ExecutionProfile()
	for trial := 0; trial < 32; trial++ {
		const inputBase, outputBase, dispatchBase = uint64(0x700000000000), uint64(0x710000000000), uint64(0x800000000)
		registers := [31]uint64{}
		for register := range registers {
			registers[register] = random.Uint64()
		}
		registers[1], registers[4], registers[10], registers[21] = inputBase, 0x2ba7c387, outputBase, dispatchBase
		addresses := []uint64{
			registers[1] + registers[4] + registers[14],
			registers[1] + (uint64(uint32(registers[9])+0x2ba76699) ^ 0xffffffffffffffad) + registers[4],
			registers[1] + registers[4] + registers[28],
			registers[1] + 0xffffffffd4583c78 + registers[4],
		}
		pages := make(map[uint64][]byte)
		for _, address := range addresses {
			pages[address&^(pageSize-1)] = make([]byte, pageSize)
		}
		for _, page := range pages {
			if _, err := random.Read(page); err != nil {
				t.Fatal(err)
			}
		}

		scout := newState(&model{pages: clonePages(pages)})
		scout.x = registers
		_, _ = runAccountBytePackStep(scout)
		outputAddress := outputBase + uint64(uint32(scout.x[12]))*4
		if pages[outputAddress&^(pageSize-1)] == nil {
			pages[outputAddress&^(pageSize-1)] = make([]byte, pageSize)
		}
		for _, index := range []uint32{uint32(registers[9]), uint32(registers[9]) ^ 0xc285} {
			address := uint64(int64(dispatchBase) + int64(int32(index))*8)
			if pages[address&^(pageSize-1)] == nil {
				pages[address&^(pageSize-1)] = make([]byte, pageSize)
			}
		}
		for ordinal, index := range []uint32{uint32(registers[9]), uint32(registers[9]) ^ 0xc285} {
			address := uint64(int64(dispatchBase) + int64(int32(index))*8)
			binary.LittleEndian.PutUint64(pages[address&^(pageSize-1)][address&(pageSize-1):], uint64(0x100000000+ordinal))
		}
		lifted, direct := newState(&model{pages: clonePages(pages)}), newState(&model{pages: clonePages(pages)})
		lifted.x, direct.x = registers, registers
		next, stopped, err := runCompactBlock(lifted, nil, profile.accountEntryBlocks[accountBytePackPass])
		if err != nil || stopped {
			t.Fatalf("trial %d lifted: next=%#x stopped=%v err=%v", trial, next, stopped, err)
		}
		directNext, err := runAccountBytePackStep(direct)
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
