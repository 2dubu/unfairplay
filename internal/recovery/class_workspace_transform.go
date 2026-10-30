package recovery

import "encoding/binary"

const classWorkspaceModulus = uint32(0x7fff)

var classWorkspaceValueColumn = [8]byte{0xfc, 0x24, 0x94, 0x78, 0xca, 0xe1, 0x7a, 0x5b}
var classWorkspaceValueOffset = [8]byte{0xc2, 0x87, 0xa5, 0xfb, 0x4c, 0x3e, 0xea, 0x35}

var legacyClassWorkspaceColumn = [8]byte{0xe1, 0x2e, 0x9b, 0xb9, 0xe8, 0x6d, 0xff, 0xdf}

// classWorkspacePolynomialWord folds a 3x3 product matrix into one semantic
// base-256 polynomial product. Limbs are semantic residues in row-major order.
func classWorkspacePolynomialWord(limbs [9]uint16) uint32 {
	diagonal0 := uint32(limbs[0]) + uint32(limbs[8])
	diagonal1 := uint32(limbs[1]) + uint32(limbs[3])
	diagonal2 := uint32(limbs[2]) + uint32(limbs[4]) + uint32(limbs[6])
	diagonal3 := uint32(limbs[5]) + uint32(limbs[7])
	return diagonal0 + diagonal1<<8 + diagonal2<<16 + diagonal3<<24
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
	return classWorkspaceModularPower32(base, uint32(exponent))
}

func classWorkspaceModularPower32(base uint16, exponent uint32) uint16 {
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

// legacyClassWorkspaceWord is the semantic 3x3 recurrence used by the older
// 630-byte enrollment layout. The native implementation transports three
// differently encoded matrices through a protected limb kernel; this direct
// form retains only the coordinates whose polynomial product is consumed.
func legacyClassWorkspaceWord(residue uint16) uint32 {
	residue %= uint16(classWorkspaceModulus)
	var coordinates, velocities [9]uint16
	for index := uint16(0); index < uint16(len(coordinates)); index++ {
		base := uint16((uint32(residue) + uint32(index) + 1) % classWorkspaceModulus)
		coordinates[index] = classWorkspaceModularPower32(base, uint32(residue)+1)
		velocities[index] = residue
	}
	for column := range 3 {
		for row := range 3 {
			position := row*3 + column
			accumulator := uint32(0x08e513f9)
			for other := range 3 {
				if other == row {
					continue
				}
				var distance uint32
				for component := range 3 {
					difference := int32(coordinates[row*3+component]) - int32(coordinates[other*3+component])
					distance += uint32(difference * difference)
				}
				coefficient := classWorkspaceModularPower32(uint16(distance%classWorkspaceModulus), 0x10001ffd)
				difference := int32(coordinates[other*3+column]) - int32(coordinates[position])
				accumulator += uint32(difference) * uint32(coefficient)
			}
			exponent := uint32(coordinates[position]) + uint32(residue) + uint32(velocities[position])
			direct := classWorkspaceModularPower32(5, exponent)
			impulse := uint16((accumulator + uint32(direct) + uint32(residue) + 0xf71aec07) % classWorkspaceModulus)
			velocities[position] = uint16((uint32(velocities[position]) + uint32(impulse)) % classWorkspaceModulus)
			coordinates[position] = uint16((uint32(coordinates[position]) + uint32(velocities[position])) % classWorkspaceModulus)
		}
	}
	return classWorkspacePolynomialWord(coordinates)
}

// legacyClassWorkspaceValue maps the final two digest-derived polynomial
// words to the eight-byte class value. It is one circulant linear map over
// GF(2^8); the native byte tables and XOR-0x1b transport cancel completely.
func legacyClassWorkspaceValue(digest [20]byte) (result [8]byte) {
	var input [8]byte
	for index := range 2 {
		residue := uint16(uint32(binary.BigEndian.Uint16(digest[16+index*2:])) % classWorkspaceModulus)
		binary.LittleEndian.PutUint32(input[index*4:], legacyClassWorkspaceWord(residue))
	}
	for inputIndex, value := range input {
		for outputIndex := range result {
			columnIndex := (outputIndex - inputIndex + len(result)) % len(result)
			result[outputIndex] ^= classWorkspaceGaloisMultiply(legacyClassWorkspaceColumn[columnIndex], value)
		}
	}
	return result
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

func classWorkspaceOutput(matrix [9]uint16) [8]byte {
	return classWorkspaceClassValue(classWorkspaceInteractionStep(matrix))
}

func classWorkspaceOutputFromSeed(seed uint16) [8]byte {
	return classWorkspaceOutput(classWorkspaceSeedMatrix(seed))
}

func classWorkspaceOutputFromDigest(digest [20]byte) [8]byte {
	return classWorkspaceOutputFromSeed(classWorkspaceSeedFromDigest(digest))
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
