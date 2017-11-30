IMPORT_PATH      := github.com/andrew-d/regex-rustgo
INSTALL_PATH     := $(shell go env GOPATH)/pkg/$(shell go env GOOS)_$(shell go env GOARCH)/$(IMPORT_PATH)
SYMBOL           := is_match
LD               ?= ld
export RUSTFLAGS ?= -Ctarget-cpu=native
TARGET           := $(shell GOOS=$(shell go env GOHOSTOS) GOARCH=$(shell go env GOHOSTARCH) \
                            go run target.go $(shell go env GOOS) $(shell go env GOARCH))

regex-rustgo/regex-rustgo.a: regex-rustgo/rustgo.go regex-rustgo/rustgo.o regex-rustgo/libregex_rustgo.o
		go tool compile -N -l -o $@ -p main -pack regex-rustgo/rustgo.go
		go tool pack r $@ regex-rustgo/rustgo.o regex-rustgo/libregex_rustgo.o

regex-rustgo/libregex_rustgo.o: target/$(TARGET)/release/libregex_rustgo.a
ifeq ($(shell go env GOOS),darwin)
		$(LD) -r -o $@ -arch x86_64 -u "_$(SYMBOL)" $^
else
		$(LD) -r -o $@ --gc-sections -u "$(SYMBOL)" $^
endif

target/$(TARGET)/release/libregex_rustgo.a: src/* Cargo.toml Cargo.lock
		cargo build --release --target $(TARGET)

regex-rustgo/rustgo.o: regex-rustgo/rustgo.s
		go tool asm -I "$(shell go env GOROOT)/pkg/include" -o $@ $^

.PHONY: install uninstall
install: regex-rustgo/regex-rustgo.a
		mkdir -p "$(INSTALL_PATH)"
		cp regex-rustgo/regex-rustgo.a "$(INSTALL_PATH)"
uninstall:
		rm -f "$(INSTALL_PATH)/regex-rustgo.a"

.PHONY: clean
clean:
		rm -rf regex-rustgo/*.[oa] target regex-rustgo
