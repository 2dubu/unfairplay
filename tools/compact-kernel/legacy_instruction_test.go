package main

import (
	"strings"
	"testing"
)

func TestTranslateLegacyInstructionForms(t *testing.T) {
	for _, test := range []struct {
		source, opcode string
	}{
		{"sub w14, w3, w14, lsr #1", "compactSubLSR32"},
		{"orr x10, x10, #217020518514230019", "compactLogicalMask64"},
		{"eor x1, x2, #217020518514230019", "compactLogicalMask64"},
		{"and x3, x4, #217020518514230019", "compactLogicalMask64"},
		{"ldr w5, [x6, w5, sxtw]", "compactLoadSXTW32"},
		{"ldr x5, [x6, w5, sxtw]", "compactLoadSXTW64"},
		{"str w5, [x6, w5, sxtw]", "compactStoreSXTW32"},
		{"tst x17, x12", "compactTestRegister64"},
		{"mneg w8, w8, w16", "compactMNeg32"},
		{"mneg x8, x8, x16", "compactMNeg64"},
	} {
		got, stop, err := compactInstruction(test.source)
		if err != nil || stop || !strings.Contains(got, test.opcode) {
			t.Fatalf("%q => %q stop=%v err=%v", test.source, got, stop, err)
		}
	}
	if _, _, err := compactInstruction("sub w0, w1, w2, lsr #32"); err == nil {
		t.Fatal("accepted an invalid 32-bit shift")
	}
	if _, _, err := compactLogicalMaskFields("orr", "0", "1", 0x123456789abcdef0); err == nil {
		t.Fatal("accepted an invalid ARM64 logical immediate")
	}
}

func TestSplitMidBlockAuthenticatedCall(t *testing.T) {
	source := `
func rawSourceBlock1000() uint64 {
	// 0x1000: mov x8, #0x2000
	// 0x1004: blraa x8, x17
	// 0x1008: mov x9, #0x3000
	// 0x100c: braa x9, x17
}
`
	blocks, err := parseAnnotatedBlocks(source, "rawSourceBlock")
	if err != nil {
		t.Fatal(err)
	}
	if len(blocks) != 2 || blocks[0].start != 0x1000 || blocks[1].start != 0x1008 {
		t.Fatalf("split blocks: %+v", blocks)
	}
	generated, skipped, _, err := generateCompactCFGSource("testBlocks", blocks, 0, true)
	if err != nil || skipped != 0 {
		t.Fatalf("generate split blocks: skipped=%d err=%v", skipped, err)
	}
	text := string(generated)
	for _, expected := range []string{"0x1000:", "compactBlockCallRegister", "0x1008:", "compactBlockRegister"} {
		if !strings.Contains(text, expected) {
			t.Fatalf("generated split profile lacks %q:\n%s", expected, text)
		}
	}
}
