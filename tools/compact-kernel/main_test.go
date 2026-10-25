package main

import (
	"strings"
	"testing"
)

func TestCompactInstruction(t *testing.T) {
	for _, test := range []struct {
		input, want string
		stop        bool
	}{
		{"mov w16, #0x58f20000", "{op: compactMovImm32High, destination: 16, immediate: 0x58f2},", false},
		{"add w15, w15, #0x131, lsl #12", "{op: compactAddImmShift32, destination: 15, left: 15, shift: 12, immediate: 0x131},", false},
		{"eor w8, w8, w9, lsl #1", "{op: compactXorLSL32, destination: 8, left: 8, right: 9, shift: 1},", false},
		{"ldr x8, [x9, w8, sxtw #3]", "{op: compactLoadSXTW64, destination: 8, left: 9, right: 8, shift: 3},", false},
		{"movk x23, #0x2484, lsl #32", "{op: compactMovK64, destination: 23, shift: 32, immediate: 0x2484},", false},
		{"adc w26, wzr, w27", "{op: compactADC32, destination: 26, left: 31, right: 27},", false},
		{"eor w16, w16, #0xffffff87", "{op: compactXorImm32, destination: 16, left: 16, right: 0xff, accumulator: 0xff, immediate: 0xff87},", false},
		{"mov w19, #-0x1080001", "{op: compactMovImm32Wide, destination: 19, right: 0xf7, accumulator: 0xfe, immediate: 0xffff},", false},
		{"mov x6, #0x25800000", "{op: compactMovImm64Shift, destination: 6, shift: 16, immediate: 0x2580},", false},
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

func TestGenerateCompactSourceKeepsDisjointSegments(t *testing.T) {
	source, err := generateCompactSource("testKernel", []annotatedInstruction{
		{pc: 0x1000, disassembly: "mov w1, #1"},
		{pc: 0x1004, disassembly: "braa x2, x3"},
		{pc: 0x2000, disassembly: "mov w4, #2"},
		{pc: 0x2004, disassembly: "braa x5, x6"},
	})
	if err != nil {
		t.Fatal(err)
	}
	text := string(source)
	for _, want := range []string{"start: 0x1000", "{index: 1, address: 0x2000}"} {
		if !strings.Contains(text, want) {
			t.Fatalf("generated source lacks %q:\n%s", want, text)
		}
	}
}
