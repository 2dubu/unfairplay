package fairplay

import (
	"archive/zip"
	"fmt"
	"io"
	"strings"

	"howett.net/plist"
)

// Account data that an output IPA does not need once its images are
// decrypted. SC_Info holds the personalized SINF/SUPF licenses, which
// identify the Store account. iTunesMetadata.plist carries the Apple ID
// written by downloaders and the purchaser record from the App Store.
var storeMetadataAccountKeys = []string{
	"apple-id", "appleId", "AppleID", "userName", "UserName",
	"purchaserID", "DSPersonID", "dsPersonId", "AltDSID", "altDsid", "FamilyID", "familyId",
}

const storeMetadataMember = "iTunesMetadata.plist"

// isLicenseMember reports whether an archive member is inside any SC_Info
// directory: the app's, an extension's, a watch app's, or a framework's.
func isLicenseMember(name string) bool {
	for _, part := range strings.Split(strings.TrimSuffix(name, "/"), "/") {
		if strings.EqualFold(part, "SC_Info") {
			return true
		}
	}
	return false
}

// scrubStoreMetadata returns iTunesMetadata.plist without account fields,
// in its original plist format, and how many fields were removed.
func scrubStoreMetadata(file *zip.File) ([]byte, int, error) {
	reader, err := file.Open()
	if err != nil {
		return nil, 0, err
	}
	data, err := io.ReadAll(io.LimitReader(reader, 16<<20))
	reader.Close()
	if err != nil {
		return nil, 0, err
	}
	var metadata map[string]any
	format, err := plist.Unmarshal(data, &metadata)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: %w", storeMetadataMember, err)
	}
	removed := 0
	for _, key := range storeMetadataAccountKeys {
		if _, ok := metadata[key]; ok {
			delete(metadata, key)
			removed++
		}
	}
	if download, ok := metadata["com.apple.iTunesStore.downloadInfo"].(map[string]any); ok {
		if _, ok := download["accountInfo"]; ok {
			delete(download, "accountInfo")
			removed++
		}
	}
	data, err = plist.Marshal(metadata, format)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: %w", storeMetadataMember, err)
	}
	return data, removed, nil
}
