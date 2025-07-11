# Build stage
FROM golang:1.21-alpine AS builder

# Install git and other dependencies
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o elitecode .

# Runtime stage
FROM alpine:latest

# Install necessary packages
RUN apk --no-cache add ca-certificates git docker

# Create non-root user
RUN addgroup -g 1001 -S elitecode && \
    adduser -S elitecode -u 1001

WORKDIR /home/elitecode

# Copy the binary from builder stage
COPY --from=builder /app/elitecode .

# Create directories
RUN mkdir -p .elitecode/cache .elitecode/problems

# Change ownership
RUN chown -R elitecode:elitecode /home/elitecode

# Switch to non-root user
USER elitecode

# Expose port (if needed for any web interface)
EXPOSE 8080

# Default command
CMD ["./elitecode"]