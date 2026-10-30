package recovery

import (
	"crypto/aes"
	"crypto/md5"
	"crypto/subtle"
	"fmt"
)

// validateSIDBContainerAuthentication checks the traced current-container
// footer using its complete wire prefix and the native selection source. It
// does not establish the codec boundary or select a decoded registry.
func validateSIDBContainerAuthentication(container, source []byte) error {
	if len(source) < selectionSourcePrefixBytes+16 || len(source) > selectionSourcePrefixBytes+16*maximumMappedSelectionModulus || (len(source)-selectionSourcePrefixBytes)%16 != 0 || len(container) < aes.BlockSize {
		return fmt.Errorf("invalid SIDB authentication input geometry")
	}
	var key [aes.BlockSize]byte
	for i := range key {
		key[i] = applyCurrentFairPlayByteMask(source[i])
	}
	block, err := aes.NewCipher(key[:])
	clear(key[:])
	if err != nil {
		return err
	}
	prefix := container[:len(container)-aes.BlockSize]
	digest := md5.Sum(prefix)
	var expected [aes.BlockSize]byte
	block.Encrypt(expected[:], digest[:])
	valid := subtle.ConstantTimeCompare(expected[:], container[len(prefix):])
	clear(expected[:])
	clear(digest[:])
	if valid != 1 {
		return fmt.Errorf("SIDB container authentication failed")
	}
	return nil
}
