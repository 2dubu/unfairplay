package recovery

func recoveryBlock(pc uint64, primary, supplemental, verifiedFallback map[uint64]compactBlock) (compactBlock, bool) {
	if block, ok := primary[pc]; ok {
		return block, true
	}
	if block, ok := supplemental[pc]; ok {
		return block, true
	}
	block, ok := verifiedFallback[pc]
	return block, ok
}

func recoverySINFEntryBlock(profile *executionProfile, pc uint64) (compactBlock, bool) {
	if block, ok := profile.sinfEntryBlock(pc); ok {
		return block, true
	}
	if block, ok := profile.supplementalBlocks[pc]; ok {
		return block, true
	}
	block, ok := profile.verifiedFallbackBlocks[pc]
	return block, ok
}
