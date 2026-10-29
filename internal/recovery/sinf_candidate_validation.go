package recovery

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"encoding/binary"
	"fmt"
	"runtime"
	"sync"
)

type authenticatedClassPrecursor struct {
	classPrecursorCandidate
	RegistryIndex int
}

func currentSINFPrivateSource(sinf []byte) ([]byte, error) {
	rootStart, rootEnd, err := enrollmentAtom(sinf, 0, len(sinf), "sinf")
	if err != nil || rootStart != 8 || rootEnd != len(sinf) {
		return nil, fmt.Errorf("SINF root is invalid")
	}
	schiStart, schiEnd, err := enrollmentAtom(sinf, rootStart, rootEnd, "schi")
	if err != nil {
		return nil, err
	}
	privateStart, privateEnd, err := enrollmentAtom(sinf, schiStart, schiEnd, "priv")
	if err != nil {
		return nil, err
	}
	return sinf[privateStart:privateEnd], nil
}

func validateCurrentSINFPrecursor(sinf []byte, precursor [16]byte) (bool, error) {
	logical, ciphertextBytes, trailerBytes, err := currentSINFLogicalSource(sinf)
	if err != nil {
		return false, err
	}
	defer clear(logical)

	plaintext := make([]byte, ciphertextBytes+trailerBytes)
	decrypted := plaintext[:ciphertextBytes]
	defer clear(plaintext)
	key := currentSINFModelKeyFromPrecursor(precursor)
	block, err := aes.NewCipher(key[:])
	clear(key[:])
	if err != nil {
		return false, err
	}
	cipher.NewCBCDecrypter(block, logical[:aes.BlockSize]).CryptBlocks(decrypted, logical[aes.BlockSize:])

	expected := [...]struct {
		typ  string
		size int
	}{
		{"key ", 24},
		{"iviv", 24},
		{"tmpr", 24},
		{"righ", 0},
		{"name", 264},
	}
	cursor := 0
	for _, atom := range expected {
		if len(plaintext)-cursor < 8 {
			return false, nil
		}
		size := int(binary.BigEndian.Uint32(plaintext[cursor : cursor+4]))
		if size < 8 || size > len(plaintext)-cursor || string(plaintext[cursor+4:cursor+8]) != atom.typ {
			return false, nil
		}
		if atom.typ == "righ" {
			if err := validateCurrentSINFRights(plaintext[cursor+8 : cursor+size]); err != nil {
				return false, nil
			}
		} else if size != atom.size {
			return false, nil
		}
		cursor += size
	}
	return cursor == len(plaintext), nil
}

func currentSINFLogicalSource(sinf []byte) ([]byte, int, int, error) {
	source, err := currentSINFPrivateSource(sinf)
	if err != nil {
		return nil, 0, 0, err
	}
	if len(source) < 48 {
		return nil, 0, 0, fmt.Errorf("SINF private source is too short")
	}
	bodyBytes := len(source) - aes.BlockSize
	trailerBytes := bodyBytes % aes.BlockSize
	ciphertextBytes := bodyBytes - trailerBytes
	if ciphertextBytes < 2*aes.BlockSize || trailerBytes != 0 && trailerBytes != 8 {
		return nil, 0, 0, fmt.Errorf("SINF private block geometry is invalid")
	}
	encodedBytes := aes.BlockSize + ciphertextBytes
	for _, value := range source[encodedBytes:] {
		if value != 0 {
			return nil, 0, 0, fmt.Errorf("SINF private trailer is not zero")
		}
	}

	logical := append([]byte(nil), source[:encodedBytes]...)
	return logical, ciphertextBytes, trailerBytes, nil
}

type currentSINFModelKeyProbe struct {
	iv         [aes.BlockSize]byte
	ciphertext [aes.BlockSize]byte
}

func newCurrentSINFModelKeyProbe(sinf []byte) (currentSINFModelKeyProbe, error) {
	var result currentSINFModelKeyProbe
	logical, _, _, err := currentSINFLogicalSource(sinf)
	if err != nil {
		return result, err
	}
	defer clear(logical)
	copy(result.iv[:], logical[:aes.BlockSize])
	copy(result.ciphertext[:], logical[aes.BlockSize:2*aes.BlockSize])
	return result, nil
}

