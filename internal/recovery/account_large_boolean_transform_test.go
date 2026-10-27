package recovery

import (
	"encoding/binary"
	"math/rand"
	"reflect"
	"testing"
)

func TestAccountLargeBooleanTransformIsSHA1FeedForward(t *testing.T) {
	const (
		frameBase    = uint64(0x700000200000)
		dispatchBase = uint64(0x710000000000)
	)
	random := rand.New(rand.NewSource(int64(accountLargeBooleanTransform)))
	for trial := 0; trial < 1024; trial++ {
		registers := [31]uint64{}
		for index := range registers {
			registers[index] = random.Uint64()
		}
		registers[21], registers[29] = dispatchBase, frameBase+0x100
		registers[24] = 0x892c
		pages := map[uint64][]byte{frameBase: make([]byte, pageSize)}
		if _, err := random.Read(pages[frameBase]); err != nil {
			t.Fatal(err)
		}
		dispatchAddress := uint64(int64(dispatchBase) + int64(int32(registers[24]))*8)
		dispatchPage := dispatchAddress &^ (pageSize - 1)
		if pages[dispatchPage] == nil {
			pages[dispatchPage] = make([]byte, pageSize)
		}
		binary.LittleEndian.PutUint64(pages[dispatchPage][dispatchAddress-dispatchPage:], 0x100000000+uint64(trial))

		circuit := newState(&model{pages: clonePages(pages)})
		semantic := newState(&model{pages: clonePages(pages)})
		circuit.x, semantic.x = registers, registers
		circuitNext, circuitErr := runAccountLargeBooleanTransformCircuit(circuit)
		semanticNext, semanticErr := runAccountLargeBooleanTransform(semantic)
		if circuitErr != nil || semanticErr != nil {
			t.Fatalf("trial %d errors: circuit=%v semantic=%v", trial, circuitErr, semanticErr)
		}
		if circuitNext != semanticNext {
			t.Fatalf("trial %d next=%#x, want %#x", trial, semanticNext, circuitNext)
		}
		for _, register := range [...]int{0, 3, 7, 9, 11, 12, 13, 15, 17, 19, 22, 24, 26} {
			if semantic.x[register] != circuit.x[register] {
				t.Fatalf("trial %d x%d=%#x, want %#x", trial, register, semantic.x[register], circuit.x[register])
			}
		}
		if circuit.n != semantic.n || circuit.z != semantic.z || circuit.c != semantic.c || circuit.v != semantic.v {
			t.Fatalf("trial %d flags differ", trial)
		}
		if !reflect.DeepEqual(circuit.writable, semantic.writable) {
			t.Fatalf("trial %d memory writes differ", trial)
		}
	}
}
