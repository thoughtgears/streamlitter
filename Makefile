ifneq (,$(wildcard ./.env))
    include .env
    export
endif

GIT_COMMIT := $(shell git rev-parse --short --verify HEAD)
GIT_SHA := $(shell git rev-parse --verify HEAD)
VERSION := $(shell cat .version)
SERVICE_NAME = "streamlitter"

.PHONY: clean lint test build

clean:
	@rm -rf builds
	@rm -rf bin

lint:
	@golangci-lint run

test: lint
	@go test -v ./...

build: clean
	@echo "Building $(SERVICE_NAME) version $(VERSION) (commit: $(GIT_COMMIT))..."
	GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags "-s -w" -o builds/$(SERVICE_NAME)-linux-amd64 .
	GOOS=linux GOARCH=arm64 go build -a -installsuffix cgo -ldflags "-s -w" -o builds/$(SERVICE_NAME)-linux-arm64 .
	GOOS=darwin GOARCH=amd64 go build -a -installsuffix cgo -ldflags "-s -w" -o builds/$(SERVICE_NAME)-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build -a -installsuffix cgo -ldflags "-s -w" -o builds/$(SERVICE_NAME)-darwin-arm64 .

deploy: build
	@mkdir -p bin
	@cp builds/$(SERVICE_NAME)-linux-amd64 bin/$(SERVICE_NAME)-$(VERSION)-$(GIT_COMMIT)
