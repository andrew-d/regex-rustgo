#include "go_asm.h"

TEXT ·isMatch(SB),0,$0-32
	// Set up first two arguments for the function-to-call.  The SysV
	// calling convention is RDI, RSI, RDX, RCX (etc.), so we load the
	// values from the stack and into our arguments.
	//
	// The Go calling convention is essentially that 0(FP) is the first
	// argument, 8(FP) is the second, etc.
	MOVQ ptr+0(FP), DI
	MOVQ len+8(FP), SI
	MOVQ out+16(FP), DX

	// RBX, RBP, and R12–R15 are callee-save, so use RBX to save the
	// existing stack pointer, and then swap to our passed one.
	MOVQ SP, BX
	MOVQ stack+24(FP), SP

	// Load the address of our function-to-call into EAX and call it.
	MOVQ ·_is_match(SB), AX
	CALL AX

	// Restore stack pointer before we exit
	MOVQ BX, SP
	RET
