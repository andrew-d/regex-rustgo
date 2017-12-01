// TODO: docs!
package regex

import (
	"fmt"
	"sync"
	"syscall"
	"unsafe"
)

// Not needed //go:binary-only-package

//go:cgo_import_static is_match
//go:cgo_import_dynamic is_match
//go:linkname is_match is_match
var is_match uintptr
var _is_match = &is_match

func isMatch(ptr unsafe.Pointer, len uint32, out *bool, stack unsafe.Pointer)

const stackSize = 2 * 1024 * 1024

var stackPool *sync.Pool

func init() {
	stackPool = &sync.Pool{
		New: func() interface{} {
			data, err := syscall.Mmap(
				-1,
				0,
				stackSize,
				syscall.PROT_READ|syscall.PROT_WRITE,
				syscall.MAP_SHARED|syscall.MAP_ANON,
			)
			if err != nil {
				panic(err)
			}

			return unsafe.Pointer(&data[0])
		},
	}
}

// IsMatch returns a boolean TKTK
func IsMatch(buf []byte) bool {
	// Get a stack buffer
	stack := stackPool.Get().(unsafe.Pointer)
	defer stackPool.Put(stack)

	// Increment it
	stack = unsafe.Pointer(uintptr(stack) + stackSize - 32)

	fmt.Printf("stack = 0x%x\n", stack)

	out := new(bool)
	isMatch(
		unsafe.Pointer(&buf[0]),
		uint32(len(buf)),
		out,
		stack,
	)
	return *out
}
