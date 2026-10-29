package fairplay

import (
	"archive/zip"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"howett.net/plist"
)

const syntheticStoreMetadata = `<?xml version="1.0" encoding="UTF-8"?>
<plist version="1.0"><dict>
<key>apple-id</key><string>buyer@example.com</string>
<key>userName</key><string>buyer@example.com</string>
<key>itemName</key><string>Test</string>
<key>bundleShortVersionString</key><string>1.0</string>
<key>com.apple.iTunesStore.downloadInfo</key><dict>
<key>accountInfo</key><dict><key>AppleID</key><string>buyer@example.com</string><key>DSPersonID</key><integer>42</integer></dict>
<key>purchaseDate</key><string>2026-01-01T00:00:00Z</string>
</dict>
</dict></plist>`

func writeOrderedArchive(t *testing.T, name string, members []string, data map[string][]byte) {
	t.Helper()
	file, err := os.Create(name)
	if err != nil {
		t.Fatal(err)
	}
	archive := zip.NewWriter(file)
	for _, member := range members {
		writer, err := archive.Create(member)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := writer.Write(data[member]); err != nil {
			t.Fatal(err)
		}
	}
	if err := archive.Close(); err != nil {
		t.Fatal(err)
	}
	if err := file.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestDecryptStripAccountData(t *testing.T) {
	directory := t.TempDir()
	mainPath := "Payload/Test.app/Test"
	widgetPath := "Payload/Test.app/PlugIns/Widget.appex/Test"
	mainFixture := syntheticCurrentSINFFixture(t, 42, 7, 88, [16]byte{1, 2, 3})
	widgetFixture := syntheticCurrentSINFFixture(t, 42, 8, 88, [16]byte{9, 8, 7})
	mainBlob, mainSUPF, mainPlain := multiSINFImage(t, mainFixture.logical, 0xd503201f)
	widgetBlob, widgetSUPF, widgetPlain := multiSINFImage(t, widgetFixture.logical, 0xd65f03c0)
	members := []string{
		storeMetadataMember,
		mainPath, adjacentSINF(mainPath), "Payload/Test.app/SC_Info/Test.supf",
		"Payload/Test.app/keep.txt",
		widgetPath, adjacentSINF(widgetPath), "Payload/Test.app/PlugIns/Widget.appex/SC_Info/Test.supf",
	}
	data := map[string][]byte{
		storeMetadataMember: []byte(syntheticStoreMetadata),
		mainPath:            mainBlob, adjacentSINF(mainPath): mainFixture.sinf,
		"Payload/Test.app/SC_Info/Test.supf": mainSUPF,
		"Payload/Test.app/keep.txt":          []byte("unchanged"),
		widgetPath:                           widgetBlob, adjacentSINF(widgetPath): widgetFixture.sinf,
		"Payload/Test.app/PlugIns/Widget.appex/SC_Info/Test.supf": widgetSUPF,
	}
	input := filepath.Join(directory, "App.ipa")
	writeOrderedArchive(t, input, members, data)
	models := filepath.Join(directory, "models")
	if err := os.Mkdir(models, 0o700); err != nil {
		t.Fatal(err)
	}
	for name, selector := range map[string]SINFSelector{
		"main.bin": mainFixture.selector, "widget.bin": widgetFixture.selector,
	} {
		if err := os.WriteFile(filepath.Join(models, name), selectorBoundModelPack(t, selector), 0o600); err != nil {
			t.Fatal(err)
		}
	}

	output, reportPath := filepath.Join(directory, "out.ipa"), filepath.Join(directory, "out.json")
	result, err := DecryptIPA(DecryptOptions{
		InputPath: input, OutputPath: output, ReportPath: reportPath, ModelPath: models, StripAccountData: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.EncryptedImages != 2 || result.RemovedLicenseFiles != 4 || result.RemovedMetadataFields != 3 {
		t.Fatalf("result %+v", result)
	}

	archive, err := zip.OpenReader(output)
	if err != nil {
		t.Fatal(err)
	}
	defer archive.Close()
	var names []string
	for _, file := range archive.File {
		names = append(names, file.Name)
		if isLicenseMember(file.Name) {
			t.Fatalf("license member kept: %s", file.Name)
		}
		got, err := zipBytes(file)
		if err != nil {
			t.Fatal(err)
		}
		switch file.Name {
		case storeMetadataMember:
			if bytes.Contains(got, []byte("buyer@example.com")) {
				t.Fatalf("Apple ID kept in metadata:\n%s", got)
			}
			var metadata map[string]any
			if _, err := plist.Unmarshal(got, &metadata); err != nil {
				t.Fatal(err)
			}
			download := metadata["com.apple.iTunesStore.downloadInfo"].(map[string]any)
			if metadata["itemName"] != "Test" || metadata["bundleShortVersionString"] != "1.0" ||
				download["accountInfo"] != nil || download["purchaseDate"] == nil {
				t.Fatalf("metadata %v", metadata)
			}
		case mainPath, widgetPath:
			want := append([]byte(nil), data[file.Name]...)
			binary.LittleEndian.PutUint32(want[48:], 0)
			plain := mainPlain
			if file.Name == widgetPath {
				plain = widgetPlain
			}
			copy(want[pageBytes:2*pageBytes], plain)
			if !bytes.Equal(got, want) {
				t.Fatalf("wrong bytes for %s", file.Name)
			}
		default:
			if !bytes.Equal(got, data[file.Name]) {
				t.Fatalf("wrong bytes for %s", file.Name)
			}
		}
	}
	if strings.Join(names, ",") != strings.Join([]string{storeMetadataMember, mainPath, "Payload/Test.app/keep.txt", widgetPath}, ",") {
		t.Fatalf("members %v", names)
	}

	reportData, err := os.ReadFile(reportPath)
	if err != nil {
		t.Fatal(err)
	}
	var report decryptionReport
	if err := json.Unmarshal(reportData, &report); err != nil {
		t.Fatal(err)
	}
	if !report.Complete || !report.AccountDataStripped || report.RemovedLicenseFiles != 4 || report.RemovedMetadataFields != 3 {
		t.Fatalf("report %+v", report)
	}

	// Without the option the output keeps every member, as before.
	plainOutput := filepath.Join(directory, "keep.ipa")
	if _, err := DecryptIPA(DecryptOptions{
		InputPath: input, OutputPath: plainOutput, ReportPath: filepath.Join(directory, "keep.json"), ModelPath: models,
	}); err != nil {
		t.Fatal(err)
	}
	kept, err := zip.OpenReader(plainOutput)
	if err != nil {
		t.Fatal(err)
	}
	defer kept.Close()
	if len(kept.File) != len(members) {
		t.Fatalf("default output has %d members, want %d", len(kept.File), len(members))
	}
}

func TestIsLicenseMember(t *testing.T) {
	for name, want := range map[string]bool{
		"Payload/A.app/SC_Info/":                         true,
		"Payload/A.app/SC_Info/A.sinf":                   true,
		"Payload/A.app/PlugIns/B.appex/sc_info/B.supf":   true,
		"Payload/A.app/Watch/W.app/SC_Info/W.sinf":       true,
		"Payload/A.app/SC_Information.txt":               false,
		"Payload/A.app/Frameworks/SC_InfoKit.framework/": false,
	} {
		if got := isLicenseMember(name); got != want {
			t.Errorf("isLicenseMember(%q) = %v", name, got)
		}
	}
}
