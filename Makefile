.PHONY: *

build:
	go build -o build/instagram main.go

fmt:
	go fmt ./...

lint:
	golangci-lint run

lint-fix:
	golangci-lint run --fix

test:
	go test ./...
