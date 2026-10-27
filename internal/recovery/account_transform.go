package recovery

import "fmt"

const (
	accountEntryByteGeneratorPass = uint64(0x100b1ddc0)
	accountEntryXORCopyPass       = uint64(0x1004bb284)
	accountEntryByteRemapPass     = uint64(0x100131cd8)
	accountEntryTableWriterPass   = uint64(0x10007e214)
	accountEntryWordMixerPass     = uint64(0x10007ba1c)
	accountEntryByteCopyPass      = uint64(0x1001f39b8)
	accountEntryFourWordPass      = uint64(0x1002376d4)
	accountEntryZeroFillPass      = uint64(0x10092144c)
	accountEntrySecondaryWordPass = uint64(0x1000fcc88)
	accountEntryIndexedBytePass   = uint64(0x1001d527c)
	accountEntryCompareSelect     = uint64(0x100af2fb8)
	accountEntryCompareTail       = uint64(0x100af3000)
	accountEntryRoundTailA        = uint64(0x100139000)
	accountEntryRoundTailB        = uint64(0x100147c00)
	accountEntryByteStoreStart    = uint64(0x100131bc0)
	accountEntryByteStoreTail     = uint64(0x100131c00)
)

func runAccountEntrySemanticBlock(s *state, pc uint64) (uint64, bool, error) {
	switch pc {
	case accountEntryByteGeneratorPass:
		next, err := runAccountEntryByteGenerator(s)
		return next, true, err
	case accountEntryXORCopyPass:
		next, err := runAccountEntryXORCopy(s)
		return next, true, err
	case accountEntryByteRemapPass:
		next, err := runAccountEntryByteRemap(s)
		return next, true, err
	case accountEntryTableWriterPass:
		next, err := runAccountEntryTableWriter(s)
		return next, true, err
	case accountEntryWordMixerPass:
		next, err := runAccountEntryWordMixer(s)
		return next, true, err
	case accountEntryByteCopyPass:
		next, err := runAccountEntryByteCopy(s)
		return next, true, err
	case accountEntryFourWordPass:
		next, err := runAccountEntryFourWordPass(s)
		return next, true, err
	case accountEntryZeroFillPass:
		next, err := runAccountEntryZeroFill(s)
		return next, true, err
	case accountEntrySecondaryWordPass:
		next, err := runAccountEntrySecondaryWordPass(s)
		return next, true, err
	case accountEntryIndexedBytePass:
		next, err := runAccountEntryIndexedBytePass(s)
		return next, true, err
	case accountEntryCompareSelect:
		next, err := runAccountEntryCompareSelect(s)
		return next, true, err
	case accountEntryRoundTailA:
		next, err := runAccountEntryRoundTailA(s)
		return next, true, err
	case accountEntryRoundTailB:
		next, err := runAccountEntryRoundTailB(s)
		return next, true, err
	case accountEntryByteStoreStart:
		next, err := runAccountEntryByteStore(s)
		return next, true, err
	case accountHotMixFront:
		next, err := runAccountHotMix(s)
		return next, true, err
	case accountHotMixSecond:
		next, err := runAccountHotMixSecond(s)
		return next, true, err
	case accountHotMixThirdFront:
		next, err := runAccountHotMixThird(s)
		return next, true, err
	case accountHotMixFourth:
		next, err := runAccountHotMixFourth(s)
		return next, true, err
	case accountBooleanWriterA:
		next, err := runAccountBooleanWriterA(s)
		return next, true, err
	case accountBooleanWriterB:
		next, err := runAccountBooleanWriterB(s)
		return next, true, err
	case accountBooleanWriterC:
		next, err := runAccountBooleanWriterC(s)
		return next, true, err
	case accountBooleanWriterD:
		next, err := runAccountBooleanWriterD(s)
		return next, true, err
	case accountSmallSelectorA:
		next, err := runAccountSmallSelectorA(s)
		return next, true, err
	case accountSmallSelectorB:
		next, err := runAccountSmallSelectorB(s)
		return next, true, err
	case accountBoundsSelector:
		next, err := runAccountBoundsSelector(s)
		return next, true, err
	case accountWordWriter:
		next, err := runAccountWordWriter(s)
		return next, true, err
	case accountRoundPermuteA:
		next, err := runAccountRoundPermuteA(s)
		return next, true, err
	case accountRoundPermuteB:
		next, err := runAccountRoundPermuteB(s)
		return next, true, err
	case accountSmallPairStart:
		next, err := runAccountSmallPair(s)
		return next, true, err
	case accountByteSelect:
		next, err := runAccountByteSelect(s)
		return next, true, err
	case accountRangeDispatch:
		next, err := runAccountRangeDispatch(s)
		return next, true, err
	case accountStoreReturn:
		next, err := runAccountStoreReturn(s)
		return next, true, err
	case accountOrderedSelect:
		next, err := runAccountOrderedSelect(s)
		return next, true, err
	case accountDescriptorMove:
		next, err := runAccountDescriptorMove(s)
		return next, true, err
	case accountLimitDispatch:
		next, err := runAccountLimitDispatch(s)
		return next, true, err
	case accountStateDispatchStart:
		next, err := runAccountStateDispatch(s)
		return next, true, err
	case accountWordRoundTail:
		next, err := runAccountWordRoundTail(s)
		return next, true, err
	case accountWordRoundStart:
		next, err := runAccountWordRound(s)
		return next, true, err
	case accountDualWordRoundFront:
		next, err := runAccountDualWordRound(s)
		return next, true, err
	case accountTripleWordRoundFront:
		next, err := runAccountTripleWordRound(s)
		return next, true, err
	case accountFourthWordRound:
		next, err := runAccountFourthWordRound(s)
		return next, true, err
	case accountScaledSelect:
		next, err := runAccountScaledSelect(s)
		return next, true, err
	case accountIndexedSelect:
		next, err := runAccountIndexedSelect(s)
		return next, true, err
	case accountFrameSetup:
		next, err := runAccountFrameSetup(s)
		return next, true, err
	case accountNestedReturn:
		next, err := runAccountNestedReturn(s)
		return next, true, err
	case accountShiftSelect:
		next, err := runAccountShiftSelect(s)
		return next, true, err
	case accountPairedSelectFront:
		next, err := runAccountPairedSelect(s)
		return next, true, err
	case accountDifferenceSelect:
		next, err := runAccountDifferenceSelect(s)
		return next, true, err
	case accountByteWriterA:
		next, err := runAccountByteWriterA(s)
		return next, true, err
	case accountByteWriterB:
		next, err := runAccountByteWriterB(s)
		return next, true, err
	case accountByteWriterCFront:
		next, err := runAccountByteWriterC(s)
		return next, true, err
	case accountNestedSetupStart:
		next, err := runAccountNestedSetup(s)
		return next, true, err
	case accountByteLimit:
		next, err := runAccountByteLimit(s)
		return next, true, err
	case accountSixByteCopy:
		next, err := runAccountSixByteCopy(s)
		return next, true, err
	case accountBitStore:
		next, err := runAccountBitStore(s)
		return next, true, err
	case accountByteAdvance:
		next, err := runAccountByteAdvance(s)
		return next, true, err
	case accountSumSelect:
		next, err := runAccountSumSelect(s)
		return next, true, err
	case accountOutputByte:
		next, err := runAccountOutputByte(s)
		return next, true, err
	case accountPointerSelect:
		next, err := runAccountPointerSelect(s)
		return next, true, err
	case accountCallReturn:
		next, err := runAccountCallReturn(s)
		return next, true, err
	case accountConditionalPick:
		next, err := runAccountConditionalPick(s)
		return next, true, err
	case accountBitPermuteWriter:
		next, err := runAccountBitPermuteWriter(s)
		return next, true, err
	case accountFrameCall:
		next, err := runAccountFrameCall(s)
		return next, true, err
	case accountByteRouteWriter:
		next, err := runAccountByteRouteWriter(s)
		return next, true, err
	case accountBytePairRoute:
		next, err := runAccountBytePairRoute(s)
		return next, true, err
	case accountIndexedByteRoute:
		next, err := runAccountIndexedByteRoute(s)
		return next, true, err
	case accountIndexedWordRoute:
		next, err := runAccountIndexedWordRoute(s)
		return next, true, err
	case accountCounterRoute:
		next, err := runAccountCounterRoute(s)
		return next, true, err
	case accountSignedByteRoute:
		next, err := runAccountSignedByteRoute(s)
		return next, true, err
	case accountSetupRoute:
		next, err := runAccountSetupRoute(s)
		return next, true, err
	case accountLengthRoute:
		next, err := runAccountLengthRoute(s)
		return next, true, err
	case accountStateCompareRoute:
		next, err := runAccountStateCompareRoute(s)
		return next, true, err
	case accountStateGuardRoute:
		next, err := runAccountStateGuardRoute(s)
		return next, true, err
	case accountRecordRoute:
		next, err := runAccountRecordRoute(s)
		return next, true, err
	case accountRecordWriterRoute:
		next, err := runAccountRecordWriterRoute(s)
		return next, true, err
	case accountSignedCompareRoute:
		next, err := runAccountSignedCompareRoute(s)
		return next, true, err
	case accountEightByteRoute:
		next, err := runAccountEightByteRoute(s)
		return next, true, err
	case accountByteValueRoute:
		next, err := runAccountByteValueRoute(s)
		return next, true, err
	case accountWideGuardRoute:
		next, err := runAccountWideGuardRoute(s)
		return next, true, err
	case accountByteGuardRoute:
		next, err := runAccountByteGuardRoute(s)
		return next, true, err
	case accountByteDispatch:
		next, err := runAccountByteDispatch(s)
		return next, true, err
	case accountByteCopyDispatch:
		next, err := runAccountByteCopyDispatch(s)
		return next, true, err
	case accountIndexedWordDispatch:
		next, err := runAccountIndexedWordDispatch(s)
		return next, true, err
	case accountIndexedWordStore:
		next, err := runAccountIndexedWordStore(s)
		return next, true, err
	case accountWordCopyRoute:
		next, err := runAccountWordCopyRoute(s)
		return next, true, err
	case accountRangeRoute:
		next, err := runAccountRangeRoute(s)
		return next, true, err
	case accountByteShiftRoute:
		next, err := runAccountByteShiftRoute(s)
		return next, true, err
	case accountByteWriteRoute:
		next, err := runAccountByteWriteRoute(s)
		return next, true, err
	case accountRuntimeSetup:
		next, err := runAccountRuntimeSetup(s)
		return next, true, err
	case accountTableOffsetRoute:
		next, err := runAccountTableOffsetRoute(s)
		return next, true, err
	case accountByteLoopRoute:
		next, err := runAccountByteLoopRoute(s)
		return next, true, err
	case accountBitByteLoopRoute:
		next, err := runAccountBitByteLoopRoute(s)
		return next, true, err
	case accountConstantRouteA:
		next, err := runAccountConstantRouteA(s)
		return next, true, err
	case accountConstantRouteB:
		next, err := runAccountConstantRouteB(s)
		return next, true, err
	case accountConstantRouteC:
		next, err := runAccountConstantRouteC(s)
		return next, true, err
	case accountFrameRouteSetup:
		next, err := runAccountFrameRouteSetup(s)
		return next, true, err
	case accountFrameRouteCheck:
		next, err := runAccountFrameRouteCheck(s)
		return next, true, err
	case accountGateSetupRoute:
		next, err := runAccountGateSetupRoute(s)
		return next, true, err
	case accountMaskGateRoute:
		next, err := runAccountMaskGateRoute(s)
		return next, true, err
	case accountDualMaskRoute:
		next, err := runAccountDualMaskRoute(s)
		return next, true, err
	case accountStackLoadRoute:
		next, err := runAccountStackLoadRoute(s)
		return next, true, err
	case accountReverseByteRouteA:
		next, err := runAccountReverseByteRoute(s, accountReverseByteRouteA, 0x426e)
		return next, true, err
	case accountReverseByteRouteB:
		next, err := runAccountReverseByteRoute(s, accountReverseByteRouteB, 0x744)
		return next, true, err
	case accountByteCompareRoute:
		next, err := runAccountByteCompareRoute(s)
		return next, true, err
	case accountByteTableRoute:
		next, err := runAccountByteTableRoute(s)
		return next, true, err
	case accountByteSignedGate:
		next, err := runAccountByteSignedGate(s)
		return next, true, err
	case accountFrameTargetSetup:
		next, err := runAccountFrameTargetSetup(s)
		return next, true, err
	case accountBitGateRoute:
		next, err := runAccountBitGateRoute(s)
		return next, true, err
	case accountDispatchSetup:
		next, err := runAccountDispatchSetup(s)
		return next, true, err
	case accountRecordGuardRoute:
		next, err := runAccountRecordGuardRoute(s)
		return next, true, err
	case accountStackGateSetup:
		next, err := runAccountStackGateSetup(s)
		return next, true, err
	case accountWideCompareRoute:
		next, err := runAccountWideCompareRoute(s)
		return next, true, err
	case accountMaskTargetRoute:
		next, err := runAccountMaskTargetRoute(s)
		return next, true, err
	case accountByteIndexRoute:
		next, err := runAccountByteIndexRoute(s)
		return next, true, err
	case accountWideOffsetRoute:
		next, err := runAccountWideOffsetRoute(s)
		return next, true, err
	case accountRouteConstant:
		next, err := runAccountRouteConstant(s)
		return next, true, err
	case accountCompareRoute:
		next, err := runAccountCompareRoute(s)
		return next, true, err
	case accountExpressionRoute:
		next, err := runAccountExpressionRoute(s)
		return next, true, err
	case accountBitTestRoute:
		next, err := runAccountBitTestRoute(s)
		return next, true, err
	case accountStackRouteSetup:
		next, err := runAccountStackRouteSetup(s)
		return next, true, err
	case accountStackRouteSave:
		next, err := runAccountStackRouteSave(s)
		return next, true, err
	case accountStackRouteLoad:
		next, err := runAccountStackRouteLoad(s)
		return next, true, err
	case accountCompoundRouteA:
		next, err := runAccountCompoundRouteA(s)
		return next, true, err
	case accountCompoundRouteB:
		next, err := runAccountCompoundRouteB(s)
		return next, true, err
	case accountCompoundRouteC:
		next, err := runAccountCompoundRouteC(s)
		return next, true, err
	case accountByteDifferenceRoute:
		next, err := runAccountByteDifferenceRoute(s)
		return next, true, err
	case accountByteSubtractRoute:
		next, err := runAccountByteSubtractRoute(s)
		return next, true, err
	case accountByteCCmpRoute:
		next, err := runAccountByteCCmpRoute(s)
		return next, true, err
	case accountTripleCompareRoute:
		next, err := runAccountTripleCompareRoute(s)
		return next, true, err
	case accountDirectReturn:
		next, err := runAccountDirectReturn(s)
		return next, true, err
	case accountZeroByteRoute:
		next, err := runAccountZeroByteRoute(s)
		return next, true, err
	case accountLimitTargetRoute:
		next, err := runAccountLimitTargetRoute(s)
		return next, true, err
	case accountBooleanTargetRoute:
		next, err := runAccountBooleanTargetRoute(s)
		return next, true, err
	case accountZeroFillRoute:
		next, err := runAccountZeroFillRoute(s)
		return next, true, err
	case accountDirectReturnSecond:
		next, err := runAccountDirectReturnAt(s, accountDirectReturnSecond)
		return next, true, err
	case accountByteTargetRoute:
		next, err := runAccountByteTargetRoute(s)
		return next, true, err
	case accountByteFourRoute:
		next, err := runAccountByteFourRoute(s)
		return next, true, err
	case accountComputedTargetRoute:
		next, err := runAccountComputedTargetRoute(s)
		return next, true, err
	case accountSignedByteTarget:
		next, err := runAccountSignedByteTarget(s)
		return next, true, err
	case accountPointerByteTarget:
		next, err := runAccountPointerByteTarget(s)
		return next, true, err
	case accountByteWriteDispatchA:
		next, err := runAccountByteWriteDispatchA(s)
		return next, true, err
	case accountByteWriteDispatchB:
		next, err := runAccountByteWriteDispatchB(s)
		return next, true, err
	case accountGreaterEqualRoute:
		next, err := runAccountGreaterEqualRoute(s)
		return next, true, err
	case accountLessThanRoute:
		next, err := runAccountLessThanRoute(s)
		return next, true, err
	case accountOffsetRoute:
		next, err := runAccountOffsetRoute(s)
		return next, true, err
	case accountSevenFoldRoute:
		next, err := runAccountSevenFoldRoute(s)
		return next, true, err
	case accountWideLimitRoute:
		next, err := runAccountWideLimitRoute(s)
		return next, true, err
	case accountEqualityDispatch:
		next, err := runAccountEqualityDispatch(s)
		return next, true, err
	case accountByteCCmpDispatch:
		next, err := runAccountByteCCmpDispatch(s)
		return next, true, err
	case accountOffsetDispatch:
		next, err := runAccountOffsetDispatch(s)
		return next, true, err
	case accountFrameGuardRoute:
		next, err := runAccountFrameGuardRoute(s)
		return next, true, err
	case accountReverseWriteRoute:
		next, err := runAccountReverseWriteRoute(s)
		return next, true, err
	case accountCCmpTargetRoute:
		next, err := runAccountCCmpTargetRoute(s)
		return next, true, err
	case accountSignedPairDispatch:
		next, err := runAccountSignedPairDispatch(s)
		return next, true, err
	case accountComputedByteDispatch:
		next, err := runAccountComputedByteDispatch(s)
		return next, true, err
	case accountRuntimeByteSetup:
		next, err := runAccountRuntimeByteSetup(s)
		return next, true, err
	case accountMultiplyTarget:
		next, err := runAccountMultiplyTarget(s)
		return next, true, err
	case accountFrameByteTarget:
		next, err := runAccountFrameByteTarget(s)
		return next, true, err
	case accountConditionalEpilogue:
		next, err := runAccountConditionalEpilogue(s)
		return next, true, err
	case accountCompareDispatchA:
		next, err := runAccountCompareDispatchA(s)
		return next, true, err
	case accountFrameDispatch:
		next, err := runAccountFrameDispatch(s)
		return next, true, err
	case accountFrameOffset:
		next, err := runAccountFrameOffset(s)
		return next, true, err
	case accountCompareDispatchB:
		next, err := runAccountCompareDispatchB(s)
		return next, true, err
	case accountLowBitTarget:
		next, err := runAccountLowBitTarget(s)
		return next, true, err
	case accountFrameCompareTarget:
		next, err := runAccountFrameCompareTarget(s)
		return next, true, err
	case accountMaskSideTarget:
		next, err := runAccountMaskSideTarget(s)
		return next, true, err
	case accountSavedWordTarget:
		next, err := runAccountSavedWordTarget(s)
		return next, true, err
	case accountArithmeticTarget:
		next, err := runAccountArithmeticTarget(s)
		return next, true, err
	case accountSavedPairTarget:
		next, err := runAccountSavedPairTarget(s)
		return next, true, err
	case accountSavedCompareDispatch:
		next, err := runAccountSavedCompareDispatch(s)
		return next, true, err
	case accountPairCompareDispatch:
		next, err := runAccountPairCompareDispatch(s)
		return next, true, err
	case accountByteMaskTarget:
		next, err := runAccountByteMaskTarget(s)
		return next, true, err
	case accountScaledCompareTarget:
		next, err := runAccountScaledCompareTarget(s)
		return next, true, err
	case accountShiftedBitTarget:
		next, err := runAccountShiftedBitTarget(s)
		return next, true, err
	case accountCCmpMaskTarget:
		next, err := runAccountCCmpMaskTarget(s)
		return next, true, err
	case accountConditionRoute:
		next, err := runAccountConditionRoute(s)
		return next, true, err
	case accountSignedIndexRoute:
		next, err := runAccountSignedIndexRoute(s)
		return next, true, err
	case accountStackCompareRoute:
		next, err := runAccountStackCompareRoute(s)
		return next, true, err
	case accountDualMaskTarget:
		next, err := runAccountDualMaskTarget(s)
		return next, true, err
	case accountSmallStoreReturn:
		next, err := runAccountSmallStoreReturn(s)
		return next, true, err
	case accountEqualityTarget:
		next, err := runAccountEqualityTarget(s)
		return next, true, err
	case accountBitWordTarget:
		next, err := runAccountBitWordTarget(s)
		return next, true, err
	case accountRoundSetupRoute:
		next, err := runAccountRoundSetupRoute(s)
		return next, true, err
	case accountMixerSetup:
		next, err := runAccountMixerSetup(s)
		return next, true, err
	case accountLargeBooleanTransform:
		next, err := runAccountLargeBooleanTransform(s)
		return next, true, err
	case accountTripleBooleanTransform:
		next, err := runAccountTripleBooleanTransform(s)
		return next, true, err
	case accountStateSetupTransform:
		next, err := runAccountStateSetupTransform(s)
		return next, true, err
	case accountDispatchSeedSetup:
		next, err := runAccountDispatchSeedSetup(s)
		return next, true, err
	case accountDispatchMixSetup:
		next, err := runAccountDispatchMixSetup(s)
		return next, true, err
	case accountByteBooleanStep:
		next, err := runAccountByteBooleanStep(s)
		return next, true, err
	case accountAccumulatorStep:
		next, err := runAccountAccumulatorStep(s)
		return next, true, err
	case accountLookupSetupTarget:
		next, err := runAccountLookupSetupTarget(s)
		return next, true, err
	case accountMultiplyCompareTarget:
		next, err := runAccountMultiplyCompareTarget(s)
		return next, true, err
	default:
		return pc, false, nil
	}
}

