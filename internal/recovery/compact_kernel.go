package recovery

import "fmt"

type compactOpcode uint8

const (
	compactAnd32 compactOpcode = iota
	compactAdd32
	compactLoadSXTW32
	compactXor64
	compactAddLSR64
	compactUMADDL
	compactLSRV64
	compactAddUXTW64
	compactAndLSL64
	compactSub64
	compactAdd64
	compactXor32
	compactStoreSXTW32
	compactAddImm32
	compactCmp32
	compactCCmp32LE
	compactCSet32LT
	compactMul32
	compactLoadSXTW64
	compactMov64
	compactMovImm64
	compactStoreSP32
	compactMovImm32
	compactMovK32
	compactUDiv32
	compactMSub32
	compactLoadSP32
	compactUMull32
	compactLoadUXTW8
	compactLSRImm64
	compactMov32
	compactOr32
	compactSub32
	compactAndLSL32
	compactAddSPShiftImm64
	compactLoadBaseImm32
	compactExtract32
	compactNot32
	compactLSRImm32
	compactLoadSP64
	compactMovImm32High
	compactAddImmShift32
	compactLSLImm32
	compactXorLSL32
	compactStoreBase32
	compactLoadUXTW32
	compactMovNegOne32
	compactMovK64
	compactAddLSL32
	compactXorLSR32
	compactLSRReg32
	compactStoreUXTW8
	compactUMulH64
	compactMAdd32
	compactMAdd64
	compactAdds64
	compactADC32
	compactCMN64
	compactCmpImm32
	compactCSet32HI
	compactStoreBaseImm32
	compactMul64
	compactSubImmShift32
	compactXorImm32
	compactSubLSL32
	compactMovImm32Wide
	compactAnd64
	compactLSLImm64
	compactMovImm64Shift
	compactXorLSL64
	compactCSet32LO
	compactOrLSL32
	compactAndImm32
	compactOrNotLSL32
	compactRORImm32
	compactXorROR32
	compactOrROR32
	compactCSELZeroEq32
	compactCSet32EQ
	compactLoadBase64
	compactLoadPostIndex8
	compactLSLReg32
	compactStoreBase64
	compactCMNImm32
	compactOpcodeCount
)

// compactInstruction stores only the operands needed by the checked recovery
// profile. It is intentionally not a general ARM64 instruction format.
type compactInstruction struct {
	op          compactOpcode
	destination uint8
	left, right uint8
	accumulator uint8
	shift       uint8
	immediate   uint16
}

type compactKernel struct {
	start        uint64
	instructions []compactInstruction
	segments     []compactSegment
}

type compactSegment struct {
	index   uint32
	address uint64
}

