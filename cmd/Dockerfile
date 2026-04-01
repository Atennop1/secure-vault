# building stage
FROM golang:1.25-alpine AS builder
WORKDIR /app

# cache modules first
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download

# copying files
COPY pkg/ ./pkg
COPY cmd/ ./cmd
COPY proto/ ./proto
COPY internal/ ./internal

# building
RUN --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o generator ./cmd/generator

RUN --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o storage ./cmd/storage

RUN --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o encoder ./cmd/encoder

RUN --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o decoder ./cmd/decoder
