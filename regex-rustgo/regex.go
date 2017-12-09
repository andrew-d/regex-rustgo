package regex

import (
	"runtime"
	"unsafe"

	"github.com/andrew-d/regex-rustgo/regex-rustgo/internal/stack"
)

type Regex struct {
	re unsafe.Pointer
	st *stack.Stack
}

func Compile(stack *stack.Stack, s string) *Regex {
	b := []byte(s)

	var re unsafe.Pointer
	rustCompile(
		stack.Bottom(),
		unsafe.Pointer(&b[0]),
		uint32(len(b)),
		&re,
	)

	runtime.KeepAlive(b)
	ret := &Regex{
		re: re,
		st: stack,
	}
	return ret
}

func (r *Regex) Free() {
	if r.re == nil {
		return
	}

	rustFree(
		r.st.Bottom(),
		r.re,
	)
	r.re = nil
}

func (r *Regex) Match(s string) bool {
	buf := []byte(s)

	var ret bool
	isMatch(
		r.st.Bottom(),
		r.re,
		unsafe.Pointer(&buf[0]),
		uint32(len(buf)),
		&ret,
	)

	runtime.KeepAlive(buf)
	return ret
}

func (r *Regex) FindIndex(s string) ([2]uint32, bool) {
	buf := []byte(s)

	var ret bool
	var match [2]uint32
	findIndex(
		r.st.Bottom(),
		r.re,
		unsafe.Pointer(&buf[0]),
		uint32(len(buf)),
		&match,
		&ret,
	)

	runtime.KeepAlive(buf)
	return match, ret
}
