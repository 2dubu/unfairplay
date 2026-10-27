package recovery

import "fmt"

const (
	accountStateDispatchStart = uint64(0x100377f2c)
	accountStateDispatchTail  = uint64(0x100378000)
)

func runAccountStateDispatch(s *state) (uint64, error) {
	s.traceBlock(accountStateDispatchStart)
	s.sp -= 0x10
	s.x[8] = 0x101074000
	first, err := s.load(s.x[8]+0x588, 8)
	if err != nil {
		return 0, fmt.Errorf("account state-dispatch first state: %w", err)
	}
	s.x[9] = first
	s.x[11] = bootstrapStateBase
	second, err := s.load(s.x[11]+0xb40, 8)
	if err != nil {
		return 0, fmt.Errorf("account state-dispatch second state: %w", err)
	}
	s.x[10] = second
	s.x[12] = s.sp + 8
	s.x[9] -= s.x[12]
	s.x[9] += s.x[10]
	s.x[10] = 0x0a5cd4cf
	s.x[12] = 0x673e3a80075f02c6
	s.x[12] = s.x[9]*s.x[10] + s.x[12]
	s.x[13] = 0x73d046208eaed1a6
	s.x[9] ^= s.x[13]
	s.x[15] = s.x[9] * s.x[10]
	if err := s.storeValue(s.x[8]+0x588, s.x[12], 8); err != nil {
		return 0, fmt.Errorf("account state-dispatch first state output: %w", err)
	}
	if err := s.storeValue(s.x[11]+0xb40, s.x[15], 8); err != nil {
		return 0, fmt.Errorf("account state-dispatch second state output: %w", err)
	}
	s.x[8] = uint64(uint32(s.x[0]) & 0x57e53a58)
	s.x[8] = uint64(uint32(s.x[8]) << 1)
	s.x[9] = 0x281ac5a7
	s.x[8] = uint64(uint32(s.x[8]) - uint32(s.x[0]))
	s.x[8] = uint64(uint32(s.x[8]) + uint32(s.x[9]))
	s.x[9] = 0xcf2817c5
	s.x[8] = uint64(uint32(s.x[8]) ^ uint32(s.x[9]))
	s.x[9] = 0xba209ef7
	s.x[11] = uint64(uint32(s.x[8]) * uint32(s.x[9]))
	value, err := s.load(s.x[0]+0x28, 8)
	if err != nil {
		return 0, fmt.Errorf("account state-dispatch first descriptor: %w", err)
	}
	s.x[9] = value
	value, err = s.load(s.x[0]+0x38, 8)
	if err != nil {
		return 0, fmt.Errorf("account state-dispatch second descriptor: %w", err)
	}
	s.x[8] = value
	word, err := s.load(value, 4)
	if err != nil {
		return 0, fmt.Errorf("account state-dispatch descriptor word: %w", err)
	}
	s.x[13] = uint64(uint32(word))
	s.x[14] = 0xcaa90d63
	s.x[16] = 0x94d8fb13
	s.x[9] = uint64(uint32(s.x[9]) + uint32(s.x[11]))
	s.x[16] = uint64(uint32(s.x[9]) + uint32(s.x[16]))
	s.x[17] = 0xf15d894a
	s.x[14] = uint64(uint32(s.x[13]) + uint32(s.x[14]))
	s.compare32(uint32(s.x[14]), uint32(s.x[17]))
	firstBound := uint32(boolValue(s.condition("lo")))
	s.x[1] = uint64(firstBound)
	s.compare32(uint32(s.x[16]), uint32(s.x[17]))
	secondBound := uint32(boolValue(s.condition("lo")))
	s.x[17] = uint64(secondBound)

	s.traceBlock(accountStateDispatchTail)
	s.compare32(uint32(s.x[14]), uint32(s.x[16]))
	word, err = s.load(s.x[0], 4)
	if err != nil {
		return 0, fmt.Errorf("account state-dispatch input word: %w", err)
	}
	s.x[14] = uint64(uint32(word))
	s.x[16] = uint64(secondBound ^ firstBound)
	ordered := uint32(boolValue(s.condition("hi")))
	s.x[17] = uint64(ordered)
	s.compare32(uint32(s.x[16]), 0)
	selected := ordered
	if s.condition("ne") {
		selected = firstBound
	}
	s.x[16] = uint64(selected)
	s.x[17] = 0xffffd5a4
	s.x[11] = uint64(uint32(s.x[14]) + uint32(s.x[11]))
	s.x[14] = uint64(uint32(s.x[11]) + uint32(s.x[17]))
	s.x[17] = 0x51a0
	s.x[14] = uint64(uint32(s.x[14]) | uint32(s.x[17]))
	s.x[17] = 0x53b8
	s.x[17] = uint64(uint32(s.x[14]) ^ uint32(s.x[17]))
	s.x[16] = uint64(uint32(s.x[16]) * uint32(s.x[17]))
	s.x[17] = 0xa6
	s.x[12] = uint64(uint32(s.x[12]) ^ uint32(s.x[17]))
	s.x[12] = uint64(uint32(s.x[12]) ^ uint32(s.x[15]))
	s.x[10] = uint64(uint32(s.x[12]) * uint32(s.x[10]))
	s.x[12] = 0x100fe9290 + uint64(byte(s.x[10]))
	value, err = s.load(s.x[12], 1)
	if err != nil {
		return 0, fmt.Errorf("account state-dispatch first byte: %w", err)
	}
	s.x[12] = uint64(byte(value)) ^ 0xf0
	s.x[15] = 0x100de29f0
	value, err = s.load(s.x[15]+uint64(uint32(s.x[12])), 1)
	if err != nil {
		return 0, fmt.Errorf("account state-dispatch second byte: %w", err)
	}
	s.x[12] = uint64(byte(value))
	s.x[10] = uint64(uint32(s.x[10]) ^ uint32(s.x[12]))
	s.x[10] &= 0xff
	s.x[10] = uint64(uint32(s.x[10]) + 0x244)
	s.x[12] = bootstrapTargetTable
	value, err = s.load(s.x[12]+uint64(uint32(s.x[10]))*8, 8)
	if err != nil {
		return 0, fmt.Errorf("account state-dispatch table: %w", err)
	}
	s.x[10] = value
	s.x[11] = uint64(uint32(s.x[16]) ^ uint32(s.x[11]))
	s.x[12] = 0xffffffff98245355
	s.x[10] += s.x[12]
	next, err := loadBootstrapPointer(s, s.x[10], uint32(s.x[11]))
	if err != nil {
		return 0, fmt.Errorf("account state-dispatch target: %w", err)
	}
	s.x[11] = next
	s.x[2] = 0x318a201b
	s.x[17] = 0xa291
	return next, nil
}
