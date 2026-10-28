package recovery

import "fmt"

const maximumCurrentSINFRightsBytes = 256

// validateCurrentSINFRights accepts the two observed rights layouts and
// bounded extensions. The complete authenticated atom sequence, not the
// opaque rights labels, identifies a valid decrypted private record.
func validateCurrentSINFRights(rights []byte) error {
	if len(rights) < 80 || len(rights) > maximumCurrentSINFRightsBytes || len(rights)%8 != 0 {
		return fmt.Errorf("current SINF rights have invalid geometry: %d bytes", len(rights))
	}
	return nil
}
