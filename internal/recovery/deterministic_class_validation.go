package recovery

import (
	"context"
	"crypto/aes"
	"crypto/sha256"
	"fmt"
)

// authenticateSelectedClassValue validates the one tagged record selected by
// the caller. It never changes registry, record, window, or key after failure.
func authenticateSelectedClassValue(ctx context.Context, sinf []byte, registries []decodedSIDBRegistryCandidate, recordIndex int, classKey [16]byte, progress func(Progress)) (authenticatedClassPrecursor, error) {
	defer clear(classKey[:])
	if ctx == nil {
		ctx = context.Background()
	}
	if err := ctx.Err(); err != nil {
		return authenticatedClassPrecursor{}, err
	}
	if len(registries) != 1 {
		return authenticatedClassPrecursor{}, fmt.Errorf("selected class requires one registry")
	}
	registry := registries[0]
	if recordIndex < 0 || recordIndex >= len(registry.Registry.Records) {
		return authenticatedClassPrecursor{}, fmt.Errorf("selected class record index is invalid: %d", recordIndex)
	}
	record := registry.Registry.Records[recordIndex]
	if record.MaskOffset < 0 || record.MaskOffset >= len(currentClassRecordMask) {
		return authenticatedClassPrecursor{}, fmt.Errorf("selected class record mask offset is invalid: %d", record.MaskOffset)
	}
	valu, err := deriveClassValuWindow(record, 0)
	if err != nil {
		return authenticatedClassPrecursor{}, err
	}
	defer clear(valu[:])
	probe, err := newCurrentSINFModelKeyProbe(sinf)
	if err != nil {
		return authenticatedClassPrecursor{}, err
	}
	event := Progress{Phase: "search-sinf", Total: 1, RegistryCandidates: 1, RecordCandidates: 1, WindowCandidates: 1}
	reportRecoveryProgress(progress, event)
	if err := ctx.Err(); err != nil {
		return authenticatedClassPrecursor{}, err
	}
	block, err := aes.NewCipher(classKey[:])
	if err != nil {
		return authenticatedClassPrecursor{}, err
	}
	precursor, accepted, err := authenticateClassValu(probe, sinf, block, valu)
	event.Completed = 1
	reportRecoveryProgress(progress, event)
	if ctxErr := ctx.Err(); ctxErr != nil {
		clear(precursor[:])
		return authenticatedClassPrecursor{}, ctxErr
	}
	if err != nil {
		clear(precursor[:])
		return authenticatedClassPrecursor{}, err
	}
	if !accepted {
		clear(precursor[:])
		return authenticatedClassPrecursor{}, errNoAuthenticatedClassPrecursor
	}
	selected := authenticatedClassPrecursor{
		classPrecursorCandidate: classPrecursorCandidate{RecordIndex: recordIndex, Record: record, Precursor: precursor},
		classKeyDigest:          sha256.Sum256(classKey[:]),
	}
	clear(precursor[:])
	event.Phase = "search-sinf-complete"
	reportRecoveryProgress(progress, event)
	return selected, nil
}
