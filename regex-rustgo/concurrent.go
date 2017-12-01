package regex

import (
	"runtime"
	"sync"
	"unsafe"

	"github.com/andrew-d/regex-rustgo/regex-rustgo/internal/stack"
)

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

func Compile(s string) *Regex {
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
	return &Regex{re}
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
