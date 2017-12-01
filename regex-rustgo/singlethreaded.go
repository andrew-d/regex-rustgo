package regex

import (
	"runtime"
	"unsafe"

	"github.com/andrew-d/regex-rustgo/regex-rustgo/internal/stack"
)

type STRegex struct {
	re    unsafe.Pointer
	stack *stack.Stack
}

func CompileST(s string) (*STRegex, error) {
	// Allocate stack for this regex
	st, err := stack.New()
	if err != nil {
		return nil, err
	}

	ret := &STRegex{
		stack: st,
	}

	b := []byte(s)

	rustCompile(
		st.Bottom(),
		unsafe.Pointer(&b[0]),
		uint32(len(b)),
		&ret.re,
	)

	runtime.KeepAlive(b)
	return ret, nil
}

func (r *STRegex) Free() {
	if r.re != nil {
		rustFree(
			r.stack.Bottom(),
			r.re,
		)
		r.re = nil
		r.stack = nil
	}
}

func (r *STRegex) Match(s string) bool {
	buf := []byte(s)

	var ret bool
	isMatch(
		r.stack.Bottom(),
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
