package recovery

const (
	currentFairPlayByteMask = byte(0x1b)
	currentFairPlayWordMask = uint32(0xd56e60de)
)

// applyCurrentFairPlayByteMask converts one byte between the native protected
// representation and its plain representation. The mapping is involutive.
func applyCurrentFairPlayByteMask(value byte) byte {
	return value ^ currentFairPlayByteMask
}

// applyCurrentFairPlayWordMask is the four-byte form of the same protected
// representation. The byte lanes use masks d5, 6e, 60, and de.
func applyCurrentFairPlayWordMask(value uint32) uint32 {
	return value ^ currentFairPlayWordMask
}
