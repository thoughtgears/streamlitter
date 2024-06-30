ifneq (,$(wildcard ./.env))
    include .env
    export
endif

GIT_COMMIT := $(shell git rev-parse --short --verify HEAD)
GIT_SHA := $(shell git rev-parse --verify HEAD)
VERSION := $(shell cat .version)
SERVICE_NAME = streamlitter

.PHONY: clean lint test build bundle

clean:
	@rm -rf builds
	@rm -rf bin

lint:
	@golangci-lint run -c .golangci.yml

test: lint
	@go test -v ./...

build: clean
	@echo "Building $(SERVICE_NAME) version $(VERSION) (commit: $(GIT_COMMIT))..."
	GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags "-s -w -X main.Version=$(VERSION)" -o builds/$(SERVICE_NAME)-linux-amd64 .
	GOOS=linux GOARCH=arm64 go build -a -installsuffix cgo -ldflags "-s -w -X main.Version=$(VERSION)" -o builds/$(SERVICE_NAME)-linux-arm64 .
	GOOS=darwin GOARCH=amd64 go build -a -installsuffix cgo -ldflags "-s -w -X main.Version=$(VERSION)" -o builds/$(SERVICE_NAME)-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build -a -installsuffix cgo -ldflags "-s -w -X main.Version=$(VERSION)" -o builds/$(SERVICE_NAME)-darwin-arm64 .

bundle: build
	@mkdir -p bin
	@cp builds/$(SERVICE_NAME)-linux-amd64 bin/main-linux-amd64-$(VERSION)
	@cp builds/$(SERVICE_NAME)-linux-arm64 bin/main-linux-arm64-$(VERSION)
	@cp builds/$(SERVICE_NAME)-darwin-amd64 bin/main-darwin-amd64-$(VERSION)
	@cp builds/$(SERVICE_NAME)-darwin-arm64 bin/main-darwin-arm64-$(VERSION)

	@git add bin/*
	@git commit -m "Bundling version $(VERSION) on commit $(GIT_COMMIT)"
