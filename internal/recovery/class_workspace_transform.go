package recovery

const (
	classWorkspaceModulus        = uint32(0x7fff)
	classWorkspaceSampleCutoff   = uint16(0x7b0f)
	classWorkspacePeerBias       = uint16(0x54ce)
	classWorkspaceOutputEncoding = uint16(0xe748)
)

type classWorkspaceLimbOutputs struct {
	Forward  uint16
	Current  uint16
	Previous uint16
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
