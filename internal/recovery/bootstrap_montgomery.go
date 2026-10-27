package recovery

import (
	"fmt"
	"math/big"
)

type montgomeryWordResult struct {
	Words []uint32
	Carry uint32
}

const bootstrapAccumulatorCycleExit = uint64(0x1007f70c4)
const bootstrapMontgomeryReductionExit = uint64(0x100992c48)

func runBootstrapMixer(s *state, _ *rawSourceRuntime, _ *executionProfile) (uint64, error) {
	return runBootstrapMontgomerySquare(s)
}

func runBootstrapAccumulatorCycle(s *state) (uint64, error) {
	// The one-round form is a final correction rather than a complete
	// multiplication; retain its already-verified semantic controller.
	if uint32(s.x[25]) != 0xb5 {
		return runBootstrapAccumulatorCycleCircuit(s)
	}
	return runBootstrapMontgomeryMultiply(s)
}

// runBootstrapMontgomeryReduction replaces the flattened final conversion out
// of the Montgomery domain. The input is an encoded (limbs+1)-word
// accumulator. The native bootstrap has already performed its first radix
// step, so this loop performs the remaining limbs-1 reduction steps.
func runBootstrapMontgomeryReduction(s *state) (uint64, error) {
	if s == nil {
		return 0, fmt.Errorf("bootstrap Montgomery reduction state is unavailable")
	}
	limbs := int(uint32(s.x[27]))
	mask := uint32(s.x[2])
	if limbs != 32 && limbs != 64 || mask != 0x224ea33b || s.x[28] == 0 ||
		s.x[7] != 0xc0c02ce911b20ca0 {
		return 0, fmt.Errorf("unsupported bootstrap Montgomery reduction ABI")
	}
	modulusPointer, err := s.load(s.sp+0x70, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap Montgomery reduction modulus: %w", err)
	}
	lastPointer, err := s.load(s.sp+0x58, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap Montgomery reduction last-word pointer: %w", err)
	}
	carryPointer, err := s.load(s.sp+0x60, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap Montgomery reduction carry pointer: %w", err)
	}
	if lastPointer != s.x[28]+uint64((limbs-1)*4) || carryPointer != s.x[28]+uint64(limbs*4) {
		return 0, fmt.Errorf("unsupported bootstrap Montgomery reduction geometry")
	}

	words, err := loadMontgomeryWords(s, s.x[28], limbs+1, mask, "reduction accumulator")
	if err != nil {
		return 0, err
	}
	modulus, err := loadMontgomeryWords(s, modulusPointer, limbs, mask, "reduction modulus")
	if err != nil {
		return 0, err
	}
	inverseValue, err := s.load(s.sp+0x7c, 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap Montgomery reduction inverse: %w", err)
	}
	inverse := uint32(inverseValue)
	lastCarry, err := montgomeryReduceWords(words, modulus, inverse, limbs-1)
	if err != nil {
		return 0, err
	}
	for index, word := range words {
		if err := s.storeValue(s.x[28]+uint64(index*4), uint64(word^mask), 4); err != nil {
			return 0, fmt.Errorf("bootstrap Montgomery reduction output %d: %w", index, err)
		}
	}

	startCounter := uint32(s.x[20])
	finalCounter := startCounter + uint32(limbs-1)
	if err := s.storeValue(s.sp+0x8c, uint64(finalCounter-1), 4); err != nil {
		return 0, fmt.Errorf("bootstrap Montgomery reduction counter: %w", err)
	}
	if err := s.storeValue(s.sp+0x98, 0x69f2, 8); err != nil {
		return 0, fmt.Errorf("bootstrap Montgomery reduction state: %w", err)
	}
	dispatchTable, err := s.load(s.sp+0xa8, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap Montgomery reduction dispatch table: %w", err)
	}

	// Preserve the compatibility ABI expected by the first block after the
	// flattened loop. These are encodings and dispatcher state, not secrets.
	s.x[1] = lastCarry ^ s.x[7]
	s.x[8] = bootstrapMontgomeryReductionExit
	s.x[10] = uint64(finalCounter)
	s.x[11] = dispatchTable
	s.x[14] = 0
	s.x[15] = 0x2cd4
	s.x[17] = 0xc98e
	s.x[20] = uint64(finalCounter)
	s.x[22] = 0x4c9e
	s.x[30] = uint64(mask)
	return bootstrapMontgomeryReductionExit, nil
}

