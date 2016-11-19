build: test
	go build -o npm-registry server.go

run:
	bra run

test:
	go test ./...

.PHONY: run build test
