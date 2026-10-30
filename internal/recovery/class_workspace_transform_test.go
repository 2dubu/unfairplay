package recovery

import (
	"encoding/binary"
	"math/rand"
	"testing"
)

func protectedClassWorkspaceResidueStep(sample, peer uint16, left, right uint32) uint16 {
	sampleWord := uint32(sample)
	if sample < 0x7b0f {
		sampleWord += 0x10000
	}
	peerWord := uint32(uint16(peer - 0x54ce))
	value := left + sampleWord + peerWord + right
	quotient := uint32((uint64(value) * 0x20005) >> 32)
	quotient = (quotient + ((value - quotient) >> 1)) >> 14
	value += quotient - quotient<<15
	value -= 0x1ca8 & (value << 1)
	value -= 0x71ac
	value ^= 0x691c
	return uint16(value)
}

func protectedClassWorkspaceReduce(value uint32) uint32 {
	quotient := uint32((uint64(value) * 0x20005) >> 32)
	quotient = (quotient + ((value - quotient) >> 1)) >> 14
	return value + quotient - quotient<<15
}

func protectedClassWorkspaceLimbStep(sample, peer, current, previous uint16, left, right uint32) classWorkspaceLimbOutputs {
	forward := protectedClassWorkspaceResidueStep(sample, peer, left, right)

	value := uint32(forward)
	a := (value ^ 0xe566e1eb) + 0x1a991e15
	b := (value ^ 0xacc1e249) + 0x533e1db7
	value = ((value ^ 0x49a7e4ea) + 0xb6581b16) ^ (a ^ b)
	value += uint32(uint16(current - 0x54ce))
	value += 0xe748
	value = protectedClassWorkspaceReduce(value)

	encodedCurrent := (value ^ uint32(0xffffd4cf)) + (uint32(0xffffa99e) & (value << 1))
	currentOutput := uint16(encodedCurrent + 0x7fff)
	currentResidue := uint16(encodedCurrent + 0x2b31)

	previousWord := uint32(previous)
	if previous < 0x7b0f {
		previousWord |= 0x10000
	}
	value = previousWord + uint32(currentResidue) + right
	value = protectedClassWorkspaceReduce(value)

	mask := uint32(0xdeb2) & (value << 1)
	value -= mask
	value += 0x6f59
	a = (value ^ 0xba8c) + 0x4015
	b = (value ^ 0x76dc) + uint32(0xffff8c45)
	a ^= b
	value = (value ^ 0xa309) + 0x5992
	value = (a ^ value) + uint32(0xffffe54d)

	return classWorkspaceLimbOutputs{
		Forward:  forward,
		Current:  currentOutput,
		Previous: uint16(value),
	}
}

func TestClassWorkspaceResidueStepMatchesProtectedArithmetic(t *testing.T) {
	controls := [][4]uint32{
		{0xce07, 0x5ca9, 0x0cb2eae0, 0xffff84f1},
		{0, 0, 0, 0},
		{0xffff, 0xffff, 0xffffffff, 0xffffffff},
		{0x7b0e, 0x54ce, 0x7fffffff, 0x80000000},
		{0x7b0f, 0x54cd, 0x80000000, 0x7fffffff},
	}
	random := rand.New(rand.NewSource(0x7fff))
	for range 1 << 16 {
		controls = append(controls, [4]uint32{
			uint32(random.Uint32()), uint32(random.Uint32()), random.Uint32(), random.Uint32(),
		})
	}
	for index, control := range controls {
		got := classWorkspaceResidueStep(uint16(control[0]), uint16(control[1]), control[2], control[3])
		want := protectedClassWorkspaceResidueStep(uint16(control[0]), uint16(control[1]), control[2], control[3])
		if got != want {
			t.Fatalf("control %d = %#x, want %#x", index, got, want)
		}
	}
	if got := classWorkspaceResidueStep(0xce07, 0x5ca9, 0x0cb2eae0, 0xffff84f1); got != 0xb851 {
		t.Fatalf("natural control = %#x, want 0xb851", got)
	}
}

func TestClassWorkspaceLimbStepMatchesProtectedArithmetic(t *testing.T) {
	random := rand.New(rand.NewSource(0x11d250))
	for index := range 1 << 16 {
		sample := uint16(random.Uint32())
		peer := uint16(random.Uint32())
		current := uint16(random.Uint32())
		previous := uint16(random.Uint32())
		left := random.Uint32()
		right := random.Uint32()
		got := classWorkspaceLimbStep(sample, peer, current, previous, left, right)
		want := protectedClassWorkspaceLimbStep(sample, peer, current, previous, left, right)
		if got != want {
			t.Fatalf("control %d = %#v, want %#v", index, got, want)
		}
	}
}

