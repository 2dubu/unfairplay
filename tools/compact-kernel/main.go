package main

import (
	"bufio"
	"flag"
	"fmt"
	"go/format"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var (
	commentPattern       = regexp.MustCompile(`^\s*// (0x[0-9a-f]+):\s+(.+)$`)
	movImmediate         = regexp.MustCompile(`^mov ([wx])(\d+), #(0x[0-9a-f]+|[0-9]+)$`)
	movNegOne32          = regexp.MustCompile(`^mov w(\d+), #-1$`)
	movNegative32        = regexp.MustCompile(`^mov w(\d+), #-(0x[0-9a-f]+|[0-9]+)$`)
	movRegister          = regexp.MustCompile(`^mov ([wx])(\d+), [wx](\d+)$`)
	movKeep              = regexp.MustCompile(`^movk ([wx])(\d+), #(0x[0-9a-f]+|[0-9]+), lsl #(\d+)$`)
	binary32             = regexp.MustCompile(`^(eor|and|add|sub|orr|mul) w(\d+), w(\d+), w(\d+)(?:, (lsl|lsr|ror) #(\d+))?$`)
	binary64             = regexp.MustCompile(`^(eor|add|mul|and) x(\d+), x(\d+), x(\d+)$`)
	binaryShift64        = regexp.MustCompile(`^(and|eor) x(\d+), x(\d+), x(\d+), lsl #(\d+)$`)
	addImmediate32       = regexp.MustCompile(`^add w(\d+), w(\d+), #(0x[0-9a-f]+|[0-9]+)(?:, lsl #(\d+))?$`)
	subImmediate32       = regexp.MustCompile(`^sub w(\d+), w(\d+), #(0x[0-9a-f]+|[0-9]+)(?:, lsl #(\d+))?$`)
	xorImmediate32       = regexp.MustCompile(`^eor w(\d+), w(\d+), #(0x[0-9a-f]+|[0-9]+)$`)
	andImmediate32       = regexp.MustCompile(`^and w(\d+), w(\d+), #(0x[0-9a-f]+|[0-9]+)$`)
	shiftLeft32          = regexp.MustCompile(`^lsl w(\d+), w(\d+), #(0x[0-9a-f]+|[0-9]+)$`)
	shiftLeft64          = regexp.MustCompile(`^lsl x(\d+), x(\d+), #(0x[0-9a-f]+|[0-9]+)$`)
	shiftLeftReg32       = regexp.MustCompile(`^lsl w(\d+), w(\d+), w(\d+)$`)
	addSP                = regexp.MustCompile(`^add x(\d+), sp, #(0x[0-9a-f]+|[0-9]+), lsl #(\d+)$`)
	loadSP               = regexp.MustCompile(`^ldr ([wx])(\d+), \[sp, #(0x[0-9a-f]+|[0-9]+)\]$`)
	loadBase             = regexp.MustCompile(`^ldr w(\d+), \[x(\d+), #(0x[0-9a-f]+|[0-9]+)\]$`)
	loadBase64           = regexp.MustCompile(`^ldr x(\d+), \[x(\d+)\]$`)
	loadSXTW64           = regexp.MustCompile(`^ldr x(\d+), \[x(\d+), w(\d+), sxtw #(\d+)\]$`)
	loadUXTW32           = regexp.MustCompile(`^ldr w(\d+), \[x(\d+), w(\d+), uxtw #(\d+)\]$`)
	loadUXTW8            = regexp.MustCompile(`^ldrb w(\d+), \[x(\d+), w(\d+), uxtw\]$`)
	loadPostIndex8       = regexp.MustCompile(`^ldrb w(\d+), \[x(\d+)\], #(0x[0-9a-f]+|[0-9]+)$`)
	storeSP              = regexp.MustCompile(`^str w(\d+), \[sp, #(0x[0-9a-f]+|[0-9]+)\]$`)
	storeBase            = regexp.MustCompile(`^str w(\d+), \[x(\d+)\]$`)
	storeBase64          = regexp.MustCompile(`^str x(\d+), \[x(\d+)\]$`)
	storeBaseImm         = regexp.MustCompile(`^str w(\d+), \[x(\d+), #(0x[0-9a-f]+|[0-9]+)\]$`)
	storeUXTW8           = regexp.MustCompile(`^strb w(\d+), \[x(\d+), w(\d+), uxtw\]$`)
	extract32            = regexp.MustCompile(`^extr w(\d+), w(\d+), w(\d+), #(0x[0-9a-f]+|[0-9]+)$`)
	not32                = regexp.MustCompile(`^mvn w(\d+), w(\d+)$`)
	orNotLSL32           = regexp.MustCompile(`^orn w(\d+), w(\d+), w(\d+), lsl #(\d+)$`)
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
	conditionSetHI       = regexp.MustCompile(`^cset w(\d+), hi$`)
	conditionSetLO       = regexp.MustCompile(`^cset w(\d+), lo$`)
	conditionSetEQ       = regexp.MustCompile(`^cset w(\d+), eq$`)
	compareRegister32    = regexp.MustCompile(`^cmp w(\d+), w(\d+)$`)
	selectZeroEq32       = regexp.MustCompile(`^csel w(\d+), wzr, w(\d+), eq$`)
	branchAuth           = regexp.MustCompile(`^braa x\d+, x\d+$`)
	identifier           = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)
)

