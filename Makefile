# Makefile for Vilpanak Real Estate

.PHONY: run build test tidy

run:
	go run ./cmd/web/main.go

build:
	go build -o bin/vilpanak ./cmd/web/main.go

test:
	go test ./...

tidy:
	go mod tidy
