package recovery

import "fmt"

const bootstrapByteMix = uint64(0x1008f3018)

func runBootstrapByteMix(s *state) (uint64, error) {
	s.traceBlock(bootstrapByteMix)
	word, err := s.load(s.x[7], 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap byte-mix first word: %w", err)
	}
	s.x[17] = uint64(uint32(word))
	word, err = s.load(s.x[22]+0xb40, 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap byte-mix second word: %w", err)
	}
	s.x[6] = uint64(uint32(word))
	s.x[28] = s.x[7]
	s.x[7] = uint64(uint32(s.x[6]) - uint32(s.x[17]))
	s.x[7] = uint64(uint32(s.x[7]) ^ uint32(s.x[12]))
	s.x[7] = uint64(uint32(s.x[7]) * uint32(s.x[13]))
	s.x[11] = s.x[19]
	s.x[19] = s.x[27] + uint64(byte(s.x[7]))
	value, err := s.load(s.x[19], 1)
	if err != nil {
		return 0, fmt.Errorf("bootstrap byte-mix first byte: %w", err)
	}
	s.x[19] = uint64(byte(value))
	s.x[19] = uint64(uint32(s.x[19]) ^ uint32(s.x[14]))
	value, err = s.load(s.x[26]+uint64(uint32(s.x[19])), 1)
	if err != nil {
		return 0, fmt.Errorf("bootstrap byte-mix second byte: %w", err)
	}
	s.x[19] = uint64(byte(value))
	s.x[7] = uint64(uint32(s.x[7]) ^ uint32(s.x[19]))
	s.x[7] &= 0xff
	s.x[7] = uint64(uint32(s.x[7]) + 0x437)
	value, err = s.load(s.x[23]+uint64(uint32(s.x[7]))*8, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap byte-mix first table: %w", err)
	}
	s.x[7] = value
	s.x[19] = s.x[8] & 0xf
	s.x[20] = s.x[19] + 1
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[16]))
	s.x[6] = uint64(uint32(s.x[17]) ^ uint32(s.x[6]))
	s.x[21] = s.x[19] + s.x[15]
	s.x[7] += s.x[21]
	s.x[6] = uint64(uint32(s.x[6]) * uint32(s.x[13]))
	s.x[21] = s.x[24] + uint64(byte(s.x[6]))
	value, err = s.load(s.x[21], 1)
	if err != nil {
		return 0, fmt.Errorf("bootstrap byte-mix third byte: %w", err)
	}
	s.x[21] = uint64(byte(value)) ^ 0xf0
	value, err = s.load(s.x[25]+uint64(uint32(s.x[21])), 1)
	if err != nil {
		return 0, fmt.Errorf("bootstrap byte-mix fourth byte: %w", err)
	}
	s.x[21] = uint64(byte(value))
	s.x[6] = uint64(uint32(s.x[6]) ^ uint32(s.x[21]))
	s.x[6] &= 0xff
	s.x[6] = uint64(uint32(s.x[6]) + 0x902)
	value, err = s.load(s.x[23]+uint64(uint32(s.x[6]))*8, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap byte-mix second table: %w", err)
	}
	s.x[6] = value
	s.x[21] = s.x[20] + s.x[0]
	s.x[6] += s.x[21]
	word, err = s.load(s.x[10], 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap byte-mix third word: %w", err)
	}
	s.x[21] = uint64(uint32(word))
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[21]))
	s.x[17] = uint64(uint32(s.x[17]) * uint32(s.x[13]))
	s.x[21] = s.x[27] + uint64(byte(s.x[17]))
	value, err = s.load(s.x[21], 1)
	if err != nil {
		return 0, fmt.Errorf("bootstrap byte-mix fifth byte: %w", err)
	}
	s.x[21] = uint64(byte(value)) ^ 4
	value, err = s.load(s.x[26]+uint64(uint32(s.x[21])), 1)
	if err != nil {
		return 0, fmt.Errorf("bootstrap byte-mix sixth byte: %w", err)
	}
	s.x[21] = uint64(byte(value))
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[21]))
	s.x[17] &= 0xff
	s.x[17] = uint64(uint32(s.x[17]) + 0x168)
	value, err = s.load(s.x[23]+uint64(uint32(s.x[17]))*8, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap byte-mix third table: %w", err)
	}
	s.x[17] = value
	s.x[20] += s.x[1]
	s.x[17] += s.x[20]
	s.x[20] = s.sp + 0xdbc
	s.x[20] += s.x[8]
	value, err = s.load(s.x[20], 1)
	if err != nil {
		return 0, fmt.Errorf("bootstrap byte-mix seventh byte: %w", err)
	}
	s.x[20] = uint64(byte(value))
	s.x[19] = uint64(uint32(s.x[19]) * uint32(s.x[2]))
	s.x[19] = uint64(uint32(s.x[20]) ^ uint32(s.x[19]))
	value, err = s.load(s.x[6], 1)
	if err != nil {
		return 0, fmt.Errorf("bootstrap byte-mix eighth byte: %w", err)
	}
	s.x[6] = uint64(byte(value))
	s.x[6] = uint64(uint32(s.x[19]) ^ uint32(s.x[6]))
	s.x[19] = s.x[11]
	value, err = s.load(s.x[7], 1)
	if err != nil {
		return 0, fmt.Errorf("bootstrap byte-mix ninth byte: %w", err)
	}
	s.x[7] = uint64(byte(value))
	s.x[6] = uint64(uint32(s.x[6]) ^ uint32(s.x[7]))
	s.x[7] = s.x[28]
	value, err = s.load(s.x[17], 1)
	if err != nil {
		return 0, fmt.Errorf("bootstrap byte-mix tenth byte: %w", err)
	}
	s.x[17] = uint64(byte(value))
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[3]))
	s.x[17] = uint64(uint32(s.x[6]) ^ uint32(s.x[17]))
	s.x[6] = s.x[9] + s.x[8]
	if err := s.storeValue(s.x[6], uint64(byte(s.x[17])), 1); err != nil {
		return 0, fmt.Errorf("bootstrap byte-mix output: %w", err)
	}
	s.compare64(s.x[8], s.x[4])
	gate := uint32(boolValue(s.condition("eq")))
	s.x[17] = uint64(gate * uint32(s.x[5]))
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[11]))
	next, err := loadBootstrapPointer(s, s.x[30], uint32(s.x[17]))
	if err != nil {
		return 0, fmt.Errorf("bootstrap byte-mix target: %w", err)
	}
	s.x[6] = next
	s.x[8]--
	s.x[17] = 0x2659
	return next, nil
}
