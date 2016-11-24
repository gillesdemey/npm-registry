all: deps test build

build:
	go build -o npm-registry server.go

deps:
	go get -t ./...

run:
	bra run

test: deps
	/bin/sh test/test.sh

.PHONY: run build test deps all
