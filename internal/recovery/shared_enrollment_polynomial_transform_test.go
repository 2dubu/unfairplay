package recovery

import (
	"math/rand"
	"testing"
)

func sharedEnrollmentPolynomialRawCircuit(left, right uint32) uint32 {
	x17, x19 := left, right
	x20 := x19*0xf707 + 0xf81ee69c
	x22 := x19 + 0x4bd39758
	x20 *= x22
	x20 = x19*0x41b024cc + x20
	x23 := x19 + 0xc363b8e9
	x24 := x20 + 0x9264fa20
	x23 = x24 * x23
	x20 = x20*0x886fde6f + x23
	x23 = x19 * 0xaf2d6c25
	x20 += 0x0bd8ba63
	x20 *= x22
	x20 = x19*0x72b7797d + x20
	x20 += 0x76381dd6
	x20 *= x22
	x19 = x19*0xf2896022 + x20
	x20 = x17 + 0xddd837a4
	x23 = x19*0xb8904dec + x23
	x19 += 0x98f319b0
	x19 = x20*x19 + x23
	x20 = x17*0xa14b + 0x6ea0c121
	x23 = x17 + 0x96688590
	x20 *= x23
	x20 = x17*0x5dfdf00f + x20
	x20 += 0xcc636f1f
	x20 *= x23
	x20 = x17*0x50436451 + x20
	x20 += 0x6f1418ff
	x20 *= x23
	x20 = x17*0xf0da6991 + x20
	x20 += 0x1874d931
	x20 *= x23
	x17 = x17*0xb426dd5f + x20
	x17 += 0x7049744b
	x17 = x17*x22 + x19
	return x17 + 0x0dbc1c05
}

func sharedEnrollmentPolynomialOutputCircuit(raw uint32) uint32 {
	product := uint64(int64(int32(raw)) * int64(-2068495821))
	high := uint32(product >> 32)
	high += raw
	quotient := uint32(int32(high) >> 15)
	quotient += high >> 31
	remainder := raw - quotient*0xf6eb
	remainder -= 0xa0 & (remainder << 1)
	return (remainder + 0x50) ^ 0x4b
}

func TestSharedEnrollmentPolynomialMatchesCircuit(t *testing.T) {
	random := rand.New(rand.NewSource(0x264c00))
	for trial := 0; trial < 100000; trial++ {
		left, right := random.Uint32(), random.Uint32()
		gotRaw := sharedEnrollmentPolynomialRaw(left, right)
		wantRaw := sharedEnrollmentPolynomialRawCircuit(left, right)
		if gotRaw != wantRaw {
			t.Fatalf("trial %d raw: got %#x want %#x", trial, gotRaw, wantRaw)
		}
		got := sharedEnrollmentPolynomialOutput(left, right)
		want := sharedEnrollmentPolynomialOutputCircuit(wantRaw)
		if got != want {
			t.Fatalf("trial %d output: got %#x want %#x", trial, got, want)
		}
	}
}
