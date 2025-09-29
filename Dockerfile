# Bitcoin Echo Node Dockerfile
# Multi-stage build for minimal, secure production image

# Build stage
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy go mod files first for better layer caching
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy source code
COPY . .

# Build the binary with optimizations
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags='-s -w -extldflags "-static"' \
    -a -installsuffix cgo \
    -o bitcoin-echo \
    ./cmd/bitcoin-echo

# Final stage - minimal runtime image
FROM scratch

# Copy certificates for HTTPS requests
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy timezone data
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Copy the binary
COPY --from=builder /app/bitcoin-echo /bitcoin-echo

# Create a non-root user (ID 65534 is 'nobody')
USER 65534

# Expose default Bitcoin ports
EXPOSE 8333 8332 18333 18332

# Health check
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD ["/bitcoin-echo", "--version"]

# Default command
ENTRYPOINT ["/bitcoin-echo"]
CMD ["--help"]

# Metadata
LABEL org.opencontainers.image.title="Bitcoin Echo Node"
LABEL org.opencontainers.image.description="Pure Bitcoin Node Implementation"
LABEL org.opencontainers.image.vendor="Bitcoin Echo Project"
LABEL org.opencontainers.image.source="https://github.com/bitcoinecho/node"
LABEL org.opencontainers.image.documentation="https://bitcoinecho.org"
LABEL org.opencontainers.image.licenses="MIT"