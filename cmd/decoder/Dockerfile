FROM alpine:latest
WORKDIR /app

# copying from previous stage
COPY /config ./config
COPY --from=vault-base /app/decoder ./decoder

# creating new user for security
RUN addgroup -S group && adduser -S user -G group
RUN chown user:group /app
USER user

# setting up entry point
ENTRYPOINT ["./decoder"]
