package recovery

import "fmt"

const (
	accountFrameWordWriterStart = uint64(0x10009d35c)
	accountFrameWordWriterTail  = uint64(0x10009d400)
	accountStackWordWriterStart = uint64(0x1000fd328)
	accountStackWordWriterTail  = uint64(0x1000fd400)
	accountWideWordWriter       = uint64(0x1001e4b10)
)

func writeAccountSHA1Digest(s *state, input, output uint64) error {
	// The native writer serializes the five encoded state words in big-endian
	// order, removes the state mask, then applies the byte-wise 0x1b transport
	// mask used by the surrounding certificate code.
	for wordIndex := uint64(0); wordIndex < 5; wordIndex++ {
		encoded, err := s.load(input+wordIndex*4, 4)
		if err != nil {
			return fmt.Errorf("account SHA-1 digest word %d: %w", wordIndex, err)
		}
		word := uint32(encoded) ^ accountSHA1StateMask
		for octet := uint64(0); octet < 4; octet++ {
			value := byte(word>>uint(24-octet*8)) ^ 0x1b
			if err := s.storeValue(output+wordIndex*4+octet, uint64(value), 1); err != nil {
				return fmt.Errorf("account SHA-1 digest byte %d: %w", wordIndex*4+octet, err)
			}
		}
	}
	return nil
}

func runAccountFrameWordWriter(s *state) (uint64, error) {
	if uint32(s.x[10]) != 0 {
		return 0, fmt.Errorf("account frame SHA-1 writer starts at %#x", uint32(s.x[10]))
	}
	limit, err := s.load(addOffset(s.x[29], -0x80), 4)
	if err != nil {
		return 0, fmt.Errorf("account frame SHA-1 limit: %w", err)
	}
	if uint32(limit) != 20 {
		return 0, fmt.Errorf("account frame SHA-1 length is %d, want 20", uint32(limit))
	}
	if err := writeAccountSHA1Digest(s, addOffset(s.x[29], -0x94), s.x[0]); err != nil {
		return 0, err
	}
	s.x[10] = 20
	s.compare32(20, uint32(limit))
	next, err := loadBootstrapPointer(s, s.x[21], uint32(s.x[9]))
	if err != nil {
		return 0, fmt.Errorf("account frame SHA-1 target: %w", err)
	}
	s.x[11] = next
	s.x[17] = 0x98cd
	return next, nil
}

func runAccountFrameWordWriterTail(_ *state) (uint64, error) {
	return 0, fmt.Errorf("unexpected account frame SHA-1 writer continuation")
}

func runAccountStackWordWriter(s *state) (uint64, error) {
	if uint32(s.x[10]) != 0 {
		return 0, fmt.Errorf("account stack SHA-1 writer starts at %#x", uint32(s.x[10]))
	}
	limit, err := s.load(s.sp+0xf8, 4)
	if err != nil {
		return 0, fmt.Errorf("account stack SHA-1 limit: %w", err)
	}
	if uint32(limit) != 20 {
		return 0, fmt.Errorf("account stack SHA-1 length is %d, want 20", uint32(limit))
	}
	if err := writeAccountSHA1Digest(s, s.sp+0xe4, s.x[16]); err != nil {
		return 0, err
	}
	s.x[10] = 20
	s.compare32(20, uint32(limit))
	dispatchIndex := uint32(s.x[9]) | 0x15
	next, err := s.load(s.x[20]+uint64(dispatchIndex)*8, 8)
	if err != nil {
		return 0, fmt.Errorf("account stack SHA-1 target: %w", err)
	}
	s.x[11] = next
	s.x[17] = 0x9986
	return next, nil
}

func runAccountStackWordWriterTail(_ *state) (uint64, error) {
	return 0, fmt.Errorf("unexpected account stack SHA-1 writer continuation")
}

func runAccountWideWordWriter(s *state) (uint64, error) {
	if uint32(s.x[10]) != 0 {
		return 0, fmt.Errorf("account wide SHA-1 writer starts at %#x", uint32(s.x[10]))
	}
	limit, err := s.load(s.sp+0x2f8, 4)
	if err != nil {
		return 0, fmt.Errorf("account wide SHA-1 limit: %w", err)
	}
	if uint32(limit) != 20 {
		return 0, fmt.Errorf("account wide SHA-1 length is %d, want 20", uint32(limit))
	}
	if err := writeAccountSHA1Digest(s, s.sp+0x2e4, s.x[16]); err != nil {
		return 0, err
	}
	s.x[10] = 20
	s.compare32(20, uint32(limit))
	dispatchIndex := uint32(s.x[9]) ^ 0x1f6
	next, err := loadBootstrapPointer(s, s.x[24], dispatchIndex)
	if err != nil {
		return 0, fmt.Errorf("account wide SHA-1 target: %w", err)
	}
	s.x[11] = next
	s.x[17] = 0x5b8c
	return next, nil
}
