package recovery

import (
	"encoding/binary"
	"math/rand"
	"reflect"
	"testing"
)

func TestBootstrapVectorExpandMatchesLiftedLoop(t *testing.T) {
	random := rand.New(rand.NewSource(int64(bootstrapVectorExpandStart)))
	profile := embedded26A428ExecutionProfile()
	for trial := 0; trial < 32; trial++ {
		pages := map[uint64][]byte{
			0x100d25000: make([]byte, pageSize),
			0x200000:    make([]byte, pageSize),
		}
		for _, page := range pages {
			if _, err := random.Read(page); err != nil {
				t.Fatal(err)
			}
		}
		lifted := newState(&model{pages: clonePages(pages)})
		direct := newState(&model{pages: clonePages(pages)})
		for register := range lifted.x {
			lifted.x[register] = random.Uint64()
		}
		for register := range lifted.vectors {
			if _, err := random.Read(lifted.vectors[register][:]); err != nil {
				t.Fatal(err)
			}
		}
		lifted.x[13] = 0x200000
		direct.x = lifted.x
		direct.vectors = lifted.vectors
		direct.n, direct.z, direct.c, direct.v = lifted.n, lifted.z, lifted.c, lifted.v

		next, stopped, err := runCompactBlock(lifted, nil, profile.bootstrapBlocks[bootstrapVectorExpandStart])
		if err != nil || stopped {
			t.Fatalf("trial %d lifted start: next=%#x stopped=%v err=%v", trial, next, stopped, err)
		}
		for next == bootstrapVectorExpandLoop {
			next, stopped, err = runCompactBlock(lifted, nil, profile.bootstrapBlocks[next])
			if err != nil || stopped {
				t.Fatalf("trial %d lifted loop: next=%#x stopped=%v err=%v", trial, next, stopped, err)
			}
		}
		directNext, err := runBootstrapVectorExpand(direct)
		if err != nil {
			t.Fatalf("trial %d direct: %v", trial, err)
		}
		assertBootstrapVectorState(t, trial, next, directNext, lifted, direct)
	}
}

func TestBootstrapVectorSubstituteMatchesLiftedLoop(t *testing.T) {
	random := rand.New(rand.NewSource(int64(bootstrapVectorSubstituteStart)))
	profile := embedded26A428ExecutionProfile()
	for trial := 0; trial < 32; trial++ {
		pages := map[uint64][]byte{
			0x100d25000: make([]byte, pageSize),
			0x100f49000: make([]byte, pageSize),
			0x200000:    make([]byte, pageSize),
			0x300000:    make([]byte, pageSize),
		}
		for _, page := range pages {
			if _, err := random.Read(page); err != nil {
				t.Fatal(err)
			}
		}
		binary.LittleEndian.PutUint64(pages[0x200000][0x48:], 0x300000)
		lifted := newState(&model{pages: clonePages(pages)})
		direct := newState(&model{pages: clonePages(pages)})
		for register := range lifted.x {
			lifted.x[register] = random.Uint64()
		}
		for register := range lifted.vectors {
			if _, err := random.Read(lifted.vectors[register][:]); err != nil {
				t.Fatal(err)
			}
		}
		lifted.x[4] = 0x200000
		lifted.x[11] = 0
		direct.x = lifted.x
		direct.vectors = lifted.vectors
		direct.n, direct.z, direct.c, direct.v = lifted.n, lifted.z, lifted.c, lifted.v

		next, stopped, err := runCompactBlock(lifted, nil, profile.bootstrapBlocks[bootstrapVectorSubstituteStart])
		if err != nil || stopped {
			t.Fatalf("trial %d lifted start: next=%#x stopped=%v err=%v", trial, next, stopped, err)
		}
		for next == bootstrapVectorSubstituteLoop {
			next, stopped, err = runCompactBlock(lifted, nil, profile.bootstrapBlocks[next])
			if err != nil || stopped {
				t.Fatalf("trial %d lifted loop: next=%#x stopped=%v err=%v", trial, next, stopped, err)
			}
		}
		directNext, err := runBootstrapVectorSubstitute(direct)
		if err != nil {
			t.Fatalf("trial %d direct: %v", trial, err)
		}
		assertBootstrapVectorState(t, trial, next, directNext, lifted, direct)
	}
}

func assertBootstrapVectorState(t *testing.T, trial int, liftedNext, directNext uint64, lifted, direct *state) {
	t.Helper()
	if liftedNext != directNext || lifted.x != direct.x {
		t.Fatalf("trial %d scalar state differs: next lifted=%#x direct=%#x\nlifted=%#v\ndirect=%#v", trial, liftedNext, directNext, lifted.x, direct.x)
	}
	if lifted.vectors != direct.vectors {
		var differences []int
		for register := range lifted.vectors {
			if lifted.vectors[register] != direct.vectors[register] {
				differences = append(differences, register)
				t.Logf("trial %d vector %d differs: lifted=%x direct=%x", trial, register, lifted.vectors[register], direct.vectors[register])
			}
		}
		t.Fatalf("trial %d vectors differ: %v", trial, differences)
	}
	if lifted.n != direct.n || lifted.z != direct.z || lifted.c != direct.c || lifted.v != direct.v {
		t.Fatalf("trial %d flags differ", trial)
	}
	if !reflect.DeepEqual(lifted.writable, direct.writable) {
		t.Fatalf("trial %d memory writes differ", trial)
	}
}
