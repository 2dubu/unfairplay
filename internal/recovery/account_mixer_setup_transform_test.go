package recovery

import (
	"encoding/binary"
	"math/rand"
	"reflect"
	"testing"
)

func TestAccountMixerSetupMatchesCircuit(t *testing.T) {
	const (
		stackBase    = uint64(0x700000100000)
		frameBase    = uint64(0x700000200000)
		dispatchBase = uint64(0x710000000000)
	)
	random := rand.New(rand.NewSource(int64(accountMixerSetup)))
	for trial := 0; trial < 1024; trial++ {
		registers := [31]uint64{}
		for index := range registers {
			registers[index] = random.Uint64()
		}
		registers[21], registers[29] = dispatchBase, frameBase+0x100
		registers[7] = 0x82835ef1
		pages := map[uint64][]byte{
			stackBase: make([]byte, pageSize),
			frameBase: make([]byte, pageSize),
		}
		for _, page := range pages {
			if _, err := random.Read(page); err != nil {
				t.Fatal(err)
			}
		}
		dispatchIndex := uint32(registers[7]) ^ 0x82832578
		if uint32(registers[7]) == 0x62e0114c {
			dispatchIndex ^= 0x8ca4
		}
		dispatchAddress := uint64(int64(dispatchBase) + int64(int32(dispatchIndex))*8)
		dispatchPage := dispatchAddress &^ (pageSize - 1)
		if pages[dispatchPage] == nil {
			pages[dispatchPage] = make([]byte, pageSize)
		}
		binary.LittleEndian.PutUint64(pages[dispatchPage][dispatchAddress-dispatchPage:], 0x100000000+uint64(trial))

		circuit := newState(&model{pages: clonePages(pages)})
		semantic := newState(&model{pages: clonePages(pages)})
		circuit.x, semantic.x = registers, registers
		circuit.sp, semantic.sp = stackBase, stackBase
		circuitNext, circuitErr := runAccountMixerSetupCircuit(circuit)
		semanticNext, semanticErr := runAccountMixerSetup(semantic)
		if circuitErr != nil || semanticErr != nil {
			t.Fatalf("trial %d errors: circuit=%v semantic=%v", trial, circuitErr, semanticErr)
		}
		if circuitNext != semanticNext || circuit.x != semantic.x {
			t.Fatalf("trial %d scalar state differs: next circuit=%#x semantic=%#x\ncircuit=%#v\nsemantic=%#v",
				trial, circuitNext, semanticNext, circuit.x, semantic.x)
		}
		if circuit.n != semantic.n || circuit.z != semantic.z || circuit.c != semantic.c || circuit.v != semantic.v {
			t.Fatalf("trial %d flags differ", trial)
		}
		if !reflect.DeepEqual(circuit.writable, semantic.writable) {
			t.Fatalf("trial %d memory writes differ: circuit=%08x semantic=%08x", trial,
				binary.LittleEndian.Uint32(circuit.writable[stackBase][0x7c:]),
				binary.LittleEndian.Uint32(semantic.writable[stackBase][0x7c:]))
		}
	}
}
