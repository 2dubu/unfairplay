package recovery

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"testing"
)

func TestSelectedEnrollmentKeyIgnoresLaterBlocks(t *testing.T) {
	source := make([]byte, 18*16)
	for i := range source {
		source[i] = byte(i*41 + i/16*13 + 5)
	}
	for nativeCount := 1; nativeCount <= 18; nativeCount++ {
		for slot := 0; slot < nativeCount; slot++ {
			key, err := deriveSelectedEnrollmentKey(source[:nativeCount*16], slot)
			if err != nil {
				t.Fatal(err)
			}
			changed := append([]byte(nil), source[:nativeCount*16]...)
			for i := (slot + 1) * 16; i < len(changed); i++ {
				changed[i] ^= 0xff
			}
			again, err := deriveSelectedEnrollmentKey(changed, slot)
			if err != nil || again != key {
				t.Fatalf("slot %d depends on later inputs: %v", slot, err)
			}
		}
	}
}
func TestSelectedEnrollmentKeyRejectsInvalidInput(t *testing.T) {
	for _, length := range []int{0, 15, 17, 19 * 16} {
		if _, err := deriveSelectedEnrollmentKey(make([]byte, length), 0); err == nil {
			t.Fatalf("accepted length %d", length)
		}
	}
	source := make([]byte, 16)
	for _, slot := range []int{-1, 1} {
		if _, err := deriveSelectedEnrollmentKey(source, slot); err == nil {
			t.Fatalf("accepted slot %d", slot)
		}
	}
	if _, err := deriveSelectedEnrollmentKey(make([]byte, selectionStateSlotCount*16), selectionStateSlotCount); err == nil {
		t.Fatalf("accepted slot %d", selectionStateSlotCount)
	}
}

func TestSelectedEnrollmentKeyMatchesCBCForAllMappedSlots(t *testing.T) {
	source := make([]byte, selectionStateSlotCount*aes.BlockSize)
	whitened := make([]byte, len(source))
	for index := range source {
		source[index] = byte(index*41 + index/aes.BlockSize*13 + 5)
		whitened[index] = source[index] ^ selectionStateInputWhitening
	}
	block, err := aes.NewCipher(selectionStateAESKey[:])
	if err != nil {
		t.Fatal(err)
	}
	states := make([]byte, len(whitened))
	cipher.NewCBCEncrypter(block, selectionStateIV[:]).CryptBlocks(states, whitened)
	for slot := 0; slot < selectionStateSlotCount; slot++ {
		got, err := deriveSelectedEnrollmentKey(source, slot)
		if err != nil {
			t.Fatalf("slot %d: %v", slot, err)
		}
		var want [aes.BlockSize]byte
		selectionSIDBCipher.Encrypt(want[:], states[slot*aes.BlockSize:(slot+1)*aes.BlockSize])
		if !bytes.Equal(got[:], want[:]) {
			t.Fatalf("slot %d differs from the CBC recurrence", slot)
		}
	}
}
