package recovery

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"sort"
	"strings"
)

const currentSIDBContainerVersion = uint32(0x0007000b)

const maximumDecodedSIDBRegistryCandidates = 128

type decodedSIDBRegistryCandidate struct {
	Registry decodedSIDBRegistry
	Data     []byte
}

type currentSIDBContainerFile struct {
	name string
	data []byte
}

const (
	selectionSourceEnvelopeBytes = 224
	selectionSourcePrefixBytes   = 32
	selectionSourceBodyBytes     = selectionSourceEnvelopeBytes - selectionSourcePrefixBytes
)

func selectionSourceLayouts(source []byte) ([][]byte, error) {
	switch len(source) {
	case selectionSourceBodyBytes:
		return [][]byte{source}, nil
	case selectionSourceEnvelopeBytes:
		// Current 224-byte account output is an envelope: the native selection
		// stage advances past its 32-byte prefix before consuming twelve real
		// state slots. Retain the full layout as a bounded alternative for the
		// two additional slots; SINF authentication decides which is coherent.
		return [][]byte{source[selectionSourcePrefixBytes:], source}, nil
	default:
		return nil, fmt.Errorf("selection source has %d bytes, want %d or %d", len(source), selectionSourceBodyBytes, selectionSourceEnvelopeBytes)
	}
}

func currentSIDBContainerFiles(files map[string][]byte) ([]currentSIDBContainerFile, error) {
	var containers []currentSIDBContainerFile
	seen := make(map[[sha256.Size]byte]bool)
	for name, data := range files {
		if !strings.HasSuffix(name, "SC Info.sidb") || len(data) != selectionSIDBContainerBytes || binary.BigEndian.Uint32(data[:4]) != currentSIDBContainerVersion {
			continue
		}
		digest := sha256.Sum256(data)
		if seen[digest] {
			continue
		}
		seen[digest] = true
		containers = append(containers, currentSIDBContainerFile{name: name, data: data})
	}
	if len(containers) == 0 {
		return nil, fmt.Errorf("typed SIDB selection found no current SIDB container")
	}
	sort.Slice(containers, func(left, right int) bool { return containers[left].name < containers[right].name })
	return containers, nil
}

func clearDecodedSIDBRegistryCandidates(candidates []decodedSIDBRegistryCandidate) {
	for index := range candidates {
		clear(candidates[index].Data)
	}
}

func appendSelectionRegistryCandidates(result []decodedSIDBRegistryCandidate, seen map[[sha256.Size]byte]bool, source []byte, container currentSIDBContainerFile) ([]decodedSIDBRegistryCandidate, error) {
	encoded, err := decodeSelectionSIDBContainer(container.data)
	if err != nil {
		return result, err
	}
	layouts, err := selectionSourceLayouts(source)
	if err != nil {
		return result, err
	}
	for _, layout := range layouts {
		frames, err := enrollmentFrameCandidates(layout)
		if err != nil {
			return result, err
		}
		for _, frame := range frames {
			decoded, err := decryptSelectionSIDBPayload(encoded, frame.Key)
			clear(frame.Key[:])
			if err != nil {
				return result, err
			}
			digest := sha256.Sum256(decoded[:])
			if seen[digest] {
				clear(decoded[:])
				continue
			}
			registry, err := parseDecodedSIDBRegistry(decoded[:])
			if err != nil {
				clear(decoded[:])
				continue
			}
			seen[digest] = true
			data := append([]byte(nil), decoded[:]...)
			clear(decoded[:])
			result = append(result, decodedSIDBRegistryCandidate{Registry: registry, Data: data})
			if len(result) > maximumDecodedSIDBRegistryCandidates {
				clearDecodedSIDBRegistryCandidates(result)
				return nil, fmt.Errorf("typed SIDB selection found more than %d registry candidates", maximumDecodedSIDBRegistryCandidates)
			}
		}
	}
	return result, nil
}

// deriveSelectionRegistryCandidates replaces native SIDB candidate selection.
// It evaluates both supported source layouts and every structurally current
// SIDB container. Personalized SINF authentication later promotes exactly one
// decoded registry/record/window tuple.
func deriveSelectionRegistryCandidates(s *state, runtime *rawSourceRuntime) ([]decodedSIDBRegistryCandidate, error) {
	if s == nil || runtime == nil {
		return nil, fmt.Errorf("typed SIDB selection state is unavailable")
	}
	profile, err := executionProfileForModel(s.model)
	if err != nil {
		return nil, err
	}
	source, err := s.outputBytes(profile.sourceAddress, 224)
	if err != nil {
		return nil, fmt.Errorf("typed SIDB selection source: %w", err)
	}
	defer clear(source)

	containers, err := currentSIDBContainerFiles(runtime.files)
	if err != nil {
		return nil, err
	}

	seenRegistries := make(map[[sha256.Size]byte]bool)
	var result []decodedSIDBRegistryCandidate
	for _, container := range containers {
		result, err = appendSelectionRegistryCandidates(result, seenRegistries, source, container)
		if err != nil {
			clearDecodedSIDBRegistryCandidates(result)
			return nil, err
		}
	}
	return result, nil
}

func deriveCurrentSelectionRegistryCandidates(files map[string][]byte, hardware []byte) ([]decodedSIDBRegistryCandidate, string, int, error) {
	containers, err := currentSIDBContainerFiles(files)
	if err != nil {
		return nil, "", 0, err
	}
	hash := sha256.New()
	seenRegistries := make(map[[sha256.Size]byte]bool)
	var result []decodedSIDBRegistryCandidate
	for _, container := range containers {
		source, iterations, err := deriveCurrentSelectionSource(container.data, hardware)
		if err != nil {
			clearDecodedSIDBRegistryCandidates(result)
			return nil, "", iterations, err
		}
		containerDigest := sha256.Sum256(container.data)
		_, _ = hash.Write(containerDigest[:])
		_, _ = hash.Write(source[:])
		result, err = appendSelectionRegistryCandidates(result, seenRegistries, source[:], container)
		clear(source[:])
		if err != nil {
			clearDecodedSIDBRegistryCandidates(result)
			return nil, "", iterations, err
		}
	}
	if len(result) == 0 {
		return nil, "", outerNaturalIterationCount, fmt.Errorf("typed SIDB selection found no decoded registry candidate")
	}
	return result, fmt.Sprintf("%x", hash.Sum(nil)), outerNaturalIterationCount, nil
}
