package recovery

const (
	classWorkspaceModulus        = uint32(0x7fff)
	classWorkspaceSampleCutoff   = uint16(0x7b0f)
	classWorkspacePeerBias       = uint16(0x54ce)
	classWorkspaceOutputEncoding = uint16(0xe748)
	classWorkspaceProductBias    = uint32(0x5a39f893)
)

var classWorkspaceValueColumn = [8]byte{0xfc, 0x24, 0x94, 0x78, 0xca, 0xe1, 0x7a, 0x5b}
var classWorkspaceValueOffset = [8]byte{0xc2, 0x87, 0xa5, 0xfb, 0x4c, 0x3e, 0xea, 0x35}

type classWorkspaceLimbOutputs struct {
	Forward  uint16
	Current  uint16
	Previous uint16
}

// classWorkspacePolynomialWord folds a 3x3 product matrix into one semantic
// base-256 polynomial product. Limbs are semantic residues in row-major order.
func classWorkspacePolynomialWord(limbs [9]uint16) uint32 {
	diagonal0 := uint32(limbs[0]) + uint32(limbs[8])
	diagonal1 := uint32(limbs[1]) + uint32(limbs[3])
	diagonal2 := uint32(limbs[2]) + uint32(limbs[4]) + uint32(limbs[6])
	diagonal3 := uint32(limbs[5]) + uint32(limbs[7])
	return diagonal0 + diagonal1<<8 + diagonal2<<16 + diagonal3<<24
}

// classWorkspaceProductWord retains the exact represented word at the native
// boundary. The following class-value transform adds the inverse bias, so the
// composed semantic path uses classWorkspacePolynomialWord directly.
func classWorkspaceProductWord(limbs [9]uint16) uint32 {
	return classWorkspacePolynomialWord(limbs) + classWorkspaceProductBias
}

func classWorkspaceGaloisMultiply(left, right byte) (result byte) {
	for range 8 {
		if right&1 != 0 {
			result ^= left
		}
		high := left & 0x80
		left <<= 1
		if high != 0 {
			left ^= 0x1b
		}
		right >>= 1
	}
	return result
}

// classWorkspacePowerFactor reduces the squared norm of three semantic
// residues, then applies the fixed exponent used by the class workspace. The
// native 65,536-entry protected decoder is exactly this x^45 operation modulo
// 32767 on its accepted representation.
func classWorkspacePowerFactor(components [3]uint16) uint16 {
	var norm uint64
	for _, component := range components {
		value := uint64(component)
		norm += value * value
	}
	norm %= uint64(classWorkspaceModulus)
	square := norm * norm % uint64(classWorkspaceModulus)
	fourth := square * square % uint64(classWorkspaceModulus)
	eighth := fourth * fourth % uint64(classWorkspaceModulus)
	sixteenth := eighth * eighth % uint64(classWorkspaceModulus)
	thirtySecond := sixteenth * sixteenth % uint64(classWorkspaceModulus)
	return uint16(thirtySecond * eighth % uint64(classWorkspaceModulus) * fourth % uint64(classWorkspaceModulus) * norm % uint64(classWorkspaceModulus))
}

func classWorkspaceModularPower(base, exponent uint16) uint16 {
	value := uint64(base)
	result := uint64(1)
	for exponent != 0 {
		if exponent&1 != 0 {
			result = result * value % uint64(classWorkspaceModulus)
		}
		value = value * value % uint64(classWorkspaceModulus)
		exponent >>= 1
	}
	return uint16(result)
}

// classWorkspaceDrift replaces the native 150-value decoder. Five has order
// 150 in the relevant residue ring, so reducing the exponent is exact.
func classWorkspaceDrift(value uint16) uint16 {
	result := uint64(classWorkspaceModularPower(5, value%150))
	return uint16((0x39a5*result + 0x0c7e) % uint64(classWorkspaceModulus))
}

func classWorkspaceSeedMatrix(seed uint16) (matrix [9]uint16) {
	for index := range matrix {
		base := uint16((uint32(seed) + uint32(index)) % classWorkspaceModulus)
		matrix[index] = classWorkspaceModularPower(base, 1000)
	}
	return matrix
}

// classWorkspaceSeedFromDigest derives the class-workspace seed from the
// final two bytes of the semantic 20-byte account digest. The native path
// transports those bytes under XOR 0x1b, builds several temporary arithmetic
// representations, and then removes every representation again.
//
// The eight suffixes whose low fourteen bits are zero or one retain one
// observable bit from the native path's bounded carry recurrence. Keeping the
// exception makes this typed form exact over the complete 16-bit domain.
func classWorkspaceSeedFromDigest(digest [20]byte) uint16 {
	suffix := uint16(digest[18])<<8 | uint16(digest[19])
	seed := uint16((uint32(suffix)%classWorkspaceModulus + 1) % classWorkspaceModulus)
	if suffix&0x3ffe == 0 {
		seed ^= 0x4000
	}
	return seed
}

