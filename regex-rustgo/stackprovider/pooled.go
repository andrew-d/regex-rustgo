package stackprovider

import (
	"github.com/andrew-d/regex-rustgo/regex-rustgo/internal/stack"
)

type PooledStackProvider struct {
	pool chan *stack.Stack
}

func NewPooledStackProvider(max int) *PooledStackProvider {
	ret := &PooledStackProvider{
		pool: make(chan *stack.Stack, max),
	}
	return ret
}

func (p *PooledStackProvider) get() (*stack.Stack, error) {
	var s *stack.Stack
	var err error

	select {
	case s = <-p.pool:
	default:
		s, err = stack.New()
	}

	if err != nil {
		return nil, err
	}

	return s, nil
}

func (p *PooledStackProvider) put(s *stack.Stack) {
	select {
	case p.pool <- s:
	default:
	}
}

func (s *PooledStackProvider) WithStack(f func(*stack.Stack)) {
	st, err := s.get()
	if err != nil {
		// TODO: signal error?
		return
	}

	defer s.put(st)

	f(st)
}
