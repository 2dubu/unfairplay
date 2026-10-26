package recovery

func recoveryBlock(pc uint64, primary map[uint64]compactBlock) (compactBlock, bool) {
	if block, ok := primary[pc]; ok {
		return block, true
	}
	block, ok := compact26A428SupplementalBlocks[pc]
	return block, ok
}

func recoverySINFEntryBlock(pc uint64) (compactBlock, bool) {
	if block, ok := compact26A428SINFEntryBlock(pc); ok {
		return block, true
	}
	block, ok := compact26A428SupplementalBlocks[pc]
	return block, ok
}
