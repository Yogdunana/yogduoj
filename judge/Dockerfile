# ---- Build stage ----
FROM golang:1.22-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /src

# Copy go mod files first for layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the judge binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-w -s" -o /judge ./cmd/judge

# ---- Runtime stage ----
FROM alpine:3.19

RUN apk add --no-cache \
    ca-certificates \
    g++ \
    openjdk17-jdk \
    python3 \
    tzdata \
    && rm -rf /var/cache/apk/*

WORKDIR /app

# Copy binary from builder
COPY --from=builder /judge /app/judge

# Create data directories
RUN mkdir -p /data/submissions /data/problems /data/rootfs

EXPOSE 50051

ENTRYPOINT ["/app/judge"]
