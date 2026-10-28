package recovery

import (
	"bytes"
	"encoding/binary"
	"testing"
)

func TestSINFAccountDigestWriterMatchesProtectedLoop(t *testing.T) {
	const stack = uint64(0x200000)
	const output = uint64(0x300000)
	const dispatch = uint64(0x400000)
	const selector = uint32(0x2a03)
	const loopMask = uint32(0x9d82)
	const exitTarget = uint64(0x123456789abcdef0)
	exitAddress := dispatch + uint64(selector)*8
	loopAddress := dispatch + uint64(selector^loopMask)*8
	exitPage := exitAddress &^ (pageSize - 1)
	loopPage := loopAddress &^ (pageSize - 1)

	pages := map[uint64][]byte{
		stack:    make([]byte, pageSize),
		output:   make([]byte, pageSize),
		exitPage: make([]byte, pageSize),
		loopPage: make([]byte, pageSize),
	}
	digest := [20]byte{
		0x41, 0x43, 0x43, 0x4f, 0x55, 0x4e, 0x54, 0x2d, 0x44, 0x49,
		0x47, 0x45, 0x53, 0x54, 0xa5, 0x5a, 0x13, 0x37, 0x24, 0x68,
	}
	for word := 0; word < 5; word++ {
		semantic := binary.BigEndian.Uint32(digest[word*4 : (word+1)*4])
		binary.LittleEndian.PutUint32(pages[stack][0x19c+word*4:], semantic^accountSHA1StateMask)
	}
	binary.LittleEndian.PutUint32(pages[stack][0x1b0:], 20)
	binary.LittleEndian.PutUint64(pages[exitPage][exitAddress-exitPage:], exitTarget)
	binary.LittleEndian.PutUint64(pages[loopPage][loopAddress-loopPage:], sinfAccountDigestWriter)

	makeState := func() *state {
		s := newState(&model{pages: clonePages(pages)})
		s.sp = stack
		s.x[3] = output
		s.x[9] = uint64(selector)
		s.x[10] = 0
		s.x[11] = uint64(loopMask)
		s.x[30] = dispatch
		return s
	}

	protected := makeState()
	var protectedNext uint64
	for iteration := 0; iteration < 5; iteration++ {
		var err error
		protectedNext, err = runSINFDirectBlock1005501bc(protected)
		if err != nil {
			t.Fatal(err)
		}
		if iteration < 4 && protectedNext != sinfAccountDigestWriter {
			t.Fatalf("protected iteration %d target = %#x", iteration, protectedNext)
		}
	}

	semantic := makeState()
	semanticNext, err := runSINFAccountDigestWriter(semantic)
	if err != nil {
		t.Fatal(err)
	}
	if protectedNext != exitTarget || semanticNext != exitTarget || protected.x[10] != 20 || semantic.x[10] != 20 {
		t.Fatalf("writer exits differ: protected=%#x semantic=%#x", protectedNext, semanticNext)
	}
	protectedOutput, err := protected.outputBytes(output, len(digest))
	if err != nil {
		t.Fatal(err)
	}
	semanticOutput, err := semantic.outputBytes(output, len(digest))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(protectedOutput, semanticOutput) {
		t.Fatalf("semantic digest output differs: %x != %x", semanticOutput, protectedOutput)
	}
	for index, value := range semanticOutput {
		if value != digest[index]^0x1b {
			t.Fatalf("digest byte %d = %#02x", index, value)
		}
	}
	if len(semantic.accountDigests) != 1 || semantic.accountDigests[0] != digest {
		t.Fatalf("semantic account digest was not retained")
	}
}

func TestSINFAccountDigestWriterRejectsMidLoopEntry(t *testing.T) {
	s := newState(&model{})
	s.x[10] = 4
	if _, err := runSINFAccountDigestWriter(s); err == nil {
		t.Fatal("mid-loop account-digest writer entry was accepted")
	}
}