func runAccountEntryByteStore(s *state) (uint64, error) {
	s.traceBlock(accountEntryByteStoreStart)
	s.x[11] = 0x545dd584
	s.x[17] = 0x00025fdf
	s.x[11] = uint64(uint32(s.x[8]) + uint32(s.x[11]))
	s.x[11] = uint64(uint32(s.x[11]) & uint32(s.x[17]))
	base, err := s.load(s.x[1]+8, 8)
	if err != nil {
		return 0, fmt.Errorf("account-entry byte-store input base: %w", err)
	}
	s.x[17] = base
	s.x[4] = uint64(uint32(s.x[3]) + 0xef50d627)
	value, err := s.load(base+uint64(uint32(s.x[4])), 1)
	if err != nil {
		return 0, fmt.Errorf("account-entry byte-store input: %w", err)
	}
	s.x[17] = uint64(byte(value))
	s.x[4] = uint64(uint32(s.x[3]) + 1)
	s.x[5] = s.x[0] + 4
	if err := s.storeValue(s.x[5]+uint64(uint32(s.x[9])), uint64(byte(s.x[17])), 1); err != nil {
		return 0, fmt.Errorf("account-entry byte-store output: %w", err)
	}
	s.x[9] = uint64(uint32(s.x[9]) + 1)
	word, err := s.load(s.x[1]+4, 4)
	if err != nil {
		return 0, fmt.Errorf("account-entry byte-store word: %w", err)
	}
	s.x[17] = uint64(uint32(word))

	s.traceBlock(accountEntryByteStoreTail)
	s.x[5] = uint64(uint32(s.x[17]) ^ uint32(s.x[12]))
	s.x[5] = uint64(uint32(s.x[5]) + uint32(s.x[13]))
	s.x[6] = uint64(uint32(s.x[17]) ^ uint32(s.x[14]))
	s.x[6] = uint64(uint32(s.x[6]) + uint32(s.x[15]))
	s.x[5] = uint64(uint32(s.x[5]) ^ uint32(s.x[6]))
	s.x[6] = 0xcd463667
	s.x[6] = uint64(uint32(s.x[11]) + uint32(s.x[6]))
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[6]))
	s.x[17] = uint64(uint32(s.x[17]) + uint32(s.x[16]))
	s.x[17] = uint64(uint32(s.x[5]) ^ uint32(s.x[17]))
	s.x[5] = 0x3489f9a3
	s.x[6] = 0x13ae9f40
	s.x[7] = 0x4539237b
	s.x[3] = uint64(uint32(s.x[3]) + uint32(s.x[5]))
	s.compare32(uint32(s.x[3]), uint32(s.x[7]))
	first := uint32(boolValue(s.condition("lo")))
	s.x[5] = uint64(first)
	s.x[17] = uint64(uint32(s.x[17]) + uint32(s.x[6]))
	s.compare32(uint32(s.x[17]), uint32(s.x[7]))
	second := uint32(boolValue(s.condition("lo")))
	s.x[6] = uint64(second)
	s.x[5] = uint64(first ^ second)
	s.compare32(uint32(s.x[3]), uint32(s.x[17]))
	ordered := uint32(boolValue(s.condition("lo")))
	s.x[17] = uint64(ordered)
	s.compare32(uint32(s.x[5]), 0)
	selected := ordered
	if s.condition("ne") {
		selected = second
	}
	s.x[17] = uint64(selected * uint32(s.x[2]))
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[8]))
	next, err := s.load(s.x[10]+uint64(uint32(s.x[17]))*8, 8)
	if err != nil {
		return 0, fmt.Errorf("account-entry byte-store target: %w", err)
	}
	s.x[5] = next
	s.x[3] = s.x[4]
	s.x[17] = 0x45bc
	return next, nil
}

