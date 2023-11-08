.PHONY: *

build:
	go build -o build/instagram main.go

docs:
	rm -rf docs/*
	go run cmd/docs/main.go

fmt:
	go fmt ./...

html-coverage:
	go tool cover -html=coverage.out

lint:
	golangci-lint run

lint-fix:
	golangci-lint run --fix

mocks:
	mockery

test:
	go test -coverpkg=./... -race -coverprofile=coverage.out -shuffle on ./...
	cat coverage.out | grep -v 'pkg/filesystem/' > coverage.temp
	mv coverage.temp coverage.out
