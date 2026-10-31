package recovery

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
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

type selectionRegistryDiagnostics struct {
	Containers      int
	ContainerBytes  int
	CodecBytes      int
	PayloadBytes    int
	SourceBytes     int
	OuterIterations int
	SourceLayouts   int
	KeyCandidates   int
	DecodedBodies   int
	TreeCandidates  int
}

func (diagnostics selectionRegistryDiagnostics) String() string {
	return fmt.Sprintf(
		"containers=%d container_bytes=%d codec_bytes=%d payload_bytes=%d source_bytes=%d iterations=%d source_layouts=%d key_candidates=%d decoded_bodies=%d tree_candidates=%d",
		diagnostics.Containers, diagnostics.ContainerBytes, diagnostics.CodecBytes,
		diagnostics.PayloadBytes, diagnostics.SourceBytes, diagnostics.OuterIterations,
		diagnostics.SourceLayouts, diagnostics.KeyCandidates, diagnostics.DecodedBodies,
		diagnostics.TreeCandidates,
	)
}

const (
	selectionSourceEnvelopeBytes = 224
	selectionSourcePrefixBytes   = 32
	selectionSourceBodyBytes     = selectionSourceEnvelopeBytes - selectionSourcePrefixBytes
)

// selectionSourceLayouts accepts the mapped native source extent. The first
// 32 bytes precede the modulus-sized state input projection.
func selectionSourceLayouts(source []byte) ([][]byte, error) {
	if len(source) < selectionSourcePrefixBytes+16 || len(source) > selectionSourcePrefixBytes+16*maximumMappedSelectionModulus || (len(source)-selectionSourcePrefixBytes)%16 != 0 {
		return nil, fmt.Errorf("unsupported selection source size %d", len(source))
	}
	return [][]byte{source[selectionSourcePrefixBytes:]}, nil
}

func currentSIDBContainerFiles(files map[string][]byte) ([]currentSIDBContainerFile, error) {
	var containers []currentSIDBContainerFile
	for name, data := range files {
		if !strings.HasSuffix(name, "SC Info.sidb") {
			continue
		}
		if len(data) < 4 || binary.BigEndian.Uint32(data[:4]) != currentSIDBContainerVersion {
			continue
		}
		if len(containers) != 0 {
			if !bytes.Equal(containers[0].data, data) {
				return nil, fmt.Errorf("typed SIDB selection requires one distinct current SIDB container")
			}
			if name < containers[0].name {
				containers[0].name = name
			}
			continue
		}
		containers = append(containers, currentSIDBContainerFile{name: name, data: data})
	}
	if len(containers) == 0 {
		return nil, fmt.Errorf("typed SIDB selection found no current SIDB container")
	}
	return containers, nil
}

func clearDecodedSIDBRegistryCandidates(candidates []decodedSIDBRegistryCandidate) {
	for index := range candidates {
		clear(candidates[index].Data)
	}
}

func appendSelectionRegistryCandidates(result []decodedSIDBRegistryCandidate, seen map[[sha256.Size]byte]bool, source []byte, container currentSIDBContainerFile, layout currentSIDBLayout, diagnostics *selectionRegistryDiagnostics) ([]decodedSIDBRegistryCandidate, error) {
	encoded, err := decodeSelectionSIDBContainerWithLayout(container.data, layout)
	if err != nil {
		return result, err
	}
	defer clear(encoded)
	layouts, err := selectionSourceLayouts(source)
	if err != nil {
		return result, err
	}
	if len(layouts[0]) != layout.SelectorModulus*16 || layout.SelectedSlot < 0 || layout.SelectedSlot >= layout.SelectorModulus {
		return result, fmt.Errorf("invalid current SIDB selected state slot")
	}
	key, err := deriveSelectedEnrollmentKey(layouts[0], layout.SelectedSlot)
	if err != nil {
		return result, err
	}
	if diagnostics != nil {
		diagnostics.SourceLayouts++
		diagnostics.KeyCandidates++
	}
	decoded, err := decryptSelectionSIDBPayloadDynamic(encoded, key)
	clear(key[:])
	if err != nil {
		return result, err
	}
	defer clear(decoded)
	if diagnostics != nil {
		diagnostics.DecodedBodies++
	}
	digest := sha256.Sum256(decoded)
	if seen[digest] {
		return result, nil
	}
	registry, err := parseDecodedSIDBRegistry(decoded, layout.RootCursor)
	if err != nil {
		return result, fmt.Errorf("selected SIDB body has invalid tree geometry: %w", err)
	}
	if diagnostics != nil {
		diagnostics.TreeCandidates++
	}
	seen[digest] = true
	result = append(result, decodedSIDBRegistryCandidate{Registry: registry, Data: append([]byte(nil), decoded...)})
	if len(result) > maximumDecodedSIDBRegistryCandidates {
		clearDecodedSIDBRegistryCandidates(result)
		return nil, fmt.Errorf("typed SIDB selection found more than %d registry candidates", maximumDecodedSIDBRegistryCandidates)
	}
	return result, nil
}

