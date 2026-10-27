package recovery

import (
	"math/rand"
	"testing"
)

func TestOuterTemporaryTransformsMatchCircuits(t *testing.T) {
	random := rand.New(rand.NewSource(0x26a428))
	for sample := 0; sample < 10000; sample++ {
		var x [8]uint32
		for index := range x {
			x[index] = random.Uint32()
		}
		if got, want := outerTablePendingS2(x[0], x[1], x[2], x[3]), outerTablePendingS2Circuit(x[0], x[1], x[2], x[3]); got != want {
			t.Fatalf("pending s2 sample %d = %08x, want %08x", sample, got, want)
		}
		if got, want := outerTableLane0PendingS1(x[0], x[1], x[2], x[3], x[4], x[5], x[6], x[7]), outerTableLane0PendingS1Circuit(x[0], x[1], x[2], x[3], x[4], x[5], x[6], x[7]); got != want {
			t.Fatalf("pending s1 sample %d = %08x, want %08x", sample, got, want)
		}
		if got, want := outerTableLane0PendingS3(x[0], x[1], x[2], x[3]), outerTableLane0PendingS3Circuit(x[0], x[1], x[2], x[3]); got != want {
			t.Fatalf("pending s3 sample %d = %08x, want %08x", sample, got, want)
		}
		if got, want := outerTableLane5PendingS4(x[0], x[1], x[2], x[3]), outerTableLane5PendingS4Circuit(x[0], x[1], x[2], x[3]); got != want {
			t.Fatalf("pending s4 sample %d = %08x, want %08x", sample, got, want)
		}
		if got, want := outerTableLane18Mix(x[0], x[1], x[2], x[3]), outerTableLane18MixCircuit(x[0], x[1], x[2], x[3]); got != want {
			t.Fatalf("lane 18 sample %d = %08x, want %08x", sample, got, want)
		}
		if got, want := outerTableLane19LateMix(x[0], x[1], x[2], x[3]), outerTableLane19LateMixCircuit(x[0], x[1], x[2], x[3]); got != want {
			t.Fatalf("late mix sample %d = %08x, want %08x", sample, got, want)
		}
	}
}
