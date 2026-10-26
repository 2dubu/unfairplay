package recovery

func recoveryBlock(pc uint64, primary, supplemental map[uint64]compactBlock) (compactBlock, bool) {
	if block, ok := primary[pc]; ok {
		return block, true
	}
	block, ok := supplemental[pc]
	return block, ok
}

func recoverySINFEntryBlock(profile *executionProfile, pc uint64) (compactBlock, bool) {
	if block, ok := profile.sinfEntryBlock(pc); ok {
		return block, true
	}
	block, ok := profile.supplementalBlocks[pc]
	return block, ok
}
