package recovery

import "fmt"

const (
	accountScaledSelect      = uint64(0x100984c48)
	accountIndexedSelect     = uint64(0x10057e764)
	accountFrameSetup        = uint64(0x100654b58)
	accountNestedReturn      = uint64(0x1008550ec)
	accountShiftSelect       = uint64(0x100911a50)
	accountPairedSelectFront = uint64(0x100a697d0)
	accountPairedSelectTail  = uint64(0x100a69800)
	accountDifferenceSelect  = uint64(0x1008d6a98)
	accountByteWriterA       = uint64(0x1000386e4)
	accountByteWriterB       = uint64(0x1001f85f4)
	accountByteWriterCFront  = uint64(0x100240fd8)
	accountByteWriterCTail   = uint64(0x100241000)
)

func runAccountScaledSelect(s *state) (uint64, error) {
	s.traceBlock(accountScaledSelect)
	s.x[8] = uint64(uint32(s.x[11]) ^ uint32(s.x[5]))
	s.x[12] = uint64(uint32(s.x[8]) + uint32(s.x[8]<<2))
	s.x[8] = uint64(uint32(s.x[11]) + uint32(s.x[20]))
	s.x[13] = uint64(uint32(s.x[8]) | uint32(s.x[22]))
	s.x[8] = 0xfb5830fc
	s.x[14] = s.x[25] + s.x[8]
	s.x[8] = uint64(int64(int32(s.x[14])))
	s.x[14] = s.x[8] + (uint64(int64(int32(s.x[14]))) << 1)
	s.x[14] <<= 4
	word, err := s.load(s.x[19]+s.x[14], 4)
	if err != nil {
		return 0, fmt.Errorf("account scaled-select word: %w", err)
	}
	s.x[14] = uint64(uint32(word))
	s.x[15] = uint64(uint32(s.x[13]) + 0x095a6efc)
	s.compare32(uint32(s.x[14]), uint32(s.x[15]))
	gate := uint32(boolValue(s.condition("ne")))
	s.x[15] = uint64(gate)
	s.x[16] = uint64(uint32(s.x[12]) ^ 0x5a41)
	s.x[15] = uint64(gate * uint32(s.x[16]))
	s.x[11] = uint64(uint32(s.x[15]) ^ uint32(s.x[11]))
	next, err := loadBootstrapPointer(s, s.x[30], uint32(s.x[11]))
	if err != nil {
		return 0, fmt.Errorf("account scaled-select target: %w", err)
	}
	s.x[11] = next
	s.x[17] = 0xde44
	return next, nil
}

func runAccountIndexedSelect(s *state) (uint64, error) {
	s.traceBlock(accountIndexedSelect)
	s.x[10] = s.x[8] + (s.x[8] << 1)
	s.x[10] = s.x[19] + (s.x[10] << 4)
	word, err := s.load(s.x[10]+0x28, 4)
	if err != nil {
		return 0, fmt.Errorf("account indexed-select word: %w", err)
	}
	s.x[12] = uint64(uint32(word))
	s.x[13] = uint64(uint32(s.x[11]) ^ 0x7d5d)
	s.x[13] = uint64(uint32(s.x[13]) & uint32(s.x[12]))
	s.x[14] = uint64(uint32(s.x[11]) + 0xffff82a4)
	s.compare32(uint32(s.x[13]), uint32(s.x[14]))
	gate := uint32(boolValue(s.condition("eq")))
	s.x[13] = uint64(gate)
	s.x[14] = uint64(gate << 1)
	s.x[13] = uint64(uint32(s.x[14]) | gate<<5)
	s.x[13] = uint64(uint32(s.x[13]) ^ uint32(s.x[11]))
	next, err := loadBootstrapPointer(s, s.x[30], uint32(s.x[13]))
	if err != nil {
		return 0, fmt.Errorf("account indexed-select target: %w", err)
	}
	s.x[13] = next
	s.x[17] = 0xe57c
	return next, nil
}

