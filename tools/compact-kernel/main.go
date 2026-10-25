package main

import (
	"bufio"
	"flag"
	"fmt"
	"go/format"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

var (
	commentPattern       = regexp.MustCompile(`^\s*// (0x[0-9a-f]+):\s+(.+)$`)
	movImmediate         = regexp.MustCompile(`^mov ([wx])(\d+), #(0x[0-9a-f]+|[0-9]+)$`)
	movNegOne32          = regexp.MustCompile(`^mov w(\d+), #-1$`)
	movNegative32        = regexp.MustCompile(`^mov w(\d+), #-(0x[0-9a-f]+|[0-9]+)$`)
	movNegative64        = regexp.MustCompile(`^mov x(\d+), #-(0x[0-9a-f]+|[0-9]+)$`)
	movRegister          = regexp.MustCompile(`^mov ([wx])(\d+), [wx](\d+)$`)
	movKeep              = regexp.MustCompile(`^movk ([wx])(\d+), #(0x[0-9a-f]+|[0-9]+)(?:, lsl #(\d+))?$`)
	binary32             = regexp.MustCompile(`^(eor|and|add|sub|orr|mul) w(\d+), w(\d+), w(\d+)(?:, (lsl|lsr|ror) #(\d+))?$`)
	binary64             = regexp.MustCompile(`^(eor|add|sub|mul|and|orr) x(\d+), x(\d+), x(\d+)$`)
	binaryShift64        = regexp.MustCompile(`^(and|eor|add) x(\d+), x(\d+), x(\d+), lsl #(\d+)$`)
	addImmediate32       = regexp.MustCompile(`^add w(\d+), w(\d+), #(0x[0-9a-f]+|[0-9]+)(?:, lsl #(\d+))?$`)
	addImmediate64       = regexp.MustCompile(`^add x(\d+), x(\d+), #(0x[0-9a-f]+|[0-9]+)$`)
	subImmediate32       = regexp.MustCompile(`^sub w(\d+), w(\d+), #(0x[0-9a-f]+|[0-9]+)(?:, lsl #(\d+))?$`)
	subImmediate64       = regexp.MustCompile(`^sub x(\d+), x(\d+), #(0x[0-9a-f]+|[0-9]+)$`)
	xorImmediate32       = regexp.MustCompile(`^eor w(\d+), w(\d+), #(0x[0-9a-f]+|[0-9]+)$`)
	andImmediate32       = regexp.MustCompile(`^and w(\d+), w(\d+), #(0x[0-9a-f]+|[0-9]+)$`)
	logicalImmediate64   = regexp.MustCompile(`^(orr|eor) x(\d+), x(\d+), #(0x[0-9a-f]+|[0-9]+)$`)
	shiftLeft32          = regexp.MustCompile(`^lsl w(\d+), w(\d+), #(0x[0-9a-f]+|[0-9]+)$`)
	shiftLeft64          = regexp.MustCompile(`^lsl x(\d+), x(\d+), #(0x[0-9a-f]+|[0-9]+)$`)
	shiftLeftReg32       = regexp.MustCompile(`^lsl w(\d+), w(\d+), w(\d+)$`)
	addSP                = regexp.MustCompile(`^add x(\d+), sp, #(0x[0-9a-f]+|[0-9]+)(?:, lsl #(\d+))?$`)
	loadSP               = regexp.MustCompile(`^ldr ([wx])(\d+), \[sp, #(0x[0-9a-f]+|[0-9]+)\]$`)
	loadBase             = regexp.MustCompile(`^ldr w(\d+), \[x(\d+), #(0x[0-9a-f]+|[0-9]+)\]$`)
	loadBase32           = regexp.MustCompile(`^ldr w(\d+), \[x(\d+)\]$`)
	loadBaseImm64        = regexp.MustCompile(`^ldr x(\d+), \[x(\d+), #(0x[0-9a-f]+|[0-9]+)\]$`)
	loadBase64           = regexp.MustCompile(`^ldr x(\d+), \[x(\d+)\]$`)
	loadSXTW64           = regexp.MustCompile(`^ldr x(\d+), \[x(\d+), w(\d+), sxtw #(\d+)\]$`)
	loadUXTW32           = regexp.MustCompile(`^ldr w(\d+), \[x(\d+), w(\d+), uxtw #(\d+)\]$`)
	loadUXTW8            = regexp.MustCompile(`^ldrb w(\d+), \[x(\d+), w(\d+), uxtw\]$`)
	loadBaseImm8         = regexp.MustCompile(`^ldrb w(\d+), \[x(\d+)(?:, #(0x[0-9a-f]+|[0-9]+))?\]$`)
	loadSPImm8           = regexp.MustCompile(`^ldrb w(\d+), \[sp, #(0x[0-9a-f]+|[0-9]+)\]$`)
	loadXIndex8          = regexp.MustCompile(`^ldrb w(\d+), \[x(\d+), x(\d+)\]$`)
	loadPostIndex8       = regexp.MustCompile(`^ldrb w(\d+), \[x(\d+)\], #(0x[0-9a-f]+|[0-9]+)$`)
	loadXIndex           = regexp.MustCompile(`^ldr ([wx])(\d+), \[x(\d+), x(\d+)\]$`)
	loadSXTW32           = regexp.MustCompile(`^ldr w(\d+), \[x(\d+), w(\d+), sxtw #(\d+)\]$`)
	loadUXTW64           = regexp.MustCompile(`^ldr x(\d+), \[x(\d+), w(\d+), uxtw #(\d+)\]$`)
	storeSP              = regexp.MustCompile(`^str w(\d+), \[sp, #(0x[0-9a-f]+|[0-9]+)\]$`)
	storeSP64            = regexp.MustCompile(`^str x(\d+), \[sp, #(0x[0-9a-f]+|[0-9]+)\]$`)
	storeBase            = regexp.MustCompile(`^str w(\d+), \[x(\d+)\]$`)
	storeBase64          = regexp.MustCompile(`^str x(\d+), \[x(\d+)\]$`)
	storeBaseImm         = regexp.MustCompile(`^str w(\d+), \[x(\d+), #(0x[0-9a-f]+|[0-9]+)\]$`)
	storeBaseImm64       = regexp.MustCompile(`^str (x(?:\d+)|xzr), \[x(\d+), #(0x[0-9a-f]+|[0-9]+)\]$`)
	storeUXTW8           = regexp.MustCompile(`^strb w(\d+), \[x(\d+), w(\d+), uxtw\]$`)
	storeBaseImm8        = regexp.MustCompile(`^strb (w(?:\d+)|wzr), \[x(\d+)(?:, #(0x[0-9a-f]+|[0-9]+))?\]$`)
	storeSPImm8          = regexp.MustCompile(`^strb (w(?:\d+)|wzr), \[sp, #(0x[0-9a-f]+|[0-9]+)\]$`)
	storeUXTW32          = regexp.MustCompile(`^str w(\d+), \[x(\d+), w(\d+), uxtw #(\d+)\]$`)
	storeSXTW32          = regexp.MustCompile(`^str (w(?:\d+)|wzr), \[x(\d+), w(\d+), sxtw #(\d+)\]$`)
	storeXIndex          = regexp.MustCompile(`^str ([wx])(\d+), \[x(\d+), x(\d+)\]$`)
	storeXIndex8         = regexp.MustCompile(`^strb w(\d+), \[x(\d+), x(\d+)\]$`)
	extract32            = regexp.MustCompile(`^extr w(\d+), w(\d+), w(\d+), #(0x[0-9a-f]+|[0-9]+)$`)
	not32                = regexp.MustCompile(`^mvn w(\d+), w(\d+)$`)
	orNotLSL32           = regexp.MustCompile(`^orn w(\d+), w(\d+), w(\d+)(?:, lsl #(\d+))?$`)
	rotateRight32        = regexp.MustCompile(`^ror w(\d+), w(\d+), #(0x[0-9a-f]+|[0-9]+)$`)
	shiftRight32         = regexp.MustCompile(`^lsr w(\d+), w(\d+), #(0x[0-9a-f]+|[0-9]+)$`)
	shiftRight64         = regexp.MustCompile(`^lsr x(\d+), x(\d+), #(0x[0-9a-f]+|[0-9]+)$`)
	shiftRightReg32      = regexp.MustCompile(`^lsr w(\d+), w(\d+), w(\d+)$`)
	unsignedLong32       = regexp.MustCompile(`^umull x(\d+), w(\d+), w(\d+)$`)
	unsignedHigh64       = regexp.MustCompile(`^umulh x(\d+), x(\d+), x(\d+)$`)
	multiplyAdd32        = regexp.MustCompile(`^madd w(\d+), w(\d+), w(\d+), w(\d+)$`)
	multiplyAdd64        = regexp.MustCompile(`^madd x(\d+), x(\d+), x(\d+), x(\d+)$`)
	multiplySub32        = regexp.MustCompile(`^msub w(\d+), w(\d+), w(\d+), w(\d+)$`)
	unsignedDivide32     = regexp.MustCompile(`^udiv w(\d+), w(\d+), w(\d+)$`)
	adds64               = regexp.MustCompile(`^adds x(\d+), x(\d+), x(\d+)$`)
	adc32                = regexp.MustCompile(`^adc w(\d+), (w\d+|wzr), (w\d+|wzr)$`)
	compareNegative64    = regexp.MustCompile(`^cmn x(\d+), x(\d+)$`)
	compareNegativeImm32 = regexp.MustCompile(`^cmn w(\d+), #(0x[0-9a-f]+|[0-9]+)$`)
	compareImmediate32   = regexp.MustCompile(`^cmp w(\d+), #(0x[0-9a-f]+|[0-9]+)$`)
	compareImmediate64   = regexp.MustCompile(`^cmp x(\d+), #(0x[0-9a-f]+|[0-9]+)$`)
	conditionSetHI       = regexp.MustCompile(`^cset w(\d+), hi$`)
	conditionSetLO       = regexp.MustCompile(`^cset w(\d+), lo$`)
	conditionSetEQ       = regexp.MustCompile(`^cset w(\d+), eq$`)
	compareRegister32    = regexp.MustCompile(`^cmp w(\d+), w(\d+)$`)
	compareRegister64    = regexp.MustCompile(`^cmp x(\d+), x(\d+)$`)
	selectZeroEq32       = regexp.MustCompile(`^csel w(\d+), wzr, w(\d+), eq$`)
	branchAuth           = regexp.MustCompile(`^braa x\d+, x\d+$`)
	addressPage          = regexp.MustCompile(`^adrp x(\d+), #(0x[0-9a-f]+)$`)
	stackAdjust          = regexp.MustCompile(`^(add|sub) sp, sp, #(0x[0-9a-f]+|[0-9]+)(?:, lsl #(\d+))?$`)
	unscaledLoad         = regexp.MustCompile(`^ldur ([wx])(\d+), \[x(\d+), #(-?(?:0x[0-9a-f]+|[0-9]+))\]$`)
	unscaledStore        = regexp.MustCompile(`^stur ([wx])(\d+), \[x(\d+), #(-?(?:0x[0-9a-f]+|[0-9]+))\]$`)
	unscaledByteLoad     = regexp.MustCompile(`^ldurb w(\d+), \[x(\d+), #(-?(?:0x[0-9a-f]+|[0-9]+))\]$`)
	signedByteLoad       = regexp.MustCompile(`^ldrsb w(\d+), \[x(\d+)(?:, #(-?(?:0x[0-9a-f]+|[0-9]+)))?\]$`)
	unscaledByteStore    = regexp.MustCompile(`^sturb w(\d+), \[x(\d+), #(-?(?:0x[0-9a-f]+|[0-9]+))\]$`)
	halfLoad             = regexp.MustCompile(`^ldrh w(\d+), \[x(\d+)(?:, #(0x[0-9a-f]+|[0-9]+))?\]$`)
	halfStore            = regexp.MustCompile(`^strh w(\d+), \[x(\d+)(?:, #(0x[0-9a-f]+|[0-9]+))?\]$`)
	pairOffset           = regexp.MustCompile(`^(ldp|stp) ([xw](?:\d+)|[xw]zr), ([xw](?:\d+)|[xw]zr), \[(sp|x\d+)(?:, #(-?(?:0x[0-9a-f]+|[0-9]+)))?\](!)?$`)
	pairPostIndex        = regexp.MustCompile(`^(ldp|stp) ([xw](?:\d+)|[xw]zr), ([xw](?:\d+)|[xw]zr), \[(sp|x\d+)\], #(-?(?:0x[0-9a-f]+|[0-9]+))$`)
	bitClear             = regexp.MustCompile(`^bic ([wx])(\d+), [wx](\d+), [wx](\d+)(?:, lsl #(\d+))?$`)
	signExtendByte       = regexp.MustCompile(`^sxtb w(\d+), w(\d+)$`)
	signExtendWord       = regexp.MustCompile(`^sxtw x(\d+), w(\d+)$`)
	negate               = regexp.MustCompile(`^neg ([wx])(\d+), [wx](\d+)$`)
	testRegister         = regexp.MustCompile(`^tst w(\d+), w(\d+)$`)
	testImmediate        = regexp.MustCompile(`^tst w(\d+), #(0x[0-9a-f]+|[0-9]+)$`)
	conditionSet         = regexp.MustCompile(`^cset w(\d+), ([a-z]+)$`)
	conditionSelect      = regexp.MustCompile(`^csel ([wx])(\d+), ([wx](?:\d+)|[wx]zr), ([wx](?:\d+)|[wx]zr), ([a-z]+)$`)
	conditionIncrement   = regexp.MustCompile(`^cinc ([wx])(\d+), [wx](\d+), ([a-z]+)$`)
	conditionNegate      = regexp.MustCompile(`^cneg w(\d+), w(\d+), ([a-z]+)$`)
	noopInstruction      = regexp.MustCompile(`^pacibsp$`)
	blockFunction        = regexp.MustCompile(`^func ([A-Za-z_][A-Za-z0-9_]*)(?:\([^)]*\))`)
	branchDirect         = regexp.MustCompile(`^b #(0x[0-9a-f]+)$`)
	branchConditional    = regexp.MustCompile(`^b\.([a-z]+) #(0x[0-9a-f]+)$`)
	branchRegister       = regexp.MustCompile(`^braa x(\d+), x\d+$`)
	callRegister         = regexp.MustCompile(`^blraa x(\d+), x\d+$`)
	returnRegister       = regexp.MustCompile(`^ret(?:ab)?$`)
	literalReturn        = regexp.MustCompile(`^\s*return (0x[0-9a-f]+), false, nil$`)
	identifier           = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)
)

