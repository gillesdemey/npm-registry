all: deps test build

build:
	go build -o npm-registry server.go

deps:
	go get -t ./...

run:
	bra run

test:
	go test ./...

.PHONY: run build test deps all
