// TODO: docs!
package regex

import (
	"fmt"
	"runtime"
	"sync"
	"unsafe"

	"github.com/andrew-d/regex-rustgo/regex-rustgo/internal/stack"
)

var _ = fmt.Printf

//go:cgo_import_static is_match
//go:cgo_import_dynamic is_match
//go:linkname is_match is_match
var is_match uintptr
var _is_match = &is_match

//go:cgo_import_static rust_compile
//go:cgo_import_dynamic rust_compile
//go:linkname rust_compile rust_compile
var rust_compile uintptr
var _rust_compile = &rust_compile

//go:cgo_import_static rust_free
//go:cgo_import_dynamic rust_free
//go:linkname rust_free rust_free
var rust_free uintptr
var _rust_free = &rust_free

func isMatch(stack, re, buf unsafe.Pointer, len uint32, out *bool)
func rustCompile(stack, buf unsafe.Pointer, len uint32, out *unsafe.Pointer)
func rustFree(stack, buf unsafe.Pointer)

var stackPool *sync.Pool

func init() {
	stackPool = &sync.Pool{
		New: func() interface{} {
			data, err := stack.New()
			if err != nil {
				panic(err)
			}

			return data
		},
	}
}

type Regex struct {
	re unsafe.Pointer
}

func Compile(s string) Regex {
	// Get a stack buffer
	stack := stackPool.Get().(*stack.Stack)
	defer stackPool.Put(stack)

	b := []byte(s)

	var re unsafe.Pointer
	rustCompile(
		stack.Bottom(),
		unsafe.Pointer(&b[0]),
		uint32(len(b)),
		&re,
	)

	runtime.KeepAlive(b)
	return Regex{re}
}

func (r *Regex) Free() {
	if r.re == nil {
		return
	}

	// Get a stack buffer
	stack := stackPool.Get().(*stack.Stack)
	defer stackPool.Put(stack)

	rustFree(
		stack.Bottom(),
		r.re,
	)
	r.re = nil
}

func (r *Regex) Match(s string) bool {
	// Get a stack buffer
	stack := stackPool.Get().(*stack.Stack)
	defer stackPool.Put(stack)

	buf := []byte(s)

	var ret bool
	isMatch(
		stack.Bottom(),
		r.re,
		unsafe.Pointer(&buf[0]),
		uint32(len(buf)),
		&ret,
	)

	// NOTE: since we put the Stack object back into our sync.Pool, we
	// don't need to worry about keeping it alive through the call to the
	// above function.  However, if we ever move away from using a
	// sync.Pool, we need to call the following to keep the value alive.
	//     runtime.KeepAlive(stack)

	runtime.KeepAlive(buf)
	return ret
}
