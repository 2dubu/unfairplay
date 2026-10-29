package recovery

import (
	"encoding/binary"
	"fmt"
	"testing"
)

// Small synthetic values have a zero width prefix and encode at full width.
const syntheticLargeTreeBytes = 64 << 10

func syntheticSIDBTreeBuffer(size int) []byte {
	data := make([]byte, size)
	for i := range data {
		data[i] = currentClassRecordMask[(i+12)%len(currentClassRecordMask)]
	}
	return data
}

func putSyntheticSIDBNode(data []byte, offset, length, count int) {
	var word [4]byte
	binary.BigEndian.PutUint32(word[:], uint32(length))
	for i, value := range word {
		data[offset+4+i] = value ^ currentClassRecordMask[(offset+4+i+12)%16]
	}
	encodedCount := uint16(count)
	if count >= 4096 {
		encodedCount = 13<<12 | uint16(count)
	}
	binary.BigEndian.PutUint16(word[:2], encodedCount)
	for i, value := range word[:2] {
		data[offset+18+i] = value ^ currentClassRecordMask[(offset+18+i+12)%16]
	}
}

func TestDecodedSIDBTreeAndRegistry(t *testing.T) {
	data := syntheticSIDBTreeBuffer(syntheticLargeTreeBytes)
	rootLen := len(data) - 15
	putSyntheticSIDBNode(data, 0, rootLen, 2)
	putSyntheticSIDBNode(data, 20, 36, 0)
	putSyntheticSIDBNode(data, 56, rootLen-56, 0)
	nodes, err := parseDecodedSIDBTree(data)
	if err != nil {
		t.Fatal(err)
	}
	if len(nodes) != 3 || nodes[1] != (decodedSIDBTreeNode{Offset: 20, Length: 36, ChildCount: 0}) {
		t.Fatalf("unexpected nodes: %+v", nodes)
	}
	registry, err := parseDecodedSIDBRegistry(data)
	if err != nil {
		t.Fatal(err)
	}
	if len(registry.Records) != 1 || registry.Records[0].SourceOffset != 40 || registry.Records[0].TailAvailable != 4 {
		t.Fatalf("unexpected records: %+v", registry.Records)
	}
}

func TestDecodedSIDBTreeAcceptsBoundedRootPadding(t *testing.T) {
	const prefix, suffix = 6, 18
	const rootBytes = decodedSIDBNodeHeaderBytes + 2*(decodedSIDBNodeHeaderBytes+decodedSIDBPayloadBytes)
	data := syntheticSIDBTreeBuffer(prefix + rootBytes + suffix)
	putSyntheticSIDBNode(data, prefix, rootBytes, 2)
	putSyntheticSIDBNode(data, prefix+decodedSIDBNodeHeaderBytes, decodedSIDBNodeHeaderBytes+decodedSIDBPayloadBytes, 0)
	putSyntheticSIDBNode(data, prefix+decodedSIDBNodeHeaderBytes+decodedSIDBNodeHeaderBytes+decodedSIDBPayloadBytes, decodedSIDBNodeHeaderBytes+decodedSIDBPayloadBytes, 0)
	nodes, err := parseDecodedSIDBTree(data)
	if err != nil {
		t.Fatal(err)
	}
	if len(nodes) != 3 || nodes[0].Offset != prefix || nodes[0].Length != rootBytes {
		t.Fatalf("unexpected padded tree: %+v", nodes)
	}
	registry, err := parseDecodedSIDBRegistry(data)
	if err != nil {
		t.Fatal(err)
	}
	if len(registry.Records) != 2 || registry.Records[0].SourceOffset != prefix+2*decodedSIDBNodeHeaderBytes {
		t.Fatalf("unexpected padded records: %+v", registry.Records)
	}
}

