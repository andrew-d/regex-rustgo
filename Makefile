IMPORT_PATH      := github.com/andrew-d/regex-rustgo
SYMBOLS          := is_match rust_compile rust_free
LD               ?= ld
export RUSTFLAGS ?= -Ctarget-cpu=native
TARGET           := $(shell GOOS=$(shell go env GOHOSTOS) GOARCH=$(shell go env GOHOSTARCH) \
                            go run target.go $(shell go env GOOS) $(shell go env GOARCH))
RUSTGO_SYSO      := regex-rustgo/libregex_rustgo_$(shell go env GOOS)_$(shell go env GOARCH).syso

regexbench: $(RUSTGO_SYSO) regex-rustgo/rustgo.go regex-rustgo/rustgo.s main.go
ifeq ($(shell go env GOOS),darwin)
	go build -ldflags '-linkmode external -s -extldflags -lresolv' -o $@
else
	go build -ldflags '-linkmode external -extldflags -ldl' -o $@
endif

$(RUSTGO_SYSO): target/$(TARGET)/release/libregex_rustgo.a
ifeq ($(shell go env GOOS),darwin)
		$(LD) -r -o $@ -arch x86_64 $(addprefix -u _,$(SYMBOLS)) $^
else
		$(LD) -r -o $@ --gc-sections $(addprefix -u ,$(SYMBOLS)) $^
endif

target/$(TARGET)/release/libregex_rustgo.a: src/* Cargo.toml Cargo.lock
		cargo build --release --target $(TARGET)

.PHONY: bench
bench: regexbench
		@./regexbench

.PHONY: clean
clean:
		rm -rf regex-rustgo/*.[oa] regex-rustgo/*.syso target regexbench

.PHONY: env
env:
		@echo "TARGET = $(TARGET)"
		@echo "RUSTGO_SYSO = $(RUSTGO_SYSO)"
