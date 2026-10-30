package recovery

import (
	"context"
	"fmt"
)

// authenticateTaggedSelection binds one parsed route to one pre-parsed
// registry and authenticates exactly the selected class tuple. The caller
// supplies a validated decoder and the native root cursor.
func authenticateTaggedSelection(ctx context.Context, decode sidbTagDecoder, sinf []byte, registries []decodedSIDBRegistryCandidate, hardware []byte, rootCursor int, user, keyID uint32, progress func(Progress)) (authenticatedClassPrecursor, error) {
	if ctx != nil {
		if err := ctx.Err(); err != nil {
			return authenticatedClassPrecursor{}, err
		}
	}
	if len(registries) != 1 || !registries[0].Registry.Tree || len(registries[0].Data) == 0 {
		return authenticatedClassPrecursor{}, fmt.Errorf("tagged SIDB selection requires one decoded tree")
	}
	parsedUser, parsedKey, err := enrollmentSINFSelector(sinf)
	if err != nil {
		return authenticatedClassPrecursor{}, err
	}
	if parsedUser != user || parsedKey != keyID {
		return authenticatedClassPrecursor{}, fmt.Errorf("SINF selector does not match requested selector")
	}
	route, err := selectTaggedSIDBRoute(decode, sinf, registries[0].Data, rootCursor)
	if err != nil {
		return authenticatedClassPrecursor{}, err
	}
	if route.RecordIndex < 0 || route.RecordIndex >= len(registries[0].Registry.Records) || registries[0].Registry.Records[route.RecordIndex] != route.Record {
		return authenticatedClassPrecursor{}, fmt.Errorf("selected SIDB record disagrees with registry")
	}
	key, err := accountClassKey(hardware, user, keyID)
	if err != nil {
		return authenticatedClassPrecursor{}, err
	}
	defer clear(key[:])
	return authenticateSelectedClassWindow(ctx, sinf, registries, 0, route.RecordIndex, 0, key, progress)
}
