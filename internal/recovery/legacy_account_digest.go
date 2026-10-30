package recovery

import (
	"crypto/sha1" // #nosec G505 -- exact FairPlay compatibility transform, not a signature primitive.
	"encoding/binary"
	"fmt"
)

var legacyAccountDigestPrefix = [17]byte{
	0x5b, 0x53, 0xba, 0xe2, 0xa6, 0x6f, 0x01, 0x67,
	0x73, 0x29, 0x10, 0xc5, 0x20, 0x99, 0xb7, 0x33, 0x81,
}

var legacyAccountDigestSuffix = [12]byte{
	0x54, 0xda, 0x7b, 0x97, 0xe1, 0xf9,
	0x0e, 0xc3, 0xae, 0x64, 0x6e, 0x3a,
}

// legacyAccountDigest constructs the semantic 45-byte identity message used
// by the older Mac enrollment family and returns its ordinary SHA-1 digest.
// The eight-byte identity is the six-byte built-in IOMACAddress followed by
// two zero bytes; the selector fields come from the personalized SINF.
func legacyAccountDigest(hardware []byte, user, keyID uint32) ([sha1.Size]byte, error) {
	if len(hardware) != 6 {
		return [sha1.Size]byte{}, fmt.Errorf("legacy class hardware identity has %d bytes, want 6", len(hardware))
	}
	if keyID > 0xff {
		return [sha1.Size]byte{}, fmt.Errorf("legacy class key selector %#x exceeds one byte", keyID)
	}
	var message [45]byte
	defer clear(message[:])
	offset := copy(message[:], legacyAccountDigestPrefix[:])
	offset += copy(message[offset:], hardware)
	offset += 2 // Complete the eight-byte hardware field with zero bytes.
	offset++    // One reserved zero byte precedes the selector.
	message[offset] = byte(keyID)
	offset++
	binary.BigEndian.PutUint32(message[offset:], user)
	offset += 4
	message[offset], message[offset+1] = 0xff, 0xff
	offset += 2
	copy(message[offset:], legacyAccountDigestSuffix[:])
	return sha1.Sum(message[:]), nil
}
