package stackprovider

import (
	"github.com/andrew-d/regex-rustgo/regex-rustgo/internal/stack"
)

// StackProvider is the interface for things that can allocate and provide a
// thread stack.
type StackProvider interface {
	WithStack(func(s *stack.Stack))
}
