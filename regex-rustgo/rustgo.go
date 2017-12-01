package regex

import (
	"fmt"
	"unsafe"
)

var _ = fmt.Printf

//go:cgo_import_static is_match
//go:cgo_import_dynamic is_match
//go:linkname is_match is_match
var is_match uintptr
var _is_match = &is_match

//go:cgo_import_static rust_compile
//go:cgo_import_dynamic rust_compile
//go:linkname rust_compile rust_compile
var rust_compile uintptr
var _rust_compile = &rust_compile

//go:cgo_import_static rust_free
//go:cgo_import_dynamic rust_free
//go:linkname rust_free rust_free
var rust_free uintptr
var _rust_free = &rust_free

func isMatch(stack, re, buf unsafe.Pointer, len uint32, out *bool)
func rustCompile(stack, buf unsafe.Pointer, len uint32, out *unsafe.Pointer)
func rustFree(stack, buf unsafe.Pointer)
