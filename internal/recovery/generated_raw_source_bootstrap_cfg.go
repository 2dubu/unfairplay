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
		if pc == bootstrapPackedBitPass && (runtime == nil || runtime.stopPC == 0) {
			next, err := runBootstrapPackedBitPass(s)
			if err != nil {
				return pc, blockCount, err
			}
			pc = next
			continue
		}
		if pc == bootstrapTableAccumulatorPass && (runtime == nil || runtime.stopPC == 0) {
			next, err := runBootstrapTableAccumulatorPass(s)
			if err != nil {
				return pc, blockCount, err
			}
			pc = next
			continue
		}
		if pc == bootstrapSparseFillSeed && (runtime == nil || runtime.stopPC == 0) {
			next, err := runBootstrapSparseFillPass(s)
			if err != nil {
				return pc, blockCount, err
			}
			pc = next
			continue
		}
		if pc == bootstrapNonlinearBytePass && (runtime == nil || runtime.stopPC == 0) {
			next, err := runBootstrapNonlinearBytePass(s)
			if err != nil {
				return pc, blockCount, err
			}
			pc = next
			continue
		}
		if pc == bootstrapByteStatePass && (runtime == nil || runtime.stopPC == 0) {
			next, err := runBootstrapByteStatePass(s)
			if err != nil {
				return pc, blockCount, err
			}
			pc = next
			continue
		}
		if pc == bootstrapAccumulatorCycleSeed && (runtime == nil || runtime.stopPC == 0) {
			next, err := runBootstrapAccumulatorCycle(s)
			if err != nil {
				return pc, blockCount, err
			}
			pc = next
			continue
		}
		if pc == bootstrapByteWriterControl && (runtime == nil || runtime.stopPC == 0) {
			next, err := runBootstrapByteWriterCycle(s)
			if err != nil {
				return pc, blockCount, err
			}
			pc = next
			continue
		}
		if pc == bootstrapSetupStart && (runtime == nil || runtime.stopPC == 0) {
			next, err := runBootstrapSetupTransaction(s)
			if err != nil {
				return pc, blockCount, err
			}
			pc = next
			continue
		}
		if pc == bootstrapPointerCycleStart && (runtime == nil || runtime.stopPC == 0) {
			next, err := runBootstrapPointerCycle(s)
			if err != nil {
				return pc, blockCount, err
			}
			pc = next
			continue
		}
		if pc == bootstrapByteClassifyPass && (runtime == nil || runtime.stopPC == 0) {
			next, err := runBootstrapByteClassify(s)
			if err != nil {
				return pc, blockCount, err
			}
			pc = next
			continue
		}
		if pc == bootstrapVectorExpandStart && (runtime == nil || runtime.stopPC == 0) {
			next, err := runBootstrapVectorExpand(s)
			if err != nil {
				return pc, blockCount, err
			}
			pc = next
			continue
		}
		if pc == bootstrapVectorSubstituteStart && (runtime == nil || runtime.stopPC == 0) {
			next, err := runBootstrapVectorSubstitute(s)
			if err != nil {
				return pc, blockCount, err
			}
			pc = next
			continue
		}
		if pc == bootstrapIndexedDispatch && (runtime == nil || runtime.stopPC == 0) {
			next, err := runBootstrapIndexedDispatch(s)
			if err != nil {
				return pc, blockCount, err
			}
			pc = next
			continue
		}
		if pc == bootstrapCountdownStart && (runtime == nil || runtime.stopPC == 0) {
			next, err := runBootstrapCountdownDispatch(s)
			if err != nil {
				return pc, blockCount, err
			}
			pc = next
			continue
		}
		if pc == bootstrapSmallFrameReturn && (runtime == nil || runtime.stopPC == 0) {
			next, err := runBootstrapSmallFrameReturn(s)
			if err != nil {
				return pc, blockCount, err
			}
			pc = next
			continue
		}
		if pc == bootstrapSavedFrameReturn && (runtime == nil || runtime.stopPC == 0) {
			next, err := runBootstrapSavedFrameReturn(s)
			if err != nil {
				return pc, blockCount, err
			}
			pc = next
			continue
		}
		if pc == bootstrapProtectedCall && (runtime == nil || runtime.stopPC == 0) {
			next, err := runBootstrapProtectedCall(s)
			if err != nil {
				return pc, blockCount, err
			}
			pc = next
			continue
		}
		if pc == bootstrapRangeDispatchA && (runtime == nil || runtime.stopPC == 0) {
			next, err := runBootstrapRangeDispatchA(s)
			if err != nil {
				return pc, blockCount, err
			}
			pc = next
			continue
		}
		if pc == bootstrapRangeDispatchB && (runtime == nil || runtime.stopPC == 0) {
			next, err := runBootstrapRangeDispatchB(s)
			if err != nil {
				return pc, blockCount, err
			}
			pc = next
			continue
		}
		if pc == bootstrapStateDispatch && (runtime == nil || runtime.stopPC == 0) {
			next, err := runBootstrapStateDispatch(s)
			if err != nil {
				return pc, blockCount, err
			}
			pc = next
			continue
		}
		if pc == bootstrapWordDispatch && (runtime == nil || runtime.stopPC == 0) {
			next, err := runBootstrapWordDispatch(s)
			if err != nil {
				return pc, blockCount, err
			}
			pc = next
			continue
		}
		if pc == bootstrapScalarRoundStart && (runtime == nil || runtime.stopPC == 0) {
			next, err := runBootstrapScalarRound(s)
			if err != nil {
				return pc, blockCount, err
			}
			pc = next
			continue
		}
		if pc == bootstrapStreamMixStart && (runtime == nil || runtime.stopPC == 0) {
			next, err := runBootstrapStreamMix(s)
			if err != nil {
				return pc, blockCount, err
			}
			pc = next
			continue
		}
		if pc == bootstrapStreamState && (runtime == nil || runtime.stopPC == 0) {
			next, err := runBootstrapStreamState(s)
			if err != nil {
				return pc, blockCount, err
			}
			pc = next
			continue
		}
		if pc == bootstrapStreamAdvance && (runtime == nil || runtime.stopPC == 0) {
			next, err := runBootstrapStreamAdvance(s)
			if err != nil {
				return pc, blockCount, err
			}
			pc = next
			continue
		}
		if pc == bootstrapByteStreamPass && (runtime == nil || runtime.stopPC == 0) {
			next, err := runBootstrapByteStream(s)
			if err != nil {
				return pc, blockCount, err
			}
			pc = next
			continue
		}
		if pc == bootstrapLimitDispatch && (runtime == nil || runtime.stopPC == 0) {
			next, err := runBootstrapLimitDispatch(s)
			if err != nil {
				return pc, blockCount, err
			}
			pc = next
			continue
		}
		if pc == bootstrapOffsetDispatch && (runtime == nil || runtime.stopPC == 0) {
			next, err := runBootstrapOffsetDispatch(s)
			if err != nil {
				return pc, blockCount, err
			}
			pc = next
			continue
		}
		if pc == bootstrapCounterSelector && (runtime == nil || runtime.stopPC == 0) {
			next, err := runBootstrapCounterSelector(s)
			if err != nil {
				return pc, blockCount, err
			}
			pc = next
			continue
		}
		if pc == bootstrapByteSelector && (runtime == nil || runtime.stopPC == 0) {
			next, err := runBootstrapByteSelector(s)
			if err != nil {
				return pc, blockCount, err
			}
			pc = next
			continue
		}
		if pc == bootstrapWordTableRound && (runtime == nil || runtime.stopPC == 0) {
			next, err := runBootstrapWordTable(s)
			if err != nil {
				return pc, blockCount, err
			}
			pc = next
			continue
		}
		if pc == bootstrapDescriptorReturn && (runtime == nil || runtime.stopPC == 0) {
			next, err := runBootstrapDescriptorReturn(s)
			if err != nil {
				return pc, blockCount, err
			}
			pc = next
			continue
		}
		if pc == bootstrapCallContinuation && (runtime == nil || runtime.stopPC == 0) {
			next, err := runBootstrapCallContinuation(s)
			if err != nil {
				return pc, blockCount, err
			}
			pc = next
			continue
		}
		if pc == bootstrapNestedReturnStart && (runtime == nil || runtime.stopPC == 0) {
			next, err := runBootstrapNestedReturn(s)
			if err != nil {
				return pc, blockCount, err
			}
			pc = next
			continue
		}
		if pc == bootstrapGuardedReturnStart && (runtime == nil || runtime.stopPC == 0) {
			next, err := runBootstrapGuardedReturn(s)
			if err != nil {
				return pc, blockCount, err
			}
			pc = next
			continue
		}
		if runtime == nil || runtime.stopPC == 0 {
			next, handled, err := runCrossStageSemanticBlock(s, pc)
			if err != nil {
				return pc, blockCount, err
			}
			if handled {
				pc = next
				continue
			}
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
