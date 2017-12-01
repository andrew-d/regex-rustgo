package main

import (
	"fmt"

	"github.com/andrew-d/regex-rustgo/regex-rustgo"
)

func main() {
	match := []byte("2017-11-12")
	out := regex.IsMatchWithStack(match)

	fmt.Printf("IsMatch returned: %+v\n", out)
}
