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
}
