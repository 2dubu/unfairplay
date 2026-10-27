package recovery

import (
	"fmt"
	"math/big"
)

type montgomeryWordResult struct {
	Words []uint32
	Carry uint32
}

func runBootstrapMixer(s *state, _ *rawSourceRuntime, _ *executionProfile) (uint64, error) {
	return runBootstrapMontgomerySquare(s)
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
	input := make([]uint32, limbs)
	modulus := make([]uint32, limbs)
	for index := 0; index < limbs; index++ {
		word, err := s.load(s.x[5]+uint64(index*4), 4)
		if err != nil {
			return 0, fmt.Errorf("bootstrap Montgomery input word %d: %w", index, err)
		}
		input[index] = uint32(word) ^ mask
		word, err = s.load(s.x[22]+uint64(index*4), 4)
		if err != nil {
			return 0, fmt.Errorf("bootstrap Montgomery modulus word %d: %w", index, err)
		}
		modulus[index] = uint32(word) ^ mask
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
