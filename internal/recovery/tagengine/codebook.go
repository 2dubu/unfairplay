package tagengine

import (
	"context"
	"encoding/binary"
	"fmt"
)

// TagCode is one semantic FourCC and its hardware-dependent serialized tag.
type TagCode struct {
	Semantic string
	Raw      uint32
}

// GenerateCodebook computes the reduced constructor for the explicit
// six-byte hardware identity. The arithmetic and public constants follow the mapped constructor contract.
func GenerateCodebook(hardware []byte) ([]TagCode, error) {
	return GenerateCodebookWithContext(context.Background(), hardware)
}

// GenerateCodebookWithContext permits cancellation throughout the constructor. Returned entries remain owned by the caller.
func GenerateCodebookWithContext(ctx context.Context, hardware []byte) ([]TagCode, error) {
	if ctx == nil {
		return nil, fmt.Errorf("nil context")
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	state, err := constructorStateWithContext(ctx, hardware)
	if err != nil {
		return nil, err
	}
	defer clear(state[:])
	// Both native arithmetic encodings reduce to (x XOR c) - c for the
	// same public c. Their inverse and transport XOR cancel across this
	// boundary, leaving an affine map of each decoded AES output word.
	var result [44]TagCode
	defer clear(result[:])
	seenRaw := map[uint32]bool{}
	defer clear(seenRaw)
	for i := range result {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		raw := binary.BigEndian.Uint32(state[4*i:])*0x75f58395 + 0x84bed201
		name := semanticTags[i]
		if seenRaw[raw] {
			return nil, fmt.Errorf("ambiguous constructor codebook")
		}
		seenRaw[raw] = true
		result[i] = TagCode{name, raw}
	}
	return append([]TagCode(nil), result[:]...), nil
}

// Semantic order is the public 176-byte reference dictionary, decoded once.
var semanticTags = [44]string{
	"sean", "vers", "grup", "head", "guid", "user", "sess", "key ",
	"valu", "type", "salt", "iviv", "sign", "tail", "itgr", "rand",
	"usag", "flur", "dong", "dval", "dtyp", "dsal", "dgui", "rbag",
	"time", "cfid", "poli", "leew", "dirs", "dstg", "filt", "base",
	"suff", "poli", "rtal", "rori", "rwri", "gval", "dval", "qhtz",
	"hi32", "drml", "mpos", "dmid",
}