func runAccountEntryRoundTailA(s *state) (uint64, error) {
	s.traceBlock(accountEntryRoundTailA)
	s.x[9] = uint64((uint32(s.x[9]) & 0xffff) | 0xd7280000)
	s.x[11] = uint64(uint32(s.x[11]) ^ uint32(s.x[9]))
	s.x[15] = uint64(uint32(s.x[30]) ^ uint32(rotateRight(s.x[2], 2, 32)))
	s.x[13]++
	s.x[9] = 0x1cc848bb
	s.compare64(s.x[13], s.x[9])
	gate := uint32(boolValue(s.condition("eq")))
	s.x[16] = uint64(gate * 0xa0e6)
	s.x[16] = uint64(uint32(s.x[16]) ^ uint32(s.x[1]))
	next, err := loadBootstrapPointer(s, s.x[20], uint32(s.x[16]))
	if err != nil {
		return 0, fmt.Errorf("account-entry round-tail A target: %w", err)
	}
	s.x[5] = next
	s.x[2] = s.x[6]
	s.x[7] = s.x[19]
	s.x[21] = 0x1cc848bb
	s.x[17] = 0x584
	return next, nil
}

func runAccountEntryRoundTailB(s *state) (uint64, error) {
	s.traceBlock(accountEntryRoundTailB)
	s.x[9] = 0x4894147e
	s.x[11] = uint64(uint32(s.x[11]) ^ uint32(s.x[9]))
	s.x[15] = uint64(uint32(s.x[30]) ^ uint32(rotateRight(s.x[19], 2, 32)))
	s.x[16]++
	s.x[9] = 0x1cc848a7
	s.compare64(s.x[16], s.x[9])
	gate := uint32(boolValue(s.condition("eq")))
	s.x[13] = uint64(gate * 0xeb)
	s.x[13] = uint64(uint32(s.x[13]) ^ uint32(s.x[0]))
	next, err := loadBootstrapPointer(s, s.x[20], uint32(s.x[13]))
	if err != nil {
		return 0, fmt.Errorf("account-entry round-tail B target: %w", err)
	}
	s.x[21] = next
	s.x[19] = s.x[2]
	s.x[6] = s.x[7]
	s.x[13] = 0x1cc848a7
	s.x[17] = 0x584
	return next, nil
}

