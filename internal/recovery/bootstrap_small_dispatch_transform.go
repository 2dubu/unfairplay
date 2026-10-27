package recovery

import "fmt"

const (
	bootstrapArithmeticDispatch = uint64(0x1001af428)
	bootstrapRangeStart         = uint64(0x100260bb4)
	bootstrapRangeTail          = uint64(0x100260c00)
	bootstrapStoreSelector      = uint64(0x1002acab4)
	bootstrapPreparedSelector   = uint64(0x100016d64)
	bootstrapBooleanFront       = uint64(0x1000623fc)
	bootstrapBooleanTail        = uint64(0x100062400)
	bootstrapPointerTailB       = uint64(0x100091158)
	bootstrapRegisterSetup      = uint64(0x100091f88)
	bootstrapRegisterSetupTail  = uint64(0x100092000)
)

func runBootstrapArithmeticDispatch(s *state) (uint64, error) {
	s.traceBlock(bootstrapArithmeticDispatch)
	s.x[11] = uint64(uint32(s.x[10]) ^ 0x744b)
	base, err := s.load(s.x[19]+8, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap arithmetic base: %w", err)
	}
	s.x[12] = base
	word, err := s.load(addOffset(base, int64(int32(s.x[9]))*4), 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap arithmetic word: %w", err)
	}
	s.x[12] = uint64(uint32(word))
	s.x[13] = 0x224ea33b
	s.compare32(uint32(s.x[12]), uint32(s.x[13]))
	gate := uint32(boolValue(s.condition("eq")))
	s.x[12] = uint64(gate)
	s.x[13] = uint64(uint32(s.x[10]) ^ 0x81d822fd)
	s.x[9] = uint64(uint32(s.x[13]) + uint32(s.x[9]))
	s.x[9] = uint64(uint32(s.x[9]) - gate)
	s.x[12] = 0x7e27fb6c
	s.x[13] = uint64(uint32(s.x[8]) + 0xb0036ebc)
	s.x[9] = uint64(uint32(s.x[9]) * uint32(s.x[13]))
	s.x[8] = uint64(uint32(s.x[8])*uint32(s.x[12]) + uint32(s.x[9]))
	s.x[9] = 0x817e7dbd
	s.x[8] = uint64(uint32(s.x[8]) + uint32(s.x[9]))
	next, err := loadBootstrapPointer(s, s.x[20], uint32(s.x[10]))
	if err != nil {
		return 0, fmt.Errorf("bootstrap arithmetic target: %w", err)
	}
	s.x[9] = next
	s.x[17] = 0x616
	return next, nil
}

func runBootstrapRange(s *state) (uint64, error) {
	s.traceBlock(bootstrapRangeStart)
	s.x[8] = 0xf7b5ffff
	s.x[8] = uint64(uint32(s.x[21]) ^ uint32(s.x[8]))
	s.x[11] = 0xef6bfffe
	s.x[10] = uint64(uint32(s.x[10]) & uint32(s.x[11]))
	s.x[11] = 0x6bdaf7f0
	s.x[8] = uint64(uint32(s.x[8]) + uint32(s.x[11]))
	s.x[10] = uint64(uint32(s.x[8]) + uint32(s.x[10]))
	word, err := s.load(s.sp+0xc, 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap range stack word: %w", err)
	}
	s.x[8] = uint64(uint32(word))
	s.x[11] = uint64(uint32(s.x[12]) + 0x3eb)
	s.x[14] = uint64(uint32(s.x[12]) + 0xf81)
	s.x[15] = 0xc68
	s.x[1] = uint64(uint32(s.x[14]) | uint32(s.x[15]))
	s.x[9] = uint64(uint32(s.x[21]) + uint32(s.x[9]))
	s.x[14] = uint64(uint32(s.x[10]) + 0x9c6f0810)
	s.x[15] = 0xffffb307

	s.traceBlock(bootstrapRangeTail)
	s.x[15] = uint64(uint32(s.x[1]) + uint32(s.x[15]))
	s.compare32(uint32(s.x[14]), uint32(s.x[15]))
	gate := uint32(boolValue(s.condition("le")))
	s.x[14] = uint64(gate * 0x2b)
	s.x[12] = uint64(uint32(s.x[14]) ^ uint32(s.x[12]))
	next, err := loadBootstrapPointer(s, s.x[20], uint32(s.x[12]))
	if err != nil {
		return 0, fmt.Errorf("bootstrap range target: %w", err)
	}
	s.x[14] = next
	s.x[12] = 0x224ea33b
	s.x[0] = s.x[9]
	s.x[21] = 0x224ea33b
	s.x[17] = 0x390b
	return next, nil
}

func runBootstrapStoreSelector(s *state) (uint64, error) {
	s.traceBlock(bootstrapStoreSelector)
	s.x[10] = uint64(uint32(s.x[11]) - 0xe92)
	base, err := s.load(s.x[19]+8, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap store-selector base: %w", err)
	}
	s.x[12] = base
	if err := s.storeValue(addOffset(base, int64(int32(s.x[0]))*4), uint64(uint32(s.x[21])), 4); err != nil {
		return 0, fmt.Errorf("bootstrap store-selector output: %w", err)
	}
	s.x[12] = uint64(uint32(s.x[0]) - 1)
	s.x[13] = uint64(uint32(s.x[11]) + 0xffffcb05)
	s.compare32(uint32(s.x[12]), uint32(s.x[13]))
	gate := uint32(boolValue(s.condition("gt")))
	s.x[13] = uint64(gate)
	s.x[14] = uint64(gate << 9)
	s.x[13] = uint64(uint32(s.x[14]) - gate<<3)
	s.x[13] = uint64(uint32(s.x[13]) ^ uint32(s.x[11]))
	next, err := loadBootstrapPointer(s, s.x[20], uint32(s.x[13]))
	if err != nil {
		return 0, fmt.Errorf("bootstrap store-selector target: %w", err)
	}
	s.x[13] = next
	s.x[17] = 0xe475
	return next, nil
}

