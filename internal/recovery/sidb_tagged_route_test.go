package recovery

import (
	"bytes"
	"encoding/binary"
	"testing"
)

func syntheticRoute() (sidbTagDecoder, []byte, []byte) {
	codes := map[uint32]string{1: "sean", 2: "vers", 3: "grup", 4: "user", 5: "key ", 6: "valu", 7: "type", 8: "divi", 9: "gval", 10: "salt", 11: "iviv", 12: "dong", 13: "xxxx", 14: "sign"}
	decode := func(c uint32) (string, bool) { s, ok := codes[c]; return s, ok }
	data := syntheticSIDBTreeBuffer(362)
	for _, n := range []struct {
		off, length, count int
		code, id           uint32
	}{
		{0, 360, 2, 1, 0}, {20, 36, 0, 2, 1}, {56, 304, 1, 3, 0}, {76, 284, 2, 4, 0x12345678},
		{96, 36, 0, 2, 1}, {132, 228, 3, 5, 7}, {152, 36, 0, 6, 0}, {188, 24, 0, 7, 1}, {212, 148, 0, 14, 1},
	} {
		putSyntheticSIDBNode(data, n.off, n.length, n.count)
		putSyntheticSIDBRawWord(data, n.off+6, n.code)
		putSyntheticSIDBRawWord(data, n.off+10, n.id)
	}
	setVersion := func(cursor int, v uint32) {
		for i := 0; i < 4; i++ {
			putSyntheticSIDBRawWord(data, cursor+20+4*i, (v>>uint(24-8*i))&255)
		}
	}
	setVersion(22, 0x10001)
	setVersion(98, 0x10002)
	var user, key [4]byte
	userID := uint32(0x12345678)
	userID = userID*0xa46fd0ff + 0xfe063531
	binary.BigEndian.PutUint32(user[:], userID)
	binary.BigEndian.PutUint32(key[:], 7)
	fields := bytes.Join([][]byte{syntheticSINFAtom("user", user[:]), syntheticSINFAtom("key ", key[:]), syntheticSINFAtom("priv", []byte{1})}, nil)
	sinf := syntheticSINFAtom("sinf", syntheticSINFAtom("schi", fields))
	return decode, sinf, data
}

func TestTaggedSIDBRouteSynthetic(t *testing.T) {
	decode, sinf, data := syntheticRoute()
	got, err := selectTaggedSIDBRoute(decode, sinf, data, 2)
	if err != nil {
		t.Fatal(err)
	}
	if got.Version != sidbSchemaVersion10001 || got.AccountVersion != sidbSchemaVersion10002 || got.Record.SourceOffset != 174 || got.TypeCursor != 190 || got.TypeLength != 24 {
		t.Fatalf("unexpected selection: %+v", got)
	}
	if !bytes.Equal(got.Window[:], data[174:190]) {
		t.Fatal("window did not use cursor+20")
	}
	reg, err := parseDecodedSIDBRegistry(data, 2)
	if err != nil || got.RecordIndex >= len(reg.Records) || reg.Records[got.RecordIndex] != got.Record {
		t.Fatal("selected record mismatch", err)
	}
}

func TestTaggedSIDBRouteRejectsUnsupported(t *testing.T) {
	for _, tc := range []struct {
		name   string
		mutate func([]byte, []byte)
	}{
		{"bad sinf", func(s, d []byte) { s[3] = 0 }},
		{"wrong root tag", func(s, d []byte) { putSyntheticSIDBRawWord(d, 6, 2) }},
		{"missing group", func(s, d []byte) { putSyntheticSIDBRawWord(d, 62, 13) }},
		{"wrong user", func(s, d []byte) { putSyntheticSIDBRawWord(d, 86, 42) }},
		{"wrong key", func(s, d []byte) { putSyntheticSIDBRawWord(d, 142, 8) }},
		{"missing type", func(s, d []byte) { putSyntheticSIDBRawWord(d, 194, 13) }},
		{"missing sign", func(s, d []byte) { putSyntheticSIDBRawWord(d, 218, 13) }},
		{"wrong sign id", func(s, d []byte) { putSyntheticSIDBRawWord(d, 222, 2) }},
		{"sign length", func(s, d []byte) { putSyntheticSIDBNode(d, 212, 144, 0) }},
		{"sign not leaf", func(s, d []byte) { putSyntheticSIDBNode(d, 212, 148, 1) }},
		{"value not first", func(s, d []byte) {
			putSyntheticSIDBRawWord(d, 158, 7)
			putSyntheticSIDBRawWord(d, 162, 1)
			putSyntheticSIDBRawWord(d, 194, 6)
		}},
		{"account version old", func(s, d []byte) { putSyntheticSIDBRawWord(d, 130, 1) }},
		{"root version unknown", func(s, d []byte) { putSyntheticSIDBRawWord(d, 54, 3) }},
		{"unsupported key salt", func(s, d []byte) { putSyntheticSIDBRawWord(d, 194, 10) }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			decode, sinf, data := syntheticRoute()
			tc.mutate(sinf, data)
			if _, err := selectTaggedSIDBRoute(decode, sinf, data, 2); err == nil {
				t.Fatal("accepted unsupported route")
			}
		})
	}
}

