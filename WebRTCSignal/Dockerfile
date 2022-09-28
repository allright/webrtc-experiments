FROM golang:1.18.3 AS builder
COPY . /app
WORKDIR /app
RUN CGO_ENABLED=0 go build

FROM scratch
# Copy our static executable.
COPY --from=builder /app/keepalive ./keepalive
# Run the hello binary.
ENTRYPOINT ["./keepalive"]