func (probe currentSINFModelKeyProbe) accepts(modelKey [aes.BlockSize]byte) (bool, error) {
	block, err := aes.NewCipher(modelKey[:])
	if err != nil {
		return false, err
	}
	var first [aes.BlockSize]byte
	block.Decrypt(first[:], probe.ciphertext[:])
	for index := range first {
		first[index] ^= probe.iv[index]
	}
	return bytes.Equal(first[:8], []byte{0, 0, 0, 24, 'k', 'e', 'y', ' '}), nil
}

func authenticateClassValu(probe currentSINFModelKeyProbe, sinf []byte, block cipher.Block, valu [aes.BlockSize]byte) ([aes.BlockSize]byte, bool, error) {
	var modelKey [aes.BlockSize]byte
	block.Decrypt(modelKey[:], valu[:])
	possible, err := probe.accepts(modelKey)
	if err != nil || !possible {
		clear(modelKey[:])
		return [aes.BlockSize]byte{}, false, err
	}
	var precursor [aes.BlockSize]byte
	for index := range precursor {
		precursor[index] = applyCurrentFairPlayByteMask(modelKey[index])
	}
	clear(modelKey[:])
	accepted, err := validateCurrentSINFPrecursor(sinf, precursor)
	if err != nil || !accepted {
		clear(precursor[:])
		return [aes.BlockSize]byte{}, false, err
	}
	return precursor, true, nil
}

type preparedClassCandidate struct {
	registryIndex int
	recordIndex   int
	windowOffset  int
	valu          [aes.BlockSize]byte
}

const (
	classSearchChunkValues         = 64
	maximumClassSearchWorkers      = 16
	maximumAuthenticatedClassPairs = uint64(1) << 31
)

func recoveryCandidateCounts(registries []decodedSIDBRegistryCandidate) (records, windows int) {
	for _, registry := range registries {
		records += len(registry.Registry.Records)
	}
	return records, records * (maximumClassWindowOffset + 1)
}

func reportRecoveryProgress(progress func(Progress), event Progress) {
	if progress != nil {
		progress(event)
	}
}

func selectAuthenticatedClassPrecursor(sinf []byte, registries []decodedSIDBRegistryCandidate) (authenticatedClassPrecursor, error) {
	return selectAuthenticatedClassPrecursorWithContext(context.Background(), sinf, registries, nil)
}