// montgomeryReduceWords applies radix-2^32 Montgomery reduction to an
// (n+1)-word accumulator. It retains the extra carry word because the native
// bootstrap intentionally postpones its final conditional subtraction.
func montgomeryReduceWords(words, modulus []uint32, inverse uint32, rounds int) (uint64, error) {
	if len(modulus) == 0 || len(words) != len(modulus)+1 || rounds < 0 || rounds > len(modulus) {
		return 0, fmt.Errorf("invalid Montgomery reduction geometry %d/%d/%d", len(words), len(modulus), rounds)
	}
	if modulus[0]&1 == 0 || modulus[0]*inverse != ^uint32(0) {
		return 0, fmt.Errorf("invalid Montgomery reduction inverse")
	}
	var lastCarry uint64
	for round := 0; round < rounds; round++ {
		factor := words[0] * inverse
		carry := uint64(0)
		for index, modulusWord := range modulus {
			total := uint64(factor)*uint64(modulusWord) + uint64(words[index]) + carry
			if index == 0 && uint32(total) != 0 {
				return 0, fmt.Errorf("Montgomery reduction did not cancel its low limb")
			}
			if index > 0 {
				words[index-1] = uint32(total)
			}
			carry = total >> 32
			lastCarry = total
		}
		total := uint64(words[len(modulus)]) + carry
		words[len(modulus)-1] = uint32(total)
		words[len(modulus)] = uint32(total >> 32)
	}
	return lastCarry, nil
}

// runBootstrapMontgomeryMultiply replaces one complete accumulator cycle with
// a typed CIOS Montgomery multiplication. The cycle's operand descriptors are
// stack-resident; its output uses the same masked 32-bit limb format as the
// recovered square path.
func runBootstrapMontgomeryMultiply(s *state) (uint64, error) {
	limbs := int(uint32(s.x[10]))
	if limbs != 32 && limbs != 64 || uint32(s.x[16]) != 0x224ea33b {
		return 0, fmt.Errorf("unsupported bootstrap Montgomery multiply ABI")
	}
	leftPointer, err := s.load(s.sp+0x70, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap Montgomery left operand: %w", err)
	}
	rightObject, err := s.load(s.sp+0x40, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap Montgomery right object: %w", err)
	}
	rightPointer, err := s.load(rightObject+8, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap Montgomery right operand: %w", err)
	}
	modulusObject, err := s.load(s.sp+0xa0, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap Montgomery modulus object: %w", err)
	}
	modulusPointer, err := s.load(modulusObject+8, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap Montgomery modulus: %w", err)
	}
	lastPointer, err := s.load(s.sp+0x50, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap Montgomery last-word pointer: %w", err)
	}
	carryPointer, err := s.load(s.sp+0x58, 8)
	if err != nil {
		return 0, fmt.Errorf("bootstrap Montgomery carry pointer: %w", err)
	}
	if lastPointer != s.x[9]+uint64((limbs-1)*4) || carryPointer != s.x[9]+uint64(limbs*4) {
		return 0, fmt.Errorf("unsupported bootstrap Montgomery output geometry")
	}
	mask := uint32(s.x[16])
	left, err := loadMontgomeryWords(s, leftPointer, limbs, mask, "left operand")
	if err != nil {
		return 0, err
	}
	right, err := loadMontgomeryWords(s, rightPointer, limbs, mask, "right operand")
	if err != nil {
		return 0, err
	}
	modulus, err := loadMontgomeryWords(s, modulusPointer, limbs, mask, "modulus")
	if err != nil {
		return 0, err
	}
	inverse, err := s.load(s.sp+0x6c, 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap Montgomery multiply inverse: %w", err)
	}
	result, err := montgomeryCIOSMultiply(left, right, modulus, uint32(inverse))
	if err != nil {
		return 0, err
	}
	for index, word := range result.Words {
		if err := s.storeValue(s.x[9]+uint64(index*4), uint64(word^mask), 4); err != nil {
			return 0, fmt.Errorf("bootstrap Montgomery multiply output %d: %w", index, err)
		}
	}
	if err := s.storeValue(carryPointer, uint64(result.Carry^mask), 4); err != nil {
		return 0, fmt.Errorf("bootstrap Montgomery multiply carry: %w", err)
	}
	finalCounter := uint32(s.x[4]) + uint32(limbs)
	if err := s.storeValue(s.sp+0x98, uint64(finalCounter-1), 4); err != nil {
		return 0, fmt.Errorf("bootstrap Montgomery multiply counter: %w", err)
	}

	s.x[4] = uint64(finalCounter)
	s.x[5] = lastPointer
	s.x[6] = 0x68ed69f7
	s.x[7] = 0x5c37
	s.x[8] = bootstrapAccumulatorCycleExit
	s.x[13] = 0x930501b9
	s.x[15] = 0x1430e33e4efdacf8
	s.x[17] = 0x902d
	s.x[19] = 0x446c
	s.x[20] = uint64(finalCounter)
	s.x[21] = ^uint64(0x5c16)
	s.x[22] = 0x3426
	s.x[23] = 0x100ff1ef0
	s.x[25] = 0x38cce9d6c53c8609
	s.x[26] = 0xc73316293ac379f7
	s.x[27] = 0xa049
	s.x[28] = 0xebcf1cc1b1025308
	return bootstrapAccumulatorCycleExit, nil
}

