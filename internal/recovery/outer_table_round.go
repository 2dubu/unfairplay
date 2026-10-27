package recovery

// outerTableRoundState is the compact semantic state consumed by the 21
// protected table writers. Pending words and mix values are produced by the
// preceding outer-state step; Source and Codec are the three byte pairs used
// while updating state word 2.
type outerTableRoundState struct {
	Words                [5]uint32
	PendingS0, PendingS1 uint32
	PendingS2, PendingS3 uint32
	PendingS4            uint32
	Mix16, Mix17         uint32
	Lane19Mix, LateMix   uint32
	Source, Codec        [3]uint32
}

func applyOuterTableRound(table *[outerTableWords]uint32, base uint32, state *outerTableRoundState) (uint32, bool) {
	indexes, ok := outerTableIndexes(base)
	if !ok || table == nil || state == nil {
		return 0, false
	}
	s0, s1, s2, s3, s4 := state.Words[0], state.Words[1], state.Words[2], state.Words[3], state.Words[4]
	tap := func(offset uint32) uint32 { return table[(base+offset)%outerTableWords] }

	table[indexes[0]] ^= outerTableLane0Mask(state.PendingS2, s0, s1, s3, s4, tap(27), tap(69), tap(111))
	table[indexes[1]] ^= outerTableLane1Mask(state.PendingS1, state.PendingS2, state.PendingS3, s4)
	s2 = state.PendingS2
	table[indexes[2]] ^= outerTableLane2Mask(state.PendingS1, state.PendingS2, state.PendingS3, s0, s4)
	s0 = outerTableLane2NextS0(state.PendingS2, state.PendingS1, state.PendingS3, s0)
	table[indexes[3]] ^= outerTableLane3Mask(s0, state.PendingS3)
	s3 = state.PendingS3
	table[indexes[4]] ^= outerTableLane4Mask(table[4], table[101], table[115], table[119])
	table[indexes[5]] ^= outerTableLane5Mask(state.PendingS1, s3, s0)
	s1 = state.PendingS1
	s2 = outerTableLane5NextS2(
		s2, state.Source[0], state.Source[1], state.Source[2],
		state.Codec[0], state.Codec[1], state.Codec[2],
	)
	table[indexes[6]] ^= outerTableLane6Mask(s1, state.PendingS4, s0)
	s4 = state.PendingS4
	table[indexes[7]] ^= outerTableLane7Mask(s0, s1, s3, s4, table[indexes[7]])
	s0 = outerTableLane7NextS0(s0, s1, s3, s4)
	table[indexes[8]] ^= outerTableLane8Mask(table[28], table[72], table[76], table[123])
	s1 = outerTableLane8NextS1(state.Mix17, s2, state.Mix16, s4)
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
	table[indexes[19]] ^= outerTableLane19Mask(state.Lane19Mix, s1, s4, tap(2), tap(44), tap(86))
	s0 = outerTableLane19NextS0(state.Lane19Mix)
	table[indexes[20]] ^= outerTableLane20Mask(state.PendingS0, s1, state.LateMix, s3)
	s0 = state.PendingS0

	state.Words = [5]uint32{s0, s1, s2, s3, s4}
	next, ok := nextOuterTableBase(base)
	return next, ok
}
