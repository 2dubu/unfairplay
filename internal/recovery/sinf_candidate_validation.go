package recovery

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	"runtime"
	"sync"
)

var errNoAuthenticatedClassPrecursor = errors.New("authenticated SINF accepted no class precursor")

type authenticatedClassPrecursor struct {
	classPrecursorCandidate
	RegistryIndex  int
	classKeyDigest [32]byte
	legacyPhase    bool
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

type preparedClassKey struct {
	key    [aes.BlockSize]byte
	digest [sha256.Size]byte
}

const (
	classSearchChunkKeys           = 64
	maximumClassSearchWorkers      = 16
	maximumAuthenticatedClassPairs = uint64(1) << 31
)

func recoveryCandidateCounts(registries []decodedSIDBRegistryCandidate) int {
	records := 0
	for _, registry := range registries {
		records += len(registry.Registry.Records)
	}
	return records
}

func reportRecoveryProgress(progress func(Progress), event Progress) {
	if progress != nil {
		progress(event)
	}
}

func prepareCurrentClassCandidates(registries []decodedSIDBRegistryCandidate) ([]preparedClassCandidate, error) {
	var candidates []preparedClassCandidate
	for registryIndex, registry := range registries {
		for recordIndex, record := range registry.Registry.Records {
			for windowOffset := 0; windowOffset <= classRecordMaximumWindowOffset(record); windowOffset++ {
				valu, err := deriveClassValuWindow(record, windowOffset)
				if err != nil {
					return nil, err
				}
				candidates = append(candidates, preparedClassCandidate{registryIndex: registryIndex, recordIndex: recordIndex, windowOffset: windowOffset, valu: valu})
			}
		}
	}
	return candidates, nil
}

func prepareLegacyClassCandidates(registries []decodedSIDBRegistryCandidate) ([]preparedClassCandidate, error) {
	var candidates []preparedClassCandidate
	for registryIndex, registry := range registries {
		for recordIndex, record := range registry.Registry.Records {
			for maskPhase := range len(currentClassRecordMask) {
				valu, err := deriveLegacyClassValu(record, maskPhase)
				if err != nil {
					return nil, err
				}
				candidates = append(candidates, preparedClassCandidate{registryIndex: registryIndex, recordIndex: recordIndex, windowOffset: maskPhase, valu: valu})
			}
		}
	}
	return candidates, nil
}

func clearPreparedClassCandidates(candidates []preparedClassCandidate) {
	for index := range candidates {
		clear(candidates[index].valu[:])
	}
}

func searchAuthenticatedClassPrecursor(ctx context.Context, sinf []byte, registries []decodedSIDBRegistryCandidate, candidates []preparedClassCandidate, keys []preparedClassKey, legacyPhase bool, progress func(Progress)) (authenticatedClassPrecursor, error) {
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
	total := uint64(len(keys)) * uint64(len(candidates))
	if total == 0 || total > maximumAuthenticatedClassPairs {
		return selected, fmt.Errorf("authenticated SINF class search has unsupported size: %d comparisons", total)
	}
	baseProgress := Progress{Phase: "search-sinf", Total: total, RegistryCandidates: len(registries), RecordCandidates: recoveryCandidateCounts(registries), WindowCandidates: len(candidates)}
	reportRecoveryProgress(progress, baseProgress)

	type searchTask struct{ start, end int }
	type searchResult struct {
		completed uint64
		matches   []authenticatedClassPrecursor
		err       error
	}
	chunks := (len(keys) + classSearchChunkKeys - 1) / classSearchChunkKeys
	workersCount := min(runtime.GOMAXPROCS(0), maximumClassSearchWorkers, chunks)
	searchContext, cancel := context.WithCancel(ctx)
	defer cancel()
	tasks := make(chan searchTask, chunks)
	results := make(chan searchResult, workersCount)
	for start := 0; start < len(keys); start += classSearchChunkKeys {
		end := min(start+classSearchChunkKeys, len(keys))
		tasks <- searchTask{start: start, end: end}
	}
	close(tasks)
	var workers sync.WaitGroup
	workers.Add(workersCount)
	for range workersCount {
		go func() {
			defer workers.Done()
			for task := range tasks {
				if searchContext.Err() != nil {
					return
				}
				result := searchResult{}
				for keyIndex := task.start; keyIndex < task.end; keyIndex++ {
					if searchContext.Err() != nil {
						break
					}
					block, blockErr := aes.NewCipher(keys[keyIndex].key[:])
					if blockErr != nil {
						result.err = blockErr
						break
					}
					for _, candidate := range candidates {
						precursor, accepted, candidateErr := authenticateClassValu(probe, sinf, block, candidate.valu)
						result.completed++
						if candidateErr != nil {
							result.err = candidateErr
							break
						}
						if accepted {
							result.matches = append(result.matches, authenticatedClassPrecursor{classPrecursorCandidate: classPrecursorCandidate{RecordIndex: candidate.recordIndex, WindowOffset: candidate.windowOffset, Record: registries[candidate.registryIndex].Registry.Records[candidate.recordIndex], Precursor: precursor}, RegistryIndex: candidate.registryIndex, classKeyDigest: keys[keyIndex].digest, legacyPhase: legacyPhase})
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
	completed, progressInterval := uint64(0), total/32
	if progressInterval == 0 {
		progressInterval = 1
	}
	nextProgress := progressInterval
	var matches []authenticatedClassPrecursor
	var firstError error
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
			return selected, fmt.Errorf("authenticated SINF accepted multiple class precursors")
		}
		return selected, errNoAuthenticatedClassPrecursor
	}
	selected = matches[0]
	event := baseProgress
	event.Phase, event.Completed = "search-sinf-complete", total
	reportRecoveryProgress(progress, event)
	return selected, nil
}

func selectAuthenticatedClassPrecursor(sinf []byte, registries []decodedSIDBRegistryCandidate, hardware []byte) (authenticatedClassPrecursor, error) {
	return selectAuthenticatedClassPrecursorWithContext(context.Background(), sinf, registries, hardware, nil)
}

func selectAuthenticatedClassPrecursorWithContext(ctx context.Context, sinf []byte, registries []decodedSIDBRegistryCandidate, hardware []byte, progress func(Progress)) (authenticatedClassPrecursor, error) {
	candidates, err := prepareCurrentClassCandidates(registries)
	if err != nil {
		return authenticatedClassPrecursor{}, err
	}
	defer clearPreparedClassCandidates(candidates)
	user, keyID, err := enrollmentSINFSelector(sinf)
	if err != nil {
		return authenticatedClassPrecursor{}, err
	}
	key, err := accountClassKey(hardware, user, keyID)
	if err != nil {
		return authenticatedClassPrecursor{}, err
	}
	defer clear(key[:])
	keys := []preparedClassKey{{key: key, digest: sha256.Sum256(key[:])}}
	return searchAuthenticatedClassPrecursor(ctx, sinf, registries, candidates, keys, false, progress)
}

// selectAuthenticatedFlatClassPrecursor retains bounded upstream compatibility
// for flat registries and deduplicates class keys across derivation paths.
func selectAuthenticatedFlatClassPrecursor(sinf []byte, registries []decodedSIDBRegistryCandidate, hardware []byte, user, keyID uint32) (authenticatedClassPrecursor, error) {
	return selectAuthenticatedFlatClassPrecursorWithContext(context.Background(), sinf, registries, hardware, user, keyID, nil)
}

func selectAuthenticatedFlatClassPrecursorWithContext(ctx context.Context, sinf []byte, registries []decodedSIDBRegistryCandidate, hardware []byte, user, keyID uint32, progress func(Progress)) (authenticatedClassPrecursor, error) {
	candidates, err := prepareCurrentClassCandidates(registries)
	if err != nil {
		return authenticatedClassPrecursor{}, err
	}
	defer clearPreparedClassCandidates(candidates)
	directKey, err := accountClassKey(hardware, user, keyID)
	if err != nil {
		return authenticatedClassPrecursor{}, err
	}
	defer clear(directKey[:])
	keys := make([]preparedClassKey, 0, classWorkspaceModulus+2)
	seen := map[[16]byte]bool{}
	appendKey := func(key [16]byte) {
		if !seen[key] {
			seen[key] = true
			keys = append(keys, preparedClassKey{key: key, digest: sha256.Sum256(key[:])})
		}
	}
	appendKey(directKey)
	appendKey((classEnrollmentParameters{ClassValue: currentClassSecondaryValue}).classKey())
	for seed := range uint32(classWorkspaceModulus) {
		appendKey((classEnrollmentParameters{ClassValue: classWorkspaceOutputFromSeed(uint16(seed))}).classKey())
	}
	defer func() {
		for index := range keys {
			clear(keys[index].key[:])
		}
	}()
	return searchAuthenticatedClassPrecursor(ctx, sinf, registries, candidates, keys, false, progress)
}

func selectAuthenticatedLegacyClassPrecursor(sinf []byte, registries []decodedSIDBRegistryCandidate, hardware []byte, user, keyID uint32) (authenticatedClassPrecursor, error) {
	return selectAuthenticatedLegacyClassPrecursorWithContext(context.Background(), sinf, registries, hardware, user, keyID, nil)
}

func selectAuthenticatedLegacyClassPrecursorWithContext(ctx context.Context, sinf []byte, registries []decodedSIDBRegistryCandidate, hardware []byte, user, keyID uint32, progress func(Progress)) (authenticatedClassPrecursor, error) {
	candidates, err := prepareLegacyClassCandidates(registries)
	if err != nil {
		return authenticatedClassPrecursor{}, err
	}
	defer clearPreparedClassCandidates(candidates)
	digest, err := legacyAccountDigest(hardware, user, keyID)
	if err != nil {
		return authenticatedClassPrecursor{}, err
	}
	key := (classEnrollmentParameters{ClassValue: legacyClassWorkspaceValue(digest)}).classKey()
	clear(digest[:])
	defer clear(key[:])
	keys := []preparedClassKey{{key: key, digest: sha256.Sum256(key[:])}}
	return searchAuthenticatedClassPrecursor(ctx, sinf, registries, candidates, keys, true, progress)
}

// selectAuthenticatedEnrollmentClassPrecursor authenticates each structural
// interpretation before deciding whether an enrollment has a unique match.
func selectAuthenticatedEnrollmentClassPrecursor(sinf []byte, registries []decodedSIDBRegistryCandidate, hardware []byte, user, keyID uint32) (authenticatedClassPrecursor, error) {
	return selectAuthenticatedEnrollmentClassPrecursorWithContext(context.Background(), sinf, registries, hardware, user, keyID, nil)
}

func selectAuthenticatedEnrollmentClassPrecursorWithContext(ctx context.Context, sinf []byte, registries []decodedSIDBRegistryCandidate, hardware []byte, user, keyID uint32, progress func(Progress)) (authenticatedClassPrecursor, error) {
	if len(registries) == 0 {
		return authenticatedClassPrecursor{}, fmt.Errorf("class precursor selection requires one decoded registry")
	}
	var selected authenticatedClassPrecursor
	matches := 0
	for i, registry := range registries {
		var candidate authenticatedClassPrecursor
		var err error
		single := registries[i : i+1]
		if !registry.Registry.Tree && registry.Registry.FooterBytes == decodedSIDBOlderFooterBytes {
			if keyID > 0xff {
				continue
			}
			candidate, err = selectAuthenticatedLegacyClassPrecursorWithContext(ctx, sinf, single, hardware, user, keyID, progress)
		} else if registry.Registry.Tree {
			candidate, err = selectAuthenticatedClassPrecursorWithContext(ctx, sinf, single, hardware, progress)
		} else {
			candidate, err = selectAuthenticatedFlatClassPrecursorWithContext(ctx, sinf, single, hardware, user, keyID, progress)
		}
		if errors.Is(err, errNoAuthenticatedClassPrecursor) {
			continue
		}
		if err != nil {
			return authenticatedClassPrecursor{}, err
		}
		candidate.RegistryIndex = i
		if matches != 0 {
			priorData, currentData := registries[selected.RegistryIndex].Data, registry.Data
			priorOffset := selected.Record.SourceOffset
			if !selected.legacyPhase {
				priorOffset += selected.WindowOffset
			}
			currentOffset := candidate.Record.SourceOffset
			if !candidate.legacyPhase {
				currentOffset += candidate.WindowOffset
			}
			same := len(priorData) != 0 && bytes.Equal(priorData, currentData) &&
				priorOffset == currentOffset &&
				selected.classKeyDigest == candidate.classKeyDigest && selected.Precursor == candidate.Precursor
			if same {
				clear(candidate.Precursor[:])
				continue
			}
			clear(selected.Precursor[:])
			clear(candidate.Precursor[:])
			return authenticatedClassPrecursor{}, fmt.Errorf("authenticated SINF accepted multiple class precursors")
		}
		matches++
		selected = candidate
	}
	if matches == 0 {
		return authenticatedClassPrecursor{}, errNoAuthenticatedClassPrecursor
	}
	return selected, nil
}