func TestTaggedSIDBRouteRejectsDuplicateAndOptionalChildren(t *testing.T) {
	for _, tc := range []struct {
		name   string
		parent string
		code   uint32
	}{
		{"duplicate group", "root", 3}, {"root divi", "root", 8},
		{"duplicate user", "group", 4}, {"duplicate key", "user", 5}, {"user gval", "user", 9}, {"key salt", "key", 10},
		{"key iviv", "key", 11}, {"key dong", "key", 12},
	} {
		t.Run(tc.name, func(t *testing.T) {
			decode, sinf, data := syntheticRoute()
			data = append(data, syntheticSIDBTreeBuffer(20)...)
			putSyntheticSIDBNode(data, 0, 380, 2)
			switch tc.parent {
			case "root":
				putSyntheticSIDBNode(data, 0, 380, 3)
			case "group":
				putSyntheticSIDBNode(data, 56, 324, 2)
			case "user":
				putSyntheticSIDBNode(data, 56, 324, 1)
				putSyntheticSIDBNode(data, 76, 304, 3)
			case "key":
				putSyntheticSIDBNode(data, 56, 324, 1)
				putSyntheticSIDBNode(data, 76, 304, 2)
				putSyntheticSIDBNode(data, 132, 248, 4)
			}
			putSyntheticSIDBNode(data, 360, 20, 0)
			putSyntheticSIDBRawWord(data, 366, tc.code)
			id := uint32(0)
			if tc.name == "duplicate user" {
				id = 0x12345678
			}
			if tc.name == "duplicate key" {
				id = 7
			}
			putSyntheticSIDBRawWord(data, 370, id)
			if _, err := selectTaggedSIDBRoute(decode, sinf, data, 2); err == nil {
				t.Fatal("accepted duplicate or unmodeled route")
			}
		})
	}
}

func TestTaggedSIDBRouteTypeGeometryAndPayload(t *testing.T) {
	for _, payload := range []uint32{0, 1, 0xffffffff} {
		decode, sinf, data := syntheticRoute()
		putSyntheticSIDBRawWord(data, 208, payload)
		if _, err := selectTaggedSIDBRoute(decode, sinf, data, 2); err != nil {
			t.Fatal("type payload unexpectedly gated", err)
		}
	}
	for _, tc := range []struct {
		name              string
		typeLength, count int
	}{
		{"short", 20, 0}, {"long", 28, 0}, {"has child", 40, 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			decode, sinf, data := syntheticRoute()
			delta := tc.typeLength - 24
			if delta < 0 {
				data = data[:len(data)+delta]
			} else {
				data = append(data, syntheticSIDBTreeBuffer(delta)...)
			}
			putSyntheticSIDBNode(data, 0, 360+delta, 2)
			putSyntheticSIDBNode(data, 56, 304+delta, 1)
			putSyntheticSIDBNode(data, 76, 284+delta, 2)
			putSyntheticSIDBNode(data, 132, 228+delta, 3)
			putSyntheticSIDBNode(data, 188, tc.typeLength, tc.count)
			if tc.count != 0 {
				putSyntheticSIDBNode(data, 208, 20, 0)
				putSyntheticSIDBRawWord(data, 214, 13)
				putSyntheticSIDBRawWord(data, 218, 0)
			}
			if _, err := selectTaggedSIDBRoute(decode, sinf, data, 2); err == nil {
				t.Fatal("accepted unsupported type geometry")
			}
		})
	}
}