func protectedClassWorkspaceProductWord(limbs [9]uint16) uint32 {
	var stored [9]uint16
	for index, limb := range limbs {
		stored[index] = uint16(uint32(limb) + uint32(0x7b0f))
	}
	value := uint32(stored[7]) + uint32(stored[5])
	middle := uint32(stored[4]) + uint32(stored[6]) + uint32(stored[2])
	low := uint32(uint16(stored[3]-0x7b0f)) + uint32(uint16(stored[1]-0x7b0f))
	low <<= 8
	low += uint32(uint16(stored[0] - 0x7b0f))
	low += uint32(uint16(stored[8] - 0x7b0f))
	return low + middle<<16 + value<<24 + 0xcb0cf893
}

func TestClassWorkspaceProductWordMatchesProtectedArithmetic(t *testing.T) {
	random := rand.New(rand.NewSource(0x359db4))
	for index := range 1 << 16 {
		var limbs [9]uint16
		for limb := range limbs {
			limbs[limb] = uint16(random.Uint32() % classWorkspaceModulus)
		}
		got := classWorkspaceProductWord(limbs)
		want := protectedClassWorkspaceProductWord(limbs)
		if got != want {
			t.Fatalf("control %d = %#x, want %#x", index, got, want)
		}
	}
}

var protectedClassWorkspaceValueColumns = [32]uint64{
	0x5b7ae1ca789424fc, 0xb6f4d98ff03348e3, 0x77f3a905fb6690dd, 0xeefd490aedcc3ba1,
	0xc7e19214c1837659, 0x95d93f28991decb2, 0x31a97e50293ac37f, 0x6249fca052749dfe,
	0x7ae1ca789424fc5b, 0xf4d98ff03348e3b6, 0xf3a905fb6690dd77, 0xfd490aedcc3ba1ee,
	0xe19214c1837659c7, 0xd93f28991decb295, 0xa97e50293ac37f31, 0x49fca052749dfe62,
	0xe1ca789424fc5b7a, 0xd98ff03348e3b6f4, 0xa905fb6690dd77f3, 0x490aedcc3ba1eefd,
	0x9214c1837659c7e1, 0x3f28991decb295d9, 0x7e50293ac37f31a9, 0xfca052749dfe6249,
	0xca789424fc5b7ae1, 0x8ff03348e3b6f4d9, 0x05fb6690dd77f3a9, 0x0aedcc3ba1eefd49,
	0x14c1837659c7e192, 0x28991decb295d93f, 0x50293ac37f31a97e, 0xa052749dfe6249fc,
}

func protectedClassWorkspaceValue(word uint32) uint64 {
	result := uint64(0x35ea3e4cfba587c2)
	for bit, column := range protectedClassWorkspaceValueColumns {
		if word&(1<<bit) != 0 {
			result ^= column
		}
	}
	return result
}

func TestClassWorkspaceClassValueMatchesRecoveredAffineMap(t *testing.T) {
	random := rand.New(rand.NewSource(0x3d92ac))
	for index := range 1 << 16 {
		var limbs [9]uint16
		for limb := range limbs {
			limbs[limb] = uint16(random.Uint32() % classWorkspaceModulus)
		}
		value := classWorkspaceClassValue(limbs)
		got := binary.LittleEndian.Uint64(value[:])
		want := protectedClassWorkspaceValue(classWorkspacePolynomialWord(limbs))
		if got != want {
			t.Fatalf("control %d = %#x, want %#x", index, got, want)
		}
	}
}

func TestClassWorkspaceNativeWordBiasCancels(t *testing.T) {
	const nativeClassInputBias = uint32(0xa5c6076d)
	sum := uint64(classWorkspaceProductBias) + uint64(nativeClassInputBias)
	if uint32(sum) != 0 {
		t.Fatalf("native biases do not cancel")
	}
}

