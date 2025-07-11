#!/bin/bash

set -e

VERSION=${1:-"v1.0.0"}

echo "Creating release $VERSION..."

# Build binaries for multiple platforms
GOOS=darwin GOARCH=amd64 go build -o bin/elitecode-darwin-amd64 main.go
GOOS=darwin GOARCH=arm64 go build -o bin/elitecode-darwin-arm64 main.go
GOOS=linux GOARCH=amd64 go build -o bin/elitecode-linux-amd64 main.go
GOOS=windows GOARCH=amd64 go build -o bin/elitecode-windows-amd64.exe main.go

# Create release archives
tar -czf bin/elitecode-darwin-amd64.tar.gz -C bin elitecode-darwin-amd64
tar -czf bin/elitecode-darwin-arm64.tar.gz -C bin elitecode-darwin-arm64
tar -czf bin/elitecode-linux-amd64.tar.gz -C bin elitecode-linux-amd64
zip -j bin/elitecode-windows-amd64.zip bin/elitecode-windows-amd64.exe

echo "Release $VERSION created!"