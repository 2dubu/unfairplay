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
	pack, err := currentEnrollmentPack(inputs.files, portableCurrentSemanticProfile(), registries)
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
