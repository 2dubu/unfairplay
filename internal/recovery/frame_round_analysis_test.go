package recovery

import "testing"

func affineByteMap(values [256]byte) bool {
	offset := values[0]
	for input := range 256 {
		want := byte(0)
		for bit := range 8 {
			if input&(1<<bit) != 0 {
				want ^= values[1<<bit] ^ offset
			}
		}
		if values[input]^offset != want {
			return false
		}
	}
	return true
}

func TestFrameGeneralizedRoundComponentRelations(t *testing.T) {
	for roundIndex, round := range frameSemantic26A428.rounds {
		for position := range 16 {
			for sourceSlot := range 4 {
				var inverse [256]byte
				var seen [256]bool
				for input := range 256 {
					value := round.tables[position][input][sourceSlot]
					if seen[value] {
						t.Fatalf("round %d input %d slot %d is not a permutation", roundIndex, position, sourceSlot)
					}
					seen[value] = true
					inverse[value] = byte(input)
				}
				for targetSlot := range 4 {
					var relation [256]byte
					for sourceValue := range 256 {
						input := inverse[sourceValue]
						relation[sourceValue] = round.tables[position][input][targetSlot]
					}
					if !affineByteMap(relation) {
						t.Fatalf("round %d input %d slots %d->%d are not affine-related", roundIndex, position, sourceSlot, targetSlot)
					}
				}
			}
		}
	}
}

func TestFrameGeneralizedRoundSubstitutionsHaveAESDDT(t *testing.T) {
	for roundIndex, round := range frameSemantic26A428.rounds {
		for position := range 16 {
			for inputDelta := 1; inputDelta < 256; inputDelta++ {
				var counts [256]uint16
				for input := range 256 {
					left := round.tables[position][input][0]
					right := round.tables[position][input^inputDelta][0]
					counts[left^right]++
				}
				zero, two, four := 0, 0, 0
				for _, count := range counts {
					switch count {
					case 0:
						zero++
					case 2:
						two++
					case 4:
						four++
					default:
						t.Fatalf("round %d input %d delta %#02x has DDT count %d", roundIndex, position, inputDelta, count)
					}
				}
				if zero != 129 || two != 126 || four != 1 {
					t.Fatalf("round %d input %d delta %#02x DDT shape = zero:%d two:%d four:%d", roundIndex, position, inputDelta, zero, two, four)
				}
			}
		}
	}
}

func TestFrameGeneralizedRoundAESAffineEquivalence(t *testing.T) {
	for roundIndex, round := range frameSemantic26A428.rounds {
		for position := range 16 {
			var target [256]byte
			for input := range 256 {
				target[input] = round.tables[position][input][0]
			}
			inputMap, outputMap, inputOffset, outputOffset, _, err := solveAffineSBox(frameSBox, target)
			if err != nil {
				t.Fatalf("round %d input %d: %v", roundIndex, position, err)
			}
			for input := range 256 {
				want := outputMap.apply(frameSBox[inputMap.apply(byte(input))^inputOffset]) ^ outputOffset
				if target[input] != want {
					t.Fatalf("round %d input %d value %#02x = %#02x, want %#02x", roundIndex, position, input, target[input], want)
				}
			}
		}
	}
}

func byteTablePermutation(table [256]byte) bool {
	var seen [256]bool
	for _, value := range table {
		if seen[value] {
			return false
		}
		seen[value] = true
	}
	return true
}

func TestFrameFinalRoundStructure(t *testing.T) {
	direct, shared := 0, 0
	for target, entry := range frameSemantic26A428.final {
		switch entry.inputCount {
		case 1:
			direct++
			if !byteTablePermutation(entry.direct) {
				t.Fatalf("final output %d direct table is not a permutation", target)
			}
		case 3:
			shared++
			if !byteTablePermutation(entry.encode) {
				t.Fatalf("final output %d outer encoding is not a permutation", target)
			}
			for slot, contribution := range entry.contributions {
				if !byteTablePermutation(contribution) {
					t.Fatalf("final output %d contribution %d is not a permutation", target, slot)
				}
				inputMap, outputMap, inputOffset, outputOffset, _, err := solveAffineSBox(frameSBox, contribution)
				if err != nil {
					t.Fatalf("final output %d contribution %d: %v", target, slot, err)
				}
				for input := range 256 {
					want := outputMap.apply(frameSBox[inputMap.apply(byte(input))^inputOffset]) ^ outputOffset
					if contribution[input] != want {
						t.Fatalf("final output %d contribution %d value %#02x = %#02x, want %#02x", target, slot, input, contribution[input], want)
					}
				}
			}
		default:
			t.Fatalf("final output %d has invalid input count %d", target, entry.inputCount)
		}
	}
	if direct != 4 || shared != 12 {
		t.Fatalf("final geometry = %d direct and %d shared outputs", direct, shared)
	}
}
