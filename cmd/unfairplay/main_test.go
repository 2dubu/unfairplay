package main

import (
	"encoding/json"

	"github.com/DrWhax/unfairplay/internal/fairplay"
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

func TestVerificationOutputSchema(t *testing.T) {
	for _, test := range []struct {
		name, shared, schema string
		hashes               map[string]string
	}{
		{"single", "key-a", "unfairplay-model-verification-v1", map[string]string{"App.sinf": "key-a"}},
		{"shared", "key-a", "unfairplay-model-verification-v1", map[string]string{"App.sinf": "key-a", "Widget.sinf": "key-a"}},
		{"distinct", "", "unfairplay-model-verification-v2", map[string]string{"App.sinf": "key-a", "Widget.sinf": "key-b"}},
	} {
		t.Run(test.name, func(t *testing.T) {
			data, err := json.Marshal(verificationSummary(fairplay.DecryptResult{
				LogicalKeySHA256: test.shared, LogicalKeySHA256BySINF: test.hashes,
			}))
			if err != nil {
				t.Fatal(err)
			}
			var got map[string]json.RawMessage
			if err := json.Unmarshal(data, &got); err != nil {
				t.Fatal(err)
			}
			var schema string
			if err := json.Unmarshal(got["schema"], &schema); err != nil || schema != test.schema {
				t.Fatalf("schema=%q err=%v", schema, err)
			}
			value, hasScalar := got["logical_key_sha256"]
			if test.shared == "" {
				if hasScalar {
					t.Fatal("distinct-key verification emitted legacy scalar")
				}
			} else {
				var hash string
				if !hasScalar || json.Unmarshal(value, &hash) != nil || hash != test.shared {
					t.Fatal("legacy scalar contract changed")
				}
			}
		})
	}
}

func TestCollectExposesArchivedExecutablePath(t *testing.T) {
	root := newRootCommand()
	command, _, err := root.Find([]string{"collect"})
	if err != nil {
		t.Fatal(err)
	}
	if command.Flags().Lookup("installed-executable-path") == nil {
		t.Fatal("collect has no archived executable path override")
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
