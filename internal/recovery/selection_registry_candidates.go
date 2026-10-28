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

	type containerFile struct {
		name string
		data []byte
	}
	var containers []containerFile
	seenContainers := make(map[[sha256.Size]byte]bool)
	for name, data := range runtime.files {
		if !strings.HasSuffix(name, "SC Info.sidb") || len(data) != selectionSIDBContainerBytes || binary.BigEndian.Uint32(data[:4]) != currentSIDBContainerVersion {
			continue
		}
		digest := sha256.Sum256(data)
		if seenContainers[digest] {
			continue
		}
		seenContainers[digest] = true
		containers = append(containers, containerFile{name: name, data: data})
	}
	if len(containers) == 0 {
		return nil, fmt.Errorf("typed SIDB selection found no current SIDB container")
	}
	sort.Slice(containers, func(left, right int) bool { return containers[left].name < containers[right].name })

	seenRegistries := make(map[[sha256.Size]byte]bool)
	var result []decodedSIDBRegistryCandidate
	for _, container := range containers {
		encoded, err := decodeSelectionSIDBContainer(container.data)
		if err != nil {
			return nil, err
		}
		layouts, err := selectionSourceLayouts(source)
		if err != nil {
			return nil, err
		}
		for _, layout := range layouts {
			frames, err := enrollmentFrameCandidates(layout)
			if err != nil {
				return nil, err
			}
			for _, frame := range frames {
				decoded, err := decryptSelectionSIDBPayload(encoded, frame.Key)
				clear(frame.Key[:])
				if err != nil {
					return nil, err
				}
				digest := sha256.Sum256(decoded[:])
				if seenRegistries[digest] {
					clear(decoded[:])
					continue
				}
				registry, err := parseDecodedSIDBRegistry(decoded[:])
				if err != nil {
					clear(decoded[:])
					continue
				}
				seenRegistries[digest] = true
				data := append([]byte(nil), decoded[:]...)
				clear(decoded[:])
				result = append(result, decodedSIDBRegistryCandidate{Registry: registry, Data: data})
				if len(result) > maximumDecodedSIDBRegistryCandidates {
					for index := range result {
						clear(result[index].Data)
					}
					return nil, fmt.Errorf("typed SIDB selection found more than %d registry candidates", maximumDecodedSIDBRegistryCandidates)
				}
			}
		}
	}
	return result, nil
}
