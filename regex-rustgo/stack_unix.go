package regex

import (
	"runtime"
	"syscall"
	"unsafe"
)

const stackSize = 2 * 1024 * 1024

type Stack struct {
	mmap []byte
}

func NewStack() (*Stack, error) {
	data, err := syscall.Mmap(
		-1,
		0,
		stackSize,
		syscall.PROT_READ|syscall.PROT_WRITE,
		syscall.MAP_SHARED|syscall.MAP_ANON,
	)
	if err != nil {
		return nil, err
	}

	// Save the stack and set up a finalizer
	s := &Stack{
		mmap: data,
	}
	runtime.SetFinalizer(s, finalizeStack)

	return s, nil
}

func (s *Stack) Pointer() unsafe.Pointer {
	return unsafe.Pointer(&s.mmap[0])
}

func (s *Stack) Bottom() unsafe.Pointer {
	return unsafe.Pointer(uintptr(s.Pointer()) + stackSize - 32)
}

func finalizeStack(s *Stack) {
	if s.mmap != nil {
		syscall.Munmap(s.mmap)
		s.mmap = nil
	}
}
