package recovery

import (
	"encoding/binary"
	"encoding/json"
	"os"
	"testing"
)

// This optional oracle is a private native trace; no enrollment material is
// retained by the test or required by the production implementation.
func TestClassWorkspaceSchedulerNativeCellOracle(t *testing.T) {
	path := os.Getenv("UNFAIRPLAY_CLASS_SCHEDULE_ORACLE")
	if path == "" {
		t.Skip("native scheduler trace not configured")
	}
	blob, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var events []struct {
		PC, SP    uint64
		Registers [31]uint64
		Matrix    []byte
	}
	if err := json.Unmarshal(blob, &events); err != nil {
		t.Fatal(err)
	}
	var before, after []struct {
		PC, SP    uint64
		Registers [31]uint64
		Matrix    []byte
	}
	for _, e := range events {
		if e.PC == 0x10011d098 {
			before = append(before, e)
		}
		if e.PC == 0x10011d250 {
			after = append(after, e)
		}
	}
	if len(before) == 0 || len(before) != len(after) {
		t.Fatal("invalid native scheduler checkpoint counts")
	}
	decode := func(data []byte) classWorkspaceScheduleState {
		var s classWorkspaceScheduleState
		s.Seed = binary.LittleEndian.Uint16(data[:2]) - classWorkspacePeerBias
		for i := range 9 {
			s.Product[i] = binary.LittleEndian.Uint16(data[2+2*i:]) - classWorkspaceSampleCutoff
			s.Current[i] = binary.LittleEndian.Uint16(data[20+2*i:]) - classWorkspacePeerBias
			s.Forward[i] = binary.LittleEndian.Uint16(data[38+2*i:]) ^ classWorkspaceOutputEncoding
		}
		return s
	}
	for i, e := range before {
		if len(e.Matrix) < 56 || len(after[i].Matrix) < 56 {
			t.Fatal("truncated native workspace")
		}
		offset := e.Registers[21] - e.Registers[19]
		if offset < 20 || offset >= 38 || offset%2 != 0 {
			t.Fatal("invalid native cell pointer")
		}
		index := int((offset - 20) / 2)
		state := decode(e.Matrix)
		if err := classWorkspaceScheduleCell(&state, index/3, index%3); err != nil {
			t.Fatal(err)
		}
		if state != decode(after[i].Matrix) {
			t.Fatalf("native scheduler cell %d differs", i)
		}
	}

	if len(before)%9 != 0 {
		t.Fatal("native trace does not contain complete passes")
	}
	for i := 0; i < len(before); i += 9 {
		state := decode(before[i].Matrix)
		if err := classWorkspaceSchedulePass(&state); err != nil {
			t.Fatal(err)
		}
		if state != decode(after[i+8].Matrix) {
			t.Fatalf("native scheduler pass %d differs", i/9)
		}
	}
	t.Log("native scheduler cells and passes match")
}
