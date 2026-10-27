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
	for trial := 0; trial < 256; trial++ {
		const inputBase, outputBase, dispatchBase = uint64(0x700000000000), uint64(0x710000000000), uint64(0x800000000)
		registers := [31]uint64{}
		for register := range registers {
			registers[register] = random.Uint64()
		}
		wordIndex := uint32(trial) % accountSHA1InputWords
		registers[1], registers[4], registers[10], registers[21] = inputBase, uint64(accountBytePackCounter+wordIndex*4), outputBase, dispatchBase
		registers[9] = 0x5d42
		registers[8] &= 0xff
		registers[14] = ^uint64(accountBytePackCounter) + 1
		registers[28] = registers[14] + 2
		pages := map[uint64][]byte{
			inputBase:  make([]byte, pageSize),
			outputBase: make([]byte, pageSize),
		}
		for _, page := range pages {
			if _, err := random.Read(page); err != nil {
				t.Fatal(err)
			}
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
		lifted := newState(&model{pages: clonePages(pages)})
		circuit := newState(&model{pages: clonePages(pages)})
		semantic := newState(&model{pages: clonePages(pages)})
		lifted.x, circuit.x, semantic.x = registers, registers, registers
		next, stopped, err := runCompactBlock(lifted, nil, profile.accountEntryBlocks[accountBytePackPass])
		if err != nil || stopped {
			t.Fatalf("trial %d lifted: next=%#x stopped=%v err=%v", trial, next, stopped, err)
		}
		circuitNext, err := runAccountBytePackStepCircuit(circuit)
		if err != nil {
			t.Fatalf("trial %d circuit: %v", trial, err)
		}
		semanticNext, err := runAccountBytePackStep(semantic)
		if err != nil {
			t.Fatalf("trial %d semantic: %v", trial, err)
		}
		if next != circuitNext || next != semanticNext || lifted.x != circuit.x || lifted.x != semantic.x {
			t.Fatalf("trial %d scalar state differs: next lifted=%#x circuit=%#x semantic=%#x\nlifted=%#v\ncircuit=%#v\nsemantic=%#v",
				trial, next, circuitNext, semanticNext, lifted.x, circuit.x, semantic.x)
		}
		if lifted.n != circuit.n || lifted.z != circuit.z || lifted.c != circuit.c || lifted.v != circuit.v ||
			lifted.n != semantic.n || lifted.z != semantic.z || lifted.c != semantic.c || lifted.v != semantic.v {
			t.Fatalf("trial %d flags differ", trial)
		}
		if !reflect.DeepEqual(lifted.writable, circuit.writable) || !reflect.DeepEqual(lifted.writable, semantic.writable) {
			t.Fatalf("trial %d memory writes differ: lifted=%08x circuit=%08x semantic=%08x", trial,
				binary.LittleEndian.Uint32(lifted.writable[outputBase][:4]),
				binary.LittleEndian.Uint32(circuit.writable[outputBase][:4]),
				binary.LittleEndian.Uint32(semantic.writable[outputBase][:4]))
		}
	}
}
