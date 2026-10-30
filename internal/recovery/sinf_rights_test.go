package recovery

import "testing"

func TestValidateCurrentSINFRightsLayouts(t *testing.T) {
	for _, size := range []int{80, 88, 96, maximumCurrentSINFRightsBytes} {
		rights := make([]byte, size)
		for index := range rights {
			rights[index] = byte(index*31 + 9)
		}
		if err := validateCurrentSINFRights(rights); err != nil {
			t.Fatalf("rights size %d: %v", size, err)
		}
	}
}

func TestValidateCurrentSINFRightsRejectsInvalidLayouts(t *testing.T) {
	for _, rights := range [][]byte{
		make([]byte, 79),
		make([]byte, 81),
		make([]byte, maximumCurrentSINFRightsBytes+8),
	} {
		if err := validateCurrentSINFRights(rights); err == nil {
			t.Fatalf("accepted invalid rights geometry of %d bytes", len(rights))
		}
	}
}