func deriveCurrentSelectionRegistryCandidates(files map[string][]byte, hardware []byte) ([]decodedSIDBRegistryCandidate, string, int, error) {
	result, digest, iterations, _, err := deriveCurrentSelectionRegistryCandidatesWithEvidence(files, hardware)
	return result, digest, iterations, err
}

func deriveCurrentSelectionRegistryCandidatesWithEvidence(files map[string][]byte, hardware []byte) ([]decodedSIDBRegistryCandidate, string, int, LayoutEvidence, error) {
	var evidence LayoutEvidence
	containers, err := currentSIDBContainerFiles(files)
	if err != nil {
		return nil, "", 0, evidence, err
	}
	hash := sha256.New()
	seenRegistries := make(map[[sha256.Size]byte]bool)
	diagnostics := selectionRegistryDiagnostics{Containers: len(containers)}
	var result []decodedSIDBRegistryCandidate
	reportedIterations := 0
	for _, container := range containers {
		source, layout, iterations, err := deriveNativeSelectionSource(container.data, hardware)
		if err != nil {
			clearDecodedSIDBRegistryCandidates(result)
			return nil, "", iterations, evidence, err
		}
		if err := validateSIDBContainerAuthentication(container.data, source); err != nil {
			clear(source)
			clearDecodedSIDBRegistryCandidates(result)
			return nil, "", iterations, evidence, err
		}
		evidence.ContainerBytes = len(container.data)
		evidence.CodecBytes = layout.CodecBytes
		evidence.PayloadBytes = layout.PayloadBytes
		evidence.SourceBytes = layout.SourceBytes
		evidence.SelectorModulus = layout.SelectorModulus
		evidence.SelectedSlot = layout.SelectedSlot
		evidence.RootCursor = layout.RootCursor
		diagnostics.ContainerBytes = len(container.data)
		diagnostics.CodecBytes = layout.CodecBytes
		diagnostics.PayloadBytes = layout.PayloadBytes
		diagnostics.SourceBytes = layout.SourceBytes
		diagnostics.OuterIterations = iterations
		if reportedIterations != 0 && iterations != reportedIterations {
			clear(source[:])
			clearDecodedSIDBRegistryCandidates(result)
			return nil, "", 0, evidence, fmt.Errorf("typed SIDB selection found mixed container layouts")
		}
		reportedIterations = iterations
		containerDigest := sha256.Sum256(container.data)
		_, _ = hash.Write(containerDigest[:])
		_, _ = hash.Write(source[:])
		result, err = appendSelectionRegistryCandidates(result, seenRegistries, source[:], container, layout, &diagnostics)
		clear(source[:])
		if err != nil {
			clearDecodedSIDBRegistryCandidates(result)
			return nil, "", iterations, evidence, err
		}
	}
	if len(result) == 0 {
		return nil, "", reportedIterations, evidence, fmt.Errorf("typed SIDB selection found no decoded registry candidate (%s)", diagnostics.String())
	}
	if len(result) == 1 {
		nodes, err := parseDecodedSIDBTree(result[0].Data, evidence.RootCursor)
		if err != nil {
			clearDecodedSIDBRegistryCandidates(result)
			return nil, "", reportedIterations, evidence, err
		}
		evidence.RootLength = nodes[0].Length
		evidence.RootChildCount = nodes[0].ChildCount
		evidence.TreeNodes = len(nodes)
		evidence.RegistryRecords = len(result[0].Registry.Records)
	}
	return result, fmt.Sprintf("%x", hash.Sum(nil)), reportedIterations, evidence, nil
}