func runAccountEntryIndexedBytePass(s *state) (uint64, error) {
	const maximumBytes = 4096
	for index := 0; index < maximumBytes; index++ {
		next, err := runAccountEntryIndexedByteStep(s)
		if err != nil || next != accountEntryIndexedBytePass {
			return next, err
		}
	}
	return 0, fmt.Errorf("account-entry indexed-byte pass exceeded %d bytes", maximumBytes)
}

func runAccountEntryIndexedByteStep(s *state) (uint64, error) {
	s.traceBlock(accountEntryIndexedBytePass)
	value, err := s.load(s.x[26]+uint64(uint32(s.x[8])), 1)
	if err != nil {
		return 0, fmt.Errorf("account-entry indexed-byte input: %w", err)
	}
	s.x[10] = uint64(byte(value))
	s.x[8] = uint64(uint32(s.x[8]) + 1)
	s.x[15] = s.x[12] + s.x[13]
	s.x[16] = uint64(uint32(s.x[0]) & uint32(s.x[10]<<1))
	s.x[10] = uint64(uint32(s.x[10]) - uint32(s.x[16]))
	s.x[10] = uint64(uint32(s.x[10]) + 0x1b)
	if err := s.storeValue(s.x[21]+s.x[15], uint64(byte(s.x[10])), 1); err != nil {
		return 0, fmt.Errorf("account-entry indexed-byte output: %w", err)
	}
	s.x[20] = uint64(uint32(s.x[4]) + uint32(s.x[11]))
	s.x[10] = s.x[20] ^ s.x[2]
	s.x[12] = s.x[10] + s.x[12]
	s.compare64(s.x[12], s.x[3])
	dispatchIndex := uint32(boolValue(s.condition("eq")))*uint32(s.x[5]) ^ uint32(s.x[4])
	next, err := s.load(s.x[24]+uint64(dispatchIndex)*8, 8)
	if err != nil {
		return 0, fmt.Errorf("account-entry indexed-byte dispatch: %w", err)
	}
	s.x[10] = next
	s.x[17] = 0xe66f
	return next, nil
}

