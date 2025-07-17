# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=elitecode
VERSION=0.1.0
COMMIT=$(shell git rev-parse --short HEAD)
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')

# Build flags
LDFLAGS=-ldflags "-X github.com/yourusername/elitecode/cmd.Version=$(VERSION) -X github.com/yourusername/elitecode/cmd.GitCommit=$(COMMIT) -X github.com/yourusername/elitecode/cmd.BuildDate=$(BUILD_TIME)"

.PHONY: all build clean test coverage deps tidy install uninstall

all: deps build

build:
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME)

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f coverage.out

test:
	$(GOTEST) -v ./...

coverage:
	$(GOTEST) -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out

deps:
	$(GOGET) -v ./...

tidy:
	$(GOMOD) tidy

install: build
	mv $(BINARY_NAME) /usr/local/bin/

uninstall:
	rm -f /usr/local/bin/$(BINARY_NAME)

# Development targets
.PHONY: dev emulators lint fmt

dev: deps
	$(GOCMD) run main.go

emulators:
	firebase emulators:start

lint:
	golangci-lint run

fmt:
	$(GOCMD) fmt ./...

# Release targets
.PHONY: release release-linux release-darwin release-windows

release: release-linux release-darwin release-windows

release-linux:
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o dist/$(BINARY_NAME)-linux-amd64
	cd dist && tar czf $(BINARY_NAME)-linux-amd64.tar.gz $(BINARY_NAME)-linux-amd64

release-darwin:
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o dist/$(BINARY_NAME)-darwin-amd64
	cd dist && tar czf $(BINARY_NAME)-darwin-amd64.tar.gz $(BINARY_NAME)-darwin-amd64
	GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o dist/$(BINARY_NAME)-darwin-arm64
	cd dist && tar czf $(BINARY_NAME)-darwin-arm64.tar.gz $(BINARY_NAME)-darwin-arm64

release-windows:
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o dist/$(BINARY_NAME)-windows-amd64.exe
	cd dist && zip $(BINARY_NAME)-windows-amd64.zip $(BINARY_NAME)-windows-amd64.exe

# Docker targets
.PHONY: docker-build docker-run docker-clean

docker-build:
	docker build -t $(BINARY_NAME) .

docker-run:
	docker run --rm -it $(BINARY_NAME)

docker-clean:
	docker rmi $(BINARY_NAME) 