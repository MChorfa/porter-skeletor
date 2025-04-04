# Dockerfile for skeletor

# Use a specific Go version for the builder stage
# Keep this in sync with the go.mod version
ARG GO_VERSION=1.23
FROM golang:${GO_VERSION}-alpine AS builder

WORKDIR /src

# Install build dependencies if any (e.g., git for private modules)
# RUN apk add --no-cache git

# Copy go module files and download dependencies first to leverage Docker cache
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the generator binary
# Using CGO_ENABLED=0 produces a static binary
# Using ldflags to strip debug information reduces binary size
ARG TARGETOS TARGETARCH
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build \
    -ldflags="-s -w" \
    -o /out/skeletor \
    ./cmd/skeletor

# Use a minimal base image like scratch or distroless for the final stage
# Using alpine as a simple, small base
FROM alpine:latest

# Copy the static binary from the builder stage
COPY --from=builder /out/skeletor /usr/local/bin/skeletor

# Set the entrypoint
ENTRYPOINT ["/usr/local/bin/skeletor"]

# Default command (optional, e.g., show help)
CMD ["--help"]