func runAccountEntryCompareSelect(s *state) (uint64, error) {
	s.traceBlock(accountEntryCompareSelect)
	s.x[8] = 0xd59bef3f
	s.x[22] = uint64(uint32(s.x[11]) + uint32(s.x[8]))
	s.x[8] = 0x41a7
	s.x[10] = uint64(uint32(s.x[11]) ^ uint32(s.x[8]))
	word, err := s.load(s.x[27], 4)
	if err != nil {
		return 0, fmt.Errorf("account-entry compare-select input: %w", err)
	}
	s.x[9] = uint64(uint32(word))
	s.x[8] = uint64(uint32(s.x[11]) + 0x785dac3f)
	s.compare32(uint32(s.x[9]), uint32(s.x[8]))
	s.x[8] = uint64(uint32(s.x[10]) + 0x04a755f9)
	if s.condition("eq") {
		s.compare32(uint32(s.x[25]), uint32(s.x[8]))
	} else {
		s.setNZCV(0)
	}
	dispatchIndex := uint32(boolValue(s.condition("ne")))*0x192 ^ uint32(s.x[11])
	s.traceBlock(accountEntryCompareTail)
	next, err := loadBootstrapPointer(s, s.x[30], dispatchIndex)
	if err != nil {
		return 0, fmt.Errorf("account-entry compare-select dispatch: %w", err)
	}
	s.x[5] = 0xce757bc5
	s.x[8] = next
	s.x[12] = 0x192
	s.x[17] = 0x8231
	return next, nil
}

func runAccountEntrySecondaryWordPass(s *state) (uint64, error) {
	const maximumWords = 4096
	for index := 0; index < maximumWords; index++ {
		next, err := runAccountEntrySecondaryWordStep(s)
		if err != nil || next != accountEntrySecondaryWordPass {
			return next, err
		}
	}
	return 0, fmt.Errorf("account-entry secondary word pass exceeded %d words", maximumWords)
}

func runAccountEntrySecondaryWordStep(s *state) (uint64, error) {
	s.traceBlock(accountEntrySecondaryWordPass)
	s.x[0] = uint64(uint32(s.x[15]) + 0x42)
	s.x[17] = s.x[3] << 2
	offsets := [...]uint64{
		0xfffffffedf91bc98,
		0xfffffffedf91bc84,
		0xfffffffedf91bc6c,
		0xfffffffedf91bc64,
	}
	registers := [...]int{2, 22, 22, 25}
	var words [4]uint32
	for index, offset := range offsets {
		value, err := s.load(s.x[20]+s.x[17]+offset, 4)
		if err != nil {
			return 0, fmt.Errorf("account-entry secondary word input %d: %w", index, err)
		}
		words[index] = uint32(value)
		s.x[registers[index]] = uint64(words[index])
	}
	s.x[2] = uint64(words[1] ^ words[0])
	s.x[22] = uint64(words[2] ^ words[3])
	s.x[2] = uint64(uint32(s.x[2]) ^ uint32(s.x[22]))
	s.x[2] = uint64(uint32(s.x[13]) ^ uint32(rotateRight(s.x[2], 31, 32)))
	s.x[17] += s.x[9]
	if err := s.storeValue(s.x[20]+s.x[17], uint64(uint32(s.x[2])), 4); err != nil {
		return 0, fmt.Errorf("account-entry secondary word output: %w", err)
	}
	s.x[3]++
	s.x[11] = 0x481b782a
	s.x[17] = s.x[0] + s.x[11]
	s.compare64(s.x[3], s.x[17])
	dispatchIndex := uint32(boolValue(s.condition("eq")))*0x46 ^ uint32(s.x[15])
	next, err := loadBootstrapPointer(s, s.x[24], dispatchIndex)
	if err != nil {
		return 0, fmt.Errorf("account-entry secondary word dispatch: %w", err)
	}
	s.x[2] = 0x481b90d7
	s.x[17] = 0xaaed
	s.x[22] = next
	return next, nil
}