func number(value string) (uint64, error) {
	base := 10
	if strings.HasPrefix(value, "0x") {
		base, value = 16, strings.TrimPrefix(value, "0x")
	}
	return strconv.ParseUint(value, base, 64)
}

func signedNumber(value string) (int64, error) {
	negative := strings.HasPrefix(value, "-")
	if negative {
		value = strings.TrimPrefix(value, "-")
	}
	result, err := number(value)
	if err != nil || result > 0x7fff+uint64(boolInt(negative)) {
		return 0, fmt.Errorf("signed immediate %q is outside int16", value)
	}
	if negative {
		return -int64(result), nil
	}
	return int64(result), nil
}

func boolInt(value bool) int {
	if value {
		return 1
	}
	return 0
}

func conditionConstant(name string) (string, bool) {
	constants := map[string]string{
		"eq": "uint8(compactConditionEQ)", "ne": "uint8(compactConditionNE)",
		"hs": "uint8(compactConditionHS)", "lo": "uint8(compactConditionLO)",
		"mi": "uint8(compactConditionMI)", "pl": "uint8(compactConditionPL)",
		"vs": "uint8(compactConditionVS)", "vc": "uint8(compactConditionVC)",
		"hi": "uint8(compactConditionHI)", "ls": "uint8(compactConditionLS)",
		"ge": "uint8(compactConditionGE)", "lt": "uint8(compactConditionLT)",
		"gt": "uint8(compactConditionGT)", "le": "uint8(compactConditionLE)",
	}
	value, ok := constants[name]
	return value, ok
}

