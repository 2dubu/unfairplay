package recovery

import "testing"

func TestOuterTableRecoveredStateVectors(t *testing.T) {
	const (
		a = uint32(0x01234567)
		b = uint32(0x89abcdef)
		c = uint32(0x0f1e2d3c)
		d = uint32(0x4b5a6978)
	)
	tests := []struct {
		name string
		got  uint32
		want uint32
	}{
		{"pending-s2", outerTablePendingS2(a, b, c, d), 0x37bdc4ba},
		{"lane0-pending-s1", outerTableLane0PendingS1(a, b, c, d, 0xdeadbeef, 0x13579bdf, 0x2468ace0, 0xfedcba98), 0xe0371293},
		{"lane0-pending-s3", outerTableLane0PendingS3(a, b, c, d), 0xab79a9e6},
		{"lane2-s0", outerTableLane2NextS0(a, b, c, d), 0xeaa93a27},
		{"lane7-s0", outerTableLane7NextS0(a, b, c, d), 0x757f6a68},
		{"lane8-s1", outerTableLane8NextS1(a, b, c, d), 0xb01075a9},
		{"lane15-s4", outerTableLane15NextS4(a, b, c, d), 0xe8712c7b},
		{"lane19-s0", outerTableLane19NextS0(a), 0xffadebe5},
		{"lane5-s2", outerTableLane5NextS2(a, b, c, d, 0xdeadbeef, 0x13579bdf, 0x2468ace0), 0xf9d7feba},
		{"lane5-pending-s4", outerTableLane5PendingS4(a, b, c, d), 0x96dff58e},
		{"lane18-mix", outerTableLane18Mix(a, b, c, d), 0xbf461a26},
		{"lane19-late-mix", outerTableLane19LateMix(a, b, c, d), 0x3c4ec24d},
	}
	for _, test := range tests {
		if test.got != test.want {
			t.Errorf("%s = %#08x, want %#08x", test.name, test.got, test.want)
		}
	}
}