func TestDecodedSIDBTreesRetainDisjointRoots(t *testing.T) {
	const rootBytes = decodedSIDBNodeHeaderBytes + decodedSIDBNodeHeaderBytes + decodedSIDBPayloadBytes
	const firstRoot = 6
	const gap = 7
	secondRoot := firstRoot + rootBytes + gap
	data := syntheticSIDBTreeBuffer(secondRoot + rootBytes + 18)
	for _, root := range []int{firstRoot, secondRoot} {
		putSyntheticSIDBNode(data, root, rootBytes, 1)
		putSyntheticSIDBNode(data, root+decodedSIDBNodeHeaderBytes, decodedSIDBNodeHeaderBytes+decodedSIDBPayloadBytes, 0)
	}
	trees, err := parseDecodedSIDBTrees(data)
	if err != nil {
		t.Fatal(err)
	}
	if len(trees) != 2 || trees[0][0].Offset != firstRoot || trees[1][0].Offset != secondRoot {
		t.Fatalf("unexpected roots: %+v", trees)
	}
	registries := decodedSIDBRegistryInterpretations(data)
	if len(registries) != 2 || len(registries[0].Records) != 1 || len(registries[1].Records) != 1 {
		t.Fatalf("unexpected registry roots: %+v", registries)
	}
	if _, err := parseDecodedSIDBTree(data); err == nil {
		t.Fatal("single-tree parser accepted multiple roots")
	}
}

func TestDecodedSIDBTreeRejectsMalformed(t *testing.T) {
	tests := []struct {
		name string
		data []byte
	}{
		{"short", syntheticSIDBTreeBuffer(19)},
	}
	add := func(name string, size, rootLen, count int) []byte {
		d := syntheticSIDBTreeBuffer(size)
		putSyntheticSIDBNode(d, 0, rootLen, count)
		tests = append(tests, struct {
			name string
			data []byte
		}{name, d})
		return d
	}
	add("short length", 20, 19, 0)
	add("length overflow", 20, 0x7fffffff, 0)
	add("too much padding", 40, 20, 0)
	add("child header truncated", 39, 39, 1)
	d := add("child exceeds parent", 40, 40, 1)
	putSyntheticSIDBNode(d, 20, 21, 0)
	d = add("child leaves gap", 41, 41, 1)
	putSyntheticSIDBNode(d, 20, 20, 0)
	add("excessive children", 20, 20, 4096)
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := parseDecodedSIDBTree(tc.data); err == nil {
				t.Fatal("accepted malformed tree")
			}
		})
	}
}

func TestDecodedSIDBTreeDepthLimit(t *testing.T) {
	const levels = maximumDecodedSIDBTreeDepth + 1
	d := syntheticSIDBTreeBuffer(levels * decodedSIDBNodeHeaderBytes)
	for i := 0; i < levels; i++ {
		count := 1
		if i == levels-1 {
			count = 0
		}
		putSyntheticSIDBNode(d, i*20, len(d)-i*20, count)
	}
	if _, err := parseDecodedSIDBTreeAt(d, 0); err == nil {
		t.Fatal("accepted excessive nesting")
	}
}

func TestDecodedSIDBTreeWidthPrefixes(t *testing.T) {
	d := syntheticSIDBTreeBuffer(20)
	// Both prefixes encode the low five bits; unknown bits between prefix and
	// payload are intentionally set to exercise native truncation.
	for i, b := range []byte{0x2f, 0xff, 0xff, 0xf4} {
		d[4+i] = b ^ currentClassRecordMask[(4+i+12)%16]
	}
	for i, b := range []byte{0x5f, 0xe0} {
		d[18+i] = b ^ currentClassRecordMask[(18+i+12)%16]
	}
	if got := sidbTreeWord(d, 4, 4, 5); got != 20 {
		t.Fatalf("length = %d", got)
	}
	if got := sidbTreeWord(d, 18, 2, 4); got != 0 {
		t.Fatalf("count = %d", got)
	}
	if _, err := parseDecodedSIDBTreeAt(d, 0); err != nil {
		t.Fatal(err)
	}
}

