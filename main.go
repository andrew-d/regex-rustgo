package main

import (
	"fmt"

	"github.com/andrew-d/regex-rustgo/regex-rustgo"
)

func main() {
	match := []byte("2017-11-12")

	var s [64]byte
	copy(s[:], match)

	out := new(bool)

	regex.IsMatch(&s, out)

	fmt.Printf("IsMatch returned: %+v\n", *out)
}
