package recovery

const (
	classWorkspaceModulus        = uint32(0x7fff)
	classWorkspaceSampleCutoff   = uint16(0x7b0f)
	classWorkspacePeerBias       = uint16(0x54ce)
	classWorkspaceOutputEncoding = uint16(0xe748)
)

// classWorkspaceResidueStep is one typed limb recurrence in the class-value
// producer. sample and peer are stored 16-bit representations; left and right
// are the two recurrence terms carried by the enclosing schedule.
func classWorkspaceResidueStep(sample, peer uint16, left, right uint32) uint16 {
	sampleWord := uint32(sample)
	if sample < classWorkspaceSampleCutoff {
		sampleWord += 1 << 16
	}
	peerWord := uint32(uint16(peer - classWorkspacePeerBias))
	sum := left + sampleWord + peerWord + right
	residue := sum % classWorkspaceModulus
	return uint16(residue) ^ classWorkspaceOutputEncoding
}