func referenceClassWorkspacePowerFactor(components [3]uint16) uint16 {
	var norm uint64
	for _, component := range components {
		value := uint64(component)
		norm += value * value
	}
	norm %= uint64(classWorkspaceModulus)
	result := uint64(1)
	for exponent := 45; exponent != 0; exponent >>= 1 {
		if exponent&1 != 0 {
			result = result * norm % uint64(classWorkspaceModulus)
		}
		norm = norm * norm % uint64(classWorkspaceModulus)
	}
	return uint16(result)
}

func TestClassWorkspacePowerFactorMatchesModularExponentiation(t *testing.T) {
	controls := [][3]uint16{{}, {1, 2, 3}, {0x7ffe, 0x7ffe, 0x7ffe}}
	random := rand.New(rand.NewSource(0x26574))
	for range 1 << 16 {
		controls = append(controls, [3]uint16{
			uint16(random.Uint32() % classWorkspaceModulus),
			uint16(random.Uint32() % classWorkspaceModulus),
			uint16(random.Uint32() % classWorkspaceModulus),
		})
	}
	for index, control := range controls {
		got := classWorkspacePowerFactor(control)
		want := referenceClassWorkspacePowerFactor(control)
		if got != want {
			t.Fatalf("control %d = %#x, want %#x", index, got, want)
		}
	}
}

func TestClassWorkspaceInteractionStepSyntheticControl(t *testing.T) {
	matrix := [9]uint16{1, 2, 3, 4, 5, 6, 7, 8, 9}
	want := [9]uint16{0x0328, 0x6a0c, 0x7164, 0x6fd2, 0x5aad, 0x5812, 0x5ce7, 0x659f, 0x11df}
	if got := classWorkspaceInteractionStep(matrix); got != want {
		t.Fatalf("interaction step = %#v, want %#v", got, want)
	}
}

func modularPowerReference(base, exponent, modulus uint64) uint64 {
	result := uint64(1)
	for exponent != 0 {
		if exponent&1 != 0 {
			result = result * base % modulus
		}
		base = base * base % modulus
		exponent >>= 1
	}
	return result
}

func TestClassWorkspaceDriftMatchesPowerFormula(t *testing.T) {
	seen := make(map[uint16]struct{})
	for value := range uint16(classWorkspaceModulus) {
		got := classWorkspaceDrift(value)
		want := uint16((0x39a5*modularPowerReference(5, uint64(value), uint64(classWorkspaceModulus)) + 0x0c7e) % uint64(classWorkspaceModulus))
		if got != want {
			t.Fatalf("drift(%#x) = %#x, want %#x", value, got, want)
		}
		seen[got] = struct{}{}
	}
	if len(seen) != 150 {
		t.Fatalf("drift image contains %d values, want 150", len(seen))
	}
}

func TestClassWorkspaceSeedMatrixSyntheticControl(t *testing.T) {
	want := [9]uint16{0x0000, 0x0001, 0x0400, 0x459b, 0x0020, 0x5ed6, 0x6e2c, 0x3751, 0x0001}
	if got := classWorkspaceSeedMatrix(0); got != want {
		t.Fatalf("seed matrix = %#v, want %#v", got, want)
	}
	wantValue := [8]byte{0xf5, 0x24, 0x95, 0x06, 0x52, 0x76, 0x20, 0xe4}
	if got := classWorkspaceOutputFromSeed(0); got != wantValue {
		t.Fatalf("enrollment value = %x, want %x", got, wantValue)
	}
}

func protectedClassWorkspaceSeedFromDigestSuffix(suffix uint16) uint16 {
	word := uint32(suffix)
	mixed := (word & 0x281f) ^ (word & 0x0c)
	mixed = ^(mixed << 1)
	mixed += word - 0x57ec
	mixed = ((mixed ^ 0x8c45) + 0x73bb) ^
		((mixed ^ 0x4143) - 0x4143) ^
		((mixed ^ 0x6515) - 0x6515)

	encoded := uint32(uint16(mixed - 0x469f))
	constant := (encoded & 0x6226) ^ 0x4ada
	zero := uint32(0)
	carry := (((zero - 0x27aa) ^ (0x5d64 & (encoded << 1))) & encoded) ^ constant
	propagate := encoded ^ 0x6eb0
	for range 12 {
		carry = (propagate & (carry << 1)) ^ constant
	}
	encoded = (encoded ^ (carry << 1) ^ 0xc2a7) & 0xffff
	peer := uint16(uint32(uint16(encoded)^0xa081) % classWorkspaceModulus)
	return uint16((uint32(peer) + 1) % classWorkspaceModulus)
}

