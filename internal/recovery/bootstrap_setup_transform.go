package recovery

import "fmt"

const (
	bootstrapSetupStart    = uint64(0x100360398)
	bootstrapSetupDecode   = uint64(0x100360400)
	bootstrapSetupStore    = uint64(0x10013233c)
	bootstrapSetupSelect   = uint64(0x1000dfa34)
	bootstrapSetupDispatch = uint64(0x10023706c)
)

func runBootstrapSetupTransaction(s *state) (uint64, error) {
	next, err := runBootstrapSetupStart(s)
	if err != nil || next != bootstrapSetupDecode {
		return next, err
	}
	next, err = runBootstrapSetupDecode(s)
	if err != nil || next != bootstrapSetupStore {
		return next, err
	}
	next, err = runBootstrapSetupStore(s)
	if err != nil || next != bootstrapSetupSelect {
		return next, err
	}
	next, err = runBootstrapSetupSelect(s)
	if err != nil || next != bootstrapSetupDispatch {
		return next, err
	}
	return runBootstrapSetupDispatch(s)
}

func runBootstrapSetupStart(s *state) (uint64, error) {
	s.traceBlock(bootstrapSetupStart)
	s.sp -= 0x10
	s.x[8] = 0x101081000
	left, err := s.load(s.x[8]+0x270, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap setup left state: %w", err)
	}
	s.x[9] = left
	s.x[10] = 0x101075000
	right, err := s.load(s.x[10]+0xb40, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap setup right state: %w", err)
	}
	s.x[11] = right
	s.x[9] = right ^ left
	s.x[11] = s.sp + 8
	s.x[12] = 0x0a5cd4cf
	s.x[13] = 0x8c2fb9df71512e5a
	s.x[14] = 0x73d046208eaed1a6
	s.x[9] += s.x[11]
	s.x[11] = s.x[9] ^ s.x[14]
	s.x[13] = s.x[9]*s.x[12] + s.x[13]
	s.x[14] = s.x[11] * s.x[12]
	if err := s.storeValue(s.x[8]+0x270, s.x[13], 8); err != nil {
		return 0, fmt.Errorf("bootstrap setup left output: %w", err)
	}
	if err := s.storeValue(s.x[10]+0xb40, s.x[14], 8); err != nil {
		return 0, fmt.Errorf("bootstrap setup right output: %w", err)
	}
	s.x[8] = uint64(uint32(s.x[0]) | 0x240d4a03)
	return bootstrapSetupDecode, nil
}

func runBootstrapSetupDecode(s *state) (uint64, error) {
	s.traceBlock(bootstrapSetupDecode)
	s.x[9] = 0xdbf2b5fc
	s.x[10] = uint64(uint32(s.x[0]) | uint32(s.x[9]))
	s.x[8] = uint64(uint32(s.x[8]) - uint32(s.x[10]))
	s.x[8] = uint64(uint32(s.x[8]) + uint32(s.x[9]))
	s.x[9] = 0xddd9d377
	s.x[8] = uint64(uint32(s.x[8]) ^ uint32(s.x[9]))
	s.x[9] = 0xda51a827
	s.x[8] = uint64(uint32(s.x[8]) * uint32(s.x[9]))
	word, err := s.load(s.x[0]+0x10, 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap setup input word: %w", err)
	}
	s.x[9] = uint64(uint32(word))
	s.x[10] = uint64(uint32(s.x[9]) ^ uint32(s.x[8]))
	left, err := s.load(s.x[0], 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap setup input left: %w", err)
	}
	right, err := s.load(s.x[0]+8, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap setup input right: %w", err)
	}
	s.x[8], s.x[9] = left, right
	s.compare64(left, right)
	equal := uint32(boolValue(s.condition("eq")))
	s.x[11] = uint64(uint32(s.x[10]) + 0x952080c1)
	s.x[15] = uint64((equal ^ uint32(s.x[11])) & 1)
	s.x[16] = 0xec
	s.x[15] = uint64(uint32(s.x[15]) * uint32(s.x[16]))
	s.x[16] = 0xa6
	s.x[13] = uint64(uint32(s.x[14]) - uint32(s.x[13]))
	s.x[13] = uint64(uint32(s.x[13]) ^ uint32(s.x[16]))
	s.x[12] = uint64(uint32(s.x[13]) * uint32(s.x[12]))
	s.x[13] = 0x100f49000 + 0x830 + uint64(byte(s.x[12]))
	value, err := s.load(s.x[13], 1)
	if err != nil {
		return 0, fmt.Errorf("bootstrap setup first substitution: %w", err)
	}
	s.x[13] = uint64(uint32(value) ^ 0xb5)
	s.x[14] = 0x100fe9000 + 0x690
	value, err = s.load(s.x[14]+uint64(uint32(s.x[13])), 1)
	if err != nil {
		return 0, fmt.Errorf("bootstrap setup second substitution: %w", err)
	}
	s.x[13] = uint64(byte(value))
	s.x[14] = uint64(uint32(s.x[15]) ^ uint32(s.x[10]))
	s.x[10] = uint64(uint32(s.x[12]) ^ uint32(s.x[13]))
	s.x[10] &= 0xff
	s.x[10] = uint64(uint32(s.x[10]) + 0x228)
	s.x[12] = 0x10105b000 + 0x6f0
	pointer, err := s.load(s.x[12]+uint64(uint32(s.x[10]))*8, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap setup pointer table: %w", err)
	}
	s.x[10] = pointer + 0xffffffff98245355
	next, err := loadBootstrapPointer(s, s.x[10], uint32(s.x[14]))
	if err != nil {
		return 0, fmt.Errorf("bootstrap setup decode dispatch: %w", err)
	}
	s.x[12] = next
	s.x[17] = 0x2a4d
	return next, nil
}