func registerNumber(value string) (string, bool) {
	if value == "wzr" || value == "xzr" {
		return "31", true
	}
	if len(value) < 2 || value[0] != 'w' && value[0] != 'x' {
		return "", false
	}
	if _, err := strconv.ParseUint(value[1:], 10, 8); err != nil {
		return "", false
	}
	return value[1:], true
}

func baseRegisterNumber(value string) (string, bool) {
	if value == "sp" {
		return "31", true
	}
	if !strings.HasPrefix(value, "x") {
		return "", false
	}
	return registerNumber(value)
}

func compactPairInstruction(match []string, mode uint8) (string, bool, error) {
	first, firstOK := registerNumber(match[2])
	second, secondOK := registerNumber(match[3])
	base, baseOK := baseRegisterNumber(match[4])
	if !firstOK || !secondOK || !baseOK || match[2][0] != match[3][0] {
		return "", false, fmt.Errorf("invalid register in pair operation")
	}
	offset := int64(0)
	var err error
	if match[5] != "" {
		offset, err = signedNumber(match[5])
	}
	if err != nil {
		return "", false, err
	}
	op := "compactPairLoad32"
	if match[1] == "stp" {
		op = "compactPairStore32"
	}
	if match[2][0] == 'x' {
		op = "compactPairLoad64"
		if match[1] == "stp" {
			op = "compactPairStore64"
		}
	}
	return fields(op, "destination", first, "left", base, "right", second, "accumulator", fmt.Sprint(mode), "immediate", fmt.Sprintf("%#x", uint16(int16(offset)))), false, nil
}

func fields(operation string, values ...string) string {
	result := "{op: " + operation
	for index := 0; index+1 < len(values); index += 2 {
		result += ", " + values[index] + ": " + values[index+1]
	}
	return result + "},"
}

func immediate32Fields(value uint64) []string {
	return []string{
		"right", fmt.Sprintf("%#x", byte(value>>16)),
		"accumulator", fmt.Sprintf("%#x", byte(value>>24)),
		"immediate", fmt.Sprintf("%#x", uint16(value)),
	}
}

func immediate48Fields(value uint64) []string {
	return []string{
		"left", fmt.Sprintf("%#x", byte(value>>16)),
		"right", fmt.Sprintf("%#x", byte(value>>24)),
		"accumulator", fmt.Sprintf("%#x", byte(value>>32)),
		"shift", fmt.Sprintf("%#x", byte(value>>40)),
		"immediate", fmt.Sprintf("%#x", uint16(value)),
	}
}

