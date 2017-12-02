package regex

import (
	"runtime"
	"unsafe"

	"github.com/andrew-d/regex-rustgo/regex-rustgo/internal/stack"
	"github.com/andrew-d/regex-rustgo/regex-rustgo/stackprovider"
)

type Regex struct {
	re unsafe.Pointer
	st stackprovider.StackProvider
}

func Compile(prov stackprovider.StackProvider, s string) *Regex {
	b := []byte(s)

	var re unsafe.Pointer

	prov.WithStack(func(stack *stack.Stack) {
		rustCompile(
			stack.Bottom(),
			unsafe.Pointer(&b[0]),
			uint32(len(b)),
			&re,
		)
	})

	runtime.KeepAlive(b)
	ret := &Regex{
		re: re,
		st: prov,
	}
	return ret
}

func (r *Regex) Free() {
	if r.re == nil {
		return
	}

	r.st.WithStack(func(stack *stack.Stack) {
		rustFree(
			stack.Bottom(),
			r.re,
		)
	})
	r.re = nil
}

func (r *Regex) Match(s string) bool {
	buf := []byte(s)

	var ret bool
	r.st.WithStack(func(stack *stack.Stack) {
		isMatch(
			stack.Bottom(),
			r.re,
			unsafe.Pointer(&buf[0]),
			uint32(len(buf)),
			&ret,
		)
	})

	runtime.KeepAlive(buf)
	return ret
}

func (r *Regex) FindIndex(s string) ([2]uint32, bool) {
	buf := []byte(s)

	var ret bool
	var match [2]uint32
	r.st.WithStack(func(stack *stack.Stack) {
		findIndex(
			stack.Bottom(),
			r.re,
			unsafe.Pointer(&buf[0]),
			uint32(len(buf)),
			&match,
			&ret,
		)
	})

	runtime.KeepAlive(buf)
	return match, ret
}
