package recovery

import "fmt"

const selectionSIDBProjectionPass = uint64(0x10006c78c)

func runSelectionSemanticBlock(s *state, pc uint64) (uint64, bool, error) {
	if pc != selectionSIDBProjectionPass {
		return pc, false, nil
	}
	next, err := runSelectionSIDBProjection(s)
	return next, true, err
}

// runSelectionSIDBProjection reduces the native SIDB container to its 944-byte
// encoded payload. The four lookup lanes are build-global byte encodings.
func runSelectionSIDBProjection(s *state) (uint64, error) {
	const maximumBytes = 2048
	for index := 0; index < maximumBytes; index++ {
		next, err := runSelectionSIDBProjectionStep(s)
		if err != nil || next != selectionSIDBProjectionPass {
			return next, err
		}
	}
	return 0, fmt.Errorf("selection SIDB projection exceeded %d bytes", maximumBytes)
}

func runSelectionSIDBProjectionStep(s *state) (uint64, error) {
	if s == nil {
		return 0, fmt.Errorf("selection SIDB projection state is unavailable")
	}
	s.traceBlock(selectionSIDBProjectionPass)
	s.x[12]--
	input, err := s.load(s.x[9]+s.x[12], 1)
	if err != nil {
		return 0, fmt.Errorf("selection SIDB projection input: %w", err)
	}
	s.x[17] = uint64(byte(input))
	lane := s.x[12] & 0xf
	s.x[1] = lane
	first, err := s.load(s.x[14]+lane, 1)
	if err != nil {
		return 0, fmt.Errorf("selection SIDB projection first lane: %w", err)
	}
	s.x[2] = uint64(byte(first))
	second, err := s.load(s.x[16]+lane+3, 1)
	if err != nil {
		return 0, fmt.Errorf("selection SIDB projection second lane: %w", err)
	}
	s.x[3] = uint64(byte(second))
	third, err := s.load(s.x[11]+lane+2, 1)
	if err != nil {
		return 0, fmt.Errorf("selection SIDB projection third lane: %w", err)
	}
	s.x[4] = uint64(byte(third))
	s.x[1] = uint64(uint32(lane) * uint32(s.x[15]))
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[2]))
	s.x[1] = uint64(uint32(s.x[3]) ^ uint32(s.x[1]))
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[1]))
	s.x[17] = uint64(uint32(s.x[17]) ^ uint32(s.x[4]))
	s.x[17] = uint64(uint32(s.x[17]) ^ 0x99999999)
	if err := s.storeValue(s.x[20]+s.x[12], uint64(byte(s.x[17])), 1); err != nil {
		return 0, fmt.Errorf("selection SIDB projection output: %w", err)
	}
	s.x[3] = uint64(uint32(s.x[10]) + uint32(s.x[13]))
	s.x[17] = s.x[3] ^ s.x[0]
	s.compare64(s.x[12], s.x[17])
	last := uint32(boolValue(s.condition("ne")))
	s.x[1] = uint64(last << 2)
	dispatchIndex := ((last << 2) | (last << 4)) ^ uint32(s.x[10])
	dispatchTable, err := s.load(addOffset(s.x[29], -0x90), 8)
	if err != nil {
		return 0, fmt.Errorf("selection SIDB projection dispatch table: %w", err)
	}
	next, err := loadBootstrapPointer(s, dispatchTable, dispatchIndex)
	if err != nil {
		return 0, fmt.Errorf("selection SIDB projection dispatch: %w", err)
	}
	s.x[1] = next
	s.x[17] = 0x4b2b
	return next, nil
}
