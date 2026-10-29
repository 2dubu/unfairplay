package recovery

import (
	"encoding/binary"
	"fmt"
)

const (
	decodedSIDBNodeHeaderBytes      = 20
	maximumDecodedSIDBTreeDepth     = 32
	maximumDecodedSIDBTreeNodes     = 4096
	maximumDecodedSIDBTreeRoots     = 128
	maximumDecodedSIDBTreeEnvelopes = 16384
)

// decodedSIDBTreeNode records native serialized geometry. Offset is relative
// to the start of the decoded body; Length includes the 20-byte header.
type decodedSIDBTreeNode struct {
	Offset     int
	Length     int
	ChildCount int
}

// sidbTreeWord removes the native position-dependent encoding and then the
// bit-width prefix. A zero prefix represents the full integer width.
func sidbTreeWord(data []byte, offset, size, prefixBits int) uint32 {
	var buf [4]byte
	for i := 0; i < size; i++ {
		buf[4-size+i] = data[offset+i] ^ currentClassRecordMask[(offset+i+12)%len(currentClassRecordMask)]
	}
	word := binary.BigEndian.Uint32(buf[:])
	bits := uint(size * 8)
	width := word >> (bits - uint(prefixBits))
	if width == 0 {
		width = uint32(bits)
	}
	if width == 32 {
		return word
	}
	return word & ((uint32(1) << width) - 1)
}

func parseDecodedSIDBTreeAt(data []byte, rootOffset int) ([]decodedSIDBTreeNode, error) {
	if rootOffset < 0 || rootOffset+decodedSIDBNodeHeaderBytes > len(data) {
		return nil, fmt.Errorf("decoded SIDB tree has a truncated root")
	}
	rootLength := int(sidbTreeWord(data, rootOffset+4, 4, 5))
	if rootLength < decodedSIDBNodeHeaderBytes || rootLength > len(data)-rootOffset {
		return nil, fmt.Errorf("decoded SIDB tree root has invalid length")
	}
	rootEnd := rootOffset + rootLength
	var nodes []decodedSIDBTreeNode
	var parse func(offset, limit, depth int) (int, error)
	parse = func(offset, limit, depth int) (int, error) {
		if depth > maximumDecodedSIDBTreeDepth {
			return 0, fmt.Errorf("decoded SIDB tree exceeds maximum depth")
		}
		if len(nodes) >= maximumDecodedSIDBTreeNodes {
			return 0, fmt.Errorf("decoded SIDB tree exceeds maximum node count")
		}
		if offset < 0 || limit > len(data) || offset > limit || limit-offset < decodedSIDBNodeHeaderBytes {
			return 0, fmt.Errorf("decoded SIDB tree has a truncated node header")
		}
		length := uint64(sidbTreeWord(data, offset+4, 4, 5))
		if length < decodedSIDBNodeHeaderBytes || length > uint64(limit-offset) {
			return 0, fmt.Errorf("decoded SIDB tree node has invalid length at offset %d", offset)
		}
		end := offset + int(length)
		count := int(sidbTreeWord(data, offset+18, 2, 4))
		if count > maximumDecodedSIDBTreeNodes-len(nodes)-1 {
			return 0, fmt.Errorf("decoded SIDB tree has an excessive child count")
		}
		nodes = append(nodes, decodedSIDBTreeNode{Offset: offset, Length: int(length), ChildCount: count})
		if count == 0 {
			return end, nil
		}
		cursor := offset + decodedSIDBNodeHeaderBytes
		for i := 0; i < count; i++ {
			next, err := parse(cursor, end, depth+1)
			if err != nil {
				return 0, err
			}
			cursor = next
		}
		if cursor != end {
			return 0, fmt.Errorf("decoded SIDB tree children do not fill node at offset %d", offset)
		}
		return end, nil
	}
	end, err := parse(rootOffset, rootEnd, 1)
	if err != nil {
		return nil, err
	}
	if end != rootEnd {
		return nil, fmt.Errorf("decoded SIDB tree root length mismatch")
	}
	return nodes, nil
}

// parseDecodedSIDBTrees discovers maximal native roots anywhere in a decoded
// body. Nested nodes are removed from the root set; disjoint roots are retained
// for later personalized-SINF authentication.
func parseDecodedSIDBTrees(data []byte) ([][]decodedSIDBTreeNode, error) {
	if len(data) < decodedSIDBNodeHeaderBytes {
		return nil, fmt.Errorf("decoded SIDB tree has a truncated root")
	}
	var candidates [][]decodedSIDBTreeNode
	type treeEnvelope struct{ start, end int }
	var envelopes []treeEnvelope
	for rootOffset := 0; rootOffset+decodedSIDBNodeHeaderBytes <= len(data); rootOffset++ {
		length := int(sidbTreeWord(data, rootOffset+4, 4, 5))
		count := int(sidbTreeWord(data, rootOffset+18, 2, 4))
		if length >= decodedSIDBNodeHeaderBytes && length <= len(data)-rootOffset && count > 0 && count <= maximumDecodedSIDBTreeNodes {
			envelopes = append(envelopes, treeEnvelope{start: rootOffset, end: rootOffset + length})
			if len(envelopes) > maximumDecodedSIDBTreeEnvelopes {
				return nil, fmt.Errorf("decoded SIDB tree has excessive envelope candidates")
			}
		}
		nodes, err := parseDecodedSIDBTreeAt(data, rootOffset)
		if err != nil {
			continue
		}
		if nodes[0].ChildCount == 0 {
			continue
		}
		candidates = append(candidates, nodes)
		if len(candidates) > maximumDecodedSIDBTreeRoots+maximumDecodedSIDBTreeNodes {
			return nil, fmt.Errorf("decoded SIDB tree has excessive root candidates")
		}
	}
	var roots [][]decodedSIDBTreeNode
	for _, candidate := range candidates {
		start := candidate[0].Offset
		end := start + candidate[0].Length
		contained := false
		for _, envelope := range envelopes {
			if envelope.start <= start && end <= envelope.end && (envelope.start < start || end < envelope.end) {
				contained = true
				break
			}
		}
		if !contained {
			roots = append(roots, candidate)
		}
	}
	if len(roots) == 0 {
		return nil, fmt.Errorf("decoded SIDB tree has no valid root")
	}
	if len(roots) > maximumDecodedSIDBTreeRoots {
		return nil, fmt.Errorf("decoded SIDB tree has more than %d roots", maximumDecodedSIDBTreeRoots)
	}
	for index := 1; index < len(roots); index++ {
		previousEnd := roots[index-1][0].Offset + roots[index-1][0].Length
		if roots[index][0].Offset < previousEnd {
			return nil, fmt.Errorf("decoded SIDB tree has overlapping roots")
		}
	}
	return roots, nil
}

// parseDecodedSIDBTree retains the single-root contract used by callers that
// require one standalone tree.
func parseDecodedSIDBTree(data []byte) ([]decodedSIDBTreeNode, error) {
	roots, err := parseDecodedSIDBTrees(data)
	if err != nil {
		return nil, err
	}
	if len(roots) != 1 {
		return nil, fmt.Errorf("decoded SIDB tree has %d roots, want one", len(roots))
	}
	return roots[0], nil
}
