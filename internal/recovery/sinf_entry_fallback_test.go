package recovery

import "testing"

func newSINFEntryTestState(t *testing.T) (*state, func(uint64, uint64, int)) {
	t.Helper()
	s := newState(&model{pages: make(map[uint64][]byte)})
	s.sp = 0x700000001000
	write := func(address, value uint64, size int) {
		t.Helper()
		for page := address &^ uint64(pageSize-1); page <= (address+uint64(size)-1)&^uint64(pageSize-1); page += pageSize {
			if s.model.pages[page] == nil {
				s.model.pages[page] = make([]byte, pageSize)
			}
		}
		if err := s.storeValue(address, value, size); err != nil {
			t.Fatal(err)
		}
	}
	return s, write
}

func TestSINFEntryDispatch(t *testing.T) {
	const table = uint64(0x710000100000)
	const destination = uint64(0x100484cd0)
	for _, test := range []struct {
		name       string
		selector   uint32
		comparison uint32
		index      int32
		wantEqual  bool
	}{
		{"matching", 0x1234, 0x4440fe62, 0x9667, true},
		{"different", 0x1234, 0x4440fe63, 0x1234, false},
		{"matching-wrap", 0xfffffff0, 0x4440ec1e, -33885, true},
		{"different-wrap", 0xfffffff0, 0x4440ec1f, -16, false},
	} {
		t.Run(test.name, func(t *testing.T) {
			block, ok := compact26A428VerifiedFallbackBlocks[0x1006dd508]
			if !ok {
				t.Fatal("missing SINF entry dispatch block")
			}
			s, write := newSINFEntryTestState(t)
			s.x[10], s.x[11] = uint64(test.comparison), uint64(test.selector)
			write(s.sp+0xdf0, table, 8)
			write(uint64(int64(table)+int64(test.index)*8), destination, 8)
			next, stopped, err := runCompactBlock(s, nil, block)
			if err != nil || stopped || next != destination || s.z != test.wantEqual {
				t.Fatalf("dispatch = %#x, stopped=%v z=%v err=%v", next, stopped, s.z, err)
			}
		})
	}
}

func TestSINFEntryFrameCallAndContinuation(t *testing.T) {
	const (
		table        = uint64(0x710000100000)
		context      = uint64(0x720000001000)
		status       = uint64(0x720000002000)
		left         = uint64(0x720000003000)
		right        = uint64(0x720000004000)
		firstTable   = uint64(0x720000005000)
		secondTable  = uint64(0x720000006000)
		callee       = uint64(0x100020000)
		continuation = uint64(0x100030000)
	)
	for _, test := range []struct {
		name                    string
		selector, first, second uint32
	}{
		{"positive-index", 0x1234, 0x520dfad5, 0x376aef66},
		{"negative-index", 0xfffffff0, 0x520e0699, 0x376afda2},
	} {
		t.Run(test.name, func(t *testing.T) {
			s, write := newSINFEntryTestState(t)
			s.x[11], s.x[20], s.x[28] = 0x10, uint64(test.selector), 0x13579bdf
			write(s.sp+0xb100, status, 8)
			write(s.sp+0xb0f8, context, 8)
			write(s.sp+0xdf0, table, 8)
			write(s.sp+0xde8, left, 8)
			write(s.sp+0xde0, right, 8)
			write(s.sp+0xda0, firstTable, 8)
			write(s.sp+0xda8, secondTable, 8)
			write(status, 0x12345678, 4)
			write(context+0x48, 0, 4)
			write(left, 0x102, 4)
			write(right, 0x107, 4)
			write(s.sp+0x11148, 0, 4)
			write(s.sp+0x11150, 0, 4)
			write(table+0x10*8, 0x100796630, 8)
			write(table+0xd03b*8, callee, 8)
			write(uint64(int64(table)+int64(int32(test.selector))*8), continuation, 8)
			for index := uint64(0); index < 256; index++ {
				write(firstTable+index, (index*3+1)&0xff, 1)
				write(secondTable+index, index^0xa5, 1)
			}
			run := func(address, destination uint64) {
				t.Helper()
				block, ok := compact26A428VerifiedFallbackBlocks[address]
				if !ok {
					t.Fatalf("missing SINF entry block %#x", address)
				}
				next, stopped, err := runCompactBlock(s, nil, block)
				if err != nil || stopped || next != destination {
					t.Fatalf("block %#x: next=%#x stopped=%v err=%v", address, next, stopped, err)
				}
			}
			run(0x100484cd0, 0x100796630)
			if got, err := s.load(status, 4); err != nil || got != 0x13345678 {
				t.Fatalf("status = %#x, %v", got, err)
			}
			if got, err := s.load(context+0x48, 4); err != nil || got != 0x13579bdf {
				t.Fatalf("context word = %#x, %v", got, err)
			}
			run(0x100796630, callee)
			if s.x[0] != s.sp+0x11148 || s.x[30] != 0x100796720 {
				t.Fatalf("call frame=%#x return=%#x", s.x[0], s.x[30])
			}
			for address, want := range map[uint64]uint32{s.sp + 0x11148: test.first, s.sp + 0x11150: test.second} {
				if got, err := s.load(address, 4); err != nil || uint32(got) != want {
					t.Fatalf("frame word at %#x = %#x, %v; want %#x", address, got, err, want)
				}
			}
			write(s.sp+0x1114c, 0xdecafbad, 4)
			run(0x100796720, continuation)
			if s.x[8] != 0xdecafbad {
				t.Fatalf("frame result = %#x", s.x[8])
			}
		})
	}
}
