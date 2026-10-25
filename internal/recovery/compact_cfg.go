package recovery

import "fmt"

type compactBlockTerminator uint8

const (
	compactBlockDirect compactBlockTerminator = iota
	compactBlockConditional
	compactBlockRegister
	compactBlockCallRegister
	compactBlockReturn
)

type compactCondition uint8

const (
	compactConditionEQ compactCondition = iota
	compactConditionNE
	compactConditionHS
	compactConditionLO
	compactConditionMI
	compactConditionPL
	compactConditionVS
	compactConditionVC
	compactConditionHI
	compactConditionLS
	compactConditionGE
	compactConditionLT
	compactConditionGT
	compactConditionLE
)

type compactBlock struct {
	start        uint64
	instructions []compactInstruction
	terminator   compactBlockTerminator
	condition    compactCondition
	register     uint8
	target       uint64
	alternate    uint64
}

func compactConditionMatches(s *state, condition compactCondition) (bool, error) {
	if condition > compactConditionLE {
		return false, fmt.Errorf("unknown compact block condition %d", condition)
	}
	name := [...]string{"eq", "ne", "hs", "lo", "mi", "pl", "vs", "vc", "hi", "ls", "ge", "lt", "gt", "le"}[condition]
	return s.condition(name), nil
}

func runCompactBlock(s *state, runtime *rawSourceRuntime, block compactBlock) (next uint64, stopped bool, err error) {
	terminalPC := block.start + uint64(len(block.instructions))*4
	if runtime != nil && runtime.stopPC >= block.start && runtime.stopPC <= terminalPC && (runtime.stopPC-block.start)%4 == 0 {
		instructionCount := int((runtime.stopPC - block.start) / 4)
		if instructionCount != 0 {
			prefix := compactKernel{start: block.start, instructions: block.instructions[:instructionCount]}
			if err := runCompactKernel(s, prefix); err != nil {
				return 0, false, err
			}
		}
		return runtime.stopPC, true, nil
	}
	if err := runCompactKernel(s, compactKernel{start: block.start, instructions: block.instructions}); err != nil {
		return 0, false, err
	}
	switch block.terminator {
	case compactBlockDirect:
		return block.target, false, nil
	case compactBlockConditional:
		matches, err := compactConditionMatches(s, block.condition)
		if err != nil {
			return 0, false, err
		}
		if matches {
			return block.target, false, nil
		}
		return block.alternate, false, nil
	case compactBlockRegister:
		return s.x[block.register], false, nil
	case compactBlockCallRegister:
		s.x[30] = terminalPC + 4
		return s.x[block.register], false, nil
	case compactBlockReturn:
		return s.x[30], false, nil
	default:
		return 0, false, fmt.Errorf("unknown compact block terminator %d", block.terminator)
	}
}
