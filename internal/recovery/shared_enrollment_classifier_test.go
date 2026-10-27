package recovery

import (
	"math/rand"
	"testing"
)

func runSharedEnrollmentClassifyCircuit(s *state) (uint64, error) {
	s.traceBlock(sharedEnrollmentClassify)
	s.x[16] = uint64(uint32(s.x[12]) - 0xae)
	s.x[17] = uint64((uint32(s.x[16]) ^ 0x3a) + 0x6d)
	s.x[0] = uint64((uint32(s.x[16]) ^ 0x584ee074) + 0x23)
	s.x[1] = uint64((uint32(s.x[16]) ^ 0x18) + 0x4f)
	s.x[19] = uint64(uint32(s.x[13]) ^ uint32(s.x[16]))
	s.x[0] = uint64(uint32(s.x[19]) ^ uint32(s.x[0]))
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[1]))
	s.x[17] = uint64(uint32(s.x[0]) ^ uint32(s.x[17]))
	s.x[0] = 0x380f
	s.x[1] = uint64(uint32(s.x[2]) + uint32(s.x[0]))
	s.x[0] = 0x5f43
	s.x[0] = uint64(uint32(s.x[1]) ^ uint32(s.x[0]))
	s.x[19] = uint64(uint32(s.x[17]) & uint32(s.x[0]))
	s.x[17] = 0xa9
	s.x[0] = s.x[19] ^ s.x[17]
	value, err := s.load(s.x[28]+s.x[0], 1)
	if err != nil {
		return 0, err
	}
	s.x[20] = uint64(byte(value))
	s.x[17] = uint64(uint32(s.x[20]) ^ 0x22e6ee17)
	s.x[22] = 0xdd1911e9
	s.x[17] = uint64(uint32(s.x[17]) + uint32(s.x[22]))
	s.x[22] = uint64(uint32(s.x[20]) ^ 0x6482ee5e)
	s.x[23] = 0x9b7d11a2
	s.x[22] = uint64(uint32(s.x[22]) + uint32(s.x[23]))
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[22]))
	s.x[22] = uint64(uint32(s.x[20]) ^ 0x46640052)
	s.x[23] = 0xb99bffae
	s.x[22] = uint64(uint32(s.x[22]) + uint32(s.x[23]))
	s.x[22] = uint64(uint32(s.x[17]) ^ uint32(s.x[22]))
	s.x[23] = uint64(0xfffffff4 | (uint32(s.x[20]) & 3))
	next, err := loadBootstrapPointer(s, s.x[21], uint32(s.x[2]))
	if err != nil {
		return 0, err
	}
	s.x[2] = next
	s.x[17] = 0x1359
	return next, nil
}

func TestSharedEnrollmentClassifierMatchesCircuit(t *testing.T) {
	random := rand.New(rand.NewSource(int64(sharedEnrollmentClassify)))
	for trial := 0; trial < 1024; trial++ {
		const (
			stateBase = uint64(0x100000)
			tableBase = uint64(0x200000)
			target    = uint64(0x123456789)
		)
		pages := map[uint64][]byte{
			stateBase:                                make([]byte, pageSize),
			(tableBase + 0x27ad*8) &^ (pageSize - 1): make([]byte, pageSize),
		}
		for _, page := range pages {
			if _, err := random.Read(page); err != nil {
				t.Fatal(err)
			}
		}
		base := newState(&model{pages: clonePages(pages)})
		for register := range base.x {
			base.x[register] = random.Uint64()
		}
		iteration := random.Uint32()
		base.x[2] = 0x27ad
		base.x[12] = uint64(iteration + 0x22f1f257)
		base.x[13] = uint64(0xdd0e0f00 - iteration)
		base.x[21] = tableBase
		base.x[28] = stateBase
		if err := base.storeValue(tableBase+0x27ad*8, target, 8); err != nil {
			t.Fatal(err)
		}

		circuit := newState(base.model)
		semantic := newState(base.model)
		for _, item := range []*state{circuit, semantic} {
			item.x = base.x
			item.n, item.z, item.c, item.v = base.n, base.z, base.c, base.v
			item.writable = clonePages(base.writable)
			item.dirty = make(map[uint64]bool, len(item.writable))
			for page := range item.writable {
				item.dirty[page] = true
			}
		}
		circuitNext, circuitErr := runSharedEnrollmentClassifyCircuit(circuit)
		semanticNext, semanticErr := runSharedEnrollmentClassify(semantic)
		if circuitErr != nil || semanticErr != nil {
			t.Fatalf("trial %d errors: circuit=%v semantic=%v", trial, circuitErr, semanticErr)
		}
		if circuitNext != semanticNext || circuit.x != semantic.x ||
			circuit.n != semantic.n || circuit.z != semantic.z || circuit.c != semantic.c || circuit.v != semantic.v {
			t.Fatalf("trial %d state differs:\ncircuit=%#v\nsemantic=%#v", trial, circuit.x, semantic.x)
		}
	}
}