// classWorkspaceInteractionStep applies one in-place 3x3 interaction pass.
// Cells are visited column-major, matching the mathematical dependency order:
// every cell observes the already-updated cells that precede it. Each of the
// other two rows contributes its coordinate difference multiplied by the
// 45th power of the rows' squared distance. Product accumulation intentionally
// wraps at 32 bits, matching the algorithm before reduction modulo 32767.
func classWorkspaceInteractionStep(matrix [9]uint16) [9]uint16 {
	for column := range 3 {
		for row := range 3 {
			position := row*3 + column
			force := uint32(0)
			for other := range 3 {
				if other == row {
					continue
				}
				var difference [3]uint16
				for coordinate := range 3 {
					left := uint32(matrix[row*3+coordinate])
					right := uint32(matrix[other*3+coordinate])
					difference[coordinate] = uint16((left + classWorkspaceModulus - right) % classWorkspaceModulus)
				}
				weight := uint32(classWorkspacePowerFactor(difference))
				delta := uint32(matrix[other*3+column]) - uint32(matrix[position])
				force += delta * weight
			}
			matrix[position] = uint16((uint64(matrix[position]) + uint64(classWorkspaceDrift(matrix[position])) + uint64(force%classWorkspaceModulus)) % uint64(classWorkspaceModulus))
		}
	}
	return matrix
}

func classWorkspaceEnrollmentValue(matrix [9]uint16) [8]byte {
	return classWorkspaceClassValue(classWorkspaceInteractionStep(matrix))
}

func classWorkspaceEnrollmentValueFromSeed(seed uint16) [8]byte {
	return classWorkspaceEnrollmentValue(classWorkspaceSeedMatrix(seed))
}

func classWorkspaceEnrollmentValueFromDigest(digest [20]byte) [8]byte {
	return classWorkspaceEnrollmentValueFromSeed(classWorkspaceSeedFromDigest(digest))
}

// classWorkspaceClassValue maps the semantic polynomial product to the
// enrollment-global eight-byte class value. It is one affine circulant matrix
// over GF(2^8), using the AES reduction polynomial, not a white-box cipher.
func classWorkspaceClassValue(limbs [9]uint16) (result [8]byte) {
	word := classWorkspacePolynomialWord(limbs)
	result = classWorkspaceValueOffset
	for inputIndex := range 4 {
		value := byte(word >> (8 * inputIndex))
		for outputIndex := range result {
			columnIndex := (outputIndex - inputIndex + len(result)) % len(result)
			result[outputIndex] ^= classWorkspaceGaloisMultiply(classWorkspaceValueColumn[columnIndex], value)
		}
	}
	return result
}

func classWorkspaceResidue(sample, peer uint16, left, right uint32) uint16 {
	sampleWord := uint32(sample)
	if sample < classWorkspaceSampleCutoff {
		sampleWord += 1 << 16
	}
	peerWord := uint32(uint16(peer - classWorkspacePeerBias))
	sum := left + sampleWord + peerWord + right
	return uint16(sum % classWorkspaceModulus)
}

// classWorkspaceResidueStep is one typed limb recurrence in the class-value
// producer. sample and peer are stored 16-bit representations; left and right
// are the two recurrence terms carried by the enclosing schedule.
func classWorkspaceResidueStep(sample, peer uint16, left, right uint32) uint16 {
	return classWorkspaceResidue(sample, peer, left, right) ^ classWorkspaceOutputEncoding
}

// classWorkspaceLimbStep is the complete three-residue kernel used by one
// limb of the class-value workspace. The returned values are the three stored
// 16-bit representations; the surrounding 3x3 schedule decides their slots.
func classWorkspaceLimbStep(sample, peer, current, previous uint16, left, right uint32) classWorkspaceLimbOutputs {
	forwardResidue := classWorkspaceResidue(sample, peer, left, right)

	currentWord := uint32(uint16(current - classWorkspacePeerBias))
	currentResidue := uint16((uint32(forwardResidue) + currentWord) % classWorkspaceModulus)

	previousWord := uint32(previous)
	if previous < classWorkspaceSampleCutoff {
		previousWord += 1 << 16
	}
	previousResidue := uint16((previousWord + uint32(currentResidue) + right) % classWorkspaceModulus)

	return classWorkspaceLimbOutputs{
		Forward:  forwardResidue ^ classWorkspaceOutputEncoding,
		Current:  uint16(uint32(currentResidue) + uint32(classWorkspacePeerBias)),
		Previous: uint16(uint32(previousResidue) + uint32(classWorkspaceSampleCutoff)),
	}
}