func runCompactKernel(s *state, kernel compactKernel) error {
	pc := kernel.start
	segment := 0
	for index, instruction := range kernel.instructions {
		if segment < len(kernel.segments) && kernel.segments[segment].index == uint32(index) {
			pc = kernel.segments[segment].address
			segment++
		}
		d, left, right := int(instruction.destination), int(instruction.left), int(instruction.right)
		switch instruction.op {
		case compactAnd32:
			s.x[d] = uint64(uint32(s.x[left]) & uint32(s.x[right]))
		case compactAdd32:
			s.x[d] = uint64(uint32(s.x[left]) + uint32(s.x[right]))
		case compactLoadSXTW32:
			value, err := s.load(s.x[left]+shiftLeft(signExtend32(uint32(s.x[right])), int(instruction.shift), 64), 4)
			if err != nil {
				return fmt.Errorf("compact kernel %#x: %w", pc, err)
			}
			s.x[d] = uint64(uint32(value))
		case compactXor64:
			s.x[d] = s.x[left] ^ s.x[right]
		case compactAddLSR64:
			s.x[d] = s.x[left] + shiftRight(s.x[right], int(instruction.shift), 64)
		case compactUMADDL:
			s.x[d] = uint64(uint32(s.x[left]))*uint64(uint32(s.x[right])) + s.x[instruction.accumulator]
		case compactLSRV64:
			s.x[d] = shiftRight(s.x[left], int(s.x[right]), 64)
		case compactAddUXTW64:
			s.x[d] = s.x[left] + uint64(uint32(s.x[right]))
		case compactAndLSL64:
			s.x[d] = s.x[left] & shiftLeft(s.x[right], int(instruction.shift), 64)
		case compactLSLImm64:
			s.x[d] = shiftLeft(s.x[left], int(instruction.shift), 64)
		case compactMovImm64Shift:
			s.x[d] = uint64(instruction.immediate) << instruction.shift
		case compactXorLSL64:
			s.x[d] = s.x[left] ^ shiftLeft(s.x[right], int(instruction.shift), 64)
		case compactCSet32LO:
			s.x[d] = uint64(uint32(boolValue(s.condition("lo"))))
		case compactOrLSL32:
			s.x[d] = uint64(uint32(s.x[left]) | uint32(shiftLeft(uint64(uint32(s.x[right])), int(instruction.shift), 32)))
		case compactAndImm32:
			immediate := uint32(instruction.immediate) | uint32(instruction.right)<<16 | uint32(instruction.accumulator)<<24
			s.x[d] = uint64(uint32(s.x[left]) & immediate)
		case compactOrNotLSL32:
			s.x[d] = uint64(uint32(s.x[left]) | ^uint32(shiftLeft(uint64(uint32(s.x[right])), int(instruction.shift), 32)))
		case compactRORImm32:
			s.x[d] = uint64(uint32(rotateRight(uint64(uint32(s.x[left])), int(instruction.shift), 32)))
		case compactXorROR32:
			s.x[d] = uint64(uint32(s.x[left]) ^ uint32(rotateRight(uint64(uint32(s.x[right])), int(instruction.shift), 32)))
		case compactOrROR32:
			s.x[d] = uint64(uint32(s.x[left]) | uint32(rotateRight(uint64(uint32(s.x[right])), int(instruction.shift), 32)))
		case compactCSELZeroEq32:
			if s.condition("eq") {
				s.x[d] = 0
			} else {
				s.x[d] = uint64(uint32(s.x[right]))
			}
		case compactCSet32EQ:
			s.x[d] = uint64(uint32(boolValue(s.condition("eq"))))
		case compactLoadBase64:
			value, err := s.load(s.x[left], 8)
			if err != nil {
				return fmt.Errorf("compact kernel %#x: %w", pc, err)
			}
			s.x[d] = value
		case compactLoadPostIndex8:
			value, err := s.load(s.x[left], 1)
			if err != nil {
				return fmt.Errorf("compact kernel %#x: %w", pc, err)
			}
			s.x[d] = uint64(uint32(value))
			s.x[left] += uint64(instruction.immediate)
		case compactLSLReg32:
			s.x[d] = uint64(uint32(shiftLeft(uint64(uint32(s.x[left])), int(s.x[right]), 32)))
		case compactStoreBase64:
			if err := s.storeValue(s.x[left], s.x[d], 8); err != nil {
				return fmt.Errorf("compact kernel %#x: %w", pc, err)
			}
		case compactCMNImm32:
			_ = s.addWithCarry(uint64(uint32(s.x[left])), uint64(uint32(instruction.immediate)), false, 32)
		case compactSub64:
			s.x[d] = s.x[left] - s.x[right]
		case compactAdd64:
			s.x[d] = s.x[left] + s.x[right]
		case compactXor32:
			s.x[d] = uint64(uint32(s.x[left]) ^ uint32(s.x[right]))
		case compactStoreSXTW32:
			address := s.x[left] + shiftLeft(signExtend32(uint32(s.x[right])), int(instruction.shift), 64)
			if err := s.storeValue(address, uint64(uint32(s.x[d])), 4); err != nil {
				return fmt.Errorf("compact kernel %#x: %w", pc, err)
			}
		case compactAddImm32:
			s.x[d] = uint64(uint32(s.x[left]) + uint32(instruction.immediate))
		case compactCmp32:
			s.compare32(uint32(s.x[left]), uint32(s.x[right]))
		case compactCCmp32LE:
			if s.condition("le") {
				s.compare32(uint32(s.x[left]), uint32(s.x[right]))
			} else {
				s.setNZCV(int(instruction.immediate))
			}
		case compactCSet32LT:
			s.x[d] = uint64(uint32(boolValue(s.condition("lt"))))
		case compactMul32:
			s.x[d] = uint64(uint32(s.x[left]) * uint32(s.x[right]))
		case compactLoadSXTW64:
			value, err := s.load(s.x[left]+shiftLeft(signExtend32(uint32(s.x[right])), int(instruction.shift), 64), 8)
			if err != nil {
				return fmt.Errorf("compact kernel %#x: %w", pc, err)
			}
			s.x[d] = value
		case compactMov64:
			s.x[d] = s.x[left]
		case compactMovImm64:
			s.x[d] = uint64(instruction.immediate)
		case compactStoreSP32:
			if err := s.storeValue(addOffset(s.sp, int64(instruction.immediate)), uint64(uint32(s.x[left])), 4); err != nil {
				return fmt.Errorf("compact kernel %#x: %w", pc, err)
			}
		case compactMovImm32:
			s.x[d] = uint64(uint32(instruction.immediate))
		case compactMovK32:
			s.x[d] = uint64(uint32(insertBits(uint64(uint32(s.x[d])), uint64(instruction.immediate), int(instruction.shift), 16, 32)))
		case compactUDiv32:
			s.x[d] = uint64(uint32(unsignedDivide(uint64(uint32(s.x[left])), uint64(uint32(s.x[right])), 32)))
		case compactMSub32:
			s.x[d] = uint64(uint32(s.x[instruction.accumulator]) - uint32(s.x[left])*uint32(s.x[right]))
		case compactLoadSP32:
			value, err := s.load(addOffset(s.sp, int64(instruction.immediate)), 4)
			if err != nil {
				return fmt.Errorf("compact kernel %#x: %w", pc, err)
			}
			s.x[d] = uint64(uint32(value))
		case compactUMull32:
			s.x[d] = uint64(uint32(s.x[left])) * uint64(uint32(s.x[right]))
		case compactLoadUXTW8:
			value, err := s.load(s.x[left]+uint64(uint32(s.x[right])), 1)
			if err != nil {
				return fmt.Errorf("compact kernel %#x: %w", pc, err)
			}
			s.x[d] = uint64(uint32(value))
		case compactLSRImm64:
			s.x[d] = shiftRight(s.x[left], int(instruction.shift), 64)
		case compactMov32:
			s.x[d] = uint64(uint32(s.x[left]))
		case compactOr32:
			s.x[d] = uint64(uint32(s.x[left]) | uint32(s.x[right]))
		case compactSub32:
			s.x[d] = uint64(uint32(s.x[left]) - uint32(s.x[right]))
		case compactAndLSL32:
			s.x[d] = uint64(uint32(s.x[left]) & uint32(shiftLeft(uint64(uint32(s.x[right])), int(instruction.shift), 32)))
		case compactAddSPShiftImm64:
			s.x[d] = s.sp + shiftLeft(uint64(instruction.immediate), int(instruction.shift), 64)
		case compactLoadBaseImm32:
			value, err := s.load(addOffset(s.x[left], int64(instruction.immediate)), 4)
			if err != nil {
				return fmt.Errorf("compact kernel %#x: %w", pc, err)
			}
			s.x[d] = uint64(uint32(value))
		case compactExtract32:
			s.x[d] = uint64(uint32(extractPair(uint64(uint32(s.x[left])), uint64(uint32(s.x[right])), int(instruction.shift), 32)))
		case compactNot32:
			s.x[d] = uint64(^uint32(s.x[left]))
		case compactLSRImm32:
			s.x[d] = uint64(uint32(shiftRight(uint64(uint32(s.x[left])), int(instruction.shift), 32)))
		case compactLoadSP64:
			value, err := s.load(addOffset(s.sp, int64(instruction.immediate)), 8)
			if err != nil {
				return fmt.Errorf("compact kernel %#x: %w", pc, err)
			}
			s.x[d] = value
		case compactMovImm32High:
			s.x[d] = uint64(uint32(instruction.immediate) << 16)
		case compactAddImmShift32:
			s.x[d] = uint64(uint32(s.x[left]) + uint32(instruction.immediate)<<instruction.shift)
		case compactLSLImm32:
			s.x[d] = uint64(uint32(shiftLeft(uint64(uint32(s.x[left])), int(instruction.shift), 32)))
		case compactXorLSL32:
			s.x[d] = uint64(uint32(s.x[left]) ^ uint32(shiftLeft(uint64(uint32(s.x[right])), int(instruction.shift), 32)))
		case compactStoreBase32:
			if err := s.storeValue(s.x[left], uint64(uint32(s.x[d])), 4); err != nil {
				return fmt.Errorf("compact kernel %#x: %w", pc, err)
			}
		case compactLoadUXTW32:
			value, err := s.load(s.x[left]+shiftLeft(uint64(uint32(s.x[right])), int(instruction.shift), 64), 4)
			if err != nil {
				return fmt.Errorf("compact kernel %#x: %w", pc, err)
			}
			s.x[d] = uint64(uint32(value))
		case compactMovNegOne32:
			s.x[d] = uint64(^uint32(0))
		case compactMovK64:
			s.x[d] = insertBits(s.x[d], uint64(instruction.immediate), int(instruction.shift), 16, 64)
		case compactAddLSL32:
			s.x[d] = uint64(uint32(s.x[left]) + uint32(shiftLeft(uint64(uint32(s.x[right])), int(instruction.shift), 32)))
		case compactXorLSR32:
			s.x[d] = uint64(uint32(s.x[left]) ^ uint32(shiftRight(uint64(uint32(s.x[right])), int(instruction.shift), 32)))
		case compactLSRReg32:
			s.x[d] = uint64(uint32(shiftRight(uint64(uint32(s.x[left])), int(s.x[right]), 32)))
		case compactStoreUXTW8:
			address := s.x[left] + uint64(uint32(s.x[right]))
			if err := s.storeValue(address, uint64(byte(s.x[d])), 1); err != nil {
				return fmt.Errorf("compact kernel %#x: %w", pc, err)
			}
		case compactUMulH64:
			s.x[d] = multiplyHigh64(s.x[left], s.x[right])
		case compactMAdd32:
			s.x[d] = uint64(uint32(s.x[left])*uint32(s.x[right]) + uint32(s.x[instruction.accumulator]))
		case compactMAdd64:
			s.x[d] = s.x[left]*s.x[right] + s.x[instruction.accumulator]
		case compactAdds64:
			s.x[d] = s.addWithCarry(s.x[left], s.x[right], false, 64)
		case compactADC32:
			leftValue, rightValue := uint32(0), uint32(0)
			if instruction.left < 31 {
				leftValue = uint32(s.x[left])
			}
			if instruction.right < 31 {
				rightValue = uint32(s.x[right])
			}
			s.x[d] = uint64(leftValue + rightValue + uint32(boolValue(s.c)))
		case compactCMN64:
			_ = s.addWithCarry(s.x[left], s.x[right], false, 64)
		case compactCmpImm32:
			s.compare32(uint32(s.x[left]), uint32(instruction.immediate))
		case compactCSet32HI:
			s.x[d] = uint64(uint32(boolValue(s.condition("hi"))))
		case compactStoreBaseImm32:
			if err := s.storeValue(addOffset(s.x[left], int64(instruction.immediate)), uint64(uint32(s.x[d])), 4); err != nil {
				return fmt.Errorf("compact kernel %#x: %w", pc, err)
			}
		case compactMul64:
			s.x[d] = s.x[left] * s.x[right]
		case compactSubImmShift32:
			s.x[d] = uint64(uint32(s.x[left]) - uint32(instruction.immediate)<<instruction.shift)
		case compactXorImm32:
			immediate := uint32(instruction.immediate) | uint32(instruction.right)<<16 | uint32(instruction.accumulator)<<24
			s.x[d] = uint64(uint32(s.x[left]) ^ immediate)
		case compactSubLSL32:
			s.x[d] = uint64(uint32(s.x[left]) - uint32(shiftLeft(uint64(uint32(s.x[right])), int(instruction.shift), 32)))
		case compactMovImm32Wide:
			s.x[d] = uint64(uint32(instruction.immediate) | uint32(instruction.right)<<16 | uint32(instruction.accumulator)<<24)
		case compactAnd64:
			s.x[d] = s.x[left] & s.x[right]
		default:
			return fmt.Errorf("compact kernel %#x has unknown opcode %d", pc, instruction.op)
		}
		pc += 4
	}
	return nil
}

