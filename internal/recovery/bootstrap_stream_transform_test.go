package recovery

import (
	"encoding/binary"
	"math/rand"
	"reflect"
	"testing"
)

func TestBootstrapStreamMixMatchesLiftedBlocks(t *testing.T) {
	random := rand.New(rand.NewSource(int64(bootstrapStreamMixStart)))
	profile := embedded26A428ExecutionProfile()
	for trial := 0; trial < 64; trial++ {
		const (
			stackBase    = uint64(0x700000)
			sourceBase   = uint64(0x800000)
			dispatchBase = uint64(0x800000000)
		)
		pages := map[uint64][]byte{
			stackBase:  make([]byte, pageSize),
			sourceBase: make([]byte, pageSize),
		}
		registers := [31]uint64{}
		for register := range registers {
			registers[register] = random.Uint64()
		}
		selector := random.Uint32()
		index := uint32(registers[22]) ^ selector
		dispatchAddress := uint64(int64(dispatchBase) + int64(int32(index))*8)
		pages[dispatchAddress&^(pageSize-1)] = make([]byte, pageSize)
		for _, page := range pages {
			if _, err := random.Read(page); err != nil {
				t.Fatal(err)
			}
		}
		binary.LittleEndian.PutUint64(pages[stackBase][0x70:], sourceBase)
		binary.LittleEndian.PutUint32(pages[stackBase][0x6c:], selector)
		binary.LittleEndian.PutUint64(pages[stackBase][0xa8:], dispatchBase)
		binary.LittleEndian.PutUint64(pages[dispatchAddress&^(pageSize-1)][dispatchAddress&(pageSize-1):], 0x100000000)

		lifted := newState(&model{pages: clonePages(pages)})
		direct := newState(&model{pages: clonePages(pages)})
		lifted.sp = stackBase
		lifted.x = registers
		direct.sp = lifted.sp
		direct.x = lifted.x
		direct.n, direct.z, direct.c, direct.v = lifted.n, lifted.z, lifted.c, lifted.v

		next, stopped, err := runCompactBlock(lifted, nil, profile.bootstrapBlocks[bootstrapStreamMixStart])
		if err != nil || stopped || next != bootstrapStreamMixTail {
			t.Fatalf("trial %d lifted start: next=%#x stopped=%v err=%v", trial, next, stopped, err)
		}
		next, stopped, err = runCompactBlock(lifted, nil, profile.bootstrapBlocks[bootstrapStreamMixTail])
		if err != nil || stopped {
			t.Fatalf("trial %d lifted tail: next=%#x stopped=%v err=%v", trial, next, stopped, err)
		}
		directNext, err := runBootstrapStreamMix(direct)
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

func TestBootstrapStreamStateMatchesLiftedBlock(t *testing.T) {
	random := rand.New(rand.NewSource(int64(bootstrapStreamState)))
	profile := embedded26A428ExecutionProfile()
	for trial := 0; trial < 64; trial++ {
		const (
			stackBase    = uint64(0x700000)
			sourceBase   = uint64(0x800000)
			dispatchBase = uint64(0x800000000)
		)
		pages := map[uint64][]byte{stackBase: make([]byte, pageSize), sourceBase: make([]byte, pageSize)}
		registers := [31]uint64{}
		for register := range registers {
			registers[register] = random.Uint64()
		}
		index := uint32(registers[11]) - 0x95d
		dispatchAddress := uint64(int64(dispatchBase) + int64(int32(index))*8)
		pages[dispatchAddress&^(pageSize-1)] = make([]byte, pageSize)
		for _, page := range pages {
			if _, err := random.Read(page); err != nil {
				t.Fatal(err)
			}
		}
		binary.LittleEndian.PutUint64(pages[stackBase][0xa8:], dispatchBase)
		binary.LittleEndian.PutUint64(pages[stackBase][0xa0:], sourceBase)
		binary.LittleEndian.PutUint64(pages[dispatchAddress&^(pageSize-1)][dispatchAddress&(pageSize-1):], 0x100000000)
		lifted, direct := bootstrapStreamStates(pages, registers, stackBase)

		next, stopped, err := runCompactBlock(lifted, nil, profile.bootstrapBlocks[bootstrapStreamState])
		if err != nil || stopped {
			t.Fatalf("trial %d lifted: next=%#x stopped=%v err=%v", trial, next, stopped, err)
		}
		directNext, err := runBootstrapStreamState(direct)
		if err != nil {
			t.Fatalf("trial %d direct: %v", trial, err)
		}
		assertBootstrapStreamState(t, trial, next, directNext, lifted, direct)
	}
}

func TestBootstrapStreamAdvanceMatchesLiftedBlock(t *testing.T) {
	random := rand.New(rand.NewSource(int64(bootstrapStreamAdvance)))
	profile := embedded26A428ExecutionProfile()
	for trial := 0; trial < 64; trial++ {
		const (
			stackBase    = uint64(0x700000)
			firstOutput  = uint64(0x800000)
			secondOutput = uint64(0x900000)
			dispatchBase = uint64(0x800000000)
		)
		pages := map[uint64][]byte{
			stackBase:    make([]byte, pageSize),
			firstOutput:  make([]byte, pageSize),
			secondOutput: make([]byte, pageSize),
		}
		registers := [31]uint64{}
		for register := range registers {
			registers[register] = random.Uint64()
		}
		state := random.Uint64()
		for _, index := range []uint32{uint32(state), uint32(state) ^ 0x17cf} {
			dispatchAddress := uint64(int64(dispatchBase) + int64(int32(index))*8)
			pages[dispatchAddress&^(pageSize-1)] = make([]byte, pageSize)
		}
		for _, page := range pages {
			if _, err := random.Read(page); err != nil {
				t.Fatal(err)
			}
		}
		binary.LittleEndian.PutUint64(pages[stackBase][0x98:], state)
		binary.LittleEndian.PutUint64(pages[stackBase][0x60:], firstOutput)
		binary.LittleEndian.PutUint64(pages[stackBase][0x58:], secondOutput)
		binary.LittleEndian.PutUint64(pages[stackBase][0xa8:], dispatchBase)
		for ordinal, index := range []uint32{uint32(state), uint32(state) ^ 0x17cf} {
			dispatchAddress := uint64(int64(dispatchBase) + int64(int32(index))*8)
			binary.LittleEndian.PutUint64(pages[dispatchAddress&^(pageSize-1)][dispatchAddress&(pageSize-1):], uint64(0x100000000+ordinal))
		}
		lifted, direct := bootstrapStreamStates(pages, registers, stackBase)

		next, stopped, err := runCompactBlock(lifted, nil, profile.bootstrapBlocks[bootstrapStreamAdvance])
		if err != nil || stopped {
			t.Fatalf("trial %d lifted: next=%#x stopped=%v err=%v", trial, next, stopped, err)
		}
		directNext, err := runBootstrapStreamAdvance(direct)
		if err != nil {
			t.Fatalf("trial %d direct: %v", trial, err)
		}
		assertBootstrapStreamState(t, trial, next, directNext, lifted, direct)
	}
}

func bootstrapStreamStates(pages map[uint64][]byte, registers [31]uint64, stack uint64) (*state, *state) {
	lifted := newState(&model{pages: clonePages(pages)})
	direct := newState(&model{pages: clonePages(pages)})
	lifted.sp = stack
	lifted.x = registers
	direct.sp = lifted.sp
	direct.x = lifted.x
	direct.n, direct.z, direct.c, direct.v = lifted.n, lifted.z, lifted.c, lifted.v
	return lifted, direct
}

func assertBootstrapStreamState(t *testing.T, trial int, liftedNext, directNext uint64, lifted, direct *state) {
	t.Helper()
	if liftedNext != directNext || lifted.x != direct.x {
		t.Fatalf("trial %d scalar state differs: next lifted=%#x direct=%#x\nlifted=%#v\ndirect=%#v", trial, liftedNext, directNext, lifted.x, direct.x)
	}
	if lifted.n != direct.n || lifted.z != direct.z || lifted.c != direct.c || lifted.v != direct.v {
		t.Fatalf("trial %d flags differ", trial)
	}
	if !reflect.DeepEqual(lifted.writable, direct.writable) {
		t.Fatalf("trial %d memory writes differ", trial)
	}
}
