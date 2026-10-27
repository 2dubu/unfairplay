package recovery

import "fmt"

const (
	bootstrapVectorLoop = uint64(0x1005df2ec)
	bootstrapVectorExit = uint64(0x1005df370)
)

func runBootstrapVectorLoop(s *state) (uint64, error) {
	for iteration := 0; iteration < 0x40; iteration++ {
		s.traceBlock(bootstrapVectorLoop)
		s.x[10] = s.x[11] + s.x[9]
		if err := s.loadVector(17, s.x[10], 16); err != nil {
			return 0, fmt.Errorf("bootstrap vector-loop input: %w", err)
		}
		s.vectorUnsignedShiftLong(18, 17, 8, 8, 0, true)
		s.vectorUnsignedShiftLong(17, 17, 8, 8, 0, false)
		s.vectorUnsignedMultiplyLong(19, 17, 5, 4, 16, false)
		s.vectorUnsignedMultiplyLong(17, 17, 6, 4, 16, true)
		s.vectorUnsignedMultiplyLong(20, 18, 5, 4, 16, false)
		s.vectorUnsignedMultiplyLong(18, 18, 6, 4, 16, true)
		for lane := 0; lane < 4; lane++ {
			s.vectorBinary("add", 21+lane, lane, 16, 4, 32)
			s.vectorBinary("sub", 25+lane, lane, 16, 4, 32)
		}
		for lane := 0; lane < 4; lane++ {
			s.vectorBinary("mul", 24-lane, 24-lane, 28-lane, 4, 32)
		}
		s.vectorEor3(28, 21, 16, 18)
		s.vectorEor3(27, 22, 16, 20)
		s.vectorEor3(26, 23, 16, 17)
		s.vectorEor3(25, 24, 16, 19)
		s.vectorTable(17, []int{25, 26, 27, 28}, 4)
		if err := s.storeVector(17, s.x[10], 16); err != nil {
			return 0, fmt.Errorf("bootstrap vector-loop output: %w", err)
		}
		s.x[9] += 0x10
		for register := 3; register >= 0; register-- {
			s.vectorBinary("add", register, register, 7, 4, 32)
		}
		s.compare64(s.x[9], 0x220)
		if s.condition("eq") {
			return bootstrapVectorExit, nil
		}
	}
	return 0, fmt.Errorf("bootstrap vector loop did not reach its natural exit")
}
