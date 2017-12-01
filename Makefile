IMPORT_PATH      := github.com/andrew-d/regex-rustgo
INSTALL_PATH     := $(shell go env GOPATH)/pkg/$(shell go env GOOS)_$(shell go env GOARCH)/$(IMPORT_PATH)
SYMBOL           := is_match
LD               ?= ld
export RUSTFLAGS ?= -Ctarget-cpu=native
TARGET           := $(shell GOOS=$(shell go env GOHOSTOS) GOARCH=$(shell go env GOHOSTARCH) \
                            go run target.go $(shell go env GOOS) $(shell go env GOARCH))

regextest: regex-rustgo/libregex_rustgo.syso
	go build -ldflags '-linkmode external -s -extldflags -lresolv' -o $@

regex-rustgo/libregex_rustgo.syso: target/$(TARGET)/release/libregex_rustgo.a
ifeq ($(shell go env GOOS),darwin)
		$(LD) -r -o $@ -arch x86_64 -u "_$(SYMBOL)" $^
else
		$(LD) -r -o $@ --gc-sections -u "$(SYMBOL)" $^
endif

target/$(TARGET)/release/libregex_rustgo.a: src/* Cargo.toml Cargo.lock
		cargo build --release --target $(TARGET)

.PHONY: clean
clean:
		rm -rf regex-rustgo/*.[oa] regex-rustgo/*.syso target test
