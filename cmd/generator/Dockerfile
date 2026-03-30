# building stage
FROM golang:1.25-alpine AS builder
WORKDIR /

# cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# copying files
COPY cmd/ ./cmd
COPY pkg/ ./pkg
COPY proto/ ./proto
COPY internal/ ./internal

# building
RUN go build -o generator ./cmd/generator


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
