# Build stage
FROM golang:1.22-alpine AS builder

# Use Chinese mirror for Alpine packages
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

RUN apk add --no-cache git ca-certificates tzdata

# Use Chinese Go proxy for faster downloads in China
ENV GOPROXY=https://goproxy.cn,direct

WORKDIR /src

# Copy both backend and judge (backend imports judge/proto via go.mod replace)
COPY backend/go.mod backend/go.sum ./backend/
COPY judge/go.mod judge/go.sum ./judge/
RUN cd backend && go mod download

# Copy source code
COPY backend/ ./backend/
COPY judge/ ./judge/

# Build the binary
WORKDIR /src/backend
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /app/bin/server ./cmd/server

# Runtime stage
FROM alpine:3.19

# Use Chinese mirror for Alpine packages
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

# Copy binary and config from builder
COPY --from=builder /app/bin/server .
COPY --from=builder /src/backend/config.yaml .
COPY --from=builder /src/backend/locales ./locales

# Create logs directory
RUN mkdir -p /app/logs

# Expose port
EXPOSE 8080

# Run as non-root user
RUN adduser -D -g '' appuser
USER appuser

ENTRYPOINT ["./server"]