func runAccountEntryZeroFill(s *state) (uint64, error) {
	const maximumBytes = 4096
	for index := 0; index < maximumBytes; index++ {
		next, err := runAccountEntryZeroFillStep(s)
		if err != nil || next != accountEntryZeroFillPass {
			return next, err
		}
	}
	return 0, fmt.Errorf("account-entry zero fill exceeded %d bytes", maximumBytes)
}

func runAccountEntryZeroFillStep(s *state) (uint64, error) {
	s.traceBlock(accountEntryZeroFillPass)
	s.x[11] = s.x[16] + s.x[10]
	if err := s.storeValue(s.x[11], 0, 1); err != nil {
		return 0, fmt.Errorf("account-entry zero-fill output: %w", err)
	}
	s.x[11] = s.x[10] - 0x48f
	s.x[13] = s.x[8] ^ s.x[14]
	s.compare64(s.x[11], s.x[13])
	dispatchIndex := uint32(boolValue(s.condition("eq")))*uint32(s.x[15]) ^ uint32(s.x[9])
	next, err := loadBootstrapPointer(s, s.x[30], dispatchIndex)
	if err != nil {
		return 0, fmt.Errorf("account-entry zero-fill dispatch: %w", err)
	}
	s.x[10]++
	s.x[11] = next
	s.x[17] = 0x345f
	return next, nil
}

func runAccountEntryByteCopy(s *state) (uint64, error) {
	const maximumBytes = 4096
	for index := 0; index < maximumBytes; index++ {
		next, err := runAccountEntryByteCopyStep(s)
		if err != nil || next != accountEntryByteCopyPass {
			return next, err
		}
	}
	return 0, fmt.Errorf("account-entry byte copy exceeded %d bytes", maximumBytes)
}

func runAccountEntryByteCopyStep(s *state) (uint64, error) {
	s.traceBlock(accountEntryByteCopyPass)
	s.compare32(uint32(s.x[15]), uint32(s.x[14]))
	gate := uint32(boolValue(s.condition("lo")))
	s.x[10] = uint64(uint32(s.x[12]) + 0xb66ead0e)
	s.x[16] = 0x4990ff67
	s.x[10] = uint64(uint32(s.x[10]) & uint32(s.x[16]))
	s.x[16] = uint64(uint32(s.x[13]))
	value, err := s.load(s.x[9]+s.x[16], 1)
	if err != nil {
		return 0, fmt.Errorf("account-entry byte copy input: %w", err)
	}
	s.x[17] = uint64(byte(value))
	s.x[16] = s.x[8] + s.x[16]
	if err := s.storeValue(s.x[16], uint64(byte(s.x[17])), 1); err != nil {
		return 0, fmt.Errorf("account-entry byte copy output: %w", err)
	}
	s.x[16] = uint64(uint32(s.x[13]) + 0xffffb09b)
	s.x[16] = uint64(uint32(s.x[16]) + uint32(s.x[10]))
	s.x[17] = 0xe21ab501
	s.x[0] = 0x1de54aff
	s.compare32(uint32(s.x[16]), uint32(s.x[0]))
	selector := gate ^ uint32(boolValue(s.condition("hi")))
	s.x[13] = uint64(uint32(s.x[13]) + uint32(s.x[17]))
	s.compare32(uint32(s.x[13]), uint32(s.x[15]))
	secondBound := uint32(boolValue(s.condition("lo")))
	s.compare32(selector, 0)
	dispatchBit := secondBound
	if s.condition("ne") {
		dispatchBit = gate
	}
	dispatchIndex := dispatchBit*0x172c ^ uint32(s.x[12])
	next, err := loadBootstrapPointer(s, s.x[20], dispatchIndex)
	if err != nil {
		return 0, fmt.Errorf("account-entry byte copy dispatch: %w", err)
	}
	s.x[11] = next
	s.x[13] = s.x[16]
	s.x[17] = 0x63d5
	return next, nil
}

func runAccountEntryFourWordPass(s *state) (uint64, error) {
	const maximumWords = 4096
	for index := 0; index < maximumWords; index++ {
		next, err := runAccountEntryFourWordStep(s)
		if err != nil || next != accountEntryFourWordPass {
			return next, err
		}
	}
	return 0, fmt.Errorf("account-entry four-word pass exceeded %d words", maximumWords)
}

func runAccountEntryFourWordStep(s *state) (uint64, error) {
	s.traceBlock(accountEntryFourWordPass)
	s.x[9] = 0xffff5fad
	s.x[16] = uint64(uint32(s.x[13]) + uint32(s.x[9]))
	s.x[9] = 0x410
	s.x[0] = uint64(uint32(s.x[16]) | uint32(s.x[9]))
	s.x[16] = s.x[1] << 2
	offsets := [...]uint64{
		0xffffffff8cdedda8,
		0xffffffff8cdedd94,
		0xffffffff8cdedd7c,
	}
	first, err := s.load(s.x[8]+s.x[16]+offsets[0], 4)
	if err != nil {
		return 0, fmt.Errorf("account-entry four-word input 0: %w", err)
	}
	s.x[17] = uint64(uint32(first))
	second, err := s.load(s.x[8]+s.x[16]+offsets[1], 4)
	if err != nil {
		return 0, fmt.Errorf("account-entry four-word input 1: %w", err)
	}
	s.x[2] = uint64(uint32(second))
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[2]))
	third, err := s.load(s.x[8]+s.x[16]+offsets[2], 4)
	if err != nil {
		return 0, fmt.Errorf("account-entry four-word input 2: %w", err)
	}
	s.x[2] = uint64(uint32(third))
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[2]))
	s.x[2] = s.x[1] + s.x[0]
	s.x[9] = 0xffffffff8cde7fb0
	s.x[2] = s.x[9] + (s.x[2] << 2)
	fourth, err := s.load(s.x[8]+s.x[2], 4)
	if err != nil {
		return 0, fmt.Errorf("account-entry four-word input 3: %w", err)
	}
	s.x[2] = uint64(uint32(fourth))
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[2]))
	s.x[17] = uint64(uint32(s.x[12]) ^ uint32(rotateRight(s.x[17], 31, 32)))
	s.x[16] += s.x[14]
	if err := s.storeValue(s.x[8]+s.x[16], uint64(uint32(s.x[17])), 4); err != nil {
		return 0, fmt.Errorf("account-entry four-word output: %w", err)
	}
	s.x[1]++
	s.x[9] = 0x1cc848e3
	s.compare64(s.x[1], s.x[9])
	dispatchIndex := uint32(boolValue(s.condition("eq")))*0x7615 ^ uint32(s.x[13])
	next, err := loadBootstrapPointer(s, s.x[20], dispatchIndex)
	if err != nil {
		return 0, fmt.Errorf("account-entry four-word dispatch: %w", err)
	}
	s.x[2] = next
	s.x[16] = 0x1cc84893
	s.x[17] = 0xe949
	return next, nil
}

