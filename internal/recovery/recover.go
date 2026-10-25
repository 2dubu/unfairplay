package recovery

import (
	"fmt"
	"os"
	"sync"
)

// Options describes the private runtime inputs used to derive one
// selector-bound enrollment model. FileMappings use SUFFIX=FILE syntax.
type Options struct {
	FileMappings   []string
	HardwarePath   string
	ExecutablePath string
}

// Result is the output of one bounded enrollment recovery.
type Result struct {
	EnrollmentModel []byte
	StateSHA256     string
	Iterations      int
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

// Recover derives a selector-bound enrollment model from a normalized native
// Mac capture. The embedded model contains build-global state only; private
// SIDB, SINF, SUPF, executable, and hardware bytes are supplied at runtime.
func Recover(options Options) (Result, error) {
	recoveryMu.Lock()
	defer recoveryMu.Unlock()

	if len(options.FileMappings) == 0 {
		return Result{}, fmt.Errorf("recovery requires captured FairPlay files")
	}
	if options.HardwarePath == "" {
		return Result{}, fmt.Errorf("recovery requires a hardware identity file")
	}
	if options.ExecutablePath == "" {
		return Result{}, fmt.Errorf("recovery requires an installed executable path")
	}
	value, err := loadModelBytes(embeddedModelBytes())
	if err != nil {
		return Result{}, fmt.Errorf("load embedded recovery model: %w", err)
	}
	defer clearPrivateRecoveryState(value)
	if err := configureRawSourceRuntime(value, "enrollment-sinf", options.FileMappings,
		options.HardwarePath, options.ExecutablePath, 0); err != nil {
		return Result{}, err
	}
	digest, iterations, err := runRawSource(value, "enrollment-sinf", false, 0)
	if err != nil {
		return Result{}, err
	}
	pack, err := currentRawSourceEnrollmentPack()
	if err != nil {
		return Result{}, err
	}
	return Result{EnrollmentModel: pack, StateSHA256: digest, Iterations: iterations}, nil
}

func clearPrivateRecoveryState(value *model) {
	if lastRawSourceState != nil {
		for _, page := range lastRawSourceState.writable {
			clear(page)
		}
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
