package fairplay

import (
	"archive/zip"
	"bytes"
	"fmt"
	"path"
	"sort"
	"strings"

	"howett.net/plist"
)

// buildSINFAssociations maps each explicitly declared SINF target to its one
// source. An undeclared image may use its adjacent SINF independently.
func buildSINFAssociations(files map[string]*zip.File) (map[string]string, error) {
	const suffix = "/SC_Info/Manifest.plist"
	var names []string
	for name := range files {
		if strings.HasSuffix(name, suffix) {
			names = append(names, name)
		}
	}
	sort.Strings(names)
	result := make(map[string]string)
	for _, name := range names {
		root := strings.TrimSuffix(name, suffix)
		if root == "" || root == ".." || strings.HasPrefix(root, "../") || path.Clean(root) != root || strings.HasPrefix(root, "/") || strings.ContainsAny(root, "\\\x00") {
			return nil, fmt.Errorf("invalid manifest member path %q", name)
		}
		data, err := zipBytes(files[name])
		if err != nil {
			return nil, fmt.Errorf("read manifest %q: %w", name, err)
		}
		var manifest macOSSINFManifest
		_, err = plist.Unmarshal(data, &manifest)
		clear(data)
		if err != nil {
			return nil, fmt.Errorf("parse manifest %q: %w", name, err)
		}
		if len(manifest.SINFPaths) != 1 {
			return nil, fmt.Errorf("manifest %q requires exactly one source SINF", name)
		}
		member := func(relative string) (string, error) {
			if strings.ContainsRune(relative, '\x00') {
				return "", fmt.Errorf("SINF path contains NUL")
			}
			return macOSSINFMember(root, relative)
		}
		source, err := member(manifest.SINFPaths[0])
		if err != nil {
			return nil, fmt.Errorf("manifest %q: %w", name, err)
		}
		sourceFile, exists := files[source]
		if !exists || sourceFile == nil {
			return nil, fmt.Errorf("manifest %q source SINF %q is missing", name, source)
		}
		sourceBytes, err := zipBytes(sourceFile)
		if err != nil {
			return nil, fmt.Errorf("read source SINF %q: %w", source, err)
		}
		for _, relative := range manifest.SINFReplicationPaths {
			target, err := member(relative)
			if err != nil {
				clear(sourceBytes)
				return nil, fmt.Errorf("manifest %q: %w", name, err)
			}
			if _, exists := result[target]; exists {
				clear(sourceBytes)
				return nil, fmt.Errorf("multiple manifest mappings for SINF %q", target)
			}
			if targetFile, exists := files[target]; exists {
				if targetFile == nil {
					clear(sourceBytes)
					return nil, fmt.Errorf("invalid target SINF %q", target)
				}
				content, readErr := zipBytes(targetFile)
				if readErr != nil {
					clear(sourceBytes)
					return nil, fmt.Errorf("read target SINF %q: %w", target, readErr)
				}
				equal := bytes.Equal(content, sourceBytes)
				clear(content)
				if !equal {
					clear(sourceBytes)
					return nil, fmt.Errorf("target SINF %q differs from source %q", target, source)
				}
			}
			result[target] = source
		}
		clear(sourceBytes)
	}
	return result, nil
}
