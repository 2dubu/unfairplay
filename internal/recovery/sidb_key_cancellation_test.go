package recovery

import (
	"crypto/sha256"
	"fmt"
	"testing"
)

func cancellationComposite(x [16]byte) [16]byte {
	return currentSIDBProducerTyped.apply(currentFrameAES.apply(x))
}
func cancellationAffine(fn func([16]byte) [16]byte) frameAffine128 {
	var a frameAffine128
	a.constant = fn([16]byte{})
	for bit := range 128 {
		var x [16]byte
		x[bit/8] = 1 << uint(bit%8)
		y := fn(x)
		for j := range y {
			a.columns[bit][j] = y[j] ^ a.constant[j]
		}
	}
	return a
}

func cancellationCompose(outer, inner frameAffine128) frameAffine128 {
	var result frameAffine128
	result.constant = outer.apply(inner.constant)
	for i, c := range inner.columns {
		result.columns[i] = outer.linear(c)
	}
	return result
}
func cancellationSandwich(t *testing.T, a frameAffine128) frameAffine128 {
	t.Helper()
	fn := func(x [16]byte) [16]byte {
		x = frameSBoxState(x)
		x = a.apply(x)
		for j := range x {
			x[j] = frameInverseSBox[x[j]]
		}
		return x
	}
	c := cancellationAffine(fn)
	seen := map[int]bool{}
	for target := range 16 {
		source := -1
		for bit, column := range a.columns {
			if column[target] != 0 {
				if source < 0 {
					source = bit / 8
				}
				if source != bit/8 {
					t.Fatalf("non-bytewise sandwich target %d", target)
				}
			}
		}
		if source < 0 || seen[source] {
			t.Fatal("not a byte permutation")
		}
		seen[source] = true
		for value := range 256 {
			var x [16]byte
			x[source] = byte(value)
			want, got := fn(x), c.apply(x)
			if got[target] != want[target] {
				t.Fatalf("non-affine sandwich target %d value %d", target, value)
			}
		}
	}
	return c
}
func cancellationForwardLayers() [10]frameAffine128 {
	f := currentFrameAES
	var layers [10]frameAffine128
	layers[0] = f.first
	layers[1] = cancellationAffine(func(x [16]byte) [16]byte {
		for j := range x {
			x[j] = f.firstPost[j].apply(x[j])
		}
		x = applyFrameAESMix(f.mixes[0], x)
		for j := range x {
			x[j] ^= f.roundKeys[0][j]
		}
		return x
	})
	for i := 1; i < 8; i++ {
		layers[i+1] = cancellationAffine(func(x [16]byte) [16]byte {
			for j := range x {
				x[j] ^= f.roundMasks[i-1][j]
			}
			x = applyFrameAESMix(f.mixes[i], x)
			for j := range x {
				x[j] ^= f.roundKeys[i][j]
			}
			return x
		})
	}
	layers[9] = cancellationAffine(func(x [16]byte) [16]byte {
		for j := range x {
			x[j] ^= f.roundMasks[7][j]
		}
		x = applyFrameAESMix(f.mixes[8], x)
		for j := range x {
			x[j] = f.lastPre[j].apply(x[j]) ^ f.lastConstant[j]
		}
		return x
	})
	return layers
}
func cancellationInverseLayers() ([9]frameAffine128, frameAffine128) {
	p := currentSIDBProducerTyped
	var layers [9]frameAffine128
	layers[0] = cancellationAffine(func(x [16]byte) [16]byte {
		var y [16]byte
		for source, target := range p.firstTargets {
			y[target] = p.firstOutput[target].apply(x[source])
		}
		y = applySIDBCanonicalInverseRound(y)
		for j := range y {
			y[j] ^= p.roundKeys[0][j]
		}
		return y
	})
	for i := 1; i < 8; i++ {
		layers[i] = cancellationAffine(func(x [16]byte) [16]byte {
			x = applySIDBCanonicalInverseRound(x)
			for j := range x {
				x[j] ^= p.roundKeys[i][j]
			}
			return x
		})
	}
	layers[8] = cancellationAffine(func(x [16]byte) [16]byte {
		x = applySIDBCanonicalInverseRound(x)
		for j := range x {
			x[j] ^= p.lastKeys[j]
		}
		return x
	})
	post := cancellationAffine(func(x [16]byte) [16]byte {
		var last, out [16]byte
		for source, target := range p.lastTargets {
			last[target] = x[source]
		}
		for target, mask := range p.finalXOR {
			for source := range 16 {
				if mask&(1<<uint(source)) != 0 {
					out[target] ^= last[source]
				}
			}
			out[target] = p.finalOutput[target].apply(out[target]) ^ p.finalConstant[target]
		}
		return out
	})
	return layers, post
}

