# building stage
FROM golang:1.25-alpine AS builder
WORKDIR /app

# cache modules first
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download

# copying files
COPY pkg/ ./pkg
COPY cmd/decoder ./cmd/decoder
COPY proto/ ./proto
COPY internal/ ./internal

# building
RUN --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o decoder ./cmd/decoder


# final image
FROM alpine:latest
WORKDIR /app

# copying from previous stage
COPY /config ./config
COPY --from=builder /app/decoder ./decoder

# creating new user for security
RUN addgroup -S group && adduser -S user -G group
RUN chown user:group /app
USER user

# setting up entry point
ENTRYPOINT ["./decoder"]

