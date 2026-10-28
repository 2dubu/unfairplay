package recovery

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"sort"
)

const maximumDecodedSIDBRegistryCandidates = 128

type decodedSIDBRegistryCandidate struct {
	Registry decodedSIDBRegistry
	Data     []byte
}

// discoverDecodedSIDBRegistries finds zero-padded decoded registry objects in
// runtime allocations. It uses only the typed record/footer geometry and the
// allocation's page boundary; native object addresses are not identities and
// are deliberately not returned.
func discoverDecodedSIDBRegistries(s *state) ([]decodedSIDBRegistryCandidate, error) {
	if s == nil {
		return nil, fmt.Errorf("decoded SIDB discovery state is unavailable")
	}
	pages := make([]uint64, 0, len(s.mapped))
	for address := range s.mapped {
		pages = append(pages, address)
	}
	sort.Slice(pages, func(left, right int) bool { return pages[left] < pages[right] })

	seen := make(map[[sha256.Size]byte]bool)
	var result []decodedSIDBRegistryCandidate
	for _, address := range pages {
		page := s.mapped[address]
		if current := s.writable[address]; current != nil {
			page = current
		}
		if len(page) != int(pageSize) {
			return nil, fmt.Errorf("decoded SIDB discovery page has invalid size %d", len(page))
		}
		for count := 1; count <= maximumDecodedSIDBRecords; count++ {
			size := count*decodedSIDBRecordBytes + decodedSIDBFooterBytes
			if size > len(page) || !allZero(page[size:]) {
				continue
			}
			valid := !allZero(page[count*decodedSIDBRecordBytes : size])
			for index := 0; valid && index < count; index++ {
				start := index * decodedSIDBRecordBytes
				valid = !allZero(page[start : start+decodedSIDBRecordBytes])
			}
			if !valid {
				continue
			}
			data := append([]byte(nil), page[:size]...)
			digest := sha256.Sum256(data)
			if seen[digest] {
				clear(data)
				continue
			}
			registry, err := parseDecodedSIDBRegistry(data)
			if err != nil {
				clear(data)
				continue
			}
			seen[digest] = true
			result = append(result, decodedSIDBRegistryCandidate{Registry: registry, Data: data})
			if len(result) > maximumDecodedSIDBRegistryCandidates {
				for index := range result {
					clear(result[index].Data)
				}
				return nil, fmt.Errorf("decoded SIDB discovery found more than %d candidates", maximumDecodedSIDBRegistryCandidates)
			}
		}
	}
	return result, nil
}

func allZero(data []byte) bool {
	return bytes.Count(data, []byte{0}) == len(data)
}
