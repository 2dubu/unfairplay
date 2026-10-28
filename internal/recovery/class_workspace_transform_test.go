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
