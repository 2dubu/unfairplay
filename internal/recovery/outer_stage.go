package recovery

import (
	"encoding/binary"
	"fmt"
)

// outerStageLayout is the small build-specific ABI around the recovered outer
// transform. The transform itself is address-independent Go; only this layout
// describes where one fairplayd build keeps its live inputs and outputs.
type outerStageLayout struct {
	tableOffset          uint64
	sourceAddress        uint64
	wordOffsets          [5]uint64
	counterSeedOffset    uint64
	selectorOffset       uint64
	optionalWordOffset   uint64
	selectionStateOffset uint64
	hardwareAddress      uint64
	codecAddress         uint64
	hardwarePointerSlot  uint64
	framePointerSlot     uint64
	codecEndPointerSlot  uint64
}

var outerStage26A428Layout = outerStageLayout{
	tableOffset:          0xcbf0,
	wordOffsets:          [5]uint64{0x1134, 0x1e98, 0x40ac, 0x4aec, 0x81b4},
	counterSeedOffset:    0x8a3c,
	selectorOffset:       0x70bc,
	optionalWordOffset:   0x824c,
	selectionStateOffset: 0xca0,
	hardwareAddress:      0x710000015b40,
	codecAddress:         0x700012500004,
	hardwarePointerSlot:  0x1218,
	framePointerSlot:     0x4cc0,
	codecEndPointerSlot:  0x6cc8,
}

func loadOuterStageState(worker *state, profile *executionProfile, layout outerStageLayout) ([outerTableWords]uint32, outerTableRoundState, outerIterationControl, error) {
	var table [outerTableWords]uint32
	var semantic outerTableRoundState
	var control outerIterationControl
	if worker == nil || profile == nil {
		return table, semantic, control, fmt.Errorf("outer stage state is unavailable")
	}

	tableBytes, err := worker.outputBytes(worker.sp+layout.tableOffset, len(table)*4)
	if err != nil {
		return table, semantic, control, fmt.Errorf("load outer table: %w", err)
	}
	for index := range table {
		table[index] = binary.LittleEndian.Uint32(tableBytes[index*4:])
	}
	for index, offset := range layout.wordOffsets {
		value, err := worker.load(worker.sp+offset, 4)
		if err != nil {
			return table, semantic, control, fmt.Errorf("load outer word %d: %w", index, err)
		}
		semantic.Words[index] = uint32(value)
	}
	hardware, err := worker.outputBytes(layout.hardwareAddress, len(semantic.Hardware))
	if err != nil {
		return table, semantic, control, fmt.Errorf("load outer hardware state: %w", err)
	}
	copy(semantic.Hardware[:], hardware)
	codec, err := worker.outputBytes(layout.codecAddress, 29)
	if err != nil {
		return table, semantic, control, fmt.Errorf("load outer codec state: %w", err)
	}
	copy(semantic.Codec[4:], codec)
	sourceAddress := layout.sourceAddress
	if sourceAddress == 0 {
		sourceAddress = profile.sourceAddress
	}
	source, err := worker.outputBytes(sourceAddress, len(semantic.Source))
	if err != nil {
		return table, semantic, control, fmt.Errorf("load outer source: %w", err)
	}
	copy(semantic.Source[:], source)
	optional, err := worker.load(worker.sp+layout.optionalWordOffset, 4)
	if err != nil {
		return table, semantic, control, fmt.Errorf("load outer optional word: %w", err)
	}
	semantic.OptionalWord = uint32(optional)
	selector, err := worker.load(worker.sp+layout.selectorOffset, 4)
	if err != nil {
		return table, semantic, control, fmt.Errorf("load outer hardware selector: %w", err)
	}
	selectorMatched := false
	for index := uint8(0); index < uint8(len(semantic.Hardware)); index++ {
		if uint32(selector) == outerSelectorFromHardwareIndex(index) {
			semantic.HardwareIndex = index
			selectorMatched = true
			break
		}
	}
	if !selectorMatched {
		return table, semantic, control, fmt.Errorf("unsupported outer hardware selector %#x", uint32(selector))
	}
	seed, err := worker.load(worker.sp+layout.counterSeedOffset, 4)
	if err != nil {
		return table, semantic, control, fmt.Errorf("load outer counter seed: %w", err)
	}
	control = outerIterationControl{X1: uint32(worker.x[1]), CounterSeed: uint32(seed)}
	return table, semantic, control, nil
}

