package recovery

const (
	outerTableWords       = uint32(167)
	outerTableInitialBase = uint32(5)
)

// outerTableLaneOffsets are the 21 table positions updated by one outer
// iteration, relative to its rolling base. The order is the natural execution
// order of the semantic writers.
var outerTableLaneOffsets = [...]uint32{
	68, 25, 149, 106, 63, 20, 144,
	101, 58, 15, 139, 96, 53, 10,
	134, 91, 48, 5, 129, 43, 0,
}

func outerTableIndexes(base uint32) ([len(outerTableLaneOffsets)]uint32, bool) {
	var indexes [len(outerTableLaneOffsets)]uint32
	if base >= outerTableWords {
		return indexes, false
	}
	for lane, offset := range outerTableLaneOffsets {
		indexes[lane] = (base + offset) % outerTableWords
	}
	return indexes, true
}

func nextOuterTableBase(base uint32) (uint32, bool) {
	if base >= outerTableWords {
		return 0, false
	}
	return (24*base + 58) % outerTableWords, true
}
