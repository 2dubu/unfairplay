package recovery

import (
	"crypto/sha256"
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
	var rejections []string
	seen := make(map[[sha256.Size]byte]bool)
	for name, data := range files {
		if !strings.HasSuffix(name, "SC Info.sidb") {
			continue
		}
		if _, err := currentSIDBLayoutForContainer(data); err != nil {
			rejections = append(rejections, fmt.Sprintf("%s (%d bytes): %v", name, len(data), err))
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
		sort.Strings(rejections)
		if len(rejections) != 0 {
			return nil, fmt.Errorf("typed SIDB selection found no current SIDB container: %s", strings.Join(rejections, "; "))
		}
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
		candidates, err := enrollmentKeyCandidates(layout)
		if err != nil {
			return result, err
		}
		for _, candidate := range candidates {
			decoded, err := decryptSelectionSIDBPayload(encoded, candidate.Key)
			clear(candidate.Key[:])
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

func deriveCurrentSelectionRegistryCandidates(files map[string][]byte, hardware []byte) ([]decodedSIDBRegistryCandidate, string, int, error) {
	containers, err := currentSIDBContainerFiles(files)
	if err != nil {
		return nil, "", 0, err
	}
	hash := sha256.New()
	seenRegistries := make(map[[sha256.Size]byte]bool)
	var result []decodedSIDBRegistryCandidate
	reportedIterations := 0
	for _, container := range containers {
		source, iterations, err := deriveCurrentSelectionSource(container.data, hardware)
		if err != nil {
			clearDecodedSIDBRegistryCandidates(result)
			return nil, "", iterations, err
		}
		if reportedIterations != 0 && iterations != reportedIterations {
			clear(source[:])
			clearDecodedSIDBRegistryCandidates(result)
			return nil, "", 0, fmt.Errorf("typed SIDB selection found mixed container layouts")
		}
		reportedIterations = iterations
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
		return nil, "", reportedIterations, fmt.Errorf("typed SIDB selection found no decoded registry candidate")
	}
	return result, fmt.Sprintf("%x", hash.Sum(nil)), reportedIterations, nil
}
