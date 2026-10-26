// Derived from the compact 26A428 selection CFG. The semantic frame dispatch
// below is maintained by hand; the remaining block data is generated.
package recovery

import "fmt"

func runRawSourceSelectionCFG(s *state, runtime *rawSourceRuntime, maxBlocks int) (uint64, int, error) {
	pc := uint64(0x1006b10a0)
	stopHits := 0
	frameCalls := 0
	for blockCount := 0; blockCount < maxBlocks; blockCount++ {
		if pc == 0x1002fbb24 {
			frameCalls++
			if frameCalls == 1 {
				// The first invocation materializes neighboring dispatcher state.
			} else if frameCalls == 2 {
				next, err := applyRecoveredEnrollmentFrame(s)
				if err != nil {
					return pc, blockCount, err
				}
				pc = next
				continue
			} else {
				return pc, blockCount, fmt.Errorf("unexpected enrollment frame invocation %d", frameCalls)
			}
		}
		if pc == 0x10008a644 {
			stopHits++
			if stopHits == 2 {
				return pc, blockCount, nil
			}
		}
		if runtime != nil && runtime.stopPC != 0 && pc == runtime.stopPC {
			return pc, blockCount, nil
		}
		if pc >= 0xf000000000000000 || pc == 0x720000000000 {
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
		if block, ok := compact26A428SelectionBlocks[pc]; ok {
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
		return pc, blockCount, fmt.Errorf("bootstrap CFG reached unknown block %#x", pc)
	}
	return pc, maxBlocks, fmt.Errorf("bootstrap CFG exceeded %d blocks", maxBlocks)
}
