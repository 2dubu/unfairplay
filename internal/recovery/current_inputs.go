package recovery

import (
	"encoding/hex"
	"fmt"
	"strings"
)

type currentRecoveryInputs struct {
	files    map[string][]byte
	hardware []byte
}

func isCurrentRecoveryFile(name string) bool {
	return strings.HasSuffix(name, ".sinf") || strings.HasSuffix(name, ".supf") || strings.HasSuffix(name, ".supp") || strings.HasSuffix(name, "SC Info.sidb")
}

func loadCurrentRecoveryInputs(fileSpecs []string, hardwarePath string) (*currentRecoveryInputs, error) {
	result := &currentRecoveryInputs{files: make(map[string][]byte)}
	for _, specification := range fileSpecs {
		suffix, path, ok := strings.Cut(specification, "=")
		if !ok || suffix == "" || path == "" {
			result.clear()
			return nil, fmt.Errorf("invalid recovery file mapping %q", specification)
		}
		if !isCurrentRecoveryFile(suffix) {
			continue
		}
		if _, exists := result.files[suffix]; exists {
			result.clear()
			return nil, fmt.Errorf("duplicate recovery file mapping %q", suffix)
		}
		data, err := readPrivateInput(path)
		if err != nil {
			result.clear()
			return nil, fmt.Errorf("read recovery file %q: %w", path, err)
		}
		result.files[suffix] = data
	}
	if hardwarePath == "" {
		result.clear()
		return nil, fmt.Errorf("recovery requires a hardware identity file")
	}
	hardware, err := readPrivateInput(hardwarePath)
	if err != nil {
		result.clear()
		return nil, err
	}
	if len(hardware) != 6 {
		raw := hardware
		encoded := strings.NewReplacer(":", "", "-", "").Replace(strings.TrimSpace(string(hardware)))
		hardware, err = hex.DecodeString(encoded)
		clear(raw)
	}
	if err != nil || len(hardware) != 6 {
		clear(hardware)
		result.clear()
		return nil, fmt.Errorf("recovery hardware address must contain six bytes")
	}
	result.hardware = hardware
	return result, nil
}

func (input *currentRecoveryInputs) clear() {
	if input == nil {
		return
	}
	for name, data := range input.files {
		clear(data)
		delete(input.files, name)
	}
	clear(input.hardware)
	input.hardware = nil
}
