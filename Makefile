# Makefile for Elitecode CLI

# Variables
BINARY_NAME=elitecode
VERSION=$(shell git describe --tags --always --dirty)
COMMIT=$(shell git rev-parse HEAD)
DATE=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS=-ldflags "-X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)"

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Build targets
.PHONY: all build clean test deps run help install docker-build docker-run dev

all: clean deps test build

build:
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME) -v

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME)-linux -v

build-windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME)-windows.exe -v

build-mac:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME)-mac -v

build-mac-arm:
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME)-mac-arm -v

build-all: build-linux build-windows build-mac build-mac-arm

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_NAME)-linux
	rm -f $(BINARY_NAME)-windows.exe
	rm -f $(BINARY_NAME)-mac
	rm -f $(BINARY_NAME)-mac-arm

test:
	$(GOTEST) -v ./...

test-coverage:
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

deps:
	$(GOMOD) tidy
	$(GOMOD) download

run:
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME) -v
	./$(BINARY_NAME)

install:
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME) -v
	mv $(BINARY_NAME) /usr/local/bin/

# Docker targets
docker-build:
	docker build -t elitecode:latest .

docker-run:
	docker run -it --rm elitecode:latest

docker-compose-up:
	docker-compose up -d

docker-compose-down:
	docker-compose down

docker-compose-logs:
	docker-compose logs -f

# Development targets
dev:
	docker-compose -f docker-compose.yml -f docker-compose.dev.yml up

dev-down:
	docker-compose -f docker-compose.yml -f docker-compose.dev.yml down

# Backend targets
backend-install:
	cd backend && npm install

backend-dev:
	cd backend && npm run dev

backend-test:
	cd backend && npm test

# Release targets
release-dry:
	goreleaser release --snapshot --rm-dist

release:
	goreleaser release --rm-dist

# Linting
lint:
	golangci-lint run

# Security check
security:
	gosec ./...

# Generate documentation
docs:
	$(GOCMD) generate ./...

# Help
help:
	@echo "Available targets:"
	@echo "  build          - Build the binary"
	@echo "  build-all      - Build for all platforms"
	@echo "  clean          - Clean build artifacts"
	@echo "  test           - Run tests"
	@echo "  test-coverage  - Run tests with coverage"
	@echo "  deps           - Download dependencies"
	@echo "  run            - Build and run"
	@echo "  install        - Install to /usr/local/bin"
	@echo "  docker-build   - Build Docker image"
	@echo "  docker-run     - Run Docker container"
	@echo "  dev            - Start development environment"
	@echo "  lint           - Run linter"
	@echo "  security       - Run security checks"
	@echo "  release        - Create release"
	@echo "  help           - Show this help"

.DEFAULT_GOAL := help