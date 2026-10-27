package recovery

import (
	"encoding/binary"
	"math/rand"
	"reflect"
	"testing"
)

func TestAccountSHA1ExpandWord(t *testing.T) {
	// First expanded word of the single padded SHA-1 block for "abc".
	previous := [4]uint32{0, 0, 0, 0x61626380}
	if got, want := accountSHA1ExpandWord(previous), uint32(0xc2c4c700); got != want {
		t.Fatalf("W[16] = %#08x, want %#08x", got, want)
	}
}

func TestAccountWordMixerMatchesCircuit(t *testing.T) {
	const (
		scheduleBase = uint64(0x710000100000)
		dispatchBase = uint64(0x710000200000)
		stackBase    = uint64(0x710000300000)
		counterEnd   = uint32(0x2ba7c3db)
		counterBase  = counterEnd - accountSHA1Rounds
		encodedBase  = uint64(0xffffffff5160f1d4)
	)
	random := rand.New(rand.NewSource(int64(accountEntryWordMixerPass)))
	for trial := 0; trial < 512; trial++ {
		pages := map[uint64][]byte{
			scheduleBase: make([]byte, pageSize),
			stackBase:    make([]byte, pageSize),
		}
		if _, err := random.Read(pages[scheduleBase]); err != nil {
			t.Fatal(err)
		}
		if _, err := random.Read(pages[stackBase]); err != nil {
			t.Fatal(err)
		}
		for _, index := range []uint32{0x5d30, 0x5d42} {
			address := dispatchBase + uint64(index)*8
			base := address &^ (pageSize - 1)
			if pages[base] == nil {
				pages[base] = make([]byte, pageSize)
			}
			binary.LittleEndian.PutUint64(pages[base][address-base:], 0x100000000+uint64(index))
		}

		registers := [31]uint64{}
		for index := range registers {
			registers[index] = random.Uint64()
		}
		word := uint32(16 + trial%64)
		registers[6], registers[10], registers[15] = encodedBase, scheduleBase, uint64(counterEnd)
		registers[20], registers[21] = uint64(counterBase+word), dispatchBase

		circuit := newState(&model{pages: clonePages(pages)})
		semantic := newState(&model{pages: clonePages(pages)})
		circuit.x, semantic.x = registers, registers
		circuit.sp, semantic.sp = stackBase, stackBase
		circuitNext, circuitErr := runAccountEntryWordMixerStepCircuit(circuit)
		semanticNext, semanticErr := runAccountEntryWordMixerStep(semantic)
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
			t.Fatalf("trial %d memory writes differ", trial)
		}
	}
}
