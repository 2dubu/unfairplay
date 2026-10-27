package recovery

import (
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"os"
	"testing"
)

func TestAccountSHA1CompressBlock(t *testing.T) {
	initial := accountSHA1State{
		A: 0x67452301, B: 0xefcdab89, C: 0x98badcfe,
		D: 0x10325476, E: 0xc3d2e1f0,
	}
	var block [64]byte
	copy(block[:], "abc")
	block[3] = 0x80
	binary.BigEndian.PutUint64(block[56:], 3*8)
	got := accountSHA1CompressBlock(initial, block, 0).State
	wantBytes, err := hex.DecodeString("a9993e364706816aba3e25717850c26c9cd0d89d")
	if err != nil {
		t.Fatal(err)
	}
	want := accountSHA1State{
		A: binary.BigEndian.Uint32(wantBytes[0:4]), B: binary.BigEndian.Uint32(wantBytes[4:8]),
		C: binary.BigEndian.Uint32(wantBytes[8:12]), D: binary.BigEndian.Uint32(wantBytes[12:16]),
		E: binary.BigEndian.Uint32(wantBytes[16:20]),
	}
	if got != want {
		t.Fatalf("SHA-1(abc) = %#v, want %#v", got, want)
	}
}

func TestAccountSHA1Round(t *testing.T) {
	state := accountSHA1State{
		A: 0x67452301, B: 0xefcdab89, C: 0x98badcfe,
		D: 0x10325476, E: 0xc3d2e1f0,
	}
	got := accountSHA1Round(state, 0x61626380, accountSHA1Choose)
	want := accountSHA1State{
		A: 0x0116fc33, B: 0x67452301, C: 0x7bf36ae2,
		D: 0x98badcfe, E: 0x10325476,
	}
	if got != want {
		t.Fatalf("SHA-1 round = %#v, want %#v", got, want)
	}
}

type accountSHAProbe struct {
	PC   uint64     `json:"pc"`
	X    [31]uint64 `json:"x"`
	Word uint32     `json:"word"`
}

func accountSHAProbeState(t *testing.T, probe accountSHAProbe, phase int) *state {
	t.Helper()
	const tableBase = uint64(0x710000100000)
	const dispatchBase = uint64(0x710000200000)
	pages := make(map[uint64][]byte)
	pages[tableBase] = make([]byte, pageSize)
	for address := dispatchBase; address < dispatchBase+0x90000; address += pageSize {
		pages[address] = make([]byte, pageSize)
	}
	for address := dispatchBase; address < dispatchBase+0x90000; address += 8 {
		base := address &^ (pageSize - 1)
		binary.LittleEndian.PutUint64(pages[base][address-base:], 0x100000000+(address-dispatchBase))
	}
	result := newState(&model{pages: pages, binaryHash: rawSourceBinaryHash})
	result.x = probe.X
	result.x[6], result.x[10], result.x[21] = 0, tableBase, dispatchBase
	counterRegister := [...]int{16, 20, 23, 20}[phase]
	result.x[counterRegister] = 5
	binary.LittleEndian.PutUint32(pages[tableBase][20:], probe.Word)
	return result
}

func TestAccountSHA1RoundsMatchCircuit(t *testing.T) {
	path := os.Getenv("UNFAIRPLAY_SHA_PROBES")
	if path == "" {
		t.Skip("private SHA-1 round probes unavailable")
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var probes []accountSHAProbe
	if err := json.Unmarshal(data, &probes); err != nil {
		t.Fatal(err)
	}
	phases := map[uint64]int{
		accountHotMixFront: 0, accountHotMixSecond: 1,
		accountHotMixThirdFront: 2, accountHotMixFourth: 3,
	}
	circuits := [...]func(*state) (uint64, error){
		runAccountHotMixCircuit, runAccountHotMixSecondCircuit,
		runAccountHotMixThirdCircuit, runAccountHotMixFourthCircuit,
	}
	semantic := [...]func(*state) (uint64, error){
		runAccountHotMix, runAccountHotMixSecond,
		runAccountHotMixThird, runAccountHotMixFourth,
	}
	for index, probe := range probes {
		phase, ok := phases[probe.PC]
		if !ok {
			t.Fatalf("probe %d has unknown phase %#x", index, probe.PC)
		}
		want := accountSHAProbeState(t, probe, phase)
		got := newState(want.model)
		got.x = want.x
		wantNext, wantErr := circuits[phase](want)
		gotNext, gotErr := semantic[phase](got)
		if wantErr != nil || gotErr != nil {
			t.Fatalf("probe %d phase %d errors: %v / %v", index, phase, wantErr, gotErr)
		}
		if gotNext != wantNext || got.x != want.x {
			for register := range got.x {
				if got.x[register] != want.x[register] {
					t.Fatalf("probe %d phase %d x%d=%#x, want %#x", index, phase, register, got.x[register], want.x[register])
				}
			}
			t.Fatalf("probe %d phase %d next=%#x, want %#x", index, phase, gotNext, wantNext)
		}
	}
}
