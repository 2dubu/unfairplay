package recovery

// outerTableRoundState is the compact semantic state consumed by one complete
// outer iteration. HardwareIndex is the first byte of the table-round window;
// the earlier digest window starts two bytes later. Codec is the 17-byte state
// shared by both stages.
type outerTableRoundState struct {
	Words          [5]uint32
	Hardware       [6]byte
	HardwareIndex  uint8
	Codec          [43]byte
	Source         [224]byte
	Tail           [140]byte
	Writer         [140]byte
	SelectionState [32]byte
	OptionalWord   uint32
	CycleRegisters [5]uint64
}

type outerIterationControl struct {
	X1          uint32
	CounterSeed uint32
}

const outerNaturalIterationCount = 317

func outerSelectorFromHardwareIndex(index uint8) uint32 {
	return 0x4fc66d49 + uint32((index+1)%6)
}

// applyOuterTableRoundReduced applies the recovered 21-lane table round in
// program order. Reads intentionally use the live table because a later lane
// can consume a slot written by an earlier lane for some rolling bases.
func applyOuterTableRoundReduced(table *[outerTableWords]uint32, base uint32, state *outerTableRoundState) (uint32, bool) {
	indexes, ok := outerTableIndexes(base)
	if !ok || table == nil || state == nil || state.HardwareIndex >= uint8(len(state.Hardware)) {
		return 0, false
	}
	s0, s1, s2, s3, s4 := state.Words[0], state.Words[1], state.Words[2], state.Words[3], state.Words[4]
	tap := func(offset uint32) uint32 { return table[(base+offset)%outerTableWords] }
	pendingS2 := outerTablePendingS2(s2, tap(70), tap(112), tap(154))
	pendingS1 := outerTableLane0PendingS1(pendingS2, s0, s1, s3, tap(27), tap(69), tap(111), tap(152))
	pendingS3 := outerTableLane0PendingS3(pendingS2, s0, s3, s4)

	table[indexes[0]] ^= outerTableLane0Mask(pendingS2, s0, s1, s3, s4, tap(27), tap(69), tap(111))
	table[indexes[1]] ^= outerTableLane1Mask(pendingS1, pendingS2, pendingS3, s4)
	s2 = pendingS2
	table[indexes[2]] ^= outerTableLane2Mask(pendingS1, pendingS2, pendingS3, s0, s4)
	s0 = outerTableLane2NextS0(pendingS2, pendingS1, pendingS3, s0)
	table[indexes[3]] ^= outerTableLane3Mask(s0, pendingS3)
	s3 = pendingS3
	table[indexes[4]] ^= outerTableLane4Mask(table[4], table[101], table[115], table[119])
	table[indexes[5]] ^= outerTableLane5Mask(pendingS1, s3, s0)
	s1 = pendingS1
	var source [3]uint32
	for index := range source {
		source[index] = uint32(state.Hardware[(int(state.HardwareIndex)+index)%len(state.Hardware)])
	}
	s2 = outerTableLane5NextS2(
		s2, source[0], source[1], source[2],
		uint32(state.Codec[4]), uint32(state.Codec[6]), uint32(state.Codec[8]),
	)
	pendingS4 := outerTableLane5PendingS4(s4, s0, pendingS1, s3)
	table[indexes[6]] ^= outerTableLane6Mask(s1, pendingS4, s0)
	s4 = pendingS4
	table[indexes[7]] ^= outerTableLane7Mask(s0, s1, s3, s4, table[indexes[7]])
	s0 = outerTableLane7NextS0(s0, s1, s3, s4)
	table[indexes[8]] ^= outerTableLane8Mask(table[28], table[72], table[76], table[123])
	s1 = outerTableLane8NextS1(s1, s2, s3, s4)
	table[indexes[9]] ^= outerTableLane9Mask(s2, s0, s4)
	table[indexes[10]] ^= outerTableLane10Mask(s2, s3, s4)
	table[indexes[11]] ^= outerTableLane11Mask(s2, s0, s3, s1)
	table[indexes[12]] ^= outerTableLane12Mask(s1, s0, s3, s4)
	table[indexes[13]] ^= outerTableLane13Mask(s3, s2, s4)
	table[indexes[14]] ^= outerTableLane14Mask(s0, s2, s3, s4)
	table[indexes[15]] ^= outerTableLane15Mask(table[27], table[56], table[88], table[141])
	s4 = outerTableLane15NextS4(s1, s4, s2, s0)
	table[indexes[16]] ^= outerTableLane16Mask(s1, s2, s3)
	table[indexes[17]] ^= outerTableLane17Mask(s0, s1, s3)
	table[indexes[18]] ^= outerTableLane18Mask(s1, s3, s0, s2, table[indexes[18]])
	lane19Mix := outerTableLane18Mix(s0, s2, s1, s3)
	table[indexes[19]] ^= outerTableLane19Mask(lane19Mix, s1, s4, tap(2), tap(44), tap(86))
	s0 = outerTableLane19NextS0(lane19Mix)
	lateMix := outerTableLane19LateMix(lane19Mix, tap(2), tap(44), tap(86))
	pendingS0 := lateMix ^ 0x71563c28
	table[indexes[20]] ^= outerTableLane20Mask(pendingS0, s1, lateMix, s3)
	s0 = pendingS0

	state.Words = [5]uint32{s0, s1, s2, s3, s4}
	state.HardwareIndex = (state.HardwareIndex + uint8(len(state.Hardware)) - 1) % uint8(len(state.Hardware))
	next, ok := nextOuterTableBase(base)
	return next, ok
}

