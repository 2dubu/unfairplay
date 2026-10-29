package recovery

import (
	"encoding/binary"
	"fmt"
)

const (
	decodedSIDBNodeHeaderBytes  = 20
	maximumDecodedSIDBTreeDepth = 32
	maximumDecodedSIDBTreeNodes = 4096
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

// parseDecodedSIDBTree verifies node lengths and child counts before exposing
// any geometry. The root may be followed by up to 15 unspecified padding bytes.
func parseDecodedSIDBTree(data []byte) ([]decodedSIDBTreeNode, error) {
	if len(data) < decodedSIDBNodeHeaderBytes {
		return nil, fmt.Errorf("decoded SIDB tree has a truncated root")
	}
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
		length := uint64(sidbTreeWord(data, offset+2, 4, 5))
		if length < decodedSIDBNodeHeaderBytes || length > uint64(limit-offset) {
			return 0, fmt.Errorf("decoded SIDB tree node has invalid length at offset %d", offset)
		}
		end := offset + int(length)
		count := int(sidbTreeWord(data, offset+16, 2, 4))
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
	end, err := parse(0, len(data), 1)
	if err != nil {
		return nil, err
	}
	if len(data)-end > 15 {
		return nil, fmt.Errorf("decoded SIDB tree has excessive trailing padding")
	}
	return nodes, nil
}
