# Stage 1: Build the Go binary
FROM golang:1.21-alpine AS builder

# Install git (needed for go modules) and upgrade dependencies
RUN apk update && apk add --no-cache git

WORKDIR /app

# Copy go.mod and go.sum first to cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN go build -o server ./cmd/main.go

# Stage 2: Run the binary with a minimal image
FROM alpine:latest

# Install timezone data (optional if you use timezones)
RUN apk add --no-cache ca-certificates tzdata

# Set working directory
WORKDIR /root/

# Copy binary from builder stage
COPY --from=builder /app/server .

# Copy any required config files if needed
COPY internal/config ./internal/config

# Expose the port your Gin app uses
EXPOSE 8080

# Command to run the app
CMD ["./server"]
