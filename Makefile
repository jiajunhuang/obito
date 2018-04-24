.PHONY: fmt vet test all

all: fmt vet test
	go build

fmt:
	go fmt ./...

vet:
	go vet -v .

test:
	go test -cover -race
