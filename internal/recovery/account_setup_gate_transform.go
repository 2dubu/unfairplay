package recovery

import "fmt"

const (
	accountByteSignedGate   = uint64(0x100105438)
	accountFrameTargetSetup = uint64(0x100129be8)
	accountFrameTargetTail  = uint64(0x100129c00)
	accountBitGateRoute     = uint64(0x100144a4c)
	accountDispatchSetup    = uint64(0x10015a84c)
)

func runAccountByteSignedGate(s *state) (uint64, error) {
	s.traceBlock(accountByteSignedGate)
	s.x[9] = 0x6c79
	s.x[17] = uint64(uint32(s.x[13]) ^ uint32(s.x[9]))
	s.x[1] = uint64(uint32(s.x[17]) + uint32(s.x[17]<<3))
	s.x[9] = 0xffffdc16
	s.x[17] = uint64(uint32(s.x[13]) + uint32(s.x[9]))
	s.compare32(uint32(s.x[16]), uint32(int32(int8(s.x[17]))))
	s.x[17] = uint64(uint32(boolValue(s.condition("lt"))))
	s.x[9] = 0xffffb8d9
	s.x[0] = uint64(uint32(s.x[1]) + uint32(s.x[9]))
	s.x[17] = uint64(uint32(s.x[17]) * uint32(s.x[0]))
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[13]))
	next, err := loadBootstrapPointer(s, s.x[20], uint32(s.x[17]))
	if err != nil {
		return 0, fmt.Errorf("account signed-gate target: %w", err)
	}
	s.x[0] = next
	s.x[17] = 0x9195
	return next, nil
}

func runAccountFrameTargetSetup(s *state) (uint64, error) {
	s.traceBlock(accountFrameTargetSetup)
	s.x[9] = 0x569157c4
	s.x[13] = uint64(uint32(s.x[19]) ^ uint32(s.x[9]))
	s.x[9] = 0xa96f22dd
	s.x[17] = uint64(uint32(s.x[19]) + uint32(s.x[9]))

	s.traceBlock(accountFrameTargetTail)
	next, err := loadBootstrapPointer(s, s.x[20], uint32(s.x[8]))
	if err != nil {
		return 0, fmt.Errorf("account frame-target first target: %w", err)
	}
	s.x[8] = next
	if err := s.storeValue(s.sp+0x40, s.x[8], 8); err != nil {
		return 0, fmt.Errorf("account frame-target saved target: %w", err)
	}
	s.x[8] = s.sp + 0x100
	s.x[8] += s.x[26]
	s.x[9] = 0x6ad4
	s.x[11] = s.x[13] ^ s.x[9]
	s.x[11] += s.x[26]
	if err := s.storeValue(s.sp+0x80, s.x[11], 8); err != nil {
		return 0, fmt.Errorf("account frame-target saved pointer: %w", err)
	}
	value, err := s.load(s.sp+0xe4, 4)
	if err != nil {
		return 0, fmt.Errorf("account frame-target first pair word: %w", err)
	}
	s.x[11] = uint64(uint32(value))
	value, err = s.load(s.sp+0xe8, 4)
	if err != nil {
		return 0, fmt.Errorf("account frame-target second pair word: %w", err)
	}
	s.x[19] = uint64(uint32(value))
	value, err = s.load(s.sp+0xec, 4)
	if err != nil {
		return 0, fmt.Errorf("account frame-target third pair word: %w", err)
	}
	s.x[15] = uint64(uint32(value))
	value, err = s.load(s.sp+0xf0, 4)
	if err != nil {
		return 0, fmt.Errorf("account frame-target fourth pair word: %w", err)
	}
	s.x[28] = uint64(uint32(value))
	value, err = s.load(s.sp+0xfc, 1)
	if err != nil {
		return 0, fmt.Errorf("account frame-target byte: %w", err)
	}
	s.x[16] = uint64(byte(value))
	s.x[0] = uint64(uint32(s.x[12]) | ^uint32(s.x[16]<<1))
	s.x[16] = uint64(uint32(s.x[16]) + uint32(s.x[0]))
	s.x[16] = uint64(uint32(s.x[16]) - 0x57)
	s.x[16] = uint64(uint32(int32(int8(s.x[16]))))
	next, err = loadBootstrapPointer(s, s.x[20], uint32(s.x[17]))
	if err != nil {
		return 0, fmt.Errorf("account frame-target target: %w", err)
	}
	s.x[0] = next
	s.x[17] = 0xcd51
	return next, nil
}

func runAccountBitGateRoute(s *state) (uint64, error) {
	s.traceBlock(accountBitGateRoute)
	s.x[9] = 0x56917b05
	s.x[0] = uint64(uint32(s.x[1]) ^ uint32(s.x[9]))
	_ = s.addWithCarry(uint64(uint32(s.x[16])), 0x56, false, 32)
	s.x[13] = uint64(uint32(boolValue(s.condition("eq"))))
	s.x[13] = uint64(uint32(s.x[0]) ^ uint32(s.x[13]))
	s.x[9] = 1
	s.x[13] = uint64(uint32(s.x[9]) &^ uint32(s.x[13]))
	s.x[9] = 0xd7
	s.x[13] = uint64(uint32(s.x[13]) * uint32(s.x[9]))
	s.x[13] = uint64(uint32(s.x[13]) ^ uint32(s.x[1]))
	next, err := loadBootstrapPointer(s, s.x[20], uint32(s.x[13]))
	if err != nil {
		return 0, fmt.Errorf("account bit-gate target: %w", err)
	}
	s.x[13] = next
	s.x[17] = 0x7c57
	return next, nil
}

func runAccountDispatchSetup(s *state) (uint64, error) {
	s.traceBlock(accountDispatchSetup)
	s.x[9] = 0xa96ee2ae
	s.x[13] = uint64(uint32(s.x[0]) + uint32(s.x[9]))
	value, err := s.load(s.sp+0xf4, 4)
	if err != nil {
		return 0, fmt.Errorf("account dispatch-setup stack word: %w", err)
	}
	s.x[6] = uint64(uint32(value))
	next, err := loadBootstrapPointer(s, s.x[20], uint32(s.x[13]))
	if err != nil {
		return 0, fmt.Errorf("account dispatch-setup target: %w", err)
	}
	s.x[1] = next
	s.x[9] = 0xa96f7764
	s.x[13] = uint64(uint32(s.x[0]) + uint32(s.x[9]))
	s.x[9] = 0xa96f7725
	s.x[16] = uint64(uint32(s.x[0]) + uint32(s.x[9]))
	s.x[0] = 0x1cc84893
	s.x[9] = 0xaf
	s.x[12] = 0xce757bc5
	s.x[5] = 0xffffffffe337b76e
	s.x[23] = 0x4e757bd5
	s.x[26] = 0x11e99fcc
	s.x[27] = 0xf2de28a2
	s.x[17] = 0xfb55
	return next, nil
}
