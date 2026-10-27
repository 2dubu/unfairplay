package recovery

import (
	"bytes"
	"encoding/binary"
	"testing"
)

func TestSemanticOuterStageBridge(t *testing.T) {
	testSemanticOuterStageBridge(t, rawSourceBinaryHash)
}

func TestSemanticOuterStageBridgeIndependentBuildHash(t *testing.T) {
	buildHash := [32]byte{0x42, 0x6f, 0x75, 0x6e, 0x64, 0x61, 0x72, 0x79}
	testSemanticOuterStageBridge(t, buildHash)
}

func testSemanticOuterStageBridge(t *testing.T, buildHash [32]byte) {
	t.Helper()
	profile := embedded26A428ExecutionProfile()
	profile.binaryHash = buildHash
	layout := outerStage26A428Layout
	const stack = uint64(0x7000000eeb80)
	pages := make(map[uint64][]byte)
	mapRange := func(address uint64, size int) {
		for page := address &^ (pageSize - 1); page < address+uint64(size); page += pageSize {
			if pages[page] == nil {
				pages[page] = make([]byte, pageSize)
			}
		}
	}
	mapRange(stack, 0xe000)
	mapRange(layout.hardwareAddress, 6)
	mapRange(layout.codecAddress, 29)
	value := &model{
		binaryHash: buildHash,
		profile:    profile,
		sp:         stack,
		pages:      pages,
	}
	value.registers[1] = 0x60319d83
	value.registers[29] = stack + 0x11250
	worker := newState(value)
	worker.reset()

	var table [outerTableWords]uint32
	for index := range table {
		table[index] = uint32(index)*0x7feb352d ^ 0x846ca68b
	}
	var tableBytes [outerTableWords * 4]byte
	for index, word := range table {
		binary.LittleEndian.PutUint32(tableBytes[index*4:], word)
	}
	if err := worker.store(stack+layout.tableOffset, tableBytes[:]); err != nil {
		t.Fatal(err)
	}
	semantic := outerTableRoundState{
		Words:    [5]uint32{0x6a09e667, 0xbb67ae85, 0x3c6ef372, 0xa54ff53a, 0x510e527f},
		Hardware: [6]byte{0x12, 0x34, 0x56, 0x78, 0x9a, 0xbc},
		Codec: [33]byte{
			4: 0x10, 5: 0x21, 6: 0x32, 7: 0x43, 8: 0x54, 9: 0x65,
			10: 0x76, 11: 0x87, 12: 0x98, 13: 0xa9, 14: 0xba, 15: 0xcb,
			16: 0xdc, 17: 0xed, 18: 0xfe, 19: 0x0f, 20: 0x11, 21: 0x22,
			22: 0x33, 23: 0x44, 24: 0x55, 25: 0x66, 26: 0x77, 27: 0x88,
			28: 0x99, 29: 0xaa, 30: 0xbb, 31: 0xcc, 32: 0xdd,
		},
	}
	for index := range semantic.Source {
		semantic.Source[index] = byte(index*29 + 7)
	}
	for index, offset := range layout.wordOffsets {
		if err := worker.storeValue(stack+offset, uint64(semantic.Words[index]), 4); err != nil {
			t.Fatal(err)
		}
	}
	if err := worker.store(layout.hardwareAddress, semantic.Hardware[:]); err != nil {
		t.Fatal(err)
	}
	if err := worker.store(layout.codecAddress, semantic.Codec[4:]); err != nil {
		t.Fatal(err)
	}
	if err := worker.store(profile.semantic.sourceAddress, semantic.Source[:]); err != nil {
		t.Fatal(err)
	}
	if err := worker.storeValue(stack+layout.selectorOffset, uint64(outerSelectorFromHardwareIndex(0)), 4); err != nil {
		t.Fatal(err)
	}
	const initialSeed = uint32(0x12345678)
	if err := worker.storeValue(stack+layout.counterSeedOffset, uint64(initialSeed), 4); err != nil {
		t.Fatal(err)
	}
	if err := worker.storeValue(stack+layout.hardwarePointerSlot, layout.hardwareAddress, 8); err != nil {
		t.Fatal(err)
	}
	if err := worker.storeValue(stack+layout.codecEndPointerSlot, layout.codecAddress+29, 8); err != nil {
		t.Fatal(err)
	}

	expectedTable := table
	expectedState := semantic
	expectedControl := outerIterationControl{X1: uint32(value.registers[1]), CounterSeed: initialSeed}
	if iterations, ok := applyOuterNaturalLoop(&expectedTable, &expectedState, &expectedControl); !ok || iterations != outerNaturalIterationCount {
		t.Fatalf("reference natural loop=%d/%v", iterations, ok)
	}
	iterations, err := runSemanticOuterStage(worker)
	if err != nil || iterations != outerNaturalIterationCount {
		t.Fatalf("semantic outer stage=%d/%v", iterations, err)
	}

	actualTable, actualState, actualControl, err := loadOuterStageState(worker, profile.semantic, layout)
	if err != nil {
		t.Fatal(err)
	}
	if actualTable != expectedTable || actualState.Words != expectedState.Words || actualState.HardwareIndex != expectedState.HardwareIndex || actualState.OptionalWord != expectedState.OptionalWord || actualControl != expectedControl {
		t.Fatal("semantic outer state does not match direct transform")
	}
	actualSelectionState, err := worker.outputBytes(stack+layout.selectionStateOffset, len(expectedState.SelectionState))
	if err != nil || !bytes.Equal(actualSelectionState, expectedState.SelectionState[:]) {
		t.Fatalf("selection state does not match direct transform: %v", err)
	}
	if actualState.Source != expectedState.Source {
		t.Fatal("source prefix does not match direct transform")
	}
	expectedSource := make([]byte, 332)
	copy(expectedSource[192:], expectedState.Tail[:])
	copy(expectedSource[92:], expectedState.Writer[:])
	copy(expectedSource[:92], expectedState.Source[:92])
	actualSource, err := worker.outputBytes(profile.semantic.sourceAddress, len(expectedSource))
	if err != nil || !bytes.Equal(actualSource, expectedSource) {
		t.Fatalf("outer source mismatch: %v", err)
	}
	wantRegisters := [5]uint64{worker.x[3], worker.x[7], worker.x[20], worker.x[22], worker.x[27]}
	if wantRegisters != expectedState.CycleRegisters || worker.x[8] != profile.semantic.selectionStart || worker.x[21] != layout.codecAddress || worker.x[29] != value.registers[29] {
		t.Fatal("selection-entry registers do not match direct transform")
	}
	if worker.n || worker.z || worker.c || worker.v {
		t.Fatal("selection-entry flags are not clear")
	}
	for address, want := range map[uint64]uint64{
		stack + layout.hardwarePointerSlot: layout.hardwareAddress,
		stack + layout.framePointerSlot:    profile.semantic.frameAddress + 0xb8,
		stack + layout.codecEndPointerSlot: layout.codecAddress + 29,
	} {
		got, err := worker.load(address, 8)
		if err != nil || got != want {
			t.Fatalf("pointer %#x=%#x/%v, want %#x", address, got, err, want)
		}
	}

	// The selection stage invokes the same transform through a neighboring
	// stack frame and a new codec allocation. Prove that the compact nested ABI
	// resolves those runtime differences without another lifted CFG.
	nestedCodecAddress := layout.codecAddress + 0x100
	if err := worker.store(nestedCodecAddress, semantic.Codec[4:]); err != nil {
		t.Fatal(err)
	}
	currentTable, err := worker.outputBytes(stack+layout.tableOffset, int(outerTableWords)*4)
	if err != nil {
		t.Fatal(err)
	}
	if err := worker.store(stack+layout.tableOffset-0x20, currentTable); err != nil {
		t.Fatal(err)
	}
	if err := worker.storeValue(stack+layout.codecEndPointerSlot, nestedCodecAddress, 8); err != nil {
		t.Fatal(err)
	}
	worker.x[1] = 0x60319d83
	iterations, err = runSemanticNestedOuterStage(worker)
	if err != nil || iterations != outerNaturalIterationCount {
		t.Fatalf("semantic nested outer stage=%d/%v", iterations, err)
	}
	if worker.x[21] != nestedCodecAddress {
		t.Fatalf("nested codec register=%#x, want %#x", worker.x[21], nestedCodecAddress)
	}
}