func number(value string) (uint64, error) {
	base := 10
	if strings.HasPrefix(value, "0x") {
		base, value = 16, strings.TrimPrefix(value, "0x")
	}
	return strconv.ParseUint(value, base, 64)
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

func compactInstruction(disassembly string) (string, bool, error) {
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
	if match := movImmediate.FindStringSubmatch(disassembly); match != nil {
		value, err := number(match[3])
		if err != nil || match[1] == "w" && value > 0xffffffff {
			return "", false, fmt.Errorf("invalid MOV immediate %q", match[3])
		}
		op := "compactMovImm32"
		if match[1] == "x" {
			if value > 0xffff {
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
				return "", false, fmt.Errorf("unsupported 32-bit MOV immediate %q", match[3])
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
		shift, shiftErr := number(match[4])
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
		opcodes := map[string]string{"eor": "compactXor64", "add": "compactAdd64", "mul": "compactMul64", "and": "compactAnd64"}
		op := opcodes[match[1]]
		return fields(op, "destination", match[2], "left", match[3], "right", match[4]), false, nil
	}
	if match := binaryShift64.FindStringSubmatch(disassembly); match != nil {
		op := "compactAndLSL64"
		if match[1] == "eor" {
			op = "compactXorLSL64"
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
		shift, shiftErr := number(match[3])
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
	if match := loadPostIndex8.FindStringSubmatch(disassembly); match != nil {
		immediate, err := number(match[3])
		if err != nil || immediate > 0xffff {
			return "", false, fmt.Errorf("invalid post-index load %q", disassembly)
		}
		return fields("compactLoadPostIndex8", "destination", match[1], "left", match[2], "immediate", fmt.Sprintf("%#x", immediate)), false, nil
	}
	if match := storeSP.FindStringSubmatch(disassembly); match != nil {
		immediate, err := number(match[2])
		if err != nil || immediate > 0xffff {
			return "", false, fmt.Errorf("invalid SP store %q", disassembly)
		}
		return fields("compactStoreSP32", "left", match[1], "immediate", fmt.Sprintf("%#x", immediate)), false, nil
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
	if match := storeUXTW8.FindStringSubmatch(disassembly); match != nil {
		return fields("compactStoreUXTW8", "destination", match[1], "left", match[2], "right", match[3]), false, nil
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
		return fields("compactOrNotLSL32", "destination", match[1], "left", match[2], "right", match[3], "shift", match[4]), false, nil
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
	if match := selectZeroEq32.FindStringSubmatch(disassembly); match != nil {
		return fields("compactCSELZeroEq32", "destination", match[1], "right", match[2]), false, nil
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
	flag.Parse()
	if *inputPath == "" || *outputPath == "" || !identifier.MatchString(*name) {
		fmt.Fprintln(os.Stderr, "--input, --output, and a valid --name are required")
		os.Exit(2)
	}
	input, err := os.Open(*inputPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer input.Close()
	var annotated []annotatedInstruction
	scanner := bufio.NewScanner(input)
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
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	formatted, err := generateCompactSource(*name, annotated)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if err := os.WriteFile(*outputPath, formatted, 0o644); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
