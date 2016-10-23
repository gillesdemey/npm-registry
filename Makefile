build:
	go build -o npm-registry server/server.go

run:
	bra run

default: build

.PHONY: run build
