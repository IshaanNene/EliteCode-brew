#!/bin/bash

set -e

echo "Building Elitecode CLI..."

# Build Go binary
go build -o bin/elitecode main.go

# Build Docker images
echo "Building Docker images..."
docker build -t elitecode/c:latest -f internal/docker/templates/c.dockerfile .
docker build -t elitecode/cpp:latest -f internal/docker/templates/cpp.dockerfile .
docker build -t elitecode/python:latest -f internal/docker/templates/python.dockerfile .
docker build -t elitecode/java:latest -f internal/docker/templates/java.dockerfile .
docker build -t elitecode/javascript:latest -f internal/docker/templates/javascript.dockerfile .

echo "Build complete!"