func runAccountEntryWordMixer(s *state) (uint64, error) {
	const maximumWords = 4096
	for index := 0; index < maximumWords; index++ {
		next, err := runAccountEntryWordMixerStep(s)
		if err != nil || next != accountEntryWordMixerPass {
			return next, err
		}
	}
	return 0, fmt.Errorf("account-entry word mixer exceeded %d words", maximumWords)
}

func runAccountEntryWordMixerStep(s *state) (uint64, error) {
	s.traceBlock(accountEntryWordMixerPass)
	s.x[9] = uint64(uint32(s.x[15]) + 0xd458908a)
	s.x[11] = s.x[20] << 2
	addresses := [...]uint64{
		0xffffffff5160f1c8,
		0xffffffff5160f1b4,
		0xffffffff5160f19c,
		0xffffffff5160f194,
	}
	registers := [...]int{12, 13, 16, 17}
	for index, offset := range addresses {
		value, err := s.load(s.x[10]+s.x[11]+offset, 4)
		if err != nil {
			return 0, fmt.Errorf("account-entry word mixer input %d: %w", index, err)
		}
		s.x[registers[index]] = uint64(uint32(value))
	}
	s.x[12] = uint64(uint32(s.x[13]) ^ uint32(s.x[12]))
	s.x[13] = uint64(uint32(s.x[16]) ^ uint32(s.x[17]))
	s.x[12] = uint64(uint32(s.x[12]) ^ uint32(s.x[13]))
	s.x[12] = uint64(uint32(s.x[9]) ^ uint32(rotateRight(s.x[12], 31, 32)))
	s.x[13] = 0xce752fa0
	s.x[12] = uint64(uint32(s.x[12]) ^ uint32(s.x[13]))
	s.x[11] += s.x[6]
	if err := s.storeValue(s.x[10]+s.x[11], uint64(uint32(s.x[12])), 4); err != nil {
		return 0, fmt.Errorf("account-entry word mixer output: %w", err)
	}
	s.x[20]++
	s.x[11] = s.x[20] - s.x[15]
	s.x[12] = s.x[15] - s.x[20]
	s.x[11] |= s.x[12]
	s.x[11] >>= 63
	s.x[11] = uint64(uint32(s.x[11]) * 0x72)
	s.x[12] = uint64(uint32(s.x[15]) + 0xd4589955)
	dispatchIndex := uint32(s.x[11]) ^ uint32(s.x[12])
	next, err := loadBootstrapPointer(s, s.x[21], dispatchIndex)
	if err != nil {
		return 0, fmt.Errorf("account-entry word mixer dispatch: %w", err)
	}
	limit, err := s.load(s.sp+0x7c, 4)
	if err != nil {
		return 0, fmt.Errorf("account-entry word mixer limit: %w", err)
	}
	s.x[4] = next
	s.x[7] = s.x[0]
	s.x[11] = s.x[3]
	s.x[12] = s.x[26]
	s.x[13] = s.x[22]
	s.x[16] = 0x2ba7c38b
	s.x[17] = 0x3601
	s.x[19] = uint64(uint32(limit))
	return next, nil
}

func runAccountEntryTableWriter(s *state) (uint64, error) {
	const maximumBytes = 4096
	for index := 0; index < maximumBytes; index++ {
		next, err := runAccountEntryTableWriterStep(s)
		if err != nil || next != accountEntryTableWriterPass {
			return next, err
		}
	}
	return 0, fmt.Errorf("account-entry table writer exceeded %d bytes", maximumBytes)
}

func runAccountEntryTableWriterStep(s *state) (uint64, error) {
	s.traceBlock(accountEntryTableWriterPass)
	s.x[11] = uint64(uint32(s.x[15]) + uint32(s.x[19]))
	s.x[13] = uint64(uint32(s.x[11]) | uint32(s.x[20]))
	value, err := s.load(s.x[2]+uint64(uint32(s.x[9])), 1)
	if err != nil {
		return 0, fmt.Errorf("account-entry table writer input: %w", err)
	}
	s.x[11] = uint64(byte(value))
	s.x[25] = uint64(uint32(s.x[9]) + 1)
	s.x[9] = s.x[16] + s.x[22]
	s.x[9] += s.x[13]
	s.x[12] = uint64(uint32(s.x[23]) & uint32(s.x[11]<<1))
	s.x[11] = uint64(uint32(s.x[11]) - uint32(s.x[12]))
	s.x[11] = uint64(uint32(s.x[11]) + 0x1b)
	if err := s.storeValue(s.x[24]+s.x[9], uint64(byte(s.x[11])), 1); err != nil {
		return 0, fmt.Errorf("account-entry table writer output: %w", err)
	}
	s.x[16]++
	s.compare64(s.x[16], s.x[26])
	dispatchIndex := uint32(boolValue(s.condition("eq")))*uint32(s.x[27]) ^ uint32(s.x[15])
	next, err := loadBootstrapPointer(s, s.x[21], dispatchIndex)
	if err != nil {
		return 0, fmt.Errorf("account-entry table writer dispatch: %w", err)
	}
	s.x[1] = s.x[24]
	s.x[4] = 0x12d6a942
	s.x[9] = s.x[25]
	s.x[11] = next
	s.x[17] = 0x71c7
	return next, nil
}

