package recovery

import "fmt"

const accountEntryByteGeneratorPass = uint64(0x100b1ddc0)

func runAccountEntrySemanticBlock(s *state, pc uint64) (uint64, bool, error) {
	if pc != accountEntryByteGeneratorPass {
		return pc, false, nil
	}
	next, err := runAccountEntryByteGenerator(s)
	return next, true, err
}

func runAccountEntryByteGenerator(s *state) (uint64, error) {
	const maximumBytes = 4096
	for index := 0; index < maximumBytes; index++ {
		next, err := runAccountEntryByteGeneratorStep(s)
		if err != nil || next != accountEntryByteGeneratorPass {
			return next, err
		}
	}
	return 0, fmt.Errorf("account-entry byte generator exceeded %d bytes", maximumBytes)
}

func runAccountEntryByteGeneratorStep(s *state) (uint64, error) {
	if s == nil {
		return 0, fmt.Errorf("account-entry byte-generator state is unavailable")
	}
	s.traceBlock(accountEntryByteGeneratorPass)
	s.x[0] = uint64(uint32(s.x[16]) ^ uint32(s.x[10]))
	s.x[15] = (s.x[0] ^ s.x[19]) + s.x[1]
	if err := s.storeValue(s.sp+0x318, s.x[15], 8); err != nil {
		return 0, fmt.Errorf("account-entry byte-generator position: %w", err)
	}
	s.x[17] = uint64(uint32(s.x[13]) ^ uint32(s.x[20]))
	s.x[1] = uint64(uint32(s.x[13]) ^ uint32(s.x[23]))
	s.x[1] = uint64(uint32(s.x[1]) + uint32(s.x[26]))
	s.x[2] = uint64(uint32(s.x[13]) ^ uint32(s.x[27]))
	s.x[3] = uint64((uint32(s.x[13]) ^ 0xfffffffb) + 0x8f)
	s.x[17] = uint64(uint32(s.x[17]) + uint32(s.x[21]))
	s.x[17] = uint64(uint32(s.x[13]) ^ uint32(s.x[17]))
	s.x[2] = uint64(uint32(s.x[2]) + uint32(s.x[28]))
	s.x[1] = uint64(uint32(s.x[1]) ^ uint32(s.x[2]))
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[1]))
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[3]))
	s.x[17] = uint64(uint32(s.x[17]) ^ 0x30)
	s.x[1] = uint64(uint32(s.x[13]) + 0x87)
	s.x[17] = uint64(uint32(s.x[17]) * uint32(s.x[1]))
	if err := s.storeValue(s.x[6]+s.x[15], uint64(byte(s.x[17])), 1); err != nil {
		return 0, fmt.Errorf("account-entry byte-generator output: %w", err)
	}
	s.x[13]++
	position, err := s.load(s.sp+0x318, 8)
	if err != nil {
		return 0, fmt.Errorf("account-entry byte-generator saved position: %w", err)
	}
	s.x[1] = position
	s.compare64(position, 0)
	dispatchIndex := uint32(boolValue(s.condition("eq")))*uint32(s.x[14]) ^ uint32(s.x[16])
	next, err := loadBootstrapPointer(s, s.x[30], dispatchIndex)
	if err != nil {
		return 0, fmt.Errorf("account-entry byte-generator dispatch: %w", err)
	}
	s.x[15] = next
	s.x[17] = 0x623c
	return next, nil
}
