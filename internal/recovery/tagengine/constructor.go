package tagengine

import (
	"context"
	"fmt"
)

// constructorStateWithContext returns the decoded AES output. The caller owns
// the returned state and must clear it after consuming the tag words.
func constructorStateWithContext(ctx context.Context, hardware []byte) ([256]byte, error) {
	var empty [256]byte
	if ctx == nil {
		return empty, fmt.Errorf("nil context")
	}
	if err := ctx.Err(); err != nil {
		return empty, err
	}
	if len(hardware) != 6 {
		return empty, fmt.Errorf("invalid hardware length")
	}
	var identity [6]byte
	copy(identity[:], hardware)
	defer clear(identity[:])
	state := initialState(identity[:])
	defer clear(state[:])
	for i := 0; i < mixRounds(identity[:]); i++ {
		if err := ctx.Err(); err != nil {
			return empty, err
		}
		mixStep(&state, i)
	}
	if err := transformConstructorState(ctx, &state); err != nil {
		return empty, err
	}
	return state, nil
}
