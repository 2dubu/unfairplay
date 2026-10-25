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
	commentPattern = regexp.MustCompile(`^\s*// (0x[0-9a-f]+):\s+(.+)$`)
	movImmediate   = regexp.MustCompile(`^mov ([wx])(\d+), #(0x[0-9a-f]+|[0-9]+)$`)
	movRegister    = regexp.MustCompile(`^mov ([wx])(\d+), [wx](\d+)$`)
	movKeep        = regexp.MustCompile(`^movk w(\d+), #(0x[0-9a-f]+|[0-9]+), lsl #(\d+)$`)
	binary32       = regexp.MustCompile(`^(eor|and|add|sub|orr|mul) w(\d+), w(\d+), w(\d+)(?:, lsl #(\d+))?$`)
	addImmediate32 = regexp.MustCompile(`^add w(\d+), w(\d+), #(0x[0-9a-f]+|[0-9]+)(?:, lsl #(\d+))?$`)
	shiftLeft32    = regexp.MustCompile(`^lsl w(\d+), w(\d+), #(0x[0-9a-f]+|[0-9]+)$`)
	addSP          = regexp.MustCompile(`^add x(\d+), sp, #(0x[0-9a-f]+|[0-9]+), lsl #(\d+)$`)
	loadSP         = regexp.MustCompile(`^ldr ([wx])(\d+), \[sp, #(0x[0-9a-f]+|[0-9]+)\]$`)
	loadBase       = regexp.MustCompile(`^ldr w(\d+), \[x(\d+), #(0x[0-9a-f]+|[0-9]+)\]$`)
	loadSXTW64     = regexp.MustCompile(`^ldr x(\d+), \[x(\d+), w(\d+), sxtw #(\d+)\]$`)
	storeSP        = regexp.MustCompile(`^str w(\d+), \[sp, #(0x[0-9a-f]+|[0-9]+)\]$`)
	storeBase      = regexp.MustCompile(`^str w(\d+), \[x(\d+)\]$`)
	extract32      = regexp.MustCompile(`^extr w(\d+), w(\d+), w(\d+), #(0x[0-9a-f]+|[0-9]+)$`)
	not32          = regexp.MustCompile(`^mvn w(\d+), w(\d+)$`)
	shiftRight32   = regexp.MustCompile(`^lsr w(\d+), w(\d+), #(0x[0-9a-f]+|[0-9]+)$`)
	branchAuth     = regexp.MustCompile(`^braa x\d+, x\d+$`)
	identifier     = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)
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

func compactInstruction(disassembly string) (string, bool, error) {
	if match := movImmediate.FindStringSubmatch(disassembly); match != nil {
		value, err := number(match[3])
		if err != nil || value > 0xffffffff {
			return "", false, fmt.Errorf("invalid MOV immediate %q", match[3])
		}
		op := "compactMovImm32"
		if match[1] == "x" {
			if value > 0xffff {
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
		value, valueErr := number(match[2])
		shift, shiftErr := number(match[3])
		if valueErr != nil || shiftErr != nil || value > 0xffff || shift > 63 {
			return "", false, fmt.Errorf("invalid MOVK %q", disassembly)
		}
		return fields("compactMovK32", "destination", match[1], "shift", fmt.Sprint(shift), "immediate", fmt.Sprintf("%#x", value)), false, nil
	}
	if match := binary32.FindStringSubmatch(disassembly); match != nil {
		opcodes := map[string]string{
			"eor": "compactXor32", "and": "compactAnd32", "add": "compactAdd32",
			"sub": "compactSub32", "orr": "compactOr32", "mul": "compactMul32",
		}
		op := opcodes[match[1]]
		extra := []string{}
		if match[5] != "" {
			if match[1] != "and" && match[1] != "eor" {
				return "", false, fmt.Errorf("unsupported shifted binary operation %q", disassembly)
			}
			if match[1] == "and" {
				op = "compactAndLSL32"
			} else {
				op = "compactXorLSL32"
			}
			extra = append(extra, "shift", match[5])
		}
		return fields(op, append([]string{"destination", match[2], "left", match[3], "right", match[4]}, extra...)...), false, nil
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
	if match := shiftLeft32.FindStringSubmatch(disassembly); match != nil {
		shift, err := number(match[3])
		if err != nil || shift > 31 {
			return "", false, fmt.Errorf("invalid LSL %q", disassembly)
		}
		return fields("compactLSLImm32", "destination", match[1], "left", match[2], "shift", fmt.Sprint(shift)), false, nil
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
	if match := loadSXTW64.FindStringSubmatch(disassembly); match != nil {
		return fields("compactLoadSXTW64", "destination", match[1], "left", match[2], "right", match[3], "shift", match[4]), false, nil
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
	if match := shiftRight32.FindStringSubmatch(disassembly); match != nil {
		shift, err := number(match[3])
		if err != nil || shift > 31 {
			return "", false, fmt.Errorf("invalid LSR %q", disassembly)
		}
		return fields("compactLSRImm32", "destination", match[1], "left", match[2], "shift", fmt.Sprint(shift)), false, nil
	}
	if branchAuth.MatchString(disassembly) {
		return "", true, nil
	}
	return "", false, fmt.Errorf("unsupported instruction %q", disassembly)
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
	start := uint64(0)
	expected := uint64(0)
	var instructions []string
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
		if start == 0 {
			start, expected = pc, pc
		}
		if pc != expected {
			err = fmt.Errorf("non-contiguous instruction at %#x, want %#x", pc, expected)
			break
		}
		literal, stop, compactErr := compactInstruction(match[2])
		if compactErr != nil {
			err = fmt.Errorf("%#x: %w", pc, compactErr)
			break
		}
		expected += 4
		if stop {
			break
		}
		instructions = append(instructions, literal)
	}
	if err == nil {
		err = scanner.Err()
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if start == 0 || len(instructions) == 0 {
		fmt.Fprintln(os.Stderr, "input contains no compact instructions")
		os.Exit(1)
	}
	source := fmt.Sprintf("// Code generated by tools/compact-kernel; DO NOT EDIT.\npackage recovery\n\nvar %s = compactKernel{\n\tstart: %#x,\n\tinstructions: []compactInstruction{\n\t\t%s\n\t},\n}\n", *name, start, strings.Join(instructions, "\n\t\t"))
	formatted, err := format.Source([]byte(source))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if err := os.WriteFile(*outputPath, formatted, 0o644); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
