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
