package recovery

import (
	"math/rand"
	"reflect"
	"strconv"
	"testing"
)

func TestSINFComplexWriterMatchesLiftedBlocks(t *testing.T) {
	random := rand.New(rand.NewSource(int64(sinfComplexWriterStart)))
	profile := embedded26A428ExecutionProfile()
	front, ok := profile.sinfEntryBlock(sinfComplexWriterStart)
	if !ok {
		t.Fatal("SINF complex-writer front is unavailable")
	}
	tail, ok := profile.sinfEntryBlock(sinfComplexWriterTail)
	if !ok {
		t.Fatal("SINF complex-writer tail is unavailable")
	}
	for trial := 0; trial < 32; trial++ {
		registers := [31]uint64{}
		for register := range registers {
			registers[register] = random.Uint64()
		}
		stack := uint64(0x700000001000)
		pages := make(map[uint64][]byte)
		for attempt := 0; attempt < 96; attempt++ {
			directScout := newState(&model{pages: clonePages(pages)})
			directScout.x, directScout.sp = registers, stack
			_, directErr := runSINFComplexWriter(directScout)
			liftedScout := newState(&model{pages: clonePages(pages)})
			liftedScout.x, liftedScout.sp = registers, stack
			next, _, liftedErr := runCompactBlock(liftedScout, nil, front)
			if liftedErr == nil && next == sinfComplexWriterTail {
				_, _, liftedErr = runCompactBlock(liftedScout, nil, tail)
			}
			if directErr == nil && liftedErr == nil {
				break
			}
			err := directErr
			if liftedErr != nil {
				err = liftedErr
			}
			match := missingOracleAddress.FindStringSubmatch(err.Error())
			if len(match) != 2 {
				t.Fatalf("trial %d discovery %d: %v", trial, attempt, err)
			}
			address, parseErr := strconv.ParseUint(match[1], 16, 64)
			if parseErr != nil {
				t.Fatal(parseErr)
			}
			pages[address&^(pageSize-1)] = make([]byte, pageSize)
			if attempt == 95 {
				t.Fatalf("trial %d did not finish memory discovery", trial)
			}
		}
		lifted, direct := newState(&model{pages: clonePages(pages)}), newState(&model{pages: clonePages(pages)})
		lifted.x, lifted.sp = registers, stack
		direct.x, direct.sp = registers, stack
		next, stopped, err := runCompactBlock(lifted, nil, front)
		if err != nil || stopped || next != sinfComplexWriterTail {
			t.Fatalf("trial %d lifted front: next=%#x stopped=%v err=%v", trial, next, stopped, err)
		}
		if _, err := runSINFComplexWriterFront(direct); err != nil {
			t.Fatalf("trial %d direct front: %v", trial, err)
		}
		if lifted.x != direct.x {
			t.Fatalf("trial %d front state differs\nlifted=%#v\ndirect=%#v", trial, lifted.x, direct.x)
		}
		next, stopped, err = runCompactBlock(lifted, nil, tail)
		if err != nil || stopped {
			t.Fatalf("trial %d lifted tail: next=%#x stopped=%v err=%v", trial, next, stopped, err)
		}
		directNext, err := runSINFComplexWriterFinish(direct)
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

func TestSINFComplexWriterNetworkMatchesLiftedSuffix(t *testing.T) {
	random := rand.New(rand.NewSource(int64(sinfComplexWriterStart ^ 0x42)))
	profile := embedded26A428ExecutionProfile()
	front, ok := profile.sinfEntryBlock(sinfComplexWriterStart)
	if !ok {
		t.Fatal("SINF complex-writer front is unavailable")
	}
	for trial := 0; trial < 1024; trial++ {
		lifted, direct := newState(&model{pages: map[uint64][]byte{}}), newState(&model{pages: map[uint64][]byte{}})
		for register := range lifted.x {
			lifted.x[register] = random.Uint64()
		}
		lifted.n, lifted.z, lifted.c, lifted.v = random.Intn(2) != 0, random.Intn(2) != 0, random.Intn(2) != 0, random.Intn(2) != 0
		direct.x = lifted.x
		direct.n, direct.z, direct.c, direct.v = lifted.n, lifted.z, lifted.c, lifted.v
		if err := runCompactKernel(lifted, compactKernel{start: sinfComplexWriterStart + 42*4, instructions: front.instructions[42:60]}); err != nil {
			t.Fatal(err)
		}
		runSINFComplexWriterNetworkFirst(direct)
		if lifted.x != direct.x {
			t.Fatalf("trial %d first network state differs\nlifted=%#v\ndirect=%#v", trial, lifted.x, direct.x)
		}
		if err := runCompactKernel(lifted, compactKernel{start: sinfComplexWriterStart + 60*4, instructions: front.instructions[60:76]}); err != nil {
			t.Fatal(err)
		}
		runSINFComplexWriterNetworkSecondFront(direct)
		if lifted.x != direct.x {
			t.Fatalf("trial %d second network front differs\nlifted=%#v\ndirect=%#v", trial, lifted.x, direct.x)
		}
		if err := runCompactKernel(lifted, compactKernel{start: sinfComplexWriterStart + 76*4, instructions: front.instructions[76:]}); err != nil {
			t.Fatal(err)
		}
		runSINFComplexWriterNetworkSecondFinish(direct)
		if lifted.x != direct.x {
			t.Fatalf("trial %d state differs\nlifted=%#v\ndirect=%#v", trial, lifted.x, direct.x)
		}
		if lifted.n != direct.n || lifted.z != direct.z || lifted.c != direct.c || lifted.v != direct.v {
			t.Fatalf("trial %d flags differ", trial)
		}
	}
}