func runBootstrapSetupStore(s *state) (uint64, error) {
	s.traceBlock(bootstrapSetupStore)
	s.x[12] = 0x6adf7dfd
	s.x[11] = uint64(uint32(s.x[11]) & uint32(s.x[12]))
	s.x[12] = uint64(uint32(s.x[11]) ^ 0x3312)
	s.x[13] = uint64(uint32(s.x[11]) + 0x716602c1)
	s.x[14] = 0x8e99fe7f
	s.x[13] = uint64(uint32(s.x[13]) & uint32(s.x[14]))
	word, err := s.load(s.x[8]+4, 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap setup stored word: %w", err)
	}
	s.x[15] = uint64(uint32(word))
	s.x[14] = 0x842a411a
	s.x[13] = uint64(uint32(s.x[13]) ^ uint32(s.x[14]))
	s.x[14] = uint64(uint32(s.x[15]) + uint32(s.x[12]))
	s.x[13] = s.addWithCarry(uint64(uint32(s.x[14])), uint64(uint32(s.x[13])), false, 32)
	s.x[14] = s.x[13]
	if s.condition("mi") {
		s.x[14] = uint64(uint32(0) - uint32(s.x[13]))
	}
	if err := s.storeValue(s.x[9]+4, uint64(uint32(s.x[15])), 4); err != nil {
		return 0, fmt.Errorf("bootstrap setup copied word: %w", err)
	}
	s.x[13] = uint64(uint32(boolValue(s.condition("eq"))) * 0x59)
	dispatchIndex := uint32(s.x[13]) ^ uint32(s.x[11])
	next, err := s.load(s.x[10]+uint64(dispatchIndex)*8, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap setup store dispatch: %w", err)
	}
	s.x[11] = next
	s.x[15] = 0x59
	s.x[17] = 0xbc07
	return next, nil
}

func runBootstrapSetupSelect(s *state) (uint64, error) {
	s.traceBlock(bootstrapSetupSelect)
	s.x[11] = uint64(uint32(s.x[12]) + 0x8eed90f0)
	s.x[13] = 0x7112b5fc
	s.x[11] = uint64(uint32(s.x[11]) & uint32(s.x[13]))
	s.x[13] = s.x[14] & 1
	s.x[15] = uint64(uint32(s.x[11]) ^ 0x9507)
	dispatchIndex := uint32(s.x[15])*uint32(s.x[13]) ^ uint32(s.x[12])
	next, err := s.load(s.x[10]+uint64(dispatchIndex)*8, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap setup select dispatch: %w", err)
	}
	s.x[13] = s.x[14]
	s.x[15] = next
	s.x[17] = 0xb24
	return next, nil
}

func runBootstrapSetupDispatch(s *state) (uint64, error) {
	s.traceBlock(bootstrapSetupDispatch)
	s.x[12] = uint64(uint32(s.x[11]) ^ 0x951d)
	s.compare32(uint32(s.x[14]), uint32(s.x[12]))
	selected := uint32(boolValue(s.condition("ne")))
	s.x[14] = uint64(selected << 2)
	dispatchIndex := ((selected << 2) | (selected << 5)) ^ uint32(s.x[11])
	next, err := loadBootstrapPointer(s, s.x[10], dispatchIndex)
	if err != nil {
		return 0, fmt.Errorf("bootstrap setup final dispatch: %w", err)
	}
	s.x[12] = next
	s.x[17] = 0xa668
	return next, nil
}