func runAccountFrameSetup(s *state) (uint64, error) {
	s.traceBlock(accountFrameSetup)
	s.x[8] = uint64(uint32(s.x[28]) + 0x20ba)
	s.x[9] = 0xb1
	s.x[11] = uint64(uint32(s.x[8]) | uint32(s.x[9]))
	word, err := s.load(s.sp+0x324, 4)
	if err != nil {
		return 0, fmt.Errorf("account frame-setup word: %w", err)
	}
	s.x[9] = uint64(uint32(word))
	s.x[8] = uint64(uint32(s.x[9]) ^ 0xc84b81ee)
	s.x[10] = 0x37b47e12
	s.x[8] = uint64(uint32(s.x[8]) + uint32(s.x[10]))
	s.x[10] = uint64(uint32(s.x[9]) ^ 0x3787152c)
	s.x[10] = uint64(uint32(s.x[10]) ^ uint32(s.x[11]))
	s.x[12] = uint64(uint32(s.x[9]) ^ 0xfb2c16c0)
	s.x[13] = 0x04d3e940
	s.x[12] = uint64(uint32(s.x[12]) + uint32(s.x[13]))
	s.x[8] = uint64(uint32(s.x[12]) ^ uint32(s.x[8]))
	s.x[12] = 0xc8782b63
	s.x[10] = uint64(uint32(s.x[10]) + uint32(s.x[12]))
	s.x[8] = uint64(uint32(s.x[8]) ^ uint32(s.x[10]))
	s.x[10] = uint64(uint32(s.x[8]) + 0xdcde939d)
	next, err := loadBootstrapPointer(s, s.x[30], uint32(s.x[28]))
	if err != nil {
		return 0, fmt.Errorf("account frame-setup target: %w", err)
	}
	s.x[8] = next
	value, err := s.load(s.sp+0x200, 8)
	if err != nil {
		return 0, fmt.Errorf("account frame-setup pointer: %w", err)
	}
	s.x[7] = value
	s.x[0] = 1
	s.x[5] = 0xd048
	s.x[20] = 0xffff505f
	s.x[22] = 0x5c5
	word, err = s.load(s.sp+0x1c0, 4)
	if err != nil {
		return 0, fmt.Errorf("account frame-setup selector: %w", err)
	}
	s.x[28] = uint64(uint32(word))
	s.x[17] = 0xd35c
	return next, nil
}

func runAccountNestedReturn(s *state) (uint64, error) {
	s.traceBlock(accountNestedReturn)
	s.x[6] = s.x[26]
	s.x[1] = 0x04a7cf03
	s.x[4] = 0x2b7ffafa
	s.x[3] = 0xce757bc5
	s.x[2] = 0x43
	value, err := s.load(s.sp+0x208, 8)
	if err != nil {
		return 0, fmt.Errorf("account nested-return saved link: %w", err)
	}
	s.x[30] = value
	word, err := s.load(addOffset(s.x[29], -0xa0), 4)
	if err != nil {
		return 0, fmt.Errorf("account nested-return status: %w", err)
	}
	s.x[5] = uint64(uint32(word))
	s.x[8] = uint64(uint32(s.x[20]) + 0xce754318)
	s.compare32(uint32(s.x[5]), uint32(s.x[8]))
	gate := uint32(boolValue(s.condition("eq")))
	s.x[8] = uint64(gate * 0x82d5)
	s.x[8] = uint64(uint32(s.x[8]) ^ uint32(s.x[20]))
	next, err := loadBootstrapPointer(s, s.x[30], uint32(s.x[8]))
	if err != nil {
		return 0, fmt.Errorf("account nested-return target: %w", err)
	}
	s.x[8] = next
	s.x[17] = 0x8231
	return next, nil
}

func runAccountShiftSelect(s *state) (uint64, error) {
	s.traceBlock(accountShiftSelect)
	s.x[9] = uint64(uint32(s.x[13]) ^ 0x04a7ca24)
	s.x[10] = 0xb
	s.x[11] = uint64(uint32(s.x[9]) * uint32(s.x[10]))
	s.x[9] = uint64(uint32(s.x[13]) ^ 0x04a7af0c)
	s.x[8] = s.x[8] + (s.x[8] << 1)
	s.x[8] = s.x[19] + (s.x[8] << 4)
	word, err := s.load(s.x[8]+0x20, 4)
	if err != nil {
		return 0, fmt.Errorf("account shift-select word: %w", err)
	}
	s.x[10] = uint64(uint32(word))
	s.x[12] = uint64(uint32(s.x[10]) - uint32(s.x[13]))
	s.x[10] = uint64(uint32(s.x[13]) - uint32(s.x[10]))
	s.x[10] = uint64(uint32(s.x[12]) | uint32(s.x[10]))
	s.x[12] = uint64(uint32(s.x[9]) + 0xffffa010)
	s.x[10] = uint64(uint32(shiftRight(s.x[10], int(s.x[12]), 32)))
	s.x[10] = uint64(uint32(s.x[11]) ^ uint32(s.x[10]))
	s.x[10] &= 1
	s.x[12] = 0x33
	s.x[10] = uint64(uint32(s.x[10]) * uint32(s.x[12]))
	s.x[12] = uint64(uint32(s.x[13]) ^ 0x04a760dd)
	s.x[10] = uint64(uint32(s.x[12]) ^ uint32(s.x[10]))
	next, err := loadBootstrapPointer(s, s.x[30], uint32(s.x[10]))
	if err != nil {
		return 0, fmt.Errorf("account shift-select target: %w", err)
	}
	s.x[10] = next
	s.x[17] = 0x66e9
	return next, nil
}