// runBootstrapMontgomerySquare replaces one complete flattened bootstrap
// mixer invocation. The native path performs either a 32- or 64-limb CIOS
// Montgomery square and retains its unreduced carry limb.
func runBootstrapMontgomerySquare(s *state) (uint64, error) {
	if s == nil {
		return 0, fmt.Errorf("bootstrap Montgomery state is unavailable")
	}
	limbs := int(uint32(s.x[8]))
	if limbs != 32 && limbs != 64 {
		return 0, fmt.Errorf("unsupported bootstrap Montgomery width %d", limbs)
	}
	if uint32(s.x[0]) != 0xe663 || uint32(s.x[1]) != 0xc1e1812b+uint32(limbs) ||
		uint32(s.x[2]) != 0x779bc723 || uint32(s.x[3]) != 0x9ecc ||
		uint32(s.x[20]) != 0x779bc723 || uint32(s.x[24]) != 0x224ea33b {
		return 0, fmt.Errorf("unsupported bootstrap Montgomery ABI")
	}
	mask := uint32(s.x[24])
	input, err := loadMontgomeryWords(s, s.x[5], limbs, mask, "square input")
	if err != nil {
		return 0, err
	}
	modulus, err := loadMontgomeryWords(s, s.x[22], limbs, mask, "modulus")
	if err != nil {
		return 0, err
	}
	inverse, err := s.load(s.sp+0x5c, 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap Montgomery inverse: %w", err)
	}
	result, err := montgomeryCIOSSquare(input, modulus, uint32(inverse))
	if err != nil {
		return 0, err
	}
	for index, word := range result.Words {
		if err := s.storeValue(s.x[4]+uint64(index*4), uint64(word^mask), 4); err != nil {
			return 0, fmt.Errorf("bootstrap Montgomery output word %d: %w", index, err)
		}
	}
	if err := s.storeValue(s.x[4]+uint64(limbs*4), uint64(result.Carry^mask), 4); err != nil {
		return 0, fmt.Errorf("bootstrap Montgomery carry: %w", err)
	}
	counter, err := s.load(s.sp+0x80, 4)
	if err != nil {
		return 0, fmt.Errorf("bootstrap Montgomery counter: %w", err)
	}
	finalCounter := uint32(counter) + uint32(limbs)
	if err := s.storeValue(s.sp+0x80, uint64(finalCounter), 4); err != nil {
		return 0, fmt.Errorf("bootstrap Montgomery final counter: %w", err)
	}

	// Compatibility ABI for the selector following the native mixer.
	s.x[0] = 0x53ceff38
	s.x[1] = 0x224ed61b
	s.x[2] = 0x779bc722 + uint64(limbs)
	s.x[3] = bootstrapMixerAfterRecurrence
	s.x[5] = 0xbdb7
	s.x[6] = 0xb3563f76
	s.x[9] = bootstrapStackSelector
	s.x[10] = 0xe4e21e1a
	s.x[11] = 0x64e21e1a + uint64(limbs)
	s.x[12] = 0x7520
	s.x[13] = 0x886438de
	s.x[14] = uint64(finalCounter)
	s.x[15] = 0x886438dd
	s.x[16] = 0x84c4d2ed
	s.x[17] = 0x8766
	s.x[19] = 0x82b6
	s.x[20] = 0x4c0d
	s.x[21] = 0x886438dd
	s.x[22] = 0x04c4d2ed + uint64(limbs)
	s.x[23] = 0x8d290bcc
	s.x[27] = 0xa9292112
	s.x[28] = uint64(finalCounter)
	return bootstrapStackSelector, nil
}

