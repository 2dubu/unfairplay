package recovery

import "fmt"

const sharedEnrollmentPolynomialReturn = uint64(0x100265000)

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
