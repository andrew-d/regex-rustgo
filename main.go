package main

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	// "github.com/andrew-d/regex-rustgo/regex-cgo"
	"github.com/andrew-d/regex-rustgo/regex-rustgo"
	"github.com/andrew-d/regex-rustgo/regex-rustgo/stackprovider"
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

	// Stack providers
	poolst := stackprovider.NewPooledStackProvider(10)
	staticst, err := stackprovider.NewStaticStackProvider()
	if err != nil {
		fmt.Printf("error allocating StaticStackProvider: %s\n", err)
		return
	}

	// Test pooled provider
	rustre := regex.Compile(poolst, `.y`)
	defer rustre.Free()

	fmt.Printf("BenchmarkPooledRustRegexp\t%v\n", testing.Benchmark(func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if !rustre.Match(x) {
				b.Fatalf("no match!")
			}
		}
	}))

	// Test static provider
	strustre := regex.Compile(staticst, `.y`)
	defer strustre.Free()

	fmt.Printf("BenchmarkStaticRustRegexp\t%v\n", testing.Benchmark(func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if !strustre.Match(x) {
				b.Fatalf("no match!")
			}
		}
	}))

	/*
		yearRe, err := regex.CompileST(`\d{4}-\d{2}-\d{2}`)
		if err != nil {
			panic(err)
		}
		defer yearRe.Free()
		text := "the date is 2017-10-03"
		match, ok := yearRe.FindIndex(text)
		fmt.Printf("match is %v; index is: (%d, %d); text = %q\n",
			ok,
			match[0],
			match[1],
			text[match[0]:match[1]],
		)
	*/
}