var compact26A428DigestPrologue = compactKernel{
	start: 0x1007A8214,
	instructions: []compactInstruction{
		{op: compactStoreSP32, left: 0, immediate: 3464},
		{op: compactStoreSP32, left: 16, immediate: 3520},
		{op: compactMovImm32, destination: 14, immediate: 0x92B8},
		{op: compactMovK32, destination: 14, shift: 16, immediate: 0xB039},
		{op: compactAdd32, destination: 14, left: 5, right: 14},
		{op: compactUDiv32, destination: 16, left: 14, right: 11},
		{op: compactMSub32, destination: 14, left: 16, right: 11, accumulator: 14},
		{op: compactLoadSP32, destination: 15, immediate: 3424},
		{op: compactAdd32, destination: 17, left: 15, right: 8},
		{op: compactUMull32, destination: 0, left: 17, right: 9},
		{op: compactLoadUXTW8, destination: 16, left: 30, right: 14},
		{op: compactLSRImm64, destination: 0, left: 0, shift: 34},
		{op: compactMSub32, destination: 17, left: 0, right: 12, accumulator: 17},
	},
}

var compact26A428BootstrapKernel = compactKernel{
	start: 0x10057AE74,
	instructions: []compactInstruction{
		{op: compactAnd32, destination: 11, left: 14, right: 10},
		{op: compactAdd32, destination: 17, left: 20, right: 13},
		{op: compactLoadSXTW32, destination: 2, left: 4, right: 17, shift: 2},
		{op: compactXor64, destination: 2, left: 2, right: 24},
		{op: compactAddLSR64, destination: 2, left: 2, right: 28, shift: 32},
		{op: compactLoadSXTW32, destination: 3, left: 5, right: 17, shift: 2},
		{op: compactXor64, destination: 3, left: 3, right: 24},
		{op: compactUMADDL, destination: 28, left: 3, right: 12, accumulator: 2},
		{op: compactXor64, destination: 2, left: 11, right: 23},
		{op: compactLSRV64, destination: 2, left: 30, right: 2},
		{op: compactXor64, destination: 2, left: 2, right: 6},
		{op: compactLoadSXTW32, destination: 17, left: 22, right: 17, shift: 2},
		{op: compactXor64, destination: 17, left: 17, right: 24},
		{op: compactUMADDL, destination: 17, left: 17, right: 19, accumulator: 2},
		{op: compactAddUXTW64, destination: 17, left: 17, right: 28},
		{op: compactAndLSL64, destination: 2, left: 7, right: 17, shift: 1},
		{op: compactSub64, destination: 17, left: 17, right: 2},
		{op: compactAdd64, destination: 30, left: 17, right: 26},
		{op: compactXor32, destination: 17, left: 30, right: 27},
		{op: compactAdd32, destination: 2, left: 20, right: 21},
		{op: compactStoreSXTW32, destination: 17, left: 4, right: 2, shift: 2},
		{op: compactAddImm32, destination: 2, left: 20, immediate: 1},
		{op: compactAdd32, destination: 17, left: 20, right: 16},
		{op: compactCmp32, left: 17, right: 15},
		{op: compactCCmp32LE, left: 17, right: 1, immediate: 8},
		{op: compactCSet32LT, destination: 17},
		{op: compactMul32, destination: 17, left: 17, right: 0},
		{op: compactXor32, destination: 17, left: 17, right: 9},
		{op: compactLoadSXTW64, destination: 3, left: 25, right: 17, shift: 3},
		{op: compactMov64, destination: 20, left: 2},
		{op: compactMovImm64, destination: 17, immediate: 0xCD95},
	},
}