// TestSIDBKeySuffixCancellation proves each nonlinear sandwich for its
// complete byte domain. Each bridge must map exactly one distinct source byte
// to each target byte; thus the exhaustive byte checks cover all 2^128 inputs.
func TestSIDBKeySuffixCancellation(t *testing.T) {
	forward := cancellationForwardLayers()
	inverse, post := cancellationInverseLayers()
	// The first producer InvSubBytes directly cancels the suffix SubBytes.
	bridge := cancellationCompose(inverse[0], forward[9])
	for step := range 9 {
		reduced := cancellationSandwich(t, bridge)
		if step == 8 {
			bridge = cancellationCompose(post, cancellationCompose(reduced, forward[0]))
			break
		}
		bridge = cancellationCompose(inverse[step+1], cancellationCompose(reduced, forward[8-step]))
	}
	expected := cancellationAffine(cancellationInverseMixColumns)
	expected.constant = selectionSIDBFinalKey
	if bridge != expected {
		t.Fatal("exact reduction is not InvMixColumns XOR the fixed final key")
	}
}

func cancellationInverseMixColumns(x [16]byte) (y [16]byte) {
	matrix := [4][4]byte{{14, 11, 13, 9}, {9, 14, 11, 13}, {13, 9, 14, 11}, {11, 13, 9, 14}}
	for column := range 4 {
		for row := range 4 {
			for source := range 4 {
				y[4*column+row] ^= frameGFMul(matrix[row][source], x[4*column+source])
			}
		}
	}
	return y
}

func TestSIDBKeyReducedMatchesComposition(t *testing.T) {
	for control := range 4096 {
		digest := sha256.Sum256([]byte(fmt.Sprint("reduced-key", control)))
		var input [16]byte
		copy(input[:], digest[:])
		want := cancellationComposite(frameCore(input))
		if got := deriveSelectionSIDBKey(input); got != want {
			t.Fatalf("control %d: %x != %x", control, got, want)
		}
	}
	// Both MixColumns maps are linear, so basis equality proves cancellation.
	for bit := range 128 {
		var input [16]byte
		input[bit/8] = 1 << uint(bit%8)
		if cancellationInverseMixColumns(cancellationMixColumns(input)) != input {
			t.Fatal("MixColumns inverse basis mismatch")
		}
	}
}

func cancellationMixColumns(x [16]byte) (y [16]byte) {
	matrix := [4][4]byte{{2, 3, 1, 1}, {1, 2, 3, 1}, {1, 1, 2, 3}, {3, 1, 1, 2}}
	for column := range 4 {
		for row := range 4 {
			for source := range 4 {
				y[4*column+row] ^= frameGFMul(matrix[row][source], x[4*column+source])
			}
		}
	}
	return y
}

func TestSIDBKeyIsNotStandardAESKeySchedule(t *testing.T) {
	// Expanding the whitened INITIAL key is incorrect; the separate interior-
	// schedule test proves the standard AES representation with the true key.
	expanded := expandSelectionSIDBAES128Key(frameRoundKeys[0])
	var second [16]byte
	copy(second[:], expanded[16:32])
	if second == frameRoundKeys[1] {
		t.Fatal("revisit independent-round-key classification")
	}
}

var cancellationBenchmarkKey [16]byte

func BenchmarkSIDBKeyOriginalComposition(b *testing.B) {
	for b.Loop() {
		cancellationBenchmarkKey = cancellationComposite(frameCore(cancellationBenchmarkKey))
	}
}
func BenchmarkSIDBKeyReduced(b *testing.B) {
	for b.Loop() {
		cancellationBenchmarkKey = deriveSelectionSIDBKey(cancellationBenchmarkKey)
	}
}

// Every interior round yields the same master key. Expansion reproduces all
// ten remaining keys, while the initial mismatch is exactly input whitening.
func TestSIDBKeyStandardAESWithInputWhitening(t *testing.T) {
	expanded := expandSelectionSIDBAES128Key(selectionSIDBAESKey)
	for round := 1; round <= 10; round++ {
		want := selectionSIDBFinalKey
		if round < 10 {
			want = frameRoundKeys[round]
		}
		var got [16]byte
		copy(got[:], expanded[round*16:(round+1)*16])
		if got != want {
			t.Fatalf("AES schedule mismatch at round %d", round)
		}
		master, err := reverseSelectionSIDBAES128RoundKey(want, round)
		if err != nil || master != selectionSIDBAESKey {
			t.Fatalf("master mismatch from round %d", round)
		}
	}
	for index := range selectionSIDBInputWhitening {
		if frameRoundKeys[0][index]^selectionSIDBAESKey[index] != selectionSIDBInputWhitening[index] {
			t.Fatalf("input whitening mismatch at byte %d", index)
		}
	}
}
