// Derived from the compact 26A428 selection CFG. The semantic frame dispatch
// below is maintained by hand; the remaining block data is generated.
package recovery

import "fmt"

func runRawSourceSelectionCFG(s *state, runtime *rawSourceRuntime, maxBlocks int) (uint64, int, error) {
	profile, err := executionProfileForModel(s.model)
	if err != nil {
		return 0, 0, err
	}
	pc := profile.selectionStart
	stopHits := 0
	frameCalls := 0
	for blockCount := 0; blockCount < maxBlocks; blockCount++ {
		if pc == profile.outerStart {
			iterations, err := runSemanticNestedOuterStage(s)
			if err != nil {
				return pc, blockCount, fmt.Errorf("semantic nested outer stage: %w", err)
			}
			if iterations != outerNaturalIterationCount || s.x[8] != profile.selectionStart {
				return pc, blockCount, fmt.Errorf("semantic nested outer stage reached an invalid exit")
			}
			pc = s.x[8]
			continue
		}
		if pc == profile.selectionFrameCall {
			s.traceBlock(pc)
			frameCalls++
			if frameCalls == 1 {
				// The first invocation materializes neighboring dispatcher state.
			} else if frameCalls == 2 {
				next, err := profile.applyFrame(s)
				if err != nil {
					return pc, blockCount, err
				}
				pc = next
				continue
			} else {
				return pc, blockCount, fmt.Errorf("unexpected enrollment frame invocation %d", frameCalls)
			}
		}
		if pc == profile.selectionStop {
			s.traceBlock(pc)
			stopHits++
			if stopHits == profile.selectionStopHits {
				return pc, blockCount, nil
			}
		}
		if runtime != nil && runtime.stopPC != 0 && pc == runtime.stopPC {
			return pc, blockCount, nil
		}
		if pc >= 0xf000000000000000 || pc == 0x720000000000 {
			s.traceBlock(pc)
			if runtime == nil {
				return pc, blockCount, fmt.Errorf("selection runtime is unavailable")
			}
			next, err := runtime.handle(s, pc)
			if err != nil {
				return pc, blockCount, err
			}
			pc = next
			continue
		}
		if runtime == nil || runtime.stopPC == 0 {
			next, handled, err := runSelectionSemanticBlock(s, pc)
			if err != nil {
				return pc, blockCount, err
			}
			if handled {
				pc = next
				continue
			}
		}
		if runtime == nil || runtime.stopPC == 0 {
			next, handled, err := runSharedEnrollmentBlock(s, pc)
			if err != nil {
				return pc, blockCount, err
			}
			if handled {
				pc = next
				continue
			}
		}
		if runtime == nil || runtime.stopPC == 0 {
			next, handled, err := runSelectionDirectBlock(s, pc)
			if err != nil {
				return pc, blockCount, err
			}
			if handled {
				pc = next
				continue
			}
			return pc, blockCount, fmt.Errorf("direct selection CFG reached unknown block %#x", pc)
		}
		// Compact blocks are retained only for explicit partial-instruction stops.
		if block, ok := recoveryBlock(s, pc, profile.selectionBlocks, profile.supplementalBlocks, profile.verifiedFallbackBlocks); ok {
			next, stopped, err := runCompactBlock(s, runtime, block)
			if err != nil {
				return pc, blockCount, fmt.Errorf("compact selection CFG %#x: %w", pc, err)
			}
			if stopped {
				return next, blockCount, nil
			}
			pc = next
			continue
		}
		return pc, blockCount, fmt.Errorf("selection CFG reached unknown block %#x", pc)
	}
	return pc, maxBlocks, fmt.Errorf("selection CFG exceeded %d blocks; pc=%#x", maxBlocks, pc)
}
