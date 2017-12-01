package main

import (
	"fmt"

	"github.com/andrew-d/regex-rustgo/regex-rustgo"
)

func main() {

	fmt.Println("started")
	defer fmt.Println("finished")

	re := regex.Compile(`^\d{4}-\d{2}-\d{2}$`)
	defer re.Free()

	const test = "2017-11-12"
	fmt.Printf("Trying to match %q: %v\n", test, re.Match(test))
}