func compactInstruction(disassembly string) (string, bool, error) {
	if noopInstruction.MatchString(disassembly) {
		return fields("compactNoop"), false, nil
	}
	if match := addressPage.FindStringSubmatch(disassembly); match != nil {
		value, err := number(match[2])
		if err != nil || value >= 1<<48 {
			return "", false, fmt.Errorf("invalid ADRP address %q", match[2])
		}
		return fields("compactMovImm48", append([]string{"destination", match[1]}, immediate48Fields(value)...)...), false, nil
	}
	if match := stackAdjust.FindStringSubmatch(disassembly); match != nil {
		immediate, immediateErr := number(match[2])
		shift := uint64(0)
		var shiftErr error
		if match[3] != "" {
			shift, shiftErr = number(match[3])
		}
		if immediateErr != nil || shiftErr != nil || immediate > 0xffff || shift > 63 {
			return "", false, fmt.Errorf("invalid stack adjustment %q", disassembly)
		}
		op := "compactAddSPImm"
		if match[1] == "sub" {
			op = "compactSubSPImm"
		}
		return fields(op, "shift", fmt.Sprint(shift), "immediate", fmt.Sprintf("%#x", immediate)), false, nil
	}
	if match := unscaledLoad.FindStringSubmatch(disassembly); match != nil {
		offset, err := signedNumber(match[4])
		if err != nil {
			return "", false, err
		}
		op := "compactLoadSignedOffset32"
		if match[1] == "x" {
			op = "compactLoadSignedOffset64"
		}
		return fields(op, "destination", match[2], "left", match[3], "immediate", fmt.Sprintf("%#x", uint16(int16(offset)))), false, nil
	}
	if match := unscaledStore.FindStringSubmatch(disassembly); match != nil {
		offset, err := signedNumber(match[4])
		if err != nil {
			return "", false, err
		}
		op := "compactStoreSignedOffset32"
		if match[1] == "x" {
			op = "compactStoreSignedOffset64"
		}
		return fields(op, "destination", match[2], "left", match[3], "immediate", fmt.Sprintf("%#x", uint16(int16(offset)))), false, nil
	}
	if match := unscaledByteLoad.FindStringSubmatch(disassembly); match != nil {
		offset, err := signedNumber(match[3])
		if err != nil {
			return "", false, err
		}
		return fields("compactLoadSignedOffset8", "destination", match[1], "left", match[2], "immediate", fmt.Sprintf("%#x", uint16(int16(offset)))), false, nil
	}
	if match := signedByteLoad.FindStringSubmatch(disassembly); match != nil {
		offset := int64(0)
		var err error
		if match[3] != "" {
			offset, err = signedNumber(match[3])
		}
		if err != nil {
			return "", false, err
		}
		return fields("compactLoadSignedByte32", "destination", match[1], "left", match[2], "immediate", fmt.Sprintf("%#x", uint16(int16(offset)))), false, nil
	}
	if match := unscaledByteStore.FindStringSubmatch(disassembly); match != nil {
		offset, err := signedNumber(match[3])
		if err != nil {
			return "", false, err
		}
		return fields("compactStoreSignedOffset8", "destination", match[1], "left", match[2], "immediate", fmt.Sprintf("%#x", uint16(int16(offset)))), false, nil
	}
	if match := pairOffset.FindStringSubmatch(disassembly); match != nil {
		mode := uint8(0)
		if match[6] == "!" {
			mode = 1
		}
		return compactPairInstruction(match, mode)
	}
	if match := pairPostIndex.FindStringSubmatch(disassembly); match != nil {
		return compactPairInstruction(match, 2)
	}
	if match := halfLoad.FindStringSubmatch(disassembly); match != nil {
		offset := uint64(0)
		var err error
		if match[3] != "" {
			offset, err = number(match[3])
		}
		if err != nil || offset > 0xffff {
			return "", false, fmt.Errorf("invalid halfword load %q", disassembly)
		}
		op := "compactLoadHalf32"
		if offset != 0 {
			op = "compactLoadSignedOffsetHalf32"
		}
		return fields(op, "destination", match[1], "left", match[2], "immediate", fmt.Sprintf("%#x", offset)), false, nil
	}
	if match := halfStore.FindStringSubmatch(disassembly); match != nil {
		offset := uint64(0)
		var err error
		if match[3] != "" {
			offset, err = number(match[3])
		}
		if err != nil || offset > 0xffff {
			return "", false, fmt.Errorf("invalid halfword store %q", disassembly)
		}
		op := "compactStoreHalf32"
		if offset != 0 {
			op = "compactStoreSignedOffsetHalf32"
		}
		return fields(op, "destination", match[1], "left", match[2], "immediate", fmt.Sprintf("%#x", offset)), false, nil
	}
	if match := bitClear.FindStringSubmatch(disassembly); match != nil {
		op := "compactBic32"
		if match[1] == "x" {
			op = "compactBic64"
		}
		shift := "0"
		if match[5] != "" {
			shift = match[5]
		}
		return fields(op, "destination", match[2], "left", match[3], "right", match[4], "shift", shift), false, nil
	}
	if match := signExtendByte.FindStringSubmatch(disassembly); match != nil {
		return fields("compactSXT8To32", "destination", match[1], "left", match[2]), false, nil
	}
	if match := signExtendWord.FindStringSubmatch(disassembly); match != nil {
		return fields("compactSXT32To64", "destination", match[1], "left", match[2]), false, nil
	}
	if match := negate.FindStringSubmatch(disassembly); match != nil {
		op := "compactNeg32"
		if match[1] == "x" {
			op = "compactNeg64"
		}
		return fields(op, "destination", match[2], "left", match[3]), false, nil
	}
	if match := testRegister.FindStringSubmatch(disassembly); match != nil {
		return fields("compactTst32", "left", match[1], "right", match[2]), false, nil
	}
	if match := testImmediate.FindStringSubmatch(disassembly); match != nil {
		immediate, err := number(match[2])
		if err != nil || immediate > 0xffffffff {
			return "", false, fmt.Errorf("invalid TST immediate %q", match[2])
		}
		return fields("compactTstImm32", append([]string{"left", match[1]}, immediate32Fields(immediate)...)...), false, nil
	}
	if match := conditionIncrement.FindStringSubmatch(disassembly); match != nil {
		condition, ok := conditionConstant(match[4])
		if !ok {
			return "", false, fmt.Errorf("unsupported CINC condition %q", match[4])
		}
		op := "compactCInc32"
		if match[1] == "x" {
			op = "compactCInc64"
		}
		return fields(op, "destination", match[2], "left", match[3], "shift", condition), false, nil
	}
	if match := conditionNegate.FindStringSubmatch(disassembly); match != nil {
		condition, ok := conditionConstant(match[3])
		if !ok {
			return "", false, fmt.Errorf("unsupported CNEG condition %q", match[3])
		}
		return fields("compactCNeg32", "destination", match[1], "left", match[2], "shift", condition), false, nil
	}
	if match := movNegOne32.FindStringSubmatch(disassembly); match != nil {
		return fields("compactMovNegOne32", "destination", match[1]), false, nil
	}
	if match := movNegative32.FindStringSubmatch(disassembly); match != nil {
		magnitude, err := number(match[2])
		if err != nil || magnitude > 0xffffffff {
			return "", false, fmt.Errorf("invalid negative MOV immediate %q", match[2])
		}
		value := uint64(uint32(0 - uint32(magnitude)))
		return fields("compactMovImm32Wide", append([]string{"destination", match[1]}, immediate32Fields(value)...)...), false, nil
	}
	if match := movNegative64.FindStringSubmatch(disassembly); match != nil {
		magnitude, err := number(match[2])
		if err != nil || magnitude > 0xffff {
			return "", false, fmt.Errorf("invalid negative 64-bit MOV immediate %q", match[2])
		}
		return fields("compactMovNegImm64", "destination", match[1], "immediate", fmt.Sprintf("%#x", magnitude)), false, nil
	}
	if match := movImmediate.FindStringSubmatch(disassembly); match != nil {
		value, err := number(match[3])
		if err != nil || match[1] == "w" && value > 0xffffffff {
			return "", false, fmt.Errorf("invalid MOV immediate %q", match[3])
		}
		op := "compactMovImm32"
		if match[1] == "x" {
			if value > 0xffff {
				repeated := uint64(uint16(value))
				if value == repeated|repeated<<16|repeated<<32|repeated<<48 {
					return fields("compactMovImm64Repeat16", "destination", match[2], "immediate", fmt.Sprintf("%#x", uint16(value))), false, nil
				}
				for shift := uint64(16); shift <= 48; shift += 16 {
					if value == uint64(uint16(value>>shift))<<shift {
						return fields("compactMovImm64Shift", "destination", match[2], "shift", fmt.Sprint(shift), "immediate", fmt.Sprintf("%#x", uint16(value>>shift))), false, nil
					}
				}
				return "", false, fmt.Errorf("unsupported 64-bit MOV immediate %q", match[3])
			}
			op = "compactMovImm64"
		} else if value > 0xffff {
			if value&0xffff != 0 {
				return fields("compactMovImm32Wide", append([]string{"destination", match[2]}, immediate32Fields(value)...)...), false, nil
			}
			op, value = "compactMovImm32High", value>>16
		}
		return fields(op, "destination", match[2], "immediate", fmt.Sprintf("%#x", value)), false, nil
	}
	if match := movRegister.FindStringSubmatch(disassembly); match != nil {
		op := "compactMov32"
		if match[1] == "x" {
			op = "compactMov64"
		}
		return fields(op, "destination", match[2], "left", match[3]), false, nil
	}
	if match := movKeep.FindStringSubmatch(disassembly); match != nil {
		value, valueErr := number(match[3])
		shift := uint64(0)
		var shiftErr error
		if match[4] != "" {
			shift, shiftErr = number(match[4])
		}
		if valueErr != nil || shiftErr != nil || value > 0xffff || shift > 63 {
			return "", false, fmt.Errorf("invalid MOVK %q", disassembly)
		}
		op := "compactMovK32"
		if match[1] == "x" {
			op = "compactMovK64"
		}
		return fields(op, "destination", match[2], "shift", fmt.Sprint(shift), "immediate", fmt.Sprintf("%#x", value)), false, nil
	}
	if match := binary32.FindStringSubmatch(disassembly); match != nil {
		opcodes := map[string]string{
			"eor": "compactXor32", "and": "compactAnd32", "add": "compactAdd32",
			"sub": "compactSub32", "orr": "compactOr32", "mul": "compactMul32",
		}
		op := opcodes[match[1]]
		extra := []string{}
		if match[5] != "" {
			if match[1] != "and" && match[1] != "eor" && match[1] != "add" && match[1] != "sub" && match[1] != "orr" {
				return "", false, fmt.Errorf("unsupported shifted binary operation %q", disassembly)
			}
			switch {
			case match[1] == "and" && match[5] == "lsl":
				op = "compactAndLSL32"
			case match[1] == "eor" && match[5] == "lsl":
				op = "compactXorLSL32"
			case match[1] == "eor" && match[5] == "lsr":
				op = "compactXorLSR32"
			case match[1] == "add" && match[5] == "lsl":
				op = "compactAddLSL32"
			case match[1] == "sub" && match[5] == "lsl":
				op = "compactSubLSL32"
			case match[1] == "orr" && match[5] == "lsl":
				op = "compactOrLSL32"
			case match[1] == "eor" && match[5] == "ror":
				op = "compactXorROR32"
			case match[1] == "orr" && match[5] == "ror":
				op = "compactOrROR32"
			default:
				return "", false, fmt.Errorf("unsupported shifted binary operation %q", disassembly)
			}
			extra = append(extra, "shift", match[6])
		}
		return fields(op, append([]string{"destination", match[2], "left", match[3], "right", match[4]}, extra...)...), false, nil
	}
	if match := binary64.FindStringSubmatch(disassembly); match != nil {
		opcodes := map[string]string{
			"eor": "compactXor64", "add": "compactAdd64", "sub": "compactSub64",
			"mul": "compactMul64", "and": "compactAnd64", "orr": "compactOr64",
		}
		op := opcodes[match[1]]
		return fields(op, "destination", match[2], "left", match[3], "right", match[4]), false, nil
	}
	if match := binaryShift64.FindStringSubmatch(disassembly); match != nil {
		op := "compactAndLSL64"
		if match[1] == "eor" {
			op = "compactXorLSL64"
		} else if match[1] == "add" {
			op = "compactAddLSL64"
		}
		return fields(op, "destination", match[2], "left", match[3], "right", match[4], "shift", match[5]), false, nil
	}
	if match := addImmediate32.FindStringSubmatch(disassembly); match != nil {
		immediate, immediateErr := number(match[3])
		shift := uint64(0)
		var shiftErr error
		if match[4] != "" {
			shift, shiftErr = number(match[4])
		}
		if immediateErr != nil || shiftErr != nil || immediate > 0xffff || shift > 31 {
			return "", false, fmt.Errorf("invalid immediate ADD %q", disassembly)
		}
		return fields("compactAddImmShift32", "destination", match[1], "left", match[2], "shift", fmt.Sprint(shift), "immediate", fmt.Sprintf("%#x", immediate)), false, nil
	}
	if match := addImmediate64.FindStringSubmatch(disassembly); match != nil {
		immediate, err := number(match[3])
		if err != nil || immediate > 0xffff {
			return "", false, fmt.Errorf("invalid 64-bit immediate ADD %q", disassembly)
		}
		return fields("compactAddImm64", "destination", match[1], "left", match[2], "immediate", fmt.Sprintf("%#x", immediate)), false, nil
	}
	if match := subImmediate32.FindStringSubmatch(disassembly); match != nil {
		immediate, immediateErr := number(match[3])
		shift := uint64(0)
		var shiftErr error
		if match[4] != "" {
			shift, shiftErr = number(match[4])
		}
		if immediateErr != nil || shiftErr != nil || immediate > 0xffff || shift > 31 {
			return "", false, fmt.Errorf("invalid immediate SUB %q", disassembly)
		}
		return fields("compactSubImmShift32", "destination", match[1], "left", match[2], "shift", fmt.Sprint(shift), "immediate", fmt.Sprintf("%#x", immediate)), false, nil
	}
	if match := subImmediate64.FindStringSubmatch(disassembly); match != nil {
		immediate, err := number(match[3])
		if err != nil || immediate > 0xffff {
			return "", false, fmt.Errorf("invalid 64-bit immediate SUB %q", disassembly)
		}
		return fields("compactSubImm64", "destination", match[1], "left", match[2], "immediate", fmt.Sprintf("%#x", immediate)), false, nil
	}
	if match := xorImmediate32.FindStringSubmatch(disassembly); match != nil {
		immediate, err := number(match[3])
		if err != nil || immediate > 0xffffffff {
			return "", false, fmt.Errorf("invalid immediate EOR %q", disassembly)
		}
		return fields("compactXorImm32", append([]string{"destination", match[1], "left", match[2]}, immediate32Fields(immediate)...)...), false, nil
	}
	if match := andImmediate32.FindStringSubmatch(disassembly); match != nil {
		immediate, err := number(match[3])
		if err != nil || immediate > 0xffffffff {
			return "", false, fmt.Errorf("invalid immediate AND %q", disassembly)
		}
		return fields("compactAndImm32", append([]string{"destination", match[1], "left", match[2]}, immediate32Fields(immediate)...)...), false, nil
	}
	if match := logicalImmediate64.FindStringSubmatch(disassembly); match != nil {
		immediate, err := number(match[4])
		if err != nil || immediate > 0xffff {
			return "", false, fmt.Errorf("invalid 64-bit logical immediate %q", disassembly)
		}
		op := "compactOrImm64"
		if match[1] == "eor" {
			op = "compactXorImm64"
		}
		return fields(op, "destination", match[2], "left", match[3], "immediate", fmt.Sprintf("%#x", immediate)), false, nil
	}
	if match := shiftLeft32.FindStringSubmatch(disassembly); match != nil {
		shift, err := number(match[3])
		if err != nil || shift > 31 {
			return "", false, fmt.Errorf("invalid LSL %q", disassembly)
		}
		return fields("compactLSLImm32", "destination", match[1], "left", match[2], "shift", fmt.Sprint(shift)), false, nil
	}
	if match := shiftLeft64.FindStringSubmatch(disassembly); match != nil {
		shift, err := number(match[3])
		if err != nil || shift > 63 {
			return "", false, fmt.Errorf("invalid 64-bit LSL %q", disassembly)
		}
		return fields("compactLSLImm64", "destination", match[1], "left", match[2], "shift", fmt.Sprint(shift)), false, nil
	}
	if match := shiftLeftReg32.FindStringSubmatch(disassembly); match != nil {
		return fields("compactLSLReg32", "destination", match[1], "left", match[2], "right", match[3]), false, nil
	}
	if match := addSP.FindStringSubmatch(disassembly); match != nil {
		immediate, immediateErr := number(match[2])
		shift := uint64(0)
		var shiftErr error
		if match[3] != "" {
			shift, shiftErr = number(match[3])
		}
		if immediateErr != nil || shiftErr != nil || immediate > 0xffff || shift > 63 {
			return "", false, fmt.Errorf("invalid SP add %q", disassembly)
		}
		return fields("compactAddSPShiftImm64", "destination", match[1], "shift", fmt.Sprint(shift), "immediate", fmt.Sprintf("%#x", immediate)), false, nil
	}
	if match := loadSP.FindStringSubmatch(disassembly); match != nil {
		immediate, err := number(match[3])
		if err != nil || immediate > 0xffff {
			return "", false, fmt.Errorf("invalid SP load %q", disassembly)
		}
		op := "compactLoadSP32"
		if match[1] == "x" {
			op = "compactLoadSP64"
		}
		return fields(op, "destination", match[2], "immediate", fmt.Sprintf("%#x", immediate)), false, nil
	}
	if match := loadBase.FindStringSubmatch(disassembly); match != nil {
		immediate, err := number(match[3])
		if err != nil || immediate > 0xffff {
			return "", false, fmt.Errorf("invalid base load %q", disassembly)
		}
		return fields("compactLoadBaseImm32", "destination", match[1], "left", match[2], "immediate", fmt.Sprintf("%#x", immediate)), false, nil
	}
	if match := loadBase32.FindStringSubmatch(disassembly); match != nil {
		return fields("compactLoadBase32", "destination", match[1], "left", match[2]), false, nil
	}
	if match := loadBaseImm64.FindStringSubmatch(disassembly); match != nil {
		immediate, err := number(match[3])
		if err != nil || immediate > 0xffff {
			return "", false, fmt.Errorf("invalid 64-bit base load %q", disassembly)
		}
		return fields("compactLoadBaseImm64", "destination", match[1], "left", match[2], "immediate", fmt.Sprintf("%#x", immediate)), false, nil
	}
	if match := loadBase64.FindStringSubmatch(disassembly); match != nil {
		return fields("compactLoadBase64", "destination", match[1], "left", match[2]), false, nil
	}
	if match := loadSXTW64.FindStringSubmatch(disassembly); match != nil {
		return fields("compactLoadSXTW64", "destination", match[1], "left", match[2], "right", match[3], "shift", match[4]), false, nil
	}
	if match := loadUXTW32.FindStringSubmatch(disassembly); match != nil {
		return fields("compactLoadUXTW32", "destination", match[1], "left", match[2], "right", match[3], "shift", match[4]), false, nil
	}
	if match := loadUXTW8.FindStringSubmatch(disassembly); match != nil {
		return fields("compactLoadUXTW8", "destination", match[1], "left", match[2], "right", match[3]), false, nil
	}
	if match := loadBaseImm8.FindStringSubmatch(disassembly); match != nil {
		immediate := uint64(0)
		var err error
		if match[3] != "" {
			immediate, err = number(match[3])
		}
		if err != nil || immediate > 0xffff {
			return "", false, fmt.Errorf("invalid byte load %q", disassembly)
		}
		return fields("compactLoadBaseImm8", "destination", match[1], "left", match[2], "immediate", fmt.Sprintf("%#x", immediate)), false, nil
	}
	if match := loadSPImm8.FindStringSubmatch(disassembly); match != nil {
		immediate, err := number(match[2])
		if err != nil || immediate > 0xffff {
			return "", false, fmt.Errorf("invalid SP byte load %q", disassembly)
		}
		return fields("compactLoadSPImm8", "destination", match[1], "immediate", fmt.Sprintf("%#x", immediate)), false, nil
	}
	if match := loadXIndex8.FindStringSubmatch(disassembly); match != nil {
		return fields("compactLoadXIndex8", "destination", match[1], "left", match[2], "right", match[3]), false, nil
	}
	if match := loadPostIndex8.FindStringSubmatch(disassembly); match != nil {
		immediate, err := number(match[3])
		if err != nil || immediate > 0xffff {
			return "", false, fmt.Errorf("invalid post-index load %q", disassembly)
		}
		return fields("compactLoadPostIndex8", "destination", match[1], "left", match[2], "immediate", fmt.Sprintf("%#x", immediate)), false, nil
	}
	if match := loadXIndex.FindStringSubmatch(disassembly); match != nil {
		op := "compactLoadXIndex32"
		if match[1] == "x" {
			op = "compactLoadXIndex64"
		}
		return fields(op, "destination", match[2], "left", match[3], "right", match[4]), false, nil
	}
	if match := loadSXTW32.FindStringSubmatch(disassembly); match != nil {
		return fields("compactLoadSXTW32", "destination", match[1], "left", match[2], "right", match[3], "shift", match[4]), false, nil
	}
	if match := loadUXTW64.FindStringSubmatch(disassembly); match != nil {
		return fields("compactLoadUXTW64", "destination", match[1], "left", match[2], "right", match[3], "shift", match[4]), false, nil
	}
	if match := storeSP.FindStringSubmatch(disassembly); match != nil {
		immediate, err := number(match[2])
		if err != nil || immediate > 0xffff {
			return "", false, fmt.Errorf("invalid SP store %q", disassembly)
		}
		return fields("compactStoreSP32", "left", match[1], "immediate", fmt.Sprintf("%#x", immediate)), false, nil
	}
	if match := storeSP64.FindStringSubmatch(disassembly); match != nil {
		immediate, err := number(match[2])
		if err != nil || immediate > 0xffff {
			return "", false, fmt.Errorf("invalid 64-bit SP store %q", disassembly)
		}
		return fields("compactStoreSP64", "left", match[1], "immediate", fmt.Sprintf("%#x", immediate)), false, nil
	}
	if match := storeBase.FindStringSubmatch(disassembly); match != nil {
		return fields("compactStoreBase32", "destination", match[1], "left", match[2]), false, nil
	}
	if match := storeBase64.FindStringSubmatch(disassembly); match != nil {
		return fields("compactStoreBase64", "destination", match[1], "left", match[2]), false, nil
	}
	if match := storeBaseImm.FindStringSubmatch(disassembly); match != nil {
		immediate, err := number(match[3])
		if err != nil || immediate > 0xffff {
			return "", false, fmt.Errorf("invalid base store %q", disassembly)
		}
		return fields("compactStoreBaseImm32", "destination", match[1], "left", match[2], "immediate", fmt.Sprintf("%#x", immediate)), false, nil
	}
	if match := storeBaseImm64.FindStringSubmatch(disassembly); match != nil {
		source, ok := registerNumber(match[1])
		immediate, err := number(match[3])
		if !ok || err != nil || immediate > 0xffff {
			return "", false, fmt.Errorf("invalid 64-bit base store %q", disassembly)
		}
		return fields("compactStoreBaseImm64", "destination", source, "left", match[2], "immediate", fmt.Sprintf("%#x", immediate)), false, nil
	}
	if match := storeUXTW8.FindStringSubmatch(disassembly); match != nil {
		return fields("compactStoreUXTW8", "destination", match[1], "left", match[2], "right", match[3]), false, nil
	}
	if match := storeUXTW32.FindStringSubmatch(disassembly); match != nil {
		return fields("compactStoreUXTW32", "destination", match[1], "left", match[2], "right", match[3], "shift", match[4]), false, nil
	}
	if match := storeSXTW32.FindStringSubmatch(disassembly); match != nil {
		source, ok := registerNumber(match[1])
		if !ok {
			return "", false, fmt.Errorf("invalid signed-index store %q", disassembly)
		}
		return fields("compactStoreSXTW32", "destination", source, "left", match[2], "right", match[3], "shift", match[4]), false, nil
	}
	if match := storeBaseImm8.FindStringSubmatch(disassembly); match != nil {
		source, ok := registerNumber(match[1])
		immediate := uint64(0)
		var err error
		if match[3] != "" {
			immediate, err = number(match[3])
		}
		if !ok || err != nil || immediate > 0xffff {
			return "", false, fmt.Errorf("invalid byte store %q", disassembly)
		}
		return fields("compactStoreBaseImm8", "destination", source, "left", match[2], "immediate", fmt.Sprintf("%#x", immediate)), false, nil
	}
	if match := storeSPImm8.FindStringSubmatch(disassembly); match != nil {
		source, ok := registerNumber(match[1])
		immediate, err := number(match[2])
		if !ok || err != nil || immediate > 0xffff {
			return "", false, fmt.Errorf("invalid SP byte store %q", disassembly)
		}
		return fields("compactStoreSPImm8", "destination", source, "immediate", fmt.Sprintf("%#x", immediate)), false, nil
	}
	if match := storeXIndex.FindStringSubmatch(disassembly); match != nil {
		op := "compactStoreXIndex32"
		if match[1] == "x" {
			op = "compactStoreXIndex64"
		}
		return fields(op, "destination", match[2], "left", match[3], "right", match[4]), false, nil
	}
	if match := storeXIndex8.FindStringSubmatch(disassembly); match != nil {
		return fields("compactStoreXIndex8", "destination", match[1], "left", match[2], "right", match[3]), false, nil
	}
	if match := extract32.FindStringSubmatch(disassembly); match != nil {
		shift, err := number(match[4])
		if err != nil || shift > 31 {
			return "", false, fmt.Errorf("invalid EXTR %q", disassembly)
		}
		return fields("compactExtract32", "destination", match[1], "left", match[2], "right", match[3], "shift", fmt.Sprint(shift)), false, nil
	}
	if match := not32.FindStringSubmatch(disassembly); match != nil {
		return fields("compactNot32", "destination", match[1], "left", match[2]), false, nil
	}
	if match := orNotLSL32.FindStringSubmatch(disassembly); match != nil {
		shift := "0"
		if match[4] != "" {
			shift = match[4]
		}
		return fields("compactOrNotLSL32", "destination", match[1], "left", match[2], "right", match[3], "shift", shift), false, nil
	}
	if match := rotateRight32.FindStringSubmatch(disassembly); match != nil {
		shift, err := number(match[3])
		if err != nil || shift > 31 {
			return "", false, fmt.Errorf("invalid ROR %q", disassembly)
		}
		return fields("compactRORImm32", "destination", match[1], "left", match[2], "shift", fmt.Sprint(shift)), false, nil
	}
	if match := shiftRight32.FindStringSubmatch(disassembly); match != nil {
		shift, err := number(match[3])
		if err != nil || shift > 31 {
			return "", false, fmt.Errorf("invalid LSR %q", disassembly)
		}
		return fields("compactLSRImm32", "destination", match[1], "left", match[2], "shift", fmt.Sprint(shift)), false, nil
	}
	if match := shiftRight64.FindStringSubmatch(disassembly); match != nil {
		shift, err := number(match[3])
		if err != nil || shift > 63 {
			return "", false, fmt.Errorf("invalid 64-bit LSR %q", disassembly)
		}
		return fields("compactLSRImm64", "destination", match[1], "left", match[2], "shift", fmt.Sprint(shift)), false, nil
	}
	if match := shiftRightReg32.FindStringSubmatch(disassembly); match != nil {
		return fields("compactLSRReg32", "destination", match[1], "left", match[2], "right", match[3]), false, nil
	}
	if match := unsignedLong32.FindStringSubmatch(disassembly); match != nil {
		return fields("compactUMull32", "destination", match[1], "left", match[2], "right", match[3]), false, nil
	}
	if match := unsignedHigh64.FindStringSubmatch(disassembly); match != nil {
		return fields("compactUMulH64", "destination", match[1], "left", match[2], "right", match[3]), false, nil
	}
	if match := multiplyAdd32.FindStringSubmatch(disassembly); match != nil {
		return fields("compactMAdd32", "destination", match[1], "left", match[2], "right", match[3], "accumulator", match[4]), false, nil
	}
	if match := multiplyAdd64.FindStringSubmatch(disassembly); match != nil {
		return fields("compactMAdd64", "destination", match[1], "left", match[2], "right", match[3], "accumulator", match[4]), false, nil
	}
	if match := multiplySub32.FindStringSubmatch(disassembly); match != nil {
		return fields("compactMSub32", "destination", match[1], "left", match[2], "right", match[3], "accumulator", match[4]), false, nil
	}
	if match := unsignedDivide32.FindStringSubmatch(disassembly); match != nil {
		return fields("compactUDiv32", "destination", match[1], "left", match[2], "right", match[3]), false, nil
	}
	if match := adds64.FindStringSubmatch(disassembly); match != nil {
		return fields("compactAdds64", "destination", match[1], "left", match[2], "right", match[3]), false, nil
	}
	if match := adc32.FindStringSubmatch(disassembly); match != nil {
		reg := func(value string) string {
			if value == "wzr" {
				return "31"
			}
			return strings.TrimPrefix(value, "w")
		}
		return fields("compactADC32", "destination", match[1], "left", reg(match[2]), "right", reg(match[3])), false, nil
	}
	if match := compareNegative64.FindStringSubmatch(disassembly); match != nil {
		return fields("compactCMN64", "left", match[1], "right", match[2]), false, nil
	}
	if match := compareNegativeImm32.FindStringSubmatch(disassembly); match != nil {
		immediate, err := number(match[2])
		if err != nil || immediate > 0xffff {
			return "", false, fmt.Errorf("invalid immediate CMN %q", disassembly)
		}
		return fields("compactCMNImm32", "left", match[1], "immediate", fmt.Sprintf("%#x", immediate)), false, nil
	}
	if match := compareImmediate32.FindStringSubmatch(disassembly); match != nil {
		immediate, err := number(match[2])
		if err != nil || immediate > 0xffff {
			return "", false, fmt.Errorf("invalid immediate CMP %q", disassembly)
		}
		return fields("compactCmpImm32", "left", match[1], "immediate", fmt.Sprintf("%#x", immediate)), false, nil
	}
	if match := compareImmediate64.FindStringSubmatch(disassembly); match != nil {
		immediate, err := number(match[2])
		if err != nil || immediate > 0xffff {
			return "", false, fmt.Errorf("invalid 64-bit immediate CMP %q", disassembly)
		}
		return fields("compactCmpImm64", "left", match[1], "immediate", fmt.Sprintf("%#x", immediate)), false, nil
	}
	if match := conditionSetHI.FindStringSubmatch(disassembly); match != nil {
		return fields("compactCSet32HI", "destination", match[1]), false, nil
	}
	if match := conditionSetLO.FindStringSubmatch(disassembly); match != nil {
		return fields("compactCSet32LO", "destination", match[1]), false, nil
	}
	if match := conditionSetEQ.FindStringSubmatch(disassembly); match != nil {
		return fields("compactCSet32EQ", "destination", match[1]), false, nil
	}
	if match := compareRegister32.FindStringSubmatch(disassembly); match != nil {
		return fields("compactCmp32", "left", match[1], "right", match[2]), false, nil
	}
	if match := compareRegister64.FindStringSubmatch(disassembly); match != nil {
		return fields("compactCmp64", "left", match[1], "right", match[2]), false, nil
	}
	if match := selectZeroEq32.FindStringSubmatch(disassembly); match != nil {
		return fields("compactCSELZeroEq32", "destination", match[1], "right", match[2]), false, nil
	}
	if match := conditionSet.FindStringSubmatch(disassembly); match != nil {
		condition, ok := conditionConstant(match[2])
		if !ok {
			return "", false, fmt.Errorf("unsupported CSET condition %q", match[2])
		}
		return fields("compactCSet32", "destination", match[1], "shift", condition), false, nil
	}
	if match := conditionSelect.FindStringSubmatch(disassembly); match != nil {
		left, leftOK := registerNumber(match[3])
		right, rightOK := registerNumber(match[4])
		condition, conditionOK := conditionConstant(match[5])
		if !leftOK || !rightOK || !conditionOK || match[1] != string(match[3][0]) || match[1] != string(match[4][0]) {
			return "", false, fmt.Errorf("invalid CSEL %q", disassembly)
		}
		op := "compactCSEL32"
		if match[1] == "x" {
			op = "compactCSEL64"
		}
		return fields(op, "destination", match[2], "left", left, "right", right, "shift", condition), false, nil
	}
	if branchAuth.MatchString(disassembly) {
		return "", true, nil
	}
	return "", false, fmt.Errorf("unsupported instruction %q", disassembly)
}

