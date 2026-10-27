package recovery

func outerPostRoundTapMix(tap82, tap41, tap124, tap166 uint32) uint32 {
	return uint32(0x8901d237) ^
		(tap82 & 0x010a2108) ^
		(tap41 & 0x40300201) ^
		(tap82 & tap41 & 0x40200100) ^
		(tap124 & 0x03050000) ^
		(tap82 & tap124 & 0x21400000) ^
		(tap41 & tap124 & 0x00240000) ^
		(tap82 & tap41 & tap124 & 0x00600000) ^
		(tap166 & 0x14800490) ^
		(tap82 & tap166 & 0x10020580) ^
		(tap41 & tap166 & 0x00000840) ^
		(tap82 & tap41 & tap166 & 0x00000100) ^
		(tap124 & tap166 & 0x10800000) ^
		(tap82 & tap124 & tap166 & 0x10400000) ^
		(tap82 & tap41 & tap124 & tap166 & 0x00400000)
}

func outerPostRoundStep(x14, x16, tap82, tap41, tap124, tap166 uint32) (uint32, uint32) {
	nextX16 := x16 ^ outerPostRoundTapMix(tap82, tap41, tap124, tap166)
	nextX14 := x14 ^ nextX16 ^ 0x0f015c57
	return nextX14, nextX16
}

func applyOuterPostRound(table *[outerTableWords]uint32, base uint32, state *outerTableRoundState) bool {
	return applyOuterPostRoundWithSeed(table, base, state, nil)
}

func applyOuterPostRoundWithSeed(table *[outerTableWords]uint32, base uint32, state *outerTableRoundState, seed *uint32) bool {
	if table == nil || state == nil || base >= outerTableWords {
		return false
	}
	x14, x16 := state.Words[2], state.Words[0]
	writerSeed := state.Words[4]
	writerTableIndex := uint32(0)
	var controls [140]uint32
	if seed != nil {
		writerTableIndex = (*seed + 0xf587d089) % outerTableWords
	}
	indexes := [4]uint32{
		(base + 77) % outerTableWords,
		(base + 36) % outerTableWords,
		(base + 119) % outerTableWords,
		(base + 161) % outerTableWords,
	}
	for iteration := range 140 {
		controls[iteration] = x16
		x14, x16 = outerPostRoundStep(x14, x16, table[indexes[0]], table[indexes[1]], table[indexes[2]], table[indexes[3]])
		if seed != nil {
			*seed = outerWriterSeedStep(*seed)
		}
		for index := range indexes {
			indexes[index] = (indexes[index] + 36) % outerTableWords
		}
	}
	state.Words[0], state.Words[2] = x16, x14
	if seed != nil {
		output, ok := writerTransform26A428(table[:], controls[:], writerSeed, writerTableIndex)
		if !ok {
			return false
		}
		state.Writer = output
		copy(state.Source[92:], output[:])
	}
	return true
}
