package recovery

import (
	"fmt"
	"os"
	"sync"
)

// Options describes the private runtime inputs used to derive one
// selector-bound enrollment model. FileMappings use SUFFIX=FILE syntax.
type Options struct {
	FileMappings []string
	HardwarePath string
}

// Result is the output of one bounded enrollment recovery.
type Result struct {
	EnrollmentModel           []byte
	StateSHA256               string
	Iterations                int
	DecodedRegistryCandidates int
}

var recoveryMu sync.Mutex

func readPrivateInput(path string) ([]byte, error) {
	info, err := os.Lstat(path)
	if err != nil {
		return nil, err
	}
	if !info.Mode().IsRegular() || info.Mode().Perm()&0o077 != 0 {
		return nil, fmt.Errorf("recovery input must be a mode-0600 regular file: %s", path)
	}
	return os.ReadFile(path)
}

// Recover derives a selector-bound enrollment model from structurally parsed
// SIDB, SINF, SUPF, and hardware inputs. Native code and the captured
// executable are not part of the production recovery path.
func Recover(options Options) (Result, error) {
	recoveryMu.Lock()
	defer recoveryMu.Unlock()

	if len(options.FileMappings) == 0 {
		return Result{}, fmt.Errorf("recovery requires captured FairPlay files")
	}
	if options.HardwarePath == "" {
		return Result{}, fmt.Errorf("recovery requires a hardware identity file")
	}
	inputs, err := loadCurrentRecoveryInputs(options.FileMappings, options.HardwarePath)
	if err != nil {
		return Result{}, err
	}
	defer inputs.clear()
	if _, err := validateCurrentSUPFFiles(inputs.files); err != nil {
		return Result{}, err
	}
	registries, digest, iterations, err := deriveCurrentSelectionRegistryCandidates(inputs.files, inputs.hardware)
	if err != nil {
		return Result{}, err
	}
	pack, err := currentEnrollmentPack(inputs.files, portableCurrentSemanticProfile(), rawSourceBinaryHash, registries)
	if err != nil {
		clearDecodedSIDBRegistryCandidates(registries)
		return Result{}, err
	}
	clearDecodedSIDBRegistryCandidates(registries)
	return Result{
		EnrollmentModel: pack, StateSHA256: digest, Iterations: iterations,
		DecodedRegistryCandidates: len(registries),
	}, nil
}

func clearPrivateRecoveryState(value *model) {
	if lastRawSourceState != nil {
		for _, page := range lastRawSourceState.writable {
			clear(page)
		}
		for _, page := range lastRawSourceState.mapped {
			clear(page)
		}
		clear(lastRawSourceState.writable)
		clear(lastRawSourceState.mapped)
		clear(lastRawSourceState.initialized)
		for index := range lastRawSourceState.accountDigests {
			clear(lastRawSourceState.accountDigests[index][:])
		}
		clear(lastRawSourceState.accountDigests)
		lastRawSourceState = nil
	}
	if activeRawSourceRuntime != nil {
		for name, data := range activeRawSourceRuntime.files {
			clear(data)
			delete(activeRawSourceRuntime.files, name)
		}
		clear(activeRawSourceRuntime.mac)
		for key, item := range activeRawSourceRuntime.fakeCFValues {
			clear(item.data)
			delete(activeRawSourceRuntime.fakeCFValues, key)
		}
		activeRawSourceRuntime = nil
	}
	if value != nil {
		for address, page := range value.pages {
			clear(page)
			delete(value.pages, address)
		}
		clear(value.registers[:])
	}
}