func selectAuthenticatedClassPrecursorWithContext(ctx context.Context, sinf []byte, registries []decodedSIDBRegistryCandidate, progress func(Progress)) (authenticatedClassPrecursor, error) {
	var selected authenticatedClassPrecursor
	if ctx == nil {
		ctx = context.Background()
	}
	if err := ctx.Err(); err != nil {
		return selected, err
	}
	probe, err := newCurrentSINFModelKeyProbe(sinf)
	if err != nil {
		return selected, err
	}
	var candidates []preparedClassCandidate
	for registryIndex, registry := range registries {
		for recordIndex, record := range registry.Registry.Records {
			for windowOffset := 0; windowOffset <= maximumClassWindowOffset; windowOffset++ {
				valu, err := deriveClassValuWindow(record, windowOffset)
				if err != nil {
					return selected, err
				}
				candidates = append(candidates, preparedClassCandidate{
					registryIndex: registryIndex,
					recordIndex:   recordIndex,
					windowOffset:  windowOffset,
					valu:          valu,
				})
			}
		}
	}
	defer func() {
		for index := range candidates {
			clear(candidates[index].valu[:])
		}
	}()
	classValues := make([][8]byte, 0, classWorkspaceModulus+1)
	classValues = append(classValues, currentClassSecondaryValue)
	seenClassValues := map[[8]byte]bool{currentClassSecondaryValue: true}
	for seed := range uint32(classWorkspaceModulus) {
		value := classWorkspaceOutputFromSeed(uint16(seed))
		if seenClassValues[value] {
			continue
		}
		seenClassValues[value] = true
		classValues = append(classValues, value)
	}
	total := uint64(len(classValues)) * uint64(len(candidates))
	if total == 0 || total > maximumAuthenticatedClassPairs {
		return selected, fmt.Errorf("authenticated SINF class search has unsupported size: %d comparisons", total)
	}
	records, windows := recoveryCandidateCounts(registries)
	baseProgress := Progress{
		Phase: "search-sinf", Total: total, RegistryCandidates: len(registries),
		RecordCandidates: records, WindowCandidates: windows,
	}
	reportRecoveryProgress(progress, baseProgress)

	type searchTask struct{ start, end int }
	type searchResult struct {
		completed uint64
		matches   []authenticatedClassPrecursor
		err       error
	}
	workerCount := runtime.GOMAXPROCS(0)
	if workerCount > maximumClassSearchWorkers {
		workerCount = maximumClassSearchWorkers
	}
	chunks := (len(classValues) + classSearchChunkValues - 1) / classSearchChunkValues
	if workerCount > chunks {
		workerCount = chunks
	}
	searchContext, cancel := context.WithCancel(ctx)
	defer cancel()
	tasks := make(chan searchTask, chunks)
	results := make(chan searchResult, workerCount)
	for start := 0; start < len(classValues); start += classSearchChunkValues {
		end := start + classSearchChunkValues
		if end > len(classValues) {
			end = len(classValues)
		}
		tasks <- searchTask{start: start, end: end}
	}
	close(tasks)
	var workers sync.WaitGroup
	workers.Add(workerCount)
	for range workerCount {
		go func() {
			defer workers.Done()
			for task := range tasks {
				if err := searchContext.Err(); err != nil {
					return
				}
				result := searchResult{}
				for classIndex := task.start; classIndex < task.end; classIndex++ {
					if searchContext.Err() != nil {
						break
					}
					parameters := classEnrollmentParameters{ClassValue: classValues[classIndex]}
					classKey := parameters.classKey()
					classBlock, blockErr := aes.NewCipher(classKey[:])
					clear(classKey[:])
					if blockErr != nil {
						result.err = blockErr
						break
					}
					for _, candidate := range candidates {
						precursor, accepted, candidateErr := authenticateClassValu(probe, sinf, classBlock, candidate.valu)
						result.completed++
						if candidateErr != nil {
							result.err = candidateErr
							break
						}
						if accepted {
							result.matches = append(result.matches, authenticatedClassPrecursor{
								classPrecursorCandidate: classPrecursorCandidate{
									RecordIndex:  candidate.recordIndex,
									WindowOffset: candidate.windowOffset,
									Record:       registries[candidate.registryIndex].Registry.Records[candidate.recordIndex],
									Precursor:    precursor,
								},
								RegistryIndex: candidate.registryIndex,
							})
						}
					}
					if result.err != nil {
						break
					}
				}
				results <- result
				if result.err != nil {
					cancel()
					return
				}
			}
		}()
	}
	go func() {
		workers.Wait()
		close(results)
	}()
	completed := uint64(0)
	var matches []authenticatedClassPrecursor
	var firstError error
	progressInterval := total / 32
	if progressInterval == 0 {
		progressInterval = 1
	}
	nextProgress := progressInterval
	for result := range results {
		completed += result.completed
		if result.err != nil && firstError == nil {
			firstError = result.err
			cancel()
		}
		matches = append(matches, result.matches...)
		if completed >= nextProgress || completed == total {
			event := baseProgress
			event.Completed = completed
			reportRecoveryProgress(progress, event)
			for nextProgress <= completed {
				nextProgress += progressInterval
			}
		}
	}
	if firstError != nil {
		for index := range matches {
			clear(matches[index].Precursor[:])
		}
		return selected, firstError
	}
	if err := ctx.Err(); err != nil {
		for index := range matches {
			clear(matches[index].Precursor[:])
		}
		return selected, err
	}
	if len(matches) != 1 {
		for index := range matches {
			clear(matches[index].Precursor[:])
		}
		if len(matches) > 1 {
			return authenticatedClassPrecursor{}, fmt.Errorf("authenticated SINF accepted multiple class precursors")
		}
		return authenticatedClassPrecursor{}, fmt.Errorf("authenticated SINF accepted no class precursor")
	}
	selected = matches[0]
	event := baseProgress
	event.Completed = total
	event.Phase = "search-sinf-complete"
	reportRecoveryProgress(progress, event)
	return selected, nil
}

func selectAuthenticatedLegacyClassPrecursor(sinf []byte, registries []decodedSIDBRegistryCandidate, hardware []byte, user, keyID uint32) (authenticatedClassPrecursor, error) {
	return selectAuthenticatedLegacyClassPrecursorWithContext(context.Background(), sinf, registries, hardware, user, keyID, nil)
}