func runBootstrapPreparedSelector(s *state) (uint64, error) {
	s.traceBlock(bootstrapPreparedSelector)
	s.x[11] = uint64(uint32(s.x[15]) ^ 0x5053)
	s.x[13] = 0xc5
	s.x[11] = uint64(uint32(s.x[11]) * uint32(s.x[13]))
	base, err := s.load(s.x[9]+8, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap prepared-selector base: %w", err)
	}
	s.x[13] = base
	next, err := loadBootstrapPointer(s, s.x[20], uint32(s.x[15]))
	if err != nil {
		return 0, fmt.Errorf("bootstrap prepared-selector target: %w", err)
	}
	s.x[2] = next
	s.x[14] = 0xdc431dae
	s.x[15] = 0x224ea33b
	s.x[16] = 0x31ee
	s.x[17] = 0x8fbb60b0
	s.x[0] = 0x7044d7bf
	s.x[17] = uint64(uint32(s.x[11]) + uint32(s.x[17]))
	s.x[17] = uint64(uint32(s.x[17]) & uint32(s.x[0]))
	s.x[0] = 0xdc436cf0
	s.x[1] = uint64(uint32(s.x[17]) - 0x2d)
	s.x[17] = 0x9b9d
	return next, nil
}

func runBootstrapBoolean(s *state) (uint64, error) {
	s.traceBlock(bootstrapBooleanFront)
	s.x[12] = 0x6d9d
	s.traceBlock(bootstrapBooleanTail)
	s.x[12] = uint64(uint32(s.x[11]) ^ uint32(s.x[12]))
	s.x[14] = 0x1d3
	s.x[21] = uint64(uint32(s.x[12]) * uint32(s.x[14]))
	s.x[14] = uint64(uint32(s.x[13]) + uint32(s.x[16]))
	s.x[15] = 0x842a8d93
	s.x[14] = uint64(uint32(s.x[14]) * uint32(s.x[15]))
	s.x[13] = uint64(uint32(s.x[13])*uint32(s.x[16]) + uint32(s.x[14]))
	s.x[14] = 0x5a874269
	s.x[13] = uint64(uint32(s.x[13]) + uint32(s.x[14]))
	s.compare32(uint32(s.x[13]), 0)
	gate := uint32(boolValue(s.condition("lt")))
	s.x[13] = uint64(gate)
	s.x[13] = uint64(uint32(s.x[13]) ^ uint32(s.x[21]))
	s.x[13] &= 1
	s.x[14] = uint64(uint32(s.x[12]) + 0x53)
	s.x[13] = uint64(uint32(s.x[13]) * uint32(s.x[14]))
	s.x[11] = uint64(uint32(s.x[13]) ^ uint32(s.x[11]))
	next, err := loadBootstrapPointer(s, s.x[19], uint32(s.x[11]))
	if err != nil {
		return 0, fmt.Errorf("bootstrap boolean target: %w", err)
	}
	s.x[11] = next
	s.x[17] = 0x8311
	return next, nil
}

func runBootstrapPointerTailB(s *state) (uint64, error) {
	s.traceBlock(bootstrapPointerTailB)
	next, err := loadBootstrapPointer(s, s.x[19], uint32(s.x[21]))
	if err != nil {
		return 0, fmt.Errorf("bootstrap pointer tail B: %w", err)
	}
	s.x[8] = next
	s.x[17] = 0x6f9c
	return next, nil
}

func runBootstrapRegisterSetup(s *state) (uint64, error) {
	s.traceBlock(bootstrapRegisterSetup)
	s.x[10] = 0xf0662331
	s.x[12] = uint64(uint32(s.x[12]) + uint32(s.x[10]))
	s.x[10] = uint64(uint32(s.x[25]) ^ 0xbe7f)
	s.x[11] = 0x92d4c62b
	s.x[13] = 0x6d2b875f
	next, err := loadBootstrapPointer(s, s.x[20], uint32(s.x[10]))
	if err != nil {
		return 0, fmt.Errorf("bootstrap register-setup target: %w", err)
	}
	s.x[14] = next
	s.x[10] = uint64(uint32(s.x[10]) + uint32(s.x[11]))
	s.x[13] = uint64(uint32(s.x[10]) & uint32(s.x[13]))
	s.x[7] = 0x23bc930e
	s.x[11] = 1
	s.x[15] = 0xeb59b62c
	s.x[16] = 0xffffe530
	s.x[0] = 0xdc436cf2
	s.x[1] = 0x224ea33b
	s.x[2] = 0xffffffff
	s.x[3] = 0x449d4676
	s.x[4] = 0x224ecfbc
	s.x[5] = 0x479d231f
	s.x[6] = 0xe4
	return bootstrapRegisterSetupTail, nil
}
