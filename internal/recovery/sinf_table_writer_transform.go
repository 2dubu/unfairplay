package recovery

import "fmt"

const (
	sinfTableWriter         = uint64(0x1000802c8)
	sinfExpandedTableWriter = uint64(0x1000f4584)
)

func runSINFTableWriter(s *state) (uint64, error) {
	s.traceBlock(sinfTableWriter)
	s.x[9] = uint64(uint32(s.x[11]) ^ 0x27c6)
	s.x[12] = 0x58a
	s.x[23] = uint64(uint32(s.x[9]) * uint32(s.x[12]))
	s.x[9] = uint64(uint32(s.x[23]) + 0xffff63ac)
	s.x[1] = uint64(uint32(s.x[9]) | 0x800)
	base, err := s.load(s.x[6], 8)
	if err != nil {
		return 0, fmt.Errorf("SINF table-writer source base: %w", err)
	}
	s.x[9] = base + s.x[5]
	value, err := s.load(s.x[9], 1)
	if err != nil {
		return 0, fmt.Errorf("SINF table-writer source byte: %w", err)
	}
	s.x[9] = uint64(byte(value))
	value, err = s.load(s.x[28]+0x210, 4)
	if err != nil {
		return 0, fmt.Errorf("SINF table-writer first word: %w", err)
	}
	s.x[12] = uint64(uint32(value))
	const stateBase = uint64(0x101075000)
	s.x[27] = stateBase
	value, err = s.load(stateBase+0xb40, 4)
	if err != nil {
		return 0, fmt.Errorf("SINF table-writer second word: %w", err)
	}
	s.x[13] = uint64(uint32(value))
	s.x[12] = uint64(uint32(s.x[12]) ^ uint32(s.x[13]))
	s.x[24] = 0xa6
	s.x[12] = uint64(uint32(s.x[12]) ^ uint32(s.x[24]))
	s.x[25] = 0x0a5cd4cf
	s.x[12] = uint64(uint32(s.x[12]) * uint32(s.x[25]))
	s.x[13] = s.x[16] + uint64(byte(s.x[12]))
	value, err = s.load(s.x[13], 1)
	if err != nil {
		return 0, fmt.Errorf("SINF table-writer first table byte: %w", err)
	}
	s.x[13] = uint64(uint32(byte(value)) ^ 0xa1)
	value, err = s.load(s.x[8]+uint64(uint32(s.x[13])), 1)
	if err != nil {
		return 0, fmt.Errorf("SINF table-writer second table byte: %w", err)
	}
	s.x[13] = uint64(byte(value))
	s.compare32(uint32(s.x[9]), 0)
	s.x[12] = uint64(uint32(s.x[12]) ^ uint32(s.x[13]))
	s.x[12] &= 0xff
	s.x[12] = uint64(uint32(s.x[12]) + 0x1a5)
	value, err = s.load(s.x[10]+uint64(uint32(s.x[12]))*8, 8)
	if err != nil {
		return 0, fmt.Errorf("SINF table-writer dispatch state: %w", err)
	}
	s.x[12] = value
	gate := uint32(boolValue(s.condition("eq")))
	s.x[13] = uint64(gate)
	s.x[15] = s.x[12] + 0xffffffffdb1cf25e
	table, err := s.load(s.x[15], 8)
	if err != nil {
		return 0, fmt.Errorf("SINF table-writer output table: %w", err)
	}
	s.x[12] = table
	if err := s.storeValue(table+(s.x[5]<<2), uint64(uint32(s.x[9])), 4); err != nil {
		return 0, fmt.Errorf("SINF table-writer lane 0: %w", err)
	}
	s.x[0] = uint64(uint32(s.x[9] << 8))
	table, err = s.load(s.x[15], 8)
	if err != nil {
		return 0, fmt.Errorf("SINF table-writer output table reload 1: %w", err)
	}
	s.x[12] = table + 0x400
	if err := s.storeValue(s.x[12]+(s.x[5]<<2), uint64(uint32(s.x[0])), 4); err != nil {
		return 0, fmt.Errorf("SINF table-writer lane 1: %w", err)
	}
	s.x[12] = uint64(uint32(s.x[9] << 16))
	table, err = s.load(s.x[15], 8)
	if err != nil {
		return 0, fmt.Errorf("SINF table-writer output table reload 2: %w", err)
	}
	s.x[17] = table + 0x800
	if err := s.storeValue(s.x[17]+(s.x[5]<<2), uint64(uint32(s.x[12])), 4); err != nil {
		return 0, fmt.Errorf("SINF table-writer lane 2: %w", err)
	}
	s.x[17] = 0xffffd818
	s.x[17] = uint64(uint32(s.x[1]) + uint32(s.x[17]))
	s.x[17] = uint64(uint32(s.x[9]) << (uint32(s.x[17]) & 31))
	table, err = s.load(s.x[15], 8)
	if err != nil {
		return 0, fmt.Errorf("SINF table-writer output table reload 3: %w", err)
	}
	s.x[15] = table + 0xc00
	if err := s.storeValue(s.x[15]+(s.x[5]<<2), uint64(uint32(s.x[17])), 4); err != nil {
		return 0, fmt.Errorf("SINF table-writer lane 3: %w", err)
	}
	s.x[15] = 0xfe4
	s.x[13] = uint64(gate * uint32(s.x[15]))
	s.x[11] = uint64(uint32(s.x[13]) ^ uint32(s.x[11]))
	next, err := loadBootstrapPointer(s, s.x[14], uint32(s.x[11]))
	if err != nil {
		return 0, fmt.Errorf("SINF table-writer dispatch: %w", err)
	}
	s.x[13] = next
	s.x[26] = 0xed43636c
	s.x[11] = 0x635fa6e0
	s.x[17] = 0x1a30
	return next, nil
}

