package recovery

import (
	"encoding/binary"
	"math/rand"
	"reflect"
	"testing"
)

func TestBootstrapByteStreamMatchesLiftedLoop(t *testing.T) {
	random := rand.New(rand.NewSource(int64(bootstrapByteStreamPass)))
	profile := embedded26A428ExecutionProfile()
	for trial := 0; trial < 64; trial++ {
		const (
			stackBase    = uint64(0x700000)
			outputBase   = uint64(0x800000)
			dispatchBase = uint64(0x800000000)
		)
		pages := map[uint64][]byte{stackBase: make([]byte, pageSize), outputBase: make([]byte, pageSize)}
		registers := [31]uint64{}
		for register := range registers {
			registers[register] = random.Uint64()
		}
		registers[1] = uint64(1 + random.Intn(256))
		registers[9] = outputBase
		registers[30] = dispatchBase
		for _, index := range []uint32{uint32(registers[8]), uint32(registers[8]) ^ uint32(registers[4])} {
			address := uint64(int64(dispatchBase) + int64(int32(index))*8)
			pages[address&^(pageSize-1)] = make([]byte, pageSize)
		}
		for _, page := range pages {
			if _, err := random.Read(page); err != nil {
				t.Fatal(err)
			}
		}
		for ordinal, index := range []uint32{uint32(registers[8]), uint32(registers[8]) ^ uint32(registers[4])} {
			address := uint64(int64(dispatchBase) + int64(int32(index))*8)
			value := bootstrapByteStreamPass
			if ordinal != 0 {
				value = 0x100000000
			}
			binary.LittleEndian.PutUint64(pages[address&^(pageSize-1)][address&(pageSize-1):], value)
		}
		lifted := newState(&model{pages: clonePages(pages)})
		direct := newState(&model{pages: clonePages(pages)})
		lifted.sp, lifted.x = stackBase, registers
		direct.sp, direct.x = lifted.sp, lifted.x
		direct.n, direct.z, direct.c, direct.v = lifted.n, lifted.z, lifted.c, lifted.v

		next := bootstrapByteStreamPass
		var stopped bool
		var err error
		for next == bootstrapByteStreamPass {
			next, stopped, err = runCompactBlock(lifted, nil, profile.bootstrapBlocks[bootstrapByteStreamPass])
			if err != nil || stopped {
				t.Fatalf("trial %d lifted: next=%#x stopped=%v err=%v", trial, next, stopped, err)
			}
		}
		directNext, err := runBootstrapByteStream(direct)
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
