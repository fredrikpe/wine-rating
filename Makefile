# Makefile for Go project

BINARY_NAME := wine-rating
GO_FILES := $(shell find . -name '*.go' -not -path "./vendor/*")

.PHONY: all build test fmt tidy vet clean

all: build

build:
	GOOS=linux GOARCH=amd64 go build -o bin/$(BINARY_NAME) ./cmd/main.go

test:
	go test ./...

fmt:
	go fmt ./...

tidy:
	go mod tidy

vet:
	go vet ./...

lint:
	golangci-lint run

clean:
	rm -rf bin/

run:
	go run ./cmd/main.go

# helpful combo
check: tidy fmt vet test

