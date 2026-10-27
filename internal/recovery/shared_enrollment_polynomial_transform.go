package recovery

import "fmt"

const sharedEnrollmentPolynomialReturn = uint64(0x100265000)

const sharedEnrollmentPackedStart = uint64(0x100264a8c)

func sharedEnrollmentPackedWeight(value byte) int32 {
	return int32(value&0xe4) - int32(value&0x1b)
}

// runSharedEnrollmentPackedPath decodes the three table bytes consumed by the
// fourth enrollment branch and packs the two polynomial operands directly.
func runSharedEnrollmentPackedPath(s *state) (uint64, error) {
	if s == nil || uint32(s.x[2]) != 0x19d6 || uint32(s.x[5]) != 0xba1c1cae ||
		uint32(s.x[6]) != 0xae || uint32(s.x[7]) != 0x52 || uint32(s.x[8]) != 0x1b ||
		uint32(s.x[11])-uint32(s.x[12]) != 0x22f1f10c ||
		uint32(s.x[12])-uint32(s.x[16]) != 0xae {
		return 0, fmt.Errorf("unsupported shared enrollment packed-path ABI")
	}
	current := byte(s.x[20])
	if current&3 != 2 {
		return 0, fmt.Errorf("unsupported shared enrollment packed class %d", current&3)
	}
	work, err := s.load(s.sp+0x58, 8)
	if err != nil {
		return 0, fmt.Errorf("shared enrollment packed table: %w", err)
	}
	index := byte(s.x[0])
	first, err := s.load(work+uint64(index+3), 1)
	if err != nil {
		return 0, fmt.Errorf("shared enrollment packed first byte: %w", err)
	}
	second, err := s.load(work+uint64(index+5), 1)
	if err != nil {
		return 0, fmt.Errorf("shared enrollment packed second byte: %w", err)
	}
	third, err := s.load(work+uint64(index+9), 1)
	if err != nil {
		return 0, fmt.Errorf("shared enrollment packed third byte: %w", err)
	}
	s.x[1] = 0x65c0
	leftBase := uint32(int64(0x6997938b) + int64(sharedEnrollmentPackedWeight(current)+2)*0x100)
	left := leftBase + uint32(sharedEnrollmentPackedWeight(byte(first)))
	right := uint32(int64(0xb42c83c3) + int64(sharedEnrollmentPackedWeight(byte(second)))*0x100 + int64(sharedEnrollmentPackedWeight(byte(third))))
	s.x[17] = uint64(left)
	s.x[19] = uint64(right)
	s.x[28] = work
	return runSharedEnrollmentPolynomialPath(s)
}

// runSharedEnrollmentPolynomialPath replaces the arithmetic tail shared by
// record selection and SINF enrollment. Its two packed inputs feed a uint32
// polynomial, signed reduction modulo 63211, and a final byte encoding.
func runSharedEnrollmentPolynomialPath(s *state) (uint64, error) {
	if s == nil || uint32(s.x[2]) != 0x19d6 {
		return 0, fmt.Errorf("unsupported shared enrollment polynomial ABI")
	}
	right := uint32(s.x[19])
	s.x[20] = uint64(sharedEnrollmentPolynomialOutput(uint32(s.x[17]), right))
	s.x[19] = 0x4b
	s.x[22] = uint64(right + 0x4bd39758)
	next, err := loadBootstrapPointer(s, s.x[21], uint32(s.x[2]))
	if err != nil {
		return 0, fmt.Errorf("shared enrollment polynomial target: %w", err)
	}
	s.x[2] = next
	s.x[17] = 0xcd20
	return sharedEnrollmentPolynomialReturn, nil
}

func sharedEnrollmentPolynomialOutput(left, right uint32) uint32 {
	remainder := uint32(int32(sharedEnrollmentPolynomialRaw(left, right)) % 63211)
	remainder -= 0xa0 & (remainder << 1)
	return (remainder + 0x50) ^ 0x4b
}

func sharedEnrollmentPolynomialRaw(left, right uint32) uint32 {
	rightStep := right + 0x4bd39758
	rightPolynomial := right*0xf707 + 0xf81ee69c
	rightPolynomial *= rightStep
	rightPolynomial += right * 0x41b024cc
	rightBranch := right + 0xc363b8e9
	rightBranch *= rightPolynomial + 0x9264fa20
	rightPolynomial = rightPolynomial*0x886fde6f + rightBranch
	rightBranch = right * 0xaf2d6c25
	rightPolynomial += 0x0bd8ba63
	rightPolynomial *= rightStep
	rightPolynomial += right * 0x72b7797d
	rightPolynomial += 0x76381dd6
	rightPolynomial *= rightStep
	rightPolynomial += right * 0xf2896022

	leftStep := left + 0x96688590
	leftPolynomial := left*0xa14b + 0x6ea0c121
	leftPolynomial *= leftStep
	leftPolynomial += left * 0x5dfdf00f
	leftPolynomial += 0xcc636f1f
	leftPolynomial *= leftStep
	leftPolynomial += left * 0x50436451
	leftPolynomial += 0x6f1418ff
	leftPolynomial *= leftStep
	leftPolynomial += left * 0xf0da6991
	leftPolynomial += 0x1874d931
	leftPolynomial *= leftStep
	leftPolynomial += left * 0xb426dd5f
	leftPolynomial += 0x7049744b

	cross := rightPolynomial*0xb8904dec + rightBranch
	cross = (left+0xddd837a4)*(rightPolynomial+0x98f319b0) + cross
	return (leftPolynomial*rightStep + cross) + 0x0dbc1c05
}
