package main

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/andrew-d/regex-rustgo/regex-rustgo"
)

func main() {
	fmt.Println("started")
	defer fmt.Println("finished")

	x := strings.Repeat("x", 50) + "y"

	gore := regexp.MustCompile(`.y`)
	fmt.Printf("BenchmarkGoRegexp\t%v\n", testing.Benchmark(func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if !gore.MatchString(x) {
				b.Fatalf("no match!")
			}
		}
	}))

	rustre := regex.Compile(`.y`)
	defer rustre.Free()

	fmt.Printf("BenchmarkRustRegexp\t%v\n", testing.Benchmark(func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if !rustre.Match(x) {
				b.Fatalf("no match!")
			}
		}
	}))
}
