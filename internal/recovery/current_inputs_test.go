package recovery

import (
	"os"
	"path/filepath"
	"testing"
)

func writeCurrentInputFixture(t *testing.T, name string, data []byte) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestLoadCurrentRecoveryInputs(t *testing.T) {
	sinf := writeCurrentInputFixture(t, "App.sinf", []byte{1, 2, 3})
	supf := writeCurrentInputFixture(t, "App.v5.supf", []byte{4, 5, 6})
	sidb := writeCurrentInputFixture(t, "SC Info.sidb", []byte{7, 8, 9})
	hardware := writeCurrentInputFixture(t, "iomac.bin", []byte{0, 1, 2, 3, 4, 5})
	inputs, err := loadCurrentRecoveryInputs([]string{
		"App.sinf=" + sinf,
		"App.v5.supf=" + supf,
		"/var/db/fpsd/SC Info/SC Info.sidb=" + sidb,
		"/unused-executable=/path/that/does/not/exist",
	}, hardware)
	if err != nil {
		t.Fatal(err)
	}
	defer inputs.clear()
	if len(inputs.files) != 3 || len(inputs.hardware) != 6 {
		t.Fatalf("loaded files=%d hardware=%d", len(inputs.files), len(inputs.hardware))
	}
}

func TestLoadCurrentRecoveryInputsAcceptsTextHardware(t *testing.T) {
	hardware := writeCurrentInputFixture(t, "iomac.txt", []byte("02:13:24:35:46:57\n"))
	inputs, err := loadCurrentRecoveryInputs(nil, hardware)
	if err != nil {
		t.Fatal(err)
	}
	defer inputs.clear()
	if got := inputs.hardware; len(got) != 6 || got[0] != 0x02 || got[5] != 0x57 {
		t.Fatalf("decoded hardware = %x", got)
	}
}

func TestLoadCurrentRecoveryInputsRejectsInvalidMappings(t *testing.T) {
	hardware := writeCurrentInputFixture(t, "iomac.bin", make([]byte, 6))
	if _, err := loadCurrentRecoveryInputs([]string{"invalid"}, hardware); err == nil {
		t.Fatal("invalid mapping was accepted")
	}
	badHardware := writeCurrentInputFixture(t, "bad.bin", []byte{1, 2, 3})
	if _, err := loadCurrentRecoveryInputs(nil, badHardware); err == nil {
		t.Fatal("invalid hardware was accepted")
	}
}
