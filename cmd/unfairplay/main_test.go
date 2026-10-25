package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCommandsAreRegistered(t *testing.T) {
	root := newRootCommand()
	for _, name := range []string{"collect", "recover", "decrypt", "verify", "doctor", "profile"} {
		command, _, err := root.Find([]string{name})
		if err != nil || command == root || command.Name() != name {
			t.Fatalf("command %s is not registered: %v", name, err)
		}
	}
}

func TestWritePrivateFileFailsIfDestinationExists(t *testing.T) {
	path := filepath.Join(t.TempDir(), "secret.bin")
	if err := writePrivateFile(path, []byte("first")); err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0o600 {
		t.Fatalf("private file mode = %04o", info.Mode().Perm())
	}
	if err := writePrivateFile(path, []byte("second")); err == nil {
		t.Fatal("private file was overwritten")
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "first" {
		t.Fatalf("private file changed to %q", data)
	}
}
