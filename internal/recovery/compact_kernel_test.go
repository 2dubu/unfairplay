package recovery

import (
	"testing"
	"unsafe"
)

func TestCompactKernelEncoding(t *testing.T) {
	if size := unsafe.Sizeof(compactInstruction{}); size != 8 {
		t.Fatalf("compact instruction size = %d, want 8", size)
	}
	for name, kernel := range map[string]compactKernel{
		"bootstrap":       compact26A428BootstrapKernel,
		"digest-prologue": compact26A428DigestPrologue,
		"cycle-final":     compact26A428CycleFinal,
		"optional":        compact26A428Optional,
		"even":            compact26A428Even,
		"odd":             compact26A428Odd,
		"tail":            compact26A428Tail,
		"writer":          compact26A428Writer,
		"digest":          compact26A428Digest,
		"cycle":           compact26A428Cycle,
		"outer":           compact26A428Outer,
	} {
		if kernel.start == 0 || len(kernel.instructions) == 0 {
			t.Fatalf("%s kernel is empty", name)
		}
		for index, instruction := range kernel.instructions {
			if instruction.op >= compactOpcodeCount {
				t.Fatalf("%s instruction %d has opcode %d", name, index, instruction.op)
			}
		}
	}
	if got := len(compact26A428BootstrapKernel.instructions); got != 31 {
		t.Fatalf("bootstrap instruction count = %d, want 31", got)
	}
	if got := len(compact26A428DigestPrologue.instructions); got != 13 {
		t.Fatalf("digest prologue instruction count = %d, want 13", got)
	}
	if got := len(compact26A428CycleFinal.instructions); got != 69 {
		t.Fatalf("cycle final instruction count = %d, want 69", got)
	}
	if got := len(compact26A428Optional.instructions); got != 117 {
		t.Fatalf("optional instruction count = %d, want 117", got)
	}
	if got := len(compact26A428Even.instructions); got != 220 {
		t.Fatalf("even instruction count = %d, want 220", got)
	}
	if got := len(compact26A428Odd.instructions); got != 311 {
		t.Fatalf("odd instruction count = %d, want 311", got)
	}
	if got := len(compact26A428Tail.instructions); got != 452 {
		t.Fatalf("tail instruction count = %d, want 452", got)
	}
	if got := len(compact26A428Writer.instructions); got != 651 {
		t.Fatalf("writer instruction count = %d, want 651", got)
	}
	if got := len(compact26A428Digest.instructions); got != 1718 {
		t.Fatalf("digest instruction count = %d, want 1718", got)
	}
	if got := len(compact26A428Cycle.instructions); got != 1716 {
		t.Fatalf("cycle instruction count = %d, want 1716", got)
	}
	if got := len(compact26A428Outer.instructions); got != 35206 {
		t.Fatalf("outer instruction count = %d, want 35206", got)
	}
	if got := len(compact26A428Outer.segments); got != 36 {
		t.Fatalf("outer segment count = %d, want 36", got)
	}
	if got := len(compact26A428OuterHooks); got != 6 {
		t.Fatalf("outer hook count = %d, want 6", got)
	}
}
