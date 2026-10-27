package recovery

// outerTableRoundState is the compact semantic state consumed by the 21
// protected table writers. Hardware is the six-byte native identity consumed
// through a rotating three-byte window. Codec is the paired three-byte state
// used while updating state word 2.
type outerTableRoundState struct {
	Words         [5]uint32
	Hardware      [6]byte
	HardwareIndex uint8
	Codec         [17]byte
}

func applyOuterTableRound(table *[outerTableWords]uint32, base uint32, state *outerTableRoundState) (uint32, bool) {
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

func applyOuterSemanticIteration(table *[outerTableWords]uint32, base uint32, state *outerTableRoundState) (uint32, bool) {
	if !applyOuterPreRound(table, base, state) {
		return 0, false
	}
	return applyOuterTableRound(table, base, state)
}
