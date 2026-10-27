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
		if pc == bootstrapStackSelector && (runtime == nil || runtime.stopPC == 0) {
			next, err := runBootstrapStackSelector(s)
			if err != nil {
				return pc, blockCount, err
			}
			pc = next
			continue
		}
		if pc == bootstrapScaledSelector && (runtime == nil || runtime.stopPC == 0) {
			next, err := runBootstrapScaledSelector(s)
			if err != nil {
				return pc, blockCount, err
			}
			pc = next
			continue
		}
		if pc == bootstrapDescriptorCall && (runtime == nil || runtime.stopPC == 0) {
			next, err := runBootstrapDescriptorCall(s)
			if err != nil {
				return pc, blockCount, err
			}
			pc = next
			continue
		}
		if pc == bootstrapNestedCallA && (runtime == nil || runtime.stopPC == 0) {
			next, err := runBootstrapNestedCallA(s)
			if err != nil {
				return pc, blockCount, err
			}
			pc = next
			continue
		}
		if pc == bootstrapNestedCallB && (runtime == nil || runtime.stopPC == 0) {
			next, err := runBootstrapNestedCallB(s)
			if err != nil {
				return pc, blockCount, err
			}
			pc = next
			continue
		}
		if pc == bootstrapSignSelector && (runtime == nil || runtime.stopPC == 0) {
			next, err := runBootstrapSignSelector(s)
			if err != nil {
				return pc, blockCount, err
			}
			pc = next
			continue
		}
		if pc == bootstrapPointerSelector && (runtime == nil || runtime.stopPC == 0) {
			next, err := runBootstrapPointerSelector(s)
			if err != nil {
				return pc, blockCount, err
			}
			pc = next
			continue
		}
		if pc == bootstrapStateSelector && (runtime == nil || runtime.stopPC == 0) {
			next, err := runBootstrapStateSelector(s)
			if err != nil {
				return pc, blockCount, err
			}
			pc = next
			continue
		}
		if pc == bootstrapCarrySelector && (runtime == nil || runtime.stopPC == 0) {
			next, err := runBootstrapCarrySelector(s)
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
		if pc == bootstrapByteCompareSelect && (runtime == nil || runtime.stopPC == 0) {
			next, err := runBootstrapByteCompareSelect(s)
			if err != nil {
				return pc, blockCount, err
			}
			pc = next
			continue
		}
		if pc == bootstrapDistanceSelect && (runtime == nil || runtime.stopPC == 0) {
			next, err := runBootstrapDistanceSelect(s)
			if err != nil {
				return pc, blockCount, err
			}
			pc = next
			continue
		}
		if pc == bootstrapSwapPass && (runtime == nil || runtime.stopPC == 0) {
			next, err := runBootstrapSwap(s)
			if err != nil {
				return pc, blockCount, err
			}
			pc = next
			continue
		}
		if pc == bootstrapReverseWordWrite && (runtime == nil || runtime.stopPC == 0) {
			next, err := runBootstrapReverseWordWrite(s)
			if err != nil {
				return pc, blockCount, err
			}
			pc = next
			continue
		}
		if pc == bootstrapSparseZeroWrite && (runtime == nil || runtime.stopPC == 0) {
			next, err := runBootstrapSparseZeroWrite(s)
			if err != nil {
				return pc, blockCount, err
			}
			pc = next
			continue
		}
		if pc == bootstrapDualZeroWrite && (runtime == nil || runtime.stopPC == 0) {
			next, err := runBootstrapDualZeroWrite(s)
			if err != nil {
				return pc, blockCount, err
			}
			pc = next
			continue
		}
		if pc == bootstrapEqualitySelect && (runtime == nil || runtime.stopPC == 0) {
			next, err := runBootstrapEqualitySelect(s)
			if err != nil {
				return pc, blockCount, err
			}
			pc = next
			continue
		}
		if pc == bootstrapParitySelect && (runtime == nil || runtime.stopPC == 0) {
			next, err := runBootstrapParitySelect(s)
			if err != nil {
				return pc, blockCount, err
			}
			pc = next
			continue
		}
		if pc == bootstrapBoundsDispatch && (runtime == nil || runtime.stopPC == 0) {
			next, err := runBootstrapBoundsDispatch(s)
			if err != nil {
				return pc, blockCount, err
			}
			pc = next
			continue
		}
		if pc == bootstrapLimitCompareStart && (runtime == nil || runtime.stopPC == 0) {
			next, err := runBootstrapLimitCompare(s)
			if err != nil {
				return pc, blockCount, err
			}
			pc = next
			continue
		}
		if pc == bootstrapCounterAdvance && (runtime == nil || runtime.stopPC == 0) {
			next, err := runBootstrapCounterAdvance(s)
			if err != nil {
				return pc, blockCount, err
			}
			pc = next
			continue
		}
		if pc == bootstrapZeroByteDispatch && (runtime == nil || runtime.stopPC == 0) {
			next, err := runBootstrapZeroByteDispatch(s)
			if err != nil {
				return pc, blockCount, err
			}
			pc = next
			continue
		}
		if pc == bootstrapPairReturn && (runtime == nil || runtime.stopPC == 0) {
			next, err := runBootstrapPairReturn(s)
			if err != nil {
				return pc, blockCount, err
			}
			pc = next
			continue
		}
		if pc == bootstrapIndexedCompare && (runtime == nil || runtime.stopPC == 0) {
			next, err := runBootstrapIndexedCompare(s)
			if err != nil {
				return pc, blockCount, err
			}
			pc = next
			continue
		}
		if pc == bootstrapLimitCompare && (runtime == nil || runtime.stopPC == 0) {
			next, err := runBootstrapLimitCompareDispatch(s)
			if err != nil {
				return pc, blockCount, err
			}
			pc = next
			continue
		}
		if pc == bootstrapGuardReturnStart && (runtime == nil || runtime.stopPC == 0) {
			next, err := runBootstrapGuardReturn(s)
			if err != nil {
				return pc, blockCount, err
			}
			pc = next
			continue
		}
		if pc == bootstrapSignedChoice && (runtime == nil || runtime.stopPC == 0) {
			next, err := runBootstrapSignedChoice(s)
			if err != nil {
				return pc, blockCount, err
			}
			pc = next
			continue
		}
		if pc == bootstrapStoreWordReturn && (runtime == nil || runtime.stopPC == 0) {
			next, err := runBootstrapStoreWordReturn(s)
			if err != nil {
				return pc, blockCount, err
			}
			pc = next
			continue
		}
		if pc == bootstrapBitLookupDispatch && (runtime == nil || runtime.stopPC == 0) {
			next, err := runBootstrapBitLookupDispatch(s)
			if err != nil {
				return pc, blockCount, err
			}
			pc = next
			continue
		}
		if pc == bootstrapSignedRange && (runtime == nil || runtime.stopPC == 0) {
			next, err := runBootstrapSignedRange(s)
			if err != nil {
				return pc, blockCount, err
			}
			pc = next
			continue
		}
		if pc == bootstrapFrameStatus && (runtime == nil || runtime.stopPC == 0) {
			next, err := runBootstrapFrameStatus(s)
			if err != nil {
				return pc, blockCount, err
			}
			pc = next
			continue
		}
		if pc == bootstrapByteMerge && (runtime == nil || runtime.stopPC == 0) {
			next, err := runBootstrapByteMerge(s)
			if err != nil {
				return pc, blockCount, err
			}
			pc = next
			continue
		}
		if pc == bootstrapBitSelect && (runtime == nil || runtime.stopPC == 0) {
			next, err := runBootstrapBitSelect(s)
			if err != nil {
				return pc, blockCount, err
			}
			pc = next
			continue
		}
		if pc == bootstrapWordLookup && (runtime == nil || runtime.stopPC == 0) {
			next, err := runBootstrapWordLookup(s)
			if err != nil {
				return pc, blockCount, err
			}
			pc = next
			continue
		}
		if pc == bootstrapByteUpdate && (runtime == nil || runtime.stopPC == 0) {
			next, err := runBootstrapByteUpdate(s)
			if err != nil {
				return pc, blockCount, err
			}
			pc = next
			continue
		}
		if pc == bootstrapByteStore && (runtime == nil || runtime.stopPC == 0) {
			next, err := runBootstrapByteStore(s)
			if err != nil {
				return pc, blockCount, err
			}
			pc = next
			continue
		}
		if pc == bootstrapPreparedDispatch && (runtime == nil || runtime.stopPC == 0) {
			next, err := runBootstrapPreparedDispatch(s)
			if err != nil {
				return pc, blockCount, err
			}
			pc = next
			continue
		}
		if pc == bootstrapStackChoice && (runtime == nil || runtime.stopPC == 0) {
			next, err := runBootstrapStackChoice(s)
			if err != nil {
				return pc, blockCount, err
			}
			pc = next
			continue
		}
		if pc == bootstrapFrameSetup && (runtime == nil || runtime.stopPC == 0) {
			next, err := runBootstrapFrameSetup(s)
			if err != nil {
				return pc, blockCount, err
			}
			pc = next
			continue
		}
		if pc == bootstrapDispatchSetup && (runtime == nil || runtime.stopPC == 0) {
			next, err := runBootstrapDispatchSetup(s)
			if err != nil {
				return pc, blockCount, err
			}
			pc = next
			continue
		}
		if pc == bootstrapRoundFront && (runtime == nil || runtime.stopPC == 0) {
			next, err := runBootstrapRound(s)
			if err != nil {
				return pc, blockCount, err
			}
			pc = next
			continue
		}
		if pc == bootstrapPreparedStart && (runtime == nil || runtime.stopPC == 0) {
			next, err := runBootstrapPreparedStart(s)
			if err != nil {
				return pc, blockCount, err
			}
			pc = next
			continue
		}
		if pc == bootstrapReducedByte && (runtime == nil || runtime.stopPC == 0) {
			next, err := runBootstrapReducedByte(s)
			if err != nil {
				return pc, blockCount, err
			}
			pc = next
			continue
		}
		if pc == bootstrapPointerTailA && (runtime == nil || runtime.stopPC == 0) {
			next, err := runBootstrapPointerTailA(s)
			if err != nil {
				return pc, blockCount, err
			}
			pc = next
			continue
		}
		if pc == bootstrapGuardReturnA && (runtime == nil || runtime.stopPC == 0) {
			next, err := runBootstrapGuardReturnA(s)
			if err != nil {
				return pc, blockCount, err
			}
			pc = next
			continue
		}
		if pc == bootstrapStoredDispatch && (runtime == nil || runtime.stopPC == 0) {
			next, err := runBootstrapStoredDispatch(s)
			if err != nil {
				return pc, blockCount, err
			}
			pc = next
			continue
		}
		if pc == bootstrapPackedSelector && (runtime == nil || runtime.stopPC == 0) {
			next, err := runBootstrapPackedSelector(s)
			if err != nil {
				return pc, blockCount, err
			}
			pc = next
			continue
		}
		if pc == bootstrapGuardReturnB && (runtime == nil || runtime.stopPC == 0) {
			next, err := runBootstrapGuardReturnB(s)
			if err != nil {
				return pc, blockCount, err
			}
			pc = next
			continue
		}
		if pc == bootstrapSignedSelector && (runtime == nil || runtime.stopPC == 0) {
			next, err := runBootstrapSignedSelector(s)
			if err != nil {
				return pc, blockCount, err
			}
			pc = next
			continue
		}
		if pc == bootstrapStackPivot && (runtime == nil || runtime.stopPC == 0) {
			next, err := runBootstrapStackPivot(s)
			if err != nil {
				return pc, blockCount, err
			}
			pc = next
			continue
		}
		if pc == bootstrapStackPairSelect && (runtime == nil || runtime.stopPC == 0) {
			next, err := runBootstrapStackPairSelect(s)
			if err != nil {
				return pc, blockCount, err
			}
			pc = next
			continue
		}
		if pc == bootstrapZeroStore && (runtime == nil || runtime.stopPC == 0) {
			next, err := runBootstrapZeroStore(s)
			if err != nil {
				return pc, blockCount, err
			}
			pc = next
			continue
		}
		if pc == bootstrapArithmeticDispatch && (runtime == nil || runtime.stopPC == 0) {
			next, err := runBootstrapArithmeticDispatch(s)
			if err != nil {
				return pc, blockCount, err
			}
			pc = next
			continue
		}
		if pc == bootstrapRangeStart && (runtime == nil || runtime.stopPC == 0) {
			next, err := runBootstrapRange(s)
			if err != nil {
				return pc, blockCount, err
			}
			pc = next
			continue
		}
		if pc == bootstrapStoreSelector && (runtime == nil || runtime.stopPC == 0) {
			next, err := runBootstrapStoreSelector(s)
			if err != nil {
				return pc, blockCount, err
			}
			pc = next
			continue
		}
		if pc == bootstrapPreparedSelector && (runtime == nil || runtime.stopPC == 0) {
			next, err := runBootstrapPreparedSelector(s)
			if err != nil {
				return pc, blockCount, err
			}
			pc = next
			continue
		}
		if pc == bootstrapBooleanFront && (runtime == nil || runtime.stopPC == 0) {
			next, err := runBootstrapBoolean(s)
			if err != nil {
				return pc, blockCount, err
			}
			pc = next
			continue
		}
		if pc == bootstrapPointerTailB && (runtime == nil || runtime.stopPC == 0) {
			next, err := runBootstrapPointerTailB(s)
			if err != nil {
				return pc, blockCount, err
			}
			pc = next
			continue
		}
		if pc == bootstrapRegisterSetup && (runtime == nil || runtime.stopPC == 0) {
			next, err := runBootstrapRegisterSetup(s)
			if err != nil {
				return pc, blockCount, err
			}
			pc = next
			continue
		}
		if pc == bootstrapProtectedFrameStart && (runtime == nil || runtime.stopPC == 0) {
			next, err := runBootstrapProtectedFrame(s)
			if err != nil {
				return pc, blockCount, err
			}
			pc = next
			continue
		}
		if pc == bootstrapDescriptorSetup && (runtime == nil || runtime.stopPC == 0) {
			next, err := runBootstrapDescriptorSetup(s)
			if err != nil {
				return pc, blockCount, err
			}
			pc = next
			continue
		}
		if pc == bootstrapNestedSetupStart && (runtime == nil || runtime.stopPC == 0) {
			next, err := runBootstrapNestedSetupStart(s)
			if err != nil {
				return pc, blockCount, err
			}
			pc = next
			continue
		}
		if pc == bootstrapNestedSetupCall && (runtime == nil || runtime.stopPC == 0) {
			next, err := runBootstrapNestedSetupCall(s)
			if err != nil {
				return pc, blockCount, err
			}
			pc = next
			continue
		}
		if pc == bootstrapNestedSetupResume && (runtime == nil || runtime.stopPC == 0) {
			next, err := runBootstrapNestedSetupResume(s)
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
