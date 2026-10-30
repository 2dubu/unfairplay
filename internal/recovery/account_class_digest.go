package recovery

import (
	"crypto/aes"
	"crypto/sha1"
	"encoding/binary"
	"fmt"
)

// Semantic constants originate in reference fairplayd SHA-256
// 1e58d37e9602dd0ca5e24b14762063b6c310248b4d3a94272771d815e4ea1784.
// Native digest-prefix/suffix sources: 0x100ea9140, 0x100f51d30;
// AES key source: 0x100fa3af0. Writer and AES wrapper: 0x1005501bc,
// 0x100535754. These are algorithm constants, not account inputs.
var accountClassDigestPrefix = [17]byte{0x5b, 0x53, 0xba, 0xe2, 0xa6, 0x6f, 0x01, 0x67, 0x73, 0x29, 0x10, 0xc5, 0x20, 0x99, 0xb7, 0x33, 0x81}
var accountClassDigestSuffix = [12]byte{0x54, 0xda, 0x7b, 0x97, 0xe1, 0xf9, 0x0e, 0xc3, 0xae, 0x64, 0x6e, 0x3a}
var accountClassAESKey = [16]byte{0x75, 0x8f, 0xa5, 0x2a, 0xe1, 0xb0, 0x73, 0x0b, 0x84, 0x47, 0x87, 0x57, 0xf0, 0x23, 0xcd, 0x38}

// accountClassDigest hashes the framed hardware identity and SINF selector.
func accountClassDigest(hardware []byte, user, keyID uint32) ([20]byte, error) {
	if len(hardware) != 6 {
		return [20]byte{}, fmt.Errorf("account class hardware has %d bytes, want 6", len(hardware))
	}
	var message [45]byte
	copy(message[:17], accountClassDigestPrefix[:])
	copy(message[17:23], hardware)
	binary.BigEndian.PutUint32(message[23:27], keyID)
	binary.BigEndian.PutUint32(message[27:31], user)
	message[31], message[32] = 0xff, 0xff
	copy(message[33:], accountClassDigestSuffix[:])
	return sha1.Sum(message[:]), nil
}

// accountClassFirstOperand is the semantic first eight AES output bytes.
func accountClassFirstOperand(digest [20]byte) ([8]byte, error) {
	var result [8]byte
	block, err := aes.NewCipher(accountClassAESKey[:])
	if err != nil {
		return result, err
	}
	var output [16]byte
	block.Encrypt(output[:], digest[:16])
	copy(result[:], output[:8])
	return result, nil
}