func runAccountEntryXORCopy(s *state) (uint64, error) {
	const maximumBytes = 4096
	for index := 0; index < maximumBytes; index++ {
		next, err := runAccountEntryXORCopyStep(s)
		if err != nil || next != accountEntryXORCopyPass {
			return next, err
		}
	}
	return 0, fmt.Errorf("account-entry XOR copy exceeded %d bytes", maximumBytes)
}

func runAccountEntryXORCopyStep(s *state) (uint64, error) {
	s.traceBlock(accountEntryXORCopyPass)
	s.x[19] = uint64(uint32(s.x[12]) | uint32(s.x[13]))
	s.x[8] = s.x[11] + s.x[14]
	s.x[11] = s.addWithCarry(s.x[8], s.x[19], false, 64)
	value, err := s.load(s.x[9]+s.x[11], 1)
	if err != nil {
		return 0, fmt.Errorf("account-entry XOR copy input: %w", err)
	}
	s.x[8] = uint64(uint32(value) ^ uint32(s.x[15]))
	if err := s.storeValue(s.x[0]+s.x[11], uint64(byte(s.x[8])), 1); err != nil {
		return 0, fmt.Errorf("account-entry XOR copy output: %w", err)
	}
	last := uint32(boolValue(s.condition("eq")))
	dispatchIndex := last*uint32(s.x[16]) ^ uint32(s.x[10])
	next, err := loadBootstrapPointer(s, s.x[20], dispatchIndex)
	if err != nil {
		return 0, fmt.Errorf("account-entry XOR copy dispatch: %w", err)
	}
	s.x[1] = next
	s.x[8] = 1
	s.x[17] = 0x56cc
	return next, nil
}

func runAccountEntryByteRemap(s *state) (uint64, error) {
	const maximumBytes = 4096
	for index := 0; index < maximumBytes; index++ {
		next, err := runAccountEntryByteRemapStep(s)
		if err != nil || next != accountEntryByteRemapPass {
			return next, err
		}
	}
	return 0, fmt.Errorf("account-entry byte remap exceeded %d bytes", maximumBytes)
}

func runAccountEntryByteRemapStep(s *state) (uint64, error) {
	s.traceBlock(accountEntryByteRemapPass)
	s.x[19] = uint64(uint32(s.x[14]) & uint32(s.x[0]))
	value, err := s.load(s.x[27]+uint64(uint32(s.x[11])), 1)
	if err != nil {
		return 0, fmt.Errorf("account-entry byte remap input: %w", err)
	}
	s.x[8] = uint64(byte(value))
	s.x[11] = uint64(uint32(s.x[11]) + 1)
	s.x[15] = s.x[13] + s.x[1]
	s.x[16] = uint64(uint32(s.x[19]) - 0x1a)
	s.x[16] = uint64(uint32(s.x[16]) & uint32(s.x[8]<<1))
	s.x[8] = uint64(uint32(s.x[8]) - uint32(s.x[16]))
	s.x[8] = uint64(uint32(s.x[8]) + 0x1b)
	if err := s.storeValue(s.x[12]+s.x[15], uint64(byte(s.x[8])), 1); err != nil {
		return 0, fmt.Errorf("account-entry byte remap output: %w", err)
	}
	s.x[13]++
	s.compare64(s.x[13], s.x[2])
	dispatchIndex := uint32(boolValue(s.condition("eq")))*uint32(s.x[5]) ^ uint32(s.x[9])
	next, err := s.load(s.x[20]+uint64(dispatchIndex)*8, 8)
	if err != nil {
		return 0, fmt.Errorf("account-entry byte remap dispatch: %w", err)
	}
	s.x[8] = next
	s.x[17] = 0xd90f
	return next, nil
}

func runAccountEntryByteGenerator(s *state) (uint64, error) {
	const maximumBytes = 4096
	for index := 0; index < maximumBytes; index++ {
		next, err := runAccountEntryByteGeneratorStep(s)
		if err != nil || next != accountEntryByteGeneratorPass {
			return next, err
		}
	}
	return 0, fmt.Errorf("account-entry byte generator exceeded %d bytes", maximumBytes)
}

func runAccountEntryByteGeneratorStep(s *state) (uint64, error) {
	if s == nil {
		return 0, fmt.Errorf("account-entry byte-generator state is unavailable")
	}
	s.traceBlock(accountEntryByteGeneratorPass)
	s.x[0] = uint64(uint32(s.x[16]) ^ uint32(s.x[10]))
	s.x[15] = (s.x[0] ^ s.x[19]) + s.x[1]
	if err := s.storeValue(s.sp+0x318, s.x[15], 8); err != nil {
		return 0, fmt.Errorf("account-entry byte-generator position: %w", err)
	}
	s.x[17] = uint64(uint32(s.x[13]) ^ uint32(s.x[20]))
	s.x[1] = uint64(uint32(s.x[13]) ^ uint32(s.x[23]))
	s.x[1] = uint64(uint32(s.x[1]) + uint32(s.x[26]))
	s.x[2] = uint64(uint32(s.x[13]) ^ uint32(s.x[27]))
	s.x[3] = uint64((uint32(s.x[13]) ^ 0xfffffffb) + 0x8f)
	s.x[17] = uint64(uint32(s.x[17]) + uint32(s.x[21]))
	s.x[17] = uint64(uint32(s.x[13]) ^ uint32(s.x[17]))
	s.x[2] = uint64(uint32(s.x[2]) + uint32(s.x[28]))
	s.x[1] = uint64(uint32(s.x[1]) ^ uint32(s.x[2]))
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[1]))
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[3]))
	s.x[17] = uint64(uint32(s.x[17]) ^ 0x30)
	s.x[1] = uint64(uint32(s.x[13]) + 0x87)
	s.x[17] = uint64(uint32(s.x[17]) * uint32(s.x[1]))
	if err := s.storeValue(s.x[6]+s.x[15], uint64(byte(s.x[17])), 1); err != nil {
		return 0, fmt.Errorf("account-entry byte-generator output: %w", err)
	}
	s.x[13]++
	position, err := s.load(s.sp+0x318, 8)
	if err != nil {
		return 0, fmt.Errorf("account-entry byte-generator saved position: %w", err)
	}
	s.x[1] = position
	s.compare64(position, 0)
	dispatchIndex := uint32(boolValue(s.condition("eq")))*uint32(s.x[14]) ^ uint32(s.x[16])
	next, err := loadBootstrapPointer(s, s.x[30], dispatchIndex)
	if err != nil {
		return 0, fmt.Errorf("account-entry byte-generator dispatch: %w", err)
	}
	s.x[15] = next
	s.x[17] = 0x623c
	return next, nil
}
