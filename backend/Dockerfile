# Build stage
FROM golang:1.22-alpine AS builder

RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

# Copy go mod files first for dependency caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /app/bin/server ./cmd/server

# Runtime stage
FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/bin/server .
COPY --from=builder /app/config.yaml .
COPY --from=builder /app/locales ./locales

# Create logs directory
RUN mkdir -p /app/logs

# Expose port
EXPOSE 8080

# Run as non-root user
RUN adduser -D -g '' appuser
USER appuser

ENTRYPOINT ["./server"]
