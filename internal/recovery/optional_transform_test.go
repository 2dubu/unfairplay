package recovery

import "testing"

func TestOptionalTransform26A428MatchesLiftedKernelDomain(t *testing.T) {
	const (
		inputBase = uint32(0x976dec00)
		selector  = uint32(0x5d5c)
		output    = uint64(0x700000100000)
	)

	// Only execute the data transform and its store. The final four lifted
	// instructions prepare the protected dispatcher continuation and are not
	// part of this value transform.
	kernel := compact26A428Optional
	for index, instruction := range kernel.instructions {
		if instruction.op == compactStoreBase32 {
			kernel.instructions = kernel.instructions[:index+1]
			break
		}
	}

	for low := uint32(0); low < 256; low++ {
		s := newState(&model{pages: map[uint64][]byte{
			output: make([]byte, pageSize),
		}})
		s.x[8] = uint64(inputBase | low)
		s.x[12] = uint64(selector)
		s.x[14] = output
		if err := runCompactKernel(s, kernel); err != nil {
			t.Fatalf("lifted kernel input %#02x: %v", low, err)
		}
		want, ok := optionalTransform26A428(inputBase|low, selector)
		if !ok {
			t.Fatalf("semantic transform rejected input %#02x", low)
		}
		got, err := s.load(output, 4)
		if err != nil {
			t.Fatalf("read lifted output for input %#02x: %v", low, err)
		}
		if uint32(got) != want {
			t.Fatalf("input %#02x: lifted=%#08x semantic=%#08x", low, uint32(got), want)
		}
	}
}

func TestOptionalTransform26A428RejectsUnknownDomain(t *testing.T) {
	for _, test := range []struct {
		input    uint32
		selector uint32
	}{
		{input: 0x976ded00, selector: 0x5d5c},
		{input: 0x976dec00, selector: 0x5d5d},
	} {
		if _, ok := optionalTransform26A428(test.input, test.selector); ok {
			t.Fatalf("accepted unknown domain input=%#08x selector=%#08x", test.input, test.selector)
		}
	}
}
