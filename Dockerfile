# Multi-stage Docker build for Email Sentinel

# Build stage
FROM golang:1.22-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /build

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
# CGO_ENABLED=0 for static binary (modernc.org/sqlite is CGO-free)
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w -X github.com/datateamsix/email-sentinel/internal/ui.AppVersion=${VERSION:-dev}" \
    -o email-sentinel .

# Runtime stage
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user
RUN addgroup -g 1000 sentinel && \
    adduser -D -u 1000 -G sentinel sentinel

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/email-sentinel .

# Copy default config templates
COPY otp_rules.yaml rules.yaml ./

# Create config directory
RUN mkdir -p /home/sentinel/.config/email-sentinel && \
    chown -R sentinel:sentinel /home/sentinel

# Switch to non-root user
USER sentinel

# Volume for persistent data
VOLUME ["/home/sentinel/.config/email-sentinel"]

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s \
    CMD /app/email-sentinel status || exit 1

# Entry point
ENTRYPOINT ["/app/email-sentinel"]

# Default command
CMD ["start", "--daemon"]
