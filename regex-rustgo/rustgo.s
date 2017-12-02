#include "go_asm.h"
#include "textflag.h"

// We use the 'NOSPLIT' flag to prevent Go from inserting a call to add more
// stack to this function.  It's not necessary, since we essentially don't use
// the stack anyway :-)
//
// Another note: in the Go calling convention, all registers are caller-saved
// except:
// - Stack pointer register (SP)
// - Zero register (if there is one)
// - G context pointer register (if there is one)
// - Frame pointer (if there is one)
TEXT ·isMatch(SB), NOSPLIT, $0-40
	// Set up first two arguments for the function-to-call.  The SysV
	// calling convention is RDI, RSI, RDX, RCX (etc.), so we load the
	// values from the stack and into our arguments.
	//
	// The Go calling convention is essentially that 0(FP) is the first
	// argument, 8(FP) is the second, etc.
	MOVQ re+8(FP), DI
	MOVQ ptr+16(FP), SI
	MOVQ len+24(FP), DX
	MOVQ out+32(FP), CX

	// In the SysV calling convention (but not Go's), RBX, RBP, and R12-R15
	// are callee-save, so use RBX to save the existing stack pointer, and
	// then swap to our passed one.
	MOVQ SP, BX
	MOVQ stack+0(FP), SP

	// Note that since both Go and the SysV calling convention have RBP
	// (frame pointer) as a callee-saved register, as long as we don't
	// touch it in this function, we don't need to save it; if the Rust
	// code we're calling touches it, it'll also be restored before the
	// call below returns.

	// TODO: do we want to save the G context pointer register?  At least
	// on OS X, the following Go assembly:
	//
	//     MOVQ (TLS), AX
	//
	// Compiles to:
	//
	//     mov rax, qword [gs:0x8a0]
	//
	// Which seems to indicate that TLS on OS X doesn't use a regular
	// register and thus doesn't need to be saved (since I'm assuming that
	// Rust doesn't screw with segment registers).

	// Load the address of our function-to-call into RAX and call it.
	MOVQ ·_is_match(SB), AX
	CALL AX

	// Restore stack pointer before we exit
	MOVQ BX, SP
	RET

TEXT ·rustCompile(SB), NOSPLIT, $0-32
	MOVQ ptr+8(FP), DI
	MOVQ len+16(FP), SI
	MOVQ out+24(FP), DX

	MOVQ SP, BX
	MOVQ stack+0(FP), SP

	MOVQ ·_rust_compile(SB), AX
	CALL AX

	MOVQ BX, SP
	RET

TEXT ·rustFree(SB), NOSPLIT, $0-16
	MOVQ ptr+8(FP), DI

	MOVQ SP, BX
	MOVQ stack+0(FP), SP

	MOVQ ·_rust_free(SB), AX
	CALL AX

	MOVQ BX, SP
	RET

TEXT ·findIndex(SB), NOSPLIT, $0-48
	MOVQ re+8(FP), DI
	MOVQ ptr+16(FP), SI
	MOVQ len+24(FP), DX
	MOVQ match+32(FP), CX
	MOVQ out+40(FP), R8

	MOVQ SP, BX
	MOVQ stack+0(FP), SP

	MOVQ ·_find_index(SB), AX
	CALL AX

	MOVQ BX, SP
	RET
