package recovery

import "fmt"

// nativeSIDBRootCursor evaluates the mapped root-offset arithmetic from the
// six hardware bytes and the public initializer. It does not inspect SIDB data.
func nativeSIDBRootCursor(hardware []byte) (int, error) {
	if len(hardware) != 6 {
		return 0, fmt.Errorf("invalid hardware identity length %d", len(hardware))
	}
	// Fold the native byte offsets and XOR/add encodings before combining
	// the affine terms. The polynomial's constant term cancels to zero.
	factors := [...]uint32{0, 13, 15, 11, 19, 9}
	multiplier := uint32(0x5f64ac89)
	var hash, polynomial uint32
	for i, value := range hardware {
		v := uint32(value)
		hash ^= v*multiplier + 0x1b4591cd
		polynomial = polynomial*factors[i] + v
		multiplier += 0x34139
	}
	return 1 + int((hash^polynomial)%12), nil
}
