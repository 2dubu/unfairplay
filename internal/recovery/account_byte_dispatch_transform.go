package recovery

import "fmt"

const (
	accountByteDispatch            = uint64(0x10045ad70)
	accountByteCopyDispatch        = uint64(0x1006d1874)
	accountIndexedWordDispatch     = uint64(0x1008e6fa4)
	accountIndexedWordDispatchTail = uint64(0x1008e7000)
)

func runAccountByteDispatch(s *state) (uint64, error) {
	s.traceBlock(accountByteDispatch)
	s.x[10] = uint64(uint32(s.x[11]))
	s.x[15] = s.x[20] + s.x[10]
	value, err := s.load(s.x[15], 1)
	if err != nil {
		return 0, fmt.Errorf("account byte-dispatch input: %w", err)
	}
	s.x[15] = uint64(byte(value))
	s.compare32(uint32(s.x[15]), uint32(byte(s.x[13])))
	s.x[15] = uint64(uint32(boolValue(s.condition("eq"))))
	s.x[11] = uint64(uint32(s.x[11]) + 1)
	s.x[16] = uint64(uint32(s.x[15]) << 4)
	s.x[15] = uint64(uint32(s.x[16]) - uint32(s.x[15]<<1))
	s.x[15] = uint64(uint32(s.x[15]) ^ uint32(s.x[12]))
	value, err = s.load(s.sp+0xdf0, 8)
	if err != nil {
		return 0, fmt.Errorf("account byte-dispatch table: %w", err)
	}
	s.x[16] = value
	next, err := loadBootstrapPointer(s, s.x[16], uint32(s.x[15]))
	if err != nil {
		return 0, fmt.Errorf("account byte-dispatch target: %w", err)
	}
	s.x[15] = next
	s.x[17] = 0xc82
	return next, nil
}

func runAccountByteCopyDispatch(s *state) (uint64, error) {
	s.traceBlock(accountByteCopyDispatch)
	s.x[17] = uint64(uint32(s.x[13]) + uint32(s.x[15]))
	s.x[22] = uint64(uint32(s.x[17]) | uint32(s.x[16]))
	s.x[17] = uint64(uint32(s.x[13]) + uint32(s.x[14]))
	s.x[23] = uint64(uint32(s.x[17]) | uint32(s.x[0]))
	value, err := s.load(s.x[24]+uint64(uint32(s.x[19])), 1)
	if err != nil {
		return 0, fmt.Errorf("account byte-copy input: %w", err)
	}
	s.x[17] = uint64(byte(value))
	if err := s.storeValue(s.x[20]+s.x[10], uint64(byte(s.x[17])), 1); err != nil {
		return 0, fmt.Errorf("account byte-copy output: %w", err)
	}
	s.x[25] = uint64(uint32(s.x[23]) - 0x14)
	s.x[25] &= 0xff
	s.compare32(uint32(s.x[17]), uint32(s.x[25]))
	s.x[17] = uint64(uint32(boolValue(s.condition("eq"))))
	s.x[25] = uint64(uint32(s.x[22]) ^ uint32(s.x[2]))
	s.x[17] = uint64(uint32(s.x[17]) * uint32(s.x[25]))
	s.x[13] = uint64(uint32(s.x[13]) + uint32(s.x[1]))
	s.x[13] = uint64(uint32(s.x[17]) ^ uint32(s.x[13]))
	value, err = s.load(s.sp+0xdf0, 8)
	if err != nil {
		return 0, fmt.Errorf("account byte-copy table: %w", err)
	}
	s.x[17] = value
	next, err := loadBootstrapPointer(s, s.x[17], uint32(s.x[13]))
	if err != nil {
		return 0, fmt.Errorf("account byte-copy target: %w", err)
	}
	s.x[13] = next
	s.x[17] = 0x9a49
	return next, nil
}

func runAccountIndexedWordDispatch(s *state) (uint64, error) {
	s.traceBlock(accountIndexedWordDispatch)
	s.x[11] = 0x77b16a0c
	s.compare32(uint32(s.x[8]), uint32(s.x[11]))
	s.x[11] = 0xffffe51b
	s.x[16] = 0x4561
	s.x[11] = uint64(uint32(s.x[9]) + uint32(s.x[11]))
	s.x[11] = uint64(uint32(s.x[11]) | uint32(s.x[16]))
	s.x[16] = uint64(uint32(s.x[15]) + 0xde3d1312)
	s.x[16] = uint64(uint32(s.x[16]) + uint32(s.x[11]))
	s.x[17] = s.sp + 0x2d0
	s.x[1] = 0x224ea33b
	if err := s.storeValue(s.x[17]+int32Offset(s.x[16], 2), uint64(uint32(s.x[1])), 4); err != nil {
		return 0, fmt.Errorf("account indexed-word dispatch output: %w", err)
	}
	s.x[16] = uint64(uint32(boolValue(s.condition("gt"))))
	s.x[17] = uint64(uint32(s.x[15]) + 1)
	s.x[1] = 0xe68c1e77
	s.x[2] = 0x884e95f3
	s.x[15] = uint64(uint32(s.x[15]) + uint32(s.x[1]))
	s.compare32(uint32(s.x[15]), uint32(s.x[2]))

	s.traceBlock(accountIndexedWordDispatchTail)
	s.x[1] = uint64(uint32(boolValue(s.condition("lt"))))
	s.x[16] = uint64(uint32(s.x[16]) ^ uint32(s.x[1]))
	s.compare32(uint32(s.x[15]), uint32(s.x[14]))
	s.x[15] = uint64(uint32(boolValue(s.condition("gt"))))
	s.compare32(uint32(s.x[16]), 0)
	if s.condition("ne") {
		s.x[15] = uint64(uint32(s.x[1]))
	}
	s.x[16] = 0x1d4
	s.x[15] = uint64(uint32(s.x[15]) * uint32(s.x[16]))
	s.x[15] = uint64(uint32(s.x[15]) ^ uint32(s.x[9]))
	next, err := loadBootstrapPointer(s, s.x[25], uint32(s.x[15]))
	if err != nil {
		return 0, fmt.Errorf("account indexed-word dispatch target: %w", err)
	}
	s.x[1] = next
	s.x[15] = s.x[17]
	s.x[17] = 0x88c4
	return next, nil
}
