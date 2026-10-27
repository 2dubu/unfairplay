package recovery

import (
	"encoding/binary"
	"testing"
)

func TestOuterDigestMatchesLiftedKernel(t *testing.T) {
	const (
		stackBase    = uint64(0x7000000e0000)
		codecBase    = uint64(0x710000000004)
		hardwareBase = uint64(0x710000010000)
		lookupBase   = uint64(0x710000100000)
	)
	pages := map[uint64][]byte{
		stackBase:                   make([]byte, pageSize),
		codecBase &^ (pageSize - 1): make([]byte, pageSize),
		hardwareBase:                make([]byte, pageSize),
	}
	for address := lookupBase; address < lookupBase+0x20000; address += pageSize {
		pages[address] = make([]byte, pageSize)
	}
	stack := pages[stackBase]
	for offset, value := range map[int]uint32{
		0xca8: 0xceb5d0de, 0xcb0: 0x9fce629e, 0xcb8: 0x9fce6293,
		0xcc0: 0x9fce6288, 0xce8: 0xd47932bc, 0xcf0: 0x9fce629e,
		0xcf8: 0x9fce6293, 0xd00: 0x9fce6288, 0xd08: 0xa905423c,
		0xd10: 0x9fce629e, 0xd18: 0x9fce6293, 0xd20: 0x9fce6288,
		0xd28: 0x1f056ae9, 0xd30: 0x9fce629e, 0xd38: 0x9fce6293,
		0xd40: 0x9fce6288, 0xd48: 0xeba72236, 0xd50: 0x9fce629e,
		0xd58: 0x9fce6293, 0xd60: 0x9fce6288, 0xd68: 0x0000340c,
		0xd70: 0x000048b4,
	} {
		binary.LittleEndian.PutUint32(stack[offset:], value)
	}
	binary.LittleEndian.PutUint64(stack[0xdf0:], lookupBase)

	hardware := [6]byte{0x12, 0x34, 0x56, 0x78, 0x9a, 0xbc}
	codec := [33]byte{
		4: 0x10, 5: 0x21, 6: 0x32, 7: 0x43, 8: 0x54, 9: 0x65,
		10: 0x76, 11: 0x87, 12: 0x98, 13: 0xa9, 14: 0xba, 15: 0xcb, 16: 0xdc,
	}
	copy(pages[hardwareBase], hardware[:])
	copy(pages[codecBase&^(pageSize-1)][codecBase&(pageSize-1):], codec[4:])

	direct := outerDigestState{
		X0: 0x01234567, X4: 0x89abcdef, X10: 0x0f1e2d3c,
		X27: 0x4b5a6978, X28: 0xdeadbeef,
		Counter: outerDigestStop + 2, HardwareSelector: 0x4fc66d4b,
		CodecSelector: 0x60319d93,
	}
	worker := newState(&model{pages: clonePages(pages)})
	worker.sp = stackBase
	worker.x[0], worker.x[4], worker.x[10] = uint64(direct.X0), uint64(direct.X4), uint64(direct.X10)
	worker.x[5], worker.x[8] = uint64(direct.HardwareSelector), uint64(direct.CodecSelector)
	worker.x[9], worker.x[11], worker.x[12] = 0x4ec4ec4f, 6, 13
	worker.x[13], worker.x[16] = codecBase, uint64(direct.Counter)
	worker.x[27], worker.x[28], worker.x[30] = uint64(direct.X27), uint64(direct.X28), hardwareBase

	if err := runRawSourceDigestPrologue(worker, fault{}); err != nil {
		t.Fatal(err)
	}
	terminated := false
	for range outerDigestMaximumIterations {
		if err := runRawSourceDigestKernel(worker, fault{}); err != nil {
			t.Fatal(err)
		}
		if worker.z {
			terminated = true
			break
		}
		if err := runRawSourceDigestPrologue(worker, fault{}); err != nil {
			t.Fatal(err)
		}
	}
	if !terminated {
		t.Fatal("lifted digest did not terminate")
	}
	if !applyOuterDigest(&direct, &hardware, &codec) {
		t.Fatal("semantic digest did not terminate")
	}
	want := outerDigestState{
		X0: uint32(worker.x[0]), X4: uint32(worker.x[4]), X10: uint32(worker.x[10]),
		X27: uint32(worker.x[27]), X28: uint32(worker.x[28]),
		Counter: uint32(worker.x[16]), HardwareSelector: uint32(worker.x[5]),
		CodecSelector: uint32(worker.x[8]),
	}
	if direct != want {
		t.Fatalf("semantic digest = %#v, lifted = %#v", direct, want)
	}
}

func TestOuterDigestRejectsInvalidInput(t *testing.T) {
	var state outerDigestState
	var hardware [6]byte
	var codec [33]byte
	if applyOuterDigest(nil, &hardware, &codec) || applyOuterDigest(&state, nil, &codec) || applyOuterDigest(&state, &hardware, nil) {
		t.Fatal("accepted an invalid digest input")
	}
}

func TestOuterDigestInitialCounter(t *testing.T) {
	vectors := []struct {
		x1, seed, want uint32
	}{
		{0x60319d83, 0x00000000, 0x62477516},
		{0x60319d83, 0x00000001, 0x62477517},
		{0x60319d83, 0x9abcdef0, 0x62477518},
		{0x60319d83, 0xffffffff, 0x6247751c},
		{0x60319e4b, 0x12345678, 0x62477517},
		{0x60319d87, 0x12345678, 0x6247751b},
	}
	for _, vector := range vectors {
		got, ok := outerDigestInitialCounter(vector.x1, vector.seed)
		if !ok || got != vector.want {
			t.Fatalf("counter(%#x, %#x) = %#x/%v, want %#x/true", vector.x1, vector.seed, got, ok, vector.want)
		}
	}
	if _, ok := outerDigestInitialCounter(0, 0); ok {
		t.Fatal("accepted a counter state outside the proved natural path")
	}
}
