package recovery

import (
	"encoding/binary"
	"math/rand"
	"reflect"
	"testing"
)

func TestBootstrapWordTableMatchesLiftedLoop(t *testing.T) {
	random := rand.New(rand.NewSource(int64(bootstrapWordTableRound)))
	profile := embedded26A428ExecutionProfile()
	for trial := 0; trial < 16; trial++ {
		const tableBase, dispatchBase = uint64(0x1000000), uint64(0x800000000)
		pages := make(map[uint64][]byte)
		for offset := uint64(0); offset < 0x40000; offset += pageSize {
			pages[tableBase+offset] = make([]byte, pageSize)
		}
		registers := [31]uint64{}
		for register := range registers {
			registers[register] = random.Uint64()
		}
		registers[9], registers[16], registers[30] = tableBase, uint64(random.Intn(64)), dispatchBase
		registers[19], registers[20], registers[21] = 0, 0, 0
		for _, index := range []uint32{uint32(registers[13]), uint32(registers[13]) ^ 0xfdd8} {
			address := uint64(int64(dispatchBase) + int64(int32(index))*8)
			pages[address&^(pageSize-1)] = make([]byte, pageSize)
		}
		for _, page := range pages {
			if _, err := random.Read(page); err != nil {
				t.Fatal(err)
			}
		}
		for ordinal, index := range []uint32{uint32(registers[13]), uint32(registers[13]) ^ 0xfdd8} {
			address := uint64(int64(dispatchBase) + int64(int32(index))*8)
			value := bootstrapWordTableRound
			if ordinal != 0 {
				value = 0x100000000
			}
			binary.LittleEndian.PutUint64(pages[address&^(pageSize-1)][address&(pageSize-1):], value)
		}
		lifted, direct := newState(&model{pages: clonePages(pages)}), newState(&model{pages: clonePages(pages)})
		lifted.x, direct.x = registers, registers
		next := bootstrapWordTableRound
		var stopped bool
		var err error
		for next == bootstrapWordTableRound {
			next, stopped, err = runCompactBlock(lifted, nil, profile.bootstrapBlocks[bootstrapWordTableRound])
			if err != nil || stopped || next != bootstrapWordTableTail {
				t.Fatalf("trial %d lifted front: next=%#x stopped=%v err=%v", trial, next, stopped, err)
			}
			next, stopped, err = runCompactBlock(lifted, nil, profile.bootstrapBlocks[bootstrapWordTableTail])
			if err != nil || stopped {
				t.Fatalf("trial %d lifted tail: next=%#x stopped=%v err=%v", trial, next, stopped, err)
			}
		}
		directNext, err := runBootstrapWordTable(direct)
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
