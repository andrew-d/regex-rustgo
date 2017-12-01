// TODO: docs!
package regex

import (
	"sync"
	"unsafe"
)

// Not needed //go:binary-only-package

//go:cgo_import_static is_match
//go:cgo_import_dynamic is_match
//go:linkname is_match is_match
var is_match uintptr
var _is_match = &is_match

func isMatch(ptr unsafe.Pointer, len uint32, out *bool, stack unsafe.Pointer)

var stackPool *sync.Pool

func init() {
	stackPool = &sync.Pool{
		New: func() interface{} {
			data, err := NewStack()
			if err != nil {
				panic(err)
			}

			return data
		},
	}
}

// IsMatch returns a boolean TKTK
func IsMatch(buf []byte) bool {
	// Get a stack buffer
	stack := stackPool.Get().(*Stack)
	defer stackPool.Put(stack)

	out := new(bool)
	isMatch(
		unsafe.Pointer(&buf[0]),
		uint32(len(buf)),
		out,
		stack.Bottom(),
	)

	// NOTE: since we put the Stack object back into our sync.Pool, we
	// don't need to worry about keeping it alive through the call to the
	// above function.  However, if we ever move away from using a
	// sync.Pool, we need to call the following to keep the value alive.
	//     runtime.KeepAlive(stack)

	// Extract and pass on return values.
	return *out
}