func runSINFExpandedTableWriter(s *state) (uint64, error) {
	s.traceBlock(sinfExpandedTableWriter)
	s.x[9] = 0x35fb004a
	s.x[13] = uint64(uint32(s.x[11]) ^ 0x799e8ff9)
	s.x[12] = uint64(uint32(s.x[13]) & uint32(s.x[12]))
	s.x[11] = uint64(uint32(s.x[11]) & 0xe921f500)
	s.x[13] = uint64(uint32(s.x[12]) ^ uint32(s.x[11]))
	value, err := s.load(s.x[28]+0x210, 4)
	if err != nil {
		return 0, fmt.Errorf("SINF expanded-table first word: %w", err)
	}
	s.x[11] = uint64(uint32(value))
	value, err = s.load(s.x[30], 4)
	if err != nil {
		return 0, fmt.Errorf("SINF expanded-table second word: %w", err)
	}
	s.x[15] = uint64(uint32(value) ^ 0xa6)
	s.x[11] = uint64(uint32(s.x[11]) ^ uint32(s.x[15]))
	s.x[15] = 0x0a5cd4cf
	s.x[11] = uint64(uint32(s.x[11]) * uint32(s.x[15]))
	base, err := s.load(s.sp+0x18, 8)
	if err != nil {
		return 0, fmt.Errorf("SINF expanded-table byte base: %w", err)
	}
	s.x[15] = base + uint64(byte(s.x[11]))
	value, err = s.load(s.x[15], 1)
	if err != nil {
		return 0, fmt.Errorf("SINF expanded-table first byte: %w", err)
	}
	s.x[15] = uint64(uint32(byte(value)) ^ 0xcd)
	value, err = s.load(s.x[24]+uint64(uint32(s.x[15])), 1)
	if err != nil {
		return 0, fmt.Errorf("SINF expanded-table second byte: %w", err)
	}
	s.x[15] = uint64(byte(value))
	s.x[11] = uint64(uint32(s.x[11]) ^ uint32(s.x[15]))
	s.x[11] &= 0xff
	s.x[11] = uint64(uint32(s.x[11]) + 0xfda)
	value, err = s.load(s.x[10]+uint64(uint32(s.x[11]))*8, 8)
	if err != nil {
		return 0, fmt.Errorf("SINF expanded-table state: %w", err)
	}
	s.x[11] = value
	s.x[15] = s.x[11] + 0xffffffffe476b2fe
	table, err := s.load(s.x[15], 8)
	if err != nil {
		return 0, fmt.Errorf("SINF expanded-table output table: %w", err)
	}
	s.x[11] = table
	s.x[17] = uint64(uint32(s.x[13]) ^ 0x96ff7aff)
	if err := s.storeValue(table+(s.x[5]<<2), uint64(uint32(s.x[17])), 4); err != nil {
		return 0, fmt.Errorf("SINF expanded-table lane 0: %w", err)
	}
	s.x[11] = uint64(uint32(s.x[13]) ^ 0x00ff7aff)
	s.x[17] = uint64(uint32(s.x[13]) ^ 0x96ff0000)
	s.x[11] = uint64(uint32(extractPair(s.x[11], s.x[17], 24, 32)))
	table, err = s.load(s.x[15], 8)
	if err != nil {
		return 0, fmt.Errorf("SINF expanded-table output table reload 1: %w", err)
	}
	s.x[17] = table + 0x400
	if err := s.storeValue(s.x[17]+(s.x[5]<<2), uint64(uint32(s.x[11])), 4); err != nil {
		return 0, fmt.Errorf("SINF expanded-table lane 1: %w", err)
	}
	s.x[11] = uint64(uint32(s.x[21]) + 0x35fb004a)
	s.x[11] = uint64(uint32(s.x[11]) & 0xca04e7ec)
	s.x[9] = uint64(uint32(s.x[13]) ^ 0x7aff)
	s.x[17] = uint64(uint32(s.x[13]) ^ 0x96ff0000)
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[11]))
	s.x[9] = uint64(uint32(extractPair(s.x[9], s.x[17], 16, 32)))
	table, err = s.load(s.x[15], 8)
	if err != nil {
		return 0, fmt.Errorf("SINF expanded-table output table reload 2: %w", err)
	}
	s.x[17] = table + 0x800
	if err := s.storeValue(s.x[17]+(s.x[5]<<2), uint64(uint32(s.x[9])), 4); err != nil {
		return 0, fmt.Errorf("SINF expanded-table lane 2: %w", err)
	}
	s.x[9] = uint64(^uint32(s.x[12]))
	s.x[12] = uint64(uint32(s.x[13]) ^ 0x96ff7a00)
	s.x[9] = uint64(uint32(extractPair(s.x[9], s.x[12], 8, 32)))
	table, err = s.load(s.x[15], 8)
	if err != nil {
		return 0, fmt.Errorf("SINF expanded-table output table reload 3: %w", err)
	}
	s.x[12] = table + 0xc00
	if err := s.storeValue(s.x[12]+(s.x[5]<<2), uint64(uint32(s.x[9])), 4); err != nil {
		return 0, fmt.Errorf("SINF expanded-table lane 3: %w", err)
	}
	s.x[5]++
	s.compare64(s.x[5], 0x100)
	gate := uint32(boolValue(s.condition("ne")))
	s.x[9] = uint64(gate)
	s.x[12] = uint64(gate << 1)
	s.x[9] = uint64(uint32(s.x[12]) | gate<<4)
	s.x[9] = uint64(uint32(s.x[21]) | uint32(s.x[9]))
	next, err := loadBootstrapPointer(s, s.x[14], uint32(s.x[9]))
	if err != nil {
		return 0, fmt.Errorf("SINF expanded-table dispatch: %w", err)
	}
	s.x[9] = next
	s.x[17] = 0x9b0b
	return next, nil
}
