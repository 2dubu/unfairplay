package recovery

import "fmt"

const (
	sharedMBAFront = uint64(0x1005515e0)
	sharedMBATail  = uint64(0x100551800)
)

// runSharedMBASemantic evaluates the two-table enrollment transform without
// executing the native register/MBA circuit. Between its table lookups the
// transform is a bivariate uint32 polynomial followed by signed reduction
// modulo 599 and a byte permutation.
func runSharedMBASemantic(s *state) (uint64, error) {
	if s == nil {
		return 0, fmt.Errorf("shared MBA state is unavailable")
	}
	s.traceBlock(sharedMBAFront)
	if uint32(s.x[19]) != 0x19d6 || uint32(s.x[14])+uint32(s.x[15]) != 0x24a2fbac {
		return 0, fmt.Errorf("unsupported shared MBA ABI")
	}
	s.x[1] = 0x65c0
	firstTable, err := s.load(s.sp+0x58, 8)
	if err != nil {
		return 0, fmt.Errorf("shared MBA first table: %w", err)
	}
	s.x[28] = firstTable
	firstValue, err := s.load(firstTable+uint64(sharedMBAFirstIndex(uint32(s.x[15]))), 1)
	if err != nil {
		return 0, fmt.Errorf("shared MBA first table byte: %w", err)
	}
	folded := sharedMBAFoldByte(byte(firstValue))
	remainder := sharedMBARemainder(folded, uint32(s.x[22]))
	s.traceBlock(sharedMBATail)
	secondTable, err := s.load(s.sp+0x50, 8)
	if err != nil {
		return 0, fmt.Errorf("shared MBA second table: %w", err)
	}
	secondValue, err := s.load(secondTable+uint64(sharedMBATailIndex(byte(remainder))), 1)
	if err != nil {
		return 0, fmt.Errorf("shared MBA second table byte: %w", err)
	}
	next, err := loadBootstrapPointer(s, s.x[21], uint32(s.x[19]))
	if err != nil {
		return 0, fmt.Errorf("shared MBA target: %w", err)
	}
	s.x[2] = next
	s.x[3] = 0x060be5d0
	s.x[17] = 0xcd20
	s.x[20] = uint64(byte(secondValue))
	s.x[27] = bootstrapTargetTable
	return next, nil
}

func sharedMBARemainder(folded, stateWord uint32) int32 {
	step := folded + 0x1b
	stateStep := stateWord + 0x1b
	statePolynomial := stateWord*0x75 + 0x48575af0
	stateProduct := stateWord * 0xb7a8b167
	statePolynomial = statePolynomial*stateStep + stateProduct
	stateProduct = statePolynomial * 0x3bd4ef30
	statePolynomial += 0x5ecab5dd
	stateProduct = stateWord*0x9e9ca0fb + stateProduct
	statePolynomial = statePolynomial*(stateWord+0xc42b10eb) + stateProduct

	left := folded*0xb + 0xa40d49ed
	left = left*step + folded*0x5bf2b73c
	left += 0x4ed6e7d8
	left = left*step + folded*0x63c26b7c
	right := (left + 0x85815614) * (folded + 0xc961596a)
	secondary := (folded+0x9d449283)*step + folded*0x62bb6d98
	left = left*0x369ea6b1 + right
	statePolynomial = secondary*0x9cba2e5b + statePolynomial
	statePolynomial = (secondary+0x69c48f08)*(folded+0x6345d1c0) + statePolynomial
	right = (folded*0xe7+0xd7876c42)*step + stateWord*0xd8f6830e + folded*0x2878ac1b

	value := (left + 0x1ca7fc33) * step
	value += folded * 0x67ed7fa1
	value += 0x576fd500
	value = value*stateStep + statePolynomial
	value += 0x27c02bb3
	value = value*stateStep + right
	raw := (value+0x26b9f953)*(stateWord+0xcda8f9e7) + value*0x32570634 + folded + 0x2ad396f7
	return int32(raw) % 599
}

// sharedMBAFirstIndex is the first table selector after applying the loop
// invariant x14+x15=0x24a2fbac. CoBRA/Z3 reduces the original circuit to the
// low byte of counter+0x54.
func sharedMBAFirstIndex(counter uint32) byte {
	return byte(counter + 0x54)
}

// sharedMBAFoldByte is the exact 32-bit value hidden by the three-term MBA
// immediately after the first table lookup.
func sharedMBAFoldByte(value byte) uint32 {
	word := uint32(value)
	return ((word + word) & 0xffffffc8) - word
}

// sharedMBATailIndex is the exact byte permutation hidden by the final MBA
// network. The preceding signed reduction bounds its input to [-598, 598]; an
// exhaustive proof over that domain shows the network depends only on the low
// byte and reduces to this add-and-toggle permutation.
func sharedMBATailIndex(value byte) byte {
	return value + (0xc9 ^ ((value & 1) << 7))
}
