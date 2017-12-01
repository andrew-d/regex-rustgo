package stack

import (
	"fmt"
	"runtime"
	"syscall"
	"unsafe"
)

var _ = fmt.Printf

const (
	stackSize = 2 * 1024 * 1024
	guardSize = 1 * 1024 * 1024
	redzone   = 4 * 1024 // A single 4KB page
)

type Stack struct {
	mmap []byte
}

func New() (*Stack, error) {
	data, err := syscall.Mmap(
		-1,
		0,
		stackSize+guardSize,
		syscall.PROT_READ|syscall.PROT_WRITE,
		syscall.MAP_SHARED|syscall.MAP_ANON,
	)
	if err != nil {
		return nil, err
	}

	// fmt.Printf("allocated stack at: 0x%x\n", unsafe.Pointer(&data[0]))

	// Since stacks grow from the top down on x86, we want to protect the
	// "bottom" of the stack to prevent a stack overflow.
	err = syscall.Mprotect(data[0:guardSize], syscall.PROT_NONE)
	if err != nil {
		syscall.Munmap(data)
		return nil, err
	}

	// In addition, we protect a "redzone" to ensure that we don't run off the top.
	err = syscall.Mprotect(data[(stackSize+guardSize)-redzone:], syscall.PROT_NONE)
	if err != nil {
		syscall.Munmap(data)
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
	return unsafe.Pointer(uintptr(s.Pointer()) + stackSize - redzone - 32)
}

func finalizeStack(s *Stack) {
	if s.mmap != nil {
		syscall.Munmap(s.mmap)
		s.mmap = nil
	}
}