type annotatedInstruction struct {
	pc          uint64
	disassembly string
}

type annotatedBlock struct {
	name         string
	start        uint64
	instructions []annotatedInstruction
	returnTarget uint64
}

func parseAnnotatedBlocks(source, prefix string) ([]annotatedBlock, error) {
	var blocks []annotatedBlock
	var current *annotatedBlock
	depth := 0
	scanner := bufio.NewScanner(strings.NewReader(source))
	for scanner.Scan() {
		line := scanner.Text()
		if current == nil {
			match := blockFunction.FindStringSubmatch(line)
			if match == nil || !strings.HasPrefix(match[1], prefix) {
				continue
			}
			current = &annotatedBlock{name: match[1]}
			depth = strings.Count(line, "{") - strings.Count(line, "}")
			continue
		}
		if match := commentPattern.FindStringSubmatch(line); match != nil {
			pc, err := strconv.ParseUint(strings.TrimPrefix(match[1], "0x"), 16, 64)
			if err != nil {
				return nil, err
			}
			if len(current.instructions) == 0 {
				current.start = pc
			}
			current.instructions = append(current.instructions, annotatedInstruction{pc: pc, disassembly: match[2]})
		}
		if match := literalReturn.FindStringSubmatch(line); match != nil {
			target, err := number(match[1])
			if err != nil {
				return nil, err
			}
			current.returnTarget = target
		}
		depth += strings.Count(line, "{") - strings.Count(line, "}")
		if depth == 0 {
			if len(current.instructions) != 0 {
				blocks = append(blocks, *current)
			}
			current = nil
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if current != nil {
		return nil, fmt.Errorf("unterminated generated block %s", current.name)
	}
	if len(blocks) == 0 {
		return nil, fmt.Errorf("input contains no functions with prefix %q", prefix)
	}
	return blocks, nil
}

func compactTerminator(instruction annotatedInstruction) (string, error) {
	if match := branchDirect.FindStringSubmatch(instruction.disassembly); match != nil {
		return "terminator: compactBlockDirect, target: " + match[1], nil
	}
	if match := branchConditional.FindStringSubmatch(instruction.disassembly); match != nil {
		if match[1] != "ne" {
			return "", fmt.Errorf("unsupported branch condition %q", match[1])
		}
		target, err := number(match[2])
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("terminator: compactBlockConditional, condition: compactConditionNE, target: %#x, alternate: %#x", target, instruction.pc+4), nil
	}
	if match := branchRegister.FindStringSubmatch(instruction.disassembly); match != nil {
		return "terminator: compactBlockRegister, register: " + match[1], nil
	}
	if match := callRegister.FindStringSubmatch(instruction.disassembly); match != nil {
		return "terminator: compactBlockCallRegister, register: " + match[1], nil
	}
	if returnRegister.MatchString(instruction.disassembly) {
		return "terminator: compactBlockReturn", nil
	}
	return "", fmt.Errorf("unsupported block terminator %q", instruction.disassembly)
}

func generateCompactCFGSource(name string, blocks []annotatedBlock, limit int) ([]byte, int, map[string]int, error) {
	var literals []string
	skipped := 0
	reasons := make(map[string]int)
	for _, block := range blocks {
		if limit > 0 && len(literals) >= limit {
			break
		}
		if len(block.instructions) == 0 {
			skipped++
			reasons["empty block"]++
			continue
		}
		terminal := block.instructions[len(block.instructions)-1]
		body := block.instructions[:len(block.instructions)-1]
		terminator, err := compactTerminator(terminal)
		if err != nil && block.returnTarget != 0 {
			terminator = fmt.Sprintf("terminator: compactBlockDirect, target: %#x", block.returnTarget)
			body = block.instructions
			err = nil
		}
		if err != nil {
			skipped++
			reasons["terminator "+strings.Fields(terminal.disassembly)[0]]++
			continue
		}
		expected := block.start
		var instructions []string
		compatible := true
		for _, entry := range body {
			if entry.pc != expected {
				compatible = false
				reasons["disjoint instruction addresses"]++
				break
			}
			literal, stop, err := compactInstruction(entry.disassembly)
			if err != nil || stop {
				compatible = false
				reasons["instruction "+strings.Fields(entry.disassembly)[0]]++
				break
			}
			instructions = append(instructions, literal)
			expected += 4
		}
		terminalPC := terminal.pc
		if len(body) == len(block.instructions) {
			terminalPC += 4
		}
		if !compatible || terminalPC != expected {
			skipped++
			if compatible {
				reasons["disjoint terminator address"]++
			}
			continue
		}
		literals = append(literals, fmt.Sprintf("%#x: {start: %#x, instructions: []compactInstruction{\n\t\t%s\n\t}, %s},", block.start, block.start, strings.Join(instructions, "\n\t\t"), terminator))
	}
	if len(literals) == 0 {
		return nil, skipped, reasons, fmt.Errorf("no generated blocks use the compact instruction set")
	}
	source := fmt.Sprintf("// Code generated by tools/compact-kernel; DO NOT EDIT.\npackage recovery\n\nvar %s = map[uint64]compactBlock{\n\t%s\n}\n", name, strings.Join(literals, "\n\t"))
	formatted, err := format.Source([]byte(source))
	return formatted, skipped, reasons, err
}

func generateCompactSource(name string, annotated []annotatedInstruction) ([]byte, error) {
	if len(annotated) == 0 {
		return nil, fmt.Errorf("input contains no compact instructions")
	}
	start, expected := annotated[0].pc, annotated[0].pc
	var instructions, segments []string
	for _, entry := range annotated {
		if entry.pc != expected {
			segments = append(segments, fmt.Sprintf("{index: %d, address: %#x},", len(instructions), entry.pc))
		}
		literal, stop, err := compactInstruction(entry.disassembly)
		if err != nil {
			return nil, fmt.Errorf("%#x: %w", entry.pc, err)
		}
		expected = entry.pc + 4
		if !stop {
			instructions = append(instructions, literal)
		}
	}
	if len(instructions) == 0 {
		return nil, fmt.Errorf("input contains no compact instructions")
	}
	segmentSource := ""
	if len(segments) != 0 {
		segmentSource = "\n\tsegments: []compactSegment{\n\t\t" + strings.Join(segments, "\n\t\t") + "\n\t},"
	}
	source := fmt.Sprintf("// Code generated by tools/compact-kernel; DO NOT EDIT.\npackage recovery\n\nvar %s = compactKernel{\n\tstart: %#x,\n\tinstructions: []compactInstruction{\n\t\t%s\n\t},%s\n}\n", name, start, strings.Join(instructions, "\n\t\t"), segmentSource)
	return format.Source([]byte(source))
}

func main() {
	inputPath := flag.String("input", "", "generated Go kernel source")
	outputPath := flag.String("output", "", "generated compact profile Go file")
	name := flag.String("name", "", "Go variable name for the compact kernel")
	cfgPrefix := flag.String("cfg-prefix", "", "generated block function prefix (enables CFG mode)")
	limit := flag.Int("limit", 0, "maximum compatible CFG blocks to emit (zero means all)")
	flag.Parse()
	if *inputPath == "" || *outputPath == "" || !identifier.MatchString(*name) {
		fmt.Fprintln(os.Stderr, "--input, --output, and a valid --name are required")
		os.Exit(2)
	}
	inputBytes, err := os.ReadFile(*inputPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	var formatted []byte
	if *cfgPrefix != "" {
		var blocks []annotatedBlock
		blocks, err = parseAnnotatedBlocks(string(inputBytes), *cfgPrefix)
		if err == nil {
			var skipped int
			var reasons map[string]int
			formatted, skipped, reasons, err = generateCompactCFGSource(*name, blocks, *limit)
			if err == nil {
				fmt.Fprintf(os.Stderr, "generated %d compact blocks; skipped %d unsupported blocks\n", strings.Count(string(formatted), "start:"), skipped)
				type reasonCount struct {
					reason string
					count  int
				}
				ordered := make([]reasonCount, 0, len(reasons))
				for reason, count := range reasons {
					ordered = append(ordered, reasonCount{reason: reason, count: count})
				}
				sort.Slice(ordered, func(i, j int) bool {
					if ordered[i].count != ordered[j].count {
						return ordered[i].count > ordered[j].count
					}
					return ordered[i].reason < ordered[j].reason
				})
				for _, item := range ordered {
					fmt.Fprintf(os.Stderr, "%4d %s\n", item.count, item.reason)
				}
			}
		}
	} else {
		var annotated []annotatedInstruction
		scanner := bufio.NewScanner(strings.NewReader(string(inputBytes)))
		for scanner.Scan() {
			match := commentPattern.FindStringSubmatch(scanner.Text())
			if match == nil {
				continue
			}
			pc, parseErr := strconv.ParseUint(strings.TrimPrefix(match[1], "0x"), 16, 64)
			if parseErr != nil {
				err = parseErr
				break
			}
			annotated = append(annotated, annotatedInstruction{pc: pc, disassembly: match[2]})
		}
		if err == nil {
			err = scanner.Err()
		}
		if err == nil {
			formatted, err = generateCompactSource(*name, annotated)
		}
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if err := os.WriteFile(*outputPath, formatted, 0o644); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
