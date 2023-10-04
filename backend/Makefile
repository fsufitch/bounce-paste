GO_BUILD_FILES := $(shell find . -type f '(' -name go.mod -o -name 'go.sum' -o -name '*.go' ')')

##########

.PHONY: default
default: all

.PHONY: all
all: hello

.PHONY: clean
clean:
	rm -rf build

##########

.PHONY: hello
hello: build/bin/hello

##########

build/bin/hello: ${GO_BUILD_FILES}
	go build -o build/bin/hello ./cmd/hello