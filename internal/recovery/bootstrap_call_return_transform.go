package recovery

import "fmt"

const (
	bootstrapDescriptorReturn   = uint64(0x1007b30e4)
	bootstrapCallContinuation   = uint64(0x10081db98)
	bootstrapNestedReturnStart  = uint64(0x10084fbd0)
	bootstrapNestedReturnTail   = uint64(0x10084fc00)
	bootstrapGuardedReturnStart = uint64(0x10090b3d0)
	bootstrapGuardedReturnFrame = uint64(0x10090b3e8)
	bootstrapGuardedReturnTail  = uint64(0x10090b400)
)

func runBootstrapDescriptorReturn(s *state) (uint64, error) {
	s.traceBlock(bootstrapDescriptorReturn)
	value, err := s.load(s.x[20], 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap descriptor-return input: %w", err)
	}
	s.x[8] = uint64(uint32(value) - uint32(s.x[21]))
	s.x[8] = (s.x[8] & 0xffff) << 4
	if err := s.storeValue(s.x[19], s.x[8], 8); err != nil {
		return 0, fmt.Errorf("bootstrap descriptor-return output: %w", err)
	}
	s.sp = s.x[29] - 0x20
	for _, pair := range [...]struct {
		first, second int
		offset        uint64
	}{{29, 30, 0x20}, {20, 19, 0x10}, {22, 21, 0}} {
		first, err := s.load(s.sp+pair.offset, 8)
		if err != nil {
			return 0, fmt.Errorf("bootstrap descriptor-return register %d: %w", pair.first, err)
		}
		second, err := s.load(s.sp+pair.offset+8, 8)
		if err != nil {
			return 0, fmt.Errorf("bootstrap descriptor-return register %d: %w", pair.second, err)
		}
		s.x[pair.first], s.x[pair.second] = first, second
	}
	s.sp += 0x30
	return s.x[30], nil
}

func runBootstrapCallContinuation(s *state) (uint64, error) {
	s.traceBlock(bootstrapCallContinuation)
	next, err := loadBootstrapPointer(s, s.x[25], uint32(s.x[19]))
	if err != nil {
		return 0, fmt.Errorf("bootstrap call continuation: %w", err)
	}
	s.x[8] = next
	s.x[17] = 0x6864
	return next, nil
}

func runBootstrapNestedReturn(s *state) (uint64, error) {
	s.traceBlock(bootstrapNestedReturnStart)
	value, err := s.load(addOffset(s.x[29], -0x70), 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap nested-return status: %w", err)
	}
	s.x[8] = uint64(uint32(value))
	s.x[9] = 0x5eec043c
	s.compare32(uint32(s.x[8]), uint32(s.x[9]))
	matched := uint32(boolValue(s.condition("eq")))
	s.x[8] = uint64(matched)
	s.x[9] = uint64(uint32(s.x[23]) + 0xffff85bf)
	s.x[9] = uint64(uint32(s.x[9]) ^ uint32(s.x[21]))
	s.x[8] = uint64(matched*uint32(s.x[9]) ^ uint32(s.x[19]))
	next, err := loadBootstrapPointer(s, s.x[25], uint32(s.x[8]))
	if err != nil {
		return 0, fmt.Errorf("bootstrap nested-return target: %w", err)
	}
	s.x[8] = next
	value, err = s.load(s.x[20]+0x3ac8, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap nested-return first base: %w", err)
	}
	s.x[9] = value
	s.traceBlock(bootstrapNestedReturnTail)
	s.x[10] = 0xffffffffc10da2f9
	s.x[27] = s.x[9] + s.x[10]
	value, err = s.load(s.x[20]+0x12f8, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap nested-return second base: %w", err)
	}
	s.x[9] = value
	s.x[10] = 0xffffffff9f79fd81
	s.x[28] = s.x[9] + s.x[10]
	s.x[20] = 1
	s.x[17] = 0xfbd4
	return next, nil
}

func runBootstrapGuardedReturn(s *state) (uint64, error) {
	s.traceBlock(bootstrapGuardedReturnStart)
	value, err := s.load(addOffset(s.x[29], -0x60), 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap guarded-return saved value: %w", err)
	}
	s.x[8] = value
	const guardRoot = uint64(0x10106d000)
	value, err = s.load(guardRoot+0x408, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap guarded-return guard pointer: %w", err)
	}
	s.x[9] = value
	value, err = s.load(value, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap guarded-return guard value: %w", err)
	}
	s.x[9] = value
	s.compare64(s.x[9], s.x[8])
	if s.condition("ne") {
		return 0x10090b408, nil
	}
	s.traceBlock(bootstrapGuardedReturnFrame)
	s.sp += 0x500
	for _, pair := range [...]struct {
		first, second int
		offset        uint64
	}{{29, 30, 0x50}, {20, 19, 0x40}, {22, 21, 0x30}, {24, 23, 0x20}, {26, 25, 0x10}} {
		first, err := s.load(s.sp+pair.offset, 8)
		if err != nil {
			return 0, fmt.Errorf("bootstrap guarded-return register %d: %w", pair.first, err)
		}
		second, err := s.load(s.sp+pair.offset+8, 8)
		if err != nil {
			return 0, fmt.Errorf("bootstrap guarded-return register %d: %w", pair.second, err)
		}
		s.x[pair.first], s.x[pair.second] = first, second
	}
	s.traceBlock(bootstrapGuardedReturnTail)
	first, err := s.load(s.sp, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap guarded-return register 28: %w", err)
	}
	second, err := s.load(s.sp+8, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap guarded-return register 27: %w", err)
	}
	s.x[28], s.x[27] = first, second
	s.sp += 0x60
	return s.x[30], nil
}
