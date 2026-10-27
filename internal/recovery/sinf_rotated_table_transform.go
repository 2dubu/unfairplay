package recovery

import "fmt"

const sinfRotatedTableWriter = uint64(0x1001aa93c)

func runSINFRotatedTableWriter(s *state) (uint64, error) {
	s.traceBlock(sinfRotatedTableWriter)
	s.x[9] = 0xffff7158
	s.x[13] = uint64(uint32(s.x[23]) + 0x025943cc)
	s.x[22] = uint64(uint32(s.x[13]) | 0x74801781)
	s.x[9] = uint64(uint32(s.x[23]) + uint32(s.x[9]))
	s.x[21] = uint64(uint32(s.x[9]) | 0x1201)
	s.x[9] = uint64(uint32(s.x[26]) ^ 0x12bc9c93)
	s.x[9] = uint64(uint32(s.x[9]) & uint32(s.x[11]))
	s.x[11] = uint64(uint32(s.x[22]) ^ 0xea7a4ebe)
	s.x[11] = uint64(uint32(s.x[11]) & uint32(s.x[26]))
	s.x[9] = uint64(uint32(s.x[9]) ^ 0xd2b270f3)
	s.x[9] = uint64(uint32(s.x[9]) ^ uint32(s.x[11]))
	s.x[9] = uint64(uint32(s.x[9]) &^ uint32(s.x[0]))
	s.x[11] = uint64(uint32(s.x[0]) & 0x6800)
	s.x[9] = uint64(uint32(s.x[9]) | uint32(s.x[11]))
	s.x[11] = uint64(uint32(s.x[21]) + 0x46544022)
	s.x[9] = uint64(uint32(s.x[9]) ^ uint32(s.x[11]))
	s.x[9] = uint64(uint32(s.x[9]) &^ uint32(s.x[12]))
	s.x[28] = s.x[4]
	value, err := s.load(s.x[4]+0x210, 4)
	if err != nil {
		return 0, fmt.Errorf("SINF rotated-table first word: %w", err)
	}
	s.x[11] = uint64(uint32(value))
	value, err = s.load(s.x[30], 4)
	if err != nil {
		return 0, fmt.Errorf("SINF rotated-table second word: %w", err)
	}
	s.x[13] = uint64(uint32(value))
	s.x[11] = uint64(uint32(s.x[13]) + uint32(s.x[11]))
	s.x[11] = uint64(uint32(s.x[11]) ^ uint32(s.x[24]))
	s.x[11] = uint64(uint32(s.x[11]) * uint32(s.x[25]))
	firstBase, err := s.load(s.sp+8, 8)
	if err != nil {
		return 0, fmt.Errorf("SINF rotated-table first byte base: %w", err)
	}
	secondBase, err := s.load(s.sp+0x10, 8)
	if err != nil {
		return 0, fmt.Errorf("SINF rotated-table second byte base: %w", err)
	}
	s.x[17], s.x[13] = firstBase, secondBase
	s.x[13] += uint64(byte(s.x[11]))
	value, err = s.load(s.x[13], 1)
	if err != nil {
		return 0, fmt.Errorf("SINF rotated-table first byte: %w", err)
	}
	s.x[13] = uint64(uint32(byte(value)) ^ 0x84)
	value, err = s.load(s.x[17]+uint64(uint32(s.x[13])), 1)
	if err != nil {
		return 0, fmt.Errorf("SINF rotated-table second byte: %w", err)
	}
	s.x[13] = uint64(byte(value))
	s.x[12] = uint64(uint32(s.x[12]) | 0x7bb9e8d0)
	s.x[11] = uint64(uint32(s.x[11]) ^ uint32(s.x[13]))
	s.x[11] &= 0xff
	s.x[11] = uint64(uint32(s.x[11]) + 0xd61)
	value, err = s.load(s.x[10]+uint64(uint32(s.x[11]))*8, 8)
	if err != nil {
		return 0, fmt.Errorf("SINF rotated-table first state: %w", err)
	}
	s.x[11] = value
	s.x[9] = uint64(uint32(s.x[9]) ^ uint32(s.x[12]))
	s.x[12] = 0xffffffff9913d2dd
	s.x[11] += s.x[12]
	table, err := s.load(s.x[11], 8)
	if err != nil {
		return 0, fmt.Errorf("SINF rotated-table first output table: %w", err)
	}
	if err := storeRotatedWordLanes(s, table, s.x[5], uint32(s.x[9])); err != nil {
		return 0, fmt.Errorf("SINF rotated-table first output: %w", err)
	}

	value, err = s.load(s.x[20], 8)
	if err != nil {
		return 0, fmt.Errorf("SINF rotated-table source: %w", err)
	}
	s.x[9] = value
	value, err = s.load(s.x[4]+0x210, 4)
	if err != nil {
		return 0, fmt.Errorf("SINF rotated-table third word: %w", err)
	}
	s.x[11] = uint64(uint32(value))
	value, err = s.load(0x101075000+0xb40, 4)
	if err != nil {
		return 0, fmt.Errorf("SINF rotated-table fourth word: %w", err)
	}
	s.x[12] = uint64(uint32(value))
	s.x[11] = uint64(uint32(s.x[12]) + uint32(s.x[11]))
	s.x[11] = uint64(uint32(s.x[11]) ^ uint32(s.x[24]))
	s.x[11] = uint64(uint32(s.x[11]) * uint32(s.x[25]))
	s.x[12] = s.x[7] + uint64(byte(s.x[11]))
	value, err = s.load(s.x[12], 1)
	if err != nil {
		return 0, fmt.Errorf("SINF rotated-table third byte: %w", err)
	}
	s.x[12] = uint64(uint32(byte(value)) ^ 0x5d)
	value, err = s.load(s.x[19]+uint64(uint32(s.x[12])), 1)
	if err != nil {
		return 0, fmt.Errorf("SINF rotated-table fourth byte: %w", err)
	}
	s.x[12] = uint64(byte(value))
	s.x[9] += s.x[5]
	s.x[11] = uint64(uint32(s.x[11]) ^ uint32(s.x[12]))
	s.x[12] = 0x1010
	s.x[11] = uint64(uint32(s.x[12]) + uint32(byte(s.x[11])))
	value, err = s.load(s.x[10]+uint64(uint32(s.x[11]))*8, 8)
	if err != nil {
		return 0, fmt.Errorf("SINF rotated-table second state: %w", err)
	}
	s.x[11] = value
	value, err = s.load(s.x[9], 1)
	if err != nil {
		return 0, fmt.Errorf("SINF rotated-table source byte: %w", err)
	}
	s.x[9] = uint64(byte(value))
	s.x[12] = 0xffffffff9c6ae4a6
	s.x[11] += s.x[12]
	table, err = s.load(s.x[11], 8)
	if err != nil {
		return 0, fmt.Errorf("SINF rotated-table second output table: %w", err)
	}
	s.compare32(uint32(s.x[9]), 0)
	if err := storeRotatedWordLanes(s, table, s.x[5], uint32(s.x[9])); err != nil {
		return 0, fmt.Errorf("SINF rotated-table second output: %w", err)
	}
	s.x[15] = table + 0x800
	gate := uint32(boolValue(s.condition("ne")))
	s.x[12] = uint64(gate)
	s.x[11] = 0xce
	s.x[11] = uint64(gate * uint32(s.x[11]))
	s.x[11] = uint64(uint32(s.x[11]) ^ uint32(s.x[23]))
	next, err := loadBootstrapPointer(s, s.x[14], uint32(s.x[11]))
	if err != nil {
		return 0, fmt.Errorf("SINF rotated-table dispatch: %w", err)
	}
	s.x[13] = next
	s.x[11] = 0x86617006
	s.x[12] = 0x16de0aff
	value, err = s.load(s.sp+0x20, 8)
	if err != nil {
		return 0, fmt.Errorf("SINF rotated-table saved register: %w", err)
	}
	s.x[24] = value
	s.x[17] = 0xc429
	return next, nil
}

func storeRotatedWordLanes(s *state, table, index uint64, word uint32) error {
	for lane, rotation := range [...]int{0, 24, 16, 8} {
		value := word
		if rotation != 0 {
			value = uint32(rotateRight(uint64(word), rotation, 32))
		}
		if err := s.storeValue(table+uint64(lane)*0x400+(index<<2), uint64(value), 4); err != nil {
			return err
		}
	}
	return nil
}
