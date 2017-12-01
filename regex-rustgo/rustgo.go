// TODO: docs!
package regex

import (
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

// IsMatch returns a boolean TKTK
func IsMatch(ptr unsafe.Pointer, len uint32, out *bool, stack unsafe.Pointer)

var stackPool *sync.Pool

func init() {
	stackPool = &sync.Pool{
		New: func() interface{} {
			data, err := syscall.Mmap(
				-1,
				0,
				2*1024*1024,
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

func IsMatchWithStack(buf []byte) bool {
	// Allocate us a stack
	stack := stackPool.Get().(unsafe.Pointer)
	defer stackPool.Put(stack)

	out := new(bool)
	IsMatch(
		unsafe.Pointer(&buf[0]),
		uint32(len(buf)),
		out,
		stack,
	)
	return *out
}
