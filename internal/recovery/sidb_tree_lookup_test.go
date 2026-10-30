package recovery

import (
	"encoding/binary"
	"testing"
)

func putSyntheticSIDBRawWord(data []byte, offset int, value uint32) {
	var b [4]byte
	binary.BigEndian.PutUint32(b[:], value)
	for i := range b {
		data[offset+i] = b[i] ^ currentClassRecordMask[(offset+i+12)%len(currentClassRecordMask)]
	}
}

func TestDecodedSIDBTaggedTreeLookup(t *testing.T) {
	// Root has a user subtree, two duplicate keys, and a separate version.
	data := syntheticSIDBTreeBuffer(142)
	for _, n := range []struct {
		offset, length, count int
		code, id              uint32
	}{
		{0, 140, 2, 10, 0}, {20, 100, 2, 11, 0x12345678},
		{40, 40, 1, 12, 7}, {60, 20, 0, 13, 0}, {80, 40, 1, 12, 7},
		{100, 20, 0, 13, 0}, {120, 20, 0, 14, 1},
	} {
		putSyntheticSIDBNode(data, n.offset, n.length, n.count)
		putSyntheticSIDBRawWord(data, n.offset+6, n.code)
		putSyntheticSIDBRawWord(data, n.offset+10, n.id)
	}
	codes := map[uint32]string{10: "root", 11: "user", 12: "key ", 13: "valu", 14: "vers"}
	decode := func(code uint32) (string, bool) { tag, ok := codes[code]; return tag, ok }
	tree, err := parseDecodedSIDBTaggedTree(data, 2, decode)
	if err != nil {
		t.Fatal(err)
	}
	userID := uint32(0x12345678)
	userID = userID*0xa46fd0ff + 0xfe063531
	for _, tc := range []struct {
		parent int
		tag    string
		id     uint32
		want   int
	}{
		{0, "user", userID, 22}, {22, "key ", 7, 42},
		{0, "vers", 1, 122}, {0, "key ", 7, 0}, {22, "vers", 1, 0},
	} {
		got, err := tree.findChild(tc.parent, tc.tag, tc.id)
		if err != nil || got != tc.want {
			t.Errorf("findChild(%d,%q,%d) = %d, %v; want %d", tc.parent, tc.tag, tc.id, got, err, tc.want)
		}
	}
	for _, tc := range []struct{ parent, offset int }{{42, 62}, {82, 102}} {
		child, err := tree.firstChild(tc.parent)
		if err != nil || child == nil || child.Offset != tc.offset || child.Tag != "valu" {
			t.Fatalf("firstChild(%d) = %+v, %v", tc.parent, child, err)
		}
	}
	if got, err := tree.findChild(0, "valu", 0); err != nil || got != 0 {
		t.Fatal("absent tag with unknown identifier policy should return zero")
	}
	if _, err := tree.findChild(42, "valu", 0); err == nil {
		t.Fatal("accepted unsupported identifier tag")
	}
	if tree.Nodes[1].RawID != 0x12345678 || tree.Nodes[3].Parent != 2 || tree.Nodes[5].Offset != 102 {
		t.Fatalf("incorrect tagged hierarchy: %+v", tree.Nodes)
	}
	if _, err := tree.findChild(21, "key ", 7); err == nil {
		t.Fatal("accepted invalid parent cursor")
	}
	// Reject unknown tags even when an earlier matching child exists.
	delete(codes, 14)
	if _, err := parseDecodedSIDBTaggedTree(data, 2, decode); err == nil {
		t.Fatal("accepted unknown tag")
	}
	if _, err := parseDecodedSIDBTaggedTree(data, 2, nil); err == nil {
		t.Fatal("accepted missing decoder")
	}
}

func TestDecodedSIDBTaggedTreeRejectsMalformedLaterChild(t *testing.T) {
	data := syntheticSIDBTreeBuffer(62)
	putSyntheticSIDBNode(data, 0, 60, 2)
	putSyntheticSIDBNode(data, 20, 20, 0)
	putSyntheticSIDBNode(data, 40, 21, 0)
	decode := func(uint32) (string, bool) { return "key ", true }
	if _, err := parseDecodedSIDBTaggedTree(data, 2, decode); err == nil {
		t.Fatal("accepted malformed later child")
	}
}

func TestDecodedSIDBTaggedTreeWidthIdentifier(t *testing.T) {
	data := syntheticSIDBTreeBuffer(42)
	putSyntheticSIDBNode(data, 0, 40, 1)
	putSyntheticSIDBNode(data, 20, 20, 0)
	// Width five, with extra bits set between the prefix and value.
	putSyntheticSIDBRawWord(data, 30, 0x2fffffe7)
	decode := func(uint32) (string, bool) { return "key ", true }
	tree, err := parseDecodedSIDBTaggedTree(data, 2, decode)
	if err != nil {
		t.Fatal(err)
	}
	got, err := tree.findChild(0, "key ", 7)
	if err != nil || got != 22 {
		t.Fatalf("lookup = %d, %v", got, err)
	}
}

func TestDecodedSIDBTaggedTreeFindChildByTag(t *testing.T) {
	data := syntheticSIDBTreeBuffer(102)
	for _, n := range []struct {
		offset, length, count int
		code, id              uint32
	}{
		{0, 100, 3, 0, 0}, {20, 40, 1, 1, 9}, {40, 20, 0, 1, 7},
		{60, 20, 0, 1, 2}, {80, 20, 0, 2, 0},
	} {
		putSyntheticSIDBNode(data, n.offset, n.length, n.count)
		putSyntheticSIDBRawWord(data, n.offset+6, n.code)
		putSyntheticSIDBRawWord(data, n.offset+10, n.id)
	}
	codes := map[uint32]string{0: "root", 1: "grup", 2: "vers"}
	tree, err := parseDecodedSIDBTaggedTree(data, 2, func(code uint32) (string, bool) { tag, ok := codes[code]; return tag, ok })
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		parent int
		tag    string
		want   int
	}{
		{0, "grup", 22}, {22, "grup", 42}, {0, "vers", 82},
		{22, "vers", 0}, {62, "grup", 0},
	} {
		got, err := tree.findChildByTag(tc.parent, tc.tag)
		if err != nil || got != tc.want {
			t.Errorf("findChildByTag(%d,%q) = %d, %v; want %d", tc.parent, tc.tag, got, err, tc.want)
		}
	}
	if tree.Nodes[1].HasIdentifier {
		t.Fatal("group unexpectedly has a decoded identifier")
	}
	if _, err := tree.findChildByTag(21, "grup"); err == nil {
		t.Fatal("accepted invalid parent")
	}
	if _, err := tree.findChildByTag(0, "bad"); err == nil {
		t.Fatal("accepted invalid tag")
	}
	var empty *decodedSIDBTaggedTree
	if _, err := empty.findChildByTag(0, "grup"); err == nil {
		t.Fatal("accepted nil tree")
	}
}
