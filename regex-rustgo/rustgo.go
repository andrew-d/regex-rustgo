//go:binary-only-package

// TODO: docs!
package regex

import _ "unsafe"

//go:cgo_import_static is_match
//go:cgo_import_dynamic is_match
//go:linkname is_match is_match
var is_match uintptr
var _is_match = &is_match

// IsMatch returns a boolean TKTK
func IsMatch(in *[64]byte, out *bool)
