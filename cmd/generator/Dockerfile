# building stage
FROM golang:1.25-alpine AS builder
WORKDIR /

# cache modules first
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download

# copying files
COPY pkg/ ./pkg
COPY cmd/generator ./cmd/generator
COPY proto/ ./proto
COPY internal/ ./internal

# building
RUN --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o generator ./cmd/generator


# final image
FROM alpine:latest
WORKDIR /app

# copying from previous stage
COPY /config ./config
COPY --from=builder /generator ./generator

# creating new user for security
RUN addgroup -S group && adduser -S user -G group
RUN chown user:group /app
USER user

# setting up entry point
ENTRYPOINT ["./generator"]
