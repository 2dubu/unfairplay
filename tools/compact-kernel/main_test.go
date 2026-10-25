package main

import "testing"

func TestCompactInstruction(t *testing.T) {
	for _, test := range []struct {
		input, want string
		stop        bool
	}{
		{"mov w16, #0x58f20000", "{op: compactMovImm32High, destination: 16, immediate: 0x58f2},", false},
		{"add w15, w15, #0x131, lsl #12", "{op: compactAddImmShift32, destination: 15, left: 15, shift: 12, immediate: 0x131},", false},
		{"eor w8, w8, w9, lsl #1", "{op: compactXorLSL32, destination: 8, left: 8, right: 9, shift: 1},", false},
		{"ldr x8, [x9, w8, sxtw #3]", "{op: compactLoadSXTW64, destination: 8, left: 9, right: 8, shift: 3},", false},
		{"braa x8, x17", "", true},
	} {
		got, stop, err := compactInstruction(test.input)
		if err != nil {
			t.Fatalf("%q: %v", test.input, err)
		}
		if got != test.want || stop != test.stop {
			t.Fatalf("%q = %q, stop=%t; want %q, stop=%t", test.input, got, stop, test.want, test.stop)
		}
	}
}

func TestCompactInstructionRejectsUnknownOperation(t *testing.T) {
	if _, _, err := compactInstruction("svc #0"); err == nil {
		t.Fatal("unsupported instruction was accepted")
	}
}
