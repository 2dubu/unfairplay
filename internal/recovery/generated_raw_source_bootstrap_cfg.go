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
		bootstrapEnd := profile.bootstrapStart + uint64(len(compact26A428BootstrapKernel.instructions))*4
		partialStop := runtime != nil && runtime.stopPC > profile.bootstrapStart && runtime.stopPC <= bootstrapEnd
		if pc == profile.bootstrapStart && !partialStop {
			if err := profile.runBootstrap(s, fault{}); err != nil {
				return pc, blockCount, err
			}
			pc = s.x[3]
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