func selectAuthenticatedLegacyClassPrecursorWithContext(ctx context.Context, sinf []byte, registries []decodedSIDBRegistryCandidate, hardware []byte, user, keyID uint32, progress func(Progress)) (authenticatedClassPrecursor, error) {
	var selected authenticatedClassPrecursor
	if ctx == nil {
		ctx = context.Background()
	}
	if err := ctx.Err(); err != nil {
		return selected, err
	}
	digest, err := legacyAccountDigest(hardware, user, keyID)
	if err != nil {
		return selected, err
	}
	parameters := classEnrollmentParameters{ClassValue: legacyClassWorkspaceValue(digest)}
	classKey := parameters.classKey()
	classBlock, err := aes.NewCipher(classKey[:])
	clear(classKey[:])
	clear(digest[:])
	if err != nil {
		return selected, err
	}
	probe, err := newCurrentSINFModelKeyProbe(sinf)
	if err != nil {
		return selected, err
	}
	records, _ := recoveryCandidateCounts(registries)
	total := uint64(records * len(currentClassRecordMask))
	baseProgress := Progress{
		Phase: "search-sinf", Total: total, RegistryCandidates: len(registries),
		RecordCandidates: records, WindowCandidates: records * len(currentClassRecordMask),
	}
	reportRecoveryProgress(progress, baseProgress)
	progressInterval := total / 32
	if progressInterval == 0 {
		progressInterval = 1
	}
	nextProgress := progressInterval
	completed := uint64(0)
	matches := 0
	for registryIndex, registry := range registries {
		for recordIndex, record := range registry.Registry.Records {
			if err := ctx.Err(); err != nil {
				clear(selected.Precursor[:])
				return authenticatedClassPrecursor{}, err
			}
			for maskPhase := range len(currentClassRecordMask) {
				valu, err := deriveLegacyClassValu(record, maskPhase)
				if err != nil {
					return selected, err
				}
				precursor, accepted, err := authenticateClassValu(probe, sinf, classBlock, valu)
				completed++
				if completed >= nextProgress || completed == total {
					event := baseProgress
					event.Completed = completed
					reportRecoveryProgress(progress, event)
					for nextProgress <= completed {
						nextProgress += progressInterval
					}
				}
				if err != nil {
					return selected, err
				}
				if !accepted {
					continue
				}
				matches++
				selected = authenticatedClassPrecursor{
					classPrecursorCandidate: classPrecursorCandidate{
						RecordIndex:  recordIndex,
						WindowOffset: maskPhase,
						Record:       record,
						Precursor:    precursor,
					},
					RegistryIndex: registryIndex,
				}
				if matches > 1 {
					clear(selected.Precursor[:])
					return authenticatedClassPrecursor{}, fmt.Errorf("authenticated SINF accepted multiple legacy class precursors")
				}
			}
		}
	}
	if matches != 1 {
		return authenticatedClassPrecursor{}, fmt.Errorf("authenticated SINF accepted no legacy class precursor")
	}
	event := baseProgress
	event.Completed = total
	event.Phase = "search-sinf-complete"
	reportRecoveryProgress(progress, event)
	return selected, nil
}

func selectAuthenticatedEnrollmentClassPrecursor(sinf []byte, registries []decodedSIDBRegistryCandidate, hardware []byte, user, keyID uint32) (authenticatedClassPrecursor, error) {
	return selectAuthenticatedEnrollmentClassPrecursorWithContext(context.Background(), sinf, registries, hardware, user, keyID, nil)
}

func selectAuthenticatedEnrollmentClassPrecursorWithContext(ctx context.Context, sinf []byte, registries []decodedSIDBRegistryCandidate, hardware []byte, user, keyID uint32, progress func(Progress)) (authenticatedClassPrecursor, error) {
	if len(registries) == 0 {
		return authenticatedClassPrecursor{}, fmt.Errorf("class precursor selection requires one decoded registry")
	}
	older := registries[0].Registry.FooterBytes == decodedSIDBOlderFooterBytes
	for _, registry := range registries[1:] {
		if (registry.Registry.FooterBytes == decodedSIDBOlderFooterBytes) != older {
			return authenticatedClassPrecursor{}, fmt.Errorf("class precursor selection found mixed registry layouts")
		}
	}
	if older {
		return selectAuthenticatedLegacyClassPrecursorWithContext(ctx, sinf, registries, hardware, user, keyID, progress)
	}
	return selectAuthenticatedClassPrecursorWithContext(ctx, sinf, registries, progress)
}
