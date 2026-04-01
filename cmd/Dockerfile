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

# building services
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o ./generator ./cmd/generator && \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o ./storage ./cmd/storage && \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o ./encoder ./cmd/encoder && \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o ./decoder ./cmd/decoder

# generator scratch
FROM scratch AS generator
COPY /config ./config
COPY --from=builder /app/generator /generator
ENTRYPOINT ["/generator"]

# storage scratch
FROM scratch AS storage
COPY /config ./config
COPY --from=builder /app/storage /storage
ENTRYPOINT ["/storage"]

# encoder scratch
FROM scratch AS encoder
COPY /config ./config
COPY --from=builder /app/encoder /encoder
ENTRYPOINT ["/encoder"]

# decoder scratch
FROM scratch AS decoder
COPY /config ./config
COPY --from=builder /app/decoder /decoder
ENTRYPOINT ["/decoder"]
