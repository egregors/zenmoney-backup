# Build stage
FROM golang:1.24-alpine AS builder

# Install build dependencies
RUN apk --no-cache add \
    ca-certificates \
    git \
    tzdata

# Set build arguments
ARG CI
ARG CI_COMMIT_BRANCH
ARG CI_MERGE_REQUEST_PROJECT_ID
ARG CI_COMMIT_SHORT_SHA

# Set environment variables
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Create working directory
WORKDIR /build

# Copy go mod files first for better caching
COPY go.mod go.sum ./

# Download dependencies (this will be cached if go.mod/go.sum don't change)
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN set -ex && \
    if [ -z "$CI" ]; then \
        version=$(git describe --tags --always --dirty 2>/dev/null || echo "dev-$(date +%Y%m%d-%H%M%S)"); \
    else \
        version="${CI_COMMIT_BRANCH}${CI_MERGE_REQUEST_PROJECT_ID}-${CI_COMMIT_SHORT_SHA:0:7}-$(date +%Y%m%d-%H:%M:%S)"; \
    fi && \
    echo "Building version: $version" && \
    go build \
        -ldflags="-X main.revision=${version} -s -w -extldflags '-static'" \
        -a -installsuffix cgo \
        -o zenb \
        ./cmd/main.go

# Final stage
FROM scratch

# Copy timezone data
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Copy SSL certificates
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy the binary
COPY --from=builder /build/zenb /zenb

# Create a non-root user (even though we're using scratch)
USER 65534:65534

# Health check (optional, remove if not needed)
# HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
#   CMD ["/zenb", "--help"]

# Expose port if your app uses one (adjust as needed)
# EXPOSE 8080

# Use exec form for better signal handling
CMD ["/zenb"]