func TestSemanticOuterStageRejectsUnknownSelector(t *testing.T) {
	profile := embedded26A428ExecutionProfile()
	layout := outerStage26A428Layout
	const stack = uint64(0x7000000eeb80)
	pages := make(map[uint64][]byte)
	for address := stack &^ (pageSize - 1); address <= stack+layout.tableOffset+uint64(outerTableWords)*4; address += pageSize {
		pages[address] = make([]byte, pageSize)
	}
	pages[layout.hardwareAddress&^(pageSize-1)] = make([]byte, pageSize)
	pages[layout.codecAddress&^(pageSize-1)] = make([]byte, pageSize)
	worker := newState(&model{binaryHash: rawSourceBinaryHash, profile: profile, sp: stack, pages: pages})
	worker.reset()
	if err := worker.storeValue(stack+layout.hardwarePointerSlot, layout.hardwareAddress, 8); err != nil {
		t.Fatal(err)
	}
	if err := worker.storeValue(stack+layout.codecEndPointerSlot, layout.codecAddress+29, 8); err != nil {
		t.Fatal(err)
	}
	if err := worker.storeValue(stack+layout.selectorOffset, 0xdeadbeef, 4); err != nil {
		t.Fatal(err)
	}
	if _, err := runSemanticOuterStage(worker); err == nil {
		t.Fatal("semantic outer stage accepted an unknown selector")
	}
}
