package recovery

import "fmt"

func runRawSourceBootstrapCFG(s *state, runtime *rawSourceRuntime, maxBlocks int) (uint64, int, error) {
	profile, err := executionProfileForModel(s.model)
	if err != nil {
		return 0, 0, err
	}
	pc := profile.bootstrapStart
	for blockCount := 0; blockCount < maxBlocks; blockCount++ {
		if pc == profile.bootstrapStop {
			s.traceBlock(pc)
			return pc, blockCount, nil
		}
		if runtime != nil && runtime.stopPC != 0 && pc == runtime.stopPC {
			return pc, blockCount, nil
		}
		if pc >= 0xf000000000000000 || pc == 0x720000000000 {
			s.traceBlock(pc)
			if runtime == nil {
				return pc, blockCount, fmt.Errorf("bootstrap runtime is unavailable")
			}
			next, err := runtime.handle(s, pc)
			if err != nil {
				return pc, blockCount, err
			}
			pc = next
			continue
		}
		// The entry block is the bootstrap's hot recurrence. Keep partial-block
		// diagnostic stops on the lifted block, but use its recovered direct Go
		// form for normal execution.
		if pc == profile.bootstrapStart && (runtime == nil || runtime.stopPC == 0) {
			next, err := runBootstrapMixer(s, runtime, profile)
			if err != nil {
				return pc, blockCount, err
			}
			pc = next
			continue
		}
		if (pc == bootstrapShortWordPass || pc == bootstrapLongWordPass) && (runtime == nil || runtime.stopPC == 0) {
			next, err := runBootstrapWordPass(s, pc)
			if err != nil {
				return pc, blockCount, err
			}
			pc = next
			continue
		}
		if (pc == bootstrapAccumulatorPassA || pc == bootstrapAccumulatorPassB) && (runtime == nil || runtime.stopPC == 0) {
			next, err := runBootstrapAccumulatorPass(s, pc)
			if err != nil {
				return pc, blockCount, err
			}
			pc = next
			continue
		}
		if pc == bootstrapShufflePass && (runtime == nil || runtime.stopPC == 0) {
			next, err := runBootstrapShufflePass(s)
			if err != nil {
				return pc, blockCount, err
			}
			pc = next
			continue
		}
		if pc == bootstrapForwardShuffleSeed && (runtime == nil || runtime.stopPC == 0) {
			next, err := runBootstrapForwardShufflePass(s)
			if err != nil {
				return pc, blockCount, err
			}
			pc = next
			continue
		}
		if pc == bootstrapReverseBytePass && (runtime == nil || runtime.stopPC == 0) {
			next, err := runBootstrapReverseBytePass(s)
			if err != nil {
				return pc, blockCount, err
			}
			pc = next
			continue
		}
		if pc == bootstrapLookupCopySeed && (runtime == nil || runtime.stopPC == 0) {
			next, err := runBootstrapLookupCopyPass(s)
			if err != nil {
				return pc, blockCount, err
			}
			pc = next
			continue
		}
		if pc == bootstrapPointerWalkSeed && (runtime == nil || runtime.stopPC == 0) {
			next, err := runBootstrapPointerWalkPass(s)
			if err != nil {
				return pc, blockCount, err
			}
			pc = next
			continue
		}
		if pc == bootstrapPairCopyPass && (runtime == nil || runtime.stopPC == 0) {
			next, err := runBootstrapPairCopyPass(s)
			if err != nil {
				return pc, blockCount, err
			}
			pc = next
			continue
		}
		if pc == bootstrapConstantFillSeed && (runtime == nil || runtime.stopPC == 0) {
			next, err := runBootstrapConstantFillPass(s)
			if err != nil {
				return pc, blockCount, err
			}
			pc = next
			continue
		}
		if block, ok := recoveryBlock(s, pc, profile.bootstrapBlocks, profile.supplementalBlocks, profile.verifiedFallbackBlocks); ok {
			next, stopped, err := runCompactBlock(s, runtime, block)
			if err != nil {
				return pc, blockCount, fmt.Errorf("compact bootstrap CFG %#x: %w", pc, err)
			}
			if stopped {
				return next, blockCount, nil
			}
			pc = next
			continue
		}
		return pc, blockCount, fmt.Errorf("bootstrap CFG reached unknown block %#x", pc)
	}
	return pc, maxBlocks, fmt.Errorf("bootstrap CFG exceeded %d blocks; pc=%#x", maxBlocks, pc)
}