func runAccountPairedSelect(s *state) (uint64, error) {
	s.traceBlock(accountPairedSelectFront)
	s.x[13] = 0x5305
	s.x[11] = uint64(uint32(s.x[11]) ^ uint32(s.x[13]))
	s.x[13] = uint64(uint32(s.x[11]) + 0x551d4b94)
	s.x[14] = 0xaae2bebf
	s.x[28] = uint64(uint32(s.x[13]) & uint32(s.x[14]))
	s.x[13] = uint64(uint32(s.x[11]) + 0x04a7a0aa)
	s.x[14] = 0xffffd1a9
	s.traceBlock(accountPairedSelectTail)
	s.x[14] = uint64(uint32(s.x[11]) + uint32(s.x[14]))
	s.x[12] = uint64(uint32(s.x[14]) & uint32(s.x[12]))
	s.x[14] = uint64(uint32(s.x[28]) + 0xffffc753)
	s.compare32(uint32(s.x[12]), uint32(s.x[14]))
	gate := uint32(boolValue(s.condition("eq")))
	s.x[12] = uint64(gate * 0x492b)
	s.x[11] = uint64(uint32(s.x[12]) ^ uint32(s.x[11]))
	next, err := loadBootstrapPointer(s, s.x[30], uint32(s.x[11]))
	if err != nil {
		return 0, fmt.Errorf("account paired-select target: %w", err)
	}
	s.x[11] = next
	s.x[17] = 0xb6ab
	return next, nil
}

func runAccountDifferenceSelect(s *state) (uint64, error) {
	s.traceBlock(accountDifferenceSelect)
	s.x[12] = uint64(uint32(s.x[13]) ^ 0x3fa)
	s.x[11] = uint64(uint32(s.x[13]) ^ 0x175b)
	s.x[15] = 0xe2
	s.x[11] = uint64(uint32(s.x[11]) * uint32(s.x[15]))
	s.x[15] = uint64(uint32(s.x[12]) << 3)
	s.x[14] = uint64(uint32(s.x[10]) - uint32(s.x[14]))
	s.x[16] = uint64(uint32(s.x[11]) + 0xcda34bbd)
	s.compare32(uint32(s.x[14]), uint32(s.x[16]))
	gate := uint32(boolValue(s.condition("eq")))
	s.x[14] = uint64(gate)
	s.x[12] = uint64(uint32(s.x[15]) - uint32(s.x[12]))
	s.x[15] = uint64(uint32(s.x[12]) + 0xffff72c6)
	s.x[14] = uint64(gate * uint32(s.x[15]))
	s.x[13] = uint64(uint32(s.x[14]) ^ uint32(s.x[13]))
	next, err := loadBootstrapPointer(s, s.x[30], uint32(s.x[13]))
	if err != nil {
		return 0, fmt.Errorf("account difference-select target: %w", err)
	}
	s.x[13] = next
	s.x[17] = 0x7b0c
	return next, nil
}