func loadMontgomeryWords(s *state, address uint64, count int, mask uint32, name string) ([]uint32, error) {
	words := make([]uint32, count)
	for index := range words {
		word, err := s.load(address+uint64(index*4), 4)
		if err != nil {
			return nil, fmt.Errorf("bootstrap Montgomery %s word %d: %w", name, index, err)
		}
		words[index] = uint32(word) ^ mask
	}
	return words, nil
}

// montgomeryCIOSSquare performs coarsely integrated operand scanning (CIOS)
// Montgomery squaring in radix 2^32. Words are little-endian limbs. The native
// bootstrap retains the extra carry limb and deliberately does not perform the
// usual final conditional subtraction.
func montgomeryCIOSSquare(input, modulus []uint32, inverse uint32) (montgomeryWordResult, error) {
	return montgomeryCIOSMultiply(input, input, modulus, inverse)
}

// montgomeryCIOSMultiply is the general multiply used by the bootstrap around
// the recovered squaring fast path.
func montgomeryCIOSMultiply(left, right, modulus []uint32, inverse uint32) (montgomeryWordResult, error) {
	if len(left) == 0 || len(left) != len(right) || len(left) != len(modulus) {
		return montgomeryWordResult{}, fmt.Errorf("invalid Montgomery geometry %d/%d/%d", len(left), len(right), len(modulus))
	}
	if modulus[0]&1 == 0 || modulus[0]*inverse != ^uint32(0) {
		return montgomeryWordResult{}, fmt.Errorf("invalid Montgomery modulus inverse")
	}

	multiplicand := littleEndianWordsBig(right)
	modulusValue := littleEndianWordsBig(modulus)
	accumulator := new(big.Int)
	product := new(big.Int)
	reduction := new(big.Int)
	for _, word := range left {
		product.Mul(multiplicand, new(big.Int).SetUint64(uint64(word)))
		accumulator.Add(accumulator, product)
		factor := uint32(accumulator.Uint64()) * inverse
		reduction.Mul(modulusValue, new(big.Int).SetUint64(uint64(factor)))
		accumulator.Add(accumulator, reduction)
		if uint32(accumulator.Uint64()) != 0 {
			return montgomeryWordResult{}, fmt.Errorf("Montgomery reduction did not cancel its low limb")
		}
		accumulator.Rsh(accumulator, 32)
	}

	words := make([]uint32, len(left))
	remaining := new(big.Int).Set(accumulator)
	mask := new(big.Int).SetUint64(0xffffffff)
	for index := range words {
		words[index] = uint32(new(big.Int).And(remaining, mask).Uint64())
		remaining.Rsh(remaining, 32)
	}
	if !remaining.IsUint64() || remaining.Uint64() > 0xffffffff {
		return montgomeryWordResult{}, fmt.Errorf("Montgomery carry exceeds one limb")
	}
	return montgomeryWordResult{Words: words, Carry: uint32(remaining.Uint64())}, nil
}

func littleEndianWordsBig(words []uint32) *big.Int {
	result := new(big.Int)
	for index := len(words) - 1; index >= 0; index-- {
		result.Lsh(result, 32)
		result.Or(result, new(big.Int).SetUint64(uint64(words[index])))
	}
	return result
}
