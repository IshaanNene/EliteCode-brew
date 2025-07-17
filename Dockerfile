# Build stage
FROM golang:1.24-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git make

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN make build

# Final stage
FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache ca-certificates docker-cli git

# Copy binary from builder
COPY --from=builder /app/elitecode /usr/local/bin/

# Create config directory
RUN mkdir -p /root/.elitecode

# Set entrypoint
ENTRYPOINT ["elitecode"] 