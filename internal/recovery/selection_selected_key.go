package recovery

import "fmt"

// deriveSelectedEnrollmentKey advances only through the specified native slot;
// sourceTail contains the source window, in complete 16-byte blocks.
func deriveSelectedEnrollmentKey(sourceTail []byte, slot int) ([16]byte, error) {
	var empty [16]byte
	if len(sourceTail) == 0 || len(sourceTail)%16 != 0 || len(sourceTail)/16 > 18 {
		return empty, fmt.Errorf("invalid selection state source length %d", len(sourceTail))
	}
	blocks := len(sourceTail) / 16
	if slot < 0 || slot >= blocks {
		return empty, fmt.Errorf("selection state slot %d is outside source block count %d", slot, blocks)
	}
	if slot >= selectionStateSlotCount {
		return empty, fmt.Errorf("selection state slot %d has no verified transform", slot)
	}
	previous := selectionStateIV
	defer clear(previous[:])
	for index := 0; index <= slot; index++ {
		advanceSelectionState(&previous, sourceTail[index*16:(index+1)*16])
	}
	var result [16]byte
	selectionSIDBCipher.Encrypt(result[:], previous[:])
	return result, nil
}
