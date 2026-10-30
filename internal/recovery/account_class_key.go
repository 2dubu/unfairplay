package recovery

import "encoding/binary"

// accountClassKey derives the native repeated-half class key from runtime
// enrollment inputs. Each workspace is initialized and scheduled separately.
func accountClassKey(hardware []byte, user, keyID uint32) ([16]byte, error) {
	var key [16]byte
	digest, err := accountClassDigest(hardware, user, keyID)
	if err != nil {
		return key, err
	}
	defer clear(digest[:])
	return accountClassKeyFromDigest(digest)
}

// accountClassKeyFromDigest is the shared, deterministic workspace and join
// path; digest construction and enrollment input validation remain upstream.
func accountClassKeyFromDigest(digest [20]byte) ([16]byte, error) {
	defer clear(digest[:])
	var key [16]byte
	first, err := accountClassFirstOperand(digest)
	if err != nil {
		return key, err
	}
	var second [8]byte
	for pass := range 2 {
		word := binary.BigEndian.Uint16(digest[16+2*pass:])
		state := classWorkspaceInitialize(word)
		if err := classWorkspaceSchedulePass(&state); err != nil {
			return key, err
		}
		binary.BigEndian.PutUint32(second[4*pass:], classWorkspacePolynomialWord(state.Product))
	}
	for i := range 8 {
		var value byte
		for j := range 8 {
			value ^= classWorkspaceGaloisMultiply(first[(i+j)%8], second[j])
		}
		key[i], key[i+8] = value, value
	}
	return key, nil
}
