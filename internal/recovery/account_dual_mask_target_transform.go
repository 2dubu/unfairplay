package recovery

import "fmt"

const (
	accountDualMaskTarget     = uint64(0x10003a3a8)
	accountDualMaskTargetTail = uint64(0x10003a400)
)

func runAccountDualMaskTarget(s *state) (uint64, error) {
	s.traceBlock(accountDualMaskTarget)
	s.x[11] = 0xf54b0e4e
	s.x[11] = uint64(uint32(s.x[13]) + uint32(s.x[11]))
	s.x[12] = 0x0ab4bfde
	s.x[11] = uint64(uint32(s.x[11]) & uint32(s.x[12]))
	s.x[12] = 0x6dfd55dd
	s.x[12] = uint64(uint32(s.x[13]) + uint32(s.x[12]))
	s.x[15] = 0x9202ffbd
	s.x[12] = uint64(uint32(s.x[12]) & uint32(s.x[15]))
	s.x[14] = uint64(uint32(s.x[9]) - uint32(s.x[14]))
	s.x[15] = uint64(uint32(s.x[12]) + 0x5bc51ad1)
	s.compare32(uint32(s.x[14]), uint32(s.x[15]))
	s.x[14] = uint64(uint32(boolValue(s.condition("eq"))))
	s.x[15] = uint64(uint32(s.x[11]) ^ 0x2dc7)
	s.x[14] = uint64(uint32(s.x[14]) * uint32(s.x[15]))
	s.x[13] = uint64(uint32(s.x[14]) ^ uint32(s.x[13]))
	s.traceBlock(accountDualMaskTargetTail)
	next, err := loadBootstrapPointer(s, s.x[20], uint32(s.x[13]))
	if err != nil {
		return 0, fmt.Errorf("account dual-mask target: %w", err)
	}
	s.x[13] = next
	s.x[17] = 0x3c02
	return next, nil
}
