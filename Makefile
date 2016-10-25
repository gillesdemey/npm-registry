build:
	go build -o npm-registry server.go

run:
	bra run

default: build

.PHONY: run build
