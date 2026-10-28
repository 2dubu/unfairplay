package recovery

import "fmt"

// semanticBuildProfile contains only genuine per-build ABI and encoding data.
// The algorithms consuming it are shared typed Go and must not vary by build.
type semanticBuildProfile struct {
	outer            outerStageLayout
	sourceAddress    uint64
	frameAddress     uint64
	selectionStart   uint64
	currentSINFWhite [16]byte
	currentSINFInput [4096]byte
}

type executionProfile struct {
	name       string
	binaryHash [32]byte
	semantic   *semanticBuildProfile

	accountEntryStart, accountEntryStop uint64
	bootstrapStart, bootstrapStop       uint64
	selectionStart, selectionStop       uint64
	sinfEntryStart, sinfEntryStop       uint64
	selectionFrameCall                  uint64
	selectionStopHits                   int

	accountEntryBlocks     map[uint64]compactBlock
	bootstrapBlocks        map[uint64]compactBlock
	selectionBlocks        map[uint64]compactBlock
	supplementalBlocks     map[uint64]compactBlock
	verifiedFallbackBlocks map[uint64]compactBlock
	sinfEntryBlock         func(uint64) (compactBlock, bool)
	applyFrame             func(*state) (uint64, error)
	runBootstrap           func(*state, fault) error
	runDigestPrologue      func(*state, fault) error
	runDigest              func(*state, fault) error
	runOptional            func(*state, fault) error
	runEven                func(*state, fault) error
	runOdd                 func(*state, fault) error
	runTail                func(*state, fault) error
	runWriter              func(*state, fault) error
	runOuter               func(*state, fault) error
	runCycle               func(*state, fault) error
	runCycleFinal          func(*state, fault) error

	onceTokens    []uint64
	sourceAddress uint64
	frameAddress  uint64
	outerStart    uint64
	outerCycle    uint64
	cycleFinal    uint64
	optionalStart uint64
}

func embedded26A428ExecutionProfile() *executionProfile {
	return &executionProfile{
		name: "macos26a428", binaryHash: rawSourceBinaryHash,
		semantic: &semanticBuildProfile{
			outer:            outerStage26A428Layout,
			sourceAddress:    0x7000000FB304,
			frameAddress:     0x7000000FB724,
			selectionStart:   0x1006b10a0,
			currentSINFWhite: generatedCurrentSINFWhitening,
			currentSINFInput: generatedCurrentSINFInputDecode,
		},
		accountEntryStart: 0x100cd1c18, accountEntryStop: 0x10057ae74,
		bootstrapStart: 0x10057ae74, bootstrapStop: 0x1005b7438,
		selectionStart: 0x1006b10a0, selectionStop: 0x10008a644, selectionStopHits: 2,
		sinfEntryStart: 0x10008a644, sinfEntryStop: 0x100234cdc,
		selectionFrameCall:     0x1002fbb24,
		accountEntryBlocks:     compact26A428AccountEntryBlocks,
		bootstrapBlocks:        compact26A428BootstrapBlocks,
		selectionBlocks:        compact26A428SelectionBlocks,
		supplementalBlocks:     compact26A428SupplementalBlocks,
		verifiedFallbackBlocks: compact26A428VerifiedFallbackBlocks,
		sinfEntryBlock:         compact26A428SINFEntryBlock,
		applyFrame:             applyRecoveredEnrollmentFrame,
		runBootstrap:           runBootstrapRecurrenceStep,
		runDigestPrologue:      runRawSourceDigestPrologue,
		runDigest:              runRawSourceDigestKernel,
		runOptional:            runRawSourceOptionalKernel,
		runEven:                runRawSourceEvenKernel,
		runOdd:                 runRawSourceOddKernel,
		runTail:                runRawSourceTailKernel,
		runWriter:              runRawSourceWriterKernel,
		runOuter:               runRawSourceOuterIteration,
		runCycle:               func(s *state, _ fault) error { return runCompactKernel(s, compact26A428Cycle) },
		runCycleFinal:          runRawSourceCycleFinal,
		onceTokens:             []uint64{0x1010746B0, 0x101075B48, 0x101081348, 0x101081358},
		sourceAddress:          0x7000000FB304, frameAddress: 0x7000000FB724,
		outerStart: 0x1005B7438, outerCycle: 0x1007822CC, cycleFinal: 0x10046119C, optionalStart: 0x1005BA520,
	}
}

func executionProfileForModel(value *model) (*executionProfile, error) {
	if value == nil {
		return nil, fmt.Errorf("recovery model is unavailable")
	}
	if value.profile != nil {
		if value.profile.binaryHash != value.binaryHash {
			return nil, fmt.Errorf("recovery execution profile does not match model binary")
		}
		return value.profile, nil
	}
	if value.binaryHash == rawSourceBinaryHash {
		value.profile = embedded26A428ExecutionProfile()
		return value.profile, nil
	}
	return nil, fmt.Errorf("no recovery execution profile for model binary")
}
