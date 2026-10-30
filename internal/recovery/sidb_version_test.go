package recovery

import (
	"encoding/binary"
	"fmt"
	"testing"
)

func syntheticSIDBVersion(cursor int, version uint32, widths [4]uint32) []byte {
	data := syntheticSIDBTreeBuffer(cursor + 36)
	for i, width := range widths {
		value := version >> (24 - 8*i) & 0xff
		word := width<<27 | value
		if width != 0 && width < 27 {
			word |= 1 << width // discarded high bit, outside the encoded width
		}
		var encoded [4]byte
		binary.BigEndian.PutUint32(encoded[:], word)
		for j, b := range encoded {
			offset := cursor + 20 + i*4 + j
			data[offset] = b ^ currentClassRecordMask[(offset+12)%len(currentClassRecordMask)]
		}
	}
	return data
}

func TestSIDBSchemaVersion(t *testing.T) {
	for _, version := range []sidbSchemaVersion{sidbSchemaVersion10001, sidbSchemaVersion10002} {
		for _, widths := range [][4]uint32{{0, 0, 0, 0}, {1, 5, 8, 31}, {31, 8, 1, 5}} {
			t.Run(fmt.Sprintf("%x-%v", version, widths), func(t *testing.T) {
				data := syntheticSIDBVersion(22, uint32(version), widths)
				got, err := readSIDBSchemaVersion(data, 22, 36)
				if err != nil || got != version {
					t.Fatalf("got %#x, %v; want %#x", got, err, version)
				}
			})
		}
	}
}

func TestSIDBSchemaVersionRejectsUnknownAndInvalidExtents(t *testing.T) {
	for _, version := range []uint32{0, 1, 0x10000, 0x10003, 0xffffffff} {
		data := syntheticSIDBVersion(22, version, [4]uint32{})
		if _, err := readSIDBSchemaVersion(data, 22, 36); err == nil {
			t.Fatalf("accepted unknown version %#x", version)
		}
	}
	data := syntheticSIDBVersion(22, 0x10001, [4]uint32{})
	for _, tc := range []struct{ cursor, length int }{
		{-1, 36}, {22, -1}, {22, 35}, {22, 37}, {59, 36},
		{int(^uint(0) >> 1), 36}, {22, int(^uint(0) >> 1)},
	} {
		if _, err := readSIDBSchemaVersion(data, tc.cursor, tc.length); err == nil {
			t.Fatalf("accepted invalid extent %+v", tc)
		}
	}
	if _, err := readSIDBSchemaVersion(data[:len(data)-1], 22, 36); err == nil {
		t.Fatal("accepted truncated body")
	}
}