func TestClassRecordLayoutMaskAndTailBounds(t *testing.T) {
	var flat decodedSIDBRecord
	flat.TailAvailable = 4
	for i := range flat.Payload {
		flat.Payload[i] = byte(i * 7)
	}
	for i := range flat.Tail {
		flat.Tail[i] = byte(i * 11)
	}
	flat.SourceOffset = 37 // Flat records retain the preexisting relative mask.
	withoutOffset := flat
	withoutOffset.SourceOffset = 0
	for i := 0; i <= maximumClassWindowOffset; i++ {
		a, err := deriveClassValuWindow(flat, i)
		if err != nil {
			t.Fatal(err)
		}
		b, err := deriveClassValuWindow(withoutOffset, i)
		if err != nil {
			t.Fatal(err)
		}
		if a != b {
			t.Fatalf("flat mask changed at window %d", i)
		}
	}
	tree := flat
	tree.MaskOffset = (tree.SourceOffset + 12) % 16
	tree.TailAvailable = 2
	var masked [16]byte
	for i := range masked {
		masked[i] = tree.Payload[i] ^ currentClassRecordMask[(tree.SourceOffset+i+12)%16]
	}
	got, err := deriveClassValuWindow(tree, 0)
	if err != nil {
		t.Fatal(err)
	}
	if got != deriveClassValuFromMaskedWords(masked) {
		t.Fatal("tree mask did not follow absolute position")
	}
	for i := 0; i <= 2; i++ {
		if _, err := deriveClassValuWindow(tree, i); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := deriveClassValuWindow(tree, 3); err == nil {
		t.Fatal("accepted window beyond available tail")
	}
	if got := classRecordMaximumWindowOffset(tree); got != 2 {
		t.Fatalf("maximum window = %d, want 2", got)
	}
}

func TestDecodedSIDBRegistryAcceptsTreesAtMultipleSizes(t *testing.T) {
	for _, size := range []int{56, 944, 1024, syntheticLargeTreeBytes} {
		t.Run(fmt.Sprint(size), func(t *testing.T) {
			data := syntheticSIDBTreeBuffer(size)
			putSyntheticSIDBNode(data, 0, size, 1)
			putSyntheticSIDBNode(data, 20, size-20, 0)
			if _, err := parseDecodedSIDBRegistry(data); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestDecodedSIDBRegistryRejectsFlatShapedInvalidTree(t *testing.T) {
	data := syntheticSIDBTreeBuffer(944)
	putSyntheticSIDBNode(data, 0, len(data), 1)
	putSyntheticSIDBNode(data, 20, 20, 0) // Does not fill the root.
	if _, err := parseDecodedSIDBTree(data); err == nil {
		t.Fatal("accepted invalid tree with flat-compatible length")
	}
}

func FuzzDecodedSIDBTrees(f *testing.F) {
	f.Add([]byte{})
	f.Add(make([]byte, decodedSIDBNodeHeaderBytes))
	seed := syntheticSIDBTreeBuffer(6 + 56 + 18)
	putSyntheticSIDBNode(seed, 6, 56, 1)
	putSyntheticSIDBNode(seed, 26, 36, 0)
	f.Add(seed)
	f.Fuzz(func(t *testing.T, data []byte) {
		if len(data) > 64<<10 {
			return
		}
		roots, err := parseDecodedSIDBTrees(data)
		if err != nil {
			return
		}
		if len(roots) == 0 || len(roots) > maximumDecodedSIDBTreeRoots {
			t.Fatalf("root count = %d", len(roots))
		}
		previousEnd := 0
		for rootIndex, nodes := range roots {
			if len(nodes) == 0 || len(nodes) > maximumDecodedSIDBTreeNodes {
				t.Fatalf("root %d node count = %d", rootIndex, len(nodes))
			}
			start := nodes[0].Offset
			end := start + nodes[0].Length
			if start < previousEnd || start < 0 || end > len(data) || end < start {
				t.Fatalf("root %d range = %d..%d after %d", rootIndex, start, end, previousEnd)
			}
			for nodeIndex, node := range nodes {
				if node.Offset < start || node.Length < decodedSIDBNodeHeaderBytes || node.Offset+node.Length > end {
					t.Fatalf("root %d node %d = %+v", rootIndex, nodeIndex, node)
				}
			}
			previousEnd = end
		}
	})
}

// syntheticSIDBRecordTree supplies 16-byte leaves followed by a small final
// leaf so every candidate has its four actual lookahead bytes in the root.
func syntheticSIDBRecordTree(records, multiplier, seed int) []byte {
	size := decodedSIDBNodeHeaderBytes + records*36 + decodedSIDBNodeHeaderBytes
	data := make([]byte, size)
	for i := range data {
		data[i] = byte(i*multiplier + seed)
	}
	putSyntheticSIDBNode(data, 0, size, records+1)
	for i := 0; i < records; i++ {
		putSyntheticSIDBNode(data, 20+i*36, 36, 0)
	}
	putSyntheticSIDBNode(data, 20+records*36, 20, 0)
	return data
}
