FROM node:18-alpine AS node-builder
WORKDIR /app/nodejs
COPY nodejs/package*.json ./
RUN npm ci --only=production

FROM golang:1.21-alpine AS go-builder
WORKDIR /app/golang
COPY golang/go.mod golang/go.sum ./
RUN go mod download
COPY golang/ ./
RUN go build -o elitecode .

FROM alpine:latest
RUN apk --no-cache add ca-certificates docker
WORKDIR /root/
COPY --from=go-builder /app/golang/elitecode .
COPY --from=node-builder /app/nodejs/node_modules ./node_modules
COPY templates/ ./templates/
CMD ["./elitecode"]