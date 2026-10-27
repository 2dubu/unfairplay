package recovery

import "fmt"

const accountBytePackPass = uint64(0x1000ac8f0)

const (
	accountSHA1InputWords  = uint32(16)
	accountBytePackCounter = uint32(0x2ba7c38b)
)

func runAccountBytePack(s *state) (uint64, error) {
	for index := 0; index < 4096; index++ {
		next, err := runAccountBytePackStep(s)
		if err != nil || next != accountBytePackPass {
			return next, err
		}
	}
	return 0, fmt.Errorf("account byte-pack exceeded 4096 words")
}

// runAccountBytePackStep decodes one big-endian input word into the native
// SHA-1 schedule representation. The native Boolean circuit only hid the word
// index and the identity (word + mask - 2*(word&mask)) == word XOR mask.
func runAccountBytePackStep(s *state) (uint64, error) {
	s.traceBlock(accountBytePackPass)
	delta := uint32(s.x[4]) - accountBytePackCounter
	if delta&3 != 0 || delta/4 >= accountSHA1InputWords {
		return 0, fmt.Errorf("account SHA-1 input counter %#x is outside the block", uint32(s.x[4]))
	}
	wordIndex := delta / 4
	var word uint32
	for octet := uint32(0); octet < 4; octet++ {
		value, err := s.load(s.x[1]+uint64(wordIndex*4+octet), 1)
		if err != nil {
			return 0, fmt.Errorf("account SHA-1 input word %d byte %d: %w", wordIndex, octet, err)
		}
		word = word<<8 | uint32(byte(value)^byte(s.x[8]))
	}
	encodedWord := word ^ accountSHA1ScheduleMask
	if err := s.storeValue(s.x[10]+uint64(wordIndex*4), uint64(encodedWord), 4); err != nil {
		return 0, fmt.Errorf("account SHA-1 input word %d: %w", wordIndex, err)
	}

	// Preserve the native dispatcher ABI until the enclosing account loop is
	// replaced. These are scratch live-outs, not part of SHA-1.
	s.x[2] = 0x159f0000
	s.x[12] = 0xc285
	s.x[13] = 0x37514285
	s.x[16] = uint64([4]uint32{0x048e0000, 0x04860000, 0x048e0000, 0x04960000}[(wordIndex>>1)&3])
	s.x[15] = uint64(accountBytePackCounter + 80)
	s.x[4] += 4
	s.compare64(uint64((wordIndex+1)*4), 0x40)
	dispatchIndex := uint32(s.x[9])
	if wordIndex+1 < accountSHA1InputWords {
		dispatchIndex ^= 0xc285
	}
	next, err := loadBootstrapPointer(s, s.x[21], dispatchIndex)
	if err != nil {
		return 0, fmt.Errorf("account byte-pack dispatch: %w", err)
	}
	s.x[11] = next
	s.x[20] = uint64(accountBytePackCounter + accountSHA1InputWords)
	s.x[17] = 0x3601
	return next, nil
}
