package recovery

import (
	"math/rand"
	"reflect"
	"regexp"
	"strconv"
	"testing"
)

var missingOracleAddress = regexp.MustCompile(`oracle memory address 0x([0-9a-f]+) is outside`)

func TestSharedEnrollmentCodecWriterMatchesLiftedBlock(t *testing.T) {
	random := rand.New(rand.NewSource(int64(sharedEnrollmentCodecWriter)))
	profile := embedded26A428ExecutionProfile()
	block, ok := profile.selectionBlocks[sharedEnrollmentCodecWriter]
	if !ok {
		block, ok = profile.supplementalBlocks[sharedEnrollmentCodecWriter]
	}
	if !ok {
		t.Fatal("shared codec-writer block is unavailable")
	}
	for trial := 0; trial < 32; trial++ {
		registers := [31]uint64{}
		for register := range registers {
			registers[register] = random.Uint64()
		}
		registers[9] = 1
		registers[14] = 0x710000000100
		registers[29] = 0x700000001000
		pages := make(map[uint64][]byte)
		for attempt := 0; attempt < 64; attempt++ {
			directScout := newState(&model{pages: clonePages(pages)})
			directScout.x = registers
			_, directErr := runSharedEnrollmentCodecWriter(directScout)
			liftedScout := newState(&model{pages: clonePages(pages)})
			liftedScout.x = registers
			_, _, liftedErr := runCompactBlock(liftedScout, nil, block)
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
			if attempt == 63 {
				t.Fatalf("trial %d did not finish memory discovery", trial)
			}
		}
		lifted, direct := newState(&model{pages: clonePages(pages)}), newState(&model{pages: clonePages(pages)})
		lifted.x, direct.x = registers, registers
		next, stopped, err := runCompactBlock(lifted, nil, block)
		if err != nil || stopped {
			t.Fatalf("trial %d lifted: next=%#x stopped=%v err=%v", trial, next, stopped, err)
		}
		directNext, err := runSharedEnrollmentCodecWriter(direct)
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