func applyOuterTableRound(table *[outerTableWords]uint32, base uint32, state *outerTableRoundState) (uint32, bool) {
	return applyOuterTableRoundReduced(table, base, state)
}

func applyOuterSemanticIteration(table *[outerTableWords]uint32, base uint32, state *outerTableRoundState) (uint32, bool) {
	if !applyOuterPreRound(table, base, state) {
		return 0, false
	}
	next, ok := applyOuterTableRound(table, base, state)
	if !ok || !applyOuterPostRound(table, base, state) {
		return 0, false
	}
	return next, true
}

func applyOuterFullIteration(table *[outerTableWords]uint32, base uint32, state *outerTableRoundState, control *outerIterationControl) (uint32, bool) {
	if control == nil {
		return 0, false
	}
	initialSeed := control.CounterSeed
	firstOptionalIndex := outerOptionalFirstIndex(initialSeed)
	if firstOptionalIndex >= outerTableWords {
		return 0, false
	}
	encodedOptionalIndex := table[firstOptionalIndex] ^ 0x224ea33b
	secondOptionalIndex := encodedOptionalIndex - uint32((uint64(encodedOptionalIndex)*0x621b97c3)>>38)*outerTableWords
	if secondOptionalIndex >= outerTableWords {
		return 0, false
	}
	optionalInput := outerOptionalDecode(table[secondOptionalIndex])
	state.OptionalWord = optionalInput
	if outerOptionalActive(optionalInput) {
		transformed, ok := optionalTransform26A428(optionalInput, 0x5d5c)
		if !ok {
			return 0, false
		}
		state.OptionalWord = transformed
	}
	state.SelectionState = [32]byte{
		0x3b, 0x49, 0xf5, 0xf1, 0xb3, 0x5c, 0x4c, 0xe9,
		0xa1, 0xb9, 0x92, 0x63, 0x6a, 0xd5, 0xb9, 0x2c,
		0x2d, 0xf3, 0xfd, 0x98, 0x4b, 0xcd, 0x22, 0x37,
		0x4b, 0x63, 0x86, 0x18, 0x23, 0xd3, 0x34, 0x66,
	}
	var hookSeeds outerHookSeeds
	if !applyOuterFront(table, base, state, control.X1, control.CounterSeed, &hookSeeds) {
		return 0, false
	}
	evenSeed := advanceOuterSeed(initialSeed, (144*base+5)%outerTableWords)
	oddSeed := advanceOuterSeed(initialSeed, (144*base+129)%outerTableWords)
	evenIndex := (evenSeed + 0xf587d089) % outerTableWords
	oddIndex := (oddSeed + 0xf587d089) % outerTableWords
	pair, ok := pairTransform26A428(table[:], hookSeeds.Even, hookSeeds.Odd, evenIndex, oddIndex)
	if !ok {
		return 0, false
	}
	copy(state.Source[32:], pair[:])
	tailSeed := advanceOuterSeed(initialSeed, (139*base)%outerTableWords)
	tailIndex := (tailSeed + 0xf587d089) % outerTableWords
	tail, ok := tailTransform26A428(table[:], hookSeeds.Tail, tailIndex)
	if !ok {
		return 0, false
	}
	state.Tail = tail
	if !applyOuterPreRound(table, base, state) {
		return 0, false
	}
	next, ok := applyOuterTableRound(table, base, state)
	if !ok {
		return 0, false
	}
	control.CounterSeed = advanceOuterSeed(control.CounterSeed, (139*base+29)%outerTableWords)
	if !applyOuterPostRoundWithSeed(table, base, state, &control.CounterSeed) {
		return 0, false
	}
	control.X1 = 0x60319e4b
	return next, true
}

func advanceOuterSeed(seed, iterations uint32) uint32 {
	for range iterations {
		seed = outerWriterSeedStep(seed)
	}
	return seed
}

func applyOuterNaturalLoop(table *[outerTableWords]uint32, state *outerTableRoundState, control *outerIterationControl) (int, bool) {
	return applyOuterLoop(table, state, control, outerNaturalIterationCount)
}

func applyOuterLoop(table *[outerTableWords]uint32, state *outerTableRoundState, control *outerIterationControl, iterations int) (int, bool) {
	if table == nil || state == nil || control == nil || iterations < 1 {
		return 0, false
	}
	base := outerTableInitialBase
	cycleIndex := uint32(0)
	for iteration := 1; iteration <= iterations; iteration++ {
		next, ok := applyOuterFullIteration(table, base, state, control)
		if !ok {
			return iteration - 1, false
		}
		base = next
		if iteration >= 9 && (iteration-9)%11 == 0 {
			cycleIndex++
			if !applyOuterCycle(state, cycleIndex) {
				return iteration, false
			}
			control.X1 = 0x60319d87
		}
	}
	return iterations, true
}