func TestClassWorkspaceSeedFromDigestMatchesProtectedDomain(t *testing.T) {
	for value := range uint32(1 << 16) {
		var digest [20]byte
		digest[18] = byte(value >> 8)
		digest[19] = byte(value)
		got := classWorkspaceSeedFromDigest(digest)
		want := protectedClassWorkspaceSeedFromDigestSuffix(uint16(value))
		if got != want {
			t.Fatalf("digest suffix %#04x seed = %#04x, want %#04x", value, got, want)
		}
	}
}

func TestClassWorkspaceOutputFromDigestComposesTypedStages(t *testing.T) {
	digest := [20]byte{0x53, 0x49, 0x4e, 0x46, 0x2d, 0x53, 0x59, 0x4e, 0x54, 0x48, 0x45, 0x54, 0x49, 0x43, 0xa5, 0x5a, 0x13, 0x37, 0x24, 0x68}
	want := classWorkspaceOutputFromSeed(classWorkspaceSeedFromDigest(digest))
	if got := classWorkspaceOutputFromDigest(digest); got != want {
		t.Fatalf("digest enrollment value = %x, want %x", got, want)
	}
}

func TestLegacyClassWorkspaceWordsMatchIndependentControls(t *testing.T) {
	controls := []struct {
		residue uint16
		word    uint32
	}{
		{0x0000, 0x2405ecc0},
		{0x0001, 0x42f36f91},
		{0x1234, 0xe00199de},
		{0x7ffe, 0x480341d7},
	}
	for _, control := range controls {
		if got := legacyClassWorkspaceWord(control.residue); got != control.word {
			t.Fatalf("legacy workspace word(%#04x) = %#08x, want %#08x", control.residue, got, control.word)
		}
	}
}

func TestLegacyClassWorkspaceValueSyntheticControl(t *testing.T) {
	var digest [20]byte
	binary.BigEndian.PutUint16(digest[16:18], 0x1234)
	binary.BigEndian.PutUint16(digest[18:20], 0x7ffe)
	want := [8]byte{0x3d, 0x34, 0x45, 0xb2, 0xf0, 0x3e, 0x53, 0x7f}
	if got := legacyClassWorkspaceValue(digest); got != want {
		t.Fatalf("legacy class value = %x, want %x", got, want)
	}
}

// Encoded native boundaries used by the upstream research controls.
const (
	classWorkspaceSampleCutoff   = uint16(0x7b0f)
	classWorkspacePeerBias       = uint16(0x54ce)
	classWorkspaceOutputEncoding = uint16(0xe748)
	classWorkspaceProductBias    = uint32(0x5a39f893)
)

type classWorkspaceLimbOutputs struct {
	Forward  uint16
	Current  uint16
	Previous uint16
}

func classWorkspaceProductWord(limbs [9]uint16) uint32 {
	return classWorkspacePolynomialWord(limbs) + classWorkspaceProductBias
}

func classWorkspaceResidue(sample, peer uint16, left, right uint32) uint16 {
	sampleWord := uint32(sample)
	if sample < classWorkspaceSampleCutoff {
		sampleWord += 1 << 16
	}
	peerWord := uint32(uint16(peer - classWorkspacePeerBias))
	sum := left + sampleWord + peerWord + right
	return uint16(sum % classWorkspaceModulus)
}

func classWorkspaceResidueStep(sample, peer uint16, left, right uint32) uint16 {
	return classWorkspaceResidue(sample, peer, left, right) ^ classWorkspaceOutputEncoding
}

func classWorkspaceLimbStep(sample, peer, current, previous uint16, left, right uint32) classWorkspaceLimbOutputs {
	forwardResidue := classWorkspaceResidue(sample, peer, left, right)

	currentWord := uint32(uint16(current - classWorkspacePeerBias))
	currentResidue := uint16((uint32(forwardResidue) + currentWord) % classWorkspaceModulus)

	previousWord := uint32(previous)
	if previous < classWorkspaceSampleCutoff {
		previousWord += 1 << 16
	}
	previousResidue := uint16((previousWord + uint32(currentResidue) + right) % classWorkspaceModulus)

	return classWorkspaceLimbOutputs{
		Forward:  forwardResidue ^ classWorkspaceOutputEncoding,
		Current:  uint16(uint32(currentResidue) + uint32(classWorkspacePeerBias)),
		Previous: uint16(uint32(previousResidue) + uint32(classWorkspaceSampleCutoff)),
	}
}
