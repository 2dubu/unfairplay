package main

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/DrWhax/unfairplay/internal/fairplay"
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
		{"ror x16, x16, #8", "compactRORImm64"},
		{"ccmp x11, #0x0, #0x4, eq", "compactCCmpImm64"},
		{"eor x0, x17, x0, ror #61", "compactXorROR64"},
		{"uxtl v4.8h, v4.8b", "compactVectorUnsignedShiftLong"},
		{"uxtl2 v18.8h, v17.16b", "compactVectorUnsignedShiftLong"},
		{"sub x11, x11, x11, lsr #16", "compactSubLSR64"},
		{"ldrb w13, [x11, w10, sxtw]", "compactLoadSXTW8"},
		{"ldr q28, [sp, #32]", "compactLoadVector"},
		{"ldp q27, q26, [sp, #720]", "compactPairLoadVector"},
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

func TestExternalMachoBlockExtraction(t *testing.T) {
	path := os.Getenv("UNFAIRPLAY_PROFILE_BINARY")
	if path == "" {
		t.Skip("external FairPlay profile fixture is not configured")
	}
	layout, err := fairplay.InspectFairPlayRecoveryCodeLayout(path)
	if err != nil {
		t.Fatal(err)
	}
	blocks, err := annotatedMachoBlocks(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(blocks) != len(layout.BasicBlocks) {
		t.Fatalf("annotated blocks=%d layout blocks=%d", len(blocks), len(layout.BasicBlocks))
	}
	wanted := map[uint64]bool{}
	switch layout.BinarySHA256 {
	case fairplay.FairPlayBuild26A428Hash:
		wanted[0x1001511f0] = false
		wanted[0x100151304] = false
	case fairplay.FairPlayBuild26B118Hash:
		wanted[0x10026d170] = false
	default:
		return
	}
	for _, block := range blocks {
		if _, ok := wanted[block.start]; !ok {
			continue
		}
		generated, skipped, _, err := generateCompactCFGSource("selected", []annotatedBlock{block}, 0, true)
		if err != nil || skipped != 0 || !strings.Contains(string(generated), fmt.Sprintf("%#x:", block.start)) {
			t.Fatalf("block %#x was not translated: skipped=%d err=%v", block.start, skipped, err)
		}
		wanted[block.start] = true
	}
	for address, found := range wanted {
		if !found {
			t.Fatalf("missing extracted block %#x", address)
		}
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