func runAccountByteWriterA(s *state) (uint64, error) {
	s.traceBlock(accountByteWriterA)
	s.x[13] = s.x[2] + s.x[14]
	s.x[14]++
	value, err := s.load(s.x[13], 1)
	if err != nil {
		return 0, fmt.Errorf("account byte-writer A input: %w", err)
	}
	s.x[12] = uint64(byte(value))
	s.x[13] = uint64(uint32(s.x[10]) & uint32(s.x[12]<<1))
	s.x[12] = uint64(uint32(s.x[12]) - uint32(s.x[13]))
	s.x[12] = uint64(uint32(s.x[12]) + 0x1b)
	if err := s.storeValue(s.x[19]+uint64(uint32(s.x[8])), uint64(byte(s.x[12])), 1); err != nil {
		return 0, fmt.Errorf("account byte-writer A output: %w", err)
	}
	s.compare32(uint32(s.x[8]), 0x3f)
	gate := uint32(boolValue(s.condition("lo")))
	s.x[12] = uint64(gate)
	s.x[13] = uint64(gate << 3)
	s.x[12] = uint64(uint32(s.x[13]) - gate)
	s.x[12] = uint64(uint32(s.x[12]) ^ uint32(s.x[11]))
	next, err := s.load(s.x[21]+uint64(uint32(s.x[12]))*8, 8)
	if err != nil {
		return 0, fmt.Errorf("account byte-writer A target: %w", err)
	}
	s.x[12] = next
	s.x[8]++
	s.x[17] = 0x2473
	return next, nil
}

func runAccountByteWriterB(s *state) (uint64, error) {
	s.traceBlock(accountByteWriterB)
	s.x[12] = s.x[27] + s.x[13]
	s.x[13]++
	value, err := s.load(s.x[12], 1)
	if err != nil {
		return 0, fmt.Errorf("account byte-writer B input: %w", err)
	}
	s.x[11] = uint64(byte(value))
	s.x[12] = uint64(uint32(s.x[10]) & uint32(s.x[11]<<1))
	s.x[11] = uint64(uint32(s.x[11]) - uint32(s.x[12]))
	s.x[11] = uint64(uint32(s.x[11]) + 0x1b)
	if err := s.storeValue(s.x[24]+uint64(uint32(s.x[8])), uint64(byte(s.x[11])), 1); err != nil {
		return 0, fmt.Errorf("account byte-writer B output: %w", err)
	}
	s.compare32(uint32(s.x[8]), 0x3f)
	gate := uint32(boolValue(s.condition("lo")))
	s.x[11] = uint64(gate)
	s.x[12] = uint64(gate << 1)
	s.x[11] = uint64(uint32(s.x[12]) | gate<<5)
	s.x[11] = uint64(uint32(s.x[11]) | uint32(s.x[9]))
	next, err := loadBootstrapPointer(s, s.x[20], uint32(s.x[11]))
	if err != nil {
		return 0, fmt.Errorf("account byte-writer B target: %w", err)
	}
	s.x[11] = next
	s.x[8]++
	s.x[17] = 0x1d0a
	return next, nil
}

func runAccountByteWriterC(s *state) (uint64, error) {
	s.traceBlock(accountByteWriterCFront)
	s.x[13] = s.x[26] + s.x[15]
	s.x[15]++
	value, err := s.load(s.x[13], 1)
	if err != nil {
		return 0, fmt.Errorf("account byte-writer C input: %w", err)
	}
	s.x[12] = uint64(byte(value))
	s.x[13] = s.x[9] & s.x[10]
	s.x[13] += s.x[22]
	s.x[14] = uint64(uint32(s.x[11]) & uint32(s.x[12]<<1))
	s.x[12] = uint64(uint32(s.x[12]) - uint32(s.x[14]))
	s.x[12] = uint64(uint32(s.x[12]) + 0x1b)
	if err := s.storeValue(s.x[13], uint64(byte(s.x[12])), 1); err != nil {
		return 0, fmt.Errorf("account byte-writer C output: %w", err)
	}
	s.compare32(uint32(s.x[10]), 0x3e)
	s.traceBlock(accountByteWriterCTail)
	gate := uint32(boolValue(s.condition("hi")))
	s.x[12] = uint64(gate)
	s.x[13] = uint64(gate << 5)
	s.x[12] = uint64(uint32(s.x[13]) - gate<<1)
	s.x[12] = uint64(uint32(s.x[12]) ^ uint32(s.x[8]))
	next, err := loadBootstrapPointer(s, s.x[24], uint32(s.x[12]))
	if err != nil {
		return 0, fmt.Errorf("account byte-writer C target: %w", err)
	}
	s.x[12] = next
	s.x[10]++
	s.x[17] = 0x6564
	return next, nil
}