func TestSharedEnrollmentClassifierRejectsUnsupportedABI(t *testing.T) {
	s := newState(&model{pages: map[uint64][]byte{}})
	if _, err := runSharedEnrollmentClassify(s); err == nil {
		t.Fatal("unsupported shared enrollment classifier ABI was accepted")
	}
}

func runSharedEnrollmentBranchCircuit(s *state) (uint64, error) {
	s.traceBlock(sharedEnrollmentBranch)
	s.x[17] = uint64(uint32(s.x[1]) + 0x4064)
	s.x[2] = 0xc88
	s.x[24] = uint64(uint32(s.x[17]) | uint32(s.x[2]))
	s.x[17] = 0x13df
	s.x[25] = uint64(uint32(s.x[1]) + uint32(s.x[17]))
	s.x[17] = uint64(uint32(s.x[25]) ^ 0x6d)
	right := uint32(int32(int8(s.x[17])))
	s.compare32(uint32(s.x[23]), right)
	selected := uint32(boolValue(s.condition("ge")))
	s.x[2] = uint64(uint32(s.x[24]) ^ 0xace3)
	dispatchIndex := selected*uint32(s.x[2]) ^ uint32(s.x[1])
	next, err := loadBootstrapPointer(s, s.x[21], dispatchIndex)
	if err != nil {
		return 0, err
	}
	s.x[1] = next
	s.x[17] = 0x8eec
	return next, nil
}

func TestSharedEnrollmentBranchMatchesCircuit(t *testing.T) {
	random := rand.New(rand.NewSource(int64(sharedEnrollmentBranch)))
	for trial := 0; trial < 1024; trial++ {
		const tableBase = uint64(0x200000)
		pages := map[uint64][]byte{}
		for _, index := range []uint32{0x5fbc, 0x5ff7} {
			address := tableBase + uint64(index)*8
			if pages[address&^(pageSize-1)] == nil {
				pages[address&^(pageSize-1)] = make([]byte, pageSize)
			}
		}
		for _, page := range pages {
			if _, err := random.Read(page); err != nil {
				t.Fatal(err)
			}
		}
		base := newState(&model{pages: clonePages(pages)})
		for register := range base.x {
			base.x[register] = random.Uint64()
		}
		current := uint32(byte(random.Uint32()))
		base.x[1] = 0x5fbc
		base.x[20] = uint64(current)
		base.x[23] = uint64(0xfffffff4 | current&3)
		base.x[21] = tableBase

		circuit := newState(base.model)
		semantic := newState(base.model)
		circuit.x, semantic.x = base.x, base.x
		circuit.n, circuit.z, circuit.c, circuit.v = base.n, base.z, base.c, base.v
		semantic.n, semantic.z, semantic.c, semantic.v = base.n, base.z, base.c, base.v
		circuitNext, circuitErr := runSharedEnrollmentBranchCircuit(circuit)
		semanticNext, semanticErr := runSharedEnrollmentBranch(semantic)
		if circuitErr != nil || semanticErr != nil {
			t.Fatalf("trial %d errors: circuit=%v semantic=%v", trial, circuitErr, semanticErr)
		}
		if circuitNext != semanticNext || circuit.x != semantic.x ||
			circuit.n != semantic.n || circuit.z != semantic.z || circuit.c != semantic.c || circuit.v != semantic.v {
			t.Fatalf("trial %d state differs:\ncircuit=%#v\nsemantic=%#v", trial, circuit.x, semantic.x)
		}
	}
}
