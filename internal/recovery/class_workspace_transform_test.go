package recovery

import (
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
