.PHONY: all test
all: build

build:
	go build -o bin/jsonconv github.com/tuan78/jsonconv/cmd/jsonconv

test:
	go test ./... -cover