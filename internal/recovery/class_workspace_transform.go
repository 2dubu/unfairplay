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
