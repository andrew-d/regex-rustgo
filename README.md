## regex-rustgo

### What Is This?

This is @andrew-d playing around with calling Rust's excellent [`regex`][re]
crate from Go... without paying overhead of [cgo][cgo] on every function call.
As mentioned in some work done by @BurntSushi ([here][rure]), the cgo function
overhead can eat up the gains you get from using the regex engine.

[re]: https://doc.rust-lang.org/regex/regex/index.html
[cgo]: https://golang.org/cmd/cgo/
[rure]: https://github.com/BurntSushi/rure-go


### How Does It Work?

First off, a huge thanks to Filippo Valsorda (@FiloSottile) for his [blog
post][rustgo] and work on [ed25519-dalek-rustgo][dal]; that project was used as
a starting point for this one.

The essentials behind how this works are very similar to those described in the
above blog post, with some modifications detailed below:

- We don't use the `no_std` feature in Rust, since the `regex` library doesn't
	support it.
- I've wrapped the underlying Rust functions in a simple Go wrapper, exposing
	an API that requires no knowledge of the fact that it's using Rust under the
	hood.  Note that this is still a bit of a work-in-progress 😃
- I've switched from using `go tool` in the Makefile to generating `.syso`
	files, which Go will properly include when running `go build`.
- For the general `Regex` type, we allocate new thread stacks with `mmap` and
	swap to these before calling our Rust code.  This lets us use
	arbitrarily-large stacks for Rust without needing to worry about Go's split
	stacks, and intelligently pools the stacks that are in-use.  As a slight
	bonus, using our own stacks also lets us mark the assembly function as
	`NOSPLIT`, which removes the prelude that checks for stack sizes, saving us a
	couple instructions :-)
- I've also added a `STRegex` type, which preallocates a single stack for use
	by the object; this is faster, since we avoid having to pay the overhead of
	`sync.Pool`, but means that each regex involves allocating a full stack's
	worth of memory (currently, stacks are 3MiB).  This also means that the
	`STRegex` cannot be used concurrently, since every single method on the
	object reuses the same stack.

[rustgo]: https://blog.filippo.io/rustgo/
[dal]: https://github.com/FiloSottile/ed25519-dalek-rustgo/
[syso]: https://github.com/golang/go/wiki/GcToolchainTricks


### How Do I Try It Out Myself?

For now, this only works on OS X and Linux, and only on x86_64 platforms (since
it requires some assembly glue per-architecture and per-calling convention).
If you're on one of these platforms, great!  After cloning the repository, you
can run `make bench` in order to run a quick **non-scientific** benchmark
demonstrating the speedup.  On my reasonably-modern Linux desktop, I see
the following results:

```
BenchmarkGoRegexp        1000000              1466 ns/op
BenchmarkRustRegexp      5000000               255 ns/op
BenchmarkSTRustRegexp   10000000               198 ns/op
```
