package recovery

import (
	"encoding/binary"
	"fmt"
)

// sidbTagDecoder translates serialized tag codes to FourCC strings. It must
// return false for codes that are not part of the validated codebook.
type sidbTagDecoder func(uint32) (string, bool)

type decodedSIDBTaggedNode struct {
	decodedSIDBTreeNode
	Tag           string
	RawID         uint32
	Identifier    uint32
	HasIdentifier bool
	Parent        int // Index in Nodes, or -1 for the root.
}

type decodedSIDBTaggedTree struct {
	Nodes    []decodedSIDBTaggedNode
	byCursor map[int]int
	children [][]int
}

// sidbTreeRawWord removes only the position mask, retaining all 32 bits.
func sidbTreeRawWord(data []byte, offset int) uint32 {
	var b [4]byte
	for i := range b {
		b[i] = data[offset+i] ^ currentClassRecordMask[(offset+i+12)%len(currentClassRecordMask)]
	}
	return binary.BigEndian.Uint32(b[:])
}

// parseDecodedSIDBTaggedTree validates the complete geometry and every tag
// before exposing lookup. A caller must supply a validated codebook.
func parseDecodedSIDBTaggedTree(data []byte, rootCursor int, decode sidbTagDecoder) (*decodedSIDBTaggedTree, error) {
	if decode == nil {
		return nil, fmt.Errorf("decoded SIDB tree requires a tag decoder")
	}
	geometry, err := parseDecodedSIDBTree(data, rootCursor)
	if err != nil {
		return nil, err
	}
	tree := &decodedSIDBTaggedTree{Nodes: make([]decodedSIDBTaggedNode, 0, len(geometry)), byCursor: make(map[int]int, len(geometry)), children: make([][]int, len(geometry))}
	var stack []int
	for _, node := range geometry {
		for len(stack) > 0 && node.Offset >= tree.Nodes[stack[len(stack)-1]].Offset+tree.Nodes[stack[len(stack)-1]].Length {
			stack = stack[:len(stack)-1]
		}
		parent := -1
		if len(stack) > 0 {
			parent = stack[len(stack)-1]
		}
		code := sidbTreeRawWord(data, node.Offset+4)
		tag, ok := decode(code)
		if !ok || len(tag) != 4 {
			return nil, fmt.Errorf("decoded SIDB tree has unknown tag code at offset %d", node.Offset)
		}
		rawID := sidbTreeRawWord(data, node.Offset+8)
		var id uint32
		var hasID bool
		switch tag {
		case "user":
			id, hasID = rawID*0xa46fd0ff+0xfe063531, true
		case "key ", "vers", "type", "sign":
			id, hasID = sidbTreeWord(data, node.Offset+8, 4, 5), true
		}
		index := len(tree.Nodes)
		tree.Nodes = append(tree.Nodes, decodedSIDBTaggedNode{decodedSIDBTreeNode: node, Tag: tag, RawID: rawID, Identifier: id, HasIdentifier: hasID, Parent: parent})
		tree.byCursor[node.Offset] = index
		if parent >= 0 {
			tree.children[parent] = append(tree.children[parent], index)
		}
		stack = append(stack, index)
	}
	return tree, nil
}

// findChild returns the length-field cursor of the first immediate child in
// serialized order, or zero when absent. The zero parent cursor names the root
// context; nonzero cursors must identify a node in this tree.
func (tree *decodedSIDBTaggedTree) findChild(parentCursor int, tag string, identifier uint32) (int, error) {
	if tree == nil || len(tree.Nodes) == 0 {
		return 0, fmt.Errorf("invalid decoded SIDB tree")
	}
	if len(tag) != 4 {
		return 0, fmt.Errorf("invalid SIDB tag")
	}
	parent := 0
	if parentCursor != 0 {
		var ok bool
		parent, ok = tree.byCursor[parentCursor]
		if !ok {
			return 0, fmt.Errorf("invalid SIDB parent cursor %d", parentCursor)
		}
	}
	for _, index := range tree.children[parent] {
		node := tree.Nodes[index]
		if node.Tag == tag {
			if !node.HasIdentifier {
				return 0, fmt.Errorf("unsupported SIDB identifier tag %q", tag)
			}
			if node.Identifier == identifier {
				return node.Offset, nil
			}
		}
	}
	return 0, nil
}

// findChildByTag returns the first immediate child with tag in serialized
// order, independent of its identifier, or zero when absent. This models the
// observed group selection primitive, not whole-recovery duplicate handling.
func (tree *decodedSIDBTaggedTree) findChildByTag(parentCursor int, tag string) (int, error) {
	if tree == nil || len(tree.Nodes) == 0 {
		return 0, fmt.Errorf("invalid decoded SIDB tree")
	}
	if len(tag) != 4 {
		return 0, fmt.Errorf("invalid SIDB tag")
	}
	parent := 0
	if parentCursor != 0 {
		var ok bool
		parent, ok = tree.byCursor[parentCursor]
		if !ok {
			return 0, fmt.Errorf("invalid SIDB parent cursor %d", parentCursor)
		}
	}
	for _, index := range tree.children[parent] {
		node := tree.Nodes[index]
		if node.Tag == tag {
			return node.Offset, nil
		}
	}
	return 0, nil
}

// firstChild returns the first serialized child, allowing callers to validate its
// tag before accessing its payload. It returns nil when there are no children.
func (tree *decodedSIDBTaggedTree) firstChild(parentCursor int) (*decodedSIDBTaggedNode, error) {
	if tree == nil || len(tree.Nodes) == 0 {
		return nil, fmt.Errorf("invalid decoded SIDB tree")
	}
	parent := 0
	if parentCursor != 0 {
		var ok bool
		parent, ok = tree.byCursor[parentCursor]
		if !ok {
			return nil, fmt.Errorf("invalid SIDB parent cursor %d", parentCursor)
		}
	}
	if len(tree.children[parent]) == 0 {
		return nil, nil
	}
	return &tree.Nodes[tree.children[parent][0]], nil
}