func storeOuterStageState(worker *state, profile *executionProfile, layout outerStageLayout, table *[outerTableWords]uint32, semantic *outerTableRoundState, control *outerIterationControl) error {
	if worker == nil || profile == nil || table == nil || semantic == nil || control == nil {
		return fmt.Errorf("outer stage result is unavailable")
	}
	tableBytes := make([]byte, len(table)*4)
	for index, value := range table {
		binary.LittleEndian.PutUint32(tableBytes[index*4:], value)
	}
	if err := worker.store(worker.sp+layout.tableOffset, tableBytes); err != nil {
		return fmt.Errorf("store outer table: %w", err)
	}
	for index, offset := range layout.wordOffsets {
		if err := worker.storeValue(worker.sp+offset, uint64(semantic.Words[index]), 4); err != nil {
			return fmt.Errorf("store outer word %d: %w", index, err)
		}
	}
	if err := worker.storeValue(worker.sp+layout.counterSeedOffset, uint64(control.CounterSeed), 4); err != nil {
		return fmt.Errorf("store outer counter seed: %w", err)
	}
	if err := worker.storeValue(worker.sp+layout.selectorOffset, uint64(outerSelectorFromHardwareIndex(semantic.HardwareIndex)), 4); err != nil {
		return fmt.Errorf("store outer hardware selector: %w", err)
	}
	if err := worker.storeValue(worker.sp+layout.optionalWordOffset, uint64(semantic.OptionalWord), 4); err != nil {
		return fmt.Errorf("store outer optional word: %w", err)
	}
	if err := worker.store(worker.sp+layout.selectionStateOffset, semantic.SelectionState[:]); err != nil {
		return fmt.Errorf("store outer selection state: %w", err)
	}

	// The three recovered writers overlap in this order in the native path.
	// Preserve the full 332-byte region, including the eight writer bytes that
	// lie beyond Source and the tail bytes that survive both later writes.
	sourceAddress := layout.sourceAddress
	if sourceAddress == 0 {
		sourceAddress = profile.sourceAddress
	}
	if err := worker.store(sourceAddress+192, semantic.Tail[:]); err != nil {
		return fmt.Errorf("store outer tail: %w", err)
	}
	if err := worker.store(sourceAddress+92, semantic.Writer[:]); err != nil {
		return fmt.Errorf("store outer writer: %w", err)
	}
	if err := worker.store(sourceAddress, semantic.Source[:92]); err != nil {
		return fmt.Errorf("store outer source: %w", err)
	}
	if err := worker.storeValue(worker.sp+layout.hardwarePointerSlot, layout.hardwareAddress, 8); err != nil {
		return fmt.Errorf("store outer hardware pointer: %w", err)
	}
	if err := worker.storeValue(worker.sp+layout.framePointerSlot, profile.frameAddress+0xb8, 8); err != nil {
		return fmt.Errorf("store outer frame pointer: %w", err)
	}
	if err := worker.storeValue(worker.sp+layout.codecEndPointerSlot, layout.codecAddress+29, 8); err != nil {
		return fmt.Errorf("store outer codec pointer: %w", err)
	}

	framePointer := worker.x[29]
	worker.x = [31]uint64{
		0xc24ef3bf1bc1dd34, 0x60319d87, 0x51d00879f89aacf1,
		semantic.CycleRegisters[0], 0x3b8c6c84, 0x035a233b,
		layout.hardwareAddress, semantic.CycleRegisters[1], profile.selectionStart,
		0x100ff1ef0, 0x0501549d, 0x83f6, 0x4d191d94, 0x5db37cfe,
		0xfa8eae82, 0x399b2000, 0xbeb1a5925f8ada8e, 0x931a, 0,
		0x8ef4a16b, semantic.CycleRegisters[2], layout.codecAddress,
		semantic.CycleRegisters[3], 0xe5d992d0, worker.sp + 0x6000,
		0xa70ab7f2, worker.sp + 0x6000, semantic.CycleRegisters[4],
		0x13b13b13b13b13b2, framePointer, worker.sp + 0x7000,
	}
	worker.vectors = [32][16]byte{}
	worker.n, worker.z, worker.c, worker.v = false, false, false, false
	return nil
}

// runSemanticOuterStage executes the recovered transform without interpreting
// fairplayd blocks. Its inputs and outputs are the same stage ABI used by the
// direct selection state machine.
func runSemanticOuterStage(worker *state) (int, error) {
	return runSemanticOuterStageWithLayout(worker, outerStage26A428Layout)
}

func runSemanticNestedOuterStage(worker *state) (int, error) {
	if worker == nil || worker.model == nil {
		return 0, fmt.Errorf("outer stage state is unavailable")
	}
	profile, err := executionProfileForModel(worker.model)
	if err != nil {
		return 0, err
	}
	layout := outerStage26A428Layout
	codecAddress, err := worker.load(worker.sp+layout.codecEndPointerSlot, 8)
	if err != nil {
		return 0, fmt.Errorf("load nested outer codec pointer: %w", err)
	}
	if codecAddress == 0 {
		return 0, fmt.Errorf("invalid nested outer codec pointer")
	}
	layout.codecAddress = codecAddress
	layout.tableOffset -= 0x20
	layout.sourceAddress = profile.sourceAddress - 0x20
	return runSemanticOuterStageWithLayout(worker, layout)
}

func runSemanticOuterStageWithLayout(worker *state, layout outerStageLayout) (int, error) {
	if worker == nil || worker.model == nil {
		return 0, fmt.Errorf("outer stage state is unavailable")
	}
	profile, err := executionProfileForModel(worker.model)
	if err != nil {
		return 0, err
	}
	if profile.binaryHash != rawSourceBinaryHash {
		return 0, fmt.Errorf("no semantic outer stage for fairplayd build")
	}
	table, semantic, control, err := loadOuterStageState(worker, profile, layout)
	if err != nil {
		return 0, err
	}
	iterations, ok := applyOuterNaturalLoop(&table, &semantic, &control)
	if !ok || iterations != outerNaturalIterationCount {
		return iterations, fmt.Errorf("semantic outer stage did not reach its natural exit")
	}
	if err := storeOuterStageState(worker, profile, layout, &table, &semantic, &control); err != nil {
		return iterations, err
	}
	return iterations, nil
}
