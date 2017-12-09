package regex

import (
	"github.com/andrew-d/regex-rustgo/regex-rustgo/internal/stack"
)

type PooledRegex struct {
	pool chan *Regex
	s    string
}

func NewPooledRegex(s string, max int) *PooledRegex {
	ret := &PooledRegex{
		pool: make(chan *Regex, max),
		s:    s,
	}
	return ret
}

func (p *PooledRegex) Get() (*Regex, error) {
	select {
	case r := <-p.pool:
		return r, nil
	default:
	}

	st, err := stack.New()
	if err != nil {
		return nil, err
	}

	return Compile(st, p.s), nil
}

func (p *PooledRegex) Put(r *Regex) {
	select {
	case p.pool <- r:
	default:
	}
}
