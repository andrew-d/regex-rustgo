package stackprovider

import (
	"github.com/andrew-d/regex-rustgo/regex-rustgo/internal/stack"
)

type StaticStackProvider struct {
	s *stack.Stack
}

func NewStaticStackProvider() (*StaticStackProvider, error) {
	s, err := stack.New()
	if err != nil {
		return nil, err
	}

	ret := &StaticStackProvider{s}
	return ret, nil
}

func (s *StaticStackProvider) WithStack(f func(*stack.Stack)) {
	f(s.s)
}
