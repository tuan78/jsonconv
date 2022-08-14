.PHONY: all
all: build

.PHONY: build
build:
	go build -o bin/jsonconv ./jsonconv/main.go

.PHONY: test
test:
	go test -cover -race ./...

.PHONY: cover
cover:
	go test -race -coverprofile=cover.out -coverpkg=./... ./...
	go tool cover -html=cover.out -o cover